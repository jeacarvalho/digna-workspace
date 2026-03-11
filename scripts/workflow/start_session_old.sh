#!/bin/bash
# start_session.sh - Inicializa uma sessão no projeto Digna
# Uso: ./start_session.sh [opcional: "quick" para modo rápido]

set -e  # Sai no primeiro erro

echo "🚀 Iniciando sessão no projeto Digna..."
echo "=========================================="

# Configurações
SESSION_ID=$(date +%Y%m%d_%H%M%S)
SESSION_FILE=".session_${SESSION_ID}"
QUICK_MODE="${1:-no}"

# 1. Criar arquivo de sessão
echo "SESSION_ID=${SESSION_ID}" > ${SESSION_FILE}
echo "START_TIME=$(date +%s)" >> ${SESSION_FILE}
echo "QUICK_MODE=${QUICK_MODE}" >> ${SESSION_FILE}

# 2. Atualizar contexto (modo rápido ou completo)
if [ "$QUICK_MODE" = "quick" ]; then
    echo "⚡ Modo rápido: verificando status básico..."
    
    # Verificação mínima
    if [ -f "./scripts/update_context.sh" ]; then
        ./scripts/update_context.sh 2>/dev/null || true
    fi
    
    # Status dos testes
    echo "📊 Verificando status dos testes..."
    if [ -d "modules" ]; then
        cd modules
        TEST_RESULT=$(./run_tests.sh 2>&1 | tail -20)
        echo "$TEST_RESULT" | grep -E "(PASS|FAIL|ok|^---)" | head -10
        cd ..
    fi
    
else
    echo "🔍 Modo completo: atualizando contexto completo..."
    
    # Atualizar contexto se o script existir
    if [ -f "./scripts/update_context.sh" ]; then
        echo "🔄 Executando update_context.sh..."
        ./scripts/update_context.sh
    else
        echo "⚠️  Script update_context.sh não encontrado, criando contexto básico..."
        
        # Criar QUICK_REFERENCE se não existir
        if [ ! -f "docs/QUICK_REFERENCE.md" ]; then
            mkdir -p docs
            cat > docs/QUICK_REFERENCE.md << 'EOF'
# 🚀 Quick Reference - Projeto Digna

**Última atualização:** $(date +%d/%m/%Y)
**Sessão iniciada:** ${SESSION_ID}

## 📋 Status Rápido
- Projeto: Digna (Economia Solidária)
- Arquitetura: Clean Architecture + DDD
- Banco: SQLite isolado por entidade
- Frontend: HTMX + Tailwind

## 🛠️ Comandos Úteis
- Iniciar sessão: `./start_session.sh`
- Processar tarefa: `./process_task.sh "descrição da tarefa"`
- Concluir tarefa: `./conclude_task.sh "aprendizados"`

## 📁 Estrutura Principal
```
modules/
├── core_lume/     # Domínio e serviços
├── ui_web/        # Interface web (HTMX/Tailwind)
├── lifecycle/     # Gerenciamento de banco
└── [outros...]
```
EOF
        fi
    fi
    
    # Status detalhado
    echo "📊 Status detalhado do projeto..."
    
    # Contar testes
    TEST_COUNT=$(find . -name "*test*.go" -type f | wc -l)
    echo "📈 Testes encontrados: ${TEST_COUNT}"
    
    # Verificar handlers
    HANDLER_COUNT=$(find modules/ui_web/internal/handler -name "*.go" -type f 2>/dev/null | wc -l || echo "0")
    echo "🎨 Handlers UI: ${HANDLER_COUNT}"
    
    # Verificar templates
    TEMPLATE_COUNT=$(find modules/ui_web/templates -name "*.html" -type f 2>/dev/null | wc -l || echo "0")
    echo "📝 Templates: ${TEMPLATE_COUNT}"
fi

# 3. Mostrar referência rápida
echo ""
echo "📚 REFERÊNCIA RÁPIDA:"
echo "===================="

if [ -f "docs/QUICK_REFERENCE.md" ]; then
    head -30 docs/QUICK_REFERENCE.md
else
    echo "ℹ️  Use './process_task.sh' para processar sua primeira tarefa."
    echo "📋 Formato recomendado:"
    echo "   ./process_task.sh \"Tipo: Feature | Módulo: ui_web | Objetivo: Implementar X | Decisões: seguir padrão Y\""
fi

