# 🔄 Fluxo de Trabalho com Agentes - ATUALIZADO

**Data:** 10/03/2026  
**Status:** ✅ COMPLETO - Pronto para uso em produção

## 🎯 **RESUMO DAS MELHORIAS IMPLEMENTADAS**

### **1. ✅ Validação E2E Integrada**
- **Script `validate_e2e.sh`**: Validação end-to-end com Playwright
- **Modo stealth (headless)**: Não abre janelas no desktop
- **Fluxo de 7 passos padrão**: Validação completa do negócio
- **Integração obrigatória**: No prompt gerado pelo `process_task.sh`

### **2. ✅ Formato Padronizado para Arquivos de Tarefa**
- **Template `task_description.md`**: Estrutura completa para documentação
- **Extração automática de metadados**: `Tipo: | Módulo: | Objetivo: | Decisões:`
- **Arquivos formatados**: `docs/tasks/[nome].md` com descrição detalhada
- **Processamento inteligente**: Script extrai informações automaticamente

### **3. ✅ Fluxo Completo Atualizado**

#### **ANTES:**
```
start_session.sh → process_task.sh "descrição" --execute → [opencode] → conclude_task.sh
```

#### **AGORA:**
```
start_session.sh → process_task.sh --file docs/tasks/tarefa.md --execute → [opencode] → validate_e2e.sh → conclude_task.sh
                                     ↑                                              ↑
                            (arquivo formatado)                          (validação E2E obrigatória)
```

## 📋 **COMO USAR O NOVO FLUXO**

### **Passo 1: Criar arquivo de tarefa formatado**
```bash
# Copiar template
cp docs/templates/task_description.md docs/tasks/minha_tarefa.md

# Editar com descrição completa
# IMPORTANTE: Incluir metadados no formato:
# Tipo: Feature | Módulo: ui_web | Objetivo: Implementar X | Decisões: seguir padrão Y
```

### **Passo 2: Processar tarefa**
```bash
# Usar arquivo formatado
./process_task.sh --file docs/tasks/minha_tarefa.md --execute

# O script irá:
# 1. Extrair metadados automaticamente
# 2. Gerar checklist pré-implementação
# 3. Criar plano de implementação
# 4. Gerar prompt para opencode com validação E2E
```

### **Passo 3: Implementar no opencode**
- Copiar prompt gerado (`.opencode_task_*.txt`)
- Colar no opencode
- Seguir instruções (inclui validação E2E obrigatória)

### **Passo 4: Validar E2E (OBRIGATÓRIO)**
```bash
# Após implementação, executar:
./scripts/dev/validate_e2e.sh --basic --headless

# Modos disponíveis:
--basic      # 7 passos padrão (recomendado)
--headless   # Stealth mode (padrão - não abre janelas)
--ui         # Com navegador visível (debug)
--full       # Todos os testes
```

### **Passo 5: Concluir (apenas se E2E passar)**
```bash
./conclude_task.sh "Aprendizados + resultado E2E: passou" --success
```

## 🧪 **EXEMPLO PRÁTICO COMPLETO**

### **1. Arquivo de tarefa:** `docs/tasks/fornecedores.md`
```markdown
# Implementar UI para Fornecedores

Tipo: Feature | Módulo: ui_web | Objetivo: Implementar gestão completa de fornecedores | Decisões: seguir MemberHandler, CNPJ opcional

[descrição detalhada...]
```

### **2. Processamento:**
```bash
./process_task.sh --file docs/tasks/fornecedores.md --execute
# Gera: checklist, plano, prompt com validação E2E
```

### **3. Prompt gerado inclui:**
```bash
## 🧪 VALIDAÇÃO E2E OBRIGATÓRIA
Após implementar, execute:
./scripts/dev/validate_e2e.sh --basic --headless
- ✅ Se passar: documentar
- ❌ Se falhar: CORRIGIR antes de marcar como completo
```

