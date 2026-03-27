package domain

import (
	"testing"
)

func TestEligibilityProfile_Validate(t *testing.T) {
	tests := []struct {
		name    string
		profile EligibilityProfile
		wantErr bool
	}{
		{
			name: "Valid profile - all required fields",
			profile: EligibilityProfile{
				EntityID:          "entity-1",
				FinalidadeCredito: FinalidadeCapitalGiro,
				TipoEntidade:      TipoEntidadeMEI,
				ValorNecessario:   100000, // R$ 1.000,00
			},
			wantErr: false,
		},
		{
			name: "Missing EntityID",
			profile: EligibilityProfile{
				FinalidadeCredito: FinalidadeCapitalGiro,
				TipoEntidade:      TipoEntidadeMEI,
			},
			wantErr: true,
		},
		{
			name: "Invalid FinalidadeCredito",
			profile: EligibilityProfile{
				EntityID:          "entity-1",
				FinalidadeCredito: "INVALIDO",
				TipoEntidade:      TipoEntidadeMEI,
			},
			wantErr: true,
		},
		{
			name: "Invalid TipoEntidade",
			profile: EligibilityProfile{
				EntityID:          "entity-1",
				FinalidadeCredito: FinalidadeCapitalGiro,
				TipoEntidade:      "INVALIDO",
			},
			wantErr: true,
		},
		{
			name: "ValorNecessario zero with specific finalidade",
			profile: EligibilityProfile{
				EntityID:          "entity-1",
				FinalidadeCredito: FinalidadeCapitalGiro,
				TipoEntidade:      TipoEntidadeMEI,
				ValorNecessario:   0,
			},
			wantErr: true,
		},
		{
			name: "Valid - OUTRO finalidade allows zero valor",
			profile: EligibilityProfile{
				EntityID:          "entity-1",
				FinalidadeCredito: FinalidadeOutro,
				TipoEntidade:      TipoEntidadeMEI,
				ValorNecessario:   0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.profile.Validate()
			if tt.wantErr && err == nil {
				t.Errorf("Validate() expected error but got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Validate() unexpected error: %v", err)
			}
		})
	}
}

