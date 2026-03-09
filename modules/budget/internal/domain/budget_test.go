package domain

import (
	"testing"
)

func TestBudgetPlan_Validate(t *testing.T) {
	tests := []struct {
		name    string
		plan    BudgetPlan
		wantErr bool
	}{
		{
			name: "valid plan",
			plan: BudgetPlan{
				EntityID: "test-entity",
				Period:   "2024-03",
				Category: "INSUMOS",
				Planned:  1000000,
			},
			wantErr: false,
		},
		{
			name: "invalid period format",
			plan: BudgetPlan{
				EntityID: "test-entity",
				Period:   "2024-3",
				Category: "INSUMOS",
				Planned:  1000000,
			},
			wantErr: true,
		},
		{
			name: "invalid category",
			plan: BudgetPlan{
				EntityID: "test-entity",
				Period:   "2024-03",
				Category: "INVALID_CATEGORY",
				Planned:  1000000,
			},
			wantErr: true,
		},
		{
			name: "negative planned value",
			plan: BudgetPlan{
				EntityID: "test-entity",
				Period:   "2024-03",
				Category: "INSUMOS",
				Planned:  -1000000,
			},
			wantErr: true,
		},
		{
			name: "zero planned value",
			plan: BudgetPlan{
				EntityID: "test-entity",
				Period:   "2024-03",
				Category: "INSUMOS",
				Planned:  0,
			},
			wantErr: true,
		},
		{
			name: "missing entity id",
			plan: BudgetPlan{
				Period:   "2024-03",
				Category: "INSUMOS",
				Planned:  1000000,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.plan.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("BudgetPlan.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBudgetPlan_CalculateExecution(t *testing.T) {
	tests := []struct {
		name           string
		plan           BudgetPlan
		executed       int64
		wantPercentage int
		wantStatus     string
	}{
		{
			name: "safe - 50%",
			plan: BudgetPlan{
				Planned: 1000000,
			},
			executed:       500000,
			wantPercentage: 50,
			wantStatus:     "SAFE",
		},
		{
			name: "safe - 70%",
			plan: BudgetPlan{
				Planned: 1000000,
			},
			executed:       700000,
			wantPercentage: 70,
			wantStatus:     "SAFE",
		},
		{
			name: "warning - 71%",
			plan: BudgetPlan{
				Planned: 1000000,
			},
			executed:       710000,
			wantPercentage: 71,
			wantStatus:     "WARNING",
		},
		{
			name: "warning - 100%",
			plan: BudgetPlan{
				Planned: 1000000,
			},
			executed:       1000000,
			wantPercentage: 100,
			wantStatus:     "WARNING",
		},
		{
			name: "exceeded - 101%",
			plan: BudgetPlan{
				Planned: 1000000,
			},
			executed:       1010000,
			wantPercentage: 100,
			wantStatus:     "EXCEEDED",
		},
		{
			name: "zero planned",
			plan: BudgetPlan{
				Planned: 0,
			},
			executed:       100000,
			wantPercentage: 0,
			wantStatus:     "SAFE",
		},
		{
			name: "negative executed",
			plan: BudgetPlan{
				Planned: 1000000,
			},
			executed:       -50000,
			wantPercentage: 0,
			wantStatus:     "SAFE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := tt.plan.CalculateExecution(tt.executed)

			if exec.Percentage != tt.wantPercentage {
				t.Errorf("CalculateExecution() Percentage = %v, want %v", exec.Percentage, tt.wantPercentage)
			}
			if exec.AlertStatus != tt.wantStatus {
				t.Errorf("CalculateExecution() AlertStatus = %v, want %v", exec.AlertStatus, tt.wantStatus)
			}
		})
	}
}

func TestGetCategoryLabel(t *testing.T) {
	tests := []struct {
		category string
		want     string
	}{
		{"INSUMOS", "Insumos"},
		{"ENERGIA", "Energia"},
		{"EQUIPAMENTOS", "Equipamentos"},
		{"TRANSPORTE", "Transporte"},
		{"MANUTENCAO", "Manutenção"},
		{"SERVICOS", "Serviços"},
		{"OUTROS", "Outros"},
		{"UNKNOWN", "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.category, func(t *testing.T) {
			if got := GetCategoryLabel(tt.category); got != tt.want {
				t.Errorf("GetCategoryLabel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAlertStatusLabel(t *testing.T) {
	tests := []struct {
		status string
		want   string
	}{
		{"SAFE", "Dentro do planejado"},
		{"WARNING", "Atenção: perto do limite"},
		{"EXCEEDED", "Ultrapassou o planejado"},
		{"UNKNOWN", "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			if got := GetAlertStatusLabel(tt.status); got != tt.want {
				t.Errorf("GetAlertStatusLabel() = %v, want %v", got, tt.want)
			}
		})
	}
}
