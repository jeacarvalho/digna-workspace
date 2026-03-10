package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

// Tipos locais para membros (simulando domain do core_lume)
type MemberRole string
type MemberStatus string

const (
	RoleCoordinator MemberRole = "COORDINATOR"
	RoleMember      MemberRole = "MEMBER"
	RoleAdvisor     MemberRole = "ADVISOR"
)

const (
	StatusActive   MemberStatus = "ACTIVE"
	StatusInactive MemberStatus = "INACTIVE"
)

type Member struct {
	ID        string
	EntityID  string
	Name      string
	Email     string
	Phone     string
	CPF       string // Opcional - LGPD
	Role      MemberRole
	Status    MemberStatus
	JoinedAt  time.Time
	Skills    []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// MemberHandler gerencia a interface web para gestão de membros/sócios
type MemberHandler struct {
	*BaseHandler
	lifecycleManager lifecycle.LifecycleManager
	tmpl             *template.Template
	// TODO: Adicionar MemberService quando tiver acesso ao core_lume
}

// NewMemberHandler cria um novo handler para gestão de membros
func NewMemberHandler(lm lifecycle.LifecycleManager) (*MemberHandler, error) {
	// Criar BaseHandler com TemplateManager
	base := NewBaseHandler(lm, true)

	// Adicionar funções específicas para membros ao TemplateManager
	base.templateManager.AddFunc("getRoleLabel", func(role MemberRole) string {
		switch role {
		case RoleCoordinator:
			return "Coordenador(a)"
		case RoleMember:
			return "Sócio(a)"
		case RoleAdvisor:
			return "Conselheiro(a)"
		default:
			return string(role)
		}
	})

	base.templateManager.AddFunc("getStatusLabel", func(status MemberStatus) string {
		switch status {
		case StatusActive:
			return "Ativo"
		case StatusInactive:
			return "Inativo"
		default:
			return string(status)
		}
	})

	base.templateManager.AddFunc("getStatusClass", func(status MemberStatus) string {
		switch status {
		case StatusActive:
			return "bg-green-100 text-green-800"
		case StatusInactive:
			return "bg-gray-100 text-gray-800"
		default:
			return "bg-gray-100 text-gray-800"
		}
	})

	base.templateManager.AddFunc("getRoleClass", func(role MemberRole) string {
		switch role {
		case RoleCoordinator:
			return "bg-blue-100 text-blue-800"
		case RoleMember:
			return "bg-green-100 text-green-800"
		case RoleAdvisor:
			return "bg-purple-100 text-purple-800"
		default:
			return "bg-gray-100 text-gray-800"
		}
	})

	base.templateManager.AddFunc("joinSkills", func(skills []string) string {
		return strings.Join(skills, ", ")
	})

	// Adicionar formatDate específico para time.Time (sobrescreve o do BaseHandler)
	base.templateManager.AddFunc("formatDate", func(t interface{}) string {
		switch v := t.(type) {
		case time.Time:
			return v.Format("02/01/2006")
		case string:
			return v
		default:
			return fmt.Sprintf("%v", t)
		}
	})

	// Criar template vazio - será carregado em tempo de execução
	tmpl := template.New("")

	return &MemberHandler{
		BaseHandler:      base,
		lifecycleManager: lm,
		tmpl:             tmpl,
	}, nil
}

// RegisterRoutes registra as rotas do handler
func (h *MemberHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/members", h.MembersPage)
	mux.HandleFunc("/members/create", h.CreateMember)
	mux.HandleFunc("/members/", h.HandleMemberActions) // Para /members/{id}/...
}

// MembersPage renderiza a página principal de gestão de membros
func (h *MemberHandler) MembersPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	entityID := r.URL.Query().Get("entity_id")
	if entityID == "" {
		http.Error(w, "entity_id é obrigatório", http.StatusBadRequest)
		return
	}

	// TODO: Buscar membros reais do serviço
	// Por enquanto, dados mock para desenvolvimento
	members := []Member{
		{
			ID:        "membro_001",
			EntityID:  entityID,
			Name:      "Maria Silva",
			Email:     "maria@cooperativa.com",
			Phone:     "(11) 99999-9999",
			CPF:       "", // Opcional - LGPD
			Role:      RoleCoordinator,
			Status:    StatusActive,
			JoinedAt:  time.Now().AddDate(0, -6, 0),
			Skills:    []string{"Gestão", "Vendas", "Produção"},
			CreatedAt: time.Now().AddDate(0, -6, 0),
			UpdatedAt: time.Now().AddDate(0, -6, 0),
		},
		{
			ID:        "membro_002",
			EntityID:  entityID,
			Name:      "João Santos",
			Email:     "joao@cooperativa.com",
			Phone:     "(11) 98888-8888",
			CPF:       "",
			Role:      RoleMember,
			Status:    StatusActive,
			JoinedAt:  time.Now().AddDate(0, -3, 0),
			Skills:    []string{"Logística", "Transporte"},
			CreatedAt: time.Now().AddDate(0, -3, 0),
			UpdatedAt: time.Now().AddDate(0, -3, 0),
		},
		{
			ID:        "membro_003",
			EntityID:  entityID,
			Name:      "Ana Oliveira",
			Email:     "ana@cooperativa.com",
			Phone:     "(11) 97777-7777",
			CPF:       "",
			Role:      RoleAdvisor,
			Status:    StatusInactive,
			JoinedAt:  time.Now().AddDate(-1, 0, 0),
			Skills:    []string{"Contabilidade", "Jurídico"},
			CreatedAt: time.Now().AddDate(-1, 0, 0),
			UpdatedAt: time.Now().AddDate(0, -1, 0),
		},
	}

	data := map[string]interface{}{
		"Title":    "Gestão de Sócios",
		"EntityID": entityID,
		"Members":  members,
		"Roles": []MemberRole{
			RoleCoordinator,
			RoleMember,
			RoleAdvisor,
		},
	}

	// Carregar template com funções específicas do MemberHandler
	tmpl, err := template.New("members_simple.html").Funcs(template.FuncMap{
		"getRoleLabel": func(role MemberRole) string {
			switch role {
			case RoleCoordinator:
				return "Coordenador(a)"
			case RoleMember:
				return "Sócio(a)"
			case RoleAdvisor:
				return "Conselheiro(a)"
			default:
				return string(role)
			}
		},
		"getStatusLabel": func(status MemberStatus) string {
			switch status {
			case StatusActive:
				return "Ativo"
			case StatusInactive:
				return "Inativo"
			default:
				return string(status)
			}
		},
		"getStatusClass": func(status MemberStatus) string {
			switch status {
			case StatusActive:
				return "bg-green-100 text-green-800"
			case StatusInactive:
				return "bg-gray-100 text-gray-800"
			default:
				return "bg-gray-100 text-gray-800"
			}
		},
		"getRoleClass": func(role MemberRole) string {
			switch role {
			case RoleCoordinator:
				return "bg-blue-100 text-blue-800"
			case RoleMember:
				return "bg-green-100 text-green-800"
			case RoleAdvisor:
				return "bg-purple-100 text-purple-800"
			default:
				return "bg-gray-100 text-gray-800"
			}
		},
		"joinSkills": func(skills []string) string {
			return strings.Join(skills, ", ")
		},
		"formatDate": func(t interface{}) string {
			switch v := t.(type) {
			case time.Time:
				return v.Format("02/01/2006")
			case string:
				return v
			default:
				return fmt.Sprintf("%v", t)
			}
		},
	}).ParseFiles("templates/members_simple.html")

	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao carregar template: %v", err), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// CreateMember processa a criação de um novo membro via POST
func (h *MemberHandler) CreateMember(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parsear formulário
	if err := r.ParseForm(); err != nil {
		renderHTMXError(w, "Erro ao processar formulário")
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	_ = r.FormValue("cpf") // CPF opcional - LGPD
	role := MemberRole(r.FormValue("role"))
	skillsStr := r.FormValue("skills")

	// Validar campos obrigatórios
	if name == "" || email == "" {
		renderHTMXError(w, "Nome e email são obrigatórios")
		return
	}

	// Validar role
	if role != RoleCoordinator && role != RoleMember && role != RoleAdvisor {
		role = RoleMember // Default
	}

	// Processar skills
	var skills []string
	if skillsStr != "" {
		skills = strings.Split(skillsStr, ",")
		for i := range skills {
			skills[i] = strings.TrimSpace(skills[i])
		}
	}

	// TODO: Chamar serviço real para criar membro
	// Por enquanto, retornar sucesso mock
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("HX-Trigger", "member-created")

	// Retornar card do novo membro (mock)
	// Usar funções de template inline
	getRoleLabel := func(role MemberRole) string {
		switch role {
		case RoleCoordinator:
			return "Coordenador(a)"
		case RoleMember:
			return "Sócio(a)"
		case RoleAdvisor:
			return "Conselheiro(a)"
		default:
			return string(role)
		}
	}

	getRoleClass := func(role MemberRole) string {
		switch role {
		case RoleCoordinator:
			return "bg-blue-100 text-blue-800"
		case RoleMember:
			return "bg-green-100 text-green-800"
		case RoleAdvisor:
			return "bg-purple-100 text-purple-800"
		default:
			return "bg-gray-100 text-gray-800"
		}
	}

	getStatusClass := func(status MemberStatus) string {
		switch status {
		case StatusActive:
			return "bg-green-100 text-green-800"
		case StatusInactive:
			return "bg-gray-100 text-gray-800"
		default:
			return "bg-gray-100 text-gray-800"
		}
	}

	html := fmt.Sprintf(`
	<div class="digna-card p-6 mb-4" id="member-mock_new">
		<div class="flex justify-between items-start mb-4">
			<div>
				<h3 class="font-bold text-lg text-digna-text">%s</h3>
				<p class="text-gray-600">%s • %s</p>
				<div class="mt-2 flex flex-wrap gap-2">
					<span class="px-3 py-1 rounded-full text-sm %s">%s</span>
					<span class="px-3 py-1 rounded-full text-sm %s">Ativo</span>
				</div>
			</div>
			<div class="flex space-x-2">
				<button class="text-blue-600 hover:text-blue-800" title="Editar">
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
					</svg>
				</button>
				<button class="text-green-600 hover:text-green-800" title="Ativar/Inativar">
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
					</svg>
				</button>
			</div>
		</div>
		<div class="border-t pt-4">
			<p class="text-sm text-gray-600"><strong>Habilidades:</strong> %s</p>
			<p class="text-sm text-gray-600 mt-1"><strong>Entrou em:</strong> %s</p>
		</div>
	</div>
	`,
		name, email, phone,
		getRoleClass(role), getRoleLabel(role),
		getStatusClass(StatusActive),
		strings.Join(skills, ", "),
		time.Now().Format("02/01/2006"))

	// Adicionar mensagem de sucesso
	html = `
	<div class="bg-green-100 border-l-4 border-green-500 text-green-700 p-4 mb-4" role="alert">
		<p class="font-bold">Sócio cadastrado com sucesso!</p>
		<p>O novo sócio foi adicionado à cooperativa.</p>
	</div>
	` + html

	fmt.Fprint(w, html)
}

// HandleMemberActions gerencia ações em membros específicos (toggle status, editar, etc.)
func (h *MemberHandler) HandleMemberActions(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/members/")
	parts := strings.Split(path, "/")

	if len(parts) < 2 {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	memberID := parts[0]
	action := parts[1]

	switch action {
	case "toggle-status":
		h.ToggleMemberStatus(w, r, memberID)
	case "edit":
		h.EditMember(w, r, memberID)
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

// ToggleMemberStatus alterna o status de um membro (ativo/inativo)
func (h *MemberHandler) ToggleMemberStatus(w http.ResponseWriter, r *http.Request, memberID string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// TODO: Implementar lógica real com serviço
	// Por enquanto, mock de sucesso/erro

	// Simular erro se for o último coordenador
	if memberID == "membro_001" {
		renderHTMXError(w, "Não é possível inativar o último coordenador da cooperativa")
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("HX-Trigger", "member-status-toggled")

	// Retornar novo status (mock)
	html := `
	<div class="bg-green-100 border-l-4 border-green-500 text-green-700 p-4 mb-4" role="alert">
		<p>Status do sócio atualizado com sucesso!</p>
	</div>
	`
	fmt.Fprint(w, html)
}

// EditMember processa a edição de um membro
func (h *MemberHandler) EditMember(w http.ResponseWriter, r *http.Request, memberID string) {
	if r.Method == http.MethodGet {
		h.renderEditForm(w, r, memberID)
	} else if r.Method == http.MethodPost {
		h.updateMember(w, r, memberID)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// renderEditForm renderiza o formulário de edição
func (h *MemberHandler) renderEditForm(w http.ResponseWriter, r *http.Request, memberID string) {
	// TODO: Buscar membro real
	// Por enquanto, mock
	w.Header().Set("Content-Type", "text/html")

	html := `
	<div class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
		<div class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white">
			<div class="flex justify-between items-center mb-4">
				<h3 class="text-lg font-bold text-digna-text">Editar Sócio</h3>
				<button onclick="this.closest('.fixed').remove()" class="text-gray-400 hover:text-gray-600">
					<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
					</svg>
				</button>
			</div>
			
			<form hx-post="/members/` + memberID + `/edit" hx-target="#member-` + memberID + `" hx-swap="outerHTML">
				<div class="space-y-4">
					<div>
						<label class="block text-sm font-medium text-gray-700 mb-1">Nome *</label>
						<input type="text" name="name" value="Maria Silva" required
							class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-digna-primary">
					</div>
					
					<div>
						<label class="block text-sm font-medium text-gray-700 mb-1">Email *</label>
						<input type="email" name="email" value="maria@cooperativa.com" required
							class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-digna-primary">
					</div>
					
					<div>
						<label class="block text-sm font-medium text-gray-700 mb-1">Telefone</label>
						<input type="tel" name="phone" value="(11) 99999-9999"
							class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-digna-primary">
					</div>
					
					<div>
						<label class="block text-sm font-medium text-gray-700 mb-1">Papel *</label>
						<select name="role" required
							class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-digna-primary">
							<option value="COORDINATOR" selected>Coordenador(a)</option>
							<option value="MEMBER">Sócio(a)</option>
							<option value="ADVISOR">Conselheiro(a)</option>
						</select>
					</div>
					
					<div>
						<label class="block text-sm font-medium text-gray-700 mb-1">Habilidades (separadas por vírgula)</label>
						<input type="text" name="skills" value="Gestão, Vendas, Produção"
							class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-digna-primary">
					</div>
				</div>
				
				<div class="mt-6 flex justify-end space-x-3">
					<button type="button" onclick="this.closest('.fixed').remove()"
						class="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50">
						Cancelar
					</button>
					<button type="submit"
						class="px-6 py-2 bg-digna-primary text-white rounded-lg font-medium hover:bg-blue-700">
						Salvar
					</button>
				</div>
			</form>
		</div>
	</div>
	`
	fmt.Fprint(w, html)
}

// updateMember atualiza os dados de um membro
func (h *MemberHandler) updateMember(w http.ResponseWriter, r *http.Request, memberID string) {
	// TODO: Implementar atualização real
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("HX-Trigger", "member-updated")

	html := `
	<div class="bg-green-100 border-l-4 border-green-500 text-green-700 p-4 mb-4" role="alert">
		<p>Sócio atualizado com sucesso!</p>
	</div>
	`
	fmt.Fprint(w, html)
}

func renderHTMXError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusBadRequest)

	html := fmt.Sprintf(`
	<div class="bg-red-100 border-l-4 border-red-500 text-red-700 p-4 mb-4" role="alert">
		<p class="font-bold">Erro</p>
		<p>%s</p>
	</div>
	`, message)

	fmt.Fprint(w, html)
}
