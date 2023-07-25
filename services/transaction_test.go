package services

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	dirmock "github.com/idanya/evm-cli/clients/directory/mocks"
	"github.com/idanya/evm-cli/clients/nodes"
	"github.com/idanya/evm-cli/clients/nodes/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTransactionService_GetTransactionReceipt(t *testing.T) {
	directoryClientMock := dirmock.NewDirectoryClient(t)
	decoder := NewDecoder(directoryClientMock)

	nodeClientMock := mocks.NewNodeClient(t)
	nodeGenerator := func() nodes.NodeClient { return nodeClientMock }

	transactionService := NewTransactionService(nodeGenerator, directoryClientMock, decoder)

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

func TestTransactionService_GetTransactionByHash(t *testing.T) {
	directoryClientMock := dirmock.NewDirectoryClient(t)
	decoder := NewDecoder(directoryClientMock)

	nodeClientMock := mocks.NewNodeClient(t)
	nodeGenerator := func() nodes.NodeClient { return nodeClientMock }

	transactionService := NewTransactionService(nodeGenerator, directoryClientMock, decoder)

	txHash := "0xec8ecd56dca115adcc8de346ffe054841f810964a68afc81faf764f8a0ae7c26"
	functionSignature := "safeTransferFrom(address,address,uint256,bytes)"
	payload := common.Hex2Bytes("b88d4fde0000000000000000000000005ab326a31b48faac615927dd7068b53423b32d8c000000000000000000000000daf57b900d6c8f9c8eea64a210fe9736a11ca9310000000000000000000000000000000000000000000000000000000000000fe200000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000000360c6ebe")
	transaction := types.NewTx(&types.DynamicFeeTx{Nonce: 1, Data: payload})

	// response := &entities.EnrichedTxInfo{DecodedData: &entities.DecodeResult{Hash: "0xb88d4fde", Method: "safeTransferFrom(address from,address to,uint256 tokenId,bytes data)"}}

	nodeClientMock.On("GetTransactionByHash", mock.Anything, txHash).Return(transaction, nil)
	directoryClientMock.On("LookupFunction", "0xb88d4fde").Return(functionSignature, nil)

	decoded, err := transactionService.GetTransactionByHash(context.Background(), txHash)
	assert.NoError(t, err)
	assert.Equal(t, functionSignature, decoded.DecodedData.Method)
	assert.Equal(t, "0xb88d4fde", decoded.DecodedData.Hash)
	assert.Equal(t, "address", decoded.DecodedData.Arguments[0].Type)
	assert.Equal(t, "0x5AB326a31b48faac615927dd7068B53423B32D8c", decoded.DecodedData.Arguments[0].Value.(common.Address).Hex())

	assert.Equal(t, "address", decoded.DecodedData.Arguments[1].Type)
	assert.Equal(t, "0xdaf57b900d6c8F9c8eeA64a210FE9736a11CA931", decoded.DecodedData.Arguments[1].Value.(common.Address).Hex())

	assert.Equal(t, "uint256", decoded.DecodedData.Arguments[2].Type)
	assert.Equal(t, big.NewInt(4066), decoded.DecodedData.Arguments[2].Value.(*big.Int))
	assert.Equal(t, "bytes", decoded.DecodedData.Arguments[3].Type)

}
