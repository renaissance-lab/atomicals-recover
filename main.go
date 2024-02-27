package main

import (
	"fmt"

	"github.com/renaissance-lab/atomicals-recover/config"
	"github.com/renaissance-lab/atomicals-recover/core"
)

func main() {
	config.Setup()

	core.CborToPayLoad("a16461726773a568626974776f726b63643030303068626974776f726b72606b6d696e745f7469636b65726670686f746f6e656e6f6e6365016474696d651a65dd204f")

	// first check address info
	lockWif := config.Cfg.Locked.Wif
	lockPubkey := config.Cfg.Locked.Pubkey
	valid, err := core.CheckWifMatchPubkey(lockWif, lockPubkey)
	if err != nil || !valid {
		fmt.Println("CheckWifMatchPubkey for locked address error")
		return
	}
	// check lock address match script and priv
	lockedScriptAddress := config.Cfg.Locked.Address
	protocol := config.Cfg.Token.Protocol
	opType := config.Cfg.Token.OpType
	bitworkr := config.Cfg.Token.Bitworkr
	bitworkc := config.Cfg.Token.Bitworkc
	mintTicker := config.Cfg.Token.MintTicker
	timeUnix := config.Cfg.LockedScript.Time
	nonce := config.Cfg.LockedScript.Nonce
	network := core.Network(config.Cfg.Network)
	valid, script, err := core.CheckWifMatchP2TRScriptAddress(lockWif, lockedScriptAddress, protocol, opType, bitworkr, bitworkc, mintTicker, timeUnix, nonce, network)
	if err != nil || !valid {
		fmt.Printf("lockWif %v, lockedScriptAddress %v, protocol %v, opType %v, bitworkr %v, bitworkc %v, mintTicker %v, timeUnix %v, nonce %v, network %v \n", lockWif, lockedScriptAddress, protocol, opType, bitworkr, bitworkc, mintTicker, timeUnix, nonce, network)
		fmt.Println("CheckWifMatchP2TRScriptAddress for locked address error")
		return
	}

	fundingWif := config.Cfg.Funding.Wif
	fundingAddress := config.Cfg.Funding.Address
	valid, err = core.CheckWifMatchP2TRAddress(fundingWif, fundingAddress, network)
	if err != nil || !valid {
		fmt.Println("CheckWifMatchP2TRAddress for funding address error")
		return
	}
	lockedHash := config.Cfg.LockedUtxo.Hash
	lockedIndex := config.Cfg.LockedUtxo.Index
	lockAmount := config.Cfg.LockedUtxo.Amount
	mintNeedAmount := config.Cfg.Token.MintNeedAmount
	outAddrStr := config.Cfg.RedeemAddress
	fundHash := config.Cfg.FundingUtxo.Hash
	fundIndex := config.Cfg.FundingUtxo.Index
	fundingAmount := config.Cfg.FundingUtxo.Amount
	feeRate := config.Cfg.FeeRate
	txHash, err := core.CreateRecoverTx(lockWif, lockedScriptAddress, fundingWif, fundingAddress, outAddrStr, bitworkr, script, lockedHash, uint32(lockedIndex), lockAmount, mintNeedAmount, fundHash, uint32(fundIndex), fundingAmount, feeRate, network)
	if err != nil {
		fmt.Printf("CreateRecoverTx err ", err)
		return
	}
	fmt.Println("unlock success, txHash ", txHash)
}