### **4. Validação após implementação:**
```bash
# Modo stealth (não abre janelas)
./scripts/dev/validate_e2e.sh --basic --headless

# Se falhar, corrigir e validar novamente
# Só concluir quando passar
```

## 🎯 **BENEFÍCIOS DO NOVO SISTEMA**

### **Para qualidade:**
1. **Validação real do negócio**: Testa fluxo completo, não apenas código
2. **Critério de aceite claro**: E2E deve passar antes de "completo"
3. **Detecção antecipada**: Problemas de integração identificados cedo

### **Para documentação:**
1. **Arquivos padronizados**: Estrutura consistente para todas as tarefas
2. **Metadados extraíveis**: Processamento automático pelo script
3. **Histórico completo**: Descrição detalhada + checklists + planos

### **Para experiência do usuário:**
1. **Modo stealth**: Validação não interfere com trabalho no desktop
2. **Relatórios claros**: Cores, detalhes, próximos passos
3. **Feedback imediato**: Resultado da validação em segundos

### **Para o agente:**
1. **Instruções completas**: Prompt inclui tudo necessário
2. **Expectativas claras**: Sabe que validação E2E será executada
3. **Critérios bem definidos**: O que significa "tarefa completa"

## 📊 **CRITÉRIOS DE SUCESSO ATUALIZADOS**

**Uma tarefa só deve ser marcada como "testada end-to-end" quando:**

1. ✅ **Testes unitários** passam (>90% cobertura)
2. ✅ **Smoke test HTTP** passa
3. ✅ **Validação E2E** passa (`validate_e2e.sh --basic --headless`)
4. ✅ **Aprendizados** documentados no `conclude_task.sh`

**Este critério garante que "completo" significa realmente funcionando no fluxo real do negócio.**

## 🚀 **PRÓXIMOS PASSOS RECOMENDADOS**

### **Imediato (próxima sessão):**
1. Criar tarefa real usando template formatado
2. Testar fluxo completo end-to-end
3. Documentar feedback e ajustar se necessário

### **Curto prazo:**
1. Adicionar `data-testid` aos componentes para seletores mais robustos
2. Expandir cobertura de testes E2E
3. Criar dashboard de métricas de qualidade

### **Longo prazo:**
1. Integração com CI/CD automática
2. Notificações de falha E2E
3. Relatórios de tendência de qualidade

## 📞 **SUPORTE E AJUDA**

### **Comandos de referência:**
```bash
# Novo fluxo completo
./process_task.sh --file docs/tasks/tarefa.md --execute
./scripts/dev/validate_e2e.sh --basic --headless
./conclude_task.sh "aprendizados" --success

# Ajuda
./process_task.sh --help
./scripts/dev/validate_e2e.sh --help
```

### **Templates disponíveis:**
- `docs/templates/task_description.md` - Descrição de tarefa
- `docs/templates/e2e_validation_checklist.md` - Checklist E2E
- `docs/templates/pre_implementation_checklist.md` - Checklist pré-impl
- `docs/templates/implementation_plan.md` - Plano de implementação

### **Documentação atualizada:**
- `docs/AGENT_WORKFLOW_GUIDE.md` - Guia completo do fluxo
- `docs/07_testing/01_test_strategy.md` - Estratégia de testes
- `docs/README_OPENCODE_WORKFLOW.md` - Visão geral do sistema

---

**🎉 PRONTO PARA COMEÇAR?**
```bash
# 1. Explore os templates
ls docs/templates/

# 2. Crie sua primeira tarefa formatada
cp docs/templates/task_description.md docs/tasks/minha_primeira_tarefa.md

# 3. Siga o fluxo completo
./process_task.sh --file docs/tasks/minha_primeira_tarefa.md --execute
```

**💡 Dica:** Comece com uma tarefa pequena para se familiarizar, depois escale para features complexas. O sistema aprende com cada tarefa concluída!