# 📋 Prompt para RF-19: Perfil de Elegibilidade (Campos Complementares)

**Tipo:** Feature  
**Módulo:** `modules/core_lume`, `modules/ui_web`  
**Prioridade:** ALTA (Habilita o Portal de Oportunidades - RF-20)  
**Estimativa:** 6-8 horas  

---

## 🎯 CONTEXTO DA TAREFA

Você vai implementar o **Perfil de Elegibilidade** no sistema Digna, adicionando campos complementares ao cadastro da entidade que habilitam o **match automático** com programas de financiamento no Portal de Oportunidades (RF-20).

**Princípio Central (PDF v1.0, Seção 3.1):** *"Nenhum usuário precisa preencher o mesmo dado duas vezes."*

Os dados do ERP (faturamento, CNAE, município, regime tributário) já são capturados automaticamente pelo uso cotidiano. Esta task adiciona **apenas os campos complementares** de preenchimento único, que serão reutilizados automaticamente pelo Portal e pela Rede Digna.

---

## 📝 DESCRIÇÃO DETALHADA

### Problema a Resolver
O Portal de Oportunidades precisa de informações específicas para executar o match com programas de financiamento (ex: "inscrito no CadÚnico?", "sócio mulher?", "finalidade do crédito?"). Sem esses campos, o match seria impreciso ou exigiria formulários repetitivos.

### Escopo da Implementação
1. **Extensão do Modelo de Domínio:** Adicionar struct `EligibilityProfile` com campos complementares
2. **Migração de Banco:** Criar tabela `eligibility_profiles` no SQLite da entidade
3. **Service Layer:** CRUD + validações de negócio
4. **Handler + Template:** Interface simples para preenchimento único (HTMX)
5. **API Pública:** Expor perfil para consumo pelo módulo `integrations` (futuro Portal)

### Fora do Escopo
- Lógica de match com programas (será RF-20)
- Integração com APIs externas de certidões (será RF-21)
- Interface do Portal de Oportunidades (será RF-20)

---

## 🏗️ ARQUITETURA TÉCNICA

### Estrutura de Output Esperada

```
modules/
├── core_lume/
│   ├── internal/
│   │   ├── domain/
│   │   │   └── eligibility_profile.go    # NOVO: Entidade + validações
│   │   ├── service/
│   │   │   └── eligibility_service.go    # NOVO: Lógica de negócio
│   │   └── repository/
│   │       └── eligibility_repository.go # NOVO: Persistência SQLite
│   └── pkg/
│       └── eligibility/
│           └── eligibility.go            # NOVO: API pública
│
└── ui_web/
    ├── internal/handler/
    │   └── eligibility_handler.go        # NOVO: Handler HTTP
    └── templates/
        └── eligibility_simple.html       # NOVO: Template cache-proof
```

### Integrações Existentes
- **LifecycleManager:** Para acesso ao SQLite da entidade
- **BaseHandler:** Para renderização de templates
- **MemberService:** Para verificar permissões de edição
- **Enterprise (core_lume):** Para vincular perfil à entidade principal

---

## 🛠️ TAREFAS DE IMPLEMENTAÇÃO

### 1. Domain Layer (`core_lume/internal/domain/eligibility_profile.go`)

