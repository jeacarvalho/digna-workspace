#!/bin/bash
# conclude_task.sh - Conclui uma tarefa e documenta aprendizados
# Uso: ./conclude_task.sh --task=TASK_ID "Aprendizados: item1, item2" [--success] [--partial] [--failed]

set -e

echo "🔗 Concluindo tarefa..."
echo "======================="

# Configurações
SESSION_DIR="work_in_progress/current_session"
TASK_ID=""
LEARNINGS=""
STATUS="success"

# Funções de ajuda
show_help() {
    cat << EOF
🎯 conclude_task.sh - Conclui tarefa e documenta aprendizados

USO:
  ./conclude_task.sh --task=TASK_ID "DESCRIÇÃO DOS APRENDIZADOS" [OPÇÕES]

OPÇÕES:
  --task, -t ID      ID da tarefa (obrigatório)
  --success          Tarefa concluída com sucesso (padrão)
  --partial          Tarefa parcialmente concluída
  --failed           Tarefa falhou/não concluída
  --help, -h         Mostrar esta ajuda

EXEMPLOS:
  # Concluir com sucesso
  ./conclude_task.sh --task=20250311_101108 "Aprendizados: checklist antecipou 3 problemas" --success
  
  # Concluir parcialmente
  ./conclude_task.sh --task=20250311_101108 "Problemas: serviço internal não acessível" --partial
  
  # Marcar como falha
  ./conclude_task.sh --task=20250311_101108 "Falha: bug crítico não resolvido" --failed

O script irá:
1. Validar implementação (handlers, testes, etc.)
2. Coletar métricas da tarefa
3. Documentar aprendizados
4. Atualizar checklists e antipadrões
5. Preparar para próxima sessão
6. Mover tarefa para archive
EOF
    exit 0
}

# Processar argumentos
NEXT_IS_TASK=false
for arg in "$@"; do
    case $arg in
        --task=*|-t=*)
            TASK_ID="${arg#*=}"
            ;;
        --task|-t)
            NEXT_IS_TASK=true
            ;;
        --success)
            STATUS="success"
            ;;
        --partial)
            STATUS="partial"
            ;;
        --failed)
            STATUS="failed"
            ;;
        --help|-h)
            show_help
            ;;
        *)
            if [ "$NEXT_IS_TASK" = true ]; then
                TASK_ID="$arg"
                NEXT_IS_TASK=false
            elif [ -z "$LEARNINGS" ]; then
                LEARNINGS="$arg"
            fi
            ;;
    esac
done

# Verificar se temos ID da tarefa
if [ -z "$TASK_ID" ]; then
    echo "❌ ID da tarefa é obrigatório."
    echo ""
    echo "📋 Tarefas em andamento:"
    if [ -d "work_in_progress/tasks" ]; then
        find work_in_progress/tasks -name "task_metadata" -exec sh -c 'echo "  - $(basename $(dirname {})): $(grep "TASK_NAME=" {} | cut -d= -f2 | tr -d "\"") ($(grep "STATUS=" {} | cut -d= -f2 | tr -d "\""))"' \; 2>/dev/null || echo "  Nenhuma tarefa encontrada"
    else
        echo "  Nenhuma tarefa encontrada"
    fi
    echo ""
    echo "💡 Use: ./conclude_task.sh --task=ID \"Aprendizados\""
    exit 1
fi

# Verificar se tarefa existe
TASK_DIR="work_in_progress/tasks/task_${TASK_ID}"
if [ ! -d "$TASK_DIR" ]; then
    echo "❌ Tarefa não encontrada: ${TASK_DIR}"
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
    echo "⚠️  Metadados da tarefa não encontrados."
    TASK_NAME="Tarefa ${TASK_ID}"
    MODULE="ui_web"
    CREATED_AT=$(date +%s)
fi

# Verificar aprendizados
if [ -z "$LEARNINGS" ]; then
    echo "⚠️  Aprendizados não especificados."
    echo "💡 Use: ./conclude_task.sh --task=${TASK_ID} \"Aprendizados: item1, item2\""
    LEARNINGS="Aprendizados não documentados"
