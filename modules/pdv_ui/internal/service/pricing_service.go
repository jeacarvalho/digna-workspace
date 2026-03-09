package service

import "fmt"

// PricingCalculation representa o resultado do cálculo pedagógico de precificação
type PricingCalculation struct {
	MaterialCost   int64 // centavos (custo do material)
	LaborMinutes   int64 // minutos de trabalho
	LaborRate      int64 // valor de 1 hora de trabalho em centavos
	LaborValue     int64 // LaborMinutes convertido em valor (centavos)
	FatesReserve   int64 // 5% do total para reserva legal/educacional (FATES)
	SuggestedPrice int64 // Soma de tudo (centavos)
}

// CalculatePrice gera a estrutura baseada na matemática exata e anti-float
func CalculatePrice(materialCost, laborMinutes, laborRatePerHour int64) PricingCalculation {
	if materialCost < 0 {
		materialCost = 0
	}
	if laborMinutes < 0 {
		laborMinutes = 0
	}
	if laborRatePerHour < 0 {
		laborRatePerHour = 0
	}

	// Calcular valor do trabalho (minutos → valor)
	// laborRatePerHour é o valor por hora em centavos
	// Para evitar float: (laborMinutes * laborRatePerHour) / 60
	var laborValue int64
	if laborMinutes > 0 && laborRatePerHour > 0 {
		laborValue = (laborMinutes * laborRatePerHour) / 60
	}

	// Calcular total antes das reservas
	totalBeforeReserves := materialCost + laborValue

	// Calcular reserva FATES (5%)
	var fatesReserve int64
	if totalBeforeReserves > 0 {
		// 5% = totalBeforeReserves * 5 / 100
		fatesReserve = (totalBeforeReserves * 5) / 100
	}

	// Calcular preço sugerido (total + reservas)
	suggestedPrice := totalBeforeReserves + fatesReserve

	return PricingCalculation{
		MaterialCost:   materialCost,
		LaborMinutes:   laborMinutes,
		LaborRate:      laborRatePerHour,
		LaborValue:     laborValue,
		FatesReserve:   fatesReserve,
		SuggestedPrice: suggestedPrice,
	}
}

// FormatToReais formata centavos para string em reais (R$ X,XX)
func FormatToReais(centavos int64) string {
	if centavos < 0 {
		centavos = -centavos
		return fmt.Sprintf("-R$ %d,%02d", centavos/100, centavos%100)
	}
	return fmt.Sprintf("R$ %d,%02d", centavos/100, centavos%100)
}

// GetPercentage calcula porcentagem sem usar float
func GetPercentage(part, total int64) int64 {
	if total == 0 {
		return 0
	}
	// Retorna porcentagem * 100 (para duas casas decimais)
	return (part * 10000) / total
}

// GetVisualBreakdown retorna a divisão visual para o gráfico
func (p *PricingCalculation) GetVisualBreakdown() (materialPct, laborPct, reservePct int64) {
	if p.SuggestedPrice == 0 {
		return 10000, 0, 0 // 100% material quando não há preço
	}

	materialPct = GetPercentage(p.MaterialCost, p.SuggestedPrice)
	laborPct = GetPercentage(p.LaborValue, p.SuggestedPrice)
	reservePct = GetPercentage(p.FatesReserve, p.SuggestedPrice)

	// Ajustar para somar 100% (pode haver arredondamento)
	totalPct := materialPct + laborPct + reservePct
	if totalPct != 10000 { // 100.00%
		// Ajustar a maior parte
		if materialPct >= laborPct && materialPct >= reservePct {
			materialPct += (10000 - totalPct)
		} else if laborPct >= materialPct && laborPct >= reservePct {
			laborPct += (10000 - totalPct)
		} else {
			reservePct += (10000 - totalPct)
		}
	}

	return materialPct, laborPct, reservePct
}
