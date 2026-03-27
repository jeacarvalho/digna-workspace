title: Architectural Decisions Record (ADR) - Ecossistema Digna
status: implemented
version: 3.0
last_updated: 2026-03-27
---

# Architectural Decisions Record - Ecossistema Digna

> **Nota:** Este documento consolida todas as decisões arquiteturais tomadas desde a Sprint 12, incluindo a expansão para o Ecossistema de 4 Módulos (PDF v1.0) e as decisões da sessão de 27/03/2026 (RF-30 - Sistema de Ajuda Educativa).

---

## 📋 Contexto da Atualização (27/03/2026)

**Motivação:** O projeto Digna evoluiu de um ERP contábil para um **Ecossistema de 4 Módulos** conforme especificação PDF v1.0. Esta atualização documenta:

1. **Decisões da Expansão do Ecossistema** (PDF v1.0, Março 2026)
   - Arquitetura de 4 módulos interdependentes
   - Princípio "Nenhum dado digitado duas vezes"
   - Separação banco central vs. banco por entidade

2. **Decisões da Sessão 27/03/2026**
   - RF-30: Sistema de Ajuda Educativa Estruturada
   - Linkagem UI → banco de ajuda (`help_topics{}`)
   - Tópicos em linguagem popular (5ª série)

3. **Preservação das Decisões Sprint 12** (ADR-001 a ADR-004)
   - Todas as decisões originais sobre Accountant Dashboard mantidas
   - Status atualizado conforme implementação real

**Versão Anterior:** 1.0 (2026-03-08)  
**Nova Versão:** 3.0 (2026-03-27)

---

## 🏛️ ADRs da Expansão do Ecossistema (PDF v1.0)

### ADR-005: Arquitetura de 4 Módulos Interdependentes

**Data:** 27/03/2026  
**Status:** ✅ ACEITO  
**Prioridade:** ALTA

#### Contexto
O Digna era documentado como um ERP único. O PDF v1.0 especifica um ecossistema de 4 módulos que compartilham dados sem duplicação de entrada.

#### Decisão
Implementar arquitetura de 4 módulos interdependentes:
1. **digna ERP** (núcleo) - Gestão financeira, fiscal e contábil
2. **Motor de Indicadores** - Coleta APIs externas (BCB, IBGE)
3. **Portal de Oportunidades** - Match automático de crédito
4. **Rede Digna** - Marketplace solidário B2B

#### Consequências
**Positivas:**
- Dados do ERP alimentam automaticamente Portal e Rede
- Elimina formulários redundantes (princípio central)
- Escalabilidade modular (cada módulo evolui independentemente)

**Negativas:**
- Complexidade de integração entre módulos
- Dependência forte do ERP (se ERP falha, todos afetados)
- Necessidade de versionamento de API interna entre módulos

#### Validação
- ✅ Modelo de dados `EligibilityProfile` copia dados do `Enterprise`
- ✅ `ProgramMatch` consome `EligibilityProfile` + `EconomicIndicator`
- ✅ `PublicProfile` deriva dados do `Enterprise` + campos públicos

---

### ADR-006: Separação Banco Central vs. Banco por Entidade

**Data:** 27/03/2026  
**Status:** ✅ ACEITO  
**Prioridade:** CRÍTICA

#### Contexto
Novos módulos requerem dados globais (indicadores, programas de crédito, tópicos de ajuda) e dados específicos por entidade (perfil, matches, transações).

#### Decisão
Manter e expandir dualidade arquitetural:
- **`central.db`**: Dados globais (indicadores, programas, help_topics, vínculos contábeis)
- **`data/entities/{entity_id}.db`**: Dados operacionais por entidade (ledger, perfil, matches)

#### Consequências
**Positivas:**
- Soberania de dados preservada (entidade leva seu `.db`)
- Dados globais atualizados uma vez, consumidos por todos
- Isolamento total entre entidades (proibido JOIN entre bancos)

**Negativas:**
- Complexidade de sincronização entre bancos
- Backup requer dois arquivos (central + entidade)
- Queries distribuídas mais complexas

#### Validação
- ✅ `help_topics` em `central.db` (tópicos globais)
- ✅ `eligibility_profiles` em `entity.db` (perfil específico)
- ✅ `economic_indicators` em `central.db` (dados externos)

---

### ADR-007: Princípio "Nenhum Dado Digitado Duas Vezes"

**Data:** 27/03/2026  
**Status:** ✅ ACEITO  
**Prioridade:** ALTA

#### Contexto
Formulários redundantes são a principal barreira de adoção para usuários de baixa escolaridade.

