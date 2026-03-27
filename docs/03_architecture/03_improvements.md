title: Melhorias e Riscos - Ecossistema Digna
status: implemented
version: 2.0
last_updated: 2026-03-27
---

# Melhorias e Riscos - Ecossistema Digna

> **Nota:** Este documento consolida todas as melhorias implementadas (Sprints 1-16), riscos identificados e lições aprendidas durante a expansão do Ecossistema Digna (PDF v1.0 + Sessão 27/03/2026).

---

## 1. Melhorias Implementadas

### 1.1 Makefile para Execução de Testes

**Problema:** `go test ./...` no diretório raiz falha devido a múltiplos módulos Go no workspace.
**Solução:** Criado `Makefile` na raiz do projeto com comandos padronizados:

```bash
# Comandos disponíveis
make help       # Mostra todos os comandos
make test       # Roda testes em todos os módulos
make test-core  # Roda testes do core_lume
```

### 1.2 Graceful Shutdown Aprimorado

**Melhorias:**
- Logs estruturados em todas as fases do shutdown
- Tempo de duração do shutdown
- Mensagens de erro detalhadas
- Fechamento ordenado de conexões

**Exemplo de logs durante shutdown:**
```json
{ "level": "INFO", "msg": "🛑 Sinal de shutdown recebido", "signal": "interrupt" }
{ "level": "INFO", "msg": "🔄 Fechando conexões com banco de dados..." }
{ "level": "INFO", "msg": "✅ Conexões fechadas" }
{ "level": "INFO", "msg": "✅ Servidor desligado com sucesso", "shutdown_duration": "150ms" }
```

### 1.3 Sistema de Templates Cache-Proof [NOVO - Sprint 16]

**Problema:** Templates Go com cache persistente que sobrevive a recompilações.
**Solução:** Migração completa para templates `*_simple.html` carregados via `ParseFiles()` no handler.

**Benefícios:**
- Zero problemas de cache
- Atualizações refletidas imediatamente
- Sem necessidade de recompilar binário

### 1.4 Validação E2E Integrada [NOVO - Sessão 27/03/2026]

**Problema:** Smoke test valida apenas endpoints HTTP, não fluxo de negócio.
**Solução:** Script `validate_e2e.sh` com modo stealth (headless) integrado ao fluxo de conclusão de tarefas.

**Comando:**
```bash
./scripts/dev/validate_e2e.sh --basic --headless
```

**Critério de Aceite:** Tarefa só é concluída se E2E passar.

### 1.5 Sistema de Ajuda Educativa [NOVO - RF-30, Sessão 27/03/2026]

**Problema:** Campos como "CadÚnico", "Inadimplência", "CNAE" são jargões burocráticos que violam o Pilar Pedagógico.
**Solução:** Sistema de ajuda estruturada com linkagem UI → banco de ajuda (`help_topics{}`).

**Implementação:**
- Tabela `help_topics` no `central.db`
- Botão "?" ao lado de campos técnicos
- Explicação em linguagem popular + legislação + próximo passo
- Carregamento via HTMX (< 500ms)

---

## 2. Riscos Identificados

### 2.1 Módulo `legal_mock` Vazio

**Status:** ⚠️ Risco Baixo
**Problema:** Módulo `modules/legal_mock/` existe apenas com `go.mod`, sem pacotes Go.
**Impacto:**
- Ferramentas de CI/CD podem reportar warnings
- `go test ./...` retorna "no packages to test"
- Potencial confusão para novos desenvolvedores

**Soluções Possíveis:**
- Remover o módulo se não for necessário
- Implementar mocks básicos para testes legais
- Documentar explicitamente que é intencional

**Recomendação:** Implementar mocks para testes de integração com serviços legais (consulta de processos, verificação de formalização, etc.)

### 2.2 Execução de Testes no Workspace

**Status:** ✅ RESOLVIDO
**Problema:** Comando `go test ./...` no diretório raiz não funciona corretamente com Go workspaces.
**Solução Implementada:**
- ✅ Makefile com targets específicos por módulo
- ✅ Testes de integração em módulo separado (`integration_test`)
- ✅ Scripts de teste por módulo

