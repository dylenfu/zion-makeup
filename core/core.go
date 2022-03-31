/*
 * Copyright (C) 2021 The Zion Authors
 * This file is part of The Zion library.
 *
 * The Zion is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The Zion is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The Zion.  If not, see <http://www.gnu.org/licenses/>.
 */

package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path"

	"github.com/dylenfu/zion-makeup/config"
	"github.com/dylenfu/zion-makeup/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
)

const (
	folder = "build"
)

var env string

func Run(dir string, n int, initAllocBalance string) {
	log.Infof("generate %d nodes", n)

	os.MkdirAll(folder, os.ModePerm)
	env = path.Join(folder, dir)

	nodes := generateNodes(n)
	sortedNodes := SortNodes(nodes)
	saveNodes(sortedNodes)
	//saveAlloc(sortedNodes, initAllocBalance)
	//saveMinerList(sortedNodes)
	//generateExtra(sortedNodes)
	saveGenesis(sortedNodes, initAllocBalance)
	generateStaticNodesFile(sortedNodes)
}

func generateNodes(n int) []*Node {
	nodes := make([]*Node, 0)

	for i := 0; i < n; i++ {
		key, _ := crypto.GenerateKey()
		addr := crypto.PubkeyToAddress(key.PublicKey)

		node := &Node{
			Address: addr,
			NodeKey: key,
		}

		nodes = append(nodes, node)
	}

	return nodes
}

func saveNodes(sortedNodes []*Node) {
	os.MkdirAll(path.Join(env, "nodes"), os.ModePerm)

	for i, v := range sortedNodes {
		sNodeIndex := fmt.Sprintf("node%d", i)
		nodeDir := path.Join(env, "nodes", sNodeIndex)
		os.MkdirAll(nodeDir, os.ModePerm)
		os.WriteFile(path.Join(nodeDir, "nodekey"), []byte(v.NodeKeyHex(false)), os.ModePerm)
		os.WriteFile(path.Join(nodeDir, "pubkey"), []byte(v.PubKeyHex()), os.ModePerm)
	}
}

func generateExtra(sortedNodes []*Node) {
	list := make([]common.Address, 0)
	for _, v := range sortedNodes {
		list = append(list, v.Address)
	}

	extra, err := Encode(list)
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(path.Join(env, "extra.dat"), []byte(extra), os.ModePerm); err != nil {
		panic(err)
	}
	log.Infof("genesis extra %s", extra)
}

func saveMinerList(sortedNodes []*Node) {
	minerlistTxt := "miners=("
	for i, v := range sortedNodes {
		minerlistTxt += v.Address.Hex()
		if i != len(sortedNodes)-1 {
			minerlistTxt += " "
		}
	}
	minerlistTxt += ")"

	if err := ioutil.WriteFile(path.Join(env, "minerlist.sh"), []byte(minerlistTxt), os.ModePerm); err != nil {
		panic(err)
	}

	log.Infof("save miner list %s", minerlistTxt)
}

func generateStaticNodesFile(sortedNodes []*Node) {
	staticNodes := make([]string, 0)
	nodesPerMachine := len(sortedNodes) / len(config.Conf.IpList)
	for i, v := range sortedNodes {
		ipIndex := i / nodesPerMachine
		port := i%nodesPerMachine + config.Conf.StartPort
		staticNodes = append(staticNodes, NodeStaticInfoTemp(v.ID(), config.Conf.IpList[ipIndex], port))
	}

	enc, err := json.MarshalIndent(staticNodes, "", "\t")
	if err != nil {
		panic(err)
	}
	log.Info(string(enc))
	if err := ioutil.WriteFile(path.Join(env, "static-nodes.json"), enc, os.ModePerm); err != nil {
		panic(err)
	}
}

type AllocInfo struct {
	PublicKey string `json:"publicKey"`
	Balance   string `json:"balance"`
}

func saveAlloc(sortedNodes []*Node, initAllocBalance string) {
	nodesMap := make(map[string]*AllocInfo)
	for _, v := range sortedNodes {
		pubkey := v.PubKeyHex()
		nodesMap[v.Address.Hex()] = &AllocInfo{
			PublicKey: pubkey,
			Balance:   initAllocBalance,
		}
	}

	enc, err := json.MarshalIndent(nodesMap, "", "\t")
	if err != nil {
		panic(err)
	}
	log.Info(string(enc))
	if err := ioutil.WriteFile(path.Join(env, "alloc-nodes.json"), enc, os.ModePerm); err != nil {
		panic(err)
	}
}

