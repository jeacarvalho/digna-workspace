package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"digna/accountant_dashboard/pkg/dashboard"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

type AccountantHandler struct {
	*BaseHandler
	dashboardService dashboard.DashboardService
	repoFactory      dashboard.RepositoryFactory
}

func NewAccountantHandler(lm lifecycle.LifecycleManager) (*AccountantHandler, error) {
	// Obter devMode do ambiente (mesmo padrão usado no middleware)
	devMode := os.Getenv("DEV") != "false" && os.Getenv("DEV") != "0"
	baseHandler := NewBaseHandler(lm, devMode)

	// Create repository factory with data directory from lifecycle manager
	dataDir := "../../data" // Default, should come from config
	repoFactory := dashboard.NewSQLiteRepositoryFactory(dataDir)

	// For now, create a service with a dummy repository
	// In production, we'd create repository per entity as needed
	dummyRepo, err := repoFactory.NewRepository("dummy")
	if err != nil {
		// If we can't create a repository, create a service without it
		// The service will fail at runtime when trying to access data
		dummyRepo = nil
	}

	dashboardService := dashboard.NewDashboardService(dummyRepo)

	return &AccountantHandler{
		BaseHandler:      baseHandler,
		dashboardService: dashboardService,
		repoFactory:      repoFactory,
	}, nil
}

func (h *AccountantHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/accountant/dashboard", h.Dashboard)
	mux.HandleFunc("/accountant/export", h.ExportFiscal)
	mux.HandleFunc("/accountant/export/", h.ExportFiscal) // Para rotas com parâmetros
}

func (h *AccountantHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	if period == "" {
		period = time.Now().Format("2006-01")
	}

	pendingEntities, _ := h.dashboardService.ListPendingEntities(r.Context(), period)

	entities := make([]map[string]interface{}, len(pendingEntities))
	for i, entityID := range pendingEntities {
		history, _ := h.dashboardService.GetExportHistory(r.Context(), entityID, period)
		entities[i] = map[string]interface{}{
			"ID":           entityID,
			"Name":         entityID,
			"Status":       "PENDING",
			"PendingMonth": period,
			"HasExports":   len(history) > 0,
		}
	}

	// For now, use default mappings - in production we'd get these from a config
	defaultMappings := []dashboard.AccountMapping{
		{LocalCode: "1.1.1.01", LocalName: "Caixa", StandardCode: "1.1.1.01", StandardName: "Caixa"},
		{LocalCode: "1.1.2.01", LocalName: "Bancos", StandardCode: "1.1.2.01", StandardName: "Bancos"},
		{LocalCode: "3.1.1.01", LocalName: "Capital Social", StandardCode: "3.1.1.01", StandardName: "Capital Social"},
		{LocalCode: "4.1.1.01", LocalName: "Receita de Vendas", StandardCode: "4.1.1.01", StandardName: "Receita de Vendas"},
		{LocalCode: "5.1.1.01", LocalName: "Custo das Vendas", StandardCode: "5.1.1.01", StandardName: "Custo das Vendas"},
		{LocalCode: "6.1.1.01", LocalName: "Despesas Administrativas", StandardCode: "6.1.1.01", StandardName: "Despesas Administrativas"},
		{LocalCode: "6.1.2.01", LocalName: "Despesas Comerciais", StandardCode: "6.1.2.01", StandardName: "Despesas Comerciais"},
		{LocalCode: "6.1.3.01", LocalName: "Despesas Financeiras", StandardCode: "6.1.3.01", StandardName: "Despesas Financeiras"},
		{LocalCode: "7.1.1.01", LocalName: "Outras Receitas", StandardCode: "7.1.1.01", StandardName: "Outras Receitas"},
		{LocalCode: "7.1.2.01", LocalName: "Outras Despesas", StandardCode: "7.1.2.01", StandardName: "Outras Despesas"},
	}

	data := map[string]interface{}{
		"Title":    "Painel do Contador Social",
		"Period":   period,
		"Entities": entities,
		"Mappings": defaultMappings,
	}

	// Carregar template cache-proof
	// Tentar múltiplos caminhos para funcionar em diferentes ambientes
	var tmpl *template.Template
	var err error

	// Tentativa 1: Caminho relativo (quando executado de modules/ui_web/)
	tmpl, err = template.ParseFiles("templates/accountant_dashboard_simple.html")
	if err != nil {
		// Tentativa 2: Caminho absoluto do projeto
		tmpl, err = template.ParseFiles("modules/ui_web/templates/accountant_dashboard_simple.html")
		if err != nil {
			// Tentativa 3: Caminho relativo alternativo
			tmpl, err = template.ParseFiles("../../templates/accountant_dashboard_simple.html")
			if err != nil {
				http.Error(w, fmt.Sprintf("template error: não foi possível carregar o template: %v", err), http.StatusInternalServerError)
				return
			}
		}
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, fmt.Sprintf("template error: %v", err), http.StatusInternalServerError)
	}
}

func (h *AccountantHandler) ExportFiscal(w http.ResponseWriter, r *http.Request) {
	// Extrair entity_id e period da URL (padrão: /accountant/export/{entity_id}/{period})
	path := r.URL.Path
	var entityID, period string

	// Parse simples da URL - em produção usaríamos um router como mux
	if len(path) > len("/accountant/export/") {
		remaining := path[len("/accountant/export/"):]
		// Encontrar a próxima barra para separar entity_id e period
		for i, char := range remaining {
			if char == '/' {
				entityID = remaining[:i]
				period = remaining[i+1:]
				break
			}
		}
	}

	// Fallback para query parameters se não encontrou na URL
	if entityID == "" || period == "" {
		entityID = r.URL.Query().Get("entity_id")
		period = r.URL.Query().Get("period")
	}

	if entityID == "" || period == "" {
		http.Error(w, "entity_id and period are required", http.StatusBadRequest)
		return
	}

	// Create a repository for this specific entity with Read-Only access
	// O parâmetro ?mode=ro garante acesso somente leitura
	repo, err := h.repoFactory.NewRepository(entityID + "?mode=ro")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to access entity database (read-only): %v", err), http.StatusInternalServerError)
		return
	}

	// Create a service with the entity-specific repository
	entityService := dashboard.NewDashboardService(repo)

	batch, data, err := entityService.TranslateAndExport(r.Context(), entityID, period)
	if err != nil {
		http.Error(w, fmt.Sprintf("export failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=fiscal_%s_%s.csv", entityID, period))
	w.Header().Set("X-Export-Hash", batch.ExportHash)
	w.Write(data)
}
