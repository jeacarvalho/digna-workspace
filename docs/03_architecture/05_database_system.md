title: Sistema de Database - Arquitetura SQLite Isolada
status: implemented
version: 2.0
last_updated: 2026-03-27
---

# Sistema de Database - Arquitetura SQLite Isolada

> **Nota:** Este documento reflete a arquitetura de database integrada do Ecossistema Digna (PDF v1.0), preservando toda a infraestrutura validada nas Sprints 1-16, incorporando os novos módulos como Fases 3 e 4, e incluindo o Sistema de Ajuda Educativa (RF-30) decidido na sessão de 27/03/2026.

---

## 📋 Visão Geral

O sistema de database do Digna é baseado no princípio de **"Soberania do Dado"** - cada entidade possui seu próprio banco SQLite isolado fisicamente. Esta arquitetura garante que o dado pertença ao produtor, não à plataforma, e permite **Exit Power** total (o usuário pode levar seu banco e migrar para outro sistema).

Com a expansão para o **Ecossistema de 4 Módulos** (PDF v1.0) e a adição do **Sistema de Ajuda Educativa** (RF-30, Sessão 27/03/2026), a arquitetura de database foi expandida para suportar:
- Dados globais compartilhados (indicadores, programas, help topics)
- Dados específicos por entidade (perfil, transações, matches)
- Isolamento rigoroso entre tenants

---

## 🎯 Objetivo da Atualização (27/03/2026)

**Contexto:** O projeto evoluiu de um ERP contábil para um **Ecossistema de 4 Módulos** conforme especificação PDF v1.0. Esta atualização documenta:

1. **Expansão da Dualidade Arquitetural** (Banco Central vs. Banco por Entidade)
2. **Novas Tabelas para Módulos 2, 3 e 4** (Indicadores, Portal, Rede)
3. **Tabela `help_topics`** para RF-30 (Sistema de Ajuda Educativa)
4. **Preservação do Schema Existente** (Sprints 1-16 validadas)

**Versão Anterior:** 1.0 (2026-03-09)  
**Nova Versão:** 2.0 (2026-03-27)

---

## 🏗️ Arquitetura de Database

### Estrutura por Entidade

```
data/entities/
├── central.db               # ✅ Banco Central (dados globais)
├── cafe_digna.db            # ✅ ENTITY (populado com dados reais)
├── cooperativa_demo.db      # ENTITY (demonstração)
└── queijaria_digna.db       # ENTITY (exemplo)
```

### Princípios de Design

| Princípio | Descrição | Status |
|-----------|-----------|--------|
| **Isolamento Total** | Cada entidade tem seu próprio arquivo `.db` | ✅ Implementado |
| **Exit Power** | Entidade pode migrar seu database facilmente | ✅ Implementado |
| **Backup Simples** | Cópia do arquivo `.db` é backup completo | ✅ Implementado |
| **Performance** | SQLite é suficiente para escala de EES (1-100 usuários) | ✅ Validado |
| **Sem JOINs Cross-Tenant** | Proibido cruzar dados entre bancos diferentes | ✅ Regra Sagrada |

---

## 📊 Dualidade Arquitetural: Central vs. Tenant

### Banco Central (`central.db`)

**O que é:** Um arquivo de banco de dados único, gerido exclusivamente pelo módulo `lifecycle`.

**Responsabilidade:** Atua como o "Agregador Central" e motor de identidade do projeto como um todo.

**O que armazena:**
- Gestão de Identidade Global (usuários Gov.br, CPFs/CNPJs)
- Mapeamento físico e chaves criptográficas dos bancos dos Tenants
- Relacionamentos Cross-Tenant (RF-12): `EnterpriseAccountant` (vínculos contábeis)
- **NOVO - Módulo 2:** Cache de indicadores econômicos (`economic_indicators{}`, `indicator_cache{}`)
- **NOVO - Módulo 3:** Catálogo de programas de financiamento (`financing_programs{}`)
- **NOVO - RF-30:** Tópicos de ajuda educativa (`help_topics{}`)
- Metadados de intercooperação institucional

**Regra Inegociável:** O Banco Central **jamais** armazena transações financeiras, itens de estoque ou detalhamento operacional das entidades.