#### Decisão
Implementar princípio arquitetural central:
- ERP captura dados no uso cotidiano (vendas, compras, caixa)
- Portal e Rede **consomem** dados do ERP, nunca exigem reentrada
- Campos complementares são **preenchimento único** (ex: CadÚnico, gênero)

#### Consequências
**Positivas:**
- Redução de fricção no onboarding
- Dados mais consistentes (única fonte da verdade)
- Usuário percebe valor imediato (sem formulários longos)

**Negativas:**
- Dependência forte da qualidade dos dados do ERP
- Migração complexa se ERP mudar estrutura
- Validação de dados mais crítica (erro se propaga)

#### Validação
- ✅ `EligibilityProfile` tem campos "copiados do ERP" + "complementares"
- ✅ `ProgramMatch` não exige novos dados do usuário
- ✅ `PublicProfile` deriva dados existentes

---

## 🎓 ADRs do Sistema de Ajuda Educativa (RF-30)

### ADR-008: Sistema de Ajuda Estruturada com Linkagem UI → Banco

**Data:** 27/03/2026  
**Status:** ✅ ACEITO  
**Prioridade:** ALTA

#### Contexto
Campos como "CadÚnico", "Inadimplência", "CNAE" são jargões burocráticos que violam o Pilar Pedagógico do Digna. Tooltips simples são insuficientes.

#### Decisão
Implementar sistema de ajuda estruturada:
- Tabela `help_topics` no `central.db`
- Botão "?" ao lado de campos técnicos na UI
- Linkagem via chave única (ex: `/help/topic/cadunico`)
- Conteúdo: explicação popular + legislação + próximo passo

#### Consequências
**Positivas:**
- Redução de abandono em formulários
- Empoderamento do usuário (aprende enquanto opera)
- Conteúdo centralizado (atualização única, uso múltiplo)

**Negativas:**
- Overhead de manutenção de conteúdo
- Risco de conteúdo desatualizado
- Performance (query adicional por campo técnico)

#### Validação
- ✅ Tabela `help_topics` com estrutura completa
- ✅ Handler `/help/topic/{key}` implementado
- ✅ Template com botão "?" linkado

---

### ADR-009: Linguagem Popular (5ª Série) para Conteúdo de Ajuda

**Data:** 27/03/2026  
**Status:** ✅ ACEITO  
**Prioridade:** CRÍTICA

#### Contexto
Usuários do Digna têm baixa escolaridade. Conteúdo técnico causa abandono e frustração.

#### Decisão
Estabelecer critério de qualidade obrigatório:
- **Linguagem:** Usuário com 5ª série entende sem ajuda externa
- **Jargão:** Zero termos técnicos sem explicação
- **Ação:** Sempre inclui "próximo passo" acionável (ex: "procure o CRAS")
- **Validação:** Teste de usabilidade com usuários reais

#### Consequências
**Positivas:**
- Inclusão digital real (não apenas retórica)
- Redução de dependência de intermediários
- Alinhamento com Pilar 2 (Tradução Cultural)

**Negativas:**
- Conteúdo mais longo (explicações detalhadas)
- Revisão por ITCPs/comunidade necessária
- Atualização mais lenta (validação humana)

#### Validação
- ✅ Seed de 6 tópicos com linguagem popular
- ✅ Critério no checklist de aceite do RF-30
- ✅ Processo de revisão documentado

---

### ADR-010: Cache de Tópicos de Ajuda com Invalidação por Atualização

**Data:** 27/03/2026  
**Status:** ✅ ACEITO  
**Prioridade:** MÉDIA

#### Contexto
Tópicos de ajuda são lidos frequentemente, mas atualizados raramente. Query direta ao banco a cada acesso é ineficiente.

#### Decisão
Implementar cache com invalidação:
- Cache em memória (1h TTL)
- Invalidação automática ao atualizar tópico
- Fallback para banco se cache falhar

#### Consequências
**Positivas:**
- Performance (< 500ms para tooltip)
- Redução de carga no `central.db`
- Experiência do usuário mais fluida

**Negativas:**
- Complexidade de invalidação
- Risco de conteúdo desatualizado no cache
- Memória adicional por instância

#### Validação
- ✅ Tabela `cache_help_topics` no schema
- ✅ Invalidação documentada no protocolo

---

## 📊 ADRs da Sprint 12 (Accountant Dashboard) - PRESERVADOS

### ADR-001: Integration via ui_web/main.go

**Data:** 08/03/2026  
**Status:** ✅ IMPLEMENTADO  
**Prioridade:** ALTA

#### Contexto
The Sprint 12 prompt suggested creating a new entry point at `cmd/digna/main.go` for the Accountant Dashboard module. However, the project already has a well-established web interface architecture centered around the `ui_web` module.

