#!/bin/bash
# quick_agent_check.sh - Validação rápida para o agente no início da sessão
# Uso: ./scripts/tools/quick_agent_check.sh [módulo]

set -e

echo "🔍 VALIDAÇÃO RÁPIDA PARA AGENTE"
echo "================================"

MODULE="${1:-all}"

case $MODULE in
    "legal")
        echo "📋 Validando módulo legal_facade..."
        echo ""
        
        # 1. Verificar se módulo existe
        if [ -d "modules/legal_facade" ]; then
            echo "✅ Módulo legal_facade existe"
            
            # 2. Verificar arquivos críticos
            CRITICAL_FILES=(
                "internal/document/generator.go"
                "internal/document/formalization.go"
                "internal/document/statute.go"
                "internal/document/legal_repository.go"
            )
            
            for file in "${CRITICAL_FILES[@]}"; do
                if [ -f "modules/legal_facade/$file" ]; then
                    echo "  ✅ $file"
                else
                    echo "  ⚠️  $file (NÃO ENCONTRADO)"
                fi
            done
            
            # 3. Verificar constante de formalização
            echo ""
            echo "📊 Constante de formalização:"
            grep -n "MinDecisionsForFormalization" modules/legal_facade/internal/document/formalization.go || echo "  ⚠️  Não encontrada"
            
            # 4. Verificar SHA256 implementations
            echo ""
            echo "🔐 Implementações SHA256:"
            grep -l "sha256.Sum256" modules/legal_facade/internal/document/*.go 2>/dev/null || echo "  ⚠️  Nenhuma encontrada"
            
        else
            echo "❌ Módulo legal_facade NÃO existe"
        fi
        ;;
        
    "core")
        echo "📋 Validando módulo core_lume..."
        echo ""
        
        if [ -d "modules/core_lume" ]; then
            echo "✅ Módulo core_lume existe"
            
            # Verificar DecisionRepository
            echo ""
            echo "🗃️ DecisionRepository:"
            if grep -q "DecisionRepository" modules/core_lume/internal/repository/interfaces.go; then
                echo "  ✅ Interface encontrada"
                grep -n "type DecisionRepository interface" modules/core_lume/internal/repository/interfaces.go
            else
                echo "  ⚠️  Interface NÃO encontrada"
            fi
            
            # Verificar SHA256
            echo ""
            echo "🔐 SHA256 em decision_service:"
            if [ -f "modules/core_lume/internal/service/decision_service.go" ]; then
                grep -n "generateHash\|sha256" modules/core_lume/internal/service/decision_service.go | head -5
            else
                echo "  ⚠️  decision_service.go não encontrado"
            fi
            
        else
            echo "❌ Módulo core_lume NÃO existe"
        fi
        ;;
        
    "ui")
        echo "📋 Validando módulo ui_web..."
        echo ""
        
        if [ -d "modules/ui_web" ]; then
            echo "✅ Módulo ui_web existe"
            
            # Verificar handlers
            echo ""
            echo "🎨 Handlers disponíveis:"
            ls modules/ui_web/internal/handler/*.go 2>/dev/null | xargs -n1 basename | head -10
            
            # Verificar file download pattern
            echo ""
            echo "📥 File download pattern (accountant_handler):"
            if [ -f "modules/ui_web/internal/handler/accountant_handler.go" ]; then
                grep -n "Content-Disposition" modules/ui_web/internal/handler/accountant_handler.go | head -2
            else
                echo "  ⚠️  accountant_handler.go não encontrado"
            fi
            
            # Verificar templates
            echo ""
            echo "📝 Templates disponíveis:"
            ls modules/ui_web/templates/*_simple.html 2>/dev/null | xargs -n1 basename | head -5
            
        else
            echo "❌ Módulo ui_web NÃO existe"
        fi
        ;;
        
    "all")
        echo "📋 Validando todos os módulos..."
        echo ""
        
        # Executar todas as validações
        $0 legal
        echo ""
        echo "---"
        echo ""
        $0 core
        echo ""
        echo "---"
        echo ""
        $0 ui
        
        # Verificar skills
        echo ""
        echo "📚 Skills disponíveis:"
        if [ -d "docs/skills" ]; then
            ls docs/skills/
        else
            echo "  ⚠️  Diretório docs/skills não encontrado"
        fi
        
        # Verificar documentação de aprendizados
        echo ""
        echo "🎓 Aprendizados recentes:"
        if [ -f "docs/learnings/SESSION_INSIGHTS_20260311.md" ]; then
            echo "  ✅ SESSION_INSIGHTS_20260311.md (última sessão)"
        else
            echo "  ℹ️  Nenhum aprendizado recente encontrado"
        fi
        ;;
        
    *)
        echo "❌ Módulo desconhecido: $MODULE"
        echo "Uso: $0 [legal|core|ui|all]"
        exit 1
        ;;
esac

echo ""
echo "✅ Validação concluída"
echo ""
echo "💡 Dica: Consulte docs/MODULES_QUICK_REFERENCE.md para referência completa"