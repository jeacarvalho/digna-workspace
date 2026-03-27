title: Product Backlog - Ecossistema Digna
status: implemented
version: 1.4
last_updated: 2026-03-27
---

# Product Backlog - Ecossistema Digna

> **Nota:** Este backlog consolida todas as funcionalidades implementadas (Sprints 1-16) e as futuras, alinhadas às Fases do Roadmap, aos Requisitos Funcionais (RF) do BRD, à Especificação v1.0 (PDF) e às decisões de design tomadas na sessão de 27/03/2026 (incluindo RF-30 - Sistema de Ajuda Educativa).

---

## Funcionalidades Concluídas

### Core Contábil e Fundação (Sprints 01 a 06)
* [x] PDV Operacional com Partidas Dobradas Invisíveis (RF-01)
* [x] Registro de Trabalho / Ponto Social ITG 2002 (RF-02)
* [x] Motor de Reservas Obrigatórias (RF-03)
* [x] Dossiê de Formalização (RF-04)
* [x] Sincronização Offline-First via Delta Tracking (RF-05)
* [x] Intercooperação B2B - Base de dados (RF-06)
* [x] Gestão de Caixa conectada ao Motor Lume (RF-09)

### Gestão de Membros (Sprint 10)
* [x] Cadastro de membros
* [x] Roles (COORDINATOR, MEMBER, ADVISOR) com hierarquia de permissões
* [x] Status (ACTIVE, INACTIVE) e trava contra desativação do último coordenador
* [x] Skills/habilidades associadas ao trabalho voluntário

### Formalização (Sprint 11)
* [x] Transição automática de status DREAM → FORMALIZED
* [x] Simulador de formalização para análise de impacto
* [x] Funcionalidade CheckFormalizationCriteria()
* [x] Funcionalidade AutoTransitionIfReady() (gatilho após 3 decisões)

### Integrações Governamentais (Mocks)
* [x] Receita Federal (CNPJ, DARF) - Mock
* [x] MTE (RAIS, CAT, eSocial) - Mock
* [x] MDS (CadÚnico, Relatório Social) - Mock
* [x] IBGE (Pesquisas, PAM, CNAE) - Mock
* [x] SEFAZ (NFe, NFS-e, Manifesto) - Mock
* [x] BNDES (Linhas de Crédito, Simulação) - Mock
* [x] SEBRAE (Cursos, Consultoria) - Mock
* [x] Providentia (Sync, Marketplace) - Mock

### Finanças Solidárias e Suprimentos (Sprints 13 e 14)
* [x] Gestão de compras e fornecedores (RF-07)
* [x] Controle de estoque por categorização obrigatória (INSUMO, PRODUTO, MERCADORIA) com baixa via PDV (RF-08)
* [x] Gestão Orçamentária e Planejamento Financeiro com alertas visuais SAFE/WARNING/EXCEEDED (RF-10)

### Estabilização, UI e Rateio (Sprints 15 e 16)
* [x] Surplus Calculator (CalculateSocialSurplus(), deduções automáticas de 10% Reserva Legal + 5% FATES)
* [x] Identidade Visual Global "Soberania e Suor" aplicada (RNF-07)
* [x] Arquitetura de renderização de Templates Cache-Proof (*_simple.html)
* [x] Integração End-to-End funcional (PDV → Estoque → Caixa)

### Aliança Contábil (Sprint 12)
* [x] Painel do Contador Social (Accountant Dashboard) - Interface Multi-tenant Read-Only (RF-11)
* [x] Exportação Fiscal (SPED/CSV) com hash SHA256 (RF-11)
* [x] Gestão de Vínculo Contábil e Delegação Temporal (RF-12) - 95% completo
* [x] Visão Analítica do Contador Social com filtro temporal (RF-13) - 95% completo

---

## Funcionalidades Futuras

### Alta Prioridade (Fase 2 - Adequação Estatal e Conformidade Digital)

