# 📋 Prompt para RF-30: Sistema de Ajuda Educativa Estruturada

**Tipo:** Feature + Infraestrutura  
**Módulo:** `modules/ui_web`, `modules/core_lume` (novo: `help_engine`)  
**Prioridade:** ALTA (Habilita adoção por baixa escolaridade - Pilar Pedagógico)  
**Estimativa:** 12-16 horas  

---

## 🎯 CONTEXTO DA TAREFA

Você vai implementar o **Sistema de Ajuda Educativa Estruturada** no sistema Digna, criando uma base de conhecimento acessível que traduz conceitos técnicos (CadÚnico, inadimplência, CNAE, etc.) em linguagem popular, com linkagem direta entre elementos de UI e registros de ajuda no banco.

**Princípio Central (Pilar 2 + Pilar 5):** *"Nenhum usuário deve se sentir humilhado por não entender um termo. O sistema ensina enquanto é operado."*

**Motivação (Decisão da Sessão 27/03/2026):** Campos como "Inscrito no CadÚnico?", "Inadimplência Ativa?", "CNAE" são jargões burocráticos que violam o Pilar Pedagógico do Digna. Tooltips simples são insuficientes — precisamos de páginas educativas completas.

---

## 📝 DESCRIÇÃO DETALHADA

### Problema a Resolver
Usuários do Digna têm baixa escolaridade. Conteúdo técnico causa abandono e frustração. O sistema precisa ensinar enquanto é operado, sem humilhar o usuário por não entender termos técnicos.

### Escopo da Implementação
1. **Base de Dados de Ajuda:** Tabela `help_topics` no `central.db` com conteúdo estruturado
2. **Handler de Ajuda:** Busca, índice, visualização de tópicos
3. **Integração UI:** Botão `?` em campos técnicos linka para tópico específico
4. **Menu de Ajuda:** Entrada acessível em todas as páginas com busca e índice
5. **Seed de Tópicos:** 6 tópicos obrigatórios iniciais (CadÚnico, Inadimplência, CNAE, DAS MEI, Reserva Legal, FATES)

### Fora do Escopo (MVP)
- Interface administrativa para editar tópicos (conteúdo hardcoded no seed)
- Busca full-text avançada (busca simples por título/chave)
- Múltiplos idiomas (apenas PT-BR inicial)
- Upload de anexos/mídia (apenas texto + links)

---

## 🏗️ ARQUITETURA TÉCNICA

### Estrutura de Output Esperada

```
modules/
├── core_lume/
│   ├── internal/
│   │   ├── domain/
│   │   │   └── help_topic.go           # NOVO: Entidade HelpTopic
│   │   ├── service/
│   │   │   └── help_service.go         # NOVO: Lógica de busca/índice
│   │   └── repository/
│   │       └── help_repository.go      # NOVO: Persistência SQLite (central.db)
│   └── pkg/
│       └── help/
│           └── help.go                 # NOVO: API pública
│
└── ui_web/
    ├── internal/handler/
    │   └── help_handler.go             # NOVO: Handler HTTP
    └── templates/
        ├── help_index_simple.html      # NOVO: Índice/busca
        └── help_topic_simple.html      # NOVO: Visualização de tópico
```

### Integrações Existentes
- **LifecycleManager:** Para acesso ao `central.db` (dados globais)
- **BaseHandler:** Para renderização de templates
- **Navegação:** Link "Ajuda" no header de todos os templates `*_simple.html`

---

## 🛠️ TAREFAS DE IMPLEMENTAÇÃO

### 1. Domain Layer (`core_lume/internal/domain/help_topic.go`)

