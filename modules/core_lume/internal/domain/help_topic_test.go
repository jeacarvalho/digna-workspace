package domain

import (
	"testing"
	"time"
)

func TestHelpTopic_Validate(t *testing.T) {
	tests := []struct {
		name    string
		topic   HelpTopic
		wantErr bool
	}{
		{
			name: "Valid topic",
			topic: HelpTopic{
				Key:         "test-key",
				Title:       "Test Title",
				Explanation: "Test explanation",
				Category:    CategoriaGeral,
			},
			wantErr: false,
		},
		{
			name: "Missing key",
			topic: HelpTopic{
				Title:       "Test Title",
				Explanation: "Test explanation",
				Category:    CategoriaGeral,
			},
			wantErr: true,
		},
		{
			name: "Missing title",
			topic: HelpTopic{
				Key:         "test-key",
				Explanation: "Test explanation",
				Category:    CategoriaGeral,
			},
			wantErr: true,
		},
		{
			name: "Missing explanation",
			topic: HelpTopic{
				Key:      "test-key",
				Title:    "Test Title",
				Category: CategoriaGeral,
			},
			wantErr: true,
		},
		{
			name: "Missing category",
			topic: HelpTopic{
				Key:         "test-key",
				Title:       "Test Title",
				Explanation: "Test explanation",
			},
			wantErr: true,
		},
		{
			name: "Invalid category",
			topic: HelpTopic{
				Key:         "test-key",
				Title:       "Test Title",
				Explanation: "Test explanation",
				Category:    "INVALID",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.topic.Validate()
			if tt.wantErr && err == nil {
				t.Errorf("Validate() expected error but got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Validate() unexpected error: %v", err)
			}
		})
	}
}

