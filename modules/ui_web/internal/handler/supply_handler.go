package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
	"github.com/providentia/digna/supply/pkg/supply"
)

type SupplyHandler struct {
	lifecycleManager lifecycle.LifecycleManager
	tmpl             *template.Template
	supplyAPI        supply.SupplyAPI
}

func NewSupplyHandler(lm lifecycle.LifecycleManager) (*SupplyHandler, error) {
	// Criar template com funções auxiliares
	funcMap := template.FuncMap{
		"formatCurrency": func(amount int64) string {
			return fmt.Sprintf("R$ %.2f", float64(amount)/100)
		},
		"formatDate": func(t time.Time) string {
			return t.Format("02/01/2006")
		},
		"stockItemTypeLabel": func(itemType string) string {
			switch itemType {
			case "INSUMO":
				return "Insumo/Matéria-prima"
			case "PRODUTO":
				return "Produto Acabado"
			case "MERCADORIA":
				return "Mercadoria para Revenda"
			default:
				return itemType
			}
		},
		"divide": func(a, b int) float64 {
			if b == 0 {
				return 0
			}
			return float64(a) / float64(b)
		},
		"isBelowMinimum": func(quantity, minQuantity int) bool {
			return quantity < minQuantity
		},
	}

	// Carregar templates
	tmpl, err := template.New("supply_templates").Funcs(funcMap).ParseGlob("templates/supply_*.html")
	if err != nil {
		// Fallback para templates embutidos se os arquivos não existirem
		tmpl = template.Must(template.New("supply_templates").Funcs(funcMap).Parse(supplyTemplates))
	}

	// Criar API de supply (sem ledgerPort por enquanto - será injetado depois se necessário)
	supplyAPI := supply.NewSupplyAPI(lm, nil)

	return &SupplyHandler{
		lifecycleManager: lm,
		tmpl:             tmpl,
		supplyAPI:        supplyAPI,
	}, nil
}

func (h *SupplyHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/supply", h.SupplyDashboard)
	mux.HandleFunc("/supply/purchase", h.PurchasePage)
	mux.HandleFunc("/supply/suppliers", h.SuppliersPage)
	mux.HandleFunc("/supply/stock", h.StockPage)

	// API endpoints
	mux.HandleFunc("/api/supply/purchase", h.RegisterPurchase)
	mux.HandleFunc("/api/supply/supplier", h.RegisterSupplier)
	mux.HandleFunc("/api/supply/stock-item", h.RegisterStockItem)
}

