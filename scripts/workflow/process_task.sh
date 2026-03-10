#!/bin/bash
# process_task.sh - Processa uma tarefa no projeto Digna
# Uso: ./process_task.sh "Descrição da tarefa" [--checklist] [--plan] [--execute]
# Ex: ./process_task.sh "Implementar UI para Fornecedores" --checklist

set -e

# Configurações
TASK_DESCRIPTION="${1}"
SESSION_FILE=$(ls .session_* 2>/dev/null | head -1)
TASK_ID=$(date +%Y%m%d_%H%M%S)
TASK_FILE=".task_${TASK_ID}"

# Funções de ajuda
show_help() {
    cat << EOF
🎯 process_task.sh - Processador Inteligente de Tarefas Digna

Uso: ./process_task.sh "DESCRIÇÃO DA TAREFA" [OPÇÕES]

Descrição da tarefa (FORMATO RECOMENDADO):
  "Tipo: Feature | Módulo: ui_web | Objetivo: Implementar X | Decisões: seguir padrão Y"

Opções:
  --checklist    Apenas gerar checklist pré-implementação
  --plan         Gerar plano de implementação completo
  --execute      Executar implementação (interage com opencode)
  --help         Mostrar esta ajuda

Exemplos:
  ./process_task.sh "Tipo: Feature | Módulo: ui_web | Objetivo: Implementar UI para Fornecedores"
  ./process_task.sh "Bug no PDV: erro ao adicionar produto" --checklist
  ./process_task.sh "Melhorar performance do dashboard" --plan

O script irá:
1. Analisar a descrição da tarefa
2. Verificar contexto existente
3. Gerar checklists/planos conforme opção
4. Preparar para execução com opencode
EOF
    exit 0
}

# Verificar se tem descrição
if [ -z "$TASK_DESCRIPTION" ] || [ "$TASK_DESCRIPTION" = "--help" ]; then
    show_help
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

echo "🔍 Processando tarefa: ${TASK_DESCRIPTION}"
echo "=========================================="

# 1. Extrair informações da descrição
echo "📋 Extraindo informações da descrição..."

# Padrões comuns de extração
TASK_TYPE=$(echo "$TASK_DESCRIPTION" | grep -oi "tipo:\s*[^|]*" | cut -d: -f2 | xargs || echo "Feature")
MODULE=$(echo "$TASK_DESCRIPTION" | grep -oi "módulo:\s*[^|]*" | cut -d: -f2 | xargs || echo "ui_web")
OBJECTIVE=$(echo "$TASK_DESCRIPTION" | grep -oi "objetivo:\s*[^|]*" | cut -d: -f2 | xargs || echo "$TASK_DESCRIPTION")
DECISIONS=$(echo "$TASK_DESCRIPTION" | grep -oi "decisões:\s*[^|]*" | cut -d: -f2 | xargs || echo "Seguir padrões estabelecidos")

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
EOF

echo "✅ Tarefa registrada: ${TASK_ID}"
echo "   Tipo: ${TASK_TYPE}"
echo "   Módulo: ${MODULE}"
echo "   Objetivo: ${OBJECTIVE}"

# 3. Modo: Checklist pré-implementação
if [ "$MODE" = "checklist" ] || [ "$MODE" = "plan" ] || [ "$MODE" = "execute" ]; then
    echo ""
    echo "📝 GERANDO CHECKLIST PRÉ-IMPLEMENTAÇÃO..."
    echo "========================================"
    
    CHECKLIST_FILE="docs/implementation_plans/${FEATURE_NAME}_pre_check.md"
    mkdir -p docs/implementation_plans
    
    # Template de checklist
    cat > ${CHECKLIST_FILE} << EOF
# 🔍 Checklist de Validação Pré-Implementação: ${FEATURE_NAME}

**Tarefa:** ${TASK_DESCRIPTION}
**Gerado em:** $(date +%d/%m/%Y %H:%M:%S)
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
    
    PLAN_FILE="docs/implementation_plans/${FEATURE_NAME}_implementation_$(date +%Y%m%d).md"
    
    # Verificar se checklist foi preenchido
    if [ ! -f "${CHECKLIST_FILE}" ]; then
        echo "⚠️  Checklist não encontrado. Gerando plano básico..."
    fi
    
    # Template de plano
    cat > ${PLAN_FILE} << EOF
# 📋 Plano de Implementação: ${FEATURE_NAME}

**Feature:** ${FEATURE_NAME}
**Tarefa ID:** ${TASK_ID}
**Gerado em:** $(date +%d/%m/%Y %H:%M:%S)
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
    
    cat > ${OPENCODE_PROMPT_FILE} << EOF
## 🎯 TAREFA PARA OPENCODE

**ID:** ${TASK_ID}
**Data:** $(date +%d/%m/%Y %H:%M:%S)
**Descrição original:** ${TASK_DESCRIPTION}

## 📋 CONTEXTO EXTRAÍDO
- **Tipo:** ${TASK_TYPE}
- **Módulo:** ${MODULE} 
- **Objetivo:** ${OBJECTIVE}
- **Decisões:** ${DECISIONS}
- **Feature Name:** ${FEATURE_NAME}

## 📁 ARQUIVOS GERADOS
1. Checklist pré-implementação: \`${CHECKLIST_FILE}\`
2. Plano de implementação: \`${PLAN_FILE}\`
3. Este prompt: \`${OPENCODE_PROMPT_FILE}\`

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
5. **QUINTO:** Documentar aprendizados com \`./conclude_task.sh\`

## 🔍 ANÁLISE PRÉVIA NECESSÁRIA

Antes de codificar, verificar:
- [ ] Backend correspondente existe em \`core_lume\`?
- [ ] É acessível (não \`internal\`)?
- [ ] Qual handler é referência (MemberHandler, CashHandler, etc.)?
- [ ] Quais funções de template são necessárias?

## 📝 FORMATO DE RESPOSTA ESPERADO

Iniciar com análise baseada no checklist, depois implementação passo a passo.
Usar todo o sistema criado (checklists, templates, antipadrões).

**Exemplo de início:**
"Analisando tarefa ${TASK_ID}. Primeiro, vou verificar se o backend existe em core_lume..."
EOF
    
    echo "✅ Prompt para opencode gerado: ${OPENCODE_PROMPT_FILE}"
    echo ""
    echo "📋 CONTEÚDO DO PROMPT:"
    echo "======================"
    cat ${OPENCODE_PROMPT_FILE}
    echo ""
    echo "🎯 AGORA COPIE E COLE O CONTEÚDO ACIMA NO OPENCODE"
    echo "   ou use: cat ${OPENCODE_PROMPT_FILE} | pbcopy (Mac)"
    echo "   ou: cat ${OPENCODE_PROMPT_FILE} | xclip -selection clipboard (Linux)"
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