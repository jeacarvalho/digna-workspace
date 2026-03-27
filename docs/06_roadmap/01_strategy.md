title: Estratégia de Release - Ecossistema Digna
status: implemented
version: 2.0
last_updated: 2026-03-27
---

# Estratégia de Release - Ecossistema Digna

> **Nota:** Este documento reflete a estratégia integrada do Ecossistema Digna (PDF v1.0), preservando toda a infraestrutura validada nas Sprints 1-16, incorporando os novos módulos como Fases 3 e 4, e incluindo o Sistema de Ajuda Educativa (RF-30) decidido na sessão de 27/03/2026.

---

## 📋 Contexto da Atualização (27/03/2026)

**Motivação:** O projeto Digna evoluiu de um ERP contábil para um **Ecossistema de 4 Módulos** conforme especificação PDF v1.0. Esta atualização documenta:

1. **Expansão da Visão de Produto:** De ERP único para ecossistema integrado (ERP + Motor + Portal + Rede)
2. **Novas Fases de Release:** Fases 3-6 adicionadas ao roadmap original
3. **RF-30 (Ajuda Educativa):** Sistema transversal de pedagogia social
4. **Preservação:** Todas as entregas validadas nas Sprints 1-16 mantidas

**Versão Anterior:** 1.0 (2026-03-13)  
**Nova Versão:** 2.0 (2026-03-27)

---

## 🎯 Visão Estratégica do Ecossistema

O **Ecossistema Digna** é uma plataforma tecnológica integrada voltada ao empreendedorismo popular no Brasil. Seu propósito central é reduzir as barreiras de acesso à informação normativa, ao crédito e à gestão financeira para microempreendedores, pequenos negócios, cooperativas e organizações da sociedade civil que atuam em contextos de vulnerabilidade socioeconômica.

### 1.1 Interseção de Conhecimentos

O projeto nasce da interseção entre quatro campos do conhecimento, cada um contribuindo com uma dimensão essencial:

| Área | Contribuição ao Projeto |
|------|------------------------|
| **Tecnologia da Informação** | Arquitetura de sistemas, APIs, automação, interfaces acessíveis e inteligência de dados |
| **Ciências Contábeis** | Gestão financeira, obrigações fiscais, escrituração, módulo de contador social e conformidade normativa |
| **Ciências Econômicas** | Indicadores macroeconômicos, análise de crédito, mapeamento de financiamentos e interpretação de política econômica |
| **Teologia Cristã** | Base ética e filosófica do design: princípios de dignidade humana, mordomia, koinonia e justiça restaurativa |

### 1.2 Problema Central

O empreendedor de baixa renda no Brasil enfrenta quatro barreiras simultâneas:
1. **Desconhecimento** das obrigações normativas do seu negócio
2. **Inacessibilidade** da linguagem burocrática dos programas de crédito
3. **Ausência de ferramentas** de gestão adequadas ao seu nível de escolaridade e realidade financeira
4. **Isolamento** — não sabe que existem redes de apoio, programas de financiamento e parceiros potenciais à sua volta

### 1.3 Solução Proposta

Uma plataforma que:
- **Esconde a complexidade normativa** por trás de interfaces simples
- **Conecta automaticamente** o empreendedor aos programas de crédito para os quais ele já é elegível
- **Cria uma rede de colaboração** entre entidades
- **Tudo sustentado por dados econômicos atualizados em tempo real**

### 1.4 Princípio Arquitetural Central

> **"Nenhum usuário precisa preencher o mesmo dado duas vezes."**

O que o ERP já sabe sobre a entidade é automaticamente aproveitado pelo Portal, pelo Motor e pela Rede. Isso reduz a fricção e elimina a principal barreira de uso: o excesso de formulários.

---

## 🗺️ Versões e Mapa de Valor [ATUALIZADO]

A estratégia do Digna respeita o tempo social dos Empreendimentos de Economia Solidária (EES). Cada fase entrega não apenas infraestrutura tecnológica, mas também ferramentas pedagógicas de emancipação e pontes de conformidade institucional.

| Versão | Fase | Descrição Sociotécnica | Status |
|--------|------|------------------------|--------|
| **v0** | Demonstração | Operação Básica, PDV Pedagógico e Validação em Campo | ✅ COMPLETE (Sprint 16) |
| **v1** | Integração | Transição Gradual Institucional, Painel do Contador Social e Exportação SPED | ✅ COMPLETE (Sprint 12) |
| **v2** | Ecossistema de Crédito | Motor de Indicadores + Portal de Oportunidades (MVP) | 🟡 EM DESENVOLVIMENTO (Fase 3) |
| **v3** | Intercooperação | Rede Digna (Marketplace Solidário) + Escala Nacional | 🔵 PLANNED (Fase 4) |
| **v4** | Finanças Territoriais | Moedas Sociais + Bancos Comunitários + Estoque Substantivo | 🔵 PLANNED (Futuro) |

