title: Status Atual - Ecossistema Digna
status: implemented
version: 2.0
last_updated: 2026-03-27
---

# Status Atual - Ecossistema Digna

> **Nota:** Este documento reflete o status integrado do Ecossistema Digna (PDF v1.0), preservando toda a infraestrutura validada nas Sprints 1-16, incorporando os novos módulos como Fases 3 e 4, e incluindo o Sistema de Ajuda Educativa (RF-30) decidido na sessão de 27/03/2026.

---

## 📋 Contexto da Atualização (27/03/2026)

**Motivação:** O projeto Digna evoluiu de um ERP contábil para um **Ecossistema de 4 Módulos** conforme especificação PDF v1.0. Esta atualização documenta:

1. **Expansão da Visão de Produto:** De ERP único para ecossistema integrado (ERP + Motor + Portal + Rede)
2. **Novos Marcos (Milestones):** Fases 3-6 adicionadas ao roadmap original
3. **RF-30 (Ajuda Educativa):** Sistema transversal de pedagogia social
4. **Preservação:** Todas as entregas validadas nas Sprints 1-16 mantidas

**Versão Anterior:** 1.0 (2026-03-13)  
**Nova Versão:** 2.0 (2026-03-27)

---

## 🎯 Status Geral do Projeto

**Última Atualização:** 2026-03-27 (Sessão de Expansão do Ecossistema + RF-30)  
**Fase Atual:** Sprint 17 (Expansão do Ecossistema - Fase 3)  
**Status:** 🟡 **EM DESENVOLVIMENTO - ECOSSISTEMA INTEGRADO**  
**Próximo Marco:** Validação do MVP do Portal de Oportunidades com entidades reais de Niterói

---

## 🗺️ Visão Geral das Fases

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    STATUS DO ECOSSISTEMA DIGNA                          │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  FASE 1 — FUNDAÇÃO (✅ COMPLETE)                                        │
│  ├── Sprint 01-06: Core Contábil + PDV + Ledger                         │
│  ├── Sprint 07-12: Contador Social + SPED                               │
│  └── Sprint 13-16: Supply + Budget + UI 100%                            │
│                                                                         │
│  FASE 2 — CONFORMIDADE ESTATAL (🟡 EM ANDAMENTO)                        │
│  ├── RF-14: EFD-Reinf + ECF (Blindagem Tributária)                     │
│  ├── RF-15: Gov.br + Assembleias Digitais                              │
│  ├── RF-16: MTSE/MAPA (Inclusão Sanitária)                             │
│  └── RF-17: CADSOL/SINAES Automático                                   │
│                                                                         │
│  FASE 3 — ECOSSISTEMA DE CRÉDITO (📋 NOVO - PDF v1.0)                   │
│  ├── RF-18: Motor de Indicadores (BCB, IBGE APIs)                      │
│  ├── RF-19: Perfil de Elegibilidade (campos complementares)            │
│  ├── RF-20: Portal de Oportunidades (MVP: 3 programas)                 │
│  └── RF-21: Checklist + Alertas de Documentos                          │
│                                                                         │
│  FASE 4 — REDE DE INTERCOOPERAÇÃO (📋 NOVO - PDF v1.0)                  │
│  ├── RF-24: Perfil Público da Entidade                                 │
│  ├── RF-25: Mural de Necessidades                                      │
│  └── RF-26: Match de Oportunidades B2B                                 │
│                                                                         │
│  FASE 5 — AJUDA E PEDAGOGIA (📋 NOVO - Decisão 27/03/2026)              │
│  └── RF-30: Sistema de Ajuda Educativa Estruturada                     │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 📊 Phase Status Overview

