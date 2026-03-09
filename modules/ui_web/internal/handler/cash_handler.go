package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/providentia/digna/cash_flow/pkg/cash_flow"
	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

type CashHandler struct {
	lifecycleManager lifecycle.LifecycleManager
	cashAPI          *cash_flow.CashFlowAPI
	tmpl             *template.Template
}

func NewCashHandler(lm lifecycle.LifecycleManager) (*CashHandler, error) {
	funcMap := template.FuncMap{
		"divide": divide,
		"formatCurrency": func(amount int64) string {
			return fmt.Sprintf("R$ %.2f", float64(amount)/100)
		},
		"formatDate": func(t time.Time) string {
			return t.Format("02/01/2006 15:04")
		},
	}

	tmpl, err := template.New("templates").Funcs(funcMap).ParseGlob("templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	return &CashHandler{
		lifecycleManager: lm,
		cashAPI:          cash_flow.NewCashFlowAPI(lm),
		tmpl:             tmpl,
	}, nil
}

func (h *CashHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/cash", h.CashPage)
	mux.HandleFunc("/api/cash/entry", h.RecordEntry)
	mux.HandleFunc("/api/cash/balance", h.GetBalance)
	mux.HandleFunc("/api/cash/flow", h.GetCashFlow)
}

func (h *CashHandler) CashPage(w http.ResponseWriter, r *http.Request) {
	entityID := "cooperativa_demo"
	if r.URL.Query().Get("entity_id") != "" {
		entityID = r.URL.Query().Get("entity_id")
	}

	balance, _ := h.cashAPI.GetBalance(entityID)
	entries, _ := h.cashAPI.GetRecentEntries(entityID, 20)

	data := map[string]interface{}{
		"Title":      "Caixa - Digna",
		"EntityID":   entityID,
		"Balance":    balance,
		"Entries":    entries,
		"Categories": []string{"SALES", "EXPENSES", "SUPPLIERS", "BANK", "OTHER_INCOME", "OTHER_EXPENSE"},
	}

	if err := h.tmpl.ExecuteTemplate(w, "cash.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *CashHandler) RecordEntry(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	entityID := r.FormValue("entity_id")
	if entityID == "" {
		entityID = "cooperativa_demo"
	}

	entryType := r.FormValue("type")
	amountStr := r.FormValue("amount")
	category := r.FormValue("category")
	description := r.FormValue("description")

	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil || amount <= 0 {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	req := cash_flow.EntryRequest{
		EntityID:    entityID,
		Type:        entryType,
		Amount:      amount,
		Category:    category,
		Description: description,
	}

	result, err := h.cashAPI.RecordEntry(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to record entry: %v", err), http.StatusInternalServerError)
		return
	}

	if !result.Success {
		http.Error(w, result.Error, http.StatusBadRequest)
		return
	}

	balance, _ := h.cashAPI.GetBalance(entityID)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<div class="bg-green-100 border-l-4 border-green-500 text-green-700 p-4 mb-4" role="alert">
			<p class="font-bold">Movimento Registrado!</p>
			<p>Tipo: %s | Valor: R$ %.2f | Saldo: R$ %.2f</p>
		</div>
	`, entryType, float64(amount)/100, float64(balance.Balance)/100)
}

func (h *CashHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	entityID := r.URL.Query().Get("entity_id")
	if entityID == "" {
		entityID = "cooperativa_demo"
	}

	balance, err := h.cashAPI.GetBalance(entityID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get balance: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `<span class="text-2xl font-bold text-green-600">R$ %.2f</span>`, float64(balance.Balance)/100)
}

func (h *CashHandler) GetCashFlow(w http.ResponseWriter, r *http.Request) {
	entityID := r.URL.Query().Get("entity_id")
	if entityID == "" {
		entityID = "cooperativa_demo"
	}

	endDate := time.Now()
	startDate := endDate.AddDate(0, -1, 0)

	flow, err := h.cashAPI.GetCashFlow(entityID, startDate, endDate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get cash flow: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"balance": %d, "credit": %d, "debit": %d}`, flow.Balance, flow.TotalCredit, flow.TotalDebit)
}
