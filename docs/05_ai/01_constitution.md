title: Constituição de IA - Projeto Digna
status: implemented
version: 2.0
last_updated: 2026-03-27
---

# Constituição de IA - Ecossistema Digna

> **Nota:** Este documento estabelece as regras sagradas e princípios que TODOS os agentes de IA devem seguir ao trabalhar no projeto Digna. Violações destas regras são consideradas erros críticos.
>
> **Versão 2.0:** Incorpora a visão de Ecossistema de 4 Módulos (PDF v1.0), o Sistema de Ajuda Educativa (RF-30 - decisão da sessão 27/03/2026), e preserva todas as regras validadas nas Sprints 1-16.

---

## 📋 Contexto da Atualização (27/03/2026)

**Motivação:** O projeto Digna evoluiu de um ERP contábil para um **Ecossistema de 4 Módulos** conforme especificação PDF v1.0. Esta atualização documenta:

1. **Novos Módulos:** Motor de Indicadores (RF-18), Portal de Oportunidades (RF-19 a RF-23), Rede Digna (RF-24 a RF-26)
2. **RF-30:** Sistema de Ajuda Educativa Estruturada (decisão da sessão 27/03/2026)
3. **Princípio Central:** "Nenhum usuário precisa preencher o mesmo dado duas vezes"
4. **Preservação:** Todas as regras sagradas das Sprints 1-16 mantidas

**Versão Anterior:** 1.0 (2026-03-13)  
**Nova Versão:** 2.0 (2026-03-27)

---

## ⚡ REGRAS SAGRADAS (NÃO NEGOCIÁVEIS)

### 1. 🚫 ANTI-FLOAT (Regra de Ouro)

**Proibido:** `float`, `float32`, `float64` para valores financeiros ou de tempo.

**Obrigatório:** `int64` para centavos (R$ 1,00 = 100) e minutos trabalhados.

**Validação:** Todo código deve passar por `grep -r "float[0-9]*" modules/` antes de commit.

**Exceção:** Apenas para exibição visual em templates (ex: `float64(amount)/100` para formatar R$).

```go
// ✅ CORRETO
type Product struct {
    PriceInCents int64  // R$ 1,99 → 199
}

// ❌ ERRADO (CRÍTICO)
type Product struct {
    Price float64  // R$ 1,99 → 1.99
}
```

**Justificativa:** Erros de arredondamento IEEE 754 podem corromper a integridade contábil e violar a norma ITG 2002.

---

### 2. 🔒 SOBERANIA DE DADOS

**Isolamento:** Cada entidade possui seu próprio arquivo SQLite isolado fisicamente (`data/entities/{entity_id}.db`).

**Banco Central:** `data/entities/central.db` apenas para relações inter-tenant (vínculos, indicadores, programas, help_topics).

**Proibido:** JOINs entre bancos de entidades diferentes.

**LifecycleManager:** Ponto único de acesso a bancos SQLite.

**Contexto:** `entity_id` extraído de `r.Context().Value("entity_id")` ou `r.URL.Query().Get("entity_id")`.

**Exit Power:** O usuário deve poder copiar seu arquivo `.db` e levar para qualquer sistema SQL compatível.

```go
// ✅ CORRETO
db, err := h.lifecycleManager.GetDatabase(entityID)

// ❌ ERRADO (CRÍTICO)
db, err := sql.Open("sqlite3", "data/entities/*.db")  // Acesso direto sem LifecycleManager
```

**Justificativa:** Garante que o dado pertence à entidade, não à plataforma, preservando a autonomia da Economia Solidária.

---

### 3. 📄 CACHE-PROOF TEMPLATES

**Nomenclatura:** Templates devem ser `*_simple.html` (documentos HTML completos).

**Carregamento:** `template.ParseFiles("templates/nome_simple.html")` NO HANDLER, não como variável global.

**Proibido:** `template.ParseGlob()`, variáveis globais de template, templates parciais com `{{define "content"}}`.

**BaseHandler:** `modules/ui_web/internal/handler/base_handler.go` gerencia TemplateManager compartilhado.

```go
// ✅ CORRETO
func (h *DashboardHandler) DashboardPage(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("templates/dashboard_simple.html")
    // ...
}

// ❌ ERRADO (CRÍTICO)
var dashboardTemplate = template.Must(template.ParseFiles("templates/dashboard.html"))  // Global
```

**Justificativa:** Resolve problemas de cache persistente do Go que sobrevivem a recompilações e impedem desenvolvimento ágil.

---

### 4. 🏛️ CLEAN ARCHITECTURE + DDD

