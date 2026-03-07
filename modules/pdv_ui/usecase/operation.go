package usecase

import (
	"fmt"
	"time"

	"github.com/providentia/digna/core_lume/pkg/governance"
	"github.com/providentia/digna/core_lume/pkg/ledger"
	"github.com/providentia/digna/core_lume/pkg/social"
	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

const (
	AccountCash      int64 = 1
	AccountSales     int64 = 2
	AccountInventory int64 = 3
	AccountCOGS      int64 = 4
)

type SaleRequest struct {
	EntityID      string
	Amount        int64
	PaymentMethod string
	Description   string
}

type SaleResult struct {
	EntryID int64
	Error   error
}

type OperationHandler struct {
	ledgerService *ledger.Service
}

func NewOperationHandler(lm lifecycle.LifecycleManager) *OperationHandler {
	return &OperationHandler{
		ledgerService: ledger.NewService(lm),
	}
}

func (oh *OperationHandler) RecordSale(req SaleRequest) (*SaleResult, error) {
	if req.Amount <= 0 {
		return &SaleResult{Error: fmt.Errorf("amount must be positive")}, fmt.Errorf("amount must be positive")
	}
	if req.EntityID == "" {
		return &SaleResult{Error: fmt.Errorf("entity_id cannot be empty")}, fmt.Errorf("entity_id cannot be empty")
	}

	txn := &ledger.Transaction{
		Date:        time.Now(),
		Description: req.Description,
		Reference:   fmt.Sprintf("SALE-%s-%d", req.PaymentMethod, time.Now().Unix()),
		Postings: []ledger.Posting{
			{
				AccountID: AccountCash,
				Amount:    req.Amount,
				Direction: ledger.Debit,
			},
			{
				AccountID: AccountSales,
				Amount:    req.Amount,
				Direction: ledger.Credit,
			},
		},
	}

	if req.Description == "" {
		txn.Description = fmt.Sprintf("Venda - %s", req.PaymentMethod)
	}

	err := oh.ledgerService.RecordTransaction(req.EntityID, txn)
	if err != nil {
		return &SaleResult{Error: err}, err
	}

	return &SaleResult{EntryID: txn.ID}, nil
}

func (oh *OperationHandler) GetCashBalance(entityID string) (int64, error) {
	return oh.ledgerService.GetAccountBalance(entityID, AccountCash)
}

type WorkRequest struct {
	EntityID     string
	MemberID     string
	Minutes      int64
	ActivityType string
	Description  string
}

type DecisionRequest struct {
	EntityID string
	Title    string
	Content  string
}

type SocialGovernanceHandler struct {
	socialService     *social.Service
	governanceService *governance.Service
}

func NewSocialGovernanceHandler(lm lifecycle.LifecycleManager) *SocialGovernanceHandler {
	return &SocialGovernanceHandler{
		socialService:     social.NewService(lm),
		governanceService: governance.NewService(lm),
	}
}

func (sgh *SocialGovernanceHandler) RecordWork(req WorkRequest) error {
	if req.Minutes <= 0 {
		return fmt.Errorf("minutes must be positive")
	}
	if req.MemberID == "" {
		return fmt.Errorf("member_id cannot be empty")
	}
	if req.EntityID == "" {
		return fmt.Errorf("entity_id cannot be empty")
	}

	record := &social.WorkRecord{
		MemberID:     req.MemberID,
		Minutes:      req.Minutes,
		ActivityType: req.ActivityType,
		LogDate:      time.Now(),
		Description:  req.Description,
	}

	return sgh.socialService.RecordWork(req.EntityID, record)
}

func (sgh *SocialGovernanceHandler) GetMemberWorkCapital(entityID string, memberID string) (int64, int64, error) {
	return sgh.socialService.GetTotalWorkByMember(entityID, memberID)
}

func (sgh *SocialGovernanceHandler) RecordDecision(req DecisionRequest) (string, error) {
	if req.Title == "" {
		return "", fmt.Errorf("title cannot be empty")
	}
	if req.Content == "" {
		return "", fmt.Errorf("content cannot be empty")
	}
	if req.EntityID == "" {
		return "", fmt.Errorf("entity_id cannot be empty")
	}

	return sgh.governanceService.RecordDecision(req.EntityID, req.Title, req.Content)
}

func (sgh *SocialGovernanceHandler) GetDecisionByHash(entityID string, hash string) (*governance.DecisionRecord, error) {
	return sgh.governanceService.GetDecisionByHash(entityID, hash)
}
