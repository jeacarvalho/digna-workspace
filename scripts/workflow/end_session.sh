#!/bin/bash
# end_session.sh - Encerra uma sessão no projeto Digna
# Uso: ./end_session.sh [opcional: "force" para forçar encerramento]

set -e

echo "🔚 Encerrando sessão no projeto Digna..."
echo "========================================="

# Configurações
SESSION_DIR="work_in_progress/current_session"
FORCE_MODE="${1:-no}"

# Verificar se há sessão ativa
if [ ! -d "$SESSION_DIR" ]; then
    echo "❌ Nenhuma sessão ativa encontrada."
    echo "💡 Execute primeiro: ./start_session.sh"
    exit 1
fi

# Carregar informações da sessão
if [ -f "${SESSION_DIR}/session_info" ]; then
    source "${SESSION_DIR}/session_info"
else
    echo "⚠️  Informações da sessão não encontradas."
    SESSION_ID="unknown"
    START_TIME=$(date +%s)
fi

# Verificar se há tarefas pendentes
PENDING_TASKS=0
if [ -d "work_in_progress/tasks" ]; then
    PENDING_TASKS=$(find work_in_progress/tasks -name "task_metadata" -exec grep -l "STATUS=\"pending\"" {} \; 2>/dev/null | wc -l)
    IN_PROGRESS_TASKS=$(find work_in_progress/tasks -name "task_metadata" -exec grep -l "STATUS=\"in_progress\"" {} \; 2>/dev/null | wc -l)
    TOTAL_PENDING=$((PENDING_TASKS + IN_PROGRESS_TASKS))
    
    if [ "$TOTAL_PENDING" -gt 0 ] && [ "$FORCE_MODE" != "force" ]; then
        echo "❌ Há ${TOTAL_PENDING} tarefa(s) pendente(s) ou em andamento."
        echo ""
        echo "📋 Tarefas pendentes:"
        find work_in_progress/tasks -name "task_metadata" -exec sh -c '
            TASK_DIR=$(dirname {})
            TASK_ID=$(basename $TASK_DIR)
            TASK_NAME=$(grep "TASK_NAME=" {} | cut -d= -f2 | tr -d "\"")
            STATUS=$(grep "STATUS=" {} | cut -d= -f2 | tr -d "\"")
            echo "  - ${TASK_ID}: ${TASK_NAME} (${STATUS})"
        ' \; 2>/dev/null || true
        echo ""
        echo "💡 Ações possíveis:"
        echo "  1. Concluir tarefas: ./conclude_task.sh --task=ID \"Aprendizados\""
        echo "  2. Forçar encerramento: ./end_session.sh force"
        echo ""
        exit 1
    fi
fi

# Verificar qualidade das tarefas concluídas (testes)
echo "🔍 VERIFICANDO QUALIDADE DAS TAREFAS CONCLUÍDAS..."
echo "================================================="

# Contar tarefas concluídas na sessão
COMPLETED_TASKS=0
TASKS_WITH_TESTS=0

if [ -d "work_in_progress/archive/session_${SESSION_ID}/tasks" ]; then
    for TASK_DIR in work_in_progress/archive/session_${SESSION_ID}/tasks/task_*; do
        if [ -d "$TASK_DIR" ]; then
            COMPLETED_TASKS=$((COMPLETED_TASKS + 1))
            
            # Verificar se há testes associados à tarefa
            TASK_ID=$(basename "$TASK_DIR" | sed 's/task_//')
            TASK_TEST_FILES=$(find modules -name "*test*.go" -newer "${TASK_DIR}/task_prompt.md" 2>/dev/null | wc -l)
            
            if [ "$TASK_TEST_FILES" -gt 0 ]; then
                TASKS_WITH_TESTS=$((TASKS_WITH_TESTS + 1))
            fi
        fi
    done
fi

