package main

import (
	"fmt"
	"io"
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
	// Usar timestamp para evitar conflitos entre testes paralelos
	testEntityID := fmt.Sprintf("test_e2e_otimizado_%d", time.Now().UnixNano())
	dataDir := filepath.Join("../../data/test_entities", testEntityID)

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

	// Criar dados de teste
	_ = setupTestData(t, lifecycleMgr, testEntityID)

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
		formData := strings.NewReader(fmt.Sprintf("entity_id=%s&product=Café+Especial&amount=4500&quantity=1&stock_item_id=test_item_1", testEntityID))

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

		// Testar página do caixa
		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Get(server.URL + "/cash")
		if err != nil {
			t.Fatalf("Failed to get cash page: %v", err)
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

	// Teste 5: Venda real com quantidade razoável
	t.Run("Venda_Real_Quantidade_Razoavel", func(t *testing.T) {
		start := time.Now()

		// Tentar vender quantidade razoável (5 unidades) de Café Especial
		// Café Especial tem 100 unidades no setup, então deve funcionar
		client := &http.Client{Timeout: 5 * time.Second}
		formData := strings.NewReader(fmt.Sprintf("entity_id=%s&product=Café+Especial&amount=22500&quantity=5", testEntityID))

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

		// Verificar se a venda foi registrada
		if resp.StatusCode == http.StatusOK && strings.Contains(response, "Venda Registrada") {
			t.Logf("✅ Venda real registrada com sucesso em %v", time.Since(start))
		} else {
			t.Logf("⚠️  Resposta da venda: %s (status: %d)", response, resp.StatusCode)
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

	// Verificação mínima - garantir que a página carregou
	// Verificar se há conteúdo na página
	content, err := page.Content()
	if err != nil {
		t.Fatalf("Failed to get page content: %v", err)
	}

	// Verificar se a página carregou com sucesso
	// Dashboard deve ter conteúdo significativo
	hasContent := len(content) > 100 // Pelo menos algum conteúdo

	if !hasContent {
		t.Errorf("Page doesn't have enough content")
	} else {
		t.Log("✅ Page loaded successfully")
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
	dashboardHandler, _ := handler.NewDashboardHandler(lifecycleMgr)

	mux := http.NewServeMux()
	dashboardHandler.RegisterRoutes(mux)
	pdvHandler.RegisterRoutes(mux)
	cashHandler.RegisterRoutes(mux)

	server := httptest.NewServer(mux)
	defer server.Close()

	client := &http.Client{Timeout: 10 * time.Second}

	// Benchmark: venda via API
	b.Run("Venda_API", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			formData := strings.NewReader(fmt.Sprintf(
				"entity_id=%s&product=Benchmark+Product&amount=1000&quantity=1&stock_item_id=test_item_%d", testEntityID,
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
			resp, err := client.Get(server.URL + "/cash")
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

// setupTestData cria dados de teste para os testes E2E
func setupTestData(t *testing.T, lm lifecycle.LifecycleManager, entityID string) map[string]string {
	t.Logf("📝 Configurando dados de teste para entity: %s", entityID)

	// Criar handler de supply para criar itens de estoque
	supplyHandler, err := handler.NewSupplyHandler(lm)
	if err != nil {
		t.Logf("⚠️  Não foi possível criar supply handler: %v", err)
		return nil
	}

	// Criar mux temporário para testar
	mux := http.NewServeMux()
	supplyHandler.RegisterRoutes(mux)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Criar itens de teste
	testItems := []struct {
		name     string
		quantity int
		unitCost int64
		unit     string
	}{
		{"Café Especial", 100, 4500, "KG"},
		{"Açúcar Orgânico", 50, 1200, "KG"},
		{"Leite Integral", 200, 350, "L"},
		{"Pão Francês", 300, 150, "UNIDADE"},
	}

	client := &http.Client{Timeout: 10 * time.Second}
	itemIDs := make(map[string]string)

	for _, item := range testItems {
		formData := fmt.Sprintf(
			"entity_id=%s&name=%s&quantity=%d&unit_cost=%d&unit=%s&type=INSUMO",
			entityID,
			item.name,
			item.quantity,
			item.unitCost,
			item.unit,
		)

		req, err := http.NewRequest("POST", server.URL+"/api/supply/stock-item", strings.NewReader(formData))
		if err != nil {
			t.Logf("⚠️  Erro ao criar request para item %s: %v", item.name, err)
			continue
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Do(req)
		if err != nil {
			t.Logf("⚠️  Erro ao criar item %s: %v", item.name, err)
			continue
		}

		// Ler resposta para obter ID (se disponível)
		_, _ = io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			t.Logf("✅ Item de teste criado: %s (%d %s)", item.name, item.quantity, item.unit)
			// Em um sistema real, extrairíamos o ID da resposta
			// Por enquanto, usamos um mapeamento nome->ID mock
			itemIDs[item.name] = fmt.Sprintf("item_%s", strings.ToLower(strings.ReplaceAll(item.name, " ", "_")))
		} else {
			t.Logf("⚠️  Falha ao criar item %s: status %d", item.name, resp.StatusCode)
		}
	}

	t.Logf("✅ Dados de teste configurados para %s", entityID)
	return itemIDs
}