# 4. Validação de Integridade do Projeto (NOVO)
echo ""
echo "🔍 VALIDAÇÃO DE INTEGRIDADE DO PROJETO:"
echo "======================================"

# Contar handlers registrados
HANDLER_COUNT=$(grep -c "RegisterRoutes" modules/ui_web/main.go 2>/dev/null || echo "0")
echo "Handlers registrados no main.go: $HANDLER_COUNT"

# Listar handlers
echo "Handlers encontrados:"
grep "New.*Handler" modules/ui_web/main.go | grep -v "//" | while read line; do
  handler=$(echo "$line" | sed 's/.*New\([A-Za-z]*\)Handler.*/\1Handler/' | head -1)
  if [ -n "$handler" ] && [ "$handler" != "line" ]; then
    echo "  ✅ $handler"
  fi
done 2>/dev/null

# Verificar templates vs handlers
echo ""
echo "Compatibilidade Templates-Handlers:"
TEMPLATE_MENTIONS=$(grep -r "ExecuteTemplate.*\.html" modules/ui_web/internal/handler/ 2>/dev/null | grep -v "test" | wc -l)
echo "  Templates referenciados em handlers: $TEMPLATE_MENTIONS"

# Verificar se smoke test script existe
if [ -f "./scripts/dev/smoke_test_new_feature.sh" ]; then
    echo "  ✅ Smoke test script disponível (scripts/dev/)"
else
    echo "  ⚠️  Smoke test script NÃO encontrado em scripts/dev/"
fi

# 5. Mostrar próximos passos sugeridos
echo ""
echo "🎯 PRÓXIMOS PASSOS SUGERIDOS:"
echo "============================="

if [ -f "docs/NEXT_STEPS.md" ]; then
    head -20 docs/NEXT_STEPS.md
else
    echo "1. Escolha uma tarefa do backlog"
    echo "2. Use './process_task.sh \"sua tarefa\"'"
    echo "3. Siga o checklist pré-implementação"
    echo "4. ✅ SEMPRE execute smoke test após implementação"
    echo "5. Documente aprendizados com './conclude_task.sh'"
fi

# 5. Criar contexto para o agente (opencode)
echo ""
echo "🤖 CRIANDO CONTEXTO PARA O AGENTE..."
echo "==================================="

# Criar arquivo de contexto do agente
cat > .agent_context.md << 'EOF'
# 🎯 CONTEXTO DO AGENTE - Projeto Digna

