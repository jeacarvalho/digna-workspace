#!/bin/bash
# migrate_to_new_structure.sh - Migra da estrutura antiga para a nova
# Uso: ./migrate_to_new_structure.sh

set -e

echo "🔄 Migrando para nova estrutura de diretórios..."
echo "================================================"

# 1. Criar estrutura básica se não existir
mkdir -p work_in_progress/{current_session,tasks,archive}
mkdir -p work_in_progress/current_session/session_learnings
mkdir -p work_in_progress/task_template

echo "✅ Estrutura de diretórios criada"

# 2. Copiar template se não existir
if [ ! -f "work_in_progress/task_template/task_prompt.md" ]; then
    if [ -f "scripts/workflow/task_prompt_template.md" ]; then
        cp "scripts/workflow/task_prompt_template.md" "work_in_progress/task_template/task_prompt.md"
    else
        # Criar template básico
        cat > "work_in_progress/task_template/task_prompt.md" << EOF
# 📋 TAREFA: [NOME_DA_TAREFA]

**Data:** [DATA]
**Prioridade:** [ALTA/MÉDIA/BAIXA]
**Estimativa:** [X] minutos/horas
**Módulo(s):** [módulo1, módulo2, ...]

---

## 🎯 OBJETIVO

[Descrição clara do que precisa ser implementado/alterado/corrigido]

---

## 📋 REQUISITOS

### Funcionais
- [ ] Requisito 1
- [ ] Requisito 2
- [ ] Requisito 3

### Técnicos
- [ ] Seguir padrões do projeto Digna
- [ ] Implementar testes unitários
- [ ] Atualizar documentação
- [ ] Validar com smoke tests

---

**Status:** PENDENTE
**Última atualização:** [DATA]
EOF
    fi
    echo "✅ Template de prompt criado"
fi

# 3. Migrar sessões existentes
echo ""
echo "📦 MIGRANDO SESSÕES EXISTENTES..."
echo "================================"

SESSION_FILES=$(ls .session_* 2>/dev/null | head -5)
if [ -n "$SESSION_FILES" ]; then
    for SESSION_FILE in $SESSION_FILES; do
        SESSION_ID=$(echo "$SESSION_FILE" | sed 's/\.session_//')
        echo "  Migrando sessão: ${SESSION_ID}"
        
        ARCHIVE_DIR="work_in_progress/archive/session_${SESSION_ID}"
        mkdir -p "${ARCHIVE_DIR}"
        
        # Copiar arquivo de sessão
        cp "${SESSION_FILE}" "${ARCHIVE_DIR}/session_info" 2>/dev/null || true
        
        # Tentar encontrar tarefas relacionadas
        TASK_FILES=$(ls .task_* 2>/dev/null | grep -v "summary" | head -10)
        for TASK_FILE in $TASK_FILES; do
            TASK_ID=$(echo "$TASK_FILE" | sed 's/\.task_//')
            
            # Verificar se a tarefa pertence a esta sessão (por timestamp)
            if [[ "$TASK_ID" == "${SESSION_ID}"* ]] || [[ "$SESSION_ID" == "${TASK_ID:0:8}"* ]]; then
                echo "    Migrando tarefa: ${TASK_ID}"
                
                TASK_DIR="${ARCHIVE_DIR}/tasks/task_${TASK_ID}"
                mkdir -p "${TASK_DIR}"
                
                # Copiar arquivo de tarefa como metadados
                cp "${TASK_FILE}" "${TASK_DIR}/task_metadata" 2>/dev/null || true
                
                # Tentar extrair nome da tarefa
                TASK_NAME=$(grep "TASK_DESCRIPTION=" "${TASK_FILE}" 2>/dev/null | cut -d= -f2 | tr -d '"' || echo "Tarefa ${TASK_ID}")
                
                # Criar prompt básico
                cat > "${TASK_DIR}/task_prompt.md" << EOF
# 📋 TAREFA: ${TASK_NAME}

**Data:** $(date +%d/%m/%Y)
**ID Original:** ${TASK_ID}
**Status:** MIGRADO

---

## 🎯 OBJETIVO