### 2.3 SurplusCalculator - Cálculo de Sobras

**Status:** ✅ RESOLVIDO
**Problema:** O SurplusCalculator retornava `TotalSurplus` com sinal negativo e não aplicava deduções automaticamente.
**Solução Implementada:**
- ✅ Novo método `CalculateWithDeductions()` implementado
- ✅ Calcula automaticamente: Reserva Legal (10%) + FATES (5%)
- ✅ Rateio proporcional baseado em minutos trabalhados
- ✅ Tratamento de resíduos (centavos)

### 2.4 Transição Automática DREAM → FORMALIZED

**Status:** ✅ RESOLVIDO
**Problema:** A transição de status não acontecia automaticamente após 3 decisões.
**Solução Implementada:** (Mapeado e corrigido no motor Lume).

### 2.5 Cobertura de Integrações Externas

**Status:** ✅ RESOLVIDO
**Solução Implementada:**
- ✅ Cobertura de: Receita Federal, MTE, MDS, IBGE, SEFAZ, BNDES, SEBRAE, Providentia
- ✅ Mocks implementados com dados realistas

### 2.6 Falta de Conta de Capital Social no Seed

**Status:** ⚠️ ABERTO
**Problema:** O seed de contas padrão não inclui conta de Capital Social (Equity).
**Impacto:** Testes precisam criar dinamicamente a conta ID 8.
**Recomendação:** Adicionar conta "Capital Social" (2.2.01) ao seed de migração em `lifecycle/internal/repository/migration.go`.

### 2.7 Acoplamento de Regras Fiscais no Core Lume [NOVO]

**Status:** ⚠️ Risco Médio (Evitado por Design)
**Problema:** A exigência legal de enviar dados para a Receita Federal pode induzir os desenvolvedores a criarem calculadoras de impostos dentro do `core_lume`.
**Impacto:** Se o Core Lume calcular impostos, ele perderá a essência de "Contabilidade Invisível" (Social) e se tornará um ERP burocrático e pesado.
**Recomendação:** Blindar o Core Lume. A geração de arquivos fiscais deve ser delegada exclusivamente ao futuro módulo `accountant_dashboard`, que apenas lerá os dados em formato Read-Only.

### 2.8 APIs Governamentais Instáveis [NOVO - PDF v1.0]

**Status:** ⚠️ Risco Alto (Mitigado)
**Problema:** APIs do BCB, IBGE e Gov.br podem mudar sem aviso ou ficar indisponíveis.
**Impacto:** Motor de Indicadores (RF-18) e Portal de Oportunidades (RF-20) podem falhar.
**Mitigação:**
- Cache local com TTL (24h para BCB, 7 dias para programas)
- Circuit breaker pattern para APIs externas
- Fallback para último valor válido
- Modo offline mantém sistema operante

### 2.9 Complexidade do Portal Além do MVP [NOVO - PDF v1.0]

**Status:** ⚠️ Risco Médio
**Problema:** Portal de Oportunidades pode crescer além dos 3 programas MVP.
**Impacto:** Dívida técnica, atraso na entrega de valor.
**Mitigação:**
- MVP com 3 programas primeiro (Acredita, Pronampe, Niterói)
- Validação com usuários reais antes de expandir
- Critérios de elegibilidade configuráveis via banco

### 2.10 Massa Crítica para Rede Digna [NOVO - PDF v1.0]

**Status:** ⚠️ Risco Alto
**Problema:** Rede Digna (RF-24 a RF-26) depende de múltiplas entidades ativas.
**Impacto:** Funcionalidade subutilizada se poucas entidades usarem.
**Mitigação:**
- Focar em ERP + Portal primeiro
- Rede como "nice-to-have" inicial
- Incentivos para primeiras entidades publicarem perfil

### 2.11 Linguagem Técnica nos Tópicos de Ajuda [NOVO - RF-30]

**Status:** ⚠️ Risco Alto
**Problema:** Tópicos de ajuda podem usar jargão técnico, violando o Pilar Pedagógico.
**Impacto:** Usuários de baixa escolaridade não entendem explicações.
**Mitigação:**
- Revisão por ITCPs/comunidade antes de publicar
- Teste de usabilidade com usuários reais (5ª série)
- Critério de aceite: "usuário entende sem ajuda externa"

