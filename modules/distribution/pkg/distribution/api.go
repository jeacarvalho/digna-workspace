// Package distribution expõe a API pública do módulo de distribuição de sobras
package distribution

import (
	"github.com/providentia/digna/distribution/internal/domain"
	"github.com/providentia/digna/distribution/internal/repository"
	"github.com/providentia/digna/distribution/internal/service"
	"github.com/providentia/digna/reporting/pkg/surplus"
)

// Tipos exportados
type (
	Distribution       = domain.Distribution
	DistributionStatus = domain.DistributionStatus
	DistributionMember = domain.DistributionMember
	SurplusCalculation = domain.SurplusCalculation
	MemberCalculation  = domain.MemberCalculation
	LedgerEntry        = domain.LedgerEntry
	Posting            = domain.Posting
)

// Status constants
const (
	StatusCalculated      = domain.StatusCalculated
	StatusPendingApproval = domain.StatusPendingApproval
	StatusApproved        = domain.StatusApproved
	StatusExecuted        = domain.StatusExecuted
	StatusRejected        = domain.StatusRejected
)

// Service exportado
type DistributionService = service.DistributionService

// NewDistributionService cria novo serviço
func NewDistributionService(calculator *surplus.Calculator) *DistributionService {
	distRepo := repository.NewMockDistributionRepository()
	assemblyRepo := repository.NewMockAssemblyRepository()
	return service.NewDistributionService(distRepo, assemblyRepo, calculator)
}

// Erros exportados
var (
	ErrInvalidSurplus  = service.ErrInvalidSurplus
	ErrPendingDecision = service.ErrPendingDecision
	ErrAlreadyExecuted = service.ErrAlreadyExecuted
)