func TestHelpTopic_IsComplete(t *testing.T) {
	tests := []struct {
		name     string
		topic    HelpTopic
		expected bool
	}{
		{
			name: "Complete topic",
			topic: HelpTopic{
				Key:         "test-key",
				Title:       "Test Title",
				Explanation: "Test explanation",
				Category:    CategoriaGeral,
			},
			expected: true,
		},
		{
			name: "Missing key",
			topic: HelpTopic{
				Title:       "Test Title",
				Explanation: "Test explanation",
				Category:    CategoriaGeral,
			},
			expected: false,
		},
		{
			name: "Missing title",
			topic: HelpTopic{
				Key:         "test-key",
				Explanation: "Test explanation",
				Category:    CategoriaGeral,
			},
			expected: false,
		},
		{
			name: "Missing explanation",
			topic: HelpTopic{
				Key:      "test-key",
				Title:    "Test Title",
				Category: CategoriaGeral,
			},
			expected: false,
		},
		{
			name: "Missing category",
			topic: HelpTopic{
				Key:         "test-key",
				Title:       "Test Title",
				Explanation: "Test explanation",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.topic.IsComplete()
			if result != tt.expected {
				t.Errorf("IsComplete() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestHelpTopic_IncrementView(t *testing.T) {
	topic := HelpTopic{
		Key:         "test-key",
		Title:       "Test Title",
		Explanation: "Test explanation",
		Category:    CategoriaGeral,
		ViewCount:   5,
	}

	initialCount := topic.ViewCount
	topic.IncrementView()

	if topic.ViewCount != initialCount+1 {
		t.Errorf("ViewCount = %d, expected %d", topic.ViewCount, initialCount+1)
	}

	if topic.UpdatedAt == 0 {
		t.Error("UpdatedAt should be set after IncrementView")
	}
}

func TestHelpTopic_GetTagsArray(t *testing.T) {
	tests := []struct {
		name     string
		tags     string
		expected []string
	}{
		{
			name:     "Multiple tags",
			tags:     "tag1,tag2,tag3",
			expected: []string{"tag1", "tag2", "tag3"},
		},
		{
			name:     "Single tag",
			tags:     "single",
			expected: []string{"single"},
		},
		{
			name:     "Empty tags",
			tags:     "",
			expected: []string{},
		},
		{
			name:     "Tags with spaces",
			tags:     "tag with space,another tag",
			expected: []string{"tag with space", "another tag"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			topic := HelpTopic{Tags: tt.tags}
			result := topic.GetTagsArray()

			if len(result) != len(tt.expected) {
				t.Errorf("GetTagsArray() returned %d items, expected %d", len(result), len(tt.expected))
				return
			}

			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("GetTagsArray()[%d] = %s, expected %s", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestHelpTopic_MatchesSearch(t *testing.T) {
	topic := HelpTopic{
		Key:         "cadunico",
		Title:       "O que é o CadÚnico?",
		Summary:     "É o cadastro do governo para programas sociais.",
		Explanation: "O Cadastro Único reúne informações sobre famílias de baixa renda.",
		Tags:        "cadastro,programa social,crédito",
		Category:    CategoriaCredito,
	}

	tests := []struct {
		name     string
		query    string
		expected bool
	}{
		{
			name:     "Match in key",
			query:    "cadunico",
			expected: true,
		},
		{
			name:     "Match in title",
			query:    "CadÚnico",
			expected: true,
		},
		{
			name:     "Match in summary",
			query:    "governo",
			expected: true,
		},
		{
			name:     "Match in explanation",
			query:    "famílias",
			expected: true,
		},
		{
			name:     "Match in tags",
			query:    "social",
			expected: true,
		},
		{
			name:     "No match",
			query:    "xyz123",
			expected: false,
		},
		{
			name:     "Empty query",
			query:    "",
			expected: true,
		},
		{
			name:     "Case insensitive",
			query:    "CADUNICO",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := topic.MatchesSearch(tt.query)
			if result != tt.expected {
				t.Errorf("MatchesSearch(%s) = %v, expected %v", tt.query, result, tt.expected)
			}
		})
	}
}

func TestHelpTopic_String(t *testing.T) {
	topic := HelpTopic{
		Key:      "test-key",
		Title:    "Test Title",
		Category: CategoriaGeral,
	}

	result := topic.String()
	if result == "" {
		t.Error("String() should not return empty string")
	}
	if result == "<nil>" {
		t.Error("String() should not return <nil>")
	}
}

func TestIsValidCategory(t *testing.T) {
	tests := []struct {
		category string
		expected bool
	}{
		{CategoriaCredito, true},
		{CategoriaTributario, true},
		{CategoriaGovernanca, true},
		{CategoriaGeral, true},
		{"INVALID", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.category, func(t *testing.T) {
			result := isValidCategory(tt.category)
			if result != tt.expected {
				t.Errorf("isValidCategory(%s) = %v, expected %v", tt.category, result, tt.expected)
			}
		})
	}
}

func TestGetCategoryLabel(t *testing.T) {
	tests := []struct {
		category string
		expected string
	}{
		{CategoriaCredito, "Crédito e Financiamento"},
		{CategoriaTributario, "Tributos e Impostos"},
		{CategoriaGovernanca, "Governança e Gestão"},
		{CategoriaGeral, "Geral"},
		{"UNKNOWN", "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.category, func(t *testing.T) {
			result := GetCategoryLabel(tt.category)
			if result != tt.expected {
				t.Errorf("GetCategoryLabel(%s) = %s, expected %s", tt.category, result, tt.expected)
			}
		})
	}
}

func TestContainsIgnoreCase(t *testing.T) {
	tests := []struct {
		s        string
		substr   string
		expected bool
	}{
		{"Hello World", "world", true},
		{"Hello World", "WORLD", true},
		{"Hello World", "Hello", true},
		{"Hello World", "xyz", false},
		{"", "test", false},
		{"test", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.s+"_"+tt.substr, func(t *testing.T) {
			result := containsIgnoreCase(tt.s, tt.substr)
			if result != tt.expected {
				t.Errorf("containsIgnoreCase(%s, %s) = %v, expected %v", tt.s, tt.substr, result, tt.expected)
			}
		})
	}
}

func TestInitialHelpTopics(t *testing.T) {
	// Verify we have exactly 6 initial topics
	if len(InitialHelpTopics) != 6 {
		t.Errorf("InitialHelpTopics has %d topics, expected 6", len(InitialHelpTopics))
	}

	// Verify all topics have required fields
	for i, topic := range InitialHelpTopics {
		if topic.Key == "" {
			t.Errorf("Topic %d has empty key", i)
		}
		if topic.Title == "" {
			t.Errorf("Topic %d has empty title", i)
		}
		if topic.Explanation == "" {
			t.Errorf("Topic %d has empty explanation", i)
		}
		if topic.Category == "" {
			t.Errorf("Topic %d has empty category", i)
		}
		if topic.WhyAsked == "" {
			t.Errorf("Topic %d has empty WhyAsked", i)
		}
		if topic.NextSteps == "" {
			t.Errorf("Topic %d has empty NextSteps", i)
		}

		// Validate the topic
		if err := topic.Validate(); err != nil {
			t.Errorf("Topic %d (%s) failed validation: %v", i, topic.Key, err)
		}
	}

	// Verify all expected keys exist
	expectedKeys := map[string]bool{
		"cadunico":      false,
		"inadimplencia": false,
		"cnae":          false,
		"das_mei":       false,
		"reserva_legal": false,
		"fates":         false,
	}

	for _, topic := range InitialHelpTopics {
		if _, exists := expectedKeys[topic.Key]; exists {
			expectedKeys[topic.Key] = true
		}
	}

	for key, found := range expectedKeys {
		if !found {
			t.Errorf("Expected topic with key '%s' not found in InitialHelpTopics", key)
		}
	}
}

func TestHelpTopic_Timestamps(t *testing.T) {
	now := time.Now().Unix()

	topic := HelpTopic{
		Key:         "test-key",
		Title:       "Test Title",
		Explanation: "Test explanation",
		Category:    CategoriaGeral,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if topic.CreatedAt != now {
		t.Errorf("CreatedAt = %d, expected %d", topic.CreatedAt, now)
	}
	if topic.UpdatedAt != now {
		t.Errorf("UpdatedAt = %d, expected %d", topic.UpdatedAt, now)
	}
}
