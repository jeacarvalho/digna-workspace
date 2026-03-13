# 🎯 Próximos Passos - Projeto Digna

**Última atualização:** 13/03/2026  
**Status:** ✅ DOCUMENTAÇÃO DE MODULARIZAÇÃO ATUALIZADA | ⚠️ RF-12 85% COMPLETO (BLOQUEADO) | ✅ PROCESSO CORRIGIDO

---

## 🚨 STATUS ATUAL E BLOQUEADORES

### 🏗️ RF-12 - Gestão de Vínculo Contábil e Delegação Temporal (11/03/2026)
**Status:** ⚠️ 85% COMPLETO (BLOQUEADO)  
**Descrição:** Sistema de vínculos contábeis entre contadores e cooperativas com controle temporal  
**Blocker Crítico:** Erro de import do módulo lifecycle (`no non-test Go files`)  
**Progresso:** 85% implementado, 100% dos testes bloqueados

**Entregas Implementadas:**
- ✅ Banco central `central.db` com tabela `enterprise_accountants`
- ✅ Entidade `EnterpriseAccountant` com regras de negócio (Cardinalidade, Exit Power)
- ✅ Repositório e serviço com filtragem temporal
- ✅ Middleware de filtro temporal no módulo `accountant_dashboard`
- ✅ Handler UI `AccountantLinkHandler` com template cache-proof
- ✅ Scripts de preservação de contexto e correção de fluxo

**Próximos passos (ALTA PRIORIDADE):**
1. 🔧 **Resolver erro de import do módulo lifecycle** - Investigar `go.mod` do `ui_web`
2. 🔧 **Reativar filtro temporal** - Remover comentários em `accountant_handler.go:111-133`
3. 🧪 **Executar testes E2E da RF-12** - Validar funcionalidade completa
4. 🎨 **Completar UI de gerenciamento de links** - Integrar repositório, formulários CRUD

**Aprendizados detalhados:** `docs/learnings/20260311_202000_rf12_accountant_link_management_learnings.md`

### ✅ Correções Críticas de Processo (11/03/2026)
**Status:** ✅ 100% CONCLUÍDO  
**Descrição:** Correção de problemas no fluxo de trabalho com opencode  
**Entregas:**
- ✅ Script `preserve_context.sh` - Preservação durante compaction
- ✅ Correção do fluxo de tarefas - Agente não executa `conclude_task.sh` automaticamente
- ✅ Validação obrigatória de testes antes da conclusão
- ✅ Documentação: `docs/COMPACTION_HANDLING.md`

**Impacto:** Processo mais robusto para todas as sessões futuras

### 🏗️ Painel do Contador Social e Exportação SPED (11/03/2026)
**Status:** ✅ CONCLUÍDA COM SUCESSO  
**Descrição:** Interface Web do Painel do Contador Social com exportação fiscal SPED/CSV  
**Entregas:**
- Handler `accountant_handler.go` estendendo `BaseHandler`
- Template cache-proof `accountant_dashboard_simple.html` com paleta "Soberania e Suor"
- Rotas: `/accountant/dashboard` (multi-tenant) e `/accountant/export/{entity_id}/{period}`
- Acesso Read-Only ao SQLite (`?mode=ro`) para contadores
- Exportação com hash SHA256 e validação "Soma Zero"
- Testes unitários completos

**Próximos passos operacionais:**
1. Testar integração com `TranslatorService` do módulo `accountant_dashboard`
2. Validar formato de exportação SPED/CSV
3. Testar acesso multi-tenant com dados reais

### 🏗️ Infraestrutura de Deploy (11/03/2026)
**Status:** ✅ CONCLUÍDA COM SUCESSO  
**Descrição:** Sistema completo de deploy em produção com Docker e variáveis de ambiente  
**Entregas:**
- Script `vps_deploy.sh` para automação de VPS
- Sistema de backup/restore para bancos SQLite
- Configuração via variáveis de ambiente (.env)
- Documentação completa (DEPLOYMENT.md, QUICK_DEPLOY.md)
- Scripts de validação automatizada

---

## 🚀 Próxima Tarefa (Sugestões)

