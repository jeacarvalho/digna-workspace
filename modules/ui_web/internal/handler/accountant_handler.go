package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"digna/accountant_dashboard/pkg/dashboard"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

type AccountantHandler struct {
	lifecycleManager lifecycle.LifecycleManager
	dashboardService dashboard.DashboardService
	repoFactory      dashboard.RepositoryFactory
	tmpl             *template.Template
}

func NewAccountantHandler(lm lifecycle.LifecycleManager) (*AccountantHandler, error) {
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

	tmpl := template.Must(template.New("accountant").Parse(accountantTemplate))

	return &AccountantHandler{
		lifecycleManager: lm,
		dashboardService: dashboardService,
		repoFactory:      repoFactory,
		tmpl:             tmpl,
	}, nil
}

func (h *AccountantHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/accountant/dashboard", h.Dashboard)
	mux.HandleFunc("/accountant/export", h.ExportFiscal)
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

	if err := h.tmpl.Execute(w, data); err != nil {
		http.Error(w, fmt.Sprintf("template error: %v", err), http.StatusInternalServerError)
	}
}

func (h *AccountantHandler) ExportFiscal(w http.ResponseWriter, r *http.Request) {
	entityID := r.URL.Query().Get("entity_id")
	period := r.URL.Query().Get("period")

	if entityID == "" || period == "" {
		http.Error(w, "entity_id and period are required", http.StatusBadRequest)
		return
	}

	// Create a repository for this specific entity
	repo, err := h.repoFactory.NewRepository(entityID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to access entity database: %v", err), http.StatusInternalServerError)
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

const accountantTemplate = `
<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-100 min-h-screen">
    <div class="container mx-auto px-4 py-8">
        <header class="mb-8">
            <h1 class="text-3xl font-bold text-gray-800">{{.Title}}</h1>
            <p class="text-gray-600 mt-2">Painel Multi-tenant para Contadores Sociais - Digna</p>
            <a href="/" class="text-blue-600 hover:underline mt-2 inline-block">← Voltar ao Digna</a>
        </header>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
            <div class="bg-white rounded-lg shadow p-6">
                <h3 class="text-lg font-semibold text-gray-700">Empreendimentos</h3>
                <p class="text-3xl font-bold text-blue-600">{{len .Entities}}</p>
            </div>
            <div class="bg-white rounded-lg shadow p-6">
                <h3 class="text-lg font-semibold text-gray-700">Período</h3>
                <p class="text-xl font-medium text-gray-800">{{.Period}}</p>
            </div>
            <div class="bg-white rounded-lg shadow p-6">
                <h3 class="text-lg font-semibold text-gray-700">Mapeamentos</h3>
                <p class="text-3xl font-bold text-green-600">{{len .Mappings}}</p>
            </div>
        </div>

        <section class="bg-white rounded-lg shadow mb-8">
            <div class="p-6 border-b">
                <h2 class="text-xl font-semibold text-gray-800">Entidades com Fechamento Pendente</h2>
            </div>
            <div class="p-6">
                <form method="get" class="mb-4">
                    <label class="block text-sm font-medium text-gray-700 mb-2">Selecione o Período:</label>
                    <input type="month" name="period" value="{{.Period}}" 
                           class="border rounded px-3 py-2 mr-2">
                    <button type="submit" class="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700">
                        Filtrar
                    </button>
                </form>

                {{if .Entities}}
                <table class="w-full">
                    <thead>
                        <tr class="border-b">
                            <th class="text-left py-3 px-4">ID</th>
                            <th class="text-left py-3 px-4">Nome</th>
                            <th class="text-left py-3 px-4">Status</th>
                            <th class="text-left py-3 px-4">Exportações</th>
                            <th class="text-left py-3 px-4">Ação</th>
                        </tr>
                    </thead>
                    <tbody>
                        {{range .Entities}}
                        <tr class="border-b hover:bg-gray-50">
                            <td class="py-3 px-4 font-mono">{{.ID}}</td>
                            <td class="py-3 px-4">{{.Name}}</td>
                            <td class="py-3 px-4">
                                 <span class="px-2 py-1 rounded text-sm {{if eq .Status "PENDING"}}bg-yellow-100 text-yellow-800{{else}}bg-green-100 text-green-800{{end}}">
                                    {{.Status}}
                                </span>
                            </td>
                            <td class="py-3 px-4">
                                {{if .HasExports}}
                                <span class="text-green-600">✓ Exportado</span>
                                {{else}}
                                <span class="text-red-600">Pendente</span>
                                {{end}}
                            </td>
                            <td class="py-3 px-4">
                                <a href="/accountant/export?entity_id={{.ID}}&period={{$.Period}}"
                                   class="inline-block bg-blue-600 text-white px-3 py-1 rounded text-sm hover:bg-blue-700">
                                    Exportar SPED
                                </a>
                            </td>
                        </tr>
                        {{end}}
                    </tbody>
                </table>
                {{else}}
                <p class="text-gray-500 text-center py-4">Nenhuma entidade com fechamento pendente para este período.</p>
                {{end}}
            </div>
        </section>

        <section class="bg-white rounded-lg shadow">
            <div class="p-6 border-b">
                <h2 class="text-xl font-semibold text-gray-800">Mapeamento de Contas (Plano de Contas Referencial)</h2>
            </div>
            <div class="p-6">
                <table class="w-full">
                    <thead>
                        <tr class="border-b">
                            <th class="text-left py-2 px-4">Código Local</th>
                            <th class="text-left py-2 px-4">Nome Local</th>
                            <th class="text-left py-2 px-4">Código Padrão</th>
                            <th class="text-left py-2 px-4">Nome Padrão</th>
                        </tr>
                    </thead>
                    <tbody>
                        {{range .Mappings}}
                        <tr class="border-b hover:bg-gray-50">
                            <td class="py-2 px-4 font-mono text-sm">{{.LocalCode}}</td>
                            <td class="py-2 px-4">{{.LocalName}}</td>
                            <td class="py-2 px-4 font-mono text-sm">{{.StandardCode}}</td>
                            <td class="py-2 px-4">{{.StandardName}}</td>
                        </tr>
                        {{end}}
                    </tbody>
                </table>
            </div>
        </section>
    </div>
</body>
</html>
`
