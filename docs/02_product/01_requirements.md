---
title: Documento de Requisitos de Negócio (BRD) e Requisitos Funcionais
status: implemented
version: 2.0
last_updated: 2026-03-11
---

# Documento de Requisitos de Negócio (BRD) - Digna

Este documento consolida todas as capacidades do sistema e serve como o Guia Antimajoração para a IA e desenvolvedores. Se um requisito não estiver aqui ou nas Sprints, ele não deve ser codificado.

---

### 1. Requisitos Funcionais (RF)

#### Core Contábil e Operacional
*   **RF-01: Gestão de Identidade (Compliance Sinaes)**
    *   *Descrição:* Suportar login unificado via portal Gov.br (OAuth2) e permitir perfis híbridos.
    *   *Critério/Regra:* O sistema deve suportar metadados tanto para grupos informais ("Sonhos" / CPFs) quanto para entidades já formalizadas (CNPJs), permitindo a transição sem perda de dados históricos.
*   **RF-02: PDV Operacional e de Impacto (Ponto de Venda)**
    *   *Descrição:* Registro simplificado de vendas monetárias e operações comerciais na ponta para o usuário que não entende de contabilidade.
    *   *Funcionalidades:* Keyboard numérico para entrada de valores, seleção rápida de produtos cadastrados.
    *   *Regra de Negócio:* Toda venda gera automaticamente um lançamento invisível de Partida Dobrada no Ledger Lume (Débito: Caixa/Ativo | Crédito: Receita de Vendas).
*   **RF-03: Registro de Trabalho / Ponto Social (ITG 2002)**
    *   *Descrição:* Captura de horas de trabalho cooperativo, militante ou voluntário, garantindo a conformidade com a ITG 2002 do Conselho Federal de Contabilidade.
    *   *Funcionalidades:* Cronômetro em tempo real e registro manual de minutos vinculados ao membro.
    *   *Regra de Negócio:* As horas devem ser convertidas em "Capital Social de Trabalho" (mensurado em minutos) para servir de base ao rateio de sobras, invertendo a lógica capitalista (o suor vale tanto ou mais que o R$).
*   **RF-04: Motor de Reservas Obrigatórias (Lei 15.068/2024)**
    *   *Descrição:* Cálculo matemático e automático de fundos estatutários/legais antes do rateio de qualquer sobra.
    *   *Regras:* O sistema aplica um bloqueio inegociável de **10% para Reserva Legal** e **5% para FATES** (Fundo de Assistência Técnica Educacional e Social) sobre o excedente, impedindo a distribuição indevida.
*   **RF-05: Dossiê de Formalização (CADSOL/DCSOL)**
    *   *Descrição:* Motor de geração automática de documentação institucional e governança para comprovação de autogestão perante o Estado.
    *   *Funcionalidades:* Geração de Atas de Assembleia, Estatutos e Relatórios de Impacto exportáveis em Markdown e PDF.
    *   *Regra de Negócio:* O documento gerado deve obrigatoriamente conter um Hash criptográfico SHA256 embutido para atestar sua imutabilidade técnica. Só é liberado após o grupo registrar no mínimo 3 decisões na plataforma (Autogestão Gradual).
*   **RF-06: Sincronização Offline-First e Intercooperação B2B**
    *   *Descrição:* Capacidade de operar perfeitamente no "Brasil Profundo" (assentamentos, feiras rurais) sem internet.
    *   *Funcionalidades:* Delta tracking local das transações. Quando há rede, sincroniza pacotes criptografados com a nuvem, além de prover um Marketplace B2B fechado para trocas entre EES.

#### Finanças Solidárias
*   **RF-07: Gestão de Compras e Fornecedores**
    *   *Descrição:* Interface simplificada para aquisição de insumos ("O que comprou? De quem? Por quanto?").
    *   *Funcionalidades:* Cadastro de fornecedores, suporte a pagamentos à vista (CASH) e a prazo (CREDIT).
    *   *Regra de Negócio:* Contabilidade invisível gerando partidas dobradas no backend (Débito em Estoque/Despesa e Crédito no Caixa/Fornecedores). Valores transitam obrigatoriamente em `int64`.
*   **RF-08: Controle de Estoque**
    *   *Descrição:* Gestão de inventário categorizada para simplificar a visão de negócio da cooperativa.
    *   *Funcionalidades:* Categorização obrigatória em INSUMO, PRODUTO ou MERCADORIA. Controle de quantidade mínima, alertas de ruptura de estoque e baixa automática assim que o PDV registrar a venda.
*   **RF-09: Gestão de Caixa**
    *   *Descrição:* Espelho financeiro real do saldo do empreendimento.
    *   *Funcionalidades:* Controle de entradas, saídas e visualização do saldo em tempo real, lendo diretamente das transações consolidadas pelo motor Lume.
*   **RF-10: Gestão Orçamentária e Planejamento**
    *   *Descrição:* Ferramenta de planejamento financeiro cruzando o "planejado vs realizado" com linguagem extremamente acessível (sem jargões como CAPEX ou Forecast).
    *   *Funcionalidades:* Categorias pré-definidas e Alertas Visuais baseados em barras de progresso: SAFE (≤70%), WARNING (71-100%), EXCEEDED (>100%).

#### Aliança Contábil e Institucional
*   **RF-11: Aliança Contábil e Exportação Fiscal (SPED)**
    *   *Descrição:* Interface Multi-tenant (Accountant Dashboard) dedicada ao Contador Social parceiro para fechamento de balanços.
    *   *Funcionalidades:* Motor de tradução que converte o plano de contas "amigável" (Gaveta) para o Plano de Contas Referencial. Geração de arquivos CSV/SPED prontos para sistemas comerciais.
    *   *Regra de Negócio:* O núcleo Digna é blindado contra o cálculo de impostos comerciais; ele apenas exporta dados contábeis íntegros (Soma Zero). Acesso da ferramenta ocorre exclusivamente em modo `Read-Only`.
