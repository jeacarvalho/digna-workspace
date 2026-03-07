---

### 📄 `docs/00_project_brief.md`

**Objetivo:** Definir a visão do produto e o diferencial competitivo (PDV).

```markdown
# Project Brief: Digna (Providentia Foundation)

**Status:** Sprint 01 - Initial Implementation ⏭️  
**Last Updated:** 2026-03-05  
**Next Phase:** Sprint 02 - Core Ledger & PDV Integration  
**Scale Target:** 1.000.000+ Entities

---

## Objective
Prover uma infraestrutura contábil soberana e invisível para a Economia Solidária brasileira. O sistema foca no **PDV (Ponto de Venda)** como ferramenta principal do empreendedor, transformando vendas e compras em registros contábeis automáticos de partidas dobradas, facilitando a jornada desde o grupo informal ("Sonho") até a entidade formalizada.

## Tech Stack
| Component | Technology |
| :--- | :--- |
| Backend | Go (Golang) 1.22+ |
| Multi-tenancy | SQLite-per-Tenant (Isolated Files) |
| Architecture | Hexagonal / Clean Architecture |
| Workspace | Go Multi-module Workspace (`go.work`) |
| Integrity | Double-Entry Bookkeeping (Strict int64) |

## Critical Business Rules

### 1. Contabilidade Invisível via PDV
O empreendedor não faz lançamentos contábeis; ele opera o negócio. O sistema traduz ações comerciais (vendas/compras) em lançamentos de débito/crédito automaticamente no motor **Lume**.

### 2. Multi-Tenancy por Arquivo
Cada entidade possui seu próprio arquivo `.sqlite`. Isso garante soberania total, permitindo que o usuário baixe seus dados e facilite backups via Litestream.

### 3. Rigor Financeiro
Todo cálculo monetário é realizado em `int64` (centavos). O uso de `float` é estritamente proibido para evitar erros de arredondamento.

---

## Execution Timeline

| Sprint | Status | Description |
| :--- | :--- | :--- |
| S01 | ⏭️ ACTIVE | Lifecycle Manager (Orquestração SQLite) & Base Ledger. |
| S02 | 📋 PLANNED | PDV Module & Invisible Accounting Mapping. |
| S03 | 📋 PLANNED | Member Registry & Social Surplus Algorithm. |
| S04 | 📋 PLANNED | Legal Facade & Formalization Mocks. |

```

---