func TestEligibilityProfile_IsComplete(t *testing.T) {
	tests := []struct {
		name     string
		profile  EligibilityProfile
		expected bool
	}{
		{
			name: "Complete profile",
			profile: EligibilityProfile{
				EntityID:          "entity-1",
				FinalidadeCredito: FinalidadeCapitalGiro,
				TipoEntidade:      TipoEntidadeMEI,
				ValorNecessario:   100000,
			},
			expected: true,
		},
		{
			name: "Missing EntityID",
			profile: EligibilityProfile{
				FinalidadeCredito: FinalidadeCapitalGiro,
				TipoEntidade:      TipoEntidadeMEI,
				ValorNecessario:   100000,
			},
			expected: false,
		},
		{
			name: "Missing FinalidadeCredito",
			profile: EligibilityProfile{
				EntityID:        "entity-1",
				TipoEntidade:    TipoEntidadeMEI,
				ValorNecessario: 100000,
			},
			expected: false,
		},
		{
			name: "Missing TipoEntidade",
			profile: EligibilityProfile{
				EntityID:          "entity-1",
				FinalidadeCredito: FinalidadeCapitalGiro,
				ValorNecessario:   100000,
			},
			expected: false,
		},
		{
			name: "Zero ValorNecessario with specific finalidade",
			profile: EligibilityProfile{
				EntityID:          "entity-1",
				FinalidadeCredito: FinalidadeCapitalGiro,
				TipoEntidade:      TipoEntidadeMEI,
				ValorNecessario:   0,
			},
			expected: false,
		},
		{
			name: "Complete - OUTRO finalidade with zero valor",
			profile: EligibilityProfile{
				EntityID:          "entity-1",
				FinalidadeCredito: FinalidadeOutro,
				TipoEntidade:      TipoEntidadeMEI,
				ValorNecessario:   0,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.profile.IsComplete()
			if result != tt.expected {
				t.Errorf("IsComplete() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestEligibilityProfile_GetCompletionPercent(t *testing.T) {
	tests := []struct {
		name     string
		profile  EligibilityProfile
		expected float64
	}{
		{
			name:     "Empty profile",
			profile:  EligibilityProfile{},
			expected: 0.0,
		},
		{
			name: "One field filled",
			profile: EligibilityProfile{
				InscritoCadUnico: true,
			},
			expected: 14.29, // 1/7
		},
		{
			name: "Half filled",
			profile: EligibilityProfile{
				InscritoCadUnico:    true,
				SocioMulher:         true,
				InadimplenciaAtiva:  true,
				ContabilidadeFormal: true,
			},
			expected: 57.14, // 4/7
		},
		{
			name: "All fields filled",
			profile: EligibilityProfile{
				InscritoCadUnico:    true,
				SocioMulher:         true,
				InadimplenciaAtiva:  true,
				ContabilidadeFormal: true,
				FinalidadeCredito:   FinalidadeCapitalGiro,
				TipoEntidade:        TipoEntidadeMEI,
				ValorNecessario:     100000,
			},
			expected: 100.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.profile.GetCompletionPercent()
			// Allow small difference for floating point
			if result < tt.expected-0.5 || result > tt.expected+0.5 {
				t.Errorf("GetCompletionPercent() = %.2f, expected %.2f", result, tt.expected)
			}
		})
	}
}

func TestEligibilityProfile_CanEdit(t *testing.T) {
	tests := []struct {
		name     string
		userRole string
		expected bool
	}{
		{"Coordinator can edit", "COORDINATOR", true},
		{"Member cannot edit", "MEMBER", false},
		{"Advisor cannot edit", "ADVISOR", false},
		{"Empty role cannot edit", "", false},
		{"Unknown role cannot edit", "UNKNOWN", false},
	}

	profile := EligibilityProfile{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := profile.CanEdit(tt.userRole)
			if result != tt.expected {
				t.Errorf("CanEdit(%s) = %v, expected %v", tt.userRole, result, tt.expected)
			}
		})
	}
}

func TestEligibilityProfile_Update(t *testing.T) {
	profile := EligibilityProfile{
		ID:                  "profile-1",
		EntityID:            "entity-1",
		InscritoCadUnico:    false,
		SocioMulher:         false,
		InadimplenciaAtiva:  false,
		ContabilidadeFormal: false,
		ValorNecessario:     0,
	}

	trueVal := true
	falseVal := false
	finalidade := string(FinalidadeCapitalGiro)
	tipo := string(TipoEntidadeMEI)
	valor := int64(50000)

	input := EligibilityInput{
		InscritoCadUnico:    &trueVal,
		SocioMulher:         &trueVal,
		InadimplenciaAtiva:  &falseVal,
		FinalidadeCredito:   &finalidade,
		TipoEntidade:        &tipo,
		ValorNecessario:     &valor,
		ContabilidadeFormal: &trueVal,
	}

	err := profile.Update(input, "user-123")
	if err != nil {
		t.Errorf("Update() unexpected error: %v", err)
	}

	// Verify updates
	if !profile.InscritoCadUnico {
		t.Error("InscritoCadUnico should be true")
	}
	if !profile.SocioMulher {
		t.Error("SocioMulher should be true")
	}
	if profile.InadimplenciaAtiva {
		t.Error("InadimplenciaAtiva should be false")
	}
	if profile.FinalidadeCredito != FinalidadeCapitalGiro {
		t.Errorf("FinalidadeCredito = %v, expected %v", profile.FinalidadeCredito, FinalidadeCapitalGiro)
	}
	if profile.TipoEntidade != TipoEntidadeMEI {
		t.Errorf("TipoEntidade = %v, expected %v", profile.TipoEntidade, TipoEntidadeMEI)
	}
	if profile.ValorNecessario != 50000 {
		t.Errorf("ValorNecessario = %d, expected 50000", profile.ValorNecessario)
	}
	if !profile.ContabilidadeFormal {
		t.Error("ContabilidadeFormal should be true")
	}
	if profile.PreenchidoPor != "user-123" {
		t.Errorf("PreenchidoPor = %s, expected user-123", profile.PreenchidoPor)
	}
	if profile.PreenchidoEm == 0 {
		t.Error("PreenchidoEm should be set")
	}
	if profile.AtualizadoEm == 0 {
		t.Error("AtualizadoEm should be set")
	}
}

func TestEligibilityProfile_Update_Partial(t *testing.T) {
	profile := EligibilityProfile{
		ID:       "profile-1",
		EntityID: "entity-1",
	}

	// Update only one field
	trueVal := true
	input := EligibilityInput{
		InscritoCadUnico: &trueVal,
	}

	err := profile.Update(input, "user-123")
	if err != nil {
		t.Errorf("Update() unexpected error: %v", err)
	}

	if !profile.InscritoCadUnico {
		t.Error("InscritoCadUnico should be true")
	}
	// Other fields should remain zero/false
	if profile.SocioMulher {
		t.Error("SocioMulher should remain false")
	}
}

func TestEligibilityProfile_GetValorNecessarioReal(t *testing.T) {
	profile := EligibilityProfile{
		ValorNecessario: 123456, // R$ 1.234,56
	}

	result := profile.GetValorNecessarioReal()
	expected := 1234.56

	if result != expected {
		t.Errorf("GetValorNecessarioReal() = %.2f, expected %.2f", result, expected)
	}
}

func TestEligibilityProfile_GetFaturamentoAnualReal(t *testing.T) {
	profile := EligibilityProfile{
		FaturamentoAnual: 5000000, // R$ 50.000,00
	}

	result := profile.GetFaturamentoAnualReal()
	expected := 50000.0

	if result != expected {
		t.Errorf("GetFaturamentoAnualReal() = %.2f, expected %.2f", result, expected)
	}
}

func TestEligibilityProfile_String(t *testing.T) {
	profile := EligibilityProfile{
		EntityID: "entity-1",
	}

	result := profile.String()
	if result == "" {
		t.Error("String() should not return empty string")
	}
	if result == "<nil>" {
		t.Error("String() should not return <nil>")
	}
}

func TestIsValidFinalidadeCredito(t *testing.T) {
	tests := []struct {
		finalidade FinalidadeCredito
		expected   bool
	}{
		{FinalidadeCapitalGiro, true},
		{FinalidadeEquipamento, true},
		{FinalidadeReforma, true},
		{FinalidadeOutro, true},
		{"INVALIDO", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.finalidade), func(t *testing.T) {
			result := isValidFinalidadeCredito(tt.finalidade)
			if result != tt.expected {
				t.Errorf("isValidFinalidadeCredito(%s) = %v, expected %v", tt.finalidade, result, tt.expected)
			}
		})
	}
}

func TestIsValidTipoEntidade(t *testing.T) {
	tests := []struct {
		tipo     TipoEntidade
		expected bool
	}{
		{TipoEntidadeMEI, true},
		{TipoEntidadeME, true},
		{TipoEntidadeEPP, true},
		{TipoEntidadeCooperativa, true},
		{TipoEntidadeOSC, true},
		{TipoEntidadeOSCIP, true},
		{TipoEntidadePF, true},
		{"INVALIDO", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.tipo), func(t *testing.T) {
			result := isValidTipoEntidade(tt.tipo)
			if result != tt.expected {
				t.Errorf("isValidTipoEntidade(%s) = %v, expected %v", tt.tipo, result, tt.expected)
			}
		})
	}
}

