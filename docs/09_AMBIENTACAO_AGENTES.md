## 📄 `09_AMBIENTACAO_AGENTES.md` (O Orientador)

```markdown
# 🚀 Contexto do Projeto: Digna (Providentia Foundation)

## 📌 Visão Geral
**Projeto:** Sistema de Gestão Contábil Público e Open Source para Economia Solidária.
**Propósito:** Transformar a contabilidade de um peso burocrático em uma ferramenta de **Dignidade Financeira** para Empreendimentos de Economia Solidária (EES).
**Diferencial:** Arquitetura *Local-First Server-Side* com isolamento físico de dados.

---

## ✅ Status Atual das Sprints

| Sprint | Status | Descrição |
| :--- | :--- | :--- |
| Sprint 00 | ✅ COMPLETE | Blueprint, Governança e Visão de Produto. |
| **Sprint 01** | ⏭️ **EM ANDAMENTO** | **Lifecycle Manager & Core Ledger (Motor Lume)** |

---

## 🎯 Sprint 01 (Tarefa Atual)
**Objetivo:** Implementar o motor em Go que gerencia o ciclo de vida dos bancos SQLite individuais e o motor de lançamentos contábeis.

**Componentes Chave:**
- **Lume (Motor):** O backend em Go.
- **Digna (App):** A interface (futura).
- **Providentia (Fundação):** A entidade de governança.

---

## 🛠️ Stack Tecnológica
- **Linguagem:** Go 1.22+ (Clean Architecture).
- **Banco de Dados:** libSQL / SQLite (Individual por Tenant).
- **Escala:** Alvo de 1 milhão de EES (Indução Serpro).
- **Regra Financeira:** Tudo em `int64` (centavos). Proibido `float`.

---

## ⚠️ Regras para o Agente
1. **Soberania de Dados:** O dado pertence à entidade. O acesso deve ser isolado por arquivo.
2. **Padrão de Lançamento:** Seguir rigorosamente o princípio de partidas dobradas (soma zero).
3. **Mocks de Formalização:** O sistema deve prever a transição `DREAM` -> `FORMALIZED` via interfaces (Facade).

## 🛠️ Regras de Execução
1. **Contexto:** Leia sempre `docs/01_architecture.md` antes de sugerir mudanças.
2. **Espaço:** Nunca use espaços em caminhos de arquivos.
3. **Finanças:** Se vir um `float` no código contábil, pare e corrija para `int64`.
4. **Isolamento:** Código do módulo `lume` não deve criar arquivos; peça ao `lifecycle`.


```


