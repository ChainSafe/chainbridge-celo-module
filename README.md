# Chainbridge celo module
<a href="https://discord.gg/ykXsJKfhgq">
  <img alt="discord" src="https://img.shields.io/discord/593655374469660673?label=Discord&logo=discord&style=flat" />
</a>

Chainbridge celo module is the part of Chainbrige-core framework. This module brings support of celo compatible client module.

*Project still in deep beta*
- Chat with us on [discord](https://discord.gg/ykXsJKfhgq).

### Table of Contents

1. [Installation](#installation)
2. [Usage](#usage)
3. [Differences Between EVM and Celo](#differences-between-evm-and-celo)

## Installation
Refer to [installation](https://github.com/ChainSafe/chainbridge-docs/blob/develop/docs/installation.md) guide for assistance in installing.

## Usage
Module should be used along with core [framework](https://github.com/ChainSafe/chainbridge-core).

Since chainbridge-celo-module is a package it will require writing some extra code to get it running alongside [chainbridge-core](https://github.com/ChainSafe/chainbridge-core). Here you can find some examples 
[Example](https://github.com/ChainSafe/chainbridge-core-example)

### Differences Between EVM and Celo

Though Celo is an EVM-compatible chain, it deviates in its implementation of the original Ethereum specifications, and therefore is deserving of its own separate module.

The differences alluded to above in how Celo constructs transactions versus those found within Ethereum can be viewed below by taking a look at the Message structs in both implementations.

[Ethereum Message Struct](https://github.com/ethereum/go-ethereum/blob/ac7baeab57405c64592b1646a91e0a2bb33d8d6c/core/types/transaction.go#L586-L598)

Here you will find fields relating to the most recent London hardfork (EIP-1559), most notably `gasFeeCap` and `gasTipCap`.

```go
Message {
   from:       from,
   to:         to,
   nonce:      nonce,
   amount:     amount,
   gasLimit:   gasLimit,
   gasPrice:   gasPrice,
   gasFeeCap:  gasFeeCap,
   gasTipCap:  gasTipCap,
   data:       data,
   accessList: accessList,
   isFake:     isFake,
}
```

[Celo Message Struct](https://github.com/ChainSafe/chainbridge-celo-module/blob/b6d7ad422a5356500d2d5cf0b98e00da86dbb42e/transaction/tx.go#L422-L435)

In Celo's struct you will notice that there are additional fields added for `feeCurrency`, `gatewayFeeRecipient` and `gatewayFee`. You may also notice the `ethCompatible` field, a boolean value we added in order to quickly determine whether the message is Ethereum compatible or not, ie, that `feeCurrency`, `gatewayFeeRecipient` and `gatewayFee` are omitted.

```go
Message {
   from:                from,
   to:                  to,
   nonce:               nonce,
   amount:              amount,
   gasLimit:            gasLimit,
   gasPrice:            gasPrice,
   feeCurrency:         feeCurrency,         // Celo-specific
   gatewayFeeRecipient: gatewayFeeRecipient, // Celo-specific
   gatewayFee:          gatewayFee,          // Celo-specific
   data:                data,
   ethCompatible:       ethCompatible,       // Bool to check presence of: feeCurrency, gatewayFeeRecipient, gatewayFee
   checkNonce:          checkNonce,
}
```

# ChainSafe Security Policy

## Reporting a Security Bug

We take all security issues seriously, if you believe you have found a security issue within a ChainSafe
project please notify us immediately. If an issue is confirmed, we will take all necessary precautions
to ensure a statement and patch release is made in a timely manner.

Please email us a description of the flaw and any related information (e.g. reproduction steps, version) to
[security at chainsafe dot io](mailto:security@chainsafe.io).

## License

_GNU Lesser General Public License v3.0_