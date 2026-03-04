## 📄 `03_CODING_STANDARDS.md`

```markdown
# Coding Standards - Digna Core

**Language:** Go (Golang)  
**Standard:** Clean Architecture / Hexagonal

---

## ⚠️ REGRAS DE OURO (NÃO NEGOCIÁVEIS)

### Regra 1: Integridade Financeira
- **PROIBIDO** o uso de `float32` ou `float64` para valores monetários.
- **OBRIGATÓRIO** o uso de `int64` (representando centavos).
- Toda transação deve passar pelo motor de validação de **Partidas Dobradas**.

### Regra 2: Isolamento de Dados
- Nenhuma query SQL pode ser executada sem um `Context` que contenha o `EntityID`.
- O acesso ao arquivo `.sqlite` deve ser gerenciado pelo `LifecycleManager`.

### Regra 3: Erros e Logs
- Erros financeiros devem ser logados com nível `CRITICAL` e conter o `trace_id`.
- Logs não devem conter dados sensíveis (CPFs) em texto claro.

---

## 1. Style Guide
- Formatação via `gofmt`.
- Nomes de pacotes curtos e em minúsculas.
- Interfaces definidas no pacote que as consome.

```


