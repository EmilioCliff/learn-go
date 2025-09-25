package block

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Transaction struct {
	senderBlockchainAddress    string
	recipientBlockchainAddress string
	timestamp                  int64
	value                      float32
}

func NewTransaction(sender string, recipient string, value float32) *Transaction {
	return &Transaction{sender, recipient, time.Now().Unix(), value}
}

type TransactionRequest struct {
	SenderBlockchainAddress    *string `json:"sender_blockchain_address"`
	RecipientBlockchainAddress *string `json:"recipient_blockchain_address"`
	SenderPublicKey            *string `json:"sender_public_key"`
	Value                      float32 `json:"value"`
	Signature                  *string `json:"signature"`
}

func (tr *TransactionRequest) Validate() bool {
	if tr.SenderBlockchainAddress == nil || tr.RecipientBlockchainAddress == nil || tr.SenderPublicKey == nil || tr.Value <= 0 || tr.Signature == nil {
		return false
	}
	return true
}

type AmountResponse struct {
	BlockchainAddress string  `json:"blockchain_address"`
	Amount            float32 `json:"amount"`
}

func (ar *AmountResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		BlockchainAddress string  `json:"blockchain_address"`
		Amount            float32 `json:"amount"`
	}{
		BlockchainAddress: ar.BlockchainAddress,
		Amount:            ar.Amount,
	})
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		SenderBlockchainAddress    string  `json:"sender_blockchain_address"`
		RecipientBlockchainAddress string  `json:"recipient_blockchain_address"`
		Value                      float32 `json:"value"`
		Timestamp                  int64   `json:"timestamp"`
	}{
		SenderBlockchainAddress:    t.senderBlockchainAddress,
		RecipientBlockchainAddress: t.recipientBlockchainAddress,
		Value:                      t.value,
		Timestamp:                  t.timestamp,
	})
}

func (t *Transaction) UnmarshalJSON(data []byte) error {
	v := &struct {
		SenderBlockchainAddress    *string  `json:"sender_blockchain_address"`
		RecipientBlockchainAddress *string  `json:"recipient_blockchain_address"`
		Value                      *float32 `json:"value"`
		Timestamp                  *int64   `json:"timestamp"`
	}{
		SenderBlockchainAddress:    &t.senderBlockchainAddress,
		RecipientBlockchainAddress: &t.recipientBlockchainAddress,
		Value:                      &t.value,
		Timestamp:                  &t.timestamp,
	}
	return json.Unmarshal(data, &v)
}

func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 60))
	fmt.Printf("%-30s %s\n", "sender_blockchain_address:", t.senderBlockchainAddress)
	fmt.Printf("%-30s %s\n", "recipient_blockchain_address:", t.recipientBlockchainAddress)
	fmt.Printf("%-30s %.1f\n", "value:", t.value)
}
