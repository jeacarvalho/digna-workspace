#!/bin/bash

# 🚀 Script de conveniência para preparar implementação
# Aponta para o script real em scripts/workflow/

echo "🔗 Executando prepare_implementation.sh de scripts/workflow/"
echo "=========================================================="

# Verificar se o script existe
if [ -f "scripts/workflow/prepare_implementation.sh" ]; then
    exec scripts/workflow/prepare_implementation.sh "$@"
else
    echo "❌ Erro: scripts/workflow/prepare_implementation.sh não encontrado!"
    echo ""
    echo "📁 Estrutura do projeto:"
    echo "  scripts/workflow/    - Scripts de workflow"
    echo "  scripts/dev/         - Scripts de desenvolvimento"
    echo "  docs/               - Documentação"
    echo ""
    echo "💡 Execute diretamente: scripts/workflow/prepare_implementation.sh"
    exit 1
fi