**Camadas:**
- `internal/domain/` - Entidades puras, interfaces Repository (SEM SQL/HTTP)
- `internal/service/` - Casos de uso, orquestração (depende de interfaces)
- `internal/repository/` - Implementações SQLite (via LifecycleManager)
- `internal/handler/` - HTTP handlers (UI Web)

**Proibido:** Importar `sql`, `http` ou frameworks externos na camada de domínio.

**Dependency Injection:** Services dependem de abstrações (interfaces), não implementações concretas.

```go
// ✅ CORRETO
type LedgerService struct {
    repo LedgerRepository  // Interface
}

// ❌ ERRADO (CRÍTICO)
type LedgerService struct {
    db *sql.DB  // Implementação concreta vazando para service
}
```

**Justificativa:** Protege o domínio de negócios de mudanças tecnológicas e facilita testabilidade.

---

### 5. 🎓 PEDAGOGIA SOCIAL (RF-30 - NOVO)

**Princípio:** O sistema ensina enquanto é operado. Nenhum usuário deve se sentir humilhado por não entender um termo.

**Linguagem:** Popular, sem jargões técnicos. Usuário com 5ª série deve entender sem ajuda externa.

**Sistema de Ajuda:** Todo campo técnico deve ter botão "?" linkado para `/help/topic/{key}`.

**Conteúdo de Ajuda:** Explicação + legislação + próximo passo acionável.

**Proibido:** Jargões como "cadastramento", "regularização fiscal", "exercício social" na interface do produtor.

```go
// ✅ CORRETO (Interface)
"Inscrito no CadÚnico?" + botão "?" → /help/topic/cadunico

// ✅ CORRETO (Conteúdo de Ajuda)
"É o cadastro do governo para programas sociais. Procure o CRAS do seu município."

// ❌ ERRADO (CRÍTICO)
"Situação de Enquadramento no Cadastro Único para Programas Sociais do Governo Federal"
```

**Justificativa:** Respeita a baixa literacia digital do público-alvo e cumpre o Pilar 2 (Tradução Cultural) e Pilar 5 (Ferramenta Pedagógica) do Digna.

---

### 6. 🌐 ECOSSISTEMA DE 4 MÓDULOS (PDF v1.0 - NOVO)

**Módulo 1 - digna ERP:** Núcleo. Gestão financeira, fiscal e contábil.

**Módulo 2 - Motor de Indicadores:** Coleta APIs externas (BCB, IBGE). Cache local com TTL.

**Módulo 3 - Portal de Oportunidades:** Match automático perfil × programas de financiamento.

**Módulo 4 - Rede Digna:** Marketplace solidário B2B entre entidades.

**Princípio Central:** "Nenhum usuário precisa preencher o mesmo dado duas vezes."

**Fluxo de Dados:**
1. ERP captura perfil no uso cotidiano
2. Motor coleta indicadores diariamente
3. Portal executa match automaticamente
4. Rede sugere conexões B2B
5. Ajuda Educativa contextualiza em todos os pontos

```go
// ✅ CORRETO (Reuso de Dados)
// EligibilityProfile copia dados do Enterprise, não exige reentrada
type EligibilityProfile struct {
    CNPJ             string  // Copiado de Enterprise
    FaturamentoAnual int64   // Copiado de Enterprise
    InscritoCadUnico bool    // Complementar (preenchimento único)
}

// ❌ ERRADO (CRÍTICO)
// Exigir que usuário digite CNPJ novamente no Portal
```

**Justificativa:** Reduz fricção, elimina barreiras de uso e respeita o tempo do empreendedor popular.

---

### 7. 🔐 LAICIDADE DO PRODUTO (PDF v1.0 - NOVO)

**Princípio:** A Teologia informa decisões de design internamente, mas o produto é acessível independente da crença do usuário.

**Proibido:** Conteúdo religioso explícito na interface do usuário final.

**Canal de Distribuição:** Igrejas e comunidades de fé são estratégicos, não doutrinários.

**Documentação:** Princípios teológicos documentados em `docs/04_governance/` sem expor ao usuário final.

```go
// ✅ CORRETO (Interno)
// Teologia informa design: Dignidade Humana → Interface que nunca humilha

// ✅ CORRETO (Interface)
// Usuário vê: "Seu negócio com dignidade" (não religioso)

// ❌ ERRADO (CRÍTICO)
// Usuário vê: "Deus abençoe seu negócio" (conteúdo religioso explícito)
```

**Justificativa:** Maximiza adoção enquanto preserva a base ética do projeto.

---

## 🎯 WORKFLOW PADRÃO PARA AGENTES

### 1. ANTES DE CODIFICAR

