# 🏛️ **Projeto Providentia: Dignidade e Soberania Financeira**

**Data de Atualização:** 09/03/2026 | **Versão:** 0.4 (Social-First)

---

## 📋 **Visão Geral**

O **Projeto Providentia** é uma iniciativa da **Fundação Providentia** para criar uma infraestrutura contábil soberana para a Economia Solidária (EES) brasileira. O produto principal, **Digna**, é uma Tecnologia Social desenhada para combater a exclusão digital e financeira através de "contabilidade invisível".

### 🎯 **Missão**
> Promover a autogestão, soberania e transformação digital dos Empreendimentos de Economia Solidária no Brasil através de tecnologia livre e acessível, atuando simultaneamente como uma **ponte tecnológica inclusiva para a conformidade legal e a classe contábil**.

### 🌟 **Princípios Fundamentais**
1. **Soberania do Dado** - Cada entidade possui seu próprio banco SQLite isolado fisicamente (Exit Power)
2. **Contabilidade Invisível** - Operações coloquiais no app geram lançamentos de partidas dobradas automaticamente
3. **Primazia do Trabalho** - Tempo de trabalho é convertido em Capital Social (Baseado na norma ITG 2002)
4. **Escala Nacional** - Arquitetura desenhada para milhões de empreendimentos informais e formais
5. **Clean Architecture** - Domínio de negócios independente de frameworks e interfaces (TDD & DDD)
6. **Aliança Contábil** - Valorização do Contador Social parceiro através de painel multi-tenant
7. **Aliança Institucional** - Envolvimento de diversos atores da sociedade para facilitar sonhos e vocação

---

## 📖 **Capítulo 1: Fundamentos, Visão e Propósito**

### **1.1. O Cenário Atual e a Necessidade de Disrupção**
A economia solidária brasileira enfrenta o desafio da "invisibilidade burocrática". Pequenos empreendimentos produzem riqueza real, mas soçobram ante a complexidade contábil e a dificuldade de acesso ao crédito. O projeto Providentia nasce para romper esse ciclo, transformando a gestão em um ato de cidadania automática.

### **1.2. A Filosofia "Digna": Contabilidade Invisível**
O produto **Digna** (app/web) é a interface de libertação do empreendedor. A filosofia central é que o trabalhador não deve "fazer contabilidade", mas sim "operar seu negócio". Através de um **PDV (Ponto de Venda)** intuitivo, cada venda de mel, cada hora de costura e cada compra de insumo é capturada e traduzida instantaneamente em linguagem contábil pelo motor **Lume**.

### **1.3. O Serpro como Indutor de Soberania Nacional**
Como braço tecnológico do Estado, o Serpro atua como o grande hospedeiro desta infraestrutura, garantindo que o dado do cooperado não seja mercadoria, mas sim um ativo de soberania para a formulação de políticas públicas de crédito e fomento.

### **1.4. A Fundação Providentia: O Modelo de Governança**
A fundação serve como a guardiã da neutralidade. Ela assegura que o software permaneça um bem público, imune a capturas de mercado, e focado exclusivamente na dignidade financeira dos empreendimentos de economia solidária (EES).

### **1.5. Objetivos Estratégicos**
* **Curto Prazo (v0 - A Vitrine):** Lançamento do MVP funcional com PDV Offline-First e registro de horas de trabalho (ITG 2002)
* **Médio Prazo (Escala e Formalização):** Automação do dossiê **CADSOL** para reconhecimento de autogestão e acesso a compras públicas (PNAE/PAA)
* **Longo Prazo (Ecossistema de Crédito):** Criação de uma rede de intercooperação e crédito baseada no histórico de impacto social e financeiro

---

## 🏛️ **Capítulo 2: Modelo de Governança e Regras de Membresia**

### **2.1. O "Apache Way" à Brasileira**
A Fundação Providentia adota o modelo de governança baseada em mérito e neutralidade. O objetivo é garantir que o projeto não seja "dono" de ninguém, mas um bem público digital gerido por quem efetivamente contribui para o seu sucesso.

**Pilares da Governança:**
* **Independência de Fornecedor:** O software deve ser funcional independentemente de qualquer empresa privada específica
* **Transparência nas Decisões:** Todas as deliberações técnicas e estratégicas são registradas e públicas
* **Comunidade sobre o Código:** A saúde da rede de desenvolvedores e usuários é mais importante que qualquer funcionalidade isolada

