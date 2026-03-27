title: Governança - Projeto Digna
status: implemented
version: 2.1
last_updated: 2026-03-27
---

# Governança - Projeto Digna

> **Nota:** Este documento reflete a governança integrada do Ecossistema Digna (PDF v1.0), preservando toda a infraestrutura validada nas Sprints 1-16, incorporando os novos módulos como Fases 3 e 4, e incluindo o Sistema de Ajuda Educativa (RF-30) decidido na sessão de 27/03/2026.

**Projeto:** Digna - Infraestrutura Contábil para Economia Solidária  
**Mantenedor:** Fundação Providentia  
**Modelo:** Apache Foundation  
**Licença:** Apache 2.0

---

## 1. Fundação Providentia

A **Fundação Providentia** é a entidade responsável por garantir a continuidade, neutralidade e missão social do projeto Digna.

### 1.1 Missão [ATUALIZADO - PDF v1.0]

Promover a autogestão, soberania e transformação digital dos Empreendimentos de Economia Solidária (EES) no Brasil, através de tecnologia livre e acessível, atuando como **ponte tecnológica inclusiva** para a conformidade legal e contábil, e como **ecossistema integrado** de 4 módulos (ERP, Motor de Indicadores, Portal de Oportunidades e Rede Digna).

### 1.2 Princípios Core [ATUALIZADO]

| Princípio | Descrição |
|-----------|-----------|
| **Soberania** | O dado pertence à entidade, nunca à plataforma (Exit Power) |
| **Transparência** | Código aberto, dados abertos (agregados) e algoritmos auditáveis visualmente |
| **Colaboração** | Intercooperação entre redes, Estado, Academia e Profissionais |
| **Transformação** | Compromisso com a mudança social real |
| **Aliança Contábil** | Valorização do Contador Social como consultor parceiro da autogestão |
| **Pedagogia Social [NOVO - RF-30]** | O sistema ensina enquanto é operado; nenhum usuário deve se sentir humilhado por não entender um termo |
| **Laicidade do Produto [NOVO - PDF v1.0]** | A Teologia informa o design internamente, mas o produto é acessível independente da crença do usuário |

### 1.3 Responsabilidades [ATUALIZADO]

| Responsabilidade | Descrição |
|-----------------|-----------|
| **Neutralidade** | Garantir que o projeto sirva aos interesses da economia solidária |
| **Missão Social** | Manter o foco em impacto social, não em lucro |
| **Conformidade Normativa e Estatal [ATUALIZADO]** | Assegurar o alinhamento tecnológico contínuo com as diretrizes do Conselho Federal de Contabilidade (CFC - ITG 2002) e a blindagem ativa da conformidade perante os órgãos estatais (MAPA, RFB, DREI e MTE) |
| **Continuidade** | Assegurar a perpetuidade do projeto independente de contribuições individuais |
| **Infraestrutura** | Coordenar com Serpro a infraestrutura de nuvem soberana |
| **Ecossistema de 4 Módulos [NOVO - PDF v1.0]** | Governar a evolução integrada dos 4 módulos (ERP, Motor, Portal, Rede) preservando o princípio "Nenhum dado digitado duas vezes" |
| **Sistema de Ajuda [NOVO - RF-30]** | Garantir que todo conteúdo educativo seja revisado por ITCPs/comunidade antes de publicação |

---

## 2. Modelo de Governança

Inspirado na **Apache Foundation**, o modelo prioriza:

- **Mérito técnico e domínio do negócio** (contábil/social) para tomada de decisões
- **Transparência nos processos**
- **Comunidade aberta e inclusiva**

### 2.1 Categorias de Membresia e Participação [ATUALIZADO - PDF v1.0]

