# 🏛️ Projeto Digna - Economia Solidária

**Status:** Sprint 17 em andamento  
**Última atualização:** 11/03/2026  
**Documentação completa:** [docs/README.md](docs/README.md)

---

## 🚀 COMEÇAR AQUI

### Para Desenvolvedores/Contribuidores:
```bash
# 1. Clonar e configurar
git clone [repositório]
cd digna-workspace

# 2. Iniciar sessão de desenvolvimento (com o uso de agentes. Tenho usado opencode)
orientar o agente: run ./start_session.sh

# 3. Criar primeira tarefa (✅ link simbólico)
run ./create_task.sh "Nome da Tarefa" [módulo]

# 4. Seguir fluxo completo (recomendado)
run ./process_task.sh --task=[ID] --checklist    # ✅ link
run ./process_task.sh --task=[ID] --plan         # ✅ link  
run ./process_task.sh --task=[ID] --execute      # ✅ link
run ./conclude_task.sh --task=[ID] "Aprendizados" # ✅ link
run ./end_session.sh                             # ✅ link

# 5. Executar testes (✅ link simbólico)
run ./run_tests.sh
```

### Para Usuários/Implementadores:
```bash
# Ver documentação de deploy
cat docs/deployment/QUICK_DEPLOY.md

# Executar testes
cd modules && ./run_tests.sh

# Iniciar sistema
cd modules/ui_web && go run main.go
```

---

## 📁 ESTRUTURA DO PROJETO

```
digna-workspace/
├── modules/                    # Módulos Go do sistema
│   ├── core_lume/             # Domínio e serviços centrais
│   ├── ui_web/                # Interface web (HTMX + Go)
│   ├── legal_facade/          # Facade jurídica (✅ dossiês, atas, etc.)
│   └── [outros módulos]/
├── docs/                      # ✅ Documentação organizada
│   ├── README.md              # Documentação principal
│   ├── deployment/            # Guias de deploy
│   ├── analysis/              # Análises e estudos
│   ├── task_conclusions/      # Conclusões de tarefas
│   ├── testing/               # Relatórios de testes
│   ├── task_prompts/          # Prompts de tarefas
│   ├── learnings/             # Aprendizados consolidados
│   ├── skills/                # Skills específicas do projeto
│   └── templates/             # Templates para desenvolvimento
├── work_in_progress/          # ✅ Trabalho em andamento
│   ├── current_session/       # Sessão atual
│   ├── tasks/                 # Tarefas ativas
│   ├── archive/               # Histórico completo
│   └── task_template/         # Template padronizado
├── scripts/                   # Scripts de automação
│   ├── workflow/              # Fluxo de trabalho (✅ links simbólicos na raiz)
│   ├── tools/                 # Ferramentas auxiliares
│   ├── dev/                   # Desenvolvimento (✅ run_tests.sh link)
│   ├── deployment/            # Configuração e deploy
│   └── testing/               # Testes e qualidade
└── data/                      # Dados (não versionado)
```

---

## 🎯 STATUS ATUAL

### ✅ Funcionalidades Implementadas:
- **Sistema completo de economia solidária**
- **Gestão de membros e trabalho**
- **Rateio social proporcional**
- **PDV (Ponto de Venda)**
- **Fluxo de caixa**
- **Orçamento e planejamento**
- **Facade jurídica** (dossiês, atas, formalização)
- **Interface web com HTMX**
- **Testes unitários e E2E**

### 🔄 Em Desenvolvimento (Sprint 17):
- Melhorias na experiência do usuário
- Otimizações de performance
- Documentação expandida
- Integrações adicionais

### 📈 Próximos Passos:
- Ver `docs/NEXT_STEPS.md` para backlog completo
- Consultar `docs/learnings/` para aprendizados recentes

---

## 🔧 TECNOLOGIAS

### Backend:
- **Go** (Golang) - Linguagem principal
- **SQLite** - Banco de dados por entidade (soberania de dados)
- **HTMX** - Interatividade no frontend
- **Tailwind CSS** - Estilização

### Padrões Arquiteturais:
- **Domain-Driven Design** (DDL adaptado)
- **Clean Architecture** (camadas separadas)
- **Event Sourcing** (para auditoria)
- **CQRS** (leitura/escrita separadas onde aplicável)

### Princípios de Design:
- **Soberania de dados** (um banco por entidade)
- **Anti-float** (nunca usar float para valores financeiros)
- **Exit power** (usuário sempre tem controle)
- **Transparência radical** (tudo auditável)

---

## 📚 DOCUMENTAÇÃO COMPLETA

### Guias Principais:
- **[docs/README.md](docs/README.md)** - Documentação detalhada
- **[docs/QUICK_REFERENCE.md](docs/QUICK_REFERENCE.md)** - Referência rápida
- **[docs/ANTIPATTERNS.md](docs/ANTIPATTERNS.md)** - O que NÃO fazer
- **[docs/NEXT_STEPS.md](docs/NEXT_STEPS.md)** - Backlog e próximos passos

### Específicos:
- **Deploy:** `docs/deployment/` - Guias de implantação
- **Análises:** `docs/analysis/` - Estudos e decisões
- **Testes:** `docs/testing/` - Relatórios e guias
- **Aprendizados:** `docs/learnings/` - Lições do desenvolvimento

### Fluxo de Trabalho:
- **Nova estrutura:** `docs/NOVA_ESTRUTURA_WORKFLOW.md`
- **Scripts:** `scripts/workflow/` - Automação completa

---

## 🤝 CONTRIBUIÇÃO

### Processo de Desenvolvimento:
1. **Sessão:** `./start_session.sh` (inicia contexto)
2. **Tarefa:** `./create_task.sh` (cria tarefa padronizada)
3. **Implementação:** Seguir padrões em `docs/QUICK_REFERENCE.md`
4. **Testes:** Executar `cd modules && ./run_tests.sh`
5. **Documentação:** `./conclude_task.sh` (aprendizados)

### Padrões Importantes:
- ✅ **Seguir antipadrões** em `docs/ANTIPATTERNS.md`
- ✅ **Usar analyze_patterns.sh** para referências
- ✅ **Implementar testes** para novas funcionalidades
- ✅ **Documentar aprendizados** ao concluir

### Canais de Apoio:
- Issues no repositório
- Documentação em `docs/`
- Aprendizados consolidados em `docs/learnings/`

---

## 📞 SUPORTE E CONTATO

### Problemas Comuns:
- **Configuração:** Ver `docs/deployment/QUICK_DEPLOY.md`
- **Desenvolvimento:** Consultar `docs/QUICK_REFERENCE.md`
- **Fluxo:** Ler `docs/NOVA_ESTRUTURA_WORKFLOW.md`

### Links Úteis:
- [Documentação completa](docs/)
- [Referência rápida](docs/QUICK_REFERENCE.md)
- [Próximos passos](docs/NEXT_STEPS.md)
- [Aprendizados recentes](docs/learnings/)

---

**Nota:** Esta é uma visão geral. Para documentação completa, acesse [docs/README.md](docs/README.md).

**Licença:** Economia Solidária - Uso livre para comunidades e cooperativas.