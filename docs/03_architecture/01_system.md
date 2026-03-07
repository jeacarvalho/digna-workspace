---
title: Arquitetura do Sistema Digna
status: implemented
version: 1.0
last_updated: 2026-03-07
---

# Arquitetura do Sistema - Digna

**Projeto:** Sistema de Gestão Contábil para Economia Solidária  
**Arquitetura:** Local-First Server-Side com Micro-databases Isolados

---

## 1. Visão Geral da Arquitetura

O Digna utiliza uma arquitetura de **Micro-databases isolados**. Em vez de um banco único, cada entidade possui sua própria instância física, garantindo soberania e escalabilidade.

### High-Level Design

```
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
```

---

## 2. Tecnologias Core

| Camada | Tecnologia | Justificativa |
|--------|------------|---------------|
| **Backend** | Go (1.22+) | Performance, concorrência e binário estático para Nuvem Serpro |
| **Database** | SQLite3 + mattn/go-sqlite3 | Isolamento total por arquivo, WAL mode, foreign keys |
| **Sync** | Change Data Capture | Sincronização assíncrona para o Agregador Central |
| **Arquitetura** | Clean Architecture | Desacoplamento de domínio |
| **Numerics** | int64 (exclusivo) | Valores financeiros e tempo sem erros IEEE 754 |

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

**Responsabilidades:**
- Criação física de arquivos `.db`
- Pool de conexões com PRAGMAs otimizados
- Migrações de schema versionado
- Lazy initialization (cria banco no primeiro acesso)

**PRAGMAs Configurados:**
- `journal_mode = WAL`
- `foreign_keys = ON`
- `synchronous = NORMAL`
- `temp_store = MEMORY`

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
└── pdv_test.go          # Testes integrados end-to-end
```

**Operações Disponíveis:**
- `RecordSale(entityID, amount, method)` → Débito Caixa + Crédito Vendas
- `RecordWork(entityID, memberID, minutes)` → Registro ITG 2002
- `RecordDecision(entityID, title, content)` → Hash + Log CADSOL

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
- Atas de Assembleia (Markdown com hash SHA256)
- Cartões de Identidade
- Simulador de Formalização

---

### 3.6 sync_engine (Sprint 04 ✅)

**Path:** `modules/sync_engine/`  
**Responsabilidade:** Sincronização offline-first e intercooperação B2B

```
sync_engine/
├── internal/
│   ├── tracker/
│   │   └── sqlite_delta.go    # Monitor de alterações
│   ├── exchange/
│   │   └── intercoop.go       # Marketplace B2B
│   └── client/
│       └── provider_sync.go   # Sync seguro
└── sprint04_test.go           # 9 testes DoD
```

---

### 3.7 ui_web (Sprint 05 ✅)

**Path:** `modules/ui_web/`  
**Responsabilidade:** Interface web mobile-first para operação diária

```
ui_web/
├── main.go                    # Servidor HTTP porta 8080
├── internal/handler/
│   ├── pdv_handler.go         # Rotas PDV e API
│   └── dashboard.go           # Dashboard e Social Clock
├── templates/
│   ├── layout.html            # Base com navegação
│   ├── pdv.html               # Teclado numérico vendas
│   ├── social_clock.html      # Registro de horas
│   └── dashboard.html         # Painel de dignidade
├── static/
│   ├── manifest.json          # Config PWA
│   └── sw.js                  # Service Worker cache
└── sprint05_test.go           # 9 testes DoD
```

**Stack Frontend:**
| Tecnologia | Versão | Uso |
|------------|--------|-----|
| Go net/http | Nativo | Servidor web |
| html/template | Nativo | Server-side rendering |
| HTMX | 1.9.10 | Atualizações parciais AJAX |
| Tailwind CSS | CDN | Design mobile-first |
| PWA | - | Instalação mobile |

---

## 4. Camadas de Segregação

### 4.1 Interface Layer (PDV UI)
- **Responsabilidade:** Receber intenções do usuário, validar formato
- **Não faz:** Acesso direto ao banco, lógica contábil
- **Exemplo:** `usecase/operation.go` - RecordSale()

### 4.2 Domain Layer (Core Lume)
- **Responsabilidade:** Regras de negócio, integridade, auditoria
- **Faz:** Validação de partidas dobradas, hash de decisões, cálculo de capital
- **Exemplo:** `pkg/ledger/service.go` - Validate(), RecordTransaction()

### 4.3 Infrastructure Layer (Lifecycle)
- **Responsabilidade:** I/O físico, conectividade, migrações
- **Faz:** Criar arquivos .db, aplicar PRAGMAs, executar DDL
- **Exemplo:** `pkg/lifecycle/sqlite.go` - GetConnection()

### 4.4 Reporting Layer
- **Responsabilidade:** Cálculos agregados, análise de dados, rateios
- **Faz:** Consultar SQLite, calcular distribuições, gerar métricas

### 4.5 Document Layer (Legal Facade)
- **Responsabilidade:** Geração de documentos institucionais
- **Faz:** Templates, formatação Markdown, hash de auditoria

### 4.6 Web Layer (UI Web)
- **Responsabilidade:** Interface humana, servidor HTTP, PWA
- **Faz:** Renderizar templates, servir assets estáticos, HTMX endpoints

---

## 5. Fluxos de Dados

### 5.1 Ciclo de Vida do Tenant

```
Usuário/App          API (Digna)        Lifecycle Mgr       SQLite (Disco)
    |                   |                   |                   |
    |-- POST /dream --> |                   |                   |
    |   {name: "Mel"}   |-- InitTenant(id) -|                   |
    |                   |                   |-- Criar mel.db -->|
    |                   |                   |                   |
    |                   |<- DB Connection --|                   |
    |                   |                   |-- RunMigrations -->|
    |                   |                   |                   |
    | <--- 201 Created -| <--- Success -----|                   |
    |    (Entity_ID)    |                   |                   |