```go
type HelpTopic struct {
    ID           string // UUID
    Key          string // Chave única (ex: "cadunico", "inadimplencia")
    Title        string // Título em linguagem popular
    Summary      string // Resumo em 1 frase (para tooltips)
    Explanation  string // Explicação completa em linguagem popular
    WhyAsked     string // "Por que perguntamos isso?"
    Legislation  string // Legislação relacionada
    NextSteps    string // Próximos passos acionáveis
    OfficialLink string // Link para fonte oficial (ex: gov.br)
    Category     string // Categoria: CREDITO, TRIBUTARIO, GOVERNANCA, GERAL
    Tags         string // Tags para busca (JSON array ou comma-separated)
    
    // Metadados
    ViewCount    int64  // Quantas vezes foi visualizado
    CreatedAt    int64  // Unix timestamp
    UpdatedAt    int64  // Unix timestamp
}

// Categorias de tópicos
const (
    CategoriaCredito    = "CREDITO"
    CategoriaTributario = "TRIBUTARIO"
    CategoriaGovernanca = "GOVERNANCA"
    CategoriaGeral      = "GERAL"
)

// Métodos de domínio
func (h *HelpTopic) Validate() error
func (h *HelpTopic) IsComplete() bool // Todos os campos obrigatórios preenchidos
```

### 2. Repository Layer (`core_lume/internal/repository/help_repository.go`)

```go
type HelpRepository interface {
    Save(topic *HelpTopic) error
    FindByKey(key string) (*HelpTopic, error)
    FindByID(id string) (*HelpTopic, error)
    ListByCategory(category string) ([]*HelpTopic, error)
    Search(query string) ([]*HelpTopic, error) // Busca simples por título/tags
    ListAll() ([]*HelpTopic, error)            // Para índice
    IncrementViewCount(id string) error
}

// Tabela SQLite (central.db):
// CREATE TABLE IF NOT EXISTS help_topics (
//     id TEXT PRIMARY KEY,
//     key TEXT NOT NULL UNIQUE,
//     title TEXT NOT NULL,
//     summary TEXT,
//     explanation TEXT NOT NULL,
//     why_asked TEXT,
//     legislation TEXT,
//     next_steps TEXT,
//     official_link TEXT,
//     category TEXT NOT NULL,
//     tags TEXT,
//     view_count INTEGER DEFAULT 0,
//     created_at INTEGER,
//     updated_at INTEGER
// )
```

### 3. Service Layer (`core_lume/internal/service/help_service.go`)

```go
type HelpService struct {
    repo HelpRepository
}

// Métodos obrigatórios:
func (s *HelpService) GetTopicByKey(key string) (*HelpTopic, error)
func (s *HelpService) GetTopicByID(id string) (*HelpTopic, error)
func (s *HelpService) ListIndex() ([]*HelpTopic, error) // Agrupado por categoria
func (s *HelpService) Search(query string) ([]*HelpTopic, error)
func (s *HelpService) GetRelatedTopics(topic *HelpTopic) ([]*HelpTopic, error)
func (s *HelpService) IncrementView(key string) error
```

### 4. Handler Layer (`modules/ui_web/internal/handler/help_handler.go`)

```go
type HelpHandler struct {
    *BaseHandler
    lifecycleManager lifecycle.LifecycleManager
    helpService      *service.HelpService
}

// Rotas obrigatórias:
// GET /help                    → Índice de tópicos (categorizado)
// GET /help/search?q={query}   → Busca de tópicos
// GET /help/topic/{key}        → Visualização de tópico específico
// GET /help/tooltip/{key}      → Fragmento HTMX para tooltip (opcional)
```

### 5. Templates

#### `help_index_simple.html` (Índice + Busca)
- Seguir padrão `*_simple.html` (documento HTML completo)
- Carregar via `ParseFiles()` no handler (cache-proof)
- Usar paleta "Soberania e Suor"
- Incluir navegação padrão (header/footer)
- Exibir:
  - Barra de busca
  - Tópicos categorizados (CRÉDITO, TRIBUTÁRIO, GOVERNANÇA, GERAL)
  - Link para cada tópico

#### `help_topic_simple.html` (Visualização de Tópico)
- Seguir padrão `*_simple.html` (documento HTML completo)
- Carregar via `ParseFiles()` no handler (cache-proof)
- Usar paleta "Soberania e Suor"
- Incluir navegação padrão (header/footer)
- Exibir:
  - Título em linguagem popular
  - Resumo em 1 frase
  - Explicação completa
  - "Por que perguntamos isso?"
  - Legislação relacionada
  - Próximos passos acionáveis
  - Link para fonte oficial
  - Botão "Voltar para Ajuda"

### 6. Seed de Tópicos Iniciais (`core_lume/internal/repository/help_seed.go`)

