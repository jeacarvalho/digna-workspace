package sync_engine_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
	"github.com/providentia/digna/pdv_ui/usecase"
	"github.com/providentia/digna/sync_engine/internal/client"
	"github.com/providentia/digna/sync_engine/internal/exchange"
	"github.com/providentia/digna/sync_engine/internal/tracker"
)

func TestSprint04_DoD(t *testing.T) {
	dataDir := "../../data/entities"
	defer os.RemoveAll(dataDir)

	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	entityID := "Entidade_Sync_Teste"
	socio1ID := "socio_sync_001"
	socio2ID := "socio_sync_002"

	opHandler := usecase.NewOperationHandler(lifecycleMgr)
	sgHandler := usecase.NewSocialGovernanceHandler(lifecycleMgr)
	deltaTracker := tracker.NewDeltaTracker(lifecycleMgr)
	syncClient := client.NewProviderSyncClient(lifecycleMgr, "https://providentia.cloud/sync")
	intercoopService := exchange.NewIntercoopService(lifecycleMgr)

	t.Run("Step1_PDV_Operation", func(t *testing.T) {
		saleReq := usecase.SaleRequest{
			EntityID:      entityID,
			Amount:        7500,
			PaymentMethod: "PIX",
			Description:   "Venda para teste de sincronização",
		}

		result, err := opHandler.RecordSale(saleReq)
		if err != nil {
			t.Fatalf("failed to record sale: %v", err)
		}

		if result.EntryID == 0 {
			t.Error("expected entry ID to be set")
		}

		t.Logf("✅ PDV Operation: Sale recorded with EntryID=%d, Amount=%d", result.EntryID, saleReq.Amount)
	})

	t.Run("Step2_Register_Work_Hours", func(t *testing.T) {
		work1 := usecase.WorkRequest{
			EntityID:     entityID,
			MemberID:     socio1ID,
			Minutes:      480,
			ActivityType: "PRODUCAO",
			Description:  "Trabalho cooperativo socio 1",
		}
		work2 := usecase.WorkRequest{
			EntityID:     entityID,
			MemberID:     socio2ID,
			Minutes:      240,
			ActivityType: "PRODUCAO",
			Description:  "Trabalho cooperativo socio 2",
		}

		if err := sgHandler.RecordWork(work1); err != nil {
			t.Fatalf("failed to record work for %s: %v", socio1ID, err)
		}
		if err := sgHandler.RecordWork(work2); err != nil {
			t.Fatalf("failed to record work for %s: %v", socio2ID, err)
		}

		t.Logf("✅ Work Hours: %s (480min), %s (240min)", socio1ID, socio2ID)
	})

	t.Run("Step3_Detect_Deltas", func(t *testing.T) {
		state, err := deltaTracker.GetCurrentState(entityID)
		if err != nil {
			t.Fatalf("failed to get current state: %v", err)
		}

		t.Logf("Current State - Pending Changes: %d, Chain Digest: %s...", state.PendingChanges, state.ChainDigest[:8])

		hasChanges, err := deltaTracker.HasChanges(entityID)
		if err != nil {
			t.Fatalf("failed to check changes: %v", err)
		}

		if !hasChanges {
			t.Error("expected changes to be detected after PDV operations")
		}

		deltas, err := deltaTracker.DetectDeltas(entityID, 0)
		if err != nil {
			t.Fatalf("failed to detect deltas: %v", err)
		}

		if len(deltas) == 0 {
			t.Error("expected deltas to be detected")
		}

		var entryDeltas, workDeltas int
		for _, d := range deltas {
			if d.TableName == "entries" {
				entryDeltas++
			}
			if d.TableName == "work_logs" {
				workDeltas++
			}
		}

		if entryDeltas == 0 {
			t.Error("expected at least one entry delta (sale)")
		}

		if workDeltas < 2 {
			t.Errorf("expected at least 2 work log deltas, got %d", workDeltas)
		}

		t.Logf("✅ Deltas Detected: %d total (%d entries, %d work logs)", len(deltas), entryDeltas, workDeltas)
	})

	t.Run("Step4_Generate_Sync_Package", func(t *testing.T) {
		pkg, err := syncClient.BuildSyncPackage(entityID)
		if err != nil {
			t.Fatalf("failed to build sync package: %v", err)
		}

		if pkg.EntityID != entityID {
			t.Errorf("expected entity_id %s, got %s", entityID, pkg.EntityID)
		}

		if pkg.DeltaCount == 0 {
			t.Error("expected delta_count > 0")
		}

		if pkg.ChainDigest == "" {
			t.Error("expected chain_digest to be set")
		}

		if pkg.Signature == "" {
			t.Error("expected signature to be set")
		}

		t.Logf("Sync Package: EntityID=%s, Deltas=%d, Chain=%s..., Sig=%s...",
			pkg.EntityID, pkg.DeltaCount, pkg.ChainDigest[:8], pkg.Signature[:8])

		t.Logf("Aggregated Metrics: Sales=%d, WorkHours=%d, Members=%d, Status=%s",
			pkg.AggregatedData.TotalSales,
			pkg.AggregatedData.TotalWorkHours,
			pkg.AggregatedData.TotalMembers,
			pkg.AggregatedData.LegalStatus)

		jsonData, err := json.MarshalIndent(pkg, "", "  ")
		if err != nil {
			t.Fatalf("failed to marshal package: %v", err)
		}

		if !strings.Contains(string(jsonData), "entity_id") {
			t.Error("JSON package missing entity_id field")
		}

		if !strings.Contains(string(jsonData), "aggregated_data") {
			t.Error("JSON package missing aggregated_data field")
		}

		t.Logf("✅ Sync Package Generated (%d bytes)", len(jsonData))
		t.Logf("JSON Preview:\n%s", string(jsonData[:500]))
	})

	t.Run("Step5_Push_Sync_Package", func(t *testing.T) {
		jsonStr, err := syncClient.PushSyncPackage(entityID)
		if err != nil {
			t.Fatalf("failed to push sync package: %v", err)
		}

		if len(jsonStr) == 0 {
			t.Error("expected non-empty JSON")
		}

		var pkg client.SyncPackage
		if err := json.Unmarshal(jsonStr, &pkg); err != nil {
			t.Fatalf("failed to unmarshal package: %v", err)
		}

		t.Logf("✅ Sync Package Pushed: %d bytes ready for transport", len(jsonStr))
	})

	t.Run("Step6_Intercoop_Marketplace", func(t *testing.T) {
		offer1, err := intercoopService.PublishOffer(
			entityID,
			"Mel Orgânico",
			100,
			2500,
			"Mel puro de flores silvestres, 500g",
		)
		if err != nil {
			t.Fatalf("failed to publish offer 1: %v", err)
		}

		offer2, err := intercoopService.PublishOffer(
			"cooperativa_b",
			"Café Especial",
			50,
			3500,
			"Café arábica torrado, 250g",
		)
		if err != nil {
			t.Fatalf("failed to publish offer 2: %v", err)
		}

		if offer1.ID == "" {
			t.Error("expected offer1 to have an ID")
		}

		if offer2.ID == "" {
			t.Error("expected offer2 to have an ID")
		}

		allOffers := intercoopService.DiscoverOffers("")
		if len(allOffers) < 2 {
			t.Errorf("expected at least 2 active offers, got %d", len(allOffers))
		}

		melOffers := intercoopService.DiscoverOffers("Mel Orgânico")
		if len(melOffers) != 1 {
			t.Errorf("expected 1 'Mel Orgânico' offer, got %d", len(melOffers))
		}

		entityOffers := intercoopService.GetEntityOffers(entityID)
		if len(entityOffers) != 1 {
			t.Errorf("expected 1 offer for entity %s, got %d", entityID, len(entityOffers))
		}

		t.Logf("✅ Intercooperação: %d ofertas ativas no marketplace", len(allOffers))
		t.Logf("   - Oferta 1: %s (%d unidades)", offer1.Product, offer1.Quantity)
		t.Logf("   - Oferta 2: %s (%d unidades)", offer2.Product, offer2.Quantity)
	})

	t.Run("Step7_Validate_Privacy", func(t *testing.T) {
		jsonStr, err := syncClient.PushSyncPackage(entityID)
		if err != nil {
			t.Fatalf("failed to get sync package: %v", err)
		}

		sensitiveFields := []string{
			"member_id",
			"member_name",
			"personal_data",
			"entry_details",
			"posting_id",
		}

		for _, field := range sensitiveFields {
			if strings.Contains(string(jsonStr), field) {
				t.Errorf("package contains sensitive field '%s' - privacy violation", field)
			}
		}

		requiredFields := []string{
			"entity_id",
			"timestamp",
			"chain_digest",
			"signature",
			"aggregated_data",
			"total_sales",
			"total_work_hours",
			"legal_status",
		}

		for _, field := range requiredFields {
			if !strings.Contains(string(jsonStr), field) {
				t.Errorf("package missing required field '%s'", field)
			}
		}

		t.Log("✅ Privacy Validation: Only aggregated data in sync package")
		t.Log("   - No sensitive member data exposed")
		t.Log("   - Only totals: Sales, Work Hours, Members Count, Status")
	})
}

func TestSync_EmptyEntity(t *testing.T) {
	dataDir := "../../data/entities"
	defer os.RemoveAll(dataDir)

	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	entityID := "Empty_Entity"
	deltaTracker := tracker.NewDeltaTracker(lifecycleMgr)

	hasChanges, err := deltaTracker.HasChanges(entityID)
	if err != nil {
		t.Fatalf("failed to check changes: %v", err)
	}

	if hasChanges {
		t.Error("empty entity should have no changes")
	}

	state, err := deltaTracker.GetCurrentState(entityID)
	if err != nil {
		t.Fatalf("failed to get state: %v", err)
	}

	if state.PendingChanges != 0 {
		t.Errorf("expected 0 pending changes, got %d", state.PendingChanges)
	}
}
