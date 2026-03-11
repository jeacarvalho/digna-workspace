# 🚀 Nova Estrutura de Workflow - Projeto Digna

**Data de implementação:** 11/03/2026  
**Status:** ✅ IMPLEMENTADA  
**Objetivo:** Organizar arquivos temporários e padronizar fluxo de trabalho

---

## 📁 ESTRUTURA DE DIRETÓRIOS

```
work_in_progress/
├── current_session/           # ✅ Sessão atual (apenas uma por vez)
│   ├── .agent_context.md     # Contexto atualizado do agente
│   ├── session_info          # Metadados da sessão (SESSION_ID, START_TIME, etc.)
│   └── session_learnings/    # Aprendizados coletados das tarefas
│
├── tasks/                    # ✅ Tarefas em andamento
│   ├── task_20250311_101108/
│   │   ├── task_prompt.md    # Prompt padronizado da tarefa
│   │   ├── checklist.md      # Checklist gerado (opcional)
│   │   ├── implementation_plan.md
│   │   ├── task_learnings.md
│   │   ├── task_metadata     # Metadados (TASK_ID, TASK_NAME, STATUS, etc.)
│   │   └── task_artifacts/   # Arquivos temporários da task
│   │
│   └── task_20250311_112301/
│       └── ... (mesma estrutura)
│
├── archive/                  # ✅ Sessões e tarefas concluídas
│   ├── session_20250311_100813/
│   │   ├── session_info
│   │   ├── session_learnings_consolidated.md
│   │   └── tasks/
│   │       ├── task_20250311_101108/
│   │       └── task_20250311_112301/
│   │
│   └── session_20250311_112108/
│       └── ... (mesma estrutura)
│
└── task_template/            # ✅ Template para novas tarefas
    └── task_prompt.md        # Template padronizado
```

---

## 🔄 NOVO FLUXO DE TRABALHO

### 📋 SEQUÊNCIA COMPLETA:

```bash
# 1. INICIAR SESSÃO (uma vez por sessão)
./start_session.sh [quick]          # "quick" para modo rápido

# 2. CRIAR TAREFA (para cada tarefa)
./create_task.sh "Nome da Tarefa" [módulo]

# 3. EDITAR PROMPT (opcional, mas recomendado)
vim work_in_progress/tasks/task_[ID]/task_prompt.md

# 4. PROCESSAR TAREFA (fluxo recomendado)
./process_task.sh --task=[ID] --checklist    # Primeiro: checklist
./process_task.sh --task=[ID] --plan         # Depois: plano
./process_task.sh --task=[ID] --execute      # Final: executar

# 5. IMPLEMENTAR (opencode faz isso)
#    - Seguir padrões identificados
#    - Usar analyze_patterns.sh para referências
#    - Validar com smoke tests

# 6. CONCLUIR TAREFA (após cada implementação)
./conclude_task.sh --task=[ID] "Aprendizados: item1, item2" --success

# 7. ENCERRAR SESSÃO (após TODAS tarefas da sessão)
./end_session.sh [force]                     # "force" para forçar com tarefas pendentes
```

### 🎯 FLUXO VISUAL:
```
┌─────────────────┐
│  start_session  │
└────────┬────────┘
         │
┌────────▼────────┐
│   create_task   │
└────────┬────────┘
         │
┌────────▼────────┐    ┌──────────────┐    ┌──────────────┐
│ process_task    │───▶│   checklist  │───▶│     plan     │
│   (--checklist) │    └──────────────┘    └────────┬─────┘
└────────┬────────┘                                  │
         │                                  ┌────────▼─────┐
┌────────▼────────┐                         │   execute    │
│ process_task    │◀────────────────────────│  (opencode)  │
│   (--execute)   │                         └────────┬─────┘
└────────┬────────┘                                  │
         │                                  ┌────────▼─────┐
┌────────▼────────┐                         │  conclude    │
│ conclude_task   │◀────────────────────────│    task      │
└────────┬────────┘                         └──────────────┘
         │
┌────────▼────────┐
│  end_session    │
└─────────────────┘
```

---

## 🛠️ SCRIPTS PRINCIPAIS

