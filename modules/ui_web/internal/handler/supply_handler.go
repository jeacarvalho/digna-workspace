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
		"stockItemUnitLabel": func(unit string) string {
			switch unit {
			case "UNIDADE":
				return "unid."
			case "KG":
				return "kg"
			case "G":
				return "g"
			case "L":
				return "L"
			case "M":
				return "m"
			case "CM":
				return "cm"
			case "PACOTE":
				return "pct"
			case "CAIXA":
				return "cx"
			case "SACO":
				return "sc"
			default:
				return unit
			}
		},
		"multiply": func(a int64, b int) int64 {
			return a * int64(b)
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
		"fdiv": func(a, b float64) float64 {
			if b == 0 {
				return 0
			}
			return a / b
		},
	}

	// Carregar templates: primeiro os embutidos (base), depois templates do disco
	baseTmpl := template.New("supply_base").Funcs(funcMap)

	// 1. Parse templates embutidos (obrigatórios)
	baseTmpl, err := baseTmpl.Parse(supplyTemplates)
	if err != nil {
		return nil, fmt.Errorf("failed to parse embedded supply templates: %w", err)
	}

	// 2. Tentar adicionar templates do disco (opcional)
	// Clone o template base para não modificar o original
	tmpl, err := baseTmpl.Clone()
	if err != nil {
		return nil, fmt.Errorf("failed to clone base template: %w", err)
	}

	// Adicionar templates do disco (se existirem)
	diskTemplates, err := tmpl.ParseGlob("templates/supply_*.html")
	if err == nil {
		// Templates do disco adicionados com sucesso
		tmpl = diskTemplates
		fmt.Printf("✅ Loaded supply templates: embedded + disk templates\n")

		// DEBUG: Listar templates carregados
		for _, t := range tmpl.Templates() {
			if t.Name() != "" {
				fmt.Printf("  - Template: %s\n", t.Name())
			}
		}
	} else {
		// Usar apenas templates embutidos
		fmt.Printf("ℹ️ Using embedded supply templates only (no disk templates found: %v)\n", err)

		// DEBUG: Listar templates embutidos
		for _, t := range baseTmpl.Templates() {
			if t.Name() != "" {
				fmt.Printf("  - Embedded template: %s\n", t.Name())
			}
		}
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
	mux.HandleFunc("/api/supply/stock-items", h.GetStockItemsAPI)
}

func (h *SupplyHandler) SupplyDashboard(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":    "Gestão de Compras e Estoque - Digna",
		"EntityID": "cooperativa_demo",
	}

	// Usar template do handler (que já tem todas as funções incluindo fdiv)
	if err := h.tmpl.ExecuteTemplate(w, "supply_dashboard_simple.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	// Usar template do handler (que já tem todas as funções incluindo fdiv)
	if err := h.tmpl.ExecuteTemplate(w, "supply_stock_simple.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	// Retornar sucesso com HTML amigável
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `
		<div class="bg-green-100 border-l-4 border-green-500 text-green-700 p-4 mb-4 rounded" role="alert">
			<div class="flex">
				<div class="flex-shrink-0">
					<svg class="h-5 w-5 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
					</svg>
				</div>
				<div class="ml-3">
					<p class="font-bold">✅ Compra registrada com sucesso!</p>
					<p class="text-sm mt-1">ID da compra: <span class="font-mono text-green-800">%s</span></p>
					<p class="text-sm mt-1">O estoque foi atualizado automaticamente.</p>
				</div>
			</div>
		</div>
	`, resp.PurchaseID)
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

	// Retornar sucesso com HTML amigável
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `
		<div class="bg-green-100 border-l-4 border-green-500 text-green-700 p-4 mb-4 rounded" role="alert">
			<div class="flex">
				<div class="flex-shrink-0">
					<svg class="h-5 w-5 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
					</svg>
				</div>
				<div class="ml-3">
					<p class="font-bold">✅ Fornecedor registrado com sucesso!</p>
					<p class="text-sm mt-1">Nome: <span class="font-semibold text-green-800">%s</span></p>
					<p class="text-sm mt-1">ID: <span class="font-mono text-green-800">%s</span></p>
				</div>
			</div>
		</div>
	`, name, resp.SupplierID)
}