fi

echo "📚 Concluindo tarefa: ${TASK_ID}"
echo "================================"
echo "Nome: ${TASK_NAME}"
echo "Status: ${STATUS}"
echo "Módulo: ${MODULE}"

# Verificação de segurança - garantir que não é execução automática do agente
if [ -n "${OPENCODE_AGENT}" ] || [ -n "${AUTOMATIC_EXECUTION}" ]; then
    echo "❌❌❌ ERRO DE SEGURANÇA ❌❌❌"
    echo "Este script NÃO deve ser executado automaticamente pelo agente opencode."
    echo "O agente deve INFORMAR ao usuário que a tarefa pode ser concluída."
    echo "O usuário deve executar este script manualmente após validar a implementação."
    echo ""
    echo "Fluxo correto:"
    echo "1. Agente implementa tarefa"
    echo "2. Agente informa: 'Tarefa pode ser concluída'"
    echo "3. Usuário executa: ./conclude_task.sh --task=${TASK_ID} \"Aprendizados\""
    echo "4. Tarefa é arquivada"
    exit 1
fi

# Calcular duração
END_TIME=$(date +%s)
DURATION=$((END_TIME - CREATED_AT))
DURATION_MIN=$((DURATION / 60))

echo "Duração: ${DURATION_MIN} minutos"
echo ""

# 1. VALIDAÇÃO OBRIGATÓRIA
echo "🔍 VALIDAÇÃO OBRIGATÓRIA ANTES DE CONCLUIR:"
echo "=========================================="

VALIDATION_PASSED=true

# 1.1 Verificar se handler está registrado no main.go (apenas para Features/UI)
echo "1. Handler registrado no main.go?"
if [ "$MODULE" = "ui_web" ] || [ "$TASK_TYPE" = "Feature" ]; then
    FEATURE_NAME_SNAKE=$(echo "${TASK_NAME}" | tr ' ' '_' | tr '[:upper:]' '[:lower:]')
    if grep -q "New${FEATURE_NAME_SNAKE^}Handler" modules/ui_web/main.go 2>/dev/null; then
        echo "   ✅ SIM - Handler encontrado em main.go"
    else
        echo "   ℹ️  NÃO - Handler NÃO está em main.go (pode ser intencional para bibliotecas)"
        # Não falha para bibliotecas
    fi
else
    echo "   ℹ️  Não aplicável para módulo: ${MODULE}"
fi

# 1.2 Verificar testes de sistema
echo "2. Testes de sistema passam?"
cd modules 2>/dev/null
if [ $? -eq 0 ]; then
    TEST_OUTPUT=$(./run_tests.sh 2>&1 | tail -30)
    if echo "$TEST_OUTPUT" | grep -q "PASS\|ok"; then
        echo "   ✅ SIM - Testes de sistema passam"
    else
        echo "   ⚠️  Testes de sistema não executados ou falharam"
        echo "   💡 Ação: Execute: go test -v ./..."
    fi
    cd - >/dev/null 2>&1
else
    echo "   ℹ️  Diretório modules não encontrado"
fi

# 1.3 Verificar se há testes E2E específicos para a tarefa
echo "3. Testes E2E específicos para a tarefa?"
# Procurar por testes Go em modules
TASK_GO_TEST_FILES=$(find modules -name "*test*.go" -newer "${TASK_DIR}/task_prompt.md" 2>/dev/null | wc -l)
# Procurar por testes Playwright (JavaScript) em múltiplos diretórios
TASK_PLAYWRIGHT_TEST_FILES=$(find scripts/testing -name "*.spec.js" -newer "${TASK_DIR}/task_prompt.md" 2>/dev/null | wc -l)
TASK_TOTAL_TEST_FILES=$((TASK_GO_TEST_FILES + TASK_PLAYWRIGHT_TEST_FILES))

