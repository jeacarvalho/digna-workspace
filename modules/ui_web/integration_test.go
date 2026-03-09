package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
	"github.com/providentia/digna/ui_web/internal/handler"
)

func TestPricingCalculatorIntegration(t *testing.T) {
	// Criar lifecycle manager
	dataDir := "../../data/entities"
	defer func() {
		// Cleanup after test
		os.RemoveAll(dataDir)
	}()

	lm := lifecycle.NewSQLiteManager()
	defer lm.CloseAll()

	// Criar handler PDV
	pdvHandler, err := handler.NewPDVHandler(lm)
	if err != nil {
		t.Fatalf("Failed to create PDV handler: %v", err)
	}

	// Criar mux e registrar rotas
	mux := http.NewServeMux()
	pdvHandler.RegisterRoutes(mux)

	// Testar se a rota de pricing está registrada
	t.Run("Pricing route registered", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/pdv/pricing/calculate?material_cost=1000&labor_minutes=60&labor_rate=2000", nil)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Pricing route returned status code = %v, want %v", rr.Code, http.StatusOK)
		}

		body := rr.Body.String()
		if !strings.Contains(body, "Preço Justo") {
			t.Errorf("Pricing calculator response missing 'Preço Justo'")
		}
		if !strings.Contains(body, "Seu Trabalho") {
			t.Errorf("Pricing calculator response missing 'Seu Trabalho'")
		}
		if !strings.Contains(body, "Fundo Cooperativo") {
			t.Errorf("Pricing calculator response missing 'Fundo Cooperativo'")
		}
	})

	// Testar se a página PDV carrega com o container da calculadora
	t.Run("PDV page includes pricing calculator", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/pdv", nil)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("PDV page returned status code = %v, want %v", rr.Code, http.StatusOK)
		}

		body := rr.Body.String()
		if !strings.Contains(body, "pricing-calculator-container") {
			t.Errorf("PDV page missing pricing calculator container")
		}
		if !strings.Contains(body, "hx-get=\"/pdv/pricing/calculate\"") {
			t.Errorf("PDV page missing HTMX trigger for pricing calculator")
		}
	})

	// Testar cálculo válido
	t.Run("Valid price calculation", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/pdv/pricing/calculate?material_cost=1500&labor_minutes=120&labor_rate=3000", nil)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Price calculation returned status code = %v, want %v", rr.Code, http.StatusOK)
		}

		body := rr.Body.String()
		// Verificar termos pedagógicos
		requiredTerms := []string{
			"Custo dos Materiais",
			"Seu Tempo",
			"Valor da Hora",
			"Preço Justo Calculado",
			"Composição do Preço",
			"Como funciona",
		}

		for _, term := range requiredTerms {
			if !strings.Contains(body, term) {
				t.Errorf("Missing required pedagogical term: %q", term)
			}
		}

		// Verificar que jargões contábeis NÃO estão presentes
		prohibitedTerms := []string{
			"Markup",
			"Net Profit",
			"COGS",
			"Cost of Goods Sold",
			"Debit",
			"Credit",
		}

		for _, term := range prohibitedTerms {
			if strings.Contains(body, term) {
				t.Errorf("Prohibited accounting jargon found: %q", term)
			}
		}

		t.Log("✅ Integração da calculadora de preços validada com sucesso")
	})
}