---

## 📦 v0 - Demonstração Operacional e Validação Cultural

**Status:** ✅ COMPLETE  
**Sprints:** 01-16  
**Objetivo:** Provar conceito técnico, arquitetura offline-first e, criticamente, a adequação sociotécnica e cultural da interface junto ao público-alvo.

### Entregas Concluídas

- [x] Lifecycle Manager (SQLite por tenant para garantia de soberania)
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

### Critérios de Release Atendidos

- ✅ 149/149 testes técnicos passando
- ✅ Operações básicas fluindo perfeitamente offline com sync posterior
- ✅ Aprovação de usabilidade e linguagem por grupos focais com baixa literacia digital
- ✅ Sistema 100% operacional (Sprint 16)

---

## 📦 v1 - Integração Institucional e Aliança Contábil

**Status:** ✅ COMPLETE  
**Sprints:** 09-16  
**Objetivo:** Oferecer os benefícios da conformidade legal (Compliance) e a ponte tecnológica com a classe contábil (CFC/CRCs), sem burocratizar ou afastar o usuário informal produtor.

### Entregas Concluídas

- [x] Módulos educativos embutidos (ex: auxílio na formação de preço no PDV)
- [x] Dossiê CADSOL automático (Apenas habilitado quando o grupo atingir maturidade política para formalização)
- [x] Geração de documentos oficiais (Atas em Markdown com hash SHA256)
- [x] Integração Gov.br (Mock → Real OAuth2 em desenvolvimento)
- [x] **Painel do Contador Social (Accountant Dashboard):** Interface Multi-tenant em modo de leitura (Read-Only) para auditores voluntários ✅ **Sprint 12**
- [x] **Exportação Fiscal (SPED):** Motor de tradução das partidas dobradas (geradas silenciosamente no Lume) para os leiautes contábeis e fiscais exigidos pela Receita Federal ✅ **Sprint 12**
- [x] **Gestão de Vínculo Contábil (RF-12):** Sistema de delegação temporal com Exit Power ✅ **95% Complete**
- [x] **Visão Analítica do Contador (RF-13):** Dashboard consultivo com filtro temporal ✅ **95% Complete**

### Entregas Pendentes (Adequação Estatal)

- [ ] **RF-14: Blindagem Tributária (EFD-Reinf e ECF)** - Módulo `tax_compliance` para mensageria de retenções via Web Service
- [ ] **RF-15: Integração Real Gov.br** - Assinatura Eletrônica Qualificada (Lei nº 14.063/2020) + anonimização de votos (IN DREI nº 79/2020)
- [ ] **RF-16: Inclusão Sanitária (MAPA)** - Gerador automatizado do Memorial Técnico Sanitário (MTSE)
- [ ] **RF-17: Integração CADSOL/SINAES** - Consumo nativo de Web Services do MTE (Decreto nº 12.784/2025)

### Critérios de Release

- ✅ Painel do Contador funcional com acesso Read-Only
- ✅ Exportação SPED validada com hash de integridade
- ✅ Vínculo contábil com controle temporal (start_date, end_date)
- 🟡 Integrações estatais reais em desenvolvimento (mocks → produção)

---

## 📦 v2 - Ecossistema de Crédito e Indicadores [NOVO - PDF v1.0]

**Status:** 🟡 EM DESENVOLVIMENTO (Fase 3)  
**Sprints:** 17-20 (Planejado)  
**Objetivo:** Conectar automaticamente o empreendedor aos programas de crédito para os quais ele já é elegível, sustentado por dados econômicos atualizados em tempo real.

### Módulo 2: Motor de Indicadores (RF-18)

| Sub-requisito | Descrição | Fonte de Dados | Prioridade |
|--------------|-----------|---------------|------------|
| RF-18.1 | Coleta de SELIC, IPCA, CDI | BCB SGS API | Alta |
| RF-18.2 | Câmbio oficial (USD/BRL, EUR/BRL) | BCB PTAX API | Alta |
| RF-18.3 | Expectativas de mercado (Focus) | BCB Focus API | Média |
| RF-18.4 | Indicadores sociais (PNAD, desocupação) | IBGE SIDRA | Média |
| RF-18.5 | Cache local e interpretação contextual | Arquitetura interna | Alta |

