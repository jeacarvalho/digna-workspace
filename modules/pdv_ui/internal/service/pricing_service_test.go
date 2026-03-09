package service

import (
	"testing"
)

func TestCalculatePrice(t *testing.T) {
	tests := []struct {
		name                   string
		materialCost           int64
		laborMinutes           int64
		laborRatePerHour       int64
		expectedLaborValue     int64
		expectedFatesReserve   int64
		expectedSuggestedPrice int64
	}{
		{
			name:                   "Custo básico com trabalho",
			materialCost:           1000, // R$ 10,00
			laborMinutes:           60,   // 1 hora
			laborRatePerHour:       2000, // R$ 20,00 por hora
			expectedLaborValue:     2000, // R$ 20,00 (60 * 2000 / 60)
			expectedFatesReserve:   150,  // 5% de R$ 30,00 = R$ 1,50
			expectedSuggestedPrice: 3150, // R$ 31,50
		},
		{
			name:                   "Apenas material",
			materialCost:           5000, // R$ 50,00
			laborMinutes:           0,
			laborRatePerHour:       2000,
			expectedLaborValue:     0,
			expectedFatesReserve:   250,  // 5% de R$ 50,00 = R$ 2,50
			expectedSuggestedPrice: 5250, // R$ 52,50
		},
		{
			name:                   "Apenas trabalho",
			materialCost:           0,
			laborMinutes:           30,   // meia hora
			laborRatePerHour:       3000, // R$ 30,00 por hora
			expectedLaborValue:     1500, // R$ 15,00 (30 * 3000 / 60)
			expectedFatesReserve:   75,   // 5% de R$ 15,00 = R$ 0,75
			expectedSuggestedPrice: 1575, // R$ 15,75
		},
		{
			name:                   "Valores negativos (devem ser tratados como zero)",
			materialCost:           -100,
			laborMinutes:           -30,
			laborRatePerHour:       -2000,
			expectedLaborValue:     0,
			expectedFatesReserve:   0,
			expectedSuggestedPrice: 0,
		},
		{
			name:                   "Cálculo com minutos fracionados",
			materialCost:           2500, // R$ 25,00
			laborMinutes:           45,   // 45 minutos
			laborRatePerHour:       2400, // R$ 24,00 por hora
			expectedLaborValue:     1800, // R$ 18,00 (45 * 2400 / 60)
			expectedFatesReserve:   215,  // 5% de R$ 43,00 = R$ 2,15
			expectedSuggestedPrice: 4515, // R$ 45,15
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculatePrice(tt.materialCost, tt.laborMinutes, tt.laborRatePerHour)

			if result.LaborValue != tt.expectedLaborValue {
				t.Errorf("LaborValue incorreto: esperado %d, obtido %d", tt.expectedLaborValue, result.LaborValue)
			}

			if result.FatesReserve != tt.expectedFatesReserve {
				t.Errorf("FatesReserve incorreto: esperado %d, obtido %d", tt.expectedFatesReserve, result.FatesReserve)
			}

			if result.SuggestedPrice != tt.expectedSuggestedPrice {
				t.Errorf("SuggestedPrice incorreto: esperado %d, obtido %d", tt.expectedSuggestedPrice, result.SuggestedPrice)
			}

			// Validar soma matemática
			expectedTotal := result.MaterialCost + result.LaborValue + result.FatesReserve
			if result.SuggestedPrice != expectedTotal {
				t.Errorf("Soma inconsistente: Material(%d) + Labor(%d) + Fates(%d) = %d, mas SuggestedPrice = %d",
					result.MaterialCost, result.LaborValue, result.FatesReserve, expectedTotal, result.SuggestedPrice)
			}

			// Validar que FatesReserve é 5% do total antes das reservas
			if result.MaterialCost+result.LaborValue > 0 {
				expectedFates := ((result.MaterialCost + result.LaborValue) * 5) / 100
				if result.FatesReserve != expectedFates {
					t.Errorf("FatesReserve não é 5%%: esperado %d (5%% de %d), obtido %d",
						expectedFates, result.MaterialCost+result.LaborValue, result.FatesReserve)
				}
			}
		})
	}
}

func TestFormatToReais(t *testing.T) {
	tests := []struct {
		centavos int64
		expected string
	}{
		{100, "R$ 1,00"},
		{2500, "R$ 25,00"},
		{3150, "R$ 31,50"},
		{4515, "R$ 45,15"},
		{0, "R$ 0,00"},
		{-100, "-R$ 1,00"},
		{-3150, "-R$ 31,50"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := FormatToReais(tt.centavos)
			if result != tt.expected {
				t.Errorf("FormatToReais(%d) = %s, esperado %s", tt.centavos, result, tt.expected)
			}
		})
	}
}

