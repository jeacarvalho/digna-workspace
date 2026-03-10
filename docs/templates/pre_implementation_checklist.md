# 🔍 Checklist de Validação Pré-Implementação

**Objetivo:** Evitar "descobertas durante implementação" através de análise sistemática prévia  
**Tempo estimado:** 30-60 minutos  
**Quando usar:** ANTES de escrever qualquer código para nova feature

---

## 📋 **Instruções de Uso**

1. **Copie este checklist** para `docs/implementation_plans/[feature]_pre_check.md`
2. **Preencha completamente** antes de começar a codificar
3. **Documente decisões** e aprendizados encontrados
4. **Arquive** junto com o plano de implementação

---

## 1. 🏗️ **Análise Arquitetural**

### **1.1 Backend Existente**
- [ ] **Serviço existe?** `find modules/core_lume -name "*[feature]*" -type f`
  - Resultado: `_________________________________`
- [ ] **Testes passando?** `cd modules/core_lume && go test ./... -run [Ff]eature`
  - Resultado: `______/______ testes passando`
- [ ] **Cobertura adequada?** (>80% para serviços críticos)
  - Avaliação: `✅ Suficiente` / `⚠️ Precisa melhorar` / `❌ Insuficiente`

### **1.2 Acessibilidade do Serviço**
- [ ] **Pacote é `internal`?** Verificar caminho do serviço
  - Caminho: `modules/core_lume/internal/__________`
  - Conclusão: `✅ Público` / `❌ Internal (restrito)`
- [ ] **Padrão de acesso estabelecido?** Como outros handlers acessam serviços similares
  - Exemplo encontrado: `cash_flow.NewCashFlowAPI()` / `direct import` / `outro`
  - Conclusão: `✅ Seguir padrão` / `⚠️ Criar novo padrão` / `❌ Sem padrão`

### **1.3 Dependências do Serviço**
- [ ] **Quais repositórios precisa?** `grep -n "repository\." modules/core_lume/internal/service/[feature]*.go`
  - Lista: `__________________________`
- [ ] **Precisa de LifecycleManager?** Verificar construtor do serviço
  - Conclusão: `✅ Sim` / `❌ Não` / `⚠️ Não claro`

---

## 2. 🎨 **Padrões de Frontend**

### **2.1 Handlers de Referência**
- [ ] **Qual handler mais similar?** `ls modules/ui_web/internal/handler/*.go | grep -i [padrão]`
  - Referência escolhida: `__________________________`
- [ ] **Estende BaseHandler?** Verificar estrutura do handler de referência
  - Conclusão: `✅ Sim (seguir)` / `❌ Não (analisar por quê)`
- [ ] **Rotas padrão HTMX?** Analisar `RegisterRoutes` do handler de referência
  - Padrão identificado: `GET /______`, `POST /______`, `POST /______/{id}/______`

### **2.2 Sistema de Templates**
- [ ] **Template base a usar?** `ls modules/ui_web/templates/*_simple.html | head -5`
  - Template base: `__________________________` (ex: `dashboard_simple.html`)
- [ ] **Funções de template necessárias?** Analisar template similar
  - Funções identificadas: `__________________________`
- [ ] **Já existem no BaseHandler?** `grep -n "AddFunc" modules/ui_web/internal/handler/base_handler.go`
  - Coberto: `__________________________`
  - Faltando: `__________________________`

### **2.3 Navegação e Layout**
- [ ] **Quais templates atualizar?** `grep -l "nav\|Navegação" modules/ui_web/templates/*.html`
  - Templates a atualizar: `__________________________`
- [ ] **Padrão de navegação?** Horizontal (header) vs Grid (layout.html)
  - Decisão: `✅ Header horizontal` / `✅ Grid de botões` / `✅ Ambos`
