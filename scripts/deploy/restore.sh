#!/bin/bash

# ============================================================================
# Digna Database Restore Script
# ============================================================================
# Script para restaurar bancos de dados SQLite do Digna a partir de backup
# 
# Funcionalidades:
# 1. Lista backups disponíveis
# 2. Restaura bancos de dados a partir de backup selecionado
# 3. Valida integridade antes da restauração
# 4. Cria backup atual antes de restaurar (safety net)
#
# Uso: ./restore.sh [--backup-file=/path/to/backup.tar.gz] [--dry-run]
# ============================================================================

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default configuration
DEFAULT_BACKUP_DIR="/var/backups/digna"
DATA_DIR="/var/lib/digna/data"

# Parse arguments
BACKUP_FILE=""
DRY_RUN=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --backup-file=*)
            BACKUP_FILE="${1#*=}"
            shift
            ;;
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            echo "Usage: $0 [--backup-file=/path/to/backup.tar.gz] [--dry-run]"
            exit 1
            ;;
    esac
done

# ============================================================================
# Helper Functions
# ============================================================================

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

confirm_action() {
    local message="$1"
    local default="${2:-n}"
    
    if [[ "${DRY_RUN}" == "true" ]]; then
        echo -e "${YELLOW}[DRY RUN] ${message}${NC}"
        return 0
    fi
    
    read -p "$message (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        return 0
    else
        return 1
    fi
}

# ============================================================================
# Restore Functions
# ============================================================================

