package domain

import (
	"time"
)

type AccountantStatus string

const (
	StatusActive   AccountantStatus = "ACTIVE"
	StatusInactive AccountantStatus = "INACTIVE"
)

type EnterpriseAccountant struct {
	ID           string           `json:"id"`
	EnterpriseID string           `json:"enterprise_id"`
	AccountantID string           `json:"accountant_id"`
	Status       AccountantStatus `json:"status"`
	StartDate    time.Time        `json:"start_date"`
	EndDate      *time.Time       `json:"end_date,omitempty"`
	DelegatedBy  string           `json:"delegated_by"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

func NewEnterpriseAccountant(enterpriseID, accountantID, delegatedBy string) *EnterpriseAccountant {
	now := time.Now().UTC()
	return &EnterpriseAccountant{
		EnterpriseID: enterpriseID,
		AccountantID: accountantID,
		Status:       StatusActive,
		StartDate:    now,
		EndDate:      nil, // Ativo, sem data de término
		DelegatedBy:  delegatedBy,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (ea *EnterpriseAccountant) IsActive() bool {
	return ea.Status == StatusActive
}

func (ea *EnterpriseAccountant) IsInactive() bool {
	return ea.Status == StatusInactive
}

func (ea *EnterpriseAccountant) Deactivate(endDate time.Time) error {
	if ea.IsInactive() {
		return nil // Já está inativo
	}

	if endDate.Before(ea.StartDate) {
		return ErrInvalidEndDate
	}

	ea.Status = StatusInactive
	ea.EndDate = &endDate
	ea.UpdatedAt = time.Now().UTC()
	return nil
}

func (ea *EnterpriseAccountant) Reactivate() {
	if ea.IsActive() {
		return // Já está ativo
	}

	ea.Status = StatusActive
	ea.EndDate = nil
	ea.UpdatedAt = time.Now().UTC()
}

func (ea *EnterpriseAccountant) IsValidForDate(checkDate time.Time) bool {
	if ea.IsInactive() && ea.EndDate != nil {
		return !checkDate.Before(ea.StartDate) && !checkDate.After(*ea.EndDate)
	}
	// Para ativos, apenas verifica se a data é após o início
	return !checkDate.Before(ea.StartDate)
}

func (ea *EnterpriseAccountant) GetDateRange() (startDate, endDate time.Time) {
	startDate = ea.StartDate
	if ea.EndDate != nil {
		endDate = *ea.EndDate
	} else {
		endDate = time.Now().UTC()
	}
	return startDate, endDate
}

// Erros de domínio
var (
	ErrInvalidEndDate  = newDomainError("end date must be after start date")
	ErrAlreadyActive   = newDomainError("accountant link is already active")
	ErrAlreadyInactive = newDomainError("accountant link is already inactive")
)

type DomainError struct {
	message string
}

func newDomainError(message string) *DomainError {
	return &DomainError{message: message}
}

func (e *DomainError) Error() string {
	return e.message
}