### Banco por Entidade (`data/entities/{entity_id}.db`)

**O que é:** O banco de dados físico isolado de cada empreendimento de economia solidária.

**Responsabilidade:** Materializar o Requisito Não Funcional de Soberania (RNF-01).

**O que armazena:**
- O **Ledger** contábil (partidas dobradas, histórico de caixa)
- Estoque, compras, vendas (PDV)
- Livro de Atas e Decisões de Assembleia (com Hashes SHA256)
- Registro de horas trabalhadas (Primazia do Trabalho - ITG 2002)
- **NOVO - Módulo 3:** Perfil de elegibilidade (`eligibility_profiles{}`)
- **NOVO - Módulo 3:** Match de programas (`program_matches{}`)
- **NOVO - Módulo 4:** Perfil público (`public_profiles{}`), mural de necessidades (`need_posts{}`)
- **NOVO - RF-27:** Cálculo de DAS MEI (`das_mei{}`)

**Regra Inegociável:** É tecnicamente impossível e proibido realizar **JOINs** (cruzamento de dados) entre o banco de um Tenant e o banco de outro Tenant.

---

## 📁 Schema Completo Atualizado (v2.0)

### 1. Banco Central (`central.db`)

```sql
-- ============================================
-- VÍNCULOS CONTÁBEIS (RF-12)
-- ============================================
CREATE TABLE IF NOT EXISTS enterprise_accountants (
    id TEXT PRIMARY KEY,
    enterprise_id TEXT NOT NULL,
    accountant_id TEXT NOT NULL,
    status TEXT NOT NULL, -- ACTIVE, INACTIVE
    start_date INTEGER NOT NULL, -- Unix timestamp
    end_date INTEGER DEFAULT 0,  -- Unix timestamp (0 = ativo)
    delegated_by TEXT NOT NULL,
    created_at INTEGER,
    updated_at INTEGER,
    UNIQUE(enterprise_id, accountant_id)
);

-- ============================================
-- MOTOR DE INDICADORES (RF-18) [NOVO - PDF v1.0]
-- ============================================
CREATE TABLE IF NOT EXISTS economic_indicators (
    id TEXT PRIMARY KEY,
    codigo TEXT NOT NULL,           -- Código da série (ex: 433 para IPCA)
    valor INTEGER NOT NULL,         -- int64 - Anti-Float (centavos de %)
    data_referencia INTEGER NOT NULL, -- Unix timestamp
    fonte TEXT NOT NULL,            -- BCB_SGS, BCB_PTAX, IBGE_SIDRA
    criado_em INTEGER,              -- Unix timestamp
    UNIQUE(codigo, data_referencia, fonte)
);

CREATE TABLE IF NOT EXISTS indicator_cache (
    id TEXT PRIMARY KEY,
    indicator_key TEXT NOT NULL UNIQUE, -- fonte + codigo
    value INTEGER NOT NULL,             -- int64 - Valor armazenado
    expires_at INTEGER NOT NULL,        -- Unix timestamp - TTL
    created_at INTEGER
);

-- ============================================
-- PORTAL DE OPORTUNIDADES (RF-20) [NOVO - PDF v1.0]
-- ============================================
CREATE TABLE IF NOT EXISTS financing_programs (
    id TEXT PRIMARY KEY,
    nome TEXT NOT NULL,
    fonte TEXT NOT NULL,            -- FEDERAL, ESTADUAL, MUNICIPAL
    valor_maximo INTEGER NOT NULL,  -- int64 - Anti-Float (centavos)
    taxa_juros INTEGER NOT NULL,    -- int64 - centavos de %
    prazo_maximo_meses INTEGER,
    carencia_meses INTEGER,
    requisitos TEXT,                -- JSON com requisitos de elegibilidade
    publico_prioritario TEXT,
    ativo INTEGER DEFAULT 1,
    link_edital TEXT,
    created_at INTEGER,
    updated_at INTEGER
);

-- ============================================
-- SISTEMA DE AJUDA EDUCATIVA (RF-30) [NOVO - Sessão 27/03/2026]
-- ============================================
CREATE TABLE IF NOT EXISTS help_topics (
    id TEXT PRIMARY KEY,
    key TEXT NOT NULL UNIQUE,           -- Chave única (ex: "cadunico", "inadimplencia")
    title TEXT NOT NULL,                -- Título em linguagem popular
    summary TEXT,                       -- Resumo em 1 frase (para tooltips)
    explanation TEXT NOT NULL,          -- Explicação completa em linguagem popular
    why_asked TEXT,                     -- "Por que perguntamos isso?"
    legislation TEXT,                   -- Legislação relacionada
    next_steps TEXT,                    -- Próximos passos acionáveis
    official_link TEXT,                 -- Link para fonte oficial
    category TEXT NOT NULL,             -- CREDITO, TRIBUTARIO, GOVERNANCA, GERAL
    tags TEXT,                          -- Tags para busca (JSON ou comma-separated)
    view_count INTEGER DEFAULT 0,       -- Quantas vezes foi visualizado
    created_at INTEGER,
    updated_at INTEGER
);

-- Índices para performance
CREATE INDEX IF NOT EXISTS idx_help_topics_category ON help_topics(category);
CREATE INDEX IF NOT EXISTS idx_help_topics_key ON help_topics(key);
```

