# 📋 Prompt para RF-27: Cálculo Automático de DAS MEI

**Tipo:** Feature  
**Módulo:** `modules/ui_web`, `modules/core_lume`  
**Prioridade:** ALTA (Baixo esforço, alto valor percebido)  
**Estimativa:** 4-6 horas  

---

## 🎯 CONTEXTO DA TAREFA

Você vai implementar o cálculo automático do DAS MEI (Documento de Arrecadação do Simples Nacional) para microempreendedores individuais que utilizam o sistema Digna. Esta funcionalidade foi identificada como de **alto valor percebido** e **baixa complexidade técnica**, servindo como preparação para o Portal de Oportunidades (RF-20).

**Regra de Negócio Principal:** O DAS MEI corresponde a **5% do salário mínimo vigente** + valores fixos de ICMS/ISS (se aplicável). O sistema deve manter uma tabela interna versionada do salário mínimo por ano.

---

## 📝 DESCRIÇÃO DETALHADA

### Problema a Resolver
Microempreendedores MEI frequentemente esquecem ou calculam incorretamente o valor mensal do DAS, gerando inadimplência e perda de benefícios. O Digna deve automatizar este cálculo e alertar sobre vencimentos.

### Escopo da Implementação
1. **Tabela de Salário Mínimo:** Estrutura interna versionada por ano (atualização via decreto presidencial)
2. **Cálculo Automático:** 5% do salário mínimo + ICMS (R$ 1,00 comércio) + ISS (R$ 5,00 serviços)
3. **Alertas de Vencimento:** Notificação até dia 20 de cada mês
4. **Histórico de Pagamentos:** Registro de quais meses foram pagos
5. **Geração de Guia:** Link para pagamento ou exportação do boleto

### Fora do Escopo
- Integração direta com Receita Federal para pagamento (apenas cálculo e alerta)
- Cálculo para outros regimes tributários (Simples Nacional ME/EPP)
- Emissão de NF-e (já definido como fora do Core Lume)

---

## 🏗️ ARQUITETURA TÉCNICA

### Estrutura de Output Esperada

```
modules/
├── core_lume/
│   ├── internal/
│   │   ├── domain/
│   │   │   └── das_mei.go              # Entidade DAS MEI
│   │   ├── service/
│   │   │   └── das_mei_service.go      # Lógica de cálculo
│   │   └── repository/
│   │       └── das_mei_repository.go   # Persistência SQLite
│   └── pkg/
│       └── das/
│           └── das.go                  # API pública
│
└── ui_web/
    ├── internal/handler/
    │   └── das_mei_handler.go          # Handler HTTP
    └── templates/
        └── das_mei_simple.html         # Template cache-proof
```

### Integrações Existentes
- **LifecycleManager:** Para acesso ao SQLite da entidade
- **BaseHandler:** Para renderização de templates
- **MemberService:** Para verificar se usuário é MEI
- **BudgetService:** Para integrar com planejamento financeiro

---

## 🛠️ TAREFAS DE IMPLEMENTAÇÃO

### 1. Domain Layer (`core_lume/internal/domain/das_mei.go`)
```go
type DASMEI struct {
    ID              string    // UUID
    EntityID        string    // Vínculo com entidade
    Competencia     string    // YYYY-MM
    ValorDevido     int64     // int64 - Anti-Float (centavos)
    ValorPago       int64     // int64 - Anti-Float (centavos)
    DataVencimento  int64     // Unix timestamp
    DataPagamento   int64     // Unix timestamp (0 se não pago)
    Status          string    // PENDENTE, PAGO, VENCIDO
    SalarioMinimo   int64     // Salário mínimo de referência (centavos)
    CreatedAt       int64     // Unix timestamp
    UpdatedAt       int64     // Unix timestamp
}

// Regras de validação
func (d *DASMEI) Validate() error
func (d *DASMEI) IsOverdue() bool
func (d *DASMEI) CalculateAmount(salarioMinimo int64, atividade string) int64
```

### 2. Service Layer (`core_lume/internal/service/das_mei_service.go`)
```go
type DASMEIService struct {
    repo DASMEIRepository
}

// Métodos obrigatórios:
func (s *DASMEIService) GenerateMonthlyDAS(entityID string, competencia string) (*DASMEI, error)
func (s *DASMEIService) GetPendingDAS(entityID string) ([]*DASMEI, error)
func (s *DASMEIService) MarkAsPaid(entityID string, dasID string) error
func (s *DASMEIService) GetMinimumWage(year int) int64
func (s *DASMEIService) CheckOverdueAlerts(entityID string) ([]Alert, error)
```

