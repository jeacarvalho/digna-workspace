# Padrão de Sessão - Digna

---

## 1. Estrutura de Sessão

Cada sessão de trabalho deve seguir um padrão estruturado para garantir consistência e rastreabilidade.

### 1.1 Início de Sessão

**Passos obrigatórios:**

1. **Verificar Contexto**
   - Ler `docs/README.md` para orientação geral
   - Verificar status atual em `06_roadmap/04_status.md`
   - Consultar `05_ai/01_constitution.md` para regras técnicas

2. **Identificar Tarefa**
   - Ler requisitos relacionados
   - Verificar dependências
   - Listar arquivos a modificar

3. **Preparar Ambiente**
   - Verificar estrutura de diretórios
   - Confirmar ferramentas disponíveis

---

### 1.2 Durante a Sessão

**Práticas recomendadas:**

1. **Implementação**
   - Seguir Clean Architecture
   - Criar testes unitários junto com código
   - Man eter interfaces pequenas focadas

2. **Validação**
   - Executar testes frequentemente
   - Verificar lint/format
   - Validar contra requisitos (RF-XX, RNF-XX)

3. **Comunicação**
   - Documentar decisões no código
   - Explicar mudanças significativas
   - Alertar sobre trade-offs

---

### 1.3 Fim de Sessão

**Passos obrigatórios:**

1. **Atualizar SESSION_LOG**
   - Criar entrada em `06_roadmap/05_session_log.md`
   - Documentar o que foi feito
   - Listar decisões tomadas

2. **Registrar Decisões**
   - Decisões arquiteturais
   - Mudanças de requisitos
   - Trade-offs aceitos

3. **Documentar Próximos Passos**
   - Tarefas pendentes
   - Dependências identificadas
   - Riscos encontrados

---

## 2. Formato de Log de Sessão

```markdown
## Session Log [NÚMERO] - [TÍTULO]

**Date:** YYYY-MM-DD
**Status:** [STATUS] | [TESTES]

### Summary
[Descrição breve do que foi realizado]

### What Was Implemented
- ✅ [Componente 1]
- ✅ [Componente 2]

### Technical Decisions
- [Decisão 1]: [Justificativa]
- [Decisão 2]: [Justificativa]

### Test Results
```
[Resultado dos testes]
```

### DoD Validated
1. ✅ [Critério 1]
2. ✅ [Critério 2]

### Next Steps
- [Próxima tarefa 1]
- [Próxima tarefa 2]
```

---

## 3. Critérios de Conclusão

A sessão deve terminar apenas quando:

- [ ] Código implementado
- [ ] Testes passando
- [ ] Documentação atualizada
- [ ]SESSION_LOG atualizado
- [ ] Próximos passos documentados

---

## 4. Boas Práticas

### 4.1 atomicidade

Cada task deve ser completada antes de iniciar outra. Não deixar código pela metade.

### 4.2 Revisão

Antes de finalizar, revisar:
- Convenção de nomes seguida
- Comentários desnecessários removidos
- Testes cobrindo caso de borda

### 4.3 Rastreabilidade

Garantir que cada mudança possa ser rastreada:
- Para requisito (RF-XX)
- Para decisão de design
- Para sessão específica

---

## 5. Handling de Erros

### 5.1 Erro de Implementação

Se encontrar erro durante implementação:
1. Documentar o erro encontrado
2. Propor solução alternativa
3. Marcar tarefa como pendente

### 5.2 Conflito de Contexto

Se contexto mudar durante sessão:
1. Parar implementação atual
2. Atualizar documentação
3. Reiniciar com novo contexto

### 5.3 Dependência Bloqueada

Se dependência estiver bloqueando:
1. Documentar bloqueio
2. Marcar como dependência externa
3. Avançar em tarefas independentes

---

## 6. Referência

- Ver `06_roadmap/05_session_log.md` para exemplos de sessão
- Ver `docs/04_governance/governance.md` para políticas de contribuição
