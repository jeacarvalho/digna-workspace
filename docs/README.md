# 📚 Documentação Completa - Projeto Digna

**Última atualização:** 11/03/2026  
**Status do projeto:** Sprint 17 em andamento

---

## 🎯 VISÃO GERAL

O **Projeto Digna** é um sistema completo de economia solidária desenvolvido em Go, com foco em soberania de dados, transparência radical e autogestão. O sistema implementa todos os aspectos necessários para uma organização de economia solidária funcionar de forma autônoma e democrática.

### Princípios Fundamentais:
1. **Soberania de dados** - Cada entidade tem seu próprio banco SQLite
2. **Transparência radical** - Todas as transações são auditáveis publicamente
3. **Exit power** - Usuários sempre têm controle total sobre seus dados
4. **Anti-float** - Valores financeiros sempre em centavos (int64)
5. **Autogestão** - Decisões coletivas documentadas e auditáveis

---

## 📁 ESTRUTURA DA DOCUMENTAÇÃO

### 📋 Guias Essenciais (Comece aqui)
- **[QUICK_REFERENCE.md](QUICK_REFERENCE.md)** - Referência rápida de arquitetura e padrões
- **[ANTIPATTERNS.md](ANTIPATTERNS.md)** - Padrões problemáticos a evitar
- **[NEXT_STEPS.md](NEXT_STEPS.md)** - Backlog e próximos passos do projeto
- **[NOVA_ESTRUTURA_WORKFLOW.md](NOVA_ESTRUTURA_WORKFLOW.md)** - Novo fluxo de trabalho organizado

### 🚀 Implantação e Configuração
- **[deployment/QUICK_DEPLOY.md](deployment/QUICK_DEPLOY.md)** - Deploy rápido em 5 minutos

### 🔍 Análises e Decisões
- **[analysis/ANALISE_FLUXO_AGENTES.md](analysis/ANALISE_FLUXO_AGENTES.md)** - Análise do fluxo com agentes
- **[analysis/FLUXO_AGENTES_ATUALIZADO.md](analysis/FLUXO_AGENTES_ATUALIZADO.md)** - Fluxo atualizado
- **[analysis/IMPROVEMENTS_SUMMARY.md](analysis/IMPROVEMENTS_SUMMARY.md)** - Resumo de melhorias

### 📝 Conclusões de Tarefas
- **[task_conclusions/TASK_CONCLUSION.md](task_conclusions/TASK_CONCLUSION.md)** - Conclusão: bugs no fluxo de caixa
- **[task_conclusions/TASK_CONCLUSION_DOCKER.md](task_conclusions/TASK_CONCLUSION_DOCKER.md)** - Conclusão: Docker e produção

### 🧪 Testes e Qualidade
- **[testing/TEST_E2E_REPORT.md](testing/TEST_E2E_REPORT.md)** - Relatório de testes end-to-end

### 🎯 Prompts de Tarefas
- **[task_prompts/legal_dossier_generation.md](task_prompts/legal_dossier_generation.md)** - Geração de dossiê jurídico

### 📚 Aprendizados Consolidados
- **[learnings/](learnings/)** - Todos os aprendizados organizados por data

### 🛠️ Skills Específicas
- **[skills/](skills/)** - Conjuntos de instruções para tarefas específicas

### 📋 Templates
- **[templates/](templates/)** - Templates para desenvolvimento

### 📦 Arquivo Histórico
- **[archive/](archive/)** - Documentação histórica

---

## 🏗️ ARQUITETURA DO SISTEMA

### Módulos Principais:
```
modules/
├── core_lume/              # Domínio e serviços centrais
│   ├── internal/domain/    # Entidades de domínio
│   ├── internal/service/   # Serviços de aplicação
│   └── internal/repository # Repositórios e interfaces
├── ui_web/                 # Interface web
│   ├── internal/handler/   # Handlers HTTP
│   ├── templates/          # Templates HTML (*_simple.html)
│   └── main.go             # Ponto de entrada
├── legal_facade/           # Facade jurídica
│   ├── internal/document/  # Geração de documentos
│   └── pkg/document/       # API pública
└── [accountant_dashboard, budget, cash_flow, distribution, 
     integrations, pdv_ui, supply, sync_engine]/
```

### Padrões Técnicos:
1. **Handlers HTTP:** Estendem `BaseHandler`, usam HTMX para interatividade
2. **Templates:** `*_simple.html` com Tailwind CSS, cache-proof
3. **Banco de dados:** SQLite por entidade, sem JOINs entre bancos
4. **Validação:** SHA256 para integridade de documentos
5. **Testes:** >90% cobertura, testes unitários e E2E

---

## 🔄 FLUXO DE DESENVOLVIMENTO

### Nova Estrutura (Recomendada):
```
work_in_progress/
├── current_session/       # Sessão atual
├── tasks/                 # Tarefas em andamento
├── archive/               # Histórico completo
└── task_template/         # Template padronizado
```

