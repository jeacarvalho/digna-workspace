package service

import (
	"context"
	"testing"
	"time"

	"github.com/providentia/digna/supply/internal/domain"
	"github.com/providentia/digna/supply/pkg/supply"
)

// MockRepository implementa repository.SupplyRepository para testes
type MockRepository struct {
	GetSupplierFunc             func(ctx context.Context, entityID, supplierID string) (*domain.Supplier, error)
	GetStockItemFunc            func(ctx context.Context, entityID, itemID string) (*domain.StockItem, error)
	BeginTxFunc                 func(ctx context.Context, entityID string) (interface{}, error)
	SavePurchaseFunc            func(ctx context.Context, entityID string, purchase *domain.Purchase) error
	CommitTxFunc                func(tx interface{}) error
	SaveStockItemFunc           func(ctx context.Context, entityID string, item *domain.StockItem) error
	SaveSupplierFunc            func(ctx context.Context, entityID string, supplier *domain.Supplier) error
	ListStockItemsFunc          func(ctx context.Context, entityID string) ([]*domain.StockItem, error)
	ListStockItemsByTypeFunc    func(ctx context.Context, entityID string, itemType domain.StockItemType) ([]*domain.StockItem, error)
	ListPurchasesFunc           func(ctx context.Context, entityID string) ([]*domain.Purchase, error)
	GetPurchaseFunc             func(ctx context.Context, entityID, purchaseID string) (*domain.Purchase, error)
	ListPurchasesBySupplierFunc func(ctx context.Context, entityID, supplierID string) ([]*domain.Purchase, error)
	UpdateStockQuantityFunc     func(ctx context.Context, entityID, itemID string, delta int) error
}

func (m *MockRepository) BeginTx(ctx context.Context, entityID string) (interface{}, error) {
	if m.BeginTxFunc != nil {
		return m.BeginTxFunc(ctx, entityID)
	}
	return nil, nil
}

func (m *MockRepository) CommitTx(tx interface{}) error {
	if m.CommitTxFunc != nil {
		return m.CommitTxFunc(tx)
	}
	return nil
}

func (m *MockRepository) RollbackTx(tx interface{}) error {
	return nil
}

func (m *MockRepository) SaveSupplier(ctx context.Context, entityID string, supplier *domain.Supplier) error {
	if m.SaveSupplierFunc != nil {
		return m.SaveSupplierFunc(ctx, entityID, supplier)
	}
	return nil
}

func (m *MockRepository) GetSupplier(ctx context.Context, entityID, supplierID string) (*domain.Supplier, error) {
	if m.GetSupplierFunc != nil {
		return m.GetSupplierFunc(ctx, entityID, supplierID)
	}
	return nil, nil
}

func (m *MockRepository) ListSuppliers(ctx context.Context, entityID string) ([]*domain.Supplier, error) {
	return nil, nil
}

func (m *MockRepository) SaveStockItem(ctx context.Context, entityID string, item *domain.StockItem) error {
	if m.SaveStockItemFunc != nil {
		return m.SaveStockItemFunc(ctx, entityID, item)
	}
	return nil
}

func (m *MockRepository) GetStockItem(ctx context.Context, entityID, itemID string) (*domain.StockItem, error) {
	if m.GetStockItemFunc != nil {
		return m.GetStockItemFunc(ctx, entityID, itemID)
	}
	return nil, nil
}

func (m *MockRepository) ListStockItems(ctx context.Context, entityID string) ([]*domain.StockItem, error) {
	if m.ListStockItemsFunc != nil {
		return m.ListStockItemsFunc(ctx, entityID)
	}
	return nil, nil
}

func (m *MockRepository) ListStockItemsByType(ctx context.Context, entityID string, itemType domain.StockItemType) ([]*domain.StockItem, error) {
	if m.ListStockItemsByTypeFunc != nil {
		return m.ListStockItemsByTypeFunc(ctx, entityID, itemType)
	}
	return nil, nil
}

func (m *MockRepository) SavePurchase(ctx context.Context, entityID string, purchase *domain.Purchase) error {
	if m.SavePurchaseFunc != nil {
		return m.SavePurchaseFunc(ctx, entityID, purchase)
	}
	return nil
}

func (m *MockRepository) ListPurchases(ctx context.Context, entityID string) ([]*domain.Purchase, error) {
	if m.ListPurchasesFunc != nil {
		return m.ListPurchasesFunc(ctx, entityID)
	}
	return nil, nil
}

