title: Padrão de Sessão - Ecossistema Digna
status: implemented
version: 2.0
last_updated: 2026-03-27
---

# Padrão de Sessão - Ecossistema Digna

> **Nota:** Este documento reflete o padrão de sessão integrado do Ecossistema Digna (PDF v1.0), preservando toda a infraestrutura validada nas Sprints 1-16, incorporando os novos módulos como Fases 3 e 4, e incluindo o Sistema de Ajuda Educativa (RF-30) decidido na sessão de 27/03/2026.

---

## 1. Estrutura de Sessão [ATUALIZADO]

Cada sessão de trabalho deve seguir um padrão estruturado para garantir consistência e rastreabilidade, agora alinhado com a arquitetura de Ecossistema de 4 Módulos.

### 1.1 Início de Sessão

**Passos obrigatórios:**

1. **Verificar Contexto**
   - Ler `docs/README.md` para orientação geral
   - Verificar status atual em `06_roadmap/04_status.md`
   - Consultar `05_ai/01_constitution.md` para regras técnicas (incluindo RF-30)
   - Verificar `06_roadmap/03_backlog.md` para tarefas prioritárias

2. **Identificar Tarefa**
   - Ler requisitos relacionados (RF-01 a RF-30)
   - Verificar dependências entre módulos (ERP → Motor → Portal → Rede)
   - Listar arquivos a modificar
   - Identificar módulo afetado (1-4 ou transversal)

3. **Preparar Ambiente**
   - Verificar estrutura de diretórios
   - Confirmar ferramentas disponíveis
   - Executar `./scripts/tools/quick_agent_check.sh all`

### 1.2 Durante a Sessão

**Práticas recomendadas:**

**Implementação:**
- Seguir Clean Architecture
- Criar testes unitários junto com código
- Manter interfaces pequenas focadas
- Respeitar princípio "Nenhum dado digitado duas vezes"
- Consultar skills relevantes em `docs/skills/`

**Validação:**
- Executar testes frequentemente
- Verificar lint/format
- Validar contra requisitos (RF-XX, RNF-XX)
- **Validar E2E:** `./scripts/dev/validate_e2e.sh --basic --headless`

**Comunicação:**
- Documentar decisões no código
- Explicar mudanças significativas
- Alertar sobre trade-offs
- Registrar aprendizados em tempo real

### 1.3 Fim de Sessão

**Passos obrigatórios:**

1. **Atualizar SESSION_LOG**
   - Criar entrada em `06_roadmap/05_session_log.md`
   - Documentar o que foi feito
   - Listar decisões tomadas

2. **Registrar Decisões**
   - Decisões arquiteturais (ADRs)
   - Mudanças de requisitos
   - Trade-offs aceitos

3. **Documentar Próximos Passos**
   - Tarefas pendentes
   - Dependências identificadas
   - Riscos encontrados

4. **Concluir Tarefas**
   - Executar `./conclude_task.sh "Aprendizados + resultado E2E" --success`
   - Validar que todos os testes passam
   - Mover tarefa para archive

---

## 2. Formato de Log de Sessão [ATUALIZADO]

```markdown
## Session Log [NÚMERO] - [TÍTULO]

**Date:** YYYY-MM-DD
**Status:** [STATUS] | [TESTES]
**Módulos Afetados:** [ERP/Motor/Portal/Rede/Transversal]

### Summary
[Descrição breve do que foi realizado]

### What Was Implemented
- ✅ [Componente 1]
- ✅ [Componente 2]
- ✅ [RF-XX implementado/atualizado]

### Technical Decisions
- [Decisão 1]: [Justificativa]
- [Decisão 2]: [Justificativa]
- [Impacto em outros módulos]: [Descrição]

### Test Results
[Resultado dos testes]
- Unitários: [X]/[Y] PASS
- E2E: [PASS/FAIL]

### DoD Validated
1. ✅ [Critério 1]
2. ✅ [Critério 2]
3. ✅ [Validação E2E passou]

### Next Steps
- [Próxima tarefa 1]
- [Próxima tarefa 2]
- [Módulo a ser desenvolvido]
```

---

## 3. Critérios de Conclusão [ATUALIZADO]

A sessão deve terminar apenas quando:

- [ ] Código implementado
- [ ] Testes passando (>90% cobertura para handlers)
- [ ] **Validação E2E passando** (`./scripts/dev/validate_e2e.sh --basic --headless`)
- [ ] Documentação atualizada
- [ ] SESSION_LOG atualizado
- [ ] Próximos passos documentados
- [ ] Aprendizados registrados em `docs/learnings/`
- [ ] Constitution de IA respeitada (Anti-Float, Cache-Proof, Soberania)
- [ ] RF-30 aplicado (se campo técnico adicionado)

---

## 4. Boas Práticas [ATUALIZADO]

### 4.1 Atomicidade

Cada task deve ser completada antes de iniciar outra. Não deixar código pela metade.

**Para Ecossistema:**
- Implementar módulo por módulo (respeitar dependências)
- ERP primeiro, depois Motor, Portal, Rede
- RF-30 é transversal - aplicar em todos os módulos

### 4.2 Revisão

Antes de finalizar, revisar:

- Convenção de nomes seguida
- Comentários desnecessários removidos
- Testes cobrindo caso de borda
- **Princípio "Nenhum dado digitado duas vezes" respeitado**
- **Campos técnicos com botão "?" (RF-30)**

### 4.3 Rastreabilidade

Garantir que cada mudança possa ser rastreada:

- Para requisito (RF-XX)
- Para decisão de design (ADR-XX)
- Para sessão específica (Session Log XX)
- Para módulo do ecossistema (1-4)

---

## 5. Handling de Erros [ATUALIZADO]

### 5.1 Erro de Implementação

Se encontrar erro durante implementação:

1. Documentar o erro encontrado
2. Propor solução alternativa
3. Marcar tarefa como pendente
4. Registrar em `docs/learnings/`

### 5.2 Conflito de Contexto

Se contexto mudar durante sessão:

1. Parar implementação atual
2. Atualizar documentação
3. Reiniciar com novo contexto
4. Preservar contexto com `./preserve_context.sh --save`

### 5.3 Dependência Bloqueada

Se dependência estiver bloqueando:

1. Documentar bloqueio
2. Marcar como dependência externa
3. Avançar em tarefas independentes
4. Notificar no SESSION_LOG

### 5.4 Compaction do Agente [NOVO - 27/03/2026]

Se agente perder contexto (compaction):

1. Usuário executa: `./preserve_context.sh --save`
2. Aguardar compaction terminar
3. Usuário executa: `./scripts/help_agent_recover_context.sh`
4. Usuário executa: `./preserve_context.sh --restore`
5. Verificar que agente recuperou contexto
6. Continuar implementação

---

## 6. Integração com Workflow de Agentes [ATUALIZADO]

### 6.1 Scripts de Workflow

| Script | Função | Quando Usar |
|--------|--------|-------------|
| `./start_session.sh` | Inicia sessão com contexto | Início de cada sessão |
| `./process_task.sh` | Processa tarefa (checklist/plan/execute) | Antes de implementar |
| `./conclude_task.sh` | Conclui tarefa + documenta aprendizados | Após implementação |
| `./end_session.sh` | Encerra sessão + consolida aprendizados | Fim de cada sessão |
| `./preserve_context.sh` | Preserva contexto durante compaction | Quando detectar perda de contexto |
| `./scripts/help_agent_recover_context.sh` | Recupera contexto após compaction | Após compaction terminar |

### 6.2 Fluxo Completo

```
┌─────────────────┐
│ start_session   │
└────────┬────────┘
         │
┌────────▼────────┐
│  process_task   │───┐
│  (--checklist)  │   │
└────────┬────────┘   │
         │            │
┌────────▼────────┐   │
│  process_task   │   │
│  (--plan)       │   │
└────────┬────────┘   │
         │            │
┌────────▼────────┐   │
│  process_task   │◀──┘
│  (--execute)    │
└────────┬────────┘
         │
┌────────▼────────┐    ┌──────────────┐
│   Implementar   │───▶│  validate_e2e│
│   (opencode)    │    │  (obrigatório)│
└────────┬────────┘    └──────────────┘
         │
┌────────▼────────┐
│  conclude_task  │
└────────┬────────┘
         │
┌────────▼────────┐
│   end_session   │
└─────────────────┘
```

---

## 7. Contexto do Ecossistema [NOVO - PDF v1.0]

### 7.1 Módulos do Ecossistema