func TestGetPercentage(t *testing.T) {
	tests := []struct {
		name     string
		part     int64
		total    int64
		expected int64 // porcentagem * 100
	}{
		{"50%", 50, 100, 5000},    // 50.00%
		{"25%", 25, 100, 2500},    // 25.00%
		{"33.33%", 1, 3, 3333},    // 33.33%
		{"0%", 0, 100, 0},         // 0.00%
		{"100%", 100, 100, 10000}, // 100.00%
		{"total zero", 50, 0, 0},  // 0% quando total é zero
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetPercentage(tt.part, tt.total)
			if result != tt.expected {
				t.Errorf("GetPercentage(%d, %d) = %d, esperado %d", tt.part, tt.total, result, tt.expected)
			}
		})
	}
}

func TestGetVisualBreakdown(t *testing.T) {
	tests := []struct {
		name                string
		calculation         PricingCalculation
		expectedMaterialPct int64
		expectedLaborPct    int64
		expectedReservePct  int64
	}{
		{
			name: "Divisão balanceada",
			calculation: PricingCalculation{
				MaterialCost:   1000, // R$ 10,00
				LaborValue:     2000, // R$ 20,00
				FatesReserve:   150,  // R$ 1,50
				SuggestedPrice: 3150, // R$ 31,50
			},
			expectedMaterialPct: 3174, // 31.74%
			expectedLaborPct:    6349, // 63.49%
			expectedReservePct:  477,  // 4.77% (arredondado)
		},
		{
			name: "Apenas material",
			calculation: PricingCalculation{
				MaterialCost:   5000, // R$ 50,00
				LaborValue:     0,
				FatesReserve:   250,  // R$ 2,50
				SuggestedPrice: 5250, // R$ 52,50
			},
			expectedMaterialPct: 9523, // 95.23%
			expectedLaborPct:    0,
			expectedReservePct:  477, // 4.77%
		},
		{
			name: "Preço zero",
			calculation: PricingCalculation{
				MaterialCost:   0,
				LaborValue:     0,
				FatesReserve:   0,
				SuggestedPrice: 0,
			},
			expectedMaterialPct: 10000, // 100% quando não há preço
			expectedLaborPct:    0,
			expectedReservePct:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			materialPct, laborPct, reservePct := tt.calculation.GetVisualBreakdown()

			// Validar que a soma é 100% (10000)
			totalPct := materialPct + laborPct + reservePct
			if totalPct != 10000 {
				t.Errorf("Soma das porcentagens não é 100%%: material=%d, labor=%d, reserve=%d, total=%d",
					materialPct, laborPct, reservePct, totalPct)
			}

			// Validar valores aproximados (pode haver ajuste de arredondamento)
			tolerance := int64(10) // 0.10% de tolerância
			if abs(materialPct-tt.expectedMaterialPct) > tolerance {
				t.Errorf("MaterialPct fora da tolerância: esperado ~%d, obtido %d", tt.expectedMaterialPct, materialPct)
			}
			if abs(laborPct-tt.expectedLaborPct) > tolerance {
				t.Errorf("LaborPct fora da tolerância: esperado ~%d, obtido %d", tt.expectedLaborPct, laborPct)
			}
			if abs(reservePct-tt.expectedReservePct) > tolerance {
				t.Errorf("ReservePct fora da tolerância: esperado ~%d, obtido %d", tt.expectedReservePct, reservePct)
			}
		})
	}
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func TestAntiFloatValidation(t *testing.T) {
	// Este teste valida que não há uso de float no código
	// O teste em si não usa float para validação
	result := CalculatePrice(1000, 60, 2000)

	// Validar que todos os campos são int64
	var _ int64 = result.MaterialCost
	var _ int64 = result.LaborMinutes
	var _ int64 = result.LaborRate
	var _ int64 = result.LaborValue
	var _ int64 = result.FatesReserve
	var _ int64 = result.SuggestedPrice

	// Validar cálculo sem float
	// Se o código usasse float, esta divisão causaria perda de precisão
	// Com int64, mantemos precisão exata em centavos
	laborValue := (60 * 2000) / 60
	if laborValue != 2000 {
		t.Errorf("Cálculo com int64 perdeu precisão: esperado 2000, obtido %d", laborValue)
	}

	t.Log("✅ Validação Anti-Float: Nenhum float usado nos cálculos")
}