### 2. Banco por Entidade (`data/entities/{entity_id}.db`)

```sql
-- ============================================
-- ENTERPRISE (Dados da Entidade)
-- ============================================
CREATE TABLE IF NOT EXISTS enterprises (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    cnpj TEXT,
    cnae TEXT,
    municipio TEXT,
    uf TEXT,
    faturamento_anual INTEGER,    -- int64 - Anti-Float (centavos)
    regime_tributario TEXT,
    data_abertura INTEGER,        -- Unix timestamp
    situacao_fiscal TEXT,
    status TEXT NOT NULL,         -- DREAM, INCUBATED, FORMALIZED
    created_at INTEGER,
    updated_at INTEGER
);

-- ============================================
-- CORE CONTÁBIL (Motor Lume)
-- ============================================
CREATE TABLE IF NOT EXISTS accounts (
    id INTEGER PRIMARY KEY,
    code TEXT NOT NULL,
    name TEXT NOT NULL,
    nature TEXT NOT NULL,         -- ASSET, LIABILITY, EQUITY, REVENUE, EXPENSE
    parent_id INTEGER,
    created_at INTEGER
);

CREATE TABLE IF NOT EXISTS entries (
    id TEXT PRIMARY KEY,
    entity_id TEXT NOT NULL,
    entry_date INTEGER NOT NULL,  -- Unix timestamp
    description TEXT,
    created_at INTEGER
);

CREATE TABLE IF NOT EXISTS postings (
    id TEXT PRIMARY KEY,
    entry_id TEXT NOT NULL,
    account_id INTEGER NOT NULL,
    amount INTEGER NOT NULL,      -- int64 - Anti-Float (centavos)
    direction TEXT NOT NULL,      -- DEBIT, CREDIT
    FOREIGN KEY (entry_id) REFERENCES entries(id),
    FOREIGN KEY (account_id) REFERENCES accounts(id)
);

-- ============================================
-- TRABALHO COOPERATIVO (ITG 2002)
-- ============================================
CREATE TABLE IF NOT EXISTS work_logs (
    id TEXT PRIMARY KEY,
    entity_id TEXT NOT NULL,
    member_id TEXT NOT NULL,
    activity TEXT,
    minutes INTEGER NOT NULL,     -- int64 - Anti-Float (minutos)
    date INTEGER NOT NULL,        -- Unix timestamp
    created_at INTEGER
);

CREATE TABLE IF NOT EXISTS members (
    id TEXT PRIMARY KEY,
    entity_id TEXT NOT NULL,
    name TEXT NOT NULL,
    email TEXT,
    phone TEXT,
    cpf TEXT,
    role TEXT NOT NULL,           -- COORDINATOR, MEMBER, ADVISOR
    status TEXT DEFAULT 'ACTIVE', -- ACTIVE, INACTIVE
    joined_at INTEGER,
    skills TEXT,                  -- JSON array
    created_at INTEGER,
    updated_at INTEGER,
    UNIQUE(entity_id, email)
);

-- ============================================
-- DECISÕES E GOVERNANÇA
-- ============================================
CREATE TABLE IF NOT EXISTS decisions (
    id TEXT PRIMARY KEY,
    entity_id TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    vote_type TEXT,               -- ABERTA, SECRETA
    result TEXT,                  -- APROVADA, REJEITADA, ADIADA
    votos_sim INTEGER,
    votos_nao INTEGER,
    votos_nulos INTEGER,
    hash_sha256 TEXT,             -- Hash de integridade
    created_at INTEGER
);

CREATE TABLE IF NOT EXISTS legal_documents (
    id TEXT PRIMARY KEY,
    entity_id TEXT NOT NULL,
    type TEXT NOT NULL,           -- ATA, ESTATUTO, DOSSIE_CADSOL, MTSE
    content TEXT,                 -- Conteúdo (Markdown)
    hash_sha256 TEXT,             -- Hash de integridade
    signature TEXT,               -- Assinatura eletrônica (Gov.br/ICP-Brasil)
    status TEXT,                  -- DRAFT, SIGNED, SUBMITTED
    created_at INTEGER
);

-- ============================================
-- FUNDOS OBRIGATÓRIOS
-- ============================================
CREATE TABLE IF NOT EXISTS funds (
    id TEXT PRIMARY KEY,
    entity_id TEXT NOT NULL,
    type TEXT NOT NULL,           -- RESERVA_LEGAL, FATES, OUTRO
    amount INTEGER NOT NULL,      -- int64 - Anti-Float (centavos)
    period TEXT NOT NULL,         -- YYYY-MM
    created_at INTEGER
);

-- ============================================
-- SUPPLY (Compras e Estoque)
-- ============================================
CREATE TABLE IF NOT EXISTS suppliers (
    id TEXT PRIMARY KEY,
    entity_id TEXT NOT NULL,
    name TEXT NOT NULL,
    contact_info TEXT,
    created_at INTEGER
);

CREATE TABLE IF NOT EXISTS stock_items (
    id TEXT PRIMARY KEY,
    entity_id TEXT NOT NULL,
    name TEXT NOT NULL,
    type TEXT NOT NULL,           -- INSUMO, PRODUTO, MERCADORIA
    unit TEXT NOT NULL,           -- KG, UN, L, etc.
    quantity INTEGER NOT NULL,    -- Quantidade atual
    min_quantity INTEGER,         -- Quantidade mínima (alerta)
    unit_cost INTEGER NOT NULL,   -- int64 - Anti-Float (centavos)
    created_at INTEGER,
    updated_at INTEGER
);

CREATE TABLE IF NOT EXISTS purchases (
    id TEXT PRIMARY KEY,
    entity_id TEXT NOT NULL,
    supplier_id TEXT NOT NULL,
    total_value INTEGER NOT NULL, -- int64 - Anti-Float (centavos)
    date INTEGER NOT NULL,        -- Unix timestamp
    payment_type TEXT,            -- CASH, CREDIT
    created_at INTEGER
);

CREATE TABLE IF NOT EXISTS purchase_items (
    id TEXT PRIMARY KEY,
    purchase_id TEXT NOT NULL,
    stock_item_id TEXT NOT NULL,
    quantity INTEGER NOT NULL,
    unit_cost INTEGER NOT NULL,   -- int64 - Anti-Float (centavos)
    total_cost INTEGER NOT NULL,  -- int64 - Anti-Float (centavos)
    FOREIGN KEY (purchase_id) REFERENCES purchases(id),
    FOREIGN KEY (stock_item_id) REFERENCES stock_items(id)
);

-- ============================================
-- BUDGET (Orçamento)
-- ============================================
CREATE TABLE IF NOT EXISTS budget_plans (
    id TEXT PRIMARY KEY,
    entity_id TEXT NOT NULL,
    category TEXT NOT NULL,
    planned_amount INTEGER NOT NULL, -- int64 - Anti-Float (centavos)
    period TEXT NOT NULL,            -- YYYY-MM
    created_at INTEGER,
    updated_at INTEGER
);

-- ============================================
-- PERFIL DE ELEGIBILIDADE (RF-19) [NOVO - PDF v1.0]
-- ============================================
CREATE TABLE IF NOT EXISTS eligibility_profiles (
    id TEXT PRIMARY KEY,
    entity_id TEXT NOT NULL UNIQUE,
    
    -- Dados do ERP (cópia para consulta rápida)
    cnpj TEXT,
    cnae TEXT,
    municipio TEXT,
    uf TEXT,
    faturamento_anual INTEGER,    -- int64 - Anti-Float
    regime_tributario TEXT,
    data_abertura INTEGER,
    situacao_fiscal TEXT,
    
    -- Campos complementares (preenchimento único)
    inscrito_cad_unico INTEGER,       -- 0/1 (bool)
    socio_mulher INTEGER,             -- 0/1 (bool)
    inadimplencia_ativa INTEGER,      -- 0/1 (bool)
    finalidade_credito TEXT,          -- CAPITAL_GIRO, EQUIPAMENTO, REFORMA, OUTRO
    valor_necessario INTEGER,         -- int64 - Anti-Float (centavos)
    tipo_entidade TEXT,               -- MEI, ME, EPP, Cooperativa, OSC, OSCIP, PF
    contabilidade_formal INTEGER,     -- 0/1 (bool)
    
    -- Metadados
    preenchido_em INTEGER,
    atualizado_em INTEGER,
    preenchido_por TEXT,
    created_at INTEGER,
    updated_at INTEGER
);

-- ============================================
-- MATCH DE PROGRAMAS (RF-20) [NOVO - PDF v1.0]
-- ============================================
CREATE TABLE IF NOT EXISTS program_matches (
    id TEXT PRIMARY KEY,
    entity_id TEXT NOT NULL,
    program_id TEXT NOT NULL,
    elegibilidade TEXT NOT NULL,    -- ELEGIVEL, NAO_ELEGIVEL, PARCIAL
    motivo TEXT,
    documentos_pendentes TEXT,      -- JSON com lista de documentos faltantes
    created_at INTEGER
);

-- ============================================
-- REDE DIGNA (RF-24, RF-25) [NOVO - PDF v1.0]
-- ============================================
CREATE TABLE IF NOT EXISTS public_profiles (
    entity_id TEXT PRIMARY KEY,
    nome_fantasia TEXT,
    missao TEXT,
    produtos TEXT,                  -- JSON array
    servicos TEXT,                  -- JSON array
    municipio TEXT,
    uf TEXT,
    contato_publico TEXT,
    foto_logo TEXT,
    created_at INTEGER,
    updated_at INTEGER
);

CREATE TABLE IF NOT EXISTS need_posts (
    id TEXT PRIMARY KEY,
    publisher_id TEXT NOT NULL,     -- Hash anonimizado
    categoria TEXT NOT NULL,        -- INSUMO, EQUIPAMENTO, SERVICO, OUTRO
    descricao TEXT,
    quantidade TEXT,
    prazo_desejado INTEGER,         -- Unix timestamp
    municipio TEXT,
    uf TEXT,
    status TEXT NOT NULL,           -- ABERTO, EM_NEGOCIACAO, CONCLUIDO
    created_at INTEGER,
    updated_at INTEGER
);

CREATE TABLE IF NOT EXISTS solidarity_transactions (
    id TEXT PRIMARY KEY,
    buyer_id TEXT NOT NULL,         -- Hash anonimizado
    seller_id TEXT NOT NULL,        -- Hash anonimizado
    valor INTEGER NOT NULL,         -- int64 - Anti-Float (centavos)
    descricao TEXT,
    data INTEGER NOT NULL,          -- Unix timestamp
    status TEXT NOT NULL,           -- COMPLETADA, CANCELADA
    created_at INTEGER
);

-- ============================================
-- DAS MEI (RF-27) [NOVO - PDF v1.0]
-- ============================================
CREATE TABLE IF NOT EXISTS das_mei (
    id TEXT PRIMARY KEY,
    entity_id TEXT NOT NULL,
    competencia TEXT NOT NULL,      -- YYYY-MM
    valor_devido INTEGER NOT NULL,  -- int64 - Anti-Float (centavos)
    valor_pago INTEGER DEFAULT 0,
    data_vencimento INTEGER NOT NULL,
    data_pagamento INTEGER DEFAULT 0,
    status TEXT NOT NULL,           -- PENDENTE, PAGO, VENCIDO
    salario_minimo INTEGER NOT NULL,
    created_at INTEGER,
    updated_at INTEGER,
    UNIQUE(entity_id, competencia)
);

-- ============================================
-- CONFORMIDADE ESTATAL (RF-14, RF-16) [NOVO - Adequação Estatal]
-- ============================================
CREATE TABLE IF NOT EXISTS reinf_events (
    id TEXT PRIMARY KEY,
    entity_id TEXT NOT NULL,
    event_type TEXT NOT NULL,       -- R-2010, R-2020, R-2099, etc.
    xml_content TEXT,
    hash_sha256 TEXT,
    status TEXT NOT NULL,           -- PENDING, SENT, CONFIRMED, ERROR
    sent_at INTEGER,
    created_at INTEGER
);

CREATE TABLE IF NOT EXISTS sanitary_dossiers (
    id TEXT PRIMARY KEY,
    entity_id TEXT NOT NULL,
    tipo_produto TEXT,              -- QUEIJO, MEL, CARNES, LATICINIOS
    capacidade_diaria INTEGER,
    fluxo_maquinario TEXT,
    origem_agua TEXT,               -- POÇO, REDE, NASCENTE
    content TEXT,                   -- Conteúdo do MTSE (Markdown/PDF)
    status TEXT NOT NULL,           -- DRAFT, SUBMITTED, APPROVED
    created_at INTEGER
);

-- ============================================
-- MIGRAÇÕES E METADADOS
-- ============================================
CREATE TABLE IF NOT EXISTS supply_migrations (
    id TEXT PRIMARY KEY,
    version TEXT NOT NULL,
    applied_at INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS sync_metadata (
    id TEXT PRIMARY KEY,
    entity_id TEXT NOT NULL,
    last_sync_at INTEGER,
    sync_status TEXT,
    updated_at INTEGER
);

-- ============================================
-- ÍNDICES PARA PERFORMANCE
-- ============================================
CREATE INDEX IF NOT EXISTS idx_entries_entity_date ON entries(entity_id, entry_date);
CREATE INDEX IF NOT EXISTS idx_postings_entry ON postings(entry_id);
CREATE INDEX IF NOT EXISTS idx_work_logs_entity ON work_logs(entity_id);
CREATE INDEX IF NOT EXISTS idx_members_entity ON members(entity_id);
CREATE INDEX IF NOT EXISTS idx_decisions_entity ON decisions(entity_id);
CREATE INDEX IF NOT EXISTS idx_stock_items_entity ON stock_items(entity_id);
CREATE INDEX IF NOT EXISTS idx_purchases_entity ON purchases(entity_id);
CREATE INDEX IF NOT EXISTS idx_eligibility_entity ON eligibility_profiles(entity_id);
CREATE INDEX IF NOT EXISTS idx_program_matches_entity ON program_matches(entity_id);
CREATE INDEX IF NOT EXISTS idx_das_mei_entity ON das_mei(entity_id, competencia);
```

