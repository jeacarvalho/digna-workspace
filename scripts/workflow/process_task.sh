#!/bin/bash
# process_task.sh - Processa uma tarefa no projeto Digna
# Uso: ./process_task.sh --task=TASK_ID [--checklist] [--plan] [--execute]

set -e

echo "🔧 Processando tarefa..."
echo "========================"

# Configurações
SESSION_DIR="work_in_progress/current_session"
TASK_ID=""
MODE="analyze"

# Funções de ajuda
show_help() {
    cat << EOF
🎯 process_task.sh - Processador de Tarefas Digna (Nova Estrutura)

USO:
  ./process_task.sh [--task=TASK_ID] [OPÇÕES]

OPÇÕES:
  --task, -t ID      ID da tarefa (opcional - usa tarefa mais nova ativa se omitido)
  --checklist, -c    Apenas gerar checklist pré-implementação
  --plan, -p         Gerar plano de implementação completo  
  --execute, -e      Executar implementação (interage com opencode)
  --help, -h         Mostrar esta ajuda

EXEMPLOS:
  # Listar tarefas disponíveis
  ls work_in_progress/tasks/
  
  # Usar tarefa mais nova ativa automaticamente
  ./process_task.sh --checklist
  ./process_task.sh --plan
  ./process_task.sh --execute
  
  # Especificar tarefa manualmente
  ./process_task.sh --task=20250311_101108 --checklist
  ./process_task.sh --task=20250311_101108 --plan
  ./process_task.sh --task=20250311_101108 --execute

FLUXO RECOMENDADO:
  1. ./create_task.sh "Nome da Tarefa" [módulo]
  2. Editar: vim work_in_progress/tasks/task_[ID]/task_prompt.md
  3. ./process_task.sh --task=[ID] --checklist
  4. ./process_task.sh --task=[ID] --plan
  5. ./process_task.sh --task=[ID] --execute
  6. ./conclude_task.sh --task=[ID] "Aprendizados"

ESTRUTURA DA TAREFA:
  work_in_progress/tasks/task_[ID]/
  ├── task_prompt.md      # Prompt da tarefa
  ├── checklist.md        # Checklist gerado
  ├── implementation_plan.md
  ├── task_learnings.md
  └── task_artifacts/     # Arquivos temporários
EOF
    exit 0
}

# Processar argumentos
for arg in "$@"; do
    case $arg in
        --task=*|-t=*)
            TASK_ID="${arg#*=}"
            ;;
        --task|-t)
            echo "❌ Use: --task=ID ou -t=ID"
            exit 1
            ;;
        --checklist|-c)
            MODE="checklist"
            ;;
        --plan|-p)
            MODE="plan"
            ;;
        --execute|-e)
            MODE="execute"
            ;;
        --help|-h)
            show_help
            ;;
        *)
            # Ignorar argumentos desconhecidos
            ;;
    esac
done

# Função simples para buscar tarefa mais nova ativa
find_latest_active_task() {
    # Verificar se o diretório existe
    [ ! -d "work_in_progress/tasks" ] && return 1
    
    local latest_task=""
    
    # Listar diretórios de tarefa
    for task_dir in work_in_progress/tasks/task_*; do
        # Pular se não for diretório
        [ ! -d "$task_dir" ] && continue
        
        if [ -f "${task_dir}/task_metadata" ]; then
            # Extrair status
            local status
            status=$(grep '^STATUS=' "${task_dir}/task_metadata" 2>/dev/null | cut -d= -f2 | tr -d '"' || echo "unknown")
            
            # Considerar apenas pending ou in_progress
            if [ "$status" = "pending" ] || [ "$status" = "in_progress" ]; then
                # Extrair ID da tarefa
                local task_id
                task_id=$(basename "$task_dir" | sed 's/task_//')
                
                # Manter a mais nova (comparação lexicográfica funciona para YYYYMMDD_HHMMSS)
                if [ -z "$latest_task" ] || [[ "$task_id" > "$latest_task" ]]; then
                    latest_task="$task_id"
                fi
            fi
        fi
    done
    
    [ -n "$latest_task" ] && echo "$latest_task" && return 0
    return 1
}

