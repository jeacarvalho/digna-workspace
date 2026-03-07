# 05 MODULARIZATION - Digna Core

**Status:** Proposed (v0.1)  
**Last Updated:** 2026-03-04  
**Strategy:** Domain-Driven Atomicity

---

## 1. Visão Geral
Para mitigar a entropia de contexto em agentes de IA, o sistema Digna é dividido em quatro módulos independentes. Cada módulo possui responsabilidades únicas e comunica-se exclusivamente através de interfaces (contratos), permitindo o desenvolvimento e teste isolado de cada componente.



---

## 2. Divisão de Módulos

### Módulo 0: `pdv_ui` (Interface de Operação)
- **Papel:** Interface de venda e compra para o empreendedor.
- **Output:** Comandos comerciais para o Motor Lume.

### Módulo 1: Lifecycle Manager (Orquestrador de Tenants)
- **Papel:** Gerenciamento físico dos bancos SQLite.
- **Componentes:** `TenantRegistry`, `MigrationEngine`.
O módulo "físico" do sistema, responsável pela infraestrutura de arquivos.
* **Responsabilidade:** Gerenciar a criação, abertura, cache de conexões e migrações dos arquivos `.sqlite` individuais.
* **Componentes Chave:** - `TenantRegistry`: Mapeamento de Entity_ID para o path físico no servidor.
    - `MigrationEngine`: Garantia de paridade de schema em todos os micro-databases.
* **Status na v0:** Crítico (Base para todos os outros módulos).

### Módulo 2: Motor Lume (Core Ledger)
- **Papel:** Motor de partidas dobradas e validação de integridade.
- **Componentes:** `JournalService`, `TransactionValidator`.
O domínio contábil puro, onde reside a inteligência financeira.
* **Responsabilidade:** Implementar o sistema de partidas dobradas e garantir a integridade das transações.
* **Componentes Chave:**
    - `TransactionValidator`: Garante a regra de soma zero (Débitos + Créditos = 0).
    - `JournalService`: Executa a escrita atômica nos registros do Tenant.
* **Regra Estrita:** Processamento exclusivo em `int64` (centavos).

### Módulo 3: Legal Facade (Simulador de Formalização)
- **Papel:** Simulação de transição de estado e geração de documentos (CNPJ).
O módulo de transição de estado e geração de documentos jurídicos.
* **Responsabilidade:** Gerir a metamorfose da entidade do status `DREAM` para `FORMALIZED`.
* **Componentes Chave:**
    - `LegalMock`: Provedor de dados simulados (CNPJ, Inscrição Estadual).
    - `DocumentGenerator`: Motor de templates para exportação de Atas e Estatutos em PDF.

### Módulo 4: Reporting Engine (Painel de Dignidade)
- **Papel:** Cálculo de rateio social e dossiê de crédito.
A camada de inteligência e saída para o usuário final.
* **Responsabilidade:** Traduzir os dados brutos do Ledger em indicadores sociais e financeiros.
* **Componentes Chave:**
    - `SurplusCalculator`: Algoritmo de rateio de sobras baseado em critérios sociais.
    - `CreditDossier`: Consolidador de histórico financeiro para análise de crédito.

---

## 3. Matriz de Dependências e Contratos

| Módulo | Consome de | Entrega para |
| :--- | :--- | :--- |
| **Lifecycle** | Sistema de Arquivos / OS | Ledger / Legal Facade |
| **Ledger** | Lifecycle (DB Connection) | Reporting Engine |
| **Legal Facade** | Lifecycle (Metadata) | Usuário Final (PDF) |
| **Reporting** | Ledger (Transaction History) | Usuário Final (UI/PDF) |

---

## 4. Estratégia de Implementação para Agentes
1. **Isolamento de Sprint:** Cada Módulo deve ter seu próprio `SESSION_LOG` e ser finalizado com testes unitários antes do início do próximo.
2. **Interface First:** Antes de codificar a lógica interna, o agente deve definir as interfaces no pacote `internal/domain`.
3. **Validation:** O agente deve validar o Módulo 1 com a criação física de arquivos antes de tentar realizar qualquer lançamento contábil (Módulo 2).


