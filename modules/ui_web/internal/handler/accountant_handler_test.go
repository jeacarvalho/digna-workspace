package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

func TestAccountantDashboard_RendersWithoutErrors(t *testing.T) {
	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	handler, err := NewAccountantHandler(lifecycleMgr)
	if err != nil {
		t.Fatalf("Failed to create accountant handler: %v", err)
	}

	req := httptest.NewRequest("GET", "/accountant/dashboard", nil)
	w := httptest.NewRecorder()

	handler.Dashboard(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if strings.Contains(body, "template error") {
		t.Error("Template rendering failed")
	}

	// Verificar se o título está presente
	if !strings.Contains(body, "Painel do Contador Social") {
		t.Error("Expected title 'Painel do Contador Social' not found")
	}
}

func TestAccountantDashboard_WithPeriodParameter(t *testing.T) {
	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	handler, err := NewAccountantHandler(lifecycleMgr)
	if err != nil {
		t.Fatalf("Failed to create accountant handler: %v", err)
	}

	req := httptest.NewRequest("GET", "/accountant/dashboard?period=2024-01", nil)
	w := httptest.NewRecorder()

	handler.Dashboard(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "2024-01") {
		t.Error("Expected period '2024-01' not found in response")
	}
}

func TestExportFiscal_RequiresEntityIDAndPeriod(t *testing.T) {
	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	handler, err := NewAccountantHandler(lifecycleMgr)
	if err != nil {
		t.Fatalf("Failed to create accountant handler: %v", err)
	}

	// Teste sem parâmetros
	req := httptest.NewRequest("GET", "/accountant/export", nil)
	w := httptest.NewRecorder()

	handler.ExportFiscal(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for missing parameters, got %d", w.Code)
	}

	// Teste com entity_id mas sem period
	req = httptest.NewRequest("GET", "/accountant/export?entity_id=test_entity", nil)
	w = httptest.NewRecorder()

	handler.ExportFiscal(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for missing period, got %d", w.Code)
	}

	// Teste com period mas sem entity_id
	req = httptest.NewRequest("GET", "/accountant/export?period=2024-01", nil)
	w = httptest.NewRecorder()

	handler.ExportFiscal(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for missing entity_id, got %d", w.Code)
	}
}

func TestExportFiscal_URLPathParameters(t *testing.T) {
	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	handler, err := NewAccountantHandler(lifecycleMgr)
	if err != nil {
		t.Fatalf("Failed to create accountant handler: %v", err)
	}

	// Teste com parâmetros na URL (novo padrão)
	req := httptest.NewRequest("GET", "/accountant/export/test_entity/2024-01", nil)
	w := httptest.NewRecorder()

	handler.ExportFiscal(w, req)

	// Como não há dados reais, esperamos um erro interno (500) ou não encontrado (404)
	// mas não um bad request (400) pois os parâmetros foram fornecidos
	if w.Code == http.StatusBadRequest {
		t.Errorf("Expected status other than 400 for valid URL parameters, got %d", w.Code)
	}
}

func TestAccountantRegisterRoutes(t *testing.T) {
	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	handler, err := NewAccountantHandler(lifecycleMgr)
	if err != nil {
		t.Fatalf("Failed to create accountant handler: %v", err)
	}

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	// Verificar se as rotas estão registradas
	// Não há uma maneira direta de testar isso sem fazer requests reais,
	// mas podemos pelo menos garantir que o método não panica
	t.Log("RegisterRoutes executed without panic")
}

func TestNewAccountantHandler_CreatesSuccessfully(t *testing.T) {
	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	handler, err := NewAccountantHandler(lifecycleMgr)
	if err != nil {
		t.Fatalf("Failed to create accountant handler: %v", err)
	}

	if handler == nil {
		t.Error("Expected non-nil handler")
	}

	if handler.BaseHandler == nil {
		t.Error("Expected BaseHandler to be initialized")
	}
}

func TestAccountantHandler_ImplementsBaseHandler(t *testing.T) {
	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	handler, err := NewAccountantHandler(lifecycleMgr)
	if err != nil {
		t.Fatalf("Failed to create accountant handler: %v", err)
	}

	// Verificar se o handler tem o campo BaseHandler
	// Isso garante que estamos seguindo o padrão de herança/composição
	if handler.BaseHandler == nil {
		t.Error("AccountantHandler should embed BaseHandler")
	}
}