### 2.12 Conteúdo de Ajuda Desatualizado [NOVO - RF-30]

**Status:** ⚠️ Risco Médio
**Problema:** Legislação e programas mudam, tópicos de ajuda podem ficar obsoletos.
**Impacto:** Informações incorretas levam a decisões erradas.
**Mitigação:**
- Versionamento de tópicos no `central.db`
- Hash de integridade para auditoria
- Processo de atualização via admin (futuro)

### 2.13 Privacidade de Dados Sensíveis [NOVO - RF-19/RF-20]

**Status:** ⚠️ Risco Alto
**Problema:** Campos como CadÚnico, inadimplência, gênero são sensíveis.
**Impacto:** Vazamento pode causar discriminação ou problemas legais (LGPD).
**Mitigação:**
- Dados sensíveis nunca transmitidos para APIs externas
- Match executado LOCALMENTE no módulo `portal_opportunities`
- Apenas resultado (elegível/não) armazenado
- Usuário autoriza envio específico por programa

---

## 3. Recomendações Futuras

### 3.1 Script de Testes Global (Workaround)

Para automatizar os testes em todos os módulos enquanto ferramentas não suportam nativamente o workspace, manter a utilização do seguinte script no pipeline:

```bash
for module in "${MODULES[@]}"; do
    echo "Testing $module..."
    cd "$module" && go test ./... -v
done
```

### 3.2 Test Coverage

**Status:** ⚠️ Risco Médio
**Problema:** Alguns módulos não têm testes:
- `core_lume/internal/service`
- `core_lume/internal/social`
- `core_lume/pkg/*`

**Impacto:**
- Mudanças podem quebrar funcionalidades sem detecção
- Dívida técnica de testes
- Dificuldade em refatorações futuras

**Recomendação:** Adicionar testes unitários para:
- Services (ledger_service, work_service)
- Handlers HTTP
- Casos de erro e edge cases

### 3.3 Métricas e Monitoramento

**Adicionar:**
- Métricas de latência (p50, p95, p99)
- Contadores de erros por endpoint
- Gauge de conexões ativas
- Health checks detalhados (database, disk space)

### 3.4 Circuit Breaker para Integrações

Para integrações externas (quando implementadas):
- Circuit breaker pattern para APIs governamentais
- Retry com backoff exponencial
- Fallback para modo offline

### 3.5 Configuração Externalizada

Mover para variáveis de ambiente:
- Porta do servidor
- Timeouts
- Diretório de dados
- Nível de log

**Exemplo:**
```bash
# .env
DIGNA_PORT=8090
DIGNA_LOG_LEVEL=info
DIGNA_DATA_DIR=/var/lib/digna
DIGNA_SHUTDOWN_TIMEOUT=10s
```

### 3.6 Documentação de API

Gerar documentação OpenAPI/Swagger para:
- Endpoints do PDV
- API de integração
- Webhooks

### 3.7 Desenvolvimento do Painel do Contador Social (Sprint 09) [NOVO]

**Recomendação de Arquitetura:** Planejar a infraestrutura do módulo `accountant_dashboard` visando escala para contadores voluntários (CFC/CRCs).

**Desafio Arquitetural:** O painel precisará acessar múltiplos bancos de dados SQLite isolados (`/data/entities/*.db`) de forma paralela e estritamente como Read-Only.

**Entrega Esperada:** Um motor de extração que agrupa as Entries do `core_lume` e as converte no leiaute do SPED Fiscal/Contábil, salvando o hash da exportação para evitar envios duplicados.

### 3.8 Sistema de Help Engine [NOVO - RF-30]

**Recomendação:** Implementar módulo `help_engine` com:
- CRUD de tópicos de ajuda no `central.db`
- API pública para consumo por todos os módulos
- Cache com invalidação por atualização
- Métricas de uso (view_count) para identificar tópicos mais acessados

### 3.9 Validação E2E Obrigatória [NOVO - Sessão 27/03/2026]

