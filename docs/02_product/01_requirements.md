title: Documento de Requisitos de Negócio (BRD) e Requisitos Funcionais
status: implemented
version: 2.1
last_updated: 2026-03-27
---

# Documento de Requisitos de Negócio (BRD) - Ecossistema Digna

> **Nota:** Este documento consolida todas as capacidades do sistema e serve como o Guia Antimajoração para a IA e desenvolvedores. Se um requisito não estiver aqui ou nas Sprints, ele não deve ser codificado.
>
> **Versão 2.1:** Incorpora a visão de Ecossistema de 4 Módulos (PDF v1.0), o Sistema de Ajuda Educativa (RF-30 - decisão de design da sessão 27/03/2026), e preserva todo o trabalho validado nas Sprints 1-16.

---

## 1. Visão Estratégica do Ecossistema

O **Ecossistema Digna** é uma plataforma tecnológica integrada voltada ao empreendedorismo popular no Brasil. Seu propósito central é reduzir as barreiras de acesso à informação normativa, ao crédito e à gestão financeira para microempreendedores, pequenos negócios, cooperativas e organizações da sociedade civil que atuam em contextos de vulnerabilidade socioeconômica.

### 1.1 Interseção de Conhecimentos

O projeto nasce da interseção entre quatro campos do conhecimento:

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

## 2. Arquitetura do Ecossistema

O Ecossistema Digna é composto por **quatro módulos interdependentes**. O **digna ERP** é o núcleo — os demais módulos são extensões que ampliam o valor gerado a partir dos dados que o ERP já captura.

| Módulo | Função Principal | Integração-Chave | Status |
|--------|------------------|------------------|--------|
| **Módulo 1: digna ERP** | Gestão financeira, fiscal e contábil do empreendedor | Alimenta todos os demais módulos com dados do perfil | ✅ 85% Completo |
| **Módulo 2: Motor de Indicadores** | Coleta e interpreta indicadores econômicos em tempo real | Consome APIs BCB/IBGE; alimenta Portal com taxas e contexto | 📋 Backlog Fase 3 |
| **Módulo 3: Portal de Oportunidades** | Match automático entre perfil e programas de financiamento | Consome ERP + Motor; gera checklist de documentos | 📋 Backlog Fase 3 |
| **Módulo 4: Rede Digna** | Marketplace solidário entre entidades do ecossistema | Consome perfil do ERP para matching de compra/venda | 🔄 Expandir sync_engine |

---

## 3. Requisitos Funcionais (RF)

### 3.1 Core Contábil e Operacional (Módulo 1 - ERP)

#### RF-01: Gestão de Identidade (Compliance Sinaes)
- **Descrição:** Suportar login unificado via portal Gov.br (OAuth2) e permitir perfis híbridos.
- **Critério/Regra:** O sistema deve suportar metadados tanto para grupos informais ("Sonhos" / CPFs) quanto para entidades já formalizadas (CNPJs), permitindo a transição sem perda de dados históricos.
- **Status:** ✅ Implementado

#### RF-02: PDV Operacional e de Impacto (Ponto de Venda)
- **Descrição:** Registro simplificado de vendas monetárias e operações comerciais na ponta para o usuário que não entende de contabilidade.
- **Funcionalidades:** Keyboard numérico para entrada de valores, seleção rápida de produtos cadastrados.
- **Regra de Negócio:** Toda venda gera automaticamente um lançamento invisível de Partida Dobrada no Ledger Lume (Débito: Caixa/Ativo | Crédito: Receita de Vendas).
- **Status:** ✅ Implementado

#### RF-03: Registro de Trabalho / Ponto Social (ITG 2002)
- **Descrição:** Captura de horas de trabalho cooperativo, militante ou voluntário, garantindo a conformidade com a ITG 2002 do Conselho Federal de Contabilidade.
- **Funcionalidades:** Cronômetro em tempo real e registro manual de minutos vinculados ao membro.
- **Regra de Negócio:** As horas devem ser convertidas em "Capital Social de Trabalho" (mensurado em minutos) para servir de base ao rateio de sobras, invertendo a lógica capitalista (o suor vale tanto ou mais que o R$).
- **Status:** ✅ Implementado

