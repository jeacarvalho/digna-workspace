package help

import (
	"github.com/providentia/digna/core_lume/internal/domain"
	"github.com/providentia/digna/core_lume/internal/repository"
	"github.com/providentia/digna/core_lume/internal/service"
	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

// HelpTopic represents the public help topic model
type HelpTopic struct {
	ID           string
	Key          string
	Title        string
	Summary      string
	Explanation  string
	WhyAsked     string
	Legislation  string
	NextSteps    string
	OfficialLink string
	Category     string
	Tags         string
	ViewCount    int64
	CreatedAt    int64
	UpdatedAt    int64
}

// Service provides help system operations
type Service struct {
	helpService *service.HelpService
}

// NewService creates a new help service
func NewService(lm lifecycle.LifecycleManager) (*Service, error) {
	helpRepo, err := repository.NewSQLiteHelpRepository(lm)
	if err != nil {
		return nil, err
	}

	helpService := service.NewHelpService(helpRepo)

	// Initialize database with seed topics
	if err := helpService.Initialize(); err != nil {
		return nil, err
	}

	return &Service{
		helpService: helpService,
	}, nil
}

// GetTopicByKey retrieves a help topic by its key
func (s *Service) GetTopicByKey(key string) (*HelpTopic, error) {
	topic, err := s.helpService.GetTopicByKey(key)
	if err != nil {
		return nil, err
	}
	return convertToPublic(topic), nil
}

// GetTopicByID retrieves a help topic by its ID
func (s *Service) GetTopicByID(id string) (*HelpTopic, error) {
	topic, err := s.helpService.GetTopicByID(id)
	if err != nil {
		return nil, err
	}
	return convertToPublic(topic), nil
}

// ListIndex returns all help topics grouped by category
func (s *Service) ListIndex() (map[string][]*HelpTopic, error) {
	index, err := s.helpService.ListIndex()
	if err != nil {
		return nil, err
	}

	publicIndex := make(map[string][]*HelpTopic)
	for category, topics := range index {
		for _, topic := range topics {
			publicIndex[category] = append(publicIndex[category], convertToPublic(topic))
		}
	}

	return publicIndex, nil
}

// Search searches for help topics
func (s *Service) Search(query string) ([]*HelpTopic, error) {
	topics, err := s.helpService.Search(query)
	if err != nil {
		return nil, err
	}

	var result []*HelpTopic
	for _, topic := range topics {
		result = append(result, convertToPublic(topic))
	}
	return result, nil
}

// GetRelatedTopics returns related topics for a given topic
func (s *Service) GetRelatedTopics(topic *HelpTopic) ([]*HelpTopic, error) {
	domainTopic := &domain.HelpTopic{Key: topic.Key, Category: topic.Category}
	related, err := s.helpService.GetRelatedTopics(domainTopic)
	if err != nil {
		return nil, err
	}

	var result []*HelpTopic
	for _, t := range related {
		result = append(result, convertToPublic(t))
	}
	return result, nil
}

// GetCategories returns all available categories with labels
func (s *Service) GetCategories() map[string]string {
	return s.helpService.GetCategories()
}

// Helper function to convert internal to public
func convertToPublic(topic *domain.HelpTopic) *HelpTopic {
	return &HelpTopic{
		ID:           topic.ID,
		Key:          topic.Key,
		Title:        topic.Title,
		Summary:      topic.Summary,
		Explanation:  topic.Explanation,
		WhyAsked:     topic.WhyAsked,
		Legislation:  topic.Legislation,
		NextSteps:    topic.NextSteps,
		OfficialLink: topic.OfficialLink,
		Category:     topic.Category,
		Tags:         topic.Tags,
		ViewCount:    topic.ViewCount,
		CreatedAt:    topic.CreatedAt,
		UpdatedAt:    topic.UpdatedAt,
	}
}
