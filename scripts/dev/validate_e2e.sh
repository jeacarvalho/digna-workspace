#!/bin/bash

# 🧪 Script de Validação E2E com Playwright
# Valida fluxos completos de negócio após implementação
# Modo stealth (headless) por padrão - não abre janelas no desktop

set -e

# Configurações
BASE_URL="http://localhost:8090"
TEST_ENTITY="cafe_digna"
TEST_PASSWORD="cd0123"
TIMEOUT=60000  # 60 segundos

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Funções de ajuda
show_help() {
    cat << EOF
${BLUE}🧪 validate_e2e.sh - Validação End-to-End com Playwright${NC}

${YELLOW}USO:${NC}
  ./validate_e2e.sh [OPÇÕES]

${YELLOW}OPÇÕES:${NC}
  --basic, -b        Validação básica (7 passos padrão) - ${GREEN}RECOMENDADO${NC}
  --full, -f         Todos os testes E2E (inclui testes completos)
  --custom, -c       Teste customizado (especificar teste)
  --ui, -u           Modo UI (abre navegador - para debug)
  --headed           Modo headed (abre navegador visível)
  --headless         Modo headless (padrão - stealth mode)
  --chrome, -C       Executar apenas no Chrome
  --firefox, -F      Executar apenas no Firefox
  --timeout N        Timeout em segundos (padrão: 60)
  --help, -h         Mostrar esta ajuda

${YELLOW}MODOS DE VALIDAÇÃO:${NC}

${GREEN}1. --basic (7 passos padrão Digna)${NC}
   ✅ Login no sistema
   ✅ Criar item de estoque (se não existir)
   ✅ Criar membro (se não existir)  
   ✅ Criar fornecedor (se não existir)
   ✅ Registrar compra do item
   ✅ Registrar venda no PDV
   ✅ Confirmar saldo e registrar horas

${GREEN}2. --full (todos os testes)${NC}
   Todos os testes em tests/digna-*.spec.js

${GREEN}3. --custom "nome do teste"${NC}
   Executa teste específico: "Login e navegação básica"

${YELLOW}EXEMPLOS:${NC}
  # Validação básica (stealth mode - padrão)
  ./validate_e2e.sh --basic

  # Validação completa com UI para debug
  ./validate_e2e.sh --full --ui

  # Teste específico com timeout maior
  ./validate_e2e.sh --custom "Fluxo simplificado" --timeout 90

  # Apenas Chrome em modo headless
  ./validate_e2e.sh --basic --chrome --headless

${YELLOW}SAÍDA:${NC}
  ✅ Sucesso: Retorna 0 e mostra relatório
  ❌ Falha: Retorna 1 e mostra detalhes do erro
  ⚠️  Aviso: Retorna 2 (problemas não críticos)

${YELLOW}INTEGRAÇÃO COM WORKFLOW:${NC}
  Use após implementação no opencode, antes de ./conclude_task.sh
EOF
    exit 0
}

# Verificar dependências
check_dependencies() {
    echo -e "${BLUE}🔍 Verificando dependências...${NC}"
    
    # Verificar servidor
    if ! curl -s "${BASE_URL}/health" > /dev/null; then
        echo -e "${RED}❌ Servidor Digna não está rodando em ${BASE_URL}${NC}"
        echo -e "💡 Execute: cd modules/ui_web && go run main.go"
        return 1
    fi
    echo -e "${GREEN}✅ Servidor rodando${NC}"
    
    # Verificar Playwright
    if ! command -v npx > /dev/null; then
        echo -e "${RED}❌ Node/npx não encontrado${NC}"
        return 1
    fi
    
    if [ ! -f "package.json" ]; then
        echo -e "${RED}❌ package.json não encontrado${NC}"
        return 1
    fi
    
    if ! npx playwright --version > /dev/null 2>&1; then
        echo -e "${YELLOW}⚠️  Playwright não instalado. Instalando...${NC}"
        npm install --save-dev @playwright/test
        npx playwright install
    fi
    
    echo -e "${GREEN}✅ Playwright disponível${NC}"
    return 0
}

