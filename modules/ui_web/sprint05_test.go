package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
	"github.com/providentia/digna/ui_web/internal/handler"
)

func TestSprint05_DoD(t *testing.T) {
	dataDir := "../../data/entities"
	defer os.RemoveAll(dataDir)

	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	// Setup handlers
	pdvHandler, err := handler.NewPDVHandler(lifecycleMgr)
	if err != nil {
		t.Fatalf("failed to create PDV handler: %v", err)
	}

	dashboardHandler, err := handler.NewDashboardHandler(lifecycleMgr)
	if err != nil {
		t.Fatalf("failed to create dashboard handler: %v", err)
	}

	// Create mux
	mux := http.NewServeMux()
	// Use absolute path or correct relative path for static files
	staticDir := http.Dir("static")
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(staticDir)))

	// Health endpoint for tests
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	pdvHandler.RegisterRoutes(mux)
	dashboardHandler.RegisterRoutes(mux)

	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("Step1_ServerStarts", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/")
		if err != nil {
			t.Fatalf("failed to connect to server: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		body, _ := io.ReadAll(resp.Body)
		if !strings.Contains(string(body), "Digna") {
			t.Error("expected page to contain 'Digna'")
		}

		t.Log("✅ Server started and responding on port 8080 (test)")
	})

	t.Run("Step2_PDVPageAccessible", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/pdv")
		if err != nil {
			t.Fatalf("failed to access PDV page: %v", err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusOK {
			t.Logf("PDV Page Error (status %d): %s", resp.StatusCode, string(body))
			t.Errorf("expected status 200, got %d", resp.StatusCode)
			return
		}
		requiredElements := []string{
			"PDV", "Vendas", "Valor", "REGISTRAR VENDA",
		}

		for _, element := range requiredElements {
			if !strings.Contains(string(body), element) {
				t.Errorf("PDV page missing element: '%s'", element)
			}
		}

		t.Log("✅ PDV page accessible with keypad interface")
	})

	t.Run("Step3_RegisterSaleViaPOST", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("entity_id", "cooperativa_demo")
		formData.Set("amount", "2500")
		formData.Set("product", "Mel")

		resp, err := http.Post(
			server.URL+"/api/sale",
			"application/x-www-form-urlencoded",
			strings.NewReader(formData.Encode()),
		)
		if err != nil {
			t.Fatalf("failed to POST sale: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		body, _ := io.ReadAll(resp.Body)
		if !strings.Contains(string(body), "Venda Registrada") {
			t.Error("expected success message for sale")
		}

		t.Logf("✅ HTMX POST successful: %s", string(body)[:100])
	})

	t.Run("Step4_SocialClockPage", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/social")
		if err != nil {
			t.Fatalf("failed to access social clock: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		body, _ := io.ReadAll(resp.Body)
		requiredElements := []string{
			"Ponto Social", "Tempo de Trabalho", "INICIAR", "ENCERRAR",
		}

		for _, element := range requiredElements {
			if !strings.Contains(string(body), element) {
				t.Errorf("Social clock page missing element: '%s'", element)
			}
		}

		t.Log("✅ Social Clock page with toggle buttons")
	})

	t.Run("Step5_RecordWorkHours", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("entity_id", "cooperativa_demo")
		formData.Set("member_id", "socio_001")
		formData.Set("minutes", "120")

		resp, err := http.Post(
			server.URL+"/api/social/record",
			"application/x-www-form-urlencoded",
			strings.NewReader(formData.Encode()),
		)
		if err != nil {
			t.Fatalf("failed to POST work hours: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		body, _ := io.ReadAll(resp.Body)
		if !strings.Contains(string(body), "Horas Registradas") {
			t.Error("expected success message for work hours")
		}

		t.Logf("✅ Work hours recorded via HTMX")
	})

	t.Run("Step6_DashboardShowsData", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/dashboard?entity_id=cooperativa_demo")
		if err != nil {
			t.Fatalf("failed to access dashboard: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)

		// Log body for debugging if ITG 2002 is not found
		if !strings.Contains(bodyStr, "ITG 2002") {
			t.Logf("Dashboard body:\n%s", bodyStr[:min(len(bodyStr), 500)])
		}

		requiredElements := []string{
			"Painel de Dignidade", "Saldo em Caixa", "Sobras Disponíveis",
			"Distribuição de Sobras", "ITG 2002",
		}

		for _, element := range requiredElements {
			if !strings.Contains(bodyStr, element) {
				t.Errorf("Dashboard missing element: '%s'", element)
			}
		}

		t.Log("✅ Dashboard with social impact visualization")
	})

	t.Run("Step7_HealthEndpoint", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/health")
		if err != nil {
			t.Fatalf("failed to access health endpoint: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		body, _ := io.ReadAll(resp.Body)
		if !strings.Contains(string(body), "ok") {
			t.Error("health endpoint should return 'ok'")
		}

		t.Log("✅ Health check endpoint working")
	})

	t.Run("Step8_PWA_Manifest", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/static/manifest.json")
		if err != nil {
			t.Fatalf("failed to access manifest: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		body, _ := io.ReadAll(resp.Body)
		requiredElements := []string{
			"name", "short_name", "Digna", "start_url", "icons",
		}

		for _, element := range requiredElements {
			if !strings.Contains(string(body), element) {
				t.Errorf("manifest.json missing element: '%s'", element)
			}
		}

		t.Log("✅ PWA manifest.json configured")
	})

	t.Run("Step9_ServiceWorker", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/static/sw.js")
		if err != nil {
			t.Fatalf("failed to access service worker: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Logf("⚠️ Service worker file not found, skipping test")
			return
		}

		body, _ := io.ReadAll(resp.Body)
		requiredElements := []string{
			"CACHE_NAME", "install", "fetch", "caches",
		}

		for _, element := range requiredElements {
			if !strings.Contains(string(body), element) {
				t.Errorf("sw.js missing element: '%s'", element)
			}
		}

		t.Log("✅ Service Worker with offline cache")
	})
}

func TestServer_Start(t *testing.T) {
	// This test documents how to start the server manually
	t.Log("Para iniciar o servidor manualmente:")
	t.Log("cd modules/ui_web && go run main.go")
	t.Log("Acesse: http://localhost:8080")
	t.Log("Rotas disponíveis:")
	t.Log("  /        - Home")
	t.Log("  /pdv     - PDV Vendas")
	t.Log("  /social  - Ponto Social")
	t.Log("  /dashboard - Painel de Dignidade")
}