```go
type EligibilityProfile struct {
    ID         string // UUID
    EntityID   string // Vínculo com entidade (único por entidade)
    
    // Dados já capturados pelo ERP (referência, não duplicação)
    CNPJ             string // Copiado de Enterprise
    CNAE             string // Copiado de Enterprise
    Municipio        string // Copiado de Enterprise
    UF               string // Copiado de Enterprise
    FaturamentoAnual int64  // int64 - Anti-Float (centavos)
    RegimeTributario string
    DataAbertura     int64 // Unix timestamp
    SituacaoFiscal   string
    
    // CAMPOS COMPLEMENTARES (preenchimento único)
    InscritoCadUnico    bool   // Habilita programas sociais
    SocioMulher         bool   // Prioridade em linhas com foco de gênero
    InadimplenciaAtiva  bool   // Direciona ao Desenrola antes de crédito novo
    FinalidadeCredito   string // Enum: CAPITAL_GIRO, EQUIPAMENTO, REFORMA, OUTRO
    ValorNecessario     int64  // int64 - Anti-Float (centavos)
    TipoEntidade        string // Enum: MEI, ME, EPP, Cooperativa, OSC, OSCIP, PF
    ContabilidadeFormal bool   // Requisito de alguns programas
    
    // Metadados
    PreenchidoEm  int64 // Unix timestamp - primeiro preenchimento
    AtualizadoEm  int64 // Unix timestamp - última atualização
    PreenchidoPor string // ID do usuário que preencheu
    
    CreatedAt int64
    UpdatedAt int64
}

// Enums para campos restritos
const (
    FinalidadeCapitalGiro   = "CAPITAL_GIRO"
    FinalidadeEquipamento   = "EQUIPAMENTO"
    FinalidadeReforma       = "REFORMA"
    FinalidadeOutro         = "OUTRO"
    
    TipoEntidadeMEI         = "MEI"
    TipoEntidadeME          = "ME"
    TipoEntidadeEPP         = "EPP"
    TipoEntidadeCooperativa = "Cooperativa"
    TipoEntidadeOSC         = "OSC"
    TipoEntidadeOSCIP       = "OSCIP"
    TipoEntidadePF          = "PF"
)

// Validações de domínio
func (e *EligibilityProfile) Validate() error
func (e *EligibilityProfile) IsComplete() bool // Todos os campos obrigatórios preenchidos?
func (e *EligibilityProfile) CanEdit(userID string) bool // Apenas coordenadores
```

### 2. Repository Layer (`core_lume/internal/repository/eligibility_repository.go`)

```go
type EligibilityRepository interface {
    Save(profile *EligibilityProfile) error // UPSERT: um perfil por entidade
    FindByEntityID(entityID string) (*EligibilityProfile, error)
    ListIncomplete() ([]*EligibilityProfile, error) // Para dashboard de preenchimento
    UpdateFields(entityID string, fields map[string]interface{}) error // Atualização parcial
}

// Tabela SQLite:
// CREATE TABLE IF NOT EXISTS eligibility_profiles (
//     id TEXT PRIMARY KEY,
//     entity_id TEXT NOT NULL UNIQUE,
//     
//     // Dados do ERP (cópia para consulta rápida)
//     cnpj TEXT,
//     cnae TEXT,
//     municipio TEXT,
//     uf TEXT,
//     faturamento_anual INTEGER,  -- int64, Anti-Float
//     regime_tributario TEXT,
//     data_abertura INTEGER,
//     situacao_fiscal TEXT,
//     
//     // Campos complementares
//     inscrito_cad_unico INTEGER,      -- 0/1 (bool)
//     socio_mulher INTEGER,            -- 0/1 (bool)
//     inadimplencia_ativa INTEGER,     -- 0/1 (bool)
//     finalidade_credito TEXT,
//     valor_necessario INTEGER,        -- int64, Anti-Float
//     tipo_entidade TEXT,
//     contabilidade_formal INTEGER,    -- 0/1 (bool)
//     
//     // Metadados
//     preenchido_em INTEGER,
//     atualizado_em INTEGER,
//     preenchido_por TEXT,
//     created_at INTEGER,
//     updated_at INTEGER
// )
```

### 3. Service Layer (`core_lume/internal/service/eligibility_service.go`)

```go
type EligibilityService struct {
    repo    EligibilityRepository
    enterpriseRepo EnterpriseRepository // Para copiar dados do ERP
}

// Métodos obrigatórios:
func (s *EligibilityService) CreateOrUpdate(entityID string, userID string, input EligibilityInput) (*EligibilityProfile, error)
func (s *EligibilityService) GetProfile(entityID string) (*EligibilityProfile, error)
func (s *EligibilityService) SyncFromEnterprise(entityID string) error // Atualiza cópia dos dados do ERP
func (s *EligibilityService) GetCompletionStatus(entityID string) (float64, error) // % de campos preenchidos
func (s *EligibilityService) ListEligibleForProgram(programCriteria ProgramCriteria) ([]string, error) // Futuro: match engine

// Input para criação/atualização (apenas campos editáveis)
type EligibilityInput struct {
    InscritoCadUnico    *bool   `json:"inscrito_cad_unico,omitempty"`
    SocioMulher         *bool   `json:"socio_mulher,omitempty"`
    InadimplenciaAtiva  *bool   `json:"inadimplencia_ativa,omitempty"`
    FinalidadeCredito   *string `json:"finalidade_credito,omitempty"`
    ValorNecessario     *int64  `json:"valor_necessario,omitempty"`  // int64
    TipoEntidade        *string `json:"tipo_entidade,omitempty"`
    ContabilidadeFormal *bool   `json:"contabilidade_formal,omitempty"`
}
```

