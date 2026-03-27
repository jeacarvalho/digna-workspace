#!/bin/bash

# Setup de ambiente de teste para Digna
# Este script prepara um ambiente isolado para testes automatizados

set -e

echo "🧪 Setup do Ambiente de Teste"
echo "=============================="

# Diretório base
BASE_DIR="/home/s015533607/Documentos/desenv/digna-workspace"
DATA_DIR="$BASE_DIR/data"
TEST_DIR="$DATA_DIR/test"

echo ""
echo "📁 Criando estrutura de diretórios..."
mkdir -p "$TEST_DIR/entities"
mkdir -p "$TEST_DIR/central"

echo ""
echo "🗄️  Preparando banco de dados de teste..."

# Criar banco central de teste (SQLite em memória ou arquivo)
CENTRAL_DB="$TEST_DIR/central/central.db"

if [ ! -f "$CENTRAL_DB" ]; then
    echo "   Criando central.db..."
    sqlite3 "$CENTRAL_DB" <<EOF
-- Tabela de usuários
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    role TEXT NOT NULL,
    status TEXT DEFAULT 'ACTIVE',
    created_at INTEGER
);

-- Tabela de entidades
CREATE TABLE IF NOT EXISTS entities (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    cnpj TEXT UNIQUE,
    status TEXT DEFAULT 'ACTIVE',
    created_at INTEGER
);

-- Tabela de vínculos usuário-entidade
CREATE TABLE IF NOT EXISTS user_entities (
    user_id TEXT NOT NULL,
    entity_id TEXT NOT NULL,
    role TEXT NOT NULL,
    created_at INTEGER,
    PRIMARY KEY (user_id, entity_id)
);

-- Tabela de membros
CREATE TABLE IF NOT EXISTS members (
    id TEXT PRIMARY KEY,
    entity_id TEXT NOT NULL,
    name TEXT NOT NULL,
    email TEXT,
    role TEXT NOT NULL,
    status TEXT DEFAULT 'ACTIVE',
    created_at INTEGER
);

-- Inserir dados de teste
INSERT OR IGNORE INTO users (id, email, name, role, status, created_at) VALUES 
('test-user-001', 'test@digna.local', 'Usuário de Teste', 'COORDINATOR', 'ACTIVE', strftime('%s', 'now'));

INSERT OR IGNORE INTO entities (id, name, cnpj, status, created_at) VALUES 
('test-entity-001', 'Cooperativa Teste', '00000000000191', 'ACTIVE', strftime('%s', 'now')),
('test-entity-002', 'MEI Teste', '00000000000272', 'ACTIVE', strftime('%s', 'now'));

INSERT OR IGNORE INTO user_entities (user_id, entity_id, role, created_at) VALUES 
('test-user-001', 'test-entity-001', 'COORDINATOR', strftime('%s', 'now')),
('test-user-001', 'test-entity-002', 'COORDINATOR', strftime('%s', 'now'));

INSERT OR IGNORE INTO members (id, entity_id, name, email, role, status, created_at) VALUES 
('test-member-001', 'test-entity-001', 'Coordenador Teste', 'coord@test.local', 'COORDINATOR', 'ACTIVE', strftime('%s', 'now')),
('test-member-002', 'test-entity-002', 'MEI Teste', 'mei@test.local', 'MEMBER', 'ACTIVE', strftime('%s', 'now'));
EOF
    echo "   ✅ Banco central criado"
else
    echo "   ℹ️  Banco central já existe"
fi

# Criar banco da entidade de teste
ENTITY_DB="$TEST_DIR/entities/test-entity-001.db"

if [ ! -f "$ENTITY_DB" ]; then
    echo ""
    echo "   Criando banco da entidade test-entity-001..."
    sqlite3 "$ENTITY_DB" <<EOF
-- Tabela de membros da entidade
CREATE TABLE IF NOT EXISTS members (
    id TEXT PRIMARY KEY,
    entity_id TEXT NOT NULL,
    name TEXT NOT NULL,
    email TEXT,
    role TEXT NOT NULL,
    status TEXT DEFAULT 'ACTIVE',
    created_at INTEGER
);

INSERT OR IGNORE INTO members (id, entity_id, name, email, role, status, created_at) VALUES 
('test-member-001', 'test-entity-001', 'Coordenador Teste', 'coord@test.local', 'COORDINATOR', 'ACTIVE', strftime('%s', 'now'));
EOF
    echo "   ✅ Banco da entidade criado"
else
    echo "   ℹ️  Banco da entidade já existe"
fi

echo ""
echo "📝 Criando arquivo de configuração de teste..."

cat > "$BASE_DIR/.env.test" <<EOF
# Configurações de Teste
DIGNA_ENV=test
DIGNA_DATA_DIR=$TEST_DIR
DIGNA_CENTRAL_DB=$CENTRAL_DB
PORT=8090

# Usuário de teste padrão
TEST_USER_ID=test-user-001
TEST_USER_EMAIL=test@digna.local
TEST_ENTITY_ID=test-entity-001
EOF

echo ""
echo "✅ Ambiente de teste configurado!"
echo ""
echo "📋 Resumo:"
echo "   Diretório de dados: $TEST_DIR"
echo "   Banco central: $CENTRAL_DB"
echo "   Entidade de teste: test-entity-001"
echo "   Usuário de teste: test-user-001"
echo ""
echo "🚀 Para iniciar em modo de teste:"
echo "   cd modules/ui_web && DIGNA_ENV=test go run ."
echo ""
echo "🧪 Para rodar smoke tests:"
echo "   ./scripts/dev/smoke_test_with_auth.sh 'Minha Feature' '/rota' test"