* [ ] **RF-14: Blindagem Tributária (EFD-Reinf e ECF)** [NOVO - PDF v1.0]
  * Módulo `tax_compliance` para mensageria de retenções via Web Service
  * Expurgo automático de receitas de Atos Cooperativos no Bloco M da ECF (Lei 5.764/71 e LC 214/2025)
  * Geração e transmissão de XMLs (série R-2000/R-4000) para Web Services da EFD-Reinf
  * Alimentação automática da DCTFWeb

* [ ] **RF-15: Integração Real Gov.br e Governança Digital** [ATUALIZADO - PDF v1.0]
  * Substituição do Mock de login unificado pelo fluxo real da Cidadania Digital (OAuth2)
  * Assinatura Eletrônica Qualificada (Lei nº 14.063/2020) para membros da mesa nas Atas de Assembleia
  * Algoritmo de anonimização sistêmica para escrutínio secreto em votações (IN DREI nº 79/2020)
  * Integração com ICP-Brasil para certificados digitais

* [ ] **RF-16: Inclusão Sanitária (MAPA)** [NOVO - PDF v1.0]
  * Módulo `sanitary_compliance` para geração automatizada do Memorial Técnico Sanitário de Estabelecimento (MTSE)
  * Conformidade com Portaria MAPA nº 393/2021 para agroindústrias
  * Parametrização de fluxogramas de maquinário, capacidade diária e potabilidade de água
  * Exportação para peticionamento no Sistema Eletrônico de Informação (SEI)

* [ ] **RF-17: Integração CADSOL/SINAES Automático** [NOVO - PDF v1.0]
  * Consumo nativo de Web Services do MTE (Decreto nº 12.784/2025)
  * Matrícula automática de entidades FORMALIZED no Cadastro Nacional de Economia Solidária
  * Substituição dos mocks atuais por integração real

* [ ] **RF-27: Cálculo Automático de DAS MEI** [NOVO - PDF v1.0, Seção 4.1]
  * Tabela de salário mínimo versionada por ano (atualização via decreto presidencial)
  * Cálculo automático: 5% do salário mínimo vigente + ICMS/ISS fixos
  * Alertas de vencimento (dia 20 de cada mês)
  * Histórico de pagamentos registrados
  * Geração de guia para pagamento (link ou exportação)
  * **Critério de Aceite:** Nenhum usuário precisa calcular manualmente o DAS
  * **Status:** 📋 Backlog (Prioridade Alta - Baixo esforço, alto valor percebido)

---

### Alta Prioridade (Fase 3 - Ecossistema de Crédito e Indicadores) [NOVO - PDF v1.0]

* [ ] **RF-18: Motor de Indicadores Econômico-Financeiros** [NOVO - PDF v1.0, Seção 5]
  * **Coleta de Dados:**
    * [ ] SELIC, IPCA, CDI via BCB SGS API (RF-18.1)
    * [ ] Câmbio oficial (USD/BRL, EUR/BRL) via BCB PTAX API (RF-18.2)
    * [ ] Expectativas de mercado (Focus) via BCB Focus API (RF-18.3)
    * [ ] Indicadores sociais (PNAD, desocupação) via IBGE SIDRA (RF-18.4)
  * **Arquitetura:**
    * [ ] Scheduler diário de coleta
    * [ ] Cache local com tabela `indicators{}` e TTL
    * [ ] Camada de interpretação contextualizada
    * [ ] API interna para consumo pelos demais módulos
  * **Critério de Aceite:** Nenhum usuário preenche dados econômicos manualmente
  * **Status:** 📋 Backlog (Prioridade Alta)