if [ "$COMPLETED_TASKS" -gt 0 ]; then
    TEST_COVERAGE_PERCENT=$((TASKS_WITH_TESTS * 100 / COMPLETED_TASKS))
    
    echo "📊 ESTATÍSTICAS DE QUALIDADE:"
    echo "  - Tarefas concluídas: $COMPLETED_TASKS"
    echo "  - Tarefas com testes: $TASKS_WITH_TESTS"
    echo "  - Cobertura de testes: $TEST_COVERAGE_PERCENT%"
    
    if [ "$TEST_COVERAGE_PERCENT" -lt 80 ]; then
        echo "⚠️  ALERTA: Baixa cobertura de testes ($TEST_COVERAGE_PERCENT%)"
        echo "💡 Recomendado: Crie mais testes antes de encerrar sessão"
        
        if [ "$FORCE_MODE" != "force" ]; then
            echo ""
            echo "❓ Deseja continuar mesmo com baixa cobertura de testes?"
            echo "  1. Continuar: ./end_session.sh force"
            echo "  2. Criar testes primeiro"
            echo ""
            exit 1
        fi
    else
        echo "✅ Boa cobertura de testes ($TEST_COVERAGE_PERCENT%)"
    fi
fi

# Calcular duração da sessão
END_TIME=$(date +%s)
DURATION=$((END_TIME - START_TIME))
DURATION_HOURS=$((DURATION / 3600))
DURATION_MIN=$(((DURATION % 3600) / 60))

echo "📊 RESUMO DA SESSÃO:"
echo "===================="
echo "ID: ${SESSION_ID}"
echo "Duração: ${DURATION_HOURS}h${DURATION_MIN}m"
echo "Tarefas concluídas: $(find work_in_progress/archive/session_${SESSION_ID}/tasks -type d -name "task_*" 2>/dev/null | wc -l)"
echo "Tarefas pendentes: ${TOTAL_PENDING:-0}"
echo ""

# 1. Consolidar aprendizados da sessão
echo "📚 CONSOLIDANDO APRENDIZADOS DA SESSÃO..."
echo "========================================="

SESSION_LEARNINGS_DIR="${SESSION_DIR}/session_learnings"
ARCHIVE_DIR="work_in_progress/archive/session_${SESSION_ID}"
mkdir -p "${ARCHIVE_DIR}"

if [ -d "$SESSION_LEARNINGS_DIR" ] && [ "$(ls -A $SESSION_LEARNINGS_DIR 2>/dev/null)" ]; then
    # Criar resumo consolidado dos aprendizados
    CONSOLIDATED_LEARNINGS="${ARCHIVE_DIR}/session_learnings_consolidated.md"
    
    cat > "${CONSOLIDATED_LEARNINGS}" << EOF
# 📚 APRENDIZADOS CONSOLIDADOS - Sessão ${SESSION_ID}

**Sessão:** ${SESSION_ID}
**Data:** $(date +%d/%m/%Y)
**Duração:** ${DURATION_HOURS}h${DURATION_MIN}m
**Tarefas concluídas:** $(find ${ARCHIVE_DIR}/tasks -type d -name "task_*" 2>/dev/null | wc -l)
**Tarefas pendentes:** ${TOTAL_PENDING:-0}

---

## 📋 RESUMO POR TAREFA