### 3. Repository Layer (`core_lume/internal/repository/das_mei_repository.go`)
```go
type DASMEIRepository interface {
    Save(das *DASMEI) error
   FindByCompetencia(entityID, competencia string) (*DASMEI, error)
    ListByEntity(entityID string) ([]*DASMEI, error)
    ListPending(entityID string) ([]*DASMEI, error)
}

// Tabela SQLite:
// CREATE TABLE das_mei (
//     id TEXT PRIMARY KEY,
//     entity_id TEXT NOT NULL,
//     competencia TEXT NOT NULL,
//     valor_devido INTEGER NOT NULL,
//     valor_pago INTEGER DEFAULT 0,
//     data_vencimento INTEGER NOT NULL,
//     data_pagamento INTEGER DEFAULT 0,
//     status TEXT NOT NULL,
//     salario_minimo INTEGER NOT NULL,
//     created_at INTEGER,
//     updated_at INTEGER,
//     UNIQUE(entity_id, competencia)
// )
```

### 4. Handler Layer (`modules/ui_web/internal/handler/das_mei_handler.go`)
```go
type DASMEIHandler struct {
    *BaseHandler
    lifecycleManager lifecycle.LifecycleManager
    dasService       *service.DASMEIService
}

// Rotas obrigatórias:
// GET /das-mei              → Lista de DAS com status
// POST /das-mei/generate    → Gera DAS do mês atual
// POST /das-mei/{id}/pay    → Marca como pago
// GET /das-mei/alerts       → Alertas de vencimento (HTMX)
```

### 5. Template (`modules/ui_web/templates/das_mei_simple.html`)
- Seguir padrão `*_simple.html` (documento HTML completo)
- Carregar via `ParseFiles()` no handler (cache-proof)
- Usar paleta "Soberania e Suor"
- Incluir navegação padrão (header/footer)
- Exibir:
  - Status do mês atual (Pendente/Pago/Vencido)
  - Valor calculado (formatado como R$)
  - Botão de ação (Gerar/Marcar Pago)
  - Histórico dos últimos 6 meses
  - Alertas visuais (cores por status)

---

## 📊 REGRAS DE NEGÓCIO DETALHADAS

### Cálculo do DAS MEI (2026)
```
Salário Mínimo 2026: R$ 1.518,00 (151800 centavos - int64)

Comércio (ICMS):  5% do SM = R$ 75,90
Serviços (ISS):   5% do SM = R$ 75,90
Comércio+Serviços: 5% do SM + ISS fixo = R$ 75,90 + R$ 5,00 = R$ 80,90

Valores em int64 (centavos):
- Comércio: 7590
- Serviços: 7590
- Misto: 8090
```

### Tabela de Salário Mínimo (Versionada)
```go
var minimumWageTable = map[int]int64{
    2024: 141200,  // R$ 1.412,00
    2025: 151800,  // R$ 1.518,00
    2026: 151800,  // R$ 1.518,00 (ajustar quando houver decreto)
}
```

### Vencimento
- **Data fixa:** Dia 20 de cada mês
- **Se fim de semana/feriado:** Antecipar para dia útil anterior
- **Alertas:** Dia 15 (5 dias antes), Dia 19 (1 dia antes), Dia 20 (vence hoje)

---

## ✅ CRITÉRIOS DE ACEITE (Definition of Done)

### Arquitetura
- [ ] Segue Clean Architecture (Domain → Service → Repository → Handler)
- [ ] Zero `float` para valores financeiros (usar `int64` para centavos)
- [ ] Templates `*_simple.html` carregados via `ParseFiles()` no handler
- [ ] Soberania mantida (dados só no `.sqlite` da entidade)
- [ ] Handler estende `BaseHandler`

### Funcionalidade
- [ ] Cálculo correto do DAS (5% do salário mínimo)
- [ ] Tabela de salário mínimo versionada por ano
- [ ] Alertas de vencimento funcionais (HTMX)
- [ ] Histórico de pagamentos registrado
- [ ] Validação de entidade MEI antes de gerar DAS

### Testes
- [ ] Testes unitários para service (cálculo, validações)
- [ ] Testes unitários para repository (CRUD)
- [ ] Testes de handler (HTTP responses)
- [ ] Cobertura >90% para service e handler
- [ ] Smoke test: `./scripts/dev/smoke_test_new_feature.sh "DAS MEI" "/das-mei"`
- [ ] Validação E2E: `./scripts/dev/validate_e2e.sh --basic --headless`

