package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
	"github.com/providentia/digna/ui_web/internal/handler"
)

// TestSystem_HandlerRoutes valida que handlers básicos estão registrados
func TestSystem_HandlerRoutes(t *testing.T) {
	// Configurar ambiente de teste
	testEntityID := "test_system_integration"
	dataDir := filepath.Join("../../data/test_entities", testEntityID)

	os.RemoveAll(dataDir)
	defer os.RemoveAll(dataDir)

	// Criar lifecycle manager
	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	// Criar servidor
	mux := http.NewServeMux()
	staticDir := http.Dir("static")
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(staticDir)))

	// Registrar handlers conhecidos
	dashboardHandler, err := handler.NewDashboardHandler(lifecycleMgr)
	if err != nil {
		t.Fatalf("Failed to create DashboardHandler: %v", err)
	}
	dashboardHandler.RegisterRoutes(mux)
	t.Log("✅ DashboardHandler registered")

	pdvHandler, err := handler.NewPDVHandler(lifecycleMgr)
	if err != nil {
		t.Fatalf("Failed to create PDVHandler: %v", err)
	}
	pdvHandler.RegisterRoutes(mux)
	t.Log("✅ PDVHandler registered")

	cashHandler, err := handler.NewCashHandler(lifecycleMgr)
	if err != nil {
		t.Fatalf("Failed to create CashHandler: %v", err)
	}
	cashHandler.RegisterRoutes(mux)
	t.Log("✅ CashHandler registered")

	// Tentar registrar MemberHandler (pode falhar se não implementado)
	memberHandler, err := handler.NewMemberHandler(lifecycleMgr)
	if err != nil {
		t.Logf("⚠️ MemberHandler not fully implemented: %v", err)
	} else {
		memberHandler.RegisterRoutes(mux)
		t.Log("✅ MemberHandler registered")
	}

	// Criar servidor de teste
	server := httptest.NewServer(mux)
	defer server.Close()

	// Testar rotas principais
	testCases := []struct {
		name        string
		path        string
		wantStatus  int
		description string
	}{
		// Dashboard
		{"Dashboard root", "/", 200, "Página principal"},
		{"Dashboard with entity", "/?entity_id=test", 200, "Dashboard com entity_id"},

		// PDV
		{"PDV page", "/pdv", 200, "Página PDV"},
		{"PDV with entity", "/pdv?entity_id=test", 200, "PDV com entity_id"},

		// Cash
		{"Cash page", "/cash", 200, "Página Cash"},
		{"Cash with entity", "/cash?entity_id=test", 200, "Cash com entity_id"},

		// Member (pode falhar se não implementado)
		{"Member page", "/members", 200, "Página Members"},
		{"Member with entity", "/members?entity_id=test", 200, "Members com entity_id"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := http.Get(server.URL + tc.path)
			if err != nil {
				t.Errorf("GET %s failed: %v", tc.path, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.wantStatus {
				// Para Member, pode ser 404 se não implementado
				if tc.name == "Member page" || tc.name == "Member with entity" {
					t.Logf("⚠️ Member route not implemented: GET %s = %d", tc.path, resp.StatusCode)
				} else {
					t.Errorf("GET %s = %d, want %d (%s)", tc.path, resp.StatusCode, tc.wantStatus, tc.description)
				}
			} else {
				t.Logf("✅ %s: %s", tc.name, tc.description)
			}
		})
	}
}

// TestSystem_TemplatesExist valida que templates obrigatórios existem
func TestSystem_TemplatesExist(t *testing.T) {
	requiredTemplates := []string{
		"dashboard_simple.html",
		"pdv_simple.html",
		"cash_simple.html",
		"members_simple.html", // NOVO - deve existir após implementação
		"layout.html",
	}

	for _, templateName := range requiredTemplates {
		t.Run("Template_"+templateName, func(t *testing.T) {
			templatePath := filepath.Join("templates", templateName)
			if _, err := os.Stat(templatePath); os.IsNotExist(err) {
				t.Errorf("❌ Template file does not exist: %s", templatePath)

				// Para novos templates, sugerir ação
				if templateName == "members_simple.html" {
					t.Log("💡 Ação: Criar template em modules/ui_web/templates/members_simple.html")
				}
			} else {
				t.Logf("✅ Template exists: %s", templateName)

				// Verificar se template não está vazio
				info, _ := os.Stat(templatePath)
				if info.Size() == 0 {
					t.Errorf("❌ Template is empty: %s", templatePath)
				}
			}
		})
	}
}

// TestSystem_NewHandlerValidation é um template para validar novos handlers
// Copiar e adaptar para cada novo handler implementado
func TestSystem_NewHandlerValidation(t *testing.T) {
	// Configurar
	testEntityID := "test_new_handler_validation"
	dataDir := filepath.Join("../../data/test_entities", testEntityID)

	os.RemoveAll(dataDir)
	defer os.RemoveAll(dataDir)

	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	// VALIDAÇÃO: MemberHandler (exemplo)
	t.Run("Validate_MemberHandler", func(t *testing.T) {
		// 1. Handler pode ser criado?
		memberHandler, err := handler.NewMemberHandler(lifecycleMgr)
		if err != nil {
			t.Errorf("❌ Failed to create MemberHandler: %v", err)
			t.Log("💡 Ação: Verificar se MemberHandler implementa NewMemberHandler corretamente")
			return
		}
		t.Log("✅ MemberHandler created")

		// 2. Rotas podem ser registradas?
		mux := http.NewServeMux()
		memberHandler.RegisterRoutes(mux)
		t.Log("✅ MemberHandler routes registered")

		// 3. Servidor responde?
		server := httptest.NewServer(mux)
		defer server.Close()

		routes := []string{
			"/members",
			"/members?entity_id=test",
		}

		for _, route := range routes {
			resp, err := http.Get(server.URL + route)
			if err != nil {
				t.Errorf("❌ GET %s failed: %v", route, err)
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				t.Errorf("❌ GET %s = %d, want 200", route, resp.StatusCode)
				t.Logf("💡 Ação: Verificar se rota %s está definida em RegisterRoutes", route)
			} else {
				t.Logf("✅ Route %s responds 200", route)
			}
		}
	})
}