# Executar testes Playwright
run_playwright() {
    local test_file="$1"
    local test_name="$2"
    local mode="$3"
    local browser="$4"
    
    echo -e "${BLUE}🚀 Executando testes E2E...${NC}"
    echo -e "📁 Teste: ${test_name}"
    echo -e "👁️  Modo: ${mode}"
    [ -n "$browser" ] && echo -e "🌐 Navegador: ${browser}"
    echo -e "⏱️  Timeout: ${TIMEOUT}ms"
    
    # Construir comando
    local cmd="npx playwright test ${test_file}"
    
    # Adicionar filtro de teste se especificado
    if [ -n "$test_name" ] && [ "$test_name" != "all" ]; then
        cmd="${cmd} --grep \"${test_name}\""
    fi
    
    # Adicionar modo
    if [ "$mode" = "headless" ]; then
        # Playwright usa --headed=false para headless
        cmd="${cmd}"
    elif [ "$mode" = "headed" ]; then
        cmd="${cmd} --headed"
    elif [ "$mode" = "ui" ]; then
        cmd="${cmd} --ui"
    fi
    
    # Adicionar browser se especificado
    if [ -n "$browser" ]; then
        cmd="${cmd} --project=${browser}"
    else
        # Por padrão, executar apenas Chrome e Firefox (evitar WebKit sem dependências)
        cmd="${cmd} --project=chromium --project=firefox"
    fi
    
    # Adicionar timeout
    cmd="${cmd} --timeout=${TIMEOUT}"
    
    # Adicionar reporter
    cmd="${cmd} --reporter=line"
    
    echo -e "${YELLOW}📝 Comando: ${cmd}${NC}"
    echo ""
    
    # Executar
    eval $cmd
    local result=$?
    
    return $result
}

# Gerar relatório
generate_report() {
    local result=$1
    local mode=$2
    local test_type=$3
    
    echo ""
    echo -e "${BLUE}📊 RELATÓRIO DE VALIDAÇÃO E2E${NC}"
    echo "================================="
    echo -e "📅 Data: $(date)"
    echo -e "🎯 Tipo: ${test_type}"
    echo -e "👁️  Modo: ${mode}"
    echo -e "🌐 URL: ${BASE_URL}"
    echo -e "👤 Entidade: ${TEST_ENTITY}"
    
    if [ $result -eq 0 ]; then
        echo -e "${GREEN}✅ RESULTADO: VALIDAÇÃO BEM-SUCEDIDA${NC}"
        echo ""
        echo -e "${GREEN}🎉 Todos os testes E2E passaram!${NC}"
        echo ""
        echo -e "📌 Próximos passos:"
        echo -e "   1. Documentar resultado no conclude_task.sh"
        echo -e "   2. Marcar tarefa como completa"
    elif [ $result -eq 2 ]; then
        echo -e "${YELLOW}⚠️  RESULTADO: VALIDAÇÃO COM AVISOS${NC}"
        echo ""
        echo -e "${YELLOW}Alguns testes passaram, mas há avisos.${NC}"
        echo -e "Verifique os logs acima para detalhes."
        echo ""
        echo -e "📌 Ação recomendada:"
        echo -e "   1. Revisar avisos no relatório"
        echo -e "   2. Corrigir se forem críticos"
        echo -e "   3. Documentar no conclude_task.sh"
    else
        echo -e "${RED}❌ RESULTADO: VALIDAÇÃO FALHOU${NC}"
        echo ""
        echo -e "${RED}Alguns testes E2E falharam.${NC}"
        echo -e "A tarefa NÃO deve ser marcada como completa."
        echo ""
        echo -e "📌 Ações necessárias:"
        echo -e "   1. Verificar erros no relatório acima"
        echo -e "   2. Corrigir os problemas"
        echo -e "   3. Executar validação novamente"
        echo -e "   4. Só concluir quando todos testes passarem"
    fi
    
    # Mostrar arquivos de resultado
    if [ -d "test-results" ]; then
        echo ""
        echo -e "${BLUE}📁 Arquivos de resultado:${NC}"
        find test-results -name "*.png" -o -name "*.webm" -o -name "*.txt" | head -5 | while read file; do
            echo "   - ${file}"
        done
        echo ""
        echo -e "💡 Para ver relatório completo: npx playwright show-report"
    fi
    
    return $result
}

# Validação básica (7 passos)
validate_basic() {
    echo -e "${BLUE}🧪 INICIANDO VALIDAÇÃO BÁSICA (7 PASSOS)${NC}"
    echo -e "${GREEN}Este é o fluxo padrão de validação Digna${NC}"
    echo ""
    
    # Executar teste básico
    run_playwright "tests/digna-basic.spec.js" "Fluxo simplificado - 7 passos" "$MODE" "$BROWSER"
    local result=$?
    
    if [ $result -eq 0 ]; then
        echo -e "${GREEN}✅ Fluxo básico validado com sucesso!${NC}"
        echo ""
        echo -e "${BLUE}📋 Resumo dos 7 passos:${NC}"
        echo "   1. ✅ Login no sistema"
        echo "   2. ✅ Dashboard carregado"  
        echo "   3. ✅ Página de estoque acessada"
        echo "   4. ✅ Página de membros acessada"
        echo "   5. ✅ Página de fornecedores acessada"
        echo "   6. ✅ PDV acessado"
        echo "   7. ✅ Caixa acessado"
    fi
    
    return $result
}

