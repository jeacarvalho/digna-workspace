# 📋 Validação Pós-Implementação

**Feature:** `{{feature_name}}`  
**Data:** `{{date}}`  
**Implementador:** `{{implementer}}`

---

## 🎯 OBJETIVO
Validar que a feature implementada **realmente funciona** no ambiente de produção/local, não apenas nos testes.

> ⚠️ **CRÍTICO:** Testes podem passar mas a aplicação quebrar. Esta validação previne isso.

---

## ✅ CRITÉRIOS DE ACEITE (DEFINITION OF DONE ATUALIZADO)

### **NÍVEL 1: Testes Unitários** (Já existente)
- [ ] Handler criado e testado (>90% coverage)
- [ ] Lógica de negócio validada
- [ ] Mocks configurados corretamente

### **NÍVEL 2: Testes de Integração** (Já existente)
- [ ] Banco de dados real (SQLite)
- [ ] Templates carregados
- [ ] Dependências injetadas

### **NÍVEL 3: Testes de Sistema** (NOVO - OBRIGATÓRIO)
- [ ] **Server Integration Test** criado e passando
  - Handler registrado no servidor
  - Rotas respondem HTTP 200
  - Templates compilam sem erro
- [ ] Teste adicionado ao `system_integration_test.go`

### **NÍVEL 4: Smoke Test Local** (NOVO - OBRIGATÓRIO)
- [ ] **Smoke test executado** com sucesso:
  ```bash
  ./scripts/smoke_test_new_feature.sh "{{feature_name}}" "{{main_route}}"
  ```
- [ ] Servidor local responde à rota
- [ ] Template existe e tem conteúdo
- [ ] Navegação atualizada (se aplicável)

### **NÍVEL 5: Testes E2E** (Opcional mas recomendado)
- [ ] Teste E2E criado ou atualizado
- [ ] Playwright valida fluxo completo
- [ ] Teste adicionado ao pipeline CI

---

## 🛠️ COMO EXECUTAR A VALIDAÇÃO

### **Passo 1: Rodar Testes de Sistema**
```bash
cd modules/ui_web
go test -v -run TestSystem
```

**Resultado esperado:** Todos os testes `TestSystem_*` passam.

### **Passo 2: Executar Smoke Test**
```bash
# Exemplo para Members:
./scripts/smoke_test_new_feature.sh "Member Management" "/members"

# Exemplo para Suppliers:
./scripts/smoke_test_new_feature.sh "Supplier Management" "/suppliers"
```

**Resultado esperado:** Todos os checks ✅ verdes.

### **Passo 3: Teste Manual (Opcional mas recomendado)**
1. Iniciar servidor: `cd modules/ui_web && go run .`
2. Acessar no browser: `http://localhost:8090{{main_route}}?entity_id=cooperativa_demo`
3. Validar:
   - Página carrega sem erros
   - Funcionalidades CRUD trabalham
   - Navegação funciona

---

## 🚨 CHECKLIST DE PROBLEMAS COMUNS

### **Problema: Rota 404**
**Sintoma:** `http://localhost:8090/feature` retorna 404  
**Causas:**
- Handler não registrado no `main.go`
- `RegisterRoutes` não chamado ou com rota errada
- Servidor não reiniciado após mudanças

**Solução:**
1. Verificar `modules/ui_web/main.go` - handler deve ser criado e registrado
2. Verificar `RegisterRoutes()` do handler - rota deve ser `/feature` (não `/api/feature`)
3. Reiniciar servidor: `Ctrl+C` e `go run .` novamente

### **Problema: Template undefined**
**Sintoma:** `html/template: "feature_simple.html" is undefined`  
**Causas:**
- Template não existe em `templates/`
- Nome do template errado no handler
- Caminho relativo errado (deve ser `templates/` relativo ao binário)

**Solução:**
1. Verificar se template existe: `ls modules/ui_web/templates/feature_simple.html`
2. Verificar nome no handler: `ExecuteTemplate(w, r, "feature_simple.html", data)`
3. Verificar `TemplateManager` no `BaseHandler` - caminho correto

### **Problema: Handler não compila**
**Sintoma:** `go run .` falha com erro de compilação  
**Causas:**
- Import circular
- Método não implementado
- Tipo errado

**Solução:**
1. Rodar `go build ./...` para ver erros
2. Verificar imports no handler
3. Verificar se implementa interface necessária

---

## 📝 TEMPLATE DE VALIDAÇÃO (Copiar e preencher)

```markdown
# Validação: {{feature_name}}

## ✅ Nível 1: Testes Unitários
- [ ] `Test{{Feature}}Handler_List` - Passa
- [ ] `Test{{Feature}}Handler_Create` - Passa  
- [ ] Coverage: ___% (meta: >90%)

## ✅ Nível 2: Testes de Integração
- [ ] Integração com banco - Passa
- [ ] Templates carregam - Passa

## ✅ Nível 3: Testes de Sistema
- [ ] `TestSystem_HandlerRoutes` - Inclui {{feature}}
- [ ] `TestSystem_TemplatesExist` - Template {{template}} existe
- [ ] `TestSystem_NewHandlerValidation` - Criado para {{feature}}

## ✅ Nível 4: Smoke Test
```bash
./scripts/smoke_test_new_feature.sh "{{feature_name}}" "{{main_route}}"
```
**Output:**
```
✅ Servidor: Rodando
✅ Rota {{main_route}}: Responde 200
✅ Template {{template}}: Existe
✅ Página: HTML válido
```

## ✅ Nível 5: Teste Manual
- [ ] Servidor inicia: `go run .`
- [ ] Browser: `http://localhost:8090{{main_route}}?entity_id=cooperativa_demo`
- [ ] Página carrega sem erros
- [ ] Funcionalidades básicas trabalham

## 🎯 RESULTADO FINAL
**Status:** ✅ PRONTO PARA PRODUÇÃO / ⚠️ PRECISA DE AJUSTES

**Problemas encontrados:**
1. [Descrição do problema] - [Solução aplicada]
2. [Descrição do problema] - [Solução aplicada]

**Aprendizados:**
- [Learning 1]
- [Learning 2]
```

---

## 🔄 INTEGRAÇÃO COM WORKFLOW

Adicione esta validação ao final de cada implementação:

1. **Após implementação:** Executar smoke test
2. **Se falhar:** Corrigir antes de marcar como "concluído"
3. **Documentar:** Adicionar validação ao `conclude_task.sh`

**Comando rápido:**
```bash
# 1. Implementar feature
# 2. Rodar testes
go test ./...

# 3. Validar com smoke test
./scripts/smoke_test_new_feature.sh "Nome da Feature" "/rota"

# 4. Documentar
./conclude_task.sh "Feature implementada e validada com smoke test"
```

---

> **Nota:** Esta validação **elimina** o problema "testes passam mas app quebra". Se o smoke test passar, a feature funciona no ambiente real.