/*
 * Copyright 2018 The openwallet Authors
 * This file is part of the openwallet library.
 *
 * The openwallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The openwallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */

package vollar

import (
	"github.com/blocktree/go-owcdrivers/btcTransaction"
	"github.com/blocktree/go-owcrypt"
)

const (
	//币种
	Symbol    = "VDS"
	CurveType = owcrypt.ECC_CURVE_SECP256K1
	Decimals  = int32(8)
)

var (
	MainNetAddressPrefix = btcTransaction.AddressPrefix{P2PKHPrefix: []byte{0x10, 0x1C}, P2WPKHPrefix: []byte{0x10, 0x41}, P2SHPrefix: nil, Bech32Prefix:"vds"}
	TestNetAddressPrefix = btcTransaction.AddressPrefix{P2PKHPrefix: []byte{0x1D, 0x25}, P2WPKHPrefix: []byte{0x1C, 0xBA}, P2SHPrefix: nil, Bech32Prefix:"vds"}
)
