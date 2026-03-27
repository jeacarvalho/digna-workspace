package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
	"github.com/providentia/digna/pdv_ui/usecase"
	"github.com/providentia/digna/reporting/pkg/surplus"
)

// Divide function for templates - handles both int64 and float64
func divide(a, b interface{}) float64 {
	var af, bf float64

	switch v := a.(type) {
	case float64:
		af = v
	case int64:
		af = float64(v)
	case int:
		af = float64(v)
	default:
		af = 0
	}

	switch v := b.(type) {
	case float64:
		bf = v
	case int64:
		bf = float64(v)
	case int:
		bf = float64(v)
	default:
		bf = 0
	}

	if bf == 0 {
		return 0
	}
	return af / bf
}

type DashboardHandler struct {
	lifecycleManager lifecycle.LifecycleManager
	tmpl             *template.Template
}

func NewDashboardHandler(lm lifecycle.LifecycleManager) (*DashboardHandler, error) {
	funcMap := template.FuncMap{
		"divide": divide,
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
		"fdiv": func(a, b float64) float64 {
			if b == 0 {
				return 0
			}
			return a / b
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

	// Criar template vazio - vamos carregar templates do disco quando necessário
	tmpl := template.New("").Funcs(funcMap)

	// Parsear templates necessários
	_, err := tmpl.ParseFiles("templates/dashboard_simple.html", "templates/social_clock.html", "templates/components/help_tooltip.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	return &DashboardHandler{
		lifecycleManager: lm,
		tmpl:             tmpl,
	}, nil
}

func (h *DashboardHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", h.HomePage)
	mux.HandleFunc("/dashboard", h.DashboardPage)
	mux.HandleFunc("/social", h.SocialClockPage)
	mux.HandleFunc("/api/social/record", h.RecordWork)
}

func (h *DashboardHandler) HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Redirecionar para login se não estiver autenticado
	// A verificação de autenticação será feita pelo middleware
	http.Redirect(w, r, "/login", http.StatusFound)
}

func (h *DashboardHandler) DashboardPage(w http.ResponseWriter, r *http.Request) {
	entityID := r.URL.Query().Get("entity_id")
	if entityID == "" {
		http.Error(w, "entity_id é obrigatório", http.StatusBadRequest)
		return
	}

	calculator := surplus.NewCalculator(h.lifecycleManager)
	calc, err := calculator.CalculateSocialSurplus(entityID)
	if err != nil {
		calc = &surplus.SurplusCalculation{
			EntityID:     entityID,
			TotalSurplus: 0,
			TotalMinutes: 0,
			Members:      []surplus.MemberShare{},
		}
	}

	opHandler := usecase.NewOperationHandler(h.lifecycleManager)
	balance, _ := opHandler.GetCashBalance(entityID)

	// Calcular sobras disponíveis após reservas (15% para Reserva Legal + FATES)
	availableSurplus := float64(calc.TotalSurplus) * 0.85 / 100 // 85% após 15% de reservas

	data := map[string]interface{}{
		"Title":            "Painel de Dignidade",
		"EntityID":         entityID,
		"CashBalance":      float64(balance) / 100,
		"TotalSurplus":     float64(calc.TotalSurplus) / 100,
		"AvailableSurplus": availableSurplus,
		"TotalHours":       calc.TotalMinutes / 60,
		"MemberCount":      len(calc.Members),
		"Members":          calc.Members,
	}

	// Usar template do handler (que já tem todas as funções incluindo fdiv)
	if err := h.tmpl.ExecuteTemplate(w, "dashboard_simple.html", data); err != nil {
		// Fallback para template antigo
		if err := h.tmpl.ExecuteTemplate(w, "dashboard.html", data); err != nil {
			http.Error(w, fmt.Sprintf("Erro ao renderizar template: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

func (h *DashboardHandler) SocialClockPage(w http.ResponseWriter, r *http.Request) {
	entityID := r.URL.Query().Get("entity_id")
	if entityID == "" {
		http.Error(w, "entity_id é obrigatório", http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"Title":    "Ponto Social - Digna",
		"EntityID": entityID,
	}

	// Usar template do handler (que já tem todas as funções incluindo fdiv)
	if err := h.tmpl.ExecuteTemplate(w, "social_clock.html", data); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao renderizar template: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *DashboardHandler) RecordWork(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	entityID := r.FormValue("entity_id")
	if entityID == "" {
		entityID = "cooperativa_demo"
	}

	memberID := r.FormValue("member_id")
	if memberID == "" {
		memberID = "socio_001"
	}

	minutesStr := r.FormValue("minutes")
	minutes, err := strconv.ParseInt(minutesStr, 10, 64)
	if err != nil || minutes <= 0 {
		http.Error(w, "Invalid minutes", http.StatusBadRequest)
		return
	}

	sgHandler := usecase.NewSocialGovernanceHandler(h.lifecycleManager)

	workReq := usecase.WorkRequest{
		EntityID:     entityID,
		MemberID:     memberID,
		Minutes:      minutes,
		ActivityType: "PRODUCAO",
		Description:  "Trabalho registrado via interface web",
	}

	if err := sgHandler.RecordWork(workReq); err != nil {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `
			<div class="bg-red-100 border-l-4 border-red-500 text-red-700 p-4 mb-4" role="alert">
				<p class="font-bold">Erro ao registrar horas!</p>
				<p>%v</p>
				<p class="text-sm mt-2">Verifique se o módulo core_lume está funcionando.</p>
			</div>
		`, err)
		return
	}

	totalMinutes, count, _ := sgHandler.GetMemberWorkCapital(entityID, memberID)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<div class="bg-blue-100 border-l-4 border-blue-500 text-blue-700 p-4 mb-4" role="alert">
			<p class="font-bold">Horas Registradas!</p>
			<p>%s: %d minutos (Total: %d min em %d registros)</p>
		</div>
	`, memberID, minutes, totalMinutes, count)
}