#### Decision
Integrate the Accountant Dashboard through the existing `modules/ui_web/main.go` instead of creating a new entry point.

#### Consequences
**Positive:**
- Maintains architectural consistency across all web interfaces
- Centralizes HTTP route management
- Simplifies deployment (single binary for all web interfaces)
- Follows DRY principle by reusing existing infrastructure

**Negative:**
- Deviates from the original prompt specification
- Creates tighter coupling between `accountant_dashboard` and `ui_web` modules

#### Current Status
✅ **IMPLEMENTED** - All handlers registered in `ui_web/main.go`

---

### ADR-002: Embedded Templates vs Separate HTML Files

**Data:** 08/03/2026  
**Status:** ✅ IMPLEMENTADO (com evolução)  
**Prioridade:** ALTA

#### Contexto
The prompt suggested creating separate HTML template files (`layout.html`, `dashboard.html`). The project has examples of both embedded templates and separate template files.

#### Decision
Use embedded Go templates within the handler code instead of separate HTML files.

#### Evolution (Sprint 16)
Na Sprint 16, evoluímos para **templates `*_simple.html` carregados do disco** via `ParseFiles()` no handler (cache-proof), mantendo a filosofia de simplicidade mas resolvendo problemas de cache persistente.

#### Consequences
**Positive:**
- Simplifies deployment (fewer files to manage)
- Improves performance (templates compiled with binary) - **Sprint 12**
- Zero cache issues (loaded from disk) - **Sprint 16**
- Easier testing (templates tested with handler code)

**Negative:**
- Less separation of concerns between logic and presentation
- Harder for non-developers to modify templates

#### Current Status
✅ **IMPLEMENTED** - All templates use `*_simple.html` + `ParseFiles()` pattern

---

### ADR-003: No Separate templates/ Directory

**Data:** 08/03/2026  
**Status:** ✅ IMPLEMENTADO  
**Prioridade:** BAIXA

#### Contexto
The prompt suggested creating a `templates/` directory within the `accountant_dashboard` module.

#### Decision
Do not create a separate `templates/` directory since templates are embedded in the code.

#### Consequences
**Positive:**
- Reduces project complexity
- Follows YAGNI principle (no unnecessary structure)
- Consistent with other modules using embedded templates

**Negative:**
- Deviates from prompt specification
- Less conventional for web development

#### Current Status
✅ **IMPLEMENTED** - Templates in `modules/ui_web/templates/`

---

### ADR-004: Public API Package Structure

**Data:** 08/03/2026  
**Status:** ✅ IMPLEMENTADO  
**Prioridade:** MÉDIA

#### Contexto
The prompt focused on internal implementation but didn't specify external API design.

#### Decision
Create a public API package (`pkg/dashboard/`) for external consumption of the Accountant Dashboard functionality.

#### Consequences
**Positive:**
- Enables integration with other systems
- Provides clean separation between internal and external APIs
- Follows Go best practices for package design

**Negative:**
- Additional complexity beyond prompt requirements
- More interfaces to maintain

#### Current Status
✅ **IMPLEMENTED** - `pkg/dashboard/` com interfaces públicas

---

## 🔗 ADRs de Integração entre Módulos

### ADR-011: API Interna entre Módulos do Ecossistema

**Data:** 27/03/2026  
**Status:** ✅ ACEITO  
**Prioridade:** ALTA

#### Contexto
Módulos do ecossistema precisam se comunicar (ERP → Portal, Motor → Portal, etc.) sem acoplamento forte.

#### Decisão
Implementar API interna via interfaces Go:
- Cada módulo expõe interface pública em `pkg/`
- Consumidores dependem de interfaces, não implementações
- Injeção de dependência no `main.go`

#### Consequências
**Positivas:**
- Baixo acoplamento entre módulos
- Testabilidade (mock de interfaces)
- Evolução independente de cada módulo

**Negativas:**
- Overhead de interfaces
- Complexidade de versionamento de API
- Necessidade de documentação clara

#### Validação
- ✅ `indicators_engine/pkg/indicators/` com API pública
- ✅ `portal_opportunities` consome via interface

---

### ADR-012: Versionamento de Schema entre Módulos

**Data:** 27/03/2026  
**Status:** ✅ ACEITO  
**Prioridade:** MÉDIA

#### Contexto
Módulos evoluem em ritmos diferentes. Mudanças no schema do ERP podem quebrar Portal ou Rede.

#### Decisão
Implementar versionamento de schema:
- Cada módulo declara versão mínima do ERP
- Migrações são backward-compatible
- API de descoberta de versão

