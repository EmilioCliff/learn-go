package block

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

type Block struct {
	nonce        int
	previousHash [32]byte
	timestamp    int64
	transactions []*Transaction
}

func NewBlock(nonce int, previousHash [32]byte, transactions []*Transaction) *Block {
	b := new(Block)
	b.timestamp = time.Now().UnixNano()
	b.nonce = nonce
	b.previousHash = previousHash
	b.transactions = transactions
	return b
}

func (b *Block) Hash() [32]byte {
	m, _ := b.MarshalJSON()
	return sha256.Sum256(m)
}

func (b *Block) PreviousHash() [32]byte {
	return b.previousHash
}

func (b *Block) Nonce() int {
	return b.nonce
}

func (b *Block) Transactions() []*Transaction {
	return b.transactions
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Nonce        int            `json:"nonce"`
		PreviousHash string         `json:"previous_hash"`
		Timestamp    int64          `json:"timestamp"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Nonce:        b.nonce,
		PreviousHash: fmt.Sprintf("%x", b.previousHash),
		Timestamp:    b.timestamp,
		Transactions: b.transactions,
	})
}

func (b *Block) UnmarshalJSON(data []byte) error {
	var previousHash string
	v := &struct {
		Nonce        *int            `json:"nonce"`
		PreviousHash *string         `json:"previous_hash"`
		Timestamp    *int64          `json:"timestamp"`
		Transactions *[]*Transaction `json:"transactions"`
	}{
		Nonce:        &b.nonce,
		PreviousHash: &previousHash,
		Timestamp:    &b.timestamp,
		Transactions: &b.transactions,
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	ph, _ := hex.DecodeString(previousHash)
	copy(b.previousHash[:], ph[:32])
	return nil
}

func (b *Block) Print() {
	fmt.Printf("Timestamp: %d\n", b.timestamp)
	fmt.Printf("Nonce: %d\n", b.nonce)
	fmt.Printf("Previous Hash: %x\n", b.previousHash)
	for _, t := range b.transactions {
		t.Print()
	}
}
