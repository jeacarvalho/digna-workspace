package handler

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
	tmpl "github.com/providentia/digna/ui_web/internal/template"
)

// BaseHandler fornece funcionalidades comuns para todos os handlers
type BaseHandler struct {
	lifecycleManager lifecycle.LifecycleManager
	templateManager  *tmpl.TemplateManager
}

// NewBaseHandler cria um novo handler base
func NewBaseHandler(lm lifecycle.LifecycleManager, devMode bool) *BaseHandler {
	// Criar template manager que carrega templates diretamente do disco
	tm := tmpl.NewTemplateManager("templates", devMode)

	// Adicionar funções comuns de template
	tm.AddFunc("formatCurrency", func(amount int64) string {
		return fmt.Sprintf("R$ %.2f", float64(amount)/100)
	})

	tm.AddFunc("formatDate", func(t interface{}) string {
		switch v := t.(type) {
		case string:
			return v
		default:
			return fmt.Sprintf("%v", t)
		}
	})

	tm.AddFunc("divide", func(a, b int64) float64 {
		if b == 0 {
			return 0
		}
		return float64(a) / float64(b)
	})

	tm.AddFunc("multiply", func(a, b int64) int64 {
		return a * b
	})

	tm.AddFunc("getAlertStatusLabel", func(status string) string {
		switch status {
		case "SAFE":
			return "Dentro do planejado"
		case "WARNING":
			return "Atenção: perto do limite"
		case "EXCEEDED":
			return "Ultrapassou o planejado"
		default:
			return "Sem dados"
		}
	})

	tm.AddFunc("getAlertStatusClass", func(status string) string {
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
	})

	tm.AddFunc("getCategoryLabel", func(category string) string {
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
	})

	// fdiv - divisão de float64 (usada nos templates para cálculos monetários)
	tm.AddFunc("fdiv", func(a, b float64) float64 {
		if b == 0 {
			return 0
		}
		return a / b
	})

	// Funções específicas para supply
	tm.AddFunc("stockItemTypeLabel", func(itemType string) string {
		switch itemType {
		case "INSUMO":
			return "Insumo/Matéria-prima"
		case "PRODUTO":
			return "Produto Acabado"
		case "MERCADORIA":
			return "Mercadoria para Revenda"
		default:
			return itemType
		}
	})

	tm.AddFunc("stockItemUnitLabel", func(unit string) string {
		switch unit {
		case "UNIDADE":
			return "unid."
		case "KG":
			return "kg"
		case "G":
			return "g"
		case "L":
			return "L"
		case "M":
			return "m"
		case "CM":
			return "cm"
		case "PACOTE":
			return "pct"
		case "CAIXA":
			return "cx"
		case "SACO":
			return "sc"
		default:
			return unit
		}
	})

	tm.AddFunc("isBelowMinimum", func(quantity, minQuantity int) bool {
		return quantity < minQuantity
	})

	// Pré-carregar templates
	if err := tm.PreloadTemplates(); err != nil {
		fmt.Printf("[WARNING] Failed to preload templates: %v\n", err)
	}

	return &BaseHandler{
		lifecycleManager: lm,
		templateManager:  tm,
	}
}

// RenderTemplate renderiza um template usando o TemplateManager
func (h *BaseHandler) RenderTemplate(w http.ResponseWriter, templateName string, data map[string]interface{}) {
	content, err := h.templateManager.ExecuteTemplate(templateName, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to render template: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(content))
}

// RenderTemplateWithLayout renderiza um template dentro do layout
func (h *BaseHandler) RenderTemplateWithLayout(w http.ResponseWriter, templateName string, data map[string]interface{}) {
	// Garantir que temos um título
	if data == nil {
		data = make(map[string]interface{})
	}

	if _, exists := data["Title"]; !exists {
		data["Title"] = "Digna"
	}

	// Primeiro renderizar o conteúdo
	content, err := h.templateManager.ExecuteTemplate(templateName, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to render content template: %v", err), http.StatusInternalServerError)
		return
	}

	// Adicionar conteúdo aos dados do layout
	layoutData := make(map[string]interface{})
	for k, v := range data {
		layoutData[k] = v
	}
	layoutData["Content"] = template.HTML(content)

	// Renderizar layout
	h.RenderTemplate(w, "layout.html", layoutData)
}
