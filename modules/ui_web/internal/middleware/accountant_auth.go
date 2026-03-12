package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/providentia/digna/ui_web/internal/handler"
)

// AccountantAuthMiddleware cria um middleware específico para contadores
type AccountantAuthMiddleware struct {
	authHandler *handler.AuthHandler
}

// NewAccountantAuthMiddleware cria um novo middleware para contadores
func NewAccountantAuthMiddleware(authHandler *handler.AuthHandler) *AccountantAuthMiddleware {
	return &AccountantAuthMiddleware{
		authHandler: authHandler,
	}
}

// Handler retorna um handler HTTP que verifica se o usuário é contador
func (m *AccountantAuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verificar se a rota é do painel do contador
		if !strings.HasPrefix(r.URL.Path, "/accountant") {
			// Se não for rota do contador, passar adiante
			next.ServeHTTP(w, r)
			return
		}

		fmt.Printf("[ACCOUNTANT MIDDLEWARE] Checking accountant route: %s\n", r.URL.Path)

		// Verificar autenticação
		entityID, valid := m.authHandler.GetCurrentEntity(r)
		if !valid {
			// Redirecionar para login
			fmt.Printf("[ACCOUNTANT MIDDLEWARE] Not authenticated, redirecting to login\n")
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// Verificar se é contador
		userType, _ := m.authHandler.GetCurrentUserType(r)
		fmt.Printf("[ACCOUNTANT MIDDLEWARE] Authenticated: entityID=%s, userType=%s\n", entityID, userType)

		if userType != "contador" {
			// Se não for contador, redirecionar para dashboard do empreendimento
			fmt.Printf("[ACCOUNTANT MIDDLEWARE] Not a contador, redirecting to dashboard\n")
			http.Redirect(w, r, "/dashboard?entity_id="+entityID, http.StatusFound)
			return
		}

		// Se for contador, permitir acesso
		fmt.Printf("[ACCOUNTANT MIDDLEWARE] Access granted to contador\n")
		next.ServeHTTP(w, r)
	})
}

// EmpreendimentoAuthMiddleware cria um middleware específico para empreendimentos
type EmpreendimentoAuthMiddleware struct {
	authHandler *handler.AuthHandler
}

// NewEmpreendimentoAuthMiddleware cria um novo middleware para empreendimentos
func NewEmpreendimentoAuthMiddleware(authHandler *handler.AuthHandler) *EmpreendimentoAuthMiddleware {
	return &EmpreendimentoAuthMiddleware{
		authHandler: authHandler,
	}
}

// Handler retorna um handler HTTP que verifica se o usuário é empreendimento
func (m *EmpreendimentoAuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Rotas que só empreendimentos podem acessar
		empreendimentoOnlyRoutes := []string{
			"/pdv",
			"/cash",
			"/supply",
			"/members",
			"/budget",
		}

		// Verificar se a rota é exclusiva de empreendimento
		isEmpreendimentoRoute := false
		for _, route := range empreendimentoOnlyRoutes {
			if strings.HasPrefix(r.URL.Path, route) {
				isEmpreendimentoRoute = true
				break
			}
		}

		if !isEmpreendimentoRoute {
			// Se não for rota exclusiva de empreendimento, passar adiante
			next.ServeHTTP(w, r)
			return
		}

		// Verificar autenticação
		_, valid := m.authHandler.GetCurrentEntity(r)
		if !valid {
			// Redirecionar para login
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// Verificar se é empreendimento
		userType, _ := m.authHandler.GetCurrentUserType(r)
		if userType != "empreendimento" {
			// Se for contador tentando acessar rota de empreendimento, redirecionar para painel do contador
			http.Redirect(w, r, "/accountant/dashboard", http.StatusFound)
			return
		}

		// Se for empreendimento, permitir acesso
		next.ServeHTTP(w, r)
	})
}
