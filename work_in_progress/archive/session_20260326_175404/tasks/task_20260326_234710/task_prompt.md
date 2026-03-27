# 📋 Prompt para RF-30 (Fase 2): Integração de Ajuda em Todas as Funcionalidades

**Tipo:** Feature + Refatoração de UX  
**Módulo:** `modules/ui_web` (todos os templates), `modules/core_lume` (help_engine)  
**Prioridade:** ALTA (Pilar Pedagógico - acessibilidade para baixa escolaridade)  
**Estimativa:** 16-24 horas (sistema inteiro)  
**Dependência:** RF-30 Fase 1 (Sistema de Ajuda Estruturada) ✅ CONCLUÍDO  

---

## 🎯 CONTEXTO DA TAREFA

O **RF-30 (Sistema de Ajuda Educativa)** já foi implementado na Fase 1 com:
- ✅ Tabela `help_topics` no `central.db`
- ✅ Handler `/help` com índice e busca
- ✅ Seed de 6 tópicos iniciais (CadÚnico, Inadimplência, CNAE, DAS MEI, Reserva Legal, FATES)
- ✅ Template `help_topic_simple.html`

**AGORA (Fase 2):** Precisamos **integrar** este sistema em TODAS as telas do sistema, adicionando botões "?" ao lado de termos técnicos que possam confundir usuários de baixa escolaridade.

**Princípio Central:** *"Nenhum usuário deve se sentir humilhado por não entender um termo. Todo conceito técnico deve ter explicação acessível em 1 clique."*

---

## 📝 ESCOPO DA INTEGRAÇÃO

### Módulos a Serem Revisados

| Módulo | Templates | Prioridade | Termos Críticos |
|--------|-----------|------------|-----------------|
| **PDV** | `pdv_simple.html` | ALTA | CNAE, DAS MEI, ICMS, ISS |
| **Caixa** | `cash_simple.html` | ALTA | Fluxo de caixa, Saldo, Entradas/Saídas |
| **Compras** | `supply_dashboard_simple.html` | ALTA | Fornecedor, Estoque, Insumo |
| **Estoque** | `supply_stock_simple.html` | ALTA | Insumo, Produto, Mercadoria |
| **Orçamento** | `budget_simple.html` | MÉDIA | Planejamento, Realizado, SAFE/WARNING |
| **Membros** | `member_simple.html` | MÉDIA | Coordenador, Status, Capital Social |
| **Contador** | `accountant_dashboard_simple.html` | MÉDIA | SPED, EFD-Reinf, ECF, ITG 2002 |
| **Formalização** | `legal_simple.html` | ALTA | CADSOL, Assembleia, Hash SHA256 |
| **Vínculo Contábil** | `accountant_link_simple.html` | ALTA | Exit Power, Read-Only, Vigência |
| **Perfil Elegibilidade** | `eligibility_simple.html` | CRÍTICA | CadÚnico, Inadimplência, Tipo Entidade |

---

## 🏗️ PADRÃO DE IMPLEMENTAÇÃO

### Componente Reutilizável: Help Tooltip

**Arquivo:** `modules/ui_web/templates/components/help_tooltip.html`

```html
<!-- Componente reutilizável para links de ajuda -->
{{define "help_tooltip"}}
<a href="/help/topic/{{.Key}}" 
   target="_blank" 
   class="help-tooltip-trigger text-digna-primary hover:text-digna-accent transition-colors"
   aria-label="Saiba mais sobre {{.Title}}"
   title="Clique para explicação">
  <svg class="w-4 h-4 inline-block ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
          d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.125 2.5-2.5 2.5-.69 0-1.25.28-1.5.75-.25.47-.25 1.03 0 1.5.25.47.81.75 1.5.75 1.375 0 2.5-1.1 2.5-2.5 0-1.657-1.79-3-4-3-1.742 0-3.223.835-3.772 2z"/>
  </svg>
</a>
{{end}}
```

### Uso nos Templates

**ANTES:**
```html
<label for="inscrito_cad_unico">Inscrito no CadÚnico?</label>
```

**DEPOIS:**
```html
<label for="inscrito_cad_unico" class="flex items-center gap-2">
  Inscrito no CadÚnico?
  {{template "help_tooltip" dict "Key" "cadunico" "Title" "CadÚnico"}}
</label>
```

---

## 🛠️ TAREFAS DE IMPLEMENTAÇÃO

### Fase 1: Auditoria de Termos (4 horas)

**Objetivo:** Identificar TODOS os termos técnicos em todos os templates.

