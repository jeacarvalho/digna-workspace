# 📋 Plano de Implementação: --file=prompt_teste_

**Feature:** --file=prompt_teste_
**Tarefa ID:** 20260311_112301
**Gerado em:** 11/03/2026 11:23:01
**Descrição:** --file=Prompt_teste_correcos.md

**📌 PRÉ-REQUISITO:** Preencher `docs/implementation_plans/--file=prompt_teste__pre_check.md` antes de implementar

---

## 0. 🔍 Fase de Descoberta (A COMPLETAR)

**Checklist:** `docs/implementation_plans/--file=prompt_teste__pre_check.md`
**Status:** □ Não iniciado □ Em progresso ✅ Completo

### **0.1 Backend Status**
- [ ] Serviço existe e testado: □ Sim □ Parcial □ Não
- [ ] Acessível do UI: □ Sim (público) □ Mock necessário □ Internal
- [ ] Padrão de acesso: [API layer / Direct import / Mock inicial]

### **0.2 Padrões Identificados**
- Handler de referência: `__________________________`
- Template base: `__________________________`
- Rotas padrão: `GET /______`, `POST /______`, `POST /______/{id}/______`

### **0.3 Riscos Principais**
1. **Risco:** [Descrição breve] → **Mitigação:** [Ação]
2. **Risco:** [Descrição breve] → **Mitigação:** [Ação]

---

## 1. 🎯 Objetivo da Tarefa

--file=Prompt_teste_correcos.md

**Contexto:** [Completar baseado na análise de descoberta]

---

## 2. 📁 Estrutura de Output Esperada

```
/modules/ui_web/internal/handler/--file=prompt_teste__handler.go
/modules/ui_web/templates/--file=prompt_teste__simple.html
/modules/ui_web/internal/handler/--file=prompt_teste__handler_test.go
/docs/implementation_plans/--file=prompt_teste__implementation_20260311.md
```

---

## 3. 🛠️ Tarefas de Implementação

### **3.1 HTTP Handler (`--file=prompt_teste_Handler`)**
- [ ] Criar controlador estendendo `BaseHandler`
- [ ] Implementar rotas HTMX:
  - `GET /--file=prompt_teste_` (renderiza página)
  - `POST /--file=prompt_teste_` (criação via formulário)
  - `POST /--file=prompt_teste_/{id}/toggle-status` (ação HTMX)
- [ ] Instanciar e consumir serviço correspondente
- [ ] Extrair `entity_id` do contexto/query

### **3.2 Template HTMX (`--file=prompt_teste__simple.html`)**
- [ ] Construir interface com paleta "Soberania e Suor"
- [ ] Incluir header/nav padrão (copiar de `dashboard_simple.html`)
- [ ] Criar formulário assíncrono (HTMX) para adição
- [ ] Implementar lista/cards com: [campos relevantes]
- [ ] Adicionar botões de ação com feedback visual via HTMX swaps

### **3.3 Atualização da Navegação**
- [ ] Inserir link para `/--file=prompt_teste_` no header de `dashboard_simple.html`
- [ ] Replicar navegação em templates principais

### **3.4 Testes TDD**
- [ ] `Test--file=prompt_teste_Handler_List--file=prompt_teste_` - Renderização
- [ ] `Test--file=prompt_teste_Handler_Create--file=prompt_teste_` - Criação via POST
- [ ] `Test--file=prompt_teste_Handler_ToggleStatus` - Alternância de status

---

## 4. ✅ Critérios de Aceite (Definition of Done)

### **Arquitetura**
- [ ] Handler utiliza abordagem cache-proof (`ExecuteTemplate` do `BaseHandler`)
- [ ] Soberania mantida: dados só acessados no arquivo `.sqlite` da entidade
- [ ] Anti-Float compliance: zero `float` para valores financeiros/tempo

### **Frontend**
- [ ] Design segue preceitos de Tecnologia Social
- [ ] Interface acessível com botões grandes e contrastes adequados
- [ ] Feedback amigável para erros

### **Funcionalidade**
- [ ] CRUD completo via HTMX (Create, Read, Update/Delete)
- [ ] Validações capturadas e exibidas como alertas amigáveis
- [ ] Navegação unificada em templates principais

### **Qualidade**
- [ ] Testes unitários com cobertura >90% para handler
- [ ] Testes de integração com banco SQLite real
- [ ] Código segue convenções do projeto

---

## 5. 📅 Cronograma Estimado

1. **Dia 1:** Implementação do Handler e testes unitários
2. **Dia 2:** Desenvolvimento do template `--file=prompt_teste__simple.html`
3. **Dia 3:** Integração com navegação e testes de integração
4. **Dia 4:** Validação final, correções, atualização de documentação

---

## 6. 📝 Código de Referência

### **Estrutura de Handler (baseado em MemberHandler)**
```go
package handler

import (
    "github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

type --file=prompt_teste_Handler struct {
    *BaseHandler
    lifecycleManager lifecycle.LifecycleManager
}

func New--file=prompt_teste_Handler(lm lifecycle.LifecycleManager) (*--file=prompt_teste_Handler, error) {
    base := NewBaseHandler(lm, true)
    return &--file=prompt_teste_Handler{
        BaseHandler:      base,
        lifecycleManager: lm,
    }, nil
}
```

---

## 🚀 PRÓXIMOS PASSOS

1. Completar fase de descoberta (checklist)
2. Revisar e ajustar este plano
3. Iniciar implementação com TDD
4. Documentar aprendizados com `./conclude_task.sh`

---

**📌 Nota:** Atualizar este plano durante a implementação.
Arquivo: docs/implementation_plans/--file=prompt_teste__implementation_20260311.md