func (h *SupplyHandler) SupplyDashboard(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":    "Gestão de Compras e Estoque - Digna",
		"EntityID": "cooperativa_demo",
	}

	if err := h.tmpl.ExecuteTemplate(w, "supply_dashboard.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *SupplyHandler) PurchasePage(w http.ResponseWriter, r *http.Request) {
	entityID := "cooperativa_demo"
	ctx := r.Context()

	// Buscar fornecedores e itens de estoque para os selects
	suppliers, _ := h.supplyAPI.GetSuppliers(ctx, entityID)
	stockItems, _ := h.supplyAPI.GetStockItems(ctx, entityID)

	data := map[string]interface{}{
		"Title":      "Nova Compra de Material - Digna",
		"EntityID":   entityID,
		"Suppliers":  suppliers,
		"StockItems": stockItems,
	}

	if err := h.tmpl.ExecuteTemplate(w, "supply_purchase.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *SupplyHandler) SuppliersPage(w http.ResponseWriter, r *http.Request) {
	entityID := "cooperativa_demo"
	ctx := r.Context()

	suppliers, err := h.supplyAPI.GetSuppliers(ctx, entityID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar fornecedores: %v", err), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":     "Fornecedores - Digna",
		"EntityID":  entityID,
		"Suppliers": suppliers,
	}

	if err := h.tmpl.ExecuteTemplate(w, "supply_suppliers.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *SupplyHandler) StockPage(w http.ResponseWriter, r *http.Request) {
	entityID := "cooperativa_demo"
	ctx := r.Context()

	stockItems, err := h.supplyAPI.GetStockItems(ctx, entityID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar estoque: %v", err), http.StatusInternalServerError)
		return
	}

	// Gerar relatório de estoque
	report, _ := h.supplyAPI.GetStockReport(ctx, entityID)

	data := map[string]interface{}{
		"Title":       "Meu Estoque - Digna",
		"EntityID":    entityID,
		"StockItems":  stockItems,
		"StockReport": report,
		"ItemTypes": []string{
			"INSUMO",
			"PRODUTO",
			"MERCADORIA",
		},
	}

	if err := h.tmpl.ExecuteTemplate(w, "supply_stock.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *SupplyHandler) RegisterPurchase(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	entityID := "cooperativa_demo"
	ctx := r.Context()

	// Parse form
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao processar formulário: %v", err), http.StatusBadRequest)
		return
	}

	supplierID := r.FormValue("supplier_id")
	paymentType := r.FormValue("payment_type")

	// Parse items (simplificado - um item por compra nesta versão)
	stockItemID := r.FormValue("stock_item_id")
	quantityStr := r.FormValue("quantity")
	unitCostStr := r.FormValue("unit_cost")

	quantity, err := strconv.Atoi(quantityStr)
	if err != nil || quantity <= 0 {
		http.Error(w, "Quantidade inválida", http.StatusBadRequest)
		return
	}

	unitCost, err := strconv.ParseInt(unitCostStr, 10, 64)
	if err != nil || unitCost <= 0 {
		http.Error(w, "Valor unitário inválido", http.StatusBadRequest)
		return
	}

	// Criar requisição de compra
	req := supply.PurchaseRequest{
		EntityID:    entityID,
		SupplierID:  supplierID,
		PaymentType: paymentType,
		Items: []supply.PurchaseItemRequest{
			{
				StockItemID: stockItemID,
				Quantity:    quantity,
				UnitCost:    unitCost,
			},
		},
	}

	// Registrar compra
	resp, err := h.supplyAPI.RegisterPurchase(ctx, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao registrar compra: %v", err), http.StatusInternalServerError)
		return
	}

	if !resp.Success {
		http.Error(w, fmt.Sprintf("Falha ao registrar compra: %s", resp.Error), http.StatusBadRequest)
		return
	}

	// Retornar sucesso
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"success": true, "purchase_id": "%s", "message": "Compra registrada com sucesso!"}`, resp.PurchaseID)
}

func (h *SupplyHandler) RegisterSupplier(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	entityID := "cooperativa_demo"
	ctx := r.Context()

	// Parse form
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao processar formulário: %v", err), http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	contactInfo := r.FormValue("contact_info")

	// Criar requisição de fornecedor
	req := supply.SupplierRequest{
		EntityID:    entityID,
		Name:        name,
		ContactInfo: contactInfo,
	}

	// Registrar fornecedor
	resp, err := h.supplyAPI.RegisterSupplier(ctx, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao registrar fornecedor: %v", err), http.StatusInternalServerError)
		return
	}

	if !resp.Success {
		http.Error(w, fmt.Sprintf("Falha ao registrar fornecedor: %s", resp.Error), http.StatusBadRequest)
		return
	}

	// Retornar sucesso
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"success": true, "supplier_id": "%s", "message": "Fornecedor registrado com sucesso!"}`, resp.SupplierID)
}

func (h *SupplyHandler) RegisterStockItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	entityID := "cooperativa_demo"
	ctx := r.Context()

	// Parse form
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao processar formulário: %v", err), http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	itemType := r.FormValue("type")
	quantityStr := r.FormValue("quantity")
	minQuantityStr := r.FormValue("min_quantity")
	unitCostStr := r.FormValue("unit_cost")

	quantity, _ := strconv.Atoi(quantityStr)
	minQuantity, _ := strconv.Atoi(minQuantityStr)
	unitCost, err := strconv.ParseInt(unitCostStr, 10, 64)
	if err != nil || unitCost <= 0 {
		http.Error(w, "Valor unitário inválido", http.StatusBadRequest)
		return
	}

	// Converter string para tipo de item
	var itemTypeEnum string
	switch itemType {
	case "INSUMO", "PRODUTO", "MERCADORIA":
		itemTypeEnum = itemType
	default:
		itemTypeEnum = "INSUMO" // default
	}

	// Criar requisição de item de estoque
	req := supply.StockItemRequest{
		EntityID:    entityID,
		Name:        name,
		Type:        itemTypeEnum,
		Quantity:    quantity,
		MinQuantity: minQuantity,
		UnitCost:    unitCost,
	}

	// Registrar item de estoque
	resp, err := h.supplyAPI.RegisterStockItem(ctx, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao registrar item de estoque: %v", err), http.StatusInternalServerError)
		return
	}

	if !resp.Success {
		http.Error(w, fmt.Sprintf("Falha ao registrar item de estoque: %s", resp.Error), http.StatusBadRequest)
		return
	}

	// Retornar sucesso
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"success": true, "stock_item_id": "%s", "message": "Item de estoque registrado com sucesso!"}`, resp.StockItemID)
}
