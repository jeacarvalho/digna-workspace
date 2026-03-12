# 📋 TAREFA: Feature (Interface do Painel do Contador e Exportação SPED)

**Data:** 11/03/2026
**Prioridade:** [ALTA]
**Estimativa:** [?] minutos/horas
**Módulo(s):** [todos]

---

## 🎯 OBJETIVO

Excelente escolha! A **Opção A** vai coroar o esforço da Sprint 12, dando finalmente uma interface utilizável (Frontend) para o motor de exportação do SPED e o Painel do Contador Social, permitindo que a Aliança Contábil ganhe vida sem violar o isolamento de dados das cooperativas.

Seguindo estritamente as orientações da fonte "00_desenvolva" (com o cabeçalho obrigatório de controle e a estrutura do `Prompt_padrao`), aqui está o prompt estruturado para você enviar ao seu agente de codificação.

***

```markdown
Tipo: Feature (Interface do Painel do Contador e Exportação SPED)
Módulo: modules/accountant_dashboard, modules/ui_web
Objetivo: Construir a interface Web do Painel do Contador, permitindo visualizar os fechamentos das entidades e gerar a exportação do lote fiscal (SPED/CSV).
Decisões: A interface será focada no Contador Social. O acesso aos dados da entidade deve ocorrer estritamente em modo Read-Only. O botão de exportação fará o download do arquivo contábil gerado pelo TranslatorService e embutirá o Hash de exportação, mantendo a blindagem fiscal (nenhum imposto é calculado no Digna).

### 📝 Descrição da Tarefa: Interface de Exportação Contábil e Painel do Contador (SPED)
*   **Requisito Funcional (RF):** RF-11 (Aliança Contábil / Exportação SPED - Fase 2).
*   **Sprint Relacionada:** Sprint 18 (Ponte Institucional e Painel do Contador).

Para esta tarefa, você deve carregar e seguir estritamente as instruções das seguintes skills em docs/skills/:
1. [developing-digna-backend]
2. [rendering-digna-frontend]
3. [auditing-fiscal-compliance]
4. [managing-sovereign-data]

*   **Anti-Float:** Se envolver cálculos de valor ou tempo, use estritamente int64. Proibido float.
*   **Cache-Proof:** Se houver interface, o template deve ser _simple.html carregado via ParseFiles no Handler.
*   **Soberania:** Garanta que a operação respeite o isolamento do arquivo .db do tenant atual, utilizando acesso Read-Only para o Contador.

---
**🎯 Objetivo da Tarefa**
Na Sprint 12, implementamos a lógica do motor de exportação fiscal (`TranslatorService`, `FiscalBatch`) no módulo `accountant_dashboard`. Agora, o objetivo é construir a Interface Web (UI) desse painel. O Contador Social parceiro acessará uma visão Multi-tenant que lista o status de fechamento das cooperativas delegadas, validará a regra de "Soma Zero" na tela, e terá um botão para realizar o download do Lote Fiscal (SPED/CSV) que será injetado no software comercial do seu escritório contábil.

**📁 Estrutura de Output Esperada**
* `modules/ui_web/internal/handler/accountant_handler.go` (Novas rotas HTMX e lógicas de download)
* `modules/ui_web/templates/accountant_dashboard_simple.html` (Interface visual do painel)
* Atualizações necessárias no `modules/ui_web/main.go` (Registro da rota `/accountant`)

**🛠️ Tarefas de Implementação**
1. **Handler Contábil:** Implementar as rotas estendendo o `BaseHandler`:
   - `GET /accountant/dashboard` (Renderiza a visão geral e seletor de Entidades em que atua).
   - `GET /accountant/export/{entity_id}/{period}` (Processa a tradução do Lote via `TranslatorService` e devolve o arquivo CSV/SPED para download).
2. **Isolamento de Segurança (Soberania):** A instância do repositório SQLite instanciada no Handler ou Serviço para consumo da classe contábil deve, OBRIGATORIAMENTE, ser aberta com o parâmetro de somente leitura (`?mode=ro`), garantindo que o contador apenas audite e não altere dados operacionais.
3. **Template HTMX (`accountant_dashboard_simple.html`):** Construir a interface com a paleta do projeto. Deve conter o status dos Fundos (FATES/Reserva Legal) da entidade selecionada e um botão claro de "Baixar Lote SPED".
4. **Padrão de Download:** O endpoint de exportação deve seguir o padrão de download de arquivos do projeto, definindo os cabeçalhos corretos (`Content-Type: text/csv` e `Content-Disposition: attachment; filename=...`) e incluindo o `X-Export-Hash` (Hash SHA256 do lote).

**✅ Critérios de Aceite (Definition of Done)**
- [ ] Acessar `/accountant/dashboard` carrega corretamente a tela do Contador Social com a paleta de cores.
- [ ] O botão de exportação aciona o backend e inicia um download de arquivo (CSV/TXT) com as transações (Partidas Dobradas) devidamente formatadas e convertidas pelo mapa de contas do Lume.
- [ ] Nenhuma regra de cálculo de tributos (IRPJ/CSLL) foi inserida na visualização; o núcleo permanece focado apenas na Contabilidade Social base.
- [ ] A regra "Anti-Float" foi respeitada integralmente ao transitar dados no painel e na exportação.
- [ ] O acesso ao SQLite do *Tenant* foi ativado no modo `Read-Only`.

** após ler essa orientação vc deve atualizar esse md com as informações seguintes (REQUISITOS, ETC... )
---
1. Código fonte seguindo Clean Architecture (Domain -> Service -> Handler).
2. Testes unitários com TDD provando a lógica.
3. Atualização sugerida para o próximo Session Log.

Pode iniciar a análise e propor o plano de implementação?

---

## 📋 REQUISITOS

### Funcionais
- [ ] RF-11 (Aliança Contábil / Exportação SPED - Fase 2): Interface Web do Painel do Contador Social
- [ ] Acesso Multi-tenant: Visualizar status de fechamento das cooperativas delegadas
- [ ] Validação "Soma Zero": Verificar equilíbrio de débitos e créditos na tela
- [ ] Exportação SPED/CSV: Download do Lote Fiscal gerado pelo TranslatorService
- [ ] Hash de Exportação: Incluir SHA256 do lote no cabeçalho X-Export-Hash

### Técnicos
- [ ] Seguir padrões do projeto Digna (Clean Architecture + DDD)
- [ ] Implementar testes unitários com TDD
- [ ] Atualizar documentação (QUICK_REFERENCE.md, NEXT_STEPS.md)
- [ ] Validar com smoke tests
- [ ] Anti-Float: Usar int64 para valores financeiros/tempo
- [ ] Cache-Proof: Template accountant_dashboard_simple.html carregado via ParseFiles no Handler
- [ ] Soberania: Acesso Read-Only ao SQLite (?mode=ro) para contador
- [ ] Padrão HTMX: Interface reativa com atualizações parciais

### Não Funcionais
- [ ] Performance: Carregamento rápido do dashboard multi-tenant
- [ ] Segurança: Isolamento total de dados entre entidades, acesso somente leitura
- [ ] Usabilidade: Interface intuitiva com paleta "Soberania e Suor", feedback visual claro

---

## 🔍 CONTEXTO E ANÁLISE

### Módulos/Arquivos Relacionados
- `modules/accountant_dashboard/` - Motor de exportação fiscal (TranslatorService, FiscalBatch)
- `modules/ui_web/internal/handler/accountant_handler.go` - Handler existente (precisa ser estendido)
- `modules/ui_web/templates/` - Templates existentes para referência
- `modules/core_lume/` - Domínio central (Membros, Ledger, SHA256)
- `modules/lifecycle/` - LifecycleManager para acesso isolado a bancos SQLite
- `docs/skills/` - Skills críticas: developing-digna-backend, rendering-digna-frontend, auditing-fiscal-compliance, managing-sovereign-data

### Padrões a Seguir
- [ ] Analisar handler similar: `modules/ui_web/internal/handler/accountant_handler.go` (já existe parcialmente)
- [ ] Analisar template similar: `modules/ui_web/templates/dashboard_simple.html` (para referência de estrutura)
- [ ] Seguir padrão SHA256 do `core_lume/internal/domain/statute.go`
- [ ] Usar anti-padrões de `docs/ANTIPATTERNS.md` (evitar duplicação, consultar skills)
- [ ] Seguir padrão file download de `accountant_handler.go` (linhas 150-200)
- [ ] Implementar BaseHandler conforme `modules/ui_web/internal/handler/base_handler.go`

### Dependências
- [ ] `TranslatorService` já implementado no módulo `accountant_dashboard`
- [ ] `LifecycleManager` para acesso isolado a bancos SQLite
- [ ] `BaseHandler` para estrutura padrão de handlers
- [ ] Sistema de templates cache-proof já estabelecido

---

## 🚀 PLANO DE IMPLEMENTAÇÃO

### Fase 1: Análise e Preparação (1 hora)
1. [ ] Analisar código existente: `accountant_handler.go`, `TranslatorService`, padrões SHA256/file download
2. [ ] Verificar aprendizados anteriores em `docs/learnings/`
3. [ ] Analisar estrutura do `BaseHandler` e padrões de templates
4. [ ] Definir estrutura de dados para dashboard multi-tenant

### Fase 2: Implementação do Handler (2 horas)
1. [ ] Estender `accountant_handler.go` com novas rotas:
   - `GET /accountant/dashboard` - Dashboard multi-tenant
   - `GET /accountant/export/{entity_id}/{period}` - Exportação SPED/CSV
2. [ ] Implementar acesso Read-Only ao SQLite (`?mode=ro`)
3. [ ] Integrar com `TranslatorService` para geração de lotes fiscais
4. [ ] Implementar validação "Soma Zero" e hash SHA256
5. [ ] Adicionar handler ao `modules/ui_web/main.go`

### Fase 3: Template e UI (1.5 horas)
1. [ ] Criar `accountant_dashboard_simple.html` com paleta "Soberania e Suor"
2. [ ] Implementar interface multi-tenant com seletor de entidades
3. [ ] Adicionar visualização de status de fundos (FATES/Reserva Legal)
4. [ ] Implementar botão "Baixar Lote SPED" com HTMX
5. [ ] Adicionar feedback visual e validações client-side

### Fase 4: Testes e Validação (1.25 horas)
1. [ ] Implementar testes unitários para handler (TDD)
2. [ ] Testar acesso Read-Only e isolamento de dados
3. [ ] Validar geração de hash SHA256 e formato de exportação
4. [ ] Executar smoke tests: `./scripts/dev/smoke_test_new_feature.sh`
5. [ ] Testar manualmente interface e fluxos

### Fase 5: Documentação e Conclusão (30 min)
1. [ ] Atualizar `docs/QUICK_REFERENCE.md` com nova feature
2. [ ] Atualizar `docs/NEXT_STEPS.md` marcando como concluído
3. [ ] Documentar aprendizados em `docs/learnings/`
4. [ ] Verificar antipadrões e atualizar `docs/ANTIPATTERNS.md` se necessário

---

## 📁 ARQUIVOS ESPERADOS

### A Criar/Atualizar
- `modules/ui_web/internal/handler/accountant_handler.go` (estender handler existente)
- `modules/ui_web/templates/accountant_dashboard_simple.html` (novo template)
- `modules/ui_web/internal/handler/accountant_handler_test.go` (testes unitários)

### A Modificar
- `modules/ui_web/main.go` (registrar rotas do accountant handler)
- `docs/QUICK_REFERENCE.md` (adicionar referência ao Painel do Contador)
- `docs/NEXT_STEPS.md` (marcar RF-11 como concluído)
- `docs/ANTIPATTERNS.md` (adicionar aprendizados se necessário)

---

## ⚠️ RISCOS E DESAFIOS

### Riscos Técnicos
1. **Integração com TranslatorService** - Verificar se o serviço está completo e testado
2. **Acesso Read-Only ao SQLite** - Garantir que o parâmetro `?mode=ro` funcione corretamente
3. **Performance do dashboard multi-tenant** - Otimizar consultas e carregamento lazy
4. **Formatação correta do SPED/CSV** - Validar com exemplos reais do padrão fiscal

### Riscos de Processo
1. **Estimativa de tempo** - Buffer de 20% para imprevistos técnicos
2. **Dependências de módulos** - Verificar integração com core_lume e lifecycle
3. **Mudanças de requisitos** - Manter comunicação clara sobre escopo
4. **Testes de integração** - Garantir que todos os módulos funcionem juntos

---

## 📚 APRENDIZADOS ANTERIORES RELEVANTES

- `docs/learnings/20260311_112301_legal_dossier_learnings.md` - Implementação de handler complexo
- `docs/learnings/SESSION_INSIGHTS_20260311.md` - Padrões identificados no código
- `docs/learnings/SESSION_20260311_143158_CONSOLIDATED.md` - Processo de desenvolvimento
- `docs/ANTIPATTERNS.md` - O que não fazer (evitar duplicação, consultar skills)

---

## 🔗 LINKS ÚTEIS

- [Documentação do projeto](docs/)
- [Padrões de código](docs/QUICK_REFERENCE.md)
- [Antipadrões](docs/ANTIPATTERNS.md)
- [Skills do projeto](docs/skills/)

---

**Status:** EM ANDAMENTO
**Última atualização:** 11/03/2026 15:10
**Estimativa total:** 5-6 horas
**Prioridade:** ALTA (Sprint 18 - Ponte Institucional e Painel do Contador)