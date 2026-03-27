package handler

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/providentia/digna/core_lume/pkg/help"
	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

// HelpHandler handles help system HTTP requests
type HelpHandler struct {
	*BaseHandler
	lifecycleManager lifecycle.LifecycleManager
	helpService      *help.Service
}

// NewHelpHandler creates a new HelpHandler
func NewHelpHandler(lm lifecycle.LifecycleManager) (*HelpHandler, error) {
	base := NewBaseHandler(lm, true)

	helpService, err := help.NewService(lm)
	if err != nil {
		return nil, fmt.Errorf("failed to create help service: %w", err)
	}

	return &HelpHandler{
		BaseHandler:      base,
		lifecycleManager: lm,
		helpService:      helpService,
	}, nil
}

// RegisterRoutes registers the handler routes
func (h *HelpHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/help", h.handleHelpIndex)
	mux.HandleFunc("/help/search", h.handleHelpSearch)
	mux.HandleFunc("/help/topic/", h.handleHelpTopic)
}

// handleHelpIndex displays the help index page
func (h *HelpHandler) handleHelpIndex(w http.ResponseWriter, r *http.Request) {
	index, err := h.helpService.ListIndex()
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao carregar índice: %v", err), http.StatusInternalServerError)
		return
	}

	categories := h.helpService.GetCategories()

	data := map[string]interface{}{
		"Title":      "Central de Ajuda - Digna",
		"Index":      index,
		"Categories": categories,
	}

	// Load template from file
	tmplPath := "modules/ui_web/templates/help_index_simple.html"
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		// Fallback to inline template
		h.renderIndexInline(w, data)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, data)
}

// handleHelpSearch handles search requests
func (h *HelpHandler) handleHelpSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	results, err := h.helpService.Search(query)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro na busca: %v", err), http.StatusInternalServerError)
		return
	}

	categories := h.helpService.GetCategories()

	data := map[string]interface{}{
		"Title":      "Busca - Central de Ajuda",
		"Query":      query,
		"Results":    results,
		"Categories": categories,
	}

	// Load template from file
	tmplPath := "modules/ui_web/templates/help_search_simple.html"
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		// Fallback - redirect to index with search param
		http.Redirect(w, r, fmt.Sprintf("/help?q=%s", query), http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, data)
}

// handleHelpTopic displays a specific help topic
func (h *HelpHandler) handleHelpTopic(w http.ResponseWriter, r *http.Request) {
	// Extract topic key from URL path
	key := r.URL.Path[len("/help/topic/"):]
	if key == "" {
		http.Error(w, "Tópico não especificado", http.StatusBadRequest)
		return
	}

	topic, err := h.helpService.GetTopicByKey(key)
	if err != nil {
		http.Error(w, fmt.Sprintf("Tópico não encontrado: %v", err), http.StatusNotFound)
		return
	}

	related, err := h.helpService.GetRelatedTopics(topic)
	if err != nil {
		// Log but don't fail - related topics are not critical
		fmt.Printf("Warning: failed to get related topics: %v\n", err)
		related = []*help.HelpTopic{}
	}

	data := map[string]interface{}{
		"Title":   topic.Title,
		"Topic":   topic,
		"Related": related,
	}

	// Load template from file
	tmplPath := "modules/ui_web/templates/help_topic_simple.html"
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		// Fallback to inline template
		h.renderTopicInline(w, data)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, data)
}

