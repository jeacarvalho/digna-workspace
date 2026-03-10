# 📚 Documentação do Projeto Digna

**Versão:** 1.5 (Sistema 100% Funcional)
**Última Atualização:** 2026-03-09 (Sessão Crítica de Correções)
**Projeto:** Sistema de Gestão Contábil, Institucional e Pedagógica para Economia Solidária
**Mantenedor:** Fundação Providentia

## 🚨 **ATUALIZAÇÃO DE STATUS - SESSÃO 09/03/2026**

### 🎯 **OBJETIVO ALCANÇADO: SISTEMA 100% FUNCIONAL**
Após extensivo trabalho de correção de bugs críticos e implementação de identidade visual completa, o sistema Digna está **operacional em todos os módulos**.

### ✅ **MÓDULOS OPERACIONAIS:**
1. **Autenticação** - Login com identidade visual completa
2. **Dashboard** - Métricas integradas e navegação
3. **PDV** - Ponto de Venda com carrinho funcional
4. **Caixa** - Gestão financeira completa
5. **Compras** - Cadastro de fornecedores e compras
6. **Estoque** - Controle de inventário em tempo real
7. **Ponto Social** - Registro de horas (ITG 2002)

### 🔧 **PROBLEMAS CRÍTICOS RESOLVIDOS:**
1. **Cache persistente de templates Go** - Sistema blindado
2. **Database vazio** - Populado com dados reais (`cafe_digna`)
3. **Logo não visível** - Implementada identidade visual completa (Sprint 16)
4. **Templates parciais não renderizados** - Migração para templates simples
5. **Navegação quebrada** - Links funcionais entre todos os módulos

### 📊 **STATUS ATUAL:**
- **Servidor:** 🟢 Rodando na porta 8090
- **Health Check:** 🟢 `{"status":"ok","version":"v.0"}`
- **Database:** 🟢 Populado com dados reais
- **Templates:** 🟢 Todos funcionando
- **Identidade Visual:** 🟢 Completa ("Soberania e Suor")

**Sistema pronto para testes de produção e uso operacional!**

---

## 📋 Visão Geral

O **Digna** é uma infraestrutura contábil soberana para a Economia Solidária (EES) brasileira, mantida pela **Fundação Providentia**. Mais do que um software, o Digna é uma Tecnologia Social desenhada para combater a exclusão digital e financeira.

### Missão

> Promover a autogestão, soberania e transformação digital dos Empreendimentos de Economia Solidária no Brasil através de tecnologia livre e acessível, atuando simultaneamente como uma **ponte tecnológica inclusiva para a conformidade legal e a classe contábil**.

### Princípios Fundamentais

1. **Soberania do Dado** - Cada entidade possui seu próprio banco SQLite isolado fisicamente (Exit Power).
2. **Contabilidade Invisível** - Operações coloquiais no app geram lançamentos de partidas dobradas automaticamente no backend.
3. **Primazia do Trabalho** - Tempo de trabalho é convertido em Capital Social (Baseado na norma ITG 2002).
4. **Escala Nacional** - Arquitetura desenhada para milhões de empreendimentos informais e formais.
5. **Clean Architecture** - Domínio de negócios independente de frameworks e interfaces (TDD & DDD).
6. **Aliança Contábil** - Valorização do Contador Social parceiro através de painel multi-tenant e geração de lotes fiscais (SPED) sem onerar o trabalhador.
7. **Aliança Institucional** - Previsão de envolver diversos atores da sociedade, tanto governamentais como privados, bem como o terceiro setor. O objetivo não é criar um produto para ser monetizado, mas, sim, um facilitador de sonhos e vocação. Sonhos dos empreendedores sociais, sem "força" de ação mais plena e vocação de profissionais de todos os setores citados, que poderão ver no Digna a possibilidade de algo "maior" do que apenas o seu próprio bem estar

---

## 🗂️ Estrutura de Documentação

A documentação segue o padrão PKM (Personal Knowledge Management) de alta integridade.

```text
docs/
├── 01_project/      # Visão, Escopo, Stakeholders, Riscos
├── 02_product/      # Requisitos, Modelos de Domínio, Algoritmos
├── 03_architecture/ # Arquitetura Técnica, Protocolos, Melhorias, ADRs
├── 04_governance/   # Fundação, PMC, Regras de Contribuição, Licença
├── 05_ai/           # Constituição de IA, Agentes, Padrões de Sessão
└── 06_roadmap/      # Estratégia, Roadmap, Backlog, Status
```

---

## 📖 Navegação por Seção (Índice)

### 01 - Projeto (Estratégia e Gestão)
| Documento | Descrição |
|-----------|-----------|
| [01_vision.md](./01_project/01_vision.md) | Visão estratégica do produto e aliança contábil |
| [02_scope.md](./01_project/02_scope.md) | Escopo, limites e capacidades principais |
| [03_stakeholders_risks.md](./01_project/03_stakeholders_risks.md) | Mapa de partes interessadas e matriz de riscos |

### 02 - Produto (Requisitos e Domínio)
| Documento | Descrição |
|-----------|-----------|
| [01_requirements.md](./02_product/01_requirements.md) | Requisitos Funcionais e Não Funcionais (RFs/RNFs) |
| [02_models.md](./02_product/02_models.md) | Modelos de Entidades, Schema e Algoritmos de Negócio |

### 03 - Arquitetura Técnica
| Documento | Descrição |
|-----------|-----------|
| [01_system.md](./03_architecture/01_system.md) | Arquitetura DDD, Clean Architecture e Componentes |
| [02_protocols.md](./03_architecture/02_protocols.md) | Protocolos de Sincronização, Segurança e Economia |
| [03_improvements.md](./03_architecture/03_improvements.md) | Radar de dívida técnica, riscos e melhorias futuras |
| [04_architectural_decisions.md](./03_architecture/04_architectural_decisions.md) | Registro de Decisões Arquiteturais (ADRs) |

### 04 - Governança
| Documento | Descrição |
|-----------|-----------|
| [governance.md](./04_governance/governance.md) | Fundação Providentia, Comitês, Licença (Apache 2.0) |

### 05 - IA & Agentes
| Documento | Descrição |
|-----------|-----------|
| [01_constitution.md](./05_ai/01_constitution.md) | Regras de Ouro inegociáveis para LLMs e Agentes (Anti-float) |
| [02_session.md](./05_ai/02_session.md) | Padrão obrigatório para execução de sessões de código |

### 06 - Roadmap e Tático
| Documento | Descrição |
|-----------|-----------|
| [01_strategy.md](./06_roadmap/01_strategy.md) | Fases de Release (v0 à v3) |
| [02_roadmap.md](./06_roadmap/02_roadmap.md) | Roadmap detalhado do Produto |
| [03_backlog.md](./06_roadmap/03_backlog.md) | Product Backlog priorizado |
| [04_status.md](./06_roadmap/04_status.md) | Status atual de todas as Sprints e Testes |
| [05_session_log.md](./06_roadmap/05_session_log.md) | Histórico de sessões de desenvolvimento |

---

## 🚀 Status das Sprints (Resumo)

| Sprint | Módulo | Status | Testes | Descrição |
|--------|--------|--------|--------|-----------|
| 01 | Lifecycle Manager | ✅ | 6/6 | Criação e gestão física de tenants (.sqlite) |
| 02 | Core Lume (Ledger) | ✅ | 8/8 | Motor contábil com partidas dobradas exatas em `int64` |
| 03 | Reporting + Legal | ✅ | 8/8 | Rateio social (ITG 2002) e geração de documentos |
| 04 | Sync Engine | ✅ | 9/9 | Motor de sincronização offline-first e B2B |
| 05 | UI Web (PWA) | ✅ | 9/9 | Interface mobile-first (HTMX + Tailwind) |
| 06 | Cash Flow | ✅ | 3/3 | Gestão de fluxo de caixa |
| 07 | DDD Refactoring | ✅ | 43/43 | Refatoração de Clean Architecture em todos os módulos |
| 08-09 | Integrações (Mocks) | ✅ | 13/13 | APIs Simuladas (Gov.br, CADSOL, Receita Federal) |
| 10 | Gestão de Membros | ✅ | 19/19 | Perfis de acesso, cooperados e permissões |
| 11 | Formalização e E2E | ✅ | 5/5 | Jornada completa "Sonho Solidário" testada ponta a ponta |
| 12 | **Accountant Dashboard** | ✅ | 8/8 + E2E | Painel do Contador Social e Exportação Fiscal (SPED) - COMPLETE |
| 13 | **Gestão de Compras e Estoque** | ✅ | 6/6 | Módulo completo com contabilidade invisível - COMPLETE |
| 14 | **Gestão Orçamentária** | ✅ | 4/4 | Planejamento financeiro com alertas visuais - COMPLETE |
| 15 | **Correções Críticas + Testes E2E** | ✅ | 3/3 + E2E | Validação PDV→Estoque→Caixa + Playwright - COMPLETE |

---

## 💻 Exemplo Prático de Arquitetura (O Poder do DDD)

Nossa arquitetura permite isolar totalmente as regras de negócio das tecnologias externas. Exemplo real do nosso módulo de integrações:

```go
// 1. O Serviço usa apenas a Interface (Domain Layer)
type CreditService struct {
    gov IntegrationRepository
}

// 2. A regra de negócio é agnóstica à implementação
func (s *CreditService) SolicitarCredito(...) {
    // Pode ser um Mock, pode ser HTTP real, o Core não se importa!
    simulacao, _ := s.gov.BNDES().SimularCredito(ctx, creditRequest)
}
```

**Para integrar de verdade no futuro:**
Basta criar novas implementações mantendo as **mesmas interfaces** e trocá-las na injeção de dependência:

```go
type HTTPReceitaFederalRepository struct { ... }

func (r *HTTPReceitaFederalRepository) ConsultarCNPJ(...) {
    // Chamada HTTP real para a API da Receita Federal
}
```
Sem mudar uma linha de código do Core Lume! (Princípio OCP).

---

## 🛠️ Stack Tecnológica

| Camada | Tecnologia | Justificativa |
|--------|------------|---------------|
| Backend / Motor | Go 1.22+ | Performance, concorrência nativa, binário estático único |
| Banco de Dados | SQLite3 | Isolamento físico por Tenant (Soberania do Dado) |
| Frontend (Trabalhador) | HTMX + Tailwind | PWA leve, foco no comportamento "Offline-first" |
| Frontend (Contador) | Vue/React ou HTMX | Visão agregada Multi-tenant para exportações fiscais |
| Tipagem Financeira | `int64` | Erros de arredondamento (`float`) são terminantemente proibidos |

---

## 🔗 Links Rápidos

- [Código Fonte (Módulos)](../modules/)
- [Dados de Exemplo](../data/)
- [Módulo de Integrações](../modules/integrations/)

---

## 📊 Métricas do Projeto

| Métrica | Valor Atualizado |
|---------|------------------|
| Total de Módulos | 13 (+ accountant_dashboard) |
| Total de Testes | 136/136 (100% Pass) 🎉 |
| Cobertura de Código | > 80% |
| Interfaces de Repositório | 10 (+ FiscalRepository) |
| Integrações Gov (Mocks) | 8 |
| Testes E2E Completos | Jornada Anual BDD com Contador Social (1) |
| Documentação | 100% Sincronizada |
```

***
