# 📋 TAREFA: Separar módulo contador social do módulo empreendimento social

**Data:** 13/03/2026
**Prioridade:** [ALTA/MÉDIA/BAIXA]
**Estimativa:** [X] minutos/horas
**Módulo(s):** ui_web

---
## CONTEXTO
- Nas etapas anteriores decidiu-se manter no mesmo módulo web as funcionalidades de contador social e associação contábil com as funcionalidades de empreendimento social. 
Mas com o crescimento do código e das funcionalidades manter no mesmo módulo, controlando sessões de acesso a bancos, login, etc, tem se tornado pesado. E mostrado que a separação será importante para a continuidade

## 🎯 OBJETIVO
Separar completamente as funcionalidades do empreendedor social (PDV, Compras, Sócios, etc) das funcionalidades do contador social, criando outro módulo web separado. Manter todos os padrões de interface e todos os padrões de desenvolvimento. Avaliar se módulos "backend" podem ser mantidos o mesmo, ou se tb devemos separar.
Vc, agente, é responsável por preencher as outras etapas desse prompt, conforme entendimento do projeto e do objetivo declarado.

---

## 📋 REQUISITOS

### Funcionais
- [ ] Requisito 1
- [ ] Requisito 2
- [ ] Requisito 3

### Técnicos
- [ ] Seguir padrões do projeto Digna
- [ ] Implementar testes unitários
- [ ] Implementar testes E2E com Playwright (fluxo completo)
- [ ] Testar todas as rotas novas/modificadas
- [ ] Validar regressões (não quebrar funcionalidades existentes)
- [ ] Atualizar documentação
- [ ] Validar com smoke tests
- [ ] **CRÍTICO:** Testar com dados reais (cafe_digna, contador_social)

### Não Funcionais
- [ ] Performance: [requisito]
- [ ] Segurança: [requisito]
- [ ] Usabilidade: [requisito]

---

## 🔍 CONTEXTO E ANÁLISE

### Módulos/Arquivos Relacionados
- `modules/[módulo]/...` - [descrição]
- `docs/...` - [documentação relevante]

### Padrões a Seguir
- [ ] Analisar handler similar: `modules/ui_web/internal/handler/[handler]_handler.go`
- [ ] Analisar template similar: `modules/ui_web/templates/[template]_simple.html`
- [ ] Seguir padrão SHA256 do `core_lume`
- [ ] Usar anti-padrões de `docs/ANTIPATTERNS.md`

### Dependências
- [ ] Feature X precisa ser implementada primeiro
- [ ] Integração com módulo Y
- [ ] Atualização de banco de dados

---

## 🚀 PLANO DE IMPLEMENTAÇÃO

### Fase 1: Análise e Preparação
1. [ ] Analisar código existente similar
2. [ ] Criar checklist de implementação
3. [ ] Verificar aprendizados anteriores

### Fase 2: Implementação
1. [ ] Criar/atualizar handler
2. [ ] Criar/atualizar template
3. [ ] Implementar lógica de negócio
4. [ ] Adicionar ao `main.go`

### Fase 3: Testes e Validação
1. [ ] Implementar testes unitários
2. [ ] Executar smoke tests
3. [ ] Validar integração
4. [ ] Testar manualmente

### Fase 4: Documentação e Conclusão
1. [ ] Atualizar documentação
2. [ ] Documentar aprendizados
3. [ ] Atualizar checklists/antipadrões

---

## 📁 ARQUIVOS ESPERADOS

### A Criar
- `modules/ui_web/internal/handler/[nome]_handler.go`
- `modules/ui_web/templates/[nome]_simple.html`
- `modules/ui_web/internal/handler/[nome]_handler_test.go`

### A Modificar
- `modules/ui_web/main.go` (registrar handler)
- `docs/QUICK_REFERENCE.md` (atualizar referência)
- `docs/NEXT_STEPS.md` (marcar como concluído)

---

## ⚠️ RISCOS E DESAFIOS

### Riscos Técnicos
1. [Risco 1] - [Mitigação]
2. [Risco 2] - [Mitigação]

### Riscos de Processo
1. [Risco 1] - [Mitigação]
2. [Risco 2] - [Mitigação]

---

## 📚 APRENDIZADOS ANTERIORES RELEVANTES

