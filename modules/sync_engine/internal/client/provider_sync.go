package client

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

type SyncPackage struct {
	EntityID       string            `json:"entity_id"`
	Timestamp      int64             `json:"timestamp"`
	PeriodStart    int64             `json:"period_start"`
	PeriodEnd      int64             `json:"period_end"`
	ChainDigest    string            `json:"chain_digest"`
	Signature      string            `json:"signature"`
	AggregatedData AggregatedMetrics `json:"aggregated_data"`
	DeltaCount     int64             `json:"delta_count"`
}

type AggregatedMetrics struct {
	TotalSales     int64  `json:"total_sales"`
	TotalWorkHours int64  `json:"total_work_hours"`
	TotalMembers   int64  `json:"total_members"`
	LegalStatus    string `json:"legal_status"`
	ActiveOffers   int64  `json:"active_offers"`
	DecisionCount  int64  `json:"decision_count"`
}

type ProviderSyncClient struct {
	lifecycleManager lifecycle.LifecycleManager
	providerEndpoint string
}

func NewProviderSyncClient(lm lifecycle.LifecycleManager, endpoint string) *ProviderSyncClient {
	return &ProviderSyncClient{
		lifecycleManager: lm,
		providerEndpoint: endpoint,
	}
}

func (psc *ProviderSyncClient) BuildSyncPackage(entityID string) (*SyncPackage, error) {
	db, err := psc.lifecycleManager.GetConnection(entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}

	state, err := psc.getSyncState(db)
	if err != nil {
		return nil, fmt.Errorf("failed to get sync state: %w", err)
	}

	metrics, err := psc.calculateAggregatedMetrics(db)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate metrics: %w", err)
	}

	chainDigest, err := psc.calculateChainDigest(db)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate chain digest: %w", err)
	}

	deltaCount, err := psc.countDeltasSince(db, state.LastSyncAt)
	if err != nil {
		return nil, fmt.Errorf("failed to count deltas: %w", err)
	}

	pkg := &SyncPackage{
		EntityID:       entityID,
		Timestamp:      time.Now().Unix(),
		PeriodStart:    state.LastSyncAt,
		PeriodEnd:      time.Now().Unix(),
		ChainDigest:    chainDigest,
		AggregatedData: *metrics,
		DeltaCount:     deltaCount,
	}

	pkg.Signature = psc.signPackage(pkg, entityID)

	return pkg, nil
}

func (psc *ProviderSyncClient) PushSyncPackage(entityID string) ([]byte, error) {
	pkg, err := psc.BuildSyncPackage(entityID)
	if err != nil {
		return nil, err
	}

	jsonData, err := json.MarshalIndent(pkg, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal package: %w", err)
	}

	db, err := psc.lifecycleManager.GetConnection(entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}

	_, err = db.Exec(
		"UPDATE sync_metadata SET last_sync_at = ? WHERE id = 1",
		time.Now().Unix(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update sync timestamp: %w", err)
	}

	return jsonData, nil
}

func (psc *ProviderSyncClient) GetPackageForTransport(entityID string) (string, error) {
	jsonData, err := psc.PushSyncPackage(entityID)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func (psc *ProviderSyncClient) getSyncState(db *sql.DB) (*syncState, error) {
	state := &syncState{}

	var lastSyncAt sql.NullInt64
	err := db.QueryRow(
		"SELECT last_sync_at FROM sync_metadata WHERE id = 1",
	).Scan(&lastSyncAt)
	if err == nil && lastSyncAt.Valid {
		state.LastSyncAt = lastSyncAt.Int64
	}

	return state, nil
}

func (psc *ProviderSyncClient) calculateAggregatedMetrics(db *sql.DB) (*AggregatedMetrics, error) {
	metrics := &AggregatedMetrics{
		LegalStatus: "DREAM",
	}

	var totalSales sql.NullInt64
	err := db.QueryRow(
		"SELECT COALESCE(SUM(amount), 0) FROM postings WHERE direction = 'CREDIT' AND account_id = 2",
	).Scan(&totalSales)
	if err == nil {
		metrics.TotalSales = totalSales.Int64
	}

	var totalMinutes sql.NullInt64
	var memberCount sql.NullInt64
	err = db.QueryRow(
		"SELECT COALESCE(SUM(minutes), 0), COUNT(DISTINCT member_id) FROM work_logs",
	).Scan(&totalMinutes, &memberCount)
	if err == nil {
		metrics.TotalWorkHours = totalMinutes.Int64 / 60
		metrics.TotalMembers = memberCount.Int64
	}

	var decisionCount sql.NullInt64
	err = db.QueryRow(
		"SELECT COUNT(*) FROM decisions_log",
	).Scan(&decisionCount)
	if err == nil {
		metrics.DecisionCount = decisionCount.Int64
	}

	if decisionCount.Int64 >= 3 {
		metrics.LegalStatus = "FORMALIZED"
	}

	return metrics, nil
}

func (psc *ProviderSyncClient) calculateChainDigest(db *sql.DB) (string, error) {
	var lastEntryRef string
	err := db.QueryRow(
		"SELECT COALESCE(MAX(reference), '') FROM entries",
	).Scan(&lastEntryRef)
	if err != nil {
		return "", err
	}

	var lastDecisionHash string
	err = db.QueryRow(
		"SELECT COALESCE(MAX(content_hash), '') FROM decisions_log",
	).Scan(&lastDecisionHash)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256([]byte(lastEntryRef + lastDecisionHash))
	return hex.EncodeToString(hash[:])[:16], nil
}

func (psc *ProviderSyncClient) countDeltasSince(db *sql.DB, since int64) (int64, error) {
	var count int64

	err := db.QueryRow(
		"SELECT COUNT(*) FROM entries WHERE created_at > ?",
		since,
	).Scan(&count)
	if err != nil {
		return 0, err
	}

	var workCount int64
	err = db.QueryRow(
		"SELECT COUNT(*) FROM work_logs WHERE created_at > ?",
		since,
	).Scan(&workCount)
	if err != nil {
		return 0, err
	}

	var decisionCount int64
	err = db.QueryRow(
		"SELECT COUNT(*) FROM decisions_log WHERE created_at > ?",
		since,
	).Scan(&decisionCount)
	if err != nil {
		return 0, err
	}

	return count + workCount + decisionCount, nil
}

func (psc *ProviderSyncClient) signPackage(pkg *SyncPackage, entityID string) string {
	data := fmt.Sprintf("%s:%d:%s:%d", pkg.EntityID, pkg.Timestamp, pkg.ChainDigest, pkg.DeltaCount)
	hash := sha256.Sum256([]byte(data + entityID))
	return hex.EncodeToString(hash[:])[:16]
}

type syncState struct {
	LastSyncAt int64
}