**Arquitetura Proposta:**
```
modules/
└── indicators_engine/
    ├── internal/
    │   ├── collector/          # Scheduler de coleta (cron diário)
    │   ├── cache/              # Tabela indicators{} com TTL
    │   ├── interpreter/        # Regras de negócio contextualizadas
    │   └── repository/         # SQLite local (central.db)
    ├── pkg/
    │   └── indicators/         # API pública para outros módulos
    └── cmd/
        └── collector/          # Binário independente para coleta
```

### Módulo 3: Portal de Oportunidades (RF-19 a RF-23)

**Programas Prioritários para MVP:**
| Programa | Valor Máximo | Público-Alvo | Diferencial |
|----------|-------------|--------------|-------------|
| **Acredita no Primeiro Passo** | ~R$ 6 mil | CadÚnico, 70% mulheres | Sem garantia, juros subsidiados |
| **Pronampe** | Até 30% da receita bruta | MEI, ME, EPP | Prazo até 84 meses, carência 24 meses |
| **Niterói Empreendedora** | Até R$ 200 mil | CNPJ ativo 12+ meses em Niterói | **JUROS ZERO**, foco em mulheres e jovens |

**Funcionalidades do MVP:**
- [ ] **Match Automático:** Cruzamento do perfil do ERP com requisitos de cada programa
- [ ] **Ranqueamento por Vantagem:** Ordenação por custo efetivo de capital (usando taxas do Motor)
- [ ] **Checklist de Documentos:** Para cada linha elegível, marca o que já está no Digna vs. o que falta
- [ ] **Alertas de Prazo:** Notificações quando editais com alto match têm prazo próximo

### Sistema Transversal: Ajuda Educativa (RF-30) [NOVO - Sessão 27/03/2026]

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

### Critérios de Release v2

- [ ] Motor de Indicadores coletando dados de BCB/IBGE diariamente
- [ ] Perfil de Elegibilidade (RF-19) com campos complementares implementados
- [ ] Portal MVP com 3 programas funcionando (match automático)
- [ ] Sistema de Ajuda Educativa (RF-30) com 10+ tópicos seed
- [ ] Validação com 5-10 entidades reais de Niterói (prova de conceito)

---

## 📦 v3 - Intercooperação Nacional (Rede Digna) [NOVO - PDF v1.0]

**Status:** 🔵 PLANNED (Fase 4)  
**Sprints:** 21-24 (Planejado)  
**Objetivo:** Conectar EES isolados em uma rede nacional de apoio e viabilidade econômica, materializando o 6º Princípio do Cooperativismo.

### Módulo 4: Rede Digna (RF-24 a RF-26)

**Funcionalidades:**
- [ ] **Perfil Público da Entidade (RF-24):** Missão, produtos, serviços e capacidades disponíveis para a rede
- [ ] **Mural de Necessidades (RF-25):** Entidades publicam demandas de compra (insumos, serviços, equipamentos) visíveis para a rede
- [ ] **Match de Oportunidades (RF-26):** Sistema sugere conexões entre quem precisa e quem oferece, priorizando proximidade geográfica e afinidade de setor
- [ ] **Histórico de Transações Solidárias:** Registro que pode ser apresentado como evidência de atividade econômica em candidaturas a financiamento

### Fundamento Teológico da Rede

O conceito de **Koinonía** — comunhão e partilha mútua na tradição cristã primitiva — é a metáfora central da Rede Digna. Não é apenas um marketplace: é uma ecclesia econômica onde entidades que compartilham valores se apoiam mutuamente, reduzindo dependência de intermediários externos e fortalecendo a economia local.

**Nota sobre Laicidade:** A Teologia informa o design internamente, mas o produto permanece acessível independentemente de crença. O canal de distribuição via igrejas e comunidades de fé é estratégico, não doutrinário.

### Critérios de Release v3

- [ ] 100+ entidades com perfil público publicado
- [ ] 500+ conexões B2B realizadas na rede
- [ ] Match geográfico e setorial funcionando
- [ ] Histórico de transações solidárias auditável
- [ ] Massa crítica mínima atingida (validação de rede)

---

## 📦 v4 - Finanças Territoriais e Solidárias [FUTURO]

**Status:** 🔵 PLANNED (Futuro)  
**Objetivo:** Suporte a múltiplas unidades de valor, reconhecendo que a riqueza na economia solidária vai além da moeda oficial (Real R$).

### Entregas Planejadas

- [ ] Integração tecnológica com Bancos Comunitários de Desenvolvimento (BCDs)
- [ ] Gestão e transação de Moedas Sociais locais
- [ ] Estoque substantivo (Controle e troca de sementes, animais, horas-trabalho)
- [ ] Rateio automático de sobras com painéis de aprovação visual para assembleias
- [ ] Score de Crédito Social baseado em histórico de trabalho (não apenas dinheiro)

---

