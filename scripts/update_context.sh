#!/bin/bash

# Script para atualizar automaticamente o contexto do projeto Digna
# Executar após implementações significativas

echo "🔄 Atualizando contexto do projeto Digna..."

# Diretório base
BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
QUICK_REF="$BASE_DIR/docs/QUICK_REFERENCE.md"

# Atualizar data da última atualização
CURRENT_DATE=$(date +"%d/%m/%Y")
sed -i "s|**Última atualização:**.*|**Última atualização:** $CURRENT_DATE|" "$QUICK_REF"

# Contar testes passando (opcional - pode ser caro)
# TEST_COUNT=$(cd "$BASE_DIR/modules" && go test ./... 2>&1 | grep -c "PASS")
# if [ $TEST_COUNT -gt 0 ]; then
#     sed -i "s|✅ PRODUCTION READY (.* testes passando)|✅ PRODUCTION READY ($TEST_COUNT testes passando)|" "$QUICK_REF"
# fi

# Verificar handlers existentes
HANDLERS=$(find "$BASE_DIR/modules/ui_web/internal/handler" -name "*.go" -not -name "*_test.go" -exec basename {} \; | sed 's/\.go//' | sort)
HANDLER_LIST=$(echo "$HANDLERS" | tr '\n' ',' | sed 's/,$//;s/,/, /g')

# Atualizar seção de handlers se necessário
if ! grep -q "Handlers existentes:" "$QUICK_REF"; then
    echo -e "\n## 🏗️ Handlers Existentes\n\n$HANDLER_LIST" >> "$QUICK_REF"
fi

echo "✅ Contexto atualizado em: $QUICK_REF"
echo "📅 Última atualização: $CURRENT_DATE"
echo "📋 Handlers: $HANDLER_LIST"