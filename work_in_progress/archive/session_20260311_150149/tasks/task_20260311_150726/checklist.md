# 📋 CHECKLIST PRÉ-IMPLEMENTAÇÃO: Feature (Interface do Painel do Contador e Exportação SPED)
# ID: 20260311_150726
# Gerado em: 

## 🔍 ANÁLISE DO CONTEXTO

### 1. Contexto do Projeto
- [x] Ler work_in_progress/current_session/.agent_context.md
- [x] Consultar docs/QUICK_REFERENCE.md
- [x] Verificar antipadrões em docs/ANTIPATTERNS.md
- [x] Revisar aprendizados anteriores em docs/learnings/
- [x] Carregar skills: developing-digna-backend, rendering-digna-frontend, auditing-fiscal-compliance, managing-sovereign-data

### 2. Análise de Similaridades
- [ ] Analisar handler existente: `modules/ui_web/internal/handler/accountant_handler.go`
- [ ] Analisar template similar: `modules/ui_web/templates/dashboard_simple.html`
- [ ] Verificar padrão SHA256: `core_lume/internal/domain/statute.go`
- [ ] Verificar padrão file download: `accountant_handler.go` (linhas 150-200)
- [ ] Analisar BaseHandler: `modules/ui_web/internal/handler/base_handler.go`
- [ ] Verificar TranslatorService: `modules/accountant_dashboard/`

### 3. Requisitos Técnicos
- [ ] Definir estrutura do handler (estender accountant_handler.go)
- [ ] Definir estrutura do template (accountant_dashboard_simple.html)
- [ ] Definir rotas HTTP: `/accountant/dashboard`, `/accountant/export/{entity_id}/{period}`
- [ ] Definir funções de template: formatCurrency, formatDate, getAlertStatusLabel
- [ ] Garantir acesso Read-Only ao SQLite (`?mode=ro`)
- [ ] Implementar validação "Soma Zero"
- [ ] Gerar hash SHA256 para exportação

## 🛠️ PREPARAÇÃO TÉCNICA

### 4. Estrutura de Arquivos
- [ ] Atualizar: modules/ui_web/internal/handler/accountant_handler.go
- [ ] Criar: modules/ui_web/templates/accountant_dashboard_simple.html
- [ ] Criar/atualizar: modules/ui_web/internal/handler/accountant_handler_test.go
- [ ] Atualizar: modules/ui_web/main.go (registrar rotas do accountant)

### 5. Dependências
- [ ] Verificar imports: lifecycle, accountant_dashboard, core_lume
- [ ] Verificar funções do lifecycle manager para acesso Read-Only
- [ ] Verificar integração com TranslatorService (accountant_dashboard)
- [ ] Verificar mapeamento de contas para SPED
- [ ] Verificar padrões de hash SHA256 do core_lume

## 🧪 TESTES E VALIDAÇÃO

### 6. Estratégia de Testes
- [ ] Testes unitários para handler (TDD)
- [ ] Testes de integração com TranslatorService
- [ ] Testes de acesso Read-Only ao SQLite
- [ ] Testes de geração de hash SHA256
- [ ] Testes de formatação SPED/CSV
- [ ] Critérios de aceitação: RF-11 completo, interface funcional, exportação válida

### 7. Validação Pós-Implementação
- [ ] Smoke test: ./scripts/dev/smoke_test_new_feature.sh
- [ ] Testes de sistema: cd modules && ./run_tests.sh
- [ ] Validação manual da UI (dashboard multi-tenant, exportação)
- [ ] Verificar paleta "Soberania e Suor"
- [ ] Testar acesso Read-Only (tentar operação de escrita)
- [ ] Validar hash SHA256 na exportação

## 📚 DOCUMENTAÇÃO

### 8. Documentação Técnica
- [ ] Atualizar docs/QUICK_REFERENCE.md (adicionar Painel do Contador)
- [ ] Atualizar docs/NEXT_STEPS.md (marcar RF-11 como concluído)
- [ ] Criar documentação da feature em docs/learnings/
- [ ] Documentar aprendizados (anti-padrões, boas práticas)
- [ ] Atualizar docs/ANTIPATTERNS.md se necessário

### 9. Checklist de Conclusão
- [ ] Todos os testes passando (149/149 + novos)
- [ ] Código revisado seguindo padrões (Anti-Float, Cache-Proof, Soberania)
- [ ] Documentação atualizada (QUICK_REFERENCE, NEXT_STEPS, aprendizados)
- [ ] Smoke test executado com sucesso
- [ ] Handler registrado no main.go
- [ ] Interface funcional com paleta "Soberania e Suor"
- [ ] Acesso Read-Only validado
- [ ] Exportação SPED/CSV com hash SHA256 funcionando

---

## 📝 NOTAS DA ANÁLISE

### Padrões Identificados:
- Handler existente accountant_handler.go (precisa ser estendido)
- Padrão SHA256 em core_lume/internal/domain/statute.go
- Padrão file download em accountant_handler.go (linhas 150-200)
- BaseHandler em modules/ui_web/internal/handler/base_handler.go
- Template cache-proof (_simple.html + ParseFiles no handler)

### Riscos Identificados:
1. Integração com TranslatorService - Verificar se está completo
2. Acesso Read-Only ao SQLite - Testar parâmetro ?mode=ro
3. Performance dashboard multi-tenant - Otimizar consultas
4. Formatação SPED/CSV - Validar com padrão fiscal

### Decisões de Design:
1. Estender accountant_handler.go existente em vez de criar novo
2. Usar BaseHandler para estrutura padrão
3. Implementar acesso Read-Only obrigatório para contador
4. Incluir hash SHA256 em todas as exportações

### Referências:
- Handler similar: accountant_handler.go (existente)
- Template similar: dashboard_simple.html
- Testes similares: accountant_handler_test.go
- Skills: developing-digna-backend, rendering-digna-frontend, auditing-fiscal-compliance, managing-sovereign-data

---

**Status do Checklist:** ✅ GERADO
**Próximo passo:** ./process_task.sh --task=20260311_150726 --plan