* [ ] **RF-19: Perfil de Elegibilidade (Campos Complementares)** [NOVO - PDF v1.0, Seção 9]
  * **Dados já disponíveis no ERP (capturados automaticamente):**
    * [ ] CNPJ/CPF, CNAE, Município, UF
    * [ ] Faturamento anual (int64), Regime Tributário
    * [ ] Data de abertura, Situação fiscal
    * [ ] Certidões negativas (integração com portais de certidões)
  * **Novos Campos (preenchimento único):**
    * [ ] Inscrito no CadÚnico (bool) - Habilita programas sociais
    * [ ] Sócio Mulher (bool) - Prioridade em linhas com foco de gênero
    * [ ] Inadimplência Ativa (bool) - Direciona ao Desenrola antes de crédito
    * [ ] Finalidade do Crédito (enum) - CAPITAL_GIRO, EQUIPAMENTO, REFORMA, OUTRO
    * [ ] Valor Necessário (int64) - Anti-Float compliance
    * [ ] Tipo de Entidade (enum) - MEI, ME, EPP, Cooperativa, OSC, OSCIP, PF
    * [ ] Contabilidade Formal (bool) - Requisito de alguns programas
  * **Critério de Aceite:** Nenhum dado digitado duas vezes
  * **Status:** 📋 Backlog (Prioridade Alta - Habilita o Portal)

* [ ] **RF-20: Portal de Oportunidades (MVP - 3 Programas)** [NOVO - PDF v1.0, Seção 6]
  * **Programas Prioritários para MVP:**
    * [ ] **Acredita no Primeiro Passo** (~R$ 6 mil, CadÚnico, 70% mulheres)
    * [ ] **Pronampe** (Até 30% da receita bruta, MEI/ME/EPP)
    * [ ] **Niterói Empreendedora** (Até R$ 200 mil, JUROS ZERO, foco em mulheres e jovens)
  * **Funcionalidades:**
    * [ ] Match Automático: Cruzamento perfil × requisitos de cada programa
    * [ ] Ranqueamento por Vantagem: Ordenação por custo efetivo de capital (usando taxas do Motor)
    * [ ] Checklist de Documentos: Marca o que já está no Digna vs. o que falta
    * [ ] Alertas de Prazo: Notificações quando editais com alto match têm prazo próximo
  * **Critério de Aceite:** Entidade descobre elegibilidade sem preencher formulários
  * **Status:** 📋 Backlog (Prioridade Alta)

* [ ] **RF-21: Checklist + Alertas de Documentos** [NOVO - PDF v1.0]
  * [ ] Integração com portais de certidões (PGFN, Estadual, Municipal) para verificação automática
  * [ ] Scraping do Diário Oficial da União (DOU) para captura de novos editais
  * [ ] Sistema de notificações push/email para prazos críticos
  * [ ] Dashboard de documentos pendentes por programa
  * **Status:** 📋 Backlog (Prioridade Média)

* [ ] **RF-22: Monitoramento de DOU** [NOVO - PDF v1.0]
  * [ ] Scraper de Diário Oficial para novos editais
  * [ ] Classificação automática por relevância (match com perfil da entidade)
  * [ ] Notificação proativa para Contador Social e entidade
  * **Status:** 📋 Backlog (Prioridade Média)

* [ ] **RF-23: Dashboard de Elegibilidade** [NOVO - PDF v1.0]
  * [ ] Visão consolidada de programas elegíveis
  * [ ] Status de cada candidatura (Elegível, Documentação Pendente, Enviada, Aprovada)
  * [ ] Histórico de candidaturas submetidas
  * **Status:** 📋 Backlog (Prioridade Média)

---

### Alta Prioridade (Fase 5 - Ajuda e Pedagogia) [NOVO - Decisão da Sessão 27/03/2026]