### UI/UX
- [ ] Design segue "Soberania e Suor" (#2A5CAA, #4A7F3E, #F57F17)
- [ ] Linguagem popular (sem jargão "competência", usar "mês de referência")
- [ ] Alertas visuais claros (verde=pago, amarelo=pendente, vermelho=vencido)
- [ ] Link adicionado na navegação principal (dashboard)

### Documentação
- [ ] Aprendizados registrados via `./conclude_task.sh`
- [ ] NEXT_STEPS.md atualizado
- [ ] Comentários no código para regras de negócio críticas

---

## 🔗 REFERÊNCIAS DE CÓDIGO

### Handlers Similares
- `modules/ui_web/internal/handler/budget_handler.go` - Estrutura de handler com service
- `modules/ui_web/internal/handler/member_handler.go` - Padrão CRUD HTMX
- `modules/ui_web/internal/handler/accountant_handler.go` - File download pattern (se precisar exportar guia)

### Templates Similares
- `modules/ui_web/templates/budget_simple.html` - Dashboard com alertas visuais
- `modules/ui_web/templates/cash_simple.html` - Tabela com status e ações

### Services Similares
- `modules/core_lume/internal/service/surplus_calculator.go` - Cálculos financeiros em int64
- `modules/budget/internal/service/budget_service.go` - Alertas e vencimentos

### Skills a Aplicar
- `docs/skills/developing-digna-backend/SKILL.md` - Anti-Float, DDD, TDD
- `docs/skills/rendering-digna-frontend/SKILL.md` - HTMX, Cache-Proof templates
- `docs/skills/applying-solidarity-logic/SKILL.md` - Linguagem popular, pedagogia
- `docs/skills/managing-sovereign-data/SKILL.md` - Isolamento SQLite

---

## 🚨 RISCOS E MITIGAÇÕES

| Risco | Probabilidade | Impacto | Mitigação |
|-------|--------------|---------|-----------|
| Salário mínimo muda durante o ano | Baixa | Médio | Tabela versionada, atualização via migration |
| Entity não é MEI | Média | Alto | Validação no service antes de gerar DAS |
| Template cache issues | Alta | Baixo | Usar `ParseFiles()` no handler, não global |
| Entity ID não encontrado | Média | Médio | Middleware de autenticação já existente |
| Cálculo incorreto | Baixa | Crítico | Testes unitários rigorosos + revisão de código |

---

## 📅 CRONOGRAMA ESTIMADO

| Fase | Duração | Entregas |
|------|---------|----------|
| **1. Domain + Repository** | 1.5 horas | Entidade, Repository interface, SQLite implementation, testes |
| **2. Service Layer** | 1.5 horas | Lógica de cálculo, tabela salário mínimo, alertas, testes |
| **3. Handler + Template** | 2 horas | Handler HTTP, template cache-proof, navegação, testes |
| **4. Integração + Validação** | 1 hora | Smoke test, E2E, documentação, conclusão |

**Total:** 6 horas

---

## 🎯 INSTRUÇÕES PARA O AGENTE

1. **ANTES DE CODIFICAR:**
   - Execute `./scripts/tools/quick_agent_check.sh all` para validar contexto
   - Preencha o checklist pré-implementação em `docs/implementation_plans/das_mei_pre_check.md`
   - Analise handlers similares (`budget_handler.go`, `member_handler.go`)

2. **DURANTE IMPLEMENTAÇÃO:**
   - Siga TDD (testes primeiro)
   - Use `int64` para TODOS os valores financeiros
   - Carregue templates via `ParseFiles()` no handler
   - Valide entity_id em todas as queries

3. **APÓS IMPLEMENTAÇÃO:**
   - Execute smoke test: `./scripts/dev/smoke_test_new_feature.sh "DAS MEI" "/das-mei"`
   - Execute validação E2E: `./scripts/dev/validate_e2e.sh --basic --headless`
   - Documente aprendizados: `./conclude_task.sh "Aprendizados: [resumo]" --success`

4. **NUNCA:**
   - Use `float` para valores financeiros
   - Crie templates globais ou use `ParseGlob()`
   - Acesse bancos de outras entidades
   - Conclua a tarefa sem testes passando

---

## 📋 CHECKLIST DE VALIDAÇÃO FINAL

```bash
# 1. Validação Anti-Float
grep -r "float[0-9]*" modules/core_lume/internal/service/das_mei_service.go
# Deve retornar apenas logs/comentários

# 2. Validação de Templates
grep -r "ParseFiles" modules/ui_web/internal/handler/das_mei_handler.go
# Deve existir no handler

# 3. Testes Unitários
cd modules/core_lume && go test ./internal/service/das_mei_service_test.go -v
cd modules/ui_web && go test ./internal/handler/das_mei_handler_test.go -v

# 4. Smoke Test
./scripts/dev/smoke_test_new_feature.sh "DAS MEI" "/das-mei"

# 5. Validação E2E
./scripts/dev/validate_e2e.sh --basic --headless
```

---

**PRONTO PARA INICIAR?**

Confirme que compreendeu:
1. [ ] Regra Anti-Float (int64 para centavos)
2. [ ] Templates cache-proof (ParseFiles no handler)
3. [ ] Soberania de dados (entity_id isolation)
4. [ ] Tabela de salário mínimo versionada
5. [ ] Validação E2E obrigatória antes de concluir

Se todas as caixas estiverem marcadas, inicie pela **Fase 1: Domain Layer**.

---
