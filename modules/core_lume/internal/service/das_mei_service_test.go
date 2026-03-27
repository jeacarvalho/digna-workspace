package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/providentia/digna/core_lume/internal/domain"
)

// MockDASMEIRepository é um mock do repositório para testes
type MockDASMEIRepository struct {
	savedDAS       []*domain.DASMEI
	findByIDFunc   func(entityID, dasID string) (*domain.DASMEI, error)
	findByCompFunc func(entityID, competencia string) (*domain.DASMEI, error)
	listByEntity   func(entityID string) ([]*domain.DASMEI, error)
	listPending    func(entityID string) ([]*domain.DASMEI, error)
	listOverdue    func(entityID string) ([]*domain.DASMEI, error)
}

func (m *MockDASMEIRepository) Save(das *domain.DASMEI) error {
	m.savedDAS = append(m.savedDAS, das)
	return nil
}

func (m *MockDASMEIRepository) FindByID(entityID, dasID string) (*domain.DASMEI, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(entityID, dasID)
	}
	return nil, nil
}

func (m *MockDASMEIRepository) FindByCompetencia(entityID, competencia string) (*domain.DASMEI, error) {
	if m.findByCompFunc != nil {
		return m.findByCompFunc(entityID, competencia)
	}
	return nil, nil
}

func (m *MockDASMEIRepository) ListByEntity(entityID string) ([]*domain.DASMEI, error) {
	if m.listByEntity != nil {
		return m.listByEntity(entityID)
	}
	return []*domain.DASMEI{}, nil
}

func (m *MockDASMEIRepository) ListPending(entityID string) ([]*domain.DASMEI, error) {
	if m.listPending != nil {
		return m.listPending(entityID)
	}
	return []*domain.DASMEI{}, nil
}

func (m *MockDASMEIRepository) ListOverdue(entityID string) ([]*domain.DASMEI, error) {
	if m.listOverdue != nil {
		return m.listOverdue(entityID)
	}
	return []*domain.DASMEI{}, nil
}

func (m *MockDASMEIRepository) Update(das *domain.DASMEI) error {
	return nil
}

func (m *MockDASMEIRepository) MarkAsPaid(entityID, dasID string, valorPago int64) error {
	return nil
}

func TestDASMEIService_GenerateMonthlyDAS(t *testing.T) {
	mockRepo := &MockDASMEIRepository{}
	service := NewDASMEIService(mockRepo)
	ctx := context.Background()

	tests := []struct {
		name    string
		req     *GenerateDASRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "Generate valid DAS - Comércio",
			req: &GenerateDASRequest{
				Competencia:  "2026-03",
				ActivityType: domain.ActivityTypeCommerce,
			},
			wantErr: false,
		},
		{
			name: "Generate valid DAS - Serviços",
			req: &GenerateDASRequest{
				Competencia:  "2026-04",
				ActivityType: domain.ActivityTypeService,
			},
			wantErr: false,
		},
		{
			name: "Missing competencia",
			req: &GenerateDASRequest{
				Competencia:  "",
				ActivityType: domain.ActivityTypeCommerce,
			},
			wantErr: true,
			errMsg:  "competencia is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			das, err := service.GenerateMonthlyDAS(ctx, "entity-test", tt.req)
			if tt.wantErr {
				if err == nil {
					t.Errorf("GenerateMonthlyDAS() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("GenerateMonthlyDAS() unexpected error: %v", err)
				}
				if das == nil {
					t.Errorf("GenerateMonthlyDAS() returned nil DAS")
					return
				}
				if das.EntityID != "entity-test" {
					t.Errorf("GenerateMonthlyDAS() EntityID = %s, expected entity-test", das.EntityID)
				}
				if das.Competencia != tt.req.Competencia {
					t.Errorf("GenerateMonthlyDAS() Competencia = %s, expected %s", das.Competencia, tt.req.Competencia)
				}
				// Verificar se o valor foi calculado corretamente
				expectedAmount := domain.CalculateDASMEIAmount(das.SalarioMinimo, tt.req.ActivityType)
				if das.ValorDevido != expectedAmount {
					t.Errorf("GenerateMonthlyDAS() ValorDevido = %d, expected %d", das.ValorDevido, expectedAmount)
				}
			}
		})
	}
}

