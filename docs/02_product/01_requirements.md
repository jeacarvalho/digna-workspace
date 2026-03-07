---
title: Requisitos do Projeto Digna
status: implemented
version: 1.0
last_updated: 2026-03-07
---

# Requisitos - Projeto Digna

**Referência Legal:** Lei nº 15.068/2024 (Lei Paul Singer) / ITG 2002  
**Projeto:** Digna - Infraestrutura Contábil para Economia Solidária

---

## 1. Requisitos Funcionais (RF)

### RF-01: PDV Operacional Soberano

**Descrição:** Interface simplificada para registro de Vendas e Compras.

**Regra de Negócio:**
- Cada venda deve gerar automaticamente uma "Entry" e dois "Postings" (Débito/Crédito) no Ledger
- Partidas dobradas: Débito em conta de Ativo / Crédito em conta de Receita

**Prioridade:** Essencial (v0)

---

### RF-02: Registro de Trabalho (ITG 2002)

**Descrição:** Captura de horas/minutos trabalhados por membros do empreendimento.

**Regra de Negócio:**
- Tempo deve ser registrado em `int64` (minutos)
- Tempo registrado constitui Capital Social de Trabalho
- Base primária para rateio de sobras sociais

**Prioridade:** Essencial (v0)

---

### RF-03: Motor de Reservas Obrigatórias

**Descrição:** Segregação automática de fundos antes da distribuição de resultados.

**Regra de Negócio:**
- Bloqueio mandatório de **10%** para Reserva Legal
- Bloqueio mandatório de **5%** para FATES (Fundo de Assistência Técnica)
- Cálculo sobre o excedente financeiro do período

**Prioridade:** Essencial (v0)

---

### RF-04: Dossiê de Formalização (CADSOL)

**Descrição:** Exportação de Atas de Assembleia e Relatórios de Impacto Social.

**Regra de Negócio:**
- Documentos gerados em formato Markdown
- Hash SHA256 para integridade
- Critérios de formalização: mínimo 3 decisões registradas

**Prioridade:** Essencial (v0)

---

### RF-05: Sincronização Offline-First

**Descrição:** Operação sem internet com sincronização posterior.

**Regra de Negócio:**
- Interface PWA deve permitir operações básicas offline
- Delta tracking para detectar alterações
- Sync apenas com dados agregados (privacidade)

**Prioridade:** Alta

---

### RF-06: Intercooperação B2B

**Descrição:** Marketplace para troca de produtos entre cooperativas.

**Regra de Negócio:**
- Ofertas entre entidades registradas
- Dados agregados apenas (sem exposição de membros)

**Prioridade:** Média

---

### RF-07: Gestão de Compras

**Descrição:** Registro de aquisições de insumos e mercadorias.

**Regra de Negócio:**
- Cada compra gera lançamento contábil (Débito em despesa/ativo, Crédito em fornecedores/caixa)
- Controle de fornecedores cadastrados
- Histórico de compras por período

**Prioridade:** Alta

---

### RF-08: Gestão de Estoque

**Descrição:** Controle de produtos e insumos.

**Regra de Negócio:**
- Cadastro de produtos com custo e preço
- Entrada e saída de estoque
- Controle de saldo mínimo
- Relatório de inventário

**Prioridade:** Média

---

### RF-09: Gestão de Caixa

**Descrição:** Controle financeiro simplificado.

**Regra de Negócio:**
- Registro de entradas e saídas
- Saldo atual em tempo real
- Extrato por período
- Conciliação bancária básica

**Prioridade:** Alta

---

### RF-10: Gestão Orçamentária

**Descrição:** Planejamento e acompanhamento orçamentário.

**Regra de Negócio:**
- Definição de orçamento por categoria
- Acompanhamento de execução
- Alertas de dépassamento
- Relatório comparativo

**Prioridade:** Média

---

## 2. Requisitos Não Funcionais (RNF)

### RNF-01: Isolamento por Tenant (Soberania)

**Descrição:** Cada "Sonho" ou Cooperativa deve ter seu próprio arquivo `.sqlite`.

**Métrica:** Tempo de criação de nova instância < 500ms

**Implementação:**
- Arquivos em `data/entities/{entity_id}.db`
- Isolamento físico total entre tenants

---

### RNF-02: Rigor Monetário (Anti-Float)

**Descrição:** Proibição de tipos de ponto flutuante em todo o Core Lume.

**Regra:** Uso obrigatório de `int64` para representar centavos

**Validação:** Qualquer código com `float32` ou `float64` em módulos contábeis deve ser rejeitado em review

---

### RNF-03: Resiliência (Offline-First)

**Descrição:** A interface (PWA) deve permitir operações básicas sem internet.

**Implementação:**
- Service Worker com Cache First
- Delta tracking local
- Background sync quando conectado

---

### RNF-04: Performance

**Descrição:** Resposta rápida para operações do dia a dia.

**Metas:**
- Criação de tenant: < 500ms
- Registro de venda: < 100ms
- Consulta de saldo: < 50ms

---

### RNF-05: Escalabilidade

**Descrição:** Suporte a milhões de empreendimentos.

**Arquitetura:**
- Banco isolado por entidade
- Sem dependência de banco central
- Serpro como infraestrutura de nuvem soberana

