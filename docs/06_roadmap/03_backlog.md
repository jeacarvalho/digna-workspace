###### title: Product Backlog
status: proposed version: 1.3 last_updated: 2026-03-13
### Product Backlog - Digna
Este documento rastreia todas as funcionalidades já implementadas e as futuras, alinhadas às Fases do Roadmap e aos Requisitos Funcionais (RF) definidos no BRD.

--------------------------------------------------------------------------------

#### Funcionalidades Concluídas
##### Core Contábil e Fundação (Sprints 01 a 06)
*  [x] PDV Operacional com Partidas Dobradas Invisíveis (RF-01)
*  [x] Registro de Trabalho / Ponto Social ITG 2002 (RF-02)
*  [x] Motor de Reservas Obrigatórias (RF-03)
*  [x] Dossiê de Formalização (RF-04)
*  [x] Sincronização Offline-First via Delta Tracking (RF-05)
*  [x] Intercooperação B2B - Base de dados (RF-06)
*  [x] Gestão de Caixa conectada ao Motor Lume (RF-09)

##### Gestão de Membros (Sprint 10)
*  [x] Cadastro de membros
*  [x] Roles (COORDINATOR, MEMBER, ADVISOR) com hierarquia de permissões
*  [x] Status (ACTIVE, INACTIVE) e trava contra desativação do último coordenador
*  [x] Skills/habilidades associadas ao trabalho voluntário

##### Formalização (Sprint 11)
*  [x] Transição automática de status DREAM → FORMALIZED
*  [x] Simulador de formalização para análise de impacto
*  [x] Funcionalidade CheckFormalizationCriteria()
*  [x] Funcionalidade AutoTransitionIfReady() (gatilho após 3 decisões)

##### Integrações Governamentais (Mocks)
*  [x] Receita Federal (CNPJ, DARF) - Mock
*  [x] MTE (RAIS, CAT, eSocial) - Mock
*  [x] MDS (CadÚnico, Relatório Social) - Mock
*  [x] IBGE (Pesquisas, PAM, CNAE) - Mock
*  [x] SEFAZ (NFe, NFS-e, Manifesto) - Mock
*  [x] BNDES (Linhas de Crédito, Simulação) - Mock
*  [x] SEBRAE (Cursos, Consultoria) - Mock
*  [x] Providentia (Sync, Marketplace) - Mock

##### Finanças Solidárias e Suprimentos (Sprints 13 e 14)
*  [x] Gestão de compras e fornecedores (RF-07)
*  [x] Controle de estoque por categorização obrigatória (INSUMO, PRODUTO, MERCADORIA) com baixa via PDV (RF-08)
*  [x] Gestão Orçamentária e Planejamento Financeiro com alertas visuais SAFE/WARNING/EXCEEDED (RF-10)

##### Estabilização, UI e Rateio (Sprints 15 e 16)
*  [x] Surplus Calculator (CalculateSocialSurplus(), deduções automáticas de 10% Reserva Legal + 5% FATES)
*  [x] Identidade Visual Global "Soberania e Suor" aplicada (RNF-07)
*  [x] Arquitetura de renderização de Templates Cache-Proof (_simple.html)
*  [x] Integração End-to-End funcional (PDV → Estoque → Caixa)

--------------------------------------------------------------------------------

