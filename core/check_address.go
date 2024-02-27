package core

import (
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
)

func CheckWifMatchPubkey(keystr, pubkeystr string) (bool, error) {
	wif, err := btcutil.DecodeWIF(keystr)
	if err != nil {
		fmt.Printf("DecodeWIF %v, err %v \n", keystr, err)
		return false, err
	}

	pubkey := wif.PrivKey.PubKey()
	// skip header 02 / 03 for check schonorr pubkey
	pubkeyCalc := fmt.Sprintf("%x", pubkey.SerializeCompressed()[1:])
	if pubkeystr == pubkeyCalc {
		return true, nil
	}
	// fmt.Printf("pubkeystr %v, pubkeyCalc %v \n", pubkeystr, pubkeyCalc)
	return false, nil
}

func CheckWifMatchP2TRScriptAddress(keyStr, expectAddress, protocol, opType, bitworkr, bitworkc, mintTicker string, timeUnix, nonce int64, net Network) (bool, string, error) {

	address, script, err := GenerateP2TRScriptAddress(keyStr, protocol, opType, bitworkr, bitworkc, mintTicker, timeUnix, nonce, net)
	if err != nil {
		return false, "", err
	}
	if address == expectAddress {
		return true, script, nil
	}
	return false, "", nil
}

func CheckWifMatchP2TRAddress(keyStr, expectAddress string, network Network) (bool, error) {
	address, err := GenerateP2TRAddress(keyStr, expectAddress, network)
	if err != nil {
		return false, err
	}
	if address == expectAddress {
		return true, nil
	}
	return false, nil
}