## 📊 Estratégia de Implementação [ATUALIZADO - PDF v1.0]

Dado que o digna ERP já está em desenvolvimento, a estratégia proposta é incremental: priorizar os módulos que ampliam o valor do ERP existente com menor esforço, reservando os mais complexos para fases posteriores.

| Fase | Escopo | Prioridade | Status |
|------|--------|------------|--------|
| **Fase 1 — Motor de Indicadores** | Implementar scheduler de coleta das APIs do BCB e IBGE. Criar tabela de cache local. Exibir indicadores contextualizados no ERP. | Alta — baixo esforço, alto impacto imediato | 📋 Backlog |
| **Fase 2 — Perfil de Elegibilidade** | Adicionar campos complementares ao cadastro do ERP (CadÚnico, gênero do sócio, inadimplência, finalidade). Consulta automática de certidões. | Alta — habilita o Portal | 📋 Backlog |
| **Fase 3 — Portal (MVP)** | Match com os 3 programas de maior alcance: Acredita no Primeiro Passo, Pronampe e Niterói Empreendedora. Checklist de documentos básico. | Alta — gera valor imediato e validável | 📋 Backlog |
| **Fase 4 — Contador Social** | Módulo de adoção de entidades por contadores. Dashboard do contador com visão consolidada. Submissão de candidaturas. | Média — estratégica para escala | 🟡 95% Complete |
| **Fase 5 — Portal completo** | Adicionar todos os programas mapeados. Monitoramento de DOU. Alertas automáticos. Ranqueamento por custo efetivo. | Média — expansão do MVP | 📋 Backlog |
| **Fase 6 — Rede Digna** | Perfil público, mural de necessidades, match de oportunidades entre entidades. | Baixa — depende de massa crítica de usuários | 📋 Backlog |
| **Fase Transversal — Ajuda Educativa (RF-30)** | Sistema de ajuda contextual em todos os módulos. Botão "?" em campos técnicos. Linguagem para 5ª série. | Alta — habilita adoção por baixa escolaridade | 📋 Backlog |

### Recomendação para Validação Inicial

Antes de construir o Portal completo, validar o motor de match com 5 a 10 entidades reais de Niterói.

O **Programa Niterói Empreendedora** (juros zero, edital ativo) é o caso de uso ideal para a prova de conceito.

Um único match bem-sucedido — entidade que não sabia que era elegível e conseguiu o crédito — é o melhor argumento para o próximo edital de captação de recursos.

---

## 🎯 Critérios de Sucesso por Versão [ATUALIZADO]

| Versão | Métrica de Sucesso | Alvo | Status |
|--------|-------------------|------|--------|
| **v0** | Entidades operando com contabilidade invisível | 100+ entidades ativas | ✅ 149/149 testes |
| **v1** | Entidades com conformidade estatal automatizada | 80% das formalizadas | ✅ Sprint 12 |
| **v2** | Entidades descobrindo elegibilidade via Portal | 50+ matches bem-sucedidos | 🟡 Em desenvolvimento |
| **v3** | Transações B2B realizadas na Rede | 500+ conexões/ano | 🔵 Planned |
| **v2 (RF-30)** | Redução de abandono em formulários | 30% redução | 📋 Backlog |
| **v2 (RF-30)** | Tópicos de ajuda visualizados/mês | 500+ visualizações | 📋 Backlog |

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

## 📝 Considerações Finais [NOVO - PDF v1.0]

O Ecossistema Digna parte de uma premissa simples e poderosa: **a informação que emancipa o empreendedor rico já existe** — está nas APIs do governo, nos editais publicados, nas linhas de crédito abertas. O que falta é uma camada de tradução, curadoria e entrega que chegue até o empreendedor pobre na linguagem e no momento certo.

O que torna este projeto distinto de iniciativas similares é a combinação de três elementos raramente encontrados juntos:
1. **Um ERP que já captura o perfil real do negócio**, eliminando formulários redundantes
2. **Um motor de dados que mantém tudo atualizado automaticamente**
3. **Uma rede de contadores e líderes comunitários como agentes de distribuição e capacitação**

O fundamento teológico não é acessório — é o que garante que as decisões de produto, mesmo sob pressão de prazo e orçamento, nunca percam de vista quem é o usuário e para que esse sistema existe.

---

**Status:** ✅ ATUALIZADO COM VISÃO DE ECOSSISTEMA (PDF v1.0) + RF-30 (Sessão 27/03/2026)  
**Próxima Ação:** Atualizar `06_roadmap/04_status.md` com status consolidado do ecossistema  
**Versão Anterior:** 1.0 (2026-03-13)  
**Versão Atual:** 2.0 (2026-03-27)