func saveGenesis(sortedNodes []*Node, initAllocBalance string) {
	nodesMap := make(map[string]*AllocInfo)
	for _, v := range sortedNodes {
		pubkey := v.PubKeyHex()
		nodesMap[v.Address.Hex()] = &AllocInfo{
			PublicKey: pubkey,
			Balance:   initAllocBalance,
		}
	}

	alloc, err := json.MarshalIndent(nodesMap, "", "\t")
	if err != nil {
		panic(err)
	}

	list := make([]common.Address, 0)
	for _, v := range sortedNodes {
		list = append(list, v.Address)
	}

	extra, err := Encode(list)
	if err != nil {
		panic(err)
	}

	data := genesisTemplate(string(alloc), extra)
	if err := ioutil.WriteFile(path.Join(env, "genesis.json"), []byte(data), os.ModePerm); err != nil {
		panic(err)
	}
}

/*
{
    "config": {
        "chainId": 60801,
        "homesteadBlock": 0,
        "eip150Block": 0,
        "eip155Block": 0,
        "eip158Block": 0,
        "byzantiumBlock": 0,
        "constantinopleBlock": 0,
        "petersburgBlock": 0,
        "istanbulBlock": 0,
        "berlinBlock": 0,
        "londonBlock": 0,
        "hotstuff": {
            "protocol": "basic"
        }
    },
    "alloc": {
        "0x258af48e28e4a6846e931ddff8e1cdf8579821e5": {"publicKey": "0x02c07fb7d48eac559a2483e249d27841c18c7ce5dbbbf2796a6963cc9cef27cabd", "balance": "100000000000000000000000000000"},
        "0x6a708455c8777630aac9d1e7702d13f7a865b27c": {"publicKey": "0x02f5135ae0853af71f017a8ecb68e720b729ab92c7123c686e75b7487d4a57ae07", "balance": "100000000000000000000000000000"},
        "0x8c09d936a1b408d6e0afaa537ba4e06c4504a0ae": {"publicKey": "0x03ecac0ebe7224cfd04056c940605a4a9d4cb0367cf5819bf7e5502bf44f68bdd4", "balance": "100000000000000000000000000000"},
        "0xad3bf5ed640cc72f37bd21d64a65c3c756e9c88c": {"publicKey": "0x03d0ecfd09db6b1e4f59da7ebde8f6c3ea3ed09f06f5190477ae4ee528ec692fa8", "balance": "100000000000000000000000000000"}
    },
    "coinbase": "0x0000000000000000000000000000000000000000",
    "difficulty": "0x1",
    "extraData": "0x0000000000000000000000000000000000000000000000000000000000000000f89bf85494258af48e28e4a6846e931ddff8e1cdf8579821e5946a708455c8777630aac9d1e7702d13f7a865b27c948c09d936a1b408d6e0afaa537ba4e06c4504a0ae94ad3bf5ed640cc72f37bd21d64a65c3c756e9c88cb8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c080",
    "gasLimit": "0xffffffff",
    "nonce": "0x4510809143055965",
    "mixhash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "timestamp": "0x00"
}
*/

var zero = big.NewInt(0)

func genesisTemplate(alloc, extra string) string {
	return fmt.Sprintf(`
{
    "config": {
        "chainId": 60801, 
        "homesteadBlock": 0,
        "eip150Block": 0,
        "eip155Block": 0,
        "eip158Block": 0,
        "byzantiumBlock": 0,
        "constantinopleBlock": 0,
        "petersburgBlock": 0,
        "istanbulBlock": 0,
        "berlinBlock": 0,
        "londonBlock": 0,
        "hotstuff": {
            "protocol": "basic"
        }
    },
    "alloc": %s,
    "coinbase": "0x0000000000000000000000000000000000000000",
    "difficulty": "0x1",
    "extraData": "%s",
    "gasLimit": "0xffffffff",
    "nonce": "0x4510809143055965",
    "mixhash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "timestamp": "0x00"
}`, alloc, extra)
}

var defaultGenesisConfig = core.Genesis{
	Config: &params.ChainConfig{
		ChainID:             big.NewInt(60801),
		HomesteadBlock:      zero,
		EIP150Block:         zero,
		EIP155Block:         zero,
		EIP158Block:         zero,
		ByzantiumBlock:      zero,
		ConstantinopleBlock: zero,
		PetersburgBlock:     zero,
		IstanbulBlock:       zero,
		BerlinBlock:         zero,
		LondonBlock:         zero,
		HotStuff: &params.HotStuffConfig{
			Protocol: "basic",
		},
	},
	Alloc:      nil,
	Coinbase:   common.Address{},
	Difficulty: big.NewInt(1),
	ExtraData:  nil,
	GasLimit:   4294967295,
	Nonce:      4976618949627435000,
	Mixhash:    common.Hash{},
	Timestamp:  0,
}
