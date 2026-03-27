package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/providentia/digna/core_lume/pkg/eligibility"
	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

// EligibilityHandler handles eligibility profile HTTP requests
type EligibilityHandler struct {
	*BaseHandler
	lifecycleManager   lifecycle.LifecycleManager
	eligibilityService *eligibility.Service
}

// NewEligibilityHandler creates a new EligibilityHandler
func NewEligibilityHandler(lm lifecycle.LifecycleManager) (*EligibilityHandler, error) {
	base := NewBaseHandler(lm, true)

	eligService := eligibility.NewService(lm)

	return &EligibilityHandler{
		BaseHandler:        base,
		lifecycleManager:   lm,
		eligibilityService: eligService,
	}, nil
}

// RegisterRoutes registers the handler routes
func (h *EligibilityHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/eligibility", h.handleEligibilityPage)
	mux.HandleFunc("/eligibility/save", h.handleSaveProfile)
	mux.HandleFunc("/eligibility/status", h.handleGetStatus)
	mux.HandleFunc("/eligibility/export", h.handleExport)
}

// handleEligibilityPage displays the eligibility profile page
func (h *EligibilityHandler) handleEligibilityPage(w http.ResponseWriter, r *http.Request) {
	entityID := getEntityID(r)
	if entityID == "" {
		http.Error(w, "Entidade não identificada", http.StatusBadRequest)
		return
	}

	// Initialize table if needed
	if err := h.eligibilityService.EnsureTableExists(entityID); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao inicializar tabela: %v", err), http.StatusInternalServerError)
		return
	}

	ctx := r.Context()

	// Get or create profile
	profile, err := h.eligibilityService.GetOrCreateProfile(ctx, entityID, "")
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao carregar perfil: %v", err), http.StatusInternalServerError)
		return
	}

	// Check completion
	completionPercent, _ := h.eligibilityService.GetCompletionStatus(ctx, entityID)
	isComplete := completionPercent >= 100.0

	data := map[string]interface{}{
		"Title":             "Perfil para Crédito - Digna",
		"EntityID":          entityID,
		"Profile":           profile,
		"CompletionPercent": completionPercent,
		"IsComplete":        isComplete,
		"FinalidadeOptions": []struct {
			Value string
			Label string
		}{
			{"CAPITAL_GIRO", "Capital de Giro"},
			{"EQUIPAMENTO", "Equipamento"},
			{"REFORMA", "Reforma"},
			{"OUTRO", "Outro"},
		},
		"TipoEntidadeOptions": []struct {
			Value string
			Label string
		}{
			{"MEI", "Microempreendedor Individual (MEI)"},
			{"ME", "Microempresa (ME)"},
			{"EPP", "Empresa de Pequeno Porte (EPP)"},
			{"Cooperativa", "Cooperativa"},
			{"OSC", "Organização da Sociedade Civil (OSC)"},
			{"OSCIP", "OSCIP"},
			{"PF", "Pessoa Física"},
		},
	}

	// Load template from file
	tmplPath := "modules/ui_web/templates/eligibility_simple.html"
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

	tmpl, err := template.New("eligibility_simple.html").Funcs(funcMap).ParseFiles(tmplPath, componentPath)
	if err != nil {
		// Fallback to inline template
		h.renderInlineTemplate(w, data)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, data)
}

