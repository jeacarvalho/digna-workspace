#### title: Arquitetura do Sistema Digna
status: implemented
version: 1.3
last_updated: 2026-03-07

### Arquitetura do Sistema - Digna
**Projeto:** Sistema de Gestão Contábil e Pedagógica para Economia Solidária
**Arquitetura:** Local-First Server-Side com Micro-databases Isolados

--------------------------------------------------------------------------------

#### 1. Visão Geral da Arquitetura
O Digna utiliza uma arquitetura de **Micro-databases isolados**. Em vez de um banco único centralizado, cada entidade (cooperativa ou grupo) possui sua própria instância física local (Soberania de Dados). Isso garante o "Poder de Saída" (Exit Power) e a operação resiliente em áreas com baixa conectividade.

##### 1.1 High-Level Design (Visão Sociotécnica)
O diagrama abaixo ilustra como a interface atua como um escudo, traduzindo a ação real do trabalhador para a complexidade contábil do backend.

    [Trabalhador/Cooperado]
              |
              | (Linguagem Coloquial / Ação Real do Dia a Dia)
              v
    +---------------------------------------------------+
    | Camada de Tradução Cultural e Pedagógica          |
    | (PDV UI - Interface amigável, sem jargões)        |
    +---------------------------------------------------+
              |
              | (API: Intenção de Negócio)
              v
    +---------------------------------------------------+
    | Motor Lume (Regras Contábeis Exatas)              |
    | - Partidas Dobradas Invisíveis (Soma Zero)        |
    | - Valoração do Trabalho (ITG 2002)                |
    +---------------------------------------------------+
              |
              | (Dados Estruturados / int64)
              v
    +---------------------------------------------------+
    | Persistência e Soberania de Dados                 |
    | (Banco SQLite Isolado Exclusivo do Tenant)        |
    +---------------------------------------------------+
              |
              | (Dados Históricos e Agregados)
              v
    +---------------------------------------------------+
    | Transparência Algorítmica e Relatórios            |
    | (Dashboard Visual para Aprovação em Assembleia)   |
    +---------------------------------------------------+

--------------------------------------------------------------------------------

#### 2. Tecnologias Core
| Camada | Tecnologia | Justificativa Técnica e Social |
| ------ | ------ | ------ |
| **Backend** | Go (1.22+) | Performance, concorrência e binário estático para Nuvem Soberana (Serpro). |
| **Database** | SQLite3 + mattn/go-sqlite3 | Isolamento total por arquivo. Independência de nuvem central para operação offline. |
| **Sync** | Change Data Capture | Sincronização assíncrona tolerante a redes rurais/periféricas instáveis. |
| **Arquitetura** | Clean Architecture | Desacoplamento profundo entre a lógica contábil (exata) e a interface (humana). |
| **Numerics** | int64 (exclusivo) | Valores financeiros e tempo (ITG 2002) exatos sem erros IEEE 754. |

--------------------------------------------------------------------------------

#### 3. Módulos Implementados e Suas Funções
##### 3.1 lifecycle (Sprint 01 ✅)
**Responsabilidade:** Orquestração de arquivos SQLite por tenant.
*   Criação física de arquivos `.db` isolados (Garantia de Soberania).
*   Pool de conexões com PRAGMAs otimizados (WAL, foreign_keys).

##### 3.2 core_lume (Sprint 02 ✅)
**Responsabilidade:** Motor contábil, governança e valoração social (Invisível ao usuário).
*   **Integridade Contábil:** Soma(Débitos) = Soma(Créditos) = 0.
*   **ITG 2002:** Tempo de trabalho convertido rigorosamente em capital social.

##### 3.3 pdv_ui (Sprint 02 ✅)
**Responsabilidade:** Interface pedagógica e **Camada de Tradução Cultural**.
*   **RecordSale:** Traduz "vendi mel" em Débito Caixa + Crédito Vendas.
*   **RecordWork:** Traduz "trabalhei 4 horas" em Registro de Capital-Trabalho.

##### 3.4 reporting (Sprint 03 ✅)
**Responsabilidade:** Cálculos agregados e Transparência Algorítmica Visual.
*   **CalculateSocialSurplus:** Distribuição proporcional de sobras exibida em painéis gráficos simples para aprovação em Assembleia.

##### 3.5 legal_facade (Sprint 03 ✅)
**Responsabilidade:** Documentação institucional gerada em Markdown (Atas com hash SHA256) respeitando a transição gradual do grupo.

##### 3.6 sync_engine e 3.7 ui_web (Sprints 04 e 05 ✅)
*   **sync_engine:** Sincronização offline-first via Delta Tracking.
*   **ui_web:** Interface mobile-first (HTMX, Tailwind CSS, PWA).

--------------------------------------------------------------------------------

#### 4. Camadas de Segregação (Clean Architecture)

    +---------------------------------------------------------+
    | INTERFACE LAYER (Web / PDV UI)                          |
    | Tradução Cultural, Pedagógica e Zero jargões contábeis  |
    +---------------------------------------------------------+
                               | 
                               | (Consome DTOs / APIs)
                               v
    +---------------------------------------------------------+
    | DOMAIN LAYER (Core Lume)                                |
    | Regras de Negócio: Ledger (Soma Zero), Social Valuation |
    +---------------------------------------------------------+
                               | 
                               | (Persiste via SQL / int64)
                               v
    +---------------------------------------------------------+
    | INFRASTRUCTURE LAYER (Lifecycle)                        |
    | Persistência e I/O: SQLite Manager (1 DB por Entidade)  |
    +---------------------------------------------------------+

--------------------------------------------------------------------------------

#### 5. Fluxos de Dados

##### 5.1 Ciclo de Vida do Tenant (Transição Respeitosa)
O sistema respeita o tempo de maturação política do empreendimento antes de exigir conformidade burocrática.

    [Início] 
       |
       v
    (DREAM) ------> Grupo Informal operando com autonomia
       |
       v
    (INCUBATED) --> Apoio de ITCP/Incubadora (Módulos Pedagógicos)
       |
       v
    (FORMALIZED) -> Mínimo 3 decisões registradas + Estatuto Base
       |
       v
    (CADSOL) -----> Sincronização Governamental / Institucional

##### 5.2 Integridade do Ledger e Tradução Cultural (Sequence)
Fluxo de como uma ação simples do usuário gera registros fiscais imutáveis.

    Trabalhador            PDV UI                 Core Lume              SQLite
        |                    |                        |                    |
        | 1. Vende produto   |                        |                    |
        |------------------->|                        |                    |
        |                    | 2. Auxilia no preço    |                    |
        |                    |----------------------->|                    |
        |                    | 3. Call RecordSale API |                    |
        |                    |----------------------->|                    |
        |                    |                        | 4. Valida int64    |
        |                    |                        |------------------->|
        |                    |                        | 5. Postings (D/C)  |
        |                    |                        |------------------->|
        |                    |                        | 6. Confirma (Hash) |
        |                    |                        |<-------------------|
        |                    | 7. Resposta Sucesso    |                    |
        |                    |<-----------------------|                    |
        | 8. "Venda Ok!"     |                        |                    |
        |<-------------------|                        |                    |

--------------------------------------------------------------------------------

#### 6. Stack Tecnológico Final
| Camada | Tecnologia | Uso |
| ------ | ------ | ------ |
| Backend | Go 1.22+ | API REST, concorrência, binário leve |
| Storage | SQLite3 | Isolamento por tenant (Soberania) |
| Front/Web| HTMX + Tailwind | Interface PWA mobile-first pedagógica |
| Hash | SHA256 | Auditoria CADSOL e imutabilidade |
| Numerics | int64 | Centavos e minutos de trabalho |
| Documents| Markdown | Atas de Assembleia |

--------------------------------------------------------------------------------

#### 7. Estratégia de Módulos para Agentes
Para mitigar a entropia de contexto em agentes de IA, o sistema Digna é dividido em módulos independentes.

| Módulo | Consome de | Entrega para | Regra Sócio-Técnica para IA |
| ------ | ------ | ------ | ------ |
| Lifecycle | Sist. de Arquivos| Lume / Facade | Nunca misturar tenants (dados cruzados). |
| Ledger | Lifecycle (DB) | Reporting | Operar apenas em int64; Garantir soma zero. |
| Legal Facade | Lifecycle (Meta) | Usuário Final | Gerar docs fáceis de ler em Assembleia. |
| Reporting | Ledger (History) | Usuário Final | Painéis visuais para aprovação democrática. |
| PDV UI | Lume (API) | Trabalhador | **PROIBIDO usar jargões contábeis na tela.** |
```

***
