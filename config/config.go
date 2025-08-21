package config

import "os"

type Config struct {
	MongoURI       string
	MongoDB        string
	EthRPCURL      string
	ContractAddress string
	VerifyURL      string
	ListenAddr     string
}

func LoadConfig() *Config {
	return &Config{
		MongoURI:       getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		MongoDB:        getEnv("MONGODB_DB", "receipts"),
		EthRPCURL:      getEnv("ETH_RPC_URL", "http://localhost:8545"),
		ContractAddress: getEnv("CONTRACT_ADDRESS", "0x0"),
		VerifyURL:      getEnv("VERIFY_URL", "http://localhost:8080/verify"),
		ListenAddr:     getEnv("LISTEN_ADDR", ":8080"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
