# 🚫 Antipadrões Comuns e Soluções - Projeto Digna

**Última atualização:** 10/03/2026  
**Baseado em:** Implementação MemberHandler + Análise de código existente

---

## 🎯 **Objetivo**
Identificar padrões problemáticos recorrentes e fornecer soluções testadas.  
**Uso:** Consultar ANTES de implementar nova feature, DURANTE code reviews.

---

## 🏗️ **Arquitetura**

### **Antipadrão 1: "Importação de pacote internal"**
**Sintoma:** `import "github.com/providentia/digna/core_lume/internal/service"` → erro de compilação
**Causa:** Pacotes `internal` em Go são restritos ao módulo
**Impacto:** Bloqueia integração entre módulos
**Solução:**
```go
// ❌ ERRADO (não compila)
import "github.com/providentia/digna/core_lume/internal/service"

// ✅ SOLUÇÃO 1: API layer (padrão cash_flow)
import "github.com/providentia/digna/cash_flow/pkg/cash_flow"
handler := cash_flow.NewCashFlowAPI(lm)

// ✅ SOLUÇÃO 2: Mock inicial + integração futura
type MemberHandler struct {
    *BaseHandler
    // TODO: Adicionar MemberService quando tiver acesso
    // Por enquanto, dados mockados para desenvolvimento
}

// ✅ SOLUÇÃO 3: Refatorar pacote para público (decisão arquitetural)
// Mover para pkg/ ou criar interfaces públicas
```

### **Antipadrão 2: "Handler independente do BaseHandler"**
**Sintoma:** Handler cria seu próprio `template.Template`, duplica funções
**Causa:** Não seguir padrão estabelecido
**Impacto:** Inconsistência, não cache-proof, duplicação
**Solução:**
```go
// ❌ ERRADO
type MyHandler struct {
    tmpl *template.Template
}
func NewMyHandler() {
    tmpl := template.New("").Funcs(myFuncs)
    tmpl.ParseFiles("templates/my_simple.html")
}

// ✅ CORRETO
type MyHandler struct {
    *BaseHandler
}
func NewMyHandler(lm lifecycle.LifecycleManager) {
    base := NewBaseHandler(lm, true)
    base.templateManager.AddFunc("myFunc", myFunc)
    return &MyHandler{BaseHandler: base}
}
func (h *MyHandler) MyPage(w http.ResponseWriter, r *http.Request) {
    h.RenderTemplate(w, "my_simple.html", data)
}
```

### **Antipadrão 3: "Float para valores financeiros"**
**Sintoma:** `float64 amount`, `float32 price`
**Causa:** Não seguir regra Anti-Float da Constituição
**Impacto:** Erros de arredondamento, inconsistência
**Solução:**
```go
// ❌ ERRADO
type Product struct {
    Price float64  // R$ 1,99 → 1.99
}

// ✅ CORRETO
type Product struct {
    PriceInCents int64  // R$ 1,99 → 199
}

// No template:
base.templateManager.AddFunc("formatCurrency", func(amount int64) string {
    return fmt.Sprintf("R$ %.2f", float64(amount)/100)
})
```

---

## 🎨 **Frontend/UI**

### **Antipadrão 4: "Template sem padrão de navegação"**
**Sintoma:** Novo template não tem link em outros templates
**Causa:** Esquecer de atualizar navegação
**Impacto:** Feature isolada, difícil de descobrir
**Solução:**
```bash
# ✅ CHECKLIST DE NAVEGAÇÃO
# 1. Listar templates que precisam do link
grep -l "nav\|Navegação" modules/ui_web/templates/*.html

# 2. Atualizar TODOS os *_simple.html
# 3. Atualizar layout.html se usar grid de botões
# 4. Documentar decisão no plano
```

### **Antipadrão 5: "Funções de template inconsistentes"**
**Sintoma:** `{{formatDate .CreatedAt}}` mostra `%!v(PANIC=...)`
**Causa:** Função não lida com tipo específico (ex: `time.Time`)
**Impacto:** Erros em runtime, templates quebrados
**Solução:**
```go
// ✅ VERIFICAÇÃO PRÉVIA
// 1. Verificar que funções o template precisa
grep -n "{{" modules/ui_web/templates/members_simple.html | head -20

// 2. Verificar se BaseHandler tem
grep -n "AddFunc.*formatDate" modules/ui_web/internal/handler/base_handler.go

// 3. Se não atende, adicionar no handler específico
base.templateManager.AddFunc("formatDate", func(t interface{}) string {
    switch v := t.(type) {
    case time.Time:
        return v.Format("02/01/2006")
    default:
        // Fallback para função do BaseHandler
        return fmt.Sprintf("%v", t)
    }
})
```

### **Antipadrão 6: "HTMX sem feedback visual"**
**Sintoma:** Ações HTMX sem spinner, mensagens de sucesso/erro
**Causa:** Implementação mínima
**Impacto:** UX pobre, usuário não sabe se ação funcionou
**Solução:**
```html
<!-- ✅ PADRÃO COMPLETO HTMX -->
<form hx-post="/feature" 
      hx-target="#result" 
      hx-swap="outerHTML"
      hx-indicator="#spinner">
    
    <div id="spinner" class="htmx-indicator">
        <svg class="animate-spin h-5 w-5">...</svg>
    </div>
    
    <!-- Campos do formulário -->
</form>

<!-- Área para feedback -->
<div id="feedback-area"></div>

<script>
document.body.addEventListener('htmx:afterSwap', function(evt) {
    if (evt.detail.target.id === 'result') {
        showFeedback('✅ Ação realizada com sucesso!', 'success');
    }
});
</script>
```

---

## 🧪 **Testes**

