#!/bin/bash
# conclude_task.sh - Conclui uma tarefa e documenta aprendizados
# Uso: ./conclude_task.sh "Aprendizados: item1, item2" [--success] [--partial] [--failed]

set -e

# Configurações
LEARNINGS="${1}"
SESSION_FILE=$(ls .session_* 2>/dev/null | head -1)
TASK_FILE=$(ls .task_* 2>/dev/null | head -1)
CONCLUSION_ID=$(date +%Y%m%d_%H%M%S)

# Status padrão
STATUS="completed"
for arg in "$@"; do
    case $arg in
        --success)
            STATUS="success"
            ;;
        --partial)
            STATUS="partial"
            ;;
        --failed)
            STATUS="failed"
            ;;
        --help)
            cat << EOF
🎯 conclude_task.sh - Conclui tarefa e documenta aprendizados

Uso: ./conclude_task.sh "DESCRIÇÃO DOS APRENDIZADOS" [OPÇÕES]

Exemplos:
  ./conclude_task.sh "Aprendizados: checklist antecipou 3 problemas, testes cobriram 90%"
  ./conclude_task.sh "Problemas: serviço internal não acessível" --partial
  ./conclude_task.sh "Falha: bug crítico não resolvido" --failed

Opções:
  --success  Tarefa concluída com sucesso (padrão)
  --partial  Tarefa parcialmente concluída
  --failed   Tarefa falhou/não concluída
  --help     Mostrar esta ajuda

O script irá:
1. Coletar métricas da tarefa
2. Documentar aprendizados
3. Atualizar checklists e antipadrões
4. Preparar para próxima sessão
EOF
            exit 0
            ;;
    esac
done

# Verificar se há tarefa ativa
if [ ! -f "$TASK_FILE" ]; then
    echo "⚠️  Nenhuma tarefa ativa encontrada (.task_*)"
    echo "💡 Crie uma tarefa primeiro: ./process_task.sh \"sua tarefa\""
    exit 1
fi

# Carregar informações da tarefa
source ${TASK_FILE}
END_TIME=$(date +%s)
DURATION=$((END_TIME - START_TIME))
DURATION_MIN=$((DURATION / 60))

echo "📚 Concluindo tarefa: ${TASK_ID}"
echo "================================"
echo "Feature: ${FEATURE_NAME}"
echo "Status: ${STATUS}"
echo "Duração: ${DURATION_MIN} minutos"

# 1. VALIDAÇÃO OBRIGATÓRIA (NOVO)
echo ""
echo "🔍 VALIDAÇÃO OBRIGATÓRIA ANTES DE CONCLUIR:"
echo "=========================================="

VALIDATION_PASSED=true

# 1.1 Verificar se handler está registrado no main.go (apenas para Features)
echo "1. Handler registrado no main.go?"
if [ "$TASK_TYPE" = "Feature" ] || [ "$TASK_TYPE" = "feature" ]; then
    if grep -q "New${FEATURE_NAME^}Handler" modules/ui_web/main.go 2>/dev/null; then
        echo "   ✅ SIM - Handler encontrado em main.go"
    else
        echo "   ❌ NÃO - Handler NÃO está em main.go!"
        echo "   💡 Ação: Adicione ao modules/ui_web/main.go"
        VALIDATION_PASSED=false
    fi
else
    echo "   ℹ️  Não aplicável para tarefas do tipo: $TASK_TYPE"
fi

# 1.2 Verificar testes de sistema
echo "2. Testes de sistema passam?"
cd modules/ui_web 2>/dev/null
SYSTEM_TEST_OUTPUT=$(go test -v -run "TestSystem.*${FEATURE_NAME}" 2>&1 | tail -20)
if echo "$SYSTEM_TEST_OUTPUT" | grep -q "PASS\|ok"; then
    echo "   ✅ SIM - Testes de sistema passam"
else
    echo "   ⚠️  Testes de sistema não executados ou falharam"
    echo "   💡 Ação: Execute: go test -v -run TestSystem"
fi
cd - >/dev/null 2>&1

# 1.3 Smoke test (se aplicável)
echo "3. Smoke test executado?"
if [ -f "./scripts/smoke_test_new_feature.sh" ] && [ -n "$FEATURE_NAME" ]; then
    echo "   ℹ️  Script disponível: ./scripts/smoke_test_new_feature.sh"
    echo "   💡 Recomendado: Execute antes de concluir"
