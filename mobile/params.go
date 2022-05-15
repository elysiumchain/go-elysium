// Copyright 2016 The go-elysium Authors
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

// Contains all the wrappers from the params package.

package gely

import (
	"encoding/json"

	"github.com/elysiumchain/go-elysium/core"
	"github.com/elysiumchain/go-elysium/p2p/enode"
	"github.com/elysiumchain/go-elysium/params"
)

// MainnetGenesis returns the JSON spec to use for the main Elysium network. It
// is actually empty since that defaults to the hard coded binary genesis block.
func MainnetGenesis() string {
	return ""
}

// ElysiumTestnetGenesis returns the JSON spec to use for the ElysiumTestnet test network.
func ElysiumTestnetGenesis() string {
	enc, err := json.Marshal(core.DefaultElysiumTestnetGenesisBlock())
	if err != nil {
		panic(err)
	}
	return string(enc)
}

// SepoliaGenesis returns the JSON spec to use for the Sepolia test network.
func SepoliaGenesis() string {
	enc, err := json.Marshal(core.DefaultSepoliaGenesisBlock())
	if err != nil {
		panic(err)
	}
	return string(enc)
}

// RinkebyGenesis returns the JSON spec to use for the Rinkeby test network
func RinkebyGenesis() string {
	enc, err := json.Marshal(core.DefaultRinkebyGenesisBlock())
	if err != nil {
		panic(err)
	}
	return string(enc)
}

// GoerliGenesis returns the JSON spec to use for the Goerli test network
func GoerliGenesis() string {
	enc, err := json.Marshal(core.DefaultGoerliGenesisBlock())
	if err != nil {
		panic(err)
	}
	return string(enc)
}

// FoundationBootnodes returns the enode URLs of the P2P bootstrap nodes operated
// by the foundation running the V5 discovery protocol.
func FoundationBootnodes() *Enodes {
	nodes := &Enodes{nodes: make([]*enode.Node, len(params.MainnetBootnodes))}
	for i, url := range params.MainnetBootnodes {
		var err error
		nodes.nodes[i], err = enode.Parse(enode.ValidSchemes, url)
		if err != nil {
			panic("invalid node URL: " + err.Error())
		}
	}
	return nodes
}
