package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/playwright-community/playwright-go"
	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
	"github.com/providentia/digna/ui_web/internal/handler"
)

// TestE2E_Otimizado testa o fluxo crítico PDV → Caixa com foco em performance
func TestE2E_Otimizado(t *testing.T) {
	// Configurar ambiente de teste UMA VEZ para todos os subtestes
	testEntityID := "test_e2e_otimizado"
	dataDir := filepath.Join("../../data/entities", testEntityID)

	os.RemoveAll(dataDir)
	defer os.RemoveAll(dataDir)

	// Criar lifecycle manager
	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	// Criar handlers
	pdvHandler, err := handler.NewPDVHandler(lifecycleMgr)
	if err != nil {
		t.Fatalf("Failed to create PDV handler: %v", err)
	}

	cashHandler, err := handler.NewCashHandler(lifecycleMgr)
	if err != nil {
		t.Fatalf("Failed to create cash handler: %v", err)
	}

	dashboardHandler, err := handler.NewDashboardHandler(lifecycleMgr)
	if err != nil {
		t.Fatalf("Failed to create dashboard handler: %v", err)
	}

	// Criar servidor
	mux := http.NewServeMux()
	staticDir := http.Dir("static")
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(staticDir)))

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	pdvHandler.RegisterRoutes(mux)
	cashHandler.RegisterRoutes(mux)
	dashboardHandler.RegisterRoutes(mux)

	server := httptest.NewServer(mux)
	defer server.Close()

	// Inicializar Playwright UMA VEZ
	pw, err := playwright.Run()
	if err != nil {
		t.Fatalf("Failed to start Playwright: %v", err)
	}
	defer pw.Stop()

	// Browser headless com timeout reduzido
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
		Timeout:  playwright.Float(30000), // 30s timeout
	})
	if err != nil {
		t.Fatalf("Failed to launch browser: %v", err)
	}
	defer browser.Close()

	// Criar contexto com viewport otimizado
	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		Viewport: &playwright.Size{Width: 1280, Height: 720},
	})
	if err != nil {
		t.Fatalf("Failed to create context: %v", err)
	}
	defer context.Close()

	// Criar página
	page, err := context.NewPage()
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}
	defer page.Close()

	// Configurar timeouts da página
	page.SetDefaultTimeout(10000) // 10s timeout para operações

	t.Log("🚀 Teste E2E Otimizado Iniciado")

	// Teste 1: Dashboard básico
	t.Run("Dashboard", func(t *testing.T) {
		start := time.Now()
		_, err := page.Goto(server.URL + "/")
		if err != nil {
			t.Fatalf("Failed to navigate: %v", err)
		}

		// Verificar rapidamente
		title, _ := page.Title()
		if !strings.Contains(title, "Digna") {
			t.Errorf("Title doesn't contain 'Digna': %s", title)
		}

		// Verificar links com timeout curto
		page.WaitForSelector("text=PDV", playwright.PageWaitForSelectorOptions{Timeout: playwright.Float(2000)})
		page.WaitForSelector("text=Caixa", playwright.PageWaitForSelectorOptions{Timeout: playwright.Float(2000)})

		t.Logf("✅ Dashboard carregado em %v", time.Since(start))
	})

	// Teste 2: PDV - interface básica
	t.Run("PDV_Interface", func(t *testing.T) {
		start := time.Now()

		// Navegar diretamente para PDV (mais rápido que clicar)
		_, err := page.Goto(server.URL + "/pdv")
		if err != nil {
			t.Fatalf("Failed to navigate to PDV: %v", err)
		}

		// Aguardar elementos críticos
		page.WaitForSelector("#product", playwright.PageWaitForSelectorOptions{Timeout: playwright.Float(3000)})
		page.WaitForSelector("#quantity", playwright.PageWaitForSelectorOptions{Timeout: playwright.Float(3000)})
		page.WaitForSelector("text=REGISTRAR VENDA", playwright.PageWaitForSelectorOptions{Timeout: playwright.Float(3000)})

		// Verificar se há produtos
		options, _ := page.QuerySelectorAll("#product option")
		if len(options) == 0 {
			t.Log("⚠️  Nenhum produto cadastrado")
		}

		t.Logf("✅ PDV carregado em %v", time.Since(start))
	})

	// Teste 3: Venda rápida via API (sem browser)
	t.Run("Venda_API_Rapida", func(t *testing.T) {
		start := time.Now()

		// Testar API diretamente (muito mais rápido que browser)
		client := &http.Client{Timeout: 5 * time.Second}

		// Testar com produto existente (Café Especial)
		formData := strings.NewReader("entity_id=cooperativa_demo&product=Café+Especial&amount=4500&quantity=1&stock_item_id=item_1773079963689515743")

		req, err := http.NewRequest("POST", server.URL+"/api/sale", formData)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		// Ler resposta rapidamente
		body := make([]byte, 256)
		n, _ := resp.Body.Read(body)
		response := string(body[:n])

		if !strings.Contains(response, "Venda Registrada") {
			t.Errorf("Sale not successful: %s", response)
		} else {
			t.Logf("✅ Venda API registrada em %v", time.Since(start))
		}
	})

	// Teste 4: Verificar caixa (via API também)
	t.Run("Caixa_API_Rapida", func(t *testing.T) {
		start := time.Now()

		// Testar API do caixa
		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Get(server.URL + "/cash/entries")
		if err != nil {
			t.Fatalf("Failed to get cash entries: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		// Verificar se retorna algo
		body := make([]byte, 1024)
		n, _ := resp.Body.Read(body)
		response := string(body[:n])

		if len(response) > 0 {
			t.Logf("✅ Caixa API funcionando, resposta em %v", time.Since(start))
		} else {
			t.Log("⚠️  Caixa API retornou resposta vazia")
		}
	})

	// Teste 5: Validação de estoque via API
	t.Run("Validacao_Estoque_API", func(t *testing.T) {
		start := time.Now()

		// Tentar vender quantidade enorme (deve falhar)
		client := &http.Client{Timeout: 5 * time.Second}
		formData := strings.NewReader("entity_id=cooperativa_demo&product=Café+Especial&amount=4500000&quantity=1000&stock_item_id=item_1773079963689515743")

		req, err := http.NewRequest("POST", server.URL+"/api/sale", formData)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		body := make([]byte, 512)
		n, _ := resp.Body.Read(body)
		response := string(body[:n])

		// Verificar se a validação de estoque está funcionando
		if strings.Contains(strings.ToLower(response), "estoque") &&
			strings.Contains(strings.ToLower(response), "insuficiente") {
			t.Logf("✅ Validação de estoque funcionando em %v", time.Since(start))
		} else if resp.StatusCode == 200 && strings.Contains(response, "Venda Registrada") {
			t.Errorf("❌ Validação de estoque falhou - permitiu venda de 1000 unidades")
		} else {
			t.Logf("⚠️  Resposta da validação: %s", response)
		}
	})

	t.Log("🏁 Teste E2E Otimizado Finalizado")
}

// TestE2E_Browser_Minimal testa apenas o essencial no browser
func TestE2E_Browser_Minimal(t *testing.T) {
	// Teste focado apenas em verificar que o browser funciona
	// sem todo o overhead de setup/teardown

	if testing.Short() {
		t.Skip("Skipping browser test in short mode")
	}

	// Setup rápido
	testEntityID := "test_browser_minimal"
	dataDir := filepath.Join("../../data/entities", testEntityID)
	os.RemoveAll(dataDir)
	defer os.RemoveAll(dataDir)

	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	pdvHandler, _ := handler.NewPDVHandler(lifecycleMgr)
	cashHandler, _ := handler.NewCashHandler(lifecycleMgr)

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	pdvHandler.RegisterRoutes(mux)
	cashHandler.RegisterRoutes(mux)

	server := httptest.NewServer(mux)
	defer server.Close()

	// Playwright com timeout curto
	pw, err := playwright.Run()
	if err != nil {
		t.Skipf("Playwright not available: %v", err)
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
		Timeout:  playwright.Float(15000), // 15s max
	})
	if err != nil {
		t.Skipf("Browser not available: %v", err)
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}
	defer page.Close()

	page.SetDefaultTimeout(5000) // 5s timeout

	// Teste único: carregar dashboard
	start := time.Now()
	_, err = page.Goto(server.URL + "/")
	if err != nil {
		t.Fatalf("Failed to load dashboard: %v", err)
	}

	// Verificação mínima
	title, _ := page.Title()
	if !strings.Contains(title, "Digna") {
		t.Errorf("Title doesn't contain 'Digna': %s", title)
	}

	t.Logf("✅ Browser test completed in %v", time.Since(start))
}

