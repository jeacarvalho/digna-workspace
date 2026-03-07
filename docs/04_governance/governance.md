---
title: Governança do Projeto Digna
status: implemented
version: 1.0
last_updated: 2026-03-07
---

# Governança - Projeto Digna

**Projeto:** Digna - Infraestrutura Contábil para Economia Solidária  
**Mantenedor:** Fundação Providentia  
**Modelo:** Apache Foundation  

---

## 1. Fundação Providentia

A **Fundação Providentia** é a entidade responsável por garantir a continuidade, neutralidade e missão social do projeto Digna.

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
| **Transparência** | clareza nos processos, dados e decisões |
| **Transformação** | Compromisso com a mudança social real |

### Responsabilidades

- **Neutralidade** - Garantir que o projeto sirva aos interesses da economia solidária
- **Missão Social** - Manter o foco em impacto social, não em lucro
- **Continuidade** - Assegurar a perpetuidade do projeto independente de contribuições individuais
- **Infraestrutura** - Coordenar com Serpro a infraestrutura de nuvem soberana

### Modelo de Governança

Inspirado na Apache Foundation, o modelo prioriza:
- Mérito técnico para tomada de decisões
- Transparência nos processos
- Comunidade aberta e inclusiva

---

## 2. Project Management Committee (PMC)

### Composição

O PMC é composto por membros com histórico de contribuição técnica significativa ao projeto.

### Responsabilidades

| Função | Descrição |
|--------|-----------|
| Decisões Técnicas | Aprovar arquitetura e mudanças fundamentais |
| Roadmap | Definir prioridades e fases do projeto |
| Revisão de Código | Assegurar qualidade e conformidade com padrões |
| Community Building | Recrutar e mentorar novos contribuidores |

### Critérios de Entrada

- Contribuição técnica comprovada
- Participação ativa em revisões
- Compromisso com os princípios do projeto

---

## 3. Regras de Contribuição

### Processo de Contribuição

```
1. Fork do repositório
2. Criar branchfeature
3. Desenvolver com testes
4. Abrir Pull Request
5. Revisão por pares
6. Aprovação do PMC
7. Merge
```

### Padrões Obrigatórios

- **Código:** Gofmt + golint
- **Testes:** Cobertura mínima 80%
- **Commits:** Conventional Commits
- **Documentação:** Atualizada junto com código

### Review Checklist

- [ ] Código segue Go conventions
- [ ] Testes passando
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

A marca "Digna" pertence exclusivamente à Fundação Providentia. O uso da marca requer autorização prévia por escrito.

### Dependências

Todas as dependências do projeto devem ser compatíveis com Apache 2.0 ou licenças permissivas similares.

---

## 5. Tomada de Decisões

### Tipos de Decisão

| Tipo | threshold | Quem |
|------|------------|------|
| Minor | 1 aprovador | Maintainer |
| Major | 2 aprovadores | PMC |
| Strategic | Consensus | Fundação |

### RFC Process

Para mudanças maiores:
1. Criar RFC em `/rfcs/`
2. Período de discussão (2 semanas)
3. Votação do PMC
4. Decisão final pela Fundação

---

## 6. Transparência

### Repositórios Públicos

- Código fonte
- Documentação
- RFCs
- Decisões de governance

### Comunicação

- Issues abertas para bugs e features
- Discussions para вопросы
- Meetings gravadas (quando aplicável)

---

## 7. Conflito de Interesses

Contribuidores devem declarar conflitos de interesse potenciais. Decisões afetadas por COI devem ser tratadas por membros não envolvidos.

---

## Referências

- [Apache Foundation Way](https://www.apache.org/foundation/how-it-works.html)
- [Contributor Covenant](https://www.contributor-covenant.org/)
- [Go Community Code of Conduct](https://go.dev/conduct)