# Verificar se temos ID da tarefa
if [ -z "$TASK_ID" ]; then
    echo "🔍 Buscando tarefa mais nova ativa..."
    
    # Tentar encontrar tarefa mais nova ativa
    if LATEST_TASK_ID=$(find_latest_active_task); then
        echo "✅ Usando tarefa mais nova ativa: ${LATEST_TASK_ID}"
        TASK_ID="$LATEST_TASK_ID"
    else
        echo "❌ Nenhuma tarefa ativa encontrada."
        echo ""
        echo "📋 Tarefas disponíveis:"
        if [ -d "work_in_progress/tasks" ]; then
            find work_in_progress/tasks -name "task_metadata" -exec sh -c 'echo "  - $(basename $(dirname {})): $(grep "TASK_NAME=" {} | cut -d= -f2 | tr -d "\"") (Status: $(grep "STATUS=" {} | cut -d= -f2 | tr -d "\"" 2>/dev/null || echo "unknown"))"' \; 2>/dev/null || echo "  Nenhuma tarefa encontrada"
        else
            echo "  Nenhuma tarefa encontrada"
        fi
        echo ""
        echo "💡 Use: ./process_task.sh --task=ID"
        echo "💡 Ou crie nova tarefa: ./create_task.sh \"Nome da Tarefa\""
        exit 1
    fi
fi



# Verificar se tarefa existe
TASK_DIR="work_in_progress/tasks/task_${TASK_ID}"
if [ ! -d "$TASK_DIR" ]; then
    echo "❌ Tarefa não encontrada: ${TASK_DIR}"
    echo "💡 Verifique o ID ou crie nova tarefa: ./create_task.sh"
    exit 1
fi

# Verificar se há sessão ativa
if [ ! -d "$SESSION_DIR" ]; then
    echo "❌ Nenhuma sessão ativa encontrada."
    echo "💡 Execute primeiro: ./start_session.sh"
    exit 1
fi

# Carregar metadados da tarefa
if [ -f "${TASK_DIR}/task_metadata" ]; then
    source "${TASK_DIR}/task_metadata"
else
    echo "⚠️  Metadados da tarefa não encontrados, usando valores padrão."
    TASK_NAME="Tarefa ${TASK_ID}"
    MODULE="ui_web"
    STATUS="pending"
fi

# Verificar prompt da tarefa
TASK_PROMPT="${TASK_DIR}/task_prompt.md"
if [ ! -f "$TASK_PROMPT" ]; then
    echo "❌ Prompt da tarefa não encontrado: ${TASK_PROMPT}"
    echo "💡 Crie o arquivo task_prompt.md no diretório da tarefa"
    exit 1
fi

echo "📋 Tarefa: ${TASK_NAME}"
echo "📁 Diretório: ${TASK_DIR}"
echo "🎯 Modo: ${MODE}"
echo ""

# Atualizar status da tarefa
sed -i "s/STATUS=\"[^\"]*\"/STATUS=\"in_progress\"/g" "${TASK_DIR}/task_metadata" 2>/dev/null || true

# Criar arquivo de tarefa no formato antigo (para compatibilidade)
TASK_FILE=".task_${TASK_ID}"
cat > ${TASK_FILE} << EOF
TASK_ID=${TASK_ID}
TASK_DESCRIPTION="${TASK_NAME}"
TASK_TYPE="Feature"
MODULE="${MODULE}"
OBJECTIVE="Implementar ${TASK_NAME}"
DECISIONS="Seguir padrões estabelecidos"
FEATURE_NAME="$(echo "${TASK_NAME}" | tr ' ' '_' | tr '[:upper:]' '[:lower:]')"
START_TIME=$(date +%s)
MODE="${MODE}"
FILE_MODE=false
FILE_PATH=""
FILE_NAME=""
EOF

# Ler conteúdo do prompt
PROMPT_CONTENT=$(cat "$TASK_PROMPT")

# Criar arquivo para opencode (se modo execute)
if [ "$MODE" = "execute" ]; then
    OPENCODE_FILE=".opencode_task_${TASK_ID}.txt"
    
    cat > ${OPENCODE_FILE} << EOF
# 📋 TAREFA: ${TASK_NAME}
# ID: ${TASK_ID}
# Data: $(date +%d/%m/%Y %H:%M:%S)
# Módulo: ${MODULE}

${PROMPT_CONTENT}

---

## 🎯 INSTRUÇÕES PARA OPENCODE

1. **Leia o contexto:** Consulte work_in_progress/current_session/.agent_context.md
2. **Siga padrões:** Use docs/QUICK_REFERENCE.md e docs/ANTIPATTERNS.md
3. **Analise similaridades:** Use ./scripts/tools/analyze_patterns.sh
4. **Implemente:** Crie handlers, templates, testes conforme necessário
5. **Valide:** Execute testes e smoke tests

### 📁 ESTRUTURA ESPERADA:
- Handler: modules/ui_web/internal/handler/[feature]_handler.go
- Template: modules/ui_web/templates/[feature]_simple.html  
- Testes: modules/ui_web/internal/handler/[feature]_handler_test.go
- Registro: modules/ui_web/main.go

