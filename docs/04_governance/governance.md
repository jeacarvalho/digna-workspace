***

```markdown
---
title: Governança do Projeto Digna
status: implemented
version: 1.1
last_updated: 2026-03-08
---

# Governança - Projeto Digna

**Projeto:** Digna - Infraestrutura Contábil para Economia Solidária
**Mantenedor:** Fundação Providentia
**Modelo:** Apache Foundation

---

## 1. Fundação Providentia

A **Fundação Providentia** é a entidade responsável por garantir a continuidade, neutralidade e missão social do projeto Digna.

### Missão

> Promover a autogestão, soberania e transformação digital dos Empreendimentos de Economia Solidária (EES) no Brasil, através de tecnologia livre e acessível, atuando como ponte tecnológica inclusiva para a conformidade legal e contábil.

### Princípios Core

| Princípio | Descrição |
|-----------|-----------|
| **Soberania** | O dado pertence à entidade, nunca à plataforma |
| **Transparência** | Código aberto, dados abertos (agregados) e algoritmos auditáveis visualmente |
| **Colaboração** | Intercooperação entre redes, Estado, Academia e Profissionais |
| **Transformação** | Compromisso com a mudança social real |
| **Aliança Contábil** | Valorização do Contador Social como consultor parceiro da autogestão |

### Responsabilidades

- **Neutralidade** - Garantir que o projeto sirva aos interesses da economia solidária.
- **Missão Social** - Manter o foco em impacto social, não em lucro.
- **Conformidade Normativa** - Assegurar o alinhamento tecnológico contínuo com as diretrizes do Conselho Federal de Contabilidade (CFC - ITG 2002).
- **Continuidade** - Assegurar a perpetuidade do projeto independente de contribuições individuais.
- **Infraestrutura** - Coordenar com Serpro a infraestrutura de nuvem soberana.

### Modelo de Governança

Inspirado na Apache Foundation, o modelo prioriza:

- Mérito técnico e domínio do negócio (contábil/social) para tomada de decisões.
- Transparência nos processos.
- Comunidade aberta e inclusiva.

---

## 2. Membresia e Papéis

### Categorias

1. **Maintainers (Fundadores)**
   - Visão estratégica e arquitetura
   - Aprovação final de RFCs estruturais
   - Governança da Fundação

2. **PMC (Project Management Committee)**
   - Gestão técnica diária
   - Aprovação de Pull Requests críticos (Core Lume e Accountant Dashboard)
   - Planejamento de Sprints

3. **Committers**
   - Desenvolvedores com acesso de escrita
   - Manutenção de módulos específicos
   - Revisão de código par-a-par

4. **Contribuidores (Desenvolvedores e Contadores Sociais)**
   - Profissionais de tecnologia ou de contabilidade que enviam código, reportam *bugs*, propõem melhorias na documentação ou validam regras de negócio fiscais (SPED/ITG 2002).

### Critérios de Entrada no PMC/Committers

- Contribuição técnica ou de domínio (negócios/contabilidade) comprovada.
- Participação ativa em revisões.
- Compromisso com os princípios do projeto e as "Regras de Ouro" (Anti-float, Soberania do Dado).

---

## 3. Regras de Contribuição

### Processo de Contribuição

```text
1. Fork do repositório
2. Criar branch feature
3. Desenvolver com testes
4. Abrir Pull Request
5. Revisão por pares
6. Aprovação do PMC
7. Merge
```

### Padrões Obrigatórios

- **Código:** Gofmt + golint.
- **Rigor Financeiro:** Proibição absoluta de variáveis `float` no Motor Lume (uso estrito de `int64`).
- **Testes:** Cobertura mínima 80% (com foco absoluto na validação de soma zero nas transações).
- **Commits:** Conventional Commits.
- **Documentação:** Atualizada junto com código.

### Review Checklist

- [ ] Código segue Go conventions
- [ ] Testes passando (com TDD para regras de negócio)
- [ ] Rigor Monetário e Contábil (uso exclusivo de `int64`, partidas dobradas e ITG 2002 mantidos)
- [ ] Documentação atualizada
- [ ] Sem dados sensíveis expostos
- [ ] Licença Apache 2.0 declarada

---

## 4. Licenciamento

### Código Fonte

**Licença:** Apache 2.0

Esta licença permite:
- Uso comercial
- Modificação
- Distribuição
- Uso privado
- Sublicenciamento

### Marca "Digna"

A marca "Digna" pertence exclusivamente à Fundação Providentia. O uso da marca requer autorização prévia por escrito, inclusive para "Selos de Certificação" concedidos a Contadores Parceiros e Incubadoras.

### Dependências
Todas as dependências do projeto devem ser compatíveis com Apache 2.0 ou licenças permissivas similares.

---

## 5. Tomada de Decisões

### Tipos de Decisão

| Tipo | threshold | Quem |
|------|------------|------|
| Minor | 1 aprovador | Maintainer |
| Major | 2 aprovadores | PMC |
| Strategic | Consensus | Fundação / Conselho Curador |

### RFC Process

Para mudanças maiores:

1. Criar RFC em `/rfcs/`
2. Período de discussão (2 semanas)
3. Votação do PMC
4. Decisão final pela Fundação
*Nota: RFCs que impactem o Core Lume ou o formato de exportação de lotes fiscais (SPED) exigem obrigatoriamente a revisão técnica de membros ligados à classe contábil/CFC.*

---

## 6. Transparência

### Repositórios Públicos

- Código fonte
- Documentação
- RFCs
- Decisões de governance

### Comunicação

- Issues abertas para bugs e features
- Discussions para dúvidas e debates sociotécnicos
- Meetings gravadas (quando aplicável)

---

## 7. Conflito de Interesses

Contribuidores devem declarar conflitos de interesse potenciais. Decisões afetadas por COI devem ser tratadas por membros não envolvidos. Profissionais (como contadores ou auditores) que atuem como validadores no projeto assumem o compromisso ético de não utilizar a governança da plataforma como meio de captação predatória de clientes, mantendo o software como bem público independente.

---

## Referências

- [Apache Foundation Way](https://www.apache.org/foundation/how-it-works.html)
- [Contributor Covenant](https://www.contributor-covenant.org/)
- [ITG 2002 (R1) - Conselho Federal de Contabilidade](https://cfc.org.br)
```

***