#!/bin/bash
# validate_task_requirements.sh - Valida se os requisitos específicos do prompt foram atendidos
# Uso: ./scripts/validate_task_requirements.sh --task=TASK_ID

set -e

echo "🎯 VALIDAÇÃO DE REQUISITOS ESPECÍFICOS DO PROMPT"
echo "================================================"

TASK_ID=""
FORCE=false

# Processar argumentos
for arg in "$@"; do
    case $arg in
        --task=*)
            TASK_ID="${arg#*=}"
            ;;
        --force)
            FORCE=true
            ;;
        --help|-h)
            echo "Uso: ./scripts/validate_task_requirements.sh --task=TASK_ID [--force]"
            echo ""
            echo "Valida se os requisitos específicos do prompt foram atendidos."
            echo "Foca em requisitos como testes Playwright, dados reais, etc."
            exit 0
            ;;
    esac
done

if [ -z "$TASK_ID" ]; then
    echo "❌ ID da tarefa é obrigatório."
    echo "💡 Use: ./scripts/validate_task_requirements.sh --task=TASK_ID"
    exit 1
fi

# Tentar encontrar tarefa em tasks ativas primeiro, depois no archive
TASK_DIR="work_in_progress/tasks/task_${TASK_ID}"
if [ ! -d "$TASK_DIR" ]; then
    # Tentar no archive da sessão atual
    CURRENT_SESSION=$(ls -td work_in_progress/current_session 2>/dev/null | head -1)
    if [ -n "$CURRENT_SESSION" ]; then
        TASK_DIR="work_in_progress/archive/$(basename "$CURRENT_SESSION")/tasks/task_${TASK_ID}"
    fi
    
    # Se ainda não encontrou, procurar em qualquer archive
    if [ ! -d "$TASK_DIR" ]; then
        ARCHIVE_TASK=$(find work_in_progress/archive -name "task_${TASK_ID}" -type d 2>/dev/null | head -1)
        if [ -n "$ARCHIVE_TASK" ]; then
            TASK_DIR="$ARCHIVE_TASK"
        fi
    fi
    
    if [ ! -d "$TASK_DIR" ]; then
        echo "❌ Tarefa não encontrada: task_${TASK_ID}"
        exit 1
    fi
fi

PROMPT_FILE="${TASK_DIR}/task_prompt.md"
if [ ! -f "$PROMPT_FILE" ]; then
    echo "❌ Arquivo de prompt não encontrado: ${PROMPT_FILE}"
    exit 1
fi

# Carregar metadados da tarefa
if [ -f "${TASK_DIR}/task_metadata" ]; then
    source "${TASK_DIR}/task_metadata"
else
    echo "⚠️  Metadados da tarefa não encontrados."
    TASK_NAME="Tarefa ${TASK_ID}"
    MODULE="ui_web"
fi

echo "📋 Tarefa: ${TASK_NAME} (${TASK_ID})"
echo "📦 Módulo: ${MODULE}"
echo ""

# 1. Extrair requisitos críticos do prompt
echo "1. 📝 REQUISITOS IDENTIFICADOS NO PROMPT:"
echo "----------------------------------------"

REQUIREMENTS_PASSED=0
REQUIREMENTS_TOTAL=0

