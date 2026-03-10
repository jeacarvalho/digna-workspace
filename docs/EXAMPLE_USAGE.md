# 🎯 Exemplo Rápido de Uso do Sistema

## **Cenário:** Você quer implementar a UI para "Gestão de Fornecedores"

### **Passo 1: Iniciar Sessão (30 segundos)**
```bash
./start_session.sh quick
```
**Saída esperada:**
```
🚀 Iniciando sessão no projeto Digna...
⚡ Modo rápido: verificando status básico...
📊 Verificando status dos testes...
PASS
ok  	github.com/providentia/digna/ui_web	0.012s
✅ Sessão iniciada com sucesso!
```

### **Passo 2: Processar Tarefa (1 minuto)**
```bash
./process_task.sh "Tipo: Feature | Módulo: ui_web | Objetivo: Implementar UI para Fornecedores | Decisões: cards, CNPJ opcional, seguir padrão MemberHandler" --execute
```

**O script vai:**
1. Extrair: Tipo=Feature, Módulo=ui_web, Objetivo=Implementar UI para Fornecedores
2. Gerar checklist: `docs/implementation_plans/fornecedores_pre_check.md`
3. Gerar plano: `docs/implementation_plans/fornecedores_implementation_20250310.md`
4. Criar prompt: `.opencode_task_20250310_101500.txt`

**Saída final:**
```
🎯 AGORA COPIE E COLE O CONTEÚDO ACIMA NO OPENCODE
   ou use: cat .opencode_task_20250310_101500.txt | pbcopy (Mac)
```

### **Passo 3: Implementar no OpenCode (15-30 minutos)**
1. **Copie** o prompt gerado
2. **Cole** no opencode
3. **Siga** as instruções do prompt (que incluem:
   - Ler o checklist
   - Seguir o plano
   - Usar padrão MemberHandler
   - Implementar com TDD)

### **Passo 4: Concluir Tarefa (2 minutos)**
```bash
./conclude_task.sh "Aprendizados: checklist antecipou problema de acesso internal, padrão MemberHandler fácil de seguir, testes cobriram 92%" --success
```

**O script vai:**
1. Coletar métricas (testes, tempo, código)
2. Criar: `docs/learnings/20250310_101500_fornecedores_learnings.md`
3. Atualizar checklists com novos itens
4. Preparar próximos passos: `docs/NEXT_STEPS.md`

**Saída final:**
```
✅ CONCLUSÃO DA TAREFA COMPLETA!
📊 RESUMO:
Tarefa: fornecedores (20250310_101500)
Status: success
Duração: 25 minutos
Testes: 12 passed, 0 failed
```

---

## 📁 **O que foi Gerado Automaticamente:**

### **1. Checklist Pré-Implementação**
`docs/implementation_plans/fornecedores_pre_check.md`
- 60+ itens de verificação
- Análise de backend existente
- Padrões a seguir
- Riscos identificados

### **2. Plano de Implementação**  
`docs/implementation_plans/fornecedores_implementation_20250310.md`
- Tarefas detalhadas
- Critérios de aceite
- Cronograma estimado
- Código de referência

### **3. Documento de Aprendizados**
`docs/learnings/20250310_101500_fornecedores_learnings.md`
- Métricas da implementação
- O que funcionou bem
- Problemas encontrados
- Melhorias para próxima

### **4. Próximos Passos**
`docs/NEXT_STEPS.md`
- Continuação desta tarefa (se parcial)
- Tarefas relacionadas
- Decisões pendentes

---

## ⚡ **Fluxo Ultra-Rápido (Para Tarefas Simples):**

```bash
# 1. Iniciar (quick)
./start_session.sh quick

# 2. Tarefa direta
./process_task.sh "Corrigir bug: campo data não formata corretamente no PDV" --execute

# 3. [Implementar no opencode]

# 4. Concluir
./conclude_task.sh "Aprendizados: bug era na função formatDate, corrigido em 10min" --success
```

**Tempo total:** ~15 minutos (vs 1-2 horas sem sistema)

---

## 🎨 **Vantagens Imediatas:**

### **✅ Antes do Sistema:**
- "Vou descobrir durante a implementação"
- Mesmos problemas repetidos
- Documentação inconsistente
- Tempo de warm-up: 5-10 minutos

### **✅ Com o Sistema:**
- Problemas antecipados no checklist
- Aprendizados acumulados e reutilizados
- Documentação automática e consistente
- Tempo de warm-up: < 2 minutos

### **✅ Para o OpenCode:**
- Prompts otimizados e completos
- Contexto sempre atualizado
- Instruções passo a passo
- Feedback loop de melhorias

---

## 🚀 **Pronto para Começar?**

**Comando inicial:**
```bash
./start_session.sh
```

**Primeira tarefa de teste:**
```bash
./process_task.sh "Tipo: Investigação | Módulo: ui_web | Objetivo: Explorar sistema de navegação" --checklist
```

**Dica:** Comece com uma tarefa pequena para se familiarizar, depois escale!

---

## 📞 **Precisa de Ajuda?**

```bash
# Ajuda dos scripts
./start_session.sh --help
./process_task.sh --help  
./conclude_task.sh --help

# Documentação completa
cat README_OPENCODE_WORKFLOW.md | head -50
```

**Problemas comuns:**
- Scripts não executáveis? `chmod +x *.sh`
- Arquivos não encontrados? Verifique se está no diretório correto
- Erros de permissão? Use `bash ./script.sh` em vez de `./script.sh`