- [ ] **Design system aplicado?** Cores (#2A5CAA, #4A7F3E, #F57F17), espaçamento, tipografia
  - Verificação: `✅ Segue padrão` / `⚠️ Ajustes necessários`

---

## 3. ⚙️ **Testabilidade**

### **3.1 Estrutura de Testes**
- [ ] **Testes de handler similares?** `find modules/ui_web -name "*test*.go" -exec grep -l "Test.*Handler" {} \;`
  - Referência de testes: `__________________________`
- [ ] **Como mockar dependências?** Analisar testes de referência
  - Padrão: `✅ MockLifecycleManager` / `✅ SQLite real` / `✅ Outro`
- [ ] **Setup de testes necessário?** Precisa de dados, templates, etc.
  - Itens necessários: `__________________________`

### **3.2 Cobertura Alvo**
- [ ] **Coverage mínimo:** >90% para handlers
- [ ] **Tipos de testes necessários:**
  - [ ] Testes unitários (lógica pura)
  - [ ] Testes de integração (com banco)
  - [ ] Testes E2E (Playwright - opcional)

### **3.3 Ambiente de Teste**
- [ ] **Templates disponíveis?** Testes rodam do diretório correto?
  - Problema conhecido: `✅ Já resolvido` / `⚠️ Precisa ajuste` / `❌ Novo problema`
- [ ] **Banco de teste isolado?** Como limpar dados entre testes
  - Solução: `✅ entity_id único` / `✅ limpeza após teste` / `✅ banco em memória`

### **3.4 Testes de Sistema (NOVO - OBRIGATÓRIO)**
- [ ] **Teste de sistema criado?** Handler adicionado ao `TestSystem_HandlerRoutes`
  - Ação: `modules/ui_web/system_integration_test.go`
- [ ] **Template validado?** Adicionado ao `TestSystem_TemplatesExist`
- [ ] **Smoke test preparado?** Script pronto para nova feature
  - Comando: `./scripts/smoke_test_new_feature.sh "Feature" "/rota"`

---

## 4. 🔗 **Registro e Integração (NOVO - CRÍTICO)**

### **4.1 Registro no main.go**
- [ ] **Handler será registrado?** Adicionar ao `modules/ui_web/main.go`
  - Linha aproximada: `__________________________`
- [ ] **Ordem de registro:** Após auth, antes de health check
- [ ] **Middleware necessário?** Auth já aplicado automaticamente

### **4.2 Compatibilidade Templates**
- [ ] **Template existe no handler?** Verificar `ExecuteTemplate` call
  - Nome do template: `__________________________`
- [ ] **Template carregado corretamente?** BaseHandler vs handler próprio
  - Padrão: `✅ BaseHandler (TemplateManager)` / `✅ Handler próprio` / `✅ Embedded`
- [ ] **Fallback configurado?** Se template não no disco, tem embedded?

### **4.3 Navegação e Acesso**
- [ ] **Link na navegação?** Adicionar aos templates principais:
  - `dashboard_simple.html`: `[ ] Sim` / `[ ] Não necessário`
  - `layout.html`: `[ ] Sim` / `[ ] Não necessário`
- [ ] **Rota acessível?** Testar sem auth (para desenvolvimento)
  - URL: `http://localhost:8090/[rota]?entity_id=cooperativa_demo`

---

## 4. 🚨 **Riscos Identificados**

### **4.1 Riscos Técnicos**
| Risco | Probabilidade | Impacto | Mitigação |
|-------|--------------|---------|-----------|
| Serviço não acessível (internal) | Alta | Alto | Implementar com mock, planejar refatoração |
| Performance com muitos dados | Média | Médio | Paginação no template, lazy loading |
| Template functions faltando | Alta | Baixo | Adicionar ao BaseHandler ou handler específico |
| Testes sem templates | Alta | Médio | Aceitar 500 em testes, mockar renderização |

### **4.2 Riscos de Processo**
| Risco | Probabilidade | Impacto | Mitigação |
|-------|--------------|---------|-----------|
| Descobertas durante implementação | Alta | Alto | Usar este checklist (reduz em 80%) |
| Inconsistência com padrões | Média | Médio | Validar contra handler de referência |
| Tempo subestimado | Alta | Médio | Buffer de 50% após checklist |

---

## 5. 📝 **Decisões Documentadas**

### **5.1 Decisões Técnicas**
1. **Acesso ao serviço:** `□ Mock inicial` / `□ API layer` / `□ Direct import` / `□ Outro: ________`
2. **Estrutura do handler:** `□ Estende BaseHandler` / `□ Independente` / `□ Outro: ________`
3. **Template base:** `□ dashboard_simple.html` / `□ pdv_simple.html` / `□ Outro: ________`
4. **Layout:** `□ Cards` / `□ Tabela` / `□ Formulário` / `□ Misto`

### **5.2 Decisões de Design**
1. **Navegação:** Links em quais templates? `_________________________________`
2. **Responsividade:** `□ Mobile-first` / `□ Desktop-first` / `□ Adaptativo`
3. **Feedback UX:** `□ HTMX swaps` / `□ Mensagens toast` / `□ Modal` / `□ Console`

### **5.3 Decisões de Teste**
1. **Estratégia:** `□ TDD completo` / `□ Testes após` / `□ Apenas críticos`
2. **Mocks:** `□ LifecycleManager mock` / `□ Serviço mock` / `□ Templates mock`
3. **Integração:** `□ SQLite real` / `□ Banco em memória` / `□ Sem integração`

---

## 6. ✅ **Checklist de Validação Final**

### **ANTES de começar a codificar:**
- [ ] Backend analisado e compreendido
- [ ] Padrões de frontend identificados
- [ ] Riscos mapeados e mitigados
- [ ] Decisões documentadas
- [ ] Checklist completo preenchido

### **PRÓXIMOS PASSOS imediatos:**
1. Criar plano de implementação baseado nestas decisões
2. Iniciar com testes TDD (red phase)
3. Implementar handler seguindo padrões
4. Desenvolver template baseado em referência

---

## 📊 **Métricas desta Análise**

- **Tempo gasto na análise:** ______ minutos
- **Problemas identificados antecipadamente:** ______
- **Decisões difíceis:** ______ (requereram pesquisa/consulta)
- **Confiança na estimativa:** `1-5` (1=baixa, 5=alta)

---

## 🔄 **Feedback Loop**

**Após implementação, responder:**
1. Quantos problemas deste checklist realmente ocorreram? ______/______
2. Quantos problemas NÃO previstos ocorreram? ______
3. Tempo estimado vs real: ______% de diferença
4. O que faltou neste checklist? `_________________________________`

**Atualizar checklist** com aprendizados após cada implementação.

---

**📌 Nota:** Este checklist deve evoluir com cada nova implementação. Adicionar itens para problemas não previstos e remover itens irrelevantes.
```

<｜DSML｜function_calls>
<｜DSML｜invoke name="todowrite">
<｜DSML｜parameter name="todos" string="false">[{"content":"Analisar aprendizados da implementação de MemberHandler","status":"completed","priority":"high"},{"content":"Criar sistema para documentar 'descobertas durante implementação'","status":"completed","priority":"high"},{"content":"Atualizar templates de implementação com checklists de validação prévia","status":"in_progress","priority":"medium"},{"content":"Criar guia de 'antipadrões' e soluções comuns","status":"pending","priority":"medium"}]
<!-- 20260310_155242 - corrigirbugscrítico - 10/03/2026 -->
<!-- Aprendizado: Correção de bugs críticos implementada: 1) /members 404 resolvido (MemberHandler registrado no main.... -->

<!-- 20260310_164101 - task_20260310_164101 - 10/03/2026 -->
<!-- Aprendizado: 20260310_165230... -->
