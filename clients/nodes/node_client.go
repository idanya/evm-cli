package nodes

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

//go:generate mockery --name NodeClient
type NodeClient interface {
	GetTransactionByHash(context context.Context, txHash string) (*types.Transaction, error)
	GetTransactionReceipt(context context.Context, txHash string) (*types.Receipt, error)
	GetAccountNonce(context context.Context, account string) (uint64, error)
	GetContractCode(context context.Context, contractAddress string) ([]byte, error)
	ExecuteReadFunction(context context.Context, contractAddress string, abi *abi.ABI, functionName string, params ...interface{}) ([]interface{}, error)
	GetContractStorageSlot(context context.Context, contractAddress string, key common.Hash) ([]byte, error)
}