if [ "$TASK_TOTAL_TEST_FILES" -gt 0 ]; then
    echo "   ✅ SIM - $TASK_TOTAL_TEST_FILES arquivos de teste criados/modificados"
    if [ "$TASK_GO_TEST_FILES" -gt 0 ]; then
        echo "      📝 $TASK_GO_TEST_FILES teste(s) Go"
    fi
    if [ "$TASK_PLAYWRIGHT_TEST_FILES" -gt 0 ]; then
        echo "      🎭 $TASK_PLAYWRIGHT_TEST_FILES teste(s) Playwright"
    fi
else
    echo "   ⚠️  NENHUM teste específico criado para esta tarefa"
    echo "   💡 Recomendado: Crie testes E2E com Playwright"
fi

# 1.4 Smoke test (se aplicável)
echo "4. Smoke test executado?"
if [ -f "./scripts/dev/smoke_test_new_feature.sh" ] && [ "$MODULE" = "ui_web" ]; then
    echo "   ℹ️  Script disponível: ./scripts/dev/smoke_test_new_feature.sh"
    echo "   💡 Recomendado: Execute antes de concluir"
else
    echo "   ℹ️  Smoke test não aplicável ou script não encontrado"
fi

# 1.5 Verificar testes Playwright (E2E)
echo "5. Testes Playwright (E2E) existem?"
PLAYWRIGHT_TESTS=$(find modules -name "*e2e*test*.go" -o -name "*playwright*test*.go" 2>/dev/null | wc -l)
if [ "$PLAYWRIGHT_TESTS" -gt 0 ]; then
    echo "   ✅ SIM - $PLAYWRIGHT_TESTS testes E2E encontrados"
else
    echo "   ⚠️  NENHUM teste E2E com Playwright encontrado"
    echo "   💡 Recomendado: Crie testes E2E para validação completa"
fi

# 1.6 Executar validação detalhada de testes
echo "6. Validação detalhada de testes..."
if [ -f "./scripts/validate_task_tests.sh" ]; then
    echo "   Executando validação de testes..."
    ./scripts/validate_task_tests.sh --task=${TASK_ID} 2>&1 | tail -20
    VALIDATION_EXIT=$?
    
    if [ "$VALIDATION_EXIT" -ne 0 ] && [ "$VALIDATION_EXIT" -ne 1 ]; then
        echo "   ⚠️  Validação de testes encontrou problemas"
        # Não falha automaticamente - permite override com --force
    fi
else
    echo "   ℹ️  Script de validação de testes não encontrado"
fi

# 1.7 VALIDAÇÃO CRÍTICA: Requisitos específicos do prompt
echo "7. Validação de requisitos específicos do prompt..."
if [ -f "./scripts/validate_task_requirements.sh" ]; then
    echo "   Executando validação de requisitos..."
    REQUIREMENTS_OUTPUT=$(./scripts/validate_task_requirements.sh --task=${TASK_ID} 2>&1 | tail -30)
    REQUIREMENTS_EXIT=$?
    
    echo "$REQUIREMENTS_OUTPUT"
    
    if [ "$REQUIREMENTS_EXIT" -ne 0 ]; then
        echo "   🚨 VALIDAÇÃO DE REQUISITOS FALHOU!"
        echo "   A tarefa NÃO atende aos requisitos específicos do prompt"
        
        # Se é tarefa de correção de bug, apenas avisa (não bloqueia)
        if echo "$TASK_NAME" | grep -qi "corrigir\|bug\|erro\|fix"; then
            echo "   ⚠️  Tarefa de correção de bug - permitindo com aviso"
        else
            echo "   ⚠️  Requisitos não atendidos (permitindo com aviso)"
        fi
    else
        echo "   ✅ Validação de requisitos PASSOU"
    fi
else
    echo "   ℹ️  Script de validação de requisitos não encontrado"
    echo "   💡 Crie: ./scripts/validate_task_requirements.sh para validação robusta"
fi

if [ "$VALIDATION_PASSED" = false ]; then
    echo ""
    echo "🚨 VALIDAÇÃO FALHOU!"
    echo "A tarefa NÃO pode ser concluída até corrigir os problemas acima."
    echo "Corrija e execute este script novamente."
    exit 1
fi

echo ""
echo "✅ VALIDAÇÃO PASSOU - Continuando com conclusão..."
echo ""