```go
var initialHelpTopics = []HelpTopic{
    {
        Key:         "cadunico",
        Title:       "O que é o CadÚnico?",
        Summary:     "É o cadastro do governo para programas sociais.",
        Explanation: "O Cadastro Único (CadÚnico) reúne informações sobre famílias de baixa renda. Estar inscrito permite acesso a programas como Bolsa Família, Tarifa Social de Energia e linhas de crédito especiais.",
        WhyAsked:    "No Digna, informamos isso para encontrar programas de crédito que só atendem quem está no CadÚnico, como o 'Acredita no Primeiro Passo'.",
        Legislation: "Decreto nº 6.135/2007",
        NextSteps:   "Se não está inscrito, procure o CRAS (Centro de Referência de Assistência Social) do seu município com documentos pessoais e comprovante de residência.",
        OfficialLink: "https://www.gov.br/cadunico",
        Category:    "CREDITO",
        Tags:        "cadastro,programa social,crédito",
    },
    {
        Key:         "inadimplencia",
        Title:       "O que é inadimplência?",
        Summary:     "É quando há dívidas não pagas registradas.",
        Explanation: "Inadimplência significa que você tem contas atrasadas com bancos, lojas ou com o governo. Isso pode aparecer em sistemas como Serasa, SPC ou Dívida Ativa da União.",
        WhyAsked:    "Alguns programas de crédito exigem que você regularize essas dívidas antes de aplicar. Outros podem ajudar você a renegociar.",
        Legislation: "Lei nº 10.820/2003 (Descontos em Folha)",
        NextSteps:   "Se tem dívidas, o Digna pode ajudar a identificar programas de renegociação como o 'Desenrola Pequenos Negócios' antes de buscar crédito novo.",
        OfficialLink: "https://www.gov.br/economia",
        Category:    "CREDITO",
        Tags:        "dívida,renegociação,crédito",
    },
    {
        Key:         "cnae",
        Title:       "O que é CNAE?",
        Summary:     "É o código que diz qual é a atividade do seu negócio.",
        Explanation: "CNAE significa Classificação Nacional de Atividades Econômicas. É um número que o governo usa para saber se você vende comida, faz costura, presta serviço, etc.",
        WhyAsked:    "Programas de crédito usam o CNAE para saber se seu negócio se enquadra nas regras deles.",
        Legislation: "Resolução CONCLA nº 1/2006",
        NextSteps:   "Se não sabe seu CNAE, consulte no cartão do CNPJ ou no site da Receita Federal.",
        OfficialLink: "https://www.gov.br/receitafederal",
        Category:    "TRIBUTARIO",
        Tags:        "atividade,cadastro,receita",
    },
    {
        Key:         "das_mei",
        Title:       "O que é o DAS MEI?",
        Summary:     "É o boleto mensal que o MEI paga.",
        Explanation: "O DAS (Documento de Arrecadação do Simples Nacional) é o imposto que o Microempreendedor Individual paga todo mês. O valor é 5% do salário mínimo + valores fixos de ICMS e ISS.",
        WhyAsked:    "O Digna calcula automaticamente o valor do DAS e avisa quando está perto de vencer, para você não pagar multa.",
        Legislation: "Lei Complementar nº 123/2006",
        NextSteps:   "O Digna gera o cálculo automaticamente. Você só precisa pagar até o dia 20 de cada mês.",
        OfficialLink: "https://www.gov.br/empresas-e-negocios",
        Category:    "TRIBUTARIO",
        Tags:        "MEI,imposto,boleto",
    },
    {
        Key:         "reserva_legal",
        Title:       "O que é Reserva Legal?",
        Summary:     "É uma parte do lucro que a lei manda guardar.",
        Explanation: "A Reserva Legal é 10% do lucro da cooperativa que deve ser guardado por lei. Esse dinheiro não pode ser distribuído aos sócios — fica guardado para proteger a cooperativa em tempos difíceis.",
        WhyAsked:    "O Digna aplica automaticamente esse bloqueio antes de distribuir as sobras, para cumprir a lei e proteger o grupo.",
        Legislation: "Lei nº 5.764/71 (Lei Geral das Cooperativas)",
        NextSteps:   "Não precisa fazer nada — o Digna calcula e guarda automaticamente.",
        OfficialLink: "https://www.planalto.gov.br/ccivil_03/leis/l5764.htm",
        Category:    "GOVERNANCA",
        Tags:        "lucro,reserva,lei",
    },
    {
        Key:         "fates",
        Title:       "O que é o FATES?",
        Summary:     "É um fundo para ajudar outros grupos a se organizarem.",
        Explanation: "O FATES (Fundo de Assistência Técnica, Educacional e Social) é 5% do lucro da cooperativa que é separado para ajudar outras cooperativas e grupos a se organizarem. É uma forma de solidariedade entre grupos.",
        WhyAsked:    "O Digna aplica automaticamente esse bloqueio antes de distribuir as sobras, para cumprir a lei da Economia Solidária.",
        Legislation: "Lei nº 15.068/2024 (Lei Paul Singer)",
        NextSteps:   "Não precisa fazer nada — o Digna calcula e guarda automaticamente.",
        OfficialLink: "https://www.gov.br/trabalho-e-emprego",
        Category:    "GOVERNANCA",
        Tags:        "fundo,solidariedade,lei",
    },
}
```