| Categoria | Perfil | Papel e Contribuição |
|-----------|--------|---------------------|
| **Membros Fundadores** | Empreendedores de sucesso (Pessoas Físicas) | Aporte de capital semente (Endowment), experiência e visão estratégica de escala |
| **Parceiro Institucional** | Serpro (Empresa Pública) | Provedor de infraestrutura de nuvem, segurança e integração com o Estado |
| **Membros Corporativos** | Empresas de Tecnologia/SaaS | Doação de horas de engenharia ou recursos financeiros (Modelo Platinum/Gold) |
| **Membros Acadêmicos** | Universidades e Centros de Pesquisa | Validação metodológica e formação de redes de tutoria |
| **Contribuidores (Desenvolvedores e Contadores Sociais)** | Profissionais de tecnologia ou de contabilidade | Ganham poder de voto através de contribuições técnicas de alta qualidade |
| **CFC / CRCs [NOVO]** | Contadores Sociais | Normatização, padronização, auditorias e apadrinhamento de empreendimentos solidários |
| **ITCPs e Incubadoras [NOVO - RF-30]** | Incubadoras Tecnológicas de Cooperativas Populares | Validação de conteúdo educativo, teste de usabilidade com usuários reais |
| **Comunidade de Fé [NOVO - PDF v1.0]** | Igrejas e líderes comunitários | Agentes de distribuição e capacitação (canal estratégico, não doutrinário) |

### 2.2 Estrutura de Tomada de Decisão [ATUALIZADO]

#### 2.2.1 Conselho Curador (Board of Trustees)
É o órgão máximo de direção. Focado na sustentabilidade financeira e no alinhamento com a missão social.

- **Voto por Assento:** Cada segmento (Fundadores, Serpro, Usuários) possui assentos definidos no estatuto para evitar a captura por um único grupo
- **Mandatos:** Rotativos para garantir renovação e evitar o "personalismo"
- **Novo:** Representante da comunidade de ITCPs para validação pedagógica (RF-30)

#### 2.2.2 Comitê de Gestão de Projeto (PMC - Project Management Committee)
Responsável pelas decisões técnicas e pelo roadmap do software.

- **Mérito Técnico:** Para entrar no PMC, o indivíduo deve ser um Committer ativo com contribuições validadas
- **Neutralidade:** O peso de voto é individual, independente da afiliação corporativa do membro
- **Novo:** Validação de conteúdo educativo requer aprovação do PMC + ITCPs (RF-30)

#### 2.2.3 Comitê de Conformidade Estatal [NOVO - PDF v1.0]
Responsável por validar adequação às normas estatais (MAPA, RFB, DREI, MTE).

- **Composição:** Contadores sociais, auditores, especialistas em conformidade
- **Responsabilidade:** Validar módulos `tax_compliance`, `sanitary_compliance`, integrações Gov.br
- **Aprovação:** Obrigatória para mudanças em exportação fiscal, assinaturas eletrônicas, MTSE

---

## 3. Regras de Contribuição

### 3.1 Processo de Contribuição

#### Padrões Obrigatórios [ATUALIZADO]

| Padrão | Descrição |
|--------|-----------|
| **Código** | Gofmt + golint |
| **Rigor Financeiro** | Proibição absoluta de variáveis float no Motor Lume (uso estrito de int64) |
| **Blindagem Tributária (Regra de Ouro) [NOVO]** | A IA e os desenvolvedores devem sempre proteger o Ato Cooperativo (Lei nº 5.764/71 e LC nº 214/2025) na base de cálculo tributária (e-Lalur/e-Lacs no Bloco M da ECF), evitando a bitributação da EES |
| **Conformidade Regulatória [NOVO]** | Observância obrigatória aos leiautes oficiais do SPED, manuais da EFD-Reinf, IN DREI nº 79/2020 e Portaria MAPA nº 393/2021 em novas automações |
| **Pedagogia (RF-30) [NOVO]** | Todo campo técnico deve ter explicação acessível via Sistema de Ajuda; linguagem para 5ª série |
| **Testes** | Cobertura mínima 80% (com foco absoluto na validação de soma zero nas transações) |
| **Commits** | Conventional Commits |
| **Documentação** | Atualizada junto com código |
| **Validação E2E [NOVO]** | Validação E2E obrigatória antes de concluir tarefas (`validate_e2e.sh --basic --headless`) |

#### Review Checklist [ATUALIZADO]

- [ ] Código segue Go conventions
- [ ] Testes passando (com TDD para regras de negócio)
- [ ] Rigor Monetário e Contábil (uso exclusivo de int64, partidas dobradas e ITG 2002 mantidos)
- [ ] Documentação atualizada
- [ ] Sem dados sensíveis expostos
- [ ] Licença Apache 2.0 declarada
- [ ] **NOVO:** Conteúdo educativo revisado por ITCPs (se aplicável - RF-30)
- [ ] **NOVO:** Conformidade estatal validada (se aplicável - módulos fiscais/sanitários)
- [ ] **NOVO:** Validação E2E passou (`validate_e2e.sh`)

