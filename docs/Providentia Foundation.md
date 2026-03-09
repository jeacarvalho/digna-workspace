# **Projeto Providentia: Dignidade e Soberania Financeira**

**Data de Atualização:** 05/03/2026 | **Versão:** 0.3 (Social-First)

## **Capítulo 1: Fundamentos, Visão e Propósito**

### **1.1. O Cenário Atual e a Necessidade de Disrupção**

A economia solidária brasileira enfrenta o desafio da "invisibilidade burocrática". Pequenos empreendimentos produzem riqueza real, mas soçobram ante a complexidade contábil e a dificuldade de acesso ao crédito. O projeto Providentia nasce para romper esse ciclo, transformando a gestão em um ato de cidadania automática.

### **1.2. A Filosofia "Digna": Contabilidade Invisível**

O produto **Digna** (app/web) é a interface de libertação do empreendedor. A filosofia central é que o trabalhador não deve "fazer contabilidade", mas sim "operar seu negócio". Através de um **PDV (Ponto de Venda)** intuitivo, cada venda de mel, cada hora de costura e cada compra de insumo é capturada e traduzida instantaneamente em linguagem contábil pelo motor **Lume**.

### **1.3. O Serpro como Indutor de Soberania Nacional**

Como braço tecnológico do Estado, o Serpro atua como o grande hospedeiro desta infraestrutura, garantindo que o dado do cooperado não seja mercadoria, mas sim um ativo de soberania para a formulação de políticas públicas de crédito e fomento.

### **1.4. A Fundação Providentia: O Modelo de Governança**

A fundação serve como a guardiã da neutralidade. Ela assegura que o software permaneça um bem público, imune a capturas de mercado, e focado exclusivamente na dignidade financeira dos empreendimentos de economia solidária (EES).

### **1.5. Objetivos Estratégicos**

* **Curto Prazo (v0 \- A Vitrine):** Lançamento do MVP funcional com PDV Offline-First e registro de horas de trabalho (ITG 2002).  
* **Médio Prazo (Escala e Formalização):** Automação do dossiê **CADSOL** para reconhecimento de autogestão e acesso a compras públicas (PNAE/PAA).  
* **Longo Prazo (Ecossistema de Crédito):** Criação de uma rede de intercooperação e crédito baseada no histórico de impacto social e financeiro.

---

## **Capítulo 2: Modelo de Governança e Regras de Membresia**

### **2.1. O "Apache Way" à Brasileira**

A Fundação Providentia adota o modelo de governança baseada em mérito e neutralidade. O objetivo é garantir que o projeto não seja "dono" de ninguém, mas um bem público digital gerido por quem efetivamente contribui para o seu sucesso.

**Pilares da Governança:**

* **Independência de Fornecedor:** O software deve ser funcional independentemente de qualquer empresa privada específica.  
* **Transparência nas Decisões:** Todas as deliberações técnicas e estratégicas são registradas e públicas.  
* **Comunidade sobre o Código:** A saúde da rede de desenvolvedores e usuários é mais importante que qualquer funcionalidade isolada.

### **2.2. Categorias de Membresia e Participação**

A fundação estrutura-se em diferentes níveis para acomodar os diversos atores do ecossistema:

| Categoria | Perfil | Papel e Contribuição |
| :---- | :---- | :---- |
| **Membros Fundadores** | Empreendedores de sucesso (Pessoas Físicas) | Aporte de capital semente (Endowment), experiência e visão estratégica de escala. |
| **Parceiro Institucional** | Serpro (Empresa Pública) | Provedor de infraestrutura de nuvem, segurança e integração com o Estado. |
| **Membros Corporativos** | Empresas de Tecnologia/SaaS | Doação de horas de engenharia ou recursos financeiros (Modelo Platinum/Gold). |
| **Membros Acadêmicos** | Universidades e Centros de Pesquisa | Validação metodológica e formação de redes de tutoria. |
| **Contribuidores** | Desenvolvedores e Contadores | Ganham poder de voto através de contribuições técnicas de alta qualidade. |
| **CFC** | Contadores Sociais | Normatização, padronização, auditorias e apadrinhamento de empreendimentos solidários |

### **2.3. Estrutura de Tomada de Decisão**

#### **2.3.1. Conselho Curador (Board of Trustees)**

É o órgão máximo de direção. Focado na sustentabilidade financeira e no alinhamento com a missão social.

* **Voto por Assento:** Cada segmento (Fundadores, Serpro, Usuários) possui assentos definidos no estatuto para evitar a captura por um único grupo.  
* **Mandatos:** Rotativos para garantir renovação e evitar o "personalismo".

