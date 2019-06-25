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
	"github.com/blocktree/bitcoin-adapter/bitcoin"
	"github.com/blocktree/openwallet/log"
)

const (
	maxAddresNum = 10000
)

type WalletManager struct {
	*bitcoin.WalletManager
}


func NewWalletManager() *WalletManager {
	wm := WalletManager{}
	wm.WalletManager = bitcoin.NewWalletManager()
	wm.Config = bitcoin.NewConfig(Symbol, CurveType, Decimals)
	wm.Config.MainNetAddressPrefix = MainNetAddressPrefix
	wm.Config.TestNetAddressPrefix = TestNetAddressPrefix
	wm.Decoder = NewAddressDecoder(&wm)
	wm.TxDecoder = NewTransactionDecoder(&wm)
	wm.Log = log.NewOWLogger(wm.Symbol())

	//不扫描内存池
	wm.Blockscanner.IsScanMemPool = false
	return &wm
}

func (wm *WalletManager) ListAddresses() ([]string, error) {
	var (
		addresses = make([]string, 0)
	)

	request := []interface{}{
		true,
	}

	result, err := wm.WalletClient.Call("v_listaddresses", request)
	if err != nil {
		return nil, err
	}

	array := result.Array()
	for _, a := range array {
		addresses = append(addresses, a.String())
	}

	return addresses, nil
}