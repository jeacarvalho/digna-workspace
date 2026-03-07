
### 📄 `docs/03_coding_standards.md`
**Objetivo:** Definir as restrições técnicas para os agentes.

```markdown
# Coding Standards - Digna Core

**Language:** Go (Golang)  
**Standard:** Clean Architecture / Multi-module Workspace

---

## ⚠️ REGRAS DE OURO (NÃO NEGOCIÁVEIS)

### Regra 1: Integridade Financeira
- **PROIBIDO** o uso de `float32/64`.
- **OBRIGATÓRIO** o uso de `int64` (centavos).
- Validação de **Soma Zero** em todos os lançamentos do Ledger.

### Regra 2: Nomenclatura de Arquivos
- **PROIBIDO** o uso de espaços em nomes de diretórios ou arquivos.
- **PADRÃO:** `kebab-case` para diretórios e `snake_case` para arquivos `.go`.

### Regra 3: Isolamento de Dados
- O acesso ao arquivo `.sqlite` é exclusivo via `LifecycleManager`.
- Cada Tenant deve ser isolado fisicamente em `data/entities/{entity_id}.db`.

---

## 1. Ferramental (VS Code)
- Extensão oficial `golang.go` obrigatória.
- Uso de `go.work` para orquestrar os módulos.
- Formatação via `gofmt`.

```

---