# Validação completa
validate_full() {
    echo -e "${BLUE}🧪 INICIANDO VALIDAÇÃO COMPLETA${NC}"
    echo -e "${YELLOW}Atenção: Esta validação pode levar vários minutos${NC}"
    echo ""
    
    # Executar todos os testes
    run_playwright "tests/" "all" "$MODE" "$BROWSER"
    return $?
}

# Validação customizada
validate_custom() {
    echo -e "${BLUE}🧪 INICIANDO VALIDAÇÃO CUSTOMIZADA${NC}"
    echo -e "Teste: ${CUSTOM_TEST}"
    echo ""
    
    # Determinar arquivo de teste
    local test_file="tests/"
    if [[ "$CUSTOM_TEST" == *"básic"* ]] || [[ "$CUSTOM_TEST" == *"simplific"* ]]; then
        test_file="tests/digna-basic.spec.js"
    elif [[ "$CUSTOM_TEST" == *"complet"* ]] || [[ "$CUSTOM_TEST" == *"fluxo complet"* ]]; then
        test_file="tests/digna-e2e.spec.js"
    fi
    
    run_playwright "$test_file" "$CUSTOM_TEST" "$MODE" "$BROWSER"
    return $?
}

# Main
main() {
    # Valores padrão
    MODE="headless"  # Modo stealth por padrão
    VALIDATION_TYPE="basic"
    BROWSER=""
    CUSTOM_TEST=""
    
    # Processar argumentos
    while [[ $# -gt 0 ]]; do
        case $1 in
            --basic|-b)
                VALIDATION_TYPE="basic"
                shift
                ;;
            --full|-f)
                VALIDATION_TYPE="full"
                shift
                ;;
            --custom|-c)
                VALIDATION_TYPE="custom"
                CUSTOM_TEST="$2"
                shift 2
                ;;
            --ui|-u)
                MODE="ui"
                shift
                ;;
            --headed)
                MODE="headed"
                shift
                ;;
            --headless)
                MODE="headless"
                shift
                ;;
            --headed)
                MODE="headed"
                shift
                ;;
            --chrome|-C)
                BROWSER="chromium"
                shift
                ;;
            --firefox|-F)
                BROWSER="firefox"
                shift
                ;;
            --timeout)
                TIMEOUT=$(($2 * 1000))  # Converter para ms
                shift 2
                ;;
            --help|-h)
                show_help
                ;;
            *)
                echo -e "${RED}❌ Argumento desconhecido: $1${NC}"
                show_help
                ;;
        esac
    done
    
    # Verificar custom test
    if [ "$VALIDATION_TYPE" = "custom" ] && [ -z "$CUSTOM_TEST" ]; then
        echo -e "${RED}❌ Modo --custom requer nome do teste${NC}"
        show_help
    fi
    
    # Banner
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}🧪 VALIDAÇÃO E2E - DIGNA${NC}"
    echo -e "${BLUE}========================================${NC}"
    echo ""
    
    # Verificar dependências
    if ! check_dependencies; then
        echo -e "${RED}❌ Dependências não atendidas. Abortando.${NC}"
        exit 1
    fi
    
    # Executar validação baseada no tipo
    local result=0
    case $VALIDATION_TYPE in
        "basic")
            validate_basic
            result=$?
            ;;
        "full")
            validate_full
            result=$?
            ;;
        "custom")
            validate_custom
            result=$?
            ;;
    esac
    
    # Gerar relatório
    generate_report $result "$MODE" "$VALIDATION_TYPE"
    
    # Retornar código de saída apropriado
    if [ $result -eq 0 ]; then
        echo -e "${GREEN}🎉 Validação E2E concluída com sucesso!${NC}"
        exit 0
    elif [ $result -eq 2 ]; then
        echo -e "${YELLOW}⚠️  Validação E2E concluída com avisos${NC}"
        exit 2
    else
        echo -e "${RED}❌ Validação E2E falhou${NC}"
        exit 1
    fi
}

# Executar main
main "$@"