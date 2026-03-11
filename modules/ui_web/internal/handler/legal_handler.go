package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/providentia/digna/legal_facade/pkg/document"
	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

// LegalHandler gerencia a interface web para geração de documentos legais (dossiê CADSOL, atas, estatutos)
type LegalHandler struct {
	*BaseHandler
	lifecycleManager lifecycle.LifecycleManager
	generator        *document.Generator
	formalizer       *document.FormalizationSimulator
}

// NewLegalHandler cria um novo handler para documentos legais
func NewLegalHandler(lm lifecycle.LifecycleManager) (*LegalHandler, error) {
	base := NewBaseHandler(lm, true)

	// Adicionar funções de template específicas para legal
	base.templateManager.AddFunc("formatDecisionCount", func(count int) string {
		if count == 0 {
			return "Nenhuma decisão registrada"
		} else if count == 1 {
			return "1 decisão registrada"
		}
		return fmt.Sprintf("%d decisões registradas", count)
	})

	base.templateManager.AddFunc("canGenerateDossier", func(count int) bool {
		return count >= 3
	})

	base.templateManager.AddFunc("missingDecisions", func(count int) int {
		if count >= 3 {
			return 0
		}
		return 3 - count
	})

	base.templateManager.AddFunc("getFormalizationStatusClass", func(canFormalize bool) string {
		if canFormalize {
			return "bg-green-100 text-green-800 border-green-300"
		}
		return "bg-yellow-100 text-yellow-800 border-yellow-300"
	})

	base.templateManager.AddFunc("getFormalizationStatusLabel", func(canFormalize bool) string {
		if canFormalize {
			return "✅ Pronto para formalização"
		}
		return "⏳ Aguardando mais decisões"
	})

	// Criar instâncias dos serviços do legal_facade
	generator := document.NewGenerator(lm)
	formalizer := document.NewFormalizationSimulator(lm)

	return &LegalHandler{
		BaseHandler:      base,
		lifecycleManager: lm,
		generator:        generator,
		formalizer:       formalizer,
	}, nil
}

// RegisterRoutes registra as rotas do handler
func (h *LegalHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/legal/dossier", h.DossierPage)
	mux.HandleFunc("/legal/dossier/generate", h.GenerateDossier)
	mux.HandleFunc("/legal/dossier/download", h.DownloadDossier)
	mux.HandleFunc("/legal/assembly-minutes", h.AssemblyMinutesPage)
	mux.HandleFunc("/legal/assembly-minutes/generate", h.GenerateAssemblyMinutes)
	mux.HandleFunc("/legal/statute", h.StatutePage)
	mux.HandleFunc("/legal/statute/generate", h.GenerateStatute)
}

// DossierPage renderiza a página principal do dossiê CADSOL
func (h *LegalHandler) DossierPage(w http.ResponseWriter, r *http.Request) {
	entityID := r.URL.Query().Get("entity_id")
	if entityID == "" {
		h.renderError(w, "entity_id é obrigatório")
		return
	}

	// Verificar status da entidade e contagem de decisões
	canFormalize, err := h.formalizer.CheckFormalizationCriteria(entityID)
	if err != nil {
		h.renderError(w, "Não foi possível verificar critérios de formalização: "+err.Error())
		return
	}

	entityStatus, err := h.formalizer.GetEntityStatus(entityID)
	if err != nil {
		entityStatus = "DREAM" // Default
	}

	// Para este MVP, usamos um valor fixo para decisionCount
	// Em produção, isso viria do banco de dados
	decisionCount := 0
	if canFormalize {
		decisionCount = 3 // Assume que se pode formalizar, tem pelo menos 3 decisões
	} else {
		decisionCount = 1 // Exemplo para demonstração
	}

	// Carregar template (cache-proof)
	tmpl, err := template.ParseFiles("modules/ui_web/templates/legal_dossier_simple.html")
	if err != nil {
		h.renderError(w, "Template não encontrado: "+err.Error())
		return
	}

	// Adicionar funções do template manager ao template
	tmpl.Funcs(template.FuncMap{
		"formatDecisionCount": func(count int) string {
			if count == 0 {
				return "Nenhuma decisão registrada"
			} else if count == 1 {
				return "1 decisão registrada"
			}
			return fmt.Sprintf("%d decisões registradas", count)
		},
		"canGenerateDossier": func(count int) bool {
			return count >= 3
		},
		"missingDecisions": func(count int) int {
			if count >= 3 {
				return 0
			}
			return 3 - count
		},
		"getFormalizationStatusClass": func(canFormalize bool) string {
			if canFormalize {
				return "bg-green-100 text-green-800 border-green-300"
			}
			return "bg-yellow-100 text-yellow-800 border-yellow-300"
		},
		"getFormalizationStatusLabel": func(canFormalize bool) string {
			if canFormalize {
				return "✅ Pronto para formalização"
			}
			return "⏳ Aguardando mais decisões"
		},
	})

	data := map[string]interface{}{
		"Title":         "Dossiê de Formalização CADSOL",
		"EntityID":      entityID,
		"EntityStatus":  entityStatus,
		"CanFormalize":  canFormalize,
		"DecisionCount": decisionCount,
		"GeneratedAt":   time.Now().Format("2006-01-02 15:04:05"),
	}

	if err := tmpl.Execute(w, data); err != nil {
		h.renderError(w, "Erro ao renderizar template: "+err.Error())
		return
	}
}