#### RF-04: Motor de Reservas Obrigatórias (Lei 15.068/2024)
- **Descrição:** Cálculo matemático e automático de fundos estatutários/legais antes do rateio de qualquer sobra.
- **Regras:** O sistema aplica um bloqueio inegociável de **10% para Reserva Legal** e **5% para FATES** (Fundo de Assistência Técnica Educacional e Social) sobre o excedente, impedindo a distribuição indevida.
- **Status:** ✅ Implementado

#### RF-05: Dossiê de Formalização (CADSOL/DCSOL)
- **Descrição:** Motor de geração automática de documentação institucional e governança para comprovação de autogestão perante o Estado.
- **Funcionalidades:** Geração de Atas de Assembleia, Estatutos e Relatórios de Impacto exportáveis em Markdown e PDF.
- **Regra de Negócio:** O documento gerado deve obrigatoriamente conter um Hash criptográfico SHA256 embutido para atestar sua imutabilidade técnica. Só é liberado após o grupo registrar no mínimo 3 decisões na plataforma (Autogestão Gradual).
- **Status:** ✅ Implementado

#### RF-06: Sincronização Offline-First e Intercooperação B2B
- **Descrição:** Capacidade de operar perfeitamente no "Brasil Profundo" (assentamentos, feiras rurais) sem internet.
- **Funcionalidades:** Delta tracking local das transações. Quando há rede, sincroniza pacotes criptografados com a nuvem, além de prover um Marketplace B2B fechado para trocas entre EES.
- **Status:** ✅ Implementado (base)

#### RF-07: Gestão de Compras e Fornecedores
- **Descrição:** Interface simplificada para aquisição de insumos ("O que comprou? De quem? Por quanto?").
- **Funcionalidades:** Cadastro de fornecedores, suporte a pagamentos à vista (CASH) e a prazo (CREDIT).
- **Regra de Negócio:** Contabilidade invisível gerando partidas dobradas no backend (Débito em Estoque/Despesa e Crédito no Caixa/Fornecedores). Valores transitam obrigatoriamente em `int64`.
- **Status:** ✅ Implementado

#### RF-08: Controle de Estoque
- **Descrição:** Gestão de inventário categorizada para simplificar a visão de negócio da cooperativa.
- **Funcionalidades:** Categorização obrigatória em INSUMO, PRODUTO ou MERCADORIA. Controle de quantidade mínima, alertas de ruptura de estoque e baixa automática assim que o PDV registrar a venda.
- **Status:** ✅ Implementado

#### RF-09: Gestão de Caixa
- **Descrição:** Espelho financeiro real do saldo do empreendimento.
- **Funcionalidades:** Controle de entradas, saídas e visualização do saldo em tempo real, lendo diretamente das transações consolidadas pelo motor Lume.
- **Status:** ✅ Implementado

#### RF-10: Gestão Orçamentária e Planejamento
- **Descrição:** Ferramenta de planejamento financeiro cruzando o "planejado vs realizado" com linguagem extremamente acessível (sem jargões como CAPEX ou Forecast).
- **Funcionalidades:** Categorias pré-definidas e Alertas Visuais baseados em barras de progresso: SAFE (≤70%), WARNING (71-100%), EXCEEDED (>100%).
- **Status:** ✅ Implementado

#### RF-27: Cálculo Automático de DAS MEI [NOVO - PDF v1.0]
- **Descrição:** Cálculo automático do DAS MEI (Documento de Arrecadação do Simples Nacional) para microempreendedores individuais.
- **Regra de Negócio:** O DAS MEI corresponde a **5% do salário mínimo vigente** + valores fixos de ICMS/ISS (se aplicável). O sistema deve manter uma tabela interna versionada do salário mínimo por ano.
- **Funcionalidades:**
  - Tabela de salário mínimo versionada por ano (atualização via decreto presidencial)
  - Alertas de vencimento (dia 20 de cada mês)
  - Histórico de pagamentos registrados
  - Geração de guia para pagamento (link ou exportação)
