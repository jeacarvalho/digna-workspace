package service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"time"

	"digna/accountant_dashboard/internal/domain"
	"digna/accountant_dashboard/internal/middleware"
)

type TranslatorService struct {
	repo   domain.FiscalRepository
	mapper domain.AccountMapper
}

func NewTranslatorService(repo domain.FiscalRepository, mapper domain.AccountMapper) *TranslatorService {
	return &TranslatorService{
		repo:   repo,
		mapper: mapper,
	}
}

type TranslationResult struct {
	Data       []byte
	Hash       string
	EntryCount int
}

func (s *TranslatorService) TranslateAndExport(ctx context.Context, entityID string, period string) (*domain.FiscalBatch, []byte, error) {
	// Check if accountant has access to this entity for the period
	validEnterprises, ok := middleware.GetValidEnterprisesFromContext(ctx)
	if ok {
		// We have temporal filtering context, validate access
		if !validEnterprises[entityID] {
			return nil, nil, fmt.Errorf("accountant does not have access to entity %s for period %s", entityID, period)
		}
	}

	entries, err := s.repo.LoadEntries(ctx, entityID, period)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load entries: %w", err)
	}

	if len(entries) == 0 {
		return nil, nil, fmt.Errorf("no entries found for period %s", period)
	}

	if err := s.validateSomaZero(entries); err != nil {
		return nil, nil, fmt.Errorf("audit validation failed: %w", err)
	}

	data, err := s.TranslateToStandardFormat(entries)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to translate format: %w", err)
	}

	hash := s.GenerateHash(data)

	batch := &domain.FiscalBatch{
		ID:           generateBatchID(entityID, period),
		EntityID:     entityID,
		Period:       period,
		TotalEntries: len(entries),
		ExportHash:   hash,
		CreatedAt:    time.Now().Unix(),
	}

	if err := s.repo.RegisterExport(ctx, entityID, batch); err != nil {
		return nil, nil, fmt.Errorf("failed to register export: %w", err)
	}

	return batch, data, nil
}

func (s *TranslatorService) validateSomaZero(entries []domain.EntryDTO) error {
	for _, entry := range entries {
		if entry.TotalDebit != entry.TotalCredit {
			return fmt.Errorf("entry %d has invalid soma zero: debit=%d, credit=%d",
				entry.ID, entry.TotalDebit, entry.TotalCredit)
		}
	}
	return nil
}

func (s *TranslatorService) TranslateToStandardFormat(entries []domain.EntryDTO) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	header := []string{
		"Data", "ID_Lancamento", "Conta_Debito", "Nome_Debito", "Conta_Credito", "Nome_Credito",
		"Valor", "Historico", "Hash_Entry",
	}
	if err := writer.Write(header); err != nil {
		return nil, fmt.Errorf("failed to write header: %w", err)
	}

	for _, entry := range entries {
		entryHash := generateEntryHash(entry)

		debitAccounts := []domain.PostingDTO{}
		creditAccounts := []domain.PostingDTO{}

		for _, p := range entry.Postings {
			if p.Debit > 0 {
				debitAccounts = append(debitAccounts, p)
			} else if p.Credit > 0 {
				creditAccounts = append(creditAccounts, p)
			}
		}

		maxLen := len(debitAccounts)
		if len(creditAccounts) > maxLen {
			maxLen = len(creditAccounts)
		}

		for i := 0; i < maxLen; i++ {
			var debitCode, debitName, creditCode, creditName string
			var amount int64

			if i < len(debitAccounts) {
				mapping := s.getAccountMapping(debitAccounts[i].AccountCode)
				debitCode = mapping.StandardCode
				debitName = mapping.StandardName
				amount = debitAccounts[i].Debit
			}

			if i < len(creditAccounts) {
				mapping := s.getAccountMapping(creditAccounts[i].AccountCode)
				creditCode = mapping.StandardCode
				creditName = mapping.StandardName
				if amount == 0 {
					amount = creditAccounts[i].Credit
				}
			}

			row := []string{
				entry.Date.Format("2006-01-02"),
				fmt.Sprintf("%d", entry.ID),
				debitCode,
				debitName,
				creditCode,
				creditName,
				fmt.Sprintf("%d", amount),
				entry.Description,
				entryHash,
			}

			if err := writer.Write(row); err != nil {
				return nil, fmt.Errorf("failed to write row: %w", err)
			}
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("flush error: %w", err)
	}

	return buf.Bytes(), nil
}

func (s *TranslatorService) getAccountMapping(localCode string) domain.AccountMapping {
	if mapping, ok := s.mapper.GetMapping(localCode); ok {
		return mapping
	}
	return domain.AccountMapping{
		LocalCode:    localCode,
		LocalName:    "Conta não mapeada",
		StandardCode: "9.9.99.99.99",
		StandardName: "Conta não mapeada",
	}
}

func (s *TranslatorService) GenerateHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func (s *TranslatorService) ListPendingEntities(ctx context.Context, period string) ([]string, error) {
	// Get all pending entities from repository
	allEntities, err := s.repo.ListPendingEntities(ctx, period)
	if err != nil {
		return nil, err
	}

	// Check if we have temporal filtering context
	validEnterprises, ok := middleware.GetValidEnterprisesFromContext(ctx)
	if !ok {
		// No filtering context, return all entities
		return allEntities, nil
	}

	// Filter entities based on valid enterprises
	var filteredEntities []string
	for _, entity := range allEntities {
		if validEnterprises[entity] {
			filteredEntities = append(filteredEntities, entity)
		}
	}

	return filteredEntities, nil
}

func (s *TranslatorService) GetExportHistory(ctx context.Context, entityID string, period string) ([]domain.FiscalExportLog, error) {
	// Check if accountant has access to this entity for the period
	validEnterprises, ok := middleware.GetValidEnterprisesFromContext(ctx)
	if ok {
		// We have temporal filtering context, validate access
		if !validEnterprises[entityID] {
			return nil, fmt.Errorf("accountant does not have access to entity %s for period %s", entityID, period)
		}
	}

	return s.repo.GetExportHistory(ctx, entityID, period)
}

func generateBatchID(entityID, period string) string {
	return fmt.Sprintf("%s_%s_%d", entityID, period, time.Now().Unix())
}

func generateEntryHash(entry domain.EntryDTO) string {
	data := fmt.Sprintf("%d|%s|%d|%d",
		entry.ID,
		entry.Date.Format("2006-01-02"),
		entry.TotalDebit,
		entry.TotalCredit)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
