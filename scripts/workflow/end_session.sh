#!/bin/bash

# 🚪 SCRIPT DE ENCERRAMENTO DE SESSÃO
# Limpa arquivos temporários e prepara para nova sessão "do zero"

set -e

# Cores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚪 ENCERRANDO SESSÃO ATUAL${NC}"
echo "================================"

# 1. Verificar sessão atual
SESSION_FILES=$(ls .session_* 2>/dev/null | wc -l)
TASK_FILES=$(ls .task_* 2>/dev/null | wc -l)
OPENCODE_FILES=$(ls .opencode_task_*.txt 2>/dev/null | wc -l)

echo "📊 Status atual:"
echo "  Sessões ativas: $SESSION_FILES"
echo "  Tarefas pendentes: $TASK_FILES"
echo "  Prompts opencode: $OPENCODE_FILES"

# 2. Backup opcional de aprendizados
echo ""
echo -e "${YELLOW}📚 BACKUP DE APRENDIZADOS (opcional)${NC}"
echo "--------------------------------"

if [ -d "docs/learnings" ] && [ "$(ls -A docs/learnings 2>/dev/null)" ]; then
    echo "Aprendizados encontrados em docs/learnings/"
    echo "Deseja criar backup antes de limpar? (s/N)"
    read -r BACKUP_CHOICE
    
    if [[ "$BACKUP_CHOICE" =~ ^[Ss]$ ]]; then
        BACKUP_DIR="backup_learnings_$(date +%Y%m%d_%H%M%S)"
        mkdir -p "$BACKUP_DIR"
        cp -r docs/learnings/* "$BACKUP_DIR/" 2>/dev/null || true
        echo -e "${GREEN}✅ Backup criado em: $BACKUP_DIR/${NC}"
    fi
fi

# 3. Limpar arquivos temporários
echo ""
echo -e "${YELLOW}🧹 LIMPANDO ARQUIVOS TEMPORÁRIOS${NC}"
echo "--------------------------------"

FILES_TO_CLEAN=0

# Remover arquivos de sessão
if [ $SESSION_FILES -gt 0 ]; then
    echo "Removendo $SESSION_FILES arquivo(s) de sessão..."
    rm -f .session_*
    FILES_TO_CLEAN=$((FILES_TO_CLEAN + SESSION_FILES))
fi

# Remover arquivos de tarefa
if [ $TASK_FILES -gt 0 ]; then
    echo "Removendo $TASK_FILES arquivo(s) de tarefa..."
    rm -f .task_*
    FILES_TO_CLEAN=$((FILES_TO_CLEAN + TASK_FILES))
fi

# Remover prompts opencode
if [ $OPENCODE_FILES -gt 0 ]; then
    echo "Removendo $OPENCODE_FILES arquivo(s) de prompt opencode..."
    rm -f .opencode_task_*.txt
    FILES_TO_CLEAN=$((FILES_TO_CLEAN + OPENCODE_FILES))
fi

# Remover planos de implementação temporários
IMPLEMENTATION_PLANS=$(find docs/implementation_plans -name "*_pre_check.md" -o -name "*_implementation_*.md" 2>/dev/null | wc -l)
if [ $IMPLEMENTATION_PLANS -gt 0 ]; then
    echo "Removendo $IMPLEMENTATION_PLANS arquivo(s) de planos de implementação..."
    find docs/implementation_plans -name "*_pre_check.md" -delete 2>/dev/null || true
    find docs/implementation_plans -name "*_implementation_*.md" -delete 2>/dev/null || true
    FILES_TO_CLEAN=$((FILES_TO_CLEAN + IMPLEMENTATION_PLANS))
fi

# 4. Limpar NEXT_STEPS.md se existir
if [ -f "docs/NEXT_STEPS.md" ]; then
    echo "Limpando docs/NEXT_STEPS.md..."
    echo "# 🎯 Próximos Passos - Projeto Digna" > docs/NEXT_STEPS.md
    echo "" >> docs/NEXT_STEPS.md
    echo "**Última atualização:** $(date +%d/%m/%Y)" >> docs/NEXT_STEPS.md
    echo "" >> docs/NEXT_STEPS.md
    echo "---" >> docs/NEXT_STEPS.md
    echo "" >> docs/NEXT_STEPS.md
    echo "## 🚀 Próxima Tarefa" >> docs/NEXT_STEPS.md
    echo "" >> docs/NEXT_STEPS.md
    echo "Escolha uma tarefa do backlog ou crie uma nova:" >> docs/NEXT_STEPS.md
    echo "" >> docs/NEXT_STEPS.md
    echo "1. Use \`./process_task.sh \"sua descrição de tarefa\"\`" >> docs/NEXT_STEPS.md
    echo "2. Siga o checklist pré-implementação" >> docs/NEXT_STEPS.md
    echo "3. Documente aprendizados com \`./conclude_task.sh\`" >> docs/NEXT_STEPS.md
    FILES_TO_CLEAN=$((FILES_TO_CLEAN + 1))
fi

# 5. Atualizar QUICK_REFERENCE.md para estado limpo
if [ -f "docs/QUICK_REFERENCE.md" ]; then
    echo "Atualizando docs/QUICK_REFERENCE.md para nova sessão..."
    # Manter apenas as seções essenciais
    head -100 docs/QUICK_REFERENCE.md > /tmp/quick_ref_tmp 2>/dev/null || true
    if [ -s /tmp/quick_ref_tmp ]; then
        # Adicionar marcador de nova sessão
        echo "" >> /tmp/quick_ref_tmp
        echo "---" >> /tmp/quick_ref_tmp
        echo "" >> /tmp/quick_ref_tmp
        echo "## 🆕 Nova Sessão" >> /tmp/quick_ref_tmp
        echo "" >> /tmp/quick_ref_tmp
        echo "**Sessão iniciada em:** $(date +%d/%m/%Y %H:%M)" >> /tmp/quick_ref_tmp
        echo "**Status:** ✅ PRONTO PARA NOVA IMPLEMENTAÇÃO" >> /tmp/quick_ref_tmp
        echo "" >> /tmp/quick_ref_tmp
        echo "Use \`./start_session.sh\` para contexto completo ou \`./process_task.sh\` para começar." >> /tmp/quick_ref_tmp
        
        mv /tmp/quick_ref_tmp docs/QUICK_REFERENCE.md
    fi
    FILES_TO_CLEAN=$((FILES_TO_CLEAN + 1))
fi

# 6. Resumo
echo ""
echo -e "${GREEN}✅ LIMPEZA CONCLUÍDA!${NC}"
echo "========================="
echo "Arquivos removidos/atualizados: $FILES_TO_CLEAN"
echo ""
echo -e "${BLUE}📋 ESTADO ATUAL:${NC}"
echo "-------------------"
echo "✅ Sessões: 0"
echo "✅ Tarefas: 0" 
echo "✅ Prompts: 0"
echo "✅ NEXT_STEPS.md: Limpo"
echo "✅ QUICK_REFERENCE.md: Atualizado"
echo ""
echo -e "${YELLOW}🚀 PRONTO PARA NOVA SESSÃO!${NC}"
echo ""
echo "Para começar do zero:"
echo "1. Execute: ./start_session.sh"
echo "2. Escolha uma tarefa: ./process_task.sh \"sua tarefa\""
echo "3. Siga o fluxo normal"
echo ""
echo -e "${GREEN}🎉 Sessão encerrada com sucesso!${NC}"

exit 0