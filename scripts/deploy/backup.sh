#!/bin/bash

# ============================================================================
# Digna Database Backup Script
# ============================================================================
# Script para backup dos bancos de dados SQLite do Digna
# 
# Funcionalidades:
# 1. Cria backup timestamped dos bancos SQLite
# 2. Compacta em arquivo .tar.gz
# 3. Mantém histórico de backups
# 4. Pode ser agendado via cron
#
# Uso: ./backup.sh [--output-dir=/path/to/backups] [--keep-days=7]
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
DEFAULT_KEEP_DAYS=7
DATA_DIR="/var/lib/digna/data"

# Parse arguments
BACKUP_DIR="${DEFAULT_BACKUP_DIR}"
KEEP_DAYS="${DEFAULT_KEEP_DAYS}"

while [[ $# -gt 0 ]]; do
    case $1 in
        --output-dir=*)
            BACKUP_DIR="${1#*=}"
            shift
            ;;
        --keep-days=*)
            KEEP_DAYS="${1#*=}"
            shift
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            echo "Usage: $0 [--output-dir=/path/to/backups] [--keep-days=7]"
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

# ============================================================================
# Backup Functions
# ============================================================================

create_backup() {
    print_header "Criando Backup dos Bancos de Dados"
    
    # Create backup directory if it doesn't exist
    mkdir -p "${BACKUP_DIR}"
    
    # Generate timestamp
    TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
    BACKUP_FILE="${BACKUP_DIR}/digna_backup_${TIMESTAMP}.tar.gz"
    
    # Check if data directory exists
    if [[ ! -d "${DATA_DIR}" ]]; then
        print_error "Diretório de dados não encontrado: ${DATA_DIR}"
        exit 1
    fi
    
    # Count SQLite files
    SQLITE_COUNT=$(find "${DATA_DIR}" -name "*.db" -type f | wc -l)
    
    if [[ ${SQLITE_COUNT} -eq 0 ]]; then
        print_warning "Nenhum arquivo .db encontrado em ${DATA_DIR}"
        return 0
    fi
    
    print_warning "Encontrados ${SQLITE_COUNT} arquivos .db para backup"
    
    # Create backup
    print_warning "Criando backup: ${BACKUP_FILE}"
    
    # Backup SQLite files
    if tar -czf "${BACKUP_FILE}" -C "${DATA_DIR}" . 2>/dev/null; then
        BACKUP_SIZE=$(du -h "${BACKUP_FILE}" | cut -f1)
        print_success "Backup criado com sucesso: ${BACKUP_FILE} (${BACKUP_SIZE})"
        
        # Verify backup integrity
        if tar -tzf "${BACKUP_FILE}" > /dev/null 2>&1; then
            print_success "Integridade do backup verificada"
        else
            print_error "Falha na verificação de integridade do backup"
            exit 1
        fi
    else
        print_error "Falha ao criar backup"
        exit 1
    fi
}

clean_old_backups() {
    print_header "Limpando Backups Antigos"
    
    if [[ ! -d "${BACKUP_DIR}" ]]; then
        print_warning "Diretório de backup não existe: ${BACKUP_DIR}"
        return 0
    fi
    
    # Find and delete old backups
    OLD_BACKUPS=$(find "${BACKUP_DIR}" -name "digna_backup_*.tar.gz" -mtime +${KEEP_DAYS})
    
    if [[ -n "${OLD_BACKUPS}" ]]; then
        OLD_COUNT=$(echo "${OLD_BACKUPS}" | wc -l)
        print_warning "Removendo ${OLD_COUNT} backups antigos (mais de ${KEEP_DAYS} dias)"
        
        echo "${OLD_BACKUPS}" | while read -r backup; do
            print_warning "  Removendo: $(basename ${backup})"
            rm -f "${backup}"
        done
        
        print_success "Backups antigos removidos"
    else
        print_success "Nenhum backup antigo para remover"
    fi
}

show_backup_info() {
    print_header "Informações do Backup"
    
    if [[ -d "${BACKUP_DIR}" ]]; then
        BACKUP_COUNT=$(find "${BACKUP_DIR}" -name "digna_backup_*.tar.gz" -type f | wc -l)
        
        if [[ ${BACKUP_COUNT} -gt 0 ]]; then
            echo -e "${YELLOW}Backups disponíveis (${BACKUP_COUNT}):${NC}"
            find "${BACKUP_DIR}" -name "digna_backup_*.tar.gz" -type f -exec ls -lh {} \; | \
                awk '{print "  " $9 " (" $5 ")"}'
            
            # Show total backup size
            TOTAL_SIZE=$(find "${BACKUP_DIR}" -name "digna_backup_*.tar.gz" -type f -exec du -ch {} + | grep total$ | cut -f1)
            echo -e "\n${YELLOW}Tamanho total dos backups: ${TOTAL_SIZE}${NC}"
        else
            echo -e "${YELLOW}Nenhum backup encontrado${NC}"
        fi
    fi
    
    echo -e "\n${YELLOW}Configuração:${NC}"
    echo "  Diretório de dados: ${DATA_DIR}"
    echo "  Diretório de backup: ${BACKUP_DIR}"
    echo "  Manter backups por: ${KEEP_DAYS} dias"
    echo ""
    echo -e "${YELLOW}Agendamento via Cron (exemplo):${NC}"
    echo "  # Backup diário às 2:00 AM"
    echo "  0 2 * * * $(pwd)/backup.sh --output-dir=${BACKUP_DIR} --keep-days=${KEEP_DAYS}"
}

# ============================================================================
# Main Execution
# ============================================================================

main() {
    print_header "🚀 Backup dos Bancos de Dados Digna"
    echo "Data: $(date)"
    echo ""
    
    # Check if running as root (for accessing /var/lib/digna/data)
    if [[ "$EUID" -ne 0 ]]; then
        print_warning "Executando como usuário não-root"
        print_warning "Certifique-se de ter permissão para acessar ${DATA_DIR}"
    fi
    
    # Execute backup steps
    create_backup
    clean_old_backups
    show_backup_info
    
    print_header "✅ Backup Concluído"
    echo "Os bancos de dados SQLite foram backupados com sucesso."
    echo "Backups são mantidos por ${KEEP_DAYS} dias."
}

# Run main function
main "$@"