### ⚠️ VALIDAÇÕES OBRIGATÓRIAS:
- [ ] Testes unitários implementados
- [ ] Testes E2E com Playwright implementados (fluxo completo)
- [ ] Handler registrado no main.go
- [ ] Smoke test executado
- [ ] Documentação atualizada
- [ ] **CRÍTICO:** Testar todas as rotas novas/modificadas
- [ ] **CRÍTICO:** Validar que não há regressões em funcionalidades existentes

### ⚠️⚠️⚠️ INSTRUÇÃO CRÍTICA DE FLUXO ⚠️⚠️⚠️
**O AGENTE OPENCODE NUNCA DEVE EXECUTAR ./conclude_task.sh AUTOMATICAMENTE**

**FLUXO CORRETO COM TESTES OBRIGATÓRIOS:**
1. Implemente a tarefa conforme instruções acima
2. **CRÍTICO:** Crie/execute testes E2E com Playwright para validar funcionalidade
3. **CRÍTICO:** Execute testes unitários e de integração
4. **CRÍTICO:** Teste manualmente todas as rotas novas/modificadas
5. Após implementação E TESTES COMPLETOS, INFORME ao usuário:
   "A tarefa ${TASK_ID} ('${TASK_NAME}') está implementada e TESTADA. PODE ser concluída"
6. AGUARDE o usuário executar:
   ./conclude_task.sh --task=${TASK_ID} "Aprendizados: [descreva aprendizados]"
7. Só após o usuário executar o comando, a tarefa será arquivada

**TESTES OBRIGATÓRIOS PARA TODAS AS TAREFAS:**
- ✅ Testes unitários para handlers/services
- ✅ Testes E2E com Playwright (fluxo completo no browser)
- ✅ Validação de rotas HTTP (status codes, responses)
- ✅ Teste de regressão (não quebrar funcionalidades existentes)
- ✅ Teste com dados reais (cafe_digna, contador_social)

**Status:** AGUARDANDO IMPLEMENTAÇÃO
EOF
    
    echo "✅ Arquivo para opencode criado: ${OPENCODE_FILE}"
    echo ""
    echo "🤖 INSTRUÇÃO PARA OPENCODE:"
    echo "============================"
    echo "Leia o arquivo: ${OPENCODE_FILE}"
    echo "Implemente seguindo as instruções acima."
    echo ""
    echo "⚠️⚠️⚠️ ATENÇÃO CRÍTICA PARA OPENCODE ⚠️⚠️⚠️"
    echo "=========================================="
    echo "1. IMPLEMENTE a tarefa conforme instruções no arquivo"
    echo "2. NÃO execute ./conclude_task.sh automaticamente"
    echo "3. APÓS implementação, INFORME ao usuário que:"
    echo "   'A tarefa está implementada e PODE ser concluída'"
    echo "4. AGUARDE o usuário executar: ./conclude_task.sh --task=${TASK_ID} \"Aprendizados\""
    echo ""
    echo "💡 Dica: O contexto atual está em: work_in_progress/current_session/.agent_context.md"
    
elif [ "$MODE" = "checklist" ]; then
    echo "📋 GERANDO CHECKLIST..."
    echo "========================"
    
    CHECKLIST_FILE="${TASK_DIR}/checklist.md"
    
    # Gerar checklist baseado no prompt
    cat > ${CHECKLIST_FILE} << EOF
# 📋 CHECKLIST PRÉ-IMPLEMENTAÇÃO: ${TASK_NAME}
# ID: ${TASK_ID}
# Gerado em: $(date +%d/%m/%Y %H:%M:%S)

## 🔍 ANÁLISE DO CONTEXTO

### 1. Contexto do Projeto
- [ ] Ler work_in_progress/current_session/.agent_context.md
- [ ] Consultar docs/QUICK_REFERENCE.md
- [ ] Verificar antipadrões em docs/ANTIPATTERNS.md
- [ ] Revisar aprendizados anteriores em docs/learnings/

### 2. Análise de Similaridades
- [ ] Encontrar handler similar: ./scripts/tools/analyze_patterns.sh [padrão]
- [ ] Analisar template similar
- [ ] Verificar padrões de testes
- [ ] Identificar funções de template usadas

### 3. Requisitos Técnicos
- [ ] Definir estrutura do handler
- [ ] Definir estrutura do template
- [ ] Definir rotas HTTP
- [ ] Definir funções de template necessárias

## 🛠️ PREPARAÇÃO TÉCNICA

