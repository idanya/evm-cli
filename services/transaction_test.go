package services

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	dirmock "github.com/idanya/evm-cli/clients/directory/mocks"
	"github.com/idanya/evm-cli/clients/nodes/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTransactionService_GetTransactionReceipt(t *testing.T) {
	directoryClientMock := dirmock.NewDirectoryClient(t)

	nodeClientMock := mocks.NewNodeClient(t)
	transactionService := NewTransactionService(nodeClientMock, directoryClientMock)

	topicHash := common.HexToHash("0xdb80dd488acf86d17c747445b0eabb5d57c541d3bd7b6b87af987858e5066b2b")
	txHash := "0xec8ecd56dca115adcc8de346ffe054841f810964a68afc81faf764f8a0ae7c26"
	receipt := &types.Receipt{Logs: []*types.Log{{Topics: []common.Hash{topicHash}}}}

	nodeClientMock.On("GetTransactionReceipt", mock.Anything, txHash).Return(receipt, nil)
	directoryClientMock.On("LookupEvent", topicHash.String()).Return("LogMessageToL2(address,uint256,uint256,uint256[],uint256,uint256)", nil)

	decoded, err := transactionService.GetTransactionReceipt(context.Background(), txHash)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(decoded.Logs))
	assert.Equal(t, "LogMessageToL2(address,uint256,uint256,uint256[],uint256,uint256)", decoded.Logs[0].FunctionSelector)
}
