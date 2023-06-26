package nodes

import (
	"context"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthereumNodeClient struct {
	client *ethclient.Client
}

func NewEthereumNodeClient(rpc string) *EthereumNodeClient {
	client, err := ethclient.Dial(rpc)

	if err != nil {
		panic(err)
	}

	return &EthereumNodeClient{client}
}

func (c *EthereumNodeClient) generateTypes(types []string) (abi.Arguments, error) {
	abiArguments := abi.Arguments{}

	for _, t := range types {
		if t == "" {
			continue
		}
		argType, err := abi.NewType(t, "", nil)
		if err != nil {
			return nil, err
		}
		abiArguments = append(abiArguments, abi.Argument{Type: argType})
	}

	return abiArguments, nil
}

func (c *EthereumNodeClient) generateMethodABI(functionName string, inputTypes []string, outputTypes []string) (*abi.ABI, error) {
	inputs, err := c.generateTypes(inputTypes)
	if err != nil {
		return nil, err
	}

	outputs, err := c.generateTypes(outputTypes)
	if err != nil {
		return nil, err
	}

	method := abi.NewMethod(functionName, functionName, abi.Function, "", false, false, inputs, outputs)
	return &abi.ABI{
		Methods: map[string]abi.Method{method.Name: method},
	}, nil
}

func (c *EthereumNodeClient) ExecuteReadFunction(context context.Context, contractAddress string, inputTypes []string, outputTypes []string, functionName string, params ...string) ([]interface{}, error) {

	abi, err := c.generateMethodABI(functionName, inputTypes, outputTypes)
	if err != nil {
		return nil, err
	}

	castedParams := make([]interface{}, len(params))
	for i, arg := range abi.Methods[functionName].Inputs {
		param := params[i]
		switch arg.Type.String() {
		case "address":
			if len(param) > 2 && param[0:2] == "0x" {
				castedParams[i] = common.HexToAddress(param)
			}
		case "uint256":
		case "int256":
			if num, err := strconv.ParseUint(param, 10, 64); err == nil {
				castedParams[i] = new(big.Int).SetUint64(num)
			}
		case "bool":
			if param == "true" {
				castedParams[i] = true
			} else {
				castedParams[i] = false
			}
		default:
			castedParams[i] = param
		}
	}

	payload, err := abi.Pack(functionName, castedParams...)
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
