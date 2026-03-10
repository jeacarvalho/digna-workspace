#!/bin/bash

# 🧪 Script para executar testes E2E com Playwright
# Valida o fluxo completo de 7 passos da Digna

set -e

echo "🚀 Iniciando testes E2E da Digna"
echo "================================="

# Verificar se servidor está rodando
echo "🔍 Verificando servidor Digna..."
if ! curl -s http://localhost:8090/health > /dev/null; then
    echo "❌ Servidor não está rodando na porta 8090"
    echo "💡 Execute: cd modules/ui_web && go run main.go"
    exit 1
fi
echo "✅ Servidor rodando"

# Verificar Playwright instalado
echo "🔍 Verificando Playwright..."
if ! command -v npx &> /dev/null; then
    echo "❌ Node/npx não encontrado"
    echo "💡 Instale Node.js primeiro"
    exit 1
fi

if [ ! -f "package.json" ]; then
    echo "❌ package.json não encontrado"
    echo "💡 Execute: npm init -y && npm install --save-dev @playwright/test"
    exit 1
fi

# Executar testes
echo "🧪 Executando testes E2E..."
echo ""

# Opções de execução
MODE="${1:-basic}"
TIMEOUT="${2:-30000}"

case "$MODE" in
    "ui")
        echo "🎮 Modo UI (interface gráfica)"
        npm run test:e2e:ui
        ;;
    "debug")
        echo "🐛 Modo Debug"
        npm run test:e2e:debug
        ;;
    "headless")
        echo "👻 Modo Headless"
        npx playwright test --timeout=$TIMEOUT
        ;;
    "chrome")
        echo "🌐 Apenas Chrome"
        npx playwright test tests/digna-basic.spec.js --project=chromium --timeout=$TIMEOUT
        ;;
    "full")
        echo "📋 Testes completos"
        npx playwright test --timeout=$TIMEOUT
        ;;
    "basic"|*)
        echo "🧪 Testes básicos de validação"
        npx playwright test tests/digna-basic.spec.js --project=chromium --headed --timeout=$TIMEOUT
        ;;
esac

# Verificar resultado
if [ $? -eq 0 ]; then
    echo ""
    echo "🎉 Testes E2E concluídos com sucesso!"
    echo ""
    echo "📊 Fluxo validado:"
    echo "   1. ✅ Login no sistema"
    echo "   2. ✅ Item de estoque criado/verificado"
    echo "   3. ✅ Membro criado/verificado"
    echo "   4. ✅ Fornecedor criado/verificado"
    echo "   5. ✅ Compra registrada"
    echo "   6. ✅ Venda registrada no PDV"
    echo "   7. ✅ Saldo e horas trabalhadas validados"
    echo ""
    echo "🔗 Para ver relatório: npm run report"
else
    echo ""
    echo "❌ Testes E2E falharam"
    echo ""
    echo "🔍 Solução de problemas:"
    echo "   - Verifique se servidor está rodando: curl http://localhost:8090/health"
    echo "   - Verifique credenciais: cafe_digna / cd0123"
    echo "   - Execute em modo debug: ./scripts/dev/run_e2e_tests.sh debug"
    echo "   - Verifique screenshots em: test-results/"
    exit 1
fi