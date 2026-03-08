// Package domain define as entidades e interfaces para distribuição de sobras
// seguindo a ITG 2002 e Lei Paul Singer
package domain

import (
	"context"
	"time"
)

// DistributionStatus representa o status da distribuição
type DistributionStatus string

const (
	StatusCalculated      DistributionStatus = "CALCULATED"       // Calculada, aguardando assembleia
	StatusPendingApproval DistributionStatus = "PENDING_APPROVAL" // Aguardando decisão da assembleia
	StatusApproved        DistributionStatus = "APPROVED"         // Aprovada em assembleia
	StatusExecuted        DistributionStatus = "EXECUTED"         // Executada (lançamentos gerados)
	StatusRejected        DistributionStatus = "REJECTED"         // Rejeitada em assembleia
)

// Distribution representa uma distribuição de sobras aprovada
type Distribution struct {
	ID                 int64
	EntityID           string
	Period             string // "2026-03" ou "2026-Q1"
	Status             DistributionStatus
	TotalSurplus       int64 // Sobra bruta (centavos)
	ReservaLegal       int64 // 10% obrigatório
	FATES              int64 // 5% obrigatório
	Distribuivel       int64 // 85% disponível para rateio
	AssemblyDecisionID int64 // Referência à decisão de assembleia
	CreatedAt          time.Time
	ExecutedAt         *time.Time
	Members            []DistributionMember
}

// DistributionMember representa a participação de cada membro
type DistributionMember struct {
	ID             int64
	DistributionID int64
	MemberID       string
	Minutes        int64
	Percentage     float64 // % das horas totais
	Amount         int64   // Valor a receber (proporcional)
	LedgerEntryID  int64   // Referência ao lançamento contábil
}

// SurplusCalculation representa o cálculo de sobras antes da aprovação
type SurplusCalculation struct {
	EntityID     string
	Period       string
	TotalSurplus int64
	ReservaLegal int64
	FATES        int64
	Distribuivel int64
	TotalMinutes int64
	Members      []MemberCalculation
}

// MemberCalculation representa o cálculo para cada membro
type MemberCalculation struct {
	MemberID   string
	Minutes    int64
	Percentage float64
	Amount     int64
}

// DistributionRepository interface para persistência
type DistributionRepository interface {
	// CRUD básico
	Save(ctx context.Context, dist *Distribution) (int64, error)
	FindByID(ctx context.Context, id int64) (*Distribution, error)
	FindByEntityAndPeriod(ctx context.Context, entityID, period string) (*Distribution, error)
	FindByStatus(ctx context.Context, entityID string, status DistributionStatus) ([]*Distribution, error)

	// Operações
	UpdateStatus(ctx context.Context, id int64, status DistributionStatus) error
	UpdateAssemblyDecisionID(ctx context.Context, id int64, decisionID int64) error
	MarkAsExecuted(ctx context.Context, id int64, executedAt time.Time) error
	SaveMember(ctx context.Context, member *DistributionMember) (int64, error)
	FindMembersByDistribution(ctx context.Context, distributionID int64) ([]DistributionMember, error)

	// Ledger integration
	SaveLedgerEntry(ctx context.Context, entityID string, entry *LedgerEntry) (int64, error)
}

// LedgerEntry representa um lançamento contábil
type LedgerEntry struct {
	ID          int64
	EntityID    string
	Reference   string
	Description string
	Date        time.Time
	Postings    []Posting
}

// Posting representa uma partida dobrada
type Posting struct {
	AccountID   int64
	Amount      int64 // positivo = crédito, negativo = débito
	Description string
}

// AssemblyDecision representa uma decisão de assembleia
type AssemblyDecision struct {
	ID           int64
	EntityID     string
	Title        string
	Content      string
	Status       string
	DecisionDate time.Time
}

// AssemblyRepository interface para acesso a decisões
type AssemblyRepository interface {
	FindByID(ctx context.Context, id int64) (*AssemblyDecision, error)
	IsApproved(ctx context.Context, id int64) bool
}