// handleSaveProfile saves the eligibility profile
func (h *EligibilityHandler) handleSaveProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	entityID := getEntityID(r)
	if entityID == "" {
		http.Error(w, "Entidade não identificada", http.StatusBadRequest)
		return
	}

	// Initialize table if needed
	if err := h.eligibilityService.EnsureTableExists(entityID); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao inicializar tabela: %v", err), http.StatusInternalServerError)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Erro ao parsear formulário", http.StatusBadRequest)
		return
	}

	// Build input from form
	input := eligibility.EligibilityInput{}

	if val := r.FormValue("inscrito_cad_unico"); val != "" {
		v := val == "true" || val == "on"
		input.InscritoCadUnico = &v
	}
	if val := r.FormValue("socio_mulher"); val != "" {
		v := val == "true" || val == "on"
		input.SocioMulher = &v
	}
	if val := r.FormValue("inadimplencia_ativa"); val != "" {
		v := val == "true" || val == "on"
		input.InadimplenciaAtiva = &v
	}
	if val := r.FormValue("contabilidade_formal"); val != "" {
		v := val == "true" || val == "on"
		input.ContabilidadeFormal = &v
	}
	if val := r.FormValue("finalidade_credito"); val != "" {
		input.FinalidadeCredito = &val
	}
	if val := r.FormValue("tipo_entidade"); val != "" {
		input.TipoEntidade = &val
	}
	if val := r.FormValue("valor_necessario"); val != "" {
		var v int64
		fmt.Sscanf(val, "%d", &v)
		// Convert from reais to centavos
		v = v * 100
		input.ValorNecessario = &v
	}

	// Get user ID from context or use a default
	userID := r.Context().Value("user_id")
	if userID == nil {
		userID = "system"
	}

	ctx := r.Context()
	_, err := h.eligibilityService.CreateOrUpdate(ctx, entityID, userID.(string), input)
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

	// Get updated status
	completionPercent, _ := h.eligibilityService.GetCompletionStatus(ctx, entityID)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<div class="bg-green-100 border-l-4 border-green-500 text-green-700 p-4" role="alert">
			<p class="font-bold">Perfil salvo!</p>
			<p>Completude: %.1f%%</p>
		</div>
	`, completionPercent)
}

// handleGetStatus returns the completion status
func (h *EligibilityHandler) handleGetStatus(w http.ResponseWriter, r *http.Request) {
	entityID := getEntityID(r)
	if entityID == "" {
		http.Error(w, "Entidade não identificada", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	percent, err := h.eligibilityService.GetCompletionStatus(ctx, entityID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"completion_percent": percent,
	})
}

// handleExport exports profile as JSON
func (h *EligibilityHandler) handleExport(w http.ResponseWriter, r *http.Request) {
	entityID := getEntityID(r)
	if entityID == "" {
		http.Error(w, "Entidade não identificada", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	profile, err := h.eligibilityService.GetProfile(ctx, entityID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao carregar perfil: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"perfil_%s.json\"", entityID))
	json.NewEncoder(w).Encode(profile)
}

// renderInlineTemplate renders an inline template as fallback
func (h *EligibilityHandler) renderInlineTemplate(w http.ResponseWriter, data map[string]interface{}) {
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
                    <span class="text-xl font-bold">Perfil para Crédito</span>
                </div>
            </div>
        </div>
    </nav>
    
    <main class="container mx-auto px-4 py-8">
        <div class="bg-white rounded-xl shadow-md p-6 mb-6">
            <h2 class="text-2xl font-bold text-gray-800 mb-4">Seu Perfil de Elegibilidade</h2>
            
            {{if .IsComplete}}
            <div class="bg-green-100 border-l-4 border-green-500 text-green-700 p-4 mb-4">
                <p class="font-bold">Perfil completo!</p>
                <p>Você está pronto para acessar oportunidades de crédito.</p>
            </div>
            {{else}}
            <div class="bg-yellow-100 border-l-4 border-yellow-500 text-yellow-700 p-4 mb-4">
                <p class="font-bold">Complete seu perfil</p>
                <p>Preencha os campos abaixo para habilitar o match com programas de crédito.</p>
            </div>
            {{end}}
            
            <div class="mb-6">
                <div class="flex justify-between mb-2">
                    <span class="text-sm font-medium text-gray-700">Completude do perfil</span>
                    <span class="text-sm font-medium text-gray-700">{{printf "%.1f" .CompletionPercent}}%</span>
                </div>
                <div class="w-full bg-gray-200 rounded-full h-2.5">
                    <div class="bg-blue-600 h-2.5 rounded-full" style="width: {{.CompletionPercent}}%"></div>
                </div>
            </div>
            
            <form hx-post="/eligibility/save" hx-target="#result" class="space-y-4">
                <input type="hidden" name="entity_id" value="{{.EntityID}}">
                
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">Finalidade do Crédito *</label>
                        <select name="finalidade_credito" required class="w-full p-2 border rounded-lg">
                            <option value="">Selecione...</option>
                            {{range .FinalidadeOptions}}
                            <option value="{{.Value}}" {{if eq $.Profile.FinalidadeCredito .Value}}selected{{end}}>{{.Label}}</option>
                            {{end}}
                        </select>
                    </div>
                    
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">Tipo de Entidade *</label>
                        <select name="tipo_entidade" required class="w-full p-2 border rounded-lg">
                            <option value="">Selecione...</option>
                            {{range .TipoEntidadeOptions}}
                            <option value="{{.Value}}" {{if eq $.Profile.TipoEntidade .Value}}selected{{end}}>{{.Label}}</option>
                            {{end}}
                        </select>
                    </div>
                </div>
                
                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">Valor Necessário (R$)</label>
                    <input type="number" name="valor_necessario" step="0.01" min="0" 
                           value="{{if .Profile.ValorNecessario}}{{printf "%.2f" (divide .Profile.ValorNecessario 100.0)}}{{end}}"
                           class="w-full p-2 border rounded-lg" placeholder="0.00">
                </div>
                
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div class="flex items-center">
                        <input type="checkbox" name="inscrito_cad_unico" id="inscrito_cad_unico" 
                               {{if .Profile.InscritoCadUnico}}checked{{end}} class="mr-2">
                        <label for="inscrito_cad_unico" class="text-sm text-gray-700">Inscrito no CadÚnico</label>
                    </div>
                    
                    <div class="flex items-center">
                        <input type="checkbox" name="socio_mulher" id="socio_mulher" 
                               {{if .Profile.SocioMulher}}checked{{end}} class="mr-2">
                        <label for="socio_mulher" class="text-sm text-gray-700">Sócio Mulher</label>
                    </div>
                    
                    <div class="flex items-center">
                        <input type="checkbox" name="inadimplencia_ativa" id="inadimplencia_ativa" 
                               {{if .Profile.InadimplenciaAtiva}}checked{{end}} class="mr-2">
                        <label for="inadimplencia_ativa" class="text-sm text-gray-700">Inadimplência Ativa</label>
                    </div>
                    
                    <div class="flex items-center">
                        <input type="checkbox" name="contabilidade_formal" id="contabilidade_formal" 
                               {{if .Profile.ContabilidadeFormal}}checked{{end}} class="mr-2">
                        <label for="contabilidade_formal" class="text-sm text-gray-700">Contabilidade Formal</label>
                    </div>
                </div>
                
                <div id="result" class="mt-4"></div>
                
                <div class="flex gap-4">
                    <button type="submit" class="flex-1 bg-blue-600 text-white py-2 px-4 rounded-lg hover:bg-blue-700">
                        Salvar Perfil
                    </button>
                    <a href="/eligibility/export?entity_id={{.EntityID}}" 
                       class="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 text-center">
                        Exportar JSON
                    </a>
                </div>
            </form>
        </div>
    </main>
</body>
</html>`

	tmpl, err := template.New("eligibility").Funcs(template.FuncMap{
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
