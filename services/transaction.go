package services

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/idanya/evm-cli/clients/directory"
	"github.com/idanya/evm-cli/entities"
)

type TransactionService struct {
	nodeClientGenerator entities.NodeClientGenerator
	directoryClient     directory.DirectoryClient
	decoder             *Decoder
}

func NewTransactionService(nodeClientGenerator entities.NodeClientGenerator,
	directoryClient directory.DirectoryClient,
	decoder *Decoder) *TransactionService {
	return &TransactionService{nodeClientGenerator, directoryClient, decoder}
}

func (ts *TransactionService) GetTransactionReceipt(context context.Context, txHash string) (*entities.EnrichedReceipt, error) {
	receipt, err := ts.nodeClientGenerator().GetTransactionReceipt(context, txHash)
	if err != nil {
		return nil, err
	}

	enrichedLogs := make([]*entities.EnrichedReceiptLog, 0)
	for _, log := range receipt.Logs {
		event, err := ts.directoryClient.LookupEvent(log.Topics[0].String())
		if err != nil {
			continue
		}

		if err != nil {
			continue
		}
		enrichedLogs = append(enrichedLogs, &entities.EnrichedReceiptLog{Log: log, FunctionSelector: event})

	}

	enrichedReceipt := &entities.EnrichedReceipt{Receipt: receipt, Logs: enrichedLogs}

	return enrichedReceipt, nil
}

func (ts *TransactionService) GetTransactionByHash(context context.Context, txHash string) (*entities.EnrichedTxInfo, error) {
	transaction, err := ts.nodeClientGenerator().GetTransactionByHash(context, txHash)
	if err != nil {
		return nil, err
	}

	decoded, err := ts.decoder.DecodeContractCallData(context, common.Bytes2Hex(transaction.Data()))
	if err == nil {
		return &entities.EnrichedTxInfo{Transaction: transaction, DecodedData: decoded}, nil
	}

	return &entities.EnrichedTxInfo{Transaction: transaction}, nil
}