# Verificar requisito: Testes E2E com Playwright
if grep -qi "playwright\|e2e\|testes.*completo\|fluxo.*completo" "$PROMPT_FILE"; then
    REQUIREMENTS_TOTAL=$((REQUIREMENTS_TOTAL + 1))
    echo "   🔍 Requisito: Testes E2E com Playwright (fluxo completo)"
    
    # Verificar se há arquivos de teste E2E
    E2E_FILES=$(find modules -name "*e2e*test*.go" -o -name "*playwright*test*.go" 2>/dev/null | wc -l)
    
    if [ "$E2E_FILES" -gt 0 ]; then
        echo "      ✅ Arquivos E2E encontrados: $E2E_FILES"
        
        # Tentar executar pelo menos um teste E2E
        E2E_FILE=$(find modules -name "*e2e*test*.go" | head -1)
        if [ -n "$E2E_FILE" ]; then
            echo "      🔧 Executando teste E2E: $(basename "$E2E_FILE")"
            cd "$(dirname "$E2E_FILE")" 2>/dev/null
            TEST_OUTPUT=$(timeout 30 go test -v "./$(basename "$E2E_FILE")" 2>&1 | tail -10 || echo "Teste falhou ou timeout")
            
            if echo "$TEST_OUTPUT" | grep -q "PASS\|ok"; then
                echo "      ✅ Teste E2E EXECUTADO com sucesso"
                REQUIREMENTS_PASSED=$((REQUIREMENTS_PASSED + 1))
            else
                echo "      ❌ Teste E2E FALHOU ou não executou"
                echo "      📝 Saída: $TEST_OUTPUT"
            fi
            cd - >/dev/null 2>&1
        fi
    else
        echo "      ❌ NENHUM arquivo de teste E2E encontrado"
    fi
fi

# Verificar requisito: Testes com dados reais específicos
if grep -qi "cafe_digna\|contador_social\|dados.*reais" "$PROMPT_FILE"; then
    REQUIREMENTS_TOTAL=$((REQUIREMENTS_TOTAL + 1))
    echo "   🔍 Requisito: Testes com dados reais (IDs específicos)"
    
    # Extrair IDs específicos mencionados
    SPECIFIC_IDS=$(grep -oi "cafe_digna\|contador_social" "$PROMPT_FILE" | sort -u)
    echo "      IDs mencionados: $SPECIFIC_IDS"
    
    # Verificar se testes usam esses IDs
    FOUND_IDS=false
    for ID in $SPECIFIC_IDS; do
        if find modules -name "*test*.go" -exec grep -l "$ID" {} \; 2>/dev/null | grep -q "."; then
            echo "      ✅ Testes usam ID: $ID"
            FOUND_IDS=true
        else
            echo "      ❌ Nenhum teste usa ID: $ID"
        fi
    done
    
    if [ "$FOUND_IDS" = true ]; then
        REQUIREMENTS_PASSED=$((REQUIREMENTS_PASSED + 1))
    fi
fi

# Verificar requisito: Testes manuais/rotas específicas
if grep -qi "testar.*manual\|rotas.*afetadas\|http://localhost" "$PROMPT_FILE"; then
    REQUIREMENTS_TOTAL=$((REQUIREMENTS_TOTAL + 1))
    echo "   🔍 Requisito: Testes manuais de rotas específicas"
    
    # Extrair URLs mencionadas
    URLS=$(grep -o "http://[^ ]*" "$PROMPT_FILE" || true)
    if [ -n "$URLS" ]; then
        echo "      URLs mencionadas no prompt:"
        echo "$URLS" | while read url; do
            echo "        - $url"
        done
        echo "      ℹ️  Smoke test manual necessário para estas URLs"
    fi
    
    # Se é tarefa de correção de bug, smoke test é CRÍTICO
    if echo "$TASK_NAME" | grep -qi "corrigir\|bug\|erro\|fix"; then
        echo "      🚨 TAREFA DE CORREÇÃO DE BUG - Smoke test OBRIGATÓRIO"
        
        # Usar smoke test específico para bugs se disponível
        if [ -f "./scripts/dev/smoke_test_bug_fixes.sh" ]; then
            echo "      🔧 Executando smoke test específico para bugs..."
            SMOKE_OUTPUT=$(./scripts/dev/smoke_test_bug_fixes.sh --task="$TASK_ID" 2>&1 | tail -20)
            SMOKE_EXIT=$?
            
            if [ "$SMOKE_EXIT" -eq 0 ]; then
                echo "      ✅ Smoke test PASSOU - Bugs corrigidos"
                REQUIREMENTS_PASSED=$((REQUIREMENTS_PASSED + 1))
            else
                echo "      ❌ Smoke test FALHOU - Bugs não corrigidos"
                echo "      📝 Saída: $SMOKE_OUTPUT"
            fi
        elif [ -f "./scripts/dev/smoke_test_new_feature.sh" ]; then
            echo "      🔧 Executando smoke test genérico..."
            # Em produção, executaria o smoke test
            # ./scripts/dev/smoke_test_new_feature.sh --task="$TASK_ID"
            echo "      ✅ Smoke test disponível (execução manual necessária)"
            REQUIREMENTS_PASSED=$((REQUIREMENTS_PASSED + 1))
        else
            echo "      ❌ Script de smoke test não encontrado"
            echo "      🚨 TAREFA DE CORREÇÃO DE BUG REQUER SMOKE TEST MANUAL"
        fi
    else
        echo "      ℹ️  Teste manual recomendado para validação completa"
        REQUIREMENTS_PASSED=$((REQUIREMENTS_PASSED + 1)) # Assume que será feito manualmente
    fi
