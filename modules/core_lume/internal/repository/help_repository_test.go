package repository

import (
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/providentia/digna/core_lume/internal/domain"
)

// mockLifecycleManagerForHelp é um mock do LifecycleManager para testes
type mockLifecycleManagerForHelp struct {
	centralDB interface{}
}

func (m *mockLifecycleManagerForHelp) GetConnection(entityID string) (interface{}, error) {
	return nil, nil
}

func (m *mockLifecycleManagerForHelp) GetCentralConnection() (interface{}, error) {
	return m.centralDB, nil
}

func (m *mockLifecycleManagerForHelp) CloseConnection(entityID string) error {
	return nil
}

func (m *mockLifecycleManagerForHelp) CloseAll() error {
	return nil
}

func (m *mockLifecycleManagerForHelp) EntityExists(entityID string) (bool, error) {
	return true, nil
}

func (m *mockLifecycleManagerForHelp) CreateEntity(entityID, entityName string) error {
	return nil
}

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

func TestMockHelpRepository_Save(t *testing.T) {
	repo := NewMockHelpRepository()

	topic := &domain.HelpTopic{
		Key:         "test-topic",
		Title:       "Test Topic",
		Explanation: "Test explanation",
		Category:    domain.CategoriaGeral,
	}

	err := repo.Save(topic)
	if err != nil {
		t.Errorf("Save() unexpected error: %v", err)
	}

	// Verify it was saved
	saved, err := repo.FindByKey("test-topic")
	if err != nil {
		t.Errorf("FindByKey() unexpected error: %v", err)
	}
	if saved.Title != "Test Topic" {
		t.Errorf("Title = %s, expected 'Test Topic'", saved.Title)
	}
}

func TestMockHelpRepository_Save_Invalid(t *testing.T) {
	repo := NewMockHelpRepository()

	// Invalid topic - missing required fields
	topic := &domain.HelpTopic{
		Key: "test-topic",
		// Missing title, explanation, category
	}

	err := repo.Save(topic)
	if err == nil {
		t.Error("Save() expected error for invalid topic but got nil")
	}
}

func TestMockHelpRepository_FindByKey_NotFound(t *testing.T) {
	repo := NewMockHelpRepository()

	_, err := repo.FindByKey("non-existent")
	if err != domain.ErrHelpTopicNotFound {
		t.Errorf("FindByKey() error = %v, expected ErrHelpTopicNotFound", err)
	}
}

func TestMockHelpRepository_ListByCategory(t *testing.T) {
	repo := NewMockHelpRepository()

	// Add topics in different categories
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
		Category:    domain.CategoriaGeral,
	})

	// List CREDITO topics
	topics, err := repo.ListByCategory(domain.CategoriaCredito)
	if err != nil {
		t.Errorf("ListByCategory() unexpected error: %v", err)
	}
	if len(topics) != 2 {
		t.Errorf("ListByCategory(CREDITO) returned %d topics, expected 2", len(topics))
	}

	// List GERAL topics
	topics, err = repo.ListByCategory(domain.CategoriaGeral)
	if err != nil {
		t.Errorf("ListByCategory() unexpected error: %v", err)
	}
	if len(topics) != 1 {
		t.Errorf("ListByCategory(GERAL) returned %d topics, expected 1", len(topics))
	}
}

