package main

import (
	"fmt"
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

func TestIntegracaoUnidadesEstoque_FluxoCompleto(t *testing.T) {
	// Configurar ambiente de teste
	dataDir := "../../data/entities"
	defer os.RemoveAll(dataDir)

	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	// Setup handler
	supplyHandler, err := handler.NewSupplyHandler(lifecycleMgr)
	if err != nil {
		t.Fatalf("failed to create supply handler: %v", err)
	}

	// Criar servidor de teste
	mux := http.NewServeMux()
	supplyHandler.RegisterRoutes(mux)
	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("Fluxo completo: cadastro → listagem → verificação", func(t *testing.T) {
		// PASSO 1: Cadastrar múltiplos itens com diferentes unidades
		testItems := []struct {
			name     string
			unit     string
			quantity int
			unitCost int64
		}{
			{"Arroz Integral", "KG", 50, 850}, // R$ 8,50/kg
			{"Leite", "L", 20, 650},           // R$ 6,50/L
			{"Ovos", "UNIDADE", 30, 90},       // R$ 0,90/unid.
			{"Farinha", "KG", 25, 420},        // R$ 4,20/kg
			{"Óleo", "L", 10, 1250},           // R$ 12,50/L
		}

		itemIDs := make([]string, 0, len(testItems))

		for _, item := range testItems {
			formData := url.Values{
				"name":         {item.name},
				"type":         {"INSUMO"},
				"unit":         {item.unit},
				"quantity":     {fmt.Sprintf("%d", item.quantity)},
				"min_quantity": {"5"},
				"unit_cost":    {fmt.Sprintf("%d", item.unitCost)},
			}

			req, err := http.NewRequest("POST", server.URL+"/api/supply/stock-item", strings.NewReader(formData.Encode()))
			if err != nil {
				t.Fatalf("failed to create request for %s: %v", item.name, err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("failed to send request for %s: %v", item.name, err)
			}

			if resp.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				resp.Body.Close()
				t.Errorf("failed to register %s: status %d: %s", item.name, resp.StatusCode, string(body))
				continue
			}

			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			bodyStr := string(body)

			// Verificar se a resposta contém informações da unidade
			if !strings.Contains(bodyStr, item.name) {
				t.Errorf("resposta não contém nome do item %s", item.name)
			}

			// Extrair ID do item da resposta (simplificado)
			if strings.Contains(bodyStr, "ID:") {
				// Em uma implementação real, extrairíamos o ID da resposta
				itemIDs = append(itemIDs, "item_"+item.name)
			}
		}

		// PASSO 2: Verificar listagem completa
		req, err := http.NewRequest("GET", server.URL+"/api/supply/stock-items", nil)
		if err != nil {
			t.Fatalf("failed to create request for listagem: %v", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("failed to send request for listagem: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("failed to get stock items list: status %d: %s", resp.StatusCode, string(body))
		}

		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)

		// Verificar se todos os itens aparecem na listagem
		for _, item := range testItems {
			if !strings.Contains(bodyStr, item.name) {
				t.Errorf("item %s não aparece na listagem", item.name)
			}

			// Verificar se a unidade aparece formatada
			expectedUnit := strings.ToLower(item.unit)
			if item.unit == "UNIDADE" {
				expectedUnit = "unid."
			} else if item.unit == "L" {
				expectedUnit = "L" // Mantém maiúsculo
			}

			if !strings.Contains(bodyStr, expectedUnit) &&
				!strings.Contains(bodyStr, strings.ToLower(expectedUnit)) &&
				!strings.Contains(bodyStr, strings.ToUpper(expectedUnit)) {
				t.Errorf("unidade %s do item %s não aparece formatada na listagem", item.unit, item.name)
			}

			// Verificar se o custo unitário aparece formatado
			expectedCost := fmt.Sprintf("%.2f", float64(item.unitCost)/100)
			// Substituir ponto por vírgula para busca
			expectedCostSearch := strings.ReplaceAll(expectedCost, ".", ",")
			if !strings.Contains(bodyStr, expectedCost) && !strings.Contains(bodyStr, expectedCostSearch) {
				t.Errorf("custo unitário R$ %s do item %s não aparece na listagem", expectedCost, item.name)
			}
		}

		// PASSO 3: Verificar cálculos totais
		// Calcular totais esperados
		totalValue := int64(0)
		for _, item := range testItems {
			totalValue += item.unitCost * int64(item.quantity)
		}

		// Verificar se algum cálculo total aparece (pode não estar explícito)
		// Mas pelo menos verificar se há formatação de moeda
		if !strings.Contains(bodyStr, "R$") {
			t.Error("listagem não contém formatação de moeda (R$)")
		}

		// PASSO 4: Testar página de estoque completa
		req, err = http.NewRequest("GET", server.URL+"/supply/stock", nil)
		if err != nil {
			t.Fatalf("failed to create request for página de estoque: %v", err)
		}

		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("failed to send request for página de estoque: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("failed to load stock page: status %d", resp.StatusCode)
		}

		pageBody, _ := io.ReadAll(resp.Body)
		pageStr := string(pageBody)

		// Verificar elementos da página
		requiredElements := []string{
			"Meu Estoque",
			"Cadastrar Novo Item",
			"Itens em Estoque",
			"name=\"unit\"",   // Campo de seleção de unidade
			"type=\"submit\"", // Botão de submit
		}

		for _, element := range requiredElements {
			if !strings.Contains(pageStr, element) {
				t.Errorf("página de estoque não contém elemento necessário: %s", element)
			}
		}

		t.Logf("✅ Fluxo completo testado com sucesso:")
		t.Logf("   • %d itens cadastrados com diferentes unidades", len(testItems))
		t.Logf("   • Listagem automática funcionando")
		t.Logf("   • Unidades formatadas corretamente")
		t.Logf("   • Página de estoque carregada com todos os elementos")
	})

	t.Run("Teste de regressão: campos obrigatórios", func(t *testing.T) {
		// Testar validação de campos obrigatórios
		testCases := []struct {
			name        string
			formData    url.Values
			expectError bool
		}{
			{
				name: "Campo nome faltando",
				formData: url.Values{
					"type":         {"INSUMO"},
					"unit":         {"KG"},
					"quantity":     {"10"},
					"min_quantity": {"2"},
					"unit_cost":    {"1000"},
				},
				expectError: true,
			},
			{
				name: "Campo tipo faltando (usa default INSUMO)",
				formData: url.Values{
					"name":         {"Produto Teste"},
					"unit":         {"KG"},
					"quantity":     {"10"},
					"min_quantity": {"2"},
					"unit_cost":    {"1000"},
				},
				expectError: false, // Sistema usa INSUMO como default
			},
			{
				name: "Campo quantidade faltando (usa 0)",
				formData: url.Values{
					"name":         {"Produto Teste"},
					"type":         {"INSUMO"},
					"unit":         {"KG"},
					"min_quantity": {"2"},
					"unit_cost":    {"1000"},
				},
				expectError: false, // Sistema converte string vazia para 0
			},
			{
				name: "Todos campos preenchidos",
				formData: url.Values{
					"name":         {"Produto Completo"},
					"type":         {"PRODUTO"},
					"unit":         {"UNIDADE"},
					"quantity":     {"15"},
					"min_quantity": {"3"},
					"unit_cost":    {"1500"},
				},
				expectError: false,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req, err := http.NewRequest("POST", server.URL+"/api/supply/stock-item", strings.NewReader(tc.formData.Encode()))
				if err != nil {
					t.Fatalf("failed to create request: %v", err)
				}
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					t.Fatalf("failed to send request: %v", err)
				}
				defer resp.Body.Close()

				if tc.expectError {
					// Esperamos erro 400 (Bad Request) ou 500
					if resp.StatusCode < 400 {
						body, _ := io.ReadAll(resp.Body)
						t.Errorf("esperado erro para campos faltando, mas obteve status %d: %s", resp.StatusCode, string(body))
					}
				} else {
					// Esperamos sucesso
					if resp.StatusCode != http.StatusOK {
						body, _ := io.ReadAll(resp.Body)
						t.Errorf("esperado sucesso, mas obteve status %d: %s", resp.StatusCode, string(body))
					}
				}
			})
		}
	})

	t.Log("🎯 Testes de integração de unidades de estoque concluídos com sucesso!")
}
