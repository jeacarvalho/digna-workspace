# 01 ARCHITECTURE - Digna (Providentia Foundation)

**Status:** Architecture Implemented (v0.3)
**Last Updated:** 2026-03-07

---

## 1. High-Level System Design (Local-First Server-Side)

O Digna utiliza uma arquitetura de **Micro-databases isolados**. Em vez de um banco único, cada entidade possui sua própria instância física, garantindo soberania e escalabilidade.

```text
┌─────────────────────────────────────────────────────────────────────────┐
│                           CORE ENGINE (LUME)                            │
│  ┌─────────────┐      ┌─────────────┐      ┌─────────────┐              │
│  │   Go API    │──────┤ Lifecycle   │──────┤  SQLite Per │              │
│  │   (REST)    │      │ Manager     │      │   Tenant    │              │
│  └──────┬──────┘      └──────┬──────┘      └──────┬──────┘              │
└─────────┼────────────────────┼────────────────────┼─────────────────────┘
          │                    │                    │
          ▼                    ▼                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                        DATA PERSISTENCE                                 │
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐       │
│  │ /data/entities/  │  │ /data/entities/  │  │ /data/entities/  │       │
│  │  ent_A.db        │  │  ent_B.db        │  │  ent_C.db        │       │
│  │  (WAL mode)      │  │  (WAL mode)      │  │  (WAL mode)      │       │
│  └────────┬─────────┘  └────────┬─────────┘  └────────┬─────────┘       │
└───────────┼─────────────────────┼─────────────────────┼─────────────────┘
            ▼                     ▼                     ▼
      ┌───────────┐         ┌───────────┐         ┌───────────┐
      │ Litestream│         │ Litestream│         │ Litestream│ (Backup S3)
      └───────────┘         └───────────┘         └───────────┘
      
  
    
---

## 2. Tecnologias Core

| Camada | Tecnologia | Justificativa |
| --- | --- | --- |
| **Backend** | Go (1.22+) | Performance, concorrência e binário estático para Nuvem Serpro. |
| **Database** | SQLite3 + mattn/go-sqlite3 | Isolamento total por arquivo, WAL mode, foreign keys. |
| **Sync** | Change Data Capture | Sincronização assíncrona para o Agregador Central da Fundação. |
| **Arquitetura** | Clean Architecture | Desacoplamento de domínio (não importa driver SQL). |
| **Numerics** | int64 (exclusivo) | Valores financeiros e tempo sem erros IEEE 754. |

---

## 3. Módulos Implementados

### 3.1 lifecycle (Sprint 01 ✅)
**Path:** `modules/lifecycle/`
**Responsabilidade:** Orquestração de arquivos SQLite por tenant

```
lifecycle/
├── internal/
│   ├── domain/
│   │   ├── entity.go        # Entity (DREAM/FORMALIZED)
│   │   └── interfaces.go    # LifecycleManager, Migrator
│   ├── manager/
│   │   └── sqlite_mgr.go    # Pool de conexões + PRAGMAs
│   └── repository/
│       └── migration.go     # DDL v0 (6 tabelas)
└── manager_test.go          # 6 testes de integridade
```

**Tabelas Criadas (Schema v0):**
- `accounts` - Plano de contas (hierárquico)
- `entries` - Lançamentos contábeis
- `postings` - Partidas dobradas (amount: int64)
- `work_logs` - ITG 2002 (minutes: int64)
- `decisions_log` - CADSOL (autogestão)
- `sync_metadata` - Versão e último sync

**Contas Padrão (Seed):**
| ID | Código | Nome | Tipo |
|----|--------|------|------|
| 1 | 1.1.01 | Caixa e Equivalentes | ASSET |
| 2 | 3.1.01 | Receita de Vendas | REVENUE |
| 3 | 1.1.02 | Bancos | ASSET |
| 4 | 2.1.01 | Fornecedores | LIABILITY |

---

### 3.2 core_lume (Sprint 02 ✅)
**Path:** `modules/core_lume/`
**Responsabilidade:** Motor contábil, governança e valoração social

```
core_lume/
├── pkg/
│   ├── ledger/          # API Pública - Journaling & Balanço
│   │   └── ledger.go    # Transaction, Posting, Service
│   ├── social/          # API Pública - ITG 2002
│   │   └── social.go    # WorkRecord, Work Capital
│   └── governance/      # API Pública - CADSOL
│       └── governance.go # DecisionRecord, Hash SHA256
└── internal/            # Implementações internas
    ├── ledger/
    ├── social/
    └── governance/
