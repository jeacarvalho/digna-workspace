#!/bin/bash
# preserve_context.sh - Preserva contexto durante compaction do opencode
# Uso: ./preserve_context.sh [--save|--restore|--clean]

set -e

SESSION_DIR="work_in_progress/current_session"
CONTEXT_FILE="${SESSION_DIR}/.compaction_context.md"
TEMP_TASK_DIR="${SESSION_DIR}/.temp_task"

show_help() {
    cat << EOF
🎯 preserve_context.sh - Preserva contexto durante compaction do opencode

⚠️⚠️⚠️ IMPORTANTE: O AGENTE OPENCODE NÃO PODE PREVER QUANDO O COMPACTION VAI ACONTECER ⚠️⚠️⚠️

O compaction do opencode é IMPREVISÍVEL e faz perder o contexto da tarefa atual.
Este script deve ser executado pelo USUÁRIO quando detectar que o agente entrou em compaction.

USO:
  ./preserve_context.sh --save      # Usuário executa quando detecta compaction
  ./preserve_context.sh --restore   # Usuário executa após compaction para restaurar
  ./preserve_context.sh --clean     # Limpa arquivos temporários
  ./preserve_context.sh --status    # Mostra status do contexto

⚠️ FLUXO CORRETO (EXECUTADO PELO USUÁRIO):
  1. Usuário detecta que opencode entrou em modo compaction
  2. Usuário executa: ./preserve_context.sh --save
  3. opencode continua em compaction (perde contexto)
  4. Após compaction, usuário executa: ./preserve_context.sh --restore
  5. opencode lê ${CONTEXT_FILE} e continua de onde parou

❌ FLUXO INCORRETO (NÃO FUNCIONA):
  Agente tenta executar --save antes do compaction ← IMPOSSÍVEL (não pode prever)

ARQUIVOS CRIADOS:
  - ${CONTEXT_FILE}: Contexto da tarefa (o que foi feito, próximos passos)
  - ${TEMP_TASK_DIR}/: Cópia temporária dos arquivos da tarefa
EOF
    exit 0
}

