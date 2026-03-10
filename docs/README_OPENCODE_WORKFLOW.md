# 🚀 Sistema de Workflow OpenCode - Projeto Digna

**Objetivo:** Automatizar todo o processo de desenvolvimento no opencode, desde a inicialização até a documentação de aprendizados.

## 📋 **Visão Geral do Sistema**

### **3 Scripts Principais:**
1. `start_session.sh` - Inicializa sessão com contexto atualizado
2. `process_task.sh` - Processa tarefas (checklists, planos, execução)
3. `conclude_task.sh` - Conclui tarefas e documenta aprendizados

### **Fluxo Completo (5-10 minutos):**
```
./start_session.sh → ./process_task.sh "tarefa" --execute → [opencode] → ./conclude_task.sh "aprendizados"
```

---

## 🎯 **Como Usar (Passo a Passo)**

### **Passo 1: Iniciar Sessão**
```bash
# Modo completo (recomendado primeira vez)
./start_session.sh

# Modo rápido (sessões subsequentes)
./start_session.sh quick
```

**Resultado:** Contexto atualizado + status do projeto + próximos passos.

### **Passo 2: Processar Tarefa**
```bash
# Formato recomendado:
./process_task.sh "Tipo: Feature | Módulo: ui_web | Objetivo: Implementar X | Decisões: seguir padrão Y"

# Opções:
./process_task.sh "descrição" --checklist    # Apenas checklist
./process_task.sh "descrição" --plan         # Checklist + plano
./process_task.sh "descrição" --execute      # Tudo + prompt para opencode
```

### **Passo 3: Implementar no OpenCode**
1. O script `--execute` gera um prompt completo
2. **Copie e cole o prompt no opencode**
3. Siga as instruções geradas

### **Passo 4: Validação Pós-Implementação** ⭐ **NOVO**
```bash
# Smoke test OBRIGATÓRIO para novas features
./scripts/smoke_test_new_feature.sh "Nome da Feature" "/rota_principal"

# Exemplo para Members:
./scripts/smoke_test_new_feature.sh "Member Management" "/members"
```

**Validação 4 níveis:**
1. ✅ Testes unitários
2. ✅ Testes de integração  
3. ✅ Testes de sistema (`go test -run TestSystem`)
4. ✅ Smoke test local (script acima)

### **Passo 5: Concluir Tarefa**
```bash
# Documentar aprendizados
./conclude_task.sh "O que aprendi + resultado do smoke test"
```

**Resultado:** Contexto atualizado + status do projeto + próximos passos.

### **Passo 2: Processar Tarefa**
```bash
# Formato recomendado:
./process_task.sh "Tipo: Feature | Módulo: ui_web | Objetivo: Implementar X | Decisões: seguir padrão Y"

# Opções:
./process_task.sh "descrição" --checklist    # Apenas checklist
./process_task.sh "descrição" --plan         # Checklist + plano
./process_task.sh "descrição" --execute      # Tudo + prompt para opencode
```

### **Passo 3: Implementar no OpenCode**
1. O script `--execute` gera um prompt completo
2. **Copie e cole o prompt no opencode**
3. Siga as instruções geradas

### **Passo 4: Concluir Tarefa**
```bash
# Documentar aprendizados
./conclude_task.sh "Aprendizados: checklist antecipou 3 problemas, testes cobriram 90%"

# Opções de status:
./conclude_task.sh "aprendizados" --success   # Concluído (padrão)
./conclude_task.sh "aprendizados" --partial   # Parcialmente concluído
./conclude_task.sh "aprendizados" --failed    # Falhou
```

---

## 📝 **Formatos de Tarefa Recomendados**

### **Para Novas Features:**
```
"Tipo: Feature | Módulo: ui_web | Objetivo: Implementar UI para [Nome] | Decisões: cards, seguir padrão MemberHandler"
```

### **Para Bugs:**
```
"Tipo: Bug | Módulo: ui_web | Objetivo: Corrigir erro no PDV ao [descrição] | Decisões: manter anti-float, testar edge cases"
```

### **Para Melhorias:**
```
"Tipo: Melhoria | Módulo: ui_web | Objetivo: Adicionar paginação no [feature] | Decisões: HTMX infinite scroll, performance"
```

### **Para Investigação:**
```
"Tipo: Investigação | Módulo: core_lume | Objetivo: Entender como [sistema] funciona | Decisões: mapear fluxo, documentar"
```

---

## 🏗️ **Estrutura de Arquivos Gerada**

### **Por Tarefa:**
```
docs/implementation_plans/
├── [feature]_pre_check.md          # Checklist pré-implementação
├── [feature]_implementation_[data].md  # Plano completo
└── (referenciado em)

docs/learnings/
└── [task_id]_[feature]_learnings.md    # Aprendizados documentados

.opencode_task_[id].txt                 # Prompt para opencode (temporário)
.task_[id]                              # Metadados da tarefa (temporário)
```

### **Documentação do Sistema:**
```
docs/templates/
├── pre_implementation_checklist.md     # Template de checklist
└── implementation_plan.md              # Template de plano

docs/antipatterns/
└── common_antipatterns_solutions.md    # Antipadrões e soluções

docs/NEXT_STEPS.md                      # Próximos passos acumulados
docs/QUICK_REFERENCE.md                 # Referência rápida do projeto
```

---

## 🔧 **Funcionalidades dos Scripts**

### **`start_session.sh`**
- ✅ Atualiza contexto do projeto
- ✅ Verifica status dos testes
- ✅ Mostra referência rápida
- ✅ Sugere próximos passos
- ✅ Modo rápido/completo

