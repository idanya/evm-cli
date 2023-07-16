package services

import (
	"context"
	"errors"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/idanya/evm-cli/clients/nodes"
	decompiler "github.com/idanya/evm-cli/decompiler"
)

type ContractService struct {
	nodeClient nodes.NodeClient
	decompiler *decompiler.Decompiler
	decoder    *Decoder
}

func NewContractService(nodeClient nodes.NodeClient, decompiler *decompiler.Decompiler, decoder *Decoder) *ContractService {
	return &ContractService{nodeClient, decompiler, decoder}
}

func (cs *ContractService) ExecuteReadFunction(context context.Context, contractAddress string, inputTypes []string, outputTypes []string, functionName string, params ...string) ([]interface{}, error) {
	abi, err := cs.generateMethodABI(functionName, inputTypes, outputTypes)
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
		case "uint256", "int256":
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

	return cs.nodeClient.ExecuteReadFunction(context, contractAddress, abi, functionName, castedParams...)
}

func (cs *ContractService) GetProxyImplementation(context context.Context, contractAddress string) (string, error) {
	address, err := cs.queryForProxyImplementation(context, contractAddress)
	if err == nil {
		return address, nil
	}

	address, err = cs.tryGetProxyImplementationByStorage(context, contractAddress)
	if err == nil {
		return address, nil
	}

	EIP_1167_BYTECODE_PREFIX := "363d3d373d3d3d363d73"
	EIP_1167_BYTECODE_SUFFIX := "5af43d82803e903d91602b57fd5bf3"

	contractCode, err := cs.nodeClient.GetContractCode(context, contractAddress)
	if err == nil {
		hexCode := common.Bytes2Hex(contractCode)
		if strings.HasPrefix(hexCode, EIP_1167_BYTECODE_PREFIX) && strings.HasSuffix(hexCode, EIP_1167_BYTECODE_SUFFIX) {
			return "0x" + hexCode[len(EIP_1167_BYTECODE_PREFIX):len(hexCode)-len(EIP_1167_BYTECODE_SUFFIX)], nil
		}
	}

	return "", errors.New("Could not find proxy contract")
}

func (cs *ContractService) queryForProxyImplementation(context context.Context, contractAddress string) (string, error) {
	methods := []string{"implementation", "masterCopy", "childImplementation", "comptrollerImplementation"}
	for _, method := range methods {
		response, err := cs.ExecuteReadFunction(context, contractAddress, []string{}, []string{"address"}, method)
		if err != nil {
			continue
		}
		return response[0].(common.Address).String(), nil
	}
	return "", errors.New("no implementation found")
}

func (cs *ContractService) tryGetProxyImplementationByStorage(context context.Context, contractAddress string) (string, error) {
	EIP_1967_LOGIC_SLOT := common.HexToHash("0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc")
	EIP_1967_BEACON_SLOT := common.HexToHash("0xa3f0ad74e5423aebfd80d3ef4346578335a9a72aeaee59ff6cb3582b35133d50")
	OPEN_ZEPPELIN_IMPLEMENTATION_SLOT := common.HexToHash("0x7050c9e0f4ca769c69bd3a8ef740bc37934f8e2c036e5a723fd8ee048ed3f8c3")
	EIP_1822_LOGIC_SLOT := common.HexToHash("0xc5f16f0fcc639fa48a6947836d9850f504798523bf8c9a3a87d5876cf622bcf7")

	storageSlots := []common.Hash{EIP_1967_LOGIC_SLOT, EIP_1967_BEACON_SLOT,
		OPEN_ZEPPELIN_IMPLEMENTATION_SLOT, EIP_1822_LOGIC_SLOT}

	for _, slot := range storageSlots {
		response, err := cs.nodeClient.GetContractStorageSlot(context, contractAddress, slot)
		if err != nil {
			continue
		}

		logicalAddress := common.BytesToAddress(response)

		if logicalAddress.String() != "0x0000000000000000000000000000000000000000" {
			return logicalAddress.String(), nil
		}
	}

	return "", errors.New("no implementation found")
}

func (cs *ContractService) generateTypes(types []string) (abi.Arguments, error) {
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

func (cs *ContractService) generateMethodABI(functionName string, inputTypes []string, outputTypes []string) (*abi.ABI, error) {
	inputs, err := cs.generateTypes(inputTypes)
	if err != nil {
		return nil, err
	}

	outputs, err := cs.generateTypes(outputTypes)
	if err != nil {
		return nil, err
	}

	method := abi.NewMethod(functionName, functionName, abi.Function, "", false, false, inputs, outputs)
	return &abi.ABI{
		Methods: map[string]abi.Method{method.Name: method},
	}, nil
}
