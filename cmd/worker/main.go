package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/myorg/worker-service/internal/config"
	infradb "github.com/myorg/worker-service/internal/infra/db"
	redisinfra "github.com/myorg/worker-service/internal/infra/redis"
	"github.com/myorg/worker-service/internal/sdk"
	"github.com/myorg/worker-service/internal/service"
	"github.com/myorg/worker-service/internal/worker"
	"github.com/myorg/worker-service/pkg"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)

	_ = godotenv.Load()

	// ── Config ─────────────────────────────────────────────────────────────
	cfg, err := config.Load()
	if err != nil {
		log.WithError(err).Fatal("failed to load config")
	}

	// ── Database ───────────────────────────────────────────────────────────
	db, err := infradb.Connect(cfg.DB)
	if err != nil {
		log.WithError(err).Fatal("failed to connect to database")
	}
	defer db.Close()

	// ── Redis ──────────────────────────────────────────────────────────────
	redisClient, err := redisinfra.NewClient(cfg.Redis)
	if err != nil {
		log.WithError(err).Fatal("failed to connect to redis")
	}
	defer redisClient.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := redisinfra.EnsureConsumerGroup(ctx, redisClient, cfg.Redis.StreamName, cfg.Redis.GroupName); err != nil {
		log.WithError(err).Fatal("failed to ensure consumer group")
	}
	log.WithFields(logrus.Fields{
		"stream": cfg.Redis.StreamName,
		"group":  cfg.Redis.GroupName,
	}).Info("consumer group ready")

	// ── Blockchain Modular SDK ─────────────────────────────────────────────
	provider, err := sdk.NewProvider(cfg.Blockchain)
	if err != nil {
		log.WithError(err).Fatal("failed to initialise blockchain provider")
	}

	didAdapter, err := sdk.NewDIDAdapter(provider, cfg.Blockchain)
	if err != nil {
		log.WithError(err).Fatal("failed to initialise DID adapter")
	}

	aaAdapter, err := sdk.NewAAAdapter(provider, cfg.Blockchain)
	if err != nil {
		log.WithError(err).Fatal("failed to initialise AA adapter")
	}

	vcAdapter, err := sdk.NewVCAdapter(provider, cfg.Blockchain)
	if err != nil {
		log.WithError(err).Fatal("failed to initialise VC adapter")
	}

	aliasAdapter, err := sdk.NewAliasAdapter(provider, cfg.Blockchain)
	if err != nil {
		log.WithError(err).Fatal("failed to initialise Alias adapter")
	}

	blockchainSvc := sdk.NewCompositeAdapter(didAdapter, aaAdapter, vcAdapter, aliasAdapter)
	log.Info("blockchain modular adapters initialised")

	// ── Repositories ───────────────────────────────────────────────────────
	jobRepo := infradb.NewPostgresJobRepository(db)

	// ── Services ───────────────────────────────────────────────────────────
	eventSvc := service.NewEventService(jobRepo, blockchainSvc, log)

	// ── Metrics ────────────────────────────────────────────────────────────
	metrics := pkg.NewMetrics()

	// ── Handler ────────────────────────────────────────────────────────────
	retryCfg := pkg.DefaultRetryConfig(cfg.Worker.MaxRetry, cfg.Worker.RetryBaseDelay)
	handler := worker.NewHandler(eventSvc, retryCfg, metrics, log)

	// ── Worker pool ────────────────────────────────────────────────────────
	pool := worker.NewPool(redisClient, handler, cfg.Worker, cfg.Redis, metrics, log)

	// ── HTTP server (health + metrics) ─────────────────────────────────────
	httpSrv := worker.NewHTTPServer(cfg.Server.Port, log)
	httpSrv.Start()

	// ── Run until signal ───────────────────────────────────────────────────
	log.Info("worker service started")
	pool.Run(ctx) // blocks until ctx cancelled

	// ── Graceful shutdown ──────────────────────────────────────────────────
	log.Info("shutting down HTTP server")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Worker.ShutdownTimeout)
	defer cancel()
	httpSrv.Shutdown(shutdownCtx)

	log.Info("worker service stopped cleanly")
}