**Sessão iniciada:** $(date +%d/%m/%Y %H:%M:%S)
**Sessão ID:** ${SESSION_ID}
**Arquivo:** \`.agent_context.md\` (gerado automaticamente)

---

## 🚀 INSTRUÇÕES PARA O AGENTE (OPENCODE)

Você está trabalhando no **Projeto Digna** - sistema de economia solidária. Siga estas instruções:

### 1. CONTEXTO OBRIGATÓRIO (LER PRIMEIRO)
- **Leia:** \`docs/QUICK_REFERENCE.md\` - Arquitetura core, padrões, antipadrões
- **Verifique:** \`docs/NEXT_STEPS.md\` - Tarefas pendentes
- **Consulte:** \`docs/implementation_plans/\` - Planos existentes (se houver)

### 2. REGRAS SAGRADAS (NUNCA VIOLAR)
1. **ANTI-FLOAT:** Nunca use \`float\`, \`float32\`, \`float64\` para valores financeiros/tempo
2. **SOBERANIA DE DADOS:** Um banco SQLite por entidade, sem JOINs entre bancos
3. **CACHE-PROOF TEMPLATES:** Use \`*_simple.html\` com \`template.ParseFiles()\` NO HANDLER
4. **DESIGN SYSTEM:** Cores #2A5CAA (azul soberania), #4A7F3E (verde suor), #F57F17 (laranja)

### 3. PROCESSO DE TRABALHO COMPLETO
\`\`\`
1. ./start_session.sh                    → Inicia sessão, cria contexto
2. ./process_task.sh --file="arquivo.md" → Processa tarefa (--checklist, --plan, --execute)
3. Implementar                           → Você (opencode) implementa seguindo padrões
4. ./conclude_task.sh "aprendizados"     → Conclui tarefa individual
5. ./end_session.sh                      → Encerra sessão completa (após TODAS tarefas)
\`\`\`

**NOTA:** Use `--execute` apenas quando estiver pronto para o opencode implementar.

### 4. VALIDAÇÕES OBRIGATÓRIAS
- **Após implementação:** \`./scripts/dev/smoke_test_new_feature.sh\`
- **Testes:** \`cd modules && ./run_tests.sh\`
- **Handler no main.go:** Verificar se está registrado

### 5. ESTRUTURA DE ARQUIVOS
\`\`\`
modules/ui_web/internal/handler/    # HTTP handlers (HTMX)
modules/ui_web/templates/           # Templates *_simple.html
modules/core_lume/                  # Domínio e serviços
modules/lifecycle/                  # LifecycleManager (banco)
\`\`\`

### 6. PADRÕES DE CÓDIGO (ANALISAR ANTES DE IMPLEMENTAR)

#### 6.1 Estrutura de Handlers
- **BaseHandler:** Todos os handlers estendem `BaseHandler` (modules/ui_web/internal/handler/base_handler.go)
- **Construtor:** `New{Feature}Handler(lifecycle.LifecycleManager) (*{Feature}Handler, error)`
- **Estrutura típica:**
  ```go
  type MemberHandler struct {
      *BaseHandler
      lifecycleManager lifecycle.LifecycleManager
      tmpl             *template.Template
  }
  ```
- **Funções de template:** Adicionadas no construtor via `base.templateManager.AddFunc()`

#### 6.2 Funções de Template Disponíveis (BaseHandler)
```go
"formatCurrency": func(amount int64) string  // R$ 1.50 (int64 centavos)
"formatDate": func(t interface{}) string     // formatação data
"divide": func(a, b int64) float64           // divisão segura (b != 0)
"multiply": func(a, b int64) int64           // multiplicação
"getAlertStatusLabel": func(status string) string
"getAlertStatusClass": func(status string) string
"getCategoryLabel": func(category string) string
"fdiv": func(a, b float64) float64           // divisão float (apenas UI)
```

#### 6.3 Padrões de Interatividade (2 abordagens)

##### Padrão HTMX (Recomendado para CRUD simples)
- **GET:** Renderização inicial da página
- **POST:** Ações de formulário (criação, atualização)
- **Target/swap:** `hx-target="#result-area" hx-swap="outerHTML"`
- **Feedback:** Status via HTMX swaps com classes Tailwind
- **Exemplos:** MemberHandler, SupplyHandler

##### Padrão API + JavaScript (Para dashboards/complexos)
- **Rotas API:** `/api/{feature}/action` (ex: `/api/cash/balance`)
- **JavaScript:** Fetch API com `async/await`
- **JSON responses:** APIs retornam JSON para atualização parcial
- **Exemplos:** CashHandler (fluxo de caixa com gráficos)

#### 6.4 Templates HTML
- **Nome:** `*_simple.html` (documento HTML completo)
- **Estrutura:**
  1. Head com Tailwind config e paleta Digna
  2. Header com logo e navegação (copiar de `dashboard_simple.html`)
  3. Main content com HTMX forms
  4. Scripts HTMX no final
- **Paleta de cores:**
  - `#2A5CAA` (azul soberania) - botões principais
  - `#4A7F3E` (verde suor) - indicadores sucesso
  - `#F57F17` (laranja energia) - alertas/destaques
  - `#F9F9F6` (fundo), `#212121` (texto)

#### 6.5 Rotas Padrão
- **Listagem:** `GET /{feature}?entity_id={id}`
- **Criação:** `POST /{feature}?entity_id={id}`
- **Ação específica:** `POST /{feature}/{id}/toggle-status?entity_id={id}`

#### 6.6 Testes
- **Cobertura:** >90% para handlers
- **Estrutura:** `Test{Feature}Handler_{Action}`
- **Mock:** Mock de `lifecycle.LifecycleManager` quando necessário
- **Setup:** Usar `httptest.NewRecorder()` para testar handlers

#### 6.7 Nomenclatura
- **Handlers:** `{feature}_handler.go` (ex: `member_handler.go`)
- **Templates:** `{feature}_simple.html` (ex: `member_simple.html`)
- **Testes:** `{feature}_handler_test.go` (ex: `member_handler_test.go`)
- **Variáveis:** camelCase, tipos exportados em PascalCase

#### 6.8 Anti-Padrões (NUNCA FAZER)
1. ❌ Usar `float` para valores financeiros/tempo (usar `int64` centavos/minutos)
2. ❌ JOINs entre bancos SQLite diferentes
3. ❌ `template.ParseGlob()` ou variáveis globais de template
4. ❌ Acessar `core_lume/internal` diretamente (usar interfaces/mocks)
5. ❌ Hardcode entity_id (extrair de `r.Context().Value("entity_id")`)