### 4. Estrutura de Arquivos
- [ ] Criar: modules/ui_web/internal/handler/[feature]_handler.go
- [ ] Criar: modules/ui_web/templates/[feature]_simple.html
- [ ] Criar: modules/ui_web/internal/handler/[feature]_handler_test.go
- [ ] Atualizar: modules/ui_web/main.go (registrar handler)

### 5. Dependências
- [ ] Verificar imports necessários
- [ ] Verificar funções do lifecycle manager
- [ ] Verificar integração com outros módulos
- [ ] Verificar atualizações de banco de dados

## 🧪 TESTES E VALIDAÇÃO

### 6. Estratégia de Testes
- [ ] Definir casos de teste unitários
- [ ] Definir testes de integração
- [ ] Preparar dados de teste
- [ ] Definir critérios de aceitação

### 7. Validação Pós-Implementação
- [ ] Smoke test: ./scripts/dev/smoke_test_new_feature.sh
- [ ] Testes de sistema: cd modules && ./run_tests.sh
- [ ] Validação manual da UI
- [ ] Verificação de acessibilidade

## 📚 DOCUMENTAÇÃO

### 8. Documentação Técnica
- [ ] Atualizar docs/QUICK_REFERENCE.md
- [ ] Atualizar docs/NEXT_STEPS.md
- [ ] Criar documentação da feature
- [ ] Documentar aprendizados

### 9. Checklist de Conclusão
- [ ] Todos os testes passando
- [ ] Código revisado seguindo padrões
- [ ] Documentação atualizada
- [ ] Smoke test executado com sucesso
- [ ] Handler registrado no main.go

---

## 📝 NOTAS DA ANÁLISE

### Padrões Identificados:
[Preencher após análise]

### Riscos Identificados:
1. [Risco 1]
2. [Risco 2]

### Decisões de Design:
1. [Decisão 1]
2. [Decisão 2]

### Referências:
- Handler similar: [nome]
- Template similar: [nome]
- Testes similares: [nome]

---

**Status do Checklist:** ✅ GERADO
**Próximo passo:** ./process_task.sh --task=${TASK_ID} --plan
EOF
    
    echo "✅ Checklist gerado: ${CHECKLIST_FILE}"
    echo ""
    echo "📋 PRÓXIMO PASSO:"
    echo "  ./process_task.sh --task=${TASK_ID} --plan"
    
elif [ "$MODE" = "plan" ]; then
    echo "📝 GERANDO PLANO DE IMPLEMENTAÇÃO..."
    echo "===================================="
    
    PLAN_FILE="${TASK_DIR}/implementation_plan.md"
    
    # Extrair seções do prompt para personalizar o plano
    OBJECTIVE=$(echo "$PROMPT_CONTENT" | grep -A5 "## 🎯 OBJETIVO" | tail -n +2 | head -5 | sed '/^---/d' | tr '\n' ' ')
    REQUIREMENTS=$(echo "$PROMPT_CONTENT" | grep -A20 "## 📋 REQUISITOS" | tail -n +2 | head -20)
    
    cat > ${PLAN_FILE} << EOF
# 🚀 PLANO DE IMPLEMENTAÇÃO: ${TASK_NAME}
# ID: ${TASK_ID}
# Gerado em: $(date +%d/%m/%Y %H:%M:%S)

## 🎯 OBJETIVO
${OBJECTIVE:-Implementar a funcionalidade conforme especificado no prompt.}

## 📋 REQUISITOS
${REQUIREMENTS:-Ver requisitos no task_prompt.md}

## 🔄 FLUXO DE IMPLEMENTAÇÃO

### Fase 1: Análise e Setup (15%)
1. **Análise de código similar** (30 min)
   - Encontrar handler/template similar com ./scripts/tools/analyze_patterns.sh
   - Extrair padrões de implementação
   - Identificar funções de template necessárias

2. **Setup do ambiente** (15 min)
   - Criar estrutura de arquivos
   - Configurar imports básicos
   - Preparar dados de teste

### Fase 2: Implementação do Handler (40%)
3. **Estrutura do Handler** (45 min)
   - Criar struct do handler (estender BaseHandler)
   - Implementar construtor New[Feature]Handler
   - Adicionar funções de template específicas

4. **Lógica de Negócio** (60 min)
   - Implementar métodos HTTP (GET/POST)
   - Integrar com lifecycle manager
   - Implementar validações
   - Tratamento de erros

5. **Rotas e Registro** (15 min)
   - Implementar RegisterRoutes
   - Adicionar handler ao main.go
   - Testar rotas básicas