---

### RNF-06: Segurança

**Descrição:** Proteção de dados e integridade.

**Implementação:**
- Isolamento de dados por entidade
- Hash SHA256 para auditoria
- Assinatura digital de pacotes de sync

---

### RNF-07: Usabilidade

**Descrição:** Interface acessível para usuários não técnicos.

**Metas:**
- Tempo de treinamento < 2 horas
- Acessibilidade (WCAG 2.1 AA)
- Design mobile-first para uso em campo
- Suporte a usuários com baixa-literacia digital

---

### RNF-08: Conformidade Legal

**Descrição:** Adequação às normas legais brasileiras.

**Requisitos:**
- Conformidade com ITG 2002 (CFC)
- Conformidade com Lei Paul Singer (15.068/2024)
- Geração de dados para CADSOL/DCSOL
- Suporte a auditoria pública

---

### RNF-09: Manutenibilidade

**Descrição:** Facilidade de evolução do sistema.

**Requisitos:**
- Arquitetura modular (Clean Architecture)
- Código testável (>80% cobertura)
- Documentação atualizada
- Logging estruturado

---

### RNF-10: Interoperabilidade

**Descrição:** Capacidade de integração com sistemas externos.

**Requisitos:**
- API REST documentada (OpenAPI)
- Formatos padrão (JSON)
- Webhooks para eventos
- Exportação de dados (CSV, JSON)

---

## 3. Matriz de Rastreabilidade

| Operação | Gatilho | Impacto no Ledger | Requisito |
|----------|----------|-------------------|------------|
| Venda no Balcão | PDV Submit | D: Ativo / C: Receita | RF-01 |
| Registro de Compra | Entrada Mercadoria | D: Despesa / C: Fornecedor | RF-07 |
| Ajuste de Estoque | Inventário | D/C conforme natureza | RF-08 |
| Fim de Turno | Log Horas | Registro de Cota-Trabalho | RF-02 |
| Fechamento Mes | Batch Job | Cálculo de Reservas (15%) | RF-03 |
| Movimento de Caixa | Entrada/Saída | Atualização de saldo | RF-09 |
| Orçamento | Definição Plano | Acompanhamento | RF-10 |
| Assembleia | Decisão | Hash em decisions_log | RF-04 |

---

## 4. Casos de Uso

### UC-01: Registro de Venda

```
Ator: Empreendedor
1. Acessa tela PDV
2. Insere valor da venda
3. Seleciona produto
4. Confirma operação
5. Sistema gera:
   - Entry contábil
   - 2 Postings (Débito/Crédito)
6. Retorna saldo atualizado
```

### UC-02: Registro de Trabalho

```
Ator: Cooperado
1. Acessa Social Clock
2. Inicia cronômetro ou insere minutos
3. Sistema registra:
   - Member ID
   - Minutos trabalhados
   - Data/Hora
4. Minutos convertidos em Capital Social
```

### UC-03: Rateio de Sobras

```
Ator: Sistema (batch)
1. Calcula total de horas do período
2. Calcula excedente financeiro
3. Aplica deduções (Reserva Legal + FATES)
4. Rateia restante proporcionalmente:
   - (Horas Sócio / Total Horas) × Excedente
5. Gera relatório de distribuição
```

### UC-04: Formalização

```
Ator: Empreendedor
1. Sistema detecta 3+ decisões registradas
2. Altera status: DREAM → FORMALIZED
3. Gera CNPJ mock (00.000.000/0000-00)
4. Disponibiliza documentos:
   - Ata de Constituição
   - Estatuto Social
```

### UC-05: Registro de Compra

```
Ator: Empreendedor
1. Acessa tela de compras
2. Seleciona fornecedor
3. Insere itens e valores
4. Confirma operação
5. Sistema gera:
   - Entry contábil
   - 2 Postings (Débito/Crédito)
6. Atualiza estoque (RF-08)
```

### UC-06: Controle de Estoque

```
Ator: Empreendedor
1. Acessa gestão de estoque
2. Realiza entrada/saída de produto
3. Sistema atualiza saldo
4. Alerta se saldo mínimo atingido
```

### UC-07: Gestão de Caixa

```
Ator: Empreendedor
1. Acessa tela de caixa
2. Registra entrada ou saída
3. Sistema atualiza saldo em tempo real
4. Gera extrato do período
```

### UC-08: Planejamento Orçamentário

```
Ator: Coordenador
1. Define orçamento por categoria
2. Sistema monitora execução
3. Alerta de dépassamento
4. Gera relatório comparativo
```

---

## 5. Glossário de Termos

| Termo | Definição |
|-------|-----------|
| **Tenant** | Entidade (cooperativa/grupo) com banco isolado |
| **Entry** | Lançamento contábil completo |
| **Posting** | Partida dobrada (débito ou crédito) |
| **ITG 2002** | Norma de contabilidade para economia solidária |
| **CADSOL** | Cadastro Nacional de Economia Solidária |
| **FATES** | Fundo de Assistência Técnica e Extensão Rural |
| **Reserva Legal** | Fundo obrigatório de 10% das sobras |
| **Chain Digest** | Hash da cadeia contábil para auditoria |
