# 📋 Template de Descrição de Tarefa

**ID da Tarefa:** `[GERADO_AUTOMATICAMENTE]`  
**Data de Criação:** `[DATA]`  
**Prioridade:** □ Alta □ Média □ Baixa  
**Estimativa:** `[ ]` horas

---

## 🎯 **INFORMAÇÕES BÁSICAS**

### **Metadados (OBRIGATÓRIOS - para extração automática):**
```
Tipo: [Feature | Bug | Melhoria | Investigação | Refatoração]
Módulo: [ui_web | core_lume | supply | pdv_ui | cash_flow | etc]
Objetivo: [Descrição concisa do que deve ser feito]
Decisões: [Padrões a seguir, decisões técnicas, restrições]
```

### **Exemplo formatado para extração:**
```
Tipo: Feature
Módulo: ui_web
Objetivo: Implementar UI para gestão de fornecedores
Decisões: Seguir padrão MemberHandler, cards com CNPJ opcional, validação anti-float
```

---

## 📝 **DESCRIÇÃO DETALHADA**

### **Contexto:**
`[Por que esta tarefa é necessária? Qual problema resolve? Qual o contexto de negócio?]`

### **Escopo:**
- `[Funcionalidade 1 a ser implementada]`
- `[Funcionalidade 2 a ser implementada]`
- `[O que NÃO está incluído no escopo]`

### **Requisitos Funcionais:**
1. `[RF1 - O que o sistema deve fazer]`
2. `[RF2 - Comportamento esperado]`
3. `[RF3 - Regras de negócio]`

### **Requisitos Não-Funcionais:**
- **Performance:** `[Tempo de resposta esperado, carga, etc.]`
- **Usabilidade:** `[Experiência do usuário esperada]`
- **Segurança:** `[Restrições de acesso, validações]`
- **Compatibilidade:** `[Navegadores, dispositivos]`

---

## 🏗️ **DESIGN TÉCNICO**

### **Arquitetura:**
- **Backend:** `[Serviços, APIs, banco de dados envolvidos]`
- **Frontend:** `[Handlers, templates, componentes]`
- **Integrações:** `[Módulos que precisam ser integrados]`

### **Componentes a Criar/Modificar:**
```
/modules/ui_web/internal/handler/[nome]_handler.go
/modules/ui_web/templates/[nome]_simple.html
/modules/ui_web/internal/handler/[nome]_handler_test.go
```

### **Dependências:**
- **Pré-requisitos:** `[O que precisa existir antes de implementar]`
- **Dependências externas:** `[Bibliotecas, serviços]`
- **Impacto em outros módulos:** `[O que pode quebrar]`

---

## 🧪 **CRITÉRIOS DE ACEITE**

### **Testes Obrigatórios:**
- [ ] **Testes unitários:** Cobertura >90% para handler
- [ ] **Testes de integração:** Com banco SQLite real
- [ ] **Smoke test:** `./scripts/dev/smoke_test_new_feature.sh`
- [ ] **Validação E2E:** `./scripts/dev/validate_e2e.sh --basic --headless`

### **Validação Funcional:**
- [ ] Login e navegação funcionam
- [ ] CRUD completo via HTMX
- [ ] Validações de dados funcionam
- [ ] Mensagens de erro amigáveis
- [ ] Design segue "Soberania e Suor"

### **Validação Técnica:**
- [ ] Anti-float compliance (zero `float` para valores)
- [ ] Soberania mantida (dados só no `.sqlite` da entidade)
- [ ] Cache-proof (usa `ExecuteTemplate` do `BaseHandler`)
- [ ] Navegação unificada em templates principais

---

## 🔗 **REFERÊNCIAS**

### **Código de Referência:**
- **Handler similar:** `[ex: MemberHandler, CashHandler]`
- **Template base:** `[ex: dashboard_simple.html]`
- **Testes de referência:** `[arquivos de teste similares]`

### **Documentação:**
- **Checklist pré-implementação:** `docs/templates/pre_implementation_checklist.md`
- **Plano de implementação:** `docs/templates/implementation_plan.md`
- **Antipadrões:** `docs/antipatterns/common_antipatterns_solutions.md`

### **Links:**
- [ ] Issue/GitHub: `[URL]`
- [ ] Design/Figura: `[URL]`
- [ ] Documentação técnica: `[URL]`

---

## 📅 **PLANO DE EXECUÇÃO**

### **Fase 1: Análise (Dia 1)**
- [ ] Preencher checklist pré-implementação
- [ ] Analisar backend correspondente
- [ ] Identificar padrões a seguir

### **Fase 2: Implementação (Dia 2-3)**
- [ ] Criar handler seguindo padrões
- [ ] Desenvolver template HTMX
- [ ] Implementar testes TDD

### **Fase 3: Integração (Dia 4)**
- [ ] Atualizar navegação
- [ ] Testes de integração
- [ ] Validação E2E

### **Fase 4: Validação (Dia 5)**
- [ ] Smoke test
- [ ] Validação E2E completa
- [ ] Documentação de aprendizados

---

## 🚨 **RISCOS IDENTIFICADOS**

| Risco | Probabilidade | Impacto | Mitigação |
|-------|--------------|---------|-----------|
| `[Risco técnico]` | □ Alta □ Média □ Baixa | □ Alto □ Médio □ Baixo | `[Ação de mitigação]` |
| `[Risco de escopo]` | □ Alta □ Média □ Baixa | □ Alto □ Médio □ Baixo | `[Ação de mitigação]` |
| `[Risco de dependência]` | □ Alta □ Média □ Baixa | □ Alto □ Médio □ Baixo | `[Ação de mitigação]` |

---

## 📊 **MÉTRICAS ESPERADAS**

### **Qualidade de Código:**
- **Cobertura de testes:** >90%
- **Complexidade ciclomática:** < 10
- **Linhas de código:** `[estimativa]`

### **Performance:**
- **Tempo de resposta:** < 2s
- **Uso de memória:** Estável
- **Tamanho do bundle:** `[estimativa]`

### **UX:**
- **Acessibilidade:** WCAG 2.1 AA
- **Responsividade:** Mobile-first
- **Feedback visual:** Imediato (< 500ms)

---

## 📝 **NOTAS ADICIONAIS**

`[Qualquer informação adicional, contexto histórico, decisões anteriores, etc.]`

---

## ✅ **CHECKLIST FINAL**

### **Antes de começar:**
- [ ] Checklist pré-implementação preenchido
- [ ] Backend analisado e compreendido
- [ ] Padrões de referência identificados
- [ ] Riscos mapeados e mitigados

### **Após implementação:**
- [ ] Testes unitários passam (>90% cobertura)
- [ ] Smoke test passa
- [ ] **Validação E2E passa** (`validate_e2e.sh --basic --headless`)
- [ ] Documentação atualizada
- [ ] Aprendizados registrados no `conclude_task.sh`

---

**📌 Como usar este template:**
1. Copie este template para `docs/tasks/[nome_da_tarefa].md`
2. Preencha todas as seções (especialmente os metadados no topo)
3. Use com: `./process_task.sh --file docs/tasks/[nome_da_tarefa].md --execute`
4. O script extrairá automaticamente Tipo, Módulo, Objetivo e Decisões

**💡 Dica:** Metadados formatados (`Tipo: X | Módulo: Y | Objetivo: Z | Decisões: W`) permitem extração automática pelo `process_task.sh`.