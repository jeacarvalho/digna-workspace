#!/bin/bash

# 🧪 Script de conveniência para executar testes
# Aponta para o script real em scripts/dev/

echo "🔗 Executando run_tests.sh de scripts/dev/"
echo "=========================================="

# Verificar se o script existe
if [ -f "scripts/dev/run_tests.sh" ]; then
    exec scripts/dev/run_tests.sh "$@"
else
    echo "❌ Erro: scripts/dev/run_tests.sh não encontrado!"
    echo ""
    echo "📁 Estrutura do projeto reorganizada:"
    echo "  scripts/workflow/    - Scripts de workflow"
    echo "  scripts/dev/         - Scripts de desenvolvimento"
    echo "  scripts/tools/       - Ferramentas auxiliares"
    echo ""
    echo "💡 Execute diretamente: scripts/dev/run_tests.sh"
    exit 1
fi