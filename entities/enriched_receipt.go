package entities

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type EnrichedReceiptLog struct {
	*types.Log
	FunctionSelector string `json:"selector"`
}

func (er *EnrichedReceiptLog) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		FunctionSelector string         `json:"selector"`
		Address          common.Address `json:"address" gencodec:"required"`
		Topics           []common.Hash  `json:"topics" gencodec:"required"`
		Data             []byte         `json:"data" gencodec:"required"`
		BlockNumber      uint64         `json:"blockNumber"`
		TxHash           common.Hash    `json:"transactionHash" gencodec:"required"`
		TxIndex          uint           `json:"transactionIndex"`
		BlockHash        common.Hash    `json:"blockHash"`
		Index            uint           `json:"logIndex"`
		Removed          bool           `json:"removed"`
	}{
		FunctionSelector: er.FunctionSelector,
		Address:          er.Log.Address,
		Topics:           er.Log.Topics,
		Data:             er.Log.Data,
		BlockNumber:      er.Log.BlockNumber,
		TxHash:           er.Log.TxHash,
		TxIndex:          er.Log.TxIndex,
		BlockHash:        er.Log.BlockHash,
		Index:            er.Log.Index,
		Removed:          er.Log.Removed,
	})
}

type EnrichedReceipt struct {
	*types.Receipt
	Logs []*EnrichedReceiptLog `json:"logs"`
}

func (er *EnrichedReceipt) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type              uint8                 `json:"type,omitempty"`
		PostState         []byte                `json:"root"`
		Status            uint64                `json:"status"`
		CumulativeGasUsed uint64                `json:"cumulativeGasUsed" gencodec:"required"`
		Bloom             types.Bloom           `json:"logsBloom"         gencodec:"required"`
		Logs              []*EnrichedReceiptLog `json:"logs"`
		TxHash            common.Hash           `json:"transactionHash" gencodec:"required"`
		ContractAddress   common.Address        `json:"contractAddress"`
		GasUsed           uint64                `json:"gasUsed" gencodec:"required"`
		EffectiveGasPrice *big.Int              `json:"effectiveGasPrice"` // required, but tag omitted for backwards compatibility
		BlockHash         common.Hash           `json:"blockHash,omitempty"`
		BlockNumber       *big.Int              `json:"blockNumber,omitempty"`
		TransactionIndex  uint                  `json:"transactionIndex"`
	}{
		Type:              er.Receipt.Type,
		PostState:         er.Receipt.PostState,
		Status:            er.Receipt.Status,
		CumulativeGasUsed: er.Receipt.CumulativeGasUsed,
		Bloom:             er.Receipt.Bloom,
		Logs:              er.Logs,
		TxHash:            er.Receipt.TxHash,
		ContractAddress:   er.Receipt.ContractAddress,
		GasUsed:           er.Receipt.GasUsed,
		EffectiveGasPrice: er.Receipt.EffectiveGasPrice,
		BlockHash:         er.Receipt.BlockHash,
		BlockNumber:       er.Receipt.BlockNumber,
		TransactionIndex:  er.Receipt.TransactionIndex,
	})
}