```

**Regras de Negócio:**
- **Integridade Contábil:** Soma(Débitos) = Soma(Créditos) = 0
- **Atomicidade:** Transações SQLite garantem consistência
- **ITG 2002:** Tempo de trabalho convertido em capital social
- **CADSOL:** Decisões imutáveis com hash SHA256

---

### 3.3 pdv_ui (Sprint 02 ✅)
**Path:** `modules/pdv_ui/`
**Responsabilidade:** Interface de operações comerciais (fachada)

```
pdv_ui/
├── usecase/
│   └── operation.go     # RecordSale, RecordWork, RecordDecision
└── pdv_test.go        # Testes integrados end-to-end
```

**Operações Disponíveis:**
- `RecordSale(entityID, amount, method)` → Débito Caixa + Crédito Vendas
- `RecordWork(entityID, memberID, minutes)` → Registro ITG 2002
- `RecordDecision(entityID, title, content)` → Hash + Log CADSOL

**Princípios:**
- Clean Architecture: Usecases não conhecem SQL
- Core Lume como Gatekeeper de integridade
- Lifecycle Manager como único ponto de I/O

---

### 3.4 reporting (Sprint 03 ✅)
**Path:** `modules/reporting/`
**Responsabilidade:** Cálculos agregados e relatórios institucionais

```
reporting/
├── internal/
│   └── surplus/
│       └── calculator.go  # Algoritmo de rateio social
└── pkg/
    └── surplus/
        └── surplus.go     # API pública para consultas
```

**Funcionalidades:**
- `CalculateSocialSurplus(entityID)` → Distribuição proporcional por horas
- Rateio: (Horas do Sócio / Total de Horas) × Excedente Financeiro
- Integração ITG 2002 + Contabilidade (valores em int64)

**Exemplo de Rateio:**
| Sócio | Horas | % do Total | Valor (R$ 100,00) |
|-------|-------|------------|-------------------|
| socio_001 | 600 | 66.7% | R$ 66.66 |
| socio_002 | 300 | 33.3% | R$ 33.33 |

---

### 3.5 legal_facade (Sprint 03 ✅)
**Path:** `modules/legal_facade/`
**Responsabilidade:** Documentação institucional e formalização

```
legal_facade/
├── internal/document/
│   ├── generator.go        # Gerador de Atas (Markdown)
│   ├── identity.go         # Cartões de identificação
│   └── formalization.go    # Simulador DREAM→FORMALIZED
└── pkg/
    └── document/
        └── document.go     # API pública
