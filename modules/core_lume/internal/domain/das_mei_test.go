package domain

import (
	"testing"
	"time"
)

func TestCalculateDASMEIAmount(t *testing.T) {
	tests := []struct {
		name          string
		salarioMinimo int64
		activity      ActivityType
		expected      int64
	}{
		{
			name:          "Comércio 2026",
			salarioMinimo: 151800, // R$ 1.518,00
			activity:      ActivityTypeCommerce,
			expected:      7690, // 5% de 151800 = 7590 + 100 (ICMS)
		},
		{
			name:          "Serviços 2026",
			salarioMinimo: 151800,
			activity:      ActivityTypeService,
			expected:      8090, // 5% de 151800 = 7590 + 500 (ISS)
		},
		{
			name:          "Misto 2026",
			salarioMinimo: 151800,
			activity:      ActivityTypeMixed,
			expected:      8190, // 5% de 151800 = 7590 + 100 + 500
		},
		{
			name:          "Comércio 2025",
			salarioMinimo: 151800,
			activity:      ActivityTypeCommerce,
			expected:      7690,
		},
		{
			name:          "Atividade inválida",
			salarioMinimo: 151800,
			activity:      ActivityType("INVALIDO"),
			expected:      7590, // Apenas 5% do SM
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateDASMEIAmount(tt.salarioMinimo, tt.activity)
			if result != tt.expected {
				t.Errorf("CalculateDASMEIAmount() = %d, expected %d", result, tt.expected)
			}
		})
	}
}

func TestGetMinimumWageForYear(t *testing.T) {
	tests := []struct {
		year     int
		expected int64
	}{
		{2024, 141200},
		{2025, 151800},
		{2026, 151800},
		{2030, 151800}, // Ano futuro, deve retornar valor mais recente
	}

	for _, tt := range tests {
		t.Run("Year_"+string(rune(tt.year)), func(t *testing.T) {
			result := GetMinimumWageForYear(tt.year)
			if result != tt.expected {
				t.Errorf("GetMinimumWageForYear(%d) = %d, expected %d", tt.year, result, tt.expected)
			}
		})
	}
}