#### Consequências
**Positivas:**
- Evita breaking changes
- Deploy independente de módulos
- Rollback seguro

**Negativas:**
- Complexidade de migrações
- Overhead de compatibilidade
- Documentação de versão necessária

#### Validação
- ✅ Tabela `schema_versions` no `central.db`
- ✅ Migrações idempotentes

---

## 📈 Matriz de Decisões por Prioridade

| ADR | Título | Prioridade | Status | Impacto |
|-----|--------|------------|--------|---------|
| ADR-005 | 4 Módulos Interdependentes | ALTA | ✅ Aceito | Arquitetura completa |
| ADR-006 | Banco Central vs. Entidade | CRÍTICA | ✅ Aceito | Soberania de dados |
| ADR-007 | Nenhum Dado Duas Vezes | ALTA | ✅ Aceito | UX e adoção |
| ADR-008 | Ajuda Estruturada | ALTA | ✅ Aceito | Pedagogia |
| ADR-009 | Linguagem 5ª Série | CRÍTICA | ✅ Aceito | Inclusão |
| ADR-010 | Cache de Ajuda | MÉDIA | ✅ Aceito | Performance |
| ADR-001 | Integration via ui_web | ALTA | ✅ Implementado | Consistência |
| ADR-002 | Embedded Templates | ALTA | ✅ Implementado | Cache-proof |
| ADR-003 | No templates/ Dir | BAIXA | ✅ Implementado | Simplicidade |
| ADR-004 | Public API Package | MÉDIA | ✅ Implementado | Integração |
| ADR-011 | API Interna entre Módulos | ALTA | ✅ Aceito | Acoplamento |
| ADR-012 | Versionamento de Schema | MÉDIA | ✅ Aceito | Evolução |

---

## 🎯 Princípios Aplicados em Todas as Decisões

### 1. KISS (Keep It Simple)
- Escolher solução mais simples que atende requisitos
- Evitar complexidade desnecessária (ex: ADR-003)

### 2. YAGNI (You Ain't Gonna Need It)
- Não implementar estrutura não imediatamente necessária
- Templates directory não criado quando embutidos funcionam

### 3. DRY (Don't Repeat Yourself)
- Reutilizar infraestrutura existente (ex: ADR-001)
- Dados do ERP reaproveitados por Portal e Rede (ADR-007)

### 4. Soberania de Dados
- Isolamento físico por entidade (ADR-006)
- Exit Power preservado em todas as decisões

### 5. Pedagogia e Inclusão
- Linguagem popular obrigatória (ADR-009)
- Ajuda contextual em todos os campos técnicos (ADR-008)

---

## 🔄 Processo de Revisão de ADRs

### Quando Revisar
- **Mensalmente:** Revisar ADRs de média/baixa prioridade
- **Por Feature:** Validar ADRs relacionados antes de implementar
- **Trimestral:** Revisão completa com PMC

### Como Atualizar
1. Criar novo ADR com número sequencial
2. Atualizar status de ADRs existentes se necessário
3. Documentar lições aprendidas em `docs/learnings/`
4. Comunicar mudanças em RFC se impactar arquitetura

### Critérios de Aceite para Novo ADR
- [ ] Contexto claramente documentado
- [ ] Decisão específica e acionável
- [ ] Consequências positivas e negativas listadas
- [ ] Validação técnica ou de negócio
- [ ] Aprovado por PMC (para ADRs críticos)

---

## 📚 Referências

### Documentos Relacionados
- `03_architecture/01_system.md` - Arquitetura do sistema
- `03_architecture/02_protocols.md` - Protocolos de integração
- `02_product/01_requirements.md` - Requisitos funcionais
- `docs/06_roadmap/02_roadmap.md` - Roadmap do ecossistema

### Sessões de Decisão
- **Sprint 12** (08/03/2026): ADR-001 a ADR-004 (Accountant Dashboard)
- **Sessão 27/03/2026**: ADR-005 a ADR-012 (Ecossistema + RF-30)

### Skills Aplicadas
- `developing-digna-backend` - Clean Architecture, DDD
- `managing-sovereign-data` - Isolamento SQLite
- `applying-solidarity-logic` - Pedagogia e inclusão
- `rendering-digna-frontend` - Cache-proof templates

---

**Status:** ✅ ATUALIZADO COM DECISÕES DO ECOSSISTEMA (PDF v1.0) + RF-30 (Sessão 27/03/2026)  
**Próxima Ação:** Atualizar `03_architecture/05_database_system.md` com schema expandido  
**Versão Anterior:** 1.0 (2026-03-08)  
**Versão Atual:** 3.0 (2026-03-27)
