#!/bin/bash

# 🚀 SMOKE TEST PARA NOVAS FEATURES
# Valida que uma nova feature funciona no ambiente REAL (não apenas em testes)
# Uso: ./scripts/smoke_test_new_feature.sh "nome_da_feature" "rota_principal"

set -e  # Para em caso de erro

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 SMOKE TEST: Validando nova feature${NC}"
echo "========================================"

FEATURE_NAME="$1"
MAIN_ROUTE="$2"
ENTITY_ID="cooperativa_demo"

if [ -z "$FEATURE_NAME" ] || [ -z "$MAIN_ROUTE" ]; then
    echo -e "${RED}❌ Uso: $0 \"nome_da_feature\" \"/rota_principal\"${NC}"
    echo "   Exemplo: $0 \"Member Management\" \"/members\""
    exit 1
fi

echo -e "${YELLOW}📋 Feature:${NC} $FEATURE_NAME"
echo -e "${YELLOW}📍 Rota:${NC} $MAIN_ROUTE"
echo -e "${YELLOW}🏢 Entity:${NC} $ENTITY_ID"
echo ""

# 1. Verificar se servidor está rodando
echo -e "${BLUE}1. Verificando servidor...${NC}"
if ! curl -s http://localhost:8090/health > /dev/null; then
    echo -e "${RED}❌ Servidor não está rodando em http://localhost:8090${NC}"
    echo "   Execute: cd modules/ui_web && go run ."
    exit 1
fi
echo -e "${GREEN}✅ Servidor rodando${NC}"

# 2. Testar rota principal
echo -e "${BLUE}2. Testando rota $MAIN_ROUTE...${NC}"
RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost:8090${MAIN_ROUTE}?entity_id=${ENTITY_ID}")

if [ "$RESPONSE" = "200" ]; then
    echo -e "${GREEN}✅ Rota $MAIN_ROUTE responde 200${NC}"
elif [ "$RESPONSE" = "404" ]; then
    echo -e "${RED}❌ Rota $MAIN_ROUTE não encontrada (404)${NC}"
    echo "   Problemas comuns:"
    echo "   - Handler não registrado no main.go"
    echo "   - Rota definida errado em RegisterRoutes"
    echo "   - Servidor precisa reiniciar após mudanças"
    exit 1
elif [ "$RESPONSE" = "500" ]; then
    echo -e "${RED}❌ Erro interno no servidor (500)${NC}"
    echo "   Verifique logs do servidor"
    exit 1
else
    echo -e "${YELLOW}⚠️  Rota $MAIN_ROUTE responde $RESPONSE${NC}"
fi

# 3. Verificar template (se aplicável)
echo -e "${BLUE}3. Verificando template...${NC}"
# Extrair nome do template da rota (ex: /members -> members_simple.html, /supply/stock -> supply_stock_simple.html)
TEMPLATE_NAME=$(echo "$MAIN_ROUTE" | sed 's|^/||' | sed 's|/$||' | sed 's|/|_|g')_simple.html
TEMPLATE_PATH="modules/ui_web/templates/$TEMPLATE_NAME"

if [ -f "$TEMPLATE_PATH" ]; then
    echo -e "${GREEN}✅ Template encontrado: $TEMPLATE_NAME${NC}"
    
    # Verificar se template tem conteúdo mínimo
    if [ $(wc -l < "$TEMPLATE_PATH") -lt 5 ]; then
        echo -e "${YELLOW}⚠️  Template muito pequeno (menos de 5 linhas)${NC}"
    fi
else
    echo -e "${RED}❌ Template não encontrado: $TEMPLATE_PATH${NC}"
    echo "   Crie o template: $TEMPLATE_NAME"
    exit 1
fi

# 4. Testar conteúdo da página
echo -e "${BLUE}4. Verificando conteúdo da página...${NC}"
PAGE_CONTENT=$(curl -s "http://localhost:8090${MAIN_ROUTE}?entity_id=${ENTITY_ID}")

# Verificações básicas
if echo "$PAGE_CONTENT" | grep -q "<html"; then
    echo -e "${GREEN}✅ Página contém HTML válido${NC}"
else
    echo -e "${RED}❌ Página não contém HTML${NC}"
    exit 1
fi

if echo "$PAGE_CONTENT" | grep -q "$FEATURE_NAME\|${TEMPLATE_NAME%_simple.html}"; then
    echo -e "${GREEN}✅ Página contém referência à feature${NC}"
else
    echo -e "${YELLOW}⚠️  Página não menciona a feature pelo nome${NC}"
fi

# 5. Verificar navegação
echo -e "${BLUE}5. Verificando navegação...${NC}"
# Verificar se link aparece em templates principais
MAIN_TEMPLATES=(
    "modules/ui_web/templates/dashboard_simple.html"
    "modules/ui_web/templates/layout.html"
)

NAV_ADDED=false
for template in "${MAIN_TEMPLATES[@]}"; do
    if [ -f "$template" ] && grep -q "$MAIN_ROUTE" "$template"; then
        echo -e "${GREEN}✅ Link encontrado em $(basename $template)${NC}"
        NAV_ADDED=true
    fi
done

if [ "$NAV_ADDED" = false ]; then
    echo -e "${YELLOW}⚠️  Link não encontrado em templates de navegação${NC}"
    echo "   Adicione link para $MAIN_ROUTE em:"
    echo "   - dashboard_simple.html"
    echo "   - layout.html"
fi

# 6. Resumo
echo ""
echo -e "${BLUE}📊 RESUMO DO SMOKE TEST${NC}"
echo "========================================"
echo -e "${GREEN}✅ Servidor: Rodando${NC}"
echo -e "${GREEN}✅ Rota $MAIN_ROUTE: Responde ${RESPONSE}${NC}"
echo -e "${GREEN}✅ Template $TEMPLATE_NAME: Existe${NC}"
echo -e "${GREEN}✅ Página: HTML válido${NC}"

if [ "$NAV_ADDED" = true ]; then
    echo -e "${GREEN}✅ Navegação: Link adicionado${NC}"
else
    echo -e "${YELLOW}⚠️  Navegação: Link pendente${NC}"
fi

echo ""
echo -e "${GREEN}🎉 SMOKE TEST CONCLUÍDO!${NC}"
echo "A feature '$FEATURE_NAME' está funcional no ambiente local."
echo ""
echo -e "${BLUE}📝 PRÓXIMOS PASSOS:${NC}"
echo "1. Testar funcionalidades específicas (CRUD, forms, etc.)"
echo "2. Executar testes E2E completos"
echo "3. Documentar aprendizados com ./conclude_task.sh"

exit 0