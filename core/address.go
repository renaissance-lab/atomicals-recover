package core

import (
	"encoding/hex"
	"fmt"

	btcec "github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
)

type Network int

const (
	MainNetWork = Network(1)
	TestNetWork = Network(2)
)

func GenerateP2TRAddress(keyStr, expectAddress string, net Network) (string, error) {

	wif, err := btcutil.DecodeWIF(keyStr)
	if err != nil {
		fmt.Printf("DecodeWIF %v, err %v ", keyStr, err)
		return "", err
	}

	pubkey := wif.PrivKey.PubKey()
	// fmt.Printf("pubkey %x\n", pubkey.SerializeCompressed())
	taprootKey := txscript.ComputeTaprootKeyNoScript(pubkey)
	netParam := &chaincfg.TestNet3Params
	if net == MainNetWork {
		netParam = &chaincfg.MainNetParams
	}
	tapScriptAddr, err := btcutil.NewAddressTaproot(
		schnorr.SerializePubKey(taprootKey), netParam,
	)
	if err != nil {
		return "", err
	}
	if tapScriptAddr.String() != expectAddress {
		errInfo := fmt.Sprintf("tapScriptAddr %v, expect %v", tapScriptAddr.String(), expectAddress)
		panic(errInfo)
	}

	return tapScriptAddr.String(), nil
}

func GenerateP2TRScriptAddress(keyStr, protocol, opType, bitworkr, bitworkc, mintTicker string, timeUnix, nonce int64, net Network) (string, string, error) {
	wif, err := btcutil.DecodeWIF(keyStr)
	if err != nil {
		fmt.Println("DecodeWIF err ", err)
		return "", "", err
	}

	pubkey := wif.PrivKey.PubKey()
	// fmt.Printf("pubkey %v\n", pubkey.EncodeToString())
	args := AtomArg{
		Time:       timeUnix,
		Nonce:      nonce,
		Bitworkc:   bitworkc,
		Bitworkr:   bitworkr,
		MintTicker: mintTicker,
	}

	return GetScriptP2TRAddress(protocol, opType, pubkey, args, net)
}

func GetScriptP2TRAddress(protocol, opType string, internalKey *btcec.PublicKey, args AtomArg, net Network) (string, string, error) {
	script, p2trPubkey, err := AppendMintUpdatedRevealScript(protocol, opType, internalKey, args)
	if err != nil {
		return "", "", err
	}

	chainNetParam := &chaincfg.TestNet3Params
	if net == MainNetWork {
		chainNetParam = &chaincfg.MainNetParams
	}

	address, err := btcutil.NewAddressTaproot(p2trPubkey, chainNetParam)
	if err != nil {
		return "", "", err
	}
	// fmt.Printf("address for net %v, is %v\n", net, address.EncodeAddress())
	return address.EncodeAddress(), hex.EncodeToString(script), nil
}

func GetPayLoad(time, nonce int64, bitworkc, bitworkr, mintTicker string) []byte {

	payLoad := AtomPayLoad{
		Args: AtomArg{
			Time:       time,
			Nonce:      nonce,
			Bitworkc:   bitworkc,
			MintTicker: mintTicker,
		},
	}
	return PayloadToCbor(payLoad)

}
func AppendMintUpdatedRevealScript(protocol, opType string, internalKey *btcec.PublicKey, args AtomArg) ([]byte, []byte, error) {

	builder := txscript.NewScriptBuilder()
	pubkey := schnorr.SerializePubKey(internalKey)
	builder.AddData(pubkey)

	builder.AddOp(txscript.OP_CHECKSIG)
	builder.AddOp(txscript.OP_0)
	builder.AddOp(txscript.OP_IF)

	builder.AddData([]byte(protocol))

	// optype
	builder.AddData([]byte(opType))
	// data
	payloadData := GetPayLoad(args.Time, args.Nonce, args.Bitworkc, args.Bitworkr, args.MintTicker)

	// fmt.Printf("payloadData %x\n", payloadData)
	builder.AddData([]byte(payloadData))

	// endif
	builder.AddOp(txscript.OP_ENDIF)

	pkScript, err := builder.Script()
	if err != nil {
		return nil, nil, err
	}
	// fmt.Printf("pkScript %x\n", pkScript)
	tapLeaf := txscript.NewBaseTapLeaf(pkScript)
	tapScriptTree := txscript.AssembleTaprootScriptTree(tapLeaf)

	tapScriptRootHash := tapScriptTree.RootNode.TapHash()
	outputKey := txscript.ComputeTaprootOutputKey(
		internalKey, tapScriptRootHash[:],
	)
	// p2trScript, err := txscript.PayToTaprootScript(outputKey)
	// if err != nil {
	// 	return nil, nil, err
	// }

	// fmt.Printf("p2trScript %x\n ", (p2trScript))
	calcPubkey := schnorr.SerializePubKey(outputKey)

	return pkScript, calcPubkey, nil

}
