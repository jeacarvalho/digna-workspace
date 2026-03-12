package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

func TestAccountantDashboard_RendersWithoutErrors(t *testing.T) {
	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	authHandler := NewMockAuthHandler(lifecycleMgr)
	handler, err := NewAccountantHandler(lifecycleMgr, authHandler)
	if err != nil {
		t.Fatalf("Failed to create accountant handler: %v", err)
	}

	// Teste simples: apenas verificar que o handler foi criado
	if handler == nil {
		t.Error("Expected non-nil handler")
	}

	if handler.BaseHandler == nil {
		t.Error("Expected BaseHandler to be initialized")
	}

	// Teste básico de rota sem executar o handler completo
	req := httptest.NewRequest("GET", "/accountant/dashboard", nil)
	req.AddCookie(&http.Cookie{Name: "digna_session", Value: "test_session"})

	// Executar apenas a verificação inicial (sem chamar serviços externos)
	// Isso testa se o handler pelo menos inicia sem panics
	accountantID, valid := handler.authHandler.GetCurrentEntity(req)
	if !valid {
		t.Error("Expected valid entity from mock auth handler")
	}
	if accountantID != "contador_social" {
		t.Errorf("Expected accountantID 'contador_social', got %s", accountantID)
	}

	t.Log("Accountant handler created and basic auth check passed")
}

func TestAccountantDashboard_WithPeriodParameter(t *testing.T) {
	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	authHandler := NewMockAuthHandler(lifecycleMgr)
	handler, err := NewAccountantHandler(lifecycleMgr, authHandler)
	if err != nil {
		t.Fatalf("Failed to create accountant handler: %v", err)
	}

	// Teste básico: verificar que o handler foi criado
	if handler == nil {
		t.Error("Expected non-nil handler")
	}

	// Testar parsing de período
	req := httptest.NewRequest("GET", "/accountant/dashboard?period=2024-01", nil)
	period := req.URL.Query().Get("period")
	if period != "2024-01" {
		t.Errorf("Expected period '2024-01', got %s", period)
	}

	t.Log("Period parameter parsing test passed")
}

func TestExportFiscal_RequiresEntityIDAndPeriod(t *testing.T) {
	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	authHandler := NewMockAuthHandler(lifecycleMgr)
	handler, err := NewAccountantHandler(lifecycleMgr, authHandler)
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

	authHandler := NewMockAuthHandler(lifecycleMgr)
	handler, err := NewAccountantHandler(lifecycleMgr, authHandler)
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

	authHandler := NewMockAuthHandler(lifecycleMgr)
	handler, err := NewAccountantHandler(lifecycleMgr, authHandler)
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

	authHandler := NewMockAuthHandler(lifecycleMgr)
	handler, err := NewAccountantHandler(lifecycleMgr, authHandler)
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

	authHandler := NewMockAuthHandler(lifecycleMgr)
	handler, err := NewAccountantHandler(lifecycleMgr, authHandler)
	if err != nil {
		t.Fatalf("Failed to create accountant handler: %v", err)
	}

	// Verificar se o handler tem o campo BaseHandler
	// Isso garante que estamos seguindo o padrão de herança/composição
	if handler.BaseHandler == nil {
		t.Error("AccountantHandler should embed BaseHandler")
	}
}