func (m *MockRepository) GetPurchase(ctx context.Context, entityID, purchaseID string) (*domain.Purchase, error) {
	if m.GetPurchaseFunc != nil {
		return m.GetPurchaseFunc(ctx, entityID, purchaseID)
	}
	return nil, nil
}

func (m *MockRepository) ListPurchasesBySupplier(ctx context.Context, entityID, supplierID string) ([]*domain.Purchase, error) {
	if m.ListPurchasesBySupplierFunc != nil {
		return m.ListPurchasesBySupplierFunc(ctx, entityID, supplierID)
	}
	return nil, nil
}

func (m *MockRepository) UpdateStockQuantity(ctx context.Context, entityID, itemID string, delta int) error {
	if m.UpdateStockQuantityFunc != nil {
		return m.UpdateStockQuantityFunc(ctx, entityID, itemID, delta)
	}
	return nil
}

// MockLedgerPort implementa supply.LedgerPort para testes
type MockLedgerPort struct {
	RecordTransactionFunc func(entityID, description string, postings []supply.LedgerPosting) error
}

func (m *MockLedgerPort) RecordTransaction(entityID, description string, postings []supply.LedgerPosting) error {
	if m.RecordTransactionFunc != nil {
		return m.RecordTransactionFunc(entityID, description, postings)
	}
	return nil
}

func TestPurchaseService_RegisterPurchase(t *testing.T) {
	ctx := context.Background()
	entityID := "test-entity"

	// Dados de teste
	supplier := &domain.Supplier{
		ID:   "supplier-1",
		Name: "Fornecedor Teste",
	}

	stockItem := &domain.StockItem{
		ID:          "item-1",
		Name:        "Cera de Abelha",
		Type:        domain.StockItemTypeRawMaterial,
		Quantity:    0,
		MinQuantity: 5,
		UnitCost:    500, // R$ 5,00
	}

	purchase := &domain.Purchase{
		ID:         "purchase-1",
		SupplierID: "supplier-1",
		Date:       time.Now(),
		Items: []domain.PurchaseItem{
			domain.NewPurchaseItem("item-1", 10, 500), // R$ 5,00
		},
		TotalValue: 5000, // R$ 50,00
	}

	// Configurar mocks
	mockRepo := &MockRepository{
		GetSupplierFunc: func(ctx context.Context, entityID, supplierID string) (*domain.Supplier, error) {
			return supplier, nil
		},
		GetStockItemFunc: func(ctx context.Context, entityID, itemID string) (*domain.StockItem, error) {
			return stockItem, nil
		},
		BeginTxFunc: func(ctx context.Context, entityID string) (interface{}, error) {
			return "tx-1", nil
		},
		SavePurchaseFunc: func(ctx context.Context, entityID string, purchase *domain.Purchase) error {
			return nil
		},
		CommitTxFunc: func(tx interface{}) error {
			return nil
		},
	}

	var ledgerCalled bool
	mockLedger := &MockLedgerPort{
		RecordTransactionFunc: func(entityID, description string, postings []supply.LedgerPosting) error {
			ledgerCalled = true
			if len(postings) != 2 {
				t.Errorf("esperado 2 postings, obtido %d", len(postings))
			}
			return nil
		},
	}

	// Criar serviço
	service := NewPurchaseService(mockRepo, mockLedger)

	// Executar teste
	err := service.RegisterPurchase(ctx, entityID, purchase, "CASH")

	// Verificar
	if err != nil {
		t.Errorf("erro inesperado: %v", err)
	}
	if !ledgerCalled {
		t.Error("ledger não foi chamado")
	}
}

func TestPurchaseService_RegisterPurchase_SupplierNotFound(t *testing.T) {
	ctx := context.Background()
	entityID := "test-entity"

	mockRepo := &MockRepository{
		GetSupplierFunc: func(ctx context.Context, entityID, supplierID string) (*domain.Supplier, error) {
			return nil, nil // Supplier não encontrado
		},
	}

	mockLedger := &MockLedgerPort{}
	service := NewPurchaseService(mockRepo, mockLedger)

	purchase := &domain.Purchase{
		ID:         "purchase-1",
		SupplierID: "supplier-1",
		Date:       time.Now(),
		Items: []domain.PurchaseItem{
			domain.NewPurchaseItem("item-1", 10, 500),
		},
		TotalValue: 5000,
	}

	err := service.RegisterPurchase(ctx, entityID, purchase, "CASH")
	if err == nil {
		t.Error("esperado erro quando fornecedor não encontrado")
	}
	if err != nil && err.Error() != "fornecedor não encontrado: supplier-1" {
		t.Errorf("erro inesperado: %v", err)
	}
}

