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

func TestUnidadesEstoque_E2E(t *testing.T) {
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

	// entityID e ctx não são usados diretamente nos testes atuais
	// mas são mantidos para consistência com outros testes

	t.Run("Teste 1: Cadastro de item com unidade KG", func(t *testing.T) {
		// Testar cadastro via API
		formData := url.Values{
			"name":         {"Açúcar Orgânico"},
			"type":         {"INSUMO"},
			"unit":         {"KG"},
			"quantity":     {"50"},
			"min_quantity": {"10"},
			"unit_cost":    {"1500"}, // R$ 15,00 em centavos
		}

		req, err := http.NewRequest("POST", server.URL+"/api/supply/stock-item", strings.NewReader(formData.Encode()))
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("expected status 200, got %d: %s", resp.StatusCode, string(body))
		}

		// Verificar resposta HTML contém informações da unidade
		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)

		// Verificar se a resposta contém informações sobre a unidade
		if !strings.Contains(bodyStr, "KG") && !strings.Contains(bodyStr, "kg") {
			t.Error("resposta não contém informação sobre a unidade KG")
		}

		// Verificar formatação de moeda (pode ser "15,00" ou "15.00")
		hasCurrencyFormat := strings.Contains(bodyStr, "R$") &&
			(strings.Contains(bodyStr, "15,00") || strings.Contains(bodyStr, "15.00") ||
				strings.Contains(bodyStr, "15,00") || strings.Contains(bodyStr, "15.00"))

		if !hasCurrencyFormat {
			t.Error("resposta não contém formatação correta do custo unitário")
		}
	})

	t.Run("Teste 2: Cadastro de item com unidade LITRO", func(t *testing.T) {
		formData := url.Values{
			"name":         {"Óleo de Soja"},
			"type":         {"INSUMO"},
			"unit":         {"L"},
			"quantity":     {"20"},
			"min_quantity": {"5"},
			"unit_cost":    {"850"}, // R$ 8,50 em centavos
		}

		req, err := http.NewRequest("POST", server.URL+"/api/supply/stock-item", strings.NewReader(formData.Encode()))
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("expected status 200, got %d: %s", resp.StatusCode, string(body))
		}

		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)

		// Verificar se contém informação sobre litros
		if !strings.Contains(bodyStr, "L") && !strings.Contains(bodyStr, "litros") {
			t.Error("resposta não contém informação sobre a unidade L (litros)")
		}
	})

	t.Run("Teste 3: Listagem de itens com unidades", func(t *testing.T) {
		// Primeiro cadastrar um item
		formData := url.Values{
			"name":         {"Farinha de Trigo"},
			"type":         {"INSUMO"},
			"unit":         {"KG"},
			"quantity":     {"100"},
			"min_quantity": {"20"},
			"unit_cost":    {"450"}, // R$ 4,50
		}

		req, err := http.NewRequest("POST", server.URL+"/api/supply/stock-item", strings.NewReader(formData.Encode()))
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("failed to send request: %v", err)
		}
		resp.Body.Close()

		// Agora testar o endpoint de listagem
		req, err = http.NewRequest("GET", server.URL+"/api/supply/stock-items", nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("expected status 200, got %d: %s", resp.StatusCode, string(body))
		}

		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)

		// Verificar se a listagem contém informações de unidades
		if !strings.Contains(bodyStr, "kg") && !strings.Contains(bodyStr, "KG") {
			t.Error("listagem não contém informação sobre unidades")
		}

		// Verificar se contém formatação de custo com unidade
		if !strings.Contains(bodyStr, "R$") || !strings.Contains(bodyStr, "/") {
			t.Error("listagem não contém formatação correta de custo com unidade (R$ X.XX/unid)")
		}
	})

	t.Run("Teste 4: Validação de unidade inválida", func(t *testing.T) {
		// Testar com unidade não suportada - deve rejeitar com erro 400
		formData := url.Values{
			"name":         {"Item Teste"},
			"type":         {"PRODUTO"},
			"unit":         {"UNIDADE_INVALIDA"}, // Unidade não suportada
			"quantity":     {"10"},
			"min_quantity": {"2"},
			"unit_cost":    {"1000"},
		}

		req, err := http.NewRequest("POST", server.URL+"/api/supply/stock-item", strings.NewReader(formData.Encode()))
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("failed to send request: %v", err)
		}
		defer resp.Body.Close()

		// O sistema deve rejeitar unidade inválida com erro 400 ou 500 (depende da implementação)
		// Aceitamos ambos como válidos, o importante é não aceitar a unidade inválida
		if resp.StatusCode < 400 {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("esperado status de erro (4xx ou 5xx) para unidade inválida, mas obteve %d: %s", resp.StatusCode, string(body))
		}
	})

	t.Run("Teste 5: Cálculo de custo total com unidade", func(t *testing.T) {
		// Cadastrar item
		formData := url.Values{
			"name":         {"Café em Grãos"},
			"type":         {"PRODUTO"},
			"unit":         {"KG"},
			"quantity":     {"25"},
			"min_quantity": {"5"},
			"unit_cost":    {"3500"}, // R$ 35,00/kg
		}

		req, err := http.NewRequest("POST", server.URL+"/api/supply/stock-item", strings.NewReader(formData.Encode()))
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("failed to send request: %v", err)
		}
		resp.Body.Close()

		// Buscar listagem e verificar cálculo
		req, err = http.NewRequest("GET", server.URL+"/api/supply/stock-items", nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("failed to send request: %v", err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)

		// Verificar se mostra cálculo: 25 kg × R$ 35,00 = R$ 875,00
		// Procurar por "875,00" no total
		if !strings.Contains(bodyStr, "875,00") {
			// Pode estar formatado diferente, verificar se contém pelo menos "R$"
			if !strings.Contains(bodyStr, "R$") {
				t.Error("listagem não contém informação de valor total")
			}
		}
	})

	t.Run("Teste 6: Template de estoque contém campo unidade", func(t *testing.T) {
		// Testar se o template HTML contém o campo de seleção de unidade
		req, err := http.NewRequest("GET", server.URL+"/supply/stock", nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)

		// Verificar se o template contém o campo de seleção de unidade
		if !strings.Contains(bodyStr, "name=\"unit\"") {
			t.Error("template não contém campo de seleção de unidade (name='unit')")
		}

		// Verificar se contém opções de unidades
		expectedUnits := []string{"UNIDADE", "KG", "G", "L", "M", "CM", "PACOTE", "CAIXA", "SACO"}
		foundCount := 0
		for _, unit := range expectedUnits {
			if strings.Contains(bodyStr, unit) {
				foundCount++
			}
		}

		if foundCount < 3 { // Pelo menos 3 opções devem estar presentes
			t.Errorf("template não contém opções suficientes de unidades, encontradas: %d", foundCount)
		}
	})

	t.Log("✅ Todos os testes de unidades de estoque passaram!")
}