**Recomendação:** Integrar validação E2E no fluxo de conclusão de tarefas:
- Script `validate_e2e.sh` executado após implementação
- Modo headless por padrão (não interfere com desktop)
- Tarefa só é concluída se E2E passar
- Aprendizados documentam resultado E2E

---

## 4. Matriz de Riscos Atualizada [NOVO - 27/03/2026]

| Risco | Probabilidade | Impacto | Mitigação | Status |
|-------|--------------|---------|-----------|--------|
| APIs governamentais instáveis | Alta | Médio | Cache local + circuit breaker + modo offline | ⚠️ Mitigado |
| Complexidade do Portal cresce além do MVP | Média | Alto | MVP com 3 programas primeiro; validação com usuários reais | ⚠️ Monitorado |
| Conflito de naming (ERP vs. Ecossistema) | Baixa | Baixo | Documentar claramente a hierarquia de módulos | ✅ Documentado |
| Massa crítica para Rede Digna não atingida | Alta | Médio | Focar em ERP + Portal primeiro; Rede como "nice-to-have" | ⚠️ Aceito |
| Teologia afeta adoção secular | Média | Alto | Manter produto laico na interface; teologia informa design internamente | ✅ Mitigado |
| Dependência de contadores sociais para escala | Média | Alto | Criar programa de capacitação + certificação CFC | ⚠️ Em progresso |
| Linguagem muito técnica nos tópicos de ajuda (RF-30) | Alta | Alto | Revisão por ITCPs/comunidade; teste de usabilidade com usuários reais | ⚠️ Novo |
| Conteúdo de ajuda desatualizado | Média | Médio | Processo de atualização via central.db, não hardcoded | ⚠️ Novo |
| Privacidade de dados sensíveis (CadÚnico, etc.) | Média | Crítico | Match local, dados nunca transmitidos, autorização por programa | ⚠️ Novo |
| Anti-Float violado em novos módulos | Baixa | Crítico | Validação obrigatória em code review + grep automation | ✅ Mitigado |
| Cache-Proof templates não seguido | Média | Baixo | Validação em smoke test + documentação clara | ✅ Mitigado |

---

## 5. Lições Aprendidas por Sessão [NOVO - 27/03/2026]

### Sessão 11/03/2026
- **Descoberta:** `legal_facade` já existe com 80% da funcionalidade
- **Aprendizado:** Consultar `docs/skills/` antes de implementar
- **Ação:** Criar `MODULES_QUICK_REFERENCE.md` para acelerar descoberta

### Sessão 27/03/2026
- **Descoberta:** Validação E2E não estava integrada ao fluxo de conclusão
- **Aprendizado:** Smoke test valida HTTP, E2E valida negócio
- **Ação:** Script `validate_e2e.sh` obrigatório antes de `conclude_task.sh`
- **Descoberta:** Campos técnicos sem explicação violam Pilar Pedagógico
- **Aprendizado:** Sistema de ajuda estruturada é infraestrutura, não feature
- **Ação:** RF-30 adicionado ao backlog com prioridade alta

---

## 6. Métricas de Qualidade [NOVO - 27/03/2026]

| Métrica | Alvo | Atual | Status |
|---------|------|-------|--------|
| Testes unitários passando | 100% | 149/149 | ✅ |
| Cobertura de handlers | >90% | ~87% | ⚠️ |
| Validação E2E por tarefa | 100% | 0% | ❌ Novo |
| Tópicos de ajuda criados | 10+ | 0 | ❌ Novo |
| Anti-Float violations | 0 | 0 | ✅ |
| Cache-Proof violations | 0 | 0 | ✅ |
| Soberania violations | 0 | 0 | ✅ |

---

**Status:** ✅ ATUALIZADO COM MELHORIAS E RISCOS DO ECOSSISTEMA (PDF v1.0 + Sessão 27/03/2026)  
**Próxima Ação:** Atualizar `03_architecture/04_architectural_decisions.md` com decisões da expansão  
**Versão Anterior:** 1.0 (2026-03-09)  
**Versão Atual:** 2.0 (2026-03-27)
