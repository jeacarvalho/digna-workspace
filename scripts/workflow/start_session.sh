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
    echo "⚠️  Sessão anterior encontrada, verificando status..."
    
    # Verificar se documentação foi lida na sessão anterior
    if [ -f "${SESSION_DIR}/docs_checklist.md" ]; then
        # Usar uma abordagem mais robusta para contar checkboxes
        DOCS_STATUS=0
        TOTAL_DOCS=0
        
        if grep -q "\[x\]" "${SESSION_DIR}/docs_checklist.md" 2>/dev/null; then
            DOCS_STATUS=$(grep -c "\[x\]" "${SESSION_DIR}/docs_checklist.md")
        fi
        
        if grep -q "\[ \]" "${SESSION_DIR}/docs_checklist.md" 2>/dev/null; then
            TOTAL_DOCS=$(grep -c "\[ \]" "${SESSION_DIR}/docs_checklist.md")
        fi
        
        # Verificar se há documentação não lida (pelo menos 4 itens totais e menos de 4 marcados)
        if [ "$DOCS_STATUS" -lt 4 ] && [ "$((DOCS_STATUS + TOTAL_DOCS))" -ge 4 ]; then
            echo "❌❌❌ ALERTA CRÍTICO: DOCUMENTAÇÃO NÃO LIDA NA SESSÃO ANTERIOR ❌❌❌"
            echo ""
            echo "📊 STATUS DA SESSÃO ANTERIOR:"
            echo "  • Documentação lida: ${DOCS_STATUS}/4"
            echo "  • Documentação pendente: $((4 - DOCS_STATUS))"
            echo ""
            echo "⚠️  CONSEQUÊNCIAS:"
            echo "  • Implementações podem ter violado padrões"
            echo "  • Possível reintrodução de antipadrões"
            echo "  • Perda de eficiência significativa"
            echo ""
            echo "🎯 AÇÕES RECOMENDADAS:"
            echo "  1. Revise a sessão anterior: ${SESSION_DIR}/"
            echo "  2. Leia a documentação obrigatória AGORA"
            echo "  3. Verifique se há antipadrões introduzidos"
            echo ""
            echo "📋 Para continuar, confirme que leu a documentação:"
            read -p "✅ Digite 'SIM' para confirmar leitura da documentação: " CONFIRMACAO
            if [ "$CONFIRMACAO" != "SIM" ]; then
                echo "❌ Confirmação negada. Encerrando sessão."
                exit 1
            fi
            echo "✅ Confirmação recebida. Continuando..."
        fi
    fi
    
    echo "⚠️  Arquivando sessão anterior..."
    
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

# 3. Criar sistema de verificação de documentação (apenas modo completo)
if [ "$QUICK_MODE" != "quick" ]; then
    echo "📚 Criando sistema de verificação de documentação obrigatória..."
    
    # Criar arquivo de verificação de documentação
    DOCS_CHECKLIST="${SESSION_DIR}/docs_checklist.md"
    cat > "${DOCS_CHECKLIST}" << EOF
# 📚 CHECKLIST DE DOCUMENTAÇÃO OBRIGATÓRIA
# Sessão: ${SESSION_ID}
# Data: $(date +%d/%m/%Y)

## ✅ DOCUMENTAÇÃO QUE O AGENTE DEVE LER ANTES DE QUALQUER AÇÃO

