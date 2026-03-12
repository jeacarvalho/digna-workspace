# 📋 TAREFA: Concluir RF-12 - Integração Temporal Filtering no UI

**Data:** 11/03/2026
**Prioridade:** ALTA
**Estimativa:** 4-6 horas
**Módulo(s):** ui_web, lifecycle, accountant_dashboard
**Tarefa ID:** 20260311_190332

---

## 🎯 OBJETIVO

Concluir a implementação da RF-12 (Gestão de Vínculo Contábil e Delegação Temporal) integrando o temporal filtering no sistema UI após perda de contexto por compaction.

---

## 📋 REQUISITOS

### Funcionais
- [ ] **RF-12.1:** Store EnterpriseAccountant relationships in Central Database (central.db)
- [ ] **RF-12.2:** Implement Exit Power - cooperatives can terminate relationships
- [ ] **RF-12.3:** Enforce temporal cardinality - only 1 active accountant per cooperative
- [ ] **RF-12.4:** Provide temporal access filtering for inactive accountants
- [ ] **RF-12.5:** Integrate temporal filtering in AccountantHandler UI

### Técnicos
- [ ] Seguir padrões do projeto Digna (Clean Architecture)
- [ ] Implementar testes de integração
- [ ] Atualizar documentação RF-12
- [ ] Validar com build completo
- [ ] Usar Anti-Float rule (int64 timestamps)

### Não Funcionais
- [ ] Performance: Temporal filtering não deve impactar performance
- [ ] Segurança: Apenas delegated_by pode desativar vínculo (Exit Power)
- [ ] Usabilidade: UI intuitiva para gerenciamento de vínculos

---

## 🔍 CONTEXTO E ANÁLISE

### 🚨 CONTEXTO DE COMPACTION
O opencode entrou em modo compaction e perdeu o contexto. Já implementamos 70% da RF-12:

### ✅ JÁ CONCLUÍDO:
1. **Extended LifecycleManager** para suportar Central Database (`central.db`)
2. **Created Domain Entity** `EnterpriseAccountant` com regras de negócio
3. **Implemented Central Repository** para `central.db`
4. **Built Business Service** com todas as regras RF-12
5. **Created Temporal Filtering Middleware** para `accountant_dashboard`
6. **Updated Translator Service** para filtrar entidades por relacionamentos
7. **Extended SQLiteManager** para implementar interface `AccountantLinkService`

### 🔄 PRÓXIMOS PASSOS (FOCOS ATUAIS):
1. **Integrar temporal filtering no `AccountantHandler` do `ui_web`**
2. **Criar UI para gerenciar vínculos contábeis**
3. **Testar integração completa**

### Módulos/Arquivos Relacionados
- `modules/lifecycle/pkg/lifecycle/interfaces.go` - Interface AccountantLinkService
- `modules/lifecycle/pkg/lifecycle/sqlite.go` - Implementação SQLiteManager
- `modules/accountant_dashboard/internal/middleware/temporal_filter.go` - Middleware
- `modules/ui_web/internal/handler/accountant_handler.go` - PRECISA SER ATUALIZADO

### Padrões a Seguir
- [ ] Analisar handler similar: `modules/ui_web/internal/handler/auth_handler.go`
- [ ] Analisar template similar: `modules/ui_web/templates/accountant_dashboard_simple.html`
- [ ] Seguir padrão Anti-Float (int64 timestamps)
- [ ] Usar anti-padrões de `docs/ANTIPATTERNS.md`

### Dependências
- [ ] `LifecycleManager` já implementa `AccountantLinkService`
- [ ] Middleware de temporal filtering já existe
- [ ] Repository com queries temporais já implementado

---

## 🚀 PLANO DE IMPLEMENTAÇÃO

### Fase 1: Integrar Temporal Filtering no AccountantHandler (2-3 horas)
1. [ ] **Modificar `accountant_handler.go`:**
   - Obter accountant ID do contexto da sessão (via `AuthHandler`)
   - Usar `LifecycleManager.GetValidEnterprisesForAccountant()`
   - Filtrar lista de entidades pendentes
   - Validar acesso em `ExportFiscal`

2. [ ] **Testar integração:**
   - Build do módulo `ui_web`
   - Testar acesso do contador
   - Verificar filtragem temporal

