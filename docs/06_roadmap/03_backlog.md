---
title: Product Backlog
status: proposed
version: 1.1
last_updated: 2026-03-08
---

# Product Backlog - Digna

---

## Funcionalidades Concluídas

### Core Contábil
- [x] PDV Operacional (RF-01)
- [x] Registro de Trabalho ITG 2002 (RF-02)
- [x] Motor de Reservas Obrigatórias (RF-03)
- [x] Dossiê de Formalização (RF-04)
- [x] Sincronização Offline-First (RF-05)
- [x] Intercooperação B2B (RF-06)
- [x] Gestão de Caixa (RF-09)

### Gestão de Membros (Sprint 10)
- [x] Cadastro de membros
- [x] Roles (COORDINATOR, MEMBER, ADVISOR)
- [x] Status (ACTIVE, INACTIVE)
- [x] Skills/habilidades

### Formalização (Sprint 11)
- [x] Transição automática DREAM → FORMALIZED
- [x] Simulador de formalização
- [x] CheckFormalizationCriteria()
- [x] AutoTransitionIfReady()

### Integrações Governamentais
- [x] Receita Federal (CNPJ, DARF) - Mock
- [x] MTE (RAIS, CAT, eSocial) - Mock
- [x] MDS (CadÚnico, Relatório Social) - Mock
- [x] IBGE (Pesquisas, PAM, CNAE) - Mock
- [x] SEFAZ (NFe, NFS-e, Manifesto) - Mock
- [x] BNDES (Linhas de Crédito, Simulação) - Mock
- [x] SEBRAE (Cursos, Consultoria) - Mock
- [x] Providentia (Sync, Marketplace) - Mock

### Surplus Calculator
- [x] CalculateSocialSurplus()
- [x] CalculateWithDeductions() - 10% Reserva Legal + 5% FATES
- [x] Rateio proporcional por minutos trabalhados
- [x] Tratamento de resíduos

### Testes
- [x] Testes unitários por módulo
- [x] Testes E2E Journey (jornada anual)
- [x] Testes de integrações governamentais

---

## Funcionalidades Futuras

### Alta Prioridade

- [ ] Autenticação Gov.br
- [ ] Exportação CADSOL completa
- [ ] Auditoria pública
- [ ] Integrações reais (APIs governo)

### Média Prioridade

- [ ] Gestão de membros (UI)
- [ ] Gestão de compras (RF-07)
- [ ] Gestão de fornecedores
- [ ] Controle de estoque (RF-08)
- [ ] Gestão orçamentária (RF-10)
- [ ] Múltiplas moedas sociais

### Baixa Prioridade

- [ ] Integração contábil fiscal
- [ ] API GraphQL
- [ ] App mobile nativo
- [ ] Relatórios em PDF