### Sequência de Comandos:
```bash
# 1. Iniciar sessão
./start_session.sh [quick]

# 2. Criar tarefa
./create_task.sh "Nome da Tarefa" [módulo]

# 3. Processar (fluxo recomendado)
./process_task.sh --task=[ID] --checklist
./process_task.sh --task=[ID] --plan
./process_task.sh --task=[ID] --execute

# 4. Concluir
./conclude_task.sh --task=[ID] "Aprendizados" --success

# 5. Encerrar sessão
./end_session.sh
```

### Validações Obrigatórias:
- ✅ Handler registrado no `main.go` (para features UI)
- ✅ Testes unitários implementados
- ✅ Smoke test executado
- ✅ Documentação atualizada

---

## 🎨 DESIGN SYSTEM

### Paleta de Cores (Soberania e Suor):
- `#2A5CAA` - Azul soberania (botões principais)
- `#4A7F3E` - Verde suor (indicadores sucesso)
- `#F57F17` - Laranja energia (alertas/destaques)
- `#F9F9F6` - Fundo
- `#212121` - Texto

### Componentes UI:
- **Headers:** Logo Digna + navegação
- **Forms:** HTMX com validação em tempo real
- **Tables:** Responsivas com ações inline
- **Cards:** Para dashboards e resumos
- **Alerts:** Feedback contextualizado

---

## 🔒 SEGURANÇA E CONFORMIDADE

### Princípios de Segurança:
1. **Dados soberanos:** Cada entidade controla seus dados
2. **Auditoria completa:** Todas as transações com hash SHA256
3. **Transparência:** Documentos públicos auditáveis
4. **Consentimento:** Ações sempre com confirmação explícita

### Conformidade:
- **CADSOL:** Sistema compatível com Cadastro de Decisões Soberanas
- **SPED:** Exportação para sistemas fiscais governamentais
- **LGPD:** Privacidade por design, dados minimizados

---

## 🧪 TESTES E QUALIDADE

### Suíte de Testes:
- **Unitários:** >90% cobertura, focados em lógica de negócio
- **Integração:** Testes entre módulos
- **E2E:** Testes completos do sistema
- **Performance:** Testes de carga e stress

### Execução:
```bash
# Todos os testes
cd modules && ./run_tests.sh

# Testes específicos
go test -v ./modules/ui_web/...

# Smoke test
./scripts/dev/smoke_test_new_feature.sh "Descrição" "/rota"
```

### Relatórios:
- Ver `docs/testing/` para relatórios completos
- Métricas coletadas automaticamente pelo fluxo de trabalho

---

## 📈 ROADMAP E EVOLUÇÃO

### Sprint Atual (17):
- Melhorias na experiência do usuário
- Otimizações de performance
- Expansão da documentação
- Novas integrações

### Próximas Sprints:
- Ver `docs/NEXT_STEPS.md` para backlog detalhado
- Priorização baseada em feedback da comunidade

### Metas de Longo Prazo:
1. **Escalabilidade:** Suporte a milhares de entidades
2. **Internacionalização:** Múltiplos idiomas e moedas
3. **Mobile:** Aplicativo nativo para dispositivos móveis
4. **Federação:** Comunicação entre instâncias Digna

---

## 🤝 CONTRIBUIÇÃO

### Como Contribuir:
1. **Estudar a documentação:** Comece por `QUICK_REFERENCE.md`
2. **Seguir o fluxo:** Use a nova estrutura `work_in_progress/`
3. **Respeitar padrões:** Consulte `ANTIPATTERNS.md`
4. **Documentar:** Sempre use `conclude_task.sh` para aprendizados

### Canais:
- **Issues:** Problemas e sugestões
- **Documentação:** Melhorias na docs/
- **Código:** Pull requests seguindo padrões

### Reconhecimento:
- Todos os contribuidores listados na documentação
- Aprendizados compartilhados publicamente
- Transparência total no processo

---

## 📞 SUPORTE

### Problemas Comuns:
- **Configuração:** Ver `deployment/QUICK_DEPLOY.md`
- **Desenvolvimento:** `QUICK_REFERENCE.md` e `skills/`
- **Fluxo:** `NOVA_ESTRUTURA_WORKFLOW.md`

### Recursos Adicionais:
- **Aprendizados:** `learnings/` - Lições de implementações anteriores
- **Análises:** `analysis/` - Decisões arquiteturais documentadas
- **Templates:** `templates/` - Modelos para acelerar desenvolvimento

### Contato:
- Issues no repositório oficial
- Documentação como fonte primária de verdade
- Comunidade de economia solidária

---

**Nota:** Esta documentação é atualizada continuamente com os aprendizados de cada sessão de desenvolvimento. Consulte `docs/learnings/` para as lições mais recentes.

**Licença:** Economia Solidária - Uso livre para comunidades, cooperativas e organizações de autogestão.