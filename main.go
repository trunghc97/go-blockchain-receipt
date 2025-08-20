package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go-blockchain-receipt/api"
	_ "go-blockchain-receipt/api/docs" // swagger docs
	"go-blockchain-receipt/internal/services"
)

func main() {
	// Connect to MongoDB
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(getEnv("MONGODB_URI", "mongodb://localhost:27017")))
	if err != nil {
		log.Fatal(err)
	}
	defer mongoClient.Disconnect(context.Background())

	db := mongoClient.Database(getEnv("MONGODB_DB", "receipts"))
	receiptsCollection := db.Collection("receipts")

	// Connect to Ethereum node
	ethClient, err := ethclient.Dial(getEnv("ETH_RPC_URL", "http://localhost:8545"))
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
		common.HexToAddress(getEnv("CONTRACT_ADDRESS", "0x0")),
		privateKey,
		getEnv("VERIFY_URL", "http://localhost:8080/verify"),
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
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Start server
	e.Logger.Fatal(e.Start(getEnv("LISTEN_ADDR", ":8080")))
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