| Módulo | Responsabilidade | RFs Relacionados |
|--------|------------------|------------------|
| **Módulo 1: digna ERP** | Gestão financeira, fiscal e contábil | RF-01 a RF-13, RF-27 |
| **Módulo 2: Motor de Indicadores** | Coleta APIs externas (BCB, IBGE) | RF-18 |
| **Módulo 3: Portal de Oportunidades** | Match automático de crédito | RF-19 a RF-23 |
| **Módulo 4: Rede Digna** | Marketplace solidário B2B | RF-24 a RF-26 |
| **Transversal: Ajuda Educativa** | Tradução de conceitos técnicos | RF-30 |

### 7.2 Dependências entre Módulos

```
ERP (Módulo 1) ──────┬──────▶ Motor (Módulo 2)
                     │
                     ├──────▶ Portal (Módulo 3)
                     │
                     └──────▶ Rede (Módulo 4)

Ajuda Educativa (RF-30) ───▶ TODOS OS MÓDULOS (transversal)
```

### 7.3 Regras de Sessão por Módulo

**Módulo 1 (ERP):**
- Preservar funcionalidades existentes (Sprints 1-16)
- Anti-Float obrigatório
- Cache-Proof templates

**Módulo 2 (Motor):**
- Cache local com TTL para APIs externas
- Circuit breaker para indisponibilidade
- Interpretação contextualizada dos dados

**Módulo 3 (Portal):**
- Princípio "Nenhum dado digitado duas vezes"
- Match automático (sem formulários)
- Dados sensíveis nunca transmitidos

**Módulo 4 (Rede):**
- Perfil público sem dados sensíveis
- Matching geográfico e setorial
- Hash anonimizado para EntityID

**Transversal (RF-30):**
- Botão "?" em campos técnicos
- Linguagem para 5ª série
- Linkagem UI → `help_topics{}`

---

## 8. Validação E2E [NOVO - 27/03/2026]

### 8.1 Critérios de Validação

Uma tarefa só deve ser marcada como "testada end-to-end" quando:

- ✅ Testes unitários passam (>90% cobertura)
- ✅ Smoke test HTTP passa
- ✅ **Validação E2E passa** (`./scripts/dev/validate_e2e.sh --basic --headless`)
- ✅ Aprendizados documentados no `conclude_task.sh`

### 8.2 Fluxo de 7 Passos Padrão Digna

1. [ ] Login no sistema
2. [ ] Criar item de estoque (se não existir)
3. [ ] Criar membro (se não existir)
4. [ ] Criar fornecedor (se não existir)
5. [ ] Registrar compra do item
6. [ ] Registrar venda no PDV
7. [ ] Confirmar saldo e registrar horas

### 8.3 Modos de Execução

```bash
# Modo stealth (padrão - não abre janelas)
./scripts/dev/validate_e2e.sh --basic --headless

# Com navegador visível (debug)
./scripts/dev/validate_e2e.sh --basic --ui

# Todos os testes
./scripts/dev/validate_e2e.sh --full --headless

# Teste específico
./scripts/dev/validate_e2e.sh --custom "fluxo" --headless
```

---

## 9. Referências [ATUALIZADO]

### 9.1 Documentos Relacionados

| Documento | Finalidade |
|-----------|------------|
| `05_ai/01_constitution.md` | Regras sagradas da IA (Anti-Float, Cache-Proof, etc.) |
| `06_roadmap/05_session_log.md` | Logs de sessões anteriores |
| `06_roadmap/03_backlog.md` | Backlog completo (RF-01 a RF-30) |
| `docs/learnings/` | Aprendizados de sessões anteriores |
| `docs/skills/` | Skills específicas do projeto |
| `COMPACTION_HANDLING.md` | Manual de preservação de contexto |

### 9.2 Skills do Projeto

| Skill | Foco | Arquivo |
|-------|------|---------|
| **developing-digna-backend** | Rigor técnico, DDD, TDD, Anti-Float | `skills/developing-digna-backend/SKILL.md` |
| **rendering-digna-frontend** | HTMX, UI "Soberania e Suor", Cache-Proof | `skills/rendering-digna-frontend/SKILL.md` |
| **managing-sovereign-data** | Isolamento SQLite, LifecycleManager | `skills/managing-sovereign-data/SKILL.md` |
| **applying-solidarity-logic** | Tradução cultural, ITG 2002, pedagogia | `skills/applying-solidarity-logic/SKILL.md` |
| **auditing-fiscal-compliance** | Accountant Dashboard, SPED, Read-Only | `skills/auditing-fiscal-compliance/SKILL.md` |

### 9.3 Comandos de Referência

