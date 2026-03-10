#!/bin/bash
# analyze_patterns.sh - Analisa padrões de código do projeto Digna
# Uso: ./analyze_patterns.sh [handler_name] [--html] [--tests] [--all]

set -e

echo "🔍 ANALISANDO PADRÕES DO PROJETO DIGNA"
echo "======================================"

# Configurações
HANDLER_NAME="${1:-member}"
ANALYZE_HTML=false
ANALYZE_TESTS=false
ANALYZE_ALL=false

# Processar argumentos
for arg in "$@"; do
    case $arg in
        --html)
            ANALYZE_HTML=true
            ;;
        --tests)
            ANALYZE_TESTS=true
            ;;
        --all)
            ANALYZE_ALL=true
            ANALYZE_HTML=true
            ANALYZE_TESTS=true
            ;;
    esac
done

# 1. Análise do BaseHandler
echo ""
echo "1. 🏗️  BASE HANDLER (modules/ui_web/internal/handler/base_handler.go)"
echo "-------------------------------------------------------------------"

if [ -f "modules/ui_web/internal/handler/base_handler.go" ]; then
    # Funções de template disponíveis
    echo "📋 Funções de template registradas:"
    grep -n "AddFunc" modules/ui_web/internal/handler/base_handler.go | sed 's/^/   /'
    
    # Estrutura do BaseHandler
    echo ""
    echo "📁 Estrutura do BaseHandler:"
    grep -n "type BaseHandler\|func NewBaseHandler" modules/ui_web/internal/handler/base_handler.go | sed 's/^/   /'
else
    echo "❌ BaseHandler não encontrado!"
fi

# 2. Análise do handler específico
echo ""
echo "2. 🎯 HANDLER ANALISADO: ${HANDLER_NAME}"
echo "--------------------------------------"

HANDLER_FILE="modules/ui_web/internal/handler/${HANDLER_NAME}_handler.go"
if [ -f "$HANDLER_FILE" ]; then
    # Estrutura do handler
    echo "📋 Estrutura:"
    grep -n "type.*Handler\|func New.*Handler" "$HANDLER_FILE" | head -5 | sed 's/^/   /'
    
    # Funções de template específicas
    echo ""
    echo "🎨 Funções de template específicas:"
    grep -n "AddFunc" "$HANDLER_FILE" | sed 's/^/   /'
    
    # Rotas registradas
    echo ""
    echo "🛣️  Rotas registradas:"
    grep -n "RegisterRoutes\|HandleFunc\|http\." "$HANDLER_FILE" | head -10 | sed 's/^/   /'
    
    # Métodos HTTP
    echo ""
    echo "🌐 Métodos HTTP implementados:"
    grep -n "func.*http" "$HANDLER_FILE" | sed 's/^/   /'
