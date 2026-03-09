***

```markdown
---
title: Melhorias e Riscos
status: draft
version: 1.1
last_updated: 2026-03-08
---

# Melhorias e Riscos - Digna

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
{"level":"INFO","msg":"🛑 Sinal de shutdown recebido","signal":"interrupt"}
{"level":"INFO","msg":"🔄 Fechando conexões com banco de dados..."}
{"level":"INFO","msg":"✅ Conexões fechadas"}
{"level":"INFO","msg":"✅ Servidor desligado com sucesso","shutdown_duration":"150ms"}
```

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
1. **Remover** o módulo se não for necessário
2. **Implementar** mocks básicos para testes legais
3. **Documentar** explicitamente que é intencional
**Recomendação:** Implementar mocks para testes de integração com serviços legais (consulta de processos, verificação de formalização, etc.)

### 2.2 Execução de Testes no Workspace
**Status:** ✅ RESOLVIDO
**Problema:** Comando `go test ./...` no diretório raiz não funciona corretamente com Go workspaces.
**Solução Implementada:**
1. ✅ Makefile com targets específicos por módulo
2. ✅ Testes de integração em módulo separado (`integration_test`)
3. ✅ Scripts de teste por módulo

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

### 2.6 Falta de Conta de Capital Social no Seed
**Status:** ⚠️ ABERTO
**Problema:** O seed de contas padrão não inclui conta de Capital Social (Equity).
**Impacto:** Testes precisam criar dinamicamente a conta ID 8.
**Recomendação:** Adicionar conta "Capital Social" (2.2.01) ao seed de migração em `lifecycle/internal/repository/migration.go`.

### 2.7 Acoplamento de Regras Fiscais no Core Lume [NOVO]
**Status:** ⚠️ Risco Médio (Evitado por Design)
**Problema:** A exigência legal de enviar dados para a Receita Federal pode induzir os desenvolvedores a criarem calculadoras de impostos dentro do `core_lume`.
**Impacto:** Se o Core Lume calcular impostos, ele perderá a essência de "Contabilidade Invisível" (Social) e se tornará um ERP burocrático e pesado.
**Recomendação:** Blindar o Core Lume. A geração de arquivos fiscais deve ser delegada exclusivamente ao futuro módulo `accountant_dashboard`, que apenas lerá os dados em formato *Read-Only*.

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
1. Services (ledger_service, work_service)
2. Handlers HTTP
3. Casos de erro e edge cases

### 3.3 Métricas e Monitoramento
Adicionar:
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
DIGNA_PORT=8080
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
- **Desafio Arquitetural:** O painel precisará acessar múltiplos bancos de dados SQLite isolados (`/data/entities/*.db`) de forma paralela e estritamente como *Read-Only*.
- **Entrega Esperada:** Um motor de extração que agrupa as *Entries* do `core_lume` e as converte no leiaute do **SPED Fiscal/Contábil**, salvando o *hash* da exportação para evitar envios duplicados.
```

***