* [ ] **RF-30: Sistema de Ajuda Educativa Estruturada** [NOVO - Decisão de Design]
  * **Descrição:** Sistema de ajuda contextual que traduz conceitos técnicos (CadÚnico, inadimplência, CNAE, etc.) em linguagem popular, com linkagem entre elementos de UI e registros de ajuda no banco.
  * **Funcionalidades:**
    * [ ] Entrada de menu "Ajuda" acessível em todas as páginas
    * [ ] Busca e índice de tópicos categorizados (CRÉDITO, TRIBUTÁRIO, GOVERNANÇA, GERAL)
    * [ ] Explicação em linguagem popular + legislação relacionada + próximo passo acionável
    * [ ] Linkagem automática: botão "?" ao lado de campos técnicos abre explicação contextual
    * [ ] Tópicos com estrutura: Título, Resumo, Explicação, Por que perguntamos, Legislação, Próximo passo, Link oficial
  * **Tópicos Educativos Obrigatórios (Seed Inicial):**
    * [ ] CadÚnico - "É o cadastro do governo para programas sociais"
    * [ ] Inadimplência - "É quando há dívidas não pagas registradas"
    * [ ] CNAE - "É o código que diz qual é a atividade do seu negócio"
    * [ ] DAS MEI - "É o boleto mensal que o MEI paga"
    * [ ] Reserva Legal - "É uma parte do lucro que a lei manda guardar"
    * [ ] FATES - "É um fundo para ajudar outros grupos a se organizarem"
  * **Critério de Aceite Pedagógico:**
    * [ ] Usuário com 5ª série consegue entender a explicação sem ajuda externa
    * [ ] Tooltip carrega em < 500ms via HTMX
    * [ ] Conteúdo não usa jargões técnicos ("cadastramento", "regularização fiscal", etc.)
    * [ ] Sempre inclui "próximo passo" acionável (ex: "procure o CRAS")
  * **Status:** 📋 Backlog (Prioridade Alta - Habilita adoção por baixa escolaridade, transversal a todos os módulos)

---

### Baixa Prioridade (Fase 4 - Intercooperação e Escala Nacional) [NOVO - PDF v1.0]

* [ ] **RF-24: Perfil Público da Entidade** [NOVO - PDF v1.0, Seção 7]
  * **Descrição:** Permitir que entidades publiquem informações visíveis para a rede, sem expor dados sensíveis.
  * **Campos do Perfil Público:**
    * [ ] EntityID (Hash anonimizado)
    * [ ] Nome Fantasia, Missão (texto livre)
    * [ ] Produtos (lista de categorias)
    * [ ] Serviços (lista de capacidades)
    * [ ] Município, UF, Contato Público
    * [ ] Foto/Logo (URL do asset)
  * **Critério de Aceite:** Dados sensíveis nunca expostos publicamente
  * **Status:** 📋 Backlog (Prioridade Baixa - Depende de massa crítica)

* [ ] **RF-25: Mural de Necessidades** [NOVO - PDF v1.0, Seção 7]
  * **Descrição:** Entidades publicam demandas de compra visíveis para a rede.
  * **Estrutura de Postagem:**
    * [ ] ID, PublisherID (Hash anonimizado)
    * [ ] Categoria (INSUMO, EQUIPAMENTO, SERVICO, OUTRO)
    * [ ] Descrição, Quantidade, Prazo Desejado
    * [ ] Município (para matching geográfico)
    * [ ] Status (ABERTO, EM_NEGOCIACAO, CONCLUIDO)
  * **Critério de Aceite:** Entidades publicam demandas visíveis para a rede
  * **Status:** 📋 Backlog (Prioridade Baixa)

* [ ] **RF-26: Match de Oportunidades B2B** [NOVO - PDF v1.0, Seção 7]
  * **Descrição:** Algoritmo que sugere conexões entre quem precisa e quem oferece.
  * **Critérios de Matching:**
    * [ ] Geográfico: Priorizar proximidade (mesmo município/UF)
    * [ ] Setorial: Afinidade de CNAE/categoria
    * [ ] Temporal: Prazos compatíveis
    * [ ] Reputação: Histórico de transações na rede (futuro)
  * **Critério de Aceite:** Sistema sugere conexões entre quem precisa e quem oferece
  * **Status:** 📋 Backlog (Prioridade Baixa)

---

### Média Prioridade (Fase 3 - Finanças Solidárias Avançadas)

* [ ] **Múltiplas Moedas Sociais**
  * [ ] Expansão do Ledger para registrar e transacionar moedas complementares de Bancos Comunitários de Desenvolvimento (BCDs)
  * [ ] Conversão automática entre moeda social e Real (R$)
  * [ ] Integração tecnológica com BCDs existentes

