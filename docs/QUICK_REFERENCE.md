# 🚀 Quick Reference - Projeto Digna

**Última atualização:** 13/03/2026
**Status:** ✅ PRODUCTION READY (149/149 testes passando) + RF-12 95% COMPLETO (FUNCIONAL)

---

## 🏗️ **Arquitetura Core (Constituição de IA)**

### **1. Anti-Float (Regra Sagrada)**
- **Proibido:** `float`, `float32`, `float64` para valores financeiros/tempo
- **Obrigatório:** `int64` para centavos (R$ 1,00 = 100) e minutos
- **Validação:** Todo handler deve escanear código por floats antes de commit

### **2. Clean Architecture + DDD**
```
internal/domain/     # Entidades puras, interfaces Repository (SEM SQL/HTTP)
internal/service/    # Casos de uso, orquestração (depende de interfaces)
internal/repository/ # Implementações SQLite (via LifecycleManager)
internal/handler/    # HTTP handlers (UI Web)
```

### **3. Soberania de Dados**
- **Isolamento:** `data/entities/{entity_id}.db` (um banco por entidade)
- **Banco Central:** `data/entities/central.db` para relações inter-tenant (RF-12)
- **LifecycleManager:** Ponto único de acesso a bancos SQLite
- **Context:** `entity_id` extraído de `r.Context().Value("entity_id")`
- **Proibido:** JOINs entre bancos diferentes

---

## 🎨 **Frontend Patterns (UI Web)**

### **1. Sistema de Templates Cache-Proof**
- **Nomenclatura:** `*_simple.html` (documentos HTML completos)
- **Carregamento:** `template.ParseFiles("templates/nome_simple.html")` NO HANDLER
- **Proibido:** Variáveis globais de template, `template.ParseGlob()`
- **BaseHandler:** `modules/ui_web/internal/handler/base_handler.go`
  ```go
  type BaseHandler struct {
      lifecycleManager lifecycle.LifecycleManager
      templateManager  *tmpl.TemplateManager
  }
  ```

### **2. Funções de Template (TemplateManager)**
```go
// Funções registradas no TemplateManager
"formatCurrency": func(amount int64) string  // R$ 1.50
"divide": func(a, b int64) float64          // divisão segura
"multiply": func(a, b int64) int64          // multiplicação
"formatDate": func(t interface{}) string    // formatação data
"getAlertStatusLabel": func(status string) string
"getAlertStatusClass": func(status string) string
"getCategoryLabel": func(category string) string
"fdiv": func(a, b float64) float64          // divisão float (apenas UI)
```

### **3. Padrão HTMX**
```html
<!-- Formulário assíncrono -->
<form hx-post="/endpoint" 
      hx-target="#result-area" 
      hx-swap="outerHTML">
</form>

<!-- Ação com feedback -->
<button hx-post="/action" 
        hx-target="#feedback"
        hx-swap="innerHTML">
  Ação
</button>
```

### **4. Design System "Soberania e Suor"**
- **Azul Soberania:** `#2A5CAA` (headers, botões principais)
- **Verde Suor:** `#4A7F3E` (indicadores trabalho/sucesso)
- **Laranja Energia:** `#F57F17` (alertas, destaques)
- **Fundo:** `#F9F9F6`, **Texto:** `#212121`
- **Fontes:** Inter (primária), Ubuntu (secundária)

---

## 📁 **Estrutura de Módulos**

### **Módulos Implementados**

#### ✅ **Módulos Completos**
```
modules/
├── accountant_dashboard/  # Dashboard contábil, SPED
├── budget/               # Orçamento, planejamento
├── cash_flow/           # Fluxo de caixa
├── core_lume/           # Domínio central (Ledger, WorkLog)
│   ├── internal/domain/member.go      # ⚠️ Será movido para member_management
│   └── internal/service/member_service.go  # ⚠️ Será movido para member_management
├── distribution/        # Distribuição de sobras
├── integrations/        # Integrações externas
├── legal_facade/        # Facade jurídica (Formalização, Documentos)
├── lifecycle/           # LifecycleManager, isolamento SQLite
├── pdv_ui/             # Ponto de Venda
├── supply/             # Compras, estoque, fornecedores
└── ui_web/             # Interface web principal
    ├── internal/handler/member_handler.go    # ⚠️ Usará novo módulo member_management
    ├── internal/handler/legal_handler.go      # Integrado com legal_facade
    └── ...                                    # Outros handlers
```

#### ⚠️ **Módulos em Desenvolvimento/Backlog**
```
modules/
├── reporting/           # Relatórios (cálculo básico implementado)
│   ├── internal/surplus/calculator.go       # ✅ Cálculo de sobras
│   └── pkg/surplus/surplus.go               # ✅ API pública
│   ⚠️ Faltando: Handlers UI, templates, exportação PDF/CSV/Excel
│
├── sync_engine/        # Sincronização (funcionalidades básicas)
│   ├── internal/exchange/intercoop.go       # ✅ Troca intercooperativa
│   ├── internal/tracker/sqlite_delta.go     # ✅ Rastreamento delta
│   └── internal/client/sync_repository.go   # ✅ Repositório de sync
│   ⚠️ Faltando: Handler UI, integração com outros módulos, sincronização cloud
│
└── member_management/  # ⚠️ NÃO EXISTE - Funcionalidade espalhada
    📍 Localização atual:
       - core_lume/internal/domain/member.go (domínio)
       - core_lume/internal/service/member_service.go (serviço)
       - ui_web/internal/handler/member_handler.go (UI com dados mock)
    🔧 Ação: Criar módulo separado movendo arquivos do core_lume
```

