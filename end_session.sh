#!/bin/bash

# 🚪 Script de conveniência para encerrar sessão
# Aponta para o script real em scripts/workflow/

echo "🔗 Executando end_session.sh de scripts/workflow/"
echo "================================================="

# Verificar se o script existe
if [ -f "scripts/workflow/end_session.sh" ]; then
    exec scripts/workflow/end_session.sh "$@"
else
    echo "❌ Erro: scripts/workflow/end_session.sh não encontrado!"
    echo ""
    echo "📁 Estrutura do projeto reorganizada:"
    echo "  scripts/workflow/    - Scripts de workflow"
    echo "  scripts/dev/         - Scripts de desenvolvimento"
    echo "  scripts/tools/       - Ferramentas auxiliares"
    echo ""
    echo "💡 Execute diretamente: scripts/workflow/end_session.sh"
    exit 1
fi