---

## 4. Licenciamento

### 4.1 Código Fonte

**Licença:** Apache 2.0

Esta licença permite:
- Uso comercial
- Modificação
- Distribuição
- Uso privado
- Sublicenciamento

### 4.2 Marca "Digna" [ATUALIZADO - PDF v1.0]

A marca "Digna" pertence exclusivamente à **Fundação Providentia**. O uso da marca requer autorização prévia por escrito, inclusive para:
- "Selos de Certificação" concedidos a Contadores Parceiros e Incubadoras
- Materiais de capacitação de líderes comunitários
- Distribuição via canais de fé (igrejas, comunidades)

**Nota sobre Laicidade:** A marca não pode ser usada de forma a associar o produto exclusivamente a uma denominação religiosa específica, preservando a acessibilidade universal.

### 4.3 Conteúdo Educativo [NOVO - RF-30]

Conteúdo do Sistema de Ajuda Educativa:
- **Licença:** CC BY-SA 4.0 (Creative Commons)
- **Revisão:** Obrigatória por ITCPs antes de publicação
- **Atualização:** Versionada com hash de integridade
- **Acesso:** Público e gratuito para todos os usuários Digna

### 4.4 Dependências

Todas as dependências do projeto devem ser compatíveis com Apache 2.0 ou licenças permissivas similares.

---

## 5. Tomada de Decisões

### 5.1 Tipos de Decisão [ATUALIZADO]

| Tipo | Threshold | Quem |
|------|-----------|------|
| **Minor** | 1 aprovador | Maintainer |
| **Major** | 2 aprovadores | PMC |
| **Strategic** | Consensus | Fundação / Conselho Curador |
| **Conformidade Estatal [NOVO]** | Validação obrigatória | Comitê de Conformidade + PMC |
| **Conteúdo Educativo [NOVO - RF-30]** | Validação por ITCPs | Comitê Pedagógico + PMC |

### 5.2 RFC Process [ATUALIZADO]

Para mudanças maiores:

1. Criar RFC em `/rfcs/`
2. Período de discussão (2 semanas)
3. Votação do PMC
4. Decisão final pela Fundação

**Nota:** RFCs que impactem o Core Lume, o formato de exportação de lotes fiscais (SPED/EFD-Reinf), governança digital (Assembleias Gov.br), conformidade sanitária (MAPA), ou conteúdo educativo (RF-30), exigem obrigatoriamente a revisão técnica de membros ligados à classe contábil/CFC, auditores de conformidade, e/ou ITCPs conforme aplicável.

---

## 6. Transparência

### 6.1 Repositórios Públicos

- Código fonte
- Documentação
- RFCs
- Decisões de governance
- **NOVO:** Conteúdo educativo do Sistema de Ajuda (revisado)

### 6.2 Comunicação

- Issues abertas para bugs e features
- Discussions para dúvidas e debates sociotécnicos
- Meetings gravadas (quando aplicável)
- **NOVO:** Relatórios de validação pedagógica (RF-30)
- **NOVO:** Relatórios de conformidade estatal (módulos fiscais/sanitários)

---

## 7. Conflito de Interesses [ATUALIZADO - PDF v1.0]

Contribuidores devem declarar conflitos de interesse potenciais. Decisões afetadas por COI devem ser tratadas por membros não envolvidos.

**Profissionais (como contadores ou auditores)** que atuem como validadores no projeto assumem o compromisso ético de:
- Não utilizar a governança da plataforma como meio de captação predatória de clientes
- Manter o software como bem público independente
- **NOVO:** Validar conteúdo de conformidade com isenção técnica (Comitê de Conformidade)
- **NOVO:** Validar conteúdo educativo com foco na inclusão, não em doutrinação (Comitê Pedagógico - RF-30)

**Canais de Distribuição via Comunidade de Fé [NOVO - PDF v1.0]:**
- São estratégicos, não doutrinários
- O produto permanece laico na interface
- A Teologia informa decisões de design internamente, mas não é exposta ao usuário final
- Líderes comunitários atuam como capacitadores, não como gatekeepers doutrinários

---

## 8. Comitês Especializados [NOVO - PDF v1.0 + RF-30]

### 8.1 Comitê de Conformidade Estatal