### **2.2. Categorias de Membresia e Participação**
| Categoria | Perfil | Papel e Contribuição |
| :---- | :---- | :---- |
| **Membros Fundadores** | Empreendedores de sucesso (Pessoas Físicas) | Aporte de capital semente (Endowment), experiência e visão estratégica de escala |
| **Parceiro Institucional** | Serpro (Empresa Pública) | Provedor de infraestrutura de nuvem, segurança e integração com o Estado |
| **Membros Corporativos** | Empresas de Tecnologia/SaaS | Doação de horas de engenharia ou recursos financeiros (Modelo Platinum/Gold) |
| **Membros Acadêmicos** | Universidades e Centros de Pesquisa | Validação metodológica e formação de redes de tutoria |
| **Contribuidores** | Desenvolvedores e Contadores | Ganham poder de voto através de contribuições técnicas de alta qualidade |
| **CFC** | Contadores Sociais | Normatização, padronização, auditorias e apadrinhamento de empreendimentos solidários |

### **2.3. Propriedade Intelectual e Licenciamento**
* **Licença Apache 2.0:** Todo o código é distribuído sob licença permissiva, garantindo que o núcleo permaneça aberto
* **Gestão de Marcas:** A marca "Providentia" pertence à Fundação, protegendo o nome contra usos indevidos

### **2.4. O Ciclo do "Giving Back" (A Contrapartida Social)**
A fundação produz um **Relatório Anual de Impacto Social**, medindo a riqueza transacionada e a valoração do trabalho voluntário (conforme a norma **ITG 2002**), gerando transparência sobre o fortalecimento da rede.

---

## ⚙️ **Capítulo 3: Arquitetura e Engenharia de Software (Lume Engine)**

### **3.1. Paradigma: Micro-databases e Offline-First**
Cada entidade possui sua própria instância física de **SQLite** isolada. Para o "Brasil Profundo", a aplicação é desenhada como **Offline-First**: as transações são capturadas localmente e sincronizadas assincronamente (Push) quando houver conectividade.

### **3.2. O Motor Lume (Regras de Ouro)**
1. **Integridade Monetária (int64):** Uso estrito de inteiros para representar centavos. Floats são proibidos
2. **Partidas Dobradas Automáticas:** Validação sistemática de soma zero (Débitos + Créditos = 0)
3. **Valoração do Trabalho:** Suporte nativo ao registro de horas de trabalho militante/voluntário como ativo social

### **3.3. Stack Tecnológica de Referência**
* **Backend:** Go (1.22+)
* **Estrutura:** Multi-module Workspace (go.work)
* **Nomenclatura:** Kebab-case para diretórios; snake_case para arquivos
* **Banco de Dados:** SQLite3 (Isolamento físico por Tenant)
* **Frontend (Trabalhador):** HTMX + Tailwind (PWA leve, Offline-first)
* **Frontend (Contador):** Vue/React ou HTMX (Visão agregada Multi-tenant)
* **Tipagem Financeira:** `int64` (Erros de arredondamento `float` são terminantemente proibidos)

### **3.4. Estratégia de Modularização para Desenvolvimento Ágil**
Para mitigar a entropia de contexto e garantir a precisão da construção assistida por IA, o sistema é segmentado em cinco módulos atômicos e independentes:

1. **Módulo 0 (pdv_ui):** Interface de operação (Venda/Compra/Horas) com suporte a cache local
2. **Módulo 1 (lifecycle):** Orquestrador de arquivos SQLite e gerenciador de estados de sincronização
3. **Módulo 2 (core_lume):** Motor contábil de partidas dobradas monetárias e sociais
4. **Módulo 3 (legal_facade):** Simulador de formalização e gerador de documentos (Atas e Estatutos)
5. **Módulo 4 (reporting):** Painel de Dignidade, Balanço Social e Dossiê de Crédito

---

## 🚀 **Capítulo 4: Status do Projeto e Roadmap**

