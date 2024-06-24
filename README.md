# Optimism Plasma DA KVStore Store implementation(irys)

## Introduction

This simple DA KVStore Store implementation based on irys. See the [Irys doc](https://docs.irys.xyz/overview/about) for more information

## Install
```bash
go get github.com/4everland/opplasmairys
```

## Start

```go
import (
  "github.com/4everland/opplasmairys"
)

config := opplasmairys.Config{
  //IrysPrivateKey
  IrysPrivateKey: os.Getenv("IRYS_PRIVATE_KEY"),
  //IrysPaymentNetwork supported networks: ethereum,matic,bnb,avalanche,arbitrum or fantom
  IrysPaymentNetwork: os.Getenv("IRYS_PAYMENT_NETWORK"),
  //IrysPaymentRPC, optional
  IrysPaymentRPC: os.Getenv("IRYS_PAYMENT_RPC"),
  //IrysFreeUpload, optional
  IrysFreeUpload: os.Getenv("IRYS_FREE_UPLOAD"),
}   
var (
  err error
  //backup kvstore, optional
  store opplasmairys.KVStore
)
//create a new irys store
store, err = opplasmairys.NewDAStore(config, store)
```

## Simple Server

This simple DA server implementation supports irys. See the [daserver](./daserver/cmd/README.md) for more information