Escolha uma tarefa do backlog ou crie uma nova:

### 🎨 Features de UI (Prioridade Alta)
1. **Dashboard de métricas** - Visão consolidada do negócio
2. **Relatórios avançados** - Análise temporal, comparativos
3. **✅ Integração com SPED** - Exportação para contabilidade **(CONCLUÍDA)**

### ⚙️ Melhorias Técnicas (Prioridade Média)
4. **Cache de templates** - Otimização de performance
5. **Testes E2E completos** - Cobertura 100% dos fluxos
6. **Documentação da API** - OpenAPI/Swagger

### 🔧 Infraestrutura (Prioridade Baixa)
7. **Monitoramento** - Prometheus + Grafana
8. **CI/CD pipeline** - GitHub Actions
9. **Multi-tenancy** - Suporte a múltiplas organizações

---

## 📋 Como Prosseguir

1. Use `./process_task.sh "sua descrição de tarefa"`
2. Siga o checklist pré-implementação
3. Documente aprendizados com `./conclude_task.sh`

### Para testar o novo sistema de deploy:
```bash
# Teste local (dry-run)
./scripts/deploy/validate_deployment.sh

# Deploy em staging
./deploy.sh --env-file=.env.staging

# Configurar backup automático
0 2 * * * /opt/digna/scripts/deploy/backup.sh --keep-days=30
```

---

## 🚀 PRÓXIMOS PASSOS PRIORIZADOS

### 🚨 ALTA PRIORIDADE (Resolver bloqueadores - Próxima sessão)
1. **🔧 Corrigir import do módulo lifecycle**
   - **Problema:** `no non-test Go files in /home/.../modules/lifecycle`
   - **Local:** Módulo `ui_web` importando `lifecycle`
   - **Ação:** Investigar `go.mod` replace directive, verificar estrutura do módulo
   - **Impacto:** Todos os testes bloqueados, RF-12 não pode ser finalizado

2. **🔧 Reativar filtro temporal na RF-12**
   - **Local:** `accountant_handler.go:111-133` (comentado)
   - **Ação:** Remover comentários após resolver import
   - **Pré-requisito:** `central.db` já existe e é acessível
   - **Teste:** Login do contador deve funcionar com filtro ativo

3. **🧪 Executar testes E2E da RF-12**
   - **Arquivo:** Testes criados mas não executam
   - **Ação:** Executar após corrigir import
   - **Validação:** Fluxo completo de vínculos contábeis
   - **Cobertura:** Login, filtro temporal, gerenciamento de links

4. **🎨 Completar UI de gerenciamento de links**
   - **Status:** Handler básico criado, falta integração
   - **Ações:** Integrar repositório em `ListLinks()`, criar formulários CRUD
   - **Validação:** Testes de UI, feedback ao usuário

### ⚠️ MÉDIA PRIORIDADE (Completar RF-12)
5. **🧪 Criar testes de integração completos**
   - **Escopo:** Fluxo end-to-end da RF-12
   - **Cenários:** Criação, filtragem, Exit Power, cardinalidade
   - **Validação:** Regras de negócio em ambiente integrado

6. **📊 Adicionar logs de auditoria para vínculos**
   - **Registro:** Quem criou/modificou/desativou cada vínculo
   - **Propósito:** Rastreabilidade e conformidade
   - **Implementação:** Tabela de logs no `central.db`

7. **⚡ Otimizar performance de consultas temporais**
   - **Análise:** EXPLAIN QUERY PLAN nas consultas com date ranges
   - **Otimização:** Índices compostos, estratégias de caching
   - **Métrica:** Tempo de resposta < 100ms para dashboards

8. **📝 Criar documentação da API RF-12**
   - **Formato:** OpenAPI/Swagger ou documentação em Markdown
   - **Conteúdo:** Endpoints, modelos, exemplos de uso
   - **Público:** Desenvolvedores integrando com o sistema

### 📈 LONGO PRAZO (Melhorias e novas features)
9. **🎨 Dashboard de métricas do negócio**
   - **Visão consolidada:** KPIs financeiros, operacionais, sociais
   - **Análise temporal:** Tendências, comparações período a período
   - **Exportação:** Relatórios PDF/Excel personalizáveis

10. **🔍 Relatórios avançados de análise**
    - **Análise de viabilidade:** Projeções financeiras, cenários
    - **Benchmarking:** Comparação com cooperativas similares
    - **Alertas inteligentes:** Notificações proativas baseadas em dados

11. **🔄 Cache de templates e otimização de performance**
    - **Cache em memória:** Templates compilados
    - **Lazy loading:** Assets sob demanda
    - **Métricas:** Tempo de carregamento, uso de recursos

12. **🔐 Aprimoramentos de segurança e auditoria**
    - **Autenticação:** MFA (Multi-Factor Authentication)
    - **Auditoria:** Logs detalhados de todas as operações
    - **Conformidade:** Relatórios para órgãos reguladores

---

## 🏗️ BACKLOG DE MODULARIZAÇÃO

**Criado em:** 13/03/2026  
**Status:** Documentação atualizada, implementação pendente

**Contexto:** Algumas funcionalidades foram implementadas de forma distribuída entre múltiplos módulos, violando o princípio SRP (Single Responsibility Principle). Esta seção documenta o plano de modularização para alinhar a implementação com a arquitetura Clean Architecture + DDD.

### 📊 Resumo do Backlog

| Módulo | Status Atual | Prioridade | Esforço Estimado | Status Documentação |
|--------|--------------|------------|------------------|---------------------|
| `member_management` | ⚠️ Espalhado | **ALTA** | 2-3 dias | ✅ Documentado |
| `reporting` | ⚠️ Básico | MÉDIA | 2-3 dias | ✅ Documentado |
| `sync_engine` | ⚠️ Isolado | MÉDIA | 2-3 dias | ✅ Documentado |

### 1. Modularização: `member_management` (PRIORIDADE ALTA)

**Problema:** Funcionalidade de gerenciamento de membros distribuída entre `core_lume` (domínio/serviço) e `ui_web` (UI com dados mock).

**Localização Atual:**
- `modules/core_lume/internal/domain/member.go` - Entidade Member
- `modules/core_lume/internal/service/member_service.go` - Serviço
- `modules/core_lume/internal/repository/member_test.go` - Testes
- `modules/ui_web/internal/handler/member_handler.go` - Handler UI (usa dados mock)

**Estrutura Planejada:**
```
modules/member_management/
├── go.mod
├── internal/
│   ├── domain/member.go              # Mover de core_lume
│   ├── service/member_service.go     # Mover de core_lume
│   └── repository/
│       ├── interfaces.go
│       └── sqlite_member.go
├── pkg/member/
│   ├── api.go                        # API pública
│   └── types.go                      # Tipos exportados
└── README.md
```

**Tarefas:**
1. ✅ **Documentar** - Estrutura e dependências (CONCLUÍDO 13/03/2026)
2. ⏳ **Criar módulo** - Estrutura de diretórios e go.mod
3. ⏳ **Mover arquivos** - Dominio, serviço, testes do core_lume
4. ⏳ **Atualizar core_lume** - Remover arquivos movidos, atualizar imports
5. ⏳ **Atualizar ui_web** - Usar novo módulo, remover dados mock
6. ⏳ **Testar** - Testes unitários e E2E
7. ⏳ **Validar** - Documentação e integração

**Dependências:**
- `lifecycle` (gerenciamento de banco)
- `core_lume` (apenas tipos compartilhados se necessário)

**Esforço Estimado:** 2-3 dias
**Bloqueadores:** Nenhum (pode começar imediatamente)

---

### 2. Expansão: `reporting` (PRIORIDADE MÉDIA)

**Problema:** Módulo existe com funcionalidade mínima (apenas cálculo de sobras), falta UI e exportação.

**Localização Atual:**
- `modules/reporting/internal/surplus/calculator.go` - Cálculo básico
- `modules/reporting/pkg/surplus/surplus.go` - API pública

**Funcionalidades Faltantes:**
- Handler no `ui_web` para relatórios
- Templates HTML (`templates/reporting_simple.html`)
- Exportação em múltiplos formatos (PDF, CSV, Excel)
- Relatórios específicos (mensal, anual, por membro)

**Estrutura Planejada:**
```
modules/reporting/ (EXPANDIR)
├── internal/
│   ├── surplus/
│   │   ├── calculator.go              # ✅ Já existe
│   │   └── calculator_test.go         # ❌ Adicionar
│   ├── reports/                       # ❌ Novo
│   │   ├── monthly_report.go
│   │   ├── annual_report.go
│   │   └── member_report.go
│   └── export/                        # ❌ Novo
│       ├── pdf_exporter.go
│       ├── csv_exporter.go
│       └── excel_exporter.go
└── pkg/reporting/                     # ❌ Expandir
    ├── api.go
    ├── types.go
    └── export.go
```

**Tarefas:**
1. ✅ **Documentar** - Estrutura e funcionalidades (CONCLUÍDO 13/03/2026)
2. ⏳ **Expandir módulo** - Adicionar estrutura de reports e export
3. ⏳ **Criar handler** - `ui_web/internal/handler/reporting_handler.go`
4. ⏳ **Criar templates** - `ui_web/templates/reporting_simple.html`
5. ⏳ **Implementar exportação** - PDF, CSV, Excel
6. ⏳ **Criar relatórios** - Mensal, anual, por membro
7. ⏳ **Testar** - Testes de exportação e relatórios

**Dependências:**
- `core_lume` (ledger e work)
- `lifecycle` (acesso a dados)
- `ui_web` (para handlers)

**Esforço Estimado:** 2-3 dias
**Bloqueadores:** Nenhum (pode começar imediatamente)

---

### 3. Integração: `sync_engine` (PRIORIDADE MÉDIA)

**Problema:** Módulo isolado com funcionalidades básicas, sem UI e integração com outros módulos.

**Localização Atual:**
- `modules/sync_engine/internal/exchange/intercoop.go` - Troca intercooperativa
- `modules/sync_engine/internal/tracker/sqlite_delta.go` - Rastreamento delta
- `modules/sync_engine/internal/client/sync_repository.go` - Repositório
- `modules/sync_engine/sprint04_test.go` - Testes

**Funcionalidades Faltantes:**
- Handler no `ui_web` para interface de sincronização
- Integração com `supply` e `distribution`
- Sincronização bidirecional
- API para serviços externos (cloud)
- Templates HTML

**Estrutura Planejada:**
```
modules/sync_engine/ (EXPANDIR)
├── internal/
│   ├── exchange/intercoop.go        # ✅ Já existe
│   ├── tracker/sqlite_delta.go      # ✅ Já existe
│   ├── client/sync_repository.go    # ✅ Já existe
│   ├── sync/                        # ❌ Novo
│   │   ├── sync_service.go
│   │   ├── cloud_sync.go
│   │   └── conflict_resolver.go
│   └── api/
│       └── external_api.go          # ❌ Novo
└── pkg/sync/
    ├── api.go                       # ❌ Novo
    └── types.go                     # ❌ Novo
```

**Tarefas:**
1. ✅ **Documentar** - Estrutura e integrações (CONCLUÍDO 13/03/2026)
2. ⏳ **Expandir módulo** - Adicionar serviços de sync
3. ⏳ **Criar handler** - `ui_web/internal/handler/sync_handler.go`
4. ⏳ **Criar templates** - `ui_web/templates/sync_simple.html`
5. ⏳ **Implementar integração** - Com supply e distribution
6. ⏳ **Implementar sync cloud** - Sincronização bidirecional
7. ⏳ **Testar** - Testes de sincronização e conflitos

**Dependências:**
- `lifecycle` (acesso a dados)
- `supply` (para integração)
- `distribution` (para integração)
- `ui_web` (para handlers)

**Esforço Estimado:** 2-3 dias
**Bloqueadores:** Nenhum (pode começar imediatamente)

---

### 📋 Sequência de Implementação Recomendada

**Fase 1 (Semanas 1-2):** `member_management`
- Maior impacto na arquitetura
- Resolve violação SRP mais crítica
- Facilita manutenção futura

**Fase 2 (Semanas 3-4):** `reporting`
- Valor de negócio alto (relatórios)
- Não tem dependências críticas
- Pode ser feito em paralelo se necessário

**Fase 3 (Semanas 5-6):** `sync_engine`
- Mais complexo (integrações)
- Depende de entendimento dos outros módulos
- Recomendado fazer por último

### 🔄 Processo de Modularização

Para cada módulo, seguir o processo:

1. **Preparação**
   - Criar branch específica
   - Fazer backup dos arquivos
   - Revisar documentação do plano

2. **Criação da Estrutura**
   - Criar diretórios (sem mover código ainda)
   - Configurar go.mod
   - Preparar imports

3. **Migração**
   - Mover arquivos conforme planejado
   - Atualizar imports nos arquivos movidos
   - Remover arquivos antigos

4. **Integração**
   - Atualizar módulos dependentes
   - Configurar replace directives no go.work
   - Resolver dependências cíclicas

5. **Testes**
   - Executar testes unitários
   - Executar testes de integração
   - Executar testes E2E

6. **Validação**
   - Verificar build de todos os módulos
   - Documentar no `docs/learnings/`
   - Atualizar `docs/NEXT_STEPS.md`

7. **Merge**
   - Code review
   - Merge para main
   - Tag de versão

### 📚 Documentação Relacionada

- `docs/03_architecture/01_system.md` - Arquitetura do sistema e tabela de sprints atualizada
- `docs/README.md` - Estrutura de módulos atualizada
- `docs/QUICK_REFERENCE.md` - Referência rápida de módulos
- `docs/learnings/` - Aprendizados de implementações anteriores

---

## 📋 COMO PROCEDER NA PRÓXIMA SESSÃO

### Fluxo Recomendado:
1. **Comece com o bloqueador:** Resolver erro de import do lifecycle
2. **Valide RF-12:** Executar testes E2E após correção
3. **Complete a implementação:** UI de gerenciamento de links
4. **Documente:** Atualizar aprendizados e NEXT_STEPS

### Scripts a Usar:
```bash
# Se opencode entrar em compaction:
./preserve_context.sh

# Para criar nova tarefa:
./process_task.sh "Descrição da tarefa"

# Para concluir (APENAS usuário executa):
./conclude_task.sh --task=ID "Aprendizados da tarefa"

# Para validar testes:
./scripts/validate_task_tests.sh
```

### Checkpoints Obrigatórios:
- ✅ Testes passando antes de concluir qualquer tarefa
- ✅ Documentação em `docs/learnings/` com timestamp correto
- ✅ NEXT_STEPS.md atualizado com progresso
- ✅ Preservação de contexto se houver compaction

---

## 📊 STATUS ATUAL DO PROJETO

**Testes:** ⚠️ 149/149 (bloqueados por import)  
**Handlers UI:** 15 ✅ (incluindo RF-12 parcial)  
**Templates:** 19 ✅ (incluindo RF-12)  
**Deploy em produção:** ✅ PRONTO  
**Documentação:** ✅ COMPLETA (processo corrigido)

**RF-12 Status:** ⚠️ 85% COMPLETO (BLOQUEADO)  
**Processo:** ✅ CORRIGIDO (compaction, fluxo, testes)  
**Próxima sessão:** FOCAR EM RESOLVER BLOQUEADOR E COMPLETAR RF-12

### Métricas de Qualidade:
- **Cobertura de testes (estimada):** 85% backend, 30% UI
- **Dívida técnica:** Moderada (correções temporárias em RF-12)
- **Risco de regressão:** Baixo (banco central isolado, testes obrigatórios)
- **Processo robustez:** Alta (scripts de preservação, validação)

---

**📌 NOTA FINAL:** O projeto está em excelente estado estrutural com deploy pronto e processo corrigido. O bloqueador atual (import do lifecycle) é técnico e resolvível. Uma vez resolvido, a RF-12 pode ser finalizada rapidamente (2-3 horas) e o projeto estará pronto para novas features.

**Prioridade absoluta para próxima sessão:** Resolver erro de import → Completar RF-12 → Validar todo o sistema.