### 7. COMO ANALISAR PADRÕES EXISTENTES

Antes de implementar, SEMPRE analise handlers e templates similares:

#### 7.1 Análise de Handler de Referência
```bash
# 1. Encontrar handler similar
ls modules/ui_web/internal/handler/*.go | grep -i [padrão]

# 2. Analisar estrutura
cat modules/ui_web/internal/handler/member_handler.go | head -100

# 3. Verificar funções de template
grep -n "AddFunc" modules/ui_web/internal/handler/member_handler.go

# 4. Analisar rotas
grep -n "RegisterRoutes\|HandleFunc" modules/ui_web/internal/handler/member_handler.go
```

#### 7.2 Análise de Template de Referência
```bash
# 1. Encontrar template similar
ls modules/ui_web/templates/*_simple.html | head -5

# 2. Analisar estrutura HTML
head -150 modules/ui_web/templates/member_simple.html

# 3. Verificar uso de HTMX
grep -n "hx-\|HTMX" modules/ui_web/templates/member_simple.html

# 4. Analisar classes Tailwind
grep -n "bg-\|text-\|border-" modules/ui_web/templates/member_simple.html | head -20
```

#### 7.3 Análise de Testes de Referência
```bash
# 1. Encontrar testes similares
find modules/ui_web -name "*test*.go" -exec grep -l "Test.*Handler" {} \;

# 2. Analisar setup de testes
grep -n "TestMain\|setup\|teardown" modules/ui_web/internal/handler/member_handler_test.go

# 3. Verificar mocks
grep -n "mock\|Mock" modules/ui_web/internal/handler/member_handler_test.go
```

### 8. EXEMPLO REAL: ANÁLISE DO MEMBERHANDLER

#### 8.1 Comandos de Análise:
```bash
# Analisar handler member
./scripts/tools/analyze_patterns.sh member

# Analisar template members (plural para template)
./scripts/tools/analyze_patterns.sh members --html

# Analisar tudo
./scripts/tools/analyze_patterns.sh member --all
```

#### 8.2 Padrões Extraídos do MemberHandler:
```go
// ESTRUTURA DO HANDLER
type MemberHandler struct {
    *BaseHandler                     // Estende BaseHandler
    lifecycleManager lifecycle.LifecycleManager
    tmpl             *template.Template
}

// CONSTRUTOR (padrão)
func NewMemberHandler(lm lifecycle.LifecycleManager) (*MemberHandler, error) {
    base := NewBaseHandler(lm, true)
    
    // Adicionar funções de template específicas
    base.templateManager.AddFunc("getRoleLabel", func(role MemberRole) string {
        switch role {
        case RoleCoordinator: return "Coordenador(a)"
        case RoleMember: return "Sócio(a)"
        case RoleAdvisor: return "Conselheiro(a)"
        default: return string(role)
        }
    })
    
    return &MemberHandler{
        BaseHandler:      base,
        lifecycleManager: lm,
    }, nil
}

// ROTAS (padrão HTMX)
func (h *MemberHandler) RegisterRoutes(mux *http.ServeMux) {
    mux.HandleFunc("/members", h.MembersPage)           // GET - renderização
    mux.HandleFunc("/members/create", h.CreateMember)   // POST - criação
    mux.HandleFunc("/members/", h.HandleMemberActions)  // /members/{id}/...
}

// TEMPLATE (members_simple.html - padrões)
// 1. Header copiado de dashboard_simple.html
// 2. Form HTMX: <form hx-post="/members/create" hx-target="#members-list">
// 3. Botões ação: <button hx-post="/members/{id}/toggle-status">
// 4. Classes: digna-card, bg-digna-primary, text-digna-text
```

#### 8.3 Padrões de Testes (member_handler_test.go):
```go
func TestMemberHandler_MembersPage(t *testing.T) {
    // Setup: mock lifecycleManager
    // Test: GET /members?entity_id=test
    // Assert: template renderizado corretamente
}

func TestMemberHandler_CreateMember(t *testing.T) {
    // Setup: mock + test data
    // Test: POST /members/create com form data
    // Assert: member criado, HTMX response correto
}
```

