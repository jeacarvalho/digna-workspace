package pricing

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewPricingCalculator(t *testing.T) {
	calculator, err := NewPricingCalculator()
	if err != nil {
		t.Fatalf("Failed to create pricing calculator: %v", err)
	}

	if calculator == nil {
		t.Fatal("Calculator should not be nil")
	}
}

func TestPricingCalculator_HandleCalculatePrice(t *testing.T) {
	calculator, err := NewPricingCalculator()
	if err != nil {
		t.Fatalf("Failed to create pricing calculator: %v", err)
	}

	tests := []struct {
		name           string
		query          string
		wantStatusCode int
		wantContains   []string
	}{
		{
			name:           "Valid calculation",
			query:          "material_cost=1500&labor_minutes=120&labor_rate=3000",
			wantStatusCode: http.StatusOK,
			wantContains:   []string{"Preço Justo", "R$", "Materiais", "Seu Trabalho", "Fundo Cooperativo"},
		},
		{
			name:           "Empty values",
			query:          "",
			wantStatusCode: http.StatusOK,
			wantContains:   []string{"Preencha os valores"},
		},
		{
			name:           "Only material cost",
			query:          "material_cost=1000&labor_minutes=0&labor_rate=0",
			wantStatusCode: http.StatusOK,
			wantContains:   []string{"R$ 10.50", "Materiais"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/pdv/pricing/calculate?"+tt.query, nil)
			rr := httptest.NewRecorder()

			calculator.HandleCalculatePrice(rr, req)

			if rr.Code != tt.wantStatusCode {
				t.Errorf("HandleCalculatePrice() status code = %v, want %v", rr.Code, tt.wantStatusCode)
			}

			body := rr.Body.String()
			for _, want := range tt.wantContains {
				if !strings.Contains(body, want) {
					t.Errorf("HandleCalculatePrice() body does not contain %q", want)
				}
			}
		})
	}
}

func TestPricingCalculator_RegisterRoutes(t *testing.T) {
	calculator, err := NewPricingCalculator()
	if err != nil {
		t.Fatalf("Failed to create pricing calculator: %v", err)
	}

	mux := http.NewServeMux()
	calculator.RegisterRoutes(mux)

	// Test that the route is registered
	req := httptest.NewRequest("GET", "/pdv/pricing/calculate?material_cost=1000", nil)
	rr := httptest.NewRecorder()

	// Serve the request through the mux
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Registered route returned status code = %v, want %v", rr.Code, http.StatusOK)
	}
}

func TestCalculateFairPrice(t *testing.T) {
	tests := []struct {
		name         string
		materialCost int64
		laborMinutes int64
		laborRate    int64
		want         int64
	}{
		{
			name:         "Basic cost with work",
			materialCost: 1500, // R$ 15,00
			laborMinutes: 120,  // 2 horas
			laborRate:    3000, // R$ 30,00/hora
			want:         7875, // R$ 78,75 (1500 + 6000 = 7500, +5% = 7875)
		},
		{
			name:         "Only material",
			materialCost: 1000, // R$ 10,00
			laborMinutes: 0,
			laborRate:    0,
			want:         1050, // R$ 10,50 (1000 + 5% = 1050)
		},
		{
			name:         "Only work",
			materialCost: 0,
			laborMinutes: 60,   // 1 hora
			laborRate:    2000, // R$ 20,00/hora
			want:         2100, // R$ 21,00 (2000 + 5% = 2100)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateFairPrice(tt.materialCost, tt.laborMinutes, tt.laborRate)
			if got != tt.want {
				t.Errorf("calculateFairPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLanguageValidation(t *testing.T) {
	calculator, err := NewPricingCalculator()
	if err != nil {
		t.Fatalf("Failed to create pricing calculator: %v", err)
	}

	req := httptest.NewRequest("GET", "/pdv/pricing/calculate?material_cost=1000&labor_minutes=60&labor_rate=2000", nil)
	rr := httptest.NewRecorder()

	calculator.HandleCalculatePrice(rr, req)
	body := rr.Body.String()

	// Verificar termos coloquiais obrigatórios
	requiredTerms := []string{
		"Custo dos Materiais",
		"Seu Tempo",
		"Seu Trabalho",
		"Preço Justo",
		"Fundo Cooperativo",
		"Como funciona",
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

	for _, term := range requiredTerms {
		if !strings.Contains(body, term) {
			t.Errorf("Missing required pedagogical term: %q", term)
		}
	}

	for _, term := range prohibitedTerms {
		if strings.Contains(body, term) {
			t.Errorf("Prohibited accounting jargon found: %q", term)
		}
	}

	t.Log("✅ Validação de linguagem: Termos coloquiais usados, jargões contábeis evitados")
}
