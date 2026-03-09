package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/providentia/digna/pdv_ui/internal/service"
)

type PricingHandler struct {
	tmpl *template.Template
}

func NewPricingHandler() (*PricingHandler, error) {
	// Carregar template do arquivo
	tmpl, err := template.ParseFiles("modules/pdv_ui/templates/components/pricing_calculator.html")
	if err != nil {
		// Fallback para template embutido se o arquivo não existir
		tmpl = template.Must(template.New("pricing_calculator").Parse(pricingCalculatorHTML))
	}

	return &PricingHandler{
		tmpl: tmpl,
	}, nil
}

// HandleCalculatePrice processa requisições HTMX para cálculo de preço
func (h *PricingHandler) HandleCalculatePrice(w http.ResponseWriter, r *http.Request) {
	// Extrair e validar parâmetros da query
	materialCostStr := r.URL.Query().Get("material_cost")
	laborMinutesStr := r.URL.Query().Get("labor_minutes")
	laborRateStr := r.URL.Query().Get("labor_rate")

	// Converter para int64 (centavos e minutos)
	materialCost, _ := strconv.ParseInt(materialCostStr, 10, 64)
	laborMinutes, _ := strconv.ParseInt(laborMinutesStr, 10, 64)
	laborRate, _ := strconv.ParseInt(laborRateStr, 10, 64)

	// Validar valores mínimos
	if materialCost < 0 {
		materialCost = 0
	}
	if laborMinutes < 0 {
		laborMinutes = 0
	}
	if laborRate < 0 {
		laborRate = 0
	}

	// Calcular preço usando o serviço
	calculation := service.CalculatePrice(materialCost, laborMinutes, laborRate)

	// Preparar dados para o template
	data := PricingTemplateData{
		MaterialCost:            calculation.MaterialCost,
		MaterialCostFormatted:   service.FormatToReais(calculation.MaterialCost),
		LaborMinutes:            calculation.LaborMinutes,
		LaborRate:               calculation.LaborRate,
		LaborRateFormatted:      service.FormatToReais(calculation.LaborRate),
		LaborValue:              calculation.LaborValue,
		LaborValueFormatted:     service.FormatToReais(calculation.LaborValue),
		FatesReserve:            calculation.FatesReserve,
		FatesReserveFormatted:   service.FormatToReais(calculation.FatesReserve),
		SuggestedPrice:          calculation.SuggestedPrice,
		SuggestedPriceFormatted: service.FormatToReais(calculation.SuggestedPrice),
	}

	// Calcular porcentagens para o gráfico visual
	materialPct, laborPct, reservePct := calculation.GetVisualBreakdown()
	data.MaterialPercentage = float64(materialPct) / 100.0
	data.LaborPercentage = float64(laborPct) / 100.0
	data.ReservePercentage = float64(reservePct) / 100.0

	// Renderizar template
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.tmpl.Execute(w, data); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao renderizar template: %v", err), http.StatusInternalServerError)
		return
	}
}

// PricingTemplateData contém os dados para o template
type PricingTemplateData struct {
	MaterialCost            int64
	MaterialCostFormatted   string
	LaborMinutes            int64
	LaborRate               int64
	LaborRateFormatted      string
	LaborValue              int64
	LaborValueFormatted     string
	FatesReserve            int64
	FatesReserveFormatted   string
	SuggestedPrice          int64
	SuggestedPriceFormatted string

	// Porcentagens para o gráfico (em float apenas para display)
	MaterialPercentage float64
	LaborPercentage    float64
	ReservePercentage  float64
}

