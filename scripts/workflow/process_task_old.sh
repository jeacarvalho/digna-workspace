#!/bin/bash
# process_task.sh - Processa uma tarefa no projeto Digna
# Uso: ./process_task.sh "Descrição da tarefa" [--checklist] [--plan] [--execute]
#      ./process_task.sh --file arquivo.md [OPÇÕES]
# Ex: ./process_task.sh "Implementar UI para Fornecedores" --checklist
# Ex: ./process_task.sh --file Prompt_teste_correcos.md --execute

set -e

# Configurações
SESSION_FILE=$(ls .session_* 2>/dev/null | head -1)
TASK_ID=$(date +%Y%m%d_%H%M%S)
TASK_FILE=".task_${TASK_ID}"

# Variáveis de controle
TASK_DESCRIPTION=""
FILE_MODE=false
FILE_PATH=""
FILE_NAME=""

# Funções de ajuda
show_help() {
    cat << EOF
🎯 process_task.sh - Processador Inteligente de Tarefas Digna

USO:
  ./process_task.sh "DESCRIÇÃO DA TAREFA" [OPÇÕES]
  ./process_task.sh --file ARQUIVO.md [OPÇÕES]
  ./process_task.sh -f ARQUIVO.md [OPÇÕES]

Descrição da tarefa (FORMATO RECOMENDADO):
  "Tipo: Feature | Módulo: ui_web | Objetivo: Implementar X | Decisões: seguir padrão Y"

Opções:
  --checklist, -c    Apenas gerar checklist pré-implementação
  --plan, -p         Gerar plano de implementação completo
  --execute, -e      Executar implementação (interage com opencode)
  --file, -f         Ler descrição da tarefa de um arquivo
  --help, -h         Mostrar esta ajuda

Exemplos:
  # Descrição direta
  ./process_task.sh "Tipo: Feature | Módulo: ui_web | Objetivo: Implementar UI para Fornecedores" --execute
  
  # Ler de arquivo
  ./process_task.sh --file Prompt_teste_correcos.md --execute
  ./process_task.sh -f docs/tasks/bug_fix.md -e
  
  # Modos diferentes
  ./process_task.sh "Bug no PDV: erro ao adicionar produto" --checklist
  ./process_task.sh "Melhorar performance do dashboard" --plan

O script irá:
1. Analisar a descrição da tarefa (diretamente ou de arquivo)
2. Verificar contexto existente
3. Gerar checklists/planos conforme opção
4. Preparar para execução com opencode
EOF
    exit 0
}

# Processar argumentos
MODE="analyze"
NEXT_IS_FILE=false

for arg in "$@"; do
    if [ "$NEXT_IS_FILE" = true ]; then
        FILE_PATH="$arg"
        FILE_MODE=true
        NEXT_IS_FILE=false
        continue
    fi
    
    case $arg in
        --checklist|-c)
            MODE="checklist"
            ;;
        --plan|-p)
            MODE="plan"
            ;;
        --execute|-e)
            MODE="execute"
            ;;
        --file|-f)
            NEXT_IS_FILE=true
            ;;
        --help|-h)
            show_help
            ;;
        *)
            # Se não é opção e não estamos em modo arquivo, é a descrição
            if [ "$FILE_MODE" = false ] && [ -z "$TASK_DESCRIPTION" ]; then
                TASK_DESCRIPTION="$arg"
            fi
            ;;
    esac
done

# Verificar se temos descrição
if [ -z "$TASK_DESCRIPTION" ] && [ "$FILE_MODE" = false ]; then
    echo "❌ Erro: Nenhuma descrição de tarefa fornecida."
    echo ""
    show_help
fi

# Se modo arquivo, ler conteúdo
if [ "$FILE_MODE" = true ]; then
    if [ -z "$FILE_PATH" ]; then
        echo "❌ Erro: Modo --file ativado mas nenhum arquivo especificado."
        show_help
    fi
    
    if [ ! -f "$FILE_PATH" ]; then
        echo "❌ Erro: Arquivo não encontrado: $FILE_PATH"
        exit 1
    fi
    
    FILE_NAME="$(basename "$FILE_PATH")"
    TASK_DESCRIPTION="$(cat "$FILE_PATH")"
    echo "📄 Lendo tarefa do arquivo: $FILE_NAME"
