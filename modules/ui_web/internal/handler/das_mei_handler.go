package handler

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/providentia/digna/core_lume/pkg/das"
	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

// DASMEIHandler handles DAS MEI related HTTP requests
type DASMEIHandler struct {
	*BaseHandler
	lifecycleManager lifecycle.LifecycleManager
	dasService       *das.Service
}

// NewDASMEIHandler creates a new DAS MEI handler
func NewDASMEIHandler(lm lifecycle.LifecycleManager) (*DASMEIHandler, error) {
	base := NewBaseHandler(lm, true)

	// Criar serviço DAS MEI
	dasService := das.NewService(lm)

	return &DASMEIHandler{
		BaseHandler:      base,
		lifecycleManager: lm,
		dasService:       dasService,
	}, nil
}

// RegisterRoutes registers the handler routes
func (h *DASMEIHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/das-mei", h.handleDASMEIPage)
	mux.HandleFunc("/das-mei/generate", h.handleGenerateDAS)
	mux.HandleFunc("/das-mei/{id}/pay", h.handleMarkAsPaid)
	mux.HandleFunc("/das-mei/alerts", h.handleAlerts)
}

// handleDASMEIPage displays the DAS MEI dashboard
func (h *DASMEIHandler) handleDASMEIPage(w http.ResponseWriter, r *http.Request) {
	entityID := getEntityID(r)
	if entityID == "" {
		http.Error(w, "Entidade não identificada", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// Inicializar tabela se necessário
	if err := h.dasService.EnsureTableExists(entityID); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao inicializar tabela: %v", err), http.StatusInternalServerError)
		return
	}

	// Atualizar status dos DAS pendentes (verificar vencimentos)
	h.dasService.UpdateDASStatus(ctx, entityID)

	// Obter todos os DAS da entidade
	dasList, err := h.dasService.GetAllDAS(ctx, entityID)
	if err != nil {
		dasList = []*das.DASMEI{}
	}

	// Obter DAS pendentes
	pendingDAS, _ := h.dasService.GetPendingDAS(ctx, entityID)

	// Obter DAS vencidos
	overdueDAS, _ := h.dasService.GetOverdueDAS(ctx, entityID)

	// Obter alertas
	alerts, _ := h.dasService.CheckOverdueAlerts(ctx, entityID)

	// Calcular total pendente
	var totalPending int64
	for _, d := range pendingDAS {
		totalPending += d.ValorDevido
	}

	// Calcular total vencido
	var totalOverdue int64
	for _, d := range overdueDAS {
		totalOverdue += d.ValorDevido
	}

	// Obter competência atual
	currentCompetencia := h.dasService.GetCurrentCompetencia()

	// Verificar se existe DAS para o mês atual
	currentDAS, _ := h.dasService.GetDASByCompetencia(ctx, entityID, currentCompetencia)
	currentDASExists := currentDAS != nil

	data := map[string]interface{}{
		"Title":              "DAS MEI - Digna",
		"EntityID":           entityID,
		"DASList":            dasList,
		"PendingDAS":         pendingDAS,
		"OverdueDAS":         overdueDAS,
		"Alerts":             alerts,
		"TotalPending":       totalPending,
		"TotalOverdue":       totalOverdue,
		"CurrentCompetencia": currentCompetencia,
		"CurrentDASExists":   currentDASExists,
		"CurrentDAS":         currentDAS,
		"ActivityTypes": []struct {
			Value string
			Label string
		}{
			{"COMERCIO", "Comércio (ICMS)"},
			{"SERVICOS", "Serviços (ISS)"},
			{"MISTO", "Comércio + Serviços"},
		},
	}

	// Carregar template do arquivo
	tmplPath := "modules/ui_web/templates/das-mei_simple.html"
	componentPath := "modules/ui_web/templates/components/help_tooltip.html"

	// Create function map with dict support
	funcMap := template.FuncMap{
		"divide": func(a, b int64) float64 {
			if b == 0 {
				return 0
			}
			return float64(a) / float64(b)
		},
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("dict requires even number of arguments")
			}
			dict := make(map[string]interface{})
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
	}

	tmpl, err := template.New("das-mei_simple.html").Funcs(funcMap).ParseFiles(tmplPath, componentPath)
	if err != nil {
		// Fallback para template inline se arquivo não existir
		h.renderInlineTemplate(w, data)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, data)
}