func TestMockHelpRepository_Search(t *testing.T) {
	repo := NewMockHelpRepository()

	// Add topics
	repo.Save(&domain.HelpTopic{
		Key:         "cadunico",
		Title:       "O que é o CadÚnico?",
		Explanation: "O Cadastro Único reúne informações sobre famílias.",
		Category:    domain.CategoriaCredito,
		Tags:        "cadastro,programa social",
	})
	repo.Save(&domain.HelpTopic{
		Key:         "cnae",
		Title:       "O que é CNAE?",
		Explanation: "É o código da atividade do seu negócio.",
		Category:    domain.CategoriaTributario,
		Tags:        "atividade,código",
	})

	// Search for "cadastro"
	results, err := repo.Search("cadastro")
	if err != nil {
		t.Errorf("Search() unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("Search('cadastro') returned %d results, expected 1", len(results))
	}

	// Search for "CNAE" (case insensitive)
	results, err = repo.Search("CNAE")
	if err != nil {
		t.Errorf("Search() unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("Search('CNAE') returned %d results, expected 1", len(results))
	}

	// Search for non-existent
	results, err = repo.Search("xyz123")
	if err != nil {
		t.Errorf("Search() unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("Search('xyz123') returned %d results, expected 0", len(results))
	}
}

func TestMockHelpRepository_ListAll(t *testing.T) {
	repo := NewMockHelpRepository()

	// Should be empty initially
	topics, err := repo.ListAll()
	if err != nil {
		t.Errorf("ListAll() unexpected error: %v", err)
	}
	if len(topics) != 0 {
		t.Errorf("ListAll() returned %d topics initially, expected 0", len(topics))
	}

	// Add topics
	repo.Save(&domain.HelpTopic{
		Key:         "topic1",
		Title:       "Topic 1",
		Explanation: "Explanation 1",
		Category:    domain.CategoriaGeral,
	})
	repo.Save(&domain.HelpTopic{
		Key:         "topic2",
		Title:       "Topic 2",
		Explanation: "Explanation 2",
		Category:    domain.CategoriaGeral,
	})

	// Should have 2 topics now
	topics, err = repo.ListAll()
	if err != nil {
		t.Errorf("ListAll() unexpected error: %v", err)
	}
	if len(topics) != 2 {
		t.Errorf("ListAll() returned %d topics, expected 2", len(topics))
	}
}

func TestMockHelpRepository_IncrementViewCount(t *testing.T) {
	repo := NewMockHelpRepository()

	topic := &domain.HelpTopic{
		ID:          "help-test",
		Key:         "test-topic",
		Title:       "Test Topic",
		Explanation: "Test explanation",
		Category:    domain.CategoriaGeral,
		ViewCount:   5,
	}
	repo.Save(topic)

	err := repo.IncrementViewCount("help-test")
	if err != nil {
		t.Errorf("IncrementViewCount() unexpected error: %v", err)
	}

	// Verify count was incremented
	saved, _ := repo.FindByKey("test-topic")
	if saved.ViewCount != 6 {
		t.Errorf("ViewCount = %d, expected 6", saved.ViewCount)
	}
}

func TestMockHelpRepository_IncrementViewCount_NotFound(t *testing.T) {
	repo := NewMockHelpRepository()

	err := repo.IncrementViewCount("non-existent")
	if err != domain.ErrHelpTopicNotFound {
		t.Errorf("IncrementViewCount() error = %v, expected ErrHelpTopicNotFound", err)
	}
}

func TestMockHelpRepository_SeedTopics(t *testing.T) {
	repo := NewMockHelpRepository()

	// Seed initial topics
	err := repo.SeedTopics()
	if err != nil {
		t.Errorf("SeedTopics() unexpected error: %v", err)
	}

	// Should have 6 topics
	topics, err := repo.ListAll()
	if err != nil {
		t.Errorf("ListAll() unexpected error: %v", err)
	}
	if len(topics) != 6 {
		t.Errorf("After seed, ListAll() returned %d topics, expected 6", len(topics))
	}

	// Verify specific topics exist
	for _, key := range []string{"cadunico", "inadimplencia", "cnae", "das_mei", "reserva_legal", "fates"} {
		_, err := repo.FindByKey(key)
		if err != nil {
			t.Errorf("After seed, topic '%s' not found", key)
		}
	}
}

func TestMockHelpRepository_Update(t *testing.T) {
	repo := NewMockHelpRepository()

	// Create initial topic
	topic := &domain.HelpTopic{
		Key:         "test-topic",
		Title:       "Original Title",
		Explanation: "Original explanation",
		Category:    domain.CategoriaGeral,
	}
	repo.Save(topic)

	// Update the topic
	topic.Title = "Updated Title"
	topic.Explanation = "Updated explanation"
	err := repo.Save(topic)
	if err != nil {
		t.Errorf("Save() update unexpected error: %v", err)
	}

	// Verify update
	saved, _ := repo.FindByKey("test-topic")
	if saved.Title != "Updated Title" {
		t.Errorf("Title = %s, expected 'Updated Title'", saved.Title)
	}
	if saved.Explanation != "Updated explanation" {
		t.Errorf("Explanation not updated")
	}
}