| Phase | Marco | Status | Conclusão |
|-------|-------|--------|-----------|
| Concepção | Marco 00 | ✅ COMPLETE | 100% |
| Foundation Setup | Marco 01 | ✅ COMPLETE | 100% |
| Core Operations | Marco 02 | ✅ COMPLETE | 100% |
| Reporting & Documents | Marco 03 | ✅ COMPLETE | 100% |
| UI & Dashboard | Marco 04 | ✅ COMPLETE | 100% |
| Integração e Aliança Contábil (Phase 2) | Marco 07 | ✅ COMPLETE | 100% |
| Gestão de Compras e Estoque (Phase 3) | Marco 08 | ✅ COMPLETE | 100% |
| Gestão Orçamentária e Planejamento Financeiro (Phase 3) | Marco 06 | ✅ COMPLETE | 100% |
| **Adequação Estatal e Conformidade (Phase 2)** | **Marco 09** | 🟡 **EM DESENVOLVIMENTO** | **0%** |
| **Ecossistema de Crédito (Phase 3 - NOVO)** | **Marco 10** | 📋 **BACKLOG** | **0%** |
| **Rede de Intercooperação (Phase 4 - NOVO)** | **Marco 11** | 📋 **BACKLOG** | **0%** |
| **Ajuda e Pedagogia (Transversal - NOVO)** | **Marco 12** | 📋 **BACKLOG** | **0%** |
| Production Deploy | Marco 05 | 🟢 READY | 95% |

---

## 🏛️ Foco Atual: Adequação Estatal e Conformidade Digital [ATUALIZADO]

*O projeto atingiu a maturidade operacional contábil (Soma Zero / int64) na Sprint 16 e foca agora na construção de pontes automáticas de conformidade institucional para proteger as EES antes de ir para a produção final.*

### ✅ Entregas Concluídas (Sprints 1-16)

- [x] Lifecycle Manager (SQLite isolado por tenant para garantia de soberania)
- [x] Core Lume (Partidas dobradas invisíveis em int64)
- [x] PDV Interface (Foco pedagógico e linguagem popular livre de jargões)
- [x] Registro de trabalho (ITG 2002 - Valoração do tempo em int64)
- [x] Dashboard de Dignidade (Transparência visual para a Assembleia)
- [x] Painel do Contador Social (Accountant Dashboard Multi-tenant) ✅ **Sprint 12**
- [x] Exportação Fiscal (SPED/CSV) ✅ **Sprint 12**
- [x] Gestão de Compras e Estoque (RF-07, RF-08) ✅ **Sprint 13-14**
- [x] Gestão Orçamentária (RF-10) ✅ **Sprint 14**
- [x] Identidade Visual "Soberania e Suor" ✅ **Sprint 16**
- [x] Templates Cache-Proof ✅ **Sprint 16**
- [x] Testes E2E com Playwright ✅ **Sprint 15-16**

### 🟡 Entregas Pendentes (Adequação Estatal)

- [ ] **RF-14: Blindagem Tributária (EFD-Reinf e ECF)** 
  * Módulo `tax_compliance` para mensageria de retenções via Web Service
  * Expurgo automático de receitas de Atos Cooperativos no Bloco M da ECF (Lei 5.764/71 e LC 214/2025)
  * Geração e transmissão de XMLs (série R-2000/R-4000) para Web Services da EFD-Reinf
  * Alimentação automática da DCTFWeb

- [ ] **RF-15: Integração Real Gov.br e Governança Digital**
  * Substituição do Mock de login unificado pelo fluxo real da Cidadania Digital (OAuth2)
  * Assinatura Eletrônica Qualificada (Lei nº 14.063/2020) para membros da mesa nas Atas de Assembleia
  * Algoritmo de anonimização sistêmica para escrutínio secreto em votações (IN DREI nº 79/2020)
  * Integração com ICP-Brasil para certificados digitais

- [ ] **RF-16: Inclusão Sanitária (MAPA)**
  * Módulo `sanitary_compliance` para geração automatizada do Memorial Técnico Sanitário de Estabelecimento (MTSE)
  * Conformidade com Portaria MAPA nº 393/2021 para agroindústrias
  * Parametrização de fluxogramas de maquinário, capacidade diária e potabilidade de água
  * Exportação para peticionamento no Sistema Eletrônico de Informação (SEI)