fi

# Processar opções
MODE="analyze"
for arg in "$@"; do
    case $arg in
        --checklist)
            MODE="checklist"
            ;;
        --plan)
            MODE="plan"
            ;;
        --execute)
            MODE="execute"
            ;;
    esac
done

echo "🔍 Processando tarefa..."
echo "========================"

if [ "$FILE_MODE" = true ]; then
    echo "📄 Fonte: Arquivo '$FILE_NAME'"
    # Mostrar preview do conteúdo
    PREVIEW=$(echo "$TASK_DESCRIPTION" | head -10 | sed 's/^/   /')
    echo "📋 Preview (primeiras 10 linhas):"
    echo "$PREVIEW"
    if [ $(echo "$TASK_DESCRIPTION" | wc -l) -gt 10 ]; then
        echo "   ... ($(echo "$TASK_DESCRIPTION" | wc -l) linhas no total)"
    fi
else
    echo "📝 Fonte: Descrição direta"
    echo "📋 Conteúdo: ${TASK_DESCRIPTION:0:100}..."
fi

# 1. Extrair informações da descrição
echo ""
echo "📋 Extraindo informações da descrição..."

# Função para extrair metadados de arquivos formatados
extract_metadata() {
    local content="$1"
    
    # Tentar extrair do formato padronizado (com |)
    local type=$(echo "$content" | grep -oi "tipo:\s*[^|]*" | cut -d: -f2 | xargs)
    local module=$(echo "$content" | grep -oi "módulo:\s*[^|]*" | cut -d: -f2 | xargs)
    local objective=$(echo "$content" | grep -oi "objetivo:\s*[^|]*" | cut -d: -f2 | xargs)
    local decisions=$(echo "$content" | grep -oi "decisões:\s*[^|]*" | cut -d: -f2 | xargs)
    
    # Se não encontrou no formato com |, tentar formato simples
    if [ -z "$type" ]; then
        type=$(echo "$content" | grep -i "^tipo:" | cut -d: -f2- | xargs)
    fi
    if [ -z "$module" ]; then
        module=$(echo "$content" | grep -i "^módulo:" | cut -d: -f2- | xargs)
    fi
    if [ -z "$objective" ]; then
        objective=$(echo "$content" | grep -i "^objetivo:" | cut -d: -f2- | xargs)
    fi
    if [ -z "$decisions" ]; then
        decisions=$(echo "$content" | grep -i "^decisões:" | cut -d: -f2- | xargs)
    fi
    
    # Valores padrão
    type="${type:-Feature}"
    module="${module:-ui_web}"
    objective="${objective:-$TASK_DESCRIPTION}"
    decisions="${decisions:-Seguir padrões estabelecidos}"
    
    echo "$type|$module|$objective|$decisions"
}

# Extrair metadados
METADATA=$(extract_metadata "$TASK_DESCRIPTION")
TASK_TYPE=$(echo "$METADATA" | cut -d'|' -f1)
MODULE=$(echo "$METADATA" | cut -d'|' -f2)
OBJECTIVE=$(echo "$METADATA" | cut -d'|' -f3)
DECISIONS=$(echo "$METADATA" | cut -d'|' -f4)

# Criar nome da feature (para arquivos)
FEATURE_NAME=$(echo "$OBJECTIVE" | tr '[:upper:]' '[:lower:]' | sed 's/implementar //g; s/ui para //g; s/ //g; s/\.//g' | cut -c1-20)
FEATURE_NAME="${FEATURE_NAME:-task_${TASK_ID}}"

# 2. Salvar informações da tarefa
cat > ${TASK_FILE} << EOF
TASK_ID=${TASK_ID}
TASK_DESCRIPTION="${TASK_DESCRIPTION}"
TASK_TYPE="${TASK_TYPE}"
MODULE="${MODULE}"
OBJECTIVE="${OBJECTIVE}"
DECISIONS="${DECISIONS}"
FEATURE_NAME="${FEATURE_NAME}"
START_TIME=$(date +%s)
MODE="${MODE}"
FILE_MODE=${FILE_MODE}
FILE_PATH="${FILE_PATH}"
FILE_NAME="${FILE_NAME}"
EOF