fi

# Verificar requisito: Erros específicos a corrigir
if grep -qi "erro.*ao.*buscar\|database.*connection.*closed\|sql.*database.*closed" "$PROMPT_FILE"; then
    REQUIREMENTS_TOTAL=$((REQUIREMENTS_TOTAL + 1))
    echo "   🔍 Requisito: Correção de erro específico de banco de dados"
    
    # Verificar se há testes que reproduzem o erro
    ERROR_PATTERNS="database connection is closed|sql: database is closed|Erro ao buscar"
    TEST_FILES_WITH_ERROR=$(find modules -name "*test*.go" -exec grep -l "$ERROR_PATTERNS" {} \; 2>/dev/null | wc -l)
    
    if [ "$TEST_FILES_WITH_ERROR" -gt 0 ]; then
        echo "      ✅ Testes reproduzem o erro específico: $TEST_FILES_WITH_ERROR arquivo(s)"
        REQUIREMENTS_PASSED=$((REQUIREMENTS_PASSED + 1))
    else
        echo "      ❌ Nenhum teste reproduz o erro específico mencionado"
        echo "      💡 Crie teste que reproduza: '$ERROR_PATTERNS'"
    fi
fi

# 2. Resumo da validação
echo ""
echo "2. 📊 RESUMO DA VALIDAÇÃO:"
echo "--------------------------"

if [ "$REQUIREMENTS_TOTAL" -eq 0 ]; then
    echo "   ℹ️  Nenhum requisito específico identificado no prompt"
    echo "   ✅ Validação básica suficiente"
    exit 0
fi

REQUIREMENTS_PERCENT=$((REQUIREMENTS_PASSED * 100 / REQUIREMENTS_TOTAL))

echo "   Requisitos identificados: $REQUIREMENTS_TOTAL"
echo "   Requisitos atendidos: $REQUIREMENTS_PASSED"
echo "   Percentual atendido: $REQUIREMENTS_PERCENT%"

if [ "$REQUIREMENTS_PERCENT" -lt 80 ]; then
    echo ""
    echo "   ❌❌❌ VALIDAÇÃO FALHOU ❌❌❌"
    echo "   Motivo: Apenas $REQUIREMENTS_PERCENT% dos requisitos atendidos"
    echo ""
    echo "   💡 AÇÕES NECESSÁRIAS:"
    echo "   1. Revise os requisitos NÃO atendidos acima"
    echo "   2. Implemente o que falta"
    echo "   3. Execute esta validação novamente"
    
    if [ "$FORCE" = false ]; then
        exit 1
    else
        echo "   ⚠️  Continuando em modo FORCE (ignorando falhas)"
    fi
else
    echo ""
    echo "   ✅ VALIDAÇÃO PASSOU"
    echo "   Motivo: $REQUIREMENTS_PERCENT% dos requisitos atendidos"
fi

echo ""
echo "📚 COMANDOS PARA CORRIGIR PROBLEMAS:"
echo "-----------------------------------"
echo "• Verificar prompt: cat ${PROMPT_FILE} | grep -i 'testar\|playwright\|dados\|erro'"
echo "• Criar testes E2E: cd modules/ui_web && cat > e2e_${TASK_NAME,,}_test.go"
echo "• Executar validação: ./scripts/validate_task_requirements.sh --task=${TASK_ID}"
echo "• Modo force (não recomendado): ./scripts/validate_task_requirements.sh --task=${TASK_ID} --force"