- [ ] **RF-17: Integração CADSOL/SINAES Automático**
  * Consumo nativo de Web Services do MTE (Decreto nº 12.784/2025)
  * Matrícula automática de entidades FORMALIZED no Cadastro Nacional de Economia Solidária
  * Substituição dos mocks atuais por integração real

---

## 🚀 Foco Novo: Ecossistema de Crédito e Indicadores [NOVO - PDF v1.0]

*O projeto expande para conectar automaticamente o empreendedor aos programas de crédito para os quais ele já é elegível, sustentado por dados econômicos atualizados em tempo real.*

### 📋 Módulo 2: Motor de Indicadores (RF-18)

| Sub-requisito | Descrição | Fonte de Dados | Prioridade | Status |
|--------------|-----------|---------------|------------|--------|
| RF-18.1 | Coleta de SELIC, IPCA, CDI | BCB SGS API | Alta | 📋 Backlog |
| RF-18.2 | Câmbio oficial (USD/BRL, EUR/BRL) | BCB PTAX API | Alta | 📋 Backlog |
| RF-18.3 | Expectativas de mercado (Focus) | BCB Focus API | Média | 📋 Backlog |
| RF-18.4 | Indicadores sociais (PNAD, desocupação) | IBGE SIDRA | Média | 📋 Backlog |
| RF-18.5 | Cache local e interpretação contextual | Arquitetura interna | Alta | 📋 Backlog |

**Arquitetura Proposta:**
```
modules/
└── indicators_engine/          # NOVO MÓDULO
    ├── internal/
    │   ├── collector/          # Scheduler de coleta (cron diário)
    │   ├── cache/              # Tabela indicators{} com TTL
    │   ├── interpreter/        # Regras de negócio contextualizadas
    │   └── repository/         # SQLite local (cache-proof)
    ├── pkg/
    │   └── indicators/         # API pública para outros módulos
    └── cmd/
        └── collector/          # Binário independente para coleta
```

### 📋 Módulo 3: Portal de Oportunidades (RF-19 a RF-23)

**Programas Prioritários para MVP:**
| Programa | Valor Máximo | Público-Alvo | Diferencial | Status |
|----------|-------------|--------------|-------------|--------|
| **Acredita no Primeiro Passo** | ~R$ 6 mil | CadÚnico, 70% mulheres | Sem garantia, juros subsidiados | 📋 Backlog |
| **Pronampe** | Até 30% da receita bruta | MEI, ME, EPP | Prazo até 84 meses, carência 24 meses | 📋 Backlog |
| **Niterói Empreendedora** | Até R$ 200 mil | CNPJ ativo 12+ meses em Niterói | **JUROS ZERO**, foco em mulheres e jovens | 📋 Backlog |

**Funcionalidades do MVP:**
- [ ] **Match Automático:** Cruzamento do perfil do ERP com requisitos de cada programa
- [ ] **Ranqueamento por Vantagem:** Ordenação por custo efetivo de capital (usando taxas do Motor)
- [ ] **Checklist de Documentos:** Para cada linha elegível, marca o que já está no Digna vs. o que falta
- [ ] **Alertas de Prazo:** Notificações quando editais com alto match têm prazo próximo

### 📋 Sistema Transversal: Ajuda Educativa (RF-30) [NOVO - Sessão 27/03/2026]

**Descrição:** Sistema de ajuda contextual que traduz conceitos técnicos (CadÚnico, inadimplência, CNAE, etc.) em linguagem popular, com linkagem entre elementos de UI e registros de ajuda no banco.

**Funcionalidades:**
- [ ] Entrada de menu "Ajuda" acessível em todas as páginas
- [ ] Busca e índice de tópicos categorizados (CRÉDITO, TRIBUTÁRIO, GOVERNANÇA, GERAL)
- [ ] Explicação em linguagem popular + legislação relacionada + próximo passo acionável
- [ ] Linkagem automática: botão "?" ao lado de campos técnicos abre explicação contextual
- [ ] Tópicos com estrutura: Título, Resumo, Explicação, Por que perguntamos, Legislação, Próximo passo, Link oficial

