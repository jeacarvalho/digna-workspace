package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

// TestConfig configuração para testes isolados
type TestConfig struct {
	EntityID     string
	DataDir      string
	LifecycleMgr lifecycle.LifecycleManager
	CleanupFunc  func()
}

// SetupTestEnvironment cria ambiente de teste isolado
func SetupTestEnvironment(t testLogger) *TestConfig {
	// Criar diretório temporário único para cada teste
	testID := fmt.Sprintf("test_%d", time.Now().UnixNano())
	dataDir := filepath.Join("../../data/test_entities", testID)

	// Garantir que o diretório existe
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		t.Fatalf("failed to create test data directory: %v", err)
	}

	// Configurar lifecycle manager
	lifecycleMgr := lifecycle.NewSQLiteManager()

	// Entity ID único para teste
	entityID := fmt.Sprintf("test_entity_%s", testID)

	// Função de cleanup
	cleanup := func() {
		lifecycleMgr.CloseAll()
		// Remover diretório de teste
		os.RemoveAll(dataDir)
	}

	return &TestConfig{
		EntityID:     entityID,
		DataDir:      dataDir,
		LifecycleMgr: lifecycleMgr,
		CleanupFunc:  cleanup,
	}
}

// Interface para logging nos testes
type testLogger interface {
	Fatalf(format string, args ...interface{})
	Logf(format string, args ...interface{})
}
