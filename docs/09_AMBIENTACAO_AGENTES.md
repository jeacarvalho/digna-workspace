## 📄 `09_AMBIENTACAO_AGENTES.md` (O Orientador)

```markdown
# 🚀 Contexto do Projeto: Digna (Providentia Foundation)

## 📌 Visão Geral
**Projeto:** Sistema de Gestão Contábil Público e Open Source para Economia Solidária.
**Propósito:** Transformar a contabilidade de um peso burocrático em uma ferramenta de **Dignidade Financeira** para Empreendimentos de Economia Solidária (EES).
**Diferencial:** Arquitetura *Local-First Server-Side* com isolamento físico de dados.

---

## ✅ Status Atual das Sprints

| Sprint | Status | Descrição |
| :--- | :--- | :--- |
| Sprint 00 | ✅ COMPLETE | Blueprint, Governança e Visão de Produto. |
| **Sprint 01** | ⏭️ **EM ANDAMENTO** | **Lifecycle Manager & Core Ledger (Motor Lume)** |

---

## 🎯 Sprint 01 (Tarefa Atual)
**Objetivo:** Implementar o motor em Go que gerencia o ciclo de vida dos bancos SQLite individuais e o motor de lançamentos contábeis.

**Componentes Chave:**
- **Lume (Motor):** O backend em Go.
- **Digna (App):** A interface (futura).
- **Providentia (Fundação):** A entidade de governança.

---

## 🛠️ Stack Tecnológica
- **Linguagem:** Go 1.22+ (Clean Architecture).
- **Banco de Dados:** libSQL / SQLite (Individual por Tenant).
- **Escala:** Alvo de 1 milhão de EES (Indução Serpro).
- **Regra Financeira:** Tudo em `int64` (centavos). Proibido `float`.

---

## ⚠️ Regras para o Agente
1. **Soberania de Dados:** O dado pertence à entidade. O acesso deve ser isolado por arquivo.
2. **Padrão de Lançamento:** Seguir rigorosamente o princípio de partidas dobradas (soma zero).
3. **Mocks de Formalização:** O sistema deve prever a transição `DREAM` -> `FORMALIZED` via interfaces (Facade).

```

---

## 📄 `00_PROJECT_BRIEF.md`

```markdown
# Project Brief: Digna (Providentia Foundation)

**Status:** Sprint 01 - Initial Implementation ⏭️  
**Next Phase:** Sprint 02 - Onboarding & Member Management  
**Scale Target:** 1.000.000+ Entities

---

## Objective
Prover uma infraestrutura contábil soberana para a Economia Solidária brasileira, hospedada na nuvem do **Serpro** e gerida pela **Fundação Providentia**. O sistema facilita a jornada desde o grupo informal ("Sonho") até a entidade formalizada.

## Tech Stack
| Component | Technology |
| :--- | :--- |
| Backend | Go (Golang) |
| Multi-tenancy | SQLite-per-Tenant (Isolated Files) |
| Architecture | Hexagonal / Clean Architecture |
| Sync Engine | CDC (Change Data Capture) via Go |
| Integrity | Double-Entry Bookkeeping (Strict) |

## Critical Business Rules

### 1. Multi-Tenancy por Arquivo
Cada cooperativa/associação tem seu próprio arquivo `.sqlite`. Isso garante que o backup e a portabilidade sejam feitos no nível do arquivo físico.

### 2. Contabilidade Social
O motor deve suportar o **Rateio de Sobras** baseado em métricas sociais (horas trabalhadas, aportes iniciais) e não apenas em capital.

### 3. Facilitação da Formalização
O sistema atua como um incubador. Ele gera Minutas de Estatutos e Atas de Fundação automaticamente a partir dos dados registrados na fase informal (`DREAM`).

---

## Execution Timeline

| Sprint | Status | Description |
| :--- | :--- | :--- |
| S01 | ⏭️ ACTIVE | Lifecycle Manager (SQLite Orchestration) & Core Ledger. |
| S02 | 📋 PLANNED | Member Registry & Initial Capital Inflow. |
| S03 | 📋 PLANNED | Social Surplus Algorithm (Motor de Rateio). |
| S04 | 📋 PLANNED | Serpro/REDESIM Integration Mock. |

```

---