**Checklist por Template:**
```markdown
- [ ] Listar todos os labels, placeholders, textos de ajuda
- [ ] Identificar termos técnicos (jargão contábil, fiscal, jurídico)
- [ ] Mapear para tópicos de ajuda existentes ou criar novos
- [ ] Documentar em `docs/implementation_plans/help_integration_audit.md`
```

**Termos Críticos a Mapear:**
| Termo | Tópico Help | Template(s) |
|-------|-------------|-------------|
| CadÚnico | `cadunico` | eligibility, portal |
| Inadimplência | `inadimplencia` | eligibility, portal |
| CNAE | `cnae` | cadastro, pdv, portal |
| DAS MEI | `das_mei` | pdv, fiscal |
| ICMS/ISS | `icms_iss` | pdv, fiscal |
| Reserva Legal | `reserva_legal` | dashboard, distribution |
| FATES | `fates` | dashboard, distribution |
| SPED | `sped` | accountant_dashboard |
| EFD-Reinf | `efd_reinf` | accountant_dashboard |
| ECF | `ecf` | accountant_dashboard |
| ITG 2002 | `itg_2002` | accountant_dashboard, work_log |
| CADSOL | `cadsol` | legal_facade |
| Exit Power | `exit_power` | accountant_link |
| Read-Only | `read_only` | accountant_link |

### Fase 2: Criar Tópicos de Ajuda Faltantes (4 horas)

**Objetivo:** Garantir que TODO termo mapeado tenha tópico no `help_topics`.

**Template de Novo Tópico:**
```go
{
    Key:         "icms_iss",
    Title:       "O que é ICMS e ISS?",
    Summary:     "São impostos sobre vendas e serviços.",
    Explanation: "ICMS é imposto sobre produtos (comércio). ISS é imposto sobre serviços. O Digna calcula automaticamente quando você registra uma venda.",
    WhyAsked:    "Precisamos saber isso para calcular corretamente os impostos do seu negócio.",
    Legislation: "Lei Complementar nº 123/2006 (Simples Nacional)",
    NextSteps:   "O Digna calcula automaticamente. Você só precisa registrar vendas corretamente no PDV.",
    OfficialLink: "https://www.gov.br/empresas-e-negocios",
    Category:    "TRIBUTARIO",
    Tags:        "imposto,venda,serviço",
}
```

**Tópicos Mínimos Obrigatórios (20+):**
1. `cadunico` ✅ (existe)
2. `inadimplencia` ✅ (existe)
3. `cnae` ✅ (existe)
4. `das_mei` ✅ (existe)
5. `reserva_legal` ✅ (existe)
6. `fates` ✅ (existe)
7. `icms_iss` ❌ (criar)
8. `sped` ❌ (criar)
9. `efd_reinf` ❌ (criar)
10. `ecf` ❌ (criar)
11. `itg_2002` ❌ (criar)
12. `cadsol` ❌ (criar)
13. `exit_power` ❌ (criar)
14. `read_only` ❌ (criar)
15. `fluxo_caixa` ❌ (criar)
16. `estoque_insumo` ❌ (criar)
17. `estoque_produto` ❌ (criar)
18. `estoque_mercadoria` ❌ (criar)
19. `capital_social` ❌ (criar)
20. `rateio_sobras` ❌ (criar)

### Fase 3: Integrar Links em Todos os Templates (8 horas)

**Objetivo:** Adicionar botão "?" ao lado de CADA termo técnico em CADA template.

**Checklist por Template:**
```markdown
- [ ] Importar componente `help_tooltip`
- [ ] Adicionar link ao lado de cada termo técnico
- [ ] Testar que todos os links funcionam
- [ ] Validar que tooltip não quebra layout mobile
- [ ] Verificar acessibilidade (aria-label, title)
```

**Templates Prioritários (ordem de implementação):**
1. `eligibility_simple.html` (CRÍTICA - muitos termos técnicos)
2. `pdv_simple.html` (ALTA - uso diário)
3. `cash_simple.html` (ALTA - uso diário)
4. `supply_stock_simple.html` (ALTA - categorização confusa)
5. `accountant_link_simple.html` (ALTA - termos jurídicos)
6. `accountant_dashboard_simple.html` (MÉDIA - uso do contador)
7. `budget_simple.html` (MÉDIA - conceitos financeiros)
8. `member_simple.html` (MÉDIA - governança)
9. `legal_simple.html` (ALTA - formalização)
10. `dashboard_simple.html` (BAIXA - apenas indicadores)

### Fase 4: Validação e Testes (4 horas)

