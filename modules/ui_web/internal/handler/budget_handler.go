package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/providentia/digna/budget/pkg/budget"
	"github.com/providentia/digna/cash_flow/pkg/cash_flow"
	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

// BudgetHandler lida com requisições relacionadas a orçamento
type BudgetHandler struct {
	budgetAPI   budget.BudgetAPI
	cashFlowAPI *cash_flow.CashFlowAPI
	tmpl        *template.Template
}

// NewBudgetHandler cria um novo handler de orçamento
func NewBudgetHandler(lm lifecycle.LifecycleManager) (*BudgetHandler, error) {
	// Criar template com funções auxiliares
	funcMap := template.FuncMap{
		"formatCurrency": func(amount int64) string {
			return fmt.Sprintf("R$ %.2f", float64(amount)/100)
		},
		"formatDate": func(t time.Time) string {
			return t.Format("02/01/2006")
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
		"fdiv": func(a, b float64) float64 {
		if b == 0 {
			return 0
		}
		return a / b
	},
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
	}

	// Template HTML básico (será substituído pelo arquivo se existir)
	htmlTemplate := `<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} - Digna</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
</head>
<body class="bg-gray-50 min-h-screen">
    <nav class="bg-blue-600 text-white shadow-lg sticky top-0 z-50">
        <div class="container mx-auto px-4 py-3">
            <div class="flex items-center justify-between">
                <div class="flex items-center space-x-2">
                    <a href="/" class="text-white hover:text-blue-100">← Voltar</a>
                    <span class="text-xl font-bold">{{.Title}}</span>
                </div>
                <div class="text-sm opacity-90">{{.Period}}</div>
            </div>
        </div>
    </nav>
    
    <main class="container mx-auto px-4 py-8">
        <div class="bg-white rounded-xl shadow-md p-6 mb-6">
            <h2 class="text-2xl font-bold text-gray-800 mb-4">Dashboard de Orçamento</h2>
            {{if .Executions}}
            <div class="overflow-x-auto">
                <table class="min-w-full divide-y divide-gray-200">
                    <thead class="bg-gray-50">
                        <tr>
                            <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Categoria</th>
                            <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Planejado</th>
                            <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Executado</th>
                            <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
                        </tr>
                    </thead>
                    <tbody class="divide-y divide-gray-200">
                        {{range .Executions}}
                        <tr>
                            <td class="px-4 py-3">{{getCategoryLabel .Plan.Category}}</td>
                            <td class="px-4 py-3">{{formatCurrency .Plan.Planned}}</td>
                            <td class="px-4 py-3">{{formatCurrency .Executed}} ({{.Percentage}}%)</td>
                            <td class="px-4 py-3">
                                <span class="px-2 py-1 text-xs font-semibold rounded-full {{getAlertStatusClass .AlertStatus}}">
                                    {{getAlertStatusLabel .AlertStatus}}
                                </span>
                            </td>
                        </tr>
                        {{end}}
                    </tbody>
                </table>
            </div>
            {{else}}
            <div class="text-center py-8 text-gray-500">
                <p class="text-lg">Nenhum planejamento encontrado para {{.Period}}</p>
                <p class="text-sm mt-2">Use o formulário abaixo para adicionar um planejamento.</p>
            </div>
            {{end}}
        </div>
        
        <div class="bg-white rounded-xl shadow-md p-6">
            <h3 class="text-lg font-semibold text-gray-700 mb-4">Adicionar Novo Planejamento</h3>
            <form hx-post="/api/budget/plan" hx-target="#result" hx-swap="innerHTML">
                <input type="hidden" name="entity_id" value="{{.EntityID}}">
                <input type="hidden" name="period" value="{{.Period}}">
                
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">Categoria</label>
                        <select name="category" required class="w-full p-2 border rounded">
                            {{range .Categories}}
                            <option value="{{.}}">{{getCategoryLabel .}}</option>
                            {{end}}
                        </select>
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">Valor (R$)</label>
                        <input type="number" name="planned" step="0.01" min="0.01" required 
                               class="w-full p-2 border rounded" placeholder="1000.00">
                    </div>
                </div>
                
                <div class="mb-4">
                    <label class="block text-sm font-medium text-gray-700 mb-2">Descrição</label>
                    <textarea name="description" rows="2" class="w-full p-2 border rounded"></textarea>
                </div>
                
                <div id="result" class="mb-4"></div>
                
                <button type="submit" class="w-full bg-blue-600 text-white py-2 px-4 rounded hover:bg-blue-700">
                    Salvar Planejamento
                </button>
            </form>
        </div>
    </main>
</body>
</html>`

	// Criar template
	tmpl, err := template.New("budget").Funcs(funcMap).Parse(htmlTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to create template: %w", err)
	}

	cashFlowAPI := cash_flow.NewCashFlowAPI(lm)
	cashFlowAdapter := budget.NewCashFlowAdapter(cashFlowAPI)
	budgetAPI := budget.NewBudgetAPI(lm, cashFlowAdapter)

	return &BudgetHandler{
		budgetAPI:   budgetAPI,
		cashFlowAPI: cashFlowAPI,
		tmpl:        tmpl,
	}, nil
}

// RegisterRoutes registra as rotas do handler
func (h *BudgetHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /budget", h.handleBudgetDashboard)
	mux.HandleFunc("GET /budget/report", h.handleBudgetReport)
	mux.HandleFunc("POST /api/budget/plan", h.handleCreatePlan)
	mux.HandleFunc("DELETE /api/budget/plan/{id}", h.handleDeletePlan)
	mux.HandleFunc("GET /api/budget/categories", h.handleGetCategories)
	mux.HandleFunc("GET /api/budget/periods", h.handleGetPeriods)
}

// handleBudgetDashboard exibe o dashboard de orçamento
func (h *BudgetHandler) handleBudgetDashboard(w http.ResponseWriter, r *http.Request) {
	entityID := getEntityID(r)
	if entityID == "" {
		http.Error(w, "Entidade não identificada", http.StatusBadRequest)
		return
	}

	// Obter período atual (YYYY-MM)
	currentPeriod := time.Now().Format("2006-01")

	// Obter relatório de execução
	ctx := r.Context()
	executionsPtr, err := h.budgetAPI.GetExecutionReport(ctx, entityID, currentPeriod)
	var executions []budget.BudgetExecution
	if err != nil || executionsPtr == nil {
		// Se não houver dados, usar array vazio
		executions = []budget.BudgetExecution{}
	} else {
		// Converter ponteiros para valores
		for _, execPtr := range executionsPtr {
			if execPtr != nil {
				executions = append(executions, *execPtr)
			}
		}
	}

	// Obter categorias disponíveis
	categoriesPtr, err := h.budgetAPI.GetCategories(ctx)
	var categories []string
	if err != nil || categoriesPtr == nil {
		categories = []string{"INSUMOS", "ENERGIA", "EQUIPAMENTOS", "TRANSPORTE", "MANUTENCAO", "SERVICOS", "OUTROS"}
	} else {
		// Extrair IDs das categorias
		for _, catPtr := range categoriesPtr {
			if catPtr != nil {
				categories = append(categories, catPtr.ID)
			}
		}
	}

	data := map[string]interface{}{
		"Title":      "Planejamento do Mês",
		"EntityID":   entityID,
		"Period":     currentPeriod,
		"Executions": executions,
		"Categories": categories,
	}

	if err := h.tmpl.Execute(w, data); err != nil {
		// Log do erro para debug
		fmt.Printf("DEBUG: Erro ao executar template: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Outros métodos permanecem iguais...
func (h *BudgetHandler) handleBudgetReport(w http.ResponseWriter, r *http.Request) {
	// Implementação simplificada
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<h1>Relatório de Orçamento</h1><p>Em desenvolvimento</p>")
}

func (h *BudgetHandler) handleCreatePlan(w http.ResponseWriter, r *http.Request) {
	entityID := getEntityID(r)
	if entityID == "" {
		http.Error(w, "Entidade não identificada", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
		return
	}

	period := r.FormValue("period")
	category := r.FormValue("category")
	description := r.FormValue("description")

	plannedStr := r.FormValue("planned")
	planned, err := strconv.ParseInt(plannedStr, 10, 64)
	if err != nil || planned <= 0 {
		// Tentar converter de decimal para centavos
		plannedFloat, err := strconv.ParseFloat(plannedStr, 64)
		if err != nil || plannedFloat <= 0 {
			http.Error(w, "Valor planejado inválido", http.StatusBadRequest)
			return
		}
		planned = int64(plannedFloat * 100)
	}

	req := budget.BudgetPlanRequest{
		EntityID:    entityID,
		Period:      period,
		Category:    category,
		Planned:     planned,
		Description: description,
	}

	ctx := r.Context()
	resp, err := h.budgetAPI.CreatePlan(ctx, req)
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<div class="bg-red-100 text-red-700 p-3 rounded">Erro: %v</div>`, err)
		return
	}

	if !resp.Success {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<div class="bg-yellow-100 text-yellow-700 p-3 rounded">%s</div>`, resp.Error)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `<div class="bg-green-100 text-green-700 p-3 rounded">
		Plano criado com sucesso! ID: %s
		<script>setTimeout(() => location.reload(), 1500)</script>
	</div>`, resp.PlanID)
}

func (h *BudgetHandler) handleDeletePlan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<p>Delete plan - em desenvolvimento</p>")
}

func (h *BudgetHandler) handleGetCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"categories": ["INSUMOS","ENERGIA","EQUIPAMENTOS","TRANSPORTE","MANUTENCAO","SERVICOS","OUTROS"]}`)
}

func (h *BudgetHandler) handleGetPeriods(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"periods": ["%s"]}`, time.Now().Format("2006-01"))
}

// Helper functions
func getEntityID(r *http.Request) string {
	entityID := r.URL.Query().Get("entity_id")
	if entityID == "" {
		entityID = "test-entity"
	}
	return entityID
}
