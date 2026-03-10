#!/bin/bash

# 🚀 PREPARE IMPLEMENTATION FROM ROADMAP/PROMPT FILE
# Recebe um arquivo de prompt/roadmap e prepara implementação completa

set -e

# Cores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

show_help() {
    echo -e "${BLUE}🚀 PREPARE IMPLEMENTATION FROM FILE${NC}"
    echo "=========================================="
    echo ""
    echo "Uso: ./prepare_implementation.sh <arquivo_prompt> [opções]"
    echo ""
    echo "Opções:"
    echo "  --checklist    Apenas gera checklist pré-implementação"
    echo "  --plan         Gera checklist + plano de implementação"
    echo "  --execute      Gera tudo + prompt final para opencode (padrão)"
    echo "  --help         Mostra esta ajuda"
    echo ""
    echo "Exemplos:"
    echo "  ./prepare_implementation.sh roadmap/suppliers.md"
    echo "  ./prepare_implementation.sh prompts/member_management.txt --plan"
    echo ""
    echo "📋 Formato esperado do arquivo:"
    echo "  # Título da Implementação"
    echo "  "
    echo "  **Tipo:** Feature | Bug Fix | Refactor"
    echo "  **Módulo:** ui_web | core_lume | supply"
    echo "  **Objetivo:** Descrição clara do que implementar"
    echo "  **Decisões:** Decisões técnicas/arquiteturais"
    echo "  "
    echo "  (opcionalmente mais detalhes)"
    echo ""
}

# Verificar argumentos
if [ $# -eq 0 ] || [ "$1" = "--help" ]; then
    show_help
    exit 0
fi

INPUT_FILE="$1"
MODE="execute"

# Processar opções
shift
while [ $# -gt 0 ]; do
    case "$1" in
        --checklist)
            MODE="checklist"
            ;;
        --plan)
            MODE="plan"
            ;;
        --execute)
            MODE="execute"
            ;;
        *)
            echo -e "${RED}❌ Opção desconhecida: $1${NC}"
            show_help
            exit 1
            ;;
    esac
    shift
done

# Verificar se arquivo existe
if [ ! -f "$INPUT_FILE" ]; then
    echo -e "${RED}❌ Arquivo não encontrado: $INPUT_FILE${NC}"
    echo ""
    echo "💡 Dica: Crie um arquivo com o formato:"
    echo "  **Tipo:** Feature"
    echo "  **Módulo:** ui_web"
    echo "  **Objetivo:** Implementar X"
    echo "  **Decisões:** seguir padrão Y"
    exit 1
fi

echo -e "${BLUE}📁 Processando arquivo: $INPUT_FILE${NC}"
echo "=========================================="

# Extrair informações do arquivo
echo -e "${YELLOW}🔍 Extraindo informações do arquivo...${NC}"

# Ler conteúdo
FILE_CONTENT=$(cat "$INPUT_FILE")

# Extrair campos usando regex
TASK_TYPE=$(echo "$FILE_CONTENT" | grep -i "^\s*\*\*Tipo:\*\*" | sed 's/.*\*\*Tipo:\*\*\s*//i' | head -1)
MODULE=$(echo "$FILE_CONTENT" | grep -i "^\s*\*\*Módulo:\*\*" | sed 's/.*\*\*Módulo:\*\*\s*//i' | head -1)
OBJECTIVE=$(echo "$FILE_CONTENT" | grep -i "^\s*\*\*Objetivo:\*\*" | sed 's/.*\*\*Objetivo:\*\*\s*//i' | head -1)
DECISIONS=$(echo "$FILE_CONTENT" | grep -i "^\s*\*\*Decisões:\*\*" | sed 's/.*\*\*Decisões:\*\*\s*//i' | head -1)

# Se não encontrar no formato **Campo:**, tentar formato simples
if [ -z "$TASK_TYPE" ]; then
    TASK_TYPE=$(echo "$FILE_CONTENT" | grep -i "^Tipo:" | sed 's/^Tipo:\s*//i' | head -1)
fi
if [ -z "$MODULE" ]; then
    MODULE=$(echo "$FILE_CONTENT" | grep -i "^Módulo:" | sed 's/^Módulo:\s*//i' | head -1)
fi
if [ -z "$OBJECTIVE" ]; then
    OBJECTIVE=$(echo "$FILE_CONTENT" | grep -i "^Objetivo:" | sed 's/^Objetivo:\s*//i' | head -1)
fi
if [ -z "$DECISIONS" ]; then
    DECISIONS=$(echo "$FILE_CONTENT" | grep -i "^Decisões:" | sed 's/^Decisões:\s*//i' | head -1)
fi

