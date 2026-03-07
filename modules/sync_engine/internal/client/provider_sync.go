package client

import (
	"crypto/sha256"
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
	syncRepo         SyncRepository
	providerEndpoint string
}

func NewProviderSyncClient(lm lifecycle.LifecycleManager, endpoint string) *ProviderSyncClient {
	return &ProviderSyncClient{
		syncRepo:         NewSQLiteSyncRepository(lm),
		providerEndpoint: endpoint,
	}
}

func (psc *ProviderSyncClient) BuildSyncPackage(entityID string) (*SyncPackage, error) {
	state, err := psc.syncRepo.GetSyncState(entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sync state: %w", err)
	}

	metrics, err := psc.calculateAggregatedMetrics(entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate metrics: %w", err)
	}

	chainDigest, err := psc.calculateChainDigest(entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate chain digest: %w", err)
	}

	deltaCount, err := psc.countDeltasSince(entityID, state)
	if err != nil {
		return nil, fmt.Errorf("failed to count deltas: %w", err)
	}

	pkg := &SyncPackage{
		EntityID:       entityID,
		Timestamp:      time.Now().Unix(),
		PeriodStart:    state,
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

	err = psc.syncRepo.UpdateLastSync(entityID, time.Now().Unix())
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

func (psc *ProviderSyncClient) calculateAggregatedMetrics(entityID string) (*AggregatedMetrics, error) {
	metrics := &AggregatedMetrics{
		LegalStatus: "DREAM",
	}

	totalSales, err := psc.syncRepo.GetTotalSales(entityID)
	if err == nil {
		metrics.TotalSales = totalSales
	}

	totalHours, totalMembers, err := psc.syncRepo.GetTotalWorkAndMembers(entityID)
	if err == nil {
		metrics.TotalWorkHours = totalHours
		metrics.TotalMembers = totalMembers
	}

	decisionCount, err := psc.syncRepo.GetDecisionCount(entityID)
	if err == nil {
		metrics.DecisionCount = decisionCount
	}

	if decisionCount >= 3 {
		metrics.LegalStatus = "FORMALIZED"
	}

	return metrics, nil
}

func (psc *ProviderSyncClient) calculateChainDigest(entityID string) (string, error) {
	lastEntryRef, err := psc.syncRepo.GetLastEntryRef(entityID)
	if err != nil {
		return "", err
	}

	lastDecisionHash, err := psc.syncRepo.GetLastDecisionHash(entityID)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256([]byte(lastEntryRef + lastDecisionHash))
	return hex.EncodeToString(hash[:])[:16], nil
}

func (psc *ProviderSyncClient) countDeltasSince(entityID string, since int64) (int64, error) {
	entriesCount, err := psc.syncRepo.GetEntriesCountSince(entityID, since)
	if err != nil {
		return 0, err
	}

	workCount, err := psc.syncRepo.GetWorkLogsCountSince(entityID, since)
	if err != nil {
		return 0, err
	}

	decisionCount, err := psc.syncRepo.GetDecisionsCountSince(entityID, since)
	if err != nil {
		return 0, err
	}

	return entriesCount + workCount + decisionCount, nil
}

func (psc *ProviderSyncClient) signPackage(pkg *SyncPackage, entityID string) string {
	data := fmt.Sprintf("%s:%d:%s:%d", pkg.EntityID, pkg.Timestamp, pkg.ChainDigest, pkg.DeltaCount)
	hash := sha256.Sum256([]byte(data + entityID))
	return hex.EncodeToString(hash[:])[:16]
}
