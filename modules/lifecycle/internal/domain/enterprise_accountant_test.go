package domain

import (
	"testing"
	"time"
)

func TestNewEnterpriseAccountant(t *testing.T) {
	enterpriseID := "ent_123"
	accountantID := "acc_456"
	delegatedBy := "user_789"

	link := NewEnterpriseAccountant(enterpriseID, accountantID, delegatedBy)

	if link.EnterpriseID != enterpriseID {
		t.Errorf("Expected EnterpriseID %s, got %s", enterpriseID, link.EnterpriseID)
	}
	if link.AccountantID != accountantID {
		t.Errorf("Expected AccountantID %s, got %s", accountantID, link.AccountantID)
	}
	if link.DelegatedBy != delegatedBy {
		t.Errorf("Expected DelegatedBy %s, got %s", delegatedBy, link.DelegatedBy)
	}
	if link.Status != StatusActive {
		t.Errorf("Expected Status ACTIVE, got %s", link.Status)
	}
	if link.EndDate != nil {
		t.Error("Expected EndDate to be nil for active link")
	}
	if link.StartDate.IsZero() {
		t.Error("Expected StartDate to be set")
	}
	if link.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
	if link.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}
}

func TestEnterpriseAccountant_Deactivate(t *testing.T) {
	link := NewEnterpriseAccountant("ent_123", "acc_456", "user_789")
	originalStartDate := link.StartDate

	// Testar desativação válida
	endDate := time.Now().UTC().Add(24 * time.Hour)
	err := link.Deactivate(endDate)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if link.Status != StatusInactive {
		t.Errorf("Expected Status INACTIVE, got %s", link.Status)
	}
	if link.EndDate == nil || !link.EndDate.Equal(endDate) {
		t.Error("Expected EndDate to be set to provided date")
	}
	if !link.UpdatedAt.After(originalStartDate) {
		t.Error("Expected UpdatedAt to be updated")
	}

	// Testar desativação com data anterior ao início
	link2 := NewEnterpriseAccountant("ent_123", "acc_456", "user_789")
	invalidEndDate := link2.StartDate.Add(-24 * time.Hour)
	err = link2.Deactivate(invalidEndDate)
	if err != ErrInvalidEndDate {
		t.Errorf("Expected ErrInvalidEndDate, got %v", err)
	}

	// Testar desativação já inativo
	link3 := NewEnterpriseAccountant("ent_123", "acc_456", "user_789")
	link3.Deactivate(time.Now().UTC().Add(24 * time.Hour))
	originalUpdatedAt := link3.UpdatedAt
	err = link3.Deactivate(time.Now().UTC().Add(48 * time.Hour))
	if err != nil {
		t.Errorf("Unexpected error when deactivating already inactive: %v", err)
	}
	if !link3.UpdatedAt.Equal(originalUpdatedAt) {
		t.Error("UpdatedAt should not change when deactivating already inactive")
	}
}

func TestEnterpriseAccountant_Reactivate(t *testing.T) {
	link := NewEnterpriseAccountant("ent_123", "acc_456", "user_789")
	link.Deactivate(time.Now().UTC().Add(24 * time.Hour))
	originalUpdatedAt := link.UpdatedAt

	link.Reactivate()
	if link.Status != StatusActive {
		t.Errorf("Expected Status ACTIVE, got %s", link.Status)
	}
	if link.EndDate != nil {
		t.Error("Expected EndDate to be nil after reactivation")
	}
	if !link.UpdatedAt.After(originalUpdatedAt) {
		t.Error("Expected UpdatedAt to be updated")
	}

	// Testar reativação já ativo
	link2 := NewEnterpriseAccountant("ent_123", "acc_456", "user_789")
	originalUpdatedAt2 := link2.UpdatedAt
	link2.Reactivate()
	if !link2.UpdatedAt.Equal(originalUpdatedAt2) {
		t.Error("UpdatedAt should not change when reactivating already active")
	}
}

func TestEnterpriseAccountant_IsValidForDate(t *testing.T) {
	now := time.Now().UTC()
	link := NewEnterpriseAccountant("ent_123", "acc_456", "user_789")
	link.StartDate = now.Add(-24 * time.Hour)

	// Testar data dentro do período ativo
	if !link.IsValidForDate(now) {
		t.Error("Expected date to be valid for active link")
	}

	// Testar data anterior ao início
	if link.IsValidForDate(now.Add(-48 * time.Hour)) {
		t.Error("Expected date before start to be invalid")
	}

	// Testar link inativo com data dentro do período
	link.Deactivate(now.Add(24 * time.Hour))
	if !link.IsValidForDate(now) {
		t.Error("Expected date to be valid for inactive link within period")
	}

	// Testar link inativo com data após o término
	if link.IsValidForDate(now.Add(48 * time.Hour)) {
		t.Error("Expected date after end to be invalid for inactive link")
	}
}

func TestEnterpriseAccountant_GetDateRange(t *testing.T) {
	now := time.Now().UTC()

	// Testar link ativo
	link := NewEnterpriseAccountant("ent_123", "acc_456", "user_789")
	link.StartDate = now.Add(-24 * time.Hour)
	start, end := link.GetDateRange()
	if !start.Equal(link.StartDate) {
		t.Error("Start date mismatch")
	}
	if !end.After(now) {
		t.Error("End date for active link should be current time")
	}

	// Testar link inativo
	link2 := NewEnterpriseAccountant("ent_123", "acc_456", "user_789")
	link2.StartDate = now.Add(-48 * time.Hour)
	endDate := now.Add(-24 * time.Hour)
	link2.Deactivate(endDate)
	start2, end2 := link2.GetDateRange()
	if !start2.Equal(link2.StartDate) {
		t.Error("Start date mismatch for inactive link")
	}
	if !end2.Equal(endDate) {
		t.Error("End date mismatch for inactive link")
	}
}

func TestEnterpriseAccountant_StatusMethods(t *testing.T) {
	link := NewEnterpriseAccountant("ent_123", "acc_456", "user_789")

	if !link.IsActive() {
		t.Error("Expected IsActive to return true for new link")
	}
	if link.IsInactive() {
		t.Error("Expected IsInactive to return false for new link")
	}

	link.Deactivate(time.Now().UTC().Add(24 * time.Hour))

	if link.IsActive() {
		t.Error("Expected IsActive to return false for inactive link")
	}
	if !link.IsInactive() {
		t.Error("Expected IsInactive to return true for inactive link")
	}
}
