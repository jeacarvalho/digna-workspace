# 🧪 Checklist de Validação E2E - Template

**Feature:** `[NOME_DA_FEATURE]`  
**Tarefa ID:** `[TASK_ID]`  
**Data da Validação:** `[DATA]`  
**Responsável:** `[NOME]`

---

## 📋 **INFORMAÇÕES DA VALIDAÇÃO**

### **Configuração:**
- [ ] **Modo de execução:** □ Headless (stealth) □ Headed (UI) □ UI Mode
- [ ] **Navegador:** □ Chrome □ Firefox □ Todos
- [ ] **Timeout:** `[ ]` segundos
- [ ] **Servidor:** `http://localhost:8090` ✅ Disponível

### **Credenciais de Teste:**
- [ ] **Entidade:** `cafe_digna` ✅ Configurada
- [ ] **Senha:** `cd0123` ✅ Disponível
- [ ] **Dados de teste:** □ Criados automaticamente □ Pré-existentes

---

## 🎯 **FLUXO DE 7 PASSOS PADRÃO DIGNA**

### **1. ✅ Login no Sistema**
- [ ] Acessa `/login`
- [ ] Seleciona empresa `cafe_digna`
- [ ] Insere senha `cd0123`
- [ ] Redireciona para `/dashboard`
- [ ] **Evidência:** Screenshot do dashboard carregado

### **2. ✅ Navegação no Dashboard**
- [ ] Menu de navegação visível
- [ ] Links funcionais (≥ 8 links)
- [ ] Título "Painel de Dignidade" presente
- [ ] **Evidência:** Contagem de links: `[ ]`

### **3. ✅ Acesso à Nova Feature**
- [ ] Link para `/[feature]` presente no menu
- [ ] Página carrega sem erros
- [ ] Título da página correto
- [ ] **Evidência:** Screenshot da página `/[feature]`

### **4. ✅ Funcionalidade Principal**
- [ ] Listagem de itens funciona
- [ ] Formulário de criação funciona (se aplicável)
- [ ] Ações (editar/excluir) funcionam (se aplicável)
- [ ] **Evidência:** Item criado/listado com sucesso

### **5. ✅ Integração com Sistema**
- [ ] Navegação de/para outras páginas funciona
- [ ] Dados persistem entre navegações
- [ ] Não quebra funcionalidades existentes
- [ ] **Evidência:** Navegação completa sem erros

### **6. ✅ Validação de Dados**
- [ ] Dados exibidos corretamente
- [ ] Formatação adequada (moeda, datas, etc.)
- [ ] Anti-float compliance (sem `float` para valores)
- [ ] **Evidência:** Valores formatados corretamente

### **7. ✅ Performance e UX**
- [ ] Carregamento rápido (< 3s)
- [ ] Feedback visual para ações
- [ ] Mensagens de erro amigáveis
- [ ] **Evidência:** Tempo de carregamento: `[ ]`s

---

## 🔧 **VALIDAÇÃO TÉCNICA**

### **Playwright Execution:**
- [ ] **Comando executado:** `./scripts/dev/validate_e2e.sh --basic --headless`
- [ ] **Status de saída:** □ 0 (sucesso) □ 1 (falha) □ 2 (avisos)
- [ ] **Tempo de execução:** `[ ]` minutos
- [ ] **Testes executados:** `[ ]` de `[ ]` passaram

### **Logs e Evidências:**
- [ ] Screenshots em caso de falha: `test-results/`
- [ ] Vídeos em caso de falha: `test-results/`
- [ ] Logs de console disponíveis
- [ ] **Arquivo de relatório:** `[CAMINHO_DO_RELATÓRIO]`

### **Problemas Encontrados:**
| Problema | Severidade | Status | Ação |
|----------|------------|--------|------|
| `[DESCRIÇÃO]` | □ Crítico □ Alto □ Médio □ Baixo | □ Corrigido □ Pendente | `[AÇÃO]` |
| `[DESCRIÇÃO]` | □ Crítico □ Alto □ Médio □ Baixo | □ Corrigido □ Pendente | `[AÇÃO]` |

---

## 📊 **MÉTRICAS DE QUALIDADE**

### **Cobertura de Testes:**
- [ ] **Testes unitários:** `[ ]`% cobertura
- [ ] **Testes integração:** `[ ]`% cobertura  
- [ ] **Testes E2E:** `[ ]`/7 passos validados

### **Performance:**
- [ ] **Tempo médio de resposta:** < 2s
- [ ] **Uso de memória:** Estável
- [ ] **Erros no console:** 0

### **Acessibilidade (WCAG):**
- [ ] **Contraste adequado:** □ Sim □ Parcial □ Não
- [ ] **Navegação por teclado:** □ Funcional □ Parcial □ Não
- [ ] **Labels para screen readers:** □ Presentes □ Parcial □ Ausentes

---

## 🚨 **CRITÉRIOS DE REJEIÇÃO**

### **REJEITAR SE:**
- [ ] **Qualquer** teste do fluxo básico falhar
- [ ] **Login** não funcionar
- [ ] **Navegação** quebrar funcionalidades existentes
- [ ] **Dados** não persistirem corretamente
- [ ] **Performance** degradar significativamente

### **APROVAR SE:**
- [ ] **Todos** os 7 passos básicos passarem
- [ ] **Nenhum** regressão introduzida
- [ ] **Documentação** atualizada
- [ ] **Aprendizados** registrados

---

## 📝 **DECISÕES E APRENDIZADOS**

### **Decisões Técnicas:**
1. **Problema:** `[DESCRIÇÃO]`  
   **Decisão:** `[SOLUÇÃO_ESCOLHIDA]`  
   **Justificativa:** `[RAZÃO]`

2. **Problema:** `[DESCRIÇÃO]`  
   **Decisão:** `[SOLUÇÃO_ESCOLHIDA]`  
   **Justificativa:** `[RAZÃO]`

### **Aprendizados:**
- `[APRENDIZADO_1]`
- `[APRENDIZADO_2]`
- `[APRENDIZADO_3]`

### **Melhorias para Próxima Validação:**
- `[MELHORIA_1]`
- `[MELHORIA_2]`

---

## ✅ **RESULTADO FINAL**

### **Status da Validação:**
- [ ] **✅ APROVADO** - Todos critérios atendidos
- [ ] **⚠️  APROVADO COM RESSALVAS** - Problemas menores documentados
- [ ] **❌ REJEITADO** - Critérios críticos não atendidos

### **Assinaturas:**
```
Validador: _________________________
Data: ______/______/______
Hora: ______:______

Revisor: ___________________________
Data: ______/______/______  
Hora: ______:______
```

### **Próximos Passos:**
1. [ ] Documentar no `conclude_task.sh`
2. [ ] Atualizar métricas de qualidade
3. [ ] Preparar para deploy (se aplicável)
4. [ ] Agendar re-validação (se necessário)

---

**📌 Notas:**  
- Este checklist deve ser preenchido após cada validação E2E
- Arquivos de evidência devem ser mantidos por 30 dias
- Problemas críticos devem ser corrigidos antes de marcar tarefa como completa
- Template disponível em: `docs/templates/e2e_validation_checklist.md`