```

**Documentos Gerados:**
- **Atas de Assembleia:** Markdown com hash SHA256 de auditoria
- **Cartões de Identidade:** Dados da entidade + status (DREAM/FORMALIZED)
- **Simulador de Formalização:** Transição automática após 3 decisões

**Critérios de Formalização:**
- Mínimo de 3 decisões registradas no `decisions_log`
- Alteração de status na tabela `sync_metadata`
- Geração de identidade com CNPJ mock (00.000.000/0000-00)

---

## 4. Fluxos de Dados (Sequência)

### 4.1 Ciclo de Vida do Tenant
O Lifecycle Manager é o único componente com permissão para realizar operações de `OS` (criação de arquivos). O sistema utiliza o padrão 'Lazy Initialization', criando o banco apenas no primeiro acesso ou no cadastro explícito.

Usuário/App          API (Digna)        Lifecycle Mgr       SQLite (Disco)
    |                   |                   |                   |
    |-- POST /dream --> |                   |                   |
    |   {name: "Mel"}   |-- InitTenant(id) -|                   |
    |                   |                   |-- Criar mel.db -->|
    |                   |                   |                   |
    |                   |<- DB Connection --|                   |
    |                   |                   |                   |
    |                   |-- RunMigrations --|                   |
    |                   |                   |-- CREATE TABLEs ->|
    |                   |                   |                   |
    | <--- 201 Created -| <--- Success -----|                   |
    |    (Entity_ID)    |                   |                   |

### 4.2 Integridade do Ledger
O motor Lume atua como um 'Gatekeeper'. Nenhuma escrita no SQLite ocorre sem que o motor de validação de partidas dobradas confirme que a soma dos lançamentos resulta em zero (Equilíbrio Contábil).

Usuário/App          API (Digna)         Lume (Ledger)      SQLite (mel.db)
    |                   |                   |                   |
    |-- POST /entry --> |                   |                   |
    | {D:100, C:100}    |-- Process(entry) -|                   |
    |                   |                   |                   |
    |                   |                   |-- Validar Soma? --|
    |                   |                   |   (D + C == 0)    |
    |                   |                   |         |         |
    |                   |                   |      OK [v]       |
    |                   |                   |         |         |
    |                   |                   |-- Begin Trans. -->|
    |                   |                   |-- Insert Entry -->|
    |                   |                   |-- Insert Posts -->|
    |                   |                   |-- Commit -------->|
    |                   |                   |                   |
    | <--- 200 OK ------| <--- Success -----|                   |
    
## 4.3 Formalização (Transição de Status)

    Usuário/App          API (Digna)        Legal Facade       Lifecycle Mgr
    |                   |                   |                   |
    |-- POST /formal -->|                   |                   |
    |                   |-- RequestCNPJ() ->|                   |
    |                   |                   |                   |
    |                   |                   |-- Gerar Mock ---->|
    |                   |                   |   (00.000.../01)  |
    |                   |                   |                   |
    |                   |<-- New Identity --|                   |
    |                   |                   |                   |
    |                   |-- UpdateStatus -->|-- Set FORMALIZED -|
    |                   |                   |    no mel.db      |
    |                   |                   |                   |
    | <--- 200 OK ------|                   |                   |
    | (Estatuto PDF)    |                   |                   |
    
### 4.4 Operação de Venda (PDV → Lume → Ledger)
Fluxo completo de uma venda no PDV: validação de dados, validação contábil, e persistência atômica.

    Usuário/App        PDV UI          Core Lume         Lifecycle Mgr    SQLite
        |                |                 |                 |             |
        |-- Venda 50,00->|                 |                 |             |
        |   (PIX)        |-- RecordSale()--|                 |             |
        |                |                 |                 |             |
        |                |                 |-- Validate(txn)--|             |
        |                |                 |   (D+C==0?)      |             |
        |                |                 |      [OK]       |             |
        |                |                 |                 |             |
        |                |                 |-- Persist(txn)---|-- INSERT ---->|
        |                |                 |                 |   entries    |
        |                |                 |                 |   postings   |
        |                |                 |                 |             |
        |                |                 |<-- Commit -------|<-- TX OK ----|
        |                |<-- EntryID -----|                 |             |
        |<-- Recibo ----|                 |                 |             |
    
### 4.5 Registro de Trabalho (ITG 2002)
Registro de horas de trabalho como capital social do cooperado.

    Usuário/App        PDV UI          Core Lume         Lifecycle Mgr    SQLite
        |                |                 |                 |             |
        |-- 8h trabalho->|                 |                 |             |
        |   membro_123   |-- RecordWork()--|                 |             |
        |                |                 |                 |             |
        |                |                 |-- Validate()----|             |
        |                |                 |   (minutes>0)    |             |
        |                |                 |      [OK]       |             |
        |                |                 |                 |             |
        |                |                 |-- Persist()------|-- INSERT ---->|
        |                |                 |                 |  work_logs   |
        |                |                 |                 |             |
        |                |<-- Confirmação -|                 |             |
        |<-- OK ---------|                 |                 |             |

### 4.6 Auditoria de Decisão (CADSOL)
Registro imutável de decisões de assembleia com hash criptográfico.

    Usuário/App        PDV UI          Core Lume         Lifecycle Mgr    SQLite
        |                |                 |                 |             |
        |-- Decisão ---->|                 |                 |             |
        |   "Aprovar     |-- RecordDecision                  |             |
        |    Orçamento"  |                 |                 |             |
        |                |                 |                 |             |
        |                |                 |-- SHA256(content)|            |
        |                |                 |   hash=abc123...  |             |
        |                |                 |                 |             |
        |                |                 |-- Persist()------|-- INSERT ---->|
        |                |                 |                 | decisions_log|
        |                |                 |                 |  (hash+status)|
        |                |                 |                 |             |
        |                |<-- hash=abc123 -|                 |             |
        |<-- Hash/Docs ---|                 |                 |             |

### 4.7 Cálculo de Rateio Social (Reporting)
Processamento do excedente financeiro distribuído proporcionalmente às horas de trabalho.

    Usuário/App        PDV UI        Reporting         Lifecycle Mgr    SQLite
        |                |               |                 |             |
        |-- Rateio? --->|               |                 |             |
        |                |-- Calculate()--|                 |             |
        |                |               |                 |             |
        |                |               |-- Query Work ----|             |
        |                |               |   (work_logs)    |             |
        |                |               |-- Query Sales ---|             |
        |                |               |   (postings)     |             |
        |                |               |                 |             |
        |                |               |-- Calcular()     |            |
        |                |               |   % = hrs/total  |             |
        |                |               |   val = % × exc  |             |
        |                |               |                 |             |
        |                |<-- Result -----|                 |             |
        |<-- socio_001:  |               |                 |             |
            R$ 66.66     |               |                 |             |

### 4.8 Geração de Ata (Legal Facade)
Transformação de decisões do banco em documento institucional Markdown.

    Usuário/App        PDV UI       Legal Facade        Lifecycle Mgr    SQLite
        |                |              |                 |             |
        |-- Gerar Ata -->|              |                 |             |
        |                |-- Generate()--|                 |             |
        |                |              |                 |             |
        |                |              |-- Query Decisions|             |
        |                |              |  (decisions_log) |             |
        |                |              |                 |             |
        |                |              |-- Template()    |             |
        |                |              |   Markdown fmt  |             |
        |                |              |   + Hash audit  |             |
        |                |              |                 |             |
        |                |<-- MD Doc ----|                 |             |
        |<-- Ata.md -----|              |                 |             |

---

## 5. Camadas de Segregação

### 5.1 Interface Layer (PDV UI)
- **Responsabilidade:** Receber intenções do usuário, validar formato
- **Não faz:** Acesso direto ao banco, lógica contábil
- **Exemplo:** `usecase/operation.go` - RecordSale()

### 5.2 Domain Layer (Core Lume)
- **Responsabilidade:** Regras de negócio, integridade, auditoria
- **Faz:** Validação de partidas dobradas, hash de decisões, cálculo de capital
- **Exemplo:** `pkg/ledger/service.go` - Validate(), RecordTransaction()

### 5.3 Infrastructure Layer (Lifecycle)
- **Responsabilidade:** I/O físico, conectividade, migrações
- **Faz:** Criar arquivos .db, aplicar PRAGMAs, executar DDL
- **Exemplo:** `pkg/lifecycle/sqlite.go` - GetConnection()

### 5.4 Reporting Layer
- **Responsabilidade:** Cálculos agregados, análise de dados, rateios
- **Faz:** Consultar SQLite, calcular distribuições, gerar métricas
- **Exemplo:** `pkg/surplus/calculator.go` - CalculateSocialSurplus()

### 5.5 Document Layer (Legal Facade)
- **Responsabilidade:** Geração de documentos institucionais
- **Faz:** Templates, formatação Markdown, hash de auditoria
- **Exemplo:** `pkg/document/generator.go` - GenerateAssemblyMinutes()

---

## 6. Stack Tecnológico Atualizado

| Camada | Tecnologia | Uso |
|--------|-----------|-----|
| **Backend** | Go 1.22+ | API REST, concorrência |
| **Storage** | SQLite3 | Isolamento por tenant |
| **Driver** | mattn/go-sqlite3 | CGO bindings |
| **Migrations** | SQL nativo | DDL versionado |
| **Templates** | Go html/template | Documentos Markdown |
| **Hash** | SHA256 | Auditoria CADSOL |
| **Numerics** | int64 | Centavos e minutos |
| **Rateio** | Proporcional por tempo | ITG 2002 + Contábil |
| **Documents** | Markdown | Atas de Assembleia |
| **Architecture** | Clean Architecture | Separação de concerns |
| **Workspace** | Go Modules | Multi-module monorepo |
