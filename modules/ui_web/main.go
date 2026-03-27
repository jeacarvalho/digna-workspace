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
	"github.com/providentia/digna/ui_web/pkg/config"
)

const (
	ShutdownTimeout = 10 * time.Second
)

func main() {
	// Load configuration from environment variables
	cfg := config.Load()

	// Configurar logger estruturado
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	logger.Info("🚀 Iniciando Digna Web Server",
		slog.String("version", "v.0 MVP"),
		slog.String("port", cfg.Port),
		slog.String("data_dir", cfg.DataDir),
		slog.String("log_level", cfg.LogLevel),
	)

	// Inicializar Lifecycle Manager com diretório configurável
	lifecycleMgr := lifecycle.NewSQLiteManagerWithDataDir(cfg.DataDir)
	defer func() {
		logger.Info("🔄 Fechando conexões com banco de dados...")
		lifecycleMgr.CloseAll()
		logger.Info("✅ Conexões fechadas")
	}()

	// Criar server configurado
	server, err := createServer(lifecycleMgr, logger, cfg)
	if err != nil {
		logger.Error("Failed to create server", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Canal para sinais de shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar server em goroutine
	go func() {
		addr := cfg.GetPortString()
		logger.Info("✅ Servidor iniciado",
			slog.String("addr", addr),
			slog.String("url", fmt.Sprintf("http://localhost:%s", cfg.Port)),
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

func createServer(lifecycleMgr lifecycle.LifecycleManager, logger *slog.Logger, cfg *config.Config) (*http.Server, error) {
	mux := http.NewServeMux()

	// Static files (PWA)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Auth handler (deve ser o primeiro)
	authHandler, err := handler.NewAuthHandler(lifecycleMgr)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth handler: %w", err)
	}
	authHandler.RegisterRoutes(mux)

	// Handlers protegidos por autenticação
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
	accountantHandler, err := handler.NewAccountantHandler(lifecycleMgr, authHandler)
	if err != nil {
		return nil, fmt.Errorf("failed to create accountant handler: %w", err)
	}
	accountantHandler.RegisterRoutes(mux)

	// Accountant link handler (RF-12)
	accountantLinkHandler, err := handler.NewAccountantLinkHandler(lifecycleMgr, authHandler)
	if err != nil {
		// Log mas não falha - pode ser implementação parcial
		fmt.Printf("⚠️ Accountant link handler creation warning: %v\n", err)
	} else {
		accountantLinkHandler.RegisterRoutes(mux)
		fmt.Println("✅ Accountant link handler registered (RF-12)")
	}

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

	// Member handler
	memberHandler, err := handler.NewMemberHandler(lifecycleMgr)
	if err != nil {
		// Log mas não falha - pode ser implementação parcial
		fmt.Printf("⚠️ Member handler creation warning: %v\n", err)
	} else {
		memberHandler.RegisterRoutes(mux)
		fmt.Println("✅ Member handler registered")
	}

	// Legal handler (dossiê CADSOL, atas, estatutos)
	legalHandler, err := handler.NewLegalHandler(lifecycleMgr)
	if err != nil {
		// Log mas não falha - pode ser implementação parcial
		fmt.Printf("⚠️ Legal handler creation warning: %v\n", err)
	} else {
		legalHandler.RegisterRoutes(mux)
		fmt.Println("✅ Legal handler registered")
	}

	// DAS MEI handler (RF-27)
	dasMEIHandler, err := handler.NewDASMEIHandler(lifecycleMgr)
	if err != nil {
		// Log mas não falha - pode ser implementação parcial
		fmt.Printf("⚠️ DAS MEI handler creation warning: %v\n", err)
	} else {
		dasMEIHandler.RegisterRoutes(mux)
		fmt.Println("✅ DAS MEI handler registered (RF-27)")
	}

	// Eligibility Profile handler (RF-19)
	eligibilityHandler, err := handler.NewEligibilityHandler(lifecycleMgr)
	if err != nil {
		// Log mas não falha - pode ser implementação parcial
		fmt.Printf("⚠️ Eligibility handler creation warning: %v\n", err)
	} else {
		eligibilityHandler.RegisterRoutes(mux)
		fmt.Println("✅ Eligibility handler registered (RF-19)")
	}

	// Help System handler (RF-30)
	helpHandler, err := handler.NewHelpHandler(lifecycleMgr)
	if err != nil {
		// Log mas não falha - pode ser implementação parcial
		fmt.Printf("⚠️ Help handler creation warning: %v\n", err)
	} else {
		helpHandler.RegisterRoutes(mux)
		fmt.Println("✅ Help handler registered (RF-30)")
	}

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

	// Adicionar middlewares
	authMiddleware := middleware.NewAuthMiddleware(authHandler)
	accountantAuthMiddleware := middleware.NewAccountantAuthMiddleware(authHandler)
	empreendimentoAuthMiddleware := middleware.NewEmpreendimentoAuthMiddleware(authHandler)
	loggerMiddleware := middleware.NewLoggerMiddleware(logger)

	// Encadeamento: auth -> logging -> accountant/empreendimento auth -> mux
	// Primeiro verifica autenticação geral, depois logging, depois tipo de usuário
	// A ORDEM É CRÍTICA: accountantAuth deve vir ANTES de empreendimentoAuth
	// para evitar que contadores sejam redirecionados para /accountant/dashboard em loop
	handler := authMiddleware.Handler(
		loggerMiddleware.Handler(
			accountantAuthMiddleware.Handler(
				empreendimentoAuthMiddleware.Handler(mux),
			),
		),
	)

	// Configure server with timeouts
	server := &http.Server{
		Addr:         cfg.GetPortString(),
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	return server, nil
}