```bash
# Iniciar sessão
./start_session.sh [quick]

# Processar tarefa
./process_task.sh "descrição" --checklist
./process_task.sh "descrição" --plan
./process_task.sh "descrição" --execute

# Validar E2E (OBRIGATÓRIO)
./scripts/dev/validate_e2e.sh --basic --headless

# Concluir tarefa
./conclude_task.sh "Aprendizados + resultado E2E" --success

# Encerrar sessão
./end_session.sh [force]

# Preservar contexto (compaction)
./preserve_context.sh --save
./preserve_context.sh --restore
./scripts/help_agent_recover_context.sh
```

---

## 10. Métricas de Sessão [NOVO - 27/03/2026]

### 10.1 Métricas a Coletar

| Métrica | Alvo | Como Medir |
|---------|------|------------|
| Tarefas por sessão | 2-4 | Contagem no SESSION_LOG |
| Tempo médio por tarefa | < 2 horas | Timestamps de início/fim |
| Taxa de sucesso E2E | 100% | Resultado do `validate_e2e.sh` |
| Cobertura de testes | >90% | `go test ./... -cover` |
| Aprendizados documentados | 100% | Arquivos em `docs/learnings/` |
| Reincidência de bugs | 0% | Comparação com sessões anteriores |

### 10.2 Dashboard de Qualidade

```
📊 QUALIDADE DA SESSÃO ${SESSION_ID}
├── ✅ Testes unitários: [X]%
├── ✅ Testes integração: [X]%
├── ✅ Smoke test HTTP: [PASS/FAIL]
├── ✅ Validação E2E: [PASS/FAIL] ([X]/7 passos)
├── ⏱️  Tempo total: [X]h[Y]min
└── 📝 Aprendizados: [N] documentos criados
```

---

## 11. Checklist de Sessão [ATUALIZADO]

### Antes de Iniciar

- [ ] `./start_session.sh` executado
- [ ] Contexto lido (`.agent_context.md`)
- [ ] Constituição de IA revisada (`05_ai/01_constitution.md`)
- [ ] Backlog consultado (`06_roadmap/03_backlog.md`)
- [ ] `./scripts/tools/quick_agent_check.sh all` executado

### Durante Implementação

- [ ] TDD seguido (testes primeiro)
- [ ] Anti-Float respeitado (zero `float`)
- [ ] Cache-Proof templates (`ParseFiles()` no handler)
- [ ] Soberania mantida (entity_id isolation)
- [ ] RF-30 aplicado (botão "?" em campos técnicos)
- [ ] Princípio "Nenhum dado digitado duas vezes" respeitado

### Antes de Concluir

- [ ] Testes unitários passando (>90% cobertura)
- [ ] **Validação E2E passando** (`validate_e2e.sh --basic --headless`)
- [ ] Documentação atualizada
- [ ] Aprendizados registrados
- [ ] SESSION_LOG atualizado
- [ ] Próximos passos documentados

### Após Encerrar Sessão

- [ ] `./end_session.sh` executado
- [ ] Aprendizados consolidados em `docs/learnings/`
- [ ] `06_roadmap/05_session_log.md` atualizado
- [ ] `06_roadmap/04_status.md` atualizado
- [ ] Contexto preservado (se necessário)

---

## 12. Lições Aprendidas por Sessão [NOVO - 27/03/2026]

### Sessão 11/03/2026

- **Descoberta:** `legal_facade` já existe com 80% da funcionalidade
- **Aprendizado:** Consultar `docs/skills/` antes de implementar
- **Ação:** Criar `MODULES_QUICK_REFERENCE.md` para acelerar descoberta

### Sessão 27/03/2026

- **Descoberta:** Validação E2E não estava integrada ao fluxo de conclusão
- **Aprendizado:** Smoke test valida HTTP, E2E valida negócio
- **Ação:** Script `validate_e2e.sh` obrigatório antes de `conclude_task.sh`
- **Descoberta:** Campos técnicos sem explicação violam Pilar Pedagógico
- **Aprendizado:** Sistema de ajuda estruturada é infraestrutura, não feature
- **Ação:** RF-30 adicionado ao backlog com prioridade alta

---

**Status:** ✅ ATUALIZADO COM ECOSSISTEMA DE 4 MÓDULOS (PDF v1.0) + RF-30 (Sessão 27/03/2026)  
**Próxima Ação:** Atualizar `06_roadmap/01_strategy.md` com estratégia de release do ecossistema  
**Versão Anterior:** 1.0 (2026-03-09)  
**Versão Atual:** 2.0 (2026-03-27)
