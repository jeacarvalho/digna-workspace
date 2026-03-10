# 📋 Template: Plano de Implementação

**Feature:** [Nome da Feature]
**Requisito Funcional:** RF-XX
**Sprint Relacionada:** Sprint XX
**Skills Aplicáveis:** [developing-digna-backend, rendering-digna-frontend, managing-sovereign-data]

**📌 PRÉ-REQUISITO:** Preencher `docs/templates/pre_implementation_checklist.md` antes deste plano

---

## 0. 🔍 **Fase de Descoberta (JÁ COMPLETADA)**

**Arquivo de análise:** `docs/implementation_plans/[feature]_pre_check.md`
**Tempo gasto:** ______ minutos
**Problemas antecipados:** ______
**Decisões críticas:** [Listar 2-3 decisões mais importantes]

### **0.1 Backend Status**
- [ ] Serviço existe e testado: `✅ Sim` / `⚠️ Parcial` / `❌ Não`
- [ ] Acessível do UI: `✅ Sim (público)` / `⚠️ Mock necessário` / `❌ Internal`
- [ ] Padrão de acesso: `[API layer / Direct import / Mock inicial]`

### **0.2 Padrões Identificados**
- Handler de referência: `__________________________`
- Template base: `__________________________`
- Rotas padrão: `GET /______`, `POST /______`, `POST /______/{id}/______`

### **0.3 Riscos Principais**
1. **Risco:** [Descrição breve] → **Mitigação:** [Ação]
2. **Risco:** [Descrição breve] → **Mitigação:** [Ação]

---

## 1. 🎯 **Objetivo da Tarefa**

[Descrição concisa baseada na análise de descoberta]

**Exemplo:** "Na Sprint 10, o backend de [Feature] foi implementado e 100% testado. Análise prévia identificou que o serviço é `internal`, portanto implementaremos com dados mockados inicialmente, seguindo padrão do `CashHandler`. Template base será `dashboard_simple.html`."

---

## 2. 📁 **Estrutura de Output Esperada**

```
/modules/ui_web/internal/handler/[feature]_handler.go
/modules/ui_web/templates/[feature]_simple.html
[Adicionais se necessário]
```

---

## 3. 🛠️ **Tarefas de Implementação**

### **3.1 HTTP Handler (`[Feature]Handler`)**
- [ ] Criar controlador estendendo `BaseHandler` (herda funções de template)
- [ ] Implementar rotas HTMX:
  - `GET /[feature]` (renderiza página)
  - `POST /[feature]` (criação via formulário)
  - `POST /[feature]/{id}/toggle-status` (ação HTMX para ativar/inativar)
- [ ] Instanciar e consumir `[Feature]Service` com base no banco do Tenant isolado
- [ ] Extrair `entity_id` do contexto: `r.Context().Value("entity_id")`

### **3.2 Template HTMX (`[feature]_simple.html`)**
- [ ] Construir interface com paleta "Soberania e Suor"
- [ ] Incluir header/nav padrão (copiar de `dashboard_simple.html`)
- [ ] Criar formulário assíncrono (HTMX) para adição sem recarregar página
- [ ] Implementar lista/tabela com: [campos relevantes]
- [ ] Adicionar botões de ação com feedback visual via HTMX swaps

### **3.3 Atualização da Navegação**
- [ ] Inserir link para `/[feature]` no header de `dashboard_simple.html`
- [ ] Replicar navegação em templates principais (`pdv_simple.html`, `cash_simple.html`, etc.)

### **3.4 Testes TDD**
- [ ] `Test[Feature]Handler_List[Feature]` - Renderização da página
- [ ] `Test[Feature]Handler_Create[Feature]` - Criação via POST HTMX
- [ ] `Test[Feature]Handler_ToggleStatus` - Alternância de status
- [ ] `Test[Feature]Handler_[RegraCrítica]` - Validação de regra de negócio específica

---

## 4. ✅ **Critérios de Aceite (Definition of Done)**

### **Arquitetura**
- [ ] Handler utiliza exclusivamente abordagem cache-proof (`ExecuteTemplate` do `BaseHandler`)
- [ ] Soberania mantida: dados só acessados no arquivo `.sqlite` da entidade atual
- [ ] Anti-Float compliance: zero `float` para valores financeiros/tempo

### **Frontend**
- [ ] Design segue preceitos de Tecnologia Social (sem jargões técnicos)
- [ ] Interface acessível com botões grandes e contrastes adequados
- [ ] Feedback amigável para erros (ex: "Não é possível inativar o último coordenador")

### **Funcionalidade**
- [ ] CRUD completo via HTMX (Create, Read, Update/Delete)
- [ ] Validações do Service capturadas e exibidas como alertas amigáveis
- [ ] Navegação unificada em todos os templates principais

### **Qualidade**
- [ ] Testes unitários com cobertura >90% para handler
- [ ] Testes de integração com banco SQLite real
- [ ] **Testes de sistema** (`TestSystem_*`) criados e passando
- [ ] **Smoke test** executado com sucesso: `./scripts/smoke_test_new_feature.sh`
- [ ] Código segue convenções do projeto (gofmt, snake_case para arquivos)

---

## 5. 🔍 **Análise do Estado Atual**

### **Backend Existente?** ✅/❌
- [ ] Domain layer (`internal/domain/[feature].go`)
- [ ] Service layer (`internal/service/[feature]_service.go`)
- [ ] Repository layer (`internal/repository/[feature]_repository.go`)
- [ ] Testes passando? (__/__ testes)

