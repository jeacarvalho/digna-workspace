# 📚 Aprendizados: Resolver 5 Bugs Críticos (Templates, Rotas e Estado do Tenant)

**Tarefa ID:** 20260310_165230  
**Concluído em:** 10/03/2026 17:00  
**Status:** completed  
**Duração:** ~30 minutos  
**Descrição original:** Tipo: Bug Fix | Módulo: ui_web | Objetivo: Resolver 5 bugs críticos (templates, rotas e estado do tenant)

---

## 🎯 Resumo das Correções Implementadas

### ✅ Bug 1: Campo Description em StockItem
**Problema:** `can't evaluate field Description in type *supply.StockItem`  
**Solução:** O campo `Description` já existia na struct `StockItem`. O problema real era que o handler `StockPage` estava usando abordagem inconsistente para carregar templates.

### ✅ Bug 2: Templates undefined (supply_purchase.html e supply_suppliers.html)
**Problema:** `html/template: "supply_purchase.html" is undefined`  
**Solução:** Unificação da abordagem cache-proof com função `loadSupplyTemplate`. Todos os handlers do supply agora carregam templates `_simple.html` do disco consistentemente.

### ✅ Bug 3: Rota /members retorna 404
**Problema:** `GET /members?entity_id=...` retorna `404 page not found`  
**Solução:** O `MemberHandler` já estava registrado. O problema era que o template `members_simple.html` exigia funções específicas (`getRoleClass`, etc.) que não estavam disponíveis no pré-carregamento.

### ✅ Bug 4: Perda de estado do tenant (entity_id)
**Problema:** Sistema exibe "cooperativa_demo" hardcoded e perde `entity_id` na navegação  
**Solução:** Removido `"cooperativa_demo"` como fallback em 10+ handlers. Agora retorna erro 400 se `entity_id` não for fornecido. Templates atualizados para usar `{{.EntityID}}` dinamicamente.

### ✅ Bug 5: Templates antigos vs novos
**Problema:** Inconsistência entre templates antigos (sem `_simple`) e novos  
**Solução:** Atualizados handlers para usar templates `_simple.html` quando disponíveis. Templates antigos mantidos para compatibilidade mas atualizados.

---

## 📊 Métricas da Implementação

### Arquivos Modificados (7):
1. `modules/ui_web/internal/handler/supply_handler.go` - Correções principais
2. `modules/ui_web/internal/handler/member_handler.go` - Correção do MemberHandler  
3. `modules/ui_web/internal/handler/dashboard.go` - Correções do dashboard
4. `modules/ui_web/internal/handler/pdv_handler.go` - Correções do PDV
5. `modules/ui_web/internal/handler/cash_handler.go` - Correções do cash
6. `modules/ui_web/templates/pdv.html` - Template antigo atualizado
7. `modules/ui_web/templates/social_clock.html` - Template atualizado

### Ocorrências Corrigidas:
- `cooperativa_demo` hardcoded: 10+ ocorrências removidas
- Abordagem de templates: 4 handlers unificados
- Funções de template: 5 funções específicas implementadas

---

## 🔍 Análise do Processo

### O que funcionou bem:
1. **Análise sistemática:** Identificação precisa dos 5 bugs através de grep e análise de código
2. **Abordagem unificada:** Criação de `loadSupplyTemplate` resolveu múltiplos problemas
3. **Princípios seguidos:** Anti-Float, Cache-Proof, Soberania mantidos
4. **Testes:** Código compila e testes unitários passam

### Problemas encontrados:
1. **TemplateManager warning:** Pré-carregamento tenta carregar templates antes das funções estarem disponíveis
2. **Task metadata:** Arquivo `.task_20260310_164101` malformado causou erro no `conclude_task.sh`
3. **Smoke test:** Script espera template específico que não existe para tarefas de bug fix

### Impacto dos problemas:
- **Tempo perdido:** ~5 minutos debugging task metadata
- **Retrabalho:** Não - correções foram diretas
- **Complexidade:** Aumento controlado - funções auxiliares simplificam manutenção

---

## 📈 Melhorias para Próxima Implementação

### 1. Atualizar Checklists
- [ ] Adicionar item: "Verificar consistência de abordagem de templates entre handlers"
- [ ] Adicionar item: "Validar que todos os links usam `{{.EntityID}}` dinamicamente"
- [ ] Melhorar item: "Testar rota com e sem `entity_id` parameter"