**Critério de Aceite Pedagógico:**
- [ ] Usuário com 5ª série consegue entender a explicação sem ajuda externa
- [ ] Tooltip carrega em < 500ms via HTMX
- [ ] Conteúdo não usa jargões técnicos ("cadastramento", "regularização fiscal", etc.)
- [ ] Sempre inclui "próximo passo" acionável (ex: "procure o CRAS")

---

## 🌐 Foco Futuro: Rede de Intercooperação [NOVO - PDF v1.0]

*Conectar EES isolados em uma rede nacional de apoio e viabilidade econômica, materializando o 6º Princípio do Cooperativismo.*

### 📋 Módulo 4: Rede Digna (RF-24 a RF-26)

**Funcionalidades:**
- [ ] **Perfil Público da Entidade (RF-24):** Missão, produtos, serviços e capacidades disponíveis para a rede
- [ ] **Mural de Necessidades (RF-25):** Entidades publicam demandas de compra (insumos, serviços, equipamentos) visíveis para a rede
- [ ] **Match de Oportunidades (RF-26):** Sistema sugere conexões entre quem precisa e quem oferece, priorizando proximidade geográfica e afinidade de setor
- [ ] **Histórico de Transações Solidárias:** Registro que pode ser apresentado como evidência de atividade econômica em candidaturas a financiamento

**Fundamento Teológico da Rede:**
O conceito de **Koinonía** — comunhão e partilha mútua na tradição cristã primitiva — é a metáfora central da Rede Digna. Não é apenas um marketplace: é uma ecclesia econômica onde entidades que compartilham valores se apoiam mutuamente, reduzindo dependência de intermediários externos e fortalecendo a economia local.

**Nota sobre Laicidade:** A Teologia informa o design internamente, mas o produto permanece acessível independentemente de crença. O canal de distribuição via igrejas e comunidades de fé é estratégico, não doutrinário.

---

## 🏃 Sprint Status [ATUALIZADO]

### Sprint 01 a 06: Core e Fundação ✅
- Lifecycle Manager (SQLite isolado)
- Ledger Service (partidas dobradas exatas em `int64`)
- Surplus Calculator (rateio social ITG 2002)
- PDV Interface (HTMX + Tailwind)
- Social Clock e Dashboard
- **Testes:** 37/37 PASS

### Sprint 07 a 12: Integração e Aliança Contábil ✅
- DDD Refactoring (43/43 PASS)
- Integrações Governamentais (Mocks) (13/13 PASS)
- Gestão de Membros (19/19 PASS)
- Formalização e E2E Journey (5/5 PASS)
- **Painel do Contador Social (Accountant Dashboard)** (8/8 + E2E PASS)
- **Exportação Fiscal (SPED)** (Validado)
- **Testes:** 93/93 PASS

### Sprint 13 a 16: Finanças Solidárias e UI 100% ✅
- Gestão de Compras e Estoque (RF-07, RF-08) (6/6 PASS)
- Gestão Orçamentária (RF-10) (4/4 PASS)
- Correções Críticas + Testes E2E (PDV→Estoque→Caixa) (3/3 + E2E PASS)
- **Identidade Visual "Soberania e Suor"** (149/149 PASS)
- **Testes:** 162/162 PASS

### Sprint 17+: Expansão do Ecossistema 🟡
- **RF-18: Motor de Indicadores** (📋 Backlog - Fase 3)
- **RF-19: Perfil de Elegibilidade** (📋 Backlog - Fase 3)
- **RF-20: Portal de Oportunidades (MVP)** (📋 Backlog - Fase 3)
- **RF-24 a RF-26: Rede Digna** (📋 Backlog - Fase 4)
- **RF-30: Sistema de Ajuda Educativa** (📋 Backlog - Transversal)

---

## 📈 Total Test Coverage [ATUALIZADO]

