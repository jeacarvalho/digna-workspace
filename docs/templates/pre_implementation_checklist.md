# 📋 Template de Checklist Pré-Implementação

**Baseado em:** Aprendizados da sessão 11/03/2026
**Última atualização:** 11/03/2026

---

## 🔄 Melhorias Identificadas

*Baseado em feedback de implementações anteriores:*

### Sessão 11/03/2026
- **ANTES de implementar:** Verificar se funcionalidade já existe (ex: `legal_facade` tem 80%)
- **ANTES de implementar:** Consultar `docs/skills/` para padrões específicos
- **ANTES de implementar:** Analisar código similar existente (ex: `member_handler.go`)
- **ANTES de implementar:** Verificar padrões já implementados (SHA256, file download)
- **PROBLEMA:** `process_task.sh --file="nome.md"` tem issues com parsing

---

## ✅ CHECKLIST COMPLETO PRÉ-IMPLEMENTAÇÃO

### 1. 🏗️ Análise Arquitetural (OBRIGATÓRIO)

#### 1.1 Verificar se já existe
- [ ] **Módulo correspondente existe?** `find modules -name "*[feature]*" -type d`
- [ ] **Funcionalidade similar já implementada?** `grep -r "[funcionalidade]" modules/`
- [ ] **Padrões similares identificados?** Analisar handlers/templates similares

#### 1.2 Skills Relevantes
- [ ] **Consultou `docs/skills/`?** `ls docs/skills/`
- [ ] **Skill de backend aplicável?** (Anti-float, DDD, TDD)
- [ ] **Skill de frontend aplicável?** (HTMX, cache-proof templates)
- [ ] **Skill de dados aplicável?** (Soberania, LifecycleManager)

### 2. 🎨 Padrões de Implementação

