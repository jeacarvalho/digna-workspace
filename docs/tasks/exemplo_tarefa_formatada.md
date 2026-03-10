# 📋 Implementar UI para Gestão de Fornecedores

**ID da Tarefa:** `TASK_EXEMPLO_001`  
**Data de Criação:** 10/03/2026  
**Prioridade:** □ Alta ☑️ Média □ Baixa  
**Estimativa:** `8` horas

---

## 🎯 **INFORMAÇÕES BÁSICAS**

### **Metadados (para extração automática):**
```
Tipo: Feature
Módulo: ui_web
Objetivo: Implementar UI para gestão completa de fornecedores
Decisões: Seguir padrão MemberHandler, cards com CNPJ opcional, validação anti-float, integração com módulo supply
```

---

## 📝 **DESCRIÇÃO DETALHADA**

### **Contexto:**
O módulo supply já tem backend para gestão de fornecedores, mas falta interface web. Usuários precisam cadastrar, listar e gerenciar fornecedores através da UI.

### **Escopo:**
- CRUD completo de fornecedores (Create, Read, Update, Delete)
- Listagem em cards com informações principais
- Busca e filtragem por nome/CNPJ
- **NÃO inclui:** Integração com sistemas externos, emissão de notas fiscais

### **Requisitos Funcionais:**
1. Cadastro de fornecedor com: nome, CNPJ (opcional), contato, endereço
2. Listagem paginada de fornecedores
3. Edição e exclusão de fornecedores
4. Busca por nome ou CNPJ
5. Validação de CNPJ quando informado

### **Requisitos Não-Funcionais:**
- **Performance:** Carregamento em < 2s com até 100 fornecedores
- **Usabilidade:** Interface intuitiva, similar a Members
- **Segurança:** Apenas usuários autenticados podem acessar
- **Compatibilidade:** Chrome, Firefox, Safari modernos

---

## 🏗️ **DESIGN TÉCNICO**

### **Arquitetura:**
- **Backend:** Serviço `SupplierService` em `modules/supply/`
- **Frontend:** Novo `SupplierHandler` em `ui_web`
- **Integrações:** Módulo supply para dados, sistema de templates existente

### **Componentes a Criar/Modificar:**
```
/modules/ui_web/internal/handler/supplier_handler.go
/modules/ui_web/templates/supplier_simple.html
/modules/ui_web/internal/handler/supplier_handler_test.go
/modules/ui_web/templates/dashboard_simple.html (adicionar link)
```

### **Dependências:**
- **Pré-requisitos:** Serviço `SupplierService` deve estar funcionando
- **Dependências externas:** Nenhuma
- **Impacto em outros módulos:** Adiciona link na navegação principal

---

## 🧪 **CRITÉRIOS DE ACEITE**

### **Testes Obrigatórios:**
- [ ] **Testes unitários:** Cobertura >90% para handler
- [ ] **Testes de integração:** Com banco SQLite real
- [ ] **Smoke test:** `./scripts/dev/smoke_test_new_feature.sh "Gestão de Fornecedores" "/suppliers"`
- [ ] **Validação E2E:** `./scripts/dev/validate_e2e.sh --basic --headless`

### **Validação Funcional:**
- [ ] Login e navegação para `/suppliers` funcionam
- [ ] CRUD completo via HTMX
- [ ] Validação de CNPJ funciona
- [ ] Busca e filtragem funcionam
- [ ] Design segue "Soberania e Suor"

### **Validação Técnica:**
- [ ] Anti-float compliance (zero `float`)
- [ ] Soberania mantida (dados só no `.sqlite` da entidade)
- [ ] Cache-proof (usa `ExecuteTemplate` do `BaseHandler`)
- [ ] Navegação unificada em templates principais

---

## 🔗 **REFERÊNCIAS**

### **Código de Referência:**
- **Handler similar:** `MemberHandler`
- **Template base:** `dashboard_simple.html`
- **Testes de referência:** `member_handler_test.go`

### **Documentação:**
- **Checklist pré-implementação:** `docs/templates/pre_implementation_checklist.md`
- **Plano de implementação:** `docs/templates/implementation_plan.md`
- **Template de tarefa:** `docs/templates/task_description.md`

---

## 📅 **PLANO DE EXECUÇÃO**

### **Fase 1: Análise (2 horas)**
- [ ] Preencher checklist pré-implementação
- [ ] Analisar `SupplierService` no backend
- [ ] Estudar padrão `MemberHandler`

### **Fase 2: Implementação (4 horas)**
- [ ] Criar `SupplierHandler` seguindo padrões
- [ ] Desenvolver template `supplier_simple.html`
- [ ] Implementar testes TDD

### **Fase 3: Integração (2 horas)**
- [ ] Adicionar link na navegação
- [ ] Testes de integração
- [ ] Validação E2E

---

## 🚨 **RISCOS IDENTIFICADOS**

| Risco | Probabilidade | Impacto | Mitigação |
|-------|--------------|---------|-----------|
| `SupplierService` não acessível | Média | Alto | Verificar acesso antes, mock inicial se necessário |
| Performance com muitos fornecedores | Baixa | Médio | Implementar paginação no template |
| Validação de CNPJ complexa | Média | Baixo | Usar biblioteca validadora existente |

---

## 📊 **MÉTRICAS ESPERADAS**

### **Qualidade de Código:**
- **Cobertura de testes:** >95%
- **Complexidade ciclomática:** < 8
- **Linhas de código:** ~300

### **Performance:**
- **Tempo de resposta:** < 1.5s
- **Uso de memória:** < 50MB adicional
- **Tamanho do template:** < 100 linhas

---

## 📝 **NOTAS ADICIONAIS**

O módulo supply já tem testes unitários para `SupplierService`. Precisamos apenas consumir a API corretamente.

CNPJ é opcional porque alguns fornecedores podem ser informais (produtores familiares).

---

## ✅ **CHECKLIST FINAL**

### **Antes de começar:**
- [ ] Checklist pré-implementação preenchido
- [ ] `SupplierService` analisado e compreendido
- [ ] Padrão `MemberHandler` estudado
- [ ] Riscos mapeados e mitigados

### **Após implementação:**
- [ ] Testes unitários passam (>90% cobertura)
- [ ] Smoke test passa
- [ ] **Validação E2E passa** (`validate_e2e.sh --basic --headless`)
- [ ] Link adicionado na navegação
- [ ] Aprendizados registrados no `conclude_task.sh`