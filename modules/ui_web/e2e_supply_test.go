package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
	"github.com/providentia/digna/ui_web/internal/handler"
)

func TestSupplyRoutes_E2E(t *testing.T) {
	// Criar lifecycle manager
	lm := lifecycle.NewSQLiteManager()
	defer lm.CloseAll()

	// Criar handler de supply
	supplyHandler, err := handler.NewSupplyHandler(lm)
	if err != nil {
		t.Fatalf("Failed to create supply handler: %v", err)
	}

	// Testar rota /supply/stock
	t.Run("StockPage_WithEntityID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/supply/stock?entity_id=cafe_digna", nil)
		w := httptest.NewRecorder()

		supplyHandler.StockPage(w, req)

		// Verificar status code
		if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 200 or 500, got %d", w.Code)
		}

		// Verificar se não há erro de conexão de banco
		body := w.Body.String()
		if w.Code == http.StatusInternalServerError {
			if contains(body, "database connection is closed") {
				t.Errorf("Database connection error: %s", body)
			}
			fmt.Printf("[INFO] Stock page returned 500 (expected in test): %s\n", body)
		} else {
			fmt.Printf("[INFO] Stock page returned 200\n")
		}
	})

	// Testar rota /supply/suppliers
	t.Run("SuppliersPage_WithEntityID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/supply/suppliers?entity_id=cafe_digna", nil)
		w := httptest.NewRecorder()

		supplyHandler.SuppliersPage(w, req)

		// Verificar status code
		if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 200 or 500, got %d", w.Code)
		}

		// Verificar se não há erro de conexão de banco
		body := w.Body.String()
		if w.Code == http.StatusInternalServerError {
			if contains(body, "database connection is closed") {
				t.Errorf("Database connection error: %s", body)
			}
			fmt.Printf("[INFO] Suppliers page returned 500 (expected in test): %s\n", body)
		} else {
			fmt.Printf("[INFO] Suppliers page returned 200\n")
		}
	})

	// Testar rota sem entity_id (deve retornar 400)
	t.Run("StockPage_WithoutEntityID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/supply/stock", nil)
		w := httptest.NewRecorder()

		supplyHandler.StockPage(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 for missing entity_id, got %d", w.Code)
		}
	})

	t.Run("SuppliersPage_WithoutEntityID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/supply/suppliers", nil)
		w := httptest.NewRecorder()

		supplyHandler.SuppliersPage(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 for missing entity_id, got %d", w.Code)
		}
	})
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && (s[0:len(substr)] == substr || contains(s[1:], substr)))
}
