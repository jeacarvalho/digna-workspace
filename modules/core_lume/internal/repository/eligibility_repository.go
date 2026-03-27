package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/providentia/digna/core_lume/internal/domain"
	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

// EligibilityRepository interface for eligibility profile operations
type EligibilityRepository interface {
	Save(profile *domain.EligibilityProfile) error
	FindByEntityID(entityID string) (*domain.EligibilityProfile, error)
	ListIncomplete() ([]*domain.EligibilityProfile, error)
	UpdateFields(entityID string, fields map[string]interface{}) error
	InitTable(entityID string) error
}

// SQLiteEligibilityRepository implements EligibilityRepository for SQLite
type SQLiteEligibilityRepository struct {
	lifecycleManager lifecycle.LifecycleManager
}

// NewSQLiteEligibilityRepository creates a new SQLiteEligibilityRepository
func NewSQLiteEligibilityRepository(lm lifecycle.LifecycleManager) *SQLiteEligibilityRepository {
	return &SQLiteEligibilityRepository{
		lifecycleManager: lm,
	}
}

// GetDB gets database connection for entity
func (r *SQLiteEligibilityRepository) GetDB(entityID string) (*sql.DB, error) {
	return r.lifecycleManager.GetConnection(entityID)
}

// Save creates or updates an eligibility profile (UPSERT)
func (r *SQLiteEligibilityRepository) Save(profile *domain.EligibilityProfile) error {
	if err := profile.Validate(); err != nil {
		return fmt.Errorf("invalid eligibility profile: %w", err)
	}

	db, err := r.GetDB(profile.EntityID)
	if err != nil {
		return fmt.Errorf("failed to get connection: %w", err)
	}

	// Convert bools to integers for SQLite
	inscritoCadUnico := 0
	if profile.InscritoCadUnico {
		inscritoCadUnico = 1
	}
	socioMulher := 0
	if profile.SocioMulher {
		socioMulher = 1
	}
	inadimplenciaAtiva := 0
	if profile.InadimplenciaAtiva {
		inadimplenciaAtiva = 1
	}
	contabilidadeFormal := 0
	if profile.ContabilidadeFormal {
		contabilidadeFormal = 1
	}

	_, err = db.Exec(`
		INSERT INTO eligibility_profiles (
			id, entity_id, cnpj, cnae, municipio, uf, faturamento_anual,
			regime_tributario, data_abertura, situacao_fiscal,
			inscrito_cad_unico, socio_mulher, inadimplencia_ativa,
			finalidade_credito, valor_necessario, tipo_entidade, contabilidade_formal,
			preenchido_em, atualizado_em, preenchido_por, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(entity_id) DO UPDATE SET
			cnpj = excluded.cnpj,
			cnae = excluded.cnae,
			municipio = excluded.municipio,
			uf = excluded.uf,
			faturamento_anual = excluded.faturamento_anual,
			regime_tributario = excluded.regime_tributario,
			data_abertura = excluded.data_abertura,
			situacao_fiscal = excluded.situacao_fiscal,
			inscrito_cad_unico = excluded.inscrito_cad_unico,
			socio_mulher = excluded.socio_mulher,
			inadimplencia_ativa = excluded.inadimplencia_ativa,
			finalidade_credito = excluded.finalidade_credito,
			valor_necessario = excluded.valor_necessario,
			tipo_entidade = excluded.tipo_entidade,
			contabilidade_formal = excluded.contabilidade_formal,
			atualizado_em = excluded.atualizado_em,
			updated_at = excluded.updated_at
	`,
		profile.ID, profile.EntityID, profile.CNPJ, profile.CNAE, profile.Municipio, profile.UF,
		profile.FaturamentoAnual, profile.RegimeTributario, profile.DataAbertura, profile.SituacaoFiscal,
		inscritoCadUnico, socioMulher, inadimplenciaAtiva,
		string(profile.FinalidadeCredito), profile.ValorNecessario, string(profile.TipoEntidade), contabilidadeFormal,
		profile.PreenchidoEm, profile.AtualizadoEm, profile.PreenchidoPor, profile.CreatedAt, profile.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save eligibility profile: %w", err)
	}

	return nil
}