### 4. Handler Layer (`modules/ui_web/internal/handler/eligibility_handler.go`)

```go
type EligibilityHandler struct {
    *BaseHandler
    lifecycleManager lifecycle.LifecycleManager
    eligibilityService *service.EligibilityService
}

// Rotas obrigatórias:
// GET /eligibility?entity_id={id}        → Renderiza formulário com dados atuais
// POST /eligibility?entity_id={id}       → Salva/atualiza perfil (HTMX)
// GET /eligibility/status?entity_id={id} → Retorna % de completude (para dashboard)
// GET /eligibility/export?entity_id={id} → Exporta perfil em JSON (para Portal futuro)
```

### 5. Template (`modules/ui_web/templates/eligibility_simple.html`)

- Seguir padrão `*_simple.html` (documento HTML completo)
- Carregar via `ParseFiles()` no handler (cache-proof)
- Usar paleta "Soberania e Suor"
- Incluir navegação padrão (header/footer)
- Exibir:
  - **Seção 1:** Dados do ERP (somente leitura, com ícone de "já capturado")
  - **Seção 2:** Campos complementares (formulário editável)
  - **Indicador visual:** Barra de progresso "% de perfil completo"
  - **Mensagem pedagógica:** "Preencha uma vez, use em todos os módulos"
  - **Botões:** Salvar, Cancelar, Exportar (JSON)

---

## 📊 REGRAS DE NEGÓCIO DETALHADAS

### Validações de Campos

| Campo | Tipo | Obrigatório? | Validação |
|-------|------|-------------|-----------|
| `InscritoCadUnico` | bool | Não | - |
| `SocioMulher` | bool | Não | - |
| `InadimplenciaAtiva` | bool | Não | - |
| `FinalidadeCredito` | enum | **Sim** | Deve ser um dos valores definidos |
| `ValorNecessario` | int64 | Condicional* | Se `FinalidadeCredito` != "OUTRO", deve ser > 0 |
| `TipoEntidade` | enum | **Sim** | Deve ser um dos valores definidos |
| `ContabilidadeFormal` | bool | Não | - |

*Condicional: `ValorNecessario` é obrigatório se a finalidade for específica (não "OUTRO")

### Sincronização com ERP

- Ao criar/atualizar `EligibilityProfile`, copiar automaticamente os dados atuais do `Enterprise`:
  - CNPJ, CNAE, Município, UF, FaturamentoAnual, RegimeTributario, DataAbertura, SituacaoFiscal
- Esta cópia é **somente para consulta rápida** — a fonte da verdade permanece no `Enterprise`
- Se os dados do ERP mudarem, o perfil **não** é atualizado automaticamente (para preservar histórico do match)

### Permissões de Edição

- Apenas usuários com role `COORDINATOR` podem editar o perfil
- Visualização é permitida para `MEMBER` e `ADVISOR` (somente leitura)
- Registro de `PreenchidoPor` e `AtualizadoEm` para auditoria

### API Pública (para consumo pelo Portal)

```go
// modules/core_lume/pkg/eligibility/eligibility.go
type EligibilityPublic struct {
    EntityID string `json:"entity_id"`
    
    // Dados do ERP (públicos para match)
    CNPJ             string `json:"cnpj"`
    CNAE             string `json:"cnae"`
    Municipio        string `json:"municipio"`
    UF               string `json:"uf"`
    FaturamentoAnual int64  `json:"faturamento_anual"` // int64
    RegimeTributario string `json:"regime_tributario"`
    
    // Campos complementares (públicos para match)
    InscritoCadUnico    bool   `json:"inscrito_cad_unico"`
    SocioMulher         bool   `json:"socio_mulher"`
    InadimplenciaAtiva  bool   `json:"inadimplencia_ativa"`
    FinalidadeCredito   string `json:"finalidade_credito"`
    ValorNecessario     int64  `json:"valor_necessario"` // int64
    TipoEntidade        string `json:"tipo_entidade"`
    ContabilidadeFormal bool   `json:"contabilidade_formal"`
    
    // Metadados (não expor dados sensíveis)
    PreenchidoEm int64 `json:"preenchido_em"`
    IsComplete   bool  `json:"is_complete"` // Todos os campos obrigatórios preenchidos?
}

// Função de conversão interna → pública (oculta campos sensíveis)
func ToPublic(profile *EligibilityProfile) *EligibilityPublic
```

