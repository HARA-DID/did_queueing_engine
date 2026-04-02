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
	RPCURLs         []string
	PrivateKey      string
	HNSName         string
	ContractAddress string
	ContractABI     string
	EntryPointAddress string
	EntryPointHNS     string
	FactoryAddress    string
	GasManagerAddress string
	VCFactoryAddress  string
	VCFactoryHNS      string
	VCStorageAddress  string
	VCStorageHNS      string
	AliasFactoryAddress string
	AliasFactoryHNS     string
	AliasStorageAddress string
	AliasStorageHNS     string
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
			RPCURLs:         splitEnv("RPC_URLS", ","),
			PrivateKey:      requireEnv("PRIVATE_KEY"),
			HNSName:         getEnvOrDefault("HNS_NAME", ""),
			ContractABI:     getEnvOrDefault("CONTRACT_ABI", ""),
			EntryPointAddress: getEnvOrDefault("ENTRYPOINT_ADDRESS", ""),
			EntryPointHNS:     getEnvOrDefault("ENTRYPOINT_HNS", ""),
			FactoryAddress:    getEnvOrDefault("FACTORY_ADDRESS", ""),
			GasManagerAddress: getEnvOrDefault("GASMANAGER_ADDRESS", ""),
			VCFactoryAddress:  getEnvOrDefault("VC_FACTORY_ADDRESS", ""),
			VCFactoryHNS:      getEnvOrDefault("VC_FACTORY_HNS", ""),
			VCStorageAddress:  getEnvOrDefault("VC_STORAGE_ADDRESS", ""),
			VCStorageHNS:      getEnvOrDefault("VC_STORAGE_HNS", ""),
			AliasFactoryAddress: getEnvOrDefault("ALIAS_FACTORY_ADDRESS", ""),
			AliasFactoryHNS:     getEnvOrDefault("ALIAS_FACTORY_HNS", ""),
			AliasStorageAddress: getEnvOrDefault("ALIAS_STORAGE_ADDRESS", ""),
			AliasStorageHNS:     getEnvOrDefault("ALIAS_STORAGE_HNS", ""),
		},
		Worker: WorkerConfig{
			ConsumerName:    getEnvOrDefault("CONSUMER_NAME", defaultConsumerName()),
			Concurrency:     getEnvInt("WORKER_CONCURRENCY", 10),
			PollInterval:    getEnvDuration("POLL_INTERVAL", 100*time.Millisecond),
			MaxRetry:        getEnvInt("MAX_RETRY", 3),
			RetryBaseDelay:  getEnvDuration("RETRY_BASE_DELAY", 1*time.Second),
			ShutdownTimeout: getEnvDuration("SHUTDOWN_TIMEOUT", 30*time.Second),
			BatchSize:       int64(getEnvInt("BATCH_SIZE", 10)),
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
	if c.Blockchain.PrivateKey == "" {
		return fmt.Errorf("PRIVATE_KEY is required")
	}
	if c.Worker.Concurrency < 1 {
		return fmt.Errorf("WORKER_CONCURRENCY must be >= 1")
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
