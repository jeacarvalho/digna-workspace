package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
	"github.com/providentia/digna/ui_web/internal/handler"
	"github.com/providentia/digna/ui_web/internal/middleware"
)

const (
	DefaultPort     = "8088"
	ShutdownTimeout = 10 * time.Second
)

func main() {
	// Configurar logger estruturado
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	logger.Info("🚀 Iniciando Digna Web Server",
		slog.String("version", "v.0 MVP"),
		slog.String("port", DefaultPort),
	)

	// Inicializar Lifecycle Manager
	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer func() {
		logger.Info("🔄 Fechando conexões com banco de dados...")
		lifecycleMgr.CloseAll()
		logger.Info("✅ Conexões fechadas")
	}()

	// Criar server configurado
	server, err := createServer(lifecycleMgr, logger)
	if err != nil {
		logger.Error("Failed to create server", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Canal para sinais de shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar server em goroutine
	go func() {
		addr := ":" + DefaultPort
		logger.Info("✅ Servidor iniciado",
			slog.String("addr", addr),
			slog.String("url", fmt.Sprintf("http://localhost:%s", DefaultPort)),
		)
		logger.Info("📱 Acesse pelo navegador ou instale o PWA")
		logger.Info("⏹️  Pressione Ctrl+C para parar")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// Aguardar sinal de shutdown
	sig := <-quit
	logger.Info("🛑 Sinal de shutdown recebido",
		slog.String("signal", sig.String()),
		slog.String("action", "iniciando desligamento gracioso"),
	)

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	shutdownStart := time.Now()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown",
			slog.String("error", err.Error()),
			slog.Duration("timeout", ShutdownTimeout),
		)
		os.Exit(1)
	}

	shutdownDuration := time.Since(shutdownStart)
	logger.Info("✅ Servidor desligado com sucesso",
		slog.Duration("shutdown_duration", shutdownDuration),
	)
}

func createServer(lifecycleMgr lifecycle.LifecycleManager, logger *slog.Logger) (*http.Server, error) {
	mux := http.NewServeMux()

	// Static files (PWA)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Handlers
	pdvHandler, err := handler.NewPDVHandler(lifecycleMgr)
	if err != nil {
		fmt.Printf("DEBUG: Erro ao criar PDV handler: %v\n", err)
		return nil, fmt.Errorf("failed to create PDV handler: %w", err)
	}
	fmt.Printf("DEBUG: PDV handler criado com sucesso\n")
	pdvHandler.RegisterRoutes(mux)

	dashboardHandler, err := handler.NewDashboardHandler(lifecycleMgr)
	if err != nil {
		return nil, fmt.Errorf("failed to create dashboard handler: %w", err)
	}
	dashboardHandler.RegisterRoutes(mux)

	cashHandler, err := handler.NewCashHandler(lifecycleMgr)
	if err != nil {
		return nil, fmt.Errorf("failed to create cash handler: %w", err)
	}
	cashHandler.RegisterRoutes(mux)

	// Accountant dashboard handler
	accountantHandler, err := handler.NewAccountantHandler(lifecycleMgr)
	if err != nil {
		return nil, fmt.Errorf("failed to create accountant handler: %w", err)
	}
	accountantHandler.RegisterRoutes(mux)

	// Supply handler
	supplyHandler, err := handler.NewSupplyHandler(lifecycleMgr)
	if err != nil {
		return nil, fmt.Errorf("failed to create supply handler: %w", err)
	}
	supplyHandler.RegisterRoutes(mux)

	// Budget handler
	budgetHandler, err := handler.NewBudgetHandler(lifecycleMgr)
	if err != nil {
		return nil, fmt.Errorf("failed to create budget handler: %w", err)
	}
	budgetHandler.RegisterRoutes(mux)

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"status":"ok","version":"v.0"}`)
	})

	// Readiness check
	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"ready":true}`)
	})

	// Adicionar middleware de logging
	loggerMiddleware := middleware.NewLoggerMiddleware(logger)
	handler := loggerMiddleware.Handler(mux)

	// Configure server with timeouts
	server := &http.Server{
		Addr:         ":" + DefaultPort,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	return server, nil
}
