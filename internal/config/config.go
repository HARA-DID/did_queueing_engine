package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Redis      RedisConfig
	DB         DBConfig
	Blockchain BlockchainConfig
	Worker     WorkerConfig
	Server     ServerConfig
}

type RedisConfig struct {
	URL        string
	StreamName string
	GroupName  string
	DLQSuffix  string
}

type DBConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type BlockchainConfig struct {
	RPCURLs             []string
	PrivateKeys         []string
	EntryPointHNS       string
	GasManagerHNS       string
	WalletHNS           string
	WalletFactoryHNS    string
	DIDRootFactoryHNS   string
	DIDRootStorageHNS   string
	DIDOrgStorageHNS    string
	VCFactoryHNS        string
	VCStorageHNS        string
	VCCertificateNFTHNS string
	VCIdentityNFTHNS    string
	AliasFactoryHNS     string
	AliasStorageHNS     string
	TxCheckChannelBuffer int
}

type WorkerConfig struct {
	ConsumerName    string
	Concurrency     int
	PollInterval    time.Duration
	MaxRetry        int
	RetryBaseDelay  time.Duration
	ShutdownTimeout time.Duration
	BatchSize       int64
}

type ServerConfig struct {
	Port string
}

func Load() (*Config, error) {
	cfg := &Config{
		Redis: RedisConfig{
			URL:        requireEnv("REDIS_URL"),
			StreamName: requireEnv("STREAM_NAME"),
			GroupName:  requireEnv("GROUP_NAME"),
			DLQSuffix:  getEnvOrDefault("DLQ_SUFFIX", ":dlq"),
		},
		DB: DBConfig{
			URL:             requireEnv("DB_URL"),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		Blockchain: BlockchainConfig{
			RPCURLs:             splitEnv("RPC_URLS", ","),
			PrivateKeys:         splitEnv("PRIVATE_KEYS", ","),
			EntryPointHNS:       getEnvOrDefault("AA_ENTRYPOINT_HNS", ""),
			GasManagerHNS:       getEnvOrDefault("AA_GAS_MANAGER_HNS", ""),
			WalletHNS:           getEnvOrDefault("AA_WALLET_HNS", ""),
			WalletFactoryHNS:    getEnvOrDefault("AA_WALLET_FACTORY_HNS", ""),
			DIDRootFactoryHNS:   getEnvOrDefault("DID_ROOT_FACTORY_HNS", ""),
			DIDRootStorageHNS:   getEnvOrDefault("DID_ROOT_STORAGE_HNS", ""),
			DIDOrgStorageHNS:    getEnvOrDefault("DID_ORG_STORAGE_HNS", ""),
			VCFactoryHNS:        getEnvOrDefault("DID_VC_FACTORY_HNS", ""),
			VCStorageHNS:        getEnvOrDefault("DID_VC_STORAGE_HNS", ""),
			VCCertificateNFTHNS: getEnvOrDefault("DID_VC_CERTIFICATE_NFT_HNS", ""),
			VCIdentityNFTHNS:    getEnvOrDefault("DID_VC_IDENTITY_NFT_HNS", ""),
			AliasFactoryHNS:     getEnvOrDefault("DID_ALIAS_FACTORY_HNS", ""),
			AliasStorageHNS:     getEnvOrDefault("DID_ALIAS_STORAGE_HNS", ""),
			TxCheckChannelBuffer: getEnvInt("TX_CHECK_CHANNEL_BUFFER", 100),
		},
		Worker: WorkerConfig{
			ConsumerName:    getEnvOrDefault("CONSUMER_NAME", defaultConsumerName()),
			BatchSize:       int64(getEnvInt("BATCH_SIZE", 10)),
			PollInterval:    getEnvDuration("POLL_INTERVAL", 100*time.Millisecond),
			MaxRetry:        getEnvInt("MAX_RETRY", 3),
			RetryBaseDelay:  getEnvDuration("RETRY_BASE_DELAY", 1*time.Second),
			ShutdownTimeout: getEnvDuration("SHUTDOWN_TIMEOUT", 30*time.Second),
		},
		Server: ServerConfig{
			Port: getEnvOrDefault("SERVER_PORT", "8080"),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}
	return cfg, nil
}

func (c *Config) validate() error {
	if c.Redis.URL == "" {
		return fmt.Errorf("REDIS_URL is required")
	}
	if c.DB.URL == "" {
		return fmt.Errorf("DB_URL is required")
	}
	if len(c.Blockchain.RPCURLs) == 0 {
		return fmt.Errorf("RPC_URLS is required")
	}
	if len(c.Blockchain.PrivateKeys) == 0 {
		return fmt.Errorf("PRIVATE_KEYS is required and cannot be empty")
	}

	if c.Worker.MaxRetry < 0 {
		return fmt.Errorf("MAX_RETRY must be >= 0")
	}
	return nil
}

// DLQStreamName returns the dead-letter queue stream name.
func (c *RedisConfig) DLQStreamName() string {
	return c.StreamName + c.DLQSuffix
}

func requireEnv(key string) string {
	v := os.Getenv(key)
	return v
}

func getEnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}

func getEnvDuration(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}

func splitEnv(key, sep string) []string {
	v := os.Getenv(key)
	if v == "" {
		return nil
	}
	parts := strings.Split(v, sep)
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			result = append(result, t)
		}
	}
	return result
}

func defaultConsumerName() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "worker-default"
	}
	return fmt.Sprintf("worker-%s", hostname)
}