- **Critério de Aceite:** Nenhum usuário precisa calcular manualmente o DAS
- **Status:** 📋 Backlog (Prioridade Alta)

---

### 3.2 Aliança Contábil e Institucional (Módulo 1 - ERP)

#### RF-11: Aliança Contábil e Exportação Fiscal (SPED)
- **Descrição:** Interface Multi-tenant (Accountant Dashboard) dedicada ao Contador Social parceiro para fechamento de balanços.
- **Funcionalidades:** Motor de tradução que converte o plano de contas "amigável" (Gaveta) para o Plano de Contas Referencial. Geração de arquivos CSV/SPED prontos para sistemas comerciais.
- **Regra de Negócio:** O núcleo Digna é blindado contra o cálculo de impostos comerciais; ele apenas exporta dados contábeis íntegros (Soma Zero). Acesso da ferramenta ocorre exclusivamente em modo `Read-Only`.
- **Status:** ✅ Implementado

#### RF-12: Gestão de Vínculo Contábil e Delegação Temporal [NOVO]
- **Descrição:** Sistema de delegação de responsabilidade técnica entre o Empreendimento (EES) e o Contador.
- **Funcionalidades:** Controle de relacionamento por meio de datas de início e fim (`start_date`, `end_date`).
- **Regra de Negócio:** A cooperativa detém o "Exit Power", podendo revogar o acesso do contador a qualquer momento. Um EES pode ter tido vários contadores no histórico, mas apenas 1 ativo. Contadores desativados mantêm acesso estrito de leitura (Read-Only) **apenas** aos dados gerados durante o período de sua vigência, garantindo respaldo legal perante o CFC.
- **Status:** ✅ 95% Implementado

#### RF-13: Visão Analítica do Contador Social (Auditoria ITG 2002) [NOVO]
- **Descrição:** Dashboard consultivo inteligente que apoia o contador na geração de relatórios e Notas Explicativas.
- **Funcionalidades:** Painel consolidado que compila os montantes em caixa destinados à Reserva Legal, FATES e a somatória da valoração econômica do trabalho cooperativo/voluntário.
- **Regra de Negócio:** Todas as consultas ao banco de dados SQLite para gerar esta visão devem ser interceptadas e filtradas obrigatoriamente pela vigência temporal do vínculo do contador (RF-12).
- **Status:** ✅ 95% Implementado

#### RF-14: Blindagem Tributária (EFD-Reinf e ECF) [NOVO - Adequação Estatal]
- **Descrição:** Módulo `tax_compliance` para mensageria de retenções via Web Service e expurgo automático de receitas de Atos Cooperativos no Bloco M da ECF.
- **Funcionalidades:**
  - Geração e transmissão de XMLs (série R-2000/R-4000) para Web Services da EFD-Reinf
  - Alimentação automática da DCTFWeb
  - Expurgo de receitas de Atos Cooperativos no Bloco M da ECF (Lei 5.764/71 e LC 214/2025)
- **Status:** 📋 Backlog (Prioridade Alta)

#### RF-15: Integração Real Gov.br e Governança Digital [ATUALIZADO]
- **Descrição:** Substituição do Mock de login unificado pelo fluxo real da Cidadania Digital (OAuth2).
- **Funcionalidades:**
  - Assinatura Eletrônica Qualificada (Lei nº 14.063/2020) para mesas de Assembleia
  - Algoritmo de anonimização sistêmica para escrutínio secreto em votações (IN DREI nº 79/2020)
  - Integração com ICP-Brasil para certificados digitais
- **Status:** 📋 Backlog (Prioridade Alta)

#### RF-16: Inclusão Sanitária (MAPA) [NOVO - Adequação Estatal]
- **Descrição:** Módulo `sanitary_compliance` para geração automatizada do Memorial Técnico Sanitário de Estabelecimento (MTSE).
- **Funcionalidades:**
  - Conformidade com Portaria MAPA nº 393/2021 para agroindústrias
  - Parametrização de fluxogramas de maquinário, capacidade diária e potabilidade de água
  - Exportação para peticionamento no Sistema Eletrônico de Informação (SEI)
