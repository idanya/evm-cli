package entities

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type EnrichedTxInfo struct {
	*types.Transaction
	DecodedData *DecodeResult
}

func (et *EnrichedTxInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		TxType      byte
		ChainID     *big.Int
		AccessList  types.AccessList
		Data        []byte
		Gas         uint64
		GasPrice    *big.Int
		GasTipCap   *big.Int
		GasFeeCap   *big.Int
		Value       *big.Int
		Nonce       uint64
		Hash        common.Hash
		To          *common.Address `json:"to"`
		DecodedData *DecodeResult   `json:"decodedData,omitempty"`
	}{
		TxType:      byte(et.Type()),
		ChainID:     et.Transaction.ChainId(),
		AccessList:  et.Transaction.AccessList(),
		Data:        et.Transaction.Data(),
		Gas:         et.Transaction.Gas(),
		GasPrice:    et.Transaction.GasPrice(),
		GasTipCap:   et.Transaction.GasTipCap(),
		GasFeeCap:   et.Transaction.GasFeeCap(),
		Value:       et.Transaction.Value(),
		Nonce:       et.Transaction.Nonce(),
		Hash:        et.Transaction.Hash(),
		To:          et.Transaction.To(),
		DecodedData: et.DecodedData,
	})
}