**Critérios de Aceite:**
- [ ] **100% dos termos técnicos** mapeados têm link "?"
- [ ] **Todos os links** funcionam e levam ao tópico correto
- [ ] **Layout não quebra** em mobile (testar em 320px width)
- [ ] **Acessibilidade** validada (leitor de tela, keyboard navigation)
- [ ] **Performance** não degradada (tooltip carrega em < 500ms)
- [ ] **Smoke test** passa: `./scripts/dev/smoke_test_new_feature.sh "Ajuda Integrada" "/help"`
- [ ] **Validação E2E** inclui verificação de links de ajuda

**Teste de Validação Manual:**
```bash
# 1. Acessar cada tela principal
# 2. Contar botões "?" visíveis
# 3. Clicar em cada um e validar que abre tópico correto
# 4. Validar que tópico tem linguagem para 5ª série

# Checklist de validação:
- [ ] /eligibility - 7+ botões "?"
- [ ] /pdv - 3+ botões "?"
- [ ] /cash - 2+ botões "?"
- [ ] /supply/stock - 3+ botões "?"
- [ ] /accountant/links - 3+ botões "?"
- [ ] /help - índice com 20+ tópicos
```

---

## ✅ CRITÉRIOS DE ACEITE (Definition of Done)

### Arquitetura
- [ ] Componente `help_tooltip` reutilizável criado
- [ ] Todos os templates importam componente corretamente
- [ ] Nenhum hardcoding de ícones/links de ajuda

### Funcionalidade
- [ ] 20+ tópicos de ajuda criados no `central.db`
- [ ] 50+ links "?" distribuídos em todos os templates
- [ ] Todos os links validados (não há 404)
- [ ] Links abrem em nova aba (`target="_blank"`)

### Pedagogia (CRÍTICO - RF-30)
- [ ] **Linguagem para 5ª série** em TODOS os tópicos
- [ ] **Zero jargões** sem explicação
- [ ] **Próximo passo acionável** em cada tópico
- [ ] **Validação com usuário real** (teste de usabilidade)