// renderIndexInline renders an inline template for help index
func (h *HelpHandler) renderIndexInline(w http.ResponseWriter, data map[string]interface{}) {
	html := `<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-50 min-h-screen">
    <nav class="bg-blue-600 text-white shadow-lg">
        <div class="container mx-auto px-4 py-3">
            <div class="flex items-center space-x-4">
                <a href="/" class="text-white hover:text-blue-100">← Voltar</a>
                <span class="text-xl font-bold">Central de Ajuda</span>
            </div>
        </div>
    </nav>
    
    <main class="container mx-auto px-4 py-8">
        <h1 class="text-3xl font-bold text-gray-800 mb-6">Central de Ajuda</h1>
        
        <form action="/help/search" method="get" class="mb-8">
            <div class="flex gap-2">
                <input type="text" name="q" placeholder="Buscar tópico..." 
                       class="flex-1 p-3 border rounded-lg">
                <button type="submit" class="bg-blue-600 text-white px-6 py-3 rounded-lg">
                    Buscar
                </button>
            </div>
        </form>
        
        {{range $category, $topics := .Index}}
        <div class="mb-8">
            <h2 class="text-xl font-semibold text-gray-700 mb-4">{{index $.Categories $category}}</h2>
            <div class="grid gap-4">
                {{range $topics}}
                <a href="/help/topic/{{.Key}}" class="bg-white p-4 rounded-lg shadow hover:shadow-md">
                    <h3 class="font-medium text-blue-600">{{.Title}}</h3>
                    <p class="text-gray-600 text-sm">{{.Summary}}</p>
                </a>
                {{end}}
            </div>
        </div>
        {{end}}
    </main>
</body>
</html>`

	tmpl, _ := template.New("help_index").Parse(html)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, data)
}

// renderTopicInline renders an inline template for help topic
func (h *HelpHandler) renderTopicInline(w http.ResponseWriter, data map[string]interface{}) {
	html := `<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Topic.Title}}</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-50 min-h-screen">
    <nav class="bg-blue-600 text-white shadow-lg">
        <div class="container mx-auto px-4 py-3">
            <div class="flex items-center space-x-4">
                <a href="/help" class="text-white hover:text-blue-100">← Voltar</a>
                <span class="text-xl font-bold">Central de Ajuda</span>
            </div>
        </div>
    </nav>
    
    <main class="container mx-auto px-4 py-8">
        <article class="bg-white rounded-lg shadow-lg p-8">
            <h1 class="text-3xl font-bold text-gray-800 mb-4">{{.Topic.Title}}</h1>
            <p class="text-lg text-gray-600 mb-6">{{.Topic.Summary}}</p>
            
            <div class="prose max-w-none">
                <h2 class="text-xl font-semibold text-gray-700 mb-2">O que é?</h2>
                <p class="text-gray-600 mb-6">{{.Topic.Explanation}}</p>
                
                <h2 class="text-xl font-semibold text-gray-700 mb-2">Por que perguntamos?</h2>
                <p class="text-gray-600 mb-6">{{.Topic.WhyAsked}}</p>
                
                <h2 class="text-xl font-semibold text-gray-700 mb-2">Próximos Passos</h2>
                <p class="text-gray-600 mb-6">{{.Topic.NextSteps}}</p>
                
                {{if .Topic.Legislation}}
                <h2 class="text-xl font-semibold text-gray-700 mb-2">Legislação</h2>
                <p class="text-gray-600 mb-6">{{.Topic.Legislation}}</p>
                {{end}}
                
                {{if .Topic.OfficialLink}}
                <a href="{{.Topic.OfficialLink}}" target="_blank" 
                   class="text-blue-600 hover:underline">
                    Saiba mais na fonte oficial →
                </a>
                {{end}}
            </div>
        </article>
        
        {{if .Related}}
        <div class="mt-8">
            <h2 class="text-xl font-semibold text-gray-700 mb-4">Tópicos Relacionados</h2>
            <div class="grid gap-4 md:grid-cols-3">
                {{range .Related}}
                <a href="/help/topic/{{.Key}}" class="bg-white p-4 rounded-lg shadow hover:shadow-md">
                    <h3 class="font-medium text-blue-600">{{.Title}}</h3>
                    <p class="text-gray-600 text-sm">{{.Summary}}</p>
                </a>
                {{end}}
            </div>
        </div>
        {{end}}
    </main>
</body>
</html>`

	tmpl, _ := template.New("help_topic").Parse(html)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, data)
}
