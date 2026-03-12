# 🚀 PLANO DE IMPLEMENTAÇÃO: Gestão de Vínculo Contábil e Delegação Temporal (RF-12)
# ID: 20260311_180051
# Gerado em: 11/03/2026

## 🎯 OBJETIVO DA TAREFA
Implementar o RF-12 criando a entidade `EnterpriseAccountant` que controla qual Contador Social tem acesso a qual Empreendimento (Tenant) e por qual período. Como esta é uma relação inter-tenant de governança global, o registro deve ocorrer obrigatoriamente no **Banco Central** do sistema (gerenciado pelo módulo `lifecycle`), protegendo a arquitetura de *micro-databases*.

## 📋 REQUISITOS ESPECÍFICOS DA RF-12

### Funcionais
- [ ] **RF-12.1:** Criar entidade `EnterpriseAccountant` no Banco Central
- [ ] **RF-12.2:** Implementar regra "Exit Power" - cooperativa pode encerrar vínculo
- [ ] **RF-12.3:** Garantir cardinalidade temporal (1 contador ATIVO por cooperativa)
- [ ] **RF-12.4:** Implementar filtro de acesso temporal para contadores desativados
- [ ] **RF-12.5:** Contadores desativados mantêm acesso Read-Only ao período de vigência

### Técnicos (Baseado nas Skills)
- [ ] **Anti-Float:** Usar `int64` para timestamps (milissegundos desde Unix epoch)
- [ ] **Soberania:** Armazenar no Banco Central (`central.db`), NÃO nos bancos dos tenants
- [ ] **Clean Architecture:** Seguir DDD (Domain → Service → Repository)
- [ ] **TDD:** Implementar testes unitários antes do código
- [ ] **Soft Delete:** Manter histórico com `EndDate` em vez de deletar fisicamente

### Não Funcionais
- [ ] **Performance:** Consultas eficientes por `EnterpriseID` e `AccountantID`
- [ ] **Segurança:** Validação de acesso baseada em período de vigência
- [ ] **Auditoria:** Manter histórico completo de vínculos
- [ ] **Consistência:** Transações atômicas para atualização de status

## 🔍 ANÁLISE DE CONTEXTO

### Módulos/Arquivos Relacionados
- `modules/lifecycle/` - Gerenciamento do Banco Central e LifecycleManager
- `modules/accountant_dashboard/` - Dashboard do contador e filtros de acesso
- `modules/ui_web/internal/handler/accountant_handler.go` - Handler existente do contador

### Padrões a Seguir
- [ ] Analisar `modules/lifecycle/internal/domain/entity.go` para padrões de entidade
- [ ] Analisar `modules/lifecycle/internal/repository/migration.go` para padrões de migração
- [ ] Seguir padrão de `Status` (ACTIVE/INACTIVE) como em outros módulos
- [ ] Usar `*time.Time` para campos opcionais (EndDate nulo quando ativo)

### Dependências
- [ ] Banco Central (`central.db`) deve existir ou ser criado
- [ ] LifecycleManager deve suportar acesso ao Banco Central
- [ ] Integração futura com middleware do `accountant_dashboard`

## 🔄 FLUXO DE IMPLEMENTAÇÃO DETALHADO

### Fase 1: Análise e Modelagem de Domínio (1 hora)
1. **Análise da estrutura atual do `lifecycle`** (30 min)
   - Estudar `entity.go` e `interfaces.go`
   - Entender padrões de migração SQLite
   - Identificar como acessar Banco Central via LifecycleManager

2. **Modelagem da entidade `EnterpriseAccountant`** (30 min)
   - Definir struct com campos: `ID`, `EnterpriseID`, `AccountantID`, `Status`, `StartDate`, `EndDate`, `DelegatedBy`
   - Definir constantes: `StatusActive`, `StatusInactive`
   - Definir métodos de validação e business rules

### Fase 2: Implementação do Repositório (2 horas)
3. **Criação do repositório central** (1 hora)
   - `modules/lifecycle/internal/repository/accountant_link_repo.go`
   - Implementar interface `EnterpriseAccountantRepository`
   - Métodos: `Create`, `Update`, `FindByEnterpriseID`, `FindByAccountantID`, `FindActiveByEnterpriseID`
   - **CRÍTICO:** Acesso exclusivo ao Banco Central (`central.db`)

4. **Migrações SQLite** (1 hora)
   - Criar migration para tabela `enterprise_accountants`
   - Definir índices: `(enterprise_id, status)`, `(accountant_id, status)`
   - Implementar rollback seguro

### Fase 3: Implementação do Serviço (2 horas)
5. **Serviço de negócio** (1.5 horas)
   - `modules/lifecycle/internal/service/accountant_link_service.go`
   - Implementar regra "Exit Power" (apenas cooperativa pode encerrar)
   - Implementar cardinalidade (1 contador ATIVO por cooperativa)
   - Implementar filtro temporal `GetValidDateRangeForAccountant()`

