package nodes

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthereumNodeClient struct {
	client *ethclient.Client
	rpc    string
}

func NewEthereumNodeClient(rpc string) *EthereumNodeClient {
	client, err := ethclient.Dial(rpc)

	if err != nil {
		panic(err)
	}

	return &EthereumNodeClient{client, rpc}
}

func (c *EthereumNodeClient) ExecuteReadFunction(context context.Context, contractAddress string, abi *abi.ABI, functionName string, params ...interface{}) ([]interface{}, error) {
	payload, err := abi.Pack(functionName, params...)
	if err != nil {
		return nil, err
	}

	addr := common.HexToAddress(contractAddress)
	response, err := c.client.CallContract(context, ethereum.CallMsg{To: &addr, Data: payload}, nil)
	if err != nil {
		return nil, err
	}

	unpacked, err := abi.Unpack(functionName, response)

	if err != nil {
		return nil, err
	}

	return unpacked, err
}

func (c *EthereumNodeClient) GetContractCode(context context.Context, contractAddress string) ([]byte, error) {
	return c.client.CodeAt(context, common.HexToAddress(contractAddress), nil)
}

func (c *EthereumNodeClient) GetContractStorageSlot(context context.Context, contractAddress string, key common.Hash) ([]byte, error) {
	return c.client.StorageAt(context, common.HexToAddress(contractAddress), key, nil)
}

func (c *EthereumNodeClient) GetAccountNonce(context context.Context, account string) (uint64, error) {
	return c.client.NonceAt(context, common.HexToAddress(account), nil)
}

func (c *EthereumNodeClient) GetTransactionByHash(context context.Context, txHash string) (*types.Transaction, error) {
	tx, _, err := c.client.TransactionByHash(context, common.HexToHash(txHash))
	return tx, err
}

func (c *EthereumNodeClient) GetTransactionReceipt(context context.Context, txHash string) (*types.Receipt, error) {
	receipt, err := c.client.TransactionReceipt(context, common.HexToHash(txHash))
	return receipt, err
}