---

## 🔄 Fluxo de Dados Validado

### 1. Compra → Estoque
```
Fornecedor → purchases → purchase_items → stock_items
```

### 2. Estoque → PDV
```
stock_items (filtro: PRODUTO/MERCADORIA) → PDV UI
```

### 3. PDV → Caixa
```
Venda no PDV → entries/postings → cash_flow
```

### 4. Caixa → Dashboard
```
entries → Cálculo Saldo → Dashboard UI
```

### 5. ERP → Portal de Oportunidades [NOVO]
```
enterprises + eligibility_profiles → program_matches
```

### 6. Motor → Portal [NOVO]
```
economic_indicators → financing_programs (ranqueamento por taxa)
```

### 7. UI → Ajuda Educativa [NOVO - RF-30]
```
Campo técnico → help_topics (via key) → Tooltip/Modal
```

---

## 📊 Database `cafe_digna` - Estado Atual

### Dados Inseridos (09/03/2026)

**1. Fornecedores**
```sql
INSERT INTO suppliers (id, name, contact_info, created_at) VALUES 
('supplier_001', 'Fazenda Café Bom', 'João Silva - (11) 99999-9999', [timestamp]);
```

**2. Itens em Estoque**

| Item | Tipo | Quantidade | Unidade | Custo Unitário | Valor Total |
|------|------|------------|---------|----------------|-------------|
| Café em Grão Arábica | INSUMO | 50kg | KG | R$ 45,00 | R$ 2.250,00 |
| Café Moído para Venda | PRODUTO | 20kg | KG | R$ 80,00 | R$ 1.600,00 |
| Café Torrado em Grão | MERCADORIA | 30kg | KG | R$ 70,00 | R$ 2.100,00 |

