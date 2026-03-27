package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/providentia/digna/core_lume/internal/domain"
	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

// HelpRepository interface for help topic operations
type HelpRepository interface {
	Save(topic *domain.HelpTopic) error
	FindByKey(key string) (*domain.HelpTopic, error)
	FindByID(id string) (*domain.HelpTopic, error)
	ListByCategory(category string) ([]*domain.HelpTopic, error)
	Search(query string) ([]*domain.HelpTopic, error)
	ListAll() ([]*domain.HelpTopic, error)
	IncrementViewCount(id string) error
	InitTable() error
	SeedTopics() error
}

// SQLiteHelpRepository implements HelpRepository for SQLite (central.db)
type SQLiteHelpRepository struct {
	centralDB *sql.DB
}

// NewSQLiteHelpRepository creates a new SQLiteHelpRepository
func NewSQLiteHelpRepository(lm lifecycle.LifecycleManager) (*SQLiteHelpRepository, error) {
	// Get central database connection
	centralDB, err := lm.GetCentralConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to get central connection: %w", err)
	}

	return &SQLiteHelpRepository{
		centralDB: centralDB,
	}, nil
}

// InitTable creates the help_topics table if not exists
func (r *SQLiteHelpRepository) InitTable() error {
	_, err := r.centralDB.Exec(`
		CREATE TABLE IF NOT EXISTS help_topics (
			id TEXT PRIMARY KEY,
			key TEXT NOT NULL UNIQUE,
			title TEXT NOT NULL,
			summary TEXT,
			explanation TEXT NOT NULL,
			why_asked TEXT,
			legislation TEXT,
			next_steps TEXT,
			official_link TEXT,
			category TEXT NOT NULL,
			tags TEXT,
			view_count INTEGER DEFAULT 0,
			created_at INTEGER,
			updated_at INTEGER
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create help_topics table: %w", err)
	}

	// Create index for faster lookups
	_, err = r.centralDB.Exec(`
		CREATE INDEX IF NOT EXISTS idx_help_key ON help_topics(key)
	`)
	if err != nil {
		return fmt.Errorf("failed to create key index: %w", err)
	}

	_, err = r.centralDB.Exec(`
		CREATE INDEX IF NOT EXISTS idx_help_category ON help_topics(category)
	`)
	if err != nil {
		return fmt.Errorf("failed to create category index: %w", err)
	}

	return nil
}

// Save creates or updates a help topic
func (r *SQLiteHelpRepository) Save(topic *domain.HelpTopic) error {
	if err := topic.Validate(); err != nil {
		return fmt.Errorf("invalid help topic: %w", err)
	}

	now := time.Now().Unix()
	if topic.CreatedAt == 0 {
		topic.CreatedAt = now
	}
	topic.UpdatedAt = now

	_, err := r.centralDB.Exec(`
		INSERT INTO help_topics (
			id, key, title, summary, explanation, why_asked, legislation,
			next_steps, official_link, category, tags, view_count, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(key) DO UPDATE SET
			title = excluded.title,
			summary = excluded.summary,
			explanation = excluded.explanation,
			why_asked = excluded.why_asked,
			legislation = excluded.legislation,
			next_steps = excluded.next_steps,
			official_link = excluded.official_link,
			category = excluded.category,
			tags = excluded.tags,
			view_count = excluded.view_count,
			updated_at = excluded.updated_at
	`,
		topic.ID, topic.Key, topic.Title, topic.Summary, topic.Explanation,
		topic.WhyAsked, topic.Legislation, topic.NextSteps, topic.OfficialLink,
		topic.Category, topic.Tags, topic.ViewCount, topic.CreatedAt, topic.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save help topic: %w", err)
	}

	return nil
}

// FindByKey finds a help topic by key
func (r *SQLiteHelpRepository) FindByKey(key string) (*domain.HelpTopic, error) {
	var topic domain.HelpTopic

	err := r.centralDB.QueryRow(`
		SELECT id, key, title, summary, explanation, why_asked, legislation,
			next_steps, official_link, category, tags, view_count, created_at, updated_at
		FROM help_topics WHERE key = ?
	`, key).Scan(
		&topic.ID, &topic.Key, &topic.Title, &topic.Summary, &topic.Explanation,
		&topic.WhyAsked, &topic.Legislation, &topic.NextSteps, &topic.OfficialLink,
		&topic.Category, &topic.Tags, &topic.ViewCount, &topic.CreatedAt, &topic.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrHelpTopicNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query help topic: %w", err)
	}

	return &topic, nil
}

// FindByID finds a help topic by ID
func (r *SQLiteHelpRepository) FindByID(id string) (*domain.HelpTopic, error) {
	var topic domain.HelpTopic

	err := r.centralDB.QueryRow(`
		SELECT id, key, title, summary, explanation, why_asked, legislation,
			next_steps, official_link, category, tags, view_count, created_at, updated_at
		FROM help_topics WHERE id = ?
	`, id).Scan(
		&topic.ID, &topic.Key, &topic.Title, &topic.Summary, &topic.Explanation,
		&topic.WhyAsked, &topic.Legislation, &topic.NextSteps, &topic.OfficialLink,
		&topic.Category, &topic.Tags, &topic.ViewCount, &topic.CreatedAt, &topic.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrHelpTopicNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query help topic: %w", err)
	}

	return &topic, nil
}

// ListByCategory lists all help topics in a category
func (r *SQLiteHelpRepository) ListByCategory(category string) ([]*domain.HelpTopic, error) {
	rows, err := r.centralDB.Query(`
		SELECT id, key, title, summary, explanation, why_asked, legislation,
			next_steps, official_link, category, tags, view_count, created_at, updated_at
		FROM help_topics WHERE category = ?
		ORDER BY title ASC
	`, category)
	if err != nil {
		return nil, fmt.Errorf("failed to query help topics: %w", err)
	}
	defer rows.Close()

	var topics []*domain.HelpTopic
	for rows.Next() {
		var topic domain.HelpTopic
		err := rows.Scan(
			&topic.ID, &topic.Key, &topic.Title, &topic.Summary, &topic.Explanation,
			&topic.WhyAsked, &topic.Legislation, &topic.NextSteps, &topic.OfficialLink,
			&topic.Category, &topic.Tags, &topic.ViewCount, &topic.CreatedAt, &topic.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan help topic: %w", err)
		}
		topics = append(topics, &topic)
	}

	return topics, nil
}

// Search searches help topics by query
func (r *SQLiteHelpRepository) Search(query string) ([]*domain.HelpTopic, error) {
	// Simple search in title, summary, explanation, and tags
	searchPattern := "%" + query + "%"

	rows, err := r.centralDB.Query(`
		SELECT id, key, title, summary, explanation, why_asked, legislation,
			next_steps, official_link, category, tags, view_count, created_at, updated_at
		FROM help_topics
		WHERE title LIKE ? OR summary LIKE ? OR explanation LIKE ? OR tags LIKE ?
		ORDER BY title ASC
	`, searchPattern, searchPattern, searchPattern, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search help topics: %w", err)
	}
	defer rows.Close()

	var topics []*domain.HelpTopic
	for rows.Next() {
		var topic domain.HelpTopic
		err := rows.Scan(
			&topic.ID, &topic.Key, &topic.Title, &topic.Summary, &topic.Explanation,
			&topic.WhyAsked, &topic.Legislation, &topic.NextSteps, &topic.OfficialLink,
			&topic.Category, &topic.Tags, &topic.ViewCount, &topic.CreatedAt, &topic.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan help topic: %w", err)
		}
		topics = append(topics, &topic)
	}

	return topics, nil
}

// ListAll lists all help topics
func (r *SQLiteHelpRepository) ListAll() ([]*domain.HelpTopic, error) {
	rows, err := r.centralDB.Query(`
		SELECT id, key, title, summary, explanation, why_asked, legislation,
			next_steps, official_link, category, tags, view_count, created_at, updated_at
		FROM help_topics
		ORDER BY category ASC, title ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query all help topics: %w", err)
	}
	defer rows.Close()

	var topics []*domain.HelpTopic
	for rows.Next() {
		var topic domain.HelpTopic
		err := rows.Scan(
			&topic.ID, &topic.Key, &topic.Title, &topic.Summary, &topic.Explanation,
			&topic.WhyAsked, &topic.Legislation, &topic.NextSteps, &topic.OfficialLink,
			&topic.Category, &topic.Tags, &topic.ViewCount, &topic.CreatedAt, &topic.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan help topic: %w", err)
		}
		topics = append(topics, &topic)
	}

	return topics, nil
}

// IncrementViewCount increments the view count for a topic
func (r *SQLiteHelpRepository) IncrementViewCount(id string) error {
	_, err := r.centralDB.Exec(`
		UPDATE help_topics SET view_count = view_count + 1, updated_at = ?
		WHERE id = ?
	`, time.Now().Unix(), id)
	if err != nil {
		return fmt.Errorf("failed to increment view count: %w", err)
	}
	return nil
}

// SeedTopics seeds initial help topics
func (r *SQLiteHelpRepository) SeedTopics() error {
	now := time.Now().Unix()

	for _, topic := range domain.InitialHelpTopics {
		topic.ID = fmt.Sprintf("help-%s", topic.Key)
		topic.CreatedAt = now
		topic.UpdatedAt = now
		topic.ViewCount = 0

		// Try to save, ignore duplicates
		if err := r.Save(&topic); err != nil {
			// Continue if duplicate, fail on other errors
			if !isDuplicateError(err) {
				return fmt.Errorf("failed to seed topic %s: %w", topic.Key, err)
			}
		}
	}

	return nil
}

// isDuplicateError checks if error is a duplicate key error
func isDuplicateError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return contains(errStr, "UNIQUE") || contains(errStr, "unique")
}

// contains checks if string contains substring
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
