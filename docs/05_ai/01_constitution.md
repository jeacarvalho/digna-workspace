---
title: Constituição de IA e Agentes
status: implemented
version: 1.0
last_updated: 2026-03-07
---

# Constituição de IA e Agentes - Digna

---

## 1. Contexto do Projeto

### 1.1 Visão Geral

**Projeto:** Sistema de Gestão Contábil Público e Open Source para Economia Solidária  
**Propósito:** Transformar a contabilidade de um peso burocrático em uma ferramenta de **Dignidade Financeira** para Empreendimentos de Economia Solidária (EES)  
**Diferencial:** Arquitetura *Local-First Server-Side* com isolamento físico de dados

### 1.2 Componentes Chave

- **Lume (Motor):** O backend em Go
- **Digna (App):** A interface (futura)
- **Providentia (Fundação):** A entidade de governança

---

## 2. Regras de Ouro (Não Negociáveis)

### Regra 1: Integridade Financeira

- **PROIBIDO** o uso de `float32/64`
- **OBRIGATÓRIO** o uso de `int64` (centavos)
- Validação de **Soma Zero** em todos os lançamentos do Ledger

### Regra 2: Nomenclatura de Arquivos

- **PROIBIDO** o uso de espaços em nomes de diretórios ou arquivos
- **PADRÃO:** `kebab-case` para diretórios e `snake_case` para arquivos `.go`

### Regra 3: Isolamento de Dados

- O acesso ao arquivo `.sqlite` é exclusivo via `LifecycleManager`
- Cada Tenant deve ser isolado fisicamente em `data/entities/{entity_id}.db`

---

## 3. Regras para Agentes

### 3.1 Princípios de Operação

1. **Soberania de Dados:** O dado pertence à entidade. O acesso deve ser isolado por arquivo
2. **Padrão de Lançamento:** Seguir rigorosamente o princípio de partidas dobradas (soma zero)
3. **Mocks de Formalização:** O sistema deve prever a transição `DREAM` → `FORMALIZED` via interfaces (Facade)

### 3.2 Regras de Execução

1. **Contexto:** Leia sempre `docs/03_architecture/01_system.md` antes de sugerir mudanças
2. **Espaço:** Nunca use espaços em caminhos de arquivos
3. **Finanças:** Se vir um `float` no código contábil, pare e corrija para `int64`
4. **Isolamento:** Código do módulo `lume` não deve criar arquivos; peça ao `lifecycle`

### 3.3 Ferramental Recomendado

- **IDE:** VS Code com extensão oficial `golang.go`
- **Workspace:** `go.work` para orquestrar os módulos
- **Formatação:** `gofmt` antes de commit

---

## 4. Stack Tecnológica

| Camada | Tecnologia |
|--------|------------|
| Linguagem | Go 1.22+ |
| Database | SQLite3 (Individual por Tenant) |
| Escala | Alvo de 1 milhão de EES |
| Regra Financeira | Tudo em `int64` (centavos). Proibido `float` |

---

## 5. Estratégia de Implementação para Agentes

### 5.1 Isolamento de Sprint

Cada Módulo deve ter seu próprio `SESSION_LOG` e ser finalizado com testes unitários antes do início do próximo.

### 5.2 Interface First

Antes de codificar a lógica interna, o agente deve definir as interfaces no pacote `internal/domain`.

### 5.3 Validation

O agente deve validar o Módulo 1 (Lifecycle) com a criação física de arquivos antes de tentar realizar qualquer lançamento contábil (Módulo 2).

---

## 6. Padrão de Sessão

### 6.1 Início de Sessão

1. Ler documentação relevante
2. Verificar status atual das sprints
3. Entender dependências

### 6.2 Durante a Sessão

1. Implementar funcionalidade
2. Executar testes
3. Validar contra requisitos

### 6.3 Fim de Sessão

1. Atualizar SESSION_LOG
2. Registrar decisões
3. Documentar mudanças arquiteturais
4. Listar próximos passos

---

## 7. Status Atual das Sprints

| Sprint | Status | Descrição |
|--------|--------|-----------|
| Sprint 00 | ✅ COMPLETE | Blueprint, Governança e Visão de Produto |
| Sprint 01 | ✅ COMPLETE | Lifecycle Manager & Core Ledger |
| Sprint 02 | ✅ COMPLETE | Operação & Contabilidade Invisível |
| Sprint 03 | ✅ COMPLETE | Reporting & Documents |
| Sprint 04 | ✅ COMPLETE | Sync Engine |
| Sprint 05 | ✅ COMPLETE | UI Web (PWA) |
| **Total** | **40/40 PASS** | |

---

## 8. Glossário para Agentes

| Termo | Significado |
|-------|-------------|
| Tenant | Entidade (cooperativa) com banco isolado |
| Ledger | Livro contábil com partidas dobradas |
| ITG 2002 | Norma de contabilidade para economia solidária |
| CADSOL | Sistema de autogestão |
| PDV | Ponto de Venda |
| PWA | Progressive Web App |
| int64 | Tipo inteiro de 64 bits (para dinheiro) |