### 1. `./start_session.sh [quick]`
- **Função:** Inicia nova sessão de trabalho
- **Parâmetros:** `quick` (modo rápido, sem atualização completa de contexto)
- **Ações:**
  - Cria `work_in_progress/current_session/`
  - Arquivar sessão anterior se existir
  - Atualiza `.agent_context.md`
  - Cria link simbólico no root para compatibilidade

### 2. `./create_task.sh "Nome da Tarefa" [módulo]`
- **Função:** Cria nova tarefa com estrutura padronizada
- **Parâmetros:** Nome da tarefa (obrigatório), módulo (opcional, padrão: `ui_web`)
- **Ações:**
  - Cria diretório `work_in_progress/tasks/task_[ID]/`
  - Copia template `task_prompt.md`
  - Personaliza template com nome da tarefa
  - Cria `task_metadata` com metadados

### 3. `./process_task.sh --task=ID [--checklist|--plan|--execute]`
- **Função:** Processa tarefa em diferentes modos
- **Parâmetros obrigatórios:** `--task=ID` (ID da tarefa)
- **Modos:**
  - `--checklist`: Gera checklist pré-implementação
  - `--plan`: Gera plano de implementação completo
  - `--execute`: Prepara para execução com opencode
- **Ações:**
  - Valida existência da tarefa
  - Lê `task_prompt.md`
  - Gera arquivos correspondentes ao modo
  - Cria arquivo para opencode (modo execute)

### 4. `./conclude_task.sh --task=ID "Aprendizados" [--success|--partial|--failed]`
- **Função:** Conclui tarefa e documenta aprendizados
- **Parâmetros obrigatórios:** `--task=ID`, descrição dos aprendizados
- **Status:** `--success` (padrão), `--partial`, `--failed`
- **Ações:**
  - Valida implementação (handlers, testes, etc.)
  - Coleta métricas
  - Cria `task_learnings.md`
  - Move tarefa para archive da sessão
  - Atualiza documentação permanente

### 5. `./end_session.sh [force]`
- **Função:** Encerra sessão completa
- **Parâmetros:** `force` (força encerramento mesmo com tarefas pendentes)
- **Ações:**
  - Verifica tarefas pendentes (a menos que `force`)
  - Consolida aprendizados da sessão
  - Move sessão para archive
  - Atualiza documentação permanente
  - Limpa diretórios temporários

---

## 🎯 VANTAGENS DA NOVA ESTRUTURA

### ✅ Para o Desenvolvedor:
1. **Diretório raiz limpo** - Nenhum arquivo temporário `.session_*`, `.task_*`
2. **Organização hierárquica** - Sessão → Tarefas → Artefatos
3. **Prompt padronizado** - Template consistente para todas as tarefas
4. **Histórico preservado** - Archive mantém todas as sessões e tarefas
5. **Aprendizados estruturados** - Coletados por tarefa e consolidados por sessão

### ✅ Para o Agente (opencode):
1. **Contexto único** - `.agent_context.md` sempre no mesmo lugar
2. **Fluxo padronizado** - Sem parâmetros complexos de arquivo
3. **Referências claras** - Estrutura de diretórios previsível
4. **Aprendizados acessíveis** - `docs/learnings/` e `session_learnings/`

### ✅ Para o Projeto:
1. **Documentação automática** - Aprendizados coletados automaticamente
2. **Métricas coletáveis** - Tempo por tarefa, sucesso/falha, etc.
3. **Processo replicável** - Fluxo claro para novos contribuidores
4. **Qualidade consistente** - Checklists e planos padronizados

---

## 🔄 MIGRAÇÃO DA ESTRUTURA ANTIGA

### Script de migração:
```bash
./scripts/workflow/migrate_to_new_structure.sh
```

### O que é migrado:
1. **Sessões existentes** (`/.session_*`) → `work_in_progress/archive/`
2. **Tarefas existentes** (`/.task_*`) → Archive com metadados
3. **Contexto do agente** (`/.agent_context.md`) → Nova sessão
4. **Arquivos opencode** (`/.opencode_task_*`) → Archive/legacy

### Backup:
- Scripts antigos copiados para `scripts/workflow/backup_[DATA]/`
- Arquivos antigos podem ser mantidos ou removidos (pergunta durante migração)

