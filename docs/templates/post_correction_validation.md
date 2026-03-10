# 🔧 Validação Pós-Correção de Bugs

**Baseado nos aprendizados das correções de:** `/members 404` e `/supply/purchase template undefined`

**Data:** `{{date}}`
**Correção aplicada por:** `{{corrector}}`

---

## 🎯 OBJETIVO
Validar que correções de bugs **realmente funcionam** e não introduzem regressões, com foco nos **problemas comuns identificados**.

> ⚠️ **CRÍTICO:** Bugs corrigidos devem ser validados com o mesmo rigor que novas features.

---

## ✅ CHECKLIST DE VALIDAÇÃO PÓS-CORREÇÃO

### **1. 🔍 Problema Original**
- [ ] **Descrição clara:** O que quebrava? Qual erro?
- [ ] **Cenário de reprodução:** Passos para reproduzir o bug
- [ ] **Impacto:** Quem/que funcionalidade era afetada?

### **2. 🛠️ Análise da Causa Raiz**
- [ ] **Causa identificada:** Por que acontecia?
  - □ Handler não registrado
  - □ Template não carregado  
  - □ Import circular
  - □ Middleware conflito
  - □ Outro: ________________
- [ ] **Solução aplicada:** O que foi feito para corrigir?

### **3. ✅ Validação da Correção**
#### **3.1 Correção Funciona?**
- [ ] **Bug reproduzível ANTES:** Sim / Não
- [ ] **Bug corrigido DEPOIS:** Sim / Não
- [ ] **Teste manual executado:** Descrever cenário testado

#### **3.2 Testes Atualizados**
- [ ] **Testes unitários atualizados:** Cobrem o cenário do bug?
- [ ] **Testes de sistema atualizados:** `TestSystem_*` inclui validação?
- [ ] **Novos testes criados:** Para prevenir regressão?

#### **3.3 Smoke Test Específico**
```bash
# Exemplo para correção de /members 404:
./scripts/smoke_test_new_feature.sh "Member Bug Fix" "/members"

# Exemplo para correção de template:
curl -I "http://localhost:8090/supply/purchase?entity_id=cooperativa_demo"
```
- [ ] Smoke test executado e passou?

### **4. 🚨 Prevenção de Regressões**
#### **4.1 Checklist Atualizado**
- [ ] **Checklist pré-implementação atualizado?** 
  - Item "Registro no main.go" adicionado?
  - Item "Compatibilidade templates" adicionado?
- [ ] **Template atualizado?** `docs/templates/pre_implementation_checklist.md`

#### **4.2 Validação no Fluxo**
- [ ] **start_session.sh:** Inclui verificação de handlers registrados?
- [ ] **conclude_task.sh:** Exige validação antes de concluir?
- [ ] **process_task.sh:** Gera checklist com itens críticos?

#### **4.3 Documentação**
- [ ] **Aprendizados documentados:** `docs/learnings/`
- [ ] **Antipattern atualizado:** `docs/antipatterns/`
- [ ] **Quick reference atualizado:** Menção ao problema/solução?

### **5. 🔄 Impacto no Sistema**
#### **5.1 Compatibilidade com Versões**
- [ ] **Backward compatible?** Não quebra funcionalidades existentes
- [ ] **Dependências atualizadas?** Nenhuma nova dependência crítica
- [ ] **Performance impact?** Negligenciável / Pequeno / Médio

#### **5.2 Cobertura de Testes**
- [ ] **Cobertura antes:** ______%
- [ ] **Cobertura depois:** ______%
- [ ] **Diferença:** ______% (deve ser ≥ 0)

---

## 📝 TEMPLATE DE RELATÓRIO DE CORREÇÃO