// FindByEntityID finds eligibility profile by entity ID
func (r *SQLiteEligibilityRepository) FindByEntityID(entityID string) (*domain.EligibilityProfile, error) {
	db, err := r.GetDB(entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}

	var profile domain.EligibilityProfile
	var finalidadeCredito, tipoEntidade string
	var inscritoCadUnico, socioMulher, inadimplenciaAtiva, contabilidadeFormal int

	err = db.QueryRow(`
		SELECT id, entity_id, cnpj, cnae, municipio, uf, faturamento_anual,
			regime_tributario, data_abertura, situacao_fiscal,
			inscrito_cad_unico, socio_mulher, inadimplencia_ativa,
			finalidade_credito, valor_necessario, tipo_entidade, contabilidade_formal,
			preenchido_em, atualizado_em, preenchido_por, created_at, updated_at
		FROM eligibility_profiles WHERE entity_id = ?
	`, entityID).Scan(
		&profile.ID, &profile.EntityID, &profile.CNPJ, &profile.CNAE, &profile.Municipio, &profile.UF,
		&profile.FaturamentoAnual, &profile.RegimeTributario, &profile.DataAbertura, &profile.SituacaoFiscal,
		&inscritoCadUnico, &socioMulher, &inadimplenciaAtiva,
		&finalidadeCredito, &profile.ValorNecessario, &tipoEntidade, &contabilidadeFormal,
		&profile.PreenchidoEm, &profile.AtualizadoEm, &profile.PreenchidoPor, &profile.CreatedAt, &profile.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrProfileNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query eligibility profile: %w", err)
	}

	// Convert integers back to bools
	profile.InscritoCadUnico = inscritoCadUnico == 1
	profile.SocioMulher = socioMulher == 1
	profile.InadimplenciaAtiva = inadimplenciaAtiva == 1
	profile.ContabilidadeFormal = contabilidadeFormal == 1
	profile.FinalidadeCredito = domain.FinalidadeCredito(finalidadeCredito)
	profile.TipoEntidade = domain.TipoEntidade(tipoEntidade)

	return &profile, nil
}

// ListIncomplete lists all incomplete eligibility profiles
func (r *SQLiteEligibilityRepository) ListIncomplete() ([]*domain.EligibilityProfile, error) {
	// This requires iterating over all entities, which is complex with the current architecture
	// For now, return empty list - in production, this would need a different approach
	return []*domain.EligibilityProfile{}, nil
}

// UpdateFields updates specific fields of an eligibility profile
func (r *SQLiteEligibilityRepository) UpdateFields(entityID string, fields map[string]interface{}) error {
	db, err := r.GetDB(entityID)
	if err != nil {
		return fmt.Errorf("failed to get connection: %w", err)
	}

	// Build dynamic query
	if len(fields) == 0 {
		return nil
	}

	query := "UPDATE eligibility_profiles SET "
	args := []interface{}{}
	i := 0

	for field, value := range fields {
		if i > 0 {
			query += ", "
		}
		query += field + " = ?"
		args = append(args, value)
		i++
	}

	query += ", updated_at = ? WHERE entity_id = ?"
	args = append(args, time.Now().Unix(), entityID)

	_, err = db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update fields: %w", err)
	}

	return nil
}

// InitTable creates the eligibility_profiles table if not exists
func (r *SQLiteEligibilityRepository) InitTable(entityID string) error {
	db, err := r.GetDB(entityID)
	if err != nil {
		return fmt.Errorf("failed to get connection: %w", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS eligibility_profiles (
			id TEXT PRIMARY KEY,
			entity_id TEXT NOT NULL UNIQUE,
			
			-- Dados do ERP (cópia para consulta rápida)
			cnpj TEXT,
			cnae TEXT,
			municipio TEXT,
			uf TEXT,
			faturamento_anual INTEGER,
			regime_tributario TEXT,
			data_abertura INTEGER,
			situacao_fiscal TEXT,
			
			-- Campos complementares
			inscrito_cad_unico INTEGER DEFAULT 0,
			socio_mulher INTEGER DEFAULT 0,
			inadimplencia_ativa INTEGER DEFAULT 0,
			finalidade_credito TEXT,
			valor_necessario INTEGER DEFAULT 0,
			tipo_entidade TEXT,
			contabilidade_formal INTEGER DEFAULT 0,
			
			-- Metadados
			preenchido_em INTEGER,
			atualizado_em INTEGER,
			preenchido_por TEXT,
			created_at INTEGER,
			updated_at INTEGER
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create eligibility_profiles table: %w", err)
	}

	// Create index for faster lookups
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_eligibility_entity ON eligibility_profiles(entity_id)
	`)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	return nil
}