```bash
# 1. Validar contexto
./scripts/tools/quick_agent_check.sh all

# 2. Preencher checklist pré-implementação
docs/implementation_plans/[feature]_pre_check.md

# 3. Analisar código similar existente
./scripts/tools/analyze_patterns.sh [handler_similar] --all

# 4. Verificar se funcionalidade já existe
find modules -name "*[feature]*" -type f
```

**Checklist Obrigatório:**
- [ ] Backend analisado e compreendido
- [ ] Padrões de frontend identificados
- [ ] Riscos mapeados e mitigados
- [ ] Decisões documentadas
- [ ] `quick_agent_check.sh` executado

---

### 2. DURANTE IMPLEMENTAÇÃO

**TDD:** Testes primeiro (Red → Green → Refactor).

**Anti-Float:** Zero `float` para valores financeiros/tempo.

**Cache-Proof:** Templates `*_simple.html` + `ParseFiles()` no handler.

**Soberania:** `entity_id` isolamento, um banco por entidade.

**Pedagogia (RF-30):** Todo campo técnico com botão "?" linkado para ajuda.

**Nunca:**
- Use `float` para valores financeiros
- Crie templates globais ou use `ParseGlob()`
- Acesse bancos de outras entidades
- Conclua tarefa sem testes passando
- Use jargão técnico na interface sem explicação

---

### 3. APÓS IMPLEMENTAÇÃO

```bash
# 1. Smoke test
./scripts/dev/smoke_test_new_feature.sh "Feature" "/rota"

# 2. Validação E2E (OBRIGATÓRIO)
./scripts/dev/validate_e2e.sh --basic --headless

# 3. Documentar aprendizados
./conclude_task.sh "Aprendizados: [resumo]" --success
```

**Critérios de Aceite:**
- [ ] Testes unitários passando (>90% cobertura para handlers)
- [ ] Smoke test passando
- [ ] **Validação E2E passando** (`validate_e2e.sh --basic --headless`)
- [ ] Aprendizados registrados via `conclude_task.sh`

---

## 📚 SKILLS DO PROJETO

Os agentes devem consultar as skills em `docs/skills/` antes de implementar:

| Skill | Foco | Arquivo |
|-------|------|---------|
| **developing-digna-backend** | Rigor técnico, DDD, TDD, Anti-Float | `skills/developing-digna-backend/SKILL.md` |
| **rendering-digna-frontend** | HTMX, UI "Soberania e Suor", Cache-Proof | `skills/rendering-digna-frontend/SKILL.md` |
| **managing-sovereign-data** | Isolamento SQLite, LifecycleManager | `skills/managing-sovereign-data/SKILL.md` |
| **applying-solidarity-logic** | Tradução cultural, ITG 2002, pedagogia | `skills/applying-solidarity-logic/SKILL.md` |
| **auditing-fiscal-compliance** | Accountant Dashboard, SPED, Read-Only | `skills/auditing-fiscal-compliance/SKILL.md` |

**Obrigatório:** Consultar skill relevante antes de implementar feature nova.

---

## 🚨 ANTIPADRÕES COMUNS

| Antipadrão | Solução |
|------------|---------|
| Importar pacote `internal` de outro módulo | Usar API layer em `pkg/` ou mock inicial |
| Handler independente do `BaseHandler` | Estender `BaseHandler` para funções de template |
| Float para valores financeiros | Usar `int64` para centavos, converter apenas para exibição |
| Template sem padrão de navegação | Atualizar TODOS os `*_simple.html` com link |
| Funções de template inconsistentes | Adicionar ao `BaseHandler` ou handler específico |
| HTMX sem feedback visual | Incluir spinner, mensagens de sucesso/erro |
| Testes dependentes de templates | Isolar testes de lógica, mockar templates |
| Testes sem TDD | Escrever teste falhando antes de implementar |
| Mock complexo desnecessário | Mock simples e focado no cenário testado |
| Implementar sem análise prévia | Fase de descoberta (30-60min) antes de codificar |
| Não documentar decisões | `docs/implementation_plans/` + `docs/learnings/` |
| Ignorar Constituição de IA | Validar contra regras sagradas antes de commit |

---

## 🔍 VALIDAÇÃO PRÉ-COMMIT

```bash
# 1. Anti-Float scan
grep -r "float[0-9]*" modules/[novo_modulo]/
# Deve retornar apenas logs/comentários

# 2. Cache-proof validation
grep -r "ParseFiles" modules/ui_web/internal/handler/[novo_handler].go
# Deve existir no handler

# 3. Soberania validation
grep -r "entity_id" modules/ui_web/internal/handler/[novo_handler].go
# Deve extrair do contexto/query

# 4. Test coverage
go test ./modules/[novo_modulo]/... -cover
# Deve ser >90% para handlers

# 5. E2E validation
./scripts/dev/validate_e2e.sh --basic --headless
# Deve passar todos os 7 passos
```