EOF
    
    # Adicionar resumo de cada tarefa
    for TASK_LEARNINGS in "${SESSION_LEARNINGS_DIR}"/*.md; do
        if [ -f "$TASK_LEARNINGS" ]; then
            TASK_NAME=$(grep "^# 📚 Aprendizados da Tarefa:" "$TASK_LEARNINGS" | sed 's/# 📚 Aprendizados da Tarefa: //')
            TASK_STATUS=$(grep "**Status:**" "$TASK_LEARNINGS" | head -1 | awk '{print $2}')
            TASK_DURATION=$(grep "**Duração:**" "$TASK_LEARNINGS" | head -1 | awk '{print $2}')
            
            if [ -n "$TASK_NAME" ]; then
                echo "### ${TASK_NAME}" >> "${CONSOLIDATED_LEARNINGS}"
                echo "- **Status:** ${TASK_STATUS:-Desconhecido}" >> "${CONSOLIDATED_LEARNINGS}"
                echo "- **Duração:** ${TASK_DURATION:-Desconhecido}" >> "${CONSOLIDATED_LEARNINGS}"
                
                # Extrair aprendizados principais (primeiras linhas após "Aprendizados Documentados")
                MAIN_LEARNINGS=$(sed -n '/## 🎯 Aprendizados Documentados/,/^---/p' "$TASK_LEARNINGS" | head -10 | tail -n +2 | sed '/^---/d')
                if [ -n "$MAIN_LEARNINGS" ]; then
                    echo "- **Aprendizados:** ${MAIN_LEARNINGS:0:100}..." >> "${CONSOLIDATED_LEARNINGS}"
                fi
                
                echo "" >> "${CONSOLIDATED_LEARNINGS}"
            fi
        fi
    done
    
    cat >> "${CONSOLIDATED_LEARNINGS}" << EOF
---

## 📈 APRENDIZADOS GERAIS DA SESSÃO

### O que funcionou bem:
1. [Listar sucessos da sessão]
2. [Processos que aceleraram o trabalho]
3. [Ferramentas úteis descobertas]

### Problemas recorrentes:
1. [Problemas que apareceram em múltiplas tarefas]
2. [Ineficiências no processo]
3. [Falhas de comunicação/entendimento]

### Melhorias identificadas:
1. [Melhorias no fluxo de trabalho]
2. [Ferramentas a adicionar/criar]
3. [Processos a otimizar]

---

## 🎯 RECOMENDAÇÕES PARA PRÓXIMA SESSÃO

### Antes de começar:
1. [Preparação recomendada]
2. [Contexto a revisar]
3. [Ferramentas a configurar]

### Durante a sessão:
1. [Boas práticas a manter]
2. [Armadilhas a evitar]
3. [Checkpoints recomendados]

### Após a sessão:
1. [Documentação prioritária]
2. [Validações obrigatórias]
3. [Métricas a coletar]

---

## 📊 MÉTRICAS DA SESSÃO

### Produtividade:
- **Tarefas/hora:** $(echo "scale=2; $(find ${ARCHIVE_DIR}/tasks -type d -name "task_*" 2>/dev/null | wc -l) / ${DURATION_HOURS}" | bc 2>/dev/null || echo "N/A")
- **Taxa de conclusão:** $(echo "scale=0; $(find ${ARCHIVE_DIR}/tasks -type d -name "task_*" 2>/dev/null | wc -l) * 100 / ($(find ${ARCHIVE_DIR}/tasks -type d -name "task_*" 2>/dev/null | wc -l) + ${TOTAL_PENDING:-0})" | bc 2>/dev/null || echo "N/A")%

### Qualidade:
- **Testes passando:** [Inserir após execução]
- **Bugs introduzidos:** [Número estimado]
- **Feedback do código:** [Positivo/Neutro/Negativo]

### Processo:
- **Tempo em análise:** [X]%
- **Tempo em implementação:** [X]%
- **Tempo em testes:** [X]%
- **Tempo em documentação:** [X]%

---

**📌 Nota:** Esta consolidação deve ser revisada antes da próxima sessão.
Use estes aprendizados para melhorar continuamente o processo de desenvolvimento.
EOF
    
    echo "✅ Aprendizados consolidados: ${CONSOLIDATED_LEARNINGS}"
    
    # Copiar para docs/learnings/ (permanente)
    PERMANENT_LEARNINGS="docs/learnings/SESSION_${SESSION_ID}_CONSOLIDATED.md"
    mkdir -p "docs/learnings"
    cp "${CONSOLIDATED_LEARNINGS}" "${PERMANENT_LEARNINGS}"
    echo "✅ Cópia permanente: ${PERMANENT_LEARNINGS}"
    
else
    echo "ℹ️  Nenhum aprendizado encontrado para consolidar"
fi

# 2. Mover tarefas pendentes para archive (se forçando)
if [ "$FORCE_MODE" = "force" ] && [ -d "work_in_progress/tasks" ]; then
    echo ""
    echo "📦 ARQUIVANDO TAREFAS PENDENTES..."
    echo "================================="
    
    mkdir -p "${ARCHIVE_DIR}/tasks_incomplete"
    mv work_in_progress/tasks/* "${ARCHIVE_DIR}/tasks_incomplete/" 2>/dev/null || true
    echo "✅ Tarefas pendentes arquivadas em: ${ARCHIVE_DIR}/tasks_incomplete/"
fi

# 3. Mover toda a sessão para archive
echo ""
echo "📦 ARQUIVANDO SESSÃO COMPLETA..."
echo "================================"

# Mover conteúdo da sessão atual
if [ -d "${SESSION_DIR}" ]; then
    # Já movemos tasks anteriormente, mover o restante
    mv "${SESSION_DIR}"/* "${ARCHIVE_DIR}/" 2>/dev/null || true
    echo "✅ Sessão arquivada em: ${ARCHIVE_DIR}"
fi

# 4. Limpar diretórios temporários
echo ""
echo "🧹 LIMPANDO DIRETÓRIOS TEMPORÁRIOS..."
echo "===================================="

rm -rf "${SESSION_DIR}" 2>/dev/null || true
rm -rf "work_in_progress/tasks" 2>/dev/null || true

# Remover link simbólico do contexto
rm -f ".agent_context.md" 2>/dev/null || true

# Remover arquivos temporários antigos no root
rm -f .session_* .task_* .opencode_task_*.txt 2>/dev/null || true

echo "✅ Diretórios temporários limpos"

# 5. Atualizar documentação permanente
echo ""
echo "📝 ATUALIZANDO DOCUMENTAÇÃO PERMANENTE..."
echo "========================================"

# Atualizar QUICK_REFERENCE com aprendizados da sessão
if [ -f "docs/QUICK_REFERENCE.md" ] && [ -f "${CONSOLIDATED_LEARNINGS}" ]; then
    # Adicionar seção de aprendizados recentes
    if ! grep -q "## 📚 Aprendizados Recentes" docs/QUICK_REFERENCE.md; then
        echo "" >> docs/QUICK_REFERENCE.md
        echo "## 📚 Aprendizados Recentes" >> docs/QUICK_REFERENCE.md
        echo "" >> docs/QUICK_REFERENCE.md
    fi
    
    # Adicionar referência a esta sessão
    SESSION_REF="**Sessão ${SESSION_ID}:** $(date +%d/%m/%Y) - ${DURATION_HOURS}h${DURATION_MIN}m, $(find ${ARCHIVE_DIR}/tasks -type d -name "task_*" 2>/dev/null | wc -l) tarefas"
    echo "- ${SESSION_REF} (ver \`docs/learnings/SESSION_${SESSION_ID}_CONSOLIDATED.md\`)" >> docs/QUICK_REFERENCE.md
    
    echo "✅ QUICK_REFERENCE.md atualizado"
fi

# Atualizar data da última sessão
sed -i "s/**Última sessão:**.*/**Última sessão:** ${SESSION_ID} ($(date +%d\/%m\/%Y))/" docs/QUICK_REFERENCE.md 2>/dev/null || true

