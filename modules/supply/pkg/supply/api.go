package supply

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
	"github.com/providentia/digna/supply/internal/domain"
	"github.com/providentia/digna/supply/internal/repository"
)

// SupplyAPIImpl implementa SupplyAPI sem causar ciclos de importação
type SupplyAPIImpl struct {
	lifecycleManager lifecycle.LifecycleManager
	repo             repository.SupplyRepository
	ledgerPort       LedgerPort
}

// NewSupplyAPI cria uma nova instância da API de supply
func NewSupplyAPI(lm lifecycle.LifecycleManager, ledgerPort LedgerPort) SupplyAPI {
	repo := repository.NewSQLiteSupplyRepository(lm)

	return &SupplyAPIImpl{
		lifecycleManager: lm,
		repo:             repo,
		ledgerPort:       ledgerPort,
	}
}

// RegisterSupplier implementa SupplyAPI.RegisterSupplier
func (api *SupplyAPIImpl) RegisterSupplier(ctx context.Context, req SupplierRequest) (*SupplierResponse, error) {
	supplier := &domain.Supplier{
		Name:        req.Name,
		ContactInfo: req.ContactInfo,
	}

	if err := supplier.Validate(); err != nil {
		return &SupplierResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	if err := api.repo.SaveSupplier(ctx, req.EntityID, supplier); err != nil {
		return &SupplierResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	return &SupplierResponse{
		SupplierID: supplier.ID,
		Success:    true,
	}, nil
}

// GetSuppliers implementa SupplyAPI.GetSuppliers
func (api *SupplyAPIImpl) GetSuppliers(ctx context.Context, entityID string) ([]*Supplier, error) {
	domainSuppliers, err := api.repo.ListSuppliers(ctx, entityID)
	if err != nil {
		return nil, err
	}

	suppliers := make([]*Supplier, len(domainSuppliers))
	for i, ds := range domainSuppliers {
		suppliers[i] = &Supplier{
			ID:          ds.ID,
			Name:        ds.Name,
			ContactInfo: ds.ContactInfo,
			CreatedAt:   ds.CreatedAt,
		}
	}

	return suppliers, nil
}

// RegisterStockItem implementa SupplyAPI.RegisterStockItem
func (api *SupplyAPIImpl) RegisterStockItem(ctx context.Context, req StockItemRequest) (*StockItemResponse, error) {
	item := &domain.StockItem{
		Name:        req.Name,
		Type:        domain.StockItemType(req.Type),
		Quantity:    req.Quantity,
		MinQuantity: req.MinQuantity,
		UnitCost:    req.UnitCost,
	}

	if err := item.Validate(); err != nil {
		return &StockItemResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	if err := api.repo.SaveStockItem(ctx, req.EntityID, item); err != nil {
		return &StockItemResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	return &StockItemResponse{
		StockItemID: item.ID,
		Success:     true,
	}, nil
}

// GetStockItems implementa SupplyAPI.GetStockItems
func (api *SupplyAPIImpl) GetStockItems(ctx context.Context, entityID string) ([]*StockItem, error) {
	domainItems, err := api.repo.ListStockItems(ctx, entityID)
	if err != nil {
		return nil, err
	}

	items := make([]*StockItem, len(domainItems))
	for i, di := range domainItems {
		items[i] = &StockItem{
			ID:          di.ID,
			Name:        di.Name,
			Type:        string(di.Type),
			Quantity:    di.Quantity,
			MinQuantity: di.MinQuantity,
			UnitCost:    di.UnitCost,
			CreatedAt:   di.CreatedAt,
		}
	}

	return items, nil
}

// GetStockItemsByType implementa SupplyAPI.GetStockItemsByType
func (api *SupplyAPIImpl) GetStockItemsByType(ctx context.Context, entityID string, itemType string) ([]*StockItem, error) {
	domainItems, err := api.repo.ListStockItemsByType(ctx, entityID, domain.StockItemType(itemType))
	if err != nil {
		return nil, err
	}

	items := make([]*StockItem, len(domainItems))
	for i, di := range domainItems {
		items[i] = &StockItem{
			ID:          di.ID,
			Name:        di.Name,
			Type:        string(di.Type),
			Quantity:    di.Quantity,
			MinQuantity: di.MinQuantity,
			UnitCost:    di.UnitCost,
			CreatedAt:   di.CreatedAt,
		}
	}

	return items, nil
}

// RegisterPurchase implementa SupplyAPI.RegisterPurchase
func (api *SupplyAPIImpl) RegisterPurchase(ctx context.Context, req PurchaseRequest) (*PurchaseResponse, error) {
	// Converter PurchaseRequest para domain.Purchase
	purchase := &domain.Purchase{
		SupplierID: req.SupplierID,
		Date:       time.Now(),
	}

	// Adicionar itens
	for _, itemReq := range req.Items {
		purchaseItem := domain.NewPurchaseItem(itemReq.StockItemID, itemReq.Quantity, itemReq.UnitCost)
		purchase.AddItem(purchaseItem)
	}

	// Validar compra
	if err := purchase.Validate(); err != nil {
		return &PurchaseResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	// Registrar compra (isso atualiza estoque)
	if err := api.repo.SavePurchase(ctx, req.EntityID, purchase); err != nil {
		return &PurchaseResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	// Registrar transação contábil se houver ledgerPort
	if api.ledgerPort != nil {
		// Determinar conta de débito (simplificado - sempre estoque para esta implementação)
		debitAccount := int64(3) // AccountInventory

		// Determinar conta de crédito
		creditAccount := int64(1) // AccountCash (default)
		if req.PaymentType == "CREDIT" || req.PaymentType == "A PRAZO" {
			creditAccount = int64(4) // AccountSuppliers
		}

		postings := []LedgerPosting{
			{
				AccountID: debitAccount,
				Amount:    purchase.TotalValue,
				Direction: "DEBIT",
			},
			{
				AccountID: creditAccount,
				Amount:    purchase.TotalValue,
				Direction: "CREDIT",
			},
		}

		description := fmt.Sprintf("Compra %s - Fornecedor: %s", purchase.ID, purchase.SupplierID)
		if err := api.ledgerPort.RecordTransaction(req.EntityID, description, postings); err != nil {
			// Logar erro mas não falhar a compra (para resiliência)
			fmt.Printf("AVISO: Falha ao registrar transação contábil: %v\n", err)
		}
	}

	return &PurchaseResponse{
		PurchaseID: purchase.ID,
		Success:    true,
	}, nil
}

// GetPurchases implementa SupplyAPI.GetPurchases
func (api *SupplyAPIImpl) GetPurchases(ctx context.Context, entityID string) ([]*Purchase, error) {
	domainPurchases, err := api.repo.ListPurchases(ctx, entityID)
	if err != nil {
		return nil, err
	}

	purchases := make([]*Purchase, len(domainPurchases))
	for i, dp := range domainPurchases {
		items := make([]PurchaseItem, len(dp.Items))
		for j, di := range dp.Items {
			items[j] = PurchaseItem{
				ID:          di.ID,
				StockItemID: di.StockItemID,
				Quantity:    di.Quantity,
				UnitCost:    di.UnitCost,
				TotalCost:   di.TotalCost,
			}
		}

		purchases[i] = &Purchase{
			ID:         dp.ID,
			SupplierID: dp.SupplierID,
			TotalValue: dp.TotalValue,
			Date:       dp.Date,
			Items:      items,
			CreatedAt:  dp.CreatedAt,
		}
	}

	return purchases, nil
}

// GetStockReport implementa SupplyAPI.GetStockReport
func (api *SupplyAPIImpl) GetStockReport(ctx context.Context, entityID string) (*StockReport, error) {
	items, err := api.GetStockItems(ctx, entityID)
	if err != nil {
		return nil, err
	}

	// Calcular valor total do estoque
	var totalValue int64
	var lowStockItems []*StockItem

	// Contar itens por tipo
	itemsByType := make(map[string]int)

	for _, item := range items {
		totalValue += item.UnitCost * int64(item.Quantity)
		itemsByType[item.Type]++

		// Verificar se está abaixo do mínimo
		if item.MinQuantity > 0 && item.Quantity < item.MinQuantity {
			lowStockItems = append(lowStockItems, item)
		}
	}

	return &StockReport{
		TotalItems:    len(items),
		TotalValue:    totalValue,
		LowStockItems: lowStockItems,
		ItemsByType:   itemsByType,
	}, nil
}

// UpdateStockQuantity implementa SupplyAPI.UpdateStockQuantity
func (api *SupplyAPIImpl) UpdateStockQuantity(ctx context.Context, entityID, itemID string, delta int) (*StockItemResponse, error) {
	slog.Info("SupplyAPI - UpdateStockQuantity chamado", "entity_id", entityID, "item_id", itemID, "delta", delta)

	// Buscar item atual
	item, err := api.repo.GetStockItem(ctx, entityID, itemID)
	if err != nil {
		slog.Error("SupplyAPI - Erro ao buscar item", "entity_id", entityID, "item_id", itemID, "erro", err)
		return &StockItemResponse{
			Success: false,
			Error:   fmt.Sprintf("erro ao buscar item: %v", err),
		}, err
	}

	if item == nil {
		slog.Error("SupplyAPI - Item não encontrado", "entity_id", entityID, "item_id", itemID)
		return &StockItemResponse{
			Success: false,
			Error:   "item não encontrado",
		}, nil
	}

	slog.Info("SupplyAPI - Item encontrado", "entity_id", entityID, "item_id", itemID, "nome", item.Name, "quantidade_atual", item.Quantity, "delta", delta)

	// Validar que não fica negativo
	newQuantity := item.Quantity + delta
	if newQuantity < 0 {
		slog.Error("SupplyAPI - Quantidade ficaria negativa", "entity_id", entityID, "item_id", itemID, "quantidade_atual", item.Quantity, "delta", delta, "nova_quantidade", newQuantity)
		return &StockItemResponse{
			Success: false,
			Error:   fmt.Sprintf("quantidade não pode ficar negativa: %d + %d = %d", item.Quantity, delta, newQuantity),
		}, nil
	}

	// Atualizar quantidade
	slog.Info("SupplyAPI - Chamando repositório para atualizar quantidade", "entity_id", entityID, "item_id", itemID, "delta", delta)
	err = api.repo.UpdateStockQuantity(ctx, entityID, itemID, delta)
	if err != nil {
		slog.Error("SupplyAPI - Erro ao atualizar estoque no repositório", "entity_id", entityID, "item_id", itemID, "delta", delta, "erro", err)
		return &StockItemResponse{
			Success: false,
			Error:   fmt.Sprintf("erro ao atualizar estoque: %v", err),
		}, err
	}

	slog.Info("SupplyAPI - Estoque atualizado com sucesso", "entity_id", entityID, "item_id", itemID, "quantidade_anterior", item.Quantity, "nova_quantidade", newQuantity)
	return &StockItemResponse{
		StockItemID: itemID,
		Success:     true,
	}, nil
}
