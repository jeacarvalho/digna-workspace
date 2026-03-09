package domain

import "time"

type StockItemType string

const (
	StockItemTypeRawMaterial StockItemType = "INSUMO"     // Matéria-prima para produção
	StockItemTypeProduct     StockItemType = "PRODUTO"    // Produto acabado para venda
	StockItemTypeMerchandise StockItemType = "MERCADORIA" // Produto para revenda
)

type StockItemUnit string

const (
	StockItemUnitUnit       StockItemUnit = "UNIDADE" // Unidades (peças, caixas, etc)
	StockItemUnitKilogram   StockItemUnit = "KG"      // Quilogramas
	StockItemUnitGram       StockItemUnit = "G"       // Gramas
	StockItemUnitLiter      StockItemUnit = "L"       // Litros
	StockItemUnitMeter      StockItemUnit = "M"       // Metros
	StockItemUnitCentimeter StockItemUnit = "CM"      // Centímetros
	StockItemUnitPackage    StockItemUnit = "PACOTE"  // Pacotes
	StockItemUnitBox        StockItemUnit = "CAIXA"   // Caixas
	StockItemUnitBag        StockItemUnit = "SACO"    // Sacos
)

type StockItem struct {
	ID          string
	Name        string
	Type        StockItemType // INSUMO|PRODUTO|MERCADORIA
	Unit        StockItemUnit // UNIDADE|KG|G|L|M|CM|PACOTE|CAIXA|SACO
	Quantity    int           // Quantidade atual em estoque
	MinQuantity int           // Quantidade mínima para alerta
	UnitCost    int64         // Custo unitário em centavos (int64)
	CreatedAt   time.Time
}

func (si *StockItem) Validate() error {
	if si.Name == "" {
		return ErrInvalidStockItemName
	}
	if si.Quantity < 0 {
		return ErrInvalidStockItemQuantity
	}
	if si.UnitCost < 0 {
		return ErrInvalidStockItemUnitCost
	}
	if si.Type != StockItemTypeRawMaterial && si.Type != StockItemTypeProduct && si.Type != StockItemTypeMerchandise {
		return ErrInvalidStockItemType
	}
	// Validar unidade
	if si.Unit == "" {
		si.Unit = StockItemUnitUnit // Default para UNIDADE
	}
	if !si.isValidUnit(si.Unit) {
		return ErrInvalidStockItemUnit
	}
	return nil
}

func (si *StockItem) isValidUnit(unit StockItemUnit) bool {
	validUnits := []StockItemUnit{
		StockItemUnitUnit,
		StockItemUnitKilogram,
		StockItemUnitGram,
		StockItemUnitLiter,
		StockItemUnitMeter,
		StockItemUnitCentimeter,
		StockItemUnitPackage,
		StockItemUnitBox,
		StockItemUnitBag,
	}
	for _, validUnit := range validUnits {
		if unit == validUnit {
			return true
		}
	}
	return false
}

func (si *StockItem) UpdateQuantity(delta int) error {
	newQuantity := si.Quantity + delta
	if newQuantity < 0 {
		return ErrInsufficientStock
	}
	si.Quantity = newQuantity
	return nil
}

func (si *StockItem) CalculateTotalCost(quantity int) int64 {
	return si.UnitCost * int64(quantity)
}

func (si *StockItem) IsBelowMinimum() bool {
	return si.MinQuantity > 0 && si.Quantity < si.MinQuantity
}

// Erros de domínio
var (
	ErrInvalidStockItemName     = newDomainError("nome do item de estoque inválido")
	ErrInvalidStockItemQuantity = newDomainError("quantidade do item de estoque inválida")
	ErrInvalidStockItemUnitCost = newDomainError("custo unitário do item de estoque inválido")
	ErrInvalidStockItemType     = newDomainError("tipo do item de estoque inválido")
	ErrInvalidStockItemUnit     = newDomainError("unidade do item de estoque inválida")
	ErrInsufficientStock        = newDomainError("estoque insuficiente")
)

type DomainError struct {
	message string
}

func newDomainError(message string) *DomainError {
	return &DomainError{message: message}
}

func (e *DomainError) Error() string {
	return e.message
}
