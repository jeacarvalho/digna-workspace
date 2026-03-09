package repository

import (
	"context"

	"github.com/providentia/digna/supply/internal/domain"
)

type SupplyRepository interface {
	// Suppliers
	SaveSupplier(ctx context.Context, entityID string, supplier *domain.Supplier) error
	GetSupplier(ctx context.Context, entityID, supplierID string) (*domain.Supplier, error)
	ListSuppliers(ctx context.Context, entityID string) ([]*domain.Supplier, error)

	// Stock Items
	SaveStockItem(ctx context.Context, entityID string, item *domain.StockItem) error
	GetStockItem(ctx context.Context, entityID, itemID string) (*domain.StockItem, error)
	ListStockItems(ctx context.Context, entityID string) ([]*domain.StockItem, error)
	ListStockItemsByType(ctx context.Context, entityID string, itemType domain.StockItemType) ([]*domain.StockItem, error)
	UpdateStockQuantity(ctx context.Context, entityID, itemID string, delta int) error

	// Purchases
	SavePurchase(ctx context.Context, entityID string, purchase *domain.Purchase) error
	GetPurchase(ctx context.Context, entityID, purchaseID string) (*domain.Purchase, error)
	ListPurchases(ctx context.Context, entityID string) ([]*domain.Purchase, error)
	ListPurchasesBySupplier(ctx context.Context, entityID, supplierID string) ([]*domain.Purchase, error)

	// Transaction management
	BeginTx(ctx context.Context, entityID string) (interface{}, error)
	CommitTx(tx interface{}) error
	RollbackTx(tx interface{}) error
}
