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
	"os"
	"path"

	"github.com/dylenfu/zion-makeup/config"

	"github.com/dylenfu/zion-makeup/log"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	folder           = "build"
	initAllocBalance = "100000000000000000000000000000"
)

func Run(n int) {
	log.Infof("generate %d nodes", n)

	nodes := generateNodes(n)
	sortedNodes := SortNodes(nodes)
	saveNodes(sortedNodes)
	saveAlloc(sortedNodes)
	generateExtra(sortedNodes)
	generateStaticNodesFile(sortedNodes)
}

func generateNodes(n int) []*Node {
	nodes := make([]*Node, 0)

	for i := 0; i < n; i++ {
		key, _ := crypto.GenerateKey()
		addr := crypto.PubkeyToAddress(key.PublicKey)
		nodekey := hexutil.Encode(crypto.FromECDSA(key))
		pubInf, _ := NodeKey2PublicInfo(nodekey)

		log.Infof("addr: %s, pubKey: %s", addr, pubInf)

		node := &Node{
			Address: addr.Hex(),
			NodeKey: nodekey,
		}
		nodes = append(nodes, node)
	}

	return nodes
}

func saveNodes(sortedNodes []*Node) {
	os.MkdirAll(path.Join(folder, "nodes"), os.ModePerm)

	for i := 0; i < len(sortedNodes); i++ {
		nodekey := sortedNodes[i].NodeKey
		pubInf, _ := NodeKey2PublicInfo(nodekey)

		sNodeIndex := fmt.Sprintf("node%d", i)
		nodeDir := path.Join(folder, "nodes", sNodeIndex)
		os.MkdirAll(nodeDir, os.ModePerm)
		os.WriteFile(path.Join(nodeDir, "nodekey"), []byte(nodekey), os.ModePerm)
		os.WriteFile(path.Join(nodeDir, "pubkey"), []byte(pubInf), os.ModePerm)
	}
}

func generateExtra(sortedNodes []*Node) {
	extra, err := Encode(NodesAddress(sortedNodes))
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(path.Join(folder, "extra"), []byte(extra), os.ModePerm); err != nil {
		panic(err)
	}
	log.Infof("genesis extra %s", extra)
}

func generateStaticNodesFile(sortedNodes []*Node) {
	staticNodes := make([]string, 0)
	nodesPerMachine := len(sortedNodes) / len(config.Conf.IpList)
	for i, v := range sortedNodes {
		nodeInf, err := NodeKey2NodeInfo(v.NodeKey)
		if err != nil {
			panic(err)
		}
		ipIndex := i / nodesPerMachine
		port := i%nodesPerMachine + config.Conf.StartPort
		staticNodes = append(staticNodes, NodeStaticInfoTemp(nodeInf, config.Conf.IpList[ipIndex], port))
	}

	enc, err := json.MarshalIndent(staticNodes, "", "\t")
	if err != nil {
		panic(err)
	}
	log.Info(string(enc))
	if err := ioutil.WriteFile(path.Join(folder, "static-nodes"), enc, os.ModePerm); err != nil {
		panic(err)
	}
}

type AllocInfo struct {
	PublicKey string `json:"PublicKey"`
	Balance   string `json:"balance"`
}

func saveAlloc(sortedNodes []*Node) {
	nodesMap := make(map[string]*AllocInfo)
	for _, v := range sortedNodes {
		pubkey, _ := NodeKey2PublicInfo(v.NodeKey)
		nodesMap[v.Address] = &AllocInfo{
			PublicKey: pubkey,
			Balance:   initAllocBalance,
		}
	}

	enc, err := json.MarshalIndent(nodesMap, "", "\t")
	if err != nil {
		panic(err)
	}
	log.Info(string(enc))
	if err := ioutil.WriteFile(path.Join(folder, "alloc-nodes"), enc, os.ModePerm); err != nil {
		panic(err)
	}
}