Tarefa migrada da estrutura antiga.

**Descrição original:** ${TASK_NAME}

---

## 📋 NOTAS

Esta tarefa foi migrada automaticamente da estrutura antiga.
Os detalhes originais podem estar no arquivo \`task_metadata\`.

---

**Status:** CONCLUÍDO (migrado)
**Última atualização:** $(date +%d/%m/%Y)
EOF
                
                # Marcar como concluída
                echo "TASK_ID=${TASK_ID}" > "${TASK_DIR}/task_metadata"
                echo "TASK_NAME=\"${TASK_NAME}\"" >> "${TASK_DIR}/task_metadata"
                echo "STATUS=\"completed\"" >> "${TASK_DIR}/task_metadata"
                echo "MIGRATED=true" >> "${TASK_DIR}/task_metadata"
            fi
        done
    done
    echo "✅ Sessões migradas"
else
    echo "ℹ️  Nenhuma sessão encontrada para migrar"
fi

# 4. Migrar arquivos de tarefas soltas
echo ""
echo "📋 MIGRANDO ARQUIVOS DE TAREFAS SOLTAS..."
echo "========================================"

# Arquivos .task_summary_
SUMMARY_FILES=$(ls .task_summary_* 2>/dev/null | head -5)
if [ -n "$SUMMARY_FILES" ]; then
    mkdir -p "work_in_progress/archive/legacy_summaries"
    for SUMMARY_FILE in $SUMMARY_FILES; do
        cp "${SUMMARY_FILE}" "work_in_progress/archive/legacy_summaries/" 2>/dev/null || true
        echo "  Migrado: ${SUMMARY_FILE}"
    done
    echo "✅ Sumários migrados"
fi

# Arquivos .opencode_task_
OPENCODE_FILES=$(ls .opencode_task_* 2>/dev/null | head -5)
if [ -n "$OPENCODE_FILES" ]; then
    mkdir -p "work_in_progress/archive/legacy_opencode"
    for OPENCODE_FILE in $OPENCODE_FILES; do
        cp "${OPENCODE_FILE}" "work_in_progress/archive/legacy_opencode/" 2>/dev/null || true
        echo "  Migrado: ${OPENCODE_FILE}"
    done
    echo "✅ Arquivos opencode migrados"
fi

# 5. Migrar .agent_context.md se existir
if [ -f ".agent_context.md" ]; then
    echo ""
    echo "🤖 MIGRANDO CONTEXTO DO AGENTE..."
    echo "================================"
    
    # Criar sessão atual com o contexto
    SESSION_ID="migrated_$(date +%Y%m%d_%H%M%S)"
    SESSION_DIR="work_in_progress/current_session"
    
    mkdir -p "${SESSION_DIR}"
    cp ".agent_context.md" "${SESSION_DIR}/.agent_context.md"
    
    # Criar info da sessão
    echo "SESSION_ID=${SESSION_ID}" > "${SESSION_DIR}/session_info"
    echo "START_TIME=$(date +%s)" >> "${SESSION_DIR}/session_info"
    echo "QUICK_MODE=no" >> "${SESSION_DIR}/session_info"
    echo "MIGRATED=true" >> "${SESSION_DIR}/session_info"
    
    # Criar link simbólico para compatibilidade
    ln -sf "${SESSION_DIR}/.agent_context.md" ".agent_context.md" 2>/dev/null || true
    
    echo "✅ Contexto do agente migrado para nova sessão: ${SESSION_ID}"
fi

# 6. Limpar arquivos antigos (perguntar primeiro)
echo ""
echo "🧹 LIMPEZA DE ARQUIVOS ANTIGOS..."
echo "================================"

OLD_FILES_COUNT=$(ls .session_* .task_* .opencode_task_* 2>/dev/null | wc -l)
if [ "$OLD_FILES_COUNT" -gt 0 ]; then
    echo "Encontrados ${OLD_FILES_COUNT} arquivos da estrutura antiga:"
    ls .session_* .task_* .opencode_task_* 2>/dev/null | head -10
    
    read -p "Remover arquivos antigos? (s/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Ss]$ ]]; then
        rm -f .session_* .task_* .opencode_task_* 2>/dev/null || true
        echo "✅ Arquivos antigos removidos"
    else
        echo "ℹ️  Arquivos antigos mantidos"
    fi
else
    echo "ℹ️  Nenhum arquivo antigo encontrado"
fi

# 7. Atualizar scripts principais (fazer backup primeiro)
echo ""
echo "🔄 ATUALIZANDO SCRIPTS PRINCIPAIS..."
echo "==================================="

# Fazer backup dos scripts antigos
mkdir -p "scripts/workflow/backup_$(date +%Y%m%d)"
cp scripts/workflow/start_session.sh "scripts/workflow/backup_$(date +%Y%m%d)/" 2>/dev/null || true
cp scripts/workflow/process_task.sh "scripts/workflow/backup_$(date +%Y%m%d)/" 2>/dev/null || true
cp scripts/workflow/conclude_task.sh "scripts/workflow/backup_$(date +%Y%m%d)/" 2>/dev/null || true
cp scripts/workflow/end_session.sh "scripts/workflow/backup_$(date +%Y%m%d)/" 2>/dev/null || true

echo "✅ Backup criado em: scripts/workflow/backup_$(date +%Y%m%d)/"

# Substituir pelos novos scripts
if [ -f "scripts/workflow/start_session_new.sh" ]; then
    mv "scripts/workflow/start_session_new.sh" "scripts/workflow/start_session.sh"
    chmod +x "scripts/workflow/start_session.sh"
    echo "✅ start_session.sh atualizado"
fi

if [ -f "scripts/workflow/process_task_new.sh" ]; then
    mv "scripts/workflow/process_task_new.sh" "scripts/workflow/process_task.sh"
    chmod +x "scripts/workflow/process_task.sh"
    echo "✅ process_task.sh atualizado"
fi

if [ -f "scripts/workflow/conclude_task_new.sh" ]; then
    mv "scripts/workflow/conclude_task_new.sh" "scripts/workflow/conclude_task.sh"
    chmod +x "scripts/workflow/conclude_task.sh"
    echo "✅ conclude_task.sh atualizado"
fi

if [ -f "scripts/workflow/end_session_new.sh" ]; then
    mv "scripts/workflow/end_session_new.sh" "scripts/workflow/end_session.sh"
    chmod +x "scripts/workflow/end_session.sh"
    echo "✅ end_session.sh atualizado"
fi

# Garantir que create_task.sh seja executável
if [ -f "scripts/workflow/create_task.sh" ]; then
    chmod +x "scripts/workflow/create_task.sh"
    echo "✅ create_task.sh configurado"
fi

echo ""
echo "✅ MIGRAÇÃO CONCLUÍDA COM SUCESSO!"
echo "=================================="
echo ""
echo "📁 NOVA ESTRUTURA:"
echo "-----------------"
echo "work_in_progress/"
echo "├── current_session/     # Sessão atual"
echo "├── tasks/              # Tarefas em andamento"
echo "├── archive/            # Sessões e tarefas concluídas"
echo "└── task_template/      # Template para novas tarefas"
echo ""
echo "🚀 NOVO FLUXO DE TRABALHO:"
echo "-------------------------"
echo "1. Iniciar sessão:      ./start_session.sh"
echo "2. Criar tarefa:        ./create_task.sh \"Nome da Tarefa\" [módulo]"
echo "3. Processar tarefa:    ./process_task.sh --task=ID [--checklist|--plan|--execute]"
echo "4. Concluir tarefa:     ./conclude_task.sh --task=ID \"Aprendizados\""
echo "5. Encerrar sessão:     ./end_session.sh"
echo ""
echo "💡 Dica: Teste o novo fluxo criando uma tarefa de teste:"
echo "       ./create_task.sh \"Tarefa de Teste\" ui_web"
echo "       ./process_task.sh --task=[ID] --checklist"
echo ""
echo "🤖 PARA OPENCODE:"
echo "----------------"
echo "O contexto agora está em: work_in_progress/current_session/.agent_context.md"
echo "Use os novos scripts para o fluxo de trabalho."