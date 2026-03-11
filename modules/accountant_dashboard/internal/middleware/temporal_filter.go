package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

// TemporalFilterMiddleware provides temporal filtering for accountant access
type TemporalFilterMiddleware struct {
	accountantLinkService lifecycle.AccountantLinkService
}

// NewTemporalFilterMiddleware creates a new temporal filter middleware
func NewTemporalFilterMiddleware(accountantLinkService lifecycle.AccountantLinkService) *TemporalFilterMiddleware {
	return &TemporalFilterMiddleware{
		accountantLinkService: accountantLinkService,
	}
}

// ContextKey is a type for context keys
type ContextKey string

const (
	// AccountantIDKey is the context key for accountant ID
	AccountantIDKey ContextKey = "accountant_id"
	// ValidEnterprisesKey is the context key for valid enterprises
	ValidEnterprisesKey ContextKey = "valid_enterprises"
	// CurrentPeriodKey is the context key for current period
	CurrentPeriodKey ContextKey = "current_period"
)

// Handler returns an HTTP handler that injects temporal filtering context
func (m *TemporalFilterMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only apply to accountant routes
		if !strings.HasPrefix(r.URL.Path, "/accountant") {
			next.ServeHTTP(w, r)
			return
		}

		// Get accountant ID from session (simplified - in real implementation, this would come from auth)
		accountantID := r.Header.Get("X-Accountant-ID")
		if accountantID == "" {
			// Try to get from query parameter (for testing)
			accountantID = r.URL.Query().Get("accountant_id")
		}

		if accountantID == "" {
			http.Error(w, "Accountant ID required", http.StatusUnauthorized)
			return
		}

		// Get period from query or use current month
		period := r.URL.Query().Get("period")
		if period == "" {
			period = time.Now().Format("2006-01")
		}

		// Parse period to get start and end timestamps
		startTime, endTime, err := parsePeriodToTime(period)
		if err != nil {
			http.Error(w, "Invalid period format", http.StatusBadRequest)
			return
		}

		// Get valid enterprises for this accountant during the period
		ctx := r.Context()
		validEnterprises, err := m.accountantLinkService.GetValidEnterprisesForAccountant(
			ctx, accountantID, startTime, endTime)
		if err != nil {
			http.Error(w, "Failed to validate accountant access", http.StatusInternalServerError)
			return
		}

		// If no valid enterprises, return empty list
		if len(validEnterprises) == 0 {
			http.Error(w, "No valid enterprises found for this period", http.StatusForbidden)
			return
		}

		// Create a map for quick lookup
		enterpriseMap := make(map[string]bool)
		for _, enterprise := range validEnterprises {
			enterpriseMap[enterprise] = true
		}

		// Add context values
		ctx = context.WithValue(ctx, AccountantIDKey, accountantID)
		ctx = context.WithValue(ctx, ValidEnterprisesKey, enterpriseMap)
		ctx = context.WithValue(ctx, CurrentPeriodKey, period)

		// Update request with new context
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// parsePeriodToTime converts a period string (YYYY-MM) to time.Time values
func parsePeriodToTime(period string) (time.Time, time.Time, error) {
	t, err := time.Parse("2006-01", period)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	// Start of month
	startTime := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)

	// End of month
	endTime := time.Date(t.Year(), t.Month()+1, 0, 23, 59, 59, 999999999, time.UTC)

	return startTime, endTime, nil
}

// FilterEntities filters a list of entities based on valid enterprises
func FilterEntities(entities []string, validEnterprises map[string]bool) []string {
	var filtered []string
	for _, entity := range entities {
		if validEnterprises[entity] {
			filtered = append(filtered, entity)
		}
	}
	return filtered
}

// GetAccountantIDFromContext retrieves accountant ID from context
func GetAccountantIDFromContext(ctx context.Context) (string, bool) {
	value := ctx.Value(AccountantIDKey)
	if value == nil {
		return "", false
	}
	accountantID, ok := value.(string)
	return accountantID, ok
}

// GetValidEnterprisesFromContext retrieves valid enterprises from context
func GetValidEnterprisesFromContext(ctx context.Context) (map[string]bool, bool) {
	value := ctx.Value(ValidEnterprisesKey)
	if value == nil {
		return nil, false
	}
	enterprises, ok := value.(map[string]bool)
	return enterprises, ok
}

// GetCurrentPeriodFromContext retrieves current period from context
func GetCurrentPeriodFromContext(ctx context.Context) (string, bool) {
	value := ctx.Value(CurrentPeriodKey)
	if value == nil {
		return "", false
	}
	period, ok := value.(string)
	return period, ok
}
