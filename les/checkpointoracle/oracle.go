// Copyright 2019 The go-elysium Authors
// This file is part of the go-elysium library.
//
// The go-elysium library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-elysium library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-elysium library. If not, see <http://www.gnu.org/licenses/>.

// Package checkpointoracle is a wrapper of checkpoint oracle contract with
// additional rules defined. This package can be used both in LES client or
// server side for offering oracle related APIs.
package checkpointoracle

import (
	"encoding/binary"
	"sync"
	"sync/atomic"
	"time"

	"github.com/elysiumchain/go-elysium/accounts/abi/bind"
	"github.com/elysiumchain/go-elysium/common"
	"github.com/elysiumchain/go-elysium/contracts/checkpointoracle"
	"github.com/elysiumchain/go-elysium/crypto"
	"github.com/elysiumchain/go-elysium/log"
	"github.com/elysiumchain/go-elysium/params"
)

// CheckpointOracle is responsible for offering the latest stable checkpoint
// generated and announced by the contract admins on-chain. The checkpoint can
// be verified by clients locally during the checkpoint syncing.
type CheckpointOracle struct {
	config   *params.CheckpointOracleConfig
	contract *checkpointoracle.CheckpointOracle

	running  int32                                 // Flag whether the contract backend is set or not
	getLocal func(uint64) params.TrustedCheckpoint // Function used to retrieve local checkpoint

	checkMu              sync.Mutex                // Mutex to sync access to the fields below
	lastCheckTime        time.Time                 // Time we last checked the checkpoint
	lastCheckPoint       *params.TrustedCheckpoint // The last stable checkpoint
	lastCheckPointHeight uint64                    // The height of last stable checkpoint
}

// New creates a checkpoint oracle handler with given configs and callback.
func New(config *params.CheckpointOracleConfig, getLocal func(uint64) params.TrustedCheckpoint) *CheckpointOracle {
	return &CheckpointOracle{
		config:   config,
		getLocal: getLocal,
	}
}

// Start binds the contract backend, initializes the oracle instance
// and marks the status as available.
func (oracle *CheckpointOracle) Start(backend bind.ContractBackend) {
	contract, err := checkpointoracle.NewCheckpointOracle(oracle.config.Address, backend)
	if err != nil {
		log.Error("Oracle contract binding failed", "err", err)
		return
	}
	if !atomic.CompareAndSwapInt32(&oracle.running, 0, 1) {
		log.Error("Already bound and listening to registrar")
		return
	}
	oracle.contract = contract
}

// IsRunning returns an indicator whether the oracle is running.
func (oracle *CheckpointOracle) IsRunning() bool {
	return atomic.LoadInt32(&oracle.running) == 1
}

// Contract returns the underlying raw checkpoint oracle contract.
func (oracle *CheckpointOracle) Contract() *checkpointoracle.CheckpointOracle {
	return oracle.contract
}

// StableCheckpoint returns the stable checkpoint which was generated by local
// indexers and announced by trusted signers.
func (oracle *CheckpointOracle) StableCheckpoint() (*params.TrustedCheckpoint, uint64) {
	oracle.checkMu.Lock()
	defer oracle.checkMu.Unlock()
	if time.Since(oracle.lastCheckTime) < 1*time.Minute {
		return oracle.lastCheckPoint, oracle.lastCheckPointHeight
	}
	// Look it up properly
	// Retrieve the latest checkpoint from the contract, abort if empty
	latest, hash, height, err := oracle.contract.Contract().GetLatestCheckpoint(nil)
	oracle.lastCheckTime = time.Now()
	if err != nil || (latest == 0 && hash == [32]byte{}) {
		oracle.lastCheckPointHeight = 0
		oracle.lastCheckPoint = nil
		return oracle.lastCheckPoint, oracle.lastCheckPointHeight
	}
	local := oracle.getLocal(latest)

	// The following scenarios may occur:
	//
	// * local node is out of sync so that it doesn't have the
	//   checkpoint which registered in the contract.
	// * local checkpoint doesn't match with the registered one.
	//
	// In both cases, no stable checkpoint will be returned.
	if local.HashEqual(hash) {
		oracle.lastCheckPointHeight = height.Uint64()
		oracle.lastCheckPoint = &local
		return oracle.lastCheckPoint, oracle.lastCheckPointHeight
	}
	return nil, 0
}

// VerifySigners recovers the signer addresses according to the signature and
// checks whether there are enough approvals to finalize the checkpoint.
func (oracle *CheckpointOracle) VerifySigners(index uint64, hash [32]byte, signatures [][]byte) (bool, []common.Address) {
	// Short circuit if the given signatures doesn't reach the threshold.
	if len(signatures) < int(oracle.config.Threshold) {
		return false, nil
	}
	var (
		signers []common.Address
		checked = make(map[common.Address]struct{})
	)
	for i := 0; i < len(signatures); i++ {
		if len(signatures[i]) != 65 {
			continue
		}
		// EIP 191 style signatures
		//
		// Arguments when calculating hash to validate
		// 1: byte(0x19) - the initial 0x19 byte
		// 2: byte(0) - the version byte (data with intended validator)
		// 3: this - the validator address
		// --  Application specific data
		// 4 : checkpoint section_index (uint64)
		// 5 : checkpoint hash (bytes32)
		//     hash = keccak256(checkpoint_index, section_head, cht_root, bloom_root)
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, index)
		data := append([]byte{0x19, 0x00}, append(oracle.config.Address.Bytes(), append(buf, hash[:]...)...)...)
		signatures[i][64] -= 27 // Transform V from 27/28 to 0/1 according to the yellow paper for verification.
		pubkey, err := crypto.Ecrecover(crypto.Keccak256(data), signatures[i])
		if err != nil {
			return false, nil
		}
		var signer common.Address
		copy(signer[:], crypto.Keccak256(pubkey[1:])[12:])
		if _, exist := checked[signer]; exist {
			continue
		}
		for _, s := range oracle.config.Signers {
			if s == signer {
				signers = append(signers, signer)
				checked[signer] = struct{}{}
			}
		}
	}
	threshold := oracle.config.Threshold
	if uint64(len(signers)) < threshold {
		log.Warn("Not enough signers to approve checkpoint", "signers", len(signers), "threshold", threshold)
		return false, nil
	}
	return true, signers
}
