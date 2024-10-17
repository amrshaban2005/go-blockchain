package blockchain

import (
	"crypto/sha256"
	//"encoding/hex"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

type Block struct {
	timestamp    int64
	nonce        int
	previousHash [32]byte
	transactions []*Transaction
}

func NewBlock(nonce int, previousHash [32]byte, transactions []*Transaction) *Block {
	b := new(Block)
	b.nonce = nonce
	b.previousHash = previousHash
	b.timestamp = time.Now().UnixNano()
	b.transactions = transactions
	return b
}

func (b *Block) Nonce() int {
	return b.nonce
}

func (b *Block) PreviousHash() [32]byte {
	return b.previousHash
}

func (b *Block) Transaction() []*Transaction {
	return b.transactions
}

func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)

	return sha256.Sum256(m)
}

func (b *Block) Print() {
	fmt.Printf("timestamp\t%d\n", b.timestamp)
	fmt.Printf("nonce\t%d\n", b.nonce)
	fmt.Printf("previous_hash\t%x\n", b.previousHash)

	for _, t := range b.transactions {
		t.Print()
	}
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64          `json:"timestamp"`
		Nonce        int            `json:"nonce"`
		PreviousHash string         `json:"previous_hash"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Timestamp:    b.timestamp,
		Nonce:        b.nonce,
		PreviousHash: fmt.Sprintf("%x", b.previousHash),
		Transactions: b.transactions,
	})
}

func (b *Block) UnmarshalJSON(data []byte) error {
	var previousHash string
	v := &struct {
		Timestamp    *int64          `json:"timestamp"`
		Nonce        *int            `json:"nonce"`
		PreviousHash *string         `json:"previous_hash"`
		Transactions *[]*Transaction `json:"transactions"`
	}{
		Timestamp:    &b.timestamp,
		Nonce:        &b.nonce,
		PreviousHash: &previousHash,
		Transactions: &b.transactions,
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	ph, _ := hex.DecodeString(*v.PreviousHash)
	copy(b.previousHash[:], ph[:32])
	return nil
}
