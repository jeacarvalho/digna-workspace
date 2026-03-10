Tipo: Bug Fix
Módulo: ui_web  
Objetivo: Resolver 5 bugs críticos (templates, rotas e estado do tenant)
Decisões: Seguir skills: developing-digna-backend, rendering-digna-frontend, managing-sovereign-data

### 📝 Descrição da Tarefa: Correção de Bugs Críticos (Templates, Rotas e Estado do Tenant)
*   **Requisito Funcional (RF):** RF-07 (Compras), RF-08 (Estoque) e RF-01 (Identidade/Navegação).
*   **Sprint Relacionada:** Sprint 17 (Estabilização).

Para esta tarefa, você deve carregar e seguir estritamente as instruções das seguintes skills em docs/skills/:
1.  **[developing-digna-backend]**
2.  **[rendering-digna-frontend]**
3.  **[managing-sovereign-data]**

*   **Anti-Float:** Se envolver cálculos de valor ou tempo, use estritamente int64. Proibido float.
*   **Cache-Proof:** Se houver interface, o template deve ser `_simple.html` carregado via `ParseFiles` no Handler.
*   **Soberania:** Garanta que a operação respeite o isolamento do arquivo `.db` do tenant atual e que o `entity_id` não seja perdido na navegação.

---
**🎯 Objetivo da Tarefa**
Resolver 5 bugs reportados no sistema que estão quebrando a navegação, renderização de templates da Fase 3 (Supply) e a manutenção do estado do tenant (Soberania do Dado). Os erros envolvem campos inexistentes em structs, chamadas para templates antigos não-simples, rotas não registradas (404) e perda do parâmetro `entity_id` na URL.

**📁 Estrutura de Output Esperada**
* `modules/supply/internal/domain/stock_item.go` (ou template correspondente)
* `modules/ui_web/internal/handler/supply_handler.go`
* `modules/ui_web/templates/supply_purchase_simple.html` e `supply_suppliers_simple.html`
* `modules/ui_web/main.go` (ou arquivo de rotas equivalente)
* Templates globais (navbar/header) que contenham links estáticos.

**🛠️ Tarefas de Implementação**

1. **Correção do Bug do Estoque (`<.Description>`):**
   - **Erro:** `can't evaluate field Description in type *supply.StockItem` na linha 221 de `supply_stock_simple.html`.
   - **Ação:** Verificar a struct `StockItem` no domínio. Se o campo não existir, adicione `Description string` à struct (e atualize repositório/banco se necessário), OU corrija o template para usar um campo válido (ex: `.Name`), dependendo do que fizer mais sentido para o negócio.

2. **Correção dos Bugs de Templates Undefined (Compras e Fornecedores):**
   - **Erro:** `html/template: "supply_purchase.html" is undefined` e `"supply_suppliers.html" is undefined`.
   - **Ação:** O `SupplyHandler` está tentando chamar os templates parciais antigos. Atualize o handler para o padrão *Cache-Proof*, carregando diretamente do disco os arquivos `supply_purchase_simple.html` e `supply_suppliers_simple.html`. Crie esses arquivos `_simple.html` herdando o layout completo se eles não existirem.

3. **Correção do Bug 404 em Members:**
   - **Erro:** `GET /members?entity_id=...` retorna `404 page not found`.
   - **Ação:** O `MemberHandler` (criado na etapa anterior) não foi injetado/registrado no roteador principal. Vá ao arquivo de rotas (`main.go` ou `routes.go` do `ui_web`) e registre a rota `/members` apontando para o `MemberHandler`.

4. **Correção do Bug de Perda de Sessão/Tenant (`entity_id`):**
   - **Erro:** O sistema exibe "cooperativa_demo" no topo em vez de "cafe_digna" e perde o parâmetro da URL na navegação.
   - **Ação:** O estado do Tenant (`entity_id`) está sendo perdido. Atualize o `BaseHandler` (ou middleware responsável) para garantir que a variável `EntityID` seja sempre injetada no contexto de dados de *todos* os templates. Varra os templates (especialmente a Navbar/Header em `_simple.html`) e garanta que *todos* os links `href` e requisições HTMX `hx-get` repassem o parâmetro dinamicamente (ex: `href="/supply?entity_id={{.EntityID}}"`), substituindo qualquer hardcode de "cooperativa_demo".

**✅ Critérios de Aceite (Definition of Done)**
- [ ] Acessar `/supply/stock` não gera mais erro de template sobre o campo `Description`.
- [ ] Acessar `/supply/purchase` e `/supply/suppliers` renderiza as páginas corretamente usando o padrão `_simple.html` (Cache-Proof).
- [ ] Acessar `/members` carrega a tela de gestão de membros (sem 404).
- [ ] Ao logar ou acessar `/supply?entity_id=cafe_digna`, o header exibe "cafe_digna", e clicar em qualquer link do menu mantém o `entity_id=cafe_digna` na URL, preservando a Soberania do Dado.

---
1. Código fonte seguindo Clean Architecture (Domain -> Service -> Handler).
2. Testes unitários com TDD provando a lógica (especialmente para o roteamento e a injeção do EntityID).
3. Atualização sugerida para o próximo Session Log.

**Pode iniciar a análise e propor o plano de implementação?**
```
