package services

import (
	"context"
	"errors"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	dirmock "github.com/idanya/evm-cli/clients/directory/mocks"
	"github.com/idanya/evm-cli/clients/nodes"
	"github.com/idanya/evm-cli/clients/nodes/mocks"
	decompiler "github.com/idanya/evm-cli/decompiler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestContractService_DetectMinimalProxyByByteCode(t *testing.T) {
	directoryClientMock := dirmock.NewDirectoryClient(t)
	decoder := NewDecoder(directoryClientMock)
	decompilerClient := decompiler.NewDecompiler(directoryClientMock)

	proxyBytecode := common.Hex2Bytes("363d3d373d3d3d363d734d11c446473105a02b5c1ab9ebe9b03f33902a295af43d82803e903d91602b57fd5bf3")
	nodeClientMock := mocks.NewNodeClient(t)

	nodeClientMock.On("GetContractCode", mock.Anything, "0x3348f2aee62a0ddb164c711b5937e4001c17080e").Return(proxyBytecode, nil)
	nodeClientMock.On("ExecuteReadFunction", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("error"))
	nodeClientMock.On("GetContractStorageSlot", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("error"))
	nodeGenerator := func() nodes.NodeClient { return nodeClientMock }

	contractService := NewContractService(nodeGenerator, decompilerClient, decoder)
	implementation, err := contractService.GetProxyImplementation(context.Background(), "0x3348f2aee62a0ddb164c711b5937e4001c17080e")
	assert.NoError(t, err)
	assert.Equal(t, "0x4d11c446473105a02b5c1ab9ebe9b03f33902a29", implementation)
}

func TestContractService_DetectProxyByImplementationMethods(t *testing.T) {
	directoryClientMock := dirmock.NewDirectoryClient(t)
	decoder := NewDecoder(directoryClientMock)
	decompilerClient := decompiler.NewDecompiler(directoryClientMock)

	methods := []string{"implementation", "masterCopy", "childImplementation", "comptrollerImplementation"}
	for _, method := range methods {
		nodeClientMock := mocks.NewNodeClient(t)

		nodeClientMock.On("ExecuteReadFunction", mock.Anything, mock.Anything, mock.Anything, method).Return([]interface{}{common.HexToAddress("0xB650eb28d35691dd1BD481325D40E65273844F9b")}, nil)
		nodeClientMock.On("ExecuteReadFunction", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("error"))
		nodeGenerator := func() nodes.NodeClient { return nodeClientMock }

		contractService := NewContractService(nodeGenerator, decompilerClient, decoder)
		implementation, err := contractService.GetProxyImplementation(context.Background(), "0x0000000000085d4780B73119b644AE5ecd22b376")
		assert.NoError(t, err)
		assert.Equal(t, "0xB650eb28d35691dd1BD481325D40E65273844F9b", implementation)
	}
}

func TestExecuteReadFunction(t *testing.T) {
	directoryClientMock := dirmock.NewDirectoryClient(t)
	decoder := NewDecoder(directoryClientMock)
	decompilerClient := decompiler.NewDecompiler(directoryClientMock)

	nodeClientMock := mocks.NewNodeClient(t)
	nodeClientMock.On("ExecuteReadFunction", mock.Anything, "0x0", mock.Anything, "func",
		common.HexToAddress("0xdac17f958d2ee523a2206206994597c13d831ec7"),
		new(big.Int).SetUint64(10), new(big.Int).SetUint64(100), mock.Anything).Return([]interface{}{"OK"}, nil)
	nodeGenerator := func() nodes.NodeClient { return nodeClientMock }

	contractService := NewContractService(nodeGenerator, decompilerClient, decoder)
	response, err := contractService.ExecuteReadFunction(context.Background(), "0x0",
		[]string{"address", "uint256", "int256", "bool"},
		[]string{"address"}, "func", "0xdac17f958d2ee523a2206206994597c13d831ec7", "10", "100", "false")

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "OK", response[0])
}

func TestGetContractStandards(t *testing.T) {
	file, _ := os.ReadFile("../assets/FiatTokenV2.hex")
	erc20Code := common.Hex2Bytes(string(file))

	directoryClientMock := dirmock.NewDirectoryClient(t)
	decoder := NewDecoder(directoryClientMock)
	decompilerClient := decompiler.NewDecompiler(directoryClientMock)

	nodeClientMock := mocks.NewNodeClient(t)
	nodeClientMock.On("GetContractCode", mock.Anything, "0x0").Return(erc20Code, nil)
	nodeGenerator := func() nodes.NodeClient { return nodeClientMock }

	contractService := NewContractService(nodeGenerator, decompilerClient, decoder)
	standards, err := contractService.GetContractStandards(context.Background(), "0x0")
	assert.Nil(t, err)
	assert.NotNil(t, standards)
	assert.Equal(t, 1, len(standards))
	assert.Equal(t, "ERC20", standards[0])
}
