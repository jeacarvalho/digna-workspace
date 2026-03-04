## 📄 `04_DATA_DICTIONARY.md`

```markdown
# Data Dictionary - Digna (Core Ledger)

**Versão:** 0.1  
**Storage:** SQLite (Individual por Tenant)

---

### Tabela: `accounts`
Armazena o plano de contas da entidade.
| Field | Type | Description |
| :--- | :--- | :--- |
| `id` | UUID | Identificador único da conta. |
| `name` | String | Nome da conta (Ex: Caixa, Estoque). |
| `type` | String | ASSET, LIABILITY, EQUITY, REVENUE, EXPENSE. |

### Tabela: `postings`
Os registros individuais de débito e crédito.
| Field | Type | Description |
| :--- | :--- | :--- |
| `amount` | Integer | Valor em centavos. Positivo para D, Negativo para C. |
| `account_id` | UUID | FK para a conta. |
| `entry_id` | UUID | FK para a transação agregadora. |

```

