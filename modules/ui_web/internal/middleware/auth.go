package middleware

import (
	"net/http"
	"strings"

	"github.com/providentia/digna/ui_web/internal/handler"
)

// AuthMiddleware cria um middleware de autenticação
type AuthMiddleware struct {
	authHandler *handler.AuthHandler
}

// NewAuthMiddleware cria um novo middleware de autenticação
func NewAuthMiddleware(authHandler *handler.AuthHandler) *AuthMiddleware {
	return &AuthMiddleware{
		authHandler: authHandler,
	}
}

// Handler retorna um handler HTTP que verifica autenticação
func (m *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Rotas públicas que não requerem autenticação
		publicRoutes := []string{
			"/login",
			"/api/login",
			"/api/auth/login",
			"/api/check-session",
			"/api/auth/check-session",
			"/static/",
			"/health",
			"/ready",
		}

		// Verificar se a rota é pública
		isPublic := false
		for _, route := range publicRoutes {
			if strings.HasPrefix(r.URL.Path, route) {
				isPublic = true
				break
			}
		}

		// Se não for rota pública, verificar autenticação
		if !isPublic {
			entityID, valid := m.authHandler.GetCurrentEntity(r)
			if !valid {
				// Se for uma requisição AJAX/HTMX, retornar erro JSON
				if r.Header.Get("HX-Request") == "true" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(`{"error": "Sessão expirada", "redirect": "/login"}`))
					return
				}
				// Caso contrário, redirecionar para login
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			// Adicionar entity_id à query string se não estiver presente
			// Isso garante que todos os handlers recebam o entity_id correto
			query := r.URL.Query()
			if query.Get("entity_id") == "" {
				query.Set("entity_id", entityID)
				r.URL.RawQuery = query.Encode()
			}
		}

		next.ServeHTTP(w, r)
	})
}

// WrapHandler envolve um handler específico com autenticação
func (m *AuthMiddleware) WrapHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entityID, valid := m.authHandler.GetCurrentEntity(r)
		if !valid {
			// Se for uma requisição AJAX/HTMX, retornar erro JSON
			if r.Header.Get("HX-Request") == "true" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error": "Sessão expirada", "redirect": "/login"}`))
				return
			}
			// Caso contrário, redirecionar para login
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// Adicionar entity_id à query string se não estiver presente
		query := r.URL.Query()
		if query.Get("entity_id") == "" {
			query.Set("entity_id", entityID)
			r.URL.RawQuery = query.Encode()
		}

		h.ServeHTTP(w, r)
	}
}
