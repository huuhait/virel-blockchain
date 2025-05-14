package walletrpc

import (
	"virel-blockchain/address"
	"virel-blockchain/util"
	"virel-blockchain/util/enc"
)

type TxInfo struct {
	Hash util.Hash `json:"hash"`
	Data *TxData   `json:"data,omitempty"`
}

type TxData struct {
	Sender    address.Integrated `json:"sender"`    // Sender
	Recipient address.Integrated `json:"recipient"` // Recipient
	Amount    uint64             `json:"amount"`
	Fee       uint64             `json:"fee"`
	Nonce     uint64             `json:"nonce"`
	Signature enc.Hex            `json:"signature"`
}

////////

type GetBalanceRequest struct {
}
type GetBalanceResponse struct {
	Balance        uint64 `json:"balance"`
	MempoolBalance uint64 `json:"mempool_balance"`
}

type GetHistoryRequest struct {
	Subaddress      string `json:"subaddress,omitempty"`
	IncludeTxData   bool   `json:"include_tx_data"`
	IncludeIncoming bool   `json:"include_incoming"`
	IncludeOutgoing bool   `json:"include_outgoing"`
}
type GetHistoryResponse struct {
	Incoming []TxInfo `json:"incoming,omitempty"`
	Outgoing []TxInfo `json:"outgoing,omitempty"`
}

type CreateTransactionRequest struct {
	Destination address.Integrated `json:"destination"`
	Amount      uint64             `json:"amount"`
}
type CreateTransactionResponse struct {
	TxBlob enc.Hex   `json:"tx_blob"`
	TXID   util.Hash `json:"txid"`
}

type SubmitTransactionRequest struct {
	TxBlob enc.Hex `json:"tx_blob"`
}
type SubmitTransactionResponse struct {
	TXID util.Hash `json:"txid"`
}
