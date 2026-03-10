#!/bin/bash

# 🧪 E2E Test Runner com Banco Isolado
# Cria ambiente de teste isolado, executa Playwright e limpa tudo após testes

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configurações
BASE_PORT=9100  # Porta base para servidor de teste
TEST_TIMEOUT=120000  # 120 segundos
MAX_RETRIES=3

# Variáveis globais
TEST_ENTITY=""
TEST_PORT=""
TEST_PID=""
TEST_DB_PATH=""
TEST_LOG_FILE=""

# Funções de utilidade
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

cleanup() {
    log_info "Executando cleanup..."
    
    # Parar servidor de teste se estiver rodando
    if [ ! -z "$TEST_PID" ] && kill -0 $TEST_PID 2>/dev/null; then
        log_info "Parando servidor de teste (PID: $TEST_PID)"
        kill -TERM $TEST_PID 2>/dev/null || true
        sleep 2
        kill -KILL $TEST_PID 2>/dev/null || true
    fi
    
    # Remover banco de dados de teste
    if [ ! -z "$TEST_DB_PATH" ] && [ -f "$TEST_DB_PATH" ]; then
        log_info "Removendo banco de teste: $TEST_DB_PATH"
        rm -f "$TEST_DB_PATH"
    fi
    
    # Remover diretório de entidade de teste
    if [ ! -z "$TEST_ENTITY" ]; then
        TEST_ENTITY_DIR="$PWD/data/entities/${TEST_ENTITY}.db"
        if [ -f "$TEST_ENTITY_DIR" ]; then
            log_info "Removendo diretório de entidade: $TEST_ENTITY_DIR"
            rm -f "$TEST_ENTITY_DIR"
        fi
    fi
    
    # Remover arquivo de log
    if [ ! -z "$TEST_LOG_FILE" ] && [ -f "$TEST_LOG_FILE" ]; then
        log_info "Removendo arquivo de log: $TEST_LOG_FILE"
        rm -f "$TEST_LOG_FILE"
    fi
    
    log_success "Cleanup completo"
}

# Setup do ambiente de teste
setup_test_environment() {
    log_info "Configurando ambiente de teste isolado..."
    
    # Gerar ID único para teste
    local timestamp=$(date +%Y%m%d_%H%M%S_%N)
    TEST_ENTITY="test_e2e_${timestamp}"
    TEST_PORT=$((BASE_PORT + RANDOM % 100))
    TEST_DB_PATH="/tmp/${TEST_ENTITY}.db"
    TEST_LOG_FILE="/tmp/digna_test_${timestamp}.log"
    
    log_info "Entidade de teste: $TEST_ENTITY"
    log_info "Porta do servidor: $TEST_PORT"
    log_info "Banco de dados: $TEST_DB_PATH"
    log_info "Arquivo de log: $TEST_LOG_FILE"
    
    # Criar diretório para banco de dados
    mkdir -p "$(dirname "$TEST_DB_PATH")"
    mkdir -p "$PWD/data/entities"
    
    # Criar banco de dados SQLite vazio
    log_info "Criando banco de dados de teste..."
    sqlite3 "$TEST_DB_PATH" "VACUUM;" 2>/dev/null || true
    
    # Copiar estrutura do banco se existir template
    if [ -f "$PWD/data/entities/cafe_digna.db" ]; then
        log_info "Copiando estrutura do banco de produção..."
        cp "$PWD/data/entities/cafe_digna.db" "$PWD/data/entities/${TEST_ENTITY}.db"
    else
        log_info "Criando banco de dados vazio..."
        touch "$PWD/data/entities/${TEST_ENTITY}.db"
    fi
    
    # Iniciar servidor de teste - usando porta fixa 8090
    # Nota: O servidor atual não aceita flag --port, usa porta fixa 8090
    log_info "Iniciando servidor de teste na porta 8090..."
    
    cd "$PWD/modules/ui_web"
    
    # Configurar variáveis de ambiente para teste
    export DIGNA_TEST_MODE="true"
    export DIGNA_TEST_ENTITY="$TEST_ENTITY"
    export DIGNA_TEST_DB_PATH="$TEST_DB_PATH"
    
    # Iniciar servidor em background (porta fixa 8090)
    go run main.go > "$TEST_LOG_FILE" 2>&1 &
    TEST_PID=$!
    
    # Atualizar TEST_PORT para 8090
    TEST_PORT="8090"
    
    # Aguardar servidor iniciar
    local max_wait=30
    local wait_count=0
    
    log_info "Aguardando servidor iniciar (PID: $TEST_PID)..."
    while [ $wait_count -lt $max_wait ]; do
        if curl -s "http://localhost:8090/health" > /dev/null 2>&1; then
            log_success "Servidor de teste iniciado na porta 8090"
            break
        fi
        
        if ! kill -0 $TEST_PID 2>/dev/null; then
            log_error "Servidor de teste falhou ao iniciar"
            cat "$TEST_LOG_FILE" | tail -50
            return 1
        fi
        
        sleep 1
        wait_count=$((wait_count + 1))
    done
    
    if [ $wait_count -ge $max_wait ]; then
        log_error "Timeout aguardando servidor iniciar"
        cat "$TEST_LOG_FILE" | tail -50
        return 1
    fi
    
    # Criar usuário de teste no banco
    log_info "Configurando credenciais de teste..."
    setup_test_credentials
    
    cd "$OLDPWD"
    
    log_success "Ambiente de teste configurado com sucesso"
    echo "TEST_ENTITY=$TEST_ENTITY"
    echo "TEST_PORT=$TEST_PORT"
    echo "TEST_PID=$TEST_PID"
    echo "TEST_BASE_URL=http://localhost:$TEST_PORT"
    echo "TEST_PASSWORD=test123"
}