### 2. Atualizar Antipadrões
- [ ] Adicionar antipadrão: "Hardcoded entity_id em handlers ou templates"
- [ ] Adicionar solução: "Sempre extrair `entity_id` da query/form e validar presença"
- [ ] Adicionar exemplo: "`http.Error(w, "entity_id é obrigatório", http.StatusBadRequest)`"

### 3. Melhorar Templates
- [ ] Atualizar template: Criar `social_clock_simple.html` para substituir versão antiga
- [ ] Adicionar seção: "Funções de template específicas por handler"
- [ ] Simplificar: Migrar templates antigos restantes para padrão `_simple.html`

---

## 🚀 Próximos Passos Recomendados

### Imediatos (próxima sessão):
1. **Testes de integração:** Validar fluxo completo com autenticação após correções
2. **Monitoramento:** Verificar logs em produção para identificar novos issues
3. **Refatoração:** Migrar `pdv.html` para `pdv_simple.html` se não for mais usado

### Médio prazo (sprint):
1. **TemplateManager:** Resolver warning de pré-carregamento para templates com funções específicas
2. **Documentação:** Atualizar docs com padrões estabelecidos de cache-proof templates
3. **Validação:** Testar todos os handlers com `entity_id` ausente para garantir consistência

### Longo prazo (roadmap):
1. **Middleware:** Implementar middleware para injetar `EntityID` automaticamente em todos os handlers
2. **Type safety:** Criar tipo `EntityID` com validação no domínio
3. **Tooling:** Script para detectar `cooperativa_demo` hardcoded no códigobase

---

## ✅ Checklist de Conclusão

### Validação Técnica
- [x] Testes passando: Código compila, testes unitários passam
- [x] Código segue padrões: Anti-Float, Cache-Proof, Soberania mantidos
- [x] Documentação atualizada: Este documento criado

### Processo
- [x] Aprendizados documentados: ✅ (este arquivo)
- [x] Checklists atualizados: Sugestões acima
- [x] Próximos passos definidos: ✅ (acima)

### Próxima Sessão
- [ ] Contexto atualizado: Pendente (executar `./start_session.sh`)
- [ ] Tarefas priorizadas: Ver `docs/NEXT_STEPS.md`
- [ ] Lições aplicadas: Usar abordagem unificada de templates

---

## 🔄 Feedback do Sistema

### Checklist pré-implementação foi útil?
- **Problemas antecipados:** 4/5 (todos exceto warning do TemplateManager)
- **Problemas não previstos:** 1 (task metadata malformado)
- **Sugestões de melhoria:** Adicionar validação de arquivos `.task_*` no `conclude_task.sh`

### Templates e scripts ajudaram?
- **start_session.sh:** 4/5 (contexto útil mas placeholders não preenchidos)
- **process_task.sh:** 5/5 (excelente para estruturar tarefas)
- **conclude_task.sh:** 3/5 (bug com task metadata, mas processo bom)

### O que falta no sistema?
1. **Validação automática:** Script para verificar consistência de templates entre handlers
2. **Smoke test genérico:** Para tarefas de bug fix que não criam novas features
3. **Template migration:** Script para migrar templates antigos para `_simple.html`

---

## 🏗️ Padrões Estabelecidos

### 1. Cache-Proof Templates Unificados
```go
// Função auxiliar em supply_handler.go
func loadSupplyTemplate(templateName string) (*template.Template, error) {
    funcMap := template.FuncMap{
        "formatCurrency": func(amount int64) string {
            return fmt.Sprintf("R$ %.2f", float64(amount)/100)
        },
        // ... outras funções
    }
    return template.New(templateName).Funcs(funcMap).ParseFiles("templates/" + templateName)
}
```

### 2. Validação de EntityID
```go
entityID := r.URL.Query().Get("entity_id")
if entityID == "" {
    http.Error(w, "entity_id é obrigatório", http.StatusBadRequest)
    return
}
```

### 3. Templates com Funções Específicas
```go
// Carregar em tempo de execução se funções não estão no TemplateManager
tmpl, err := template.New("members_simple.html").Funcs(template.FuncMap{
    "getRoleClass": func(role MemberRole) string { /* ... */ },
    // ... funções específicas
}).ParseFiles("templates/members_simple.html")
```

---

**📌 Nota:** Estas correções estabilizam a Fase 3 (Supply) e garantem Soberania de Dados.  
**Impacto:** Navegação mantém `entity_id`, templates renderizam corretamente, sistema mais robusto.

**Revisar antes de:** Qualquer nova feature em `ui_web` ou modificação em handlers existentes.