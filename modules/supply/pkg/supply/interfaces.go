package supply

import (
	"context"
	"time"
)

// LedgerPosting representa um posting contábil
type LedgerPosting struct {
	AccountID int64
	Amount    int64
	Direction string // "DEBIT" ou "CREDIT"
}

// LedgerPort interface para integração com o core_lume
type LedgerPort interface {
	RecordTransaction(entityID string, description string, postings []LedgerPosting) error
}

// PurchaseRequest representa uma requisição de compra
type PurchaseRequest struct {
	EntityID    string
	SupplierID  string
	Items       []PurchaseItemRequest
	PaymentType string // "CASH" ou "CREDIT"
}

type PurchaseItemRequest struct {
	StockItemID string
	Quantity    int
	UnitCost    int64
}

// PurchaseResponse representa a resposta de uma compra
type PurchaseResponse struct {
	PurchaseID string
	Success    bool
	Error      string
}

// StockItem representa um item de estoque para a API pública
type StockItem struct {
	ID          string
	Name        string
	Type        string // "INSUMO", "PRODUTO", "MERCADORIA"
	Quantity    int
	MinQuantity int
	UnitCost    int64
	CreatedAt   time.Time
}

// Supplier representa um fornecedor para a API pública
type Supplier struct {
	ID          string
	Name        string
	ContactInfo string
	CreatedAt   time.Time
}

// Purchase representa uma compra para a API pública
type Purchase struct {
	ID         string
	SupplierID string
	TotalValue int64
	Date       time.Time
	Items      []PurchaseItem
	CreatedAt  time.Time
}

type PurchaseItem struct {
	ID          string
	StockItemID string
	Quantity    int
	UnitCost    int64
	TotalCost   int64
}

// StockItemRequest representa uma requisição para criar/atualizar item de estoque
type StockItemRequest struct {
	EntityID    string
	Name        string
	Type        string // "INSUMO", "PRODUTO", "MERCADORIA"
	Quantity    int
	MinQuantity int
	UnitCost    int64
}

// SupplierRequest representa uma requisição para criar/atualizar fornecedor
type SupplierRequest struct {
	EntityID    string
	Name        string
	ContactInfo string
}

// SupplyAPI interface pública do módulo supply
type SupplyAPI interface {
	// Suppliers
	RegisterSupplier(ctx context.Context, req SupplierRequest) (*SupplierResponse, error)
	GetSuppliers(ctx context.Context, entityID string) ([]*Supplier, error)

	// Stock Items
	RegisterStockItem(ctx context.Context, req StockItemRequest) (*StockItemResponse, error)
	GetStockItems(ctx context.Context, entityID string) ([]*StockItem, error)
	GetStockItemsByType(ctx context.Context, entityID string, itemType string) ([]*StockItem, error)
	UpdateStockQuantity(ctx context.Context, entityID, itemID string, delta int) (*StockItemResponse, error)

	// Purchases
	RegisterPurchase(ctx context.Context, req PurchaseRequest) (*PurchaseResponse, error)
	GetPurchases(ctx context.Context, entityID string) ([]*Purchase, error)

	// Reports
	GetStockReport(ctx context.Context, entityID string) (*StockReport, error)
}

type SupplierResponse struct {
	SupplierID string
	Success    bool
	Error      string
}

type StockItemResponse struct {
	StockItemID string
	Success     bool
	Error       string
}

type StockReport struct {
	TotalItems    int
	TotalValue    int64
	LowStockItems []*StockItem
	ItemsByType   map[string]int
}
