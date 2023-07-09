# evm-cli
A command line tool for inspecting Ethereum smart contracts, transactions and accounts.

### Install
go install github.com/idanya/evm-cli@latest

### TODO
- [X] Detect minimal proxy by comparing on-chain bytecode to the minimal proxy template
- [ ] optimize concurrency when detecting proxy
- [ ] mock node responses for tests
- [ ] `inspect` command for complete analysis of a contract
- [ ] Extract dispatched logs (PUSH32 ?)
