***

```markdown
#### title: Constituição de IA e Agentes
status: implemented
version: 1.3
last_updated: 2026-03-08

### Constituição de IA e Agentes - Digna

#### 1. Contexto do Projeto

##### 1.1 Visão Geral
**Projeto:** Sistema de Gestão Contábil Público, Open Source e Participativo para Economia Solidária.
**Propósito:** Transformar a contabilidade de um peso burocrático em uma ferramenta pedagógica e de **Dignidade Financeira** para Empreendimentos de Economia Solidária (EES), atuando também como ponte tecnológica para a classe contábil (Contadores Sociais).
**Diferencial:** Arquitetura *Local-First Server-Side* com isolamento físico de dados, unida a uma interface de Tecnologia Social construída participativamente.

##### 1.2 Componentes Chave
* **Lume (Motor):** O backend em Go (rigoroso, exato e com contabilidade invisível ao usuário).
* **Digna (App):** A interface (pedagógica, construída com linguagem popular e acessível).
* **Accountant Dashboard:** A interface multi-tenant (modo leitura) para o Contador Social validar a conformidade e exportar obrigações fiscais.
* **Providentia (Fundação):** A entidade de governança sociotécnica (híbrida entre desenvolvedores, contadores e movimentos sociais).

--------------------------------------------------------------------------------

#### 2. Regras de Ouro (Não Negociáveis)

##### Regra 1: Integridade Financeira
* **PROIBIDO** o uso de `float32` ou `float64` em qualquer cálculo.
* **OBRIGATÓRIO** o uso de `int64` (centavos e minutos de trabalho).
* Validação de **Soma Zero** em todos os lançamentos do Ledger.

##### Regra 2: Adequação Sociotécnica e Linguagem
* **PROIBIDO** vazar jargões contábeis (Débito, Crédito, Provisão, Partidas Dobradas) para a interface do usuário (pdv_ui).
* **OBRIGATÓRIO** o uso de linguagem popular focada na ação real do dia a dia (ex: "Vendi um produto", "Trabalhei 4 horas", "Dinheiro no Caixa").

##### Regra 3: Nomenclatura de Arquivos
* **PROIBIDO** o uso de espaços em nomes de diretórios ou arquivos.
* **PADRÃO:** `kebab-case` para diretórios e `snake_case` para arquivos `.go`.

##### Regra 4: Isolamento de Dados (Soberania)
* O acesso ao arquivo `.sqlite` é exclusivo via `LifecycleManager`.
* Cada Tenant deve ser isolado fisicamente em `data/entities/{entity_id}.db`. Nunca faça *JOINs* entre entidades.

##### Regra 5: Blindagem do Core Lume (Separação Gerencial/Fiscal) [NOVO]
* **PROIBIDO** codificar calculadoras de impostos complexas (ICMS, Simples Nacional, etc.) ou regras fiscais dentro do motor central.
* **OBRIGATÓRIO** delegar as obrigações estatais à exportação de dados estruturados (SPED/Lotes Fiscais) para softwares comerciais externos utilizados pelos Contadores.

--------------------------------------------------------------------------------

#### 3. Regras para Agentes

##### 3.1 Princípios de Operação
1. **Soberania de Dados:** O dado pertence à entidade. O acesso deve ser isolado por arquivo.
2. **Padrão de Lançamento:** Seguir rigorosamente o princípio de partidas dobradas (soma zero) no backend, mantendo a tradução invisível no frontend.
3. **Transição Respeitosa:** O sistema deve prever a transição gradual e pedagógica do grupo informal (`DREAM`) para o formal (`FORMALIZED`) via interfaces (Facade), respeitando o tempo de maturação do coletivo.

##### 3.2 Regras de Execução
1. **Contexto e Empatia:** Leia sempre `docs/03_architecture/01_system.md` antes de sugerir mudanças. Ao propor UIs, coloque-se no lugar de um usuário com baixa literacia digital.
2. **Espaço:** Nunca use espaços em caminhos de arquivos.
3. **Finanças:** Se vir um `float` no código contábil, pare imediatamente e corrija para `int64`.
4. **Isolamento Read-Only:** Quando atuar no módulo do contador (`accountant_dashboard`), certifique-se de que as consultas aos micro-databases operem estritamente em modo de leitura.
5. **Criação de Arquivos:** Código do módulo lume não deve criar arquivos físicos no disco de forma avulsa; peça essa orquestração ao `lifecycle`.

--------------------------------------------------------------------------------

#### 4. Ferramental Recomendado
* **IDE:** VS Code com extensão oficial `golang.go`.
* **Workspace:** `go.work` para orquestrar os múltiplos módulos.
* **Formatação:** Comando `gofmt` obrigatório antes de qualquer commit.

--------------------------------------------------------------------------------

#### 5. Estratégia de Módulos (Mitigação de Entropia da IA)
Para evitar que o agente perca o foco, atue sempre por módulos fechados:

| Módulo | Consome de | Entrega para | Foco Central |
| ------ | ------ | ------ | ------ |
| **Lifecycle** | Sistema de Arquivos | Ledger / Legal Facade | Criação do `.db` físico |
| **Core Lume** | Lifecycle | Reporting | Soma zero e `int64` |
| **PDV UI** | Core Lume | Usuário Final | Linguagem popular e HTMX |
| **Legal Facade**| Lifecycle | CADSOL | Hashes SHA256 e Markdown |
| **Accountant**| Lifecycle (Read-Only)| Fisco / Contador | Exportação SPED |

--------------------------------------------------------------------------------

#### 6. Padrão de Sessão (Protocolo do Agente)

##### 6.1 Início de Sessão
1. Ler documentação relevante (`01_system.md` e `01_requirements.md`).
2. Verificar status atual das sprints (`04_status.md`).
3. Entender dependências atômicas.

##### 6.2 Durante a Sessão
1. Implementar funcionalidade e adequar linguagem à realidade local.
2. Executar testes unitários (TDD).
3. Validar contra os requisitos técnicos e sociais desta constituição.

##### 6.3 Fim de Sessão
1. Atualizar o `05_session_log.md`.
2. Registrar decisões técnicas.
3. Documentar eventuais "Brechas/Gaps Arquiteturais Identificados".

--------------------------------------------------------------------------------

#### 7. Glossário de Termos

| Termo | Significado |
| ------ | ------ |
| **Tenant** | Entidade (cooperativa/grupo informal) com banco SQLite isolado. |
| **Ledger** | Livro contábil com partidas dobradas gerido silenciosamente no backend. |
| **ITG 2002** | Norma de contabilidade para economia solidária emitida pelo CFC. |
| **CADSOL** | Sistema governamental (Cadastro Nacional de Economia Solidária). |
| **PDV** | Ponto de Venda (A interface pedagógica front-end do produtor). |
| **int64** | Tipo inteiro de 64 bits (exclusivo para dinheiro e controle de horas). |
| **Tecnologia Social**| Ferramenta construída *com* o usuário para a sua própria emancipação. |
| **Contador Social** | Profissional parceiro que audita a ITG 2002 e gera obrigações sem atuar como "digitador de notas". |
| **SPED / Lote Fiscal**| Formato de extração de dados gerado pelo Digna para cumprimento de obrigações legais em sistemas externos. |
```

***
