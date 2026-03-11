#!/bin/bash
# validate_task_tests.sh - Valida se uma tarefa tem testes adequados
# Uso: ./scripts/validate_task_tests.sh --task=TASK_ID

set -e

echo "🧪 VALIDAÇÃO DE TESTES PARA TAREFA"
echo "=================================="

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
            echo "Uso: ./scripts/validate_task_tests.sh --task=TASK_ID [--force]"
            echo ""
            echo "Valida se uma tarefa tem testes adequados antes de conclusão."
            echo ""
            echo "Opções:"
            echo "  --task=TASK_ID    ID da tarefa (obrigatório)"
            echo "  --force           Ignorar avisos e continuar"
            echo "  --help, -h        Mostrar esta ajuda"
            exit 0
            ;;
    esac
done

if [ -z "$TASK_ID" ]; then
    echo "❌ ID da tarefa é obrigatório."
    echo "💡 Use: ./scripts/validate_task_tests.sh --task=TASK_ID"
    exit 1
fi

TASK_DIR="work_in_progress/tasks/task_${TASK_ID}"
if [ ! -d "$TASK_DIR" ]; then
    echo "❌ Tarefa não encontrada: ${TASK_DIR}"
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

# 1. Verificar arquivos de teste criados/modificados após a tarefa
echo "1. 📁 ARQUIVOS DE TESTE CRIADOS/MODIFICADOS:"
echo "--------------------------------------------"

TEST_FILES=$(find modules -name "*test*.go" -newer "${TASK_DIR}/task_prompt.md" 2>/dev/null || true)
TEST_COUNT=$(echo "$TEST_FILES" | wc -w)

if [ "$TEST_COUNT" -eq 0 ]; then
    echo "   ❌ NENHUM arquivo de teste criado/modificado"
    echo "   🚨 CRÍTICO: Tarefa deve incluir testes"
else
    echo "   ✅ $TEST_COUNT arquivo(s) de teste:"
    echo "$TEST_FILES" | while read file; do
        echo "      - $file"
    done
fi

# 2. Verificar testes E2E (Playwright)
echo ""
echo "2. 🎭 TESTES E2E (PLAYWRIGHT):"
echo "-----------------------------"

E2E_TESTS=$(find modules -name "*e2e*test*.go" -o -name "*playwright*test*.go" 2>/dev/null || true)
E2E_COUNT=$(echo "$E2E_TESTS" | wc -w)

if [ "$E2E_COUNT" -eq 0 ]; then
    echo "   ⚠️  NENHUM teste E2E encontrado"
    echo "   💡 Recomendado: Crie testes E2E para validação completa no browser"
else
    echo "   ✅ $E2E_COUNT teste(s) E2E encontrado(s)"
fi

# 3. Verificar cobertura de rotas
echo ""
echo "3. 🌐 COBERTURA DE ROTAS:"
echo "------------------------"

# Extrair rotas do handler (simplificado)
if [ -f "modules/ui_web/internal/handler/${TASK_NAME,,}_handler.go" ]; then
    HANDLER_FILE="modules/ui_web/internal/handler/${TASK_NAME,,}_handler.go"
    ROUTES=$(grep -o '"/[^"]*"' "$HANDLER_FILE" 2>/dev/null || true)
    
    if [ -n "$ROUTES" ]; then
        echo "   ✅ Rotas identificadas no handler:"
        echo "$ROUTES" | while read route; do
            echo "      - $route"
        done
    else
        echo "   ℹ️  Não foi possível identificar rotas no handler"
    fi
else
    echo "   ℹ️  Handler não encontrado (pode ser biblioteca)"
fi

# 4. Executar testes existentes
echo ""
echo "4. 🚀 EXECUTANDO TESTES EXISTENTES:"
echo "----------------------------------"

cd modules 2>/dev/null
if [ $? -eq 0 ]; then
    echo "   Executando testes de sistema..."
    TEST_OUTPUT=$(timeout 30 ./run_tests.sh 2>&1 | tail -20 || echo "Testes falharam ou timeout")
    
    if echo "$TEST_OUTPUT" | grep -q "PASS\|ok"; then
        echo "   ✅ Testes de sistema PASSAM"
    else
        echo "   ❌ Testes de sistema FALHARAM ou não executados"
        echo "   📝 Saída:"
        echo "$TEST_OUTPUT" | while read line; do
            echo "      $line"
        done
    fi
    cd - >/dev/null 2>&1
else
    echo "   ℹ️  Diretório modules não encontrado"
fi

# 5. Resumo e recomendação
echo ""
echo "5. 📊 RESUMO E RECOMENDAÇÃO:"
echo "---------------------------"

if [ "$TEST_COUNT" -eq 0 ]; then
    echo "   ❌❌❌ TAREFA NÃO PODE SER CONCLUÍDA ❌❌❌"
    echo "   Motivo: Nenhum teste criado"
    echo ""
    echo "   💡 AÇÕES NECESSÁRIAS:"
    echo "   1. Crie testes unitários para handlers/services"
    echo "   2. Crie testes E2E com Playwright"
    echo "   3. Teste manualmente todas as rotas"
    echo "   4. Execute ./scripts/validate_task_tests.sh novamente"
    
    if [ "$FORCE" = false ]; then
        exit 1
    else
        echo "   ⚠️  Continuando em modo FORCE (ignorando avisos)"
    fi
elif [ "$TEST_COUNT" -lt 3 ]; then
    echo "   ⚠️  TAREFA COM TESTES INSUFICIENTES"
    echo "   Motivo: Apenas $TEST_COUNT teste(s) criado(s)"
    echo ""
    echo "   💡 RECOMENDAÇÕES:"
    echo "   1. Adicione mais testes unitários"
    echo "   2. Considere criar testes E2E"
    echo "   3. Teste casos de erro/borda"
    
    if [ "$FORCE" = false ]; then
        echo ""
        echo "   ❓ Deseja continuar mesmo com testes insuficientes?"
        echo "      Use --force para ignorar"
        exit 1
    fi
else
    echo "   ✅ TAREFA COM TESTES ADEQUADOS"
    echo "   Motivo: $TEST_COUNT teste(s) criado(s)"
    echo ""
    echo "   🎯 PRÓXIMO PASSO:"
    echo "   Execute: ./conclude_task.sh --task=${TASK_ID} \"Aprendizados: ...\""
fi

echo ""
echo "📚 COMANDOS ÚTEIS PARA CRIAR TESTES:"
echo "-----------------------------------"
echo "• Testes unitários: go test -v ./modules/ui_web/internal/handler/"
echo "• Testes E2E: cd modules/ui_web && go test -v ./e2e_*test.go"
echo "• Smoke test: ./scripts/dev/smoke_test_new_feature.sh"
echo "• Validar novamente: ./scripts/validate_task_tests.sh --task=${TASK_ID}"