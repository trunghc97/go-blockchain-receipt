# Blockchain Receipt Service

Dịch vụ REST API để tạo và xác thực biên lai giao dịch trên blockchain sử dụng Go, MongoDB và Ethereum.

## Tính năng

- Tạo biên lai với chữ ký JWS ES256
- Lưu trữ biên lai trong MongoDB
- Anchor hash biên lai lên blockchain EVM
- Xác thực biên lai và trạng thái on-chain
- Cung cấp JWKS cho xác thực chữ ký
- Kiểm tra sức khỏe hệ thống
- Docker hóa toàn bộ stack

## Cài đặt

### Yêu cầu

- Docker
- Docker Compose

### Chạy dịch vụ

1. Clone repository:
```bash
git clone https://github.com/hct97/go-blockchain-receipt.git
cd go-blockchain-receipt
```

2. Khởi động các services:
```bash
docker-compose up -d
```

Các services sẽ chạy tại:
- API: http://localhost:8080
- MongoDB Express: http://localhost:8081
- Anvil (EVM): http://localhost:8545

## API Endpoints

### POST /receipts
Tạo biên lai mới từ payload giao dịch.

Request:
```json
{
  "amount": 1000,
  "currency": "USD",
  "description": "Payment for services"
}
```

Response:
```json
{
  "jws": "eyJhbGciOiJFUzI1...",
  "hash": "1234abcd...",
  "kid": "demo-key-1",
  "verifyUrl": "http://localhost:8080/verify?jws=...",
  "qrPng": "base64...",
  "anchorTx": "0x..."
}
```

### GET /verify?rid=...&jws=...
Xác thực biên lai và trạng thái on-chain.

Response:
```json
{
  "ok": true,
  "status": "ANCHORED",
  "kid": "demo-key-1",
  "ts": "2024-03-20T10:00:00Z"
}
```

### GET /jwks.json
Lấy public JWKS cho xác thực chữ ký.

### GET /healthz
Kiểm tra sức khỏe hệ thống.

## Phát triển

### Cấu trúc project

```
.
├── api/            # API handlers
├── config/         # Configuration
├── contracts/      # Smart contracts
├── internal/       # Internal packages
│   ├── models/     # Data models
│   └── services/   # Business logic
└── pkg/           # Public packages
```

### Build và test

```bash
# Build
go build

# Test
go test ./...
```

## License

MIT