---

## 📊 REGRAS DE NEGÓCIO DETALHADAS

### Categorias de Tópicos
| Categoria | Descrição | Exemplos |
|-----------|-----------|----------|
| `CREDITO` | Programas de financiamento, elegibilidade | CadÚnico, inadimplência, Pronampe |
| `TRIBUTARIO` | Impostos, obrigações fiscais | CNAE, DAS MEI, Simples Nacional |
| `GOVERNANCA` | Assembleias, decisões, formalização | Reserva Legal, FATES, CADSOL |
| `GERAL` | Conceitos gerais do sistema | Soberania de dados, contabilidade invisível |

### Linkagem UI → Tópico
Em templates, campos técnicos devem ter botão `?` linkando para `/help/topic/{key}`:

```html
<label class="flex items-center gap-2">
  Inscrito no CadÚnico?
  <a href="/help/topic/cadunico" target="_blank" 
     class="help-tooltip-trigger text-digna-primary hover:underline"
     aria-label="Saiba mais sobre CadÚnico">
    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.125 2.5-2.5 2.5-.69 0-1.25.28-1.5.75-.25.47-.25 1.03 0 1.5.25.47.81.75 1.5.75 1.375 0 2.5-1.1 2.5-2.5 0-1.657-1.79-3-4-3-1.742 0-3.223.835-3.772 2z"/>
      </svg>
  </a>
</label>
```

**Alternativa (HTMX Tooltip):**
```html
<button type="button" 
        hx-get="/help/tooltip/cadunico"
        hx-target="#tooltip-cadunico"
        hx-swap="outerHTML">
  ?
</button>
<div id="tooltip-cadunico"></div>
```

### Permissões de Acesso
- **Leitura:** Pública (não requer autenticação) — maximiza acesso à informação
- **Escrita/Admin:** Apenas usuários com role `ADMIN` ou `FOUNDATION` (futuro)

---

## ✅ CRITÉRIOS DE ACEITE (Definition of Done)

### Arquitetura
- [ ] Segue Clean Architecture (Domain → Service → Repository → Handler)
- [ ] Tabela `help_topics` no `central.db` (dados globais)
- [ ] Handler estende `BaseHandler`
- [ ] Templates `*_simple.html` carregados via `ParseFiles()` no handler

### Funcionalidade
- [ ] Índice de tópicos categorizado acessível em `/help`
- [ ] Busca funcional por título/tags
- [ ] Visualização de tópico em `/help/topic/{key}`
- [ ] Link "Ajuda" no header de todos os templates principais
- [ ] Botão `?` em campos técnicos linka para tópico específico
- [ ] Seed de 6 tópicos iniciais implementado

### Pedagogia (CRÍTICO - RF-30)
- [ ] **Linguagem para 5ª série:** Usuário com 5ª série consegue entender sem ajuda externa
- [ ] **Zero jargões:** Conteúdo não usa termos técnicos sem explicação ("cadastramento", "regularização fiscal", etc.)
- [ ] **Próximo passo acionável:** Sempre inclui ação concreta (ex: "procure o CRAS")
- [ ] **Tooltip carrega em < 500ms** via HTMX

