package das

import (
	"context"
	"time"

	"github.com/providentia/digna/core_lume/internal/domain"
	"github.com/providentia/digna/core_lume/internal/repository"
	"github.com/providentia/digna/core_lume/internal/service"
	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

// DASMEI represents the public DAS MEI model
type DASMEI struct {
	ID             string
	EntityID       string
	Competencia    string
	ValorDevido    int64
	ValorPago      int64
	DataVencimento int64
	DataPagamento  int64
	Status         string
	SalarioMinimo  int64
	ActivityType   string
	CreatedAt      int64
	UpdatedAt      int64
}

// ActivityType represents the business activity type
type ActivityType string

const (
	ActivityTypeCommerce ActivityType = "COMERCIO"
	ActivityTypeService  ActivityType = "SERVICOS"
	ActivityTypeMixed    ActivityType = "MISTO"
)

// Alert represents a DAS MEI alert
type Alert struct {
	DASID       string
	Competencia string
	Message     string
	Severity    string
}

// Service provides DAS MEI operations
type Service struct {
	dasService *service.DASMEIService
}

// NewService creates a new DAS MEI service
func NewService(lm lifecycle.LifecycleManager) *Service {
	dasRepo := repository.NewSQLiteDASMEIRepository(lm)
	return &Service{
		dasService: service.NewDASMEIService(dasRepo),
	}
}

// EnsureTableExists ensures the das_mei table exists for the entity
func (s *Service) EnsureTableExists(entityID string) error {
	return s.dasService.InitTableForEntity(entityID)
}

// GenerateDASRequest represents a request to generate DAS
type GenerateDASRequest struct {
	Competencia  string
	ActivityType ActivityType
}

// GenerateMonthlyDAS generates a new DAS MEI for a specific month
func (s *Service) GenerateMonthlyDAS(ctx context.Context, entityID string, req *GenerateDASRequest) (*DASMEI, error) {
	svcReq := &service.GenerateDASRequest{
		Competencia:  req.Competencia,
		ActivityType: domain.ActivityType(req.ActivityType),
	}

	das, err := s.dasService.GenerateMonthlyDAS(ctx, entityID, svcReq)
	if err != nil {
		return nil, err
	}

	return convertToPublic(das), nil
}

// GetPendingDAS returns all pending DAS MEI
func (s *Service) GetPendingDAS(ctx context.Context, entityID string) ([]*DASMEI, error) {
	dasList, err := s.dasService.GetPendingDAS(ctx, entityID)
	if err != nil {
		return nil, err
	}

	result := make([]*DASMEI, len(dasList))
	for i, das := range dasList {
		result[i] = convertToPublic(das)
	}
	return result, nil
}

// GetOverdueDAS returns all overdue DAS MEI
func (s *Service) GetOverdueDAS(ctx context.Context, entityID string) ([]*DASMEI, error) {
	dasList, err := s.dasService.GetOverdueDAS(ctx, entityID)
	if err != nil {
		return nil, err
	}

	result := make([]*DASMEI, len(dasList))
	for i, das := range dasList {
		result[i] = convertToPublic(das)
	}
	return result, nil
}

// GetAllDAS returns all DAS MEI for an entity
func (s *Service) GetAllDAS(ctx context.Context, entityID string) ([]*DASMEI, error) {
	dasList, err := s.dasService.GetAllDAS(ctx, entityID)
	if err != nil {
		return nil, err
	}

	result := make([]*DASMEI, len(dasList))
	for i, das := range dasList {
		result[i] = convertToPublic(das)
	}
	return result, nil
}

// GetDASByCompetencia returns a DAS MEI by competencia
func (s *Service) GetDASByCompetencia(ctx context.Context, entityID, competencia string) (*DASMEI, error) {
	das, err := s.dasService.GetDASByCompetencia(ctx, entityID, competencia)
	if err != nil {
		return nil, err
	}
	return convertToPublic(das), nil
}

// MarkAsPaid marks a DAS MEI as paid
func (s *Service) MarkAsPaid(ctx context.Context, entityID, dasID string) error {
	return s.dasService.MarkAsPaid(ctx, entityID, dasID)
}

// GetMinimumWage returns the minimum wage for a year
func (s *Service) GetMinimumWage(year int) int64 {
	return s.dasService.GetMinimumWage(year)
}

// GetCurrentMinimumWage returns the minimum wage for current year
func (s *Service) GetCurrentMinimumWage() int64 {
	return s.dasService.GetCurrentMinimumWage()
}

// CheckOverdueAlerts checks for overdue DAS MEI
func (s *Service) CheckOverdueAlerts(ctx context.Context, entityID string) ([]Alert, error) {
	alerts, err := s.dasService.CheckOverdueAlerts(ctx, entityID)
	if err != nil {
		return nil, err
	}

	result := make([]Alert, len(alerts))
	for i, alert := range alerts {
		result[i] = Alert{
			DASID:       alert.DASID,
			Competencia: alert.Competencia,
			Message:     alert.Message,
			Severity:    alert.Severity,
		}
	}
	return result, nil
}

// CalculateDASAmount calculates DAS MEI amount
func (s *Service) CalculateDASAmount(year int, activity ActivityType) int64 {
	return s.dasService.CalculateDASAmount(year, domain.ActivityType(activity))
}

// UpdateDASStatus updates DAS MEI status
func (s *Service) UpdateDASStatus(ctx context.Context, entityID string) error {
	return s.dasService.UpdateDASStatus(ctx, entityID)
}

// GetCurrentCompetencia returns current month
func (s *Service) GetCurrentCompetencia() string {
	return s.dasService.GetCurrentCompetencia()
}

// GenerateCurrentMonthDAS generates DAS for current month
func (s *Service) GenerateCurrentMonthDAS(ctx context.Context, entityID string, activity ActivityType) (*DASMEI, error) {
	das, err := s.dasService.GenerateCurrentMonthDAS(ctx, entityID, domain.ActivityType(activity))
	if err != nil {
		return nil, err
	}
	return convertToPublic(das), nil
}

// GetDueDate returns the due date as time.Time
func (d *DASMEI) GetDueDate() time.Time {
	return time.Unix(d.DataVencimento, 0)
}

// GetPaymentDate returns the payment date as time.Time
func (d *DASMEI) GetPaymentDate() *time.Time {
	if d.DataPagamento == 0 {
		return nil
	}
	t := time.Unix(d.DataPagamento, 0)
	return &t
}

// IsOverdue checks if DAS is overdue
func (d *DASMEI) IsOverdue() bool {
	if d.Status == "PAGO" || d.Status == "CANCELADO" {
		return false
	}
	return time.Now().Unix() > d.DataVencimento
}

// Helper function to convert internal DASMEI to public
func convertToPublic(das *domain.DASMEI) *DASMEI {
	return &DASMEI{
		ID:             das.ID,
		EntityID:       das.EntityID,
		Competencia:    das.Competencia,
		ValorDevido:    das.ValorDevido,
		ValorPago:      das.ValorPago,
		DataVencimento: das.DataVencimento,
		DataPagamento:  das.DataPagamento,
		Status:         string(das.Status),
		SalarioMinimo:  das.SalarioMinimo,
		ActivityType:   string(das.ActivityType),
		CreatedAt:      das.CreatedAt,
		UpdatedAt:      das.UpdatedAt,
	}
}