### Fase 2: Criar UI para Gerenciamento de Vínculos (1-2 horas)
1. [ ] **Criar `accountant_link_handler.go`:**
   - Handler para criar vínculos
   - Handler para desativar vínculos (Exit Power)
   - Listar vínculos ativos/inativos

2. [ ] **Criar template `accountant_link_simple.html`:**
   - Formulário para criar vínculo
   - Lista de vínculos com ações
   - Validação client-side

3. [ ] **Registrar handler no `main.go`**

### Fase 3: Testes e Validação (1 hora)
1. [ ] **Testar cardinality rule:** Apenas 1 contador ativo por cooperativa
2. [ ] **Testar Exit Power:** Apenas cooperativa pode desativar vínculo
3. [ ] **Testar temporal filtering:** Contador só vê entidades do período ativo
4. [ ] **Build completo:** Todos os módulos compilam

---

## 📁 ARQUIVOS ESPERADOS

### A Criar
- `modules/ui_web/internal/handler/accountant_link_handler.go`
- `modules/ui_web/templates/accountant_link_simple.html`
- `modules/ui_web/internal/handler/accountant_link_handler_test.go`

### A Modificar
- `modules/ui_web/internal/handler/accountant_handler.go` (integração temporal filtering)
- `modules/ui_web/main.go` (registrar handler)
- `docs/QUICK_REFERENCE.md` (atualizar referência RF-12)
- `docs/NEXT_STEPS.md` (marcar RF-12 como concluída)

---

## ⚠️ RISCOS E DESAFIOS

### Riscos Técnicos
1. **Autenticação:** O `AccountantHandler` precisa acessar `AuthHandler` para obter accountant ID
   - **Mitigação:** Verificar como outros handlers fazem (ex: `dashboard_handler.go`)

2. **Sessão:** Informações do usuário devem estar no contexto da request
   - **Mitigação:** Usar `AuthHandler.GetCurrentEntity()` e `GetCurrentUserType()`

3. **Integração:** Múltiplos módulos precisam trabalhar juntos
   - **Mitigação:** Testar build incremental de cada módulo

### Riscos de Processo
1. **Compaction:** opencode pode entrar em compaction novamente
   - **Mitigação:** Usar `./preserve_context.sh --save` antes

2. **Contexto perdido:** Já perdemos contexto uma vez
   - **Mitigação:** Documentar tudo em `session_learnings/`

---

## 📚 APRENDIZADOS ANTERIORES RELEVANTES

### Preservação de Contexto Durante Compaction
- **Arquivo:** `work_in_progress/current_session/session_learnings/COMPACTION_CONTEXT_PRESERVATION.md`
- **Aprendizado:** Sempre usar `./preserve_context.sh --save` antes do compaction

### Implementação RF-12 Parcial
- **Arquivo:** `work_in_progress/current_session/.compaction_context.md`
- **Aprendizado:** 70% da RF-12 já implementado, foco em integração UI

---

## 🔗 LINKS ÚTEIS

- [Documentação do projeto](docs/)
- [Padrões de código](docs/QUICK_REFERENCE.md)
- [Antipadrões](docs/ANTIPATTERNS.md)
- [Skills do projeto](docs/skills/)
- [Contexto preservado](work_in_progress/current_session/.compaction_context.md)

---

**Status:** EM ANDAMENTO
**Última atualização:** 11/03/2026 19:05

---

## 🤖 INSTRUÇÃO PARA OPENCODE

**LEIA PRIMEIRO ESTES ARQUIVOS:**
1. `work_in_progress/current_session/.compaction_context.md` - Contexto preservado
2. `work_in_progress/current_session/session_learnings/COMPACTION_CONTEXT_PRESERVATION.md` - Aprendizado sobre compaction

**ARQUIVOS JÁ IMPLEMENTADOS (VERIFICAR):**
1. `modules/accountant_dashboard/internal/middleware/temporal_filter.go`
2. `modules/lifecycle/pkg/lifecycle/sqlite.go`
3. `modules/lifecycle/internal/service/accountant_link_service.go`

**PRÓXIMO PASSO IMEDIATO:**
Modificar `modules/ui_web/internal/handler/accountant_handler.go` para usar temporal filtering.

**COMANDOS DE VALIDAÇÃO:**
```bash
# Build todos os módulos
cd modules && ./run_tests.sh

# Build específico
cd modules/ui_web && go build ./...
cd modules/lifecycle && go build ./...
cd modules/accountant_dashboard && go build ./...
```