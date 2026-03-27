# Estratégia de Testes para Digna

## Visão Geral

Este documento descreve a estratégia de testes automatizados para o projeto Digna, permitindo rodar smoke tests e testes E2E com dados consistentes e isolados.

## Ambiente de Teste

### Estrutura de Diretórios

```
data/
├── test/                    # Dados de teste (isolados)
│   ├── central/
│   │   └── central.db      # Banco central de teste
│   └── entities/
│       └── test-entity-001.db  # Banco da entidade de teste
└── entities/               # Dados de produção (não tocar)
```

### Dados de Teste Padrão

**Usuário de Teste:**
- ID: `test-user-001`
- Email: `test@digna.local`
- Role: `COORDINATOR`
- Status: `ACTIVE`

**Entidades de Teste:**
1. **test-entity-001** (Cooperativa Teste)
   - CNPJ: 00000000000191
   - Tipo: Cooperativa
   
2. **test-entity-002** (MEI Teste)
   - CNPJ: 00000000000272
   - Tipo: MEI

**Membros:**
- Coordenador Teste (test-member-001)
- MEI Teste (test-member-002)

## Scripts Disponíveis

### 1. Setup do Ambiente de Teste

```bash
./scripts/dev/setup_test_env.sh
```

Este script:
- Cria a estrutura de diretórios em `data/test/`
- Inicializa o banco central com dados de teste
- Cria bancos de entidades de teste
- Gera arquivo `.env.test` com configurações

### 2. Smoke Test com Autenticação

```bash
# Uso básico
./scripts/dev/smoke_test_with_auth.sh "Nome da Feature" "/rota" test

# Exemplo
./scripts/dev/smoke_test_with_auth.sh "Central de Ajuda" "/help" test test-entity-001
```

Parâmetros:
1. Nome da feature (ex: "Central de Ajuda")
2. Rota a testar (ex: "/help")
3. Modo: `test` ou `prod`
4. Entity ID (opcional, padrão: test-entity-001)

### 3. Smoke Test Original

```bash
./scripts/dev/smoke_test_new_feature.sh "Nome" "/rota"
```

## Fluxo de Testes

### Para Novas Features

1. **Desenvolver com TDD:**
   ```bash
   # Criar testes unitários primeiro
   go test -v ./internal/domain/...
   go test -v ./internal/service/...
   go test -v ./internal/handler/...
   ```

2. **Configurar ambiente de teste:**
   ```bash
   ./scripts/dev/setup_test_env.sh
   ```

3. **Iniciar servidor em modo teste:**
   ```bash
   cd modules/ui_web
   DIGNA_ENV=test go run .
   ```

4. **Rodar smoke test:**
   ```bash
   ./scripts/dev/smoke_test_with_auth.sh "Minha Feature" "/minha-rota" test
   ```

5. **Verificar templates:**
   - Template deve existir: `modules/ui_web/templates/{feature}_simple.html`
   - Deve conter HTML válido

## Dados de Teste SQL

Para adicionar mais dados de teste, edite:
- `scripts/dev/test_data.sql` - Dados SQL puros
- `scripts/dev/setup_test_env.sh` - Setup programático

## Boas Práticas

### 1. Isolamento
- Nunca use dados de produção em testes
- Sempre use `data/test/` para testes
- Resetar banco de teste quando necessário

### 2. Autenticação em Testes
- Use `test-user-001` como usuário padrão
- Passe `entity_id` nas URLs para testes
- Adicione `test_mode=true` quando suportado

### 3. Templates
- Sempre crie `{feature}_simple.html`
- Use inline templates como fallback no handler
- Siga o padrão "Soberania e Suor"

### 4. Testes Unitários
- **Obrigatório:** Testes para Domain, Repository, Service e Handler
- Cobertura mínima: 80%
- Use mocks para isolamento

## Resolução de Problemas

### Servidor não encontra dados de teste

```bash
# Verifique se o setup foi executado
ls -la data/test/

# Se não existir, rode setup
./scripts/dev/setup_test_env.sh
```

### Smoke test falha com 404

```bash
# Verifique se o handler está registrado em main.go
grep -n "New.*Handler" modules/ui_web/main.go

# Verifique se a rota existe
grep -n "HandleFunc.*rota" modules/ui_web/internal/handler/*.go
```

### Template não encontrado

```bash
# Liste templates existentes
ls modules/ui_web/templates/*_simple.html

# Verifique se o nome segue o padrão: {feature}_simple.html
```

## CI/CD

Para integração contínua:

```yaml
# Exemplo de pipeline
steps:
  - name: Setup Test Environment
    run: ./scripts/dev/setup_test_env.sh
  
  - name: Start Server (Test Mode)
    run: cd modules/ui_web && DIGNA_ENV=test go run . &
  
  - name: Run Unit Tests
    run: go test -v ./...
  
  - name: Run Smoke Tests
    run: |
      ./scripts/dev/smoke_test_with_auth.sh "Feature 1" "/rota1" test
      ./scripts/dev/smoke_test_with_auth.sh "Feature 2" "/rota2" test
```

## Checklist de Testes

Antes de concluir uma tarefa:

- [ ] Testes unitários criados (Domain, Service, Handler)
- [ ] Todos os testes passando (`go test ./...`)
- [ ] Ambiente de teste configurado
- [ ] Smoke test executado com sucesso
- [ ] Template `{feature}_simple.html` criado
- [ ] Handler registrado em `main.go`
- [ ] Dados de teste documentados

## Contato

Dúvidas sobre testes? Consulte:
- `docs/QUICK_REFERENCE.md` - Arquitetura
- `docs/ANTIPATTERNS.md` - O que evitar
- Este arquivo para estratégia de testes