// handleGenerateDAS generates a new DAS MEI
func (h *DASMEIHandler) handleGenerateDAS(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	entityID := getEntityID(r)
	if entityID == "" {
		http.Error(w, "Entidade não identificada", http.StatusBadRequest)
		return
	}

	// Inicializar tabela se necessário
	if err := h.dasService.EnsureTableExists(entityID); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao inicializar tabela: %v", err), http.StatusInternalServerError)
		return
	}

	competencia := r.FormValue("competencia")
	activityType := r.FormValue("activity_type")

	if competencia == "" {
		competencia = h.dasService.GetCurrentCompetencia()
	}

	ctx := r.Context()

	req := &das.GenerateDASRequest{
		Competencia:  competencia,
		ActivityType: das.ActivityType(activityType),
	}

	d, err := h.dasService.GenerateMonthlyDAS(ctx, entityID, req)
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `
			<div class="bg-red-100 border-l-4 border-red-500 text-red-700 p-4" role="alert">
				<p class="font-bold">Erro!</p>
				<p>%s</p>
			</div>
		`, err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<div class="bg-green-100 border-l-4 border-green-500 text-green-700 p-4" role="alert">
			<p class="font-bold">DAS MEI Gerado!</p>
			<p>Competência: %s | Valor: R$ %.2f</p>
		</div>
	`, d.Competencia, float64(d.ValorDevido)/100.0)
}

// handleMarkAsPaid marks a DAS MEI as paid
func (h *DASMEIHandler) handleMarkAsPaid(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	entityID := getEntityID(r)
	if entityID == "" {
		http.Error(w, "Entidade não identificada", http.StatusBadRequest)
		return
	}

	// Inicializar tabela se necessário
	if err := h.dasService.EnsureTableExists(entityID); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao inicializar tabela: %v", err), http.StatusInternalServerError)
		return
	}

	// Extrair ID da URL
	dasID := r.PathValue("id")
	if dasID == "" {
		http.Error(w, "ID do DAS não informado", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	err := h.dasService.MarkAsPaid(ctx, entityID, dasID)
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `
			<div class="bg-red-100 border-l-4 border-red-500 text-red-700 p-4" role="alert">
				<p class="font-bold">Erro!</p>
				<p>%s</p>
			</div>
		`, err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<div class="bg-green-100 border-l-4 border-green-500 text-green-700 p-4" role="alert">
			<p class="font-bold">DAS MEI Marcado como Pago!</p>
			<p>O pagamento foi registrado com sucesso.</p>
		</div>
	`)
}