---

## 📚 INTEGRAÇÃO COM SISTEMA EXISTENTE

### Documentação:
- **`docs/QUICK_REFERENCE.md`** - Atualizado automaticamente com referência a sessões
- **`docs/NEXT_STEPS.md`** - Atualizado com conclusão de tarefas
- **`docs/learnings/`** - Recebe aprendizados consolidados de cada sessão
- **`docs/ANTIPATTERNS.md`** - Pode ser atualizado com aprendizados (manual)

### Ferramentas existentes:
- **`./scripts/tools/analyze_patterns.sh`** - Continua funcionando normalmente
- **`./scripts/dev/smoke_test_new_feature.sh`** - Integrado com validação
- **`./scripts/update_context.sh`** - Chamado automaticamente

### Padrões de código:
- **Handlers:** `modules/ui_web/internal/handler/[feature]_handler.go`
- **Templates:** `modules/ui_web/templates/[feature]_simple.html`
- **Testes:** `modules/ui_web/internal/handler/[feature]_handler_test.go`
- **Registro:** `modules/ui_web/main.go`

---

## ⚠️ CONSIDERAÇÕES IMPORTANTES

### 1. Compatibilidade com Versão Anterior
- Link simbólico: `.agent_context.md` → `work_in_progress/current_session/.agent_context.md`
- Scripts antigos mantidos em backup
- Estrutura antiga pode coexistir durante transição

### 2. Gerenciamento de Espaço
- Archive pode crescer com o tempo
- Considerar limpeza periódica de sessões muito antigas
- Aprendizados importantes movidos para `docs/learnings/`

### 3. Trabalho em Equipe
- Estrutura local por desenvolvedor
- Archive não deve ser commitado no git
- Aprendizados consolidados podem ser compartilhados via `docs/learnings/`

### 4. Personalização
- Template `task_prompt.md` pode ser customizado
- Scripts podem ser estendidos para necessidades específicas
- Estrutura de archive pode ser adaptada

---

## 🚀 COMEÇANDO COM A NOVA ESTRUTURA

### Primeiros passos:
```bash
# 1. Migrar estrutura existente (opcional)
./scripts/workflow/migrate_to_new_structure.sh

# 2. Iniciar primeira sessão
./start_session.sh

# 3. Criar tarefa de teste
./create_task.sh "Tarefa de Teste" ui_web

# 4. Explorar fluxo completo
./process_task.sh --task=[ID] --checklist
./process_task.sh --task=[ID] --plan
./process_task.sh --task=[ID] --execute

# 5. Concluir tarefa
./conclude_task.sh --task=[ID] "Aprendizados: teste do novo fluxo" --success

# 6. Encerrar sessão
./end_session.sh
```

### Dicas rápidas:
- **Listar tarefas:** `ls work_in_progress/tasks/`
- **Ver contexto:** `cat work_in_progress/current_session/.agent_context.md`
- **Ver archive:** `ls -la work_in_progress/archive/`
- **Modo rápido:** `./start_session.sh quick` para sessões curtas

---

## 📞 SUPORTE E PROBLEMAS CONHECIDOS

### Problemas comuns:
1. **Tarefa não encontrada** - Verificar ID com `ls work_in_progress/tasks/`
2. **Sessão não ativa** - Executar `./start_session.sh` primeiro
3. **Permissões de script** - `chmod +x scripts/workflow/*.sh`
4. **Espaço em nomes** - Usar quotes: `./create_task.sh "Nome com Espaço"`

### Logs e debug:
- **Sessão atual:** `work_in_progress/current_session/session_info`
- **Metadados da tarefa:** `work_in_progress/tasks/task_[ID]/task_metadata`
- **Archive:** `work_in_progress/archive/session_[ID]/`

### Melhorias futuras:
1. Interface web para gerenciamento de tarefas
2. Integração com git para rastreamento de mudanças
3. Métricas avançadas (velocidade, qualidade, etc.)
4. Templates específicos por tipo de tarefa (bug, feature, refatoração)

---

**Status:** ✅ IMPLEMENTADA E PRONTA PARA USO  
**Última atualização:** 11/03/2026  
**Próxima revisão:** Após 5 sessões de uso