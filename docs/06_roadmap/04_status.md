***

```markdown
---
title: Status Atual
status: implemented
version: 1.4
last_updated: 2026-03-09
---

# Status Atual - Digna

**Última Atualização:** 2026-03-09
**Fase Atual:** Sprint 15 (Correções Críticas e Testes E2E) ✅ COMPLETE
**Próximo Marco:** Production Deploy (Marco 05)

---

## Phase Status Overview

| Phase | Marco | Status | Conclusão |
|-------|-------|--------|-----------|
| Concepção | Marco 00 | ✅ COMPLETE | 100% |
| Foundation Setup | Marco 01 | ✅ COMPLETE | 100% |
| Core Operations | Marco 02 | ✅ COMPLETE | 100% |
| Reporting & Documents | Marco 03 | ✅ COMPLETE | 100% |
| UI & Dashboard | Marco 04 | ✅ COMPLETE | 100% |
| Integração e Aliança Contábil (Phase 2) | Marco 07 | ✅ COMPLETE | 100% |
| Gestão de Compras e Estoque (Phase 3) | Marco 08 | ✅ COMPLETE | 100% |
| Gestão Orçamentária e Planejamento Financeiro (Phase 3) | Marco 06 | ✅ COMPLETE | 100% |
| Production Deploy | Marco 05 | 📋 PLANNED | 0% |

---

## Sprint Status

### Sprint 01 a 03: Core, Ledger e Reporting ✅
- Lifecycle Manager (SQLite isolado)
- Ledger Service (partidas dobradas exatas em `int64`)
- Surplus Calculator (rateio social ITG 2002)
- **Testes:** 22/22 PASS

### Sprint 04: Sincronização & Intercooperação ✅
- Delta Tracker
- Sync Package
- Marketplace B2B
- **Testes:** 9/9 PASS

### Sprint 05: Interface Humana Dignidade ✅
- Servidor HTTP porta 8088
- PDV Screen (HTMX)
- Social Clock e Dashboard
- PWA (manifest + service worker)
- **Testes:** 9/9 PASS

### Sprint 06: Gestão de Caixa (RF-09) ✅
- Módulo cash_flow criado
- Registro de entradas e saídas
- Saldo em tempo real e Extrato por período
- Interface web /cash
- **Testes:** 3/3 PASS

### Sprint 07: DDD Refactoring ✅
- Centralizado validação de transações (EntryValidator)
- Removido erros ignorados (result.LastInsertId)
- Adicionado rows.Err() checks em todas as queries
- Implementado graceful shutdown no servidor HTTP
- **Testes:** 8/8 PASS (novos) + regressão 35 PASS

### Sprint 08 e 09: Integrações e Testes de Base ✅
- 8 Interfaces Governamentais (Mock) implementadas.
- Cobertura expandida para testagem de fluxos internos.
- **Testes:** 13/13 PASS

### Sprint 10: Gestão de Membros ✅
- Entity Member com roles (COORDINATOR, MEMBER, ADVISOR)
- MemberRepository com UPSERT, FindByID, ListByEntity
- MemberService com Register, Update, Deactivate
- Validação: não permite desativar último coordenador
- **Testes:** 19/19 PASS

### Sprint 11: Formalização e E2E Journey ✅
- **SurplusCalculator:** Novo método CalculateWithDeductions() (15% bloqueados para FATES e Reserva Legal).
- **FormalizationSimulator:** Novo método AutoTransitionIfReady() (DREAM -> FORMALIZED após 3 decisões).
- **E2E:** `journey_e2e_test.go` finalizado simulando a jornada "Sonho Solidário".
- **Testes:** 5/5 PASS

### Sprint 12: Painel do Contador Social (Accountant Dashboard) ✅ COMPLETE
- **Objetivo:** Interface Multi-tenant para profissionais contábeis parceiros.
- **Isolamento:** Acesso estritamente *Read-Only* aos micro-databases `.sqlite` das entidades autorizadas (`?mode=ro`).
- **Exportação:** Motor de Tradução Fiscal (Geração de Lotes SPED a partir das partidas dobradas).
- **Anti-Float:** Todos os valores monetários usam `int64`, sem `float`.
- **Decisões Arquiteturais:**
  - ✅ **Integração via `ui_web`:** Em vez de `cmd/digna/main.go`, seguindo consistência arquitetural
  - ✅ **Templates Embutidos:** Em vez de arquivos `.html` separados, simplificando deploy
  - ✅ **Princípios Aplicados:** KISS, YAGNI, DRY, Consistência
- **Implementado:**
  - [x] Domain Layer (FiscalBatch, EntryDTO, AccountMapper) - 100% coverage
  - [x] Repository Layer (SQLite Read-Only Adapter) - 87.2% coverage
  - [x] Service Layer (Translator Service com Soma Zero validation) - 91.3% coverage
  - [x] Handler Layer (Dashboard + Export com HTMX/Tailwind) - 97.1% coverage
  - [x] Integration with ui_web module (accountant_handler.go)
  - [x] Public API for external consumption - 26.7% coverage
  - [x] Integration tests covering complete workflow
  - [x] **E2E Journey Test Updated:** Jornada "Sonho Solidário" atualizada com auditorias do Contador Social
- **Testes:** Todos os testes PASS com cobertura total de 69.0% (core packages: 93.9% average) ✅
- **E2E Validation:** Teste de jornada anual atualizado e validado com sucesso ✅

### Sprint 13: Gestão de Compras e Controle de Estoque (RF-07 e RF-08) ✅ COMPLETE
- **Objetivo:** Módulo completo para registro de compras, gestão de fornecedores e controle de estoque com **contabilidade invisível**.
- **Paradigma:** Usuário final NÃO FAZ CONTABILIDADE - apenas informa "Comprei X de Y por Z reais"
- **Categorização:** Tipos de itens (INSUMO, PRODUTO, MERCADORIA) para:
  - Interface PDV: mostrar apenas "Produto Acabado" na venda
  - Contabilidade: comprar "Insumo" → despesa/estoque; saída "Produto" → receita
- **Implementado:**
  - [x] **Módulo `supply`** com Clean Architecture + DDD
  - [x] **Domínio:** Supplier, StockItem, Purchase, PurchaseItem
  - [x] **Repository:** SQLiteSupplyRepository com DDL completo
  - [x] **Service:** PurchaseService com integração core_lume
  - [x] **Contabilidade Invisível:** Partidas dobradas automáticas baseadas no tipo do item
  - [x] **UI Web:** Handler em `ui_web` com templates embutidos
  - [x] **Rotas:** `/supply`, `/supply/purchase`, `/supply/suppliers`, `/supply/stock`
  - [x] **API:** Endpoints REST para todas as operações
  - [x] **Testes Unitários:** Cobertura completa do módulo supply
  - [x] **Anti-Float:** Validação completa - nenhum uso de `float`
- **Integração Contábil:**
  - INSUMO/MERCADORIA: Débito em `AccountInventory` (3)
  - Pagamento à vista: Crédito em `AccountCash` (1)
  - Pagamento a prazo: Crédito em `AccountSuppliers` (4)
- **Testes:** Todos os testes PASS, compilação completa do projeto ✅

### Sprint 14: Gestão Orçamentária e Planejamento Financeiro (RF-10) ✅ COMPLETE
- **Objetivo:** Módulo completo para planejamento financeiro com acompanhamento "planejado vs realizado" e interface pedagógica.
- **Paradigma:** Usuário final NÃO USA termos técnicos como "Budget", "Forecast", "CAPEX". Usa linguagem popular: "Planejamento do Mês", "O que combinamos de gastar", "O que já gastamos", "Aviso de limite".
- **Contabilidade Invisível:** Sistema cruza automaticamente planejado com realizado (transações reais do Ledger/Caixa).
- **Alertas Visuais:** Status SAFE (≤70%), WARNING (71-100%), EXCEEDED (>100%) com barras de progresso coloridas.
- **Anti-Float Obrigatório:** Todos os valores em `int64` (centavos). Proibido uso de `float`.
- **Implementado:**
  - [x] **Módulo `budget`** com Clean Architecture + DDD
  - [x] **Domínio:** BudgetPlan, BudgetExecution, BudgetAlertStatus, BudgetCategory
  - [x] **Repository:** SQLiteBudgetRepository com DDL completo
  - [x] **Service:** BudgetService com integração cash_flow (CashFlowPort)
  - [x] **Contabilidade Invisível:** Cálculo automático planejado vs realizado
  - [x] **UI Web:** Handler em `ui_web` com templates embutidos
  - [x] **Rotas:** `/budget`, `/budget/report`, `/budget/create`
  - [x] **API:** BudgetAPI com CashFlowAdapter para integração real
  - [x] **Testes Unitários:** Cobertura completa do módulo budget
  - [x] **Anti-Float:** Validação completa - nenhum uso de `float`
- **Categorias Pré-definidas:** INSUMOS, ENERGIA, EQUIPAMENTOS, TRANSPORTE, MANUTENCAO, SERVICOS, OUTROS
- **Integração:** Conectado automaticamente ao cash_flow para buscar transações reais
- **Testes:** Todos os testes PASS, compilação completa do projeto ✅

### Sprint 15: Correções Críticas e Testes E2E (PDV → Estoque → Caixa) ✅ COMPLETE
- **Objetivo:** Corrigir três problemas críticos reportados e implementar testes de integração E2E completos com Playwright.
- **Problemas Resolvidos:**
  1. **Vendas registradas no PDV não aparecem na tela do caixa** ✅
  2. **Sistema permite vender mais itens do que existem em estoque** ✅  
  3. **Sistema não atualiza o estoque após vendas** ✅
- **Implementado:**
  - [x] **Correção PDV → Caixa:** Implementado `getEntriesFromDatabase` no cash handler para buscar transações diretamente do banco
  - [x] **Validação de Estoque:** Adicionada validação no PDV handler que verifica `quantidade ≤ estoque disponível`
  - [x] **Atualização de Estoque:** Implementado `UpdateStockQuantity` na API do supply e integração no PDV
  - [x] **Integração Frontend:** Corrigido JavaScript no template PDV para passar `stock_item_id` corretamente
  - [x] **Testes E2E com Playwright:** Configurado ambiente completo com servidor real + browser headless
  - [x] **Testes de Fluxo Completo:** Criado testes que simulam usuário interagindo com a aplicação
  - [x] **Validação de Integridade:** Testes verificam fluxo PDV → Estoque → Caixa
  - [x] **Interface Gestão de Estoque:** Dashboard atualizado com link para `/supply`
  - [x] **Testes Otimizados:** Criados testes mais rápidos sem browser (API-only)
  - [x] **Testes de Validação:** Testes específicos para validação de estoque insuficiente
- **Arquivos Modificados/Criados:**
  - `modules/ui_web/internal/handler/cash_handler.go` - Adicionada busca de transações do banco
  - `modules/ui_web/internal/handler/pdv_handler.go` - Adicionada validação e atualização de estoque
  - `modules/ui_web/templates/pdv.html` - Corrigido JavaScript para passar stock_item_id
  - `modules/ui_web/templates/dashboard.html` - Adicionado link para gestão de estoque
  - `modules/supply/pkg/supply/api.go` - Implementado UpdateStockQuantity
  - `modules/supply/pkg/supply/interfaces.go` - Adicionado UpdateStockQuantity à interface
  - `modules/ui_web/e2e_pdv_estoque_caixa_test.go` - Teste E2E completo com Playwright
  - `modules/ui_web/e2e_simplificado_test.go` - Teste E2E simplificado
  - `modules/ui_web/e2e_otimizado_test.go` - Testes E2E otimizados (mais rápidos)
  - `modules/ui_web/test_fluxo_completo_test.go` - Teste de fluxo completo PDV→Estoque→Caixa
  - `modules/ui_web/test_validacao_estoque_test.go` - Teste específico de validação de estoque
- **Testes:** Todos os testes PASS, validação completa da integração ✅

---

## Total Test Coverage

| Sprint | Testes | Status |
|--------|--------|--------|
| 01 | 6/6 | ✅ PASS |
| 02 | 8/8 | ✅ PASS |
| 03 | 8/8 | ✅ PASS |
| 04 | 9/9 | ✅ PASS |
| 05 | 9/9 | ✅ PASS |
| 06 | 3/3 | ✅ PASS |
| 07 | 43/43 | ✅ PASS |
| 08 | 5/5 | ✅ PASS |
| 09 | 8/8 | ✅ PASS |
| 10 | 19/19 | ✅ PASS |
| 11 | 5/5 | ✅ PASS |
| 12 | 8/8 | ✅ PASS |
| 13 | 6/6 | ✅ PASS |
| 14 | 4/4 | ✅ PASS |
| 15 | 3/3 | ✅ PASS |
| **Total** | **149/149** | **100% PASS** 🎉 |

---

## DDD Architecture Status

| Módulo | Interface Repository | Implementação | Status |
|--------|---------------------|-----------------|--------|
| core_lume | LedgerRepository, WorkRepository, DecisionRepository | SQLite | ✅ COMPLETE |
| reporting | SurplusRepository | Adapter Pattern | ✅ COMPLETE |
| sync_engine | SyncRepository | SQLite | ✅ COMPLETE |
| legal_facade | LegalRepository | SQLite | ✅ COMPLETE |
| integrations | 8 interfaces governamentais | Mock | ✅ COMPLETE |
| accountant_dashboard| FiscalRepository | Read-Only SQLite Adapter | ✅ COMPLETE |
| supply | SupplyRepository | SQLite | ✅ COMPLETE (com UpdateStockQuantity) |
| budget | BudgetRepository | SQLite | ✅ COMPLETE |
| ui_web | CashHandler, PDVHandler | HTTP Handlers | ✅ COMPLETE (com integração PDV→Estoque→Caixa) |

---

## Próximos Passos

1. **Production Deploy (Marco 05):** Preparar para deploy em produção
2. **Integração Real:** Iniciar substituição da autenticação simulada pelo OAuth2 real do Gov.br.
3. **Testes de Usabilidade:** Levar o PWA e o Motor Lume para campo com cooperativas reais e Incubadoras (ITCPs).
4. **Documentação Técnica:** Gerar API Docs / Swagger para permitir intercooperação com BCDs (Bancos Comunitários).
5. **Expansão Testes E2E:** Adicionar mais cenários de teste e interface para gestão de estoque.
```

***