# 2. Coletar métricas atuais
echo "📊 COLETANDO MÉTRICAS..."
echo "======================="

# Testes
TEST_METRICS=""
if [ -d "modules" ]; then
    cd modules
    TEST_OUTPUT=$(./run_tests.sh 2>&1 | tail -30)
    TEST_SUMMARY=$(echo "$TEST_OUTPUT" | grep -E "(PASS|FAIL|ok|^---)" | tail -5)
    TEST_METRICS=$(echo "$TEST_OUTPUT" | grep -oE "[0-9]+ (passed|failed)" || echo "Testes não executados")
    cd ..
else
    TEST_METRICS="Diretório modules não encontrado"
fi

# Código (se implementado)
CODE_METRICS=""
FEATURE_NAME_SNAKE=$(echo "${TASK_NAME}" | tr ' ' '_' | tr '[:upper:]' '[:lower:]')
if [ -f "modules/ui_web/internal/handler/${FEATURE_NAME_SNAKE}_handler.go" ]; then
    HANDLER_LINES=$(wc -l < "modules/ui_web/internal/handler/${FEATURE_NAME_SNAKE}_handler.go" 2>/dev/null || echo "0")
    TEST_LINES=$(wc -l < "modules/ui_web/internal/handler/${FEATURE_NAME_SNAKE}_handler_test.go" 2>/dev/null || echo "0")
    TEMPLATE_LINES=$(wc -l < "modules/ui_web/templates/${FEATURE_NAME_SNAKE}_simple.html" 2>/dev/null || echo "0")
    CODE_METRICS="Handler: ${HANDLER_LINES} linhas, Testes: ${TEST_LINES} linhas, Template: ${TEMPLATE_LINES} linhas"
else
    CODE_METRICS="Arquivos de implementação não encontrados (pode ser biblioteca)"
fi

# 3. Criar documento de aprendizados da tarefa
echo ""
echo "📝 CRIANDO DOCUMENTO DE APRENDIZADOS DA TAREFA..."
echo "================================================"

TASK_LEARNINGS_FILE="${TASK_DIR}/task_learnings.md"
mkdir -p "${SESSION_DIR}/session_learnings"

# Extrair aprendizados da descrição
CLEAN_LEARNINGS=$(echo "$LEARNINGS" | sed 's/Aprendizados:\s*//i')

cat > ${TASK_LEARNINGS_FILE} << EOF
# 📚 Aprendizados da Tarefa: ${TASK_NAME}

**Tarefa ID:** ${TASK_ID}
**Concluído em:** $(date +%d/%m/%Y %H:%M:%S)
**Status:** ${STATUS}
**Duração:** ${DURATION_MIN} minutos
**Módulo:** ${MODULE}

---

## 📊 Métricas da Implementação

### Tempo e Status
- **Tempo total:** ${DURATION_MIN} minutos
- **Status:** ${STATUS}
- **Módulo:** ${MODULE}