#### 2.1 Para CRUD/HTMX features:
- [ ] **Handler estende `BaseHandler`?** (padrão `member_handler.go`)
- [ ] **Template `*_simple.html`?** (cache-proof)
- [ ] **Rotas HTMX padrão?** `GET /{feature}`, `POST /{feature}/create`
- [ ] **Design system aplicado?** (#2A5CAA, #4A7F3E, #F57F17)

#### 2.2 Para File Download/Export:
- [ ] **Analisou `accountant_handler.go`?** (padrão file download)
- [ ] **Headers corretos?** `Content-Type`, `Content-Disposition`
- [ ] **Hash para integridade?** `X-Export-Hash` (opcional)

#### 2.3 Para Document Generation:
- [ ] **Verificou `legal_facade/generator.go`?** (padrão existente)
- [ ] **SHA256 já implementado?** `sha256.Sum256([]byte(data))`
- [ ] **FormalizationSimulator?** `MinDecisionsForFormalization = 3`

### 3. ⚙️ Testabilidade

#### 3.1 Estrutura de Testes:
- [ ] **Testes similares existem?** `find modules/ui_web -name "*test*.go"`
- [ ] **Mock de `LifecycleManager`?** (padrão em testes de handlers)
- [ ] **Cobertura >90%?** (requisito para handlers)

#### 3.2 Setup de Testes:
- [ ] **`httptest.NewRecorder()`?** (padrão para testar handlers)
- [ ] **Entity ID mockado?** `r.URL.Query().Get("entity_id")`
- [ ] **Templates carregados?** `ParseFiles()` no teste

### 4. 🚨 Validações Específicas

#### 4.1 Anti-Float (REGRA SAGRADA):
- [ ] **Nenhum `float` para valores financeiros/tempo?**
- [ ] **`int64` para centavos/minutos?** (R$ 1,00 = 100)
- [ ] **Validação de soma zero?** (se aplicável)

#### 4.2 Soberania de Dados:
- [ ] **Acesso apenas ao banco da entidade atual?**
- [ ] **Sem JOINs entre bancos diferentes?**
- [ ] **`entity_id` extraído do contexto?** `r.Context().Value("entity_id")`

#### 4.3 Cache-Proof Templates:
- [ ] **Template `*_simple.html`?** (documento HTML completo)
- [ ] **`ParseFiles()` NO HANDLER?** (não variável global)
- [ ] **Sem `ParseGlob()` ou templates globais?**

### 5. 📝 Decisões Documentadas

#### 5.1 Decisões Técnicas:
- [ ] **Handler structure:** □ Estende BaseHandler □ Independente
- [ ] **Template base:** □ dashboard_simple.html □ Outro: ________
- [ ] **Interatividade:** □ HTMX □ API+JavaScript □ Misto

#### 5.2 Dependências:
- [ ] **Backend service existe?** □ Sim □ Não □ Precisa criar
- [ ] **Acessível do UI?** □ Público □ Internal □ Mock necessário
- [ ] **Repository interface existe?** □ Sim □ Não □ Precisa criar

### 6. 🔍 Comandos de Validação Rápida

#### Execute ANTES de começar:
```bash
# 1. Validação rápida dos módulos
./scripts/tools/quick_agent_check.sh all

# 2. Verificar se funcionalidade já existe
find modules -name "*[feature]*" -type f

# 3. Analisar handler similar
./scripts/tools/analyze_patterns.sh [handler_similar] --all

# 4. Verificar padrões específicos
grep -r "sha256.Sum256\|Content-Disposition" modules/
```

### 7. 📊 Riscos Identificados

| Risco | Probabilidade | Impacto | Mitigação |
|-------|--------------|---------|-----------|
| Funcionalidade já existe | Alta | Alto | Verificar com `quick_agent_check.sh` |
| Performance issues | Média | Médio | Paginação, lazy loading |
| Template cache issues | Alta | Baixo | Usar `ParseFiles()` no handler |
| Entity ID não encontrado | Média | Médio | Validação no middleware |

---

## 🎯 CRITÉRIOS DE PRONTO PARA IMPLEMENTAR

### ✅ TODAS estas condições devem ser atendidas:
- [ ] Backend analisado e compreendido
- [ ] Padrões de frontend identificados  
- [ ] Riscos mapeados e mitigados
- [ ] Decisões documentadas
- [ ] Checklist completo preenchido
- [ ] `quick_agent_check.sh` executado e analisado

### ⏱️ Tempo Estimado para Esta Análise:
- **Sem documentação:** 40-60 minutos
- **Com este checklist:** 10-15 minutos
- **Economia:** 70-75%

---

## 📌 NOTAS IMPORTANTES

### Arquivos Temporários vs. Permanentes:
- **PERMANENTE:** `docs/QUICK_REFERENCE.md`, `docs/ANTIPATTERNS.md`, `docs/learnings/`
- **TEMPORÁRIO:** `.agent_context.md` (excluído no `end_session.sh`)
- **SEMPRE:** Documentar aprendizados em arquivos permanentes

### Fluxo Correto:
1. `./start_session.sh` (cria `.agent_context.md` novo)
2. Executar `quick_agent_check.sh`
3. Preencher ESTE checklist
4. Implementar
5. `./conclude_task.sh` (aprendizados vão para arquivos permanentes)
6. `./end_session.sh` (limpa arquivos temporários, atualiza documentação)

### Problemas Conhecidos:
- **`process_task.sh --file="nome.md"`**: Pode falhar com caracteres especiais
- **Solução:** Ler arquivo manualmente antes: `cat "nome.md"`

---

**📚 REFERÊNCIAS:**
- `docs/MODULES_QUICK_REFERENCE.md` - Mapa completo de módulos
- `docs/learnings/SESSION_INSIGHTS_20260311.md` - Aprendizados detalhados
- `docs/skills/` - Skills específicas do projeto
- `scripts/tools/quick_agent_check.sh` - Validação rápida
<!-- 20260311_112301 - legal_dossier - 11/03/2026 -->
<!-- Aprendizado: Implementação da função GenerateDossier() no módulo legal_facade: geração de dossiê jurídico com has... -->
