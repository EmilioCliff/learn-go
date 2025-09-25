package block

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/EmilioCliff/learn-go/blockchain/utils"
	"github.com/EmilioCliff/learn-go/blockchain/wallet"
)

type Blockchain struct {
	chain             []*Block
	transactionPool   []*Transaction
	blockchainAddress string
	port              uint16
	mux               sync.Mutex
	cancelMining      *time.Timer

	config utils.Config

	neighbors   []string
	neighborMux sync.Mutex

	wallets map[string]*wallet.Wallet
}

func NewBlockchain(blockchainAddress string, port uint16) *Blockchain {
	bc := &Blockchain{
		blockchainAddress: blockchainAddress,
		port:              port,
		wallets:           make(map[string]*wallet.Wallet),
		transactionPool:   []*Transaction{},
		chain:             []*Block{},
		cancelMining:      nil,
	}
	bc.config, _ = utils.LoanConfig()
	bc.neighbors = bc.config.NEIGHBORS

	b := &Block{}
	bc.CreateBlock(0, b.previousHash)

	return bc
}

func (bc *Blockchain) Run() {
	bc.ResolveConflicts()
	// bc.StartMining() // comment auto mining out
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}

	for _, n := range bc.neighbors {
		if n == fmt.Sprintf("%s:%d", bc.config.HOST, bc.port) {
			continue
		}
		endpoint := fmt.Sprintf("%s/transactions", n)
		req, _ := http.NewRequest("DELETE", endpoint, nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("ERROR: Send  transaction to %s : %s\n", n, err.Error())
			continue
		}
		if resp.StatusCode == http.StatusOK {
			log.Printf("INFO: Send transaction to %s\n", n)
		} else {
			log.Printf("ERROR: Send transaction to %s\n", n)
		}
	}
	return b
}

func (bc *Blockchain) CreateTransaction(sender string, recipient string, value float32, senderPublicKey *ecdsa.PublicKey, s *utils.Signature) bool {
	isTransacted := bc.AddTransaction(sender, recipient, value, senderPublicKey, s)
	if isTransacted {
		for _, n := range bc.neighbors {
			if n == fmt.Sprintf("%s:%d", bc.config.HOST, bc.port) {
				continue
			}
			publickKeyStr := fmt.Sprintf("%064x%064x", senderPublicKey.X.Bytes(), senderPublicKey.Y.Bytes())
			signatureStr := s.String()
			bt := &TransactionRequest{
				SenderBlockchainAddress:    &sender,
				RecipientBlockchainAddress: &recipient,
				Value:                      value,
				SenderPublicKey:            &publickKeyStr,
				Signature:                  &signatureStr,
			}
			m, _ := json.Marshal(bt)
			endpoint := fmt.Sprintf("%s/transactions", n)
			req, _ := http.NewRequest("PUT", endpoint, bytes.NewBuffer(m))
			client := &http.Client{}
			resp, _ := client.Do(req)
			if resp.StatusCode == http.StatusOK {
				log.Printf("INFO: Send transaction to %s\n", n)
			} else {
				log.Printf("ERROR: Send transaction to %s\n", n)
			}
		}
	}

	return isTransacted
}

func (bc *Blockchain) AddTransaction(sender string, recipient string, value float32, senderPublicKey *ecdsa.PublicKey, s *utils.Signature) bool {
	t := NewTransaction(sender, recipient, value)

	if sender == bc.config.MINING_SENDER {
		log.Println("INFO: Mining reward")
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	}

	if bc.VerifyTransactionSignature(senderPublicKey, s, t) {
		log.Println(bc.CalculateTotalAmount(sender))
		log.Println(value)
		if bc.CalculateTotalAmount(sender) < value {
			log.Println("ERROR: Not enough balance in a wallet")
			return false
		}
		log.Println("INFO: Verified transaction")
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	} else {
		log.Println("ERROR: Verify Transaction")
	}

	return false
}

func (bc *Blockchain) RegisterWallet(w *wallet.Wallet) {
	bc.wallets[w.BlockchainAddress()] = w
}

func (bc *Blockchain) GetWallet(blockchainAddress string) *wallet.Wallet {
	return bc.wallets[blockchainAddress]
}

func (bc *Blockchain) VerifyTransactionSignature(senderPublicKey *ecdsa.PublicKey, s *utils.Signature, t *Transaction) bool {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	return ecdsa.Verify(senderPublicKey, h[:], s.R, s.S)
}

func (bc *Blockchain) ProofOfWork() int {
	previousHash := bc.LastBlock().Hash()
	transaction := bc.CopyTransactionPool()
	nonce := 0
	for !bc.ValidProof(nonce, previousHash, transaction) {
		nonce++
	}
	return nonce
}

func (bc *Blockchain) ValidProof(nonce int, previousHash [32]byte, transactions []*Transaction) bool {
	b := Block{nonce, previousHash, 0, transactions}
	hash := fmt.Sprintf("%x", b.Hash())
	return hash[:bc.config.MINING_DIFFICULTY] == strings.Repeat("0", bc.config.MINING_DIFFICULTY)
}

