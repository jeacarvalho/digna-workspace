package handler

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/providentia/digna/cash_flow/pkg/cash_flow"
	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

type CashHandler struct {
	*BaseHandler
	lifecycleManager lifecycle.LifecycleManager
	cashAPI          *cash_flow.CashFlowAPI
}

func NewCashHandler(lm lifecycle.LifecycleManager) (*CashHandler, error) {
	// Criar BaseHandler com TemplateManager
	base := NewBaseHandler(lm, true)

	return &CashHandler{
		BaseHandler:      base,
		lifecycleManager: lm,
		cashAPI:          cash_flow.NewCashFlowAPI(lm),
	}, nil
}

func (h *CashHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/cash", h.CashPage)
	mux.HandleFunc("/api/cash/entry", h.RecordEntry)
	mux.HandleFunc("/api/cash/balance", h.GetBalance)
	mux.HandleFunc("/api/cash/flow", h.GetCashFlow)
}

func (h *CashHandler) CashPage(w http.ResponseWriter, r *http.Request) {
	entityID := r.URL.Query().Get("entity_id")
	if entityID == "" {
		http.Error(w, "entity_id é obrigatório", http.StatusBadRequest)
		return
	}

	balanceResp, _ := h.cashAPI.GetBalance(entityID)

	var balance int64
	if balanceResp != nil {
		balance = balanceResp.Balance
	}

	// Sempre buscar do banco diretamente já que a API retorna lista vazia
	entries := h.getEntriesFromDatabase(entityID, 20)

	data := map[string]interface{}{
		"Title":      "Caixa - Digna",
		"EntityID":   entityID,
		"Balance":    float64(balance), // Converter para float64 para compatibilidade com fdiv
		"Entries":    entries,
		"Categories": []string{"VENDAS", "DESPESAS", "FORNECEDORES", "BANCO", "OUTRAS ENTRADAS", "OUTRAS SAÍDAS"},
	}

	// Usar TemplateManager do BaseHandler
	h.RenderTemplate(w, "cash_simple.html", data)
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

// CashEntry representa uma entrada de caixa
type CashEntry struct {
	ID          int64
	EntityID    string
	Type        string
	Amount      int64
	Description string
	Category    string
	Date        time.Time
	CreatedAt   time.Time
}

// getEntriesFromDatabase busca entradas diretamente do banco de dados
func (h *CashHandler) getEntriesFromDatabase(entityID string, limit int) []CashEntry {
	dbPath := fmt.Sprintf("../../data/entities/%s.db", entityID)
	slog.Info("CashHandler - Buscando entradas do banco", "db_path", dbPath, "entity_id", entityID, "limit", limit)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		slog.Error("CashHandler - Failed to open database", "error", err, "entity_id", entityID, "db_path", dbPath)
		return []CashEntry{}
	}
	defer db.Close()

	// Testar conexão
	if err := db.Ping(); err != nil {
		slog.Error("CashHandler - Failed to ping database", "error", err, "entity_id", entityID)
		return []CashEntry{}
	}

	// Buscar movimentos de caixa (vendas e compras)
	query := `
		SELECT 
			e.id,
			e.entry_date,
			e.description,
			e.reference,
			p.amount,
			a.name as account_name,
			p.direction,
			a.code
		FROM entries e
		JOIN postings p ON e.id = p.entry_id
		JOIN accounts a ON p.account_id = a.id
		WHERE (
			-- Vendas: crédito em receitas (3.x)
			(a.code LIKE '3.%' AND p.direction = 'CREDIT' AND e.description LIKE 'Venda PDV:%')
			OR
			-- Compras à vista: débito em caixa (1)
			(a.code = '1' AND p.direction = 'DEBIT' AND e.description LIKE 'Compra %')
			OR
			-- Saídas de caixa manualmente registradas
			(a.code = '1' AND p.direction = 'DEBIT' AND e.description NOT LIKE 'Venda PDV:%' AND e.description NOT LIKE 'Compra %')
			OR
			-- Entradas de caixa manualmente registradas  
			(a.code = '1' AND p.direction = 'CREDIT' AND e.description NOT LIKE 'Venda PDV:%' AND e.description NOT LIKE 'Compra %')
		)
		ORDER BY e.entry_date DESC
		LIMIT ?
	`

	slog.Info("CashHandler - Executando query", "query", query, "limit", limit)
	rows, err := db.Query(query, limit)
	if err != nil {
		slog.Error("CashHandler - Failed to query entries", "error", err, "entity_id", entityID, "query", query)
		return []CashEntry{}
	}
	defer rows.Close()

	var entries []CashEntry
	rowCount := 0
	for rows.Next() {
		var id int64
		var entryDate int64
		var description, reference string
		var amount int64
		var accountName, direction, accountCode string

		err := rows.Scan(&id, &entryDate, &description, &reference, &amount, &accountName, &direction, &accountCode)
		if err != nil {
			slog.Error("CashHandler - Failed to scan row", "error", err, "row_count", rowCount)
			continue
		}

		// Determinar tipo baseado na direção e conta
		entryType := "CREDIT"
		category := "OUTROS"

		// Lógica para determinar tipo e categoria
		if strings.Contains(strings.ToUpper(description), "VENDA PDV") {
			category = "VENDAS"
			entryType = "CREDIT" // Vendas são entradas (crédito)
		} else if strings.Contains(strings.ToUpper(description), "COMPRA") {
			category = "COMPRAS"
			// Compras à vista: débito em caixa = saída
			if accountCode == "1" && direction == "DEBIT" {
				entryType = "DEBIT" // Saída de caixa
			} else if accountCode == "1" && direction == "CREDIT" {
				entryType = "CREDIT" // Entrada de caixa (devolução?)
			}
		} else if accountCode == "1" {
			// Movimentos manuais em caixa
			if direction == "CREDIT" {
				entryType = "CREDIT"
				category = "ENTRADA"
			} else {
				entryType = "DEBIT"
				category = "SAÍDA"
			}
		}

		entry := CashEntry{
			ID:          id,
			EntityID:    entityID,
			Type:        entryType,
			Amount:      amount,
			Description: description,
			Category:    category,
			Date:        time.Unix(entryDate, 0),
			CreatedAt:   time.Unix(entryDate, 0),
		}
		entries = append(entries, entry)
		rowCount++
		slog.Info("CashHandler - Linha processada", "id", id, "description", description, "amount", amount, "date", time.Unix(entryDate, 0))
	}

	if err := rows.Err(); err != nil {
		slog.Error("CashHandler - Error iterating rows", "error", err)
	}

	slog.Info("CashHandler - Entradas recuperadas do banco", "count", len(entries), "entity_id", entityID, "rows_processed", rowCount)
	return entries
}