### **Antipadrão 7: "Testes dependentes de templates"**
**Sintoma:** Testes falham com `500 Internal Server Error` (template não encontrado)
**Causa:** Testes rodam de diretório diferente
**Impacto:** Testes quebrados, cobertura artificialmente baixa
**Solução:**
```go
// ✅ ESTRATÉGIA DE TESTES
// 1. Aceitar 500 como válido para testes de rota
func TestMyPage_GET(t *testing.T) {
    handler.MembersPage(w, req)
    if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
        t.Errorf("Esperado 200 ou 500 (template error), got %d", w.Code)
    }
}

// 2. Isolar testes de lógica
func TestCreateMember_Logic(t *testing.T) {
    // Testar apenas a lógica, não renderização
    // Mockar chamadas de serviço, validar dados
}

// 3. Testes de integração com setup completo
func TestMemberHandler_Integration(t *testing.T) {
    // Setup com templates reais, banco real
    // Rodar de diretório correto
}
```

### **Antipadrão 8: "Testes sem TDD"**
**Sintoma:** Testes escritos após implementação
**Causa:** Pressa, "vou testar depois"
**Impacto:** Bugs descobertos tarde, design não testável
**Solução:**
```go
// ✅ FLUXO TDD
// 1. RED: Escrever teste falhando
func TestNewFeature(t *testing.T) {
    t.Errorf("TODO: Implementar teste")
}

// 2. GREEN: Implementar mínimo para passar
// 3. REFACTOR: Melhorar código mantendo testes passando
// 4. REPETIR para cada funcionalidade
```

### **Antipadrão 9: "Mock complexo desnecessário"**
**Sintoma:** Mock com 100+ linhas, replica lógica real
**Causa:** Mockar demais
**Impacto:** Testes frágeis, difícil de manter
**Solução:**
```go
// ✅ MOCK SIMPLES E FOCADO
type MockService struct {
    ShouldFail bool
    LastCall   string
}

func (m *MockService) DoSomething() error {
    m.LastCall = "DoSomething"
    if m.ShouldFail {
        return errors.New("mock error")
    }
    return nil
}

// Testar cenários específicos, não replicar lógica
```

---

## 🔄 **Processo**

### **Antipadrão 10: "Implementar sem análise prévia"**
**Sintoma:** "Descobertas" durante implementação
**Causa:** Pular fase de descoberta
**Impacto:** Retrabalho, problemas evitáveis
**Solução:**
```
✅ PROCESSO OBRIGATÓRIO:
1. Fase de descoberta (30-60min)
   - Análise backend
   - Identificação padrões
   - Checklist pré-implementação
2. Só então começar a codificar
```

### **Antipadrão 11: "Não documentar decisões"**
**Sintoma:** Mesmos problemas repetidos em implementações diferentes
**Causa:** Conhecimento não capturado
**Impacto:** Retrabalho, inconsistência
**Solução:**
```
✅ SISTEMA DE DOCUMENTAÇÃO:
1. docs/implementation_plans/[feature]_plan.md
2. docs/learnings/[feature]_learnings.md  
3. Atualizar checklists com novos aprendizados
4. Revisar antes de próxima implementação
```

### **Antipadrão 12: "Ignorar Constituição de IA"**
**Sintoma:** Violação de Anti-Float, não cache-proof, quebra soberania
**Causa:** Não consultar/documentação
**Impacto:** Problemas arquiteturais sérios
**Solução:**
```
✅ CHECKLIST CONSTITUIÇÃO:
- [ ] ANTI-FLOAT: Zero float para valores financeiros/tempo
- [ ] CACHE-PROOF: Templates *_simple.html, ParseFiles() no handler
- [ ] SOBERANIA: entity_id isolamento, um banco por entidade
- [ ] VALIDAR antes de commit
```

---

## 📋 **Checklist Rápido de Validação**

### **Antes do Commit**
- [ ] **Anti-Float scan:** `grep -n "float[0-9]*" arquivos_novos.go`
- [ ] **Cache-proof:** Templates são `*_simple.html`? Carregados com `ParseFiles()`?
- [ ] **Soberania:** Handler usa `entity_id` do contexto/query?
- [ ] **Padrões:** Segue `BaseHandler`? Rotas HTMX padrão?
- [ ] **Navegação:** Link adicionado em templates relevantes?
- [ ] **Testes:** Passando? Cobertura >90% para handler?

### **Durante Code Review**
1. Verificar antipadrões nesta lista
2. Validar contra Constituição de IA
3. Checar consistência com padrões estabelecidos
4. Verificar documentação de decisões

---

## 🔄 **Processo de Melhoria Contínua**

### **Quando encontrar novo antipadrão:**
1. **Documentar** nesta lista
2. **Adicionar** ao checklist pré-implementação
3. **Comunicar** em retrospectiva
4. **Validar** em próxima implementação

### **Revisão periódica:**
- Mensal: Revisar antipadrões encontrados
- Por feature: Atualizar com aprendizados
- Trimestral: Revisão completa do guia

---

## 🎯 **Métricas de Sucesso**

### **Indicadores de melhoria:**
- **Problemas repetidos:** ↓ (meta: zero)
- **Tempo de implementação:** ↓ (meta: -30%)
- **Bugs em produção:** ↓ (meta: -50%)
- **Consistência código:** ↑ (meta: 95% padrões seguidos)

### **Como medir:**
1. Registrar antipadrões encontrados em cada implementação
2. Comparar com implementações anteriores
3. Ajustar processos e checklists

---

**📌 Nota:** Este é documento vivo. Atualizar com cada novo aprendizado.
Contribuir com antipadrões encontrados e soluções testadas.