* [ ] **Estoque Substantivo**
  * [ ] Suporte à contabilidade não-monetária
  * [ ] Gestão de "Fundos Rotativos Solidários" (troca e controle genético de sementes, animais para repasse)
  * [ ] Valoração de bens substantivos (sementes, animais, horas-trabalho)

* [ ] **Rateio de Sobras na Interface (UI)**
  * [ ] Painel visual para que a Assembleia aprove a divisão justa do excedente
  * [ ] Transparência algorítmica baseada nas horas trabalhadas
  * [ ] Exportação de relatório de rateio para ata de assembleia

---

### Baixa Prioridade (Fase 4 - Intercooperação e Escala Nacional)

* [ ] **Integração Contábil Fiscal Definitiva**
  * [ ] Conexão direta por API (sem arquivos intermediários) com softwares comerciais de contabilidade
  * [ ] Via Contador Social, preservando isolamento de dados

* [ ] **Score de Crédito Social**
  * [ ] Motor que calcula a reputação da entidade baseada no trabalho e autogestão (não apenas em dinheiro)
  * [ ] Histórico de transações solidárias como evidência de atividade econômica
  * [ ] Apresentação em candidaturas a financiamento

* [ ] **API Pública Restrita (OpenAPI/Swagger)**
  * [ ] Documentação técnica e geração de endpoints para ecossistemas parceiros (Serpro, Governos Estaduais)
  * [ ] Rate limiting e autenticação OAuth2 para APIs públicas

* [ ] **App Mobile Nativo / Relatórios PDF Avulsos**
  * [ ] Versão PWA aprimorada ou aplicativo nativo
  * [ ] Geração de relatórios PDF sob demanda

---

## Matriz de Priorização (Impacto × Esforço) [ATUALIZADO - Sessão 27/03/2026]

```
                    ALTO IMPACTO
                         │
         ┌───────────────┼───────────────┐
         │  RF-19        │  RF-18        │
         │  (Perfil)     │  (Indicadores)│
         │               │               │
ESFORÇO  │   🎯 RF-30    │               │  ESFORÇO
BAIXO    │   (Ajuda)     │               │  ALTO
         │   RF-27       │               │
         │   (DAS MEI)   │               │
         ├───────────────┼───────────────┤
         │               │  RF-20        │
         │               │  (Portal)     │
         │               │               │
         └───────────────┼───────────────┘
                         │
                    BAIXO IMPACTO
```

**Recomendação de Sequência de Implementação:**
1. **RF-27 (DAS MEI)** — Baixo esforço, alto valor percebido, prepara terreno para Portal
2. **RF-30 (Sistema de Ajuda)** — Baixo esforço, habilita adoção por baixa escolaridade, transversal a todos os módulos
3. **RF-19 (Perfil de Elegibilidade)** — Habilita o match automático do Portal
4. **RF-18 (Motor de Indicadores)** — Fornece contexto macroeconômico para decisões
5. **RF-20 (Portal MVP)** — Entrega valor imediato com 3 programas validados

---

## Dependências entre Requisitos [ATUALIZADO]

```
RF-18 (Motor) ──────┐
                    ├──▶ RF-20 (Portal)
RF-19 (Perfil) ─────┘
                         │
                         ▼
RF-27 (DAS MEI) ────────▶ RF-20 (Portal)
                         │
                         ▼
RF-21-23 (Alertas) ◀────┘

RF-24 (Perfil Público) ─┐
                        ├──▶ RF-26 (Match B2B)
RF-25 (Mural) ──────────┘

RF-30 (Ajuda) ──────────▶ TODOS OS MÓDULOS (transversal)
```