- **Status:** 📋 Backlog (Prioridade Média)

#### RF-17: Integração CADSOL/SINAES Automático [NOVO - Adequação Estatal]
- **Descrição:** Consumo nativo de Web Services do MTE (Decreto nº 12.784/2025).
- **Funcionalidades:**
  - Matrícula automática de entidades FORMALIZED no Cadastro Nacional de Economia Solidária
  - Substituição dos mocks atuais por integração real
- **Status:** 📋 Backlog (Prioridade Média)

---

### 3.3 Motor de Indicadores Econômico-Financeiros (Módulo 2) [NOVO - PDF v1.0]

#### RF-18: Motor de Indicadores Econômico-Financeiros
- **Descrição:** Coleta, armazenamento e interpretação de dados econômicos externos em tempo real.
- **Diferencial:** Não é apenas um painel de indicadores — o sistema interpreta o dado no contexto do negócio: "A SELIC subiu — veja como isso afeta seu custo de crédito".
- **Fontes de Dados:**
  - **BCB SGS:** SELIC, IPCA, CDI, INPC, TJLP (Diária/Mensal)
  - **BCB PTAX:** Câmbio USD/BRL, EUR/BRL (Diária)
  - **BCB Focus:** Expectativas de mercado (Semanal)
  - **IBGE SIDRA:** IPCA, INPC, PNAD, PIB trimestral (Mensal/Trimestral)
  - **Portal da Transparência:** Convênios federais, CAGED
- **Arquitetura Técnica:**
  - **Coleta:** Scheduler diário que executa chamadas às APIs externas
  - **Cache:** Tabela `indicators{}` com TTL para evitar chamadas redundantes
  - **Interpretação:** Regras de negócio que transformam valor bruto em orientação contextualizada
  - **API Interna:** Endpoint REST consumido pelo ERP, Portal e interface do usuário
- **Critério de Aceite:** Nenhum usuário preenche dados econômicos manualmente
- **Status:** 📋 Backlog (Prioridade Alta)

---

### 3.4 Perfil de Elegibilidade e Portal de Oportunidades (Módulo 3) [NOVO - PDF v1.0]

#### RF-19: Perfil de Elegibilidade (Campos Complementares)
- **Descrição:** Expandir o modelo de dados do ERP para capturar informações necessárias ao match de crédito, com preenchimento único e reutilização automática.
- **Dados já disponíveis no ERP (capturados automaticamente):**
  - CNPJ/CPF, CNAE, Município, UF
  - Faturamento anual (int64), Regime Tributário
  - Data de abertura, Situação fiscal
  - Certidões negativas (integração com portais de certidões)
- **Novos Campos (preenchimento único):**
  - Inscrito no CadÚnico (bool) - Habilita programas sociais
  - Sócio Mulher (bool) - Prioridade em linhas com foco de gênero
  - Inadimplência Ativa (bool) - Direciona ao Desenrola antes de crédito
  - Finalidade do Crédito (enum) - CAPITAL_GIRO, EQUIPAMENTO, REFORMA, OUTRO
  - Valor Necessário (int64) - Anti-Float compliance
  - Tipo de Entidade (enum) - MEI, ME, EPP, Cooperativa, OSC, OSCIP, PF
  - Contabilidade Formal (bool) - Requisito de alguns programas
- **Critério de Aceite:** Nenhum dado digitado duas vezes
- **Status:** 📋 Backlog (Prioridade Alta)

#### RF-20: Portal de Oportunidades (MVP - 3 Programas)
- **Descrição:** Match automático entre perfil da entidade e programas de financiamento.
- **Programas Prioritários para MVP:**
  - **Acredita no Primeiro Passo** (~R$ 6 mil, CadÚnico, 70% mulheres)
  - **Pronampe** (Até 30% da receita bruta, MEI/ME/EPP)
  - **Niterói Empreendedora** (Até R$ 200 mil, JUROS ZERO, foco em mulheres e jovens)
