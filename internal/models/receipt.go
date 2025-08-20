package models

import (
	"time"
)

type Receipt struct {
	ID          string    `bson:"_id,omitempty"`
	Payload     string    `bson:"payload"`      // Canonicalized payload (JCS)
	JWS         string    `bson:"jws"`          // Compact JWS
	Hash        string    `bson:"hash"`         // SHA-256 hash
	KID         string    `bson:"kid"`          // Key ID used for signing
	AnchorTx    string    `bson:"anchorTx"`     // Transaction hash on EVM
	Status      string    `bson:"status"`       // PENDING, ANCHORED, SKIPPED
	SkipReason  string    `bson:"skipReason"`   // Reason if status is SKIPPED
	CreatedAt   time.Time `bson:"createdAt"`
	UpdatedAt   time.Time `bson:"updatedAt"`
}

type VerifyResponse struct {
	OK     bool      `json:"ok"`
	Status string    `json:"status"`
	KID    string    `json:"kid"`
	TS     time.Time `json:"ts"`
}

type CreateReceiptResponse struct {
	JWS       string `json:"jws"`
	Hash      string `json:"hash"`
	KID       string `json:"kid"`
	VerifyURL string `json:"verifyUrl"`
	QRPng     string `json:"qrPng"` // Base64 encoded PNG
	AnchorTx  string `json:"anchorTx"`
}
