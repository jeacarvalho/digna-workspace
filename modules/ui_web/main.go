package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
	"github.com/providentia/digna/ui_web/internal/handler"
)

const (
	DefaultPort = "8080"
)

func main() {
	fmt.Println("🚀 Iniciando Digna Web Server...")
	fmt.Println("📍 Versão: v.0 MVP")
	fmt.Printf("🔗 http://localhost:%s\n\n", DefaultPort)

	// Inicializar Lifecycle Manager
	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	// Criar mux principal
	mux := http.NewServeMux()

	// Static files (PWA)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Handlers
	pdvHandler, err := handler.NewPDVHandler(lifecycleMgr)
	if err != nil {
		log.Fatalf("Failed to create PDV handler: %v", err)
	}
	pdvHandler.RegisterRoutes(mux)

	dashboardHandler, err := handler.NewDashboardHandler(lifecycleMgr)
	if err != nil {
		log.Fatalf("Failed to create dashboard handler: %v", err)
	}
	dashboardHandler.RegisterRoutes(mux)

	// Cash Handler
	cashHandler, err := handler.NewCashHandler(lifecycleMgr)
	if err != nil {
		log.Fatalf("Failed to create cash handler: %v", err)
	}
	cashHandler.RegisterRoutes(mux)

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"status":"ok","version":"v.0"}`)
	})

	// Start server
	addr := ":" + DefaultPort
	fmt.Printf("✅ Servidor iniciado em %s\n", addr)
	fmt.Println("📱 Acesse pelo navegador ou instale o PWA")

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