### **`process_task.sh`**
- ✅ Extrai informações da descrição (Tipo, Módulo, Objetivo, Decisões)
- ✅ Gera checklist pré-implementação personalizado
- ✅ Cria plano de implementação completo
- ✅ Prepara prompt otimizado para opencode
- ✅ Segue padrões estabelecidos (MemberHandler, etc.)

### **`conclude_task.sh`**
- ✅ Coleta métricas automáticas (testes, código, tempo)
- ✅ Documenta aprendizados estruturados
- ✅ Atualiza checklists com novos itens
- ✅ Atualiza antipadrões com problemas encontrados
- ✅ Prepara próximos passos para próxima sessão

---

## 🎨 **Exemplo Completo de Uso**

### **Sessão de 30 minutos:**
```bash
# 1. Iniciar (2 min)
./start_session.sh quick

# 2. Processar tarefa (3 min)
./process_task.sh "Tipo: Feature | Módulo: ui_web | Objetivo: Implementar UI para Fornecedores | Decisões: cards, CNPJ opcional, seguir MemberHandler" --execute

# 3. [No opencode] Implementar seguindo prompt (20 min)

# 4. Concluir (5 min)
./conclude_task.sh "Aprendizados: checklist útil, padrão MemberHandler fácil de seguir, testes cobriram 95%" --success
```

### **Arquivos gerados:**
1. `docs/implementation_plans/fornecedores_pre_check.md`
2. `docs/implementation_plans/fornecedores_implementation_20250310.md`
3. `docs/learnings/20250310_101500_fornecedores_learnings.md`
4. `docs/NEXT_STEPS.md` atualizado

---

## 📊 **Benefícios do Sistema**

### **Para o Desenvolvedor:**
- ✅ **Menos "descobertas durante implementação"** (checklists antecipam problemas)
- ✅ **Consistência** (segue padrões estabelecidos)
- ✅ **Documentação automática** (aprendizados não se perdem)
- ✅ **Onboarding rápido** (novas features seguem mesmo processo)

### **Para o Projeto:**
- ✅ **Conhecimento acumulado** (checklists melhoram a cada tarefa)
- ✅ **Qualidade consistente** (validações automáticas)
- ✅ **Métricas reais** (tempo, problemas, cobertura)
- ✅ **Processo reproduzível** (qualquer dev pode seguir)

### **Para o OpenCode:**
- ✅ **Prompts otimizados** (contexto completo + instruções claras)
- ✅ **Menor tempo de warm-up** (contexto já atualizado)
- ✅ **Resultados previsíveis** (segue planos estruturados)
- ✅ **Feedback loop** (aprendizados melhoram próximas interações)

---

## 🚨 **Solução de Problemas Comuns**

### **"Script não encontrado"**
```bash
# Tornar executáveis
chmod +x start_session.sh process_task.sh conclude_task.sh

# Se scripts não existirem, recrie:
curl -O https://raw.githubusercontent.com/seu-repo/digna/main/start_session.sh
```

### **"Checklist não preenchido"**
```bash
# Gerar apenas checklist primeiro
./process_task.sh "descrição" --checklist

# Preencher manualmente o arquivo gerado
# Depois gerar plano
./process_task.sh "descrição" --plan
```

### **"Prompt muito longo para opencode"**
```bash
# Usar modo --plan primeiro, revisar, depois --execute
# Ou dividir tarefa em subtarefas menores
```

### **"Esqueci de documentar aprendizados"**
```bash
# Sempre rode conclude_task.sh, mesmo para tarefas parciais
./conclude_task.sh "Aprendizados: parei na etapa X, problema Y" --partial
```

---

## 🔄 **Melhoria Contínua do Sistema**

### **Após cada tarefa:**
1. Revisar `docs/learnings/[tarefa]_learnings.md`
2. Atualizar checklists com itens faltantes
3. Adicionar antipadrões novos encontrados
4. Melhorar templates baseado no feedback

### **Métricas a acompanhar:**
- **Tempo médio por tarefa** (meta: < 60 minutos)
- **Problemas antecipados vs reais** (meta: > 80%)
- **Cobertura de testes** (meta: > 90%)
- **Satisfação com o sistema** (1-5 no conclude_task.sh)

---

## 🎯 **Primeira Vez Usando o Sistema?**

### **Fluxo Inicial Recomendado:**
```bash
# 1. Explore o sistema
./start_session.sh
./process_task.sh --help
./conclude_task.sh --help

# 2. Tarefa de teste simples
./process_task.sh "Tipo: Investigação | Módulo: ui_web | Objetivo: Entender sistema de templates" --checklist

# 3. Preencher checklist manualmente
# 4. Gerar plano
./process_task.sh "mesma descrição" --plan

# 5. Concluir (mesmo sem implementar)
./conclude_task.sh "Aprendizados: sistema fácil de usar, checklists completos" --success
```

### **Próximos Passos Imediatos:**
1. Escolher uma feature real do backlog
2. Usar `--execute` e implementar no opencode
3. Documentar aprendizados reais
4. Melhorar sistema baseado na experiência

---

## 📞 **Suporte e Feedback**

### **Problemas com os scripts:**
```bash
# Verificar logs
tail -f .session_* .task_* 2>/dev/null

# Debug mode
bash -x ./process_task.sh "teste" --checklist
```

### **Sugestões de melhoria:**
1. Documentar em `docs/learnings/[data]_system_feedback.md`
2. Atualizar templates em `docs/templates/`
3. Compartilhar com a equipe

### **Status do Sistema:**
- ✅ Scripts criados e testados
- ✅ Documentação completa
- ✅ Pronto para uso em produção
- ✅ Melhoria contínua habilitada

---

**🎉 Pronto para começar?**
```bash
./start_session.sh
```

**💡 Dica:** Comece com uma tarefa pequena para se familiarizar, depois escale para features complexas. O sistema aprende com cada tarefa concluída!