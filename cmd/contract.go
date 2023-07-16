package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/idanya/evm-cli/clients/directory/openchain"
	decompiler "github.com/idanya/evm-cli/decompiler"
	"github.com/idanya/evm-cli/services"
	"github.com/spf13/cobra"
)

const DefaultEditor = "vim"

type ContractCommands struct {
	decompiler *decompiler.Decompiler
}

func NewContractCommands(decompiler *decompiler.Decompiler) *ContractCommands {
	return &ContractCommands{decompiler}
}

func (cc *ContractCommands) GetRootCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "contract",
		Short: "Contract related commands",
	}

	command.AddCommand(cc.GetContractOpCodeCommand())
	command.AddCommand(cc.GetContractFunctionListCommand())
	command.AddCommand(cc.GetContractExecCommand())
	command.AddCommand(cc.GetDecodeCallDataCommand())
	command.AddCommand(cc.GetContractProxyImplementationCommand())

	return command
}

func (cc *ContractCommands) GetContractOpCodeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "opcode <address>",
		Short: "Get opcode",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("Fetching contract bytecode")
			bytecode, err := NodeClientFromViper().GetContractCode(context.Background(), args[0])
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("Decoding...")
			opcode, err := cc.decompiler.Disassemble(bytecode)
			if err != nil {
				log.Fatal(err)
			}

			openInEditor([]byte(strings.Join(opcode, "\n")))
		},
	}
}

func (cc *ContractCommands) printContractFunctions(contractAddress string) {
	log.Printf("Fetching contract bytecode")
	bytecode, err := NodeClientFromViper().GetContractCode(context.Background(), contractAddress)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Decoding 4byte function list")
	funcList, err := cc.decompiler.Decompile(bytecode)
	if err != nil {
		log.Fatal(err)
	}

	for _, function := range funcList {
		fmt.Println(function.String())
	}
}

func (cc *ContractCommands) GetContractFunctionListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "func-list <address>",
		Short: "Get function list",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			decoder := services.NewDecoder(openchain.NewClient())
			contractService := services.NewContractService(NodeClientFromViper(), cc.decompiler, decoder)

			log.Printf("Checking if contract is proxy...")
			implementationAddress, err := contractService.GetProxyImplementation(context.Background(), args[0])
			if err == nil && implementationAddress != "" {
				log.Printf("Contract is proxy to %s, getting implementation functions...", implementationAddress)
				cc.printContractFunctions(implementationAddress)
				log.Printf("Getting proxy functions...")
			} else {
				log.Print("Contract is not proxy, getting functions...")
			}

			cc.printContractFunctions(args[0])

		},
	}
}

func (cc *ContractCommands) GetContractProxyImplementationCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "proxy <address>",
		Short: "Get proxy implementation address",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			decoder := services.NewDecoder(openchain.NewClient())
			contractService := services.NewContractService(NodeClientFromViper(), cc.decompiler, decoder)
			implementationAddress, err := contractService.GetProxyImplementation(context.Background(), args[0])
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("Implementation address: %s", implementationAddress)
		},
	}
}

func (cc *ContractCommands) GetDecodeCallDataCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "decode <calldata>",
		Short: "Decode contract call data",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			decoder := services.NewDecoder(openchain.NewClient())
			decoded, err := decoder.DecodeContractCallData(context.Background(), args[0])
			if err != nil {
				log.Fatal(err)
			}

			decodedJson, err := json.MarshalIndent(decoded, "", "  ")
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Decoded call data:\n\n%s", decodedJson)
		},
	}
}

func (cc *ContractCommands) GetContractExecCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "exec",
		Short:   "Run contract readonly method",
		Example: "exec <address> \"<methodName(inType1,inType2,...)(outType1)>\" <method args>",
		Args:    cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			contractAddress := args[0]
			contractMethod := args[1]
			methodParams := args[2:]

			var re = regexp.MustCompile(`(?m)(.*)?\((.*)?\)\((.*)?\)`)
			matches := re.FindStringSubmatch(contractMethod)

			if len(matches) != 4 {
				log.Fatal("Invalid method format. Should be methodName(inType1,inType2,...)(outType1)")
			}

			methodName := matches[1]
			methodTypes := matches[2]
			outputTypes := matches[3]

			decoder := services.NewDecoder(openchain.NewClient())
			contractService := services.NewContractService(NodeClientFromViper(), cc.decompiler, decoder)

			response, err := contractService.ExecuteReadFunction(context.Background(), contractAddress, strings.Split(methodTypes, ","), strings.Split(outputTypes, ","), methodName, methodParams...)
			if err != nil {
				log.Fatal(err)
			}

			data, _ := json.MarshalIndent(response, "", "  ")
			log.Println(string(data))
		},
	}
	return cmd
}
