
#### title: Modelos de Domínio e Dados
status: implemented
version: 1.3
last_updated: 2026-03-08

### Modelos - Projeto Digna
**Projeto:** Sistema de Gestão Contábil e Pedagógica para Economia Solidária

--------------------------------------------------------------------------------

#### 1. Domain Model (Modelo de Domínio)

O domínio do Digna reflete os princípios da autogestão e da contabilidade invisível, priorizando as relações humanas sobre o capital financeiro, mas agora atuando também como **ponte institucional** para a conformidade legal.

##### 1.1 Entidades Principais

###### Enterprise (Empreendimento de Economia Solidária - EES)
Representa o coletivo produtivo. Pode transitar gradualmente por três estados, respeitando o tempo político do grupo:
*   **DREAM (Sonho):** Grupo informal, focado na união produtiva inicial.
*   **INCUBATED (Incubado):** Em processo de estruturação, recebendo apoio pedagógico (ITCPs, ONGs).
*   **FORMALIZED (Formalizado):** Cooperativa ou Associação com CNPJ e estatuto base (Pronto para CADSOL e obrigações fiscais).

###### Member (Trabalhador/Cooperado)
Pessoa participante do empreendimento. Suas horas dedicadas são o lastro do capital social (Princípio da Primazia do Trabalho).

###### SocialAccountant (Contador Social) [NOVO]
Entidade parceira (externa). Não é dona do dado, mas possui permissão delegada de leitura para auditar a conformidade (ITG 2002) e extrair os balancetes através de um Painel Multi-tenant.

###### EnterpriseAccountant (Vínculo Contábil e Responsabilidade Técnica) [NOVO]**
Entidade de relacionamento temporal que consolida o **RF-12**. Garante o "Exit Power" da cooperativa e protege a responsabilidade do contador.
*   `ID`: Identificador único (UUID).
*   `EnterpriseID`: Referência ao Empreendimento (Tenant).
*   `AccountantID`: Referência ao Contador Social.
*   `Status`: `ACTIVE` (Apenas 1 permitido no tempo presente) ou `INACTIVE`.
*   `StartDate`: Data de início da responsabilidade técnica.
*   `EndDate`: Data do encerramento (Nulo se for o contador atual).
*   `DelegatedBy`: ID do membro coordenador que aprovou o vínculo (Auditoria).

###### Transaction (Operação Comercial)
Evento econômico do dia a dia (venda na feira, compra de insumo). Traduzido internamente para partidas dobradas.

###### WorkLog (Registro ITG 2002)
Registro de trabalho cooperativo. Converte o suor em Capital Social de Trabalho (mensurado em minutos).

###### Decision (Decisão Democrática)
Decisão coletiva tomada e registrada em Assembleia. Base para a geração das Atas em Markdown (CADSOL).

###### LegalDocument (Dossiê CADSOL)**
Agrupamento das `Decisions` exportado em formato Markdown/PDF para fins legais, comprovando a gestão democrática perante o Estado.


###### Fund (Fundos Obrigatórios)
Reservas estatutárias e legais blindadas pelo sistema (Ex: Reserva Legal e FATES).

###### FiscalBatch (Lote Fiscal) [NOVO]
Conjunto imutável de transações agregadas e exportadas para o formato exigido pela Receita Federal (ex: SPED) ou softwares de escrituração contábil externos.

##### 1.2 Value Objects

| Value Object | Formato Técnico | Justificativa Sociotécnica |
| ------ | ------ | ------ |
| **Money** | `int64` (centavos) | Evita erros de arredondamento capitalista (IEEE 754). Garante exatidão total para o trabalhador. |
| **Time/Labor** | `int64` (minutos) | Unidade de medida do Capital Social. |
| **AccountCode**| `string` (ex: 1.1.01)| Padronização invisível ao usuário, usada para tradução fiscal. |
| **Period** | `YYYY-MM` | Ciclo contábil e de prestação de contas. |

--------------------------------------------------------------------------------

#### 2. Data Model (Schema v1)

