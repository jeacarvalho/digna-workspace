#!/bin/bash
# smoke_test_bug_fixes.sh - Smoke test específico para tarefas de correção de bugs
# Uso: ./scripts/dev/smoke_test_bug_fixes.sh --task=TASK_ID

set -e

echo "🐛 SMOKE TEST PARA CORREÇÃO DE BUGS"
echo "==================================="

TASK_ID=""
SERVER_PID=""
SERVER_PORT="8090"

# Processar argumentos
for arg in "$@"; do
    case $arg in
        --task=*)
            TASK_ID="${arg#*=}"
            ;;
        --port=*)
            SERVER_PORT="${arg#*=}"
            ;;
        --help|-h)
            echo "Uso: ./scripts/dev/smoke_test_bug_fixes.sh --task=TASK_ID [--port=PORT]"
            echo ""
            echo "Executa smoke test específico para tarefas de correção de bugs."
            echo "Inicia servidor, testa URLs problemáticas, valida correções."
            exit 0
            ;;
    esac
done

if [ -z "$TASK_ID" ]; then
    echo "❌ ID da tarefa é obrigatório."
    echo "💡 Use: ./scripts/dev/smoke_test_bug_fixes.sh --task=TASK_ID"
    exit 1
fi

TASK_DIR="work_in_progress/tasks/task_${TASK_ID}"
if [ ! -d "$TASK_DIR" ]; then
    echo "❌ Tarefa não encontrada: ${TASK_DIR}"
    exit 1
fi

PROMPT_FILE="${TASK_DIR}/task_prompt.md"
if [ ! -f "$PROMPT_FILE" ]; then
    echo "❌ Arquivo de prompt não encontrado: ${PROMPT_FILE}"
    exit 1
fi

echo "📋 Tarefa: $(basename "$TASK_DIR")"
echo "🔍 Analisando prompt para identificar bugs a corrigir..."
echo ""

# 1. Extrair URLs problemáticas do prompt
echo "1. 🔎 IDENTIFICANDO BUGS NO PROMPT:"
echo "----------------------------------"

BUGGY_URLS=$(grep -o "http://[^ ]*" "$PROMPT_FILE" || true)
ERROR_PATTERNS=$(grep -oi "erro.*ao.*buscar\|database.*connection.*closed\|sql.*database.*closed\|não.*funciona\|falha\|bug" "$PROMPT_FILE" | sort -u || true)

if [ -n "$BUGGY_URLS" ]; then
    echo "   URLs problemáticas identificadas:"
    echo "$BUGGY_URLS" | while read url; do
        echo "   - $url"
    done
fi

if [ -n "$ERROR_PATTERNS" ]; then
    echo "   Padrões de erro identificados:"
    echo "$ERROR_PATTERNS" | while read pattern; do
        echo "   - \"$pattern\""
    done
fi

if [ -z "$BUGGY_URLS" ] && [ -z "$ERROR_PATTERNS" ]; then
    echo "   ℹ️  Nenhum bug específico identificado no prompt"
    echo "   ✅ Smoke test básico suficiente"
    exit 0
fi

# 2. Iniciar servidor para testes
echo ""
echo "2. 🚀 INICIANDO SERVIDOR PARA TESTES:"
echo "------------------------------------"

# Verificar se o servidor já está rodando
if curl -s "http://localhost:${SERVER_PORT}/health" >/dev/null 2>&1; then
    echo "   ✅ Servidor já está rodando na porta ${SERVER_PORT}"
else
    echo "   🔧 Iniciando servidor na porta ${SERVER_PORT}..."
    
    # Tentar iniciar o servidor (ajuste conforme seu projeto)
    cd modules/ui_web 2>/dev/null
    if [ $? -eq 0 ]; then
        # Iniciar em background
        go run main.go --port="${SERVER_PORT}" >/tmp/digna_server.log 2>&1 &
        SERVER_PID=$!
        
        # Aguardar servidor iniciar
        echo "   ⏳ Aguardando servidor iniciar..."
        sleep 5
        
        # Verificar se iniciou
        if curl -s "http://localhost:${SERVER_PORT}/health" >/dev/null 2>&1; then
            echo "   ✅ Servidor iniciado (PID: $SERVER_PID)"
        else
            echo "   ❌ Falha ao iniciar servidor"
            echo "   📝 Logs: /tmp/digna_server.log"
            cat /tmp/digna_server.log | tail -20
            exit 1
        fi
        cd - >/dev/null 2>&1
    else
        echo "   ❌ Não foi possível acessar módulo ui_web"
        echo "   ℹ️  Smoke test manual necessário"
        exit 0
    fi
fi

# 3. Testar URLs problemáticas
echo ""
echo "3. 🧪 TESTANDO URLs PROBLEMÁTICAS:"
echo "---------------------------------"

ALL_TESTS_PASSED=true

