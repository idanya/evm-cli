package nodes

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
)

type NodeClient interface {
	GetTransactionByHash(context context.Context, txHash string) (*types.Transaction, error)
	GetTransactionReceipt(context context.Context, txHash string) (*types.Receipt, error)
	GetAccountNonce(context context.Context, account string) (uint64, error)
	GetContractCode(context context.Context, contractAddress string) ([]byte, error)
	ExecuteReadFunction(context context.Context, contractAddress string, inputTypes []string, outputTypes []string, functionName string, params ...string) ([]interface{}, error)
}