func TestDASMEIService_GetMinimumWage(t *testing.T) {
	mockRepo := &MockDASMEIRepository{}
	service := NewDASMEIService(mockRepo)

	tests := []struct {
		year     int
		expected int64
	}{
		{2024, 141200},
		{2025, 151800},
		{2026, 151800},
	}

	for _, tt := range tests {
		t.Run("Year_"+string(rune(tt.year)), func(t *testing.T) {
			result := service.GetMinimumWage(tt.year)
			if result != tt.expected {
				t.Errorf("GetMinimumWage(%d) = %d, expected %d", tt.year, result, tt.expected)
			}
		})
	}
}

func TestDASMEIService_CalculateDASAmount(t *testing.T) {
	mockRepo := &MockDASMEIRepository{}
	service := NewDASMEIService(mockRepo)

	tests := []struct {
		year     int
		activity domain.ActivityType
		expected int64
	}{
		{2026, domain.ActivityTypeCommerce, 7690},
		{2026, domain.ActivityTypeService, 8090},
		{2026, domain.ActivityTypeMixed, 8190},
	}

	for _, tt := range tests {
		t.Run(string(tt.activity), func(t *testing.T) {
			result := service.CalculateDASAmount(tt.year, tt.activity)
			if result != tt.expected {
				t.Errorf("CalculateDASAmount() = %d, expected %d", result, tt.expected)
			}
		})
	}
}

func TestDASMEIService_GetCurrentCompetencia(t *testing.T) {
	mockRepo := &MockDASMEIRepository{}
	service := NewDASMEIService(mockRepo)

	result := service.GetCurrentCompetencia()
	now := time.Now()
	expected := now.Format("2006-01")

	if result != expected {
		t.Errorf("GetCurrentCompetencia() = %s, expected %s", result, expected)
	}
}

func TestDASMEIService_GetPreviousCompetencia(t *testing.T) {
	mockRepo := &MockDASMEIRepository{}
	service := NewDASMEIService(mockRepo)

	result := service.GetPreviousCompetencia()
	prevMonth := time.Now().AddDate(0, -1, 0)
	expected := prevMonth.Format("2006-01")

	if result != expected {
		t.Errorf("GetPreviousCompetencia() = %s, expected %s", result, expected)
	}
}

func TestDASMEIService_CheckOverdueAlerts(t *testing.T) {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1).Unix()
	tomorrow := now.AddDate(0, 0, 1).Unix()
	nextWeek := now.AddDate(0, 0, 7).Unix()

	mockRepo := &MockDASMEIRepository{
		listPending: func(entityID string) ([]*domain.DASMEI, error) {
			return []*domain.DASMEI{
				{
					ID:             "das-1",
					Competencia:    "2026-02",
					DataVencimento: yesterday, // Vencido
					Status:         domain.DASMEIStatusPending,
				},
				{
					ID:             "das-2",
					Competencia:    "2026-03",
					DataVencimento: tomorrow, // Vence amanhã
					Status:         domain.DASMEIStatusPending,
				},
				{
					ID:             "das-3",
					Competencia:    "2026-04",
					DataVencimento: nextWeek, // Vence próxima semana
					Status:         domain.DASMEIStatusPending,
				},
			}, nil
		},
	}

	service := NewDASMEIService(mockRepo)
	ctx := context.Background()

	alerts, err := service.CheckOverdueAlerts(ctx, "entity-test")
	if err != nil {
		t.Errorf("CheckOverdueAlerts() unexpected error: %v", err)
		return
	}

	if len(alerts) == 0 {
		t.Errorf("CheckOverdueAlerts() expected alerts but got none")
		return
	}

	// Verificar se temos alertas de diferentes severidades
	var hasCritical bool
	for _, alert := range alerts {
		if alert.Severity == "CRITICAL" {
			hasCritical = true
			break
		}
	}

	if !hasCritical {
		t.Errorf("CheckOverdueAlerts() expected CRITICAL alert for overdue DAS")
	}
}