list_backups() {
    print_header "Backups Disponíveis"
    
    if [[ ! -d "${DEFAULT_BACKUP_DIR}" ]]; then
        print_error "Diretório de backup não encontrado: ${DEFAULT_BACKUP_DIR}"
        return 1
    fi
    
    local backups=($(find "${DEFAULT_BACKUP_DIR}" -name "digna_backup_*.tar.gz" -type f | sort -r))
    
    if [[ ${#backups[@]} -eq 0 ]]; then
        print_error "Nenhum backup encontrado em ${DEFAULT_BACKUP_DIR}"
        return 1
    fi
    
    echo -e "${YELLOW}Backups disponíveis:${NC}"
    for i in "${!backups[@]}"; do
        local backup="${backups[$i]}"
        local filename=$(basename "${backup}")
        local size=$(du -h "${backup}" | cut -f1)
        local date=$(echo "${filename}" | grep -oE '[0-9]{8}_[0-9]{6}' | sed 's/_/ /' | sed 's/\([0-9]\{4\}\)\([0-9]\{2\}\)\([0-9]\{2\}\)/\1-\2-\3/')
        
        echo "  $((i+1)). ${filename}"
        echo "      Data: ${date}"
        echo "      Tamanho: ${size}"
        echo ""
    done
    
    return 0
}

select_backup() {
    if [[ -n "${BACKUP_FILE}" ]]; then
        if [[ ! -f "${BACKUP_FILE}" ]]; then
            print_error "Arquivo de backup não encontrado: ${BACKUP_FILE}"
            return 1
        fi
        echo "${BACKUP_FILE}"
        return 0
    fi
    
    list_backups || return 1
    
    read -p "Selecione o número do backup para restaurar (ou 0 para cancelar): " selection
    
    if [[ ! "${selection}" =~ ^[0-9]+$ ]]; then
        print_error "Seleção inválida"
        return 1
    fi
    
    if [[ "${selection}" -eq 0 ]]; then
        print_warning "Restauração cancelada"
        return 1
    fi
    
    local backups=($(find "${DEFAULT_BACKUP_DIR}" -name "digna_backup_*.tar.gz" -type f | sort -r))
    
    if [[ "${selection}" -gt ${#backups[@]} ]]; then
        print_error "Número de backup inválido"
        return 1
    fi
    
    echo "${backups[$((selection-1))]}"
}

validate_backup() {
    local backup_file="$1"
    
    print_header "Validando Backup"
    
    # Check if file exists
    if [[ ! -f "${backup_file}" ]]; then
        print_error "Arquivo de backup não encontrado: ${backup_file}"
        return 1
    fi
    
    # Check file integrity
    print_warning "Verificando integridade do backup..."
    if ! tar -tzf "${backup_file}" > /dev/null 2>&1; then
        print_error "Backup corrompido ou inválido"
        return 1
    fi
    
    # List contents
    print_warning "Conteúdo do backup:"
    tar -tzf "${backup_file}" | head -20 | while read -r file; do
        echo "  ${file}"
    done
    
    local file_count=$(tar -tzf "${backup_file}" | wc -l)
    print_success "Backup válido contendo ${file_count} arquivos"
    
    return 0
}

create_pre_restore_backup() {
    print_header "Criando Backup de Segurança"
    
    if [[ "${DRY_RUN}" == "true" ]]; then
        print_warning "[DRY RUN] Criaria backup atual antes da restauração"
        return 0
    fi
    
    # Create emergency backup of current data
    local timestamp=$(date +"%Y%m%d_%H%M%S")
    local emergency_backup="${DEFAULT_BACKUP_DIR}/emergency_pre_restore_${timestamp}.tar.gz"
    
    if [[ -d "${DATA_DIR}" ]] && [[ -n "$(ls -A ${DATA_DIR} 2>/dev/null)" ]]; then
        print_warning "Criando backup de emergência dos dados atuais..."
        
        if tar -czf "${emergency_backup}" -C "${DATA_DIR}" . 2>/dev/null; then
            local size=$(du -h "${emergency_backup}" | cut -f1)
            print_success "Backup de emergência criado: $(basename ${emergency_backup}) (${size})"
            echo "${emergency_backup}"
        else
            print_error "Falha ao criar backup de emergência"
            return 1
        fi
    else
        print_warning "Diretório de dados vazio ou não existe, pulando backup de emergência"
    fi
    
    return 0
}

restore_backup() {
    local backup_file="$1"
    
    print_header "Restaurando Backup"
    
    if [[ "${DRY_RUN}" == "true" ]]; then
        print_warning "[DRY RUN] Restauraria backup: ${backup_file}"
        print_warning "[DRY RUN] Para diretório: ${DATA_DIR}"
        return 0
    fi
    
    # Stop Digna service if running
    print_warning "Parando serviço Digna..."
    if docker-compose ps | grep -q "digna-app"; then
        docker-compose stop
        print_success "Serviço parado"
    else
        print_warning "Serviço Digna não está rodando"
    fi
    
    # Create data directory if it doesn't exist
    mkdir -p "${DATA_DIR}"
    
    # Clear existing data (with confirmation)
    if confirm_action "⚠️  ATENÇÃO: Isso irá APAGAR todos os dados atuais em ${DATA_DIR}. Continuar?"; then
        print_warning "Limpando dados atuais..."
        rm -rf "${DATA_DIR}"/*
        print_success "Dados atuais removidos"
    else
        print_error "Restauração cancelada pelo usuário"
        return 1
    fi
    
    # Restore from backup
    print_warning "Restaurando dados do backup..."
    if tar -xzf "${backup_file}" -C "${DATA_DIR}"; then
        local restored_count=$(find "${DATA_DIR}" -name "*.db" -type f | wc -l)
        print_success "Backup restaurado com sucesso"
        print_success "${restored_count} bancos de dados SQLite restaurados"
    else
        print_error "Falha ao restaurar backup"
        return 1
    fi
    
    # Fix permissions
    print_warning "Ajustando permissões..."
    chmod -R 755 "${DATA_DIR}"
    print_success "Permissões ajustadas"
    
    return 0
}

start_service() {
    print_header "Iniciando Serviço"
    
    if [[ "${DRY_RUN}" == "true" ]]; then
        print_warning "[DRY RUN] Iniciaria serviço Digna"
        return 0
    fi
    
    print_warning "Iniciando serviço Digna..."
    if docker-compose up -d; then
        print_success "Serviço iniciado"
        
        # Wait for service to be healthy
        print_warning "Aguardando serviço ficar saudável..."
        for i in {1..30}; do
            if curl -s -f "http://localhost:${DIGNA_PORT:-8090}/health" > /dev/null; then
                print_success "Serviço está saudável!"
                break
            fi
            echo -n "."
            sleep 2
        done
    else
        print_error "Falha ao iniciar serviço"
        return 1
    fi
    
    return 0
}

# ============================================================================
# Main Execution
# ============================================================================

main() {
    print_header "🚀 Restauração de Bancos de Dados Digna"
    echo "Data: $(date)"
    echo ""
    
    # Check if running in project directory
    if [[ ! -f "docker-compose.yml" ]]; then
        print_error "Execute este script do diretório raiz do projeto Digna"
        exit 1
    fi
    
    # Select backup file
    local selected_backup=$(select_backup)
    if [[ -z "${selected_backup}" ]]; then
        exit 1
    fi
    
    # Validate backup
    validate_backup "${selected_backup}" || exit 1
    
    # Show warning
    print_header "⚠️  AVISO IMPORTANTE"
    echo "Esta operação irá:"
    echo "1. PARAR o serviço Digna"
    echo "2. APAGAR todos os dados atuais em ${DATA_DIR}"
    echo "3. RESTAURAR dados do backup: $(basename ${selected_backup})"
    echo ""
    echo "Um backup de emergência será criado antes da restauração."
    echo ""
    
    if ! confirm_action "Deseja continuar com a restauração?"; then
        print_warning "Restauração cancelada"
        exit 0
    fi
    
    # Create emergency backup
    local emergency_backup=$(create_pre_restore_backup)
    
    # Restore backup
    restore_backup "${selected_backup}" || exit 1
    
    # Start service
    start_service || exit 1
    
    print_header "✅ Restauração Concluída"
    echo "Backup restaurado com sucesso: $(basename ${selected_backup})"
    if [[ -n "${emergency_backup}" ]]; then
        echo "Backup de emergência criado: $(basename ${emergency_backup})"
    fi
    echo ""
    echo "Serviço Digna está rodando e saudável."
    echo "Acesse http://localhost:${DIGNA_PORT:-8090} para verificar."
}

# Run main function
main "$@"