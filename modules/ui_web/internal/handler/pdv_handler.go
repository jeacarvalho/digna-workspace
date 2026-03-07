package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
	"github.com/providentia/digna/pdv_ui/usecase"
)

type PDVHandler struct {
	lifecycleManager lifecycle.LifecycleManager
	tmpl             *template.Template
}

func NewPDVHandler(lm lifecycle.LifecycleManager) (*PDVHandler, error) {
	// Create template with custom functions
	funcMap := template.FuncMap{
		"divide": func(a, b float64) float64 {
			if b == 0 {
				return 0
			}
			return a / b
		},
	}

	tmpl, err := template.New("templates").Funcs(funcMap).ParseGlob("templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	return &PDVHandler{
		lifecycleManager: lm,
		tmpl:             tmpl,
	}, nil
}

func (h *PDVHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/pdv", h.PDVPage)
	mux.HandleFunc("/api/sale", h.RecordSale)
	mux.HandleFunc("/api/balance", h.GetBalance)
}

func (h *PDVHandler) PDVPage(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":    "PDV - Digna",
		"EntityID": "cooperativa_demo",
		"Products": []string{"Mel", "Artesanato", "Serviços"},
	}

	if err := h.tmpl.ExecuteTemplate(w, "pdv.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *PDVHandler) RecordSale(w http.ResponseWriter, r *http.Request) {
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

	opHandler := usecase.NewOperationHandler(h.lifecycleManager)

	saleReq := usecase.SaleRequest{
		EntityID:      entityID,
		Amount:        amount,
		PaymentMethod: "CASH",
		Description:   fmt.Sprintf("Venda: %s", product),
	}

	result, err := opHandler.RecordSale(saleReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to record sale: %v", err), http.StatusInternalServerError)
		return
	}

	balance, _ := opHandler.GetCashBalance(entityID)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<div class="bg-green-100 border-l-4 border-green-500 text-green-700 p-4 mb-4" role="alert">
			<p class="font-bold">Venda Registrada!</p>
			<p>EntryID: %d | Saldo Caixa: R$ %.2f</p>
		</div>
	`, result.EntryID, float64(balance)/100)
}

func (h *PDVHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	entityID := r.URL.Query().Get("entity_id")
	if entityID == "" {
		entityID = "cooperativa_demo"
	}

	opHandler := usecase.NewOperationHandler(h.lifecycleManager)
	balance, err := opHandler.GetCashBalance(entityID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get balance: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `<span class="text-2xl font-bold text-green-600">R$ %.2f</span>`, float64(balance)/100)
}