---

## ✅ CRITÉRIOS DE ACEITE (Definition of Done)

### Arquitetura
- [ ] Segue Clean Architecture (Domain → Service → Repository → Handler)
- [ ] Zero `float` para valores financeiros (usar `int64` para centavos)
- [ ] Templates `*_simple.html` carregados via `ParseFiles()` no handler
- [ ] Soberania mantida (dados só no `.sqlite` da entidade)
- [ ] Handler estende `BaseHandler`
- [ ] API pública em `pkg/eligibility/` para consumo externo

### Funcionalidade
- [ ] CRUD completo do perfil de elegibilidade
- [ ] Sincronização automática dos dados do ERP na criação
- [ ] Validação de enums e campos obrigatórios
- [ ] Permissões de edição baseadas em role (apenas COORDINATOR)
- [ ] Indicador visual de "% de perfil completo"
- [ ] Exportação em JSON para consumo pelo Portal (futuro)

### Testes
- [ ] Testes unitários para service (validações, sincronização, permissões)
- [ ] Testes unitários para repository (CRUD, UPSERT)
- [ ] Testes de handler (HTTP responses, autenticação)
- [ ] Cobertura >90% para service e handler
- [ ] Smoke test: `./scripts/dev/smoke_test_new_feature.sh "Perfil de Elegibilidade" "/eligibility"`
- [ ] Validação E2E: `./scripts/dev/validate_e2e.sh --basic --headless`

