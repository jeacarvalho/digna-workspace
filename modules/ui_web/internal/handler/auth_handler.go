package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"sync"
	"time"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

// AuthHandler gerencia autenticação e sessões
type AuthHandler struct {
	lifecycleManager lifecycle.LifecycleManager
	tmpl             *template.Template
	sessions         map[string]Session
	sessionMutex     sync.RWMutex
}

// Session representa uma sessão de usuário
type Session struct {
	EntityID   string
	EntityName string
	CreatedAt  time.Time
	LastAccess time.Time
}

// Empresas de teste pré-configuradas
var testCompanies = map[string]struct {
	Name     string
	Password string
}{
	"cafe_digna": {
		Name:     "Café Digna",
		Password: "cd0123",
	},
	"queijaria_digna": {
		Name:     "Queijaria Digna",
		Password: "qd321",
	},
}

// NewAuthHandler cria um novo handler de autenticação
func NewAuthHandler(lm lifecycle.LifecycleManager) (*AuthHandler, error) {
	funcMap := template.FuncMap{
		"formatCurrency": func(amount int64) string {
			return fmt.Sprintf("R$ %.2f", float64(amount)/100)
		},
		"formatDate": func(t time.Time) string {
			return t.Format("02/01/2006 15:04")
		},
		"formatDateShort": func(t time.Time) string {
			return t.Format("02/01/2006")
		},
		"divide": func(a, b int64) float64 {
			if b == 0 {
				return 0
			}
			return float64(a) / float64(b)
		},
		"multiply": func(a, b int64) int64 {
			return a * b
		},
		"getAlertStatusLabel": func(status string) string {
			switch status {
			case "SAFE":
				return "Dentro do planejado"
			case "WARNING":
				return "Atenção: perto do limite"
			case "EXCEEDED":
				return "Ultrapassou o planejado"
			default:
				return "Sem dados"
			}
		},
		"getCategoryLabel": func(category string) string {
			labels := map[string]string{
				"INSUMOS":      "Insumos",
				"ENERGIA":      "Energia",
				"EQUIPAMENTOS": "Equipamentos",
				"TRANSPORTE":   "Transporte",
				"MANUTENCAO":   "Manutenção",
				"SERVICOS":     "Serviços",
				"OUTROS":       "Outros",
			}
			if label, ok := labels[category]; ok {
				return label
			}
			return category
		},
		"getAlertStatusClass": func(status string) string {
			switch status {
			case "SAFE":
				return "bg-green-100 text-green-800 border-green-300"
			case "WARNING":
				return "bg-yellow-100 text-yellow-800 border-yellow-300"
			case "EXCEEDED":
				return "bg-red-100 text-red-800 border-red-300"
			default:
				return "bg-gray-100 text-gray-800 border-gray-300"
			}
		},
	}

	// Parsear apenas o template simples de login
	tmpl, err := template.New("login_simple.html").Funcs(funcMap).ParseFiles("templates/login_simple.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse login template: %w", err)
	}

	return &AuthHandler{
		lifecycleManager: lm,
		tmpl:             tmpl,
		sessions:         make(map[string]Session),
	}, nil
}

// RegisterRoutes registra as rotas de autenticação
func (h *AuthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/login", h.LoginPage)
	mux.HandleFunc("/logout", h.Logout)
	mux.HandleFunc("/api/login", h.HandleLogin)
	mux.HandleFunc("/api/check-session", h.CheckSession)
}