O banco de dados é instanciado fisicamente de forma isolada por `Enterprise` (Soberania do Dado local). 
*Nota Arquitetural:* O Painel do Contador Social **não possui** banco de transações próprio; ele atua apenas consumindo dados em modo de leitura (Read-Only) dos micro-databases autorizados. E como um contador social pode estar cuidando de várias entidades, as funcionalidades a ele relacionadas devem estar preparadas para pesquisar em diversos db sqlite

    TABELAS PRINCIPAIS (SQLite)
    
    [ accounts ]          --> Plano de contas hierárquico (Padrão ITG 2002)
    [ entries ]           --> Lançamentos contábeis (O evento principal)
    [ postings ]          --> Partidas dobradas (Débito e Crédito associados à Entry)
    [ work_logs ]         --> Tabela de valoração social (Registro de minutos trabalhados)
    [ decisions_log ]     --> Registro de governança em Assembleia (Gera a Ata)
    [ sync_metadata ]     --> Delta tracking para resiliência Offline-First
    [ fiscal_exports ]    --> [NOVO] Log de extrações SPED realizadas pelo Contador Social (evita envios duplicados)

--------------------------------------------------------------------------------

#### 3. Algoritmos de Negócio e Governança

##### 3.1 Algoritmo de Rateio Social (Transparência Algorítmica Visual)
**Objetivo:** Distribuir o excedente financeiro de forma justa, baseada na Primazia do Trabalho.
**Regra Sociotécnica:** O cálculo não deve ser obscuro. O algoritmo DEVE emitir uma saída (gráfico/tabela) didática para ser lida em Assembleia Geral.
**Entrada:** `totalSurplus` (int64), `memberHours` (Mapa de member_id -> minutos).

    Exemplo Didático Gerado pelo Algoritmo para a Assembleia:
    Sobras Totais do Mês: R$ 100,00
    --------------------------------------------------------------
    Trabalhador   | Suor (Minutos) | % do Total | Valor a Receber
    --------------------------------------------------------------
    Maria (001)   | 600 min        | 66.67%     | R$ 66,66
    João  (002)   | 300 min        | 33.33%     | R$ 33,33
    --------------------------------------------------------------

##### 3.2 Algoritmo de Reservas Obrigatórias (Segregação de Fundos)
**Objetivo:** Garantir a conformidade legal (Lei Paul Singer) e a sustentabilidade de longo prazo antes de qualquer rateio.
**Processo:** 
1. Apura o resultado positivo do período.
2. Bloqueia 10% para o Fundo de Reserva Legal.
3. Bloqueia 5% para o FATES (Assistência Técnica Educacional).
4. Libera o saldo (85%) para a lógica do Algoritmo 3.1.

##### 3.3 Algoritmo de Partidas Dobradas Invisíveis
**Objetivo:** Validar a integridade financeira (Soma Zero).
**Processo:** `soma(débitos) + soma(créditos) == 0`. Acionado silenciosamente no backend a cada ação comercial (venda/compra) feita no PDV.

##### 3.4 Algoritmo de Formalização Gradual
**Objetivo:** Avaliar a maturidade institucional para permitir a transição `DREAM` -> `FORMALIZED`.
**Critérios Automatizados:**
*   Mínimo de 3 registros de `Decision` (Assembleias realizadas provando autogestão).
*   Mínimo de 1 membro ativo com histórico de `WorkLog`.
*   Criação automática do Dossiê Hash SHA256.

##### 3.5 Algoritmo de Tradução Fiscal (Ponte do Contador) [NOVO]
**Objetivo:** Converter a contabilidade social gerada invisivelmente pelo produtor em linguagem de conformidade estatal (SPED/Lotes Fiscais).
**Processo:**
1. Painel do Contador solicita dados de um `Period` fechado.
2. O algoritmo compila todas as `entries` de soma zero.
3. Mapeia as contas locais amigáveis ("Gaveta") para o Plano de Contas Referencial da Receita Federal.
4. Gera o pacote CSV/SPED e salva o evento de exportação na tabela `fiscal_exports`.

--------------------------------------------------------------------------------

#### 4. Seed Data (Carga Inicial Padrão)

Toda nova base SQLite de um EES nasce com este plano de contas enxuto e adaptado:

| ID | Código | Nome Amigável | Natureza Contábil (Invisível) | Mapeamento Fiscal [NOVO] |
| -- | ------ | ------------- | ----------------------------- | ------------------------ |
| 1  | 1.1.01 | Gaveta / Caixa| ASSET (Ativo)                 | Disponibilidades (Ativo) |
| 2  | 3.1.01 | Nossas Vendas | REVENUE (Receita)             | Receita Bruta            |
| 3  | 1.1.02 | Banco / Conta | ASSET (Ativo)                 | Contas Bancárias         |
| 4  | 2.1.01 | Quem Fornece  | LIABILITY (Passivo)           | Fornecedores a Pagar     |
| 5  | 3.2.01 | Fundo FATES   | EQUITY (Patrimônio Líquido)   | Reservas Estatutárias    |
```

***