func TestDASMEI_Validate(t *testing.T) {
	now := time.Now().Unix()

	tests := []struct {
		name    string
		das     DASMEI
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid DAS MEI",
			das: DASMEI{
				ID:             "das-test-1",
				EntityID:       "entity-1",
				Competencia:    "2026-03",
				ValorDevido:    7690,
				ValorPago:      0,
				DataVencimento: now + 86400,
				Status:         DASMEIStatusPending,
				SalarioMinimo:  151800,
				ActivityType:   ActivityTypeCommerce,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			wantErr: false,
		},
		{
			name: "Missing EntityID",
			das: DASMEI{
				ID:          "das-test-2",
				Competencia: "2026-03",
			},
			wantErr: true,
			errMsg:  "entity ID is required",
		},
		{
			name: "Invalid Competencia",
			das: DASMEI{
				ID:          "das-test-3",
				EntityID:    "entity-1",
				Competencia: "invalid",
			},
			wantErr: true,
			errMsg:  "competencia is required",
		},
		{
			name: "Invalid Status",
			das: DASMEI{
				ID:          "das-test-4",
				EntityID:    "entity-1",
				Competencia: "2026-03",
				Status:      DASMEIStatus("INVALID"),
			},
			wantErr: true,
			errMsg:  "invalid DAS MEI status",
		},
		{
			name: "Invalid Activity Type",
			das: DASMEI{
				ID:           "das-test-5",
				EntityID:     "entity-1",
				Competencia:  "2026-03",
				Status:       DASMEIStatusPending,
				ActivityType: ActivityType("INVALID"),
			},
			wantErr: true,
			errMsg:  "invalid activity type",
		},
		{
			name: "Zero ValorDevido",
			das: DASMEI{
				ID:            "das-test-6",
				EntityID:      "entity-1",
				Competencia:   "2026-03",
				Status:        DASMEIStatusPending,
				ActivityType:  ActivityTypeCommerce,
				ValorDevido:   0,
				SalarioMinimo: 151800,
			},
			wantErr: true,
			errMsg:  "valor devido must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.das.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Validate() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestDASMEI_IsOverdue(t *testing.T) {
	now := time.Now().Unix()
	yesterday := now - 86400
	tomorrow := now + 86400

	tests := []struct {
		name     string
		status   DASMEIStatus
		dueDate  int64
		expected bool
	}{
		{
			name:     "Pending and overdue",
			status:   DASMEIStatusPending,
			dueDate:  yesterday,
			expected: true,
		},
		{
			name:     "Pending and not overdue",
			status:   DASMEIStatusPending,
			dueDate:  tomorrow,
			expected: false,
		},
		{
			name:     "Paid (not overdue)",
			status:   DASMEIStatusPaid,
			dueDate:  yesterday,
			expected: false,
		},
		{
			name:     "Canceled (not overdue)",
			status:   DASMEIStatusCanceled,
			dueDate:  yesterday,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			das := DASMEI{
				Status:         tt.status,
				DataVencimento: tt.dueDate,
			}
			result := das.IsOverdue()
			if result != tt.expected {
				t.Errorf("IsOverdue() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestDASMEI_MarkAsPaid(t *testing.T) {
	now := time.Now().Unix()

	tests := []struct {
		name    string
		das     DASMEI
		wantErr bool
		errMsg  string
	}{
		{
			name: "Mark pending as paid",
			das: DASMEI{
				ID:            "das-1",
				Status:        DASMEIStatusPending,
				ValorDevido:   7690,
				DataPagamento: 0,
			},
			wantErr: false,
		},
		{
			name: "Try to mark already paid",
			das: DASMEI{
				ID:            "das-2",
				Status:        DASMEIStatusPaid,
				ValorDevido:   7690,
				ValorPago:     7690,
				DataPagamento: now,
			},
			wantErr: true,
			errMsg:  "DAS MEI is already paid",
		},
		{
			name: "Try to pay canceled",
			das: DASMEI{
				ID:     "das-3",
				Status: DASMEIStatusCanceled,
			},
			wantErr: true,
			errMsg:  "cannot pay a canceled DAS MEI",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.das.MarkAsPaid()
			if tt.wantErr {
				if err == nil {
					t.Errorf("MarkAsPaid() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("MarkAsPaid() unexpected error: %v", err)
				}
				if tt.das.Status != DASMEIStatusPaid {
					t.Errorf("MarkAsPaid() status = %v, expected PAID", tt.das.Status)
				}
				if tt.das.ValorPago != tt.das.ValorDevido {
					t.Errorf("MarkAsPaid() valorPago = %d, expected %d", tt.das.ValorPago, tt.das.ValorDevido)
				}
				if tt.das.DataPagamento == 0 {
					t.Errorf("MarkAsPaid() dataPagamento should be set")
				}
			}
		})
	}
}

func TestDASMEI_Cancel(t *testing.T) {
	tests := []struct {
		name    string
		status  DASMEIStatus
		wantErr bool
	}{
		{
			name:    "Cancel pending",
			status:  DASMEIStatusPending,
			wantErr: false,
		},
		{
			name:    "Try to cancel paid",
			status:  DASMEIStatusPaid,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			das := DASMEI{
				Status: tt.status,
			}
			err := das.Cancel()
			if tt.wantErr {
				if err == nil {
					t.Errorf("Cancel() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Cancel() unexpected error: %v", err)
				}
				if das.Status != DASMEIStatusCanceled {
					t.Errorf("Cancel() status = %v, expected CANCELED", das.Status)
				}
			}
		})
	}
}

func TestParseCompetencia(t *testing.T) {
	tests := []struct {
		name        string
		competencia string
		wantYear    int
		wantMonth   int
		wantErr     bool
	}{
		{
			name:        "Valid 2026-03",
			competencia: "2026-03",
			wantYear:    2026,
			wantMonth:   3,
			wantErr:     false,
		},
		{
			name:        "Valid 2025-12",
			competencia: "2025-12",
			wantYear:    2025,
			wantMonth:   12,
			wantErr:     false,
		},
		{
			name:        "Invalid format",
			competencia: "invalid",
			wantErr:     true,
		},
		{
			name:        "Invalid month 13",
			competencia: "2026-13",
			wantErr:     true,
		},
		{
			name:        "Invalid month 0",
			competencia: "2026-00",
			wantErr:     true,
		},
		{
			name:        "Empty string",
			competencia: "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			year, month, err := ParseCompetencia(tt.competencia)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseCompetencia() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("ParseCompetencia() unexpected error: %v", err)
				}
				if year != tt.wantYear {
					t.Errorf("ParseCompetencia() year = %d, expected %d", year, tt.wantYear)
				}
				if month != tt.wantMonth {
					t.Errorf("ParseCompetencia() month = %d, expected %d", month, tt.wantMonth)
				}
			}
		})
	}
}

func TestCalculateDueDate(t *testing.T) {
	// Teste para garantir que dia 20 é usado (ou dia útil anterior)
	tests := []struct {
		year        int
		month       int
		description string
	}{
		{2026, 3, "March 2026 (Friday 20th)"},
		{2026, 4, "April 2026 (Monday 20th)"},
		{2026, 5, "May 2026 (Wednesday 20th)"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			dueDate := CalculateDueDate(tt.year, tt.month)
			dueTime := time.Unix(dueDate, 0)

			// Verificar se é dia 20 ou anterior (se caiu em fim de semana)
			if dueTime.Day() > 20 {
				t.Errorf("CalculateDueDate() day = %d, should be <= 20", dueTime.Day())
			}

			// Verificar se não é fim de semana
			weekday := dueTime.Weekday()
			if weekday == time.Saturday || weekday == time.Sunday {
				t.Errorf("CalculateDueDate() fell on weekend: %v", weekday)
			}
		})
	}
}

func TestGenerateCompetencia(t *testing.T) {
	tests := []struct {
		date     time.Time
		expected string
	}{
		{time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC), "2026-03"},
		{time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC), "2026-12"},
		{time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), "2025-01"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := GenerateCompetencia(tt.date)
			if result != tt.expected {
				t.Errorf("GenerateCompetencia() = %s, expected %s", result, tt.expected)
			}
		})
	}
}

func TestDASMEI_GetValorDevidoReal(t *testing.T) {
	das := DASMEI{
		ValorDevido: 7690, // 76,90 em centavos
	}

	result := das.GetValorDevidoReal()
	expected := 76.90

	if result != expected {
		t.Errorf("GetValorDevidoReal() = %f, expected %f", result, expected)
	}
}

func TestDASMEI_GetSalarioMinimoReal(t *testing.T) {
	das := DASMEI{
		SalarioMinimo: 151800, // 1.518,00 em centavos
	}

	result := das.GetSalarioMinimoReal()
	expected := 1518.00

	if result != expected {
		t.Errorf("GetSalarioMinimoReal() = %f, expected %f", result, expected)
	}
}

func TestDASMEI_String(t *testing.T) {
	das := DASMEI{
		ID:          "das-test",
		Competencia: "2026-03",
		ValorDevido: 7690,
		Status:      DASMEIStatusPending,
	}

	result := das.String()
	expected := "DASMEI{ID: das-test, Competencia: 2026-03, Valor: 7690, Status: PENDENTE}"

	if result != expected {
		t.Errorf("String() = %s, expected %s", result, expected)
	}
}
