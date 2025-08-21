package services

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/golang-jwt/jwt/v5"
	"github.com/skip2/go-qrcode"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"go-blockchain-receipt/internal/models"
)

type ReceiptService struct {
	db          *mongo.Collection
	ethClient   *ethclient.Client
	contract    *common.Address
	privateKey  *ecdsa.PrivateKey
	verifyURL   string
}

func NewReceiptService(
	db *mongo.Collection,
	ethClient *ethclient.Client,
	contract common.Address,
	privateKey *ecdsa.PrivateKey,
	verifyURL string,
) *ReceiptService {
	return &ReceiptService{
		db:          db,
		ethClient:   ethClient,
		contract:    &contract,
		privateKey:  privateKey,
		verifyURL:   verifyURL,
	}
}

func (s *ReceiptService) CreateReceipt(ctx context.Context, payload interface{}) (*models.CreateReceiptResponse, error) {
	// Canonicalize payload using JCS
	canonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Calculate SHA-256 hash
	hash := sha256.Sum256(canonBytes)
	hashHex := hex.EncodeToString(hash[:])

	// Create JWS using ES256
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"payload": string(canonBytes),
		"hash":    hashHex,
		"iat":     time.Now().Unix(),
	})
	token.Header["kid"] = "demo-key-1" // Demo KID

	jws, err := token.SignedString(s.privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign JWS: %w", err)
	}

	// Generate QR code
	verifyURL := fmt.Sprintf("%s?jws=%s", s.verifyURL, jws)
	qr, err := qrcode.Encode(verifyURL, qrcode.Medium, 256)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	receipt := &models.Receipt{
		Payload:   string(canonBytes),
		JWS:      jws,
		Hash:     hashHex,
		KID:      "demo-key-1",
		Status:   "PENDING",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save to MongoDB
	_, err = s.db.InsertOne(ctx, receipt)
	if err != nil {
		return nil, fmt.Errorf("failed to save receipt: %w", err)
	}

	// Anchor to blockchain (simplified)
	// TODO: Implement actual contract call
	receipt.Status = "ANCHORED"
	receipt.AnchorTx = "0x..." // Mock transaction hash

	// Update receipt status
	_, err = s.db.UpdateOne(ctx, 
		bson.M{"_id": receipt.ID},
		bson.M{"$set": bson.M{
			"status": receipt.Status,
			"anchorTx": receipt.AnchorTx,
			"updatedAt": time.Now(),
		}},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update receipt status: %w", err)
	}

	return &models.CreateReceiptResponse{
		JWS:       jws,
		Hash:      hashHex,
		KID:       "demo-key-1",
		VerifyURL: verifyURL,
		QRPng:     base64.StdEncoding.EncodeToString(qr),
		AnchorTx:  receipt.AnchorTx,
	}, nil
}

func (s *ReceiptService) VerifyReceipt(ctx context.Context, rid, jws string) (*models.VerifyResponse, error) {
	// Parse and verify JWS
	token, err := jwt.Parse(jws, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return &s.privateKey.PublicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid signature: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	// Get receipt from MongoDB
	var receipt models.Receipt
	err = s.db.FindOne(ctx, bson.M{"_id": rid}).Decode(&receipt)
	if err != nil {
		return nil, fmt.Errorf("receipt not found: %w", err)
	}

	// Verify hash matches
	if claims["hash"] != receipt.Hash {
		return nil, fmt.Errorf("hash mismatch")
	}

	return &models.VerifyResponse{
		OK:     true,
		Status: receipt.Status,
		KID:    receipt.KID,
		TS:     receipt.CreatedAt,
	}, nil
}

func (s *ReceiptService) GetJWKS() map[string]interface{} {
	// Demo JWKS - in production, this would be properly generated from the key
	return map[string]interface{}{
		"keys": []map[string]interface{}{
			{
				"kty": "EC",
				"kid": "demo-key-1",
				"use": "sig",
				"alg": "ES256",
				"crv": "P-256",
				"x":   "demo-x",
				"y":   "demo-y",
			},
		},
	}
}