echo ""
echo "✅ SESSÃO ENCERRADA COM SUCESSO!"
echo "================================"
echo ""
echo "📊 RESUMO FINAL:"
echo "---------------"
echo "Sessão: ${SESSION_ID}"
echo "Duração: ${DURATION_HOURS}h${DURATION_MIN}m"
echo "Tarefas concluídas: $(find ${ARCHIVE_DIR}/tasks -type d -name "task_*" 2>/dev/null | wc -l)"
echo "Aprendizados: ${PERMANENT_LEARNINGS:-Nenhum}"
echo ""
echo "📁 ARQUIVOS CRIADOS:"
echo "-------------------"
echo "1. Archive da sessão: ${ARCHIVE_DIR}"
echo "2. Aprendizados consolidados: ${PERMANENT_LEARNINGS:-Nenhum}"
echo "3. QUICK_REFERENCE.md atualizado"
echo ""
echo "🚀 PRÓXIMA SESSÃO:"
echo "-----------------"
echo "1. Iniciar nova sessão: ./start_session.sh"
echo "2. Revisar aprendizados: cat ${PERMANENT_LEARNINGS:-docs/learnings/}"
echo "3. Escolher próxima tarefa do backlog: docs/NEXT_STEPS.md"
echo ""
echo "💡 Dica: Os aprendizados desta sessão serão usados para melhorar"
echo "       o processo da próxima sessão!"
echo ""
echo "🤖 INSTRUÇÃO PARA OPENCODE:"
echo "---------------------------"
echo "A sessão foi encerrada. Para nova sessão, execute: ./start_session.sh"