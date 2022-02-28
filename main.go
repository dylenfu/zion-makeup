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

package main

import (
	"flag"

	"github.com/dylenfu/zion-makeup/config"
	"github.com/dylenfu/zion-makeup/core"
)

var (
	nodes    int
	filePath string
)

func init() {
	flag.IntVar(&nodes, "nodes", 7, "denotes nodes number")
	flag.StringVar(&filePath, "config", "config.json", "configuration file path")
	flag.Parse()
}

func main() {
	config.LoadConfig(filePath)
	core.Run(nodes)
}