### **4.1. Status das Sprints (Resumo)**
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
| 12 | **Accountant Dashboard** | ✅ | 8/8 + E2E | Painel do Contador Social e Exportação Fiscal (SPED) |
| 13 | **Gestão de Compras e Estoque** | ✅ | 6/6 | Módulo completo com contabilidade invisível |
| 14 | **Gestão Orçamentária** | ✅ | 4/4 | Planejamento financeiro com alertas visuais |
| 15 | **Correções Críticas + Testes E2E** | ✅ | 3/3 + E2E | Validação PDV→Estoque→Caixa + Playwright |

### **4.2. Métricas do Projeto**
| Métrica | Valor Atualizado |
|---------|------------------|
| Total de Módulos | 13 (+ accountant_dashboard) |
| Total de Testes | 136/136 (100% Pass) 🎉 |
| Cobertura de Código | > 80% |
| Interfaces de Repositório | 10 (+ FiscalRepository) |
| Integrações Gov (Mocks) | 8 |
| Testes E2E Completos | Jornada Anual BDD com Contador Social (1) |
| Documentação | 100% Sincronizada |

### **4.3. Roadmap da Versão 0 (A Vitrine)**
**Objetivo:** Demonstrar a "Contabilidade Invisível" através de um PDV funcional que registra vendas e horas de trabalho, gerando automaticamente um balanço de impacto social.

**Funcionalidades Principais:**
1. **Gestão de Contexto (Lifecycle Manager):** Orquestração automática de bancos SQLite por entidade
2. **PDV Operacional:** Registro offline de transações comerciais
3. **Dossiê CADSOL:** Geração automática de atas de assembleia e logs de decisão para provar a autogestão

---

## 📚 **Estrutura de Documentação**

A documentação segue o padrão PKM (Personal Knowledge Management) de alta integridade:

```text
docs/
├── 01_project/              # Visão, Escopo, Stakeholders, Riscos
├── 02_product/              # Requisitos, Modelos de Domínio, Algoritmos
├── 03_architecture/         # Arquitetura Técnica, Protocolos, Melhorias, ADRs
├── 04_governance/           # Fundação, PMC, Regras de Contribuição, Licença
├── 05_ai/                   # Constituição de IA, Agentes, Padrões de Sessão
├── 06_roadmap/              # Estratégia, Roadmap, Backlog, Status
│   └── sprints/             # Documentação histórica de sprints
├── 07_testing/              # Estratégia de testes, padrões, scripts
├── 08_references/           # Referências externas e legais
│   ├── external/            # Documentos de orientadores externos
│   └── legal/               # Legislação sobre economia solidária
└── Providentia Foundation.md # Documento fundacional completo
```

### **Links Rápidos**
- [Documentação Completa](./docs/README.md)
- [Código Fonte (Módulos)](./modules/)
- [Dados de Exemplo](./data/)
- [Módulo de Integrações](./modules/integrations/)

---

## 💻 **Exemplo Prático de Arquitetura (O Poder do DDD)**

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

## 🧪 **Testes e Qualidade**

### **Estratégia de Testes Implementada**
1. **Testes Unitários (por módulo):** Testam funcionalidades específicas com mocks
2. **Testes de Integração (módulo ui_web):** Testam integração entre módulos via APIs
3. **Testes E2E (opcional):** Testam fluxos completos com browser real

### **Script de Execução de Testes**
```bash
# Executar testes principais de forma segura
./run_tests.sh

# Ou executar testes específicos
go test ./modules/ui_web -v -run "TestUnidadesEstoque" -timeout 30s
```

### **Princípios de Testes**
- **Isolamento:** Cada teste cria seus próprios dados com IDs únicos
- **Tolerância:** Testes verificam comportamento observável, não implementação interna
- **Setup robusto:** Funções criam dados reais via API
- **Verificações realistas:** Foco em "o sistema funciona" vs implementação detalhada

---

## 🤝 **Como Contribuir**

O Projeto Providentia é mantido pela **Fundação Providentia** e segue o modelo de governança "Apache Way à Brasileira". Contribuições são bem-vindas seguindo as regras estabelecidas na documentação de governança.

**Principais Canais:**
1. **Documentação:** [docs/04_governance/](./docs/04_governance/)
2. **Código:** [modules/](./modules/)
3. **Issues:** Verifique a documentação para o processo de contribuição

---

**Fundação Providentia** · **Tecnologia Social para Economia Solidária** · **Soberania Digital Brasileira**