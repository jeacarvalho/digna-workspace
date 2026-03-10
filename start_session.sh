#!/bin/bash

# 🚀 Script de conveniência para iniciar sessão
# Aponta para o script real em scripts/workflow/

echo "🔗 Executando start_session.sh de scripts/workflow/"
echo "=================================================="

# Verificar se o script existe
if [ -f "scripts/workflow/start_session.sh" ]; then
    exec scripts/workflow/start_session.sh "$@"
else
    echo "❌ Erro: scripts/workflow/start_session.sh não encontrado!"
    echo ""
    echo "📁 Estrutura do projeto reorganizada:"
    echo "  scripts/workflow/    - Scripts de workflow"
    echo "  scripts/dev/         - Scripts de desenvolvimento"
    echo "  scripts/tools/       - Ferramentas auxiliares"
    echo "  docs/               - Documentação"
    echo "  modules/            - Código fonte Go"
    echo ""
    echo "💡 Execute diretamente: scripts/workflow/start_session.sh"
    exit 1
fi