echo "✅ Tarefa registrada: ${TASK_ID}"
echo "   Tipo: ${TASK_TYPE}"
echo "   Módulo: ${MODULE}"
echo "   Objetivo: ${OBJECTIVE}"
if [ "$FILE_MODE" = true ]; then
    echo "   Fonte: Arquivo '${FILE_NAME}'"
fi

# 3. Modo: Checklist pré-implementação
if [ "$MODE" = "checklist" ] || [ "$MODE" = "plan" ] || [ "$MODE" = "execute" ]; then
    echo ""
    echo "📝 GERANDO CHECKLIST PRÉ-IMPLEMENTAÇÃO..."
    echo "========================================"
    
    CHECKLIST_FILE="docs/implementation_plans/${FEATURE_NAME}_pre_check.md"
    mkdir -p docs/implementation_plans
    
    # Template de checklist
    CURRENT_DATE=$(date +%d/%m/%Y\ %H:%M:%S)
    cat > ${CHECKLIST_FILE} << EOF
# 🔍 Checklist de Validação Pré-Implementação: ${FEATURE_NAME}

**Tarefa:** ${TASK_DESCRIPTION}
**Gerado em:** ${CURRENT_DATE}
**Tarefa ID:** ${TASK_ID}

---

## 📋 Informações Extraídas
- **Tipo:** ${TASK_TYPE}
- **Módulo:** ${MODULE}
- **Objetivo:** ${OBJECTIVE}
- **Decisões:** ${DECISIONS}

---

## 1. 🏗️ Análise Arquitetural

