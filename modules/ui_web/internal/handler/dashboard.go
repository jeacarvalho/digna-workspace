package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

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
	// Create template with custom functions
	funcMap := template.FuncMap{
		"divide": divide,
	}

	tmpl, err := template.New("templates").Funcs(funcMap).ParseGlob("templates/*.html")
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

	data := map[string]interface{}{
		"Title":    "Digna - Providentia Foundation",
		"EntityID": "cooperativa_demo",
	}

	if err := h.tmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *DashboardHandler) DashboardPage(w http.ResponseWriter, r *http.Request) {
	entityID := r.URL.Query().Get("entity_id")
	if entityID == "" {
		entityID = "cooperativa_demo"
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

	data := map[string]interface{}{
		"Title":        "Painel de Dignidade",
		"EntityID":     entityID,
		"CashBalance":  float64(balance) / 100,
		"TotalSurplus": float64(calc.TotalSurplus) / 100,
		"TotalHours":   calc.TotalMinutes / 60,
		"MemberCount":  len(calc.Members),
		"Members":      calc.Members,
	}

	if err := h.tmpl.ExecuteTemplate(w, "dashboard.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *DashboardHandler) SocialClockPage(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":    "Ponto Social - ITG 2002",
		"EntityID": "cooperativa_demo",
	}

	if err := h.tmpl.ExecuteTemplate(w, "social_clock.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, fmt.Sprintf("Failed to record work: %v", err), http.StatusInternalServerError)
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