```markdown
# Correção: {{nome_do_bug}}

## 📋 Informações
- **Data da correção:** {{data}}
- **Responsável:** {{nome}}
- **Issue/Ticket:** {{referencia}}

## 🔍 Problema Original
**Descrição:** {{descrição_do_problema}}

**Erro:**
```
{{código_do_erro}}
```

**Cenário de reprodução:**
1. {{passo_1}}
2. {{passo_2}}
3. {{passo_3}}

## 🛠️ Análise e Solução
**Causa raiz:** {{explicação_da_causa}}

**Solução aplicada:**
```{{linguagem}}
{{código_da_solução}}
```

**Arquivos modificados:**
- {{arquivo_1}}: {{alteração}}
- {{arquivo_2}}: {{alteração}}

## ✅ Validação
### Testes Executados:
- [ ] Teste manual: {{cenário}}
- [ ] Testes unitários: {{quantidade}} criados/atualizados
- [ ] Testes de sistema: `TestSystem_*` atualizado
- [ ] Smoke test: Executado com sucesso

### Resultados:
- **Status antes:** ❌ Falhava com {{erro}}
- **Status depois:** ✅ Funciona corretamente
- **Performance:** {{impacto_performance}}

## 🚨 Prevenção de Regressões
### Atualizações no Fluxo:
- [ ] Checklist pré-implementação atualizado
- [ ] Validação obrigatória no conclude_task.sh
- [ ] Verificação no start_session.sh

### Novos Mecanismos:
1. {{mecanismo_1}}
2. {{mecanismo_2}}

## 📊 Métricas
- **Linhas de código:** +{{adicionadas}} -{{removidas}}
- **Tempo de correção:** {{horas}} horas
- **Cobertura de testes:** {{antes}}% → {{depois}}%

## 🎯 Aprendizados
1. {{aprendizado_1}}
2. {{aprendizado_2}}
3. {{aprendizado_3}}

**Status final:** ✅ CORRIGIDO E VALIDADO
```

---

## 🚀 INTEGRAÇÃO COM WORKFLOW EXISTENTE

### **Quando usar esta validação:**
1. **Após correção de bug crítico** (404, template errors, etc.)
2. **Após refatoração que afeta múltiplos handlers**
3. **Após atualização de dependências principais**
4. **Sempre que o bug era "testes passam mas app quebra"**

### **Comandos Rápidos:**
```bash
# 1. Corrigir bug
# 2. Executar validação
./scripts/smoke_test_new_feature.sh "Bug Fix" "/rota_problema"

# 3. Executar testes de sistema
cd modules/ui_web && go test -v -run TestSystem

# 4. Validar registro no main.go
grep -n "New.*Handler" modules/ui_web/main.go

# 5. Documentar correção
./conclude_task.sh "Correção de bug: descrição. Aprendidos: ..."
```

### **Checklist Rápido Pós-Correção:**
```bash
✅ Bug reproduzível antes
✅ Bug corrigido depois  
✅ Handler registrado (se aplicável)
✅ Template carregado (se aplicável)
✅ Testes atualizados/criados
✅ Smoke test executado
✅ Nenhuma regressão introduzida
✅ Aprendizados documentados
```

---

## 📈 MELHORIAS NO PROCESSO (BASEADO NOS APRENDIZADOS)

### **Problemas Comuns Identificados:**
1. **Handler não registrado no main.go** → ✅ Agora checklist obrigatório
2. **Template referenciado mas não carregado** → ✅ Validação template-handler
3. **Testes passam mas app quebra** → ✅ Smoke test obrigatório

### **Novos Mecanismos Implementados:**
1. **Validação no start_session.sh** - Verifica handlers registrados
2. **Checklist expandido** - Inclui registro e compatibilidade
3. **Validação obrigatória no conclude_task.sh** - Não deixa concluir sem validar
4. **Testes de sistema** - Valida integração real

### **Fluxo Atualizado:**
```
ANTES: Implementar → Testes unitários → Marcar como concluído
AGORA: Implementar → Testes unitários → ✅ Registrar no main.go → 
       ✅ Validar templates → ✅ Testes de sistema → 
       ✅ Smoke test → ✅ Validação obrigatória → Concluir
```

> **Nota:** Esta validação **elimina** bugs do tipo "esqueci de registrar no main.go" ou "template não carregado". Se passar por todas as etapas, a correção está **garantidamente funcional**.