### 1.1 Backend Existente
- [ ] **Serviço existe?** \`find modules/core_lume -name "*${FEATURE_NAME}*" -type f\`
- [ ] **Testes passando?** \`cd modules/core_lume && go test ./... -run [Ff]${FEATURE_NAME}\`
- [ ] **Cobertura adequada?** (>80% para serviços críticos)

### 1.2 Acessibilidade do Serviço
- [ ] **Pacote é \`internal\`?** Verificar caminho do serviço
- [ ] **Padrão de acesso estabelecido?** Como outros handlers acessam serviços similares
- [ ] **Dependências do Serviço:** Quais repositórios precisa?

---

## 2. 🎨 Padrões de Frontend

### 2.1 Handlers de Referência
- [ ] **Qual handler mais similar?** \`ls modules/ui_web/internal/handler/*.go | grep -i [padrão]\`
- [ ] **Estende BaseHandler?** Verificar estrutura do handler de referência
- [ ] **Rotas padrão HTMX?** Analisar \`RegisterRoutes\` do handler de referência

### 2.2 Sistema de Templates
- [ ] **Template base a usar?** \`ls modules/ui_web/templates/*_simple.html | head -5\`
- [ ] **Funções de template necessárias?** Analisar template similar
- [ ] **Já existem no BaseHandler?** \`grep -n "AddFunc" modules/ui_web/internal/handler/base_handler.go\`

### 2.3 Navegação e Layout
- [ ] **Quais templates atualizar?** \`grep -l "nav\\\\|Navegação" modules/ui_web/templates/*.html\`
- [ ] **Padrão de navegação?** Horizontal (header) vs Grid (layout.html)
- [ ] **Design system aplicado?** Cores (#2A5CAA, #4A7F3E, #F57F17)

---

## 3. ⚙️ Testabilidade

### 3.1 Estrutura de Testes
- [ ] **Testes de handler similares?** \`find modules/ui_web -name "*test*.go" -exec grep -l "Test.*Handler" {} \\;\`
- [ ] **Como mockar dependências?** Analisar testes de referência
- [ ] **Setup de testes necessário?** Precisa de dados, templates, etc.

### 3.2 Cobertura Alvo
- [ ] **Coverage mínimo:** >90% para handlers
- [ ] **Tipos de testes necessários:**
  - [ ] Testes unitários (lógica pura)
  - [ ] Testes de integração (com banco)
  - [ ] Testes E2E (Playwright - opcional)

---

## 4. 🚨 Riscos Identificados

| Risco | Probabilidade | Impacto | Mitigação |
|-------|--------------|---------|-----------|
| Serviço não acessível (internal) | Alta | Alto | Implementar com mock, planejar refatoração |
| Performance com muitos dados | Média | Médio | Paginação no template, lazy loading |
| Template functions faltando | Alta | Baixo | Adicionar ao BaseHandler ou handler específico |

---

## 5. 📝 Decisões Documentadas

### 5.1 Decisões Técnicas
1. **Acesso ao serviço:** □ Mock inicial □ API layer □ Direct import □ Outro: ________
2. **Estrutura do handler:** □ Estende BaseHandler □ Independente □ Outro: ________
3. **Template base:** □ dashboard_simple.html □ pdv_simple.html □ Outro: ________

### 5.2 Decisões de Design
1. **Navegação:** Links em quais templates? \`_________________________________\`
2. **Responsividade:** □ Mobile-first □ Desktop-first □ Adaptativo

---

## ✅ Checklist de Validação Final

### ANTES de começar a codificar:
- [ ] Backend analisado e compreendido
- [ ] Padrões de frontend identificados
- [ ] Riscos mapeados e mitigados
- [ ] Decisões documentadas
- [ ] Checklist completo preenchido

---

## 📊 Métricas desta Análise

- **Tempo gasto na análise:** ______ minutos
- **Problemas identificados antecipadamente:** ______
- **Confiança na estimativa:** 1-5 (1=baixa, 5=alta)

---

**📌 Nota:** Preencher este checklist antes de criar plano de implementação.
Arquivo: ${CHECKLIST_FILE}
EOF
    
    echo "✅ Checklist gerado: ${CHECKLIST_FILE}"
    echo "   ℹ️  Preencha antes de prosseguir."
fi

# 4. Modo: Plano de implementação
if [ "$MODE" = "plan" ] || [ "$MODE" = "execute" ]; then
    echo ""
    echo "📋 GERANDO PLANO DE IMPLEMENTAÇÃO..."
    echo "==================================="
    
    PLAN_DATE=$(date +%Y%m%d)
    PLAN_FILE="docs/implementation_plans/${FEATURE_NAME}_implementation_${PLAN_DATE}.md"
    
    # Verificar se checklist foi preenchido
    if [ ! -f "${CHECKLIST_FILE}" ]; then
        echo "⚠️  Checklist não encontrado. Gerando plano básico..."
    fi
    
    # Template de plano
    PLAN_DATE_FULL=$(date +%d/%m/%Y\ %H:%M:%S)
    cat > ${PLAN_FILE} << EOF
# 📋 Plano de Implementação: ${FEATURE_NAME}

**Feature:** ${FEATURE_NAME}
**Tarefa ID:** ${TASK_ID}
**Gerado em:** ${PLAN_DATE_FULL}
**Descrição:** ${TASK_DESCRIPTION}

**📌 PRÉ-REQUISITO:** Preencher \`${CHECKLIST_FILE}\` antes de implementar

---

## 0. 🔍 Fase de Descoberta (A COMPLETAR)

**Checklist:** \`${CHECKLIST_FILE}\`
**Status:** □ Não iniciado □ Em progresso ✅ Completo

### **0.1 Backend Status**
- [ ] Serviço existe e testado: □ Sim □ Parcial □ Não
- [ ] Acessível do UI: □ Sim (público) □ Mock necessário □ Internal
- [ ] Padrão de acesso: [API layer / Direct import / Mock inicial]

### **0.2 Padrões Identificados**
- Handler de referência: \`__________________________\`
- Template base: \`__________________________\`
- Rotas padrão: \`GET /______\`, \`POST /______\`, \`POST /______/{id}/______\`

### **0.3 Riscos Principais**
1. **Risco:** [Descrição breve] → **Mitigação:** [Ação]
2. **Risco:** [Descrição breve] → **Mitigação:** [Ação]

---

## 1. 🎯 Objetivo da Tarefa

${OBJECTIVE}

**Contexto:** [Completar baseado na análise de descoberta]

---

## 2. 📁 Estrutura de Output Esperada

\`\`\`
/modules/ui_web/internal/handler/${FEATURE_NAME}_handler.go
/modules/ui_web/templates/${FEATURE_NAME}_simple.html
/modules/ui_web/internal/handler/${FEATURE_NAME}_handler_test.go
/docs/implementation_plans/${FEATURE_NAME}_implementation_$(date +%Y%m%d).md
\`\`\`

---

## 3. 🛠️ Tarefas de Implementação

### **3.1 HTTP Handler (\`${FEATURE_NAME}Handler\`)**
- [ ] Criar controlador estendendo \`BaseHandler\`
- [ ] Implementar rotas HTMX:
  - \`GET /${FEATURE_NAME}\` (renderiza página)
  - \`POST /${FEATURE_NAME}\` (criação via formulário)
  - \`POST /${FEATURE_NAME}/{id}/toggle-status\` (ação HTMX)
- [ ] Instanciar e consumir serviço correspondente
- [ ] Extrair \`entity_id\` do contexto/query

### **3.2 Template HTMX (\`${FEATURE_NAME}_simple.html\`)**
- [ ] Construir interface com paleta "Soberania e Suor"
- [ ] Incluir header/nav padrão (copiar de \`dashboard_simple.html\`)
- [ ] Criar formulário assíncrono (HTMX) para adição
- [ ] Implementar lista/cards com: [campos relevantes]
- [ ] Adicionar botões de ação com feedback visual via HTMX swaps

### **3.3 Atualização da Navegação**
- [ ] Inserir link para \`/${FEATURE_NAME}\` no header de \`dashboard_simple.html\`
- [ ] Replicar navegação em templates principais

### **3.4 Testes TDD**
- [ ] \`Test${FEATURE_NAME^}Handler_List${FEATURE_NAME^}\` - Renderização
- [ ] \`Test${FEATURE_NAME^}Handler_Create${FEATURE_NAME^}\` - Criação via POST
- [ ] \`Test${FEATURE_NAME^}Handler_ToggleStatus\` - Alternância de status

---

## 4. ✅ Critérios de Aceite (Definition of Done)

### **Arquitetura**
- [ ] Handler utiliza abordagem cache-proof (\`ExecuteTemplate\` do \`BaseHandler\`)
- [ ] Soberania mantida: dados só acessados no arquivo \`.sqlite\` da entidade
- [ ] Anti-Float compliance: zero \`float\` para valores financeiros/tempo

### **Frontend**
- [ ] Design segue preceitos de Tecnologia Social
- [ ] Interface acessível com botões grandes e contrastes adequados
- [ ] Feedback amigável para erros

### **Funcionalidade**
- [ ] CRUD completo via HTMX (Create, Read, Update/Delete)
- [ ] Validações capturadas e exibidas como alertas amigáveis
- [ ] Navegação unificada em templates principais

### **Qualidade**
- [ ] Testes unitários com cobertura >90% para handler
- [ ] Testes de integração com banco SQLite real
- [ ] Código segue convenções do projeto

---

## 5. 📅 Cronograma Estimado

1. **Dia 1:** Implementação do Handler e testes unitários
2. **Dia 2:** Desenvolvimento do template \`${FEATURE_NAME}_simple.html\`
3. **Dia 3:** Integração com navegação e testes de integração
4. **Dia 4:** Validação final, correções, atualização de documentação

---

## 6. 📝 Código de Referência

### **Estrutura de Handler (baseado em MemberHandler)**
\`\`\`go
package handler

import (
    "github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

type ${FEATURE_NAME^}Handler struct {
    *BaseHandler
    lifecycleManager lifecycle.LifecycleManager
}

func New${FEATURE_NAME^}Handler(lm lifecycle.LifecycleManager) (*${FEATURE_NAME^}Handler, error) {
    base := NewBaseHandler(lm, true)
    return &${FEATURE_NAME^}Handler{
        BaseHandler:      base,
        lifecycleManager: lm,
    }, nil
}
\`\`\`

---

## 🚀 PRÓXIMOS PASSOS

1. Completar fase de descoberta (checklist)
2. Revisar e ajustar este plano
3. Iniciar implementação com TDD
4. Documentar aprendizados com \`./conclude_task.sh\`

---

**📌 Nota:** Atualizar este plano durante a implementação.
Arquivo: ${PLAN_FILE}
EOF
    
    echo "✅ Plano de implementação gerado: ${PLAN_FILE}"
fi

# 5. Modo: Executar (preparar para opencode)
if [ "$MODE" = "execute" ]; then
    echo ""
    echo "🚀 PREPARANDO EXECUÇÃO COM OPENCODE..."
    echo "====================================="
    
    # Criar prompt para opencode
    OPENCODE_PROMPT_FILE=".opencode_task_${TASK_ID}.txt"
    PROMPT_DATE=$(date +%d/%m/%Y\ %H:%M:%S)
    
    cat > ${OPENCODE_PROMPT_FILE} << EOF
## 🎯 TAREFA PARA OPENCODE

**ID:** ${TASK_ID}
**Data:** ${PROMPT_DATE}
$(if [ "$FILE_MODE" = true ]; then echo "**Fonte:** Arquivo \`${FILE_NAME}\`"; else echo "**Descrição original:** ${TASK_DESCRIPTION}"; fi)

## 📋 CONTEXTO EXTRAÍDO
- **Tipo:** ${TASK_TYPE}
- **Módulo:** ${MODULE} 
- **Objetivo:** ${OBJECTIVE}
- **Decisões:** ${DECISIONS}
- **Feature Name:** ${FEATURE_NAME}
$(if [ "$FILE_MODE" = true ]; then echo "- **Arquivo fonte:** \`${FILE_PATH}\`"; fi)

## 📁 ARQUIVOS GERADOS
1. Checklist pré-implementação: \`${CHECKLIST_FILE}\`
2. Plano de implementação: \`${PLAN_FILE}\`
3. Este prompt: \`${OPENCODE_PROMPT_FILE}\`
$(if [ "$FILE_MODE" = true ]; then echo "4. Arquivo fonte: \`${FILE_PATH}\`"; fi)

## 🚀 INSTRUÇÕES PARA OPENCODE

1. **PRIMEIRO:** Ler e entender o checklist (\`${CHECKLIST_FILE}\`)
2. **SEGUNDO:** Seguir o plano de implementação (\`${PLAN_FILE}\`)
3. **TERCEIRO:** Implementar usando TDD e padrões estabelecidos
4. **QUARTO:** **VALIDAÇÃO OBRIGATÓRIA:** Após implementar, executar:
   \`\`\`bash
   ./scripts/dev/smoke_test_new_feature.sh "${OBJECTIVE}" "/${FEATURE_NAME}"
   \`\`\`
   - Se smoke test passar: continuar
   - Se falhar: corrigir antes de prosseguir
5. **QUINTO:** **VALIDAÇÃO E2E OBRIGATÓRIA:** Executar validação end-to-end:
   \`\`\`bash
   ./scripts/dev/validate_e2e.sh --basic --headless
   \`\`\`
   - ✅ Se passar: documentar resultado
   - ❌ Se falhar: CORRIGIR antes de marcar como completo
   - ⚠️  Modos disponíveis:
     - \`--basic\` (7 passos padrão - recomendado)
     - \`--full\` (todos testes)
     - \`--ui\` (com navegador visível para debug)
     - \`--headless\` (stealth mode - padrão)
6. **SEXTO:** Documentar aprendizados com \`./conclude_task.sh\`

## 🔍 ANÁLISE PRÉVIA NECESSÁRIA

Antes de codificar, verificar:
- [ ] Backend correspondente existe em \`core_lume\`?
- [ ] É acessível (não \`internal\`)?
- [ ] Qual handler é referência (MemberHandler, CashHandler, etc.)?
- [ ] Quais funções de template são necessárias?

## 🧪 VALIDAÇÃO E2E ESPERADA

**Fluxo de 7 passos que deve funcionar após implementação:**
1. [ ] Login no sistema (\`cafe_digna\` / \`cd0123\`)
2. [ ] Acesso à nova feature (\`/${FEATURE_NAME}\`)
3. [ ] Criação de item (se aplicável)
4. [ ] Listagem de itens
5. [ ] Edição/atualização (se aplicável)
6. [ ] Integração com navegação existente
7. [ ] Smoke test passa (\`./scripts/dev/smoke_test_new_feature.sh\`)

**Comando de validação E2E (executar APÓS implementação):**
\`\`\`bash
./scripts/dev/validate_e2e.sh --basic --headless
\`\`\`

**Critério de aceite:** Validação E2E deve passar antes de marcar tarefa como completa.

## 📝 FORMATO DE RESPOSTA ESPERADO

Iniciar com análise baseada no checklist, depois implementação passo a passo.
Usar todo o sistema criado (checklists, templates, antipadrões).

**Exemplo de início:**
"Analisando tarefa ${TASK_ID}. Primeiro, vou verificar se o backend existe em core_lume..."
EOF
    
    echo "✅ Prompt para opencode gerado: ${OPENCODE_PROMPT_FILE}"
    
    # Atualizar contexto do agente com referência à tarefa
    if [ -f ".agent_context.md" ]; then
        sed -i "/^## 🎯 COMO PROCEDER AGORA/i\\
## 📋 TAREFA ATIVA\\
**ID:** ${TASK_ID}\\
**Objetivo:** ${OBJECTIVE}\\
$(if [ "$FILE_MODE" = true ]; then echo "**Fonte:** Arquivo \`${FILE_NAME}\`\\\\"; fi)\\
**Plano:** \`${PLAN_FILE}\`\\
**Prompt:** \`${OPENCODE_PROMPT_FILE}\`\\
**Status:** AGUARDANDO IMPLEMENTAÇÃO\\
\\
**INSTRUÇÃO PARA AGENTE:** Implemente seguindo o plano acima. Use o prompt gerado como guia." .agent_context.md
        echo "✅ Contexto do agente atualizado com tarefa ativa"
    fi
    
    echo ""
    echo "📋 CONTEÚDO DO PROMPT:"
    echo "======================"
    cat ${OPENCODE_PROMPT_FILE}
    echo ""
    echo "🎯 AGORA COPIE E COLE O CONTEÚDO ACIMA NO OPENCODE"
    echo "   ou use: cat ${OPENCODE_PROMPT_FILE} | pbcopy (Mac)"
    echo "   ou: cat ${OPENCODE_PROMPT_FILE} | xclip -selection clipboard (Linux)"
    echo ""
    echo "💡 O contexto do agente (.agent_context.md) também foi atualizado"
fi

# 6. Resumo final
echo ""
echo "✅ PROCESSAMENTO CONCLUÍDO!"
echo "==========================="
echo "Tarefa ID: ${TASK_ID}"
echo "Feature: ${FEATURE_NAME}"
echo "Modo: ${MODE}"
echo ""
echo "📁 ARQUIVOS CRIADOS:"
[ "$MODE" = "checklist" ] || [ "$MODE" = "plan" ] || [ "$MODE" = "execute" ] && echo "  - ${CHECKLIST_FILE}"
[ "$MODE" = "plan" ] || [ "$MODE" = "execute" ] && echo "  - ${PLAN_FILE}"
[ "$MODE" = "execute" ] && echo "  - ${OPENCODE_PROMPT_FILE}"
echo "  - ${TASK_FILE} (metadados)"
echo ""
echo "🚀 PRÓXIMOS PASSOS:"

case $MODE in
    "checklist")
        echo "  1. Preencher checklist: ${CHECKLIST_FILE}"
        echo "  2. Rodar novamente: ./process_task.sh \"${TASK_DESCRIPTION}\" --plan"
        ;;
    "plan")
        echo "  1. Revisar plano: ${PLAN_FILE}"
        echo "  2. Ajustar conforme checklist preenchido"
        echo "  3. Rodar: ./process_task.sh \"${TASK_DESCRIPTION}\" --execute"
        ;;
    "execute")
        echo "  1. Copiar prompt: ${OPENCODE_PROMPT_FILE}"
        echo "  2. Colar no opencode"
        echo "  3. Seguir implementação"
        ;;
    *)
        echo "  1. Escolher modo: --checklist, --plan, ou --execute"
        ;;
esac

echo ""
echo "💡 Dica: Use './conclude_task.sh' ao final para documentar aprendizados."

exit 0