**Notas Arquiteturais Críticas:**
1. **Soberania do Dado:** Cada novo módulo deve preservar o isolamento físico dos bancos SQLite por entidade. Nenhum módulo pode violar este princípio.
2. **Anti-Float:** Todos os cálculos financeiros e de tempo devem usar `int64`. Esta regra é transversal a todas as fases.
3. **Cache-Proof Templates:** A interface web deve continuar carregando templates do disco via `ParseFiles()` no handler, garantindo atualizações imediatas.
4. **Nenhum Dado Digitado Duas Vezes:** Novos módulos devem **consumir** dados do ERP, nunca exigir reentrada. Este é o princípio central do PDF v1.0.
5. **Pedagogia Contextual (RF-30):** Todo campo técnico deve ter explicação acessível via botão "?" linkado ao help_engine.

---

## Backlog de Modularização [MANTIDO]

**Contexto:** Algumas funcionalidades foram implementadas de forma distribuída entre múltiplos módulos, violando o princípio SRP.

| Módulo | Status Atual | Prioridade | Esforço Estimado |
|--------|-------------|------------|------------------|
| member_management | ⚠️ Espalhado | **ALTA** | 2-3 dias |
| reporting | ⚠️ Básico | MÉDIA | 2-3 dias |
| sync_engine | ⚠️ Isolado | MÉDIA | 2-3 dias |

**Ver `docs/NEXT_STEPS.md` para detalhes completos do backlog de modularização.**

---

## Riscos do Backlog [ATUALIZADO - Sessão 27/03/2026]

| Risco | Probabilidade | Impacto | Mitigação |
|-------|--------------|---------|-----------|
| APIs governamentais instáveis | Alta | Médio | Cache local + circuit breaker + modo offline |
| Complexidade do Portal cresce além do MVP | Média | Alto | MVP com 3 programas primeiro; validação com usuários reais |
| Conflito de naming (ERP vs. Ecossistema) | Baixa | Baixo | Documentar claramente a hierarquia de módulos |
| Massa crítica para Rede Digna não atingida | Alta | Médio | Focar em ERP + Portal primeiro; Rede como "nice-to-have" |
| Teologia afeta adoção secular | Média | Alto | Manter produto laico na interface; teologia informa design internamente |
| Dependência de contadores sociais para escala | Média | Alto | Criar programa de capacitação + certificação CFC |
| **Linguagem muito técnica nos tópicos de ajuda (RF-30)** | Alta | Alto | Revisão por ITCPs/comunidade; teste de usabilidade com usuários reais |
| **Conteúdo de ajuda desatualizado** | Média | Médio | Processo de atualização via central.db, não hardcoded |

---

## Critérios de Aceite Gerais para Novos Requisitos [NOVO - Sessão 27/03/2026]

| Critério | Descrição | Validação |
|----------|-----------|-----------|
| **Anti-Float** | Zero `float` para valores financeiros/tempo | `grep -r "float[0-9]*" modules/` retorna apenas logs/comentários |
| **Soberania** | Dados isolados por entidade (`entity_id`) | Nenhum JOIN entre bancos diferentes |
| **Cache-Proof** | Templates `*_simple.html` carregados no handler | `ParseFiles()` no handler, não variáveis globais |
| **Reuso de Dados** | Nenhum dado digitado duas vezes | Perfil do ERP alimenta Portal/Rede automaticamente |
| **Testes** | Cobertura >90% para handlers, >80% para services | `go test ./... -cover` |
| **E2E** | Validação com `validate_e2e.sh --basic --headless` | 7 passos padrão Digna passam |
| **Documentação** | Aprendizados em `docs/learnings/` | `conclude_task.sh` executado |
| **Pedagogia (RF-30)** | Linguagem acessível (5ª série) | Teste de usabilidade com usuários reais |

---

**Status:** ✅ ATUALIZADO COM VISÃO DE ECOSSISTEMA (PDF v1.0) + RF-30 (Decisão de Design 27/03/2026)  
**Próxima Ação:** Gerar prompts iniciais para RF-27 (DAS MEI), RF-30 (Ajuda), RF-19 (Perfil)  
**Versão Anterior:** 1.3 (2026-03-13)  
**Versão Atual:** 1.4 (2026-03-27)