else
    echo "   ℹ️  Smoke test não aplicável ou script não encontrado"
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
if [ -f "modules/ui_web/internal/handler/${FEATURE_NAME}_handler.go" ]; then
    HANDLER_LINES=$(wc -l < "modules/ui_web/internal/handler/${FEATURE_NAME}_handler.go" 2>/dev/null || echo "0")
    TEST_LINES=$(wc -l < "modules/ui_web/internal/handler/${FEATURE_NAME}_handler_test.go" 2>/dev/null || echo "0")
    TEMPLATE_LINES=$(wc -l < "modules/ui_web/templates/${FEATURE_NAME}_simple.html" 2>/dev/null || echo "0")
    CODE_METRICS="Handler: ${HANDLER_LINES} linhas, Testes: ${TEST_LINES} linhas, Template: ${TEMPLATE_LINES} linhas"
else
    CODE_METRICS="Arquivos de implementação não encontrados"
fi

# 2. Criar documento de aprendizados
echo ""
echo "📝 CRIANDO DOCUMENTO DE APRENDIZADOS..."
echo "======================================"

LEARNINGS_FILE="docs/learnings/${TASK_ID}_${FEATURE_NAME}_learnings.md"
mkdir -p docs/learnings

# Extrair aprendizados da descrição
if [ -z "$LEARNINGS" ]; then
    LEARNINGS="Aprendizados não especificados. Use: ./conclude_task.sh \"Aprendizados: item1, item2\""
fi

CLEAN_LEARNINGS=$(echo "$LEARNINGS" | sed 's/Aprendizados:\s*//i')

cat > ${LEARNINGS_FILE} << EOF
# 📚 Aprendizados: ${FEATURE_NAME}

**Tarefa ID:** ${TASK_ID}
**Concluído em:** $(date +%d/%m/%Y %H:%M:%S)
**Status:** ${STATUS}
**Duração:** ${DURATION_MIN} minutos
**Descrição original:** ${TASK_DESCRIPTION}

---

## 📊 Métricas da Implementação

### Tempo e Status
- **Tempo total:** ${DURATION_MIN} minutos
- **Status:** ${STATUS}
- **Modo usado:** ${MODE}