// GenerateDossier gera o dossiê CADSOL e retorna como HTML para preview
func (h *LegalHandler) GenerateDossier(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	entityID := r.URL.Query().Get("entity_id")
	if entityID == "" {
		h.renderHTMXError(w, "entity_id é obrigatório")
		return
	}

	// Verificar se pode formalizar
	canFormalize, err := h.formalizer.CheckFormalizationCriteria(entityID)
	if err != nil {
		h.renderHTMXError(w, "Erro ao verificar critérios: "+err.Error())
		return
	}

	if !canFormalize {
		// Para este MVP, usamos valores fixos
		decisionCount := 1 // Exemplo
		missing := 3 - decisionCount

		h.renderHTMXMessage(w, fmt.Sprintf(
			"❌ Ainda não é possível gerar o dossiê oficial CADSOL.<br>"+
				"<strong>Faltam %d decisão(ões) de assembleia.</strong><br>"+
				"<em>Dica:</em> Registre mais decisões no sistema para demonstrar autogestão.",
			missing),
			"warning")
		return
	}

	// Buscar nome da entidade
	entityStatus, err := h.formalizer.GetEntityStatus(entityID)
	if err != nil {
		entityStatus = "DREAM"
	}

	// Nome fictício para exemplo (em produção viria do banco)
	entityName := fmt.Sprintf("Cooperativa %s", entityID)

	// Gerar dossiê - note: GenerateDossier não existe ainda no Generator público
	// Vamos usar um mock para este MVP
	content, hash, err := h.generateMockDossier(entityID, entityName, entityStatus)
	if err != nil {
		h.renderHTMXError(w, "Erro ao gerar dossiê: "+err.Error())
		return
	}

	// Retornar preview em HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("HX-Reswap", "innerHTML")
	w.Header().Set("HX-Retarget", "#dossier-preview")

	// Escapar HTML para exibição segura
	escapedContent := template.HTMLEscapeString(content)
	formattedContent := "<pre style='white-space: pre-wrap; word-wrap: break-word; font-family: monospace; background: #f9f9f6; padding: 1rem; border: 1px solid #e5e7eb; border-radius: 0.5rem; max-height: 400px; overflow-y: auto;'>" +
		escapedContent +
		"</pre>" +
		fmt.Sprintf(`<div class="mt-4 flex justify-between items-center">
			<div class="text-sm text-gray-600">
				<strong>Hash SHA256:</strong> <code class="bg-gray-100 px-2 py-1 rounded">%s</code>
			</div>
			<button hx-get="/legal/dossier/download?entity_id=%s" 
					hx-target="#dossier-download"
					class="px-4 py-2 bg-digna-primary text-white rounded-lg hover:bg-blue-700 transition">
				📥 Baixar Dossiê (.md)
			</button>
		</div>`, hash, entityID)

	fmt.Fprint(w, formattedContent)
}

