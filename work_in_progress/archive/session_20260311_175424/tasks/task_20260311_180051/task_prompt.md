# 📋 TAREFA: [Implementar a Gestão de Vínculo Contábil e Delegação Temporal (RF-12)]

**Data:** [DATA]
**Prioridade:** [ALTA]
**Estimativa:** [X] minutos/horas
**Módulo(s):** [VOCÊ, AGENTE, DEVE DEFINIR, CONFORME AS ORIENTAÇÕES NA SEÇÃO "OBJETIVO"]

---

## 🎯 OBJETIVO

Implementar a Gestão de Vínculo Contábil e Delegação Temporal (RF-12), garantindo o controle de responsabilidade técnica.
Decisões: A entidade de associação `EnterpriseAccountant` DEVE ser armazenada no Banco de Dados Central (gerido pelo módulo `lifecycle`), e NÃO no banco do Tenant. O empreendimento detém o "Exit Power" (Soberania) para encerrar o vínculo. Contadores desativados perdem acesso aos lançamentos correntes, mantendo apenas acesso Read-Only aos dados do seu período de vigência técnica (`start_date` a `end_date`).

### 📝 Descrição da Tarefa: Implementar a Gestão de Vínculo Contábil no Banco Central
*   **Requisito Funcional (RF):** RF-12 (Gestão de Vínculo Contábil e Delegação Temporal).
*   **Sprint Relacionada:** Sprint 18 (Ponte Institucional e Aliança Contábil).

Para esta tarefa, você deve carregar e seguir estritamente as instruções das seguintes skills em docs/skills/:
1. [developing-digna-backend]
2. [managing-sovereign-data]

*   **Anti-Float:** Se envolver cálculos de valor ou tempo, use estritamente int64. Proibido float.
*   **Cache-Proof:** Se houver interface, o template deve ser _simple.html carregado via ParseFiles no Handler.
*   **Soberania:** Garanta que a operação respeite o isolamento do arquivo .db do tenant atual. NENHUMA tabela de relação global pode ser salva no banco isolado da entidade.

---
**🎯 Objetivo da Tarefa**
Materializar o RF-12 criando a entidade `EnterpriseAccountant` que controla qual Contador Social tem acesso a qual Empreendimento (Tenant) e por qual período. Como esta é uma relação inter-tenant de governança global, o registro deve ocorrer obrigatoriamente no **Banco Central** do sistema (gerenciado pelo módulo `lifecycle`), protegendo a arquitetura de *micro-databases*. Você deve implementar as rotinas de criação de vínculo, encerramento de vínculo (Exit Power da cooperativa) e a regra de filtro temporal para acesso.

**📁 Estrutura de Output Esperada**
* `modules/lifecycle/internal/domain/enterprise_accountant.go`
* `modules/lifecycle/internal/repository/accountant_link_repo.go` (Acesso exclusivo ao Banco Central)
* `modules/lifecycle/internal/service/accountant_link_service.go`

**🛠️ Tarefas de Implementação**
1. **Modelagem de Domínio (`lifecycle`):** Crie a struct `EnterpriseAccountant` contendo `ID`, `EnterpriseID`, `AccountantID`, `Status` (ACTIVE/INACTIVE), `StartDate`, `EndDate` (*pointer* para time, nulo se ativo) e `DelegatedBy`.
2. **Repositório Central:** Implemente as funções no banco de dados **Central** (ex: `central.db` e não nos bancos dos tenants) para gravar o vínculo, atualizar o `EndDate` (quando inativado) e consultar vínculos por Contador ou por Entidade.
3. **Regra de Negócio (Exit Power):** No serviço correspondente, garanta que um Empreendimento só possa ter **1 (um) contador com status ACTIVE**. Se um novo vínculo for ativado, o anterior deve automaticamente ser inativado (preenchendo o `EndDate` com a data atual).
4. **Filtro de Acesso Temporal:** Crie uma função (ex: `GetValidDateRangeForAccountant(accountantID, enterpriseID)`) que retorne o `StartDate` e `EndDate` exatos do período em que aquele contador possui/possuiu responsabilidade técnica. Essa função será injetada posteriormente como middleware no `accountant_dashboard` para filtrar consultas.

**✅ Critérios de Aceite (Definition of Done)**
- [ ] A entidade `EnterpriseAccountant` foi criada e as migrações do SQLite apontam estritamente para a conexão do Banco Central, mantendo os bancos dos *Tenants* intactos.
- [ ] A regra de cardinalidade temporal impede que uma mesma cooperativa tenha dois contadores "ativos" ao mesmo tempo.
- [ ] O encerramento do vínculo grava a data atual no `EndDate` em vez de deletar o registro físico (Soft Delete / Histórico para auditoria).
- [ ] A arquitetura segue Clean Architecture e o princípio de Inversão de Dependência (DIP).

---
1. Código fonte seguindo Clean Architecture (Domain -> Service -> Handler).
2. Testes unitários com TDD provando a lógica.
3. Atualização sugerida para o próximo Session Log.
4. Você, agente, deve preencher todas as seções seguintes a partir do entendimento que teve do "objetivo"

Pode iniciar a análise e propor o plano de implementação?
```

---

## 📋 REQUISITOS

### Funcionais
- [ ] Requisito 1
- [ ] Requisito 2
- [ ] Requisito 3

### Técnicos
- [ ] Seguir padrões do projeto Digna
- [ ] Implementar testes unitários
- [ ] Atualizar documentação
- [ ] Validar com smoke tests

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
**Última atualização:** [DATA]