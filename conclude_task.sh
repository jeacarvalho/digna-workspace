#!/bin/bash

# 📚 Script de conveniência para concluir tarefas
# Aponta para o script real em scripts/workflow/

echo "🔗 Executando conclude_task.sh de scripts/workflow/"
echo "=================================================="

# Verificar se o script existe
if [ -f "scripts/workflow/conclude_task.sh" ]; then
    exec scripts/workflow/conclude_task.sh "$@"
else
    echo "❌ Erro: scripts/workflow/conclude_task.sh não encontrado!"
    echo ""
    echo "📁 Estrutura do projeto reorganizada:"
    echo "  scripts/workflow/    - Scripts de workflow"
    echo "  scripts/dev/         - Scripts de desenvolvimento"
    echo "  scripts/tools/       - Ferramentas auxiliares"
    echo ""
    echo "💡 Execute diretamente: scripts/workflow/conclude_task.sh"
    exit 1
fi