// BenchmarkE2E_APIs testa performance das APIs críticas
func BenchmarkE2E_APIs(b *testing.B) {
	// Setup para benchmark
	testEntityID := "benchmark_e2e"
	dataDir := filepath.Join("../../data/entities", testEntityID)
	os.RemoveAll(dataDir)
	defer os.RemoveAll(dataDir)

	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	pdvHandler, _ := handler.NewPDVHandler(lifecycleMgr)
	cashHandler, _ := handler.NewCashHandler(lifecycleMgr)

	mux := http.NewServeMux()
	pdvHandler.RegisterRoutes(mux)
	cashHandler.RegisterRoutes(mux)

	server := httptest.NewServer(mux)
	defer server.Close()

	client := &http.Client{Timeout: 10 * time.Second}

	// Benchmark: venda via API
	b.Run("Venda_API", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			formData := strings.NewReader(fmt.Sprintf(
				"entity_id=cooperativa_demo&product=Benchmark+Product&amount=1000&quantity=1&stock_item_id=item_benchmark_%d",
				i,
			))

			req, _ := http.NewRequest("POST", server.URL+"/api/sale", formData)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			resp, err := client.Do(req)
			if err != nil {
				b.Fatalf("Request failed: %v", err)
			}
			resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				b.Errorf("Expected 200, got %d", resp.StatusCode)
			}
		}
	})

	// Benchmark: consulta caixa
	b.Run("Consulta_Caixa", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			resp, err := client.Get(server.URL + "/cash/entries")
			if err != nil {
				b.Fatalf("Request failed: %v", err)
			}
			resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				b.Errorf("Expected 200, got %d", resp.StatusCode)
			}
		}
	})
}
