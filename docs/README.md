# 📚 Documentação do Projeto Digna

**Versão:** 1.1  
**Última Atualização:** 2026-03-07  
**Projeto:** Sistema de Gestão Contábil para Economia Solidária  
**Mantenedor:** Fundação Providentia

---

## 📋 Visão Geral

O **Digna** é uma infraestrutura contábil soberana para Economia Solidária brasileira, mantido pela **Fundação Providentia**.

### Missão

> Promover a autogestão, soberania e transformação digital dos Empreendimentos de Economia Solidária (EES) no Brasil, através de tecnologia livre e acessível.

### Visão

> Ser a principal infraestrutura digital de conexão e gestão para uma rede nacional de Empreendimentos de Economia Solidária, contribuindo para a transformação social e econômica do país.

### Valores

| Valor | Descrição |
|-------|-----------|
| **Autogestão** | Respeito à decisão coletiva e管理模式 autónomo dos EES |
| **Soberania** | Controle dos dados pelos próprios empreendimentos |
| **Inclusão** | Acessibilidade para grupos historicamente marginalizados |
| **Transparência** | Clareza nos processos, dados e decisões |
| **Transformação** | Compromisso com a mudança social real |

### Princípios Fundamentais

1. **Soberania do Dado** - Cada entidade possui seu próprio banco SQLite
2. **Contabilidade Invisível** - Operações geram lançamentos automaticamente
3. **Primazia do Trabalho** - Tempo de trabalho = Capital Social (ITG 2002)
4. **Escala Nacional** - Arquitetura para milhões de empreendimentos
5. **Clean Architecture** - Domínio independente de frameworks (DDD)

---

## 🗂️ Estrutura de Documentação

```
docs/
├── 01_project/        # Visão, Escopo, Stakeholders, Riscos
├── 02_product/       # Requisitos, Modelos, Algoritmos
├── 03_architecture/  # Arquitetura Técnica, Protocolos, DDD
├── 04_governance/    # Fundação, PMC, Contribuição, Licença
├── 05_ai/            # Constituição de IA, Agentes
└── 06_roadmap/      # Estratégia, Roadmap, Backlog, Status
```

---

## 📖 Navegação por Seção

### 01 - Projeto (Gestão)
| Documento | Descrição |
|-----------|-----------|
| [01_vision.md](./01_project/01_vision.md) | Visão estratégica do produto |
| [02_scope.md](./01_project/02_scope.md) | Escopo e capacidades principais |
| [03_stakeholders_risks.md](./01_project/03_stakeholders_risks.md) | Stakeholders e riscos |

### 02 - Produto (Requisitos)
| Documento | Descrição |
|-----------|-----------|
| [01_requirements.md](./02_product/01_requirements.md) | BRD + NFR consolidados |
| [02_models.md](./02_product/02_models.md) | Domain Model + Data Model + Algoritmos |

### 03 - Arquitetura Técnica
| Documento | Descrição |
|-----------|-----------|
| [01_system.md](./03_architecture/01_system.md) | Arquitetura DDD + Clean Architecture |
| [02_protocols.md](./03_architecture/02_protocols.md) | Sync, Security, Economic |

### 04 - Governança
| Documento | Descrição |
|-----------|-----------|
| [governance.md](./04_governance/governance.md) | Fundação, PMC, Contribuição, Licença |

### 05 - IA & Agentes
| Documento | Descrição |
|-----------|-----------|
| [01_constitution.md](./05_ai/01_constitution.md) | Constituição + Agentes |
| [02_session.md](./05_ai/02_session.md) | Padrão de sessão |

### 06 - Roadmap
| Documento | Descrição |
|-----------|-----------|
| [01_strategy.md](./06_roadmap/01_strategy.md) | Estratégia de release |
| [02_roadmap.md](./06_roadmap/02_roadmap.md) | Roadmap de produto |
| [03_backlog.md](./06_roadmap/03_backlog.md) | Product Backlog |
| [04_status.md](./06_roadmap/04_status.md) | Status atual |
| [05_session_log.md](./06_roadmap/05_session_log.md) | Histórico de sessões |

---

## 🚀 Status das Sprints

| Sprint | Módulo | Status | Testes | Descrição |
|--------|--------|--------|--------|-----------|
| 01 | Lifecycle Manager | ✅ | 6/6 | Criação e gestão de tenants |
| 02 | Core Lume (Ledger) | ✅ | 8/8 | Motor contábil com partidas dobradas |
| 03 | Reporting + Legal | ✅ | 8/8 | Rateio social e documentação |
| 04 | Sync Engine | ✅ | 9/9 | Sincronização offline-first |
| 05 | UI Web (PWA) | ✅ | 9/9 | Interface mobile-first |
| 06 | Cash Flow | ✅ | 3/3 | Gestão de caixa |
| 07 | **DDD Refactoring** | ✅ | 43/43 | Aplicado DDD a todos os módulos |
| 08 | **Integrações** | ✅ | 5/5 | Interfaces gov (mock) |
| **Total** | | | **91/91** | **100% PASS** 🎉 |

---

## 🏛️ Arquitetura DDD

O projeto segue rigorosamente **Domain-Driven Design (DDD)** e **Clean Architecture**:

### Repository Pattern
```go
// Domain Layer - Interface pura
type LedgerRepository interface {
    SaveEntry(entry *Entry) (int64, error)
    GetBalance(accountID int64) (int64, error)
}

// Infrastructure Layer - Implementação SQLite
type SQLiteLedgerRepository struct { ... }
```

### Integrações Externas (Módulo integrations/)

O sistema possui **8 interfaces de integração governamental**, todas implementadas com mocks realistas:

| Órgão | Serviços | Status |
|-------|----------|--------|
| **Receita Federal** | Consultar CNPJ, Emitir DARF | ✅ Mock |
| **MTE** | CAT, RAIS, eSocial | ✅ Mock |
| **MDS** | CadÚnico, Relatório Social | ✅ Mock |
| **IBGE** | Pesquisas, PAM, CNAE | ✅ Mock |
| **SEFAZ** | NFe, NFS-e, Manifesto | ✅ Mock |
| **BNDES** | Linhas de Crédito, Simulação | ✅ Mock |
| **SEBRAE** | Cursos, Consultoria | ✅ Mock |
| **Providentia** | Sync, Marketplace | ✅ Mock |

**Como usar:**
```go
service, _ := integrations.NewMockIntegrationService(db)

// Consultar CNPJ
cnpjData, _ := service.ReceitaFederal().ConsultarCNPJ(ctx, "12345678000190")

// Emitir NFe
nfe, _ := service.SEFAZ().EmitirNFe(ctx, nfeRequest)

// Simular crédito
simulacao, _ := service.BNDES().SimularCredito(ctx, creditRequest)
```

**Para integrar de verdade no futuro:**
Basta criar novas implementações mantendo as **mesmas interfaces**:
```go
type HTTPReceitaFederalRepository struct { ... }
func (r *HTTPReceitaFederalRepository) ConsultarCNPJ(...) { 
    // Chamada HTTP real para API da Receita Federal
}
```

Sem mudar uma linha do código que usa as integrações! (Princípio OCP)

---

## 🛠️ Stack Tecnológica

| Camada | Tecnologia | Justificativa |
|--------|------------|---------------|
| Backend | Go 1.22+ | Performance, concorrência, binário leve |
| Database | SQLite3 | Isolamento por tenant (Soberania) |
| Numerics | int64 (centavos) | Precisão financeira exata |
| Frontend | HTMX + Tailwind CSS | PWA mobile-first pedagógica |
| Architecture | Clean Arch + DDD | Domínio independente de frameworks |
| Hash | SHA256 | Auditoria CADSOL |
| Documents | Markdown | Atas de Assembleia |

---

## 📚 Referências Legais e Normativas

### Legislação
- **Lei nº 15.068/2024** (Lei Paul Singer) - Marco Legal da Economia Solidária
- **Constituição Federal** - Artigos 3º e 4º (Dignidade da pessoa humana, trabalho)
- **CLT** - Artigos 442-A a 442-O (Cooperativas)

### Normas Contábeis
- **ITG 2002** (CFC) - NBC T 19.51 - Contabilidade para Entidades sem Finalidade de Lucro
- **NBC TG 1000** - Contabilidade para Pequenas e Médias Empresas

### Registros Públicos
- **CADSOL** - Cadastro Nacional de Economia Solidária (MTE)
- **DCSOL** - Declaração de Economia Solidária

### Infraestrutura
- **Serpro** - Infraestrutura de nuvem soberana do governo federal

---

## 🔗 Links Rápidos

- [Código Fonte](../modules/)
- [Dados de Exemplo](../data/)
- [Módulo de Integrações](../modules/integrations/)

---

## 📊 Métricas do Projeto

| Métrica | Valor |
|---------|-------|
| Total de Módulos | 9 |
| Total de Testes | 91/91 (100%) |
| Linhas de Código | ~15.000 |
| Interfaces Repository | 8 |
| Integrações Gov | 8 (mock) |
| Documentação | 100% atualizada |

---

*Esta documentação segue o padrão PKM (Personal Knowledge Management) de alta integridade.*
*Última atualização: Session 009 - DDD Refactoring & Integrações*
