# 🗂️ REFERÊNCIA RÁPIDA DE MÓDULOS - Projeto Digna

**Última atualização:** 11/03/2026  
**Baseado em:** Análise da sessão 20260311

---

## 📊 VISÃO GERAL DOS MÓDULOS

| Módulo | Estado | Funcionalidade Principal | Dependências |
|--------|--------|--------------------------|--------------|
| **core_lume** | ✅ PRODUCTION | Domínio central (Membros, Decisões, Ledger) | - |
| **legal_facade** | ✅ 80% COMPLETA | Geração de documentos, formalização CADSOL | core_lume |
| **ui_web** | ✅ PRODUCTION | Interface web (HTMX + Tailwind) | core_lume, lifecycle |
| **lifecycle** | ✅ PRODUCTION | LifecycleManager, isolamento SQLite | - |
| **accountant_dashboard** | ✅ PRODUCTION | Dashboard contábil, SPED | core_lume |
| **budget** | ✅ PRODUCTION | Orçamento, planejamento | core_lume |
| **cash_flow** | ✅ PRODUCTION | Fluxo de caixa | core_lume |
| **supply** | ✅ PRODUCTION | Compras, estoque, fornecedores | core_lume |
| **pdv_ui** | ✅ PRODUCTION | Ponto de Venda | core_lume |
| **distribution** | ✅ PRODUCTION | Distribuição de sobras | core_lume |
| **integrations** | ✅ PRODUCTION | Integrações externas | core_lume |
| **sync_engine** | ✅ PRODUCTION | Sincronização delta | core_lume |

---

## 🏗️ MÓDULO: `core_lume` (Domínio Central)

### **Estrutura:**
```
core_lume/
├── internal/domain/entities.go          # ✅ Decision, Member, LedgerEntry
├── internal/repository/interfaces.go    # ✅ DecisionRepository, MemberRepository
├── internal/repository/sqlite.go        # ✅ Implementações SQLite
├── internal/service/decision_service.go # ✅ Lógica de negócio + SHA256
└── pkg/                                 # APIs públicas
```

### **Funcionalidades Implementadas:**
- **Decision management**: CRUD completo de decisões de assembleia
- **SHA256 hashing**: `generateHash()` em `decision_service.go`
- **Repository pattern**: Interfaces para isolamento de banco
- **Domain entities**: Structs puras sem dependências externas

### **Padrões Específicos:**
```go
// SHA256 Pattern (já implementado)
func generateHash(content string, entityID string) string {
    data := content + "|" + entityID + "|" + time.Now().Format(time.RFC3339)
    hash := sha256.Sum256([]byte(data))
    return hex.EncodeToString(hash[:])
}
```

---

## ⚖️ MÓDULO: `legal_facade` (Formalização)

### **Estrutura:**
```
legal_facade/
├── internal/document/generator.go       # ✅ Gera atas (precisa extender)
├── internal/document/formalization.go   # ✅ FormalizationSimulator
├── internal/document/statute.go         # ✅ StatuteGenerator + SHA256
├── internal/document/identity.go        # ✅ IdentityGenerator
├── internal/document/legal_repository.go # ✅ Interface LegalRepository
└── pkg/document/document.go             # ✅ API pública
```

### **Funcionalidades Implementadas:**
- **FormalizationSimulator**: `MinDecisionsForFormalization = 3` (já tem!)
- **StatuteGenerator**: Geração de estatuto com SHA256
- **LegalRepository**: Interface para acessar decisões
- **Generator**: Geração básica de atas (precisa extender para dossiê)

### **Código Reutilizável:**
```go
// Já existe em formalization.go
const MinDecisionsForFormalization = 3

func (s *FormalizationSimulator) CheckFormalizationCriteria(entityID string) (bool, int, error) {
    count, err := s.repo.GetDecisionCount(entityID)
    return count >= MinDecisionsForFormalization, count, err
}
```

---

## 🎨 MÓDULO: `ui_web` (Interface Web)

### **Estrutura:**
```
ui_web/
├── internal/handler/
│   ├── base_handler.go                  # ✅ BaseHandler com TemplateManager
│   ├── member_handler.go                # ✅ Padrão CRUD HTMX
│   ├── accountant_handler.go            # ✅ File download pattern
│   ├── cash_handler.go                  # ✅ API + JavaScript pattern
│   └── [auth, budget, dashboard, pdv, supply]_handler.go
├── templates/
│   ├── dashboard_simple.html            # ✅ Template base (copiar header)
│   ├── member_simple.html               # ✅ Padrão HTMX
│   └── *_simple.html                    # ✅ Cache-proof templates
└── main.go                              # ✅ Registro de handlers
```

### **Padrões Implementados:**

#### **1. File Download (accountant_handler.go):**
```go
w.Header().Set("Content-Type", "text/csv")
w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=fiscal_%s_%s.csv", entityID, period))
w.Header().Set("X-Export-Hash", batch.ExportHash)
w.Write(data)
```

#### **2. BaseHandler Template Functions:**
```go
// Já disponíveis:
"formatCurrency", "formatDate", "divide", "multiply",
"getAlertStatusLabel", "getAlertStatusClass", "getCategoryLabel"
```

