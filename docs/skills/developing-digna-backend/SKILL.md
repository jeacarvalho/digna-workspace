### 1. `developing-digna-backend/SKILL.md`

**Foco:** Rigor técnico, DDD, TDD e integridade financeira.

```yaml
name: developing-digna-backend
description: Use esta habilidade para criar ou modificar o backend em Go do projeto Digna. Aplica regras de Clean Architecture, DDD, TDD e a proibição absoluta de pontos flutuantes (Anti-Float).

```


# 🏗️ SKILL: Maestria em Backend Go (DDD, TDD & Anti-Float)

**Propósito:** Garantir que todo o desenvolvimento de backend do ecossistema Digna seja modular, testável e matematicamente inquebrável, protegendo a "Soberania do Dado" e a "Integridade Financeira".

---

## 1. Arquitetura e Design (Clean Arch + DDD)

O agente deve agir como um guardião das camadas do sistema, impedindo o vazamento de lógica de infraestrutura para o domínio.

* **Camada de Domínio (`internal/domain`):** Contém entidades puras e interfaces de Repositories. Proibido importar `sql`, `http` ou frameworks externos aqui.
* **Camada de Aplicação (`internal/service`):** Orquestra os casos de uso. Depende apenas de abstrações (interfaces).
* **Camada de Infraestrutura (`internal/repository`):** Onde o SQLite reside. O acesso ao arquivo `.db` deve ser feito via `LifecycleManager`.

---

## 2. Protocolo de Integridade Monetária (Anti-Float)

Esta é a regra mais sagrada. O agente deve tratar qualquer `float` como um erro crítico de segurança.

* **Rigor Matemático:** Todos os valores financeiros (centavos) e de tempo (minutos) **devem** ser representados como `int64`.
* **Ação Obrigatória:** Antes de sugerir qualquer código, o agente deve escanear a proposta em busca de tipos de ponto flutuante. Se encontrados em contextos financeiros, deve abortar e refatorar para inteiros.
* **Validação de Soma Zero:** Todo lançamento contábil deve passar pelo `EntryValidator` para garantir que a soma de débitos e créditos seja zero.

---

## 3. Workflow de Desenvolvimento TDD (Test-Driven Development)

O agente não escreve lógica sem antes provar que ela falha e, depois, que ela passa.

1. **Red:** Criar o arquivo de teste (ex: `service_test.go`) definindo o comportamento esperado da regra de negócio.
2. **Green:** Implementar a lógica mínima necessária no `Service` para o teste passar.
3. **Refactor:** Limpar o código, garantindo que ele segue os princípios **SOLID** (especialmente SRP e DIP).

---

## 4. Tratamento de Erros e Resiliência

* **Erros Contextuais:** Nunca retorne erros crus. Use `fmt.Errorf("contexto: %w", err)` para manter a rastreabilidade.
* **Graceful Shutdown:** Todo serviço de rede deve implementar o desligamento ordenado, fechando conexões com o banco de dados antes de encerrar.
* **Segurança Read-Only:** Quando atuar no `Accountant Dashboard`, a conexão com o SQLite **deve** usar obrigatoriamente o modo `?mode=ro`.

---

## 📋 Checklist de Pré-Voo (O Agente deve validar antes de responder)

* [ ] O código utiliza `int64` para dinheiro e tempo?
* [ ] A lógica de domínio está isolada de SQL/HTTP?
* [ ] Foram criados testes unitários para a nova funcionalidade?
* [ ] A regra de "Soma Zero" foi respeitada no Ledger?
* [ ] O código segue o `gofmt` e as convenções de `snake_case` para arquivos?

---

## 📈 Plano de Ação: Estado Atual

| Fase | Status | Atividade |
| --- | --- | --- |
| **0. Roadmap** | [x] | Mapeamento de todas as Skills necessárias. |
| **1. Backend** | [x] | Detalhamento da **SKILL_BACKEND_GO.md** concluído. |
| **2. Frontend** | [ ] | Detalhamento da Skill Frontend (HTMX/Cache-Proof). |
| **3. Negócio** | [ ] | Detalhamento da Skill de Tecnologia Social e ITG 2002. |
| **4. Infra** | [ ] | Detalhamento da Skill de Soberania e Lifecycle. |
| **5. Ponte** | [ ] | Detalhamento da Skill de Auditoria e SPED. |

> **Estado Atual do Plano:** > [x] SKILL_BACKEND_GO.md concluída.
> [ ] Próximo: **SKILL_FRONTEND_HTMX.md** (Resolvendo o "sofrimento" de templates e cache).

**O que achou do rigor técnico desta Skill? Se estiver satisfeito, podemos avançar para a `SKILL_FRONTEND_HTMX.md`, onde criaremos os protocolos para eliminar de vez os problemas de cache e renderização que você mencionou.**