### 9. ARQUIVOS DE SESSÃO ATUAL
- **Sessão:** \`.session_*\` (timestamp)
- **Tarefa ativa:** \`.task_*\` (se existir)
- **Prompt atual:** \`.opencode_task_*\` (se process_task.sh --execute foi usado)

---

## 📋 STATUS ATUAL DO PROJETO

**Última atualização:** $(date +%d/%m/%Y)
**Testes:** ${TEST_COUNT} encontrados
**Handlers:** ${HANDLER_COUNT} UI handlers
**Templates:** ${TEMPLATE_COUNT} templates HTML

**Backlog:** Ver \`docs/NEXT_STEPS.md\`
**Aprendizados recentes:** Ver \`docs/learnings/\`

---

## 🎯 FLUXO COMPLETO DE TRABALHO

### 📋 SEQUÊNCIA CORRETA DE COMANDOS:
\`\`\`bash
# 1. INICIAR SESSÃO (uma vez por sessão)
./start_session.sh

# 2. PROCESSAR TAREFA (para cada tarefa)
./process_task.sh --file="prompt_tarefa.md" --checklist    # Primeiro: checklist
./process_task.sh --file="prompt_tarefa.md" --plan         # Depois: plano
./process_task.sh --file="prompt_tarefa.md" --execute      # Final: executar

# 3. IMPLEMENTAR (opencode faz isso)
#    - Seguir padrões identificados
#    - Usar analyze_patterns.sh para referências
#    - Validar com smoke tests

# 4. CONCLUIR TAREFA (após cada implementação)
./conclude_task.sh "Aprendizados: item1, item2" --success

# 5. ENCERRAR SESSÃO (após TODAS tarefas da sessão)
./end_session.sh
\`\`\`

### 🔄 CICLO POR TAREFA:
1. **Análise:** \`analyze_patterns.sh\` + checklist
2. **Planejamento:** Plano de implementação  
3. **Execução:** Implementação seguindo padrões
4. **Validação:** Smoke tests + testes unitários
5. **Documentação:** \`conclude_task.sh\`

### ⚠️ IMPORTANTE:
- **start_session.sh** → Uma vez por sessão
- **process_task.sh** → Para cada tarefa (pode rodar múltiplas vezes)
- **conclude_task.sh** → Após cada tarefa concluída
- **end_session.sh** → Apenas ao final de TODAS tarefas

**Referência rápida de comandos:**
\`\`\`bash
# 📋 FLUXO COMPLETO:
# 1. Iniciar sessão (uma vez)
./start_session.sh

# 2. Analisar padrões (antes de cada tarefa)
./scripts/tools/analyze_patterns.sh member --all
./scripts/tools/analyze_patterns.sh [handler] [--html|--tests|--all]

# 3. Processar tarefa (para cada tarefa)
./process_task.sh --file="prompt.md" --checklist    # Primeiro
./process_task.sh --file="prompt.md" --plan         # Depois  
./process_task.sh --file="prompt.md" --execute      # Executar

# 4. Validar implementação
./scripts/dev/smoke_test_new_feature.sh "Descrição" "/rota"
./scripts/dev/validate_e2e.sh --basic --headless

# 5. Concluir tarefa (após cada tarefa)
./conclude_task.sh "Aprendizados: item1, item2" --success

# 6. Encerrar sessão (após TODAS tarefas)
./end_session.sh
\`\`\`

---

**Este arquivo é atualizado automaticamente por \`start_session.sh\`.**
**O agente DEVE consultá-lo no início de cada interação.**
EOF

echo "✅ Contexto do agente criado: .agent_context.md"
echo "   ℹ️  O opencode deve ler este arquivo primeiro"

# 6. Informações da sessão
echo ""
echo "📋 INFORMAÇÕES DA SESSÃO:"
echo "========================="
echo "ID: ${SESSION_ID}"
echo "Data: $(date)"
echo "Modo: ${QUICK_MODE}"
echo "Arquivo de sessão: ${SESSION_FILE}"
echo "Contexto do agente: .agent_context.md"
echo ""
echo "✅ Sessão iniciada com sucesso!"
echo ""
echo "💡 Instrução para opencode: LEIA .agent_context.md PRIMEIRO"
echo "💡 Depois: './process_task.sh \"sua tarefa\" --execute'"

# Tornar scripts executáveis se não forem
chmod +x process_task.sh 2>/dev/null || true
chmod +x conclude_task.sh 2>/dev/null || true

exit 0