#### **3. Estrutura de Handler:**
```go
type MemberHandler struct {
    *BaseHandler
    lifecycleManager lifecycle.LifecycleManager
    tmpl             *template.Template
}

func NewMemberHandler(lm lifecycle.LifecycleManager) (*MemberHandler, error) {
    base := NewBaseHandler(lm, true)
    // Adicionar funções específicas
    base.templateManager.AddFunc("getRoleLabel", ...)
    return &MemberHandler{BaseHandler: base, lifecycleManager: lm}, nil
}
```

---

## 🔄 MÓDULO: `lifecycle` (Gerenciamento de Banco)

### **Funcionalidade:**
- **LifecycleManager**: Ponto único de acesso a bancos SQLite
- **Isolamento**: `data/entities/{entity_id}.db` (um banco por entidade)
- **Migrações**: Versionadas e idempotentes
- **Graceful shutdown**: Fecha conexões corretamente

### **Uso em Handlers:**
```go
// Extrair entity_id do contexto (já feito pelo middleware)
entityID := r.URL.Query().Get("entity_id")

// Acessar banco via LifecycleManager
db, err := h.lifecycleManager.GetDatabase(entityID)
```

---

## 🛠️ PADRÕES DE IMPLEMENTAÇÃO POR TIPO DE FEATURE

### **1. CRUD Simples (ex: MemberHandler)**
```
1. Handler estende BaseHandler
2. Template: {feature}_simple.html
3. Rotas: GET /{feature}, POST /{feature}/create
4. HTMX para updates parciais
```

### **2. File Download (ex: AccountantHandler)**
```
1. Handler estende BaseHandler  
2. Rota específica para download
3. Headers: Content-Type, Content-Disposition
4. Opcional: X-Export-Hash para integridade
```

### **3. Document Generation (ex: Legal Dossier)**
```
1. Backend: Extender generator.go existente
2. Usar FormalizationSimulator para validação
3. Implementar SHA256 pattern já existente
4. Frontend: File download pattern
```

### **4. Dashboard Complexo (ex: CashHandler)**
```
1. API routes: /api/{feature}/data
2. JavaScript fetch para updates
3. JSON responses
4. Charts/visualizações
```

---

## 🔍 COMANDOS DE ANÁLISE POR MÓDULO

### **Para `legal_facade`:**
```bash
# Verificar estrutura existente
find modules/legal_facade -name "*.go" -type f

# Analisar FormalizationSimulator
grep -n "MinDecisionsForFormalization" modules/legal_facade/internal/document/formalization.go

# Verificar SHA256 implementations
grep -n "sha256.Sum256" modules/legal_facade/internal/document/
```

### **Para `core_lume`:**
```bash
# Verificar DecisionRepository
grep -n "DecisionRepository" modules/core_lume/internal/repository/interfaces.go

# Analisar SHA256 pattern
grep -n "generateHash" modules/core_lume/internal/service/decision_service.go

# Verificar domain entities
head -100 modules/core_lume/internal/domain/entities.go
```

### **Para `ui_web`:**
```bash
# Analisar handler similar
./scripts/tools/analyze_patterns.sh member --all

# Verificar file download pattern
grep -n "Content-Disposition" modules/ui_web/internal/handler/accountant_handler.go

# Analisar template base
head -100 modules/ui_web/templates/dashboard_simple.html
```

---

## ✅ CHECKLIST ANTES DE IMPLEMENTAR

### **Se feature envolve `legal_facade`:**
- [ ] Verificar se `generator.go` já tem funcionalidade similar
- [ ] Usar `FormalizationSimulator.CheckFormalizationCriteria()`
- [ ] Seguir padrão SHA256 já implementado
- [ ] Extender, não recriar

### **Se feature envolve file download:**
- [ ] Copiar padrão de `accountant_handler.go`
- [ ] Incluir headers: Content-Type, Content-Disposition
- [ ] Considerar hash para integridade (X-Export-Hash)

### **Se feature é CRUD simples:**
- [ ] Analisar `member_handler.go` como referência
- [ ] Usar HTMX para interatividade
- [ ] Template: `{feature}_simple.html`
- [ ] Estender `BaseHandler`

### **Se feature precisa de SHA256:**
- [ ] Usar padrão já implementado: `sha256.Sum256([]byte(data))`
- [ ] Retornar `hex.EncodeToString(hash[:])`
- [ ] Incluir timestamp/entityID no hash para unicidade

---

## 🚀 FLUXO DE IMPLEMENTAÇÃO OTIMIZADO

1. **Descobrir:** Usar comandos acima para verificar o que já existe
2. **Reutilizar:** Extender código existente em vez de criar novo
3. **Seguir padrões:** Copiar padrões já estabelecidos
4. **Integrar:** Usar interfaces e serviços já disponíveis
5. **Validar:** Smoke tests + testes existentes

---

## 📈 ESTIMATIVA DE ESFORÇO POR TIPO

| Tipo de Feature | Esforço (sem doc) | Esforço (com esta doc) | Redução |
|-----------------|-------------------|------------------------|---------|
| CRUD Simples | 2-3 horas | 1 hora | 50-60% |
| File Download | 3-4 horas | 1.5 horas | 50-60% |
| Document Generation | 4-6 horas | 2 horas | 60-70% |
| Dashboard Complexo | 5-8 horas | 3 horas | 40-60% |

**Economia estimada:** 50-70% do tempo de descoberta

---

**📌 MANTER ATUALIZADO:**  
Adicionar novos padrões e módulos conforme descobertas em sessões futuras.