| Sprint | Testes | Status | Notas |
|--------|--------|--------|-------|
| 01 | 6/6 | ✅ PASS | Lifecycle Manager |
| 02 | 8/8 | ✅ PASS | Core Lume + PDV |
| 03 | 8/8 | ✅ PASS | Reporting + Legal |
| 04 | 9/9 | ✅ PASS | Sync Engine |
| 05 | 9/9 | ✅ PASS | UI Web (PWA) |
| 06 | 3/3 | ✅ PASS | Cash Flow |
| 07 | 43/43 | ✅ PASS | DDD Refactoring |
| 08-09 | 13/13 | ✅ PASS | Integrações (Mocks) |
| 10 | 19/19 | ✅ PASS | Gestão de Membros |
| 11 | 5/5 | ✅ PASS | Formalização e E2E |
| 12 | 8/8 | ✅ PASS | Accountant Dashboard + SPED |
| 13 | 6/6 | ✅ PASS | Gestão de Compras e Estoque |
| 14 | 4/4 | ✅ PASS | Gestão Orçamentária |
| 15 | 3/3 | ✅ PASS | Correções Críticas + E2E |
| 16 | 149/149 | ✅ PASS | Identidade Visual + Sistema 100% |
| **Total** | **311/311** | **100% PASS** 🎉 | **Base sólida para expansão** |

---

## 🏗️ DDD Architecture Status [ATUALIZADO]

| Módulo | Interface Repository | Implementação | Status |
|--------|---------------------|-----------------|--------|
| core_lume | LedgerRepository, WorkRepository, DecisionRepository | SQLite | ✅ COMPLETE |
| reporting | SurplusRepository | Adapter Pattern | ✅ COMPLETE |
| sync_engine | SyncRepository | SQLite | ✅ COMPLETE |
| legal_facade | LegalRepository | SQLite | ✅ COMPLETE |
| integrations | 8 interfaces governamentais | Mock | ✅ COMPLETE |
| accountant_dashboard | FiscalRepository | Read-Only SQLite Adapter | ✅ COMPLETE |
| supply | SupplyRepository | SQLite | ✅ COMPLETE |
| budget | BudgetRepository | SQLite | ✅ COMPLETE |
| ui_web | CashHandler, PDVHandler, etc. | HTTP Handlers | ✅ COMPLETE |
| **indicators_engine** | **IndicatorRepository** | **SQLite (central.db)** | 📋 BACKLOG |
| **portal_opportunities** | **ProgramRepository, MatchRepository** | **SQLite (entity.db)** | 📋 BACKLOG |
| **rede_digna** | **PublicProfileRepository, NeedRepository** | **SQLite (entity.db)** | 📋 BACKLOG |
| **help_engine** | **HelpTopicRepository** | **SQLite (central.db)** | 📋 BACKLOG |

---

## 🎯 Critérios de Sucesso por Fase [ATUALIZADO]

| Fase | Métrica de Sucesso | Alvo | Status |
|------|-------------------|------|--------|
| Fase 1 | Entidades operando com contabilidade invisível | 100+ entidades ativas | ✅ 311/311 testes |
| Fase 2 | Entidades com conformidade estatal automatizada | 80% das formalizadas | 🟡 Em desenvolvimento |
| Fase 3 | Entidades descobrindo elegibilidade via Portal | 50+ matches bem-sucedidos | 📋 Backlog |
| Fase 4 | Transações B2B realizadas na Rede | 500+ conexões/ano | 📋 Backlog |
| Fase 5 (RF-30) | Redução de abandono em formulários | 30% redução | 📋 Backlog |
| Fase 5 (RF-30) | Tópicos de ajuda visualizados/mês | 500+ visualizações | 📋 Backlog |

---

## ⚠️ Riscos e Mitigações [ATUALIZADO]

| Risco | Probabilidade | Impacto | Mitigação |
|-------|--------------|---------|-----------|
| APIs governamentais instáveis | Alta | Médio | Cache local + circuit breaker + modo offline |
| Complexidade do Portal cresce além do MVP | Média | Alto | MVP com 3 programas primeiro; validação com usuários reais |
| Conflito de naming (ERP vs. Ecossistema) | Baixa | Baixo | Documentar claramente a hierarquia de módulos |
| Massa crítica para Rede Digna não atingida | Alta | Médio | Focar em ERP + Portal primeiro; Rede como "nice-to-have" |
| Teologia afeta adoção secular | Média | Alto | Manter produto laico na interface; teologia informa design internamente |
| Dependência de contadores sociais para escala | Média | Alto | Criar programa de capacitação + certificação CFC |
| Linguagem muito técnica nos tópicos de ajuda (RF-30) | Alta | Alto | Revisão por ITCPs/comunidade; teste de usabilidade com usuários reais |
| Conteúdo de ajuda desatualizado | Média | Médio | Processo de atualização via central.db, não hardcoded |

