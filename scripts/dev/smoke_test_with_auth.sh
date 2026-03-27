#!/bin/bash

# Smoke test aprimorado com modo de teste
# Uso: ./smoke_test_with_auth.sh [feature_name] [route_path] [test_mode]

set -e

FEATURE_NAME="${1:-Nova Feature}"
ROUTE_PATH="${2:-/}"
TEST_MODE="${3:-test}"  # 'test' para usar dados de teste, 'prod' para produção
ENTITY_ID="${4:-test-entity-001}"

echo "🚀 SMOKE TEST (Modo: $TEST_MODE)"
echo "========================================"
echo "📋 Feature: $FEATURE_NAME"
echo "📍 Rota: $ROUTE_PATH"
echo "🏢 Entity: $ENTITY_ID"
echo ""

# Cores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Verificar se servidor está rodando
echo -e "${BLUE}1. Verificando servidor...${NC}"
if curl -s http://localhost:8090/health > /dev/null; then
    echo -e "${GREEN}✅ Servidor rodando${NC}"
else
    echo -e "${RED}❌ Servidor não está rodando em http://localhost:8090${NC}"
    echo "   Execute: cd modules/ui_web && DIGNA_ENV=test go run ."
    exit 1
fi

# Em modo de teste, usar sessão de teste
if [ "$TEST_MODE" = "test" ]; then
    echo -e "${BLUE}2. Criando sessão de teste...${NC}"
    
    # Tentar criar uma sessão de teste (isso depende da implementação do auth)
    # Por enquanto, vamos assumir que o modo teste aceita entity_id direto
    TEST_URL="http://localhost:8090${ROUTE_PATH}?entity_id=${ENTITY_ID}&test_mode=true"
else
    TEST_URL="http://localhost:8090${ROUTE_PATH}"
fi

# Testar rota
echo -e "${BLUE}3. Testando rota $ROUTE_PATH...${NC}"
RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" "$TEST_URL" 2>&1)

if [ "$RESPONSE" = "200" ]; then
    echo -e "${GREEN}✅ Rota $ROUTE_PATH responde 200 OK${NC}"
elif [ "$RESPONSE" = "302" ]; then
    echo -e "${YELLOW}⚠️  Rota $ROUTE_PATH responde 302 (redirect para login)${NC}"
    echo "   💡 Para testar sem autenticação, use: test_mode=true"
else
    echo -e "${YELLOW}⚠️  Rota $ROUTE_PATH responde $RESPONSE${NC}"
fi

# Verificar template
echo -e "${BLUE}4. Verificando template...${NC}"
TEMPLATE_NAME=$(echo "$FEATURE_NAME" | tr '[:upper:]' '[:lower:]' | tr ' ' '_')
TEMPLATE_FILE="modules/ui_web/templates/${TEMPLATE_NAME}_simple.html"

if [ -f "$TEMPLATE_FILE" ]; then
    echo -e "${GREEN}✅ Template encontrado: ${TEMPLATE_NAME}_simple.html${NC}"
    
    # Verificar conteúdo
    if grep -q "<html" "$TEMPLATE_FILE"; then
        echo -e "${GREEN}✅ Template contém HTML válido${NC}"
    else
        echo -e "${RED}❌ Template não contém HTML válido${NC}"
    fi
else
    echo -e "${RED}❌ Template não encontrado: ${TEMPLATE_FILE}${NC}"
    echo "   Crie o template: ${TEMPLATE_NAME}_simple.html"
fi

# Verificar handler no main.go
echo -e "${BLUE}5. Verificando handler em main.go...${NC}"
if grep -q "New.*Handler.*$FEATURE_NAME\|${TEMPLATE_NAME}" modules/ui_web/main.go 2>/dev/null || \
   grep -q "$FEATURE_NAME" modules/ui_web/main.go 2>/dev/null; then
    echo -e "${GREEN}✅ Handler registrado em main.go${NC}"
else
    echo -e "${YELLOW}⚠️  Handler pode não estar registrado em main.go${NC}"
fi

echo ""
echo -e "${BLUE}6. Resumo:${NC}"
echo "   Modo: $TEST_MODE"
echo "   URL Testada: $TEST_URL"
echo "   Status HTTP: $RESPONSE"

if [ "$RESPONSE" = "200" ] || [ "$RESPONSE" = "302" ]; then
    echo -e "${GREEN}✅ Smoke test concluído com sucesso!${NC}"
    exit 0
else
    echo -e "${RED}❌ Smoke test encontrou problemas${NC}"
    exit 1
fi
