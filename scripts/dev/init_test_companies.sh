#!/bin/bash

# Script para inicializar bancos de dados das empresas de teste

echo "🔧 Inicializando empresas de teste do sistema Digna"
echo "=================================================="

cd /home/s015533607/Documentos/desenv/digna-workspace

# Criar diretório de dados se não existir
mkdir -p data/entities

echo ""
echo "1. Verificando/Criando bancos de dados..."
echo "-----------------------------------------"

# Empresas de teste
declare -A companies=(
    ["cafe_digna"]="Café Digna"
    ["queijaria_digna"]="Queijaria Digna"
)

for entity_id in "${!companies[@]}"; do
    company_name="${companies[$entity_id]}"
    db_file="data/entities/${entity_id}.db"
    
    if [ -f "$db_file" ]; then
        echo "✅ $company_name: Banco já existe ($db_file)"
    else
        echo "📦 $company_name: Criando banco de dados..."
        
        # Criar arquivo de banco de dados vazio
        touch "$db_file"
        
        # Configurar pragmas básicos do SQLite
        sqlite3 "$db_file" << EOF
PRAGMA journal_mode=WAL;
PRAGMA foreign_keys=ON;
PRAGMA synchronous=NORMAL;
PRAGMA temp_store=MEMORY;
PRAGMA mmap_size=268435456;
EOF
        
        if [ $? -eq 0 ]; then
            echo "   ✅ Banco criado com sucesso"
        else
            echo "   ❌ Erro ao criar banco"
        fi
    fi
done

echo ""
echo "2. Verificando estrutura..."
echo "---------------------------"

# Verificar se os bancos têm tabelas
for entity_id in "${!companies[@]}"; do
    db_file="data/entities/${entity_id}.db"
    
    if [ -f "$db_file" ]; then
        table_count=$(sqlite3 "$db_file" "SELECT COUNT(*) FROM sqlite_master WHERE type='table';" 2>/dev/null || echo "0")
        echo "   $entity_id: $table_count tabelas encontradas"
    fi
done

echo ""
echo "3. Credenciais de acesso..."
echo "---------------------------"
echo "   Café Digna:"
echo "     • Usuário: cafe_digna"
echo "     • Senha: cd0123"
echo ""
echo "   Queijaria Digna:"
echo "     • Usuário: queijaria_digna"
echo "     • Senha: qd321"
echo ""
echo "4. Iniciando servidor..."
echo "------------------------"
echo "   Para iniciar o servidor:"
echo "   cd modules/ui_web && go run main.go"
echo ""
echo "   Ou usar o Makefile:"
echo "   make run"
echo ""
echo "   Acesse: http://localhost:8088/login"
echo ""
echo "✅ Inicialização concluída!"