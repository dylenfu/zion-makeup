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
	"testing"

	"github.com/dylenfu/zion-makeup/pkg/files"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

func TestPrintMinerList(t *testing.T) {
	filepath := "/Users/dylen/software/hotstuff/zion-makeup/build/extra"
	nodesList := []string{}

	enc, err := files.ReadFile(filepath)
	if err != nil {
		t.Fatal(err)
	}

	raw, err := hexutil.Decode(string(enc))
	if err != nil {
		t.Fatal(err)
	}

	extra, err := types.ExtractHotstuffExtraPayload(raw)
	if err != nil {
		t.Fatal(err)
	}

	ret := "("
	for i, v := range extra.Validators {
		ret += v.Hex()
		if i != len(nodesList)-1 {
			ret += " "
		}
	}
	ret += ")"
	t.Log(ret)
}
