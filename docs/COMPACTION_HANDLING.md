# 🚨 MANUAL: Como Lidar com Compaction do Opencode

## 🎯 O Problema do Compaction

### ❌ O que é Compaction?
Compaction é um mecanismo interno do opencode que:
1. **É IMPREVISÍVEL** - acontece quando o sistema decide economizar tokens
2. **PERDE CONTEXTO** - o agente esquece tudo sobre a tarefa atual
3. **REQUER INTERVENÇÃO** - usuário precisa ajudar o agente a recuperar

### ⚠️ Sintomas de Compaction:
- Agente pergunta: "O que estamos fazendo?"
- Agente esqueceu a tarefa atual
- Agente parece "travado" ou confuso
- Respostas são genéricas, sem contexto específico

## 🔄 Fluxo Correto para Lidar com Compaction

### 📋 FLUXO COMPLETO (EXECUTADO PELO USUÁRIO):

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Detecção do    │    │  Durante        │    │  Após           │
│  Compaction     │    │  Compaction     │    │  Compaction     │
├─────────────────┤    ├─────────────────┤    ├─────────────────┤
│ 1. Usuário      │    │ 1. opencode     │    │ 1. Usuário      │
│    detecta      │    │    trabalha     │    │    executa:     │
│    compaction   │    │    sem contexto │    │    ./scripts/   │
│                 │    │ 2. Perde        │    │    help_agent_  │
│ 2. Usuário      │    │    progresso    │    │    recover_     │
│    executa:     │    │                 │    │    context.sh   │
│    ./preserve_  │    │                 │    │                 │
│    context.sh   │    │                 │    │ 2. opencode lê  │
│    --save       │    │                 │    │    contexto     │
│                 │    │                 │    │ 3. Continua de  │
│ 3. Salva        │    │                 │    │    onde parou   │
│    contexto     │    │                 │    │                 │
│    em arquivo   │    │                 │    │ 4. Usuário      │
│                 │    │                 │    │    executa      │
│                 │    │                 │    │    ./conclude_  │
│                 │    │                 │    │    task.sh      │
│                 │    │                 │    │    quando       │
│                 │    │                 │    │    terminar     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🛠️ Scripts Disponíveis

### 1. **`./scripts/help_agent_recover_context.sh`**
**Quando usar:** Sempre que o agente parecer perdido após compaction
**O que faz:** 
- Verifica se há contexto salvo
- Lista tarefas ativas
- Sugere próxima ação
- Mostra comandos úteis

```bash
./scripts/help_agent_recover_context.sh
```

### 2. **`./preserve_context.sh --save`**
**Quando usar:** Quando detectar que o agente entrou em compaction
**O que faz:** Salva contexto da tarefa atual em arquivo
**IMPORTANTE:** Executado pelo USUÁRIO, não pelo agente

```bash
./preserve_context.sh --save
```

### 3. **`./preserve_context.sh --restore`**
**Quando usar:** Após o compaction terminar
**O que faz:** Restaura contexto salvo para o agente
**IMPORTANTE:** Executado pelo USUÁRIO

```bash
./preserve_context.sh --restore
```

### 4. **`./preserve_context.sh --status`**
**Quando usar:** Para verificar se há contexto salvo
**O que faz:** Mostra status do contexto preservado

```bash
./preserve_context.sh --status
```

## 🎯 Cenários Comuns e Soluções

### Cenário 1: Agente pergunta "O que estamos fazendo?"
```
Sintoma: Agente esqueceu completamente a tarefa
Solução:
  1. ./scripts/help_agent_recover_context.sh
  2. Identificar tarefa ativa
  3. ./process_task.sh --task=ID --execute (se necessário)
  4. Guiar agente de volta
```

### Cenário 2: Agente está "travado" sem progresso
```
Sintoma: Respostas genéricas, sem avanço na implementação
Solução:
  1. ./preserve_context.sh --save (se possível)
  2. Aguardar alguns minutos
  3. ./preserve_context.sh --restore
  4. ./scripts/help_agent_recover_context.sh
```

### Cenário 3: Contexto foi perdido sem salvamento
```
Sintoma: Nenhum arquivo de contexto salvo
Solução:
  1. ./scripts/help_agent_recover_context.sh
  2. Identificar tarefa mais recente
  3. ./process_task.sh --task=ID --execute
  4. Reexplicar contexto se necessário
```