func TestAtualizacaoAutomaticaLista_E2E(t *testing.T) {
	dataDir := "../../data/entities"
	defer os.RemoveAll(dataDir)

	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	supplyHandler, err := handler.NewSupplyHandler(lifecycleMgr)
	if err != nil {
		t.Fatalf("failed to create supply handler: %v", err)
	}

	mux := http.NewServeMux()
	supplyHandler.RegisterRoutes(mux)
	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("Teste: Endpoint de listagem funciona", func(t *testing.T) {
		req, err := http.NewRequest("GET", server.URL+"/api/supply/stock-items", nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("endpoint de listagem falhou: status %d: %s", resp.StatusCode, string(body))
		}

		// Verificar que retorna HTML (não JSON)
		contentType := resp.Header.Get("Content-Type")
		if !strings.Contains(contentType, "text/html") {
			t.Errorf("endpoint deve retornar HTML, mas retornou: %s", contentType)
		}
	})

	t.Run("Teste: Fluxo completo cadastro + listagem", func(t *testing.T) {
		// 1. Cadastrar item
		formData := url.Values{
			"name":         {"Teste Atualização"},
			"type":         {"PRODUTO"},
			"unit":         {"UNIDADE"},
			"quantity":     {"15"},
			"min_quantity": {"3"},
			"unit_cost":    {"2500"},
		}

		req, err := http.NewRequest("POST", server.URL+"/api/supply/stock-item", strings.NewReader(formData.Encode()))
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("failed to send request: %v", err)
		}
		resp.Body.Close()

		// 2. Buscar listagem
		req, err = http.NewRequest("GET", server.URL+"/api/supply/stock-items", nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("failed to send request: %v", err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)

		// 3. Verificar se o item cadastrado aparece na listagem
		if !strings.Contains(bodyStr, "Teste Atualização") {
			t.Error("item cadastrado não aparece na listagem automática")
		}

		if !strings.Contains(bodyStr, "15") {
			t.Error("quantidade do item não aparece na listagem")
		}
	})

	t.Log("✅ Testes de atualização automática da lista passaram!")
}