**Total Estoque:** 3 itens × R$ 5.950,00

**3. Compras Registradas**
```sql
INSERT INTO purchases (id, supplier_id, total_value, date, created_at) VALUES 
('purchase_001', 'supplier_001', 385000, [timestamp], [timestamp]);
-- Total: R$ 3.850,00 (50kg × R$45 + 20kg × R$80)
```

---

## 🛠️ Scripts de População

### `test_cafe_digna_fixed.sql`

Script completo para popular o database com dados reais:

```sql
-- 1. Adicionar fornecedor
INSERT INTO suppliers (id, name, contact_info, created_at) VALUES 
('supplier_001', 'Fazenda Café Bom', 'João Silva - (11) 99999-9999', strftime('%s', 'now'));

-- 2. Adicionar itens ao estoque
INSERT INTO stock_items (id, name, type, unit, quantity, min_quantity, unit_cost, created_at) VALUES 
('item_001', 'Café em Grão Arábica', 'INSUMO', 'KG', 50, 10, 4500, strftime('%s', 'now')),
('item_002', 'Café Moído para Venda', 'PRODUTO', 'KG', 20, 5, 8000, strftime('%s', 'now')),
('item_003', 'Café Torrado em Grão', 'MERCADORIA', 'KG', 30, 8, 7000, strftime('%s', 'now'));

-- 3. Adicionar compra
INSERT INTO purchases (id, supplier_id, total_value, date, created_at) VALUES 
('purchase_001', 'supplier_001', 385000, strftime('%s', 'now'), strftime('%s', 'now'));

-- 4. Adicionar itens da compra
INSERT INTO purchase_items (id, purchase_id, stock_item_id, quantity, unit_cost, total_cost) VALUES 
('pi_001', 'purchase_001', 'item_001', 50, 4500, 225000),
('pi_002', 'purchase_001', 'item_002', 20, 8000, 160000);
```