```

### 5.2 Integridade do Ledger

```
Usuário/App          API (Digna)         Lume (Ledger)      SQLite (mel.db)
    |                   |                   |                   |
    |-- POST /entry --> |                   |                   |
    | {D:100, C:100}    |-- Process(entry) -|                   |
    |                   |                   |-- Validar Soma? --|
    |                   |                   |   (D + C == 0)    |
    |                   |                   |      [OK]         |
    |                   |                   |-- Persist(txn)----|
    |                   |                   |                   |
    | <--- 200 OK ------| <--- Success -----|                   |
```

---

## 6. Stack Tecnológico Final

| Camada | Tecnologia | Uso |
|--------|-----------|-----|
| Backend | Go 1.22+ | API REST, concorrência |
| Storage | SQLite3 | Isolamento por tenant |
| Driver | mattn/go-sqlite3 | CGO bindings |
| Migrations | SQL nativo | DDL versionado |
| Templates | Go html/template | Documentos Markdown |
| Hash | SHA256 | Auditoria CADSOL |
| Numerics | int64 | Centavos e minutos |
| Rateio | Proporcional por tempo | ITG 2002 + Contábil |
| Documents | Markdown | Atas de Assembleia |
| Architecture | Clean Architecture | Separação de concerns |
| Workspace | Go Modules | Multi-module monorepo |

---

## 7. Estratégia de Módulos para Agentes

Para mitigar a entropia de contexto em agentes de IA, o sistema Digna é dividido em módulos independentes:

| Módulo | Consome de | Entrega para |
|--------|------------|--------------|
| Lifecycle | Sistema de Arquivos | Ledger / Legal Facade |
| Ledger | Lifecycle (DB Connection) | Reporting Engine |
| Legal Facade | Lifecycle (Metadata) | Usuário Final (PDF) |
| Reporting | Ledger (Transaction History) | Usuário Final (UI) |

**Estratégia de Implementação:**
1. **Isolamento de Sprint:** Cada módulo com SESSION_LOG próprio
2. **Interface First:** Definir interfaces em `internal/domain`
3. **Validation:** Validar Módulo 1 antes do Módulo 2