func (h *SupplyHandler) GetStockItemsAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	entityID := "cooperativa_demo"
	ctx := r.Context()

	stockItems, err := h.supplyAPI.GetStockItems(ctx, entityID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar itens de estoque: %v", err), http.StatusInternalServerError)
		return
	}

	// Renderizar apenas a lista de itens
	w.Header().Set("Content-Type", "text/html")

	if len(stockItems) == 0 {
		fmt.Fprintf(w, `
			<div class="text-center py-8 text-gray-500">
				<svg class="w-12 h-12 mx-auto text-gray-300 mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"></path>
				</svg>
				<p>Nenhum item cadastrado ainda.</p>
				<p class="text-sm mt-1">Cadastre seu primeiro item usando o formulário ao lado.</p>
			</div>
		`)
		return
	}

	fmt.Fprintf(w, `<div class="space-y-4">`)
	for _, item := range stockItems {
		// Função auxiliar para traduzir tipo de item
		typeLabel := item.Type
		switch item.Type {
		case "INSUMO":
			typeLabel = "Insumo/Matéria-prima"
		case "PRODUTO":
			typeLabel = "Produto Acabado"
		case "MERCADORIA":
			typeLabel = "Mercadoria para Revenda"
		}

		// Função auxiliar para unidade
		unitLabel := item.Unit
		switch item.Unit {
		case "UNIDADE":
			unitLabel = "unid."
		case "KG":
			unitLabel = "kg"
		case "G":
			unitLabel = "g"
		case "L":
			unitLabel = "L"
		case "M":
			unitLabel = "m"
		case "CM":
			unitLabel = "cm"
		case "PACOTE":
			unitLabel = "pct"
		case "CAIXA":
			unitLabel = "cx"
		case "SACO":
			unitLabel = "sc"
		}

		// Verificar se está abaixo do mínimo
		isBelowMinimum := item.Quantity < item.MinQuantity

		fmt.Fprintf(w, `
			<div class="border border-gray-200 rounded-xl p-4 hover:bg-gray-50 transition %s">
				<div class="flex justify-between items-start">
					<div>
						<div class="flex items-center space-x-2">
							<h3 class="font-semibold text-gray-800">%s</h3>
							<span class="text-xs px-2 py-1 rounded-full %s">
								%s
							</span>
						</div>
						<div class="mt-2 grid grid-cols-4 gap-3 text-sm">
							<div>
								<span class="text-gray-500">Qtde:</span>
								<span class="font-semibold ml-1 %s">
									%d %s
								</span>
							</div>
							<div>
								<span class="text-gray-500">Mín:</span>
								<span class="font-semibold ml-1">%d %s</span>
							</div>
							<div>
								<span class="text-gray-500">Custo:</span>
								<span class="font-semibold ml-1">R$ %.2f/%s</span>
							</div>
							<div>
								<span class="text-gray-500">Total:</span>
								<span class="font-semibold ml-1">R$ %.2f</span>
							</div>
						</div>
						%s
					</div>
					<span class="text-xs text-gray-500">%s</span>
				</div>
			</div>
		`,
			func() string {
				if isBelowMinimum {
					return "border-red-200 bg-red-50"
				}
				return ""
			}(),
			item.Name,
			func() string {
				switch item.Type {
				case "INSUMO":
					return "bg-blue-100 text-blue-800"
				case "PRODUTO":
					return "bg-green-100 text-green-800"
				default:
					return "bg-purple-100 text-purple-800"
				}
			}(),
			typeLabel,
			func() string {
				if isBelowMinimum {
					return "text-red-600"
				}
				return "text-gray-800"
			}(),
			item.Quantity, unitLabel,
			item.MinQuantity, unitLabel,
			float64(item.UnitCost)/100, unitLabel,
			float64(item.UnitCost*int64(item.Quantity))/100,
			func() string {
				if isBelowMinimum {
					return `<div class="mt-2 text-xs text-red-600 font-medium">⚠️ Estoque abaixo do mínimo</div>`
				}
				return ""
			}(),
			item.CreatedAt.Format("02/01/2006"),
		)
	}
	fmt.Fprintf(w, `</div>`)
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
	unit := r.FormValue("unit")
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

	// Validar unidade
	if unit == "" {
		unit = "UNIDADE" // default
	}

	// Criar requisição de item de estoque
	req := supply.StockItemRequest{
		EntityID:    entityID,
		Name:        name,
		Type:        itemTypeEnum,
		Unit:        unit,
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

	// Retornar sucesso com HTML amigável
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	// Função auxiliar para traduzir tipo de item
	typeLabel := itemTypeEnum
	switch itemTypeEnum {
	case "INSUMO":
		typeLabel = "Insumo/Matéria-prima"
	case "PRODUTO":
		typeLabel = "Produto Acabado"
	case "MERCADORIA":
		typeLabel = "Mercadoria para Revenda"
	}

	// Função auxiliar para unidade
	unitLabel := unit
	switch unit {
	case "UNIDADE":
		unitLabel = "unidades"
	case "KG":
		unitLabel = "kg"
	case "G":
		unitLabel = "g"
	case "L":
		unitLabel = "litros"
	case "M":
		unitLabel = "metros"
	case "CM":
		unitLabel = "cm"
	case "PACOTE":
		unitLabel = "pacotes"
	case "CAIXA":
		unitLabel = "caixas"
	case "SACO":
		unitLabel = "sacos"
	}

	fmt.Fprintf(w, `
		<div class="bg-green-100 border-l-4 border-green-500 text-green-700 p-4 mb-4 rounded" role="alert">
			<div class="flex">
				<div class="flex-shrink-0">
					<svg class="h-5 w-5 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
					</svg>
				</div>
				<div class="ml-3">
					<p class="font-bold">✅ Item de estoque registrado com sucesso!</p>
					<p class="text-sm mt-1">Nome: <span class="font-semibold text-green-800">%s</span></p>
					<p class="text-sm mt-1">Tipo: <span class="font-semibold text-green-800">%s</span></p>
					<p class="text-sm mt-1">Unidade: <span class="font-semibold text-green-800">%s</span></p>
					<p class="text-sm mt-1">Quantidade: <span class="font-semibold text-green-800">%d %s</span></p>
					<p class="text-sm mt-1">Custo unitário: <span class="font-semibold text-green-800">R$ %.2f/%s</span></p>
					<p class="text-sm mt-1">ID: <span class="font-mono text-green-800">%s</span></p>
				</div>
			</div>
		</div>
	`, name, typeLabel, unitLabel, quantity, unitLabel, float64(unitCost)/100, unitLabel, resp.StockItemID)
}
