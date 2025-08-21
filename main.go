package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go-blockchain-receipt/api"
	"go-blockchain-receipt/config"
	"go-blockchain-receipt/internal/services"
)

func main() {
	cfg := config.LoadConfig()

	// Connect to MongoDB
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatal(err)
	}
	defer mongoClient.Disconnect(context.Background())

	db := mongoClient.Database(cfg.MongoDB)
	receiptsCollection := db.Collection("receipts")

	// Connect to Ethereum node
	ethClient, err := ethclient.Dial(cfg.EthRPCURL)
	if err != nil {
		log.Fatal(err)
	}

	// Generate demo private key (in production, this would be loaded from secure storage)
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize services
	receiptService := services.NewReceiptService(
		receiptsCollection,
		ethClient,
		common.HexToAddress(cfg.ContractAddress),
		privateKey,
		cfg.VerifyURL,
	)

	// Initialize Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize handlers
	h := api.NewHandler(receiptService)

	// Routes
	e.POST("/receipts", h.CreateReceipt)
	e.GET("/verify", h.VerifyReceipt)
	e.GET("/jwks.json", h.GetJWKS)
	e.GET("/healthz", h.HealthCheck)

	// Swagger UI
	e.GET("/swagger/*", echo.WrapHandler(http.StripPrefix("/swagger/", http.FileServer(http.Dir("api/docs")))))

	// Start server
	e.Logger.Fatal(e.Start(cfg.ListenAddr))
}