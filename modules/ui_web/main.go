package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
	"github.com/providentia/digna/ui_web/internal/handler"
)

const (
	DefaultPort     = "8080"
	ShutdownTimeout = 10 * time.Second
)

func main() {
	fmt.Println("🚀 Iniciando Digna Web Server...")
	fmt.Println("📍 Versão: v.0 MVP")
	fmt.Printf("🔗 http://localhost:%s\n\n", DefaultPort)

	// Inicializar Lifecycle Manager
	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	// Criar server configurado
	server, err := createServer(lifecycleMgr)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Canal para sinais de shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar server em goroutine
	go func() {
		addr := ":" + DefaultPort
		fmt.Printf("✅ Servidor iniciado em %s\n", addr)
		fmt.Println("📱 Acesse pelo navegador ou instale o PWA")
		fmt.Println("⏹️  Pressione Ctrl+C para parar")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Aguardar sinal de shutdown
	<-quit
	fmt.Println("🛑 Sinal de shutdown recebido, iniciando desligamento gracioso...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("✅ Servidor desligado com sucesso")
}

func createServer(lifecycleMgr lifecycle.LifecycleManager) (*http.Server, error) {
	mux := http.NewServeMux()

	// Static files (PWA)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Handlers
	pdvHandler, err := handler.NewPDVHandler(lifecycleMgr)
	if err != nil {
		return nil, fmt.Errorf("failed to create PDV handler: %w", err)
	}
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

	// Configure server with timeouts
	server := &http.Server{
		Addr:         ":" + DefaultPort,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return server, nil
}