else
    echo "❌ Handler ${HANDLER_NAME} não encontrado!"
    echo "   Handlers disponíveis:"
    ls modules/ui_web/internal/handler/*_handler.go 2>/dev/null | xargs -n1 basename | sed 's/_handler.go//' | sed 's/^/   - /'
fi

# 3. Análise do template HTML
if [ "$ANALYZE_HTML" = true ] || [ "$ANALYZE_ALL" = true ]; then
    echo ""
    echo "3. 🎨 TEMPLATE HTML (modules/ui_web/templates/${HANDLER_NAME}_simple.html)"
    echo "-----------------------------------------------------------------------"
    
    # Tentar diferentes padrões de nome
    TEMPLATE_FILE="modules/ui_web/templates/${HANDLER_NAME}_simple.html"
    if [ ! -f "$TEMPLATE_FILE" ]; then
        # Tentar plural
        TEMPLATE_FILE="modules/ui_web/templates/${HANDLER_NAME}s_simple.html"
    fi
    if [ ! -f "$TEMPLATE_FILE" ]; then
        # Tentar nome exato
        TEMPLATE_FILE="modules/ui_web/templates/${HANDLER_NAME}.html"
    fi
    if [ -f "$TEMPLATE_FILE" ]; then
        # Estrutura básica
        echo "📏 Tamanho: $(wc -l < "$TEMPLATE_FILE") linhas"
        
        # Elementos HTMX
        echo ""
        echo "⚡ Elementos HTMX encontrados:"
        grep -c "hx-" "$TEMPLATE_FILE" | sed 's/^/   Total: /'
        grep -n "hx-post\|hx-get\|hx-target\|hx-swap" "$TEMPLATE_FILE" | head -10 | sed 's/^/   /'
        
        # Classes Tailwind (padrão Digna)
        echo ""
        echo "🎨 Classes Tailwind (paleta Digna):"
        grep -n "bg-digna-\|text-digna-\|border-digna-\|digna-" "$TEMPLATE_FILE" | head -10 | sed 's/^/   /'
        
        # Formulários
        echo ""
        echo "📝 Formulários encontrados:"
        grep -c "<form" "$TEMPLATE_FILE" | sed 's/^/   Total: /'
        grep -n "<form" "$TEMPLATE_FILE" | sed 's/^/   /'
    else
        echo "❌ Template ${HANDLER_NAME}_simple.html não encontrado!"
        echo "   Templates disponíveis:"
        ls modules/ui_web/templates/*_simple.html 2>/dev/null | xargs -n1 basename | sed 's/^/   - /'
    fi
fi

# 4. Análise de testes
if [ "$ANALYZE_TESTS" = true ] || [ "$ANALYZE_ALL" = true ]; then
    echo ""
    echo "4. 🧪 TESTES (modules/ui_web/internal/handler/${HANDLER_NAME}_handler_test.go)"
    echo "---------------------------------------------------------------------------"
    
    TEST_FILE="modules/ui_web/internal/handler/${HANDLER_NAME}_handler_test.go"
    if [ -f "$TEST_FILE" ]; then
        # Testes disponíveis
        echo "📋 Testes implementados:"
        grep -n "func Test" "$TEST_FILE" | sed 's/^/   /'
        
        # Cobertura aproximada
        echo ""
        echo "📊 Métricas:"
        TEST_COUNT=$(grep -c "func Test" "$TEST_FILE")
        echo "   Total de testes: $TEST_COUNT"
        
        # Setup/teardown
        echo ""
        echo "⚙️  Setup/Teardown:"
        grep -n "TestMain\|setup\|teardown" "$TEST_FILE" | sed 's/^/   /'
    else
        echo "❌ Testes para ${HANDLER_NAME} não encontrados!"
        echo "   Testes disponíveis:"
        ls modules/ui_web/internal/handler/*_test.go 2>/dev/null | xargs -n1 basename | sed 's/^/   - /'
    fi
fi

# 5. Padrões comuns extraídos
echo ""
echo "5. 📋 PADRÕES COMUNS IDENTIFICADOS"
echo "---------------------------------"

echo "🎯 Estrutura de Handler:"
echo "   1. type {Feature}Handler struct { *BaseHandler, lifecycleManager }"
echo "   2. func New{Feature}Handler(lm) (*{Feature}Handler, error)"
echo "   3. func (h *{Feature}Handler) RegisterRoutes(mux)"
echo "   4. Métodos: List{Feature}(), Create{Feature}(), Update{Feature}(), etc."

echo ""
echo "🎨 Padrões de Template:"
echo "   1. <!DOCTYPE html> com Tailwind config"
echo "   2. Header com navegação (copiar de dashboard_simple.html)"
echo "   3. Main content com HTMX forms"
echo "   4. Scripts HTMX no final"
echo "   5. Classes: digna-card, digna-gradient, bg-digna-*, text-digna-*"

echo ""
echo "⚡ Padrões HTMX:"
echo "   - GET: Renderização inicial"
echo "   - POST: Ações de formulário"
echo "   - hx-target + hx-swap para feedback"
echo "   - URL patterns: /{feature}, /{feature}/{id}/action"

echo ""
echo "🧪 Padrões de Testes:"
echo "   - Test{Feature}Handler_{Action}()"
echo "   - httptest.NewRecorder()"
echo "   - Mock lifecycle.LifecycleManager"
echo "   - Testes de integração com banco real"

# 6. Comandos úteis para análise mais profunda
echo ""
echo "6. 🔧 COMANDOS ÚTEIS PARA ANÁLISE PROFUNDA"
echo "-----------------------------------------"

cat << 'EOF'
# Analisar todos os handlers
for h in $(ls modules/ui_web/internal/handler/*_handler.go | xargs -n1 basename | sed 's/_handler.go//'); do
  echo "=== $h ==="
  grep -n "type.*Handler\|func New" "modules/ui_web/internal/handler/${h}_handler.go" | head -2
done

# Analisar uso de HTMX em todos templates
grep -r "hx-" modules/ui_web/templates/ | sort | uniq -c | sort -rn

# Analisar funções de template
grep -r "AddFunc" modules/ui_web/internal/handler/ | sort

# Analisar padrões de rotas
grep -r "HandleFunc\|http\." modules/ui_web/internal/handler/ | grep -v "test" | sort
EOF

echo ""
echo "✅ ANÁLISE CONCLUÍDA"
echo "==================="
echo "Use estas informações para entender padrões antes de implementar."
echo "Sempre consulte handlers/templates similares como referência."

exit 0