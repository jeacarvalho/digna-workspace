package service

import (
	"testing"
	"time"

	"github.com/providentia/digna/core_lume/internal/domain"
)

// MockHelpRepository é um mock do HelpRepository para testes
type MockHelpRepository struct {
	topics map[string]*domain.HelpTopic
}

func NewMockHelpRepository() *MockHelpRepository {
	return &MockHelpRepository{
		topics: make(map[string]*domain.HelpTopic),
	}
}

func (m *MockHelpRepository) Save(topic *domain.HelpTopic) error {
	if err := topic.Validate(); err != nil {
		return err
	}
	m.topics[topic.Key] = topic
	return nil
}

func (m *MockHelpRepository) FindByKey(key string) (*domain.HelpTopic, error) {
	topic, exists := m.topics[key]
	if !exists {
		return nil, domain.ErrHelpTopicNotFound
	}
	return topic, nil
}

func (m *MockHelpRepository) FindByID(id string) (*domain.HelpTopic, error) {
	for _, topic := range m.topics {
		if topic.ID == id {
			return topic, nil
		}
	}
	return nil, domain.ErrHelpTopicNotFound
}

func (m *MockHelpRepository) ListByCategory(category string) ([]*domain.HelpTopic, error) {
	var result []*domain.HelpTopic
	for _, topic := range m.topics {
		if topic.Category == category {
			result = append(result, topic)
		}
	}
	return result, nil
}

func (m *MockHelpRepository) Search(query string) ([]*domain.HelpTopic, error) {
	var result []*domain.HelpTopic
	for _, topic := range m.topics {
		if topic.MatchesSearch(query) {
			result = append(result, topic)
		}
	}
	return result, nil
}

func (m *MockHelpRepository) ListAll() ([]*domain.HelpTopic, error) {
	var result []*domain.HelpTopic
	for _, topic := range m.topics {
		result = append(result, topic)
	}
	return result, nil
}

func (m *MockHelpRepository) IncrementViewCount(id string) error {
	for _, topic := range m.topics {
		if topic.ID == id {
			topic.ViewCount++
			return nil
		}
	}
	return domain.ErrHelpTopicNotFound
}

func (m *MockHelpRepository) InitTable() error {
	return nil
}

func (m *MockHelpRepository) SeedTopics() error {
	now := time.Now().Unix()
	for _, topic := range domain.InitialHelpTopics {
		topic.ID = "help-" + topic.Key
		topic.CreatedAt = now
		topic.UpdatedAt = now
		m.topics[topic.Key] = &topic
	}
	return nil
}

func TestHelpService_GetTopicByKey(t *testing.T) {
	repo := NewMockHelpRepository()
	service := NewHelpService(repo)

	// Seed topics
	repo.SeedTopics()

	// Test getting existing topic
	topic, err := service.GetTopicByKey("cadunico")
	if err != nil {
		t.Errorf("GetTopicByKey() unexpected error: %v", err)
	}
	if topic == nil {
		t.Fatal("GetTopicByKey() returned nil topic")
	}
	if topic.Key != "cadunico" {
		t.Errorf("Key = %s, expected 'cadunico'", topic.Key)
	}

	// Verify view count was incremented
	if topic.ViewCount != 1 {
		t.Errorf("ViewCount = %d, expected 1", topic.ViewCount)
	}
}

func TestHelpService_GetTopicByKey_NotFound(t *testing.T) {
	repo := NewMockHelpRepository()
	service := NewHelpService(repo)

	_, err := service.GetTopicByKey("non-existent")
	if err == nil {
		t.Error("GetTopicByKey() expected error for non-existent topic but got nil")
	}
}

func TestHelpService_GetTopicByID(t *testing.T) {
	repo := NewMockHelpRepository()
	service := NewHelpService(repo)

	// Add a topic with known ID
	topic := &domain.HelpTopic{
		ID:          "help-test",
		Key:         "test-topic",
		Title:       "Test Topic",
		Explanation: "Test explanation",
		Category:    domain.CategoriaGeral,
	}
	repo.Save(topic)

	// Get by ID
	found, err := service.GetTopicByID("help-test")
	if err != nil {
		t.Errorf("GetTopicByID() unexpected error: %v", err)
	}
	if found == nil {
		t.Fatal("GetTopicByID() returned nil")
	}
	if found.Key != "test-topic" {
		t.Errorf("Key = %s, expected 'test-topic'", found.Key)
	}
}

func TestHelpService_ListIndex(t *testing.T) {
	repo := NewMockHelpRepository()
	service := NewHelpService(repo)

	// Add topics in different categories
	repo.Save(&domain.HelpTopic{
		Key:         "topic1",
		Title:       "A Topic",
		Explanation: "Explanation 1",
		Category:    domain.CategoriaCredito,
	})
	repo.Save(&domain.HelpTopic{
		Key:         "topic2",
		Title:       "B Topic",
		Explanation: "Explanation 2",
		Category:    domain.CategoriaCredito,
	})
	repo.Save(&domain.HelpTopic{
		Key:         "topic3",
		Title:       "C Topic",
		Explanation: "Explanation 3",
		Category:    domain.CategoriaGeral,
	})

	// List index
	index, err := service.ListIndex()
	if err != nil {
		t.Errorf("ListIndex() unexpected error: %v", err)
	}

	// Check CREDITO category has 2 topics
	if len(index[domain.CategoriaCredito]) != 2 {
		t.Errorf("CREDITO category has %d topics, expected 2", len(index[domain.CategoriaCredito]))
	}

	// Check GERAL category has 1 topic
	if len(index[domain.CategoriaGeral]) != 1 {
		t.Errorf("GERAL category has %d topics, expected 1", len(index[domain.CategoriaGeral]))
	}

	// Verify sorting (A should come before B)
	creditoTopics := index[domain.CategoriaCredito]
	if len(creditoTopics) >= 2 && creditoTopics[0].Title != "A Topic" {
		t.Error("Topics not sorted alphabetically")
	}
}