### Como Executar

```bash
# Conectar ao database
sqlite3 data/entities/cafe_digna.db

# Executar script
.read test_cafe_digna_fixed.sql

# Verificar dados
SELECT * FROM stock_items;
SELECT * FROM purchases;
```

---

## 📈 Impacto no Sistema

### Antes da Sessão (09/03/2026)
- ❌ Database vazio
- ❌ PDV sem produtos para venda
- ❌ Módulo de compras vazio
- ❌ Estoque vazio
- ❌ Fluxo completo impossível de testar

### Após a Sessão (09/03/2026)
- ✅ Database populado com dados reais
- ✅ PDV com 2 produtos disponíveis (50kg total)
- ✅ Módulo de compras com histórico
- ✅ Estoque com 3 itens (R$ 5.950,00)
- ✅ Fluxo completo validado

### Após Expansão do Ecossistema (27/03/2026)
- ✅ Schema expandido para 4 módulos
- ✅ Tabelas para Indicadores, Portal, Rede
- ✅ Tabela `help_topics` para RF-30
- ✅ Tabela `das_mei` para RF-27
- ✅ Tabelas de conformidade estatal (RF-14, RF-16)

---

## 🔍 Análise de Dados

### Itens Disponíveis no PDV

O PDV filtra itens do estoque por:
- **Tipo:** `PRODUTO` ou `MERCADORIA`
- **Quantidade:** `> 0`