**Responsabilidade:** Validar módulos de adequação estatal (EFD-Reinf, ECF, MAPA, Gov.br, CADSOL).

**Composição:**
- Contadores sociais certificados
- Auditores de conformidade fiscal
- Especialistas em legislação cooperativa (Lei 5.764/71, LC 214/2025)

**Aprovações Obrigatórias:**
- Exportação SPED/EFD-Reinf
- Assinaturas eletrônicas Gov.br
- Memorial Técnico Sanitário (MTSE)
- Integração CADSOL/SINAES

### 8.2 Comitê Pedagógico [NOVO - RF-30]

**Responsabilidade:** Validar conteúdo do Sistema de Ajuda Educativa.

**Composição:**
- Representantes de ITCPs
- Educadores de economia solidária
- Usuários finais (testadores de usabilidade)

**Aprovações Obrigatórias:**
- Novos tópicos de ajuda
- Atualizações de conteúdo educativo
- Mudanças de linguagem em campos técnicos

**Critério de Validação:**
- Usuário com 5ª série consegue entender sem ajuda externa
- Zero jargões técnicos sem explicação
- Sempre inclui "próximo passo" acionável

### 8.3 Comitê de Ecossistema [NOVO - PDF v1.0]

**Responsabilidade:** Governar integração entre os 4 módulos (ERP, Motor, Portal, Rede).

**Composição:**
- Arquitetos de software
- Representantes de cada módulo
- Especialistas em integração de dados

**Aprovações Obrigatórias:**
- Mudanças no princípio "Nenhum dado digitado duas vezes"
- Novas integrações entre módulos
- Mudanças no schema de dados compartilhados

---

## 9. Referências [ATUALIZADO - PDF v1.0]

| Referência | Aplicação no Digna |
|------------|-------------------|
| Apache Foundation Way | Modelo de governança |
| Contributor Covenant | Código de conduta |
| ITG 2002 (R1) - Conselho Federal de Contabilidade | Norma de contabilidade para EES |
| Lei nº 5.764/71 (Lei Geral das Cooperativas) e LC 214/2025 [NOVO] | Atos Cooperativos e tributação |
| Lei nº 15.068/2024 (Lei Paul Singer) [NOVO] | Economia Solidária |
| Lei nº 14.063/2020 (Assinaturas Eletrônicas Gov.br) [NOVO] | Governança digital |
| IN DREI nº 79/2020 (Reuniões e Assembleias Digitais) [NOVO] | Assembleias virtuais |
| Decreto nº 12.784/2025 (SINAES/CADSOL) [NOVO] | Cadastro nacional de EES |
| Portaria MAPA nº 393/2021 [NOVO] | Memorial Técnico Sanitário (MTSE) |
| Lei 12.249/2010 (Reserva de Mercado Contábil) [NOVO] | Base legal para Contador Social |

---

## 10. Próximos Passos de Governança [NOVO - PDF v1.0 + RF-30]

### 10.1 Imediatos (Próximo Trimestre)

- [ ] Formalizar Comitê de Conformidade Estatal
- [ ] Formalizar Comitê Pedagógico (RF-30)
- [ ] Definir processo de validação de conteúdo educativo
- [ ] Documentar critérios de certificação para Contadores Sociais

### 10.2 Médio Prazo (6 meses)

- [ ] Estabelecer parcerias com 5+ ITCPs para validação pedagógica
- [ ] Criar programa de certificação CFC para Contadores Sociais
- [ ] Validar primeiro conteúdo educativo com usuários reais
- [ ] Estabelecer canal de distribuição via 3+ comunidades de fé

### 10.3 Longo Prazo (1 ano+)

- [ ] Expandir Comitê de Ecossistema para governança de 4 módulos
- [ ] Estabelecer métricas de impacto social (além de métricas técnicas)
- [ ] Criar relatório anual de impacto social (Giving Back)
- [ ] Avaliar habilitação como IMPO no PNMPO

---

**Status:** ✅ ATUALIZADO COM GOVERNANÇA DO ECOSSISTEMA (PDF v1.0) + RF-30 (Sessão 27/03/2026)  
**Próxima Ação:** Atualizar `05_ai/01_constitution.md` com nova constituição de IA para agentes  
**Versão Anterior:** 1.2 (2026-03-13)  
**Versão Atual:** 2.1 (2026-03-27)
