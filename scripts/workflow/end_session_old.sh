#!/bin/bash

# 🚪 SCRIPT DE ENCERRAMENTO DE SESSÃO
# Analisa aprendizados e atualiza documentação antes de limpar

set -e

# Cores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Função para analisar aprendizados da sessão
analyze_session_learnings() {
    echo -e "${PURPLE}🔍 ANALISANDO APRENDIZADOS DA SESSÃO${NC}"
    echo "======================================"
    
    if [ ! -d "docs/learnings" ] || [ -z "$(ls -A docs/learnings 2>/dev/null)" ]; then
        echo "ℹ️  Nenhum aprendizado encontrado em docs/learnings/"
        return 0
    fi
    
    # Contar aprendizados
    LEARNING_FILES=$(find docs/learnings -name "*.md" 2>/dev/null | wc -l)
    echo "📚 Aprendizados encontrados: $LEARNING_FILES arquivo(s)"
    
    # Coletar estatísticas
    TOTAL_SUCCESS=0
    TOTAL_PARTIAL=0
    TOTAL_FAILED=0
    TOTAL_DURATION=0
    KEY_LEARNINGS=""
    
    # Analisar cada arquivo de aprendizado
    for learning_file in docs/learnings/*.md; do
        if [ -f "$learning_file" ]; then
            # Extrair informações básicas
            STATUS=$(grep -i "^\*\*Status:\*\*" "$learning_file" | cut -d: -f2- | xargs)
            DURATION=$(grep -i "^\*\*Duração:\*\*" "$learning_file" | grep -o "[0-9]\+" | head -1)
            TASK_ID=$(grep -i "^\*\*Tarefa ID:\*\*" "$learning_file" | sed 's/\*\*Tarefa ID:\*\*//; s/\*\*//g' | xargs)
            FEATURE=$(grep -i "^\*\*Feature:\*\*" "$learning_file" 2>/dev/null | sed 's/\*\*Feature:\*\*//; s/\*\*//g' | xargs)
            if [ -z "$FEATURE" ]; then
                # Tentar extrair do título
                FEATURE=$(grep "^# 📚 Aprendizados:" "$learning_file" | sed 's/# 📚 Aprendizados: //' | xargs)
            fi
            FEATURE="${FEATURE:-Desconhecida}"
            
            # Contar por status
            case "$STATUS" in
                *success*|*Success*)
                    TOTAL_SUCCESS=$((TOTAL_SUCCESS + 1))
                    ;;
                *partial*|*Partial*)
                    TOTAL_PARTIAL=$((TOTAL_PARTIAL + 1))
                    ;;
                *failed*|*Failed*)
                    TOTAL_FAILED=$((TOTAL_FAILED + 1))
                    ;;
            esac
            
            # Somar duração
            if [ -n "$DURATION" ] && [ "$DURATION" -gt 0 ]; then
                TOTAL_DURATION=$((TOTAL_DURATION + DURATION))
            fi
            
        # Extrair aprendizados chave (melhor parsing)
        if grep -q "## 🎯 Aprendizados Documentados" "$learning_file"; then
            LEARNING_SECTION=$(sed -n '/## 🎯 Aprendizados Documentados/,/^---/p' "$learning_file" | sed '1d' | head -5)
            # Remover markdown e limpar
            CLEAN_LEARNING=$(echo "$LEARNING_SECTION" | grep -v "^#" | grep -v "^\*\*" | grep -v "^$" | head -2 | tr '\n' ' ' | sed 's/  / /g; s/^ //; s/ $//')
            if [ -n "$CLEAN_LEARNING" ] && [ ${#CLEAN_LEARNING} -gt 10 ]; then
                KEY_LEARNINGS="${KEY_LEARNINGS}- **$FEATURE** ($TASK_ID): ${CLEAN_LEARNING:0:80}...\n"
            fi
        fi
        fi
    done
    
    # Mostrar resumo
    echo ""
    echo -e "${CYAN}📊 RESUMO DA SESSÃO:${NC}"
    echo "-------------------"
    echo "✅ Sucessos: $TOTAL_SUCCESS"
    echo "⚠️  Parciais: $TOTAL_PARTIAL"
    echo "❌ Falhas: $TOTAL_FAILED"
    echo "⏱️  Tempo total: ${TOTAL_DURATION} minutos"
    echo ""
    
    if [ -n "$KEY_LEARNINGS" ]; then
        echo -e "${CYAN}🎯 APRENDIZADOS CHAVE:${NC}"
        echo "-------------------"
        echo -e "$KEY_LEARNINGS"
    fi
    
    return 0
}

# Função para atualizar documentação baseada nos aprendizados
update_documentation_from_learnings() {
    echo -e "${PURPLE}📝 ATUALIZANDO DOCUMENTAÇÃO${NC}"
    echo "=============================="
    
    UPDATES_MADE=0
    
    # 1. Atualizar QUICK_REFERENCE.md com aprendizados relevantes
    if [ -f "docs/QUICK_REFERENCE.md" ]; then
        echo "📋 Atualizando QUICK_REFERENCE.md..."
        
        # Verificar se há aprendizados sobre padrões
        PATTERN_LEARNINGS=$(grep -r -i "padrão\|pattern" docs/learnings/ 2>/dev/null | grep -v "antipadrão" | head -3 || true)
        
        if [ -n "$PATTERN_LEARNINGS" ]; then
            # Primeiro, remover seção anterior de aprendizados se existir
            TEMP_FILE=$(mktemp)
            # Remover seções de "Aprendizados Recentes" anteriores
            grep -v "## 🔄 Aprendizados Recentes da Sessão" docs/QUICK_REFERENCE.md | \
                sed '/^## 🔄 Aprendizados Recentes da Sessão/,/^##/d' | \
                sed '/^## 🔄 Aprendizados Recentes da Sessão/,/^## 🆕 Nova Sessão/d' > "$TEMP_FILE"
            
            # Extrair o primeiro aprendizado útil
            FIRST_LEARNING=$(echo "$PATTERN_LEARNINGS" | head -1 | sed 's/.*learnings\///' | cut -d: -f2- | xargs)
            if [ -n "$FIRST_LEARNING" ] && [ ${#FIRST_LEARNING} -gt 20 ]; then
                CLEAN_LEARNING="${FIRST_LEARNING:0:80}..."
            else
                CLEAN_LEARNING="Padrões identificados em $LEARNING_FILES tarefa(s)"
            fi
            
            # Adicionar nova seção antes de "Nova Sessão"
            if grep -q "## 🆕 Nova Sessão" "$TEMP_FILE"; then
                sed -i '/## 🆕 Nova Sessão/i\
\
## 🔄 Aprendizados Recentes da Sessão\
**Baseado em:** '"$LEARNING_FILES"' tarefa(s) concluída(s)\
**Período:** '"$(date +%d/%m/%Y)"'\
\
### Insights:\
- '"$CLEAN_LEARNING"'\
\
### Status:\
- Taxa de sucesso: '"$(if [ $LEARNING_FILES -gt 0 ]; then echo "$((TOTAL_SUCCESS * 100 / LEARNING_FILES))%"; else echo "N/A"; fi)"'\
- Tempo médio: '"$(if [ $LEARNING_FILES -gt 0 ]; then echo "$((TOTAL_DURATION / LEARNING_FILES)) minutos"; else echo "N/A"; fi)"'\
\
💡 **Dica:** Consulte `docs/learnings/` e `docs/ANTIPATTERNS.md` para detalhes.' "$TEMP_FILE"
            else
                # Adicionar no final
                echo "" >> "$TEMP_FILE"
                echo "## 🔄 Aprendizados Recentes da Sessão" >> "$TEMP_FILE"
                echo "**Baseado em:** $LEARNING_FILES tarefa(s) concluída(s)" >> "$TEMP_FILE"
                echo "**Período:** $(date +%d/%m/%Y)" >> "$TEMP_FILE"
                echo "" >> "$TEMP_FILE"
                echo "### Insights:" >> "$TEMP_FILE"
                echo "- $CLEAN_LEARNING" >> "$TEMP_FILE"
                echo "" >> "$TEMP_FILE"
                echo "💡 **Dica:** Consulte \`docs/learnings/\` para aprendizados detalhados." >> "$TEMP_FILE"
            fi
            
            mv "$TEMP_FILE" docs/QUICK_REFERENCE.md
            UPDATES_MADE=$((UPDATES_MADE + 1))
            echo "  ✅ QUICK_REFERENCE.md atualizado com aprendizados"
        else
            echo "  ℹ️  Nenhum aprendizado relevante para QUICK_REFERENCE.md"
        fi
    fi
    
    # 2. Atualizar ou criar ANTIPATTERNS.md se houver aprendizados relevantes
    # Procurar por antipadrões específicos nas seções de aprendizados
    ANTIPATTERNS_FOUND=""
    for learning_file in docs/learnings/*.md; do
        if [ -f "$learning_file" ]; then
            # Extrair seção de antipadrões
            if grep -q "## 📈 Melhorias para Próxima Implementação" "$learning_file"; then
                ANTIPATTERNS_SECTION=$(sed -n '/## 📈 Melhorias para Próxima Implementação/,/^---/p' "$learning_file" | grep -i "antipadrão\|não fazer\|erro\|problema" | head -3)
                if [ -n "$ANTIPATTERNS_SECTION" ]; then
                    # Limpar e adicionar
                    CLEAN_ANTI=$(echo "$ANTIPATTERNS_SECTION" | sed 's/^- \[ \] //; s/^- //; s/\*\*//g' | head -1)
                    if [ -n "$CLEAN_ANTI" ] && [ ${#CLEAN_ANTI} -gt 15 ]; then
                        ANTIPATTERNS_FOUND="${ANTIPATTERNS_FOUND}$CLEAN_ANTI\n"
                    fi
                fi
            fi
        fi
    done
    
    if [ -n "$ANTIPATTERNS_FOUND" ]; then
        echo "🚫 Atualizando antipadrões..."
        
        ANTIPATTERNS_FILE="docs/ANTIPATTERNS.md"
        if [ ! -f "$ANTIPATTERNS_FILE" ]; then
            mkdir -p docs
            cat > "$ANTIPATTERNS_FILE" << EOF
# 🚫 Antipadrões - Projeto Digna

**Última atualização:** $(date +%d/%m/%Y)
**Fonte:** Aprendizados de sessões anteriores

---

## ❌ O que NÃO fazer

Esta lista é baseada em erros comuns identificados durante implementações.

EOF
        fi
        
        # Adicionar novos antipadrões
        echo "" >> "$ANTIPATTERNS_FILE"
        echo "## 🔄 Sessão $(date +%d/%m/%Y)" >> "$ANTIPATTERNS_FILE"
        echo "" >> "$ANTIPATTERNS_FILE"
        
        counter=1
        echo -e "$ANTIPATTERNS_FOUND" | while read -r line; do
            if [ -n "$line" ] && [ ${#line} -gt 15 ]; then
                echo "$counter. **ANTIPADRÃO:** $line" >> "$ANTIPATTERNS_FILE"
                counter=$((counter + 1))
            fi
        done
        
        UPDATES_MADE=$((UPDATES_MADE + 1))
        echo "  ✅ ANTIPATTERNS.md atualizado"
    fi
    
    # 3. Atualizar CHECKLIST_TEMPLATES se houver melhorias identificadas
    CHECKLIST_IMPROVEMENTS=$(grep -r -i "checklist\|validação\|verificar" docs/learnings/ 2>/dev/null | grep -i "melhorar\|faltou\|sugestão" | head -2 || true)
    
    if [ -n "$CHECKLIST_IMPROVEMENTS" ]; then
        echo "📋 Atualizando templates de checklist..."
        
        TEMPLATES_DIR="docs/templates"
        mkdir -p "$TEMPLATES_DIR"
        
        CHECKLIST_TEMPLATE="$TEMPLATES_DIR/pre_implementation_checklist.md"
        if [ ! -f "$CHECKLIST_TEMPLATE" ]; then
            cat > "$CHECKLIST_TEMPLATE" << 'EOF'
# 📋 Template de Checklist Pré-Implementação

**Baseado em:** Aprendizados de sessões anteriores
**Última atualização:** $(date +%d/%m/%Y)

---

## 🔄 Melhorias Identificadas

*Baseado em feedback de implementações anteriores:*

EOF
        fi
        
        # Adicionar melhorias
        echo "" >> "$CHECKLIST_TEMPLATE"
        echo "### Sessão $(date +%d/%m/%Y)" >> "$CHECKLIST_TEMPLATE"
        
        echo "$CHECKLIST_IMPROVEMENTS" | while read -r line; do
            CLEAN_LINE=$(echo "$line" | sed 's/.*learnings\///' | cut -d: -f2- | xargs)
            if [ -n "$CLEAN_LINE" ]; then
                echo "- $CLEAN_LINE" >> "$CHECKLIST_TEMPLATE"
            fi
        done
        
        UPDATES_MADE=$((UPDATES_MADE + 1))
        echo "  ✅ Checklist templates atualizados"
    fi
    
    # 4. Criar resumo da sessão em docs/session_summaries/
    if [ "$LEARNING_FILES" -gt 0 ]; then
        echo "📊 Criando resumo da sessão..."
        
        SESSION_SUMMARY_DIR="docs/session_summaries"
        mkdir -p "$SESSION_SUMMARY_DIR"
        
        SESSION_DATE=$(date +%Y%m%d)
        SUMMARY_FILE="$SESSION_SUMMARY_DIR/session_${SESSION_DATE}.md"
        
        # Coletar dados para o resumo
        SESSION_ID=$(ls .session_* 2>/dev/null | head -1 | sed 's/\.session_//' || echo "unknown")
        
        cat > "$SUMMARY_FILE" << EOF
# 📊 Resumo da Sessão - $(date +%d/%m/%Y)

**Sessão ID:** ${SESSION_ID:-Desconhecida}
**Data:** $(date +%d/%m/%Y\ %H:%M)
**Tarefas concluídas:** $LEARNING_FILES

---

## 📈 Métricas

### Status das Tarefas
- ✅ **Sucessos:** $TOTAL_SUCCESS
- ⚠️  **Parciais:** $TOTAL_PARTIAL  
- ❌ **Falhas:** $TOTAL_FAILED
- ⏱️  **Tempo total:** ${TOTAL_DURATION} minutos

### Eficiência
- **Taxa de sucesso:** $(if [ $LEARNING_FILES -gt 0 ]; then echo "$((TOTAL_SUCCESS * 100 / LEARNING_FILES))%"; else echo "N/A"; fi)
- **Tempo médio por tarefa:** $(if [ $LEARNING_FILES -gt 0 ]; then echo "$((TOTAL_DURATION / LEARNING_FILES)) minutos"; else echo "N/A"; fi)

---

## 🎯 Aprendizados Principais

$(echo -e "$KEY_LEARNINGS" | head -10)

---

## 🔧 Impacto na Documentação

### Arquivos Atualizados:
$(if [ $UPDATES_MADE -gt 0 ]; then
    [ -f "docs/QUICK_REFERENCE.md" ] && echo "- \`QUICK_REFERENCE.md\` (aprendizados integrados)"
    [ -f "docs/ANTIPATTERNS.md" ] && echo "- \`ANTIPATTERNS.md\` (novos antipadrões)"
    [ -f "docs/templates/pre_implementation_checklist.md" ] && echo "- \`pre_implementation_checklist.md\` (melhorias checklist)"
    echo "- \`$SUMMARY_FILE\` (este resumo)"
else
    echo "- Nenhum arquivo atualizado (aprendizados insuficientes)"
fi)

### Próximos Passos Sugeridos:
1. Revisar aprendizados em \`docs/learnings/\`
2. Aplicar antipadrões identificados na próxima sessão
3. Melhorar checklists baseado no feedback

---

## 📁 Arquivos de Aprendizado

$(find docs/learnings -name "*.md" 2>/dev/null | xargs -n1 basename | sed 's/^/- /' | head -10)

$(if [ $(find docs/learnings -name "*.md" 2>/dev/null | wc -l) -gt 10 ]; then
    echo "... e mais $(($(find docs/learnings -name "*.md" 2>/dev/null | wc -l) - 10)) arquivos"
fi)

---

**📌 Nota:** Este resumo foi gerado automaticamente pelo \`end_session.sh\`.
Consulte os arquivos individuais em \`docs/learnings/\` para detalhes completos.
EOF
        
        UPDATES_MADE=$((UPDATES_MADE + 1))
        echo "  ✅ Resumo da sessão criado: $SUMMARY_FILE"
    fi
    
    echo ""
    echo -e "${GREEN}✅ DOCUMENTAÇÃO ATUALIZADA!${NC}"
    echo "  Total de atualizações: $UPDATES_MADE"
    
    return $UPDATES_MADE
}

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

# 2. Analisar aprendizados da sessão
analyze_session_learnings

# 3. Atualizar documentação baseada nos aprendizados
update_documentation_from_learnings
DOC_UPDATES=$?

# 4. Backup opcional de aprendizados
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

# Remover contexto do agente
if [ -f ".agent_context.md" ]; then
    echo "Removendo contexto do agente (.agent_context.md)..."
    rm -f .agent_context.md
    FILES_TO_CLEAN=$((FILES_TO_CLEAN + 1))
fi

# Remover planos de implementação temporários
IMPLEMENTATION_PLANS=$(find docs/implementation_plans -name "*.md" 2>/dev/null | wc -l)
if [ $IMPLEMENTATION_PLANS -gt 0 ]; then
    echo "Removendo $IMPLEMENTATION_PLANS arquivo(s) de planos de implementação..."
    find docs/implementation_plans -name "*.md" -delete 2>/dev/null || true
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
        echo "**Sessão iniciada em:** $(date '+%d/%m/%Y %H:%M')" >> /tmp/quick_ref_tmp
        echo "**Status:** ✅ PRONTO PARA NOVA IMPLEMENTAÇÃO" >> /tmp/quick_ref_tmp
        echo "" >> /tmp/quick_ref_tmp
        echo "Use \`./start_session.sh\` para contexto completo ou \`./process_task.sh\` para começar." >> /tmp/quick_ref_tmp
        
        mv /tmp/quick_ref_tmp docs/QUICK_REFERENCE.md
    fi
    FILES_TO_CLEAN=$((FILES_TO_CLEAN + 1))
fi

# 6. Resumo
echo ""
echo -e "${GREEN}✅ SESSÃO ENCERRADA COM SUCESSO!${NC}"
echo "======================================"
echo ""
echo -e "${CYAN}📊 RESUMO FINAL:${NC}"
echo "-------------------"
echo "📁 Arquivos temporários removidos: $FILES_TO_CLEAN"
echo "📝 Documentação atualizada: $DOC_UPDATES arquivo(s)"
echo "📚 Aprendizados processados: $(find docs/learnings -name "*.md" 2>/dev/null | wc -l)"
echo ""
echo -e "${BLUE}📋 ESTADO ATUAL DO PROJETO:${NC}"
echo "-----------------------------"
echo "✅ Sessões: 0 (limpas)"
echo "✅ Tarefas: 0 (limpas)" 
echo "✅ Prompts: 0 (limpos)"
echo "✅ NEXT_STEPS.md: Limpo e pronto"
echo "✅ QUICK_REFERENCE.md: Atualizado com aprendizados"
echo "✅ ANTIPATTERNS.md: $(if [ -f "docs/ANTIPATTERNS.md" ]; then echo "Atualizado"; else echo "Não criado"; fi)"
echo "✅ Session Summary: $(if [ -d "docs/session_summaries" ] && [ "$(ls -A docs/session_summaries 2>/dev/null)" ]; then echo "Criado"; else echo "Não criado"; fi)"
echo ""
echo -e "${YELLOW}🚀 PRONTO PARA NOVA SESSÃO!${NC}"
echo ""
echo "Para começar do zero:"
echo "1. Execute: ./start_session.sh (contexto atualizado)"
echo "2. Analise: ./scripts/tools/analyze_patterns.sh [handler]"
echo "3. Processe: ./process_task.sh -f \"tarefa.md\" --execute"
echo "4. Conclua: ./conclude_task.sh \"aprendizados\" --success"
echo ""
echo -e "${GREEN}💡 APRENDIZADOS APLICADOS:${NC}"
echo "Os insights desta sessão foram integrados na documentação."
echo "A próxima sessão começará com conhecimento acumulado!"
echo ""
echo -e "${PURPLE}🎉 Sessão encerrada com sucesso!${NC}"

exit 0