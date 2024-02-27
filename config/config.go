package config

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

var Cfg Config

type WalletInfo struct {
	Wif     string `yaml:"WIF"`
	Address string `yaml:"address"`
	Pubkey  string `yaml:"pubkey"`
}
type Utxo struct {
	Hash   string `yaml:"hash"`
	Index  int32  `yaml:"index"`
	Amount int64  `yaml:"amount"`
}

type LockedScriptInfo struct {
	Time  int64 `yaml:"time"`
	Nonce int64 `yaml:"nonce"`
}
type AtomTokenArg struct {
	OpType         string `yaml:"opType"`
	Protocol       string `yaml:"protocol"`
	Bitworkc       string `yaml:"bitworkc"`
	Bitworkr       string `yaml:"bitworkr"`
	MintTicker     string `yaml:"mint_ticker"`
	MintNeedAmount int64  `yaml:"mint_need_amount"`
}

type Config struct {
	Locked        WalletInfo       `yaml:"locked"`
	LockedUtxo    Utxo             `yaml:"lockedUtxo"`
	LockedScript  LockedScriptInfo `yaml:"lockedScript"`
	Funding       WalletInfo       `yaml:"funding"`
	FundingUtxo   Utxo             `yaml:"fundingUtxo"`
	RedeemAddress string           `yaml:"redeemAddress"`
	UnlockType    int              `yaml:"unlockType"`
	Network       int              `yaml:"network"`
	FeeRate       int64            `yaml:"feeRate"`
	Token         AtomTokenArg     `yaml:"token"`
}

func Setup() {
	fileContent, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		// config
		fmt.Println("readconfig file error")
		return
	}

	err = yaml.Unmarshal(fileContent, &Cfg)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("config:%#v\n", Cfg)
}