func TestPurchaseService_CreateStockItem(t *testing.T) {
	ctx := context.Background()
	entityID := "test-entity"

	var saveCalled bool
	mockRepo := &MockRepository{
		SaveStockItemFunc: func(ctx context.Context, entityID string, item *domain.StockItem) error {
			saveCalled = true
			if item.Name != "Test Item" {
				t.Errorf("nome do item incorreto: %s", item.Name)
			}
			return nil
		},
	}

	mockLedger := &MockLedgerPort{}
	service := NewPurchaseService(mockRepo, mockLedger)

	item := &domain.StockItem{
		ID:          "item-1",
		Name:        "Test Item",
		Type:        domain.StockItemTypeProduct,
		Quantity:    10,
		MinQuantity: 5,
		UnitCost:    1000,
	}

	err := service.CreateStockItem(ctx, entityID, item)
	if err != nil {
		t.Errorf("erro inesperado: %v", err)
	}
	if !saveCalled {
		t.Error("SaveStockItem não foi chamado")
	}
}

func TestPurchaseService_CreateSupplier(t *testing.T) {
	ctx := context.Background()
	entityID := "test-entity"

	var saveCalled bool
	mockRepo := &MockRepository{
		SaveSupplierFunc: func(ctx context.Context, entityID string, supplier *domain.Supplier) error {
			saveCalled = true
			if supplier.Name != "Test Supplier" {
				t.Errorf("nome do fornecedor incorreto: %s", supplier.Name)
			}
			return nil
		},
	}

	mockLedger := &MockLedgerPort{}
	service := NewPurchaseService(mockRepo, mockLedger)

	supplier := &domain.Supplier{
		ID:   "supplier-1",
		Name: "Test Supplier",
	}

	err := service.CreateSupplier(ctx, entityID, supplier)
	if err != nil {
		t.Errorf("erro inesperado: %v", err)
	}
	if !saveCalled {
		t.Error("SaveSupplier não foi chamado")
	}
}

func TestDetermineDebitAccount(t *testing.T) {
	// Teste com apenas insumos
	itemTypes1 := map[string]domain.StockItemType{
		"item-1": domain.StockItemTypeRawMaterial,
		"item-2": domain.StockItemTypeRawMaterial,
	}
	if got := determineDebitAccount(itemTypes1); got != AccountInventory {
		t.Errorf("determineDebitAccount() = %d, esperado %d", got, AccountInventory)
	}

	// Teste com apenas mercadorias
	itemTypes2 := map[string]domain.StockItemType{
		"item-1": domain.StockItemTypeMerchandise,
	}
	if got := determineDebitAccount(itemTypes2); got != AccountInventory {
		t.Errorf("determineDebitAccount() = %d, esperado %d", got, AccountInventory)
	}

	// Teste com apenas produtos
	itemTypes3 := map[string]domain.StockItemType{
		"item-1": domain.StockItemTypeProduct,
	}
	if got := determineDebitAccount(itemTypes3); got != AccountInventory {
		t.Errorf("determineDebitAccount() = %d, esperado %d", got, AccountInventory)
	}

	// Teste com mistura de tipos
	itemTypes4 := map[string]domain.StockItemType{
		"item-1": domain.StockItemTypeRawMaterial,
		"item-2": domain.StockItemTypeProduct,
	}
	if got := determineDebitAccount(itemTypes4); got != AccountExpenses {
		t.Errorf("determineDebitAccount() = %d, esperado %d", got, AccountExpenses)
	}
}

func TestDetermineCreditAccount(t *testing.T) {
	tests := []struct {
		name        string
		paymentType string
		want        int64
	}{
		{"CASH", "CASH", AccountCash},
		{"À VISTA", "À VISTA", AccountCash},
		{"CREDIT", "CREDIT", AccountSuppliers},
		{"A PRAZO", "A PRAZO", AccountSuppliers},
		{"default", "OUTRO", AccountCash},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := determineCreditAccount(tt.paymentType); got != tt.want {
				t.Errorf("determineCreditAccount() = %d, esperado %d", got, tt.want)
			}
		})
	}
}
