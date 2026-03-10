# 🧪 Estratégia de Validação - Projeto Digna

**Última atualização:** 10/03/2026  
**Criado durante:** Correção do bug HTMX no estoque

---

## 🎯 Objetivo

Prevenir conclusão prematura de tarefas garantindo que:
1. **Bugs são realmente corrigidos** (não apenas sintaticamente)
2. **Testes E2E passam** com o fluxo completo
3. **Ambientes são isolados** (sem contaminação de dados)
4. **Regressões são detectadas** antes de marcar como concluído

---

## 📋 Checklist de Validação

### ✅ **ANTES de marcar tarefa como concluída:**

#### 1. **Compilação e Build**
- [ ] Servidor compila sem erros: `go build`
- [ ] Templates carregam sem erros
- [ ] Dependências estão instaladas

#### 2. **Testes Unitários**
- [ ] Testes existentes passam: `go test ./...`
- [ ] Novos testes foram criados para o bug/funcionalidade
- [ ] Cobertura de casos de borda

#### 3. **Testes E2E com Isolamento**
- [ ] Usar `e2e_test_runner.sh` com bancos isolados
- [ ] Testar fluxo completo (login → ação → verificação)
- [ ] Validar persistência de dados
- [ ] Verificar cleanup automático

#### 4. **Validação Manual**
- [ ] Acessar funcionalidade no navegador
- [ ] Testar cenários positivos e negativos
- [ ] Verificar mensagens de erro/sucesso
- [ ] Testar recarregamento da página

#### 5. **Prevenção de Regressão**
- [ ] Executar testes relacionados
- [ ] Verificar funcionalidades adjacentes
- [ ] Documentar casos de teste

---

## 🛠️ Ferramentas de Validação

### 1. **E2E Test Runner** (`scripts/dev/e2e_test_runner.sh`)
```bash
# Executar todos testes
./scripts/dev/e2e_test_runner.sh

# Teste específico
./scripts/dev/e2e_test_runner.sh --spec tests/stock-crud.spec.js

# Com navegador visível (debug)
./scripts/dev/e2e_test_runner.sh --headed

# Manter ambiente (debug)
./scripts/dev/e2e_test_runner.sh --skip-cleanup
```

**Características:**
- ✅ Cria banco de dados isolado (`/tmp/test_*.db`)
- ✅ Configura entidade de teste única
- ✅ Inicia servidor em porta dedicada
- ✅ Executa testes Playwright
- ✅ Limpa tudo automaticamente

### 2. **Testes Playwright** (`tests/`)
- `stock-crud.spec.js` - CRUD de itens de estoque
- `digna-e2e.spec.js` - Fluxos principais do sistema

### 3. **Testes Unitários Go** (`*_test.go`)
- Executar: `go test ./modules/ui_web/...`
- Foco em lógica de negócio e handlers

### 4. **Validação Manual**
```bash
# Verificar servidor
curl http://localhost:8090/health

# Verificar template (com autenticação)
# 1. Login via navegador
# 2. Testar funcionalidade
# 3. Verificar console por erros
```

---

## 🚨 **Sinais de Alerta (NÃO concluir se...)**

### ❌ **Problemas de Compilação**
- Erros de tipo no template
- Handlers não compilando
- Dependências faltando

### ❌ **Testes Falhando**
- E2E tests não passam
- Unit tests quebrando
- Regressões em funcionalidades existentes

### ❌ **Problemas de Isolamento**
- Dados de teste contaminando produção
- Servidores conflitando em portas
- Cleanup não funcionando

### ❌ **Fluxo Incompleto**
- Autenticação não funcionando
- Redirecionamentos incorretos
- Mensagens de erro/sucesso ausentes

---

## 📝 **Documentação Obrigatória**

### 1. **Aprendizados** (`docs/learnings/`)
- Root cause do bug
- Solução implementada
- Lições aprendidas
- Decisões técnicas

### 2. **Próximos Passos** (`docs/NEXT_STEPS.md`)
- Status atual (✅/⚠️/❌)
- Correções implementadas
- Próximas ações
- Decisões pendentes

### 3. **Casos de Teste**
- Cenários testados
- Resultados esperados vs. obtidos
- Dados de teste usados
- Screenshots (se aplicável)

---

## 🔄 **Processo de Validação**

### **Fase 1: Preparação**
```bash
# 1. Verificar ambiente
./scripts/dev/e2e_test_runner.sh --help

# 2. Limpar dados antigos
rm -f /tmp/test_*.db

# 3. Compilar servidor
cd modules/ui_web && go build
```

### **Fase 2: Execução**
```bash
# 1. Executar testes E2E
./scripts/dev/e2e_test_runner.sh --spec tests/RELEVANTE.spec.js

# 2. Executar testes unitários
go test ./modules/ui_web/...

# 3. Teste manual (se necessário)
# - Acessar localhost:8090
# - Testar fluxo manualmente
```

### **Fase 3: Verificação**
```bash
# 1. Verificar logs
tail -f /tmp/digna_test_*.log

# 2. Verificar bancos de teste
ls -la /tmp/test_*.db

# 3. Verificar cleanup
# (não deve haver arquivos após teste completo)
```

### **Fase 4: Documentação**
```bash
# 1. Atualizar NEXT_STEPS.md
# 2. Criar/atualizar aprendizados
# 3. Atualizar casos de teste
```

---

## 🎓 **Lições do Bug HTMX**

### **O que deu errado inicialmente:**
1. ✅ Corrigimos sintaxe, mas não testamos fluxo completo
2. ✅ Testes unitários passaram, mas E2E falharam
3. ✅ Handler retornava HTML completo (duplicação)

### **O que aprendemos:**
1. **HTMX requer 3 coisas**: target ID + rota correta + biblioteca
2. **Handlers HTMX** devem retornar apenas fragmentos
3. **Testes E2E** devem incluir autenticação
4. **Isolamento** é crítico para testes confiáveis

### **Solução implementada:**
1. ✅ Runner com bancos isolados
2. ✅ Testes abrangentes de CRUD
3. ✅ Cleanup automático
4. ✅ Validação de fluxo completo

---

## 📞 **Suporte**

### **Problemas comuns:**
1. **Autenticação falhando**: Verificar credenciais de teste
2. **Porta em uso**: Matar processos antigos `pkill -f "go run"`
3. **Banco bloqueado**: Remover arquivos `/tmp/test_*.db`

### **Debug:**
```bash
# Logs do servidor
tail -f /tmp/digna_server.log

# Logs do teste
tail -f /tmp/digna_test_*.log

# Banco de teste
sqlite3 /tmp/test_*.db ".tables"
```

---

**Nota:** Esta estratégia foi criada após identificar que a correção inicial do bug HTMX estava incompleta. A validação adequada teria detectado que os testes E2E falhavam devido a problemas de autenticação e fluxo incompleto.