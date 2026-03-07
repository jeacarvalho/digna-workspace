#### title: Requisitos do Projeto Digna
status: implemented
version: 1.2
last_updated: 2026-03-07

### Requisitos - Projeto Digna
**Referência Legal:** Lei nº 15.068/2024 (Lei Paul Singer) / ITG 2002
**Projeto:** Digna - Infraestrutura Contábil e Pedagógica para Economia Solidária

--------------------------------------------------------------------------------

#### 1. Requisitos Funcionais (RF)

##### RF-01: PDV Operacional Soberano e Pedagógico
**Descrição:** Interface simplificada para registro de Vendas e Compras que atua, também, como ferramenta de educação financeira.
**Regra de Negócio:**
*   Cada venda deve gerar automaticamente uma "Entry" e dois "Postings" (Débito/Crédito) no Ledger (Partidas dobradas invisíveis ao usuário).
*   **Sócio-Técnica:** O PDV deve auxiliar visualmente o empreendedor na formação de preço, demonstrando de forma gráfica o custo do insumo versus o valor da hora trabalhada.
**Prioridade:** Essencial (v0)

##### RF-02: Registro de Trabalho (ITG 2002)
**Descrição:** Captura de horas/minutos trabalhados por membros do empreendimento.
**Regra de Negócio:**
*   Tempo deve ser registrado em int64 (minutos).
*   Tempo registrado constitui Capital Social de Trabalho (Primazia do Trabalho).
*   Base primária para rateio de sobras sociais.
**Prioridade:** Essencial (v0)

##### RF-03: Motor de Reservas Obrigatórias e Transparência Algorítmica
**Descrição:** Segregação automática de fundos antes da distribuição de resultados, com output visual para a Assembleia.
**Regra de Negócio:**
*   Bloqueio mandatório de **10%** para Reserva Legal e **5%** para FATES (Fundo de Assistência Técnica).
*   **Sócio-Técnica:** O cálculo deve gerar um painel gráfico simples e didático explicando a divisão, para ser lido, compreendido e aprovado em Assembleia Geral pelos trabalhadores.
**Prioridade:** Essencial (v0)

##### RF-04: Dossiê de Formalização Gradual (CADSOL)
**Descrição:** Exportação de Atas de Assembleia e Relatórios de Impacto Social de forma progressiva.
**Regra de Negócio:**
*   Respeitar o tempo de maturação do grupo informal (DREAM); o sistema não deve forçar a formalização.
*   Documentos gerados em formato Markdown com hash SHA256 para integridade.
*   Critérios de formalização: mínimo 3 decisões registradas.
**Prioridade:** Essencial (v0)

##### RF-05: Sincronização Offline-First
**Descrição:** Operação sem internet com sincronização posterior, respeitando a realidade de conectividade de áreas rurais/periféricas.
**Regra de Negócio:**
*   Interface PWA deve permitir operações básicas offline.
*   Delta tracking para detectar alterações.
*   Sync apenas com dados agregados (privacidade).
**Prioridade:** Alta

##### RF-06: Intercooperação B2B
**Descrição:** Marketplace para troca de produtos entre cooperativas.
**Regra de Negócio:** Ofertas entre entidades registradas (dados agregados sem exposição de membros).
**Prioridade:** Média

##### RF-07 a RF-10: Gestão Complementar
*   **RF-07 Gestão de Compras:** Registro de aquisições gerando lançamento contábil (D: Despesa / C: Caixa).
*   **RF-08 Gestão de Estoque:** Entrada, saída e saldo mínimo.
*   **RF-09 Gestão de Caixa:** Entradas, saídas e conciliação simplificada.
*   **RF-10 Gestão Orçamentária:** Planejamento e acompanhamento visual.

--------------------------------------------------------------------------------

#### 2. Requisitos Não Funcionais (RNF)

##### RNF-01: Isolamento por Tenant (Soberania)
**Descrição:** Cada "Sonho" ou Cooperativa deve ter seu próprio arquivo .sqlite.
**Implementação:** Arquivos em `data/entities/{entity_id}.db` com isolamento físico total para garantir "Poder de Saída" (Exit Power) dos dados.

##### RNF-02: Rigor Monetário (Anti-Float)
**Descrição:** Proibição de tipos de ponto flutuante em todo o Core Lume.
**Regra:** Uso obrigatório de int64 para representar centavos e minutos.

##### RNF-03: Resiliência (Offline-First)
**Descrição:** A interface (PWA) deve operar sem internet (Cache First, Delta tracking local, Background sync).

##### RNF-04 a RNF-06: Performance, Escalabilidade e Segurança
*   **Performance:** Criação tenant < 500ms; Registro venda < 100ms.
*   **Escala:** Suporte a milhões de EES; Serpro como infraestrutura de nuvem soberana.
*   **Segurança:** Hash SHA256 para auditoria e assinatura digital de pacotes.

##### RNF-07: Usabilidade e Design Participativo (Novo)
**Descrição:** Interface acessível co-criada com usuários não técnicos e sem jargões contábeis.
**Metas:** 
*   Validação obrigatória das telas com grupos reais de EES e Incubadoras Tecnológicas (ITCPs).
*   Tempo de treinamento < 2 horas.
*   Acessibilidade (WCAG 2.1 AA) e Design mobile-first.
*   Suporte ativo a usuários com baixa literacia digital.

##### RNF-08: Conformidade Legal (Contabilidade Invisível)
**Descrição:** Adequação às normas processadas silenciosamente no backend.
**Requisitos:** ITG 2002, Lei Paul Singer e geração de dados para CADSOL/DCSOL.

##### RNF-09 e RNF-10: Manutenibilidade e Interoperabilidade
*   Arquitetura modular (Clean Architecture), cobertura de testes >80%, API REST e Webhooks.

--------------------------------------------------------------------------------

#### 3. Matriz de Rastreabilidade
| Operação | Gatilho | Impacto no Ledger | Requisito |
| ------ | ------ | ------ | ------ |
| Venda no Balcão | PDV Submit | D: Ativo / C: Receita | RF-01 |
| Registro de Compra | Entrada Mercadoria | D: Despesa / C: Fornecedor | RF-07 |
| Fim de Turno | Log Horas | Registro de Cota-Trabalho | RF-02 |
| Fechamento Mês | Batch Job | Cálculo de Reservas (15%) + Gráfico Social | RF-03 |
| Assembleia | Decisão | Hash em decisions_log | RF-04 |

--------------------------------------------------------------------------------

#### 4. Casos de Uso Principais
*   **UC-01:** Registro de Venda (Operação coloquial gerando partida dobrada)
*   **UC-02:** Registro de Trabalho (Conversão de Tempo em Capital)
*   **UC-03:** Rateio de Sobras (Aprovação visual em Assembleia)
*   **UC-04:** Formalização Gradual (DREAM -> FORMALIZED)

--------------------------------------------------------------------------------

#### 5. Glossário de Termos
| Termo | Definição |
| ------ | ------ |
| **Tenant** | Entidade (cooperativa/grupo) com banco isolado |
| **Entry / Posting** | Lançamento contábil e Partidas dobradas (invisíveis na UI) |
| **ITG 2002** | Norma de contabilidade para economia solidária |
| **FATES** | Fundo de Assistência Técnica e Extensão Rural |
| **Reserva Legal** | Fundo obrigatório de 10% das sobras |
| **Design Participativo** | Construção do software *com* o usuário final |
```

***