// HTML do componente (será movido para arquivo separado depois)
const pricingCalculatorHTML = `
<div class="bg-white rounded-lg shadow p-6 mb-6">
	<h3 class="text-xl font-semibold text-gray-800 mb-4">Calculadora de Preço Justo</h3>
	
	<div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
		<div>
			<label class="block text-sm font-medium text-gray-700 mb-2">
				Custo do Material (R$)
			</label>
			<input type="number" 
				   name="material_cost" 
				   value="{{.MaterialCost}}"
				   hx-get="/pdv/pricing/calculate" 
				   hx-trigger="keyup changed delay:500ms"
				   hx-target="#pricing-result"
				   class="w-full border rounded px-3 py-2"
				   placeholder="Ex: 1000 (R$ 10,00)">
		</div>
		
		<div>
			<label class="block text-sm font-medium text-gray-700 mb-2">
				Seu Tempo (minutos)
			</label>
			<input type="number" 
				   name="labor_minutes" 
				   value="{{.LaborMinutes}}"
				   hx-get="/pdv/pricing/calculate" 
				   hx-trigger="keyup changed delay:500ms"
				   hx-target="#pricing-result"
				   class="w-full border rounded px-3 py-2"
				   placeholder="Ex: 60 (1 hora)">
		</div>
		
		<div>
			<label class="block text-sm font-medium text-gray-700 mb-2">
				Valor da Hora (R$)
			</label>
			<input type="number" 
				   name="labor_rate" 
				   value="{{.LaborRate}}"
				   hx-get="/pdv/pricing/calculate" 
				   hx-trigger="keyup changed delay:500ms"
				   hx-target="#pricing-result"
				   class="w-full border rounded px-3 py-2"
				   placeholder="Ex: 2000 (R$ 20,00/hora)">
		</div>
	</div>

	<div id="pricing-result">
		{{if gt .SuggestedPrice 0}}
		<div class="border-t pt-6">
			<h4 class="text-lg font-medium text-gray-800 mb-4">Preço Sugerido: <span class="text-green-600">{{.SuggestedPriceFormatted}}</span></h4>
			
			<!-- Gráfico visual -->
			<div class="mb-6">
				<div class="flex h-8 rounded-lg overflow-hidden mb-2">
					<div class="bg-blue-500" style="width: {{.MaterialPercentage}}%"></div>
					<div class="bg-green-500" style="width: {{.LaborPercentage}}%"></div>
					<div class="bg-yellow-500" style="width: {{.ReservePercentage}}%"></div>
				</div>
				
				<div class="flex justify-between text-sm">
					<div class="flex items-center">
						<div class="w-3 h-3 bg-blue-500 rounded mr-2"></div>
						<span>Material: {{.MaterialCostFormatted}}</span>
					</div>
					<div class="flex items-center">
						<div class="w-3 h-3 bg-green-500 rounded mr-2"></div>
						<span>Seu Trabalho: {{.LaborValueFormatted}}</span>
					</div>
					<div class="flex items-center">
						<div class="w-3 h-3 bg-yellow-500 rounded mr-2"></div>
						<span>Fundo da Cooperativa: {{.FatesReserveFormatted}}</span>
					</div>
				</div>
			</div>
			
			<!-- Detalhamento -->
			<div class="bg-gray-50 rounded p-4">
				<h5 class="font-medium text-gray-700 mb-2">Como chegamos nesse valor:</h5>
				<ul class="space-y-1 text-sm text-gray-600">
					<li>• Custo do material: <span class="font-medium">{{.MaterialCostFormatted}}</span></li>
					<li>• Valor do seu trabalho ({{.LaborMinutes}} minutos × {{.LaborRateFormatted}}/hora): <span class="font-medium">{{.LaborValueFormatted}}</span></li>
					<li>• Contribuição para o fundo da cooperativa (5%): <span class="font-medium">{{.FatesReserveFormatted}}</span></li>
					<li class="pt-2 border-t font-medium">= Preço justo sugerido: <span class="text-green-600">{{.SuggestedPriceFormatted}}</span></li>
				</ul>
			</div>
			
			<div class="mt-4 text-sm text-gray-500">
				<p>💡 <strong>Dica importante:</strong> Este preço garante que seu trabalho seja valorizado e que sua cooperativa tenha recursos para crescer!</p>
			</div>
		</div>
		{{else}}
		<div class="text-center py-8 text-gray-500">
			<p>Preencha os campos acima para calcular o preço justo do seu produto.</p>
			<p class="text-sm mt-2">💡 Lembre-se: seu tempo e trabalho têm valor!</p>
		</div>
		{{end}}
	</div>
</div>
`
