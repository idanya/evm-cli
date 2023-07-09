package decomplier

import (
	"fmt"
	"log"
	"sync"

	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/idanya/evm-cli/clients/openchain"
)

type TranslatedFunction struct {
	Hash      string
	Signature string
}

func (t *TranslatedFunction) String() string {
	return fmt.Sprintf("[%s] %s", t.Hash, t.Signature)
}

type Decompiler struct {
	openchainClient *openchain.Client
}

func NewDecompiler(openchainClient *openchain.Client) *Decompiler {
	return &Decompiler{openchainClient}
}

func (d *Decompiler) Decompile(bytecode []byte) ([]*TranslatedFunction, error) {
	translated := make([]*TranslatedFunction, 0)

	funcs, err := d.extractPush4bytes(bytecode)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup

	log.Printf("Found %d functions, matching 4bytes with openchain...", len(funcs))
	for _, f := range funcs {
		wg.Add(1)

		go func(f string) {
			defer wg.Done()

			lookupFunc, err := d.openchainClient.LookupFunction(f)
			if err == nil {
				translated = append(translated, &TranslatedFunction{Hash: f, Signature: lookupFunc.Name})
			} else {
				translated = append(translated, &TranslatedFunction{Hash: f, Signature: "unknown"})
			}
		}(f)
	}

	wg.Wait()

	return translated, nil
}

func (d *Decompiler) extractPush4bytes(bytecode []byte) ([]string, error) {
	hashes := make(map[string]string, 0)

	it := NewInstructionIterator(bytecode)
	for it.Next() {
		functionHash := fmt.Sprintf("%#x", it.Arg())

		if it.op == vm.PUSH4 && functionHash != "0xffffffff" {
			if _, ok := hashes[functionHash]; !ok {
				hashes[functionHash] = functionHash
			}
		}
	}
	// if err := it.Error(); err != nil {
	// 	return nil, err
	// }

	keys := make([]string, 0, len(hashes))
	for k := range hashes {
		keys = append(keys, k)
	}

	return keys, nil
}

func (d *Decompiler) Disassemble(script []byte) ([]string, error) {
	instrs := make([]string, 0)

	it := NewInstructionIterator(script)
	for it.Next() {
		if it.Arg() != nil && 0 < len(it.Arg()) {
			instrs = append(instrs, fmt.Sprintf("%05x: %v %#x\n", it.PC(), it.Op(), it.Arg()))
		} else {
			instrs = append(instrs, fmt.Sprintf("%05x: %v\n", it.PC(), it.Op()))
		}
	}
	if err := it.Error(); err != nil {
		instrs = append(instrs, fmt.Sprintf("[ERROR] %05x: %v err: %v\n", it.PC(), it.Op(), err))
	}
	return instrs, nil
}
