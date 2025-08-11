package walletrpc

import (
	"github.com/virel-project/virel-blockchain/address"
	"github.com/virel-project/virel-blockchain/rpc/daemonrpc"
	"github.com/virel-project/virel-blockchain/util"
	"github.com/virel-project/virel-blockchain/util/enc"
)

type TxInfo struct {
	Hash util.Hash                         `json:"hash"`
	Data *daemonrpc.GetTransactionResponse `json:"data,omitempty"`
}

type TxData struct {
	Sender    *address.Integrated `json:"sender"`    // Sender
	Recipient address.Integrated  `json:"recipient"` // Recipient
	Amount    uint64              `json:"amount"`
	Fee       uint64              `json:"fee"`
	Nonce     uint64              `json:"nonce"`
	Signature enc.Hex             `json:"signature"`
}

////////

type GetBalanceRequest struct {
}
type GetBalanceResponse struct {
	Balance        uint64 `json:"balance"`
	MempoolBalance uint64 `json:"mempool_balance"`
}

type GetHistoryRequest struct {
	IncludeTxData             bool   `json:"include_tx_data"`
	IncludeIncoming           bool   `json:"include_incoming"`
	IncludeOutgoing           bool   `json:"include_outgoing"`
	FilterIncomingByPaymentId uint64 `json:"filter_incoming_by_payment_id"`
}
type GetHistoryResponse struct {
	Incoming []TxInfo `json:"incoming,omitempty"`
	Outgoing []TxInfo `json:"outgoing,omitempty"`
}

type Output struct {
	Amount    uint64             `json:"amount"`
	Recipient address.Integrated `json:"recipient"`
}

type CreateTransactionRequest struct {
	Outputs []Output `json:"outputs"`
}
type CreateTransactionResponse struct {
	TxBlob enc.Hex   `json:"tx_blob"`
	TXID   util.Hash `json:"txid"`
	Fee    uint64    `json:"fee"`
}

type SubmitTransactionRequest struct {
	TxBlob enc.Hex `json:"tx_blob"`
}
type SubmitTransactionResponse struct {
	TXID util.Hash `json:"txid"`
}

type RefreshRequest struct {
}
type RefreshResponse struct {
	Success bool `json:"success"`
}

type GetSubaddressRequest struct {
	PaymentId  uint64             `json:"payment_id"`
	Subaddress address.Integrated `json:"subaddress"`
}
type GetSubaddressResponse struct {
	PaymentId            uint64             `json:"payment_id"`
	Subaddress           address.Integrated `json:"subaddress"`
	TotalReceived        uint64             `json:"total_received"`
	MempoolTotalReceived uint64             `json:"mempool_total_received"`
}