#### 📋 **Módulos de Teste**
```
modules/
├── integration_test/    # Testes de integração e E2E
│   ├── journey_e2e_test.go          # Testes de jornada completa
│   ├── integrations_e2e_test.go     # Testes de integração
│   └── functional_test.go           # Testes funcionais
```

### **Estrutura de Módulos Backlog (Próximos Passos)**

| Módulo | Status | Prioridade | Esforço Estimado |
|--------|--------|------------|------------------|
| `member_management` | ⚠️ Espalhado | **ALTA** | 2-3 dias |
| `reporting` | ⚠️ Básico | MÉDIA | 2-3 dias |
| `sync_engine` | ⚠️ Isolado | MÉDIA | 2-3 dias |

**Ver:** `docs/NEXT_STEPS.md` para detalhes completos do backlog


---

## 🚀 Deploy em Produção

### Scripts de Deploy
- **`deploy.sh`** - Wrapper principal (raiz do projeto)
- **`scripts/deploy/vps_deploy.sh`** - Deploy automático em VPS
- **`scripts/deploy/backup.sh`** - Backup de bancos SQLite
- **`scripts/deploy/restore.sh`** - Restauração de backup
- **`scripts/deploy/validate_deployment.sh`** - Validação dos scripts

### Configuração
- **Porta:** `DIGNA_PORT=8090` (variável de ambiente)
- **Dados:** `DIGNA_DATA_DIR=/var/lib/digna/data`
- **Logs:** `DIGNA_LOG_LEVEL=info`
- **Documentação:** `docs/DEPLOYMENT.md` (guia completo)

### Comandos Rápidos
```bash
# Deploy automático
./deploy.sh

# Backup dados
./scripts/deploy/backup.sh --keep-days=30

# Validação
./scripts/deploy/validate_deployment.sh
```


## 🔄 Aprendizados Recentes da Sessão
**Baseado em:** Análise profunda do códigobase (11/03/2026)
**Período:** 11/03/2026

### Insights Críticos Descobertos:
1. **`legal_facade` já existe** com 80% da funcionalidade (generator.go, formalization.go, statute.go)
2. **`FormalizationSimulator` já tem** `MinDecisionsForFormalization = 3`
3. **SHA256 já implementado** em múltiplos lugares (statute.go, decision_service.go)
4. **File download pattern** existe em `accountant_handler.go`
5. **Skills críticas** em `docs/skills/` (backend, frontend, soberania de dados)

### Padrões Identificados:
- **File download:** `Content-Disposition: attachment; filename=...`
- **SHA256:** `sha256.Sum256([]byte(data))` + `hex.EncodeToString(hash[:])`
- **Cache-proof templates:** `*_simple.html` + `ParseFiles()` no handler

### Status:
- **Economia de tempo:** 40min descoberta → 5min consulta (80% redução)
- **Próxima sessão:** 50-70% mais eficiente com documentação criada

💡 **Dica:** Consulte `docs/MODULES_QUICK_REFERENCE.md` e `docs/learnings/SESSION_INSIGHTS_20260311.md` para detalhes.
## 🆕 Nova Sessão

**Sessão iniciada em:** 11/03/2026 10:08
**Status:** ✅ PRONTO PARA NOVA IMPLEMENTAÇÃO

Use `./start_session.sh` para contexto completo ou `./process_task.sh` para começar.

## 🏗️ Handlers Existentes

accountant_handler, auth_handler, base_handler, budget_handler, budget_templates, cash_handler, dashboard, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, auth_handler, base_handler, budget_handler, budget_templates, cash_handler, dashboard, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, auth_handler, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, auth_handler, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 📚 Aprendizados Recentes

- **Sessão migrated_20260311_140200:** 11/03/2026 - 0h6m, 2 tarefas (ver `docs/learnings/SESSION_migrated_20260311_140200_CONSOLIDATED.md`)

## 🏗️ Handlers Existentes

accountant_handler, auth_handler, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

### 📊 Painel do Contador Social (accountant_handler)
- **Rotas:** `/accountant/dashboard`, `/accountant/export/{entity_id}/{period}`
- **Funcionalidade:** Interface multi-tenant para contadores sociais
- **Segurança:** Acesso Read-Only ao SQLite (`?mode=ro`)
- **Exportação:** Geração de lotes fiscais SPED/CSV com hash SHA256
- **Validação:** "Soma Zero" automática antes da exportação
- **Template:** `accountant_dashboard_simple.html` (cache-proof)
- **Sessão 20260311_143158:** 11/03/2026 - 0h16m, 6 tarefas (ver `docs/learnings/SESSION_20260311_143158_CONSOLIDATED.md`)