func TestHelpService_Search(t *testing.T) {
	repo := NewMockHelpRepository()
	service := NewHelpService(repo)

	// Add topics
	repo.Save(&domain.HelpTopic{
		Key:         "cadunico",
		Title:       "O que é o CadÚnico?",
		Explanation: "O Cadastro Único reúne informações sobre famílias.",
		Category:    domain.CategoriaCredito,
		Tags:        "cadastro,programa",
	})
	repo.Save(&domain.HelpTopic{
		Key:         "cnae",
		Title:       "O que é CNAE?",
		Explanation: "É o código da atividade.",
		Category:    domain.CategoriaTributario,
		Tags:        "atividade,código",
	})

	// Search for "cadastro"
	results, err := service.Search("cadastro")
	if err != nil {
		t.Errorf("Search() unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("Search('cadastro') returned %d results, expected 1", len(results))
	}

	// Search for "CNAE"
	results, err = service.Search("CNAE")
	if err != nil {
		t.Errorf("Search() unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("Search('CNAE') returned %d results, expected 1", len(results))
	}

	// Empty search should return all
	results, err = service.Search("")
	if err != nil {
		t.Errorf("Search() unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("Search('') returned %d results, expected 2", len(results))
	}
}

func TestHelpService_GetRelatedTopics(t *testing.T) {
	repo := NewMockHelpRepository()
	service := NewHelpService(repo)

	// Add topics in same category
	repo.Save(&domain.HelpTopic{
		Key:         "topic1",
		Title:       "Topic 1",
		Explanation: "Explanation 1",
		Category:    domain.CategoriaCredito,
	})
	repo.Save(&domain.HelpTopic{
		Key:         "topic2",
		Title:       "Topic 2",
		Explanation: "Explanation 2",
		Category:    domain.CategoriaCredito,
	})
	repo.Save(&domain.HelpTopic{
		Key:         "topic3",
		Title:       "Topic 3",
		Explanation: "Explanation 3",
		Category:    domain.CategoriaCredito,
	})
	repo.Save(&domain.HelpTopic{
		Key:         "topic4",
		Title:       "Topic 4",
		Explanation: "Explanation 4",
		Category:    domain.CategoriaGeral,
	})

	currentTopic := &domain.HelpTopic{
		Key:      "topic1",
		Category: domain.CategoriaCredito,
	}

	related, err := service.GetRelatedTopics(currentTopic)
	if err != nil {
		t.Errorf("GetRelatedTopics() unexpected error: %v", err)
	}

	// Should return 2 related topics (excluding current)
	if len(related) != 2 {
		t.Errorf("GetRelatedTopics() returned %d topics, expected 2", len(related))
	}

	// Should not include the current topic
	for _, topic := range related {
		if topic.Key == "topic1" {
			t.Error("GetRelatedTopics() should not include current topic")
		}
	}
}

func TestHelpService_IncrementView(t *testing.T) {
	repo := NewMockHelpRepository()
	service := NewHelpService(repo)

	// Add topic
	repo.Save(&domain.HelpTopic{
		ID:          "help-test",
		Key:         "test-topic",
		Title:       "Test Topic",
		Explanation: "Test explanation",
		Category:    domain.CategoriaGeral,
		ViewCount:   5,
	})

	err := service.IncrementView("test-topic")
	if err != nil {
		t.Errorf("IncrementView() unexpected error: %v", err)
	}

	// Verify count was incremented
	topic, _ := repo.FindByKey("test-topic")
	if topic.ViewCount != 6 {
		t.Errorf("ViewCount = %d, expected 6", topic.ViewCount)
	}
}

func TestHelpService_GetCategories(t *testing.T) {
	repo := NewMockHelpRepository()
	service := NewHelpService(repo)

	categories := service.GetCategories()

	// Should have 4 categories
	if len(categories) != 4 {
		t.Errorf("GetCategories() returned %d categories, expected 4", len(categories))
	}

	// Check specific categories exist
	expectedCategories := []string{domain.CategoriaCredito, domain.CategoriaTributario, domain.CategoriaGovernanca, domain.CategoriaGeral}
	for _, cat := range expectedCategories {
		if _, exists := categories[cat]; !exists {
			t.Errorf("Category %s not found", cat)
		}
	}

	// Check labels are not empty
	for _, label := range categories {
		if label == "" {
			t.Error("Category label should not be empty")
		}
	}
}

func TestHelpService_Initialize(t *testing.T) {
	repo := NewMockHelpRepository()
	service := NewHelpService(repo)

	err := service.Initialize()
	if err != nil {
		t.Errorf("Initialize() unexpected error: %v", err)
	}

	// Verify topics were seeded
	topics, _ := repo.ListAll()
	if len(topics) != 6 {
		t.Errorf("After Initialize(), ListAll() returned %d topics, expected 6", len(topics))
	}
}