---

## 🚀 Próximos Passos Recomendados [ATUALIZADO]

### Imediatos (Próximo Trimestre)

1. **RF-27 (DAS MEI):** Cálculo automático, tabela versionada de salário mínimo — baixo esforço, alto valor percebido
2. **RF-30 (Ajuda Educativa):** Sistema de ajuda estruturada, seed de 10+ tópicos — habilita adoção por baixa escolaridade
3. **RF-19 (Perfil de Elegibilidade):** Campos complementares, preenchimento único — habilita o Portal
4. **RF-18 (Motor de Indicadores):** Coleta BCB/IBGE, cache local, interpretação — fornece contexto macroeconômico
5. **RF-20 (Portal MVP):** Match com 3 programas, checklist de documentos — entrega valor imediato

### Médio Prazo (6 meses)

1. **RF-14 a RF-17:** Adequação Estatal completa (EFD-Reinf, MAPA, Gov.br, CADSOL)
2. **RF-21 a RF-23:** Portal completo com monitoramento de DOU e alertas
3. **Validação de Campo:** 5-10 entidades reais em Niterói para prova de conceito

### Longo Prazo (1 ano+)

1. **RF-24 a RF-26:** Rede Digna com massa crítica de usuários
2. **v4:** Finanças Territoriais (moedas sociais, BCDs, estoque substantivo)
3. **Escala Nacional:** Expansão para além de Niterói/RJ

---

## 📊 Métricas de Qualidade do Projeto [ATUALIZADO]

| Métrica | Valor Atual | Meta | Status |
|---------|-------------|------|--------|
| Testes unitários passando | 311/311 (100%) | >95% | ✅ |
| Cobertura de handlers | ~87% | >90% | ⚠️ |
| Validação E2E por feature | 0% (novo) | 100% | ❌ Novo |
| Tópicos de ajuda criados | 0 (novo) | 10+ | ❌ Novo |
| Anti-Float violations | 0 | 0 | ✅ |
| Cache-Proof violations | 0 | 0 | ✅ |
| Soberania violations | 0 | 0 | ✅ |
| Tempo de resposta (local) | <50ms | <200ms | ✅ |

---

## 📝 Considerações Finais [NOVO - PDF v1.0]

O Ecossistema Digna parte de uma premissa simples e poderosa: **a informação que emancipa o empreendedor rico já existe** — está nas APIs do governo, nos editais publicados, nas linhas de crédito abertas. O que falta é uma camada de tradução, curadoria e entrega que chegue até o empreendedor pobre na linguagem e no momento certo.

O que torna este projeto distinto de iniciativas similares é a combinação de três elementos raramente encontrados juntos:
1. **Um ERP que já captura o perfil real do negócio**, eliminando formulários redundantes
2. **Um motor de dados que mantém tudo atualizado automaticamente**
3. **Uma rede de contadores e líderes comunitários como agentes de distribuição e capacitação**

O fundamento teológico não é acessório — é o que garante que as decisões de produto, mesmo sob pressão de prazo e orçamento, nunca percam de vista quem é o usuário e para que esse sistema existe.

---

**Status:** ✅ ATUALIZADO COM VISÃO DE ECOSSISTEMA (PDF v1.0) + RF-30 (Sessão 27/03/2026)  
**Próxima Ação:** Atualizar `06_roadmap/05_session_log.md` com log da sessão de expansão do ecossistema  
**Versão Anterior:** 1.0 (2026-03-13)  
**Versão Atual:** 2.0 (2026-03-27)