### Testes
\`\`\`
${TEST_SUMMARY}
\`\`\`
**Resumo:** ${TEST_METRICS}

### Código Produzido
${CODE_METRICS}

### Arquivos Gerados
- Checklist: \`docs/implementation_plans/${FEATURE_NAME}_pre_check.md\`
- Plano: \`docs/implementation_plans/${FEATURE_NAME}_implementation_*.md\`
- Este documento: \`${LEARNINGS_FILE}\`

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
- **start_session.sh:** [1-5] (1=não, 5=muito)
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

echo "✅ Documento de aprendizados criado: ${LEARNINGS_FILE}"

# 3. Atualizar checklists com novos aprendizados
echo ""
echo "🔄 ATUALIZANDO CHECKLISTS..."
echo "============================"

if [ -f "docs/templates/pre_implementation_checklist.md" ]; then
    # Adicionar aprendizado como comentário no checklist
    echo "" >> docs/templates/pre_implementation_checklist.md
    echo "<!-- ${TASK_ID} - ${FEATURE_NAME} - $(date +%d/%m/%Y) -->" >> docs/templates/pre_implementation_checklist.md
    echo "<!-- Aprendizado: ${CLEAN_LEARNINGS:0:100}... -->" >> docs/templates/pre_implementation_checklist.md
    echo "✅ Checklist atualizado com referência ao aprendizado"
else
    echo "⚠️  Checklist template não encontrado, criando..."
    mkdir -p docs/templates
    # Criar básico se não existir
fi

# 4. Atualizar antipadrões se houver problemas significativos
if [ "$STATUS" = "failed" ] || [ "$STATUS" = "partial" ]; then
    echo ""
    echo "🚫 ATUALIZANDO ANTIPADRÕES..."
    echo "============================"
    
    if [ -f "docs/antipatterns/common_antipatterns_solutions.md" ]; then
        echo "" >> docs/antipatterns/common_antipatterns_solutions.md
        echo "## 🔄 Aprendizado de ${FEATURE_NAME} ($(date +%d/%m/%Y))" >> docs/antipatterns/common_antipatterns_solutions.md
        echo "**Problema:** Tarefa com status ${STATUS}" >> docs/antipatterns/common_antipatterns_solutions.md
        echo "**Aprendizado:** ${CLEAN_LEARNINGS:0:200}" >> docs/antipatterns/common_antipatterns_solutions.md
        echo "✅ Antipadrões atualizados"
    fi
fi

# 5. Atualizar NEXT_STEPS.md
echo ""
echo "⏭️  ATUALIZANDO PRÓXIMOS PASSOS..."
echo "================================="

NEXT_STEPS_FILE="docs/NEXT_STEPS.md"
mkdir -p docs

if [ ! -f "$NEXT_STEPS_FILE" ]; then
    cat > ${NEXT_STEPS_FILE} << EOF
# 🎯 Próximos Passos - Projeto Digna

**Última atualização:** $(date +%d/%m/%Y)

---

## 🔄 Continuar de Onde Paramoss

EOF
fi

# Adicionar entrada para esta tarefa
cat >> ${NEXT_STEPS_FILE} << EOF

## ${FEATURE_NAME} (${TASK_ID})
**Status:** ${STATUS}
**Concluído em:** $(date +%d/%m/%Y)
**Duração:** ${DURATION_MIN} minutos

### Próximas Ações:
1. [Baseado no status ${STATUS} - ajustar]
2. [Revisar aprendizados em ${LEARNINGS_FILE}]
3. [Aplicar melhorias no processo]

### Decisões Pendentes:
- [Decisão 1]
- [Decisão 2]

### Links:
- Aprendizados: \`${LEARNINGS_FILE}\`
- Checklist: \`docs/implementation_plans/${FEATURE_NAME}_pre_check.md\`
- Plano: \`docs/implementation_plans/${FEATURE_NAME}_implementation_*.md\`

EOF

echo "✅ Próximos passos atualizados: ${NEXT_STEPS_FILE}"

# 6. Limpar arquivos temporários (opcional)
echo ""
echo "🧹 LIMPANDO ARQUIVOS TEMPORÁRIOS..."
echo "=================================="

# Manter por padrão, mas oferecer opção
echo "Arquivos temporários:"
echo "  - ${TASK_FILE} (tarefa)"
SESSION_FILES=$(ls .session_* 2>/dev/null | wc -l)
echo "  - .session_* (${SESSION_FILES} arquivos de sessão)"
echo ""
echo "💡 Manter para referência? (s/N)"
read -t 5 -n 1 KEEP_FILES || KEEP_FILES="n"

if [[ "$KEEP_FILES" =~ ^[Nn]$ ]]; then
    rm -f ${TASK_FILE}
    echo "✅ Arquivo de tarefa removido"
else
    echo "✅ Arquivos mantidos para referência"
fi

# 7. Atualizar contexto para próxima sessão
echo ""
echo "🔄 ATUALIZANDO CONTEXTO..."
echo "=========================="

if [ -f "./scripts/update_context.sh" ]; then
    ./scripts/update_context.sh 2>/dev/null || echo "⚠️  Update context falhou, continuando..."
else
    echo "ℹ️  Script update_context.sh não encontrado"
fi

# 8. Resumo final
echo ""
echo "✅ CONCLUSÃO DA TAREFA COMPLETA!"
echo "================================"
echo ""
echo "📊 RESUMO:"
echo "----------"
echo "Tarefa: ${FEATURE_NAME} (${TASK_ID})"
echo "Status: ${STATUS}"
echo "Duração: ${DURATION_MIN} minutos"
echo "Testes: ${TEST_METRICS}"
echo ""
echo "📁 ARQUIVOS CRIADOS:"
echo "-------------------"
echo "1. Aprendizados: ${LEARNINGS_FILE}"
echo "2. Próximos passos: ${NEXT_STEPS_FILE}"
echo "3. Checklists atualizados"
echo ""
echo "🚀 PRÓXIMA SESSÃO:"
echo "-----------------"
echo "1. Iniciar: ./start_session.sh"
echo "2. Escolher próxima tarefa do backlog"
echo "3. Revisar aprendizados: cat ${LEARNINGS_FILE}"
echo ""
echo "💡 Dica: Os aprendizados desta tarefa serão usados para melhorar"
echo "       o processo da próxima implementação!"

exit 0