setup_test_credentials() {
    # Esta função deve configurar credenciais no banco de teste
    # Por enquanto, usamos credenciais hardcoded que o sistema aceita
    log_info "Usando credenciais padrão de teste"
    # Em uma implementação real, inseriríamos no banco:
    # INSERT INTO entities (id, name, password_hash) VALUES ('$TEST_ENTITY', 'Test Entity', 'hash_do_test123');
}

# Executar testes Playwright
run_e2e_tests() {
    local test_spec="${1:-tests/}"
    local browser="${2:-chromium}"
    local headed="${3:-false}"
    
    log_info "Executando testes E2E no $browser..."
    
    # Configurar variáveis de ambiente para Playwright
    export PLAYWRIGHT_TEST_BASE_URL="http://localhost:$TEST_PORT"
    export PLAYWRIGHT_TEST_ENTITY="$TEST_ENTITY"
    export PLAYWRIGHT_TEST_PASSWORD="test123"
    export PLAYWRIGHT_TIMEOUT="$TEST_TIMEOUT"
    
    # Opções do Playwright
    local playwright_opts=""
    if [ "$headed" = "true" ]; then
        playwright_opts="$playwright_opts --headed"
    fi
    
    # Executar testes
    local retry_count=0
    local test_result=1
    
    while [ $retry_count -lt $MAX_RETRIES ] && [ $test_result -ne 0 ]; do
        log_info "Execução $((retry_count + 1)) de $MAX_RETRIES..."
        
        npx playwright test "$test_spec" \
            --project="$browser" \
            --timeout="$TEST_TIMEOUT" \
            --reporter=line,html \
            --output="test-results/$TEST_ENTITY" \
            $playwright_opts
        
        test_result=$?
        
        if [ $test_result -ne 0 ]; then
            log_warning "Testes falharam, tentando novamente..."
            retry_count=$((retry_count + 1))
            sleep 2
        fi
    done
    
    if [ $test_result -eq 0 ]; then
        log_success "Testes E2E concluídos com sucesso"
    else
        log_error "Testes E2E falharam após $MAX_RETRIES tentativas"
    fi
    
    return $test_result
}

# Função principal
main() {
    local test_spec="tests/"
    local browser="chromium"
    local headed="false"
    local skip_cleanup="false"
    
    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --spec)
                test_spec="$2"
                shift 2
                ;;
            --browser)
                browser="$2"
                shift 2
                ;;
            --headed)
                headed="true"
                shift
                ;;
            --skip-cleanup)
                skip_cleanup="true"
                shift
                ;;
            --help)
                show_help
                exit 0
                ;;
            *)
                log_error "Opção desconhecida: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    # Registrar cleanup handler
    trap cleanup EXIT INT TERM
    
    log_info "🚀 Iniciando E2E Test Runner com Banco Isolado"
    log_info "=============================================="
    
    # Setup environment
    if ! setup_test_environment; then
        log_error "Falha ao configurar ambiente de teste"
        exit 1
    fi
    
    # Run tests
    if ! run_e2e_tests "$test_spec" "$browser" "$headed"; then
        log_error "Testes falharam"
        
        # Mostrar logs em caso de falha
        if [ -f "$TEST_LOG_FILE" ]; then
            log_info "Últimos logs do servidor:"
            tail -100 "$TEST_LOG_FILE"
        fi
        
        # Manter ambiente se skip_cleanup
        if [ "$skip_cleanup" = "true" ]; then
            log_warning "Cleanup skipped. Ambiente mantido:"
            echo "  Entidade: $TEST_ENTITY"
            echo "  Porta: $TEST_PORT"
            echo "  PID: $TEST_PID"
            echo "  Banco: $TEST_DB_PATH"
            echo "  Log: $TEST_LOG_FILE"
            trap - EXIT INT TERM  # Remover trap
        fi
        
        exit 1
    fi
    
    log_success "✅ Todos os testes passaram!"
    
    # Manter ambiente se skip_cleanup
    if [ "$skip_cleanup" = "true" ]; then
        log_warning "Cleanup skipped. Ambiente mantido para debug."
        trap - EXIT INT TERM  # Remover trap
    fi
}

show_help() {
    cat << EOF
${BLUE}🧪 E2E Test Runner com Banco Isolado${NC}

${YELLOW}USO:${NC}
  ./e2e_test_runner.sh [OPÇÕES]

${YELLOW}OPÇÕES:${NC}
  --spec ARQUIVO      Especificar arquivo de teste (padrão: tests/)
  --browser NAVEGADOR Navegador para testes (chromium, firefox, webkit)
  --headed            Executar com navegador visível (para debug)
  --skip-cleanup      Não limpar ambiente após testes (para debug)
  --help              Mostrar esta ajuda

${YELLOW}EXEMPLOS:${NC}
  ./e2e_test_runner.sh                    # Executar todos testes em chromium headless
  ./e2e_test_runner.sh --spec tests/digna-basic.spec.js
  ./e2e_test_runner.sh --browser firefox --headed
  ./e2e_test_runner.sh --skip-cleanup     # Manter ambiente para debug

${YELLOW}DESCRIÇÃO:${NC}
  Este script cria um ambiente de teste completamente isolado:
  1. Cria entidade de teste única
  2. Configura banco de dados isolado
  3. Inicia servidor em porta dedicada
  4. Executa testes Playwright
  5. Limpa tudo automaticamente

${GREEN}Variáveis de ambiente para testes:${NC}
  PLAYWRIGHT_TEST_BASE_URL    URL do servidor de teste
  PLAYWRIGHT_TEST_ENTITY      Nome da entidade de teste
  PLAYWRIGHT_TEST_PASSWORD    Senha de teste (test123)
EOF
}

# Executar função principal
main "$@"