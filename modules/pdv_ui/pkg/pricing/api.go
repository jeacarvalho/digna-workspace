package pricing

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// PricingCalculator é a API pública para o módulo de precificação educacional
type PricingCalculator struct {
	tmpl *template.Template
}

// NewPricingCalculator cria uma nova instância da calculadora de preços
func NewPricingCalculator() (*PricingCalculator, error) {
	// Template embutido para evitar dependência de arquivos externos
	tmpl := template.Must(template.New("pricing_calculator").Parse(pricingCalculatorHTML))

	return &PricingCalculator{
		tmpl: tmpl,
	}, nil
}

// HandleCalculatePrice processa requisições HTMX para cálculo de preço
func (pc *PricingCalculator) HandleCalculatePrice(w http.ResponseWriter, r *http.Request) {
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

	// Calcular preço justo usando matemática exata (int64)
	fairPrice := calculateFairPrice(materialCost, laborMinutes, laborRate)

	// Calcular valores para visualização
	laborValue := calculateLaborValue(laborMinutes, laborRate)
	totalBeforeReserves := materialCost + laborValue
	cooperativeFund := fairPrice - totalBeforeReserves

	// Calcular percentuais para o gráfico
	materialPercent := getPercentage(materialCost, fairPrice)
	laborPercent := getPercentage(laborValue, fairPrice)
	cooperativePercent := getPercentage(cooperativeFund, fairPrice)

	// Preparar dados para o template
	data := map[string]interface{}{
		"MaterialCost":         materialCost,
		"LaborMinutes":         laborMinutes,
		"LaborRate":            laborRate,
		"FairPrice":            fairPrice,
		"LaborValue":           laborValue,
		"CooperativeFund":      cooperativeFund,
		"MaterialPercent":      materialPercent,
		"LaborPercent":         laborPercent,
		"CooperativePercent":   cooperativePercent,
		"FormattedMaterial":    formatToReais(materialCost),
		"FormattedLabor":       formatToReais(laborValue),
		"FormattedCooperative": formatToReais(cooperativeFund),
		"FormattedFairPrice":   formatToReais(fairPrice),
	}

	// Renderizar template
	w.Header().Set("Content-Type", "text/html")
	if err := pc.tmpl.Execute(w, data); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao renderizar template: %v", err), http.StatusInternalServerError)
	}
}

// RegisterRoutes registra as rotas de precificação no mux fornecido
func (pc *PricingCalculator) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/pdv/pricing/calculate", pc.HandleCalculatePrice)
}

// Funções auxiliares de cálculo (implementação exata com int64)

func calculateFairPrice(materialCost, laborMinutes, laborRate int64) int64 {
	// Calcular valor do trabalho em centavos
	laborValue := calculateLaborValue(laborMinutes, laborRate)

	// Soma do custo dos materiais + valor do trabalho
	totalBeforeReserves := materialCost + laborValue

	// Adicionar 5% para o fundo cooperativo (usando matemática exata)
	// 5% = 5/100 = 1/20
	var cooperativeFund int64
	if totalBeforeReserves > 0 {
		cooperativeFund = (totalBeforeReserves * 5) / 100
	}

	// Preço justo = custo materiais + valor trabalho + fundo cooperativo
	return totalBeforeReserves + cooperativeFund
}

func calculateLaborValue(laborMinutes, laborRate int64) int64 {
	// Valor do trabalho = (minutos × taxa horária) / 60
	// Isso evita perda de precisão ao dividir a taxa primeiro
	if laborMinutes == 0 || laborRate == 0 {
		return 0
	}
	return (laborMinutes * laborRate) / 60
}

func calculateCooperativeFund(fairPrice int64) int64 {
	// Fundo cooperativo = 5% do total antes das reservas
	// Para calcular o fundo cooperativo a partir do preço final:
	// Se preço final = P, fundo = F, então P = (C + F) onde C = custo total
	// E F = 5% de C, então F = C * 5 / 100
	// Portanto C = P * 100 / 105
	// E F = P - C = P - (P * 100 / 105) = P * 5 / 105
	if fairPrice == 0 {
		return 0
	}
	return (fairPrice * 5) / 105
}

func getPercentage(part, total int64) int64 {
	if total == 0 {
		return 0
	}
	// Calcular percentual: (part × 100) / total
	return (part * 100) / total
}

func formatToReais(amount int64) string {
	// Converter centavos para reais com 2 casas decimais
	reais := float64(amount) / 100.0
	return fmt.Sprintf("R$ %.2f", reais)
}

