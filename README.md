![Build](https://github.com/idanya/evm-cli/actions/workflows/go.yml/badge.svg?branch=main)
![release](https://img.shields.io/github/v/release/idanya/evm-cli)


# evm-cli
A command line tool for inspecting Ethereum smart contracts, transactions and accounts.

### Install

#### Homebrew
```
brew tap idanya/tools
brew install evm-cli
```
#### Go
```
go install github.com/idanya/evm-cli@latest
```

### Usage
```
$ evm-cli --help
A CLI tool to interact with the EVM blockchains via JSON-RPC

Usage:
  evm-cli [flags]
  evm-cli [command]

Available Commands:
  account     Account related commands
  completion  Generate the autocompletion script for the specified shell
  contract    Contract related commands
  help        Help about any command
  tx          Transaction related commands

Flags:
  -c, --chain-id uint    Chain ID of the blockchain (default 1)
  -h, --help             help for evm-cli
      --rpc-url string   node RPC endpoint (overrides the chain ID)

Use "evm-cli [command] --help" for more information about a command.
```

### Account commands
```
nonce       Get account nonce
```

### Transactions commands
```
info        Get transaction data by hash
receipt     Get transaction receipt by hash
```

### Contract commands
```
decode      Decode contract call data
exec        Run contract readonly method
func-list   Get function list
opcode      Get opcode
proxy       Get proxy implementation address
```
