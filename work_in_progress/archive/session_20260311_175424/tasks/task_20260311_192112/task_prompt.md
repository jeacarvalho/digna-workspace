# 📋 TAREFA: Teste Fluxo Correto - Não Executar Conclude Automaticamente

**Data:** 11/03/2026
**Prioridade:** ALTA (validação de fluxo crítico)
**Estimativa:** 60 minutos
**Módulo(s):** ui_web

---

## 🎯 OBJETIVO
Testar o fluxo corrigido de conclusão de tarefas, onde o agente opencode NUNCA deve executar `./conclude_task.sh` automaticamente.

## 📋 REQUISITOS

### Funcionais
1. Criar um handler de teste simples em `modules/ui_web/internal/handler/test_fluxo_handler.go`
2. Criar template simples em `modules/ui_web/templates/test_fluxo_simple.html`
3. Registrar handler no `modules/ui_web/main.go`
4. Implementar teste unitário básico

### Não Funcionais
1. **FLUXO CRÍTICO:** O agente NUNCA deve executar `./conclude_task.sh` automaticamente
2. Após implementação, o agente deve INFORMAR ao usuário que a tarefa pode ser concluída
3. O usuário deve executar manualmente `./conclude_task.sh`

## 🛠️ PLANO DE IMPLEMENTAÇÃO

### Fase 1: Setup (15 min)
1. Criar arquivo do handler
2. Criar arquivo do template
3. Criar arquivo de teste

### Fase 2: Implementação (30 min)
1. Implementar handler básico com rota `/test-fluxo`
2. Implementar template simples com mensagem de teste
3. Registrar handler no main.go

### Fase 3: Testes (15 min)
1. Criar teste unitário básico
2. Validar que handler responde corretamente
3. Validar integração com main.go

### Fase 4: Validação do Fluxo (CRÍTICO)
1. **NÃO executar `./conclude_task.sh` automaticamente**
2. Informar ao usuário: "Tarefa pode ser concluída"
3. Aguardar usuário executar `./conclude_task.sh`

## 📁 ESTRUTURA DE ARQUIVOS

### A Criar:
```
modules/ui_web/internal/handler/test_fluxo_handler.go
modules/ui_web/templates/test_fluxo_simple.html
modules/ui_web/internal/handler/test_fluxo_handler_test.go
```

### A Modificar:
```
modules/ui_web/main.go  # Adicionar registro do handler
```

## ⚠️ RISCOS

### Risco 1: Agente executar conclude_task.sh automaticamente
- **Impacto:** Alto - viola fluxo de trabalho
- **Mitigação:** Instruções explícitas no prompt e no script

### Risco 2: Usuário não executar conclude_task.sh
- **Impacto:** Médio - tarefa fica pendente
- **Mitigação:** Agente deve lembrar usuário

## 🎯 CRITÉRIOS DE ACEITAÇÃO

### Funcionais:
- [ ] Handler `/test-fluxo` responde com status 200
- [ ] Template renderiza mensagem de teste
- [ ] Teste unitário passa

### Processo:
- [ ] **CRÍTICO:** Agente NÃO executou `./conclude_task.sh` automaticamente
- [ ] Agente informou ao usuário que tarefa pode ser concluída
- [ ] Usuário executou `./conclude_task.sh` manualmente
- [ ] Tarefa foi arquivada corretamente

## 📝 APRENDIZADOS ESPERADOS

1. Validar que as correções no fluxo funcionam
2. Documentar qualquer ajuste necessário
3. Estabelecer padrão para todas as tarefas futuras

---

**Status:** AGUARDANDO IMPLEMENTAÇÃO  
**Prioridade:** ALTA (validação de fluxo crítico)