// Copyright 2015 The go-elysium Authors
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

package params

import "github.com/elysiumchain/go-elysium/common"

// MainnetBootnodes are the enode URLs of the P2P bootstrap nodes running on
// the main Elysium network.
var MainnetBootnodes = []string{
	// Elysium Foundation Go Bootnodes
}

// ElysiumTestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// ElysiumTestnet test network.
var ElysiumTestnetBootnodes = []string{
	// AWS bootnode
	"enode://ab9eee973d818a81474d16e1ec8194b0ae1e4e2bdf3fdf091a6c45f9d460f46827e0d161ecb24e9e004cc370a811178f32555860e467428909e9a2752aab2d58@15.164.246.71:20202",

	// YJ node
	"enode://7f433d8e6e4c2e07ff575b63f66ea85a6a2db108ee4cee93ce87fb7a6a923f17b2dd6894db926723aacaa040790cef2b323be108a61d1b9dfd8668875303f78b@203.229.158.227:20202",

	// Elixir node
	"enode://892658ad7f9aacf72fe2dac0d8eda914e58ee474b656eaa81d883ae3ed55b9df598c77a131b3620f3eb399915075bf12bcecd205426bcb41a17050674816d60c@112.159.68.43:20202",
}

// SepoliaBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Sepolia test network.
var SepoliaBootnodes = []string{
	// 
}

// RinkebyBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Rinkeby test network.
var RinkebyBootnodes = []string{
	// 
}

// GoerliBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// G??rli test network.
var GoerliBootnodes = []string{
	// 
}

var KilnBootnodes = []string{
	// 
}

var V5Bootnodes = []string{
	// 
}

const dnsPrefix = "enrtree://AKA3AM6LPBYEUDMVNU3BSVQJ5AD45Y7YPOHJLEF6W26QOE4VTUDPE@"

// KnownDNSNetwork returns the address of a public DNS-based node list for the given
// genesis hash and protocol. See https://github.com/elysiumchain/discv4-dns-lists for more
// information.
func KnownDNSNetwork(genesis common.Hash, protocol string) string {
	var net string
	switch genesis {
	case MainnetGenesisHash:
		net = "mainnet"
	case ElysiumTestnetGenesisHash:
		net = "elysiumTestnet"
	case RinkebyGenesisHash:
		net = "rinkeby"
	case GoerliGenesisHash:
		net = "goerli"
	default:
		return ""
	}
	return dnsPrefix + protocol + "." + net + ".ethdisco.net"
}