// handleAlerts returns DAS MEI alerts for HTMX
func (h *DASMEIHandler) handleAlerts(w http.ResponseWriter, r *http.Request) {
	entityID := getEntityID(r)
	if entityID == "" {
		http.Error(w, "Entidade não identificada", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	alerts, err := h.dasService.CheckOverdueAlerts(ctx, entityID)
	if err != nil || len(alerts) == 0 {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<span class="text-gray-500">Nenhum alerta</span>`)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	for _, alert := range alerts {
		var bgClass string
		switch alert.Severity {
		case "CRITICAL":
			bgClass = "bg-red-100 border-red-500 text-red-700"
		case "WARNING":
			bgClass = "bg-yellow-100 border-yellow-500 text-yellow-700"
		default:
			bgClass = "bg-blue-100 border-blue-500 text-blue-700"
		}

		fmt.Fprintf(w, `
			<div class="%s border-l-4 p-3 mb-2 text-sm" role="alert">
				<p>%s</p>
			</div>
		`, bgClass, alert.Message)
	}
}

// renderInlineTemplate renders an inline template as fallback
func (h *DASMEIHandler) renderInlineTemplate(w http.ResponseWriter, data map[string]interface{}) {
	htmlTemplate := `<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
</head>
<body class="bg-gray-50 min-h-screen">
    <nav class="bg-blue-600 text-white shadow-lg sticky top-0 z-50">
        <div class="container mx-auto px-4 py-3">
            <div class="flex items-center justify-between">
                <div class="flex items-center space-x-2">
                    <a href="/" class="text-white hover:text-blue-100">← Voltar</a>
                    <span class="text-xl font-bold">DAS MEI</span>
                </div>
            </div>
        </div>
    </nav>
    
    <main class="container mx-auto px-4 py-8">
        <div class="bg-white rounded-xl shadow-md p-6 mb-6">
            <h2 class="text-2xl font-bold text-gray-800 mb-4">Cálculo do DAS MEI</h2>
            <p class="text-gray-600 mb-4">Sistema de cálculo automático do Documento de Arrecadação do Simples Nacional para Microempreendedores Individuais.</p>
            
            {{if .CurrentDAS}}
            <div class="bg-blue-100 border-l-4 border-blue-500 text-blue-700 p-4 mb-4">
                <p class="font-bold">DAS MEI do Mês Atual ({{.CurrentCompetencia}})</p>
                <p>Valor: R$ {{printf "%.2f" (divide .CurrentDAS.ValorDevido 100.0)}}</p>
                <p>Status: {{.CurrentDAS.Status}}</p>
                {{if eq .CurrentDAS.Status "PENDENTE"}}
                <button hx-post="/das-mei/{{.CurrentDAS.ID}}/pay" hx-target="#result" class="mt-2 bg-green-600 text-white py-1 px-3 rounded hover:bg-green-700">
                    Marcar como Pago
                </button>
                {{end}}
            </div>
            {{else}}
            <div class="bg-yellow-100 border-l-4 border-yellow-500 text-yellow-700 p-4 mb-4">
                <p class="font-bold">DAS MEI não gerado para {{.CurrentCompetencia}}</p>
                <form hx-post="/das-mei/generate" hx-target="#result" class="mt-2">
                    <input type="hidden" name="competencia" value="{{.CurrentCompetencia}}">
                    <select name="activity_type" class="p-2 border rounded mr-2">
                        {{range .ActivityTypes}}
                        <option value="{{.Value}}">{{.Label}}</option>
                        {{end}}
                    </select>
                    <button type="submit" class="bg-blue-600 text-white py-1 px-3 rounded hover:bg-blue-700">
                        Gerar DAS
                    </button>
                </form>
            </div>
            {{end}}
            
            <div id="result" class="mb-4"></div>
        </div>
        
        {{if .Alerts}}
        <div class="bg-white rounded-xl shadow-md p-6 mb-6">
            <h3 class="text-lg font-semibold text-gray-700 mb-4">Alertas</h3>
            {{range .Alerts}}
            <div class="{{if eq .Severity "CRITICAL"}}bg-red-100 border-red-500 text-red-700{{else if eq .Severity "WARNING"}}bg-yellow-100 border-yellow-500 text-yellow-700{{else}}bg-blue-100 border-blue-500 text-blue-700{{end}} border-l-4 p-3 mb-2">
                {{.Message}}
            </div>
            {{end}}
        </div>
        {{end}}
        
        {{if .DASList}}
        <div class="bg-white rounded-xl shadow-md p-6">
            <h3 class="text-lg font-semibold text-gray-700 mb-4">Histórico</h3>
            <table class="min-w-full divide-y divide-gray-200">
                <thead class="bg-gray-50">
                    <tr>
                        <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Competência</th>
                        <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Valor</th>
                        <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
                    </tr>
                </thead>
                <tbody class="divide-y divide-gray-200">
                    {{range .DASList}}
                    <tr>
                        <td class="px-4 py-3">{{.Competencia}}</td>
                        <td class="px-4 py-3">R$ {{printf "%.2f" (divide .ValorDevido 100.0)}}</td>
                        <td class="px-4 py-3">
                            <span class="px-2 py-1 text-xs font-semibold rounded-full {{if eq .Status "PAGO"}}bg-green-100 text-green-800{{else if eq .Status "VENCIDO"}}bg-red-100 text-red-800{{else}}bg-yellow-100 text-yellow-800{{end}}">
                                {{.Status}}
                            </span>
                        </td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
        </div>
        {{end}}
    </main>
</body>
</html>`

	tmpl, err := template.New("das_mei").Funcs(template.FuncMap{
		"divide": func(a, b float64) float64 {
			if b == 0 {
				return 0
			}
			return a / b
		},
	}).Parse(htmlTemplate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create template: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, data)
}
