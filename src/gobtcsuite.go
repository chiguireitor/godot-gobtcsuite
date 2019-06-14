package main

import (
  "fmt"
  "bytes"
  "encoding/base64"

	"github.com/shadowapex/godot-go/gdnative"

  "github.com/btcsuite/btcd/btcec"
  "github.com/btcsuite/btcd/chaincfg"
  "github.com/btcsuite/btcd/chaincfg/chainhash"
  "github.com/btcsuite/btcd/wire"
  "github.com/btcsuite/btcutil/hdkeychain"
)

type GoBtcSuite struct {
	base gdnative.Object
  network *chaincfg.Params
  seed *hdkeychain.ExtendedKey
}

var Instances = map[string]*GoBtcSuite{}

func init() {
  gdnative.SetNativeScriptInit(nativeScriptInit)
}

func nativeScriptInit() {

  createFunc := gdnative.InstanceCreateFunc{
		CreateFunc: gobtcsuiteConstructor,
		MethodData: "GoBtcSuite",
		FreeFunc:   func(methodData string) {},
  }

	destroyFunc := gdnative.InstanceDestroyFunc{
		DestroyFunc: gobtcsuiteDestructor,
		MethodData:  "GoBtcSuite",
		FreeFunc:    func(methodData string) {},
	}

	gdnative.NativeScript.RegisterClass(
		"GoBtcSuite",
		"Reference",
		&createFunc,
		&destroyFunc,
	)

	gdnative.NativeScript.RegisterMethod(
		"GoBtcSuite",
		"gen_masterseed",
		&gdnative.MethodAttributes{
			RPCType: gdnative.MethodRpcModeDisabled,
		},
		&gdnative.InstanceMethod{
			Method:     genMasterSeed,
			MethodData: "gen_masterseed",
			FreeFunc:   func(methodData string) {},
		},
  )

  gdnative.NativeScript.RegisterMethod(
		"GoBtcSuite",
		"load_masterseed",
		&gdnative.MethodAttributes{
			RPCType: gdnative.MethodRpcModeDisabled,
		},
		&gdnative.InstanceMethod{
			Method:     masterSeedFromString,
			MethodData: "load_masterseed",
			FreeFunc:   func(methodData string) {},
		},
  )

  gdnative.NativeScript.RegisterMethod(
		"GoBtcSuite",
		"set_network",
		&gdnative.MethodAttributes{
			RPCType: gdnative.MethodRpcModeDisabled,
		},
		&gdnative.InstanceMethod{
			Method:     setNetwork,
			MethodData: "set_network",
			FreeFunc:   func(methodData string) {},
		},
  )

  gdnative.NativeScript.RegisterMethod(
		"GoBtcSuite",
		"get_address",
		&gdnative.MethodAttributes{
			RPCType: gdnative.MethodRpcModeDisabled,
		},
		&gdnative.InstanceMethod{
			Method:     genAddress,
			MethodData: "get_address",
			FreeFunc:   func(methodData string) {},
		},
  )

  gdnative.NativeScript.RegisterMethod(
		"GoBtcSuite",
		"sign_message",
		&gdnative.MethodAttributes{
			RPCType: gdnative.MethodRpcModeDisabled,
		},
		&gdnative.InstanceMethod{
			Method:     signMessage,
			MethodData: "sign_message",
			FreeFunc:   func(methodData string) {},
		},
  )
}

func gobtcsuiteConstructor(object gdnative.Object, methodData string) string {
	// Create a new instance of our struct.
	instance := &GoBtcSuite{
		base: object,
	}

	// Use the pointer address as the instance ID
	instanceID := fmt.Sprintf("%p", instance)
	/*Instances[instanceID] = instance*/
  Instances[object.ID()] = instance

	// Return the instanceID
	return instanceID
}

func gobtcsuiteDestructor(object gdnative.Object, methodData, userData string) {
	// Delete the instance from our map of instances
	delete(Instances, userData)
}

func setNetwork(object gdnative.Object, methodData, userData string, numArgs int, args []gdnative.Variant) gdnative.Variant {
  instance := Instances[object.ID()]

  net := string(args[0].AsString())

  if (net == "mainnet") {
    instance.network = &chaincfg.MainNetParams
  } else if (net == "testnet" || net == "testnet3") {
    instance.network = &chaincfg.TestNet3Params
  } else if (net == "regtest") {
    instance.network = &chaincfg.RegressionNetParams
  }

  data := gdnative.NewStringWithWideString(instance.network.Name)
	ret := gdnative.NewVariantWithString(data)

  return ret
}

func genMasterSeed(object gdnative.Object, methodData, userData string, numArgs int, args []gdnative.Variant) gdnative.Variant {
  instance := Instances[object.ID()]

  seed, err := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
	if err != nil {
		fmt.Println(err)
		return gdnative.NewVariantNil()
	}

	// Generate a new master node using the seed.
	key, err := hdkeychain.NewMaster(seed, instance.network)
	if err != nil {
		fmt.Println(err)
		return gdnative.NewVariantNil()
	}

  instance.seed = key

  data := gdnative.NewStringWithWideString(key.String())
	ret := gdnative.NewVariantWithString(data)

	return ret
}

func masterSeedFromString(object gdnative.Object, methodData, userData string, numArgs int, args []gdnative.Variant) gdnative.Variant {
  instance := Instances[object.ID()]

	// Generate a new master node using the seed.
	key, err := hdkeychain.NewKeyFromString(string(args[0].AsString()))
	if err != nil {
		fmt.Println(err)
		return gdnative.NewVariantNil()
	}

  instance.seed = key

  data := gdnative.NewStringWithWideString(key.String())
	ret := gdnative.NewVariantWithString(data)

	return ret
}

func genAddress(object gdnative.Object, methodData, userData string, numArgs int, args []gdnative.Variant) gdnative.Variant {
  instance := Instances[object.ID()]

  addr, err := instance.seed.Address(instance.network)
  if err != nil {
    fmt.Println(err)
    return gdnative.NewVariantNil()
  }

  data := gdnative.NewStringWithWideString(addr.String())
	ret := gdnative.NewVariantWithString(data)

	return ret
}

func signMessage(object gdnative.Object, methodData, userData string, numArgs int, args []gdnative.Variant) gdnative.Variant {
  instance := Instances[object.ID()]
  msg := string(args[0].AsString())
  
  var buf bytes.Buffer
	wire.WriteVarString(&buf, 0, "Bitcoin Signed Message:\n")
	wire.WriteVarString(&buf, 0, msg)
  messageHash := chainhash.DoubleHashB(buf.Bytes())

  key, err := instance.seed.ECPrivKey()
  if err != nil {
		return gdnative.NewVariantNil()
	}

  sigbytes, err := btcec.SignCompact(btcec.S256(), key, messageHash, true)
	if err != nil {
		return gdnative.NewVariantNil()
	}

  str := base64.StdEncoding.EncodeToString(sigbytes)

  data := gdnative.NewStringWithWideString(str)
	ret := gdnative.NewVariantWithString(data)

	return ret
}

func main() {
}
