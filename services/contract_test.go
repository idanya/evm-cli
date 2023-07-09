package services

import (
	"context"
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/idanya/evm-cli/clients/nodes/mocks"
	"github.com/idanya/evm-cli/clients/openchain"
	decompiler "github.com/idanya/evm-cli/decompiler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	decompilerClient = decompiler.NewDecompiler(openchain.NewClient())
)

func TestContractService_DetectMinimalProxyByByteCode(t *testing.T) {
	proxyBytecode := common.Hex2Bytes("363d3d373d3d3d363d734d11c446473105a02b5c1ab9ebe9b03f33902a295af43d82803e903d91602b57fd5bf3")
	nodeClientMock := mocks.NewNodeClient(t)

	nodeClientMock.On("GetContractCode", mock.Anything, "0x3348f2aee62a0ddb164c711b5937e4001c17080e").Return(proxyBytecode, nil)
	nodeClientMock.On("ExecuteReadFunction", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("error"))
	nodeClientMock.On("GetContractStorageSlot", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("error"))

	contractService := NewContractService(nodeClientMock, decompilerClient)
	implementation, err := contractService.GetProxyImplementation(context.Background(), "0x3348f2aee62a0ddb164c711b5937e4001c17080e")
	assert.NoError(t, err)
	assert.Equal(t, "0x4d11c446473105a02b5c1ab9ebe9b03f33902a29", implementation)
}

func TestContractService_DetectProxyByImplementationMethods(t *testing.T) {

	methods := []string{"implementation", "masterCopy", "childImplementation", "comptrollerImplementation"}
	for _, method := range methods {
		nodeClientMock := mocks.NewNodeClient(t)

		nodeClientMock.On("ExecuteReadFunction", mock.Anything, mock.Anything, mock.Anything, method).Return([]interface{}{common.HexToAddress("0xB650eb28d35691dd1BD481325D40E65273844F9b")}, nil)
		nodeClientMock.On("ExecuteReadFunction", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("error"))

		contractService := NewContractService(nodeClientMock, decompilerClient)
		implementation, err := contractService.GetProxyImplementation(context.Background(), "0x0000000000085d4780B73119b644AE5ecd22b376")
		assert.NoError(t, err)
		assert.Equal(t, "0xB650eb28d35691dd1BD481325D40E65273844F9b", implementation)
	}
}
