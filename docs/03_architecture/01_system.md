#### title: Arquitetura do Sistema Digna
status: implemented
version: 1.4
last_updated: 2026-03-07

### Arquitetura do Sistema - Digna
**Projeto:** Sistema de Gestão Contábil e Pedagógica para Economia Solidária
**Arquitetura:** Local-First Server-Side com Micro-databases Isolados + Domain-Driven Design

--------------------------------------------------------------------------------

#### 1. Visão Geral da Arquitetura
O Digna utiliza uma arquitetura de **Micro-databases isolados** combinada com **Domain-Driven Design (DDD)**. Cada entidade possui sua própria instância física local (Soberania de Dados), enquanto o código segue rigorosamente os princípios de Clean Architecture e DDD.

##### 1.1 High-Level Design (Visão Sociotécnica)

    [Trabalhador/Cooperado]
              |
              | (Linguagem Coloquial)
              v
    +---------------------------------------------------+
    | INTERFACE LAYER (Web / PDV UI)                    |
    | Tradução Cultural, Zero jargões contábeis          |
    +---------------------------------------------------+
              |
              | (Use Cases / DTOs)
              v
    +---------------------------------------------------+
    | APPLICATION LAYER (Services)                      |
    | Coordenação de casos de uso                       |
    +---------------------------------------------------+
              |
              | (Repository Interfaces)
              v
    +---------------------------------------------------+
    | DOMAIN LAYER (Core Lume)                          |
    | - Entities (Entry, Posting, WorkLog)              |
    | - Domain Services (LedgerService)                 |
    | - Repository Interfaces (Contracts)             |
    +---------------------------------------------------+
              |
              | (Repository Implementations)
              v
    +---------------------------------------------------+
    | INFRASTRUCTURE LAYER                              |
    | - SQLite Repositories                             |
    | - Mock Integration Repositories                   |
    +---------------------------------------------------+
              |
              | (Dados Estruturados / int64)
              v
    +---------------------------------------------------+
    | Persistência e Soberania de Dados                |
    | (SQLite Isolado por Tenant)                       |
    +---------------------------------------------------+

--------------------------------------------------------------------------------

#### 2. Tecnologias Core
| Camada | Tecnologia | Justificativa |
| ------ | ------ | ------ |
| **Backend** | Go (1.22+) | Performance, concorrência, binário estático |
| **Database** | SQLite3 | Isolamento total por arquivo |
| **Arquitetura** | Clean Arch + DDD | Domínio independente de frameworks |
| **Numerics** | int64 (exclusivo) | Valores financeiros exatos |

--------------------------------------------------------------------------------

#### 3. Arquitetura DDD

##### 3.1 Repository Pattern

```go
// Domain Layer - Interface pura
type LedgerRepository interface {
    SaveEntry(entry *Entry) (int64, error)
    SavePosting(posting *Posting) error
    GetBalance(accountID int64) (int64, error)
}

// Infrastructure Layer - Implementação SQLite
type SQLiteLedgerRepository struct {
    lifecycleManager lifecycle.LifecycleManager
}
```

##### 3.2 Camadas de Responsabilidade

1. **Domain Layer:** Entidades, Value Objects, Domain Services, Repository Interfaces
2. **Application Layer:** Use Cases, Coordenação, DTOs
3. **Infrastructure Layer:** Repository Implementations (SQLite), External APIs
4. **Interface Layer:** HTTP Handlers, Templates, API REST

--------------------------------------------------------------------------------

#### 4. Módulos Implementados

| Sprint | Módulo | Status | Testes |
|--------|--------|--------|--------|
| 01 | lifecycle | ✅ | 6/6 |
| 02 | core_lume + pdv_ui | ✅ | 8/8 |
| 03 | reporting + legal_facade | ✅ | 8/8 |
| 04 | sync_engine | ✅ | 9/9 |
| 05 | ui_web | ✅ | 9/9 |
| 06 | cash_flow | ✅ | 3/3 |
| 07 | DDD Refactoring | ✅ | 43/43 |
| 08 | integrations | ✅ | 5/5 |

##### 4.9 integrations (Sprint 08) ✅
**Responsabilidade:** Integrações externas governamentais.
- **8 Interfaces:** Receita Federal, MTE, MDS, IBGE, SEFAZ, BNDES, SEBRAE, Providentia
- **Mock Implementations:** Todas funcionando com dados realistas
- **DDD Pattern:** Domínio independente, fácil substituir mocks por HTTP real

--------------------------------------------------------------------------------

#### 5. Princípios SOLID Aplicados

- **SRP:** Cada módulo tem uma única responsabilidade
- **OCP:** Sistema aberto para extensão (novas integrações sem mudar código)
- **LSP:** Implementações de Repository são intercambiáveis (SQLite ↔ Mock ↔ HTTP)
- **ISP:** Interfaces pequenas e específicas
- **DIP:** Services dependem de abstrações (interfaces), não de implementações

--------------------------------------------------------------------------------

#### 6. Stack Tecnológico Final
| Camada | Tecnologia | Uso |
| ------ | ------ | ------ |
| Backend | Go 1.22+ | API REST, binário leve |
| Storage | SQLite3 | Isolamento por tenant |
| Front/Web| HTMX + Tailwind | PWA mobile-first |
| Hash | SHA256 | Auditoria CADSOL |
| Numerics | int64 | Centavos e minutos |
| Architecture | Clean Arch + DDD | Domínio protegido |

***

### Resumo Versão 1.4

**Adicionado:**
- ✅ Arquitetura DDD completa
- ✅ Repository Pattern em todos os módulos
- ✅ Módulo integrations/ com 8 interfaces governamentais
- ✅ Princípios SOLID documentados
- ✅ 91 testes automatizados (100% PASS)

**Total de Testes:** 91/91 ✅ (6+8+8+9+9+3+43+5)
