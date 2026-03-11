#!/bin/bash
# start_session.sh - Inicializa uma sessão no projeto Digna
# Uso: ./start_session.sh [opcional: "quick" para modo rápido]

set -e  # Sai no primeiro erro

echo "🚀 Iniciando sessão no projeto Digna..."
echo "=========================================="

# Configurações
SESSION_ID=$(date +%Y%m%d_%H%M%S)
SESSION_DIR="work_in_progress/current_session"
SESSION_INFO="${SESSION_DIR}/session_info"
QUICK_MODE="${1:-no}"

# Limpar sessão anterior se existir
if [ -d "${SESSION_DIR}" ]; then
    echo "⚠️  Sessão anterior encontrada, arquivando..."
    
    # Verificar se há tarefas não concluídas
    if [ -d "work_in_progress/tasks" ]; then
        PENDING_TASKS=$(find work_in_progress/tasks -name "task_metadata" -exec grep -l "STATUS=\"pending\"" {} \; 2>/dev/null | wc -l)
        if [ "$PENDING_TASKS" -gt 0 ]; then
            echo "❌ Há ${PENDING_TASKS} tarefa(s) pendente(s). Conclua-as primeiro."
            echo "💡 Use: ./conclude_task.sh para cada tarefa pendente"
            exit 1
        fi
    fi
    
    # Arquivar sessão anterior
    ARCHIVE_DIR="work_in_progress/archive/session_$(cat ${SESSION_INFO} 2>/dev/null | grep "SESSION_ID=" | cut -d= -f2 || echo "unknown")"
    mkdir -p "${ARCHIVE_DIR}"
    
    # Mover arquivos da sessão
    if [ -f "${SESSION_INFO}" ]; then
        cp "${SESSION_INFO}" "${ARCHIVE_DIR}/"
    fi
    
    # Mover aprendizados da sessão
    if [ -d "${SESSION_DIR}/session_learnings" ]; then
        cp -r "${SESSION_DIR}/session_learnings" "${ARCHIVE_DIR}/" 2>/dev/null || true
    fi
    
    # Mover tarefas concluídas
    if [ -d "work_in_progress/tasks" ]; then
        mkdir -p "${ARCHIVE_DIR}/tasks"
        find work_in_progress/tasks -type d -name "task_*" -exec mv {} "${ARCHIVE_DIR}/tasks/" \; 2>/dev/null || true
    fi
    
    # Limpar diretório atual
    rm -rf "${SESSION_DIR}"/*
    rm -rf "work_in_progress/tasks" 2>/dev/null || true
    
    echo "✅ Sessão anterior arquivada em: ${ARCHIVE_DIR}"
fi

# Criar estrutura de diretórios
mkdir -p "${SESSION_DIR}"
mkdir -p "${SESSION_DIR}/session_learnings"
mkdir -p "work_in_progress/tasks"
mkdir -p "work_in_progress/archive"

# 1. Criar arquivo de sessão
echo "SESSION_ID=${SESSION_ID}" > ${SESSION_INFO}
echo "START_TIME=$(date +%s)" >> ${SESSION_INFO}
echo "QUICK_MODE=${QUICK_MODE}" >> ${SESSION_INFO}

echo "✅ Sessão criada: ${SESSION_ID}"
echo "📁 Diretório: ${SESSION_DIR}"

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
            echo "# 🚀 QUICK REFERENCE - Projeto Digna" > docs/QUICK_REFERENCE.md
            echo "" >> docs/QUICK_REFERENCE.md
            echo "**Última atualização:** $(date +%d/%m/%Y)" >> docs/QUICK_REFERENCE.md
            echo "" >> docs/QUICK_REFERENCE.md
            echo "## 📋 Handlers Registrados" >> docs/QUICK_REFERENCE.md
            echo "" >> docs/QUICK_REFERENCE.md
            echo "Nenhum handler registrado ainda." >> docs/QUICK_REFERENCE.md
        fi
    fi
fi

# 3. Copiar/criar .agent_context.md na sessão atual
if [ -f ".agent_context.md" ]; then
    cp ".agent_context.md" "${SESSION_DIR}/.agent_context.md"
    echo "✅ Contexto do agente copiado para a sessão"
else
    echo "⚠️  .agent_context.md não encontrado, criando novo..."
    
    # Criar contexto básico
    cat > "${SESSION_DIR}/.agent_context.md" << EOF
# 🎯 CONTEXTO DO AGENTE - Projeto Digna

**Sessão iniciada:** $(date +%d/%m/%Y %H:%M:%S)
**Sessão ID:** ${SESSION_ID}
**Arquivo:** work_in_progress/current_session/.agent_context.md

---

## 🚀 INSTRUÇÕES PARA O AGENTE (OPENCODE)

Você está trabalhando no **Projeto Digna** - sistema de economia solidária.

### 📁 ESTRUTURA DO PROJETO
\`\`\`
work_in_progress/
├── current_session/          # ✅ Você está aqui
│   ├── .agent_context.md     # Este arquivo
│   ├── session_info          # Metadados da sessão
│   └── session_learnings/    # Aprendizados coletados
├── tasks/                    # Tarefas em andamento
│   └── task_[ID]/           # Cada tarefa tem seu diretório
└── archive/                 # Sessões e tarefas concluídas
\`\`\`

### 🔄 FLUXO DE TRABALHO
1. \`./start_session.sh\`              → Inicia sessão (já feito)
2. \`./create_task.sh "Nome"\`         → Cria nova tarefa
3. \`./process_task.sh --task=[ID]\`   → Processa tarefa
4. Implementar                       → Você (opencode) implementa
5. \`./conclude_task.sh\`              → Conclui tarefa
6. \`./end_session.sh\`                → Encerra sessão (após TODAS tarefas)

### 📚 DOCUMENTAÇÃO IMPORTANTE
- \`docs/QUICK_REFERENCE.md\` - Arquitetura e padrões
- \`docs/ANTIPATTERNS.md\` - O que NÃO fazer
- \`docs/NEXT_STEPS.md\` - Tarefas pendentes
- \`docs/learnings/\` - Aprendizados anteriores

---

**Status da sessão:** ATIVA
**Tarefas ativas:** 0
**Última atualização:** $(date +%d/%m/%Y %H:%M:%S)
EOF
fi

# 4. Criar link simbólico para .agent_context.md no root (para compatibilidade)
ln -sf "${SESSION_DIR}/.agent_context.md" ".agent_context.md" 2>/dev/null || true

echo ""
echo "✅ SESSÃO INICIADA COM SUCESSO!"
echo "================================"
echo "ID: ${SESSION_ID}"
echo "Modo: ${QUICK_MODE}"
echo ""
echo "📋 PRÓXIMOS PASSOS:"
echo "1. Criar tarefa: ./create_task.sh \"Nome da Tarefa\" [módulo]"
echo "2. Listar tarefas: ls work_in_progress/tasks/"
echo "3. Processar tarefa: ./process_task.sh --task=[ID] --checklist"
echo ""
echo "💡 Dica: Use modo rápido para sessões curtas: ./start_session.sh quick"