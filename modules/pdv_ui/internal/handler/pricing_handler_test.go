package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPricingHandler_HandleCalculatePrice(t *testing.T) {
	handler, err := NewPricingHandler()
	if err != nil {
		t.Fatalf("Failed to create pricing handler: %v", err)
	}

	tests := []struct {
		name             string
		queryParams      string
		expectedStatus   int
		expectedContains []string
	}{
		{
			name:           "Cálculo válido",
			queryParams:    "material_cost=1000&labor_minutes=60&labor_rate=2000",
			expectedStatus: http.StatusOK,
			expectedContains: []string{
				"R$ 31,50", // Preço sugerido
				"R$ 10,00", // Material
				"R$ 20,00", // Trabalho
				"R$ 1,50",  // Fundo
				"bg-blue-500",
				"bg-green-500",
				"bg-yellow-500",
			},
		},
		{
			name:           "Valores vazios",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			expectedContains: []string{
				"Preencha os campos",
				"seu tempo e trabalho",
			},
		},
		{
			name:           "Valores negativos (devem ser tratados como zero)",
			queryParams:    "material_cost=-100&labor_minutes=-30&labor_rate=-2000",
			expectedStatus: http.StatusOK,
			expectedContains: []string{
				"Preencha os campos",
			},
		},
		{
			name:           "Apenas material",
			queryParams:    "material_cost=5000&labor_minutes=0&labor_rate=2000",
			expectedStatus: http.StatusOK,
			expectedContains: []string{
				"R$ 52,50", // Preço sugerido
				"R$ 50,00", // Material
				"R$ 0,00",  // Trabalho
				"R$ 2,50",  // Fundo
			},
		},
		{
			name:           "Apenas trabalho",
			queryParams:    "material_cost=0&labor_minutes=45&labor_rate=2400",
			expectedStatus: http.StatusOK,
			expectedContains: []string{
				"R$ 18,90", // Preço sugerido (R$ 18,00 + 5%)
				"R$ 0,00",  // Material
				"R$ 18,00", // Trabalho
				"R$ 0,90",  // Fundo
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/pdv/pricing/calculate?"+tt.queryParams, nil)
			rr := httptest.NewRecorder()

			handler.HandleCalculatePrice(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Status code incorreto: esperado %d, obtido %d", tt.expectedStatus, rr.Code)
			}

			body := rr.Body.String()
			for _, expected := range tt.expectedContains {
				if !strings.Contains(body, expected) {
					t.Errorf("Resposta não contém '%s'", expected)
				}
			}

			// Validar que a resposta é HTML
			contentType := rr.Header().Get("Content-Type")
			if !strings.Contains(contentType, "text/html") {
				t.Errorf("Content-Type incorreto: esperado text/html, obtido %s", contentType)
			}

			// Validar que não há erros de template
			if strings.Contains(body, "Erro ao renderizar template") {
				t.Errorf("Erro ao renderizar template: %s", body)
			}
		})
	}
}

func TestNewPricingHandler(t *testing.T) {
	handler, err := NewPricingHandler()
	if err != nil {
		t.Fatalf("NewPricingHandler() error = %v", err)
	}
	if handler == nil {
		t.Fatal("NewPricingHandler() retornou nil")
	}
	if handler.tmpl == nil {
		t.Fatal("Template não foi carregado")
	}
}

func TestPricingHandler_LanguageValidation(t *testing.T) {
	// Validar que a linguagem é coloquial e não usa jargões contábeis
	handler, err := NewPricingHandler()
	if err != nil {
		t.Fatalf("Failed to create pricing handler: %v", err)
	}

	req := httptest.NewRequest("GET", "/pdv/pricing/calculate?material_cost=1000&labor_minutes=60&labor_rate=2000", nil)
	rr := httptest.NewRecorder()

	handler.HandleCalculatePrice(rr, req)

	body := rr.Body.String()

	// Termos PROIBIDOS (jargões contábeis)
	prohibitedTerms := []string{
		"Markup",
		"Lucro Líquido",
		"CPV",
		"Custo do Produto Vendido",
		"Débito",
		"Crédito",
		"Profit",
		"Margin",
		"ROI",
	}

	for _, term := range prohibitedTerms {
		if strings.Contains(strings.ToLower(body), strings.ToLower(term)) {
			t.Errorf("Linguagem proibida encontrada: '%s'", term)
		}
	}

	// Termos OBRIGATÓRIOS (linguagem coloquial)
	requiredTerms := []string{
		"Custo do Material",
		"Seu Tempo",
		"Seu Trabalho",
		"Preço Sugerido",
		"Preço justo",
		"Fundo da Cooperativa",
		"Valor do seu trabalho",
	}

	for _, term := range requiredTerms {
		if !strings.Contains(body, term) {
			t.Errorf("Linguagem coloquial não encontrada: '%s'", term)
		}
	}

	t.Log("✅ Validação de linguagem: Termos coloquiais usados, jargões contábeis evitados")
}

func TestPricingHandler_HTMXIntegration(t *testing.T) {
	// Validar que o componente está preparado para HTMX
	handler, err := NewPricingHandler()
	if err != nil {
		t.Fatalf("Failed to create pricing handler: %v", err)
	}

	req := httptest.NewRequest("GET", "/pdv/pricing/calculate?material_cost=1000", nil)
	rr := httptest.NewRecorder()

	handler.HandleCalculatePrice(rr, req)

	body := rr.Body.String()

	// Validar atributos HTMX
	htmxAttributes := []string{
		"hx-get=",
		"hx-trigger=",
		"hx-target=",
		"keyup changed delay",
	}

	for _, attr := range htmxAttributes {
		if !strings.Contains(body, attr) {
			t.Errorf("Atributo HTMX não encontrado: '%s'", attr)
		}
	}

	// Validar classes Tailwind
	tailwindClasses := []string{
		"bg-white",
		"rounded-lg",
		"shadow",
		"grid",
		"grid-cols-",
		"bg-blue-500",
		"bg-green-500",
		"bg-yellow-500",
	}

	for _, class := range tailwindClasses {
		if !strings.Contains(body, class) {
			t.Errorf("Classe Tailwind não encontrada: '%s'", class)
		}
	}

	t.Log("✅ Validação HTMX/Tailwind: Componente interativo pronto")
}