*   **RF-12: Gestão de Vínculo Contábil e Delegação Temporal [NOVO]**
    *   *Descrição:* Sistema de delegação de responsabilidade técnica entre o Empreendimento (EES) e o Contador.
    *   *Funcionalidades:* Controle de relacionamento por meio de datas de início e fim (`start_date`, `end_date`).
    *   *Regra de Negócio:* A cooperativa detém o "Exit Power", podendo revogar o acesso do contador a qualquer momento. Um EES pode ter tido vários contadores no histórico, mas apenas 1 ativo. Contadores desativados mantêm acesso estrito de leitura (Read-Only) **apenas** aos dados gerados durante o período de sua vigência, garantindo respaldo legal perante o CFC.
*   **RF-13: Visão Analítica do Contador Social (Auditoria ITG 2002) [NOVO]**
    *   *Descrição:* Dashboard consultivo inteligente que apoia o contador na geração de relatórios e Notas Explicativas.
    *   *Funcionalidades:* Painel consolidado que compila os montantes em caixa destinados à Reserva Legal, FATES e a somatória da valoração econômica do trabalho cooperativo/voluntário.
    *   *Regra de Negócio:* Todas as consultas ao banco de dados SQLite para gerar esta visão devem ser interceptadas e filtradas obrigatoriamente pela vigência temporal do vínculo do contador (RF-12).

---

### 2. Requisitos Não Funcionais (RNF)

*   **RNF-01: Soberania de Dados (SQLite-per-tenant)**
    *   O dado não pertence à "nuvem", pertence ao usuário. Cada entidade possui um arquivo físico isolado (`/data/entities/{entity_id}.db`). O backend orquestra a conexão via `LifecycleManager`. É tecnicamente impossível e proibido cruzar dados (JOINs) entre bancos de entidades diferentes.
*   **RNF-02: Integridade Financeira (Anti-Float)**
    *   Todos os cálculos monetários e de tempo devem utilizar inteiros de 64 bits (`int64`). O uso da tipagem `float32/64` é expressamente proibido nas camadas de Domínio, Serviço e Banco de Dados para evitar erros de dízima/arredondamento padrão IEEE 754.
*   **RNF-03: Resiliência Offline (PWA)**
    *   A interface deve permitir a operação contínua mesmo sem internet. A aplicação utiliza Manifest e Service Workers, gravando o cache local e realizando a sincronização de deltas (`sync_metadata`) de forma transparente quando há rede.
*   **RNF-04: Adequação Sociotécnica e Linguagem**
    *   É terminantemente proibido vazar jargões contábeis (ex: Débito, Crédito, Provisão) para o *frontend* do produtor. O sistema atua como tradutor cultural. O design deve atender a baixa literacia digital, empregando botões amplos para o toque (`min-h-[44px]`), alto contraste (WCAG 2.1 AA) e a paleta "Soberania e Suor".

---

### 3. Matriz de Regras Contábeis (Motor Lume)

| Evento | Conta Débito | Conta Crédito | Observação |
| :--- | :--- | :--- | :--- |
| **Venda de Produto** | Caixa / Ativo (1.1.01) | Receita de Vendas (3.1.01) | Lançamento Automático PDV |
| **Trabalho (Sócio)** | Despesa Social | Receita (Capital Trabalho) | Registro de Soma Zero (ITG 2002) |
| **Rateio de Sobras** | Sobras/Excedentes | Crédito p/ Sócio | Proporcional a Horas + Capital |
| **Compra de Insumo** | Estoque/Despesa | Caixa/Fornecedores | Lançamento Automático Supply |

---

### 4. Cláusula Anti-Alucinação para a IA

> "Qualquer implementação de código deve consultar primeiro o RF (Requisito Funcional) correspondente neste documento. Se a funcionalidade proposta não estiver no BRD, o agente deve solicitar confirmação do Operador Humano antes de prosseguir. A prioridade máxima é a **Soberania do Dado (RNF-01)** e o **Rigor Matemático em int64 (RNF-02)** sobre a estética ou funcionalidades 'extras'."
```

***

### 📋 Próximos Documentos a Atualizar (Fila de Trabalho)

Para manter o nosso sistema PKM (Personal Knowledge Management) 100% íntegro com a introdução do vínculo temporal do Contador, os seguintes documentos precisarão de ajustes em sequência:

1.  **`docs/06_roadmap/03_backlog.md`:** (Inserir os novos RF-12 e RF-13 na "Alta Prioridade" e marcá-los nas fases correspondentes).
2.  **`docs/02_product/02_models.md`:** (Modelagem de dados: inserir a entidade `EnterpriseAccountant` e detalhar o algoritmo de "Filtro de Vigência" para as consultas fiscais).
3.  **`docs/01_project/01_vision.md`:** (Expansão do Pilar "Aliança Contábil", formalizando o compromisso de proteger a responsabilidade técnica do contador parceiro).
4.  **`docs/03_architecture/01_system.md`:** (Atualizar as responsabilidades do módulo `accountant_dashboard` para incluir a regra de negócio do filtro temporal *Read-Only*).

O novo arquivo de Requisitos (BRD) está pronto. **Aguardo o seu sinal verde** para enviar o conteúdo completo e atualizado do próximo da fila: o **`03_backlog.md`**!