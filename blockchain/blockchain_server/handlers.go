package main

import (
	"github.com/EmilioCliff/learn-go/blockchain/block"
	"github.com/EmilioCliff/learn-go/blockchain/utils"
	"github.com/EmilioCliff/learn-go/blockchain/wallet"
	"github.com/gin-gonic/gin"
)

func (bcs *BlockchainServer) createWallet(c *gin.Context) {
	myWallet := wallet.NewWallet()
	bc := bcs.GetBlockchain()
	bc.RegisterWallet(myWallet)
	m, _ := myWallet.MarshalJSON()
	c.Data(200, "application/json", m)
}

func (bcs *BlockchainServer) getWallet(c *gin.Context) {
	blockchainAddress := c.Param("blockchain_address")
	bc := bcs.GetBlockchain()
	myWallet := bc.GetWallet(blockchainAddress)
	if myWallet == nil {
		c.JSON(400, gin.H{"message": "failed", "error": "wallet not found"})
		return
	}
	m, _ := myWallet.MarshalJSON()
	c.Data(200, "application/json", m)
}

func (bcs *BlockchainServer) getChain(c *gin.Context) {
	bc := bcs.GetBlockchain()
	m, _ := bc.MarshalJSON()
	c.Data(200, "application/json", m)
}

func (bcs *BlockchainServer) getWalletAmount(c *gin.Context) {
	blockchainAddress := c.Param("blockchain_address")
	bc := bcs.GetBlockchain()
	amount := bc.CalculateTotalAmount(blockchainAddress)
	ar := &block.AmountResponse{BlockchainAddress: blockchainAddress, Amount: amount}
	m, _ := ar.MarshalJSON()
	c.Data(200, "application/json", m)
}

func (bcs *BlockchainServer) listTransactionPool(c *gin.Context) {
	bc := bcs.GetBlockchain()
	c.JSON(200, bc.TransactionsPool())
}

type transactionRequest struct {
	SenderPrivateKey           string  `json:"sender_private_key"`
	SenderPublicKey            string  `json:"sender_public_key"`
	SenderBlockchainAddress    string  `json:"sender_blockchain_address"`
	RecipientBlockchainAddress string  `json:"recipient_blockchain_address"`
	Value                      float32 `json:"value"`
}

func (tr *transactionRequest) Validate() bool {
	if tr.SenderPrivateKey == "" || tr.SenderPublicKey == "" || tr.SenderBlockchainAddress == "" || tr.RecipientBlockchainAddress == "" || tr.Value <= 0 {
		return false
	}
	return true
}

// CreateTransactionHandler handles POST /transactions from clients
func (bcs *BlockchainServer) createTransaction(c *gin.Context) {
	var tr transactionRequest
	if err := c.ShouldBindJSON(&tr); err != nil {
		c.JSON(400, gin.H{"message": "failed", "error": err.Error()})
		return
	}

	if !tr.Validate() {
		c.JSON(400, gin.H{"message": "failed", "error": "missing or invalid fields"})
		return
	}

	publicKey := utils.PublicKeyFromString(tr.SenderPublicKey)
	privateKey := utils.PrivateKeyFromString(tr.SenderPrivateKey, publicKey)

	transaction := wallet.NewTransaction(
		privateKey,
		publicKey,
		tr.SenderBlockchainAddress,
		tr.RecipientBlockchainAddress,
		tr.Value,
	)

	signature := transaction.GenerateSignature()
	signatureStr := signature.String()

	bt := &block.TransactionRequest{
		SenderBlockchainAddress:    &tr.SenderBlockchainAddress,
		RecipientBlockchainAddress: &tr.RecipientBlockchainAddress,
		SenderPublicKey:            &tr.SenderPublicKey,
		Value:                      tr.Value,
		Signature:                  &signatureStr,
	}

	bc := bcs.GetBlockchain()
	isCreated := bc.CreateTransaction(
		*bt.SenderBlockchainAddress,
		*bt.RecipientBlockchainAddress,
		bt.Value,
		publicKey,
		signature,
	)

	if !isCreated {
		c.JSON(400, gin.H{"message": "failed", "error": "failed to create a transaction"})
		return
	}
	c.JSON(200, gin.H{"message": "success"})
}

// AddTransactionHandler handles PUT /transactions from other nodes
func (bcs *BlockchainServer) addTransaction(c *gin.Context) {
	var tr block.TransactionRequest
	if err := c.ShouldBindJSON(&tr); err != nil {
		c.JSON(400, gin.H{"message": "failed", "error": err.Error()})
		return
	}

	if !tr.Validate() {
		c.JSON(400, gin.H{"message": "failed", "error": "missing or invalid fields"})
		return
	}

	publicKey := utils.PublicKeyFromString(*tr.SenderPublicKey)
	signature := utils.SignatureFromString(*tr.Signature)
	bc := bcs.GetBlockchain()
	isCreated := bc.AddTransaction(
		*tr.SenderBlockchainAddress,
		*tr.RecipientBlockchainAddress,
		tr.Value,
		publicKey,
		signature,
	)
	if !isCreated {
		c.JSON(400, gin.H{"message": "failed", "error": "failed to add a transaction"})
		return
	}
	c.JSON(200, gin.H{"message": "success"})
}

func (bcs *BlockchainServer) clearTransaction(c *gin.Context) {
	bc := bcs.GetBlockchain()
	bc.ClearTransactionPool()
	c.JSON(200, gin.H{"message": "success"})
}

func (bcs *BlockchainServer) mine(c *gin.Context) {
	bc := bcs.GetBlockchain()
	isMined := bc.Mining()
	if !isMined {
		c.JSON(400, gin.H{"message": "failed", "error": "mining is already in progress"})
		return
	}
	c.JSON(200, gin.H{"message": "success"})
}

func (bcs *BlockchainServer) startMining(c *gin.Context) {
	bc := bcs.GetBlockchain()
	bc.StartMining()
	c.JSON(200, gin.H{"message": "success"})
}

func (bcs *BlockchainServer) stopMining(c *gin.Context) {
	bc := bcs.GetBlockchain()
	bc.StopMining()
	c.JSON(200, gin.H{"message": "success"})
}

func (bcs *BlockchainServer) consensusResolve(c *gin.Context) {
	bc := bcs.GetBlockchain()
	isReplaced := bc.ResolveConflicts()
	if isReplaced {
		c.JSON(200, gin.H{"message": "success", "replaced": true})
	}
	c.JSON(200, gin.H{"message": "success", "replaced": false})
}
