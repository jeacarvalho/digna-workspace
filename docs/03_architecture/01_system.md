
#### title: Arquitetura do Sistema Digna
status: implemented
version: 1.5
last_updated: 2026-03-08

### Arquitetura do Sistema - Digna
**Projeto:** Sistema de Gestão Contábil e Pedagógica para Economia Solidária
**Arquitetura:** Local-First Server-Side com Micro-databases Isolados + Domain-Driven Design

--------------------------------------------------------------------------------

#### 1. Visão Geral da Arquitetura
O Digna utiliza uma arquitetura de **Micro-databases isolados** combinada com **Domain-Driven Design (DDD)**. Cada entidade possui sua própria instância física local (Soberania de Dados). O código segue rigorosamente os princípios de Clean Architecture, garantindo que o domínio central seja protegido, ao mesmo tempo em que atua como uma **ponte tecnológica** para a classe contábil.

##### 1.1 High-Level Design (Visão Sociotécnica e Contábil)
O diagrama abaixo ilustra como a arquitetura atende a dois mundos distintos: a linguagem coloquial do trabalhador (operação) e o rigor técnico do contador parceiro (auditoria), ambos convergindo no mesmo motor de domínio.

    [Trabalhador/Cooperado]                         [Contador Social (Parceiro)]
              |                                                 |
    (Linguagem Coloquial/Ação)                      (Auditoria Fiscal e ITG 2002)
              v                                                 v
    +---------------------------------+       +-----------------------------------+
    | INTERFACE LAYER (PDV UI)        |       | ACCOUNTANT DASHBOARD (Novo)       |
    | Tradução Cultural, Zero jargões |       | Visão Multi-tenant e Export SPED  |
    +---------------------------------+       +-----------------------------------+
              |                                                 |
              | (Use Cases / DTOs)                              | (Read-Only API / DTOs)
              v                                                 v
    +-----------------------------------------------------------------------------+
    |                         APPLICATION LAYER (Services)                        |
    |                         Coordenação de casos de uso                         |
    +-----------------------------------------------------------------------------+
                                      |
                                      | (Repository Interfaces)
                                      v
    +-----------------------------------------------------------------------------+
    |                         DOMAIN LAYER (Core Lume)                            |
    | - Entities (Entry, Posting, WorkLog, FiscalBatch)                           |
    | - Domain Services (LedgerService, TaxTranslatorService)                     |
    | - Repository Interfaces (Contracts)                                         |
    +-----------------------------------------------------------------------------+
                                      |
                                      | (Repository Implementations)
                                      v
    +-----------------------------------------------------------------------------+
    |                         INFRASTRUCTURE LAYER                                |
    | - SQLite Repositories (Persistência)                                        |
    | - Mock Integration Repositories / SPED File Generators                      |
    +-----------------------------------------------------------------------------+
                                      |
                                      | (Dados Estruturados / int64)
                                      v
    +-----------------------------------------------------------------------------+
    |                  Persistência e Soberania de Dados                          |
    |         (Banco SQLite Isolado Exclusivo por Tenant / Entidade)              |
    +-----------------------------------------------------------------------------+

--------------------------------------------------------------------------------

#### 2. Tecnologias Core

| Camada | Tecnologia | Justificativa |
| ------ | ------ | ------ |
| **Backend** | Go (1.22+) | Performance, concorrência, binário estático |
| **Database** | SQLite3 | Isolamento total por arquivo (Soberania) |
| **Arquitetura** | Clean Arch + DDD | Domínio independente de frameworks |
| **Numerics** | `int64` (exclusivo) | Valores financeiros e horas (ITG 2002) exatos |

--------------------------------------------------------------------------------

#### 3. Arquitetura DDD

##### 3.1 Repository Pattern
Exemplo de como o domínio contábil é blindado da infraestrutura física:

```go
// Domain Layer - Interface pura (Não sabe o que é banco de dados)
type LedgerRepository interface {
    SaveEntry(entry *Entry) (int64, error)
    SavePosting(posting *Posting) error
    GetBalance(accountID int64) (int64, error)
}

// Infrastructure Layer - Implementação SQLite
type SQLiteLedgerRepository struct {
    lifecycleManager lifecycle.LifecycleManager
}
// ... implementação dos métodos SQL injetando a entidade conectada
```