**Resultado para `cafe_digna`:**
- `Café Moído para Venda` (PRODUTO) - 20kg disponíveis
- `Café Torrado em Grão` (MERCADORIA) - 30kg disponíveis

**Total disponível para venda:** 50kg

### Valorização do Estoque

| Categoria | Itens | Quantidade | Valor |
|-----------|-------|------------|-------|
| INSUMO | 1 | 50kg | R$ 2.250,00 |
| PRODUTO | 1 | 20kg | R$ 1.600,00 |
| MERCADORIA | 1 | 30kg | R$ 2.100,00 |
| **TOTAL** | **3** | **100kg** | **R$ 5.950,00** |

### Compras Registradas

- **Fornecedor:** Fazenda Café Bom
- **Valor Total:** R$ 3.850,00
- **Itens:** 2 (café em grão + café moído)
- **Status:** Database pronto para operações

---

## 🚀 Próximos Passos para Database

### 1. Scripts de Migração [NOVO - Expandido]
- [ ] Migrações versionadas para novos módulos (RF-18 a RF-30)
- [ ] Rollback seguro
- [ ] Validação de integridade pós-migração

### 2. Backup Automatizado
- [ ] Backup periódico (diário)
- [ ] Compressão (.tar.gz)
- [ ] Restauração simplificada