func TestGetFinalidadeCreditoLabel(t *testing.T) {
	tests := []struct {
		finalidade FinalidadeCredito
		expected   string
	}{
		{FinalidadeCapitalGiro, "Capital de Giro"},
		{FinalidadeEquipamento, "Equipamento"},
		{FinalidadeReforma, "Reforma"},
		{FinalidadeOutro, "Outro"},
		{"UNKNOWN", "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(string(tt.finalidade), func(t *testing.T) {
			result := GetFinalidadeCreditoLabel(tt.finalidade)
			if result != tt.expected {
				t.Errorf("GetFinalidadeCreditoLabel(%s) = %s, expected %s", tt.finalidade, result, tt.expected)
			}
		})
	}
}

func TestGetTipoEntidadeLabel(t *testing.T) {
	tests := []struct {
		tipo     TipoEntidade
		expected string
	}{
		{TipoEntidadeMEI, "Microempreendedor Individual (MEI)"},
		{TipoEntidadeME, "Microempresa (ME)"},
		{TipoEntidadeEPP, "Empresa de Pequeno Porte (EPP)"},
		{TipoEntidadeCooperativa, "Cooperativa"},
		{TipoEntidadeOSC, "Organização da Sociedade Civil (OSC)"},
		{TipoEntidadeOSCIP, "OSCIP"},
		{TipoEntidadePF, "Pessoa Física"},
		{"UNKNOWN", "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(string(tt.tipo), func(t *testing.T) {
			result := GetTipoEntidadeLabel(tt.tipo)
			if result != tt.expected {
				t.Errorf("GetTipoEntidadeLabel(%s) = %s, expected %s", tt.tipo, result, tt.expected)
			}
		})
	}
}
