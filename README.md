# 🏛️ Projeto Digna - Economia Solidária

**Status:** Sprint 16 completa, Sprint 17 em andamento  
**Arquitetura:** Clean Architecture + DDD + HTMX  
**Banco:** SQLite isolado por entidade (cooperativa)  
**Última atualização:** 10/03/2026

---

## 📁 ESTRUTURA DO PROJETO (REORGANIZADA)

```
digna-workspace/
├── 📁 modules/                 # Código fonte Go
│   ├── core_lume/             # Core business logic
│   ├── ui_web/                # Web interface (HTMX + Tailwind)
│   └── [outros módulos]/
├── 📁 docs/                   # Documentação completa
│   ├── learnings/             # Aprendizados por tarefa
│   ├── templates/             # Templates para workflow
│   ├── antipatterns/          # Antipadrões e soluções
│   └── implementation_plans/  # Planos de implementação
├── 📁 scripts/                # Scripts organizados
│   ├── workflow/              # Fluxo de trabalho principal
│   │   ├── start_session.sh   # Iniciar sessão
│   │   ├── process_task.sh    # Processar tarefa
│   │   ├── conclude_task.sh   # Concluir tarefa
│   │   └── end_session.sh     # Encerrar sessão
│   ├── dev/                   # Desenvolvimento
│   │   ├── run_tests.sh       # Executar testes
│   │   ├── smoke_test_new_feature.sh
│   │   └── init_test_companies.sh
│   └── tools/                 # Ferramentas auxiliares
├── 📁 data/                   # Dados (banco SQLite por entidade)
├── 📁 tmp/                    # Arquivos temporários (.gitignore)
├── 📁 backups/                # Backups de sessões (.gitignore)
├── 📄 README.md               # Este arquivo
├── 📄 Makefile                # Build automation
├── 📄 go.work                 # Go workspace
└── 📄 .gitignore              # Arquivos ignorados pelo Git
```

---

## 🚀 COMEÇAR A TRABALHAR

### **Fluxo 1: Com arquivo de prompt/roadmap** (RECOMENDADO)
```bash
# 1. Iniciar sessão
./start_session.sh          # Modo completo
./start_session.sh quick    # Modo rápido

# 2. Preparar implementação a partir de arquivo
./prepare_implementation.sh prompts/suppliers-ui.md --execute

# 3. O script gera:
#    - Checklist pré-implementação
#    - Plano de implementação  
#    - Prompt final para opencode

# 4. Copie o prompt (.opencode_task_*.txt) e cole no opencode

# 5. Implemente seguindo as instruções

# 6. Validação pós-implementação (OBRIGATÓRIA)
./scripts/dev/smoke_test_new_feature.sh "Suppliers" "/suppliers"

# 7. Concluir tarefa
./conclude_task.sh "Suppliers implementado. Aprendi: ..."

# 8. Encerrar sessão (opcional)
./end_session.sh
```

### **Fluxo 2: Direto com descrição**
```bash
# 1. Iniciar sessão
./start_session.sh

# 2. Processar tarefa com descrição
./process_task.sh "Tipo: Feature | Módulo: ui_web | Objetivo: Implementar X | Decisões: seguir padrão Y"

# Opções:
./process_task.sh "descrição" --checklist    # Apenas checklist
./process_task.sh "descrição" --plan         # Checklist + plano
./process_task.sh "descrição" --execute      # Tudo + prompt para opencode

# 3. Implementar no OpenCode
# Copie o prompt gerado e cole no opencode

# 4. Validação pós-implementação (OBRIGATÓRIA)
./scripts/dev/smoke_test_new_feature.sh "Nome da Feature" "/rota"

# 5. Concluir tarefa
./conclude_task.sh "O que aprendi com esta implementação"
```

---

## 🛠️ DESENVOLVIMENTO

### **Executar Testes**
```bash
./run_tests.sh              # Todos os testes
cd modules && ./run_tests.sh # Alternativa
```