func (bc *Blockchain) Mining() bool {
	bc.mux.Lock()
	defer bc.mux.Unlock()

	bc.AddTransaction(bc.config.MINING_SENDER, bc.blockchainAddress, bc.config.MINING_REWARD, nil, nil)
	nonce := bc.ProofOfWork()
	previousHash := bc.LastBlock().Hash()
	bc.CreateBlock(nonce, previousHash)

	for _, n := range bc.neighbors {
		if n == fmt.Sprintf("%s:%d", bc.config.HOST, bc.port) {
			continue
		}
		enpoint := fmt.Sprintf("%s/consensus", n)
		client := &http.Client{}
		req, _ := http.NewRequest("PUT", enpoint, nil)
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("ERROR: Send consensus to %s : %s\n", n, err.Error())
			continue
		}
		log.Printf("INFO: Send consensus to %s %d\n", n, resp.StatusCode)
	}
	return true
}

func (bc *Blockchain) ValidChain(chain []*Block) bool {
	preBlock := chain[0]
	currentIndex := 1
	for currentIndex < len(chain) {
		b := chain[currentIndex]
		if b.previousHash != preBlock.Hash() {
			return false
		}
		if !bc.ValidProof(b.Nonce(), b.PreviousHash(), b.Transactions()) {
			return false
		}
		preBlock = b
		currentIndex++
	}
	return true
}

func (bc *Blockchain) ResolveConflicts() bool {
	var longestChain []*Block = nil
	maxLength := len(bc.chain)

	for _, n := range bc.neighbors {
		if n == fmt.Sprintf("%s:%d", bc.config.HOST, bc.port) {
			continue
		}
		endpoint := fmt.Sprintf("%s/chain", n)
		resp, err := http.Get(endpoint)
		if err != nil {
			log.Printf("ERROR: Get chain from %s\n", n)
			continue
		}
		if resp.StatusCode == http.StatusOK {
			var bcResp Blockchain
			json.NewDecoder(resp.Body).Decode(&bcResp)

			chain := bcResp.Chain()
			if len(chain) > maxLength && bc.ValidChain(chain) {
				maxLength = len(chain)
				longestChain = chain
			}
		}
	}
	if longestChain != nil {
		bc.chain = longestChain
		log.Printf("INFO: Replace chain with the longest chain from neighbors\n")
		return true
	}
	log.Printf("INFO: No conflicts found\n")
	return false
}

func (bc *Blockchain) StartMining() {
	log.Println("INFO: Start mining...")
	bc.Mining()
	bc.cancelMining = time.AfterFunc(bc.config.MINING_TIMER, bc.StartMining)
}

func (bc *Blockchain) StopMining() {
	if bc.cancelMining != nil {
		bc.cancelMining.Stop()
		bc.cancelMining = nil
	}
	log.Println("INFO: Stop mining")
}

func (bc *Blockchain) ClearTransactionPool() {
	bc.transactionPool = bc.transactionPool[:0]
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, len(bc.transactionPool))
	for i, t := range bc.transactionPool {
		transactions[i] = NewTransaction(t.senderBlockchainAddress, t.recipientBlockchainAddress, t.value)
	}

	return transactions
}

func (bc *Blockchain) CalculateTotalAmount(blockchainAddress string) float32 {
	var totalAmount float32 = 0.0
	for _, b := range bc.chain {
		for _, t := range b.transactions {
			value := t.value
			if blockchainAddress == t.recipientBlockchainAddress {
				totalAmount += value
			}
			if blockchainAddress == t.senderBlockchainAddress {
				totalAmount -= value
			}
		}
	}

	for _, t := range bc.transactionPool {
		if blockchainAddress == t.senderBlockchainAddress {
			totalAmount -= t.value
		}
		if blockchainAddress == t.recipientBlockchainAddress {
			totalAmount += t.value
		}
	}
	return totalAmount
}

func (bc *Blockchain) BlockchainAddress() string {
	return bc.blockchainAddress
}

func (bc *Blockchain) Neighbors() []string {
	return bc.neighbors
}

func (bc *Blockchain) Chain() []*Block {
	return bc.chain
}

func (bc *Blockchain) TransactionsPool() []*Transaction {
	return bc.transactionPool
}

func (bc *Blockchain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks            []*Block                  `json:"chain"`
		ChainLenght       int                       `json:"chain_length"`
		TransactionPool   []*Transaction            `json:"transaction_pool"`
		BlockchainAddress string                    `json:"blockchain_address"`
		Port              uint16                    `json:"port"`
		Host              string                    `json:"host"`
		Mining            bool                      `json:"mining"`
		MiningDifficulty  int                       `json:"mining_difficulty"`
		MiningReward      float32                   `json:"mining_reward"`
		Neighbors         []string                  `json:"neighbors"`
		Wallets           map[string]*wallet.Wallet `json:"wallets"`
	}{
		Blocks:            bc.chain,
		ChainLenght:       len(bc.chain),
		TransactionPool:   bc.transactionPool,
		BlockchainAddress: bc.blockchainAddress,
		Host:              bc.config.HOST,
		Port:              bc.port,
		Mining:            bc.cancelMining != nil,
		MiningDifficulty:  bc.config.MINING_DIFFICULTY,
		MiningReward:      bc.config.MINING_REWARD,
		Neighbors:         bc.neighbors,
		Wallets:           bc.wallets,
	})
}

func (bc *Blockchain) UnmarshalJSON(data []byte) error {
	v := struct {
		Blocks *[]*Block `json:"chain"`
	}{
		Blocks: &bc.chain,
	}
	return json.Unmarshal(data, &v)
}

func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Block %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 60))
}