### UI/UX
- [ ] Design consistente em todos os templates
- [ ] Mobile-responsive (não quebra layout em 320px)
- [ ] Acessibilidade (aria-label, keyboard navigation)
- [ ] Cores seguem "Soberania e Suor" (#2A5CAA para links)

### Testes
- [ ] Smoke test passando
- [ ] Validação E2E inclui verificação de ajuda
- [ ] Validação manual de TODOS os links
- [ ] Documentação atualizada em `docs/02_product/help_topics.md`

---

## 📊 MÉTRICAS DE SUCESSO

| Métrica | Alvo | Como Medir |
|---------|------|------------|
| Tópicos de ajuda criados | 20+ | `SELECT COUNT(*) FROM help_topics` |
| Links "?" em templates | 50+ | `grep -r "help_tooltip" modules/ui_web/templates/` |
| Telas com ajuda integrada | 100% (10/10) | Checklist de templates prioritários |
| Validação de links | 100% funcionais | Teste manual de cada link |
| Redução de abandono em formulários | 30% | Analytics de `eligibility_simple.html` |
| Visualizações de ajuda/mês | 1000+ | `SUM(view_count)` no primeiro mês |

---

## 🚨 RISCOS E MITIGAÇÕES

| Risco | Probabilidade | Impacto | Mitigação |
|-------|--------------|---------|-----------|
| Links quebrados (404) | Alta | Médio | Validação manual de TODOS os links antes de concluir |
| Layout quebrado em mobile | Média | Alto | Teste em 320px width antes de merge |
| Tópicos com linguagem técnica | Alta | Crítico | Revisão por ITCPs antes de publicar |
| Performance degradada | Baixa | Baixo | Cache de tópicos já implementado no RF-30 Fase 1 |
| Links não visíveis | Média | Médio | Contraste validado (WCAG 2.1 AA) |
| Acessibilidade não validada | Média | Alto | Teste com leitor de tela (NVDA/VoiceOver) |

---

## 📅 CRONOGRAMA ESTIMADO

| Fase | Duração | Entregas |
|------|---------|----------|
| **1. Auditoria de Termos** | 4 horas | Lista completa de termos + mapeamento |
| **2. Criar Tópicos Faltantes** | 4 horas | 14+ novos tópicos no `central.db` |
| **3. Integrar em Templates** | 8 horas | 50+ links em 10 templates |
| **4. Validação e Testes** | 4 horas | Smoke test, E2E, validação manual |
| **TOTAL** | **20 horas** | **Sistema 100% integrado** |

---

## 🎯 INSTRUÇÕES PARA O AGENTE

### ANTES DE COMEÇAR:
1. [ ] Executar `./scripts/tools/quick_agent_check.sh all`
2. [ ] Ler RF-30 Fase 1 (já implementado)
3. [ ] Consultar `docs/skills/applying-solidarity-logic/SKILL.md`
4. [ ] Preencher checklist em `docs/implementation_plans/help_integration_pre_check.md`

### DURANTE IMPLEMENTAÇÃO:
1. [ ] Começar pela auditoria (Fase 1) - NÃO pular para codificação
2. [ ] Criar tópicos de ajuda ANTES de adicionar links
3. [ ] Usar componente reutilizável, não hardcoding
4. [ ] Testar em mobile a cada template modificado
5. [ ] Validar linguagem com critério "5ª série"

### APÓS IMPLEMENTAÇÃO:
1. [ ] Smoke test: `./scripts/dev/smoke_test_new_feature.sh "Ajuda Integrada" "/help"`
2. [ ] Validação E2E: `./scripts/dev/validate_e2e.sh --basic --headless`
3. [ ] Validação manual de TODOS os links (checklist acima)
4. [ ] Documentar aprendizados: `./conclude_task.sh "Aprendizados: [resumo]" --success`

### NUNCA:
- ❌ Adicionar links sem tópico de ajuda correspondente
- ❌ Usar linguagem técnica nos tópicos
- ❌ Quebrar layout mobile
- ❌ Concluir sem validação manual de todos os links

---

## 📋 CHECKLIST DE VALIDAÇÃO FINAL

```bash
# 1. Contar tópicos de ajuda
sqlite3 data/entities/central.db "SELECT COUNT(*) FROM help_topics;"
# Deve retornar: 20+

# 2. Contar links em templates
grep -r "help_tooltip" modules/ui_web/templates/ | wc -l
# Deve retornar: 50+

# 3. Validar que todos os templates prioritários têm ajuda
for template in eligibility pdv cash supply_stock accountant_link; do
  echo "=== ${template}_simple.html ==="
  grep -c "help_tooltip" modules/ui_web/templates/${template}_simple.html
done

# 4. Smoke test
./scripts/dev/smoke_test_new_feature.sh "Ajuda Integrada" "/help"

# 5. Validação E2E
./scripts/dev/validate_e2e.sh --basic --headless

# 6. Validação manual (CRÍTICO)
# Acessar cada tela e clicar em CADA link "?"
# Validar que abre tópico correto com linguagem acessível
```

---

## 🔄 INTEGRAÇÃO COM RFs EXISTENTES

Esta task **NÃO cria novas funcionalidades**, apenas **integra ajuda** nas existentes:

| RF | Funcionalidade | Links de Ajuda Adicionados |
|----|----------------|---------------------------|
| RF-19 | Perfil de Elegibilidade | 7 links (CadÚnico, Inadimplência, Tipo Entidade, etc.) |
| RF-27 | DAS MEI | 3 links (DAS, ICMS, ISS) |
| RF-11 | Contador Social | 5 links (SPED, EFD-Reinf, ECF, ITG 2002, Read-Only) |
| RF-12 | Vínculo Contábil | 3 links (Exit Power, Vigência, Delegação) |
| RF-07/08 | Compras/Estoque | 4 links (Insumo, Produto, Mercadoria, Fornecedor) |
| RF-10 | Orçamento | 3 links (Planejamento, SAFE/WARNING, Realizado) |

---

**PRONTO PARA INICIAR?**

Confirme que compreendeu:
1. [ ] RF-30 Fase 1 já existe (sistema de ajuda)
2. [ ] Esta task é Fase 2 (integração em TODAS as telas)
3. [ ] Componente reutilizável `help_tooltip` deve ser criado
4. [ ] 20+ tópicos de ajuda devem existir
5. [ ] 50+ links "?" devem ser distribuídos
6. [ ] Validação manual de TODOS os links é obrigatória
7. [ ] Linguagem para 5ª série é critério crítico

Se todas as caixas estiverem marcadas, inicie pela **Fase 1: Auditoria de Termos**.

---

**ID da Tarefa:** `RF-30-FASE2-HELP-INTEGRATION-20260327`  
**Gerado em:** 27/03/2026  
**Próxima Revisão:** Após Fase 1 (Auditoria)  
**Dependência:** RF-30 Fase 1 ✅ CONCLUÍDO  
**Habilita:** Acessibilidade completa do sistema para baixa escolaridade