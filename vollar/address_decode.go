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
	"fmt"
	"github.com/blocktree/bitcoin-adapter/bitcoin"
	"github.com/assetsadapterstore/vollar-adapter/vollar_addrdec"
	"github.com/blocktree/go-owcdrivers/addressEncoder"
	"github.com/blocktree/go-owcrypt"
)

func init() {

}

//var (
//	AddressDecoder = &openwallet.AddressDecoder{
//		PrivateKeyToWIF:    PrivateKeyToWIF,
//		PublicKeyToAddress: PublicKeyToAddress,
//		WIFToPrivateKey:    WIFToPrivateKey,
//	}
//)

type addressDecoder struct {
	wm *WalletManager //钱包管理者
}

//NewAddressDecoder 地址解析器
func NewAddressDecoder(wm *WalletManager) *addressDecoder {
	decoder := addressDecoder{}
	decoder.wm = wm
	return &decoder
}

//PrivateKeyToWIF 私钥转WIF
func (decoder *addressDecoder) PrivateKeyToWIF(priv []byte, isTestnet bool) (string, error) {

	cfg := vollar_addrdec.VDS_mainnetPrivateWIFCompressed
	if decoder.wm.Config.IsTestNet {
		cfg = vollar_addrdec.VDS_testnetPrivateWIFCompressed
	}

	wif, _ := vollar_addrdec.Default.AddressEncode(priv, cfg)

	return wif, nil

}

//PublicKeyToAddress 公钥转地址
func (decoder *addressDecoder) PublicKeyToAddress(pub []byte, isTestnet bool) (string, error) {
	vollar_addrdec.Default.IsTestNet = decoder.wm.Config.IsTestNet
	address, err := vollar_addrdec.Default.AddressEncode(pub)
	if err != nil {
		return "", err
	}

	if decoder.wm.Config.RPCServerType == bitcoin.RPCServerCore {
		//如果使用core钱包作为全节点，需要导入地址到core，这样才能查询地址余额和utxo
		err := decoder.wm.ImportAddress(address, "")
		if err != nil {
			return "", err
		}
	}

	return address, nil

}

//RedeemScriptToAddress 多重签名赎回脚本转地址
func (decoder *addressDecoder) RedeemScriptToAddress(pubs [][]byte, required uint64, isTestnet bool) (string, error) {

	cfg := vollar_addrdec.VDS_mainnetAddressP2SH
	if decoder.wm.Config.IsTestNet {
		cfg = vollar_addrdec.VDS_testnetAddressP2SH
	}

	redeemScript := make([]byte, 0)

	for _, pub := range pubs {
		redeemScript = append(redeemScript, pub...)
	}

	pkHash := owcrypt.Hash(redeemScript, 0, owcrypt.HASH_ALG_HASH160)

	address, _ := vollar_addrdec.Default.AddressEncode(pkHash, cfg)

	return address, nil

}

//WIFToPrivateKey WIF转私钥
func (decoder *addressDecoder) WIFToPrivateKey(wif string, isTestnet bool) ([]byte, error) {

	cfg := vollar_addrdec.VDS_mainnetPrivateWIFCompressed
	if decoder.wm.Config.IsTestNet {
		cfg = vollar_addrdec.VDS_testnetPrivateWIFCompressed
	}

	priv, err := vollar_addrdec.Default.AddressDecode(wif, cfg)
	if err != nil {
		return nil, err
	}

	return priv, err

}

//ScriptPubKeyToBech32Address scriptPubKey转Bech32地址
func (decoder *addressDecoder) ScriptPubKeyToBech32Address(scriptPubKey []byte) (string, error) {
	var (
		hash []byte
	)

	cfg := addressEncoder.BTC_mainnetAddressBech32V0
	if decoder.wm.Config.IsTestNet {
		cfg = addressEncoder.BTC_testnetAddressBech32V0
	}

	if len(scriptPubKey) == 22 || len(scriptPubKey) == 34 {

		hash = scriptPubKey[2:]

		address := addressEncoder.AddressEncode(hash, cfg)

		return address, nil

	} else {
		return "", fmt.Errorf("scriptPubKey length is invalid")
	}

}