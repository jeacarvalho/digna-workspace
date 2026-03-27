package service

import (
	"fmt"
	"sort"

	"github.com/providentia/digna/core_lume/internal/domain"
	"github.com/providentia/digna/core_lume/internal/repository"
)

// HelpService implements business logic for help topics
type HelpService struct {
	repo repository.HelpRepository
}

// NewHelpService creates a new HelpService
func NewHelpService(repo repository.HelpRepository) *HelpService {
	return &HelpService{
		repo: repo,
	}
}

// GetTopicByKey retrieves a help topic by its key
func (s *HelpService) GetTopicByKey(key string) (*domain.HelpTopic, error) {
	topic, err := s.repo.FindByKey(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get topic: %w", err)
	}

	// Increment view count
	if err := s.repo.IncrementViewCount(topic.ID); err != nil {
		// Log but don't fail - view count is not critical
		fmt.Printf("Warning: failed to increment view count for topic %s: %v\n", key, err)
	}

	return topic, nil
}

// GetTopicByID retrieves a help topic by its ID
func (s *HelpService) GetTopicByID(id string) (*domain.HelpTopic, error) {
	topic, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get topic: %w", err)
	}
	return topic, nil
}

// ListIndex returns all help topics grouped by category
func (s *HelpService) ListIndex() (map[string][]*domain.HelpTopic, error) {
	topics, err := s.repo.ListAll()
	if err != nil {
		return nil, fmt.Errorf("failed to list topics: %w", err)
	}

	// Group by category
	grouped := make(map[string][]*domain.HelpTopic)
	for _, topic := range topics {
		grouped[topic.Category] = append(grouped[topic.Category], topic)
	}

	// Sort topics within each category by title
	for _, topics := range grouped {
		sort.Slice(topics, func(i, j int) bool {
			return topics[i].Title < topics[j].Title
		})
	}

	return grouped, nil
}

// Search searches for help topics
func (s *HelpService) Search(query string) ([]*domain.HelpTopic, error) {
	if query == "" {
		return s.repo.ListAll()
	}

	topics, err := s.repo.Search(query)
	if err != nil {
		return nil, fmt.Errorf("failed to search topics: %w", err)
	}

	return topics, nil
}

// GetRelatedTopics returns related topics for a given topic
func (s *HelpService) GetRelatedTopics(topic *domain.HelpTopic) ([]*domain.HelpTopic, error) {
	// Get all topics in the same category
	allTopics, err := s.repo.ListByCategory(topic.Category)
	if err != nil {
		return nil, fmt.Errorf("failed to get related topics: %w", err)
	}

	// Filter out the current topic and limit to 3
	var related []*domain.HelpTopic
	for _, t := range allTopics {
		if t.Key != topic.Key {
			related = append(related, t)
			if len(related) >= 3 {
				break
			}
		}
	}

	return related, nil
}

// IncrementView increments the view count for a topic
func (s *HelpService) IncrementView(key string) error {
	topic, err := s.repo.FindByKey(key)
	if err != nil {
		return fmt.Errorf("failed to find topic: %w", err)
	}

	if err := s.repo.IncrementViewCount(topic.ID); err != nil {
		return fmt.Errorf("failed to increment view: %w", err)
	}

	return nil
}

// GetCategories returns all available categories with labels
func (s *HelpService) GetCategories() map[string]string {
	return map[string]string{
		domain.CategoriaCredito:    domain.GetCategoryLabel(domain.CategoriaCredito),
		domain.CategoriaTributario: domain.GetCategoryLabel(domain.CategoriaTributario),
		domain.CategoriaGovernanca: domain.GetCategoryLabel(domain.CategoriaGovernanca),
		domain.CategoriaGeral:      domain.GetCategoryLabel(domain.CategoriaGeral),
	}
}

// Initialize seeds the database with initial topics
func (s *HelpService) Initialize() error {
	if err := s.repo.InitTable(); err != nil {
		return fmt.Errorf("failed to init table: %w", err)
	}

	if err := s.repo.SeedTopics(); err != nil {
		return fmt.Errorf("failed to seed topics: %w", err)
	}

	return nil
}
