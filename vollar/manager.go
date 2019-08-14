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
	"github.com/shopspring/decimal"
	"math"
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


//EstimateFee 预估手续费
func (wm *WalletManager) EstimateFee(inputs, outputs int64, feeRate decimal.Decimal) (decimal.Decimal, error) {

	var piece int64 = 1

	//UTXO如果大于设定限制，则分拆成多笔交易单发送
	if inputs > int64(wm.Config.MaxTxInputs) {
		piece = int64(math.Ceil(float64(inputs) / float64(wm.Config.MaxTxInputs)))
	}

	//双倍费率
	feeRate = feeRate.Mul(decimal.New(2, 0))

	//计算公式如下：148 * 输入数额 + 67 * 输出数额 + 79
	trx_bytes := decimal.New(inputs*148+outputs*67+piece*79, 0)
	trx_fee := trx_bytes.Div(decimal.New(1000, 0)).Mul(feeRate)
	trx_fee = trx_fee.Round(wm.Decimal())
	//wm.Log.Debugf("trx_fee: %s", trx_fee.String())
	//wm.Log.Debugf("MinFees: %s", wm.Config.MinFees.String())
	//是否低于最小手续费
	if trx_fee.LessThan(wm.Config.MinFees) {
		trx_fee = wm.Config.MinFees
	}

	return trx_fee, nil
}