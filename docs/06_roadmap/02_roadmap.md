***

```markdown
---
title: Roadmap de Produto
status: implemented
version: 1.3
last_updated: 2026-03-09
---

# Roadmap de Produto - Digna

O roadmap segue as 4 fases estratégicas definidas no DVP (Documento de Visão e Escopo) e incorpora a aliança com a classe contábil para ganho de escala na Economia Solidária.

---

## Phase 1: Foundational (Fundação e Operação)

**Status:** ✅ COMPLETE
**Objetivo:** Prover a infraestrutura básica de operação contábil e a "Contabilidade Invisível".

**Entregas:**
- [x] PDV operacional com interface pedagógica
- [x] Motor de partidas dobradas e Soma Zero em `int64` (Core Lume)
- [x] Registro de horas e valoração do suor ITG 2002
- [x] Lifecycle Manager (SQLite isolado por tenant)
- [x] Dashboard de dignidade e transparência algorítmica
- [x] Interface web PWA (Offline-first)
- [x] Gestão de membros
- [x] Testes E2E BDD (Jornada Anual "Sonho Solidário")

**Maturidade:** Sistema operacional empoderando grupos informais a gerirem seus negócios com rigor, mas sem jargões.

---

## Phase 2: Integração Institucional e Aliança Contábil

**Status:** ✅ COMPLETE (Sprint 12 Finalizada)
**Objetivo:** Automatizar a transição de grupos informais para entidades formalizadas e criar a ponte tecnológica estrutural com os Contadores Sociais (CFC/CRCs).

**Entregas:**
- [x] Simulador de formalização e algoritmos de gatilho
- [x] Geração de atas (Markdown com Hash SHA256)
- [x] Transição automática DREAM → FORMALIZED
- [x] Integrações governamentais (Mocks via Clean Architecture)
- [ ] Integração Gov.br (Autenticação real via OAuth2)
- [x] **[NOVO] Painel do Contador Social (Accountant Dashboard):** Interface Multi-tenant *Read-Only* para auditores voluntários. ✅ **SPRINT 12**
- [x] **[NOVO] Exportação Fiscal (SPED):** Motor de tradução das partidas dobradas geradas pelo Core Lume para os leiautes contábeis e fiscais exigidos pela Receita Federal. ✅ **SPRINT 12**
- [x] **[NOVO] Testes E2E Atualizados:** Jornada "Sonho Solidário" inclui auditorias do Contador Social ✅ **SPRINT 12 E2E**
- [ ] **[NOVO] Módulos educativos embutidos:** Auxílio na formação de preço considerando a hora trabalhada.

**Maturidade:** Grupos informais tornam-se visíveis ao Estado e amparados legalmente pela classe contábil, sem a necessidade de o produtor virar um "digitador de notas". A arquitetura garante Soberania do Dado com acesso read-only para contadores.

---

## Phase 3: Finanças Solidárias e Territoriais

**Status:** 🟡 EM DESENVOLVIMENTO (Sprint 13 Concluída)
**Objetivo:** Suporte a múltiplas unidades de valor e fortalecimento da economia local.

**Entregas:**
- [x] **Gestão de Compras (RF-07):** Módulo supply com contabilidade invisível ✅ **SPRINT 13**
- [x] **Controle de Estoque (RF-08):** Categorização INSUMO/PRODUTO/MERCADORIA ✅ **SPRINT 13**
- [ ] Gestão orçamentária (RF-10)
- [ ] Integração tecnológica com Bancos Comunitários de Desenvolvimento (BCDs)
- [ ] Gestão e transação de Moedas Sociais locais
- [ ] Estoque substantivo (Controle e troca de sementes, animais, horas-trabalho)

**Maturidade:** Retenção da riqueza no território e independência do sistema financeiro tradicional.

---

## Phase 4: Intercooperação Nacional

**Status:** 🔵 PLANNED
**Objetivo:** Conectar EES isolados em uma rede nacional de apoio e viabilidade econômica.

**Entregas:**
- [ ] Marketplace B2B restrito a EES (Privacidade preservada)
- [ ] Score de Crédito Social baseado em histórico de trabalho (ITG 2002)
- [ ] Integração de crédito com BNDES/Políticas Públicas

**Maturidade:** O ecossistema Digna atua como espinha dorsal do SINAES.

---

## Marcos (Milestones)

| Marco | Fase | Status | Previsão |
|-------|------|--------|----------|
| MVP Operacional | 1 | ✅ COMPLETE | 03/2026 |
| Formalização Beta | 2 | ✅ COMPLETE | 03/2026 |
| E2E Journey Tests | 1-2 | ✅ COMPLETE | 03/2026 |
| **Painel do Contador & SPED** | 2 | ✅ COMPLETE | 03/2026 |
| **Gestão de Compras e Estoque** | 3 | ✅ COMPLETE | 03/2026 |
| Gestão Financeira Territorial | 3 | 🟡 | 07/2026 |
| Rede Nacional | 4 | 🔵 | 2027 |

---

## Dependências entre Fases

```text
Phase 1 (Foundational)
       ↓
Phase 2 (Integração e Aliança Contábil)
       ↓
Phase 3 (Financial Territorial)
       ↓
Phase 4 (National Network)
```

**Nota de Arquitetura:** Cada fase depende restritamente da conclusão da anterior. A Fase 2 (especificamente o Painel do Contador) depende que o *Motor Lume* da Fase 1 mantenha a integridade inquebrável de Soma Zero e `int64`, pois o fisco e o SPED rejeitarão dados com erros de dízima.
```

***