// LoginPage exibe a página de login
func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Verificar se já está autenticado
	if sessionID := h.getSessionID(r); sessionID != "" {
		h.sessionMutex.RLock()
		session, exists := h.sessions[sessionID]
		h.sessionMutex.RUnlock()

		if exists && time.Since(session.LastAccess) < 24*time.Hour {
			// Redirecionar para dashboard
			http.Redirect(w, r, fmt.Sprintf("/dashboard?entity_id=%s", session.EntityID), http.StatusFound)
			return
		}
	}

	data := map[string]interface{}{
		"Title": "Login - Digna",
		"Companies": []map[string]string{
			{"id": "cafe_digna", "name": "Café Digna"},
			{"id": "queijaria_digna", "name": "Queijaria Digna"},
		},
	}

	if err := h.tmpl.ExecuteTemplate(w, "login_simple.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// HandleLogin processa o login
func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var entityID, password string

	// Tentar ler como JSON primeiro
	contentType := r.Header.Get("Content-Type")
	if contentType == "application/json" {
		var loginData struct {
			EntityID string `json:"entity_id"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"message": fmt.Sprintf("Erro ao ler dados de login: %s", err.Error()),
			})
			return
		}
		entityID = loginData.EntityID
		password = loginData.Password
	} else {
		// Form data tradicional
		entityID = r.FormValue("entity_id")
		password = r.FormValue("password")
	}

	// Validar credenciais
	company, exists := testCompanies[entityID]

	if !exists || company.Password != password {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Credenciais inválidas",
		})
		return
	}

	// Criar banco de dados se não existir
	if err := h.ensureDatabaseExists(entityID, company.Name); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": fmt.Sprintf("Erro ao criar banco de dados: %s", err.Error()),
		})
		return
	}

	// Criar sessão
	sessionID := h.createSession(entityID, company.Name)

	// Configurar cookie
	cookie := &http.Cookie{
		Name:     "digna_session",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Em produção deve ser true (HTTPS)
		SameSite: http.SameSiteStrictMode,
		MaxAge:   86400, // 24 horas
	}
	http.SetCookie(w, cookie)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"redirect": fmt.Sprintf("/dashboard?entity_id=%s", entityID),
	})
}

// ensureDatabaseExists cria o banco de dados da empresa se não existir
func (h *AuthHandler) ensureDatabaseExists(entityID, entityName string) error {
	// Verificar se o banco já existe
	exists, err := h.lifecycleManager.EntityExists(entityID)
	if err != nil {
		return fmt.Errorf("erro ao verificar existência da entidade: %w", err)
	}

	if !exists {
		// Criar nova entidade
		if err := h.lifecycleManager.CreateEntity(entityID, entityName); err != nil {
			return fmt.Errorf("erro ao criar entidade: %w", err)
		}

		// Inicializar dados básicos se necessário
		// (isso pode ser expandido para criar contas padrão, etc.)
		fmt.Printf("✅ Banco de dados criado para: %s (%s)\n", entityName, entityID)
	}

	return nil
}

// createSession cria uma nova sessão
func (h *AuthHandler) createSession(entityID, entityName string) string {
	h.sessionMutex.Lock()
	defer h.sessionMutex.Unlock()

	// Gerar ID único para sessão
	timestamp := time.Now().UnixNano()
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s-%d", entityID, timestamp)))
	sessionID := hex.EncodeToString(hash[:16])

	h.sessions[sessionID] = Session{
		EntityID:   entityID,
		EntityName: entityName,
		CreatedAt:  time.Now(),
		LastAccess: time.Now(),
	}

	// Limpar sessões antigas (mais de 24 horas)
	go h.cleanupOldSessions()

	return sessionID
}

// getSessionID obtém o ID da sessão do cookie
func (h *AuthHandler) getSessionID(r *http.Request) string {
	cookie, err := r.Cookie("digna_session")
	if err != nil {
		return ""
	}
	return cookie.Value
}

// GetCurrentEntity obtém a entidade atual da sessão
func (h *AuthHandler) GetCurrentEntity(r *http.Request) (string, bool) {
	sessionID := h.getSessionID(r)
	if sessionID == "" {
		return "", false
	}

	h.sessionMutex.RLock()
	session, exists := h.sessions[sessionID]
	h.sessionMutex.RUnlock()

	if !exists {
		return "", false
	}

	// Atualizar último acesso
	h.sessionMutex.Lock()
	session.LastAccess = time.Now()
	h.sessions[sessionID] = session
	h.sessionMutex.Unlock()

	return session.EntityID, true
}

// CheckSession verifica se a sessão é válida
func (h *AuthHandler) CheckSession(w http.ResponseWriter, r *http.Request) {
	entityID, valid := h.GetCurrentEntity(r)

	w.Header().Set("Content-Type", "application/json")
	if valid {
		fmt.Fprintf(w, `{"valid": true, "entity_id": "%s"}`, entityID)
	} else {
		fmt.Fprintf(w, `{"valid": false}`)
	}
}

// Logout encerra a sessão
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	sessionID := h.getSessionID(r)
	if sessionID != "" {
		h.sessionMutex.Lock()
		delete(h.sessions, sessionID)
		h.sessionMutex.Unlock()
	}

	// Limpar cookie
	cookie := &http.Cookie{
		Name:     "digna_session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/login", http.StatusFound)
}

// cleanupOldSessions remove sessões antigas
func (h *AuthHandler) cleanupOldSessions() {
	h.sessionMutex.Lock()
	defer h.sessionMutex.Unlock()

	cutoff := time.Now().Add(-24 * time.Hour)
	for sessionID, session := range h.sessions {
		if session.LastAccess.Before(cutoff) {
			delete(h.sessions, sessionID)
		}
	}
}

// Middleware para proteger rotas
func (h *AuthHandler) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, valid := h.GetCurrentEntity(r)
		if !valid {
			// Se for uma requisição AJAX/HTMX, retornar erro JSON
			if r.Header.Get("HX-Request") == "true" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintf(w, `{"error": "Sessão expirada", "redirect": "/login"}`)
				return
			}
			// Caso contrário, redirecionar para login
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		next(w, r)
	}
}
