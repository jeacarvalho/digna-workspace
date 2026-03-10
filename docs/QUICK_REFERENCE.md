# 🚀 Quick Reference - Projeto Digna

**Última atualização:** 10/03/2026
**Status:** ✅ PRODUCTION READY (149/149 testes passando)

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

### **Módulos Principais (100% testados)**
```
modules/
├── accountant_dashboard/  # Dashboard contábil, SPED
├── budget/               # Orçamento, planejamento
├── cash_flow/           # Fluxo de caixa
├── core_lume/           # Domínio central (Membros, Ledger)
├── distribution/        # Distribuição de sobras
├── integrations/        # Integrações externas
├── lifecycle/           # LifecycleManager, isolamento SQLite
├── pdv_ui/             # Ponto de Venda
├── supply/             # Compras, estoque, fornecedores
├── sync_engine/        # Sincronização delta
└── ui_web/             # Interface web principal
```

---

## 🆕 Nova Sessão

**Sessão iniciada em:** 10/03/2026 20:20
**Status:** ✅ PRONTO PARA NOVA IMPLEMENTAÇÃO

Use `./start_session.sh` para contexto completo ou `./process_task.sh` para começar.

## 🏗️ Handlers Existentes

accountant_handler, auth_handler, base_handler, budget_handler, budget_templates, cash_handler, dashboard, member_handler, pdv_handler, supply_handler, supply_templates