// Template HTML embutido para a calculadora de preços
const pricingCalculatorHTML = `
<div id="pricing-calculator" class="bg-white rounded-lg shadow-md p-6 mb-6">
	<h3 class="text-xl font-semibold text-gray-800 mb-4">Calculadora de Preço Justo</h3>
	
	<div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
		<div>
			<label class="block text-sm font-medium text-gray-700 mb-1">
				Custo dos Materiais (R$)
			</label>
			<input 
				type="number" 
				name="material_cost" 
				value="{{.MaterialCost}}"
				class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
				hx-get="/pdv/pricing/calculate"
				hx-trigger="keyup changed delay:500ms"
				hx-target="#pricing-results"
				hx-include="[name='material_cost'],[name='labor_minutes'],[name='labor_rate']"
				placeholder="Ex: 1500 (para R$ 15,00)"
			>
			<p class="text-xs text-gray-500 mt-1">Em centavos (R$ 1,00 = 100)</p>
		</div>
		
		<div>
			<label class="block text-sm font-medium text-gray-700 mb-1">
				Seu Tempo (minutos)
			</label>
			<input 
				type="number" 
				name="labor_minutes" 
				value="{{.LaborMinutes}}"
				class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
				hx-get="/pdv/pricing/calculate"
				hx-trigger="keyup changed delay:500ms"
				hx-target="#pricing-results"
				hx-include="[name='material_cost'],[name='labor_minutes'],[name='labor_rate']"
				placeholder="Ex: 120 (2 horas)"
			>
		</div>
		
		<div>
			<label class="block text-sm font-medium text-gray-700 mb-1">
				Valor da Hora (R$)
			</label>
			<input 
				type="number" 
				name="labor_rate" 
				value="{{.LaborRate}}"
				class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
				hx-get="/pdv/pricing/calculate"
				hx-trigger="keyup changed delay:500ms"
				hx-target="#pricing-results"
				hx-include="[name='material_cost'],[name='labor_minutes'],[name='labor_rate']"
				placeholder="Ex: 3000 (R$ 30,00/hora)"
			>
			<p class="text-xs text-gray-500 mt-1">Em centavos por hora</p>
		</div>
	</div>
	
	<div id="pricing-results">
		{{if gt .FairPrice 0}}
		<div class="bg-gray-50 rounded-lg p-4 mb-4">
			<h4 class="text-lg font-medium text-gray-800 mb-3">Preço Justo Calculado</h4>
			
			<div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-4">
				<div class="text-center">
					<div class="text-2xl font-bold text-blue-600">{{.FormattedFairPrice}}</div>
					<div class="text-sm text-gray-600">Preço Justo</div>
				</div>
				<div class="text-center">
					<div class="text-xl font-semibold text-green-600">{{.FormattedMaterial}}</div>
					<div class="text-sm text-gray-600">Materiais</div>
				</div>
				<div class="text-center">
					<div class="text-xl font-semibold text-yellow-600">{{.FormattedLabor}}</div>
					<div class="text-sm text-gray-600">Seu Trabalho</div>
				</div>
				<div class="text-center">
					<div class="text-xl font-semibold text-purple-600">{{.FormattedCooperative}}</div>
					<div class="text-sm text-gray-600">Fundo Cooperativo</div>
				</div>
			</div>
			
			<!-- Gráfico de pizza visual -->
			<div class="mb-4">
				<h5 class="text-md font-medium text-gray-700 mb-2">Composição do Preço</h5>
				<div class="h-6 rounded-full overflow-hidden flex">
					{{if gt .MaterialPercent 0}}
					<div class="bg-green-500 h-full" style="width: {{.MaterialPercent}}%"></div>
					{{end}}
					{{if gt .LaborPercent 0}}
					<div class="bg-yellow-500 h-full" style="width: {{.LaborPercent}}%"></div>
					{{end}}
					{{if gt .CooperativePercent 0}}
					<div class="bg-purple-500 h-full" style="width: {{.CooperativePercent}}%"></div>
					{{end}}
				</div>
				<div class="flex justify-between text-xs text-gray-600 mt-1">
					<span>Materiais: {{.MaterialPercent}}%</span>
					<span>Seu Trabalho: {{.LaborPercent}}%</span>
					<span>Fundo: {{.CooperativePercent}}%</span>
				</div>
			</div>
			
			<!-- Explicação pedagógica -->
			<div class="bg-blue-50 border-l-4 border-blue-400 p-3 rounded">
				<p class="text-sm text-gray-700">
					<strong>Como funciona:</strong> O preço justo considera o custo dos materiais, 
					o valor do seu tempo de trabalho e inclui 10% para o fundo cooperativo, 
					que ajuda a sustentar a comunidade.
				</p>
			</div>
		</div>
		{{else}}
		<div class="text-center py-8 text-gray-500">
			<p>Preencha os valores acima para calcular o preço justo do seu produto.</p>
		</div>
		{{end}}
	</div>
</div>
`