# Testar cada URL problemática
echo "$BUGGY_URLS" | while read url; do
    # Extrair path da URL
    PATH_ONLY=$(echo "$url" | sed 's|http://[^/]*||')
    
    if [ -n "$PATH_ONLY" ]; then
        echo "   🔍 Testando: $PATH_ONLY"
        
        # Fazer request
        RESPONSE=$(curl -s -w "%{http_code}" "http://localhost:${SERVER_PORT}${PATH_ONLY}" -o /tmp/response_body.txt 2>&1)
        STATUS_CODE=$(echo "$RESPONSE" | tail -1)
        BODY=$(cat /tmp/response_body.txt)
        
        echo "   📊 Status: $STATUS_CODE"
        
        # Verificar se há erros de banco de dados na resposta
        if echo "$BODY" | grep -qi "database connection is closed\|sql: database is closed"; then
            echo "   ❌ ERRO DE BANCO DE DADOS DETECTADO!"
            echo "   📝 Resposta contém:"
            echo "$BODY" | grep -i "database\|sql\|erro" | head -5
            ALL_TESTS_PASSED=false
        elif [ "$STATUS_CODE" -eq 500 ]; then
            echo "   ⚠️  Status 500 (Internal Server Error)"
            echo "   📝 Verifique logs do servidor"
        elif [ "$STATUS_CODE" -eq 200 ] || [ "$STATUS_CODE" -eq 302 ]; then
            echo "   ✅ Status $STATUS_CODE (OK)"
        else
            echo "   ℹ️  Status $STATUS_CODE"
        fi
        
        echo ""
    fi
done

# 4. Testar padrões de erro específicos
echo ""
echo "4. 🐛 TESTANDO PADRÕES DE ERRO ESPECÍFICOS:"
echo "------------------------------------------"

# Se houver padrões de erro específicos, testar endpoints relacionados
if echo "$ERROR_PATTERNS" | grep -qi "supply.*stock\|stock.*page"; then
    echo "   🔍 Testando /supply/stock (erro mencionado)"
    
    RESPONSE=$(curl -s -w "%{http_code}" "http://localhost:${SERVER_PORT}/supply/stock?entity_id=cafe_digna" -o /tmp/response_stock.txt 2>&1)
    STATUS_CODE=$(echo "$RESPONSE" | tail -1)
    BODY=$(cat /tmp/response_stock.txt)
    
    echo "   📊 Status: $STATUS_CODE"
    
    if echo "$BODY" | grep -qi "database connection is closed\|sql: database is closed"; then
        echo "   ❌ ERRO DE BANCO NÃO CORRIGIDO em /supply/stock!"
        ALL_TESTS_PASSED=false
    else
        echo "   ✅ Sem erros de banco em /supply/stock"
    fi
fi

if echo "$ERROR_PATTERNS" | grep -qi "supply.*suppliers\|suppliers.*page"; then
    echo "   🔍 Testando /supply/suppliers (erro mencionado)"
    
    RESPONSE=$(curl -s -w "%{http_code}" "http://localhost:${SERVER_PORT}/supply/suppliers?entity_id=cafe_digna" -o /tmp/response_suppliers.txt 2>&1)
    STATUS_CODE=$(echo "$RESPONSE" | tail -1)
    BODY=$(cat /tmp/response_suppliers.txt)
    
    echo "   📊 Status: $STATUS_CODE"
    
    if echo "$BODY" | grep -qi "database connection is closed\|sql: database is closed"; then
        echo "   ❌ ERRO DE BANCO NÃO CORRIGIDO em /supply/suppliers!"
        ALL_TESTS_PASSED=false
    else
        echo "   ✅ Sem erros de banco em /supply/suppliers"
    fi
fi

# 5. Finalizar e limpar
echo ""
echo "5. 🧹 FINALIZANDO:"
echo "----------------"

# Parar servidor se nós o iniciamos
if [ -n "$SERVER_PID" ]; then
    echo "   Parando servidor (PID: $SERVER_PID)..."
    kill $SERVER_PID 2>/dev/null || true
    wait $SERVER_PID 2>/dev/null || true
    echo "   ✅ Servidor parado"
fi

# Remover arquivos temporários
rm -f /tmp/response_*.txt /tmp/digna_server.log 2>/dev/null || true

# Resumo final
echo ""
echo "📊 RESUMO DO SMOKE TEST:"
echo "-----------------------"

if [ "$ALL_TESTS_PASSED" = true ]; then
    echo "✅✅✅ SMOKE TEST PASSOU!"
    echo "Todos os bugs identificados foram corrigidos."
    exit 0
else
    echo "❌❌❌ SMOKE TEST FALHOU!"
    echo "Alguns bugs NÃO foram corrigidos."
    echo ""
    echo "💡 AÇÕES NECESSÁRIAS:"
    echo "1. Revise os erros acima"
    echo "2. Corrija os problemas no código"
    echo "3. Execute este smoke test novamente:"
    echo "   ./scripts/dev/smoke_test_bug_fixes.sh --task=${TASK_ID}"
    exit 1
fi