### Testes
- [ ] Testes unitários para service (busca, listagem)
- [ ] Testes unitários para repository (CRUD)
- [ ] Testes de handler (HTTP responses)
- [ ] Cobertura >90% para service e handler
- [ ] Smoke test: `./scripts/dev/smoke_test_new_feature.sh "Central de Ajuda" "/help"`
- [ ] Validação E2E: `./scripts/dev/validate_e2e.sh --basic --headless`

### UI/UX
- [ ] Design segue "Soberania e Suor" (#2A5CAA, #4A7F3E, #F57F17)
- [ ] Navegação clara (voltar para índice, tópicos relacionados)
- [ ] Acessível via teclado e leitor de tela
- [ ] Responsivo (mobile-first)

### Documentação
- [ ] Aprendizados registrados via `./conclude_task.sh`
- [ ] NEXT_STEPS.md atualizado
- [ ] Lista de tópicos seed documentada em `docs/02_product/help_topics.md`

---

## 🔗 REFERÊNCIAS DE CÓDIGO

### Handlers Similares
- `modules/ui_web/internal/handler/member_handler.go` - Padrão CRUD HTMX
- `modules/ui_web/internal/handler/budget_handler.go` - Busca + listagem
- `modules/ui_web/internal/handler/legal_handler.go` - Download de documentos

### Templates Similares
- `modules/ui_web/templates/dashboard_simple.html` - Header com navegação
- `modules/ui_web/templates/member_simple.html` - Lista categorizada

### Services Similares
- `modules/core_lume/internal/service/member_service.go` - CRUD + validações
- `modules/core_lume/internal/service/decision_service.go` - Seed de dados

### Skills a Aplicar
- `docs/skills/developing-digna-backend/SKILL.md` - Anti-Float, DDD, TDD
- `docs/skills/rendering-digna-frontend/SKILL.md` - HTMX, Cache-Proof templates
- `docs/skills/applying-solidarity-logic/SKILL.md` - **CRÍTICO:** Linguagem popular, pedagogia
- `docs/skills/managing-sovereign-data/SKILL.md` - Isolamento SQLite

---

## 🚨 RISCOS E MITIGAÇÕES

| Risco | Probabilidade | Impacto | Mitigação |
|-------|--------------|---------|-----------|
| Linguagem muito técnica nos tópicos | Alta | Alto | Revisão por ITCPs/comunidade; teste de usabilidade com usuários reais |
| Conteúdo de ajuda desatualizado | Média | Médio | Versionamento de tópicos + hash de integridade; processo de atualização via admin (futuro) |
| Performance de busca com muitos tópicos | Baixa | Baixo | Índice no banco; limite de resultados na busca |
| Links quebrados em fontes oficiais | Média | Baixo | Validação periódica; usar Wayback Machine como fallback |
| Cache de tópicos desatualizado | Média | Baixo | Invalidação de cache ao atualizar tópico |

---

## 📅 CRONOGRAMA ESTIMADO

| Fase | Duração | Entregas |
|------|---------|----------|
| **1. Domain + Repository** | 3 horas | Entidade, Repository interface, SQLite implementation (central.db), seed, testes |
| **2. Service Layer** | 3 horas | Lógica de busca, listagem, categorias, testes |
| **3. Handler + Templates** | 4 horas | Handler HTTP, 2 templates cache-proof, navegação, testes |
| **4. Integração UI** | 2 horas | Link no header, botões `?` em campos existentes (RF-19) |
| **5. Validação + Seed** | 2 horas | Seed de 6 tópicos, smoke test, E2E, documentação |

**Total:** 14 horas

---

## 🎯 INSTRUÇÕES PARA O AGENTE

1. **ANTES DE CODIFICAR:**
   - Execute `./scripts/tools/quick_agent_check.sh all` para validar contexto
   - Preencha o checklist pré-implementação em `docs/implementation_plans/help_pre_check.md`
   - Analise handlers similares (`member_handler.go`, `budget_handler.go`)
   - **CRÍTICO:** Consulte `docs/skills/applying-solidarity-logic/SKILL.md` para garantir linguagem popular

2. **DURANTE IMPLEMENTAÇÃO:**
   - Siga TDD (testes primeiro)
   - Use `int64` para `ViewCount` (Anti-Float)
   - Carregue templates via `ParseFiles()` no handler (cache-proof)
   - Garanta que `central.db` seja criado se não existir
   - Escreva tópicos seed em linguagem popular (teste: uma criança de 10 anos entenderia?)

3. **APÓS IMPLEMENTAÇÃO:**
   - Execute smoke test: `./scripts/dev/smoke_test_new_feature.sh "Central de Ajuda" "/help"`
   - Execute validação E2E: `./scripts/dev/validate_e2e.sh --basic --headless`
   - Documente aprendizados: `./conclude_task.sh "Aprendizados: [resumo]" --success`

4. **NUNCA:**
   - Use jargão técnico nos tópicos de ajuda
   - Crie templates globais ou use `ParseGlob()`
   - Coloque tópicos de ajuda no banco da entidade (deve ser central)
   - Conclua a tarefa sem seed de tópicos inicial

---

## 📋 CHECKLIST DE VALIDAÇÃO FINAL

```bash
# 1. Validação Anti-Float
grep -r "float[0-9]*" modules/core_lume/internal/service/help_service.go
# Deve retornar apenas logs/comentários

# 2. Validação de Templates
grep -r "ParseFiles" modules/ui_web/internal/handler/help_handler.go
# Deve existir no handler

# 3. Validação de Banco Central
ls -la data/entities/central.db
# Deve existir

# 4. Testes Unitários
cd modules/core_lume && go test ./internal/service/help_service_test.go -v
cd modules/ui_web && go test ./internal/handler/help_handler_test.go -v

# 5. Smoke Test
./scripts/dev/smoke_test_new_feature.sh "Central de Ajuda" "/help"

# 6. Validação E2E
./scripts/dev/validate_e2e.sh --basic --headless

# 7. Validação de Conteúdo (Manual)
# Acessar /help e verificar:
# - 6 tópicos listados
# - Busca funcional
# - Linguagem popular (sem jargão)
# - Próximo passo acionável em cada tópico
```

---

## 🔄 INTEGRAÇÃO COM RF-19 (Perfil de Elegibilidade)

Após esta task ser concluída, atualizar RF-19 para adicionar botões `?` nos campos:

```html
<!-- Em eligibility_simple.html -->
<label class="flex items-center gap-2">
  Inscrito no CadÚnico?
  <a href="/help/topic/cadunico" target="_blank" class="text-digna-primary">?</a>
</label>

<label class="flex items-center gap-2">
  Inadimplência Ativa?
  <a href="/help/topic/inadimplencia" target="_blank" class="text-digna-primary">?</a>
</label>
```

---

## 📊 MÉTRICAS DE SUCESSO

| Métrica | Alvo (3 meses) | Como Medir |
|---------|---------------|------------|
| Tópicos Criados | 10+ | `SELECT COUNT(*) FROM help_topics` |
| Visualizações/Mês | 500+ | `SUM(view_count)` |
| Tópicos Mais Acessados | Top 5 identificados | `ORDER BY view_count DESC` |
| Redução de Abandono em Formulários | 30% | Analytics de formulários RF-19 |
| Feedback de Usabilidade | 4/5 estrelas | Pesquisa com usuários reais |

---

**PRONTO PARA INICIAR?**

Confirme que compreendeu:
1. [ ] Regra Anti-Float (int64 para ViewCount)
2. [ ] Templates cache-proof (ParseFiles no handler)
3. [ ] Banco central para tópicos globais (não por entidade)
4. [ ] Linguagem popular em todos os tópicos (skill: applying-solidarity-logic)
5. [ ] Link "Ajuda" no header de todas as páginas
6. [ ] Validação E2E obrigatória antes de concluir

Se todas as caixas estiverem marcadas, inicie pela **Fase 1: Domain Layer**.

---

**Gerado em:** 27/03/2026  
**Próxima Revisão:** Após implementação da Fase 1  
**Dependência:** RF-19 (Perfil de Elegibilidade) ⏳ EM ANDAMENTO  
**Habilita:** Melhoria de UX em todos os formulários do sistema