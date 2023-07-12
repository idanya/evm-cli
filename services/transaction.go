package services

import (
	"context"

	"github.com/idanya/evm-cli/clients/directory"
	"github.com/idanya/evm-cli/clients/nodes"
	"github.com/idanya/evm-cli/entities"
)

type TransactionService struct {
	nodeClient      nodes.NodeClient
	directoryClient directory.DirectoryClient
}

func NewTransactionService(nodeClient nodes.NodeClient, directoryClient directory.DirectoryClient) *TransactionService {
	return &TransactionService{nodeClient, directoryClient}
}

func (ts *TransactionService) GetTransactionReceipt(context context.Context, txHash string) (*entities.EnrichedReceipt, error) {
	receipt, err := ts.nodeClient.GetTransactionReceipt(context, txHash)
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
		enrichedLogs = append(enrichedLogs, &entities.EnrichedReceiptLog{log, event})

	}

	enrichedReceipt := &entities.EnrichedReceipt{receipt, enrichedLogs}

	return enrichedReceipt, nil
}
