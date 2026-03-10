#!/bin/bash
# start_session.sh - Inicializa uma sessão no projeto Digna
# Uso: ./start_session.sh [opcional: "quick" para modo rápido]

set -e  # Sai no primeiro erro

echo "🚀 Iniciando sessão no projeto Digna..."
echo "=========================================="

# Configurações
SESSION_ID=$(date +%Y%m%d_%H%M%S)
SESSION_FILE=".session_${SESSION_ID}"
QUICK_MODE="${1:-no}"

# 1. Criar arquivo de sessão
echo "SESSION_ID=${SESSION_ID}" > ${SESSION_FILE}
echo "START_TIME=$(date +%s)" >> ${SESSION_FILE}
echo "QUICK_MODE=${QUICK_MODE}" >> ${SESSION_FILE}

# 2. Atualizar contexto (modo rápido ou completo)
if [ "$QUICK_MODE" = "quick" ]; then
    echo "⚡ Modo rápido: verificando status básico..."
    
    # Verificação mínima
    if [ -f "./scripts/update_context.sh" ]; then
        ./scripts/update_context.sh 2>/dev/null || true
    fi
    
    # Status dos testes
    echo "📊 Verificando status dos testes..."
    if [ -d "modules" ]; then
        cd modules
        TEST_RESULT=$(./run_tests.sh 2>&1 | tail -20)
        echo "$TEST_RESULT" | grep -E "(PASS|FAIL|ok|^---)" | head -10
        cd ..
    fi
    
else
    echo "🔍 Modo completo: atualizando contexto completo..."
    
    # Atualizar contexto se o script existir
    if [ -f "./scripts/update_context.sh" ]; then
        echo "🔄 Executando update_context.sh..."
        ./scripts/update_context.sh
    else
        echo "⚠️  Script update_context.sh não encontrado, criando contexto básico..."
        
        # Criar QUICK_REFERENCE se não existir
        if [ ! -f "docs/QUICK_REFERENCE.md" ]; then
            mkdir -p docs
            cat > docs/QUICK_REFERENCE.md << 'EOF'
# 🚀 Quick Reference - Projeto Digna

**Última atualização:** $(date +%d/%m/%Y)
**Sessão iniciada:** ${SESSION_ID}

## 📋 Status Rápido
- Projeto: Digna (Economia Solidária)
- Arquitetura: Clean Architecture + DDD
- Banco: SQLite isolado por entidade
- Frontend: HTMX + Tailwind

## 🛠️ Comandos Úteis
- Iniciar sessão: `./start_session.sh`
- Processar tarefa: `./process_task.sh "descrição da tarefa"`
- Concluir tarefa: `./conclude_task.sh "aprendizados"`

## 📁 Estrutura Principal
```
modules/
├── core_lume/     # Domínio e serviços
├── ui_web/        # Interface web (HTMX/Tailwind)
├── lifecycle/     # Gerenciamento de banco
└── [outros...]
```
EOF
        fi
    fi
    
    # Status detalhado
    echo "📊 Status detalhado do projeto..."
    
    # Contar testes
    TEST_COUNT=$(find . -name "*test*.go" -type f | wc -l)
    echo "📈 Testes encontrados: ${TEST_COUNT}"
    
    # Verificar handlers
    HANDLER_COUNT=$(find modules/ui_web/internal/handler -name "*.go" -type f 2>/dev/null | wc -l || echo "0")
    echo "🎨 Handlers UI: ${HANDLER_COUNT}"
    
    # Verificar templates
    TEMPLATE_COUNT=$(find modules/ui_web/templates -name "*.html" -type f 2>/dev/null | wc -l || echo "0")
    echo "📝 Templates: ${TEMPLATE_COUNT}"
fi

# 3. Mostrar referência rápida
echo ""
echo "📚 REFERÊNCIA RÁPIDA:"
echo "===================="

if [ -f "docs/QUICK_REFERENCE.md" ]; then
    head -30 docs/QUICK_REFERENCE.md
else
    echo "ℹ️  Use './process_task.sh' para processar sua primeira tarefa."
    echo "📋 Formato recomendado:"
    echo "   ./process_task.sh \"Tipo: Feature | Módulo: ui_web | Objetivo: Implementar X | Decisões: seguir padrão Y\""
fi

# 4. Validação de Integridade do Projeto (NOVO)
echo ""
echo "🔍 VALIDAÇÃO DE INTEGRIDADE DO PROJETO:"
echo "======================================"

# Contar handlers registrados
HANDLER_COUNT=$(grep -c "RegisterRoutes" modules/ui_web/main.go 2>/dev/null || echo "0")
echo "Handlers registrados no main.go: $HANDLER_COUNT"

# Listar handlers
echo "Handlers encontrados:"
grep "New.*Handler" modules/ui_web/main.go | grep -v "//" | while read line; do
  handler=$(echo "$line" | sed 's/.*New\([A-Za-z]*\)Handler.*/\1Handler/' | head -1)
  if [ -n "$handler" ] && [ "$handler" != "line" ]; then
    echo "  ✅ $handler"
  fi
done 2>/dev/null

# Verificar templates vs handlers
echo ""
echo "Compatibilidade Templates-Handlers:"
TEMPLATE_MENTIONS=$(grep -r "ExecuteTemplate.*\.html" modules/ui_web/internal/handler/ 2>/dev/null | grep -v "test" | wc -l)
echo "  Templates referenciados em handlers: $TEMPLATE_MENTIONS"

# Verificar se smoke test script existe
if [ -f "./scripts/dev/smoke_test_new_feature.sh" ]; then
    echo "  ✅ Smoke test script disponível (scripts/dev/)"
else
    echo "  ⚠️  Smoke test script NÃO encontrado em scripts/dev/"
fi

# 5. Mostrar próximos passos sugeridos
echo ""
echo "🎯 PRÓXIMOS PASSOS SUGERIDOS:"
echo "============================="

if [ -f "docs/NEXT_STEPS.md" ]; then
    head -20 docs/NEXT_STEPS.md
else
    echo "1. Escolha uma tarefa do backlog"
    echo "2. Use './process_task.sh \"sua tarefa\"'"
    echo "3. Siga o checklist pré-implementação"
    echo "4. ✅ SEMPRE execute smoke test após implementação"
    echo "5. Documente aprendizados com './conclude_task.sh'"
fi

# 5. Informações da sessão
echo ""
echo "📋 INFORMAÇÕES DA SESSÃO:"
echo "========================="
echo "ID: ${SESSION_ID}"
echo "Data: $(date)"
echo "Modo: ${QUICK_MODE}"
echo "Arquivo de sessão: ${SESSION_FILE}"
echo ""
echo "✅ Sessão iniciada com sucesso!"
echo ""
echo "💡 Dica: Agora use './process_task.sh \"sua descrição de tarefa\"'"
echo "       ou './process_task.sh --help' para ver opções"

# Tornar scripts executáveis se não forem
chmod +x process_task.sh 2>/dev/null || true
chmod +x conclude_task.sh 2>/dev/null || true

exit 0