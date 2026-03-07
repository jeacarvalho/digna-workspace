---
title: Stakeholders e Riscos
status: implemented
version: 1.0
last_updated: 2026-03-07
---

# Stakeholders e Riscos - Digna

---

## 1. Stakeholder Map

### 1.1 Stakeholders Primários

| Stakeholder | Role |
|-------------|------|
| Cooperativas | Usuários principais |
| Associações | Gestão coletiva |
| Grupos informais | Entrada no sistema |

### 1.2 Stakeholders Institucionais

| Stakeholder | Interest |
|-------------|----------|
| Ministério do Trabalho | Política pública |
| Senaes | Economia solidária |
| Serpro | Infraestrutura tecnológica |

### 1.3 Stakeholders de Suporte

| Stakeholder | Role |
|-------------|------|
| Universidades | Pesquisa |
| ONGs | Apoio territorial |
| Incubadoras | Formação cooperativa |

### 1.4 Stakeholders de Governança

| Stakeholder | Role |
|-------------|------|
| Fundação Providentia | Governança |
| PMC | Decisões técnicas |
| Comunidade dev | Evolução do software |

---

## 2. Risk Register

### 2.1 Riscos de Projeto

| Risco | Probabilidade | Impacto | Mitigação |
|-------|---------------|---------|-----------|
| Perda de dados locais | Medium | High | Backup automatizado (Litestream) |
| Uso incorreto do sistema | Medium | Medium | UX simplificada |
| Bugs contábeis | Low | Critical | Testes contábeis rigorosos |
| Complexidade institucional | Medium | High | Documentação clara |
| Dependência tecnológica | Low | High | Open source |

### 2.2 Riscos Técnicos

| Risco | Probabilidade | Impacto | Mitigação |
|-------|---------------|---------|-----------|
| Inconsistência no ledger | Low | Critical | Validação soma zero |
| Falhas na sincronização | Medium | Medium | Retry com backoff |
| Corrupção de banco SQLite | Low | High | WAL mode + backup |

### 2.3 Riscos de Governança

| Risco | Probabilidade | Impacto | Mitigação |
|-------|---------------|---------|-----------|
| Captura institucional | Low | High | Modelo Apache |
| Fragmentação da comunidade | Medium | Medium | Comunicação ativa |

---

## 3. Matriz de Responsabilidade

| Atividade | Fundação | PMC | Comunidade |
|-----------|----------|-----|------------|
| Roadmap estratégico | Decisão | Input | Input |
| Decisões técnicas | Veto | Decisão | Input |
| Implementação | - | Review | Execução |
| Documentação | - | Aprovação | Contribuição |
| Suporte | Oversight | - | Community |
