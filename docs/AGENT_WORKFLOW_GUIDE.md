# 🤖 Guia de Workflow para o Agente (OpenCode)

Este documento explica como o agente deve interagir com o sistema de workflow do Projeto Digna.

## 📋 Fluxo Completo do Agente

### FASE 1: INICIALIZAÇÃO DA SESSÃO
```
Usuário executa → ./start_session.sh
                ↓
Agente deve ler → .agent_context.md (CRIADO AUTOMATICAMENTE)
                ↓
Agente ganha contexto completo do projeto
```

**O que o agente faz:**
1. Usuário executa `./start_session.sh`
2. Script cria `.agent_context.md` com:
   - Regras sagradas (anti-float, soberania, etc.)
   - Status atual do projeto
   - Instruções específicas para o agente
3. **Agente DEVE ler `.agent_context.md` primeiro**
4. Agente segue instruções no arquivo

### FASE 2: PROCESSAMENTO DE TAREFA
```
Usuário executa → ./process_task.sh "descrição" --execute
                ↓
Script gera → .opencode_task_*.txt (prompt estruturado)
            ↓
Script atualiza → .agent_context.md (com tarefa ativa)
                ↓
Usuário copia prompt → Cola no opencode
                      ↓
Agente implementa → Seguindo prompt e contexto
```

**O que o agente faz:**
1. Usuário executa `./process_task.sh "Tipo: Feature | Módulo: ui_web | Objetivo: X" --execute`
2. Script gera:
   - Checklist pré-implementação
   - Plano de implementação  
   - Prompt estruturado (`.opencode_task_*.txt`)
3. Script atualiza `.agent_context.md` com tarefa ativa
4. Usuário copia prompt para opencode
5. **Agente lê `.agent_context.md` para contexto**
6. Agente implementa seguindo prompt + regras do projeto

### FASE 3: IMPLEMENTAÇÃO
```
Agente segue → Plano em docs/implementation_plans/
             ↓
Agente implementa → Com TDD e padrões
                   ↓
Agente valida → ./scripts/dev/smoke_test_new_feature.sh
               ↓
Se passar → Validação E2E → ./scripts/dev/validate_e2e.sh --basic --headless
Se falhar → Corrige        ↓
                          Se passar → Continua
                          Se falhar → Corrige fluxo E2E
```

**Regras durante implementação:**
1. **Sempre** seguir regras sagradas do `.agent_context.md`
2. **Sempre** usar TDD (testes primeiro)
3. **Sempre** validar com smoke test após implementação
4. **Nunca** ignorar antipadrões documentados

### FASE 4: VALIDAÇÃO E2E (OBRIGATÓRIA)
```
Agente implementa → Feature completa
                   ↓
Usuário executa → ./scripts/dev/validate_e2e.sh --basic --headless
                 ↓
Se passar → Documenta resultado
Se falhar → Agente corrige problemas E2E
           ↓
           Valida novamente → Até passar
```

### FASE 5: CONCLUSÃO
```
Validação E2E → Passou
               ↓
Usuário executa → ./conclude_task.sh "aprendizados + resultado E2E" --success
                 ↓
Script documenta → Aprendizados em docs/learnings/
                   ↓
Script atualiza → .agent_context.md (tarefa concluída)
                   ↓
Sistema pronto → Para próxima tarefa
```

**O que o agente faz (indiretamente via usuário):**
1. Implementação completa e validada
2. Usuário executa `./conclude_task.sh "Aprendizados: X, Y, Z" --success`
3. Script:
   - Cria documento de aprendizados
   - Atualiza checklists e antipadrões
   - Atualiza `.agent_context.md` com conclusão
   - Prepara NEXT_STEPS.md para próxima sessão
4. Sistema fica pronto para nova tarefa

## 🎯 Arquivos Chave para o Agente

### 1. `.agent_context.md` (MAIS IMPORTANTE)
- **Gerado por:** `start_session.sh`
- **Atualizado por:** `process_task.sh`, `conclude_task.sh`
- **Conteúdo:** Instruções específicas para o agente
- **Quando ler:** SEMPRE no início da interação

### 2. `docs/QUICK_REFERENCE.md`
- Arquitetura core, padrões, antipadrões
- Referência rápida de comandos e estrutura

### 3. `docs/implementation_plans/`
- Planos de implementação específicos
- Checklists pré-implementação

### 4. `.opencode_task_*.txt`
- Prompt estruturado para tarefa específica
- Gerado por `process_task.sh --execute`

## 🔄 Exemplo de Interação Completa

### Sessão 1: Nova Feature
```
Usuário: ./start_session.sh
Agente: [Lê .agent_context.md automaticamente]
        → Entende projeto, regras, status

Usuário: ./process_task.sh "Tipo: Feature | Módulo: ui_web | Objetivo: Implementar UI para Fornecedores" --execute
Agente: [Usuário copia prompt para opencode]
        [Agente lê .agent_context.md → vê tarefa ativa]
        [Agente implementa seguindo prompt + regras]
        [Agente valida com smoke test]

Usuário: ./conclude_task.sh "Aprendizados: checklist antecipou 2 problemas, testes 95% cobertura" --success
Agente: [Sistema documenta aprendizados]
        [.agent_context.md atualizado com conclusão]
```

### Sessão 2: Continuação
```
Usuário: ./start_session.sh
Agente: [Lê .agent_context.md atualizado]
        → Vê tarefa anterior concluída
        → Vê aprendizados documentados
        → Pronto para nova tarefa
```

## ⚠️ Armadilhas Comuns a Evitar

### ❌ NÃO FAZER:
- Ignorar `.agent_context.md`
- Violar regras sagradas (anti-float, etc.)
- Pular validação com smoke test
- Não seguir padrões de templates (`*_simple.html`)
- Esquecer de registrar handler no `main.go`

### ✅ SEMPRE FAZER:
- Ler `.agent_context.md` primeiro
- Seguir checklists pré-implementação
- Usar TDD (testes primeiro)
- Validar com smoke test
- **VALIDAR COM E2E** (`./scripts/dev/validate_e2e.sh --basic --headless`)
- Documentar aprendizados via `conclude_task.sh`

## 🛠️ Comandos de Referência para o Agente

```bash
# Durante implementação (agente executa):
cd modules && ./run_tests.sh           # Validar todos testes
./scripts/dev/smoke_test_new_feature.sh "Desc" "/rota"  # Validar feature

# Validação E2E (usuário executa APÓS implementação):
./scripts/dev/validate_e2e.sh --basic --headless  # Validação stealth
./scripts/dev/validate_e2e.sh --basic --ui        # Com navegador (debug)

# Para usuário executar:
./start_session.sh                     # Iniciar sessão
./process_task.sh "desc" --execute     # Processar tarefa  
./conclude_task.sh "aprend" --success  # Concluir tarefa
```

## 📈 Benefícios deste Sistema

1. **Contexto persistente:** `.agent_context.md` mantém estado entre interações
2. **Padronização:** Checklists garantem qualidade consistente
3. **Aprendizado contínuo:** Documentação automática de aprendizados
4. **Validação sistemática:** Smoke tests + E2E tests obrigatórios
5. **Qualidade real:** Valida fluxo de negócio, não apenas código
6. **Stealth mode:** E2E em modo headless não interfere com desktop
7. **Orientação clara para agente:** Instruções específicas no contexto

---

**Última atualização:** $(date +%d/%m/%Y)

**Status do sistema:** ✅ FUNCIONAL - Pronto para uso com opencode

**Próximo passo:** Testar o fluxo completo com uma tarefa real.