--------------------------------------------------------------------------------

#### 4. Módulos Implementados e Sprints

| Sprint | Módulo | Status | Testes |
| ------ | ------ | ------ | ------ |
| 01 | lifecycle | ✅ | 6/6 |
| 02 | core_lume + pdv_ui | ✅ | 8/8 |
| 03 | reporting + legal_facade | ✅ | 8/8 |
| 04 | sync_engine | ✅ | 9/9 |
| 05 | ui_web | ✅ | 9/9 |
| 06 | cash_flow | ✅ | 3/3 |
| 07 | DDD Refactoring | ✅ | 43/43 |
| 08 | integrations | ✅ | 5/5 |
| 09 | accountant_dashboard | ✅ | 8/8 |
| 10 | member_management | ✅ | 19/19 |
| 11 | formalization_e2e | ✅ | 5/5 |
| 12 | accountant_dashboard_complete | ✅ | 8/8 |

##### 4.8 integrations (Sprint 08) ✅
**Responsabilidade:** Integrações externas governamentais.
- **8 Interfaces:** Receita Federal, MTE, MDS, IBGE, SEFAZ, BNDES, SEBRAE, Providentia.
- **Mock Implementations:** Todas funcionando com dados realistas.
- **DDD Pattern:** Domínio independente, fácil substituir mocks por HTTP real.

##### 4.9 accountant_dashboard (Sprint 09-12) ✅ [COMPLETE]
**Responsabilidade:** Interface Multi-tenant para a classe contábil e exportações fiscais.
- **Visão Agregada:** Permite ao Contador Social visualizar o status de fechamento de múltiplos arquivos SQLite (Tenants) simultaneamente.
- **Isolamento Read-Only:** O módulo acessa os micro-databases apenas para leitura (`?mode=ro`) e auditoria da ITG 2002 e dos bloqueios de FATES/Reserva Legal.
- **Tradução Fiscal:** Mapeia as *Entries* geradas pelo Core Lume para os leiautes padrão da Receita Federal (SPED/CSV), acabando com o trabalho braçal de digitação do contador.
- **Anti-Float:** Todos os valores monetários usam `int64`, sem `float`.
- **Test Coverage:** Domain: 100%, Handler: 97.1%, Repository: 87.2%, Service: 91.3%, Public API: 26.7%

--------------------------------------------------------------------------------

#### 5. Princípios SOLID Aplicados

- **SRP (Single Responsibility):** Cada módulo tem uma única responsabilidade (ex: `pdv_ui` traduz cultura, `accountant_dashboard` traduz obrigações fiscais).
- **OCP (Open/Closed):** Sistema aberto para extensão (novas integrações governamentais sem mudar o Core Lume).
- **LSP (Liskov Substitution):** Implementações de Repository são intercambiáveis (SQLite ↔ Mock ↔ HTTP).
- **ISP (Interface Segregation):** Interfaces pequenas e específicas na camada de domínio.
- **DIP (Dependency Inversion):** Services (ex: LedgerService) dependem de abstrações (interfaces), não de implementações concretas do SQLite.

--------------------------------------------------------------------------------

#### 6. Stack Tecnológico Final

| Camada | Tecnologia | Uso |
| ------ | ------ | ------ |
| Backend | Go 1.22+ | API REST, binário leve, concorrência |
| Storage | SQLite3 | Isolamento por tenant |
| Front/Web| HTMX + Tailwind | PWA mobile-first para o trabalhador |
| Front/Dashboard| Vue/React ou HTMX | Visão Multi-tenant para o Contador |
| Hash | SHA256 | Auditoria CADSOL e imutabilidade |
| Numerics | `int64` | Centavos monetários e minutos trabalhados |
| Architecture | Clean Arch + DDD | Domínio protegido |
| Fiscal | SPED / CSV Export | Ponte com sistemas contábeis comerciais |
```
