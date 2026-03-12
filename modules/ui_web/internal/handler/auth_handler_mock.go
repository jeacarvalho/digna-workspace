package handler

import (
	"html/template"
	"net/http"
	"time"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

// MockAuthHandler é um mock do AuthHandler para testes
type MockAuthHandler struct {
	lifecycleManager lifecycle.LifecycleManager
	tmpl             *template.Template
}

// NewMockAuthHandler cria um novo mock do AuthHandler
func NewMockAuthHandler(lm lifecycle.LifecycleManager) *AuthHandler {
	// Criar um AuthHandler real mas com comportamento mockado
	authHandler := &AuthHandler{
		lifecycleManager: lm,
		tmpl:             template.New("mock"),
		sessions:         make(map[string]Session),
	}

	// Adicionar sessão mock para testes
	sessionID := "test_session"
	authHandler.sessions[sessionID] = Session{
		EntityID:   "contador_social",
		EntityName: "Contador Social Test",
		UserType:   "contador",
		CreatedAt:  time.Now(),
		LastAccess: time.Now(),
	}

	return authHandler
}

// RegisterRoutes implementa interface mínima para testes
func (h *MockAuthHandler) RegisterRoutes(mux *http.ServeMux) {
	// Não faz nada em testes
}

// LoginPage implementa interface mínima para testes
func (h *MockAuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	// Não faz nada em testes
}

// Login implementa interface mínima para testes
func (h *MockAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Não faz nada em testes
}

// Logout implementa interface mínima para testes
func (h *MockAuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Não faz nada em testes
}

// GetSession implementa interface mínima para testes
func (h *MockAuthHandler) GetSession(r *http.Request) (Session, bool) {
	// Retorna sessão mock para testes
	return Session{
		EntityID:   "test_coop",
		EntityName: "Test Cooperative",
		UserType:   "empreendimento",
		CreatedAt:  time.Now(),
		LastAccess: time.Now(),
	}, true
}

// GetCurrentEntity implementa interface para testes
func (h *MockAuthHandler) GetCurrentEntity(r *http.Request) (string, bool) {
	// Retorna entityID mock para testes
	return "contador_social", true
}

// GetCurrentUserType implementa interface para testes
func (h *MockAuthHandler) GetCurrentUserType(r *http.Request) (string, bool) {
	// Retorna userType mock para testes (contador)
	return "contador", true
}

// RequireAuth implementa interface mínima para testes
func (h *MockAuthHandler) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Em testes, sempre passa para o próximo handler
		next(w, r)
	}
}

// RequireEntityID implementa interface mínima para testes
func (h *MockAuthHandler) RequireEntityID(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Em testes, sempre passa para o próximo handler
		next(w, r)
	}
}
