---
title: Melhorias e Riscos
status: draft
version: 1.0
last_updated: 2026-03-07
---

# Melhorias e Riscos - Digna

## Melhorias Implementadas

### 1. Makefile para Execução de Testes

**Problema:** `go test ./...` no diretório raiz falha devido a múltiplos módulos Go no workspace.

**Solução:** Criado `Makefile` na raiz do projeto com comandos padronizados:

```bash
# Comandos disponíveis
make help          # Mostra todos os comandos
make test          # Roda testes em todos os módulos
make test-core     # Roda testes do core_lume
make test-integration  # Roda testes de integração
make test-coverage   # Roda testes com cobertura
make build         # Builda o servidor
make clean         # Limpa arquivos de build
make lint          # Roda linter em todos os módulos
```

**Uso recomendado:**
```bash
# Testes rápidos durante desenvolvimento
make test-core

# Testes completos antes de commit
make test

# Build para produção
make build
```

### 2. Logging Estruturado (Observabilidade)

**Problema:** Logs em formato livre dificultam troubleshooting em produção.

**Solução:** Implementado `log/slog` do Go 1.21+ com:
- Formato JSON estruturado
- Níveis de log (INFO, ERROR)
- Campos contextualizados (method, path, status, duration)
- Middleware de logging HTTP

**Exemplo de log:**
```json
{
  "time": "2026-03-07T23:30:00Z",
  "level": "INFO",
  "msg": "HTTP request",
  "method": "POST",
  "path": "/api/sale",
  "status": 200,
  "duration": "45.2ms",
  "remote_addr": "127.0.0.1:54321"
}
```

**Localização:**
- Middleware: `modules/ui_web/internal/middleware/logger.go`
- Configuração: `modules/ui_web/main.go`

### 3. Graceful Shutdown Aprimorado

**Melhorias:**
- Logs estruturados em todas as fases do shutdown
- Tempo de duração do shutdown
- Mensagens de erro detalhadas
- Fechamento ordenado de conexões

**Exemplo de logs durante shutdown:**
```
{"level":"INFO","msg":"🛑 Sinal de shutdown recebido","signal":"interrupt"}
{"level":"INFO","msg":"🔄 Fechando conexões com banco de dados..."}
{"level":"INFO","msg":"✅ Conexões fechadas"}
{"level":"INFO","msg":"✅ Servidor desligado com sucesso","shutdown_duration":"150ms"}
```

## Riscos Identificados

### 1. Módulo `legal_mock` Vazio

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

### 2. Execução de Testes no Workspace

**Status:** ✅ RESOLVIDO

**Problema:** Comando `go test ./...` no diretório raiz não funciona corretamente com Go workspaces.

**Solução Implementada:**
1. ✅ Makefile com targets específicos por módulo
2. ✅ Testes de integração em módulo separado (`integration_test`)
3. ✅ Scripts de teste por módulo

---

### 3. SurplusCalculator - Cálculo de Sobras

**Status:** ✅ RESOLVIDO

**Problema:** O SurplusCalculator retornava `TotalSurplus` com sinal negativo e não aplicava deduções automaticamente.

**Solução Implementada:**
- ✅ Novo método `CalculateWithDeductions()` implementado
- ✅ Calcula automaticamente: Reserva Legal (10%) + FATES (5%)
- ✅ Rateio proporcional baseado em minutos trabalhados
- ✅ Tratamento de resíduos (centavos)

---

### 4. Transição Automática DREAM → FORMALIZED

**Status:** ✅ RESOLVIDO

**Problema:** A transição de status não acontecia automaticamente após 3 decisões.

**Solução Implementada:**
- ✅ Novo método `AutoTransitionIfReady()` implementado
- ✅ Transiciona automaticamente após 3 decisões registradas
- ✅ Integração com `CheckFormalizationCriteria`

---

### 5. Testes E2E de Integrações Governamentais

**Status:** ✅ RESOLVIDO

**Problema:** Não existiam testes E2E que validassem as integrações governamentais.

**Solução Implementada:**
- ✅ `journey_e2e_test.go` - Teste BDD da jornada anual
- ✅ `integrations_e2e_test.go` - Teste de todas as 8 integrações mock
- ✅ Cobertura de: Receita Federal, MTE, MDS, IBGE, SEFAZ, BNDES, SEBRAE, Providentia

---

### 6. Falta de Conta de Capital Social no Seed

**Status:** ⚠️ ABERTO

**Problema:** O seed de contas padrão não inclui conta de Capital Social (Equity).

**Impacto:** Testes precisam criar dinamicamente a conta ID 8.

**Recomendação:** Adicionar conta "Capital Social" (2.2.01) ao seed de migração em `lifecycle/internal/repository/migration.go`.

---

## Recomendações Futuras

for module in "${MODULES[@]}"; do
    echo "Testing $module..."
    cd "$module" && go test ./... -v
done
```

### 3. Test Coverage

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

## Recomendações Futuras

### 1. Métricas e Monitoramento

Adicionar:
- Métricas de latência (p50, p95, p99)
- Contadores de erros por endpoint
- Gauge de conexões ativas
- Health checks detalhados (database, disk space)

### 2. Circuit Breaker para Integrações

Para integrações externas (quando implementadas):
- Circuit breaker pattern para APIs governamentais
- Retry com backoff exponencial
- Fallback para modo offline

### 3. Configuração Externalizada

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

### 4. Documentação de API

Gerar documentação OpenAPI/Swagger para:
- Endpoints do PDV
- API de integração
- Webhooks

---

*Última atualização: 2026-03-07*