func TestDASMEIService_MarkAsPaid(t *testing.T) {
	mockRepo := &MockDASMEIRepository{
		findByIDFunc: func(entityID, dasID string) (*domain.DASMEI, error) {
			return &domain.DASMEI{
				ID:            dasID,
				EntityID:      entityID,
				Status:        domain.DASMEIStatusPending,
				ValorDevido:   7690,
				ValorPago:     0,
				DataPagamento: 0,
			}, nil
		},
	}

	service := NewDASMEIService(mockRepo)
	ctx := context.Background()

	err := service.MarkAsPaid(ctx, "entity-test", "das-1")
	if err != nil {
		t.Errorf("MarkAsPaid() unexpected error: %v", err)
	}
}

func TestDASMEIService_GenerateCurrentMonthDAS(t *testing.T) {
	mockRepo := &MockDASMEIRepository{
		findByCompFunc: func(entityID, competencia string) (*domain.DASMEI, error) {
			// Simula que não existe DAS para o mês atual (retorna erro)
			return nil, fmt.Errorf("DAS MEI not found")
		},
	}

	service := NewDASMEIService(mockRepo)
	ctx := context.Background()

	das, err := service.GenerateCurrentMonthDAS(ctx, "entity-test", domain.ActivityTypeCommerce)
	if err != nil {
		t.Errorf("GenerateCurrentMonthDAS() unexpected error: %v", err)
		return
	}

	if das == nil {
		t.Errorf("GenerateCurrentMonthDAS() returned nil DAS")
		return
	}

	// Verificar se a competência é o mês atual
	currentComp := service.GetCurrentCompetencia()
	if das.Competencia != currentComp {
		t.Errorf("GenerateCurrentMonthDAS() Competencia = %s, expected %s", das.Competencia, currentComp)
	}
}

func TestDASMEIService_GenerateCurrentMonthDAS_AlreadyExists(t *testing.T) {
	mockRepo := &MockDASMEIRepository{
		findByCompFunc: func(entityID, competencia string) (*domain.DASMEI, error) {
			// Simula que já existe DAS para o mês atual
			return &domain.DASMEI{
				ID:          "existing-das",
				Competencia: competencia,
			}, nil
		},
	}

	service := NewDASMEIService(mockRepo)
	ctx := context.Background()

	_, err := service.GenerateCurrentMonthDAS(ctx, "entity-test", domain.ActivityTypeCommerce)
	if err == nil {
		t.Errorf("GenerateCurrentMonthDAS() expected error for existing DAS but got nil")
	}
}

func TestDASMEIService_UpdateDASStatus(t *testing.T) {
	now := time.Now().Unix()
	yesterday := now - 86400

	mockRepo := &MockDASMEIRepository{
		listPending: func(entityID string) ([]*domain.DASMEI, error) {
			return []*domain.DASMEI{
				{
					ID:             "das-1",
					EntityID:       entityID,
					Competencia:    "2026-02",
					DataVencimento: yesterday,
					Status:         domain.DASMEIStatusPending,
				},
			}, nil
		},
	}

	service := NewDASMEIService(mockRepo)
	ctx := context.Background()

	err := service.UpdateDASStatus(ctx, "entity-test")
	if err != nil {
		t.Errorf("UpdateDASStatus() unexpected error: %v", err)
	}
}

func TestDASMEIService_GetCurrentMinimumWage(t *testing.T) {
	mockRepo := &MockDASMEIRepository{}
	service := NewDASMEIService(mockRepo)

	result := service.GetCurrentMinimumWage()
	expected := int64(151800) // 2026

	if result != expected {
		t.Errorf("GetCurrentMinimumWage() = %d, expected %d", result, expected)
	}
}