### 3. Análise de Dados
- [ ] Relatórios avançados
- [ ] Tendências
- [ ] Previsões

### 4. Sincronização
- [ ] Sync entre dispositivos
- [ ] Resolução de conflitos
- [ ] Offline-first

### 5. Segurança
- [ ] Criptografia opcional (SQLCipher)
- [ ] Controle de acesso refinado
- [ ] Auditoria de acessos

### 6. Seed de Tópicos de Ajuda [NOVO - RF-30]
- [ ] Popular `help_topics` com 10+ tópicos iniciais
- [ ] Categorias: CRÉDITO, TRIBUTÁRIO, GOVERNANÇA, GERAL
- [ ] Validação de linguagem popular (5ª série)

---

## ⚠️ Considerações Importantes

### 1. Performance SQLite

| Vantagens | Limitações | Adequação |
|-----------|------------|-----------|
| Simplicidade | Concorrência escrita limitada | Perfeito para escala de EES (1-100 usuários) |
| Portabilidade | Tamanho máximo (140TB teórico) | Mais que suficiente para EES |
| Zero configuração | Sem users/permissions nativos | Gerenciado via aplicação |

### 2. Isolamento vs. Relatórios Consolidados

| Isolamento | Desafio | Solução |
|------------|---------|---------|
| Soberania do dado | Relatórios multi-entidade | Sync para nuvem (opcional) com consentimento |
| Exit Power | Consolidação para contador | Painel multi-tenant com acesso Read-Only |

### 3. Backup e Recovery

- **Backup:** Cópia do arquivo `.db`
- **Recovery:** Restaurar arquivo
- **Versioning:** Git LFS para histórico (opcional)

### 4. Migração para Outros Sistemas

- **Export:** SQL dump padrão
- **Import:** Qualquer sistema que suporte SQL
- **Interoperabilidade:** Máxima possível

---

## 📊 Validação Anti-Float

**Regra Sagrada:** Todos os campos monetários e de tempo devem usar `int64`.

**Validação:**
```bash
grep -r "float[0-9]*" modules/
# Deve retornar apenas logs/comentários
```

**Campos Críticos:**
- `faturamento_anual INTEGER` — Centavos
- `minutes INTEGER` — Minutos trabalhados
- `valor_necessario INTEGER` — Centavos
- `valor_maximo INTEGER` — Centavos (programas de crédito)
- `view_count INTEGER` — Contagem de visualizações (Help System)

---

## ✅ Conclusão

O database do Digna está completamente atualizado para suportar o **Ecossistema de 4 Módulos** (PDF v1.0) e o **Sistema de Ajuda Educativa** (RF-30, Sessão 27/03/2026).

**Status do Database:**
- ✅ Schema expandido para todos os módulos
- ✅ `cafe_digna.db` populado com dados reais
- ✅ Dualidade Central vs. Tenant preservada
- ✅ Anti-Float compliance em todas as tabelas
- ✅ Índices de performance criados

**Próximas Ações:**
1. Popular `help_topics` com seed inicial (RF-30)
2. Testar fluxo completo com dados reais
3. Validar relatórios com dados populados
4. Documentar procedimentos de backup/migração

---

**Status:** ✅ ATUALIZADO COM VISÃO DE ECOSSISTEMA (PDF v1.0) + RF-30 (Sessão 27/03/2026)  
**Próxima Ação:** Atualizar `04_governance/governance.md` com novas responsabilidades  
**Versão Anterior:** 1.0 (2026-03-09)  
**Versão Atual:** 2.0 (2026-03-27)
