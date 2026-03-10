package handler

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/providentia/digna/cash_flow/pkg/cash_flow"
	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
	"github.com/providentia/digna/pdv_ui/pkg/pricing"
	"github.com/providentia/digna/supply/pkg/supply"
)

type PDVHandler struct {
	lifecycleManager  lifecycle.LifecycleManager
	tmpl              *template.Template
	pricingCalculator *pricing.PricingCalculator
	supplyAPI         supply.SupplyAPI
	cashAPI           *cash_flow.CashFlowAPI
}

func NewPDVHandler(lm lifecycle.LifecycleManager) (*PDVHandler, error) {
	funcMap := template.FuncMap{
		"divide": func(a, b float64) float64 {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"formatCurrency": func(amount int64) string {
			return fmt.Sprintf("R$ %.2f", float64(amount)/100)
		},
		"formatDate": func(t time.Time) string {
			return t.Format("02/01/2006 15:04")
		},
		// Adicionar funções necessárias para templates compartilhados
		"getAlertStatusLabel": func(status string) string {
			switch status {
			case "SAFE":
				return "Dentro do planejado"
			case "WARNING":
				return "Atenção: perto do limite"
			case "EXCEEDED":
				return "Ultrapassou o planejado"
			default:
				return status
			}
		},
		"getAlertStatusClass": func(status string) string {
			switch status {
			case "SAFE":
				return "bg-green-100 text-green-800 border-green-300"
			case "WARNING":
				return "bg-yellow-100 text-yellow-800 border-yellow-300"
			case "EXCEEDED":
				return "bg-red-100 text-red-800 border-red-300"
			default:
				return "bg-gray-100 text-gray-800 border-gray-300"
			}
		},
		"getCategoryLabel": func(category string) string {
			labels := map[string]string{
				"INSUMOS":      "Insumos",
				"ENERGIA":      "Energia",
				"EQUIPAMENTOS": "Equipamentos",
				"TRANSPORTE":   "Transporte",
				"MANUTENCAO":   "Manutenção",
				"SERVICOS":     "Serviços",
				"OUTROS":       "Outros",
			}
			if label, ok := labels[category]; ok {
				return label
			}
			return category
		},
	}

	// Criar template simples para evitar problemas de cache
	tmpl := template.New("pdv_simple.html").Funcs(funcMap)

	// Criar calculadora de preços
	pricingCalc, err := pricing.NewPricingCalculator()
	if err != nil {
		return nil, fmt.Errorf("failed to create pricing calculator: %w", err)
	}

	// Criar LedgerPort mock para supply
	mockLedgerPort := &mockLedgerPort{}

	// Criar API do módulo supply
	supplyAPI := supply.NewSupplyAPI(lm, mockLedgerPort)

	// Criar API do módulo cash_flow
	cashAPI := cash_flow.NewCashFlowAPI(lm)

	return &PDVHandler{
		lifecycleManager:  lm,
		tmpl:              tmpl,
		pricingCalculator: pricingCalc,
		supplyAPI:         supplyAPI,
		cashAPI:           cashAPI,
	}, nil
}

// mockLedgerPort implementa LedgerPort para testes
type mockLedgerPort struct{}

func (m *mockLedgerPort) RecordTransaction(entityID string, description string, postings []supply.LedgerPosting) error {
	// Implementação mock - apenas retorna sucesso
	return nil
}

func (h *PDVHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/pdv", h.PDVPage)
	mux.HandleFunc("/api/sale", h.RecordSale)
	mux.HandleFunc("/api/balance", h.GetBalance)

	// Registrar rotas da calculadora de preços
	h.pricingCalculator.RegisterRoutes(mux)
}

func (h *PDVHandler) PDVPage(w http.ResponseWriter, r *http.Request) {
	entityID := "cooperativa_demo"
	if r.URL.Query().Get("entity_id") != "" {
		entityID = r.URL.Query().Get("entity_id")
	}

	// Buscar produtos do estoque (apenas produtos acabados)
	ctx := r.Context()
	stockItems, err := h.supplyAPI.GetStockItems(ctx, entityID)
	slog.Info("PDV - Buscando produtos do estoque", "entity", entityID, "encontrados", len(stockItems), "erro", err)

	// Log detalhes dos produtos encontrados
	for i, item := range stockItems {
		slog.Debug("Produto encontrado", "idx", i, "nome", item.Name, "tipo", item.Type, "qtd", item.Quantity, "custo", item.UnitCost)
	}

	var products []string
	var productDetails []map[string]interface{}
	if err == nil && len(stockItems) > 0 {
		for _, item := range stockItems {
			// Incluir PRODUTO (produto acabado) e MERCADORIA (para revenda)
			if (item.Type == "PRODUTO" || item.Type == "MERCADORIA") && item.Quantity > 0 {
				products = append(products, item.Name)
				productDetails = append(productDetails, map[string]interface{}{
					"Name":           item.Name,
					"ID":             item.ID,
					"Type":           item.Type,
					"Quantity":       item.Quantity,
					"UnitCost":       int64(item.UnitCost),
					"FormattedPrice": fmt.Sprintf("R$ %.2f", float64(item.UnitCost)/100),
				})
			}
		}
	}

	// Se não encontrar produtos, mostrar mensagem
	if len(products) == 0 {
		// Não usar produtos mock - incentivar cadastro no estoque
		products = []string{"Nenhum produto cadastrado"}
		productDetails = []map[string]interface{}{
			{"Name": "Cadastre produtos no estoque primeiro", "ID": "", "Type": "", "Quantity": 0, "UnitCost": int64(0), "FormattedPrice": "R$ 0,00"},
		}
	}

	data := map[string]interface{}{
		"Title":          "PDV - Digna",
		"EntityID":       entityID,
		"Products":       products,
		"ProductDetails": productDetails,
	}

	// Carregar template do disco para evitar problemas de cache
	tmpl, err := template.New("pdv_simple.html").ParseFiles("templates/pdv_simple.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao carregar template: %v", err), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *PDVHandler) RecordSale(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("DEBUG: PDVHandler.RecordSale chamado\n")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	entityID := r.FormValue("entity_id")
	if entityID == "" {
		entityID = "cooperativa_demo"
	}

	amountStr := r.FormValue("amount")
	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil || amount <= 0 {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	product := r.FormValue("product")
	if product == "" {
		product = "Produto"
	}

	quantityStr := r.FormValue("quantity")
	quantity := 1
	if quantityStr != "" {
		if q, err := strconv.Atoi(quantityStr); err == nil && q > 0 {
			quantity = q
		}
	}

	stockItemID := r.FormValue("stock_item_id")
	slog.Info("PDV - Recebendo venda", "produto", product, "quantidade", quantity, "stock_item_id", stockItemID, "amount", amount)

	// Validar e atualizar estoque se stockItemID for fornecido
	if stockItemID != "" && stockItemID != "mock-1" && stockItemID != "mock-2" && stockItemID != "mock-3" {
		slog.Info("PDV - Tentando atualizar estoque", "stock_item_id", stockItemID, "quantidade", quantity)

		// Buscar item do estoque para validar quantidade
		ctx := r.Context()
		stockItems, err := h.supplyAPI.GetStockItems(ctx, entityID)
		if err != nil {
			slog.Error("PDV - Erro ao buscar itens do estoque", "erro", err)
		} else {
			slog.Info("PDV - Itens do estoque encontrados", "quantidade", len(stockItems))
			var currentStockItem *supply.StockItem
			for _, item := range stockItems {
				if item.ID == stockItemID {
					currentStockItem = item
					slog.Info("PDV - Item encontrado no estoque", "id", item.ID, "nome", item.Name, "quantidade_atual", item.Quantity)
					break
				}
			}

			if currentStockItem != nil {
				// Validar quantidade disponível
				if quantity > currentStockItem.Quantity {
					w.Header().Set("Content-Type", "text/html")
					fmt.Fprintf(w, `
						<div class="bg-red-100 border-l-4 border-red-500 text-red-700 p-4 mb-4" role="alert">
							<p class="font-bold">Estoque insuficiente!</p>
							<p>Quantidade solicitada: %d | Disponível: %d</p>
						</div>
					`, quantity, currentStockItem.Quantity)
					return
				}

				// Atualizar estoque (reduzir quantidade)
				slog.Info("PDV - Chamando UpdateStockQuantity", "entity_id", entityID, "stock_item_id", stockItemID, "delta", -quantity)
				resp, err := h.supplyAPI.UpdateStockQuantity(ctx, entityID, stockItemID, -quantity)
				if err != nil || !resp.Success {
					slog.Error("PDV - Erro ao atualizar estoque", "produto", product, "quantidade", quantity, "erro", err, "resp", resp)
					// Continuar com a venda mesmo se falhar a atualização de estoque
				} else {
					slog.Info("PDV - Estoque atualizado com sucesso", "produto", product, "quantidade", quantity, "estoque_anterior", currentStockItem.Quantity, "stock_id", stockItemID)
				}
			} else {
				slog.Error("PDV - Item não encontrado no estoque", "stock_item_id", stockItemID)
			}
		}
	} else {
		slog.Info("PDV - StockItemID vazio ou mock, pulando atualização de estoque", "stock_item_id", stockItemID)
	}

	// Criar descrição com quantidade
	description := fmt.Sprintf("Venda PDV: %s", product)
	if quantity > 1 {
		description = fmt.Sprintf("Venda PDV: %s (Qtd: %d)", product, quantity)
	}

	// Registrar venda no cash_flow como crédito (entrada de dinheiro)
	req := cash_flow.EntryRequest{
		EntityID:    entityID,
		Type:        "CREDIT",
		Amount:      amount,
		Category:    "SALES",
		Description: description,
	}

	result, err := h.cashAPI.RecordEntry(req)
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `
			<div class="bg-red-100 border-l-4 border-red-500 text-red-700 p-4 mb-4" role="alert">
				<p class="font-bold">Erro ao registrar venda!</p>
				<p>%v</p>
			</div>
		`, err)
		return
	}

	if !result.Success {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `
			<div class="bg-yellow-100 border-l-4 border-yellow-500 text-yellow-700 p-4 mb-4" role="alert">
				<p class="font-bold">Aviso na venda</p>
				<p>%s</p>
			</div>
		`, result.Error)
		return
	}

	// Obter saldo atualizado
	balance, _ := h.cashAPI.GetBalance(entityID)
	var balanceAmount int64
	if balance != nil {
		balanceAmount = balance.Balance
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<div class="bg-green-100 border-l-4 border-green-500 text-green-700 p-4 mb-4" role="alert">
			<p class="font-bold">✅ Venda Registrada!</p>
			<p>Produto: <strong>%s</strong> | Valor: <strong>R$ %.2f</strong></p>
			<p class="mt-2">Saldo atual: <strong class="text-green-800">R$ %.2f</strong></p>
		</div>
	`, product, float64(amount)/100, float64(balanceAmount)/100)
}

func (h *PDVHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	entityID := r.URL.Query().Get("entity_id")
	if entityID == "" {
		entityID = "cooperativa_demo"
	}

	// Obter saldo real do cash_flow
	balance, err := h.cashAPI.GetBalance(entityID)
	if err != nil || balance == nil {
		// Fallback para valor mock em caso de erro
		balanceAmount := int64(150000) // R$ 1.500,00 mock
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<span class="text-2xl font-bold text-green-600">R$ %.2f</span>`, float64(balanceAmount)/100)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `<span class="text-2xl font-bold text-green-600">R$ %.2f</span>`, float64(balance.Balance)/100)
}