### **Iniciar Servidor**
```bash
cd modules/ui_web && go run .
# Acesse: http://localhost:8090
```

### **Smoke Test para Novas Features**
```bash
./scripts/dev/smoke_test_new_feature.sh "Member Management" "/members"
```

---

## 📚 DOCUMENTAÇÃO

### **Documentos Principais**
- `docs/QUICK_REFERENCE.md` - Referência rápida de padrões
- `docs/README_OPENCODE_WORKFLOW.md` - Fluxo completo do opencode
- `docs/EXAMPLE_USAGE.md` - Exemplos de uso

### **Templates**
- `docs/templates/pre_implementation_checklist.md` - Checklist pré-implementação
- `docs/templates/implementation_plan.md` - Template de plano
- `docs/templates/post_correction_validation.md` - Validação pós-correção

### **Aprendizados**
- `docs/learnings/` - Aprendizados documentados por tarefa
- `docs/antipatterns/` - Antipadrões e soluções

---

## 🏗️ ARQUITETURA

### **Princípios Fundamentais**
1. **Soberania do Dado** - Banco SQLite isolado por entidade
2. **Anti-Float** - Zero `float` para valores financeiros/tempo
3. **Clean Architecture** - Domínio independente de frameworks
4. **Contabilidade Invisível** - Operações geram lançamentos automáticos

### **Módulos Principais**
- **core_lume/** - Lógica de negócio central
- **ui_web/** - Interface web (HTMX + Tailwind)
- **cash_flow/** - Gestão de fluxo de caixa
- **supply/** - Compras e estoque
- **lifecycle/** - Gerenciamento de ciclo de vida

### **Handlers Registrados**
- AuthHandler, DashboardHandler, PDVHandler, CashHandler
- SupplyHandler, BudgetHandler, MemberHandler, AccountantHandler

---

## ✅ SISTEMA DE VALIDAÇÃO 4-NÍVEIS

### **Nível 1: Testes Unitários**
- Lógica pura do handler
- Mock de dependências
- Cobertura >90%

### **Nível 2: Testes de Integração**
- Banco real (SQLite)
- Templates carregados
- Dependências injetadas

### **Nível 3: Testes de Sistema**
- Handler registrado no servidor
- Rotas respondem HTTP 200
- Templates compilam sem erro

### **Nível 4: Smoke Test Local**
- Testa feature no ambiente REAL
- Valida servidor, rotas, templates
- Output com ações corretivas

---

## 🔧 SISTEMA DE PREVENÇÃO DE BUGS

### **Problemas Comuns Prevenidos:**
1. **Handler não registrado** - Checklist obrigatório + validação automática
2. **Template não carregado** - Validação de compatibilidade
3. **Testes passam mas app quebra** - Smoke test obrigatório

### **Mecanismos Implementados:**
- ✅ Validação no `start_session.sh`
- ✅ Checklist expandido com itens críticos
- ✅ Validação obrigatória no `conclude_task.sh`
- ✅ Testes de sistema (`TestSystem_*`)
- ✅ Smoke test script

---

## 🎯 MISSÃO DO PROJETO

> Promover a autogestão, soberania e transformação digital dos Empreendimentos de Economia Solidária no Brasil através de tecnologia livre e acessível, atuando simultaneamente como uma **ponte tecnológica inclusiva para a conformidade legal e a classe contábil**.

---

## 📞 CONTATO E SUPORTE

### **Para ajuda com opencode:**
- Reportar issues: https://github.com/anomalyco/opencode/issues
- Documentação: `docs/README_OPENCODE_WORKFLOW.md`

### **Próximos Passos Sugeridos:**
1. Revisar `docs/NEXT_STEPS.md`
2. Escolher tarefa do backlog
3. Seguir fluxo automatizado

---

**Status:** 🟢 **PRODUCTION READY** - Sistema 100% funcional com workflow automatizado