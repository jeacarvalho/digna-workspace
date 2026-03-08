---
title: Roadmap de Produto
status: implemented
version: 1.1
last_updated: 2026-03-08
---

# Roadmap de Produto - Digna

O roadmap segue as 4 fases estratégicas definidas no DVP (Documento de Visão e Escopo).

---

## Phase 1: Foundational (Fundação e Operação)

**Status:** ✅ COMPLETE

**Objetivo:** Prover a infraestrutura básica de operação contábil.

**Entregas:**
- [x] PDV operacional
- [x] Motor de partidas dobradas (Core Lume)
- [x] Registro de horas ITG 2002
- [x] Lifecycle Manager (SQLite por tenant)
- [x] Dashboard de dignidade
- [x] Interface web PWA
- [x] Gestão de membros
- [x] Testes E2E BDD

**Maturidade:** Sistema operacional para grupos informais

---

## Phase 2: Formalization (Formalização Automática)

**Status:** ✅ COMPLETE

**Objetivo:** Automatar a transição de grupos informais para entidades formalizadas.

**Entregas:**
- [x] Simulador de formalização
- [x] Geração de atas (Markdown)
- [x] Transição automática DREAM → FORMALIZED
- [x] Integrações governamentais (mock)
- [ ] Integração Gov.br (real)
- [ ] Dossiê CADSOL/DCSOL automático
- [ ] Geração de documentos oficiais
- [ ] Integração com órgãos públicos (real)

**Maturidade:** Grupos informais tornam-se visíveis ao Estado

---

## Phase 3: Financial (Finanças Solidárias)

**Status:** 🟡 EM DESENVOLVIMENTO

**Objetivo:** Suporte a múltiplas unidades de valor e gestão financeira.

**Entregas:**
- [ ] Gestão de compras (RF-07)
- [ ] Gestão de estoque (RF-08)
- [x] Gestão de caixa (RF-09) ✅
- [x] SurplusCalculator com deduções ✅
- [x] Rateio automático de sobras ✅
- [ ] Gestão orçamentária (RF-10)
- [ ] Moedas sociais
- [ ] Estoque substantivo (sementes, animais)

**Maturidade:** Sistema gerencia riqueza além do Real (R$)

---

## Phase 4: Network (Intercooperação Nacional)

**Status:** 🔵 PLANEJADO

**Objetivo:** Criar rede nacional de economia solidária interconectada.

**Entregas:**
- [ ] Marketplace B2B
- [ ] Score de crédito social
- [ ] Integração BNDES/Serpro
- [ ] API pública
- [ ] Rede nacional de EES conectadas

**Maturidade:** Ecossistema completo de economia solidária

---

## Marcos (Milestones)

| Marco | Fase | Status | Previsão |
|-------|------|--------|----------|
| MVP Operacional | 1 | ✅ COMPLETE | 03/2026 |
| Formalização Beta | 2 | ✅ COMPLETE | 03/2026 |
| E2E Journey Tests | 1-2 | ✅ COMPLETE | 03/2026 |
| Gestão Financeira | 3 | 🟡 | 06/2026 |
| Rede Nacional | 4 | 🔵 | 2027 |

---

## Dependências entre Fases

```
Phase 1 (Foundational)
        ↓
Phase 2 (Formalization)
        ↓
Phase 3 (Financial)
        ↓
Phase 4 (Network)
```

**Nota:** Cada fase depende da conclusão da anterior. A Phase 4 só é viável após:
- Phase 1: Infraestrutura funcionando
- Phase 2: Entidades formalizadas
- Phase 3: Gestão financeira madura
