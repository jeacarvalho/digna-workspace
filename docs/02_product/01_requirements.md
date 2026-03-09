***

```markdown
#### title: Requisitos do Projeto Digna
status: implemented
version: 1.3
last_updated: 2026-03-08

### Requisitos - Projeto Digna
**Referência Legal:** Lei nº 15.068/2024 (Lei Paul Singer) / ITG 2002 (CFC)
**Projeto:** Digna - Infraestrutura Contábil e Pedagógica para Economia Solidária

--------------------------------------------------------------------------------

#### 1. Requisitos Funcionais (RF)

##### RF-01: PDV Operacional Soberano e Pedagógico
**Descrição:** Interface simplificada para registro de Vendas e Compras que atua, também, como ferramenta de educação financeira.
**Regra de Negócio:**
* Cada venda deve gerar automaticamente uma "Entry" e dois "Postings" (Débito/Crédito) no Ledger (Partidas dobradas invisíveis ao usuário).
* **Sócio-Técnica:** O PDV deve auxiliar visualmente o empreendedor na formação de preço, demonstrando de forma gráfica o custo do insumo versus o valor da hora trabalhada.
**Prioridade:** Essencial (v0)

##### RF-02: Registro de Trabalho (ITG 2002)
**Descrição:** Captura de horas/minutos trabalhados por membros do empreendimento.
**Regra de Negócio:**
* Tempo deve ser registrado estritamente em `int64` (minutos).
* O tempo registrado constitui o Capital Social de Trabalho, servindo de lastro.
**Prioridade:** Essencial (v0)

##### RF-03: Motor de Reservas Obrigatórias e Rateio
**Descrição:** Segregação automática de fundos antes da distribuição de resultados (Sobras).
**Regra de Negócio:**
* Bloqueio mandatório de **10%** para Reserva Legal e **5%** para FATES (Fundo de Assistência Técnica).
* Rateio do saldo restante deve ser proporcional aos minutos trabalhados (RF-02).
* **Sócio-Técnica:** O cálculo deve gerar um painel gráfico simples e didático explicando a divisão, para ser lido, compreendido e aprovado em Assembleia Geral pelos trabalhadores.
**Prioridade:** Essencial (v0)

##### RF-04: Dossiê de Formalização Gradual (CADSOL)
**Descrição:** Exportação de Atas de Assembleia e Relatórios de Impacto Social de forma progressiva.
**Regra de Negócio:**
* Respeitar o tempo de maturação do grupo informal (DREAM); o sistema não deve forçar a formalização.
* Documentos gerados em formato Markdown com hash SHA256 para integridade.
* Critérios de formalização: mínimo 3 decisões registradas.
**Prioridade:** Essencial (v0)

##### RF-05: Sincronização Offline-First
**Descrição:** Operação sem internet com sincronização posterior, respeitando a realidade de conectividade de áreas rurais/periféricas.
**Regra de Negócio:**
* Interface PWA deve permitir operações básicas offline.
* Delta tracking para detectar alterações.
* Sync apenas com dados agregados (privacidade).
**Prioridade:** Alta

##### RF-06: Intercooperação B2B
**Descrição:** Marketplace para troca de produtos entre cooperativas.
**Regra de Negócio:**
* Ofertas entre entidades registradas mantendo o isolamento dos dados.
**Prioridade:** Média

##### RF-11: Painel do Contador Social e Exportação Fiscal [NOVO]
**Descrição:** Interface web segregada (Multi-tenant) para contadores parceiros auditarem a conformidade e exportarem obrigações acessórias.
**Regra de Negócio:**
* O contador deve conseguir visualizar o status de fechamento de múltiplos empreendimentos solidários em uma única tela.
* O sistema **não calcula guias de impostos**, mas agrega as partidas dobradas geradas pelo `Core Lume` e as exporta em leiaute padrão (SPED ou importação para sistemas contábeis comerciais).
* Auditoria visual automática da norma ITG 2002 e dos bloqueios de FATES/Reserva Legal.
**Prioridade:** Alta (Fase 1)

--------------------------------------------------------------------------------

#### 2. Requisitos Não Funcionais (RNF)

##### RNF-01: Isolamento por Tenant (Soberania)
**Descrição:** Cada "Sonho" ou Cooperativa deve ter seu próprio arquivo `.sqlite`.
**Implementação:** Arquivos em `data/entities/{entity_id}.db` com isolamento físico total para garantir "Poder de Saída" (Exit Power) dos dados. Nenhum dado transacional de um tenant pode ser misturado com outro.

##### RNF-02: Rigor Monetário (Anti-Float)
**Descrição:** Proibição de tipos de ponto flutuante em todo o Core Lume.
**Regra:** Uso obrigatório de `int64` para representar centavos e minutos. Qualquer variável `float` no backend contábil causará falha na esteira de integração (CI).

##### RNF-03: Resiliência (Offline-First)
**Descrição:** A interface (PWA) deve operar sem internet (Cache First, Delta tracking local, Background sync).

##### RNF-04 a RNF-06: Performance, Escalabilidade e Segurança
* **Performance:** Criação tenant < 500ms; Registro venda < 100ms.
* **Escala:** Suporte a milhões de EES; Serpro como infraestrutura de nuvem soberana.
* **Segurança:** Hash SHA256 para auditoria e assinatura digital de pacotes.

##### RNF-07: Usabilidade e Design Participativo
**Descrição:** Interface acessível co-criada com usuários não técnicos e sem jargões contábeis.
**Metas:**
* Validação obrigatória das telas com grupos reais de EES e Incubadoras Tecnológicas (ITCPs).
* Tempo de treinamento < 2 horas.
* Acessibilidade (WCAG 2.1 AA) e Design mobile-first.
* Suporte ativo a usuários com baixa literacia digital.

##### RNF-08: Conformidade Legal (Contabilidade Invisível)
**Descrição:** Adequação às normas processadas silenciosamente no backend, criando a ponte para a classe contábil.
**Requisitos:** 
* Conformidade rigorosa com a ITG 2002 do Conselho Federal de Contabilidade (CFC).
* Conformidade com a Lei Paul Singer (15.068/2024).
* Geração de dados e hashes institucionais para o CADSOL/SINAES.

##### RNF-09 e RNF-10: Manutenibilidade e Interoperabilidade
* **Manutenibilidade:** Arquitetura modular (Clean Architecture) e cobertura de testes >80%.
* **Interoperabilidade:** Geração de arquivos fiscais estruturados (CSV/TXT/SPED) para consumo em softwares de escrituração contábil externos.

--------------------------------------------------------------------------------

#### 3. Matriz de Rastreabilidade

| Operação | Gatilho | Impacto no Ledger | Requisito |
| ------ | ------ | ------ | ------ |
| Venda no Balcão | PDV Submit | D: Ativo / C: Receita | RF-01 |
| Registro de Compra | Entrada Mercadoria | D: Despesa / C: Fornecedor | RF-07 |
| Fim de Turno | Log Horas | Registro de Cota-Trabalho | RF-02 |
| Fechamento Mês | Batch Job | Cálculo de Reservas (15%) + Gráfico Social | RF-03 |
| Assembleia | Decisão | Hash em decisions_log | RF-04 |
| Fechamento Fiscal | Painel do Contador | Geração de SPED/Lote (Apenas Leitura) | RF-11 |

--------------------------------------------------------------------------------

#### 4. Casos de Uso Principais

* **UC-01:** Registro de Venda (Operação coloquial gerando partida dobrada).
* **UC-02:** Registro de Trabalho (Conversão de Tempo em Capital).
* **UC-03:** Rateio de Sobras (Aprovação visual e democrática em Assembleia).
* **UC-04:** Formalização Gradual (DREAM -> FORMALIZED via aprovação de atas).
* **UC-05:** Auditoria e Exportação Contábil (Contador Social acessa painel multi-tenant, valida a ITG 2002 e exporta lote fiscal para a Receita/Sistemas terceiros).

--------------------------------------------------------------------------------

#### 5. Glossário de Termos

| Termo | Definição |
| ------ | ------ |
| **Tenant** | Entidade (cooperativa/grupo) com banco SQLite isolado. |
| **Entry / Posting** | Lançamento contábil e Partidas dobradas (invisíveis na UI humana). |
| **ITG 2002** | Norma de contabilidade para economia solidária emitida pelo CFC. |
| **FATES** | Fundo de Assistência Técnica e Extensão Rural. |
| **Reserva Legal** | Fundo obrigatório de 10% das sobras. |
| **Design Participativo**| Construção do software *com* o usuário final e entidades de apoio. |
| **Contador Social** | Profissional parceiro que audita o EES via dashboard, sem atuar como "digitador". |
| **SPED / Lote Fiscal** | Formato de arquivo exportado pelo Digna para que o contador faça o recolhimento de tributos. |
```