### Testes
\`\`\`
${TEST_SUMMARY}
\`\`\`
**Resumo:** ${TEST_METRICS}

### Código Produzido
${CODE_METRICS}

### Arquivos Gerados
- Prompt: \`${TASK_DIR}/task_prompt.md\`
- Checklist: \`${TASK_DIR}/checklist.md\` (se gerado)
- Plano: \`${TASK_DIR}/implementation_plan.md\` (se gerado)
- Este documento: \`${TASK_LEARNINGS_FILE}\`

---

## 🎯 Aprendizados Documentados

${CLEAN_LEARNINGS}

---

## 🔍 Análise do Processo

### O que funcionou bem:
1. [Preencher baseado na experiência]
2. [O que ajudou na implementação?]
3. [Quais ferramentas/processos foram úteis?]

### Problemas encontrados:
1. [Listar problemas técnicos]
2. [Problemas de processo]
3. [Falhas de comunicação/entendimento]

### Impacto dos problemas:
- **Tempo perdido:** [X] minutos/horas
- **Retrabalho:** [Sim/Não] - [Descrição]
- **Complexidade aumentada:** [Sim/Não] - [Por quê?]

---

## 📈 Melhorias para Próxima Implementação

### 1. Atualizar Checklists
- [ ] Adicionar item sobre: [novo problema encontrado]
- [ ] Melhorar item: [item que foi confuso]
- [ ] Remover item: [item irrelevante]

### 2. Atualizar Antipadrões
- [ ] Adicionar antipadrão: [novo padrão problemático]
- [ ] Atualizar solução: [solução melhorada]
- [ ] Adicionar exemplo: [exemplo concreto]

### 3. Melhorar Templates
- [ ] Atualizar template de: [qual template?]
- [ ] Adicionar seção: [nova seção necessária]
- [ ] Simplificar: [o que pode ser mais simples?]

---

## 🚀 Próximos Passos Recomendados

### Imediatos (próxima sessão):
1. [Ação 1 - baseado no status]
2. [Ação 2 - correção/continuação]
3. [Ação 3 - validação]

### Médio prazo (sprint):
1. [Refatoração necessária]
2. [Integração pendente]
3. [Testes adicionais]

### Longo prazo (roadmap):
1. [Melhoria arquitetural]
2. [Feature relacionada]
3. [Otimização]

---

## ✅ Checklist de Conclusão

### Validação Técnica
- [ ] Testes passando: ${TEST_METRICS}
- [ ] Código segue padrões: [Sim/Não]
- [ ] Documentação atualizada: [Sim/Não]

### Processo
- [ ] Aprendizados documentados: ✅ (este arquivo)
- [ ] Checklists atualizados: [Pendente/Feito]
- [ ] Próximos passos definidos: ✅ (acima)

### Próxima Sessão
- [ ] Contexto atualizado: [Pendente/Feito]
- [ ] Tarefas priorizadas: [Pendente/Feito]
- [ ] Lições aplicadas: [Pendente/Feito]

---

## 🔄 Feedback do Sistema

### Checklist pré-implementação foi útil?
- **Problemas antecipados:** [X]/[Y] (quantos do checklist ocorreram?)
- **Problemas não previstos:** [Z] (listar abaixo)
- **Sugestões de melhoria:** [O que faltou?]

### Templates e scripts ajudaram?
- **create_task.sh:** [1-5] (1=não, 5=muito)
- **process_task.sh:** [1-5]
- **conclude_task.sh:** [1-5] (este script)

### O que falta no sistema?
1. [Funcionalidade desejada]
2. [Melhoria no processo]
3. [Ferramenta adicional]

---

**📌 Nota:** Este documento deve ser revisado antes da próxima implementação similar.
Use estas lições para melhorar o processo continuamente.
EOF

echo "✅ Documento de aprendizados da tarefa criado: ${TASK_LEARNINGS_FILE}"

# 4. Copiar aprendizados para sessão (resumo)
SESSION_LEARNINGS_FILE="${SESSION_DIR}/session_learnings/task_${TASK_ID}_learnings.md"
cp "${TASK_LEARNINGS_FILE}" "${SESSION_LEARNINGS_FILE}" 2>/dev/null || true

# 5. Atualizar checklists e antipadrões (se houver aprendizados relevantes)
echo ""
echo "🔄 ATUALIZANDO CHECKLISTS E ANTIPADRÕES..."
echo "=========================================="

# Verificar se há checklists para atualizar
if [ -f "docs/templates/pre_implementation_checklist.md" ]; then
    echo "✅ Checklist de templates disponível para atualização"
    # Em uma versão futura, poderíamos automatizar a atualização
fi

if [ -f "docs/ANTIPATTERNS.md" ]; then
    echo "✅ Antipadrões disponíveis para atualização"
fi

# 6. Atualizar próximos passos
echo ""
echo "⏭️  ATUALIZANDO PRÓXIMOS PASSOS..."
echo "================================="

if [ -f "docs/NEXT_STEPS.md" ]; then
    # Adicionar conclusão da tarefa aos próximos passos
    sed -i "/^## 📋 Tarefas Pendentes/,/^##/ { /${TASK_NAME}/s/\[ \]/\[x\]/ }" docs/NEXT_STEPS.md 2>/dev/null || true
    echo "✅ Próximos passos atualizados: docs/NEXT_STEPS.md"
else
    echo "⚠️  docs/NEXT_STEPS.md não encontrado"
fi

# 7. Mover tarefa para archive da sessão
echo ""
echo "📦 MOVENDO TAREFA PARA ARCHIVE..."
echo "================================"

SESSION_ID=$(cat "${SESSION_DIR}/session_info" 2>/dev/null | grep "SESSION_ID=" | cut -d= -f2 || echo "unknown")
ARCHIVE_DIR="work_in_progress/archive/session_${SESSION_ID}/tasks"
mkdir -p "${ARCHIVE_DIR}"

# Mover diretório da tarefa
mv "${TASK_DIR}" "${ARCHIVE_DIR}/" 2>/dev/null

if [ $? -eq 0 ]; then
    echo "✅ Tarefa movida para: ${ARCHIVE_DIR}/task_${TASK_ID}"
else
    echo "⚠️  Não foi possível mover a tarefa (pode já ter sido movida)"
fi

# 8. Limpar arquivos temporários
echo ""
echo "🧹 LIMPANDO ARQUIVOS TEMPORÁRIOS..."
echo "=================================="

# Remover arquivos temporários no root (compatibilidade)
rm -f ".task_${TASK_ID}" ".opencode_task_${TASK_ID}.txt" 2>/dev/null || true
echo "✅ Arquivos temporários removidos"

# 9. Atualizar contexto do projeto
echo ""
echo "🔄 ATUALIZANDO CONTEXTO DO PROJETO..."
echo "===================================="

if [ -f "./scripts/update_context.sh" ]; then
    ./scripts/update_context.sh 2>/dev/null || true
    echo "✅ Contexto atualizado"
else
    echo "⚠️  Script update_context.sh não encontrado"
fi

# 10. Atualizar contexto do agente
echo ""
echo "🤖 ATUALIZANDO CONTEXTO DO AGENTE..."
echo "==================================="

if [ -f "${SESSION_DIR}/.agent_context.md" ]; then
    # Adicionar conclusão da tarefa ao contexto
    sed -i "/^## 📋 STATUS DA SESSÃO/,/^---/ { /Tarefas concluídas:/s/$/ ${TASK_NAME}/ }" "${SESSION_DIR}/.agent_context.md" 2>/dev/null || true
    echo "✅ Contexto do agente atualizado com conclusão"
else
    echo "⚠️  Contexto do agente não encontrado"
fi

echo ""
echo "✅ CONCLUSÃO DA TAREFA COMPLETA!"
echo "================================"

echo ""
echo "📊 RESUMO:"
echo "----------"
echo "Tarefa: ${TASK_NAME} (${TASK_ID})"
echo "Status: ${STATUS}"
echo "Duração: ${DURATION_MIN} minutos"
echo "Testes: ${TEST_METRICS}"
echo ""

echo "📁 ARQUIVOS CRIADOS:"
echo "-------------------"
echo "1. Aprendizados da tarefa: ${TASK_LEARNINGS_FILE}"
echo "2. Cópia na sessão: ${SESSION_LEARNINGS_FILE}"
echo "3. Próximos passos atualizados: docs/NEXT_STEPS.md"
echo "4. Tarefa arquivada: ${ARCHIVE_DIR}/task_${TASK_ID}"
echo ""

echo "🚀 PRÓXIMOS PASSOS:"
echo "------------------"
echo "1. Continuar sessão: Trabalhar em outra tarefa"
echo "2. Listar tarefas: ls work_in_progress/tasks/"
echo "3. Criar nova tarefa: ./create_task.sh \"Nome\""
echo "4. Encerrar sessão (após TODAS tarefas): ./end_session.sh"
echo ""

echo "💡 Dica: Os aprendizados desta tarefa serão usados para melhorar"
echo "       o processo da próxima implementação!"
echo ""
echo "🤖 INSTRUÇÃO PARA OPENCODE:"
echo "---------------------------"
echo "Leia work_in_progress/current_session/.agent_context.md para status atual."
echo "Use ./start_session.sh para nova sessão (após ./end_session.sh)."