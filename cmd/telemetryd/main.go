package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpadapter "github.com/example/telemetry-rule-service/internal/adapter/http"
	"github.com/example/telemetry-rule-service/internal/adapter/memory"
	"github.com/example/telemetry-rule-service/internal/application/events"
	"github.com/example/telemetry-rule-service/internal/application/ingest"
	"github.com/example/telemetry-rule-service/internal/application/rules"
	"github.com/example/telemetry-rule-service/internal/application/subscriptions"
	"github.com/example/telemetry-rule-service/internal/infrastructure/config"
	"github.com/example/telemetry-rule-service/internal/infrastructure/logging"
)

func main() {
	cfg := config.Load()
	logger := logging.New(cfg.Environment)
	store := memory.NewStore()
	sourceRepo := memory.SourceRepository{Store: store}
	metricRepo := memory.MetricRepository{Store: store}
	ruleRepo := memory.RuleRepository{Store: store}
	eventRepo := memory.EventRepository{Store: store}
	subscriptionRepo := memory.SubscriptionRepository{Store: store}
	server := &http.Server{Addr: cfg.HTTPAddr, Handler: httpadapter.NewServer(sourceRepo, metricRepo, ruleRepo, events.NewService(eventRepo), subscriptions.NewService(subscriptionRepo), ingest.NewService(metricRepo), rules.NewService(metricRepo, ruleRepo, eventRepo), logger).Handler(), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go func() {
		logger.Info("telemetry service started", "addr", cfg.HTTPAddr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server stopped unexpectedly", "error", err)
			os.Exit(1)
		}
	}()
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.ShutdownSeconds)*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
	}
}
