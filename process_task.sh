#!/bin/bash

# 📋 Script de conveniência para processar tarefas
# Aponta para o script real em scripts/workflow/

echo "🔗 Executando process_task.sh de scripts/workflow/"
echo "=================================================="

# Verificar se o script existe
if [ -f "scripts/workflow/process_task.sh" ]; then
    exec scripts/workflow/process_task.sh "$@"
else
    echo "❌ Erro: scripts/workflow/process_task.sh não encontrado!"
    echo ""
    echo "📁 Estrutura do projeto reorganizada:"
    echo "  scripts/workflow/    - Scripts de workflow"
    echo "  scripts/dev/         - Scripts de desenvolvimento"
    echo "  scripts/tools/       - Ferramentas auxiliares"
    echo ""
    echo "💡 Execute diretamente: scripts/workflow/process_task.sh"
    exit 1
fi