---

## 📊 MATRIZ DE RESPONSABILIDADE

| Regra | Quem Valida | Quando |
|-------|-------------|--------|
| Anti-Float | Agente + Code Review | Antes de cada commit |
| Cache-Proof | Agente + Smoke Test | Após implementação |
| Soberania | Agente + Code Review | Antes de cada commit |
| Pedagogia (RF-30) | Agente + ITCP | Antes de merge |
| E2E Validation | Agente + Usuário | Antes de concluir tarefa |
| Laicidade | PMC + Governança | Antes de release |

---

## 📝 DECISÕES ARQUITETURAIS CRÍTICAS

### 1. Emissão de NF-e/NFC-e

**Decisão:** Manter fora do Core Lume. Criar módulo `fiscal_bridge` separado.

**Justificativa:** Não acoplar ao Motor Lume para preservar essência de "Contabilidade Invisível".

**Implementação:** Integrações de terceiros (ex: eNotas, NFe.io) via módulo separado.

---

### 2. Arquitetura de Módulos

**Decisão:** Manter arquitetura atual (Go workspace, monolito modular).

**Justificativa:** Go workspace permite separação lógica sem complexidade de deploy.

**Implementação:** Extrair `indicators_engine` como módulo separado, mas mesmo binário.

---

### 3. Banco de Dados para Novos Módulos

| Módulo | Banco | Justificativa |
|--------|-------|---------------|
| indicators_engine | central.db | Dados globais, não específicos por entidade |
| financing_programs | central.db | Catálogo global de programas |
| eligibility_profiles | entity.db | Perfil específico de cada entidade |
| program_matches | entity.db | Match específico de cada entidade |
| public_profiles | entity.db + sync | Perfil público com sincronização controlada |
| help_topics | central.db | Tópicos de ajuda globais, reutilizáveis |

---

### 4. Sistema de Ajuda Educativa (RF-30)

**Decisão:** Implementar sistema de ajuda estruturada com linkagem UI → banco de ajuda.

**Justificativa:** Campos como "CadÚnico", "Inadimplência", "CNAE" são jargões burocráticos que violam o Pilar Pedagógico.

**Implementação:**
- Tabela `help_topics` no `central.db`
- Botão "?" ao lado de campos técnicos na UI
- Explicação em linguagem popular + legislação + próximo passo
- Carregamento via HTMX (< 500ms)

---

## 🎯 PRÓXIMOS PASSOS PARA AGENTES

1. **RF-27 (DAS MEI):** Cálculo automático, tabela versionada de salário mínimo
2. **RF-30 (Ajuda Educativa):** Sistema de ajuda estruturada, seed de 10+ tópicos
3. **RF-19 (Perfil de Elegibilidade):** Campos complementares, preenchimento único
4. **RF-18 (Motor de Indicadores):** Coleta BCB/IBGE, cache local, interpretação
5. **RF-20 (Portal MVP):** Match com 3 programas, checklist de documentos

**Prioridade:** RF-27 e RF-30 primeiro (baixo esforço, alto impacto, habilitam outros módulos).

---

## 📞 SUPORTE E REFERÊNCIAS

**Documentação:**
- `02_product/01_requirements.md` - Requisitos completos
- `03_architecture/01_system.md` - Arquitetura do sistema
- `06_roadmap/02_roadmap.md` - Roadmap com Fases 1-5
- `06_roadmap/03_backlog.md` - Backlog detalhado

**Skills:**
- `docs/skills/` - 5 skills específicas do projeto

**Scripts:**
- `./scripts/tools/quick_agent_check.sh` - Validação rápida
- `./scripts/dev/smoke_test_new_feature.sh` - Smoke test
- `./scripts/dev/validate_e2e.sh` - Validação E2E

**Aprendizados:**
- `docs/learnings/` - Aprendizados de sessões anteriores
- `docs/ANTIPATTERNS.md` - Antipadrões e soluções

---

**Status:** ✅ ATUALIZADO COM ECOSSISTEMA DE 4 MÓDULOS (PDF v1.0) + RF-30 (Sessão 27/03/2026)  
**Próxima Ação:** Agentes devem consultar esta constituição antes de cada implementação  
**Versão Anterior:** 1.0 (2026-03-13)  
**Versão Atual:** 2.0 (2026-03-27)
