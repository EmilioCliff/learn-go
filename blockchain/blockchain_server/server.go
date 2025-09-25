package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/EmilioCliff/learn-go/blockchain/block"
	"github.com/EmilioCliff/learn-go/blockchain/wallet"
	"github.com/gin-gonic/gin"
)

var cache map[string]*block.Blockchain = make(map[string]*block.Blockchain)

type BlockchainServer struct {
	port    uint16
	gateway string
	router  *gin.Engine
	ln      net.Listener
	srv     *http.Server
}

func NewBlockchainServer(port uint16) *BlockchainServer {
	r := gin.Default()

	bcs := &BlockchainServer{port: port, router: r}
	bcs.setUpRoutes()
	return bcs
}

func (bcs *BlockchainServer) setUpRoutes() {
	bcs.router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})
	bcs.router.GET("/wallet", bcs.createWallet)
	bcs.router.GET("/wallet/:blockchain_address", bcs.getWallet)
	bcs.router.GET("/address/:blockchain_address/amount", bcs.getWalletAmount)
	bcs.router.GET("/chain", bcs.getChain)
	// bcs.router.DELETE("/wallet/:blockchain_address", bcs.deleteWallet)
	bcs.router.GET("/transactions", bcs.listTransactionPool)
	bcs.router.POST("/transactions", bcs.createTransaction)

	// internal
	bcs.router.PUT("/transactions", bcs.addTransaction)
	bcs.router.DELETE("/transactions", bcs.clearTransaction)
	bcs.router.GET("/mine", bcs.mine)
	bcs.router.GET("/mine/start", bcs.startMining)
	bcs.router.GET("/mine/stop", bcs.stopMining)
	bcs.router.PUT("/consensus", bcs.consensusResolve)

	bcs.srv = &http.Server{
		Addr:         bcs.PortAddress(),
		Handler:      bcs.router.Handler(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func (bcs *BlockchainServer) Start() error {
	bcs.GetBlockchain().Run()
	var err error
	if bcs.ln, err = net.Listen("tcp", bcs.PortAddress()); err != nil {
		return err
	}

	go func(s *BlockchainServer) {
		err := s.srv.Serve(s.ln)
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}(bcs)

	return nil
}

func (bcs *BlockchainServer) Stop() error {
	log.Println("Shutting down http server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return bcs.srv.Shutdown(ctx)
}

// initialize blockchain if not already done
func (bcs *BlockchainServer) GetBlockchain() *block.Blockchain {
	bc, ok := cache["blockchain"]
	if !ok {
		minersWallet := wallet.NewWallet()
		bc = block.NewBlockchain(minersWallet.BlockchainAddress(), bcs.Port())
		cache["blockchain"] = bc
		bc.RegisterWallet(minersWallet)
		log.Printf("blockchain_address %v", minersWallet.BlockchainAddress())
	}
	return bc
}

func (bcs *BlockchainServer) Port() uint16 {
	return bcs.port
}

func (bcs *BlockchainServer) PortAddress() string {
	return fmt.Sprintf("0.0.0.0:%d", bcs.port)
}
