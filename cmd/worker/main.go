package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	"github.com/HARA-DID/did-queueing-engine/internal/callback"
	"github.com/HARA-DID/did-queueing-engine/internal/config"
	infradb "github.com/HARA-DID/did-queueing-engine/internal/infra/db"
	redisinfra "github.com/HARA-DID/did-queueing-engine/internal/infra/redis"
	"github.com/HARA-DID/did-queueing-engine/internal/sdk"
	"github.com/HARA-DID/did-queueing-engine/internal/service"
	"github.com/HARA-DID/did-queueing-engine/internal/worker"
	"github.com/HARA-DID/did-queueing-engine/pkg"
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

	// ── Metrics, Repositories, Callbacks ───────────────────────────────────
	jobRepo := infradb.NewPostgresJobRepository(db)
	callbackRegistry := callback.NewDefaultRegistry()
	metrics := pkg.NewMetrics(prometheus.DefaultRegisterer)
	retryCfg := pkg.DefaultRetryConfig(cfg.Worker.MaxRetry, cfg.Worker.RetryBaseDelay)

	// ── HTTP server (health + metrics) ─────────────────────────────────────
	httpSrv := worker.NewHTTPServer(cfg.Server.Port, log)
	httpSrv.Start()

	log.WithField("worker_count", len(cfg.Blockchain.PrivateKeys)).Info("starting worker pools per identity")

	// ── Worker pool per private key ────────────────────────────────────────
	g, gCtx := errgroup.WithContext(ctx)
	for i, pk := range cfg.Blockchain.PrivateKeys {
		workerIndex := i + 1
		workerConsumerName := fmt.Sprintf("%s-%d", cfg.Worker.ConsumerName, workerIndex)

		workerConsumerNameCaptured := workerConsumerName

		g.Go(func() error {
			workerLog := log.WithField("consumer", workerConsumerNameCaptured)

			provider, err := sdk.NewProvider(cfg.Blockchain, pk)
			if err != nil {
				workerLog.WithError(err).Fatal("failed to initialise blockchain provider")
			}

			didAdapter, err := sdk.NewDIDAdapter(provider, cfg.Blockchain)
			if err != nil {
				workerLog.WithError(err).Fatal("failed to initialise DID adapter")
			}

			aaAdapter, err := sdk.NewAAAdapter(provider, cfg.Blockchain)
			if err != nil {
				workerLog.WithError(err).Fatal("failed to initialise AA adapter")
			}

			vcAdapter, err := sdk.NewVCAdapter(provider, cfg.Blockchain)
			if err != nil {
				workerLog.WithError(err).Fatal("failed to initialise VC adapter")
			}

			aliasAdapter, err := sdk.NewAliasAdapter(provider, cfg.Blockchain)
			if err != nil {
				workerLog.WithError(err).Fatal("failed to initialise Alias adapter")
			}

			blockchainSvc := sdk.NewCompositeAdapter(didAdapter, aaAdapter, vcAdapter, aliasAdapter)

			// Initialize EventService first to populate its callback map
			eventSvc := service.NewEventService(jobRepo, blockchainSvc, log)

			// Initialize and start background transaction confirmation service
			txCheckSvc := service.NewTxCheckService(jobRepo, provider.Chain, callbackRegistry, eventSvc.EventCallbacks, log, cfg.Blockchain.TxCheckChannelBuffer)
			g.Go(func() error {
				txCheckSvc.Start(gCtx)
				return nil
			})

			// Link them
			eventSvc.SetTxCheckService(txCheckSvc)

			handler := worker.NewHandler(eventSvc, retryCfg, metrics, log)

			workerCfg := cfg.Worker
			workerCfg.ConsumerName = workerConsumerNameCaptured

			pool := worker.NewPool(redisClient, handler, jobRepo, workerCfg, cfg.Redis, metrics, log)

			workerLog.Info("worker pool starting")
			pool.Run(gCtx)
			return nil
		})
	}

	_ = g.Wait()

	// ── Graceful shutdown ──────────────────────────────────────────────────
	log.Info("shutting down HTTP server")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Worker.ShutdownTimeout)
	defer cancel()
	httpSrv.Shutdown(shutdownCtx)

	log.Info("worker service stopped cleanly")
}