### Fase 3: Template e UI (25%)
6. **Template HTML** (60 min)
   - Copiar estrutura de template similar
   - Adaptar para nova feature
   - Implementar forms HTMX
   - Estilizar com Tailwind (paleta Digna)

7. **Interatividade** (30 min)
   - Implementar ações HTMX
   - Adicionar feedback visual
   - Validação client-side

### Fase 4: Testes e Validação (15%)
8. **Testes Unitários** (45 min)
   - Criar testes para handler
   - Testar casos de sucesso/erro
   - Mock de lifecycle manager

9. **Validação Integrada** (30 min)
   - Smoke test: ./scripts/dev/smoke_test_new_feature.sh
   - Testes de sistema
   - Validação manual

10. **Documentação** (15 min)
    - Atualizar QUICK_REFERENCE.md
    - Atualizar NEXT_STEPS.md
    - Documentar aprendizados

### Fase 5: Revisão e Conclusão (5%)
11. **Revisão Final** (15 min)
    - Revisar código seguindo padrões
    - Verificar antipadrões
    - Validar integridade

## 📁 ESTRUTURA DE ARQUIVOS

### A Criar:
\`\`\`
modules/ui_web/internal/handler/$(echo "${TASK_NAME}" | tr ' ' '_' | tr '[:upper:]' '[:lower:]')_handler.go
modules/ui_web/templates/$(echo "${TASK_NAME}" | tr ' ' '_' | tr '[:upper:]' '[:lower:]')_simple.html
modules/ui_web/internal/handler/$(echo "${TASK_NAME}" | tr ' ' '_' | tr '[:upper:]' '[:lower:]')_handler_test.go
\`\`\`

### A Modificar:
\`\`\`
modules/ui_web/main.go  # Registrar handler
docs/QUICK_REFERENCE.md # Adicionar referência
docs/NEXT_STEPS.md      # Marcar como concluído
\`\`\`

## ⚠️ RISCOS E MITIGAÇÕES

### Riscos Técnicos:
1. **Complexidade inesperada** - Mitigar com análise detalhada prévia
2. **Integração com módulos existentes** - Mitigar com testes de integração
3. **Performance** - Mitigar com profiling e otimizações

### Riscos de Processo:
1. **Estimativa imprecisa** - Mitigar com buffer de 20% no tempo
2. **Dependências externas** - Mitigar identificando early
3. **Mudanças de requisitos** - Mitigar com validação contínua

## 🎯 CRITÉRIOS DE ACEITAÇÃO

### Funcionais:
- [ ] Feature implementada conforme requisitos
- [ ] UI funcional e responsiva
- [ ] Integração com sistema existente

### Técnicos:
- [ ] Testes unitários com cobertura >80%
- [ ] Código segue padrões do projeto
- [ ] Documentação atualizada
- [ ] Smoke test passa

### Qualidade:
- [ ] Sem regressões identificadas
- [ ] Performance aceitável
- [ ] Código revisado e limpo

## 📊 ESTIMATIVA DE TEMPO

**Total estimado:** 5-6 horas
**Buffer recomendado:** 1 hora (20%)

### Breakdown:
- Análise: 45 min
- Handler: 2 horas
- Template: 1.5 horas
- Testes: 1.25 horas
- Documentação: 30 min
- Revisão: 15 min

---

**Status do Plano:** ✅ GERADO
**Próximo passo:** ./process_task.sh --task=${TASK_ID} --execute
**Ou implementar manualmente seguindo este plano.**
EOF
    
    echo "✅ Plano de implementação gerado: ${PLAN_FILE}"
    echo ""
    echo "📋 PRÓXIMO PASSO:"
    echo "  ./process_task.sh --task=${TASK_ID} --execute"
    
else
    echo "🔍 MODO ANÁLISE (padrão)"
    echo "========================"
    echo ""
    echo "📋 Tarefa: ${TASK_NAME}"
    echo "📁 Diretório: ${TASK_DIR}"
    echo "📄 Prompt: ${TASK_PROMPT}"
    echo ""
    echo "📋 OPÇÕES DISPONÍVEIS:"
    echo "  1. Gerar checklist:   ./process_task.sh --task=${TASK_ID} --checklist"
    echo "  2. Gerar plano:       ./process_task.sh --task=${TASK_ID} --plan"
    echo "  3. Executar:          ./process_task.sh --task=${TASK_ID} --execute"
    echo ""
    echo "💡 Recomendado: Siga a sequência checklist → plan → execute"
fi

echo ""
echo "✅ PROCESSAMENTO CONCLUÍDO"
echo "=========================="