#### Funcionalidades Futuras
##### Alta Prioridade (Fase 1 e Fase 2 - Institucional e Aliança Contábil)
*  [ ]  **Gestão de Vínculo Contábil e Delegação Temporal (RF-12) [NOVO]:**  Sistema de associação temporal (start_date, end_date) entre Empreendimentos (EES) e Contadores Sociais. Garante o "Exit Power" (Soberania) da cooperativa e assegura o direito de consulta retroativo (Read-Only) do contador aos períodos em que foi Responsável Técnico.
*  [ ]  **Visão Analítica do Contador Social (RF-13) [NOVO]:**  Dashboard consultivo exclusivo e filtrado pelo vínculo temporal. Compila a saúde dos fundos obrigatórios (FATES e Reserva Legal) e a somatória do capital social de trabalho (ITG 2002) para elaboração de Notas Explicativas.
*  [ ]  **Interface de Gestão de Membros (UI):**  Criação das telas HTMX para o motor de CRUD de membros já finalizado no backend.
*  [ ]  **Módulo tax_compliance (EFD-Reinf e ECF) [RF-14] [NOVO]:** Desenvolver motor de geração e mensageria de XMLs (série R-2000/R-4000) para Web Services da EFD-Reinf. Garantir contábil e sistemicamente o preenchimento do Bloco M da ECF (e-Lalur/e-Lacs), executando o expurgo matemático das receitas de Atos Cooperativos da base de cálculo, blindando o produtor contra a bitributação (Lei nº 5.764/71 e LC nº 214/2025).
*  [ ]  **Integração Real Gov.br e Governança Digital (RF-15) [ATUALIZADO]:** Substituição do Mock de login unificado pelo fluxo real da Cidadania Digital (OAuth2). Atualizar o módulo `legal_facade` substituindo o aceite simples do aplicativo pelo consumo de APIs de Assinatura Eletrônica Avançada/Qualificada (Lei nº 14.063/2020) para os membros da mesa nas Atas de Assembleia. Desenvolver algoritmo de anonimização sistêmica para garantir escrutínio secreto em votações eletrônicas (IN DREI nº 79/2020).
*  [ ]  **Exportação CADSOL Oficial:**  Motor em PDF/Markdown que agrupa as decisões e estatutos para protocolo governamental.
*  [ ]  **Auditoria Pública (Cidadania Criptográfica):**  Tela para que qualquer cidadão audite as atas geradas através da validação de Hashes SHA256.

##### Média Prioridade (Fase 3 - Finanças Solidárias Avançadas e Adequação Estatal)
*  [ ]  **Múltiplas Moedas Sociais:**  Expansão do Ledger para registrar e transacionar moedas complementares de Bancos Comunitários de Desenvolvimento (BCDs).
*  [ ]  **Estoque Substantivo:**  Suporte à contabilidade não-monetária, gerindo "Fundos Rotativos Solidários" (troca e controle genético de sementes, animais para repasse).
*  [ ]  **Rateio de Sobras na Interface (UI):**  Painel visual para que a Assembleia aprove a divisão justa do excedente com base nas horas trabalhadas (transparência algorítmica).
*  [ ]  **Módulo sanitary_compliance (MAPA) [RF-16] [NOVO]:** Criar motor gerador do Memorial Técnico Sanitário de Estabelecimento (MTSE) em conformidade com a Portaria MAPA nº 393/2021. O sistema deverá parametrizar fluxogramas e capacidades da agroindústria, exportando a documentação exigida para peticionamento no SIF/SEI.
*  [ ]  **Integração SINAES / CADSOL Automático [RF-17] [NOVO]:** Expandir o módulo `integrations` para abandonar os *mocks* e consumir ativamente as novas APIs do MTE (Decreto nº 12.784/2025). O sistema deve matricular a entidade formalizada diretamente no Cadastro Nacional de Economia Solidária.

##### Baixa Prioridade (Fase 4 - Intercooperação e Escala Nacional)
*  [ ]  **Integração Contábil Fiscal Definitiva:**  Conexão direta por API (sem arquivos intermediários) com softwares comerciais de contabilidade (via Contador Social).
*  [ ]  **Marketplace B2B Solidário:**  Interface para trocas e vendas em lote apenas entre entidades autenticadas na rede Digna.
*  [ ]  **Score de Crédito Social:**  Motor que calcula a reputação da entidade baseada no trabalho e autogestão (não apenas em dinheiro) para solicitar crédito ao BNDES.
*  [ ]  **API Pública Restrita (OpenAPI/Swagger):**  Documentação técnica e geração de endpoints para ecossistemas parceiros (Serpro, Governos Estaduais).
*  [ ]  **App Mobile Nativo / Relatórios PDF Avulsos.**