### **Padrões Aplicáveis**
- [ ] Segue estrutura DDD (Domain → Service → Repository → Handler)
- [ ] Usa `LifecycleManager` para isolamento de banco
- [ ] Implementa regras de negócio críticas no Service

---

## 6. ⚠️ **Riscos e Mitigações**

| Risco | Probabilidade | Impacto | Mitigação |
|-------|--------------|---------|-----------|
| Conflito com autenticação | Baixa | Alto | Usar `middleware.AuthMiddleware` já configurado |
| Performance com muitos registros | Média | Médio | Implementar paginação no template (futuro) |
| Validação de dados duplicados | Alta | Baixo | Service já valida, handler captura e exibe erro |
| Acesso cross-tenant | Baixa | Crítico | EntityID do contexto garante isolamento |

---

## 7. 📅 **Cronograma Estimado**

1. **Dia 1:** Implementação do Handler e testes unitários
2. **Dia 2:** Desenvolvimento do template `[feature]_simple.html`
3. **Dia 3:** Integração com navegação e testes de integração
4. **Dia 4:** Validação final, correções, atualização de documentação

---

## 8. 📝 **Código de Referência**

### **Estrutura de Handler (exemplo)**
```go
type [Feature]Handler struct {
    *BaseHandler
    [feature]Service *service.[Feature]Service
}

func (h *[Feature]Handler) List[Feature](w http.ResponseWriter, r *http.Request) {
    entityID := getEntityIDFromContext(r.Context())
    items, _ := h.[feature]Service.ListByEntity(entityID)
    
    data := map[string]interface{}{
        "EntityID": entityID,
        "Items":    items,
        "Title":    "[Título Amigável]",
    }
    
    h.templateManager.ExecuteTemplate(w, "[feature]_simple.html", data)
}
```

### **Template Base (exemplo)**
```html
<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <!-- Configuração padrão (copiar de dashboard_simple.html) -->
</head>
<body class="min-h-screen">
    <!-- Header padrão com navegação -->
    
    <main class="container mx-auto px-4 py-8">
        <h1 class="text-3xl font-bold text-digna-text mb-8">{{.Title}}</h1>
        
        <!-- Formulário HTMX -->
        <form hx-post="/[feature]" hx-target="#items-list" hx-swap="outerHTML">
            <!-- Campos do formulário -->
        </form>
        
        <!-- Lista de itens -->
        <div id="items-list">
            {{range .Items}}
            <!-- Item com ações HTMX -->
            {{end}}
        </div>
    </main>
</body>
</html>
```

---

## 9. 🔄 **Atualizações Pós-Implementação**

### **9.1 Documentação**
- [ ] Atualizar `docs/QUICK_REFERENCE.md` com novos padrões
- [ ] Atualizar `docs/06_roadmap/05_session_log.md` com resumo da implementação
- [ ] Verificar se novas skills são necessárias em `docs/skills/`

### **9.2 Validação Técnica** (ATUALIZADO)
- [ ] **Testes de sistema:** `go test -v -run TestSystem` (NOVO - obrigatório)
- [ ] **Smoke test:** `./scripts/smoke_test_new_feature.sh` (NOVO - obrigatório)
- [ ] Executar todos os testes: `cd modules && ./run_tests.sh`
- [ ] Validar cache-proof: templates carregados do disco
- [ ] Validar soberania: isolamento por entity_id
- [ ] Validar anti-float: zero floats em código novo

### **9.3 Retrospectiva e Aprendizados**
- [ ] **Problemas previstos vs reais:** ______/______ (quantos do checklist ocorreram?)
- [ ] **Problemas não previstos:** ______ (listar abaixo)
- [ ] **Tempo estimado vs real:** ______% diferença
- [ ] **Decisões acertadas:** [Listar 2-3]
- [ ] **Decisões para revisar:** [Listar 1-2]

#### **Aprendizados para próxima implementação:**
1. **Descoberta:** [O que aprendemos?]
   - **Impacto:** [Como afetou o projeto?]
   - **Solução:** [O que fizemos?]
   - **Prevenção:** [Como evitar no futuro?]

2. **Descoberta:** [O que aprendemos?]
   - **Impacto:** [Como afetou o projeto?]
   - **Solução:** [O que fizemos?]
   - **Prevenção:** [Como evitar no futuro?]

### **9.4 Atualização de Checklists**
- [ ] Revisar `docs/templates/pre_implementation_checklist.md`
  - [ ] Adicionar itens para problemas não previstos
  - [ ] Remover itens irrelevantes
  - [ ] Melhorar itens confusos
- [ ] Atualizar `docs/learnings/[feature]_implementation_learnings.md`

---

## 10. 📈 **Métricas da Implementação**

### **Código**
- Linhas de código: ______ (handler: ______, template: ______, testes: ______)
- Tempo total: ______ horas (descoberta: ______, implementação: ______, testes: ______)
- Commits: ______

### **Qualidade**
- Testes: ______ unitários + ______ integração
- Cobertura: >______%
- Bugs encontrados: ______ (críticos: ______, médios: ______, menores: ______)

### **Processo**
- Checklist útil? `1-5` (1=não, 5=muito)
- Problemas antecipados: ______%
- Melhoria estimada para próxima: ______%

---

**📌 Nota:** Este template deve ser preenchido no início de cada nova implementação significativa. 
1. **Primeiro:** Preencher checklist pré-implementação
2. **Depois:** Preencher este plano baseado na análise
3. **Após:** Completar retrospectiva e atualizar checklists

Arquivar em `docs/implementation_plans/[feature]_implementation_[data].md`