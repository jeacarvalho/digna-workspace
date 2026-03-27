title: Arquitetura do Sistema Digna
status: implemented
version: 2.0
last_updated: 2026-03-27
---

# Arquitetura do Sistema - Ecossistema Digna

> **Nota:** Este documento reflete a arquitetura integrada do Ecossistema Digna (PDF v1.0), preservando toda a infraestrutura validada nas Sprints 1-16, incorporando os novos módulos como Fases 3 e 4, e incluindo o Sistema de Ajuda Educativa (RF-30) decidido na sessão de 27/03/2026.

**Projeto:** Sistema de Gestão Contábil e Pedagógica para Economia Solidária  
**Arquitetura:** Local-First Server-Side com Micro-databases Isolados + Domain-Driven Design + Ecossistema de 4 Módulos

---

## 1. Topologia de Dados: A Dualidade Arquitetural

A premissa fundamental do Digna é que o dado pertence ao produtor, não à plataforma. Para garantir isso e, ao mesmo tempo, permitir a gestão do ecossistema, a arquitetura de banco de dados é estritamente dividida em duas esferas:

### 1.1. O Banco de Dados Central (Governança do Ecossistema)

**O que é:** Um arquivo de banco de dados único (ex: `data/central.db`), gerido exclusivamente pelo módulo `lifecycle`.

**Responsabilidade:** Atua como o "Agregador Central" e motor de identidade do projeto como um todo.

**O que armazena:**
- Gestão de Identidade Global (usuários Gov.br, CPFs/CNPJs)
- Mapeamento físico e chaves criptográficas dos bancos dos Tenants
- Relacionamentos Cross-Tenant (RF-12): Tabelas estruturais do ecossistema, como a `EnterpriseAccountant`, que define qual Contador Social (Identidade Global) tem permissão de acesso a qual Empreendimento (Tenant), com suas respectivas datas de vigência
- **NOVO - Módulo 2:** Cache de indicadores econômicos (`indicators{}`)
- **NOVO - Módulo 3:** Catálogo de programas de financiamento (`financing_programs{}`)
- **NOVO - RF-30:** Tópicos de ajuda educativa (`help_topics{}`)
- Metadados de intercooperação institucional

**Regra Inegociável:** O Banco Central **jamais** armazena transações financeiras, itens de estoque ou detalhamento operacional das entidades.

### 1.2. O Banco de Dados do Tenant (Soberania do Empreendimento)

**O que é:** O banco de dados físico isolado de cada empreendimento de economia solidária (ex: `data/entities/{entity_id}.db`).

**Responsabilidade:** Materializar o Requisito Não Funcional de Soberania (RNF-01). Se o grupo decidir deixar a rede, ele simplesmente leva o seu arquivo `.db` embora.

**O que armazena:**
- O **Ledger** contábil (partidas dobradas, histórico de caixa)
- Estoque, compras, vendas (PDV)
- Livro de Atas e Decisões de Assembleia (com Hashes SHA256)
- Registro de horas trabalhadas (Primazia do Trabalho - ITG 2002)
- **NOVO - Módulo 3:** Perfil de elegibilidade (`eligibility_profiles{}`)
- **NOVO - Módulo 3:** Match de programas (`program_matches{}`)
- **NOVO - Módulo 4:** Perfil público (`public_profiles{}`), mural de necessidades (`need_posts{}`)

**Regra Inegociável:** É tecnicamente impossível e proibido realizar **JOINs** (cruzamento de dados) entre o banco de um Tenant e o banco de outro Tenant.

---

## 2. Arquitetura do Ecossistema (4 Módulos) [ATUALIZADO - PDF v1.0]

O Ecossistema Digna é composto por **quatro módulos interdependentes**. O **digna ERP** é o núcleo — os demais módulos são extensões que ampliam o valor gerado a partir dos dados que o ERP já captura.

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    ECOSSISTEMA DIGNA — ARQUITETURA INTEGRADA            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐   │
│  │   MÓDULO 1      │     │   MÓDULO 2      │     │   MÓDULO 3      │   │
│  │   digna ERP     │────▶│   Motor de      │────▶│   Portal de     │   │
│  │   (✅ 85%)      │     │   Indicadores   │     │   Oportunidades │   │
│  │                 │     │   (📋 NOVO)     │     │   (📋 NOVO)     │   │
│  └─────────────────┘     └─────────────────┘     └─────────────────┘   │
│           │                       │                       │           │
│           │                       │                       │           │
│           ▼                       ▼                       ▼           │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                    MÓDULO 4 — Rede Digna                        │   │
│  │                    (Marketplace Solidário)                      │   │
│  │                    (🔄 Expandir sync_engine)                    │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │              SISTEMA TRANSVERSAL: Ajuda Educativa (RF-30)       │   │
│  │              (📋 NOVO - Decisão 27/03/2026)                     │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

| Módulo | Função Principal | Integração-Chave | Status |
|--------|------------------|------------------|--------|
| **Módulo 1: digna ERP** | Gestão financeira, fiscal e contábil do empreendedor | Alimenta todos os demais módulos com dados do perfil | ✅ 85% Completo |
| **Módulo 2: Motor de Indicadores** | Coleta e interpreta indicadores econômicos em tempo real | Consome APIs BCB/IBGE; alimenta Portal com taxas e contexto | 📋 Backlog Fase 3 |
| **Módulo 3: Portal de Oportunidades** | Match automático entre perfil e programas de financiamento | Consome ERP + Motor; gera checklist de documentos | 📋 Backlog Fase 3 |
| **Módulo 4: Rede Digna** | Marketplace solidário entre entidades do ecossistema | Consome perfil do ERP para matching de compra/venda | 🔄 Expandir sync_engine |
| **Sistema Transversal: Ajuda Educativa** | Tradução de conceitos técnicos em linguagem popular | Linkagem UI → banco de ajuda (`help_topics{}`) | 📋 Backlog RF-30 |

---

## 3. Módulos Principais (Módulos de Domínio)

### 3.1 Módulos Implementados (Sprints 1-16) ✅

| Módulo | Responsabilidade | Status | Sprint |
|--------|------------------|--------|--------|
| **core_lume** | Motor contábil invisível (Ledger, WorkLog, Decision) | ✅ PRODUCTION | 01-07 |
| **lifecycle** | Gerenciador do ecossistema, isolamento SQLite, banco central | ✅ PRODUCTION | 01 |
| **ui_web** | Interface web principal (HTMX + Tailwind, cache-proof) | ✅ PRODUCTION | 05, 15-16 |
| **accountant_dashboard** | Aliança contábil, exportação SPED, visão multi-tenant | ✅ PRODUCTION | 09-12 |
| **legal_facade** | Governança, formalização, documentos (CADSOL, atas) | ✅ PRODUCTION | 03, 11 |
| **supply** | Compras, estoque, fornecedores | ✅ PRODUCTION | 13-14 |
| **budget** | Gestão orçamentária e planejamento financeiro | ✅ PRODUCTION | 14 |
| **cash_flow** | Gestão de caixa | ✅ PRODUCTION | 06 |
| **pdv_ui** | Ponto de Venda pedagógico | ✅ PRODUCTION | 02 |
| **distribution** | Rateio de sobras (Reserva Legal 10% + FATES 5%) | ✅ PRODUCTION | 03 |
| **integrations** | Mocks de APIs governamentais | ✅ PRODUCTION | 08 |
| **sync_engine** | Sincronização offline-first | ✅ PRODUCTION | 04 |

### 3.2 Novos Módulos (Backlog - PDF v1.0 + Decisões da Sessão) 📋

| Módulo | Responsabilidade | Status | Fase |
|--------|------------------|--------|------|
| **indicators_engine** [NOVO] | Motor de indicadores econômicos (BCB, IBGE APIs) | 📋 BACKLOG | Fase 3 |
| **portal_opportunities** [NOVO] | Match automático de programas de financiamento | 📋 BACKLOG | Fase 3 |
| **rede_digna** [NOVO] | Marketplace solidário, perfil público, mural | 🔄 EXPANDIR | Fase 4 |
| **tax_compliance** [NOVO] | EFD-Reinf, ECF, blindagem tributária | 📋 BACKLOG | Fase 2 |
| **sanitary_compliance** [NOVO] | MTSE/MAPA para agroindústrias | 📋 BACKLOG | Fase 2 |
| **help_engine** [NOVO - RF-30] | Sistema de ajuda educativa estruturada | 📋 BACKLOG | Transversal |

---

## 4. Arquitetura DDD

### 4.1 Repository Pattern

Exemplo de como o domínio contábil é blindado da infraestrutura física:

```go
// Domain Layer (core_lume/internal/domain/)
type LedgerEntry struct {
    ID          string
    EntityID    string
    Amount      int64  // Anti-Float: centavos
    Description string
    CreatedAt   int64
}

// Repository Interface (core_lume/internal/repository/)
type LedgerRepository interface {
    Save(entry *LedgerEntry) error
    FindByEntity(entityID string) ([]*LedgerEntry, error)
}

// Repository Implementation (core_lume/internal/repository/sqlite.go)
type SQLiteLedgerRepository struct {
    lm lifecycle.LifecycleManager
}

func (r *SQLiteLedgerRepository) Save(entry *LedgerEntry) error {
    db, err := r.lm.GetDatabase(entry.EntityID) // Isolamento por entidade
    // ...
}
```

### 4.2 Fluxo de Dados entre Módulos [NOVO - PDF v1.0]

O fluxo parte sempre do **digna ERP**, que é a fonte de verdade sobre o perfil de cada entidade.

1. **ERP captura o perfil:** No uso cotidiano (vendas, compras, caixa), o Digna constrói automaticamente o perfil completo da entidade: faturamento, CNAE, município, regime tributário, situação fiscal.
2. **Motor coleta indicadores:** Diariamente, o Motor atualiza SELIC, IPCA, câmbio, taxas de crédito e expectativas de mercado via APIs públicas. Esses dados ficam disponíveis para todos os módulos.
3. **Portal executa o match:** Com o perfil do ERP e as taxas do Motor, o Portal calcula automaticamente para quais programas de financiamento a entidade é elegível e em que condições.
4. **Contador Social consolida:** O Contador Social adotado pela entidade tem visão consolidada do perfil de elegibilidade e pode submeter candidaturas aos programas identificados pelo Portal.
5. **Rede cria conexões:** Com o perfil estabelecido, a Rede identifica potenciais parceiros comerciais dentro do ecossistema e sugere conexões de compra e venda.
6. **Ajuda Educativa contextualiza:** Em qualquer ponto do fluxo, o usuário pode acessar explicações em linguagem popular via botão "?" linkado ao `help_engine`.
7. **Alertas fecham o ciclo:** Quando um edital novo é identificado ou um prazo se aproxima, o sistema notifica a entidade e o Contador Social responsável.

---

## 5. Princípio Arquitetural Central [NOVO - PDF v1.0]

> **"Nenhum usuário precisa preencher o mesmo dado duas vezes."**

O que o ERP já sabe sobre a entidade é automaticamente aproveitado pelo Portal, pelo Motor e pela Rede. Isso reduz a fricção e elimina a principal barreira de uso: o excesso de formulários.

**Implicações Técnicas:**
- `EligibilityProfile` copia dados do `Enterprise` (não duplica entrada manual)
- `PublicProfile` deriva dados do `Enterprise` + campos públicos explícitos
- `ProgramMatch` consome `EligibilityProfile` + `EconomicIndicator`
- Banco central armazena dados globais; banco da entidade armazena dados operacionais
- `HelpTopic` armazena explicações em linguagem popular linkadas a campos técnicos

---

## 6. Tecnologias Core

| Camada | Tecnologia | Justificativa |
|--------|------------|---------------|
| Backend | Go (1.22+) | Performance, concorrência, binário estático |
| Database | SQLite3 | Isolamento total por arquivo (Soberania) |
| Arquitetura | Clean Arch + DDD | Domínio independente de frameworks |
| Numerics | int64 (exclusivo) | Valores financeiros e horas (ITG 2002) exatos |
| Frontend | HTMX + Tailwind | PWA mobile-first para o trabalhador |
| Hash | SHA256 | Auditoria CADSOL e imutabilidade |
| Gov Integration [NOVO] | API Gov.br / ICP-Brasil | Assinaturas Eletrônicas Qualificadas e Autenticação |
| Tax Compliance [NOVO] | XML / Web Services | Processamento síncrono/assíncrono da EFD-Reinf |
| Economic APIs [NOVO] | BCB SGS, PTAX, Focus / IBGE SIDRA | Motor de Indicadores (Módulo 2) |
| Help System [NOVO - RF-30] | HTMX + central.db | Sistema de ajuda educativa com linkagem UI |

---

## 7. Módulos Implementados e Sprints

| Sprint | Módulo | Status | Testes | Observações |
|--------|--------|--------|--------|-------------|
| 01 | lifecycle | ✅ | 6/6 | - |
| 02 | core_lume + pdv_ui | ✅ | 8/8 | - |
| 03 | reporting + legal_facade | ⚠️ | 8/8 | Funcionalidades básicas implementadas, modularização pendente |
| 04 | sync_engine | ⚠️ | 9/9 | Funcionalidades básicas implementadas, UI e integração pendente |
| 05 | ui_web | ✅ | 9/9 | - |
| 06 | cash_flow | ✅ | 3/3 | - |
| 07 | DDD Refactoring | ✅ | 43/43 | - |
| 08 | integrations | ✅ | 5/5 | - |
| 09 | accountant_dashboard | ✅ | 8/8 | - |
| 10 | member_management | ⚠️ | 19/19 | IMPLEMENTADO DE FORMA ESPALHADA - Ver backlog para modularização |
| 11 | formalization_e2e | ✅ | 5/5 | Implementado no módulo legal_facade |
| 12 | accountant_dashboard_complete | ✅ | 8/8 | - |
| 13 | supply (Gestão de Compras e Estoque) | ✅ | 6/6 | RF-07, RF-08 |
| 14 | budget (Gestão Orçamentária) | ✅ | 4/4 | RF-10 |
| 15 | Correções Críticas + Testes E2E | ✅ | 3/3 + E2E | Validação PDV→Estoque→Caixa + Playwright |
| 16 | Identidade Visual e Sistema 100% Funcional | ✅ | 149/149 | RNF-07 completo |

---

## 8. Princípios SOLID Aplicados

| Princípio | Aplicação no Digna |
|-----------|-------------------|
| **SRP (Single Responsibility)** | Cada módulo tem uma única responsabilidade (ex: pdv_ui traduz cultura, accountant_dashboard traduz obrigações fiscais) |
| **OCP (Open/Closed)** | Sistema aberto para extensão (novas integrações governamentais sem mudar o Core Lume) |
| **LSP (Liskov Substitution)** | Implementações de Repository são intercambiáveis (SQLite ↔ Mock ↔ HTTP) |
| **ISP (Interface Segregation)** | Interfaces pequenas e específicas na camada de domínio |
| **DIP (Dependency Inversion)** | Services (ex: LedgerService) dependem de abstrações (interfaces), não de implementações concretas do SQLite |

---

## 9. Stack Tecnológico Final

| Camada | Tecnologia | Uso |
|--------|------------|-----|
| Backend | Go 1.22+ | API REST, binário leve, concorrência |
| Storage | SQLite3 | Isolamento por tenant |
| Front/Web | HTMX + Tailwind | PWA mobile-first para o trabalhador |
| Front/Dashboard | Vue/React ou HTMX | Visão Multi-tenant para o Contador |
| Hash | SHA256 | Auditoria CADSOL e imutabilidade |
| Numerics | int64 | Centavos monetários e minutos trabalhados |
| Architecture | Clean Arch + DDD | Domínio protegido |
| Fiscal | SPED / CSV Export | Ponte com sistemas contábeis comerciais |
| Gov Integration [NOVO] | API Gov.br / ICP-Brasil | Assinaturas Eletrônicas Qualificadas e Autenticação |
| Tax Compliance [NOVO] | XML / Web Services | Processamento síncrono/assíncrono da EFD-Reinf |
| Economic APIs [NOVO] | BCB SGS, PTAX, Focus / IBGE SIDRA | Motor de Indicadores (Módulo 2) |
| Help System [NOVO - RF-30] | HTMX + central.db | Sistema de ajuda educativa com linkagem UI |

---

## 10. Segurança e Soberania

### 10.1 Isolamento de Dados

Cada entidade possui banco próprio:
- **Path:** `data/entities/{entity_id}.db`
- **Isolamento:** Físico total
- **Acesso Cruzado:** Proibido entre tenants
- **Banco Central:** `data/entities/central.db` para relações inter-tenant (RF-12, indicadores, programas, help_topics)

### 10.2 Acesso do Contador Social (Painel Multi-tenant)

O acesso do contador parceiro aos dados do empreendimento ocorre estritamente em modo de leitura (Read-Only) e mediante delegação de acesso prévia. O painel apenas consulta e compila as transações, sem nunca quebrar o isolamento do arquivo local `.sqlite`.

### 10.3 Integridade

- **Hash SHA256 para auditoria:** Cada decisão gera hash do conteúdo
- **Chain digest:** Cada bloco contábil gera hash de integridade
- **Imutabilidade:** Garantida por design

### 10.4 Transporte

- **Pacotes assinados digitalmente:** Assinatura com EntityID
- **Verificação de integridade:** Hash validation
- **Non-repudiation:** Timestamp + nonce

---

## 11. Protocolo de Sincronização [ATUALIZADO]

### 11.1 Modelo

**Estratégia:** Delta-based synchronization

O sistema detecta alterações desde a última sincronização e transmite apenas os deltas, não os dados completos.

### 11.2 Estrutura do Pacote de Sync

```json
{
   "entity_id": "cooperativa_mel",
   "timestamp": 1772856840,
   "chain_digest": "d51e6eb402a6984e",
   "signature": "f802343da66e8396",
   "aggregated_data": {
     "total_sales": 7500,
     "total_work_hours": 12,
     "total_members": 2,
     "legal_status": "DREAM",
     "decision_count": 0,
     "eligibility_complete": false,
     "public_profile_published": false
  },
   "delta_count": 3,
   "module_versions": {
     "erp": "1.0",
     "indicators": "1.0",
     "portal": "1.0",
     "rede": "1.0",
     "help": "1.0"
   }
}
```

### 11.3 Processo de Sincronização

1. **DETECT:** Query deltas desde last_sync_at
   - entries: alterações em lançamentos
   - work_logs: novos registros de trabalho
   - decisions_log: novas decisões
   - fiscal_exports: novos lotes extraídos pelo Contador Social
   - eligibility_profile: atualizações no perfil de elegibilidade [NOVO]
   - public_profile: publicações na Rede Digna [NOVO]
   - need_posts: necessidades publicadas na Rede [NOVO]

2. **AGGREGATE:** Calcular métricas agregadas
   - Soma de vendas (total_sales)
   - Soma de horas (total_work_hours)
   - Contagem de membros (total_members)
   - Status atual (legal_status)
   - Completude do perfil de elegibilidade (eligibility_complete) [NOVO]
   - Perfil público publicado (public_profile_published) [NOVO]

3. **HASH:** Gerar chain digest
   - SHA256 da cadeia contábil atual
   - Inclui todos os hashes de decisões
   - Inclui hash do perfil de elegibilidade [NOVO]

4. **SIGN:** Assinar pacote
   - Usar entity_id como chave
   - Gera signature para autenticidade

5. **TRANSMIT:** Enviar para agregador
   - JSON ~500 bytes (expandido para novos módulos)
   - Apenas dados agregados
   - Dados sensíveis nunca transmitidos

### 11.4 Privacidade - Campos Incluídos vs Protegidos [ATUALIZADO]

| Campo | Incluído | Descrição |
|-------|----------|-----------|
| entity_id | ✅ | ID da entidade |
| total_sales | ✅ | Total vendas (int64) |
| total_work_hours | ✅ | Total horas |
| total_members | ✅ | Quantidade sócios |
| legal_status | ✅ | DREAM ou FORMALIZED |
| chain_digest | ✅ | Hash de integridade |
| signature | ✅ | Assinatura digital |
| fiscal_batch_hash | ✅ | Hash de integridade do último Lote SPED |
| eligibility_complete | ✅ | Perfil de elegibilidade completo (bool) |
| public_profile_published | ✅ | Perfil público publicado na Rede (bool) |
| credit_matches_count | ✅ | Quantidade de matches de crédito encontrados |
| **member_id** | ❌ | Dados sensíveis protegidos |
| **entry_details** | ❌ | Transações detalhadas |
| **posting_id** | ❌ | IDs internos |
| **cadunico_status** | ❌ | Status CadÚnico (sensível) |
| **inadimplencia** | ❌ | Status de inadimplência (sensível) |
| **credit_purpose** | ❌ | Finalidade do crédito (sensível) |

---

## 12. Modularização Pendente ⚠️

**Status:** Algumas funcionalidades foram implementadas de forma distribuída entre múltiplos módulos, violando o princípio SRP.

| Módulo | Status Atual | Prioridade | Esforço Estimado |
|--------|-------------|------------|------------------|
| member_management | ⚠️ Espalhado | **ALTA** | 2-3 dias |
| reporting | ⚠️ Básico | MÉDIA | 2-3 dias |
| sync_engine | ⚠️ Isolado | MÉDIA | 2-3 dias |

**Ver `docs/NEXT_STEPS.md` para detalhes completos do backlog de modularização.**

---

## 13. Decisões Arquiteturais Críticas [NOVO - PDF v1.0 + Sessão 27/03/2026]

### 13.1 Laicidade do Produto vs. Fundamento Teológico

**Decisão:** Manter abordagem secular na interface. A teologia informa decisões de design internamente, mas o produto permanece acessível independentemente de crença.

**Implementação:** Documentar princípios em `docs/04_governance/` sem expor ao usuário final.

### 13.2 Emissão de NF-e/NFC-e

**Decisão:** Manter fora do Core Lume. Criar módulo `fiscal_bridge` separado que pode ser substituído por integrações de terceiros (ex: eNotas, NFe.io).

**Justificativa:** Não acoplar ao Motor Lume para preservar essência de "Contabilidade Invisível".

### 13.3 Arquitetura de Módulos

**Decisão:** Manter arquitetura atual (Go workspace, monolito modular).

**Justificativa:** Go workspace permite separação lógica sem complexidade de deploy. Extrair `indicators_engine` como módulo separado, mas mesmo binário.

### 13.4 Banco de Dados para Novos Módulos

| Módulo | Banco | Justificativa |
|--------|-------|---------------|
| indicators_engine | central.db | Dados globais, não específicos por entidade |
| financing_programs | central.db | Catálogo global de programas |
| eligibility_profiles | entity.db | Perfil específico de cada entidade |
| program_matches | entity.db | Match específico de cada entidade |
| public_profiles | entity.db + sync | Perfil público com sincronização controlada |
| help_topics | central.db | Tópicos de ajuda globais, reutilizáveis |

### 13.5 Sistema de Ajuda Educativa (RF-30) [NOVO - Decisão 27/03/2026]

**Decisão:** Implementar sistema de ajuda estruturada com linkagem UI → banco de ajuda.

**Justificativa:** Campos como "CadÚnico", "Inadimplência", "CNAE" são jargões burocráticos que violam o Pilar Pedagógico do Digna.

**Implementação:**
- Tabela `help_topics` no `central.db`
- Botão "?" ao lado de campos técnicos na UI
- Explicação em linguagem popular + legislação + próximo passo acionável
- Carregamento via HTMX (< 500ms)

---

## 14. Riscos Arquiteturais e Mitigações

| Risco | Probabilidade | Impacto | Mitigação |
|-------|--------------|---------|-----------|
| APIs governamentais instáveis | Alta | Médio | Cache local + circuit breaker + modo offline |
| Complexidade do Portal cresce além do MVP | Média | Alto | MVP com 3 programas primeiro; validação com usuários reais |
| Conflito de naming (ERP vs. Ecossistema) | Baixa | Baixo | Documentar claramente a hierarquia de módulos |
| Massa crítica para Rede Digna não atingida | Alta | Médio | Focar em ERP + Portal primeiro; Rede como "nice-to-have" |
| Teologia afeta adoção secular | Média | Alto | Manter produto laico na interface; teologia informa design internamente |
| Dependência de contadores sociais para escala | Média | Alto | Criar programa de capacitação + certificação CFC |
| **Linguagem muito técnica nos tópicos de ajuda (RF-30)** | Alta | Alto | Revisão por ITCPs/comunidade; teste de usabilidade com usuários reais |
| **Conteúdo de ajuda desatualizado** | Média | Médio | Processo de atualização via central.db, não hardcoded |

---

## 15. Próximos Passos Arquiteturais

1. **Criar módulo `indicators_engine`** (RF-18) - Estrutura base com collector, cache, interpreter
2. **Adicionar tabela `eligibility_profiles`** (RF-19) - Migration no lifecycle
3. **Expandir `sync_engine`** para Rede Digna (RF-24 a RF-26)
4. **Implementar `tax_compliance`** (RF-14) - EFD-Reinf, ECF
5. **Implementar `sanitary_compliance`** (RF-16) - MTSE/MAPA
6. **Implementar `help_engine`** (RF-30) - Sistema de ajuda educativa com linkagem UI

---

**Status:** ✅ ATUALIZADO COM VISÃO DE ECOSSISTEMA (PDF v1.0) + RF-30 (Decisão de Design 27/03/2026)  
**Próxima Ação:** Atualizar `03_architecture/02_protocols.md` com novos protocolos de integração  
**Versão Anterior:** 1.6 (2026-03-13)  
**Versão Atual:** 2.0 (2026-03-27)
