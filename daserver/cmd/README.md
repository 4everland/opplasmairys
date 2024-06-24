# Plasma DA Server

## Introduction

This simple DA server implementation supports irys.
See the [Irys doc](https://docs.irys.xyz/overview/about) for more information on how to configure the irys client.

The DA server implementation supports local storage via file based storage and remote via S3 refer to [op-optimism](https://github.com/ethereum-optimism/optimism/blob/develop/op-plasma/cmd/daserver/README.md) for more information

## Irys Configuration

Depending on your cloud provider a wide array of configurations are available. The S3 client will
load configurations from the environment, shared credentials and shared config files.
Sample environment variables are provided below:

```bash
export IRYS_PRIVATE_KEY=YOUR_IRYS_PRIVATE_KEY
export IRYS_PAYMENT_NETWORK=YOUR_IRYS_PAYMENT_NETWORK
export IRYS_PAYMENT_RPC=YOUR_IRYS_PAYMENT_RPC
export IRYS_PAYMENT_RPC=true
```

