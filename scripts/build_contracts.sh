#!/usr/bin/env bash
# Copyright 2020 ChainSafe Systems
# SPDX-License-Identifier: LGPL-3.0-only

CONTRACTS_REPO="https://github.com/ChainSafe/chainbridge-celo-solidity/"
CONTRACTS_BRANCH="main"
CONTRACTS_COMMIT="70c97947f6b05397c2887fe8213a3d49d3d06de7"
CONTRACTS_DIR="./solidity"
DEST_DIR="./bindings"

CELO_BRANCH="v1.3.2"


base_path="./build/bindings"
BIN_DIR="$base_path/bin"
ABI_DIR="$base_path/abi"
RUNTIME_DIR="$base_path/runtime"
GO_DIR="$base_path/go"

set -eux

case $TARGET in
	"build")
		git clone -b $CONTRACTS_BRANCH $CONTRACTS_REPO $CONTRACTS_DIR
#    git checkout $CONTRACTS_COMMIT

    if [ -x "$(command -v celo-ganache)" ]
	then
	  echo "celo-ganache found, skipping install"
	else
	  git clone --depth 1 https://github.com/celo-org/ganache-cli.git
	  npm install --prefix ./ganache-cli
 	  mkdir -p ~/.local/bin
	  ln -f -s  $PWD/ganache-cli/cli.js  ~/.local/bin/celo-ganache
	fi

	if [ -x "$(command -v celo-abigen)" ]
	then
	  echo "celo-abigen found, skipping install"
	else
	  git clone --depth 1 https://github.com/celo-org/celo-blockchain.git --branch $CELO_BRANCH --single-branch
	  cd celo-blockchain
	  make geth
	  env GOBIN=$PWD/build/bin go install ./cmd/abigen
	  mkdir -p ~/.local/bin
	  ln -f -s  $PWD/build/bin/abigen  ~/.local/bin/celo-abigen
	fi

	pushd $CONTRACTS_DIR
    make bindings

    popd

    mkdir $DEST_DIR
    cp -r $CONTRACTS_DIR/build/bindings/go/* $DEST_DIR
		;;

	"cli-only")
		git clone -b $CONTRACTS_BRANCH $CONTRACTS_REPO $CONTRACTS_DIR
    pushd $CONTRACTS_DIR
    git checkout $CONTRACTS_COMMIT
		;;

esac
