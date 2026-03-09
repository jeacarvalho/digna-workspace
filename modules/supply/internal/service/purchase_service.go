package service

import (
	"context"
	"fmt"
	"time"

	"github.com/providentia/digna/supply/internal/domain"
	"github.com/providentia/digna/supply/internal/repository"
	"github.com/providentia/digna/supply/pkg/supply"
)

// Constantes para contas contábeis (deve ser consistente com outros módulos)
const (
	AccountCash      int64 = 1 // Caixa
	AccountInventory int64 = 3 // Estoque
	AccountSuppliers int64 = 4 // Fornecedores (Contas a Pagar)
	AccountExpenses  int64 = 5 // Despesas
)

type PurchaseService struct {
	repo       repository.SupplyRepository
	ledgerPort supply.LedgerPort
}

func NewPurchaseService(repo repository.SupplyRepository, ledgerPort supply.LedgerPort) *PurchaseService {
	return &PurchaseService{
		repo:       repo,
		ledgerPort: ledgerPort,
	}
}

// RegisterPurchase registra uma compra e gera a partida dobrada correspondente
func (s *PurchaseService) RegisterPurchase(ctx context.Context, entityID string, purchase *domain.Purchase, paymentType string) error {
	// Validar compra
	if err := purchase.Validate(); err != nil {
		return fmt.Errorf("validação da compra falhou: %w", err)
	}

	// Verificar se o fornecedor existe
	supplier, err := s.repo.GetSupplier(ctx, entityID, purchase.SupplierID)
	if err != nil {
		return fmt.Errorf("erro ao buscar fornecedor: %w", err)
	}
	if supplier == nil {
		return fmt.Errorf("fornecedor não encontrado: %s", purchase.SupplierID)
	}

	// Verificar se todos os itens existem e obter seus tipos
	itemTypes := make(map[string]domain.StockItemType)
	for _, item := range purchase.Items {
		stockItem, err := s.repo.GetStockItem(ctx, entityID, item.StockItemID)
		if err != nil {
			return fmt.Errorf("erro ao buscar item de estoque %s: %w", item.StockItemID, err)
		}
		if stockItem == nil {
			return fmt.Errorf("item de estoque não encontrado: %s", item.StockItemID)
		}
		itemTypes[item.StockItemID] = stockItem.Type
	}

	// Iniciar transação
	tx, err := s.repo.BeginTx(ctx, entityID)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}

	// Defer rollback em caso de erro
	success := false
	defer func() {
		if !success {
			s.repo.RollbackTx(tx)
		}
	}()

	// Salvar compra (isso também atualiza o estoque)
	if err := s.repo.SavePurchase(ctx, entityID, purchase); err != nil {
		return fmt.Errorf("erro ao salvar compra: %w", err)
	}

	// Registrar transação contábil
	if err := s.recordAccountingTransaction(ctx, entityID, purchase, itemTypes, paymentType); err != nil {
		return fmt.Errorf("erro ao registrar transação contábil: %w", err)
	}

	// Commit da transação
	if err := s.repo.CommitTx(tx); err != nil {
		return fmt.Errorf("erro ao commitar transação: %w", err)
	}

	success = true
	return nil
}

// recordAccountingTransaction registra a partida dobrada para a compra
func (s *PurchaseService) recordAccountingTransaction(ctx context.Context, entityID string, purchase *domain.Purchase, itemTypes map[string]domain.StockItemType, paymentType string) error {
	if s.ledgerPort == nil {
		// Se não houver ledger port, apenas logar (para testes)
		fmt.Printf("DEBUG: Compra registrada sem contabilidade (ledgerPort=nil): %s - R$ %d.%02d\n",
			purchase.ID, purchase.TotalValue/100, purchase.TotalValue%100)
		return nil
	}

	// Determinar conta de débito baseada no tipo dos itens
	debitAccount := determineDebitAccount(itemTypes)

	// Determinar conta de crédito baseada no tipo de pagamento
	creditAccount := determineCreditAccount(paymentType)

	// Criar postings para a transação
	postings := []supply.LedgerPosting{
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

	// Descrição da transação
	description := fmt.Sprintf("Compra %s - Fornecedor: %s", purchase.ID, purchase.SupplierID)

	// Registrar transação no ledger
	return s.ledgerPort.RecordTransaction(entityID, description, postings)
}

// determineDebitAccount determina a conta de débito baseada nos tipos dos itens
func determineDebitAccount(itemTypes map[string]domain.StockItemType) int64 {
	// Se todos os itens são do mesmo tipo, usar conta correspondente
	// Se houver mistura, usar AccountExpenses como fallback

	// Contar tipos
	typeCount := make(map[domain.StockItemType]int)
	for _, itemType := range itemTypes {
		typeCount[itemType]++
	}

	// Se houver apenas um tipo
	if len(typeCount) == 1 {
		for itemType := range typeCount {
			switch itemType {
			case domain.StockItemTypeRawMaterial, domain.StockItemTypeMerchandise:
				return AccountInventory // Estoque
			case domain.StockItemTypeProduct:
				// Produto acabado comprado (raro, mas possível)
				return AccountInventory // Estoque
			}
		}
	}

	// Fallback para despesas (mistura de tipos ou tipo desconhecido)
	return AccountExpenses
}

// determineCreditAccount determina a conta de crédito baseada no tipo de pagamento
func determineCreditAccount(paymentType string) int64 {
	switch paymentType {
	case "CASH", "À VISTA":
		return AccountCash
	case "CREDIT", "A PRAZO":
		return AccountSuppliers
	default:
		// Default para pagamento à vista
		return AccountCash
	}
}

// GetStockItemsByType retorna itens de estoque filtrados por tipo
func (s *PurchaseService) GetStockItemsByType(ctx context.Context, entityID string, itemType domain.StockItemType) ([]*domain.StockItem, error) {
	return s.repo.ListStockItemsByType(ctx, entityID, itemType)
}

// GetLowStockItems retorna itens com estoque abaixo do mínimo
func (s *PurchaseService) GetLowStockItems(ctx context.Context, entityID string) ([]*domain.StockItem, error) {
	items, err := s.repo.ListStockItems(ctx, entityID)
	if err != nil {
		return nil, err
	}

	var lowStock []*domain.StockItem
	for _, item := range items {
		if item.IsBelowMinimum() {
			lowStock = append(lowStock, item)
		}
	}

	return lowStock, nil
}

// CreateStockItem cria um novo item de estoque
func (s *PurchaseService) CreateStockItem(ctx context.Context, entityID string, item *domain.StockItem) error {
	if err := item.Validate(); err != nil {
		return err
	}

	return s.repo.SaveStockItem(ctx, entityID, item)
}

// CreateSupplier cria um novo fornecedor
func (s *PurchaseService) CreateSupplier(ctx context.Context, entityID string, supplier *domain.Supplier) error {
	if err := supplier.Validate(); err != nil {
		return err
	}

	return s.repo.SaveSupplier(ctx, entityID, supplier)
}

// GetPurchaseHistory retorna histórico de compras
func (s *PurchaseService) GetPurchaseHistory(ctx context.Context, entityID string, startDate, endDate time.Time) ([]*domain.Purchase, error) {
	allPurchases, err := s.repo.ListPurchases(ctx, entityID)
	if err != nil {
		return nil, err
	}

	var filtered []*domain.Purchase
	for _, purchase := range allPurchases {
		if (purchase.Date.After(startDate) || purchase.Date.Equal(startDate)) &&
			(purchase.Date.Before(endDate) || purchase.Date.Equal(endDate)) {
			filtered = append(filtered, purchase)
		}
	}

	return filtered, nil
}