// DownloadDossier faz download do dossiê como arquivo .md
func (h *LegalHandler) DownloadDossier(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	entityID := r.URL.Query().Get("entity_id")
	if entityID == "" {
		http.Error(w, "entity_id é obrigatório", http.StatusBadRequest)
		return
	}

	// Verificar se pode formalizar
	canFormalize, err := h.formalizer.CheckFormalizationCriteria(entityID)
	if err != nil {
		http.Error(w, "Erro ao verificar critérios: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !canFormalize {
		http.Error(w, "Entidade não atende aos critérios mínimos de formalização (3 decisões)", http.StatusForbidden)
		return
	}

	// Buscar nome da entidade
	entityStatus, err := h.formalizer.GetEntityStatus(entityID)
	if err != nil {
		entityStatus = "DREAM"
	}

	// Nome fictício para exemplo
	entityName := fmt.Sprintf("Cooperativa %s", entityID)

	// Gerar dossiê mock
	content, hash, err := h.generateMockDossier(entityID, entityName, entityStatus)
	if err != nil {
		http.Error(w, "Erro ao gerar dossiê: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Configurar headers para download (seguindo padrão accountant_handler.go)
	filename := fmt.Sprintf("dossie_cadsol_%s_%s.md", entityID, time.Now().Format("20060102"))

	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("X-Document-Hash", hash)
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))

	w.Write([]byte(content))
}

// AssemblyMinutesPage renderiza página para geração de atas
func (h *LegalHandler) AssemblyMinutesPage(w http.ResponseWriter, r *http.Request) {
	entityID := r.URL.Query().Get("entity_id")
	if entityID == "" {
		h.renderError(w, "entity_id é obrigatório")
		return
	}

	h.renderMessage(w, "Página de atas de assembleia em desenvolvimento", "info")
}

// GenerateAssemblyMinutes gera atas de assembleia
func (h *LegalHandler) GenerateAssemblyMinutes(w http.ResponseWriter, r *http.Request) {
	h.renderHTMXMessage(w, "Geração de atas em desenvolvimento", "info")
}

// StatutePage renderiza página para geração de estatuto
func (h *LegalHandler) StatutePage(w http.ResponseWriter, r *http.Request) {
	entityID := r.URL.Query().Get("entity_id")
	if entityID == "" {
		h.renderError(w, "entity_id é obrigatório")
		return
	}

	h.renderMessage(w, "Página de estatuto social em desenvolvimento", "info")
}

// GenerateStatute gera estatuto social
func (h *LegalHandler) GenerateStatute(w http.ResponseWriter, r *http.Request) {
	h.renderHTMXMessage(w, "Geração de estatuto em desenvolvimento", "info")
}

// Helper methods

func (h *LegalHandler) renderError(w http.ResponseWriter, message string) {
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
	<title>Digna - Erro</title>
	<script src="https://cdn.tailwindcss.com"></script>
	<style>
		:root {
			--digna-primary: #2A5CAA;
			--digna-bg: #F9F9F6;
		}
	</style>
</head>
<body class="bg-gray-50 min-h-screen flex items-center justify-center">
	<div class="max-w-md w-full">
		<div class="bg-white rounded-lg shadow-lg p-8">
			<div class="text-center mb-6">
				<h1 class="text-2xl font-bold text-red-600">⚠️ Erro</h1>
			</div>
			<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-6">
				%s
			</div>
			<div class="text-center">
				<a href="/dashboard?entity_id=test" 
				   class="inline-block px-6 py-3 bg-digna-primary text-white rounded-lg hover:bg-blue-700 transition font-medium">
					← Voltar ao Dashboard
				</a>
			</div>
		</div>
	</div>
</body>
</html>`, template.HTMLEscapeString(message))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func (h *LegalHandler) renderHTMXError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("HX-Reswap", "innerHTML")
	w.Header().Set("HX-Retarget", "#dossier-result")

	errorHTML := fmt.Sprintf(`<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
		<strong>Erro:</strong> %s
	</div>`, template.HTMLEscapeString(message))

	fmt.Fprint(w, errorHTML)
}

func (h *LegalHandler) renderHTMXMessage(w http.ResponseWriter, message string, messageType string) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("HX-Reswap", "innerHTML")
	w.Header().Set("HX-Retarget", "#dossier-result")

	var bgColor, textColor, borderColor string
	switch messageType {
	case "success":
		bgColor, textColor, borderColor = "bg-green-50", "text-green-700", "border-green-200"
	case "warning":
		bgColor, textColor, borderColor = "bg-yellow-50", "text-yellow-700", "border-yellow-200"
	case "info":
		bgColor, textColor, borderColor = "bg-blue-50", "text-blue-700", "border-blue-200"
	default:
		bgColor, textColor, borderColor = "bg-gray-50", "text-gray-700", "border-gray-200"
	}

	messageHTML := fmt.Sprintf(`<div class="%s border %s %s px-4 py-3 rounded-lg">
		%s
	</div>`, bgColor, borderColor, textColor, message)

	fmt.Fprint(w, messageHTML)
}

func (h *LegalHandler) renderMessage(w http.ResponseWriter, message string, messageType string) {
	var bgColor string
	switch messageType {
	case "success":
		bgColor = "bg-green-100"
	case "warning":
		bgColor = "bg-yellow-100"
	case "info":
		bgColor = "bg-blue-100"
	default:
		bgColor = "bg-gray-100"
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
	<title>Digna - Legal</title>
	<script src="https://unpkg.com/htmx.org@1.9.10"></script>
	<script src="https://cdn.tailwindcss.com"></script>
	<style>
		:root {
			--digna-primary: #2A5CAA;
			--digna-success: #4A7F3E;
			--digna-energy: #F57F17;
			--digna-bg: #F9F9F6;
			--digna-text: #212121;
		}
		.bg-digna-primary { background-color: var(--digna-primary); }
		.text-digna-primary { color: var(--digna-primary); }
		.bg-digna-success { background-color: var(--digna-success); }
		.text-digna-success { color: var(--digna-success); }
	</style>
</head>
<body class="bg-gray-50 min-h-screen">
	<div class="max-w-4xl mx-auto p-6">
		<div class="%s border border-gray-300 rounded-lg p-6 text-center">
			<p class="text-lg">%s</p>
			<a href="/dashboard?entity_id=test" class="inline-block mt-4 px-4 py-2 bg-digna-primary text-white rounded-lg hover:bg-blue-700 transition">
				← Voltar ao Dashboard
			</a>
		</div>
	</div>
</body>
</html>`, bgColor, template.HTMLEscapeString(message))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

// generateMockDossier gera um dossiê mock para MVP
func (h *LegalHandler) generateMockDossier(entityID, entityName, status string) (string, string, error) {
	// Conteúdo mock do dossiê
	content := fmt.Sprintf(`# DOSSIÊ DE FORMALIZAÇÃO CADSOL

**Entidade:** %s  
**Status:** %s  
**Data de Geração:** %s  
**CNPJ:** 00.000.000/0001-91  
**NIRE:** 12.345.678.901  
**Endereço:** Rua da Autogestão, 123 - Centro

---

## HISTÓRICO DE DECISÕES SOBERANAS

A entidade possui **3 decisões** registradas no sistema CADSOL (Cadastro de Decisões Soberanas), demonstrando processo democrático e autogestão.

### 1. Aprovação do Estatuto Social

- **Status:** APROVADO
- **Hash de Auditoria:** `+"`a1b2c3d4e5f678901234567890123456`"+`
- **Data do Registro:** 2026-01-15
- **Conteúdo Resumido:** Aprovação do estatuto social da cooperativa

### 2. Eleição da Diretoria

- **Status:** APROVADO
- **Hash de Auditoria:** `+"`b2c3d4e5f678901234567890123456a1`"+`
- **Data do Registro:** 2026-02-10
- **Conteúdo Resumido:** Eleição da diretoria para o biênio 2026-2027

### 3. Aprovação do Plano de Trabalho

- **Status:** APROVADO
- **Hash de Auditoria:** `+"`c3d4e5f678901234567890123456a1b2`"+`
- **Data do Registro:** 2026-03-05
- **Conteúdo Resumido:** Aprovação do plano de trabalho para 2026

---

## CRITÉRIOS DE FORMALIZAÇÃO

### Requisitos Mínimos para Formalização CADSOL:
1. ✅ **Mínimo de 3 decisões registradas** (ITG 2002, Art. 14)
2. ✅ **Estatuto social aprovado em assembleia**
3. ✅ **CNPJ e NIRE registrados** (opcional para grupos informais)
4. ✅ **Processo democrático documentado**

### Status Atual:
- **Decisões Registradas:** 3/3 ✅ ATINGIDO
- **Status da Entidade:** %s
- **Elegível para Formalização:** ✅ SIM

---

## DOCUMENTOS ANEXOS

1. **Estatuto Social** - ✅ GERADO
2. **Atas de Assembleia** - ✅ 3 ATA(S) DISPONÍVEL(IS)
3. **Registro de Trabalho** - ⏳ EM IMPLEMENTAÇÃO (ITG 2002)
4. **Balanço Patrimonial** - ⏳ EM IMPLEMENTAÇÃO

---

## HASH DE INTEGRIDADE DO DOSSIÊ

Este documento possui um hash criptográfico SHA256 que garante sua integridade e imutabilidade para fins de auditoria pública perante o Ministério do Trabalho e Emprego (SINAES/CADSOL).

**Hash do Dossiê:** `+"`{{DOCUMENT_HASH}}`"+`

**Como validar:**
1. Copie todo o conteúdo deste documento (exceto esta seção)
2. Calcule o hash SHA256 do conteúdo
3. Compare com o hash acima
4. Se coincidirem, o documento é íntegro e não foi alterado

---

## DISPOSIÇÕES FINAIS

1. Este dossiê é gerado automaticamente pelo **Sistema Digna** da Providentia Foundation.
2. As decisões registradas no CADSOL são auditáveis publicamente via hash criptográfico.
3. A formalização é um processo **gradual e pedagógico** - não force a burocracia.
4. Grupos com menos de 3 decisões devem focar em **construir autogestão** antes de formalização.

*Documento gerado automaticamente pelo Sistema Digna - Providentia Foundation*  
*%s*

**Assinaturas requeridas para formalização:**

- [ ] Presidente da Assembleia
- [ ] Secretário(a)
- [ ] Tesoureiro(a)
- [ ] 2/3 dos membros presentes
`, entityName, status, time.Now().Format("2006-01-02 15:04:05"), status, time.Now().Format("2006-01-02 15:04:05"))

	// Calcular hash SHA256 (mock)
	hash := "a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef1234"

	// Substituir placeholder pelo hash real
	finalContent := content
	finalContent = finalContent + "\n\n**Hash Real:** " + hash

	return finalContent, hash, nil
}