# Se ainda não encontrar, usar primeiras linhas como objetivo
if [ -z "$OBJECTIVE" ]; then
    # Pegar primeira linha não vazia como objetivo
    OBJECTIVE=$(echo "$FILE_CONTENT" | grep -v "^$" | head -1 | sed 's/^#\s*//')
fi

# Validações
if [ -z "$OBJECTIVE" ]; then
    echo -e "${RED}❌ Não foi possível extrair 'Objetivo' do arquivo${NC}"
    echo "💡 Adicione uma linha com '**Objetivo:**' ou 'Objetivo:'"
    exit 1
fi

# Valores padrão
TASK_TYPE=${TASK_TYPE:-"Feature"}
MODULE=${MODULE:-"ui_web"}
DECISIONS=${DECISIONS:-"seguir padrões estabelecidos no projeto"}

# Criar nome da feature (simplificado do objetivo)
FEATURE_NAME=$(echo "$OBJECTIVE" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9]/-/g' | sed 's/-\+/-/g' | sed 's/^-\|-$//g' | cut -c1-30)
if [ -z "$FEATURE_NAME" ]; then
    FEATURE_NAME="implementation-$(date +%s)"
fi

echo -e "${GREEN}✅ Informações extraídas:${NC}"
echo "  Tipo: $TASK_TYPE"
echo "  Módulo: $MODULE"
echo "  Objetivo: $OBJECTIVE"
echo "  Decisões: $DECISIONS"
echo "  Feature: $FEATURE_NAME"
echo ""

# Criar descrição formatada para process_task.sh
TASK_DESCRIPTION="Tipo: $TASK_TYPE | Módulo: $MODULE | Objetivo: $OBJECTIVE | Decisões: $DECISIONS"

echo -e "${YELLOW}🚀 Preparando implementação via process_task.sh...${NC}"

# Executar process_task.sh com as informações extraídas
case "$MODE" in
    "checklist")
        echo "📋 Gerando apenas checklist..."
        ./process_task.sh "$TASK_DESCRIPTION" --checklist
        ;;
    "plan")
        echo "📋📝 Gerando checklist + plano..."
        ./process_task.sh "$TASK_DESCRIPTION" --plan
        ;;
    "execute")
        echo "🚀 Gerando implementação completa..."
        ./process_task.sh "$TASK_DESCRIPTION" --execute
        
        # Encontrar arquivo de prompt gerado
        PROMPT_FILE=$(ls -t .opencode_task_*.txt 2>/dev/null | head -1)
        
        if [ -n "$PROMPT_FILE" ] && [ -f "$PROMPT_FILE" ]; then
            echo ""
            echo -e "${GREEN}🎯 PROMPT FINAL GERADO!${NC}"
            echo "=================================="
            echo ""
            echo "📋 O que fazer agora:"
            echo "1. Copie o conteúdo de: $PROMPT_FILE"
            echo "2. Cole no opencode"
            echo "3. Siga as instruções do prompt"
            echo ""
            echo "💡 Comandos rápidos:"
            echo "   cat $PROMPT_FILE | pbcopy       # Mac"
            echo "   cat $PROMPT_FILE | xclip -sel c # Linux"
            echo "   cat $PROMPT_FILE                # Ver conteúdo"
            echo ""
            
            # Mostrar preview do prompt
            echo -e "${YELLOW}📄 Preview do prompt (primeiras 10 linhas):${NC}"
            head -10 "$PROMPT_FILE"
            echo "..."
        else
            echo -e "${YELLOW}⚠️  Prompt não gerado automaticamente.${NC}"
            echo "💡 Execute manualmente: ./process_task.sh \"$TASK_DESCRIPTION\" --execute"
        fi
        ;;
esac

echo ""
echo -e "${GREEN}✅ PREPARAÇÃO CONCLUÍDA!${NC}"
echo "=============================="

# Mostrar próximos passos
echo ""
echo -e "${BLUE}📋 PRÓXIMOS PASSOS:${NC}"
echo "-------------------"
case "$MODE" in
    "checklist")
        echo "1. Preencha o checklist: docs/implementation_plans/${FEATURE_NAME}_pre_check.md"
        echo "2. Execute: ./prepare_implementation.sh \"$INPUT_FILE\" --plan"
        ;;
    "plan")
        echo "1. Revise o plano: docs/implementation_plans/${FEATURE_NAME}_implementation_*.md"
        echo "2. Execute: ./prepare_implementation.sh \"$INPUT_FILE\" --execute"
        ;;
    "execute")
        echo "1. Copie o prompt (.opencode_task_*.txt) para o opencode"
        echo "2. Implemente seguindo as instruções"
        echo "3. Após implementar: ./scripts/dev/smoke_test_new_feature.sh"
        echo "4. Conclua: ./conclude_task.sh \"Aprendizados da implementação\""
        ;;
esac

exit 0