### 📋 CHECKLIST DE LEITURA (MARQUE ✅ APÓS LER)
- [ ] **docs/QUICK_REFERENCE.md** - Arquitetura, padrões e handlers
- [ ] **docs/ANTIPATTERNS.md** - O que NÃO fazer no projeto  
- [ ] **docs/NEXT_STEPS.md** - Tarefas pendentes e prioridades
- [ ] **docs/learnings/** - TODOS os aprendizados anteriores (CRÍTICO)

### 📊 STATUS DE LEITURA
**Data da verificação:** 
**Agente:** opencode
**Status:** PENDENTE (leia AGORA)

### 🎯 INSTRUÇÕES
1. Leia TODOS os documentos acima
2. Marque cada item com ✅ após ler
3. Atualize o status para "CONCLUÍDO"
4. Só então prossiga com implementações

### ⚠️ CONSEQUÊNCIAS DE NÃO LER
- Implementações podem violar padrões do projeto
- Pode recriar funcionalidades já existentes
- Pode introduzir antipadrões documentados
- Perda de tempo significativa (40min+ por sessão)

EOF

    echo "✅ Sistema de verificação de documentação criado"
else
    echo "⚡ Modo rápido: sistema de verificação de documentação desativado"
    DOCS_CHECKLIST=""
fi

# 3. Criar sistema de verificação de documentação
echo "📚 Criando sistema de verificação de documentação obrigatória..."

# Criar arquivo de verificação de documentação
DOCS_CHECKLIST="${SESSION_DIR}/docs_checklist.md"
cat > "${DOCS_CHECKLIST}" << EOF
# 📚 CHECKLIST DE DOCUMENTAÇÃO OBRIGATÓRIA
# Sessão: ${SESSION_ID}
# Data: $(date +%d/%m/%Y)

## ✅ DOCUMENTAÇÃO QUE O AGENTE DEVE LER ANTES DE QUALQUER AÇÃO

### 📋 CHECKLIST DE LEITURA (MARQUE ✅ APÓS LER)
- [ ] **docs/QUICK_REFERENCE.md** - Arquitetura, padrões e handlers
- [ ] **docs/ANTIPATTERNS.md** - O que NÃO fazer no projeto  
- [ ] **docs/NEXT_STEPS.md** - Tarefas pendentes e prioridades
- [ ] **docs/learnings/** - TODOS os aprendizados anteriores (CRÍTICO)

### 📊 STATUS DE LEITURA
**Data da verificação:** 
**Agente:** opencode
**Status:** PENDENTE (leia AGORA)

### 🎯 INSTRUÇÕES
1. Leia TODOS os documentos acima
2. Marque cada item com ✅ após ler
3. Atualize o status para "CONCLUÍDO"
4. Só então prossiga com implementações

### ⚠️ CONSEQUÊNCIAS DE NÃO LER
- Implementações podem violar padrões do projeto
- Pode recriar funcionalidades já existentes
- Pode introduzir antipadrões documentados
- Perda de tempo significativa (40min+ por sessão)

EOF

# 4. Copiar/criar .agent_context.md na sessão atual
if [ "$QUICK_MODE" != "quick" ]; then
    echo "📝 Criando contexto do agente com sistema de verificação..."
    
    # Criar contexto com sistema de verificação (modo completo)
    cat > "${SESSION_DIR}/.agent_context.md" << EOF
# 🎯 CONTEXTO DO AGENTE - Projeto Digna

**Sessão iniciada:** $(date +%d/%m/%Y %H:%M:%S)
**Sessão ID:** ${SESSION_ID}
**Arquivo:** work_in_progress/current_session/.agent_context.md
**Verificação docs:** ${DOCS_CHECKLIST}
**Modo:** COMPLETO (com bloqueio)

---

## 🚀 INSTRUÇÕES OBRIGATÓRIAS PARA O AGENTE (OPENCODE)

Você está trabalhando no **Projeto Digna** - sistema de economia solidária.

### ⚠️ BLOQUEIO DE IMPLEMENTAÇÃO ATIVADO
**O AGENTE ESTÁ BLOQUEADO PARA IMPLEMENTAÇÕES ATÉ CONFIRMAR LEITURA DA DOCUMENTAÇÃO.**

### ✅ PRIMEIRO PASSO OBRIGATÓRIO (BLOQUEANTE)
1. 📚 **LEIA** a documentação obrigatória: \`${DOCS_CHECKLIST}\`
2. ✅ **MARQUE** cada item como lido no checklist
3. 🔓 **SÓ ENTÃO** o agente está liberado para implementar

### 📚 DOCUMENTAÇÃO OBRIGATÓRIA (BLOQUEANTE)
**LEIA AGORA ANTES DE QUALQUER AÇÃO:**
1. \`docs/QUICK_REFERENCE.md\` - Arquitetura, padrões e handlers
2. \`docs/ANTIPATTERNS.md\` - O que NÃO fazer no projeto
3. \`docs/NEXT_STEPS.md\` - Tarefas pendentes e prioridades
4. \`docs/learnings/\` - TODOS os aprendizados anteriores (CRÍTICO)

### 📁 ESTRUTURA DO PROJETO
\`\`\`
work_in_progress/
├── current_session/          # ✅ Você está aqui
│   ├── .agent_context.md     # Este arquivo (instruções obrigatórias)
│   ├── docs_checklist.md     # ✅ Verificação de documentação (BLOQUEANTE)
│   ├── session_info          # Metadados da sessão
│   └── session_learnings/    # Aprendizados coletados (OBRIGATÓRIO)
├── tasks/                    # Tarefas em andamento
│   └── task_[ID]/           # Cada tarefa tem seu diretório
└── archive/                 # Sessões e tarefas concluídas
\`\`\`

### 🔄 FLUXO DE TRABALHO (OBRIGATÓRIO)
1. ✅ \`./start_session.sh\`              → Sessão iniciada
2. ⚠️  **VERIFIQUE BLOQUEIO**            → Sistema ativo: docs_checklist.md
3. 📚 **LEIA DOCUMENTAÇÃO OBRIGATÓRIA**  → BLOQUEANTE: leia AGORA
4. ✅ **MARQUE CHECKLIST**               → Confirme leitura
5. 🔓 **BLOQUEIO REMOVIDO**              → Agora pode implementar
6. 📋 \`./create_task.sh "Nome"\`         → Cria nova tarefa
7. 🔍 \`./process_task.sh --task=[ID]\`   → Processa tarefa com checklist
8. 💻 **Implementar**                    → Você (opencode) implementa
9. 📝 **Registrar aprendizado**          → Em \`session_learnings/\`
10. ✅ \`./conclude_task.sh\`              → Conclui tarefa
11. 🔚 \`./end_session.sh\`                → Encerra sessão (após TODAS tarefas)

### 🎯 PRÓXIMOS PASSOS (BLOQUEADO ATÉ LEITURA)
**IMPLEMENTAÇÕES BLOQUEADAS ATÉ CONFIRMAR LEITURA:**
1. 📚 **LEIA E MARQUE** \`${DOCS_CHECKLIST}\` (BLOQUEANTE)
2. 📋 **CRIE** tarefa: \`./create_task.sh "Nome da Tarefa" [módulo]\`
3. 🔍 **PROCESSE** tarefa: \`./process_task.sh --task=[ID] --checklist\`

---

**Status da sessão:** ATIVA ⚠️ BLOQUEADA
**Tarefas ativas:** 0
**Documentação lida:** NÃO (BLOQUEANTE - LEIA AGORA)
**Verificação:** ${DOCS_CHECKLIST}
**Última atualização:** $(date +%d/%m/%Y %H:%M:%S)

**⚠️  ALERTA:** O agente está BLOQUEADO para implementações até marcar a documentação como lida.
**💡 Para sessões rápidas sem bloqueio:** \`./start_session.sh quick\`
EOF

    echo "✅ Contexto do agente criado com bloqueio ativo"
else
    echo "📝 Criando contexto do agente (modo rápido)..."
    
    # Criar contexto sem sistema de verificação (modo rápido)
    cat > "${SESSION_DIR}/.agent_context.md" << EOF
# 🎯 CONTEXTO DO AGENTE - Projeto Digna

**Sessão iniciada:** $(date +%d/%m/%Y %H:%M:%S)
**Sessão ID:** ${SESSION_ID}
**Arquivo:** work_in_progress/current_session/.agent_context.md
**Modo:** RÁPIDO (sem bloqueio)
**Aviso:** Documentação não verificada - risco de antipadrões

---

## 🚀 INSTRUÇÕES PARA O AGENTE (OPENCODE)

Você está trabalhando no **Projeto Digna** - sistema de economia solidária.

### ⚠️ MODO RÁPIDO - AVISO DE RISCO
**Esta sessão NÃO verificou leitura da documentação obrigatória.**
**Risco aumentado de:** violar padrões, recriar funcionalidades, introduzir antipadrões.

### 📚 DOCUMENTAÇÃO RECOMENDADA (NÃO BLOQUEANTE)
**Considere ler antes de implementar:**
1. \`docs/QUICK_REFERENCE.md\` - Arquitetura, padrões e handlers
2. \`docs/ANTIPATTERNS.md\` - O que NÃO fazer no projeto
3. \`docs/NEXT_STEPS.md\` - Tarefas pendentes e prioridades

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
1. ✅ \`./start_session.sh quick\`        → Sessão rápida iniciada
2. 📋 \`./create_task.sh "Nome"\`         → Cria nova tarefa
3. 🔍 \`./process_task.sh --task=[ID]\`   → Processa tarefa com checklist
4. 💻 **Implementar**                    → Você (opencode) implementa
5. 📝 **Registrar aprendizado**          → Em \`session_learnings/\`
6. ✅ \`./conclude_task.sh\`              → Conclui tarefa

### 🎯 PRÓXIMOS PASSOS
1. 📋 **CRIE** tarefa: \`./create_task.sh "Nome da Tarefa" [módulo]\`
2. 🔍 **PROCESSE** tarefa: \`./process_task.sh --task=[ID] --checklist\`

---

**Status da sessão:** ATIVA ⚡ RÁPIDA
**Tarefas ativas:** 0
**Documentação verificada:** NÃO (modo rápido)
**Risco:** ALTO (documentação não verificada)
**Última atualização:** $(date +%d/%m/%Y %H:%M:%S)

**⚠️  AVISO:** Modo rápido não verifica documentação. Use apenas para tarefas simples.
**✅ Para sessões completas com verificação:** \`./start_session.sh\` (sem "quick")
EOF

    echo "✅ Contexto do agente criado (modo rápido)"
fi

# 4. Criar link simbólico para .agent_context.md no root (para compatibilidade)
ln -sf "${SESSION_DIR}/.agent_context.md" ".agent_context.md" 2>/dev/null || true

echo ""
if [ "$QUICK_MODE" != "quick" ]; then
    echo "✅ SESSÃO INICIADA COM SUCESSO! ⚠️  BLOQUEIO ATIVO"
    echo "=================================================="
    echo "ID: ${SESSION_ID}"
    echo "Modo: ${QUICK_MODE} (COMPLETO)"
    echo "Status: ⚠️  IMPLEMENTAÇÕES BLOQUEADAS"
    echo ""
    echo "❌❌❌ SISTEMA DE BLOQUEIO ATIVADO ❌❌❌"
    echo "======================================"
    echo "O AGENTE ESTÁ BLOQUEADO PARA IMPLEMENTAÇÕES"
    echo "ATÉ CONFIRMAR LEITURA DA DOCUMENTAÇÃO OBRIGATÓRIA"
    echo ""
    echo "📚 DOCUMENTAÇÃO OBRIGATÓRIA (BLOQUEANTE - LEIA AGORA):"
    echo "1. 📋 docs/QUICK_REFERENCE.md - Arquitetura, padrões e handlers"
    echo "2. 🚫 docs/ANTIPATTERNS.md - O que NÃO fazer no projeto"
    echo "3. 🎯 docs/NEXT_STEPS.md - Tarefas pendentes e prioridades"
    echo "4. 📚 docs/learnings/ - TODOS os aprendizados anteriores (CRÍTICO)"
    echo ""
    echo "🎯 FLUXO DE DESBLOQUEIO OBRIGATÓRIO:"
    echo "1. 📚 LEIA documentação acima (AGORA - BLOQUEANTE)"
    echo "2. ✅ MARQUE checklist: ${DOCS_CHECKLIST}"
    echo "3. 🔓 BLOQUEIO REMOVIDO - Agora pode implementar"
    echo "4. 📋 CRIE tarefa: ./create_task.sh \"Nome da Tarefa\" [módulo]"
    echo "5. 🔍 PROCESSE tarefa: ./process_task.sh --task=[ID] --checklist"
    echo "6. 💻 IMPLEMENTE seguindo padrões do projeto"
    echo "7. 📝 REGISTRE aprendizado em session_learnings/"
    echo "8. ✅ CONCLUA tarefa: ./conclude_task.sh"
    echo ""
    echo "📁 Sistema de verificação: ${DOCS_CHECKLIST}"
    echo "📁 Contexto do agente: .agent_context.md"
    echo "💡 Modo rápido (sem bloqueio): ./start_session.sh quick"
else
    echo "✅ SESSÃO RÁPIDA INICIADA! ⚡ SEM BLOQUEIO"
    echo "=========================================="
    echo "ID: ${SESSION_ID}"
    echo "Modo: ${QUICK_MODE} (RÁPIDO)"
    echo "Status: ⚡ IMPLEMENTAÇÕES LIBERADAS"
    echo ""
    echo "⚠️⚠️⚠️ AVISO: MODO RÁPIDO - SEM VERIFICAÇÃO ⚠️⚠️⚠️"
    echo "=============================================="
    echo "O AGENTE NÃO VERIFICOU LEITURA DA DOCUMENTAÇÃO"
    echo "RISCO AUMENTADO DE VIOLAR PADRÕES DO PROJETO"
    echo ""
    echo "📚 DOCUMENTAÇÃO RECOMENDADA (NÃO BLOQUEANTE):"
    echo "1. 📋 docs/QUICK_REFERENCE.md - Arquitetura, padrões"
    echo "2. 🚫 docs/ANTIPATTERNS.md - O que NÃO fazer"
    echo "3. 🎯 docs/NEXT_STEPS.md - Tarefas pendentes"
    echo ""
    echo "🎯 FLUXO DE TRABALHO RÁPIDO:"
    echo "1. 📋 CRIE tarefa: ./create_task.sh \"Nome da Tarefa\" [módulo]"
    echo "2. 🔍 PROCESSE tarefa: ./process_task.sh --task=[ID] --checklist"
    echo "3. 💻 IMPLEMENTE com cuidado (risco de antipadrões)"
    echo "4. 📝 REGISTRE aprendizado em session_learnings/"
    echo "5. ✅ CONCLUA tarefa: ./conclude_task.sh"
    echo ""
    echo "📁 Contexto do agente: .agent_context.md"
    echo "💡 Modo completo (com verificação): ./start_session.sh"
fi