save_context() {
    echo "💾 SALVANDO CONTEXTO (EXECUTADO PELO USUÁRIO)..."
    echo "================================================"
    echo "⚠️  Este comando deve ser executado pelo USUÁRIO quando detectar"
    echo "    que o agente opencode entrou em modo compaction."
    echo ""
    
    # 1. Verificar se há sessão ativa
    if [ ! -d "${SESSION_DIR}" ]; then
        echo "❌ Nenhuma sessão ativa encontrada."
        echo "💡 Execute primeiro: ./start_session.sh"
        exit 1
    fi
    
    # 2. Verificar se há tarefa ativa
    TASK_COUNT=$(find work_in_progress/tasks -maxdepth 1 -type d -name "task_*" 2>/dev/null | wc -l)
    if [ "$TASK_COUNT" -eq 0 ]; then
        echo "⚠️  Nenhuma tarefa ativa encontrada."
        echo "💡 Crie uma tarefa primeiro: ./create_task.sh \"Nome da Tarefa\""
        exit 1
    fi
    
    # 3. Encontrar a tarefa mais recente
    LATEST_TASK=$(find work_in_progress/tasks -maxdepth 1 -type d -name "task_*" -exec ls -dt {} + | head -1)
    TASK_ID=$(basename "${LATEST_TASK}" | sed 's/task_//')
    
    if [ -z "${TASK_ID}" ]; then
        echo "❌ Não foi possível identificar a tarefa atual."
        exit 1
    fi
    
    echo "📋 Tarefa identificada: ${TASK_ID}"
    echo "📁 Diretório: ${LATEST_TASK}"
    
    # 4. Carregar metadados da tarefa
    if [ -f "${LATEST_TASK}/task_metadata" ]; then
        source "${LATEST_TASK}/task_metadata"
        echo "📝 Nome da tarefa: ${TASK_NAME}"
        echo "🏷️  Tipo: ${TASK_TYPE}"
        echo "📦 Módulo: ${MODULE}"
    else
        echo "⚠️  Metadados da tarefa não encontrados."
        TASK_NAME="Tarefa ${TASK_ID}"
        TASK_TYPE="Feature"
        MODULE="ui_web"
    fi
    
    # 5. Criar diretório temporário
    mkdir -p "${TEMP_TASK_DIR}"
    
    # 6. Copiar arquivos importantes da tarefa
    echo "📦 Copiando arquivos da tarefa..."
    cp -r "${LATEST_TASK}"/* "${TEMP_TASK_DIR}/" 2>/dev/null || true
    
    # 7. Criar arquivo de contexto
    echo "📝 Criando arquivo de contexto..."
    
    # Coletar informações do que foi feito
    IMPLEMENTED_FILES=$(find modules -name "*.go" -newer "${LATEST_TASK}/task_prompt.md" 2>/dev/null | head -10 || echo "Nenhum arquivo novo encontrado")
    
    cat > "${CONTEXT_FILE}" << EOF
# 🔄 CONTEXTO PRESERVADO PARA RECUPERAÇÃO APÓS COMPACTION

## 📋 INFORMAÇÕES DA TAREFA
- **ID da Tarefa:** ${TASK_ID}
- **Nome:** ${TASK_NAME}
- **Tipo:** ${TASK_TYPE}
- **Módulo:** ${MODULE}
- **Data/Hora do save:** $(date +"%d/%m/%Y %H:%M:%S")

## 📁 ARQUIVOS DA TAREFA (cópia em ${TEMP_TASK_DIR}/)
- task_prompt.md: $(wc -l < "${TEMP_TASK_DIR}/task_prompt.md" 2>/dev/null || echo "0") linhas
- checklist.md: $(wc -l < "${TEMP_TASK_DIR}/checklist.md" 2>/dev/null || echo "0") linhas
- implementation_plan.md: $(wc -l < "${TEMP_TASK_DIR}/implementation_plan.md" 2>/dev/null || echo "0") linhas

## 🎯 O QUE JÁ FOI IMPLEMENTADO
\`\`\`
${IMPLEMENTED_FILES}
\`\`\`

## 📝 RESUMO DO TRABALHO ATÉ AGORA
$(cat "${TEMP_TASK_DIR}/task_prompt.md" 2>/dev/null | head -20 || echo "Não foi possível ler o prompt da tarefa")

## 🔄 PRÓXIMOS PASSOS (RECOMENDADO APÓS RESTAURAÇÃO)
1. Ler este arquivo para entender o contexto
2. Verificar arquivos em ${TEMP_TASK_DIR}/
3. Continuar implementação de onde parou
4. Usar ./conclude_task.sh quando terminar

## 🤖 INSTRUÇÃO PARA OPENCODE APÓS RESTAURAÇÃO
1. Leia este arquivo completamente
2. Verifique os arquivos em ${TEMP_TASK_DIR}/
3. Continue a implementação da tarefa ${TASK_ID}: ${TASK_NAME}
4. Quando terminar, execute: ./conclude_task.sh --task=${TASK_ID} "Aprendizados: ..."

## ⚠️ ALERTA DE COMPACTION
Este arquivo foi criado porque o opencode entrou em modo compaction e perdeu o contexto.
Use estas informações para continuar de onde parou.

EOF
    
    echo "✅ Contexto salvo com sucesso!"
    echo ""
    echo "📊 RESUMO:"
    echo "----------"
    echo "Arquivo de contexto: ${CONTEXT_FILE}"
    echo "Cópia da tarefa: ${TEMP_TASK_DIR}/"
    echo "Tarefa ID: ${TASK_ID}"
    echo ""
    echo "💡 PRÓXIMOS PASSOS PARA O USUÁRIO:"
    echo "1. opencode está em modo compaction (perdeu contexto)"
    echo "2. Aguarde compaction terminar"
    echo "3. Execute: ./preserve_context.sh --restore"
    echo "4. opencode lerá ${CONTEXT_FILE} e continuará de onde parou"
    echo ""
    echo "⚠️  LEMBRETE: O agente opencode NÃO pode prever quando o compaction acontece."
    echo "    Este script deve ser executado pelo USUÁRIO quando detectar compaction."
}

restore_context() {
    echo "🔄 RESTAURANDO CONTEXTO APÓS COMPACTION (EXECUTADO PELO USUÁRIO)..."
    echo "=================================================================="
    echo "⚠️  Este comando deve ser executado pelo USUÁRIO após o compaction terminar."
    
    # 1. Verificar se há arquivo de contexto
    if [ ! -f "${CONTEXT_FILE}" ]; then
        echo "❌ Arquivo de contexto não encontrado: ${CONTEXT_FILE}"
        echo "💡 Execute primeiro: ./preserve_context.sh --save"
        exit 1
    fi
    
    # 2. Ler informações do contexto
    TASK_ID=$(grep "ID da Tarefa:" "${CONTEXT_FILE}" | cut -d: -f2 | xargs)
    TASK_NAME=$(grep "Nome:" "${CONTEXT_FILE}" | cut -d: -f2 | xargs)
    
    if [ -z "${TASK_ID}" ]; then
        echo "❌ Não foi possível extrair ID da tarefa do contexto."
        exit 1
    fi
    
    echo "📋 Restaurando tarefa: ${TASK_ID}"
    echo "📝 Nome: ${TASK_NAME}"
    
    # 3. Verificar se a tarefa ainda existe
    TASK_DIR="work_in_progress/tasks/task_${TASK_ID}"
    if [ -d "${TASK_DIR}" ]; then
        echo "✅ Tarefa ainda existe em: ${TASK_DIR}"
        echo "💡 Contexto já está disponível."
    else
        echo "⚠️  Tarefa não encontrada em work_in_progress/tasks/"
        echo "💡 Restaurando da cópia temporária..."
        
        # Restaurar da cópia temporária
        if [ -d "${TEMP_TASK_DIR}" ]; then
            mkdir -p "work_in_progress/tasks"
            cp -r "${TEMP_TASK_DIR}" "${TASK_DIR}"
            echo "✅ Tarefa restaurada para: ${TASK_DIR}"
        else
            echo "❌ Cópia temporária não encontrada: ${TEMP_TASK_DIR}"
            echo "💡 Crie uma nova tarefa: ./create_task.sh \"${TASK_NAME}\""
        fi
    fi
    
    # 4. Mostrar instruções para opencode
    echo ""
    echo "🤖 INSTRUÇÕES PARA OPENCODE:"
    echo "============================"
    echo "1. Leia o arquivo de contexto: ${CONTEXT_FILE}"
    echo "2. Verifique a tarefa em: ${TASK_DIR}"
    echo "3. Continue a implementação de onde parou"
    echo "4. Quando terminar: ./conclude_task.sh --task=${TASK_ID} \"Aprendizados: ...\""
    echo ""
    echo "📝 CONTEÚDO DO ARQUIVO DE CONTEXTO (primeiras 20 linhas):"
    echo "--------------------------------------------------------"
    head -20 "${CONTEXT_FILE}"
}

clean_context() {
    echo "🧹 LIMPANDO ARQUIVOS DE CONTEXTO TEMPORÁRIOS..."
    echo "=============================================="
    
    rm -f "${CONTEXT_FILE}" 2>/dev/null || true
    rm -rf "${TEMP_TASK_DIR}" 2>/dev/null || true
    
    echo "✅ Arquivos temporários removidos:"
    echo "   - ${CONTEXT_FILE}"
    echo "   - ${TEMP_TASK_DIR}/"
}

show_status() {
    echo "📊 STATUS DO CONTEXTO DE COMPACTION"
    echo "==================================="
    
    if [ -f "${CONTEXT_FILE}" ]; then
        echo "✅ Arquivo de contexto encontrado: ${CONTEXT_FILE}"
        echo "📝 Conteúdo (resumo):"
        grep -E "(ID da Tarefa:|Nome:|Data/Hora do save:)" "${CONTEXT_FILE}" | head -5
    else
        echo "❌ Nenhum arquivo de contexto encontrado."
    fi
    
    if [ -d "${TEMP_TASK_DIR}" ]; then
        echo "✅ Cópia temporária da tarefa encontrada: ${TEMP_TASK_DIR}/"
        echo "📁 Conteúdo:"
        ls -la "${TEMP_TASK_DIR}/" 2>/dev/null | head -10 || echo "   (vazio)"
    else
        echo "❌ Nenhuma cópia temporária encontrada."
    fi
    
    echo ""
    echo "💡 Use:"
    echo "  ./preserve_context.sh --save      # Salvar contexto antes do compaction"
    echo "  ./preserve_context.sh --restore   # Restaurar contexto após o compaction"
    echo "  ./preserve_context.sh --clean     # Limpar arquivos temporários"
}

# Processar argumentos
case "$1" in
    --save|-s)
        save_context
        ;;
    --restore|-r)
        restore_context
        ;;
    --clean|-c)
        clean_context
        ;;
    --status|-st)
        show_status
        ;;
    --help|-h)
        show_help
        ;;
    *)
        echo "❌ Argumento inválido: $1"
        echo ""
        show_help
        exit 1
        ;;
esac