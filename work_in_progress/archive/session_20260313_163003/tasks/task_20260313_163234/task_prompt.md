# 📋 TAREFA: Avaliar e corrigir RF 12

**Data:** 13/03/2026
**Prioridade:** [ALTA/MÉDIA/BAIXA]
**Estimativa:** [X] minutos/horas
**Módulo(s):** ui_web

---

## 🎯 OBJETIVO

Finalizamos a implementação do RF 12 (veja docs/02_product/01_requirements.md). 
Mas não encontro onde associar um contador a um empreendimento. 
Também separamos os módulos. Agora temos 3: o autenticador, o empreendedor e o contador. Um empreendedor deve poder informar qual contador está cuidando da empresa. Para isso ele precisa abrir o banco central para ler os contadores disponível e informar qual irá cuidar de seu empreendimento. Essa relaçaõ precisa estarn o banco central pq de lá será carregado as empresas de um contador, quando ele logar. 
Penso que tudo isso foi implementado, mas faltou a tela no módulo empreendedor (ui_web). Sua tarefa é deixar o RF 12 concluído nessa nova arquitetura de 3 módulos.
Você, agente, tb é responsável de concluir as seções desse prompt com as informações necessárias. 
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