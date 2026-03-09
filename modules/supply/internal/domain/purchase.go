package domain

import "time"

type Purchase struct {
	ID         string
	SupplierID string
	TotalValue int64 // Valor total em centavos (int64)
	Date       time.Time
	Items      []PurchaseItem
	CreatedAt  time.Time
}

type PurchaseItem struct {
	ID          string
	PurchaseID  string
	StockItemID string
	Quantity    int
	UnitCost    int64 // Custo unitário em centavos (int64)
	TotalCost   int64 // Quantity * UnitCost
}

func (p *Purchase) Validate() error {
	if p.SupplierID == "" {
		return ErrInvalidPurchaseSupplier
	}
	if p.TotalValue <= 0 {
		return ErrInvalidPurchaseTotalValue
	}
	if len(p.Items) == 0 {
		return ErrInvalidPurchaseItems
	}

	// Validar cada item
	for _, item := range p.Items {
		if err := item.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (p *Purchase) CalculateTotal() int64 {
	var total int64
	for _, item := range p.Items {
		total += item.TotalCost
	}
	return total
}

func (p *Purchase) AddItem(item PurchaseItem) {
	p.Items = append(p.Items, item)
	p.TotalValue = p.CalculateTotal()
}

func (pi *PurchaseItem) Validate() error {
	if pi.StockItemID == "" {
		return ErrInvalidPurchaseItemStockItem
	}
	if pi.Quantity <= 0 {
		return ErrInvalidPurchaseItemQuantity
	}
	if pi.UnitCost <= 0 {
		return ErrInvalidPurchaseItemUnitCost
	}
	if pi.TotalCost != pi.UnitCost*int64(pi.Quantity) {
		return ErrInvalidPurchaseItemTotalCost
	}
	return nil
}

func NewPurchaseItem(stockItemID string, quantity int, unitCost int64) PurchaseItem {
	return PurchaseItem{
		StockItemID: stockItemID,
		Quantity:    quantity,
		UnitCost:    unitCost,
		TotalCost:   unitCost * int64(quantity),
	}
}

// Erros de domínio para compras
var (
	ErrInvalidPurchaseSupplier      = newDomainError("fornecedor da compra inválido")
	ErrInvalidPurchaseTotalValue    = newDomainError("valor total da compra inválido")
	ErrInvalidPurchaseItems         = newDomainError("compra deve ter pelo menos um item")
	ErrInvalidPurchaseItemStockItem = newDomainError("item de estoque inválido na compra")
	ErrInvalidPurchaseItemQuantity  = newDomainError("quantidade do item de compra inválida")
	ErrInvalidPurchaseItemUnitCost  = newDomainError("custo unitário do item de compra inválido")
	ErrInvalidPurchaseItemTotalCost = newDomainError("custo total do item de compra inválido")
)