#### **2.3.2. Comitê de Gestão de Projeto (PMC \- Project Management Committee)**

Responsável pelas decisões técnicas e pelo roadmap do software.

* **Mérito Técnico:** Para entrar no PMC, o indivíduo deve ser um Committer ativo com contribuições validadas.  
* **Neutralidade:** O peso de voto é individual, independente da afiliação corporativa do membro.

### **2.4. Propriedade Intelectual e Licenciamento**

* **Licença Apache 2.0:** Todo o código é distribuído sob licença permissiva, garantindo que o núcleo permaneça aberto.  
* **Gestão de Marcas:** A marca "Providentia" pertence à Fundação, protegendo o nome contra usos indevidos.

### **2.5. O Ciclo do "Giving Back" (A Contrapartida Social)**

A fundação produz um **Relatório Anual de Impacto Social**, medindo a riqueza transacionada e a valoração do trabalho voluntário (conforme a norma **ITG 2002**), gerando transparência sobre o fortalecimento da rede.

### **2.6. Resolução de Conflitos e Soberania**

Em caso de divergência, a solução deve sempre maximizar a Soberania do Usuário Final. A arquitetura **Local-First** é uma regra de governança: o dado pertence ao usuário, garantindo a liberdade de saída (*Exit Power*) caso a missão seja comprometida.

---

## **Capítulo 3: Arquitetura e Engenharia de Software (Lume Engine)**

### **3.1. Paradigma: Micro-databases e Offline-First**

Cada entidade possui sua própria instância física de **SQLite** isolada. Para o "Brasil Profundo", a aplicação é desenhada como **Offline-First**: as transações são capturadas localmente e sincronizadas assincronamente (Push) quando houver conectividade.

### **3.2. O Motor Lume (Regras de Ouro)**

1. **Integridade Monetária (int64):** Uso estrito de inteiros para representar centavos. Floats são proibidos.  
2. **Partidas Dobradas Automáticas:** Validação sistemática de soma zero (Débitos \+ Créditos \= 0).  
3. **Valoração do Trabalho:** Suporte nativo ao registro de horas de trabalho militante/voluntário como ativo social.

### **3.3. Sincronização e o Grande Agregador**

O sistema utiliza um motor de sincronização por lotes que alimenta o Agregador Central da Fundação para fins de auditoria social e geração de dossiês de crédito.

### **3.4. Ciclo de Vida e Formalização (Legal Facade)**

Utilização de interfaces (Mocks) para gerir a transição do estado informal (**DREAM**) para o formalizado (**FORMALIZED**), gerando automaticamente atas de fundação e estatutos compatíveis com o **CADSOL**.

### **3.5. Stack Tecnológica de Referência**

* **Backend:** Go (1.22+).  
* **Estrutura:** Multi-module Workspace (go.work).  
* **Nomenclatura:** Kebab-case para diretórios; snake\_case para arquivos.

---

## **Capítulo 4: Roadmap e Estratégia da v0**

### **4.1. Objetivo da Versão 0 (A Vitrine)**

Demonstrar a "Contabilidade Invisível" através de um PDV funcional que registra vendas e horas de trabalho, gerando automaticamente um balanço de impacto social.

### **4.2. Detalhamento das Funcionalidades**

1. **Gestão de Contexto (Lifecycle Manager):** Orquestração automática de bancos SQLite por entidade.  
2. **PDV Operacional:** Registro offline de transações comerciais.  
3. **Dossiê CADSOL:** Geração automática de atas de assembleia e logs de decisão para provar a autogestão.

---

## **Capítulo 5: Estratégia de Modularização para Desenvolvimento Ágil**

Para mitigar a entropia de contexto e garantir a precisão da construção assistida por IA, o sistema é segmentado em cinco módulos atômicos e independentes:

1. **Módulo 0 (pdv\_ui):** Interface de operação (Venda/Compra/Horas) com suporte a cache local.  
2. **Módulo 1 (lifecycle):** Orquestrador de arquivos SQLite e gerenciador de estados de sincronização.  
3. **Módulo 2 (core\_lume):** Motor contábil de partidas dobradas monetárias e sociais.  
4. **Módulo 3 (legal\_facade):** Simulador de formalização e gerador de documentos (Atas e Estatutos).  
5. **Módulo 4 (reporting):** Painel de Dignidade, Balanço Social e Dossiê de Crédito.

Esta modularização permite que o desenvolvimento ocorra em ciclos fechados, onde cada componente é testado independentemente, garantindo que a complexidade do sistema cresça de forma controlada e sem alucinações do agente de codificação.