[Consultar `docs/learnings/` para tarefas similares]
- `docs/learnings/[arquivo].md` - [resumo do aprendizado]

---

## 🔗 LINKS ÚTEIS

- [Documentação do projeto](docs/)
- [Padrões de código](docs/QUICK_REFERENCE.md)
- [Antipadrões](docs/ANTIPATTERNS.md)
- [Skills do projeto](docs/skills/)

---

**Status:** [PENDENTE/EM ANDAMENTO/CONCLUÍDA]
**Última atualização:** 13/03/2026

---

## 📋 REQUISITOS

### Funcionais
- [ ] Requisito 1
- [ ] Requisito 2
- [ ] Requisito 3

### Técnicos
- [ ] Seguir padrões do projeto Digna
- [ ] Implementar testes unitários
- [ ] Implementar testes E2E com Playwright (fluxo completo)
- [ ] Testar todas as rotas novas/modificadas
- [ ] Validar regressões (não quebrar funcionalidades existentes)
- [ ] Atualizar documentação
- [ ] Validar com smoke tests
- [ ] **CRÍTICO:** Testar com dados reais (cafe_digna, contador_social)

### Não Funcionais
- [ ] Performance: [requisito]
- [ ] Segurança: [requisito]
- [ ] Usabilidade: [requisito]

---

## 🔍 CONTEXTO E ANÁLISE

### Módulos/Arquivos Relacionados
- `modules/[módulo]/...` - [descrição]
- `docs/...` - [documentação relevante]

### Padrões a Seguir
- [ ] Analisar handler similar: `modules/ui_web/internal/handler/[handler]_handler.go`
- [ ] Analisar template similar: `modules/ui_web/templates/[template]_simple.html`
- [ ] Seguir padrão SHA256 do `core_lume`
- [ ] Usar anti-padrões de `docs/ANTIPATTERNS.md`

### Dependências
- [ ] Feature X precisa ser implementada primeiro
- [ ] Integração com módulo Y
- [ ] Atualização de banco de dados

---

## 🚀 PLANO DE IMPLEMENTAÇÃO

### Fase 1: Análise e Preparação
1. [ ] Analisar código existente similar
2. [ ] Criar checklist de implementação
3. [ ] Verificar aprendizados anteriores

### Fase 2: Implementação
1. [ ] Criar/atualizar handler
2. [ ] Criar/atualizar template
3. [ ] Implementar lógica de negócio
4. [ ] Adicionar ao `main.go`

### Fase 3: Testes e Validação
1. [ ] Implementar testes unitários
2. [ ] Executar smoke tests
3. [ ] Validar integração
4. [ ] Testar manualmente

### Fase 4: Documentação e Conclusão
1. [ ] Atualizar documentação
2. [ ] Documentar aprendizados
3. [ ] Atualizar checklists/antipadrões

---

## 📁 ARQUIVOS ESPERADOS

### A Criar
- `modules/ui_web/internal/handler/[nome]_handler.go`
- `modules/ui_web/templates/[nome]_simple.html`
- `modules/ui_web/internal/handler/[nome]_handler_test.go`

### A Modificar
- `modules/ui_web/main.go` (registrar handler)
- `docs/QUICK_REFERENCE.md` (atualizar referência)
- `docs/NEXT_STEPS.md` (marcar como concluído)

---

## ⚠️ RISCOS E DESAFIOS

### Riscos Técnicos
1. [Risco 1] - [Mitigação]
2. [Risco 2] - [Mitigação]

### Riscos de Processo
1. [Risco 1] - [Mitigação]
2. [Risco 2] - [Mitigação]

---

## 📚 APRENDIZADOS ANTERIORES RELEVANTES

[Consultar `docs/learnings/` para tarefas similares]
- `docs/learnings/[arquivo].md` - [resumo do aprendizado]

---

## 🔗 LINKS ÚTEIS

- [Documentação do projeto](docs/)
- [Padrões de código](docs/QUICK_REFERENCE.md)
- [Antipadrões](docs/ANTIPATTERNS.md)
- [Skills do projeto](docs/skills/)

---

**Status:** [PENDENTE/EM ANDAMENTO/CONCLUÍDA]
**Última atualização:** 13/03/2026