## 📋 Checklist para Usuário

### ✅ Antes de Começar Qualquer Tarefa:
- [ ] Verificar `work_in_progress/current_session/.agent_context.md`
- [ ] Executar `./preserve_context.sh --clean` para limpar contexto antigo
- [ ] Confirmar que agente leu documentação obrigatória

### ✅ Durante Implementação (Monitoramento):
- [ ] Observar se agente mantém contexto da tarefa
- [ ] Estar atento a sinais de compaction
- [ ] Ter `./preserve_context.sh --save` pronto para executar

### ✅ Após Detectar Compaction:
- [ ] Executar `./preserve_context.sh --save` imediatamente
- [ ] Aguardar compaction terminar
- [ ] Executar `./scripts/help_agent_recover_context.sh`
- [ ] Executar `./preserve_context.sh --restore`
- [ ] Verificar que agente recuperou contexto

### ✅ Para Prevenir Problemas:
- [ ] Tarefas curtas (1-2 horas máximo)
- [ ] Checkpoints frequentes (salvar progresso)
- [ ] Documentar aprendizados após cada sessão
- [ ] Usar `./conclude_task.sh` apenas quando usuário validar

## ⚠️ Limitações do Sistema

### ❌ O que NÃO é possível:
1. **Prever compaction** - Agente não sabe quando vai acontecer
2. **Executar `--save` automaticamente** - Agente não pode prever
3. **Prevenir perda total** - Sem salvamento, contexto é perdido

### ✅ O que É possível:
1. **Detecção pelo usuário** - Usuário pode perceber compaction
2. **Salvamento manual** - Usuário executa `--save` quando detecta
3. **Recuperação guiada** - Scripts ajudam a recuperar contexto

## 🚀 Melhores Práticas

### 1. **Tarefas Pequenas e Focadas**
- Divida tarefas grandes em subtarefas
- Conclua uma tarefa antes de começar outra
- Use `./conclude_task.sh` após cada implementação

### 2. **Monitoramento Ativo**
- Fique atento a mudanças no comportamento do agente
- Pergunte periodicamente: "Em que ponto estamos?"
- Verifique se agente mantém contexto

### 3. **Documentação Contínua**
- Documente aprendizados em `task_learnings.md`
- Atualize checklists com problemas encontrados
- Compartilhe soluções para compaction

### 4. **Fluxo Disciplinado**
- Sempre siga: create → process → implement → conclude
- Nunca pule etapas
- Sempre valide antes de concluir

## 📚 Comandos de Referência Rápida

```bash
# Quando agente está perdido:
./scripts/help_agent_recover_context.sh

# Quando detectar compaction:
./preserve_context.sh --save

# Após compaction terminar:
./preserve_context.sh --restore

# Para ver status:
./preserve_context.sh --status

# Para limpar:
./preserve_context.sh --clean

# Fluxo normal de trabalho:
./create_task.sh "Nome da Tarefa"
# Editar task_prompt.md
./process_task.sh --task=ID --execute
# Agente implementa
./conclude_task.sh --task=ID "Aprendizados: ..."
```

## 🔄 Integração com Fluxo Existente

O sistema de handling de compaction se integra com:

1. **`create_task.sh`** - Identifica tarefa ativa
2. **`process_task.sh`** - Processa tarefa após recuperação
3. **`conclude_task.sh`** - Conclui apenas quando usuário valida
4. **`start_session.sh`** - Inicia sessão com contexto limpo
5. **`end_session.sh`** - Limpa todos os arquivos temporários

## 🎯 Conclusão

**Compaction é inevitável, mas gerenciável.**  
Com os scripts e fluxos corretos, o usuário pode:

1. **Detectar** quando compaction acontece
2. **Preservar** contexto antes que seja perdido
3. **Recuperar** contexto após compaction
4. **Continuar** implementação de onde parou

**Lembrete crítico:** O agente opencode NUNCA pode prever compaction.  
A responsabilidade de detectar e gerenciar compaction é do USUÁRIO.

---

**📌 Status:** ✅ MANUAL COMPLETO  
**Última atualização:** 11/03/2026  
**Próxima revisão:** Após próxima ocorrência de compaction