- **Funcionalidades:**
  - Match Automático: Cruzamento perfil × requisitos de cada programa
  - Ranqueamento por Vantagem: Ordenação por custo efetivo de capital (usando taxas do Motor)
  - Checklist de Documentos: Marca o que já está no Digna vs. o que falta
  - Alertas de Prazo: Notificações quando editais com alto match têm prazo próximo
- **Critério de Aceite:** Entidade descobre elegibilidade sem preencher formulários
- **Status:** 📋 Backlog (Prioridade Alta)

#### RF-21: Checklist + Alertas de Documentos [NOVO]
- **Descrição:** Sistema de acompanhamento de documentos necessários para candidaturas.
- **Funcionalidades:**
  - Integração com portais de certidões (PGFN, Estadual, Municipal) para verificação automática
  - Scraping do Diário Oficial da União (DOU) para captura de novos editais
  - Sistema de notificações push/email para prazos críticos
  - Dashboard de documentos pendentes por programa
- **Status:** 📋 Backlog (Prioridade Média)

#### RF-22: Monitoramento de DOU [NOVO]
- **Descrição:** Scraper de Diário Oficial para novos editais.
- **Funcionalidades:**
  - Classificação automática por relevância (match com perfil da entidade)
  - Notificação proativa para Contador Social e entidade
- **Status:** 📋 Backlog (Prioridade Média)

#### RF-23: Dashboard de Elegibilidade [NOVO]
- **Descrição:** Visão consolidada de programas elegíveis.
- **Funcionalidades:**
  - Status de cada candidatura (Elegível, Documentação Pendente, Enviada, Aprovada)
  - Histórico de candidaturas submetidas
- **Status:** 📋 Backlog (Prioridade Média)

---

### 3.5 Rede Digna - Marketplace Solidário (Módulo 4) [NOVO - PDF v1.0]

#### RF-24: Perfil Público da Entidade
- **Descrição:** Permitir que entidades publiquem informações visíveis para a rede, sem expor dados sensíveis.
- **Campos do Perfil Público:**
  - EntityID (Hash anonimizado)
  - Nome Fantasia, Missão (texto livre)
  - Produtos (lista de categorias)
  - Serviços (lista de capacidades)
  - Município, UF, Contato Público
  - Foto/Logo (URL do asset)
- **Critério de Aceite:** Dados sensíveis nunca expostos publicamente
- **Status:** 📋 Backlog (Prioridade Baixa)

#### RF-25: Mural de Necessidades
- **Descrição:** Entidades publicam demandas de compra visíveis para a rede.
- **Estrutura de Postagem:**
  - ID, PublisherID (Hash anonimizado)
  - Categoria (INSUMO, EQUIPAMENTO, SERVICO, OUTRO)
  - Descrição, Quantidade, Prazo Desejado
  - Município (para matching geográfico)
  - Status (ABERTO, EM_NEGOCIACAO, CONCLUIDO)
- **Critério de Aceite:** Entidades publicam demandas visíveis para a rede
- **Status:** 📋 Backlog (Prioridade Baixa)

#### RF-26: Match de Oportunidades B2B
- **Descrição:** Algoritmo que sugere conexões entre quem precisa e quem oferece.
- **Critérios de Matching:**
  - Geográfico: Priorizar proximidade (mesmo município/UF)
  - Setorial: Afinidade de CNAE/categoria
  - Temporal: Prazos compatíveis
  - Reputação: Histórico de transações na rede (futuro)
- **Critério de Aceite:** Sistema sugere conexões entre quem precisa e quem oferece
- **Status:** 📋 Backlog (Prioridade Baixa)

---

### 3.6 Sistema de Ajuda Educativa Estruturada [NOVO - Decisão de Design 27/03/2026]

#### RF-30: Sistema de Ajuda Educativa Estruturada
- **Descrição:** Sistema de ajuda contextual que traduz conceitos técnicos (CadÚnico, inadimplência, CNAE, etc.) em linguagem popular, com linkagem entre elementos de UI e registros de ajuda no banco.
- **Funcionalidades:**
  - Entrada de menu "Ajuda" acessível em todas as páginas
  - Busca e índice de tópicos categorizados (CRÉDITO, TRIBUTÁRIO, GOVERNANÇA, GERAL)
  - Explicação em linguagem popular + legislação relacionada + próximo passo acionável
  - Linkagem automática: botão "?" ao lado de campos técnicos abre explicação contextual
  - Tópicos com estrutura: Título, Resumo, Explicação, Por que perguntamos, Legislação, Próximo passo, Link oficial
- **Critério de Aceite:** Usuário com 5ª série consegue entender a explicação sem ajuda externa; tooltip carrega em < 500ms via HTMX; linguagem sem jargões técnicos.
- **Status:** 📋 Backlog (Prioridade Média - Habilita adoção por baixa escolaridade)

---

## 4. Requisitos Não Funcionais (RNF)

*   **RNF-01: Soberania de Dados (SQLite-per-tenant)**
    *   O dado não pertence à "nuvem", pertence ao usuário. Cada entidade possui um arquivo físico isolado (`/data/entities/{entity_id}.db`). O backend orquestra a conexão via `LifecycleManager`. É tecnicamente impossível e proibido cruzar dados (JOINs) entre bancos de entidades diferentes.
*   **RNF-02: Integridade Financeira (Anti-Float)**
    *   Todos os cálculos monetários e de tempo devem utilizar inteiros de 64 bits (`int64`). O uso da tipagem `float32/64` é expressamente proibido nas camadas de Domínio, Serviço e Banco de Dados para evitar erros de dízima/arredondamento padrão IEEE 754.
*   **RNF-03: Resiliência Offline (PWA)**
    *   A interface deve permitir a operação contínua mesmo sem internet. A aplicação utiliza Manifest e Service Workers, gravando o cache local e realizando a sincronização de deltas (`sync_metadata`) de forma transparente quando há rede.
*   **RNF-04: Adequação Sociotécnica e Linguagem**
    *   É terminantemente proibido vazar jargões contábeis (ex: Débito, Crédito, Provisão) para o *frontend* do produtor. O sistema atua como tradutor cultural. O design deve atender a baixa literacia digital, empregando botões amplos para o toque (`min-h-[44px]`), alto contraste (WCAG 2.1 AA) e a paleta "Soberania e Suor".
*   **RNF-05: Cache-Proof Templates [NOVO]**
    *   Templates `*_simple.html` devem ser carregados via `ParseFiles()` no handler, não como variáveis globais. Isso garante atualizações imediatas sem necessidade de recompilar o binário.
*   **RNF-06: Laicidade do Produto [NOVO - PDF v1.0]**
    *   Embora o fundamento filosófico do projeto inclua princípios da Teologia Cristã (Dignidade Humana, Mordomia, Koinonia, Justiça Restaurativa), a plataforma **não apresenta conteúdo religioso explícito ao usuário final**. A Teologia informa as decisões de design e ética, mas o produto é acessível e útil independente da crença do usuário.
*   **RNF-07: Performance e Escalabilidade [NOVO]**
    *   **Tempo de resposta:** < 200ms para consultas locais SQLite
    *   **Concorrência:** Suporte a 100+ usuários simultâneos por instância
    *   **Cache:** Indicadores econômicos com TTL de 24h para APIs externas

---

## 5. Matriz de Regras Contábeis (Motor Lume)

| Evento | Conta Débito | Conta Crédito | Observação |
| :--- | :--- | :--- | :--- |
| **Venda de Produto** | Caixa / Ativo (1.1.01) | Receita de Vendas (3.1.01) | Lançamento Automático PDV |
| **Trabalho (Sócio)** | Despesa Social | Receita (Capital Trabalho) | Registro de Soma Zero (ITG 2002) |
| **Rateio de Sobras** | Sobras/Excedentes | Crédito p/ Sócio | Proporcional a Horas + Capital |
| **Compra de Insumo** | Estoque/Despesa | Caixa/Fornecedores | Lançamento Automático Supply |

---

## 6. Perfil de Elegibilidade — Modelo de Dados [NOVO - PDF v1.0]

O Perfil de Elegibilidade é o conjunto de campos que o Portal precisa para executar o match com programas de financiamento.

### 6.1 Dados já disponíveis no ERP

| Campo | Como é Capturado |
|-------|------------------|
| CNPJ/CPF | Cadastro inicial da entidade |
| CNAE principal e secundários | Cadastro inicial — consulta automática na Receita Federal |
| Município e UF | Endereço cadastrado |
| Faturamento bruto anual | Agregado dos lançamentos de vendas do exercício |
| Regime tributário (MEI, Simples, etc.) | Configuração inicial do perfil fiscal |
| Data de abertura do CNPJ | Consulta automática na Receita Federal |
| Situação cadastral (ativa/irregular) | Consulta automática periódica na Receita Federal |
| Certidões negativas | Integração com portais de certidões (PGFN, Estadual, Municipal) |

### 6.2 Dados complementares — preenchimento único

| Campo | Por que é Necessário |
|-------|---------------------|
| Inscrito no CadÚnico? (s/n) | Habilita acesso ao Acredita no Primeiro Passo e outros programas sociais |
| Sócio(a) mulher responsável? (s/n) | Prioridade em Pronampe, Niterói Empreendedora e outros |
| Possui inadimplência ativa? (s/n) | Direciona ao Desenrola antes de candidaturas a crédito novo |
| Finalidade do crédito buscado | Capital de giro/ Investimento em equipamentos/ Reforma/ Outro |
| Valor aproximado necessário | Filtra programas pelo teto de crédito |
| Tipo de entidade | MEI/ ME/ EPP/ Cooperativa/ OSC/ OSCIP/ Pessoa física |
| Possui contabilidade formal? (s/n) | Requisito de alguns programas (balanço dos últimos 2 anos) |

---

## 7. Estratégia de Implementação [NOVO - PDF v1.0]

Dado que o digna ERP já está em desenvolvimento, a estratégia proposta é incremental: priorizar os módulos que ampliam o valor do ERP existente com menor esforço, reservando os mais complexos para fases posteriores.

| Fase | Escopo | Prioridade |
|------|--------|------------|
| **Fase 1 — Motor de Indicadores** | Implementar o scheduler de coleta das APIs do BCB e IBGE. Criar a tabela de cache local. Exibir indicadores contextualizados no ERP. | Alta — baixo esforço, alto impacto imediato |
| **Fase 2 — Perfil de Elegibilidade** | Adicionar os campos complementares ao cadastro do ERP (CadÚnico, gênero do sócio, inadimplência, finalidade). Consulta automática de certidões. | Alta — habilita o Portal |
| **Fase 3 — Portal (MVP)** | Match com os 3 programas de maior alcance: Acredita no Primeiro Passo, Pronampe e Niterói Empreendedora. Checklist de documentos básico. | Alta — gera valor imediato e validável |
| **Fase 4 — Contador Social** | Módulo de adoção de entidades por contadores. Dashboard do contador com visão consolidada. Submissão de candidaturas. | Média — estratégica para escala |
| **Fase 5 — Portal completo** | Adicionar todos os programas mapeados. Monitoramento de DOU. Alertas automáticos. Ranqueamento por custo efetivo. | Média — expansão do MVP |
| **Fase 6 — Rede Digna** | Perfil público, mural de necessidades, match de oportunidades entre entidades. | Baixa — depende de massa crítica de usuários |

---

## 8. Cláusula Anti-Alucinação para a IA

> "Qualquer implementação de código deve consultar primeiro o RF (Requisito Funcional) correspondente neste documento. Se a funcionalidade proposta não estiver no BRD, o agente deve solicitar confirmação do Operador Humano antes de prosseguir. A prioridade máxima é a **Soberania do Dado (RNF-01)** e o **Rigor Matemático em int64 (RNF-02)** sobre a estética ou funcionalidades 'extras'."

---

## 9. Fundamento Filosófico-Teológico [NOVO - PDF v1.0]

A Teologia Cristã não é um elemento decorativo do projeto — é a base filosófica que orienta decisões concretas de produto, design e relacionamento com o usuário.

| Princípio | Aplicação Prática no Digna |
|-----------|---------------------------|
| **Dignidade Humana (Imago Dei)** | Interface que nunca humilha. Linguagem sem jargão excludente. Design que pressupõe competência, não ignorância. Ausência de mensagens punitivas por inadimplência |
| **Mordomia (Stewardship)** | O sistema incentiva gestão responsável dos recursos. Relatórios que mostram saúde financeira de forma construtiva. Ferramentas de planejamento, não de maximização de lucro a qualquer custo |
| **Koinonía (Comunhão)** | A funcionalidade de Rede Digna não é apenas negócio — é comunidade econômica que se sustenta mutuamente. Igrejas e líderes comunitários como agentes de distribuição e capacitação |
| **Justiça Restaurativa (Shalom)** | Acesso à informação normativa como ato de equidade. O que o contador caro faz pelo rico, o sistema faz pelo pobre. Métricas de bem-estar que vão além do lucro |

---

## 10. Próximos Documentos a Atualizar (Fila de Trabalho)

Para manter o nosso sistema PKM (Personal Knowledge Management) 100% íntegro com esta expansão de requisitos, os seguintes documentos precisarão de ajustes em sequência:

1. `docs/06_roadmap/02_roadmap.md`: (Incorporar Fases 3-4 do Ecossistema)
2. `docs/06_roadmap/03_backlog.md`: (Inserir RF-18 a RF-27 e RF-30 na "Alta Prioridade")
3. `docs/02_product/02_models.md`: (Modelagem de dados: inserir `EligibilityProfile`, `Indicator`, `PublicProfile`, `HelpTopic`)
4. `docs/03_architecture/01_system.md`: (Atualizar arquitetura para 4 módulos + Sistema de Ajuda)

---

**Status:** ✅ ATUALIZADO COM VISÃO DE ECOSSISTEMA (PDF v1.0) + RF-30 (Decisão de Design 27/03/2026)  
**Próxima Ação:** Atualizar `06_roadmap/02_roadmap.md` e `06_roadmap/03_backlog.md` com RF-30 incluído  
**Versão Anterior:** 2.0 (2026-03-11)  
**Versão Atual:** 2.1 (2026-03-27)

===== END FILE: 02_product/01_requirements.md =====
```

---

## 📋 Resumo das Alterações Realizadas (Merge Real)

| Seção | Alteração | Justificativa |
|-------|-----------|---------------|
| **Cabeçalho** | Versão 2.1 + `last_updated: 2026-03-27` | Versionamento corrigido com data real |
| **Visão Estratégica** | Expandida com 4 áreas de conhecimento | PDF v1.0, Seção 1 |
| **Arquitetura** | Expandida para 4 módulos interdependentes | Visão de ecossistema do PDF |
| **RF-27** | Adicionado (DAS MEI) | PDF v1.0, Seção 4.1 |
| **RF-14 a RF-17** | Adicionados (Adequação Estatal) | Backlog atual mais detalhado que PDF |
| **RF-18** | Adicionado (Motor de Indicadores) | PDF v1.0, Seção 5 |
| **RF-19 a RF-23** | Adicionados (Portal de Oportunidades) | PDF v1.0, Seção 6 |
| **RF-24 a RF-26** | Adicionados (Rede Digna) | PDF v1.0, Seção 7 |
| **RF-30** | **Adicionado (Sistema de Ajuda Educativa)** | **Decisão de design da sessão 27/03/2026** |
| **RNF-05 a RNF-07** | Adicionados (Cache-Proof, Laicidade, Performance) | Decisões arquiteturais críticas |
| **Perfil de Elegibilidade** | Nova seção com modelo de dados completo | PDF v1.0, Seção 9 |
| **Estratégia de Implementação** | Nova seção com 6 fases | PDF v1.0, Seção 10 |
| **Fundamento Teológico** | Nova seção com 4 princípios | PDF v1.0, Seção 2 |
