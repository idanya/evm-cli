package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"gitlab.com/fireblocks/web3/utils/evm-cli/clients/nodes"
	decompiler "gitlab.com/fireblocks/web3/utils/evm-cli/decompiler"
)

const DefaultEditor = "vim"

type ContractCommands struct {
	clientFactory nodes.NodeClientFactoryFunc
	decompiler    *decompiler.Decompiler
}

func NewContractCommands(clientFactory nodes.NodeClientFactoryFunc, decompiler *decompiler.Decompiler) *ContractCommands {
	return &ContractCommands{clientFactory, decompiler}
}

func (tx *ContractCommands) GetRootCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "contract",
		Short: "Contract related commands",
	}
}

func openInEditor(text []byte) error {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "evm-cli-")
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}

	defer os.Remove(tmpFile.Name())

	if _, err = tmpFile.Write(text); err != nil {
		log.Fatal("Failed to write to temporary file", err)
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = DefaultEditor
	}

	executable, err := exec.LookPath(editor)
	if err != nil {
		log.Fatal("Cannot find editor", err)
	}

	cmd := exec.Command(executable, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (cc *ContractCommands) GetContractOpCodeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "opcode <address>",
		Short: "Get opcode",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			chainId, err := cmd.Flags().GetUint("chain-id")
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("Fetching contract bytecode")
			bytecode, err := cc.clientFactory(chainId).GetContractCode(context.Background(), args[0])
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

func (cc *ContractCommands) GetContractFunctionListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "func-list <address>",
		Short: "Get function list",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			chainId, err := cmd.Flags().GetUint("chain-id")
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("Fetching contract bytecode")
			bytecode, err := cc.clientFactory(chainId).GetContractCode(context.Background(), args[0])
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
		},
	}
}

// func (cc *ContractCommands) GetContractABICommand() *cobra.Command {
// 	return &cobra.Command{
// 		Use:   "abi <address>",
// 		Short: "Get contract ABI",
// 		Args:  cobra.ExactArgs(1),
// 		Run: func(cmd *cobra.Command, args []string) {
// 			contractABI, err := cc.client.GetContractABI(context.Background(), args[0])
// 			if err != nil {
// 				log.Fatal(err)
// 			}

// 			var prettyJSON bytes.Buffer
// 			err = json.Indent(&prettyJSON, []byte(contractABI), "", "  ")

// 			if err != nil {
// 				log.Fatal(err)
// 			}

// 			err = openInEditor(prettyJSON.Bytes())
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 		},
// 	}
// }

func (cc *ContractCommands) GetContractExecCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "exec <address> <function> <args>",
		Short: "Run contract readonly method",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			inputTypes, _ := cmd.Flags().GetString("input-types")
			outputTypes, _ := cmd.Flags().GetString("output-types")

			chainId, err := cmd.Flags().GetUint("chain-id")
			if err != nil {
				log.Fatal(err)
			}
			response, err := cc.clientFactory(chainId).ExecuteReadFunction(context.Background(), args[0], strings.Split(inputTypes, ","), strings.Split(outputTypes, ","), args[1], args[2:]...)
			if err != nil {
				log.Fatal(err)
			}

			data, _ := json.MarshalIndent(response, "", "  ")
			fmt.Println(string(data))
		},
	}

	cmd.PersistentFlags().String("input-types", "", "Input types (comma separated)")
	cmd.PersistentFlags().String("output-types", "", "Output types (comma separated)")

	return cmd
}
