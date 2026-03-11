#!/bin/bash
# help_agent_recover_context.sh - Ajuda o agente a recuperar contexto após compaction
# Uso: ./scripts/help_agent_recover_context.sh

set -e

echo "🤖 AJUDANDO AGENTE A RECUPERAR CONTEXTO APÓS COMPACTION"
echo "========================================================"
echo ""
echo "⚠️  Este script é para quando o agente opencode:"
echo "    1. Entrou em modo compaction (imprevisível)"
echo "    2. Perdeu o contexto da tarefa atual"
echo "    3. Precisa de ajuda para continuar"
echo ""

SESSION_DIR="work_in_progress/current_session"
CONTEXT_FILE="${SESSION_DIR}/.compaction_context.md"

# 1. Verificar se há contexto salvo
if [ -f "${CONTEXT_FILE}" ]; then
    echo "✅ Contexto salvo encontrado!"
    echo "📝 Resumo do contexto:"
    grep -E "(ID da Tarefa:|Nome:|Data/Hora do save:)" "${CONTEXT_FILE}" | head -5
    echo ""
    echo "💡 Para restaurar contexto:"
    echo "   ./preserve_context.sh --restore"
    echo ""
else
    echo "❌ Nenhum contexto salvo encontrado."
    echo ""
fi

# 2. Verificar tarefas ativas
echo "🔍 BUSCANDO TAREFAS ATIVAS..."
echo "=============================="

TASK_COUNT=$(find work_in_progress/tasks -maxdepth 1 -type d -name "task_*" 2>/dev/null | wc -l)

if [ "$TASK_COUNT" -eq 0 ]; then
    echo "❌ Nenhuma tarefa ativa encontrada."
    echo "💡 Crie uma nova tarefa: ./create_task.sh \"Nome da Tarefa\""
    exit 1
fi

# 3. Listar tarefas ativas
echo "📋 TAREFAS ATIVAS ENCONTRADAS:"
echo "-----------------------------"

for TASK_DIR in work_in_progress/tasks/task_*; do
    if [ -d "${TASK_DIR}" ]; then
        TASK_ID=$(basename "${TASK_DIR}" | sed 's/task_//')
        
        if [ -f "${TASK_DIR}/task_metadata" ]; then
            source "${TASK_DIR}/task_metadata" 2>/dev/null || true
            echo "  - ${TASK_ID}: ${TASK_NAME:-Nome não encontrado} (${STATUS:-status desconhecido})"
        else
            echo "  - ${TASK_ID}: Metadados não encontrados"
        fi
    fi
done

echo ""

# 4. Sugerir próxima ação
LATEST_TASK=$(find work_in_progress/tasks -maxdepth 1 -type d -name "task_*" -exec ls -dt {} + | head -1)
TASK_ID=$(basename "${LATEST_TASK}" | sed 's/task_//')

if [ -n "${TASK_ID}" ]; then
    echo "🎯 SUGESTÃO DE AÇÃO:"
    echo "-------------------"
    echo "Tarefa mais recente: ${TASK_ID}"
    echo ""
    
    if [ -f "${CONTEXT_FILE}" ]; then
        echo "1. Restaurar contexto: ./preserve_context.sh --restore"
        echo "2. Agente continua tarefa ${TASK_ID}"
    else
        echo "1. Processar tarefa: ./process_task.sh --task=${TASK_ID} --execute"
        echo "2. Agente implementa tarefa ${TASK_ID}"
    fi
    
    echo "3. Quando terminar, usuário executa: ./conclude_task.sh --task=${TASK_ID} \"Aprendizados\""
else
    echo "❌ Não foi possível identificar tarefa mais recente."
    echo "💡 Crie nova tarefa: ./create_task.sh \"Nome da Tarefa\""
fi

echo ""
echo "📚 COMANDOS ÚTEIS:"
echo "-----------------"
echo "• Listar tarefas: ls work_in_progress/tasks/"
echo "• Criar tarefa: ./create_task.sh \"Nome\""
echo "• Processar tarefa: ./process_task.sh --task=ID --execute"
echo "• Concluir tarefa: ./conclude_task.sh --task=ID \"Aprendizados\""
echo "• Preservar contexto: ./preserve_context.sh --save (quando detectar compaction)"
echo "• Restaurar contexto: ./preserve_context.sh --restore (após compaction)"
echo "• Ver status: ./preserve_context.sh --status"

echo ""
echo "⚠️  LEMBRETE IMPORTANTE:"
echo "----------------------"
echo "O agente opencode NÃO pode prever quando o compaction acontece."
echo "Compaction é controlado internamente pelo sistema para economizar tokens."
echo ""
echo "Quando o agente parecer 'travado' ou perguntar 'o que estamos fazendo?':"
echo "1. Execute este script para ver opções"
echo "2. Use ./preserve_context.sh --save se possível"
echo "3. Guie o agente de volta à tarefa atual"