title: Visão Estratégica
status: implemented
version: 2.0
last_updated: 2026-03-13
---

# Visão Estratégica - Ecossistema Digna

## 1. Introdução

O **Digna** é um ecossistema de soberania financeira desenhado para a Economia Solidária (EES) e o empreendedorismo popular no Brasil. Ele não é um "ERP" (Enterprise Resource Planning) imposto de cima para baixo com lógica extrativista; é uma **Tecnologia Social e um Protocolo de Emancipação**.

Seu propósito é transformar a contabilidade — historicamente vista como um fardo burocrático — em um subproduto invisível da operação diária, atuando simultaneamente como uma ferramenta pedagógica que transforma a atividade produtiva em cidadania digital.

Além de proteger e educar o trabalhador, o Digna atua como uma **Ponte Tecnológica Institucional**, conectando a realidade informal da base produtiva às exigências da Receita Federal e do Conselho Federal de Contabilidade (CFC), através da viabilização da "Contabilidade Popular" gerida por contadores parceiros.

**Diferencial Estratégico:** O Ecossistema Digna combina três elementos raramente encontrados juntos: (1) um ERP que já captura o perfil real do negócio, eliminando formulários redundantes; (2) um motor de dados que mantém tudo atualizado automaticamente; e (3) uma rede de contadores e líderes comunitários como agentes de distribuição e capacitação.

---

## 1.1 Declaração de Posicionamento Estratégico

> "A informação que emancipa o empreendedor rico já existe — está nas APIs do governo, nos editais publicados, nas linhas de crédito abertas. O que falta é uma camada de tradução, curadoria e entrega que chegue até o empreendedor pobre na linguagem e no momento certo."

O Ecossistema Digna reduz as barreiras de acesso à informação normativa, ao crédito e à gestão financeira para microempreendedores, pequenos negócios, cooperativas e organizações da sociedade civil que atuam em contextos de vulnerabilidade socioeconômica.

---

## 2. Pilares de Design (As Leis Sociotécnicas do Sistema)

### Pilar 1: Soberania do Dado e Poder de Saída (Exit Power)
O dado não pertence à "nuvem" de uma corporação, pertence à entidade produtiva. O dado reside em um arquivo SQLite isolado fisicamente por empreendimento. O usuário detém o poder absoluto de auditar, copiar ou sair do sistema levando toda a sua história com ele.

### Pilar 2: Contabilidade Invisível e Tradução Cultural
A interface humana (Frontend) foca na ação coloquial (vender, comprar, trabalhar) e atua como uma barreira contra jargões contábeis. O débito e o crédito (Partidas Dobradas) são subprodutos gerados automaticamente pelo Motor Lume no backend.

### Pilar 3: Primazia do Trabalho (ITG 2002)
O sistema inverte a lógica capitalista: o suor (tempo/horas trabalhadas) vale tanto ou mais que o capital investido (R$). O tempo registrado (em minutos/int64) constitui o Capital Social de Trabalho e é a base para o rateio justo de sobras.

### Pilar 4: Transição Institucional Gradual (Sem Burocracia Forçada)
O Digna respeita o tempo social do grupo. Ele atua como um facilitador da conformidade (CADSOL/Sinaes), gerando atas e relatórios, mas não impõe a formalização precoce a grupos informais ("Sonhos") que ainda estão construindo sua confiança e coesão política.

### Pilar 5: Ferramenta Pedagógica e Design Participativo
O software ensina enquanto é operado. Ele auxilia visualmente o trabalhador na formação correta do seu preço (custo de insumos + hora trabalhada). Todo o seu desenvolvimento deve ser validado **com** os trabalhadores e Incubadoras (ITCPs).

### Pilar 6: Aliança Contábil e Responsabilidade Técnica [NOVO]
O sistema não exclui o contador; ele o eleva. O Digna é a materialização tecnológica da norma ITG 2002 do Conselho Federal de Contabilidade (CFC). Ao automatizar a digitação, o sistema permite que "Contadores Sociais" voluntários ou de baixo custo atendam dezenas de cooperativas simultaneamente de forma viável, transformando conformidade legal em inclusão. O Digna estabelece o Contador Social como um parceiro estratégico. A ferramenta blinda a Responsabilidade Técnica do contador (exigência do CFC) através da **delegação temporal do acesso**. Um contador inativo perde o acesso gerencial aos dados atuais da cooperativa, mas preserva o direito inalienável de consulta (Read-Only) ao período em que gerou o SPED, garantindo segurança jurídica total para os profissionais que apoiam o movimento.

### Pilar 7: Adequação Estatal e Blindagem Digital [NOVO]
O Digna opera como um tradutor implacável entre a fome regulatória do Estado e a simplicidade da autogestão. O software absorve a fricção burocrática garantindo que o trabalhador lide com uma interface cotidiana, enquanto os motores em segundo plano resolvem a técnica:

- **Soberania nas Assembleias (Governança):** Garantir, via integração Gov.br, que as decisões políticas tenham validade imediata com Assinatura Eletrônica Qualificada (Lei nº 14.063/2020), além de assegurar a anonimização sistêmica de votos, protegendo a liberdade democrática em assembleias virtuais (IN DREI nº 79/2020).
- **Blindagem do Ato Cooperativo (Tributário):** Proteger o produtor contra a bitributação e multas de retenção de terceiros. O sistema opera a EFD-Reinf nos bastidores e realiza o expurgo matemático das transações entre cooperados no Bloco M da ECF (e-Lalur/e-Lacs), preservando a isenção prevista na Lei 5.764/71 e LC 214/2025.
- **Inclusão Sanitária (MAPA):** A agricultura familiar fica frequentemente refém do SIF para comercializar alimentos de origem animal. O Digna atua como engenheiro, possuindo um gerador automatizado do Memorial Técnico Sanitário de Estabelecimento (MTSE) para atender à Portaria MAPA nº 393/2021, pavimentando o "caderno da fábrica".
- **Interoperabilidade com Políticas Públicas:** Acesso desburocratizado a fomento estatal por meio do consumo nativo de APIs do Cadastro Nacional de Empreendimentos Econômicos Solidários (CADSOL / SINAES), fortalecendo os propósitos da Lei Paul Singer (15.068/2024).

### Pilar 8: Nenhum Dado Digitado Duas Vezes [NOVO - PDF v1.0]
**Princípio Arquitetural Central:** O que o ERP já sabe sobre a entidade é automaticamente aproveitado pelo Portal de Oportunidades, pelo Motor de Indicadores e pela Rede Digna. Isso reduz a fricção e elimina a principal barreira de uso: o excesso de formulários. O perfil de elegibilidade para crédito é construído naturalmente pelo uso cotidiano do sistema.

---

## 3. Princípios Centrais de Operação

O trabalhador da Economia Solidária **não faz contabilidade tradicional**, ele pratica a autogestão. O Contador parceiro **não digita notas**, ele audita e orienta.

**Nota sobre Laicidade do Produto:** Embora o fundamento filosófico do projeto inclua princípios da Teologia Cristã (Dignidade Humana, Mordomia, Koinonia, Justiça Restaurativa), a plataforma **não apresenta conteúdo religioso explícito ao usuário final**. A Teologia informa as decisões de design e ética, mas o produto é acessível e útil independente da crença do usuário. O canal de distribuição via igrejas e comunidades de fé é estratégico, não doutrinário.

---

## 4. Arquitetura do Ecossistema [ATUALIZADO - PDF v1.0]

O Ecossistema Digna é composto por **quatro módulos interdependentes**. O **digna ERP** é o núcleo — os demais módulos são extensões que ampliam o valor gerado a partir dos dados que o ERP já captura.

| Módulo | Função Principal | Integração-Chave | Status |
|--------|------------------|------------------|--------|
| **Módulo 1: digna ERP** | Gestão financeira, fiscal e contábil do empreendedor | Alimenta todos os demais módulos com dados do perfil | ✅ 85% Completo |
| **Módulo 2: Motor de Indicadores** | Coleta e interpreta indicadores econômicos em tempo real | Consome APIs BCB/IBGE; alimenta Portal com taxas e contexto | 📋 Backlog Fase 3 |
| **Módulo 3: Portal de Oportunidades** | Match automático entre perfil e programas de financiamento | Consome ERP + Motor; gera checklist de documentos | 📋 Backlog Fase 3 |
| **Módulo 4: Rede Digna** | Marketplace solidário entre entidades do ecossistema | Consome perfil do ERP para matching de compra/venda | 🔄 Expandir sync_engine |

### 4.1 Fluxo de Dados entre Módulos

O fluxo parte sempre do **digna ERP**, que é a fonte de verdade sobre o perfil de cada entidade.

1. **ERP captura o perfil:** No uso cotidiano (vendas, compras, caixa), o digna constrói automaticamente o perfil completo da entidade: faturamento, CNAE, município, regime tributário, situação fiscal.
2. **Motor coleta indicadores:** Diariamente, o Motor atualiza SELIC, IPCA, câmbio, taxas de crédito e expectativas de mercado via APIs públicas. Esses dados ficam disponíveis para todos os módulos.
3. **Portal executa o match:** Com o perfil do ERP e as taxas do Motor, o Portal calcula automaticamente para quais programas de financiamento a entidade é elegível e em que condições.
4. **Contador Social consolida:** O Contador Social adotado pela entidade tem visão consolidada do perfil de elegibilidade e pode submeter candidaturas aos programas identificados pelo Portal.
5. **Rede cria conexões:** Com o perfil estabelecido, a Rede identifica potenciais parceiros comerciais dentro do ecossistema e sugere conexões de compra e venda.
6. **Alertas fecham o ciclo:** Quando um edital novo é identificado ou um prazo se aproxima, o sistema notifica a entidade e o Contador Social responsável.

---

## 5. Roadmap Estratégico de Longo Prazo [ATUALIZADO]

### Fase 0: Demonstração e Validação Cultural (✅ COMPLETE)
**Foco:** O grupo informal passa a operar com rigor contábil, mas com interface amigável.

**Milestones:**
- ✅ Motor Lume Exato + PDV Pedagógico
- ✅ Testes de Usabilidade em Campo com as EES
- ✅ Sprint 16: Sistema 100% Funcional (149/149 testes)

---

### Fase 1: Integração, O Trilho da Formalização e o Contador Social (🟡 EM ANDAMENTO)
**Foco:** Oferecer os benefícios do Estado sem o peso da burocracia e aproximar a classe contábil.

**Milestones:**
- ✅ Integração Gov.br (Mock → Real OAuth2)
- ✅ Dossiê CADSOL automático
- ✅ Painel do Contador (Accountant Dashboard Multi-tenant) + SPED export
- 🔄 Criação do programa "Imposto de Renda Solidário" com Faculdades de Contabilidade/CRCs
- 🔄 Adequação Estatal: EFD-Reinf, ECF, MAPA (MTSE), CADSOL/SINAES

---

### Fase 2: Ecossistema de Crédito e Indicadores [NOVO - PDF v1.0]
**Foco:** Conectar automaticamente o empreendedor aos programas de crédito para os quais ele já é elegível, sustentado por dados econômicos atualizados em tempo real.

**Milestones:**
- 📋 **Motor de Indicadores (RF-18):** Coleta diária de SELIC, IPCA, câmbio via APIs BCB/IBGE
- 📋 **Perfil de Elegibilidade (RF-19):** Campos complementares (CadÚnico, gênero, inadimplência, finalidade do crédito)
- 📋 **Portal de Oportunidades (RF-20):** Match automático com 3 programas MVP (Pronampe, PNMPO, Niterói Empreendedora)
- 📋 **Checklist + Alertas de Documentos (RF-21-23):** Monitoramento de DOU, prazos de editais

---

### Fase 3: Intercooperação e Escala [EXPANDIDO - PDF v1.0]
**Foco:** Uma rede nacional de apoio mútuo (O 6º Princípio do Cooperativismo).

**Milestones:**
- 📋 **Perfil Público da Entidade (RF-24):** Missão, produtos, capacidades visíveis na rede
- 📋 **Mural de Necessidades (RF-25):** Entidades publicam demandas de compra
- 🔄 **Marketplace B2B fechado para EES:** Expandir sync_engine existente
- 🔄 **Score de Crédito Social baseado no trabalho:** Integrar com histórico do Motor Lume
- 🔄 **Integração com BNDES e políticas públicas via Serpro**

---

### Fase 4: Finanças Solidárias e Territoriais [REORDENADO]
**Foco:** Gerenciar riquezas além da moeda oficial (Real R$).

**Milestones:**
- 🔄 Integração tecnológica com Bancos Comunitários de Desenvolvimento (BCDs)
- 🔄 Moedas Sociais Locais
- 🔄 Estoque Substantivo (Troca de Sementes, Animais e Horas)

---

## 6. Fundamento Filosófico [NOVO - PDF v1.0]

O projeto nasce da interseção entre quatro campos do conhecimento, cada um contribuindo com uma dimensão essencial:

| Área | Contribuição ao Projeto |
|------|------------------------|
| **Tecnologia da Informação** | Arquitetura de sistemas, APIs, automação, interfaces acessíveis e inteligência de dados |
| **Ciências Contábeis** | Gestão financeira, obrigações fiscais, escrituração, módulo de contador social e conformidade normativa |
| **Ciências Econômicas** | Indicadores macroeconômicos, análise de crédito, mapeamento de financiamentos e interpretação de política econômica |
| **Teologia Cristã** | Base ética e filosófica do design: princípios de dignidade humana, mordomia, koinonia e justiça restaurativa |

### 6.1 Princípios Teológicos Aplicados ao Design

| Princípio | Aplicação Prática no Digna |
|-----------|---------------------------|
| **Dignidade Humana (Imago Dei)** | Interface que nunca humilha. Linguagem sem jargão excludente. Design que pressupõe competência, não ignorância. Ausência de mensagens punitivas por inadimplência |
| **Mordomia (Stewardship)** | O sistema incentiva gestão responsável dos recursos. Relatórios que mostram saúde financeira de forma construtiva. Ferramentas de planejamento, não de maximização de lucro a qualquer custo |
| **Koinonia (Comunhão)** | A funcionalidade de Rede Digna não é apenas negócio — é comunidade econômica que se sustenta mutuamente. Igrejas e líderes comunitários como agentes de distribuição e capacitação |
| **Justiça Restaurativa (Shalom)** | Acesso à informação normativa como ato de equidade. O que o contador caro faz pelo rico, o sistema faz pelo pobre. Métricas de bem-estar que vão além do lucro |

---

## 7. Problema Central e Solução [NOVO - PDF v1.0]

### 7.1 Problema Central

O empreendedor de baixa renda no Brasil enfrenta quatro barreiras simultâneas:
1. **Desconhecimento das obrigações normativas** do seu negócio
2. **Inacessibilidade da linguagem burocrática** dos programas de crédito
3. **Ausência de ferramentas de gestão** adequadas ao seu nível de escolaridade e realidade financeira
4. **Isolamento** — não sabe que existem redes de apoio, programas de financiamento e parceiros potenciais à sua volta

### 7.2 Solução Proposta

Uma plataforma que:
- **Esconde a complexidade normativa** por trás de interfaces simples
- **Conecta automaticamente** o empreendedor aos programas de crédito para os quais ele já é elegível
- **Cria uma rede de colaboração** entre entidades
- **Tudo sustentado por dados econômicos atualizados em tempo real**

---

## 8. Constituição Técnica [MANTIDO]

A arquitetura técnica permanece inalterada, garantindo consistência e qualidade:

| Princípio | Implementação |
|-----------|--------------|
| **Soberania do Dado** | SQLite isolado por entidade: `data/entities/{entity_id}.db` |
| **Anti-Float** | `int64` obrigatório para centavos/minutos (proibido `float`) |
| **Cache-Proof Templates** | `*_simple.html` + `ParseFiles()` no handler (não globais) |
| **Clean Architecture + DDD** | Domain → Service → Repository → Handler |
| **Offline-First** | PWA com sync delta assíncrono |

**Stack:** Go 1.22+ • SQLite3 • HTMX + Tailwind • Multi-module Workspace (`go.work`)

---

## 9. Métricas de Sucesso do Ecossistema [NOVO]

| Métrica | Descrição | Alvo |
|---------|-----------|------|
| **Redução de Formulários** | Dados reaproveitados do ERP vs. formulários manuais | 80% redução |
| **Match de Crédito** | Entidades que descobrem elegibilidade via Portal | 100+ no primeiro ano |
| **Adoção por Contadores** | Contadores Sociais ativos na plataforma | 50+ no primeiro ano |
| **Transações na Rede** | Conexões B2B realizadas via Rede Digna | 500+ no primeiro ano |
| **Conformidade Automatizada** | Entidades com EFD-Reinf/ECF gerada automaticamente | 90% das formalizadas |

---

## 10. Notas de Governança [MANTIDO]

- **Fundação Providentia:** Guardiã da neutralidade e missão social
- **Licença:** Apache 2.0 (código aberto, uso comercial permitido)
- **Marca "Digna":** Pertence exclusivamente à Fundação Providentia
- **Modelo de Governança:** Inspirado na Apache Foundation, com mérito técnico e domínio do negócio

---

**Status:** ✅ ATUALIZADO COM VISÃO DE ECOSSISTEMA (PDF v1.0)  
**Próxima Ação:** Atualizar `06_roadmap/02_roadmap.md` e `06_roadmap/03_backlog.md` com Fases 2-4 expandidas  
**Versão Anterior:** 1.5 (2026-03-13)  
**Versão Atual:** 2.0 (2026-03-13)