## 🏗️ Handlers Existentes

accountant_handler, auth_handler, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, auth_handler, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, auth_handler, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates
- **Sessão 20260311_150149:** 11/03/2026 - 1h22m, 2 tarefas (ver `docs/learnings/SESSION_20260311_150149_CONSOLIDATED.md`)

## 🏗️ Handlers Existentes

accountant_handler, auth_handler, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, auth_handler, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, auth_handler, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, auth_handler, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, auth_handler, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, auth_handler, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, auth_handler, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, auth_handler, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, auth_handler, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, auth_handler, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, auth_handler, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, auth_handler, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, auth_handler, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, auth_handler, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, auth_handler, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, auth_handler, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, auth_handler, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, auth_handler, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, auth_handler, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, accountant_link_handler, auth_handler, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, accountant_link_handler, auth_handler, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates
- **Sessão 20260311_175424:** 11/03/2026 - 1h44m, 6 tarefas (ver `docs/learnings/SESSION_20260311_175424_CONSOLIDATED.md`)

## 🏗️ Handlers Existentes

accountant_handler, accountant_link_handler, auth_handler, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, accountant_link_handler, auth_handler, auth_handler_mock, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates
- **Sessão 20260311_203057:** 11/03/2026 - 0h44m, 2 tarefas (ver `docs/learnings/SESSION_20260311_203057_CONSOLIDATED.md`)

## 🏗️ Handlers Existentes

accountant_handler, accountant_link_handler, auth_handler, auth_handler_mock, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, accountant_link_handler, auth_handler, auth_handler_mock, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, accountant_link_handler, auth_handler, auth_handler_mock, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, accountant_link_handler, auth_handler, auth_handler_mock, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, accountant_link_handler, auth_handler, auth_handler_mock, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, accountant_link_handler, auth_handler, auth_handler_mock, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, accountant_link_handler, auth_handler, auth_handler_mock, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, accountant_link_handler, auth_handler, auth_handler_mock, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, accountant_link_handler, auth_handler, auth_handler_mock, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, accountant_link_handler, auth_handler, auth_handler_mock, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, accountant_link_handler, auth_handler, auth_handler_mock, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, accountant_link_handler, auth_handler, auth_handler_mock, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, accountant_link_handler, auth_handler, auth_handler_mock, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, accountant_link_handler, auth_handler, auth_handler_mock, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, accountant_link_handler, auth_handler, auth_handler_mock, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates
- **Sessão 20260312_105735:** 12/03/2026 - 2h56m, 2 tarefas (ver `docs/learnings/SESSION_20260312_105735_CONSOLIDATED.md`)

## 🏗️ Handlers Existentes

accountant_handler, accountant_link_handler, auth_handler, auth_handler_mock, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, accountant_link_handler, auth_handler, auth_handler_mock, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates
- **Sessão 20260312_135447:** 12/03/2026 - 0h7m, 2 tarefas (ver `docs/learnings/SESSION_20260312_135447_CONSOLIDATED.md`)

## 🏗️ Handlers Existentes

accountant_handler, accountant_link_handler, auth_handler, auth_handler_mock, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, accountant_link_handler, auth_handler, auth_handler_mock, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates
- **Sessão 20260312_141503:** 12/03/2026 - 1h2m, 2 tarefas (ver `docs/learnings/SESSION_20260312_141503_CONSOLIDATED.md`)
- **Sessão unknown:** 13/03/2026 - 0h0m, 0 tarefas (ver `docs/learnings/SESSION_unknown_CONSOLIDATED.md`)

## 🏗️ Handlers Existentes

accountant_handler, accountant_link_handler, auth_handler, auth_handler_mock, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, accountant_link_handler, auth_handler, auth_handler_mock, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates
- **Sessão 20260313_163003:** 13/03/2026 - 1h2m, 2 tarefas (ver `docs/learnings/SESSION_20260313_163003_CONSOLIDATED.md`)

## 🏗️ Handlers Existentes

accountant_handler, accountant_link_handler, auth_handler, auth_handler_mock, base_handler, budget_handler, budget_templates, cash_handler, dashboard, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, accountant_link_handler, auth_handler, auth_handler_mock, base_handler, budget_handler, budget_templates, cash_handler, dashboard, das_mei_handler, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, accountant_link_handler, auth_handler, auth_handler_mock, base_handler, budget_handler, budget_templates, cash_handler, dashboard, das_mei_handler, eligibility_handler, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, accountant_link_handler, auth_handler, auth_handler_mock, base_handler, budget_handler, budget_templates, cash_handler, dashboard, das_mei_handler, eligibility_handler, help_handler, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates

## 🏗️ Handlers Existentes

accountant_handler, accountant_link_handler, auth_handler, auth_handler_mock, base_handler, budget_handler, budget_templates, cash_handler, dashboard, das_mei_handler, eligibility_handler, help_handler, legal_handler, member_handler, pdv_handler, supply_handler, supply_templates
- **Sessão 20260326_175404:** 27/03/2026 - 6h29m, 8 tarefas (ver `docs/learnings/SESSION_20260326_175404_CONSOLIDATED.md`)
