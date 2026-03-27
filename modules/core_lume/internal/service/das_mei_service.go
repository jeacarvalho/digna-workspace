package service

import (
	"context"
	"fmt"
	"time"

	"github.com/providentia/digna/core_lume/internal/domain"
	"github.com/providentia/digna/core_lume/internal/repository"
)

// DASMEIService implements application logic for DAS MEI management
type DASMEIService struct {
	dasRepo repository.DASMEIRepository
}

// NewDASMEIService creates a new DASMEIService
func NewDASMEIService(dasRepo repository.DASMEIRepository) *DASMEIService {
	return &DASMEIService{
		dasRepo: dasRepo,
	}
}

// DASMEIRepositoryWithInit interface that supports table initialization
type DASMEIRepositoryWithInit interface {
	repository.DASMEIRepository
	InitTable(entityID string) error
}

// InitTableForEntity initializes the das_mei table for a specific entity
func (s *DASMEIService) InitTableForEntity(entityID string) error {
	if repoWithInit, ok := s.dasRepo.(DASMEIRepositoryWithInit); ok {
		return repoWithInit.InitTable(entityID)
	}
	return nil
}

// GenerateDASRequest represents the request to generate a new DAS MEI
type GenerateDASRequest struct {
	Competencia  string
	ActivityType domain.ActivityType
}