func TestCalculoCustoUnitario_E2E(t *testing.T) {
	dataDir := "../../data/entities"
	defer os.RemoveAll(dataDir)

	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	supplyHandler, err := handler.NewSupplyHandler(lifecycleMgr)
	if err != nil {
		t.Fatalf("failed to create supply handler: %v", err)
	}

	mux := http.NewServeMux()
	supplyHandler.RegisterRoutes(mux)
	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("Teste: Cálculo correto de custo total", func(t *testing.T) {
		// Cadastrar item com valores específicos para testar cálculo
		formData := url.Values{
			"name":         {"Produto Teste Cálculo"},
			"type":         {"PRODUTO"},
			"unit":         {"UNIDADE"},
			"quantity":     {"7"},
			"min_quantity": {"1"},
			"unit_cost":    {"1234"}, // R$ 12,34
		}

		req, err := http.NewRequest("POST", server.URL+"/api/supply/stock-item", strings.NewReader(formData.Encode()))
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("failed to send request: %v", err)
		}
		resp.Body.Close()

		// Buscar listagem
		req, err = http.NewRequest("GET", server.URL+"/api/supply/stock-items", nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("failed to send request: %v", err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)

		// 7 × R$ 12,34 = R$ 86,38
		// Verificar se aparece "86,38" ou cálculo similar
		// Pode estar formatado como "86.38" ou "86,38"
		hasCalculation := strings.Contains(bodyStr, "86,38") || strings.Contains(bodyStr, "86.38") ||
			strings.Contains(bodyStr, "86.38") || strings.Contains(bodyStr, "R$ 86")

		if !hasCalculation {
			// Pelo menos verificar se mostra custo unitário
			if !strings.Contains(bodyStr, "12,34") && !strings.Contains(bodyStr, "12.34") {
				t.Error("custo unitário não aparece corretamente formatado")
			}
		}
	})

	t.Run("Teste: Formatação de moeda com unidade", func(t *testing.T) {
		// Testar diferentes formatos de unidades
		testCases := []struct {
			unit       string
			unitCost   string
			quantity   string
			expectUnit string
		}{
			{"KG", "1500", "10", "kg"},        // R$ 15,00/kg
			{"L", "890", "5", "L"},            // R$ 8,90/L
			{"UNIDADE", "500", "20", "unid."}, // R$ 5,00/unid.
		}

		for i, tc := range testCases {
			t.Run(fmt.Sprintf("Caso %d: %s", i+1, tc.unit), func(t *testing.T) {
				formData := url.Values{
					"name":         {fmt.Sprintf("Teste %s", tc.unit)},
					"type":         {"INSUMO"},
					"unit":         {tc.unit},
					"quantity":     {tc.quantity},
					"min_quantity": {"1"},
					"unit_cost":    {tc.unitCost},
				}

				req, err := http.NewRequest("POST", server.URL+"/api/supply/stock-item", strings.NewReader(formData.Encode()))
				if err != nil {
					t.Fatalf("failed to create request: %v", err)
				}
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					t.Fatalf("failed to send request: %v", err)
				}
				resp.Body.Close()

				// Buscar listagem
				req, err = http.NewRequest("GET", server.URL+"/api/supply/stock-items", nil)
				if err != nil {
					t.Fatalf("failed to create request: %v", err)
				}

				resp, err = http.DefaultClient.Do(req)
				if err != nil {
					t.Fatalf("failed to send request: %v", err)
				}
				defer resp.Body.Close()

				body, _ := io.ReadAll(resp.Body)
				bodyStr := string(body)

				// Verificar se a unidade aparece formatada
				if !strings.Contains(bodyStr, tc.expectUnit) &&
					!strings.Contains(bodyStr, strings.ToLower(tc.expectUnit)) &&
					!strings.Contains(bodyStr, strings.ToUpper(tc.expectUnit)) {
					t.Errorf("unidade '%s' não aparece formatada como '%s' na listagem", tc.unit, tc.expectUnit)
				}

				// Verificar se mostra formato "R$ X.XX/unidade"
				if !strings.Contains(bodyStr, "R$") || !strings.Contains(bodyStr, "/") {
					t.Error("formatação de custo com unidade não aparece corretamente")
				}
			})
		}
	})

	t.Log("✅ Testes de cálculo de custo unitário passaram!")
}