### UI/UX
- [ ] Design segue "Soberania e Suor" (#2A5CAA, #4A7F3E, #F57F17)
- [ ] Linguagem popular (sem jargão "elegibilidade", usar "seu perfil para crédito")
- [ ] Dados do ERP em seção "somente leitura" com ícone explicativo
- [ ] Barra de progresso visual para completude do perfil
- [ ] Link adicionado na navegação principal (dashboard → "Perfil para Crédito")

### Documentação
- [ ] Aprendizados registrados via `./conclude_task.sh`
- [ ] NEXT_STEPS.md atualizado
- [ ] Comentários no código para regras de negócio críticas
- [ ] Atualizar `docs/02_product/02_models.md` com novo modelo `EligibilityProfile`

---

## 🔗 REFERÊNCIAS DE CÓDIGO

### Handlers Similares
- `modules/ui_web/internal/handler/member_handler.go` - Padrão CRUD HTMX com validação de role
- `modules/ui_web/internal/handler/budget_handler.go` - Formulário com campos condicionais
- `modules/ui_web/internal/handler/legal_handler.go` - Exportação em JSON para consumo externo

### Templates Similares
- `modules/ui_web/templates/member_simple.html` - Formulário com seções somente leitura + editáveis
- `modules/ui_web/templates/budget_simple.html` - Barra de progresso visual
- `modules/ui_web/templates/dashboard_simple.html` - Navegação e cards informativos

### Services Similares
- `modules/core_lume/internal/service/member_service.go` - Validação de permissões por role
- `modules/core_lume/internal/service/enterprise_service.go` - Sincronização entre entidades

### Skills a Aplicar
- `docs/skills/developing-digna-backend/SKILL.md` - Anti-Float, DDD, TDD
- `docs/skills/rendering-digna-frontend/SKILL.md` - HTMX, Cache-Proof templates
- `docs/skills/applying-solidarity-logic/SKILL.md` - Linguagem popular, pedagogia
- `docs/skills/managing-sovereign-data/SKILL.md` - Isolamento SQLite

---

## 🚨 RISCOS E MITIGAÇÕES

| Risco | Probabilidade | Impacto | Mitigação |
|-------|--------------|---------|-----------|
| Duplicação de dados entre Enterprise e EligibilityProfile | Alta | Médio | Documentar claramente que cópia é para consulta rápida; fonte da verdade é Enterprise |
| Validação de enums falha | Média | Alto | Usar constantes definidas + validação no service + testes unitários |
| Permissões de edição vazadas | Baixa | Crítico | Testes específicos para cada role + middleware de autenticação existente |
| Template cache issues | Alta | Baixo | Usar `ParseFiles()` no handler, não global |
| Entity ID não encontrado | Média | Médio | Middleware de autenticação já existente |
| Perfil incompleto afeta match futuro | Alta | Médio | Indicador visual de completude + mensagem pedagógica incentivando preenchimento |

---


## 🎯 INSTRUÇÕES PARA O AGENTE

1. **ANTES DE CODIFICAR:**
   - Execute `./scripts/tools/quick_agent_check.sh all` para validar contexto
   - Preencha o checklist pré-implementação em `docs/implementation_plans/eligibility_pre_check.md`
   - Analise handlers similares (`member_handler.go`, `budget_handler.go`)
   - Consulte `docs/02_product/02_models.md` para entender modelo Enterprise existente

2. **DURANTE IMPLEMENTAÇÃO:**
   - Siga TDD (testes primeiro)
   - Use `int64` para TODOS os valores financeiros (Anti-Float)
   - Carregue templates via `ParseFiles()` no handler (cache-proof)
   - Valide entity_id em todas as queries (soberania)
   - Use enums com constantes definidas, não strings soltas

3. **APÓS IMPLEMENTAÇÃO:**
   - Execute smoke test: `./scripts/dev/smoke_test_new_feature.sh "Perfil de Elegibilidade" "/eligibility"`
   - Execute validação E2E: `./scripts/dev/validate_e2e.sh --basic --headless`
   - Documente aprendizados: `./conclude_task.sh "Aprendizados: [resumo]" --success`

4. **NUNCA:**
   - Use `float` para valores financeiros
   - Crie templates globais ou use `ParseGlob()`
   - Acesse bancos de outras entidades
   - Conclua a tarefa sem testes passando
   - Exponha dados sensíveis na API pública `EligibilityPublic`

---

## 📋 CHECKLIST DE VALIDAÇÃO FINAL

```bash
# 1. Validação Anti-Float
grep -r "float[0-9]*" modules/core_lume/internal/service/eligibility_service.go
# Deve retornar apenas logs/comentários

# 2. Validação de Templates
grep -r "ParseFiles" modules/ui_web/internal/handler/eligibility_handler.go
# Deve existir no handler

# 3. Testes Unitários
cd modules/core_lume && go test ./internal/service/eligibility_service_test.go -v
cd modules/ui_web && go test ./internal/handler/eligibility_handler_test.go -v

# 4. Smoke Test
./scripts/dev/smoke_test_new_feature.sh "Perfil de Elegibilidade" "/eligibility"

# 5. Validação E2E
./scripts/dev/validate_e2e.sh --basic --headless

# 6. Validação de API Pública
grep -r "EligibilityPublic" modules/core_lume/pkg/eligibility/
# Deve existir e não expor campos sensíveis
```

---

## 🔄 PRÉ-REQUISITOS PARA RF-20 (Portal de Oportunidades)

Esta task RF-19 habilita o match automático do Portal. Após conclusão, o Portal poderá:

1. Consultar `EligibilityPublic` via API interna
2. Cruzar com critérios de programas (ex: `InscritoCadUnico == true` → habilita "Acredita no Primeiro Passo")
3. Gerar checklist de documentos baseado em campos faltantes



**PRONTO PARA INICIAR?**

Confirme que compreendeu:
1. [ ] Regra Anti-Float (int64 para centavos)
2. [ ] Templates cache-proof (ParseFiles no handler)
3. [ ] Soberania de dados (entity_id isolation)
4. [ ] Sincronização com Enterprise (cópia para consulta, não duplicação)
5. [ ] API pública segura (EligibilityPublic sem dados sensíveis)
6. [ ] Validação E2E obrigatória antes de concluir

Se todas as caixas estiverem marcadas, inicie pela **Fase 1: Domain Layer**.

