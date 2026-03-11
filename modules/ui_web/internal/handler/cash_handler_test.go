package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

func TestCashPage_RendersWithoutErrors(t *testing.T) {
	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	handler, err := NewCashHandler(lifecycleMgr)
	if err != nil {
		t.Fatalf("Failed to create cash handler: %v", err)
	}

	req := httptest.NewRequest("GET", "/cash?entity_id=test_cash", nil)
	w := httptest.NewRecorder()

	handler.CashPage(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if strings.Contains(body, "Failed to render template") {
		t.Error("Template rendering failed")
	}
}

func TestCashPage_RequiresEntityID(t *testing.T) {
	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	handler, err := NewCashHandler(lifecycleMgr)
	if err != nil {
		t.Fatalf("Failed to create cash handler: %v", err)
	}

	req := httptest.NewRequest("GET", "/cash", nil)
	w := httptest.NewRecorder()

	handler.CashPage(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for missing entity_id, got %d", w.Code)
	}
}
