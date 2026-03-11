#!/bin/bash

# ============================================================================
# Deployment Script Validation
# ============================================================================
# Validates deployment scripts without requiring running server
# ============================================================================

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_header() {
    echo -e "\n${BLUE}========================================${NC}"
    echo -e "${BLUE}  $1${NC}"
    echo -e "${BLUE}========================================${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

validate_script() {
    local script="$1"
    local description="$2"
    
    echo -n "Validando ${description}... "
    
    if [[ ! -f "${script}" ]]; then
        print_error "Arquivo não encontrado: ${script}"
        return 1
    fi
    
    if [[ ! -x "${script}" ]]; then
        print_warning "Script não executável, corrigindo..."
        chmod +x "${script}"
    fi
    
    # Check for syntax errors
    if bash -n "${script}" 2>/dev/null; then
        print_success "OK"
        return 0
    else
        print_error "Erro de sintaxe"
        bash -n "${script}"
        return 1
    fi
}

validate_docker_files() {
    print_header "Validando Arquivos Docker"
    
    local errors=0
    
    # Check Dockerfile
    if [[ -f "Dockerfile" ]]; then
        echo -n "Validando Dockerfile... "
        if docker build --quiet . > /dev/null 2>&1; then
            print_success "OK"
        else
            print_error "Erro no Dockerfile"
            errors=$((errors + 1))
        fi
    else
        print_error "Dockerfile não encontrado"
        errors=$((errors + 1))
    fi
    
    # Check docker-compose.yml
    if [[ -f "docker-compose.yml" ]]; then
        echo -n "Validando docker-compose.yml... "
        if docker-compose config > /dev/null 2>&1; then
            print_success "OK"
        else
            print_error "Erro no docker-compose.yml"
            errors=$((errors + 1))
        fi
    else
        print_error "docker-compose.yml não encontrado"
        errors=$((errors + 1))
    fi
    
    # Check docker-compose.prod.yml
    if [[ -f "docker-compose.prod.yml" ]]; then
        echo -n "Validando docker-compose.prod.yml... "
        if docker-compose -f docker-compose.prod.yml config > /dev/null 2>&1; then
            print_success "OK"
        else
            print_error "Erro no docker-compose.prod.yml"
            errors=$((errors + 1))
        fi
    else
        print_warning "docker-compose.prod.yml não encontrado (opcional)"
    fi
    
    return ${errors}
}

validate_env_files() {
    print_header "Validando Arquivos de Ambiente"
    
    local errors=0
    
    # Check .env.example
    if [[ -f ".env.example" ]]; then
        echo -n "Validando .env.example... "
        
        # Check if it contains required variables
        local required_vars=("DIGNA_PORT" "DIGNA_DATA_DIR" "DIGNA_LOG_LEVEL")
        local missing_vars=()
        
        for var in "${required_vars[@]}"; do
            if ! grep -q "^${var}=" .env.example; then
                missing_vars+=("${var}")
            fi
        done
        
        if [[ ${#missing_vars[@]} -eq 0 ]]; then
            print_success "OK"
        else
            print_error "Variáveis faltando: ${missing_vars[*]}"
            errors=$((errors + 1))
        fi
    else
        print_error ".env.example não encontrado"
        errors=$((errors + 1))
    fi
    
    # Check config package
    if [[ -f "modules/ui_web/pkg/config/config.go" ]]; then
        echo -n "Validando config package... "
        
        # Check if config loads environment variables
        if grep -q "getEnv" "modules/ui_web/pkg/config/config.go" && \
           grep -q "DIGNA_PORT" "modules/ui_web/pkg/config/config.go"; then
            print_success "OK"
        else
            print_warning "Config package pode não estar usando variáveis de ambiente"
        fi
    else
        print_warning "Config package não encontrado (pode ser esperado)"
    fi
    
    return ${errors}
}

validate_deployment_scripts() {
    print_header "Validando Scripts de Deploy"
    
    local errors=0
    
    # Main deployment script
    validate_script "scripts/deploy/vps_deploy.sh" "vps_deploy.sh" || errors=$((errors + 1))
    
    # Backup script
    validate_script "scripts/deploy/backup.sh" "backup.sh" || errors=$((errors + 1))
    
    # Restore script  
    validate_script "scripts/deploy/restore.sh" "restore.sh" || errors=$((errors + 1))
    
    # Wrapper script
    validate_script "deploy.sh" "deploy.sh (wrapper)" || errors=$((errors + 1))
    
    return ${errors}
}

validate_documentation() {
    print_header "Validando Documentação"
    
    local errors=0
    
    local docs=("docs/DEPLOYMENT.md" "QUICK_DEPLOY.md")
    
    for doc in "${docs[@]}"; do
        if [[ -f "${doc}" ]]; then
            echo -n "Validando ${doc}... "
            
            # Check if file has content
            if [[ -s "${doc}" ]]; then
                print_success "OK"
            else
                print_error "Arquivo vazio"
                errors=$((errors + 1))
            fi
        else
            print_error "${doc} não encontrado"
            errors=$((errors + 1))
        fi
    done
    
    return ${errors}
}

main() {
    print_header "🚀 Validação de Scripts de Deploy"
    echo "Data: $(date)"
    echo ""
    
    local total_errors=0
    
    # Run validations
    validate_docker_files || total_errors=$((total_errors + $?))
    validate_env_files || total_errors=$((total_errors + $?))
    validate_deployment_scripts || total_errors=$((total_errors + $?))
    validate_documentation || total_errors=$((total_errors + $?))
    
    print_header "📊 Resultado da Validação"
    
    if [[ ${total_errors} -eq 0 ]]; then
        echo -e "${GREEN}✅ TODOS os testes passaram!${NC}"
        echo ""
        echo "Os scripts de deploy estão prontos para uso."
        echo "Execute './deploy.sh' para iniciar o deploy na VPS."
    else
        echo -e "${RED}❌ ${total_errors} teste(s) falharam${NC}"
        echo ""
        echo "Corrija os problemas acima antes de prosseguir."
        exit 1
    fi
    
    print_header "🎯 Próximos Passos"
    echo "1. Teste em ambiente staging: ./deploy.sh --env-file=.env.test"
    echo "2. Configure backup automático via cron"
    echo "3. Documente procedimentos específicos da sua organização"
    echo ""
    echo "Para deploy em produção:"
    echo "  ./deploy.sh --env-file=.env.production"
}

# Run main function
main "$@"