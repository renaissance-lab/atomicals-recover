package core

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

const (
	DustAmount = int64(546) // dust
)

func TxSize(expectOut int) int64 {
	return int64(150) + int64(120) + int64(expectOut-1)*32
}

func CreateRecoverTx(ctrlKeyStr, lockedAddrStr, fundingKeyStr, fundingAddrStr, outAddrStr, bitworkR, spenderScriptStr string, lockedHash string, lockedIndex uint32, amount, mintTokenStashiAmount int64, fundingHash string, fundingIndex uint32, fundingAmount, feeRate int64, net Network) (string, error) {

	wif, err := btcutil.DecodeWIF(ctrlKeyStr)
	if err != nil {
		fmt.Println("DecodeWIF ctrlKeyStr err ", err)
		return "", err
	}
	fundingwif, err := btcutil.DecodeWIF(fundingKeyStr)
	if err != nil {
		fmt.Println("DecodeWIF fundingKeyStr err ", err)
		return "", err
	}
	chainParam := &chaincfg.TestNet3Params
	if net == MainNetWork {
		chainParam = &chaincfg.MainNetParams
	}
	// lockedAddrStr := "tb1pdksr7shjlwm55sqct3gqsfz2vzszpfpxg0j8tv3dgf3xv7wrupzs0h376y"
	spendAddr, err := btcutil.DecodeAddress(lockedAddrStr, chainParam)
	if err != nil {
		log.Println("DecodeAddress spendAddr err", err)
		return "", err
	}

	spenderAddrByte, err := txscript.PayToAddrScript(spendAddr)
	if err != nil {
		log.Println("spendAddr PayToAddrScript err", err)
		return "", err
	}

	fundingAddr, err := btcutil.DecodeAddress(fundingAddrStr, chainParam)
	if err != nil {
		log.Println("DecodeAddress fundingAddrStr err", err)
		return "", err
	}
	fundingAddrByte, err := txscript.PayToAddrScript(fundingAddr)
	if err != nil {
		log.Println("fundingAddr PayToAddrScript err", err)
		return "", err
	}

	// outAddrStr := "tb1px8ap86kjz2cesyd6sy4r6cmfhpxe2lhd0zvshatm73zlz4wrnpesqfpxes"
	outAddr, err := btcutil.DecodeAddress(outAddrStr, chainParam)
	if err != nil {
		log.Println("DecodeAddress outAddrStr err", err)
		return "", err
	}

	outAddrByte, err := txscript.PayToAddrScript(outAddr)
	if err != nil {
		log.Println("outAddr PayToAddrScript err", err)
		return "", err
	}

	pkScript, err := hex.DecodeString(spenderScriptStr)
	if err != nil {
		log.Println("pkScript err", err)
		return "", err
	}
	tapLeaf := txscript.NewBaseTapLeaf(pkScript)
	tapScriptTree := txscript.AssembleTaprootScriptTree(tapLeaf)

	internalKey := wif.PrivKey.PubKey()
	ctrlBlock := tapScriptTree.LeafMerkleProofs[0].ToControlBlock(
		internalKey,
	)

	// tapScriptRootHash := tapScriptTree.RootNode.TapHash()
	// outputKey := txscript.ComputeTaprootOutputKey(
	// 	internalKey, tapScriptRootHash[:],
	// )
	// p2trScript, err := txscript.PayToTaprootScript(outputKey)
	// if err != nil {
	// 	log.Println("p2trScript", err)
	// 	return "", err
	// }
	// fmt.Printf("p2trScript      %x\n", p2trScript)
	// fmt.Printf("spenderAddrByte %x\n", spenderAddrByte)

	chainHashFrom, err := chainhash.NewHashFromStr(lockedHash)
	if err != nil {
		fmt.Println("invalid lockedHash")
		return "", err
	}
	outPoint := &wire.OutPoint{
		Hash:  *chainHashFrom,
		Index: lockedIndex,
	}

	chainHashFunding, err := chainhash.NewHashFromStr(fundingHash)
	if err != nil {
		fmt.Println("invalid fundingHash")
		return "", err
	}
	fundingPoint := &wire.OutPoint{
		Hash:  *chainHashFunding,
		Index: fundingIndex,
	}
	txIn0 := &wire.TxIn{
		PreviousOutPoint: *outPoint,
	}
	txIn1 := &wire.TxIn{
		PreviousOutPoint: *fundingPoint,
	}

	txIns := []*wire.TxIn{}
	txIns = append(txIns, txIn0)
	txIns = append(txIns, txIn1)

	var txOuts []*wire.TxOut
	// out0 origin try to mint token
	txOut0 := &wire.TxOut{
		PkScript: outAddrByte,
		Value:    mintTokenStashiAmount,
	}
	txOuts = append(txOuts, txOut0)
	// out1 remain satashi
	if amount > mintTokenStashiAmount+DustAmount {
		txOut1 := &wire.TxOut{
			PkScript: outAddrByte,
			Value:    amount - mintTokenStashiAmount,
		}
		txOuts = append(txOuts, txOut1)
	}

	// out2 redeem
	size := TxSize(3)
	fee := size * feeRate
	if fundingAmount > fee+DustAmount {
		txOut2 := &wire.TxOut{
			PkScript: outAddrByte,
			Value:    fundingAmount - fee,
		}
		txOuts = append(txOuts, txOut2)
	}
	tx := &wire.MsgTx{
		TxIn:  txIns,
		TxOut: txOuts,
	}

	tx.Version = 1

	//
	var w bytes.Buffer
	err = tx.SerializeNoWitness(&w)
	if err != nil {
		fmt.Println("SerializeNoWitness err ", err)
		return "", err
	}
	// serializedBlock := w.Bytes()
	// fmt.Printf("serializedBlock %x\n", serializedBlock)

	// var seqNum uint32
	for i := uint32(0); i < 4294967295; i++ {
		tx.TxIn[0].Sequence = i
		txHash := tx.TxHash().String()
		if bitworkR == "" || strings.HasPrefix(txHash, bitworkR) {
			// seqNum = i
			break
		}

	}

	//fmt.Printf("get num %d received stop signal\n", num)
	// fmt.Println("get seq ", seqNum)

	fundingOutPoint := &wire.TxOut{
		Value:    fundingAmount,
		PkScript: fundingAddrByte,
	}
	a := txscript.NewMultiPrevOutFetcher(map[wire.OutPoint]*wire.TxOut{
		*outPoint: &wire.TxOut{
			Value:    amount,
			PkScript: spenderAddrByte,
		},
	})

	a.AddPrevOut(*fundingPoint, fundingOutPoint)
	sigHashes := txscript.NewTxSigHashes(tx, a)

	// Now that we have the sig, we'll make a valid witness
	// including the control block.

	signature, err := txscript.RawTxInTapscriptSignature(tx, sigHashes, 0, amount, spenderAddrByte, tapLeaf, txscript.SigHashDefault, wif.PrivKey)
	if err != nil {
		panic(err)
	}
	ctrlBlockBytes, err := ctrlBlock.ToBytes()
	if err != nil {
		fmt.Println("ctrlBlock.ToBytes() err ", err)
		return "", err
	}
	// calc input1 signature
	signature1, err := txscript.TaprootWitnessSignature(tx, sigHashes, 1, fundingAmount, fundingAddrByte, txscript.SigHashDefault, fundingwif.PrivKey)
	if err != nil {
		fmt.Println("TaprootWitnessSignature err", err)
		return "", err
	}

	txCopy := tx.Copy()
	txCopy.TxIn[0].Witness = wire.TxWitness{
		signature, pkScript, ctrlBlockBytes,
	}
	txCopy.TxIn[1].Witness = signature1

	var signedTx bytes.Buffer
	txCopy.Serialize(&signedTx)

	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	fmt.Println("rawtx: ", hexSignedTx)

	return txCopy.TxHash().String(), nil
}