6. **Validações e tratamento de erros** (30 min)
   - Validação de datas (StartDate <= EndDate quando inativo)
   - Validação de status transitions
   - Erros descritivos com contexto

### Fase 4: Testes Unitários (1.5 horas)
7. **Testes do repositório** (45 min)
   - Testes com banco SQLite em memória
   - Testar todos os métodos CRUD
   - Testar constraints e índices

8. **Testes do serviço** (45 min)
   - Testar regra de cardinalidade
   - Testar "Exit Power"
   - Testar filtro temporal
   - Testar edge cases (datas inválidas, status inválidos)

### Fase 5: Integração e Validação (1 hora)
9. **Integração com sistema existente** (30 min)
   - Verificar que não quebra builds existentes
   - Testar com LifecycleManager real
   - Validar acesso ao Banco Central

10. **Validação de critérios de aceite** (30 min)
    - Banco Central isolado dos tenants ✓
    - Cardinalidade temporal respeitada ✓
    - Soft Delete com histórico ✓
    - Clean Architecture seguida ✓

## 📁 ESTRUTURA DE ARQUIVOS ESPERADA

### A Criar no módulo `lifecycle`:
```
modules/lifecycle/internal/domain/enterprise_accountant.go
modules/lifecycle/internal/domain/enterprise_accountant_test.go
modules/lifecycle/internal/repository/accountant_link_repo.go
modules/lifecycle/internal/repository/accountant_link_repo_test.go
modules/lifecycle/internal/repository/migrations/003_enterprise_accountants.sql
modules/lifecycle/internal/service/accountant_link_service.go
modules/lifecycle/internal/service/accountant_link_service_test.go
```

### A Modificar:
```
modules/lifecycle/internal/repository/migration.go (adicionar nova migration)
modules/lifecycle/pkg/lifecycle/sqlite.go (adicionar acesso ao central.db se necessário)
```

### Para integração futura (nesta tarefa apenas preparação):
```
modules/accountant_dashboard/internal/service/translator_service.go (futura injeção do filtro)
modules/ui_web/internal/middleware/accountant_auth.go (futuro middleware temporal)
```

## ⚠️ RISCOS E MITIGAÇÕES

### Riscos Técnicos:
1. **Acesso ao Banco Central não implementado** - Mitigar: Verificar se `LifecycleManager` já suporta `central.db`
2. **Conflito com migrações existentes** - Mitigar: Usar versionamento incremental (003_...)
3. **Performance de consultas** - Mitigar: Índices otimizados e query planning

### Riscos de Arquitetura:
1. **Vazamento de lógica para tenants** - Mitigar: Validação rigorosa que NADA é salvo nos bancos dos tenants
2. **Dependência circular com outros módulos** - Mitigar: Interfaces bem definidas e injeção de dependência

### Riscos de Negócio:
1. **"Exit Power" mal implementado** - Mitigar: Testes específicos para validação de permissões
2. **Filtro temporal ineficiente** - Mitigar: Benchmark de queries e otimização

## 🎯 CRITÉRIOS DE ACEITAÇÃO (Definition of Done)

### Funcionais:
- [ ] Entidade `EnterpriseAccountant` criada no Banco Central
- [ ] Regra: 1 contador ATIVO por cooperativa (cardinalidade)
- [ ] "Exit Power": Cooperativa pode inativar contador (preenche EndDate)
- [ ] Soft Delete: Histórico mantido para auditoria
- [ ] Filtro temporal `GetValidDateRangeForAccountant()` implementado

### Técnicos:
- [ ] 100% testes unitários passando
- [ ] Código segue Clean Architecture (Domain → Service → Repository)
- [ ] Anti-Float respeitado (int64 para timestamps)
- [ ] Soberania respeitada (apenas Banco Central)
- [ ] Migrações SQLite versionadas e rollback seguras

### Qualidade:
- [ ] Sem regressões no build atual
- [ ] Documentação em código (godoc)
- [ ] Tratamento de erros descritivo
- [ ] Performance aceitável (consultas < 100ms)

## 📊 ESTIMATIVA DE TEMPO

**Total estimado:** 7.5 horas
**Buffer recomendado:** 1.5 horas (20%)

### Breakdown detalhado:
- Fase 1 (Análise): 1 hora
- Fase 2 (Repositório): 2 horas
- Fase 3 (Serviço): 2 horas
- Fase 4 (Testes): 1.5 horas
- Fase 5 (Integração): 1 hora

## 🔧 PRÓXIMOS PASSOS

1. **Executar implementação:** `./process_task.sh --task=20260311_180051 --execute`
2. **Seguir TDD:** Criar testes primeiro, depois implementação
3. **Validar com skills:** Consultar `developing-digna-backend` e `managing-sovereign-data`
4. **Documentar aprendizados:** Ao final, registrar em `session_learnings/`

---

**Status do Plano:** ✅ ESPECÍFICO PARA RF-12
**Complexidade:** ALTA (integração com arquitetura de micro-databases)
**Risco:** MÉDIO (crítico para soberania de dados)
**Próximo passo:** Executar implementação seguindo este plano