package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

type AccountantLinkHandler struct {
	*BaseHandler
	lifecycleManager lifecycle.LifecycleManager
	authHandler      *AuthHandler
}

func NewAccountantLinkHandler(lm lifecycle.LifecycleManager, authHandler *AuthHandler) (*AccountantLinkHandler, error) {
	// Obter devMode do ambiente
	devMode := os.Getenv("DEV") != "false" && os.Getenv("DEV") != "0"
	baseHandler := NewBaseHandler(lm, devMode)

	return &AccountantLinkHandler{
		BaseHandler:      baseHandler,
		lifecycleManager: lm,
		authHandler:      authHandler,
	}, nil
}

func (h *AccountantLinkHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/accountant/links", h.ListLinks)
	mux.HandleFunc("/accountant/links/create", h.CreateLink)
	mux.HandleFunc("/accountant/links/deactivate", h.DeactivateLink)
}

// ListLinks mostra a lista de vínculos contábeis
func (h *AccountantLinkHandler) ListLinks(w http.ResponseWriter, r *http.Request) {
	// Verificar autenticação
	entityID, valid := h.authHandler.GetCurrentEntity(r)
	if !valid {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Verificar tipo de usuário
	userType, _ := h.authHandler.GetCurrentUserType(r)

	data := map[string]interface{}{
		"Title":    "Gerenciar Vínculos Contábeis",
		"UserType": userType,
		"EntityID": entityID,
		"Links":    []map[string]interface{}{}, // Será preenchido quando tivermos repositório
	}

	// Carregar template
	tmpl, err := h.loadTemplate("accountant_link_simple.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("template error: %v", err), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, fmt.Sprintf("template error: %v", err), http.StatusInternalServerError)
	}
}

// CreateLink cria um novo vínculo contábil
func (h *AccountantLinkHandler) CreateLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Verificar autenticação
	entityID, valid := h.authHandler.GetCurrentEntity(r)
	if !valid {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Verificar se é empreendimento (apenas empreendimentos podem criar vínculos)
	userType, _ := h.authHandler.GetCurrentUserType(r)
	if userType != "empreendimento" {
		http.Error(w, "Apenas empreendimentos podem criar vínculos contábeis", http.StatusForbidden)
		return
	}

	// Parse form
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao processar formulário: %v", err), http.StatusBadRequest)
		return
	}

	accountantID := r.FormValue("accountant_id")
	if accountantID == "" {
		http.Error(w, "ID do contador é obrigatório", http.StatusBadRequest)
		return
	}

	// Verificar se lifecycleManager implementa AccountantLinkService
	accountantLinkService, ok := h.lifecycleManager.(lifecycle.AccountantLinkService)
	if !ok {
		http.Error(w, "Serviço de vínculos contábeis não disponível", http.StatusInternalServerError)
		return
	}

	// Criar vínculo
	_, err := accountantLinkService.CreateLink(entityID, accountantID, entityID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao criar vínculo: %v", err), http.StatusInternalServerError)
		return
	}

	// Redirecionar para lista com mensagem de sucesso
	http.Redirect(w, r, "/accountant/links?success=Vínculo criado com sucesso", http.StatusFound)
}

// DeactivateLink desativa um vínculo contábil (Exit Power)
func (h *AccountantLinkHandler) DeactivateLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Verificar autenticação
	entityID, valid := h.authHandler.GetCurrentEntity(r)
	if !valid {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Parse form
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao processar formulário: %v", err), http.StatusBadRequest)
		return
	}

	linkID := r.FormValue("link_id")
	if linkID == "" {
		http.Error(w, "ID do vínculo é obrigatório", http.StatusBadRequest)
		return
	}

	// Verificar se lifecycleManager implementa AccountantLinkService
	accountantLinkService, ok := h.lifecycleManager.(lifecycle.AccountantLinkService)
	if !ok {
		http.Error(w, "Serviço de vínculos contábeis não disponível", http.StatusInternalServerError)
		return
	}

	// Desativar vínculo (Exit Power - apenas a cooperativa que criou pode desativar)
	err := accountantLinkService.DeactivateLink(linkID, entityID, entityID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao desativar vínculo: %v", err), http.StatusInternalServerError)
		return
	}

	// Redirecionar para lista com mensagem de sucesso
	http.Redirect(w, r, "/accountant/links?success=Vínculo desativado com sucesso", http.StatusFound)
}

// loadTemplate carrega template cache-proof
func (h *AccountantLinkHandler) loadTemplate(filename string) (*template.Template, error) {
	// Tentar múltiplos caminhos para funcionar em diferentes ambientes
	var tmpl *template.Template
	var err error

	// Tentativa 1: Caminho relativo (quando executado de modules/ui_web/)
	tmpl, err = template.ParseFiles("templates/" + filename)
	if err != nil {
		// Tentativa 2: Caminho absoluto do projeto
		tmpl, err = template.ParseFiles("modules/ui_web/templates/" + filename)
		if err != nil {
			// Tentativa 3: Caminho relativo alternativo
			tmpl, err = template.ParseFiles("../../templates/" + filename)
			if err != nil {
				return nil, fmt.Errorf("não foi possível carregar o template: %v", err)
			}
		}
	}

	return tmpl, nil
}