// GenerateMonthlyDAS generates a DAS MEI for a specific month
func (s *DASMEIService) GenerateMonthlyDAS(ctx context.Context, entityID string, req *GenerateDASRequest) (*domain.DASMEI, error) {
	// Validate request
	if req.Competencia == "" {
		return nil, fmt.Errorf("competencia is required")
	}

	// Parse competencia to extract year and month
	year, month, err := domain.ParseCompetencia(req.Competencia)
	if err != nil {
		return nil, fmt.Errorf("invalid competencia format: %w", err)
	}

	// Check if DAS already exists for this competencia
	existingDAS, err := s.dasRepo.FindByCompetencia(entityID, req.Competencia)
	if err == nil && existingDAS != nil {
		return nil, fmt.Errorf("DAS MEI already exists for competencia %s", req.Competencia)
	}

	// Get minimum wage for the year
	salarioMinimo := domain.GetMinimumWageForYear(year)

	// Calculate DAS amount based on activity type
	valorDevido := domain.CalculateDASMEIAmount(salarioMinimo, req.ActivityType)

	// Calculate due date (20th or previous business day)
	dataVencimento := domain.CalculateDueDate(year, month)

	now := time.Now().Unix()
	das := &domain.DASMEI{
		ID:             fmt.Sprintf("das-%d", now),
		EntityID:       entityID,
		Competencia:    req.Competencia,
		ValorDevido:    valorDevido,
		ValorPago:      0,
		DataVencimento: dataVencimento,
		DataPagamento:  0,
		Status:         domain.DASMEIStatusPending,
		SalarioMinimo:  salarioMinimo,
		ActivityType:   req.ActivityType,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.dasRepo.Save(das); err != nil {
		return nil, fmt.Errorf("failed to save DAS MEI: %w", err)
	}

	return das, nil
}

// GetPendingDAS returns all pending DAS MEI for an entity
func (s *DASMEIService) GetPendingDAS(ctx context.Context, entityID string) ([]*domain.DASMEI, error) {
	return s.dasRepo.ListPending(entityID)
}

// GetOverdueDAS returns all overdue DAS MEI for an entity
func (s *DASMEIService) GetOverdueDAS(ctx context.Context, entityID string) ([]*domain.DASMEI, error) {
	return s.dasRepo.ListOverdue(entityID)
}

// GetAllDAS returns all DAS MEI for an entity
func (s *DASMEIService) GetAllDAS(ctx context.Context, entityID string) ([]*domain.DASMEI, error) {
	return s.dasRepo.ListByEntity(entityID)
}

// GetDASByCompetencia returns a specific DAS MEI by competencia
func (s *DASMEIService) GetDASByCompetencia(ctx context.Context, entityID, competencia string) (*domain.DASMEI, error) {
	return s.dasRepo.FindByCompetencia(entityID, competencia)
}

// MarkAsPaid marks a DAS MEI as paid
func (s *DASMEIService) MarkAsPaid(ctx context.Context, entityID, dasID string) error {
	das, err := s.dasRepo.FindByID(entityID, dasID)
	if err != nil {
		return fmt.Errorf("DAS MEI not found: %w", err)
	}

	if err := das.MarkAsPaid(); err != nil {
		return fmt.Errorf("cannot mark DAS MEI as paid: %w", err)
	}

	return s.dasRepo.Save(das)
}

// GetMinimumWage returns the minimum wage for a specific year
func (s *DASMEIService) GetMinimumWage(year int) int64 {
	return domain.GetMinimumWageForYear(year)
}

// GetCurrentMinimumWage returns the minimum wage for the current year
func (s *DASMEIService) GetCurrentMinimumWage() int64 {
	return domain.GetMinimumWageForYear(time.Now().Year())
}

// Alert represents a DAS MEI alert
type Alert struct {
	DASID       string
	Competencia string
	Message     string
	Severity    string // INFO, WARNING, CRITICAL
}

// CheckOverdueAlerts checks for overdue DAS MEI and returns alerts
func (s *DASMEIService) CheckOverdueAlerts(ctx context.Context, entityID string) ([]Alert, error) {
	pendingDAS, err := s.dasRepo.ListPending(entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending DAS: %w", err)
	}

	var alerts []Alert
	now := time.Now()

	for _, das := range pendingDAS {
		dueDate := das.GetDueDate()
		daysUntilDue := int(dueDate.Sub(now).Hours() / 24)

		if daysUntilDue < 0 {
			// Overdue
			alerts = append(alerts, Alert{
				DASID:       das.ID,
				Competencia: das.Competencia,
				Message:     fmt.Sprintf("DAS MEI %s está vencido há %d dias", das.Competencia, -daysUntilDue),
				Severity:    "CRITICAL",
			})
		} else if daysUntilDue <= 1 {
			// Due today or tomorrow
			alerts = append(alerts, Alert{
				DASID:       das.ID,
				Competencia: das.Competencia,
				Message:     fmt.Sprintf("DAS MEI %s vence em %d dia(s)", das.Competencia, daysUntilDue),
				Severity:    "WARNING",
			})
		} else if daysUntilDue <= 5 {
			// Due within 5 days
			alerts = append(alerts, Alert{
				DASID:       das.ID,
				Competencia: das.Competencia,
				Message:     fmt.Sprintf("DAS MEI %s vence em %d dias", das.Competencia, daysUntilDue),
				Severity:    "INFO",
			})
		}
	}

	return alerts, nil
}

// CalculateDASAmount calculates the DAS MEI amount for a given year and activity type
func (s *DASMEIService) CalculateDASAmount(year int, activity domain.ActivityType) int64 {
	salarioMinimo := domain.GetMinimumWageForYear(year)
	return domain.CalculateDASMEIAmount(salarioMinimo, activity)
}

// GetDASDetails returns detailed information about a DAS MEI
func (s *DASMEIService) GetDASDetails(ctx context.Context, entityID, dasID string) (*domain.DASMEI, error) {
	return s.dasRepo.FindByID(entityID, dasID)
}

// CancelDAS cancels a DAS MEI
func (s *DASMEIService) CancelDAS(ctx context.Context, entityID, dasID string) error {
	das, err := s.dasRepo.FindByID(entityID, dasID)
	if err != nil {
		return fmt.Errorf("DAS MEI not found: %w", err)
	}

	if err := das.Cancel(); err != nil {
		return fmt.Errorf("cannot cancel DAS MEI: %w", err)
	}

	return s.dasRepo.Save(das)
}

// UpdateDASStatus updates the status of pending DAS MEI based on current date
func (s *DASMEIService) UpdateDASStatus(ctx context.Context, entityID string) error {
	pendingDAS, err := s.dasRepo.ListPending(entityID)
	if err != nil {
		return fmt.Errorf("failed to get pending DAS: %w", err)
	}

	for _, das := range pendingDAS {
		das.UpdateStatus()
		if das.Status == domain.DASMEIStatusOverdue {
			if err := s.dasRepo.Update(das); err != nil {
				return fmt.Errorf("failed to update DAS MEI status: %w", err)
			}
		}
	}

	return nil
}

// GetCurrentCompetencia returns the current month in YYYY-MM format
func (s *DASMEIService) GetCurrentCompetencia() string {
	return domain.GenerateCompetencia(time.Now())
}

// GetPreviousCompetencia returns the previous month in YYYY-MM format
func (s *DASMEIService) GetPreviousCompetencia() string {
	t := time.Now().AddDate(0, -1, 0)
	return domain.GenerateCompetencia(t)
}

// GenerateCurrentMonthDAS generates DAS for the current month if it doesn't exist
func (s *DASMEIService) GenerateCurrentMonthDAS(ctx context.Context, entityID string, activity domain.ActivityType) (*domain.DASMEI, error) {
	competencia := s.GetCurrentCompetencia()

	// Check if already exists
	_, err := s.dasRepo.FindByCompetencia(entityID, competencia)
	if err == nil {
		return nil, fmt.Errorf("DAS MEI for current month already exists")
	}

	req := &GenerateDASRequest{
		Competencia:  competencia,
		ActivityType: activity,
	}

	return s.GenerateMonthlyDAS(ctx, entityID, req)
}
