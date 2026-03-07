---
title: Modelos de Domínio e Dados
status: implemented
version: 1.0
last_updated: 2026-03-07
---

# Modelos - Projeto Digna

---

## 1. Domain Model

### 1.1 Entidades

#### Enterprise
Empreendimento de economia solidária. Pode estar em dois estados:
- **DREAM** - Grupo informal/em formação
- **FORMALIZED** - Cooperativa/Associação formalizada

```
Enterprise
├── id: string (UUID)
├── name: string
├── status: enum (DREAM | FORMALIZED)
├── cnpj: string (nullable)
├── created_at: timestamp
└── sync_metadata: SyncMetadata
```

#### Member
Pessoa participante do empreendimento.

```
Member
├── id: string (UUID)
├── enterprise_id: string (FK)
├── name: string
├── role: enum (MEMBER | COORDINATOR)
└── work_logs: WorkLog[]
```

#### Transaction
Evento econômico (venda, compra, despesa).

```
Transaction
├── id: string (UUID)
├── enterprise_id: string (FK)
├── type: enum (SALE | PURCHASE | EXPENSE)
├── amount: int64 (centavos)
├── date: timestamp
├── description: string
└── postings: Posting[]
```

#### WorkLog
Registro de trabalho cooperativo (ITG 2002).

```
WorkLog
├── id: string (UUID)
├── member_id: string (FK)
├── enterprise_id: string (FK)
├── minutes: int64
├── activity_type: string
├── date: timestamp
└── hash: string (SHA256)
```

#### Decision
Decisão coletiva registrada em assembleia (CADSOL).

```
Decision
├── id: string (UUID)
├── enterprise_id: string (FK)
├── title: string
├── content: string
├── content_hash: string (SHA256)
├── status: enum (PENDING | APPROVED | REJECTED)
├── decided_at: timestamp
└── decided_by: string
```

#### Fund
Fundos obrigatórios constituídos.

```
Fund
├── id: string (UUID)
├── enterprise_id: string (FK)
├── type: enum (RESERVE_LEGAL | FATES)
├── amount: int64
├── created_at: timestamp
└── period: string (YYYY-MM)
```

### 1.2 Aggregates

```
Enterprise (Aggregate Root)
├── Members[]
├── Transactions[]
├── WorkLogs[]
├── Decisions[]
└── Funds[]
```

### 1.3 Value Objects

| Value Object | Descrição |
|--------------|------------|
| Money | Valor monetário em centavos (int64) |
| AccountCode | Código de conta contábil (ex: 1.1.01) |
| Period | Período contábil (YYYY-MM) |

---

## 2. Data Model

### 2.1 Tabelas do Schema v0

#### accounts
Plano de contas hierárquico.

```sql
CREATE TABLE accounts (
    id INTEGER PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    parent_id INTEGER REFERENCES accounts(id),
    account_type TEXT NOT NULL CHECK(account_type IN ('ASSET', 'LIABILITY', 'REVENUE', 'EXPENSE', 'EQUITY'))
);
```

#### entries
Lançamentos contábeis.

```sql
CREATE TABLE entries (
    id INTEGER PRIMARY KEY,
    entry_date TEXT NOT NULL,
    description TEXT,
    reference TEXT,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP
);
```

#### postings
Partidas dobradas.

```sql
CREATE TABLE postings (
    id INTEGER PRIMARY KEY,
    entry_id INTEGER NOT NULL REFERENCES entries(id),
    account_id INTEGER NOT NULL REFERENCES accounts(id),
    amount INTEGER NOT NULL,
    direction TEXT NOT NULL CHECK(direction IN ('DEBIT', 'CREDIT'))
);
```

#### work_logs
Registro de trabalho (ITG 2002).

```sql
CREATE TABLE work_logs (
    id INTEGER PRIMARY KEY,
    member_id TEXT NOT NULL,
    enterprise_id TEXT NOT NULL,
    minutes INTEGER NOT NULL,
    activity_type TEXT,
    work_date TEXT NOT NULL,
    hash TEXT NOT NULL,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP
);
```

#### decisions_log
Registro de decisões (CADSOL).

```sql
CREATE TABLE decisions_log (
    id INTEGER PRIMARY KEY,
    enterprise_id TEXT NOT NULL,
    title TEXT NOT NULL,
    content TEXT,
    content_hash TEXT NOT NULL,
    status TEXT DEFAULT 'APPROVED',
    decided_at TEXT DEFAULT CURRENT_TIMESTAMP
);
```

#### sync_metadata
Metadados de sincronização.

```sql
CREATE TABLE sync_metadata (
    id INTEGER PRIMARY KEY,
    enterprise_id TEXT NOT NULL UNIQUE,
    status TEXT DEFAULT 'DREAM',
    last_sync_at TEXT,
    version INTEGER DEFAULT 1,
    chain_digest TEXT
);
```

### 2.2 Relacionamentos

```
entries
    └── postings (1:N)
    
enterprise
    └── work_logs (1:N)
    └── decisions_log (1:N)
    └── sync_metadata (1:1)
```

---

## 3. Algoritmos

### 3.1 Algoritmo de Rateio Social

**Objetivo:** Distribuir excedente financeiro proporcionalmente às horas trabalhadas.

**Entrada:**
- `totalSurplus`: Excedente financeiro disponível (int64)
- `memberHours`: Mapa de member_id → minutos trabalhados

**Processo:**

```
1. Calcular total de minutos:
   totalMinutes = SUM(memberHours.values())

2. Para cada membro:
   percentage = memberMinutes / totalMinutes
   share = (percentage * totalSurplus) / 100
   
3. Retornar distribuição:
   { member_id: { percentage, share } }
```

**Exemplo:**

| Sócio | Minutos | % do Total | Valor (R$ 100,00) |
|-------|---------|------------|-------------------|
| socio_001 | 600 | 66.67% | R$ 66.66 |
| socio_002 | 300 | 33.33% | R$ 33.33 |

**Implementação:**
```go
func CalculateSocialSurplus(totalSurplus int64, workLogs []WorkLog) map[string]SurplusShare {
    totalMinutes := int64(0)
    memberMinutes := make(map[string]int64)
    
    for _, log := range workLogs {
        memberMinutes[log.MemberID] += log.Minutes
        totalMinutes += log.Minutes
    }
    
    result := make(map[string]SurplusShare)
    for memberID, minutes := range memberMinutes {
        percentage := float64(minutes) / float64(totalMinutes) * 100
        share := (int64(percentage * float64(totalSurplus))) / 100
        result[memberID] = SurplusShare{
            Percentage: percentage,
            Share:      share,
        }
    }
    return result
}
```

---

### 3.2 Algoritmo de Partidas Dobradas

**Objetivo:** Validar que soma(débitos) = soma(créditos) = 0.

**Entrada:**
- `postings`: Lista de partidas

**Processo:**

```
1. Separar débitos e créditos:
   debits = SUM(postings where direction = DEBIT)
   credits = SUM(postings where direction = CREDIT)

2. Validar equilíbrio:
   IF (debits + credits) != 0 THEN
       REJECT transaction
   END IF

3. Persistir atomicamente
```

---

### 3.3 Algoritmo de Delta Detection

**Objetivo:** Detectar alterações desde última sincronização.

**Entrada:**
- `entityID`: ID da entidade
- `lastSyncTimestamp`: Timestamp da última sincronização

**Processo:**

```
1. Query alterações por tabela:
   entries = SELECT COUNT(*) FROM entries WHERE created_at > lastSync
   work_logs = SELECT COUNT(*) FROM work_logs WHERE created_at > lastSync
   decisions = SELECT COUNT(*) FROM decisions_log WHERE decided_at > lastSync

2. Calcular chain digest:
   digest = SHA256(concat(all_entry_hashes, all_decision_hashes))

3. Retornar pacote:
   { delta_count, digest, timestamp }
```

---

### 3.4 Algoritmo de Formalização

**Objetivo:** Transicionar entidade de DREAM para FORMALIZED.

**Critérios:**
- Mínimo de 3 decisões registradas
- Mínimo de 1 membro
- Histórico mínimo de operações

**Processo:**

```
1. Verificar critérios:
   decisionCount = SELECT COUNT(*) FROM decisions_log WHERE enterprise_id = ?
   
   IF decisionCount >= 3 THEN
       UPDATE sync_metadata SET status = 'FORMALIZED'
       Generate mock CNPJ
       RETURN formalization_success
   ELSE
       RETURN formalization_pending
   END IF
```

---

### 3.5 Algoritmo de Geração de Ata

**Objetivo:** Gerar documento de ata de assembleia em Markdown.

**Entrada:**
- `enterprise`: Dados da entidade
- `decisions`: Lista de decisões da assembleia

**Processo:**

```
1. Header com dados da entidade:
   # Ata de Assembleia
   Entity: {name}
   Date: {current_date}
   Status: {status}

2. Listar decisões:
   ## Pautas
   - {decision.title} - {decision.status}
     {decision.content}

3. Gerar hash de auditoria:
   content_hash = SHA256(concat(all_decisions))
   
4. Assinatura final:
   ---
   Hash de Auditoria: {content_hash}
```

---

## 4. Seed Data

### 4.1 Contas Padrão

| ID | Código | Nome | Tipo |
|----|--------|------|------|
| 1 | 1.1.01 | Caixa e Equivalentes | ASSET |
| 2 | 3.1.01 | Receita de Vendas | REVENUE |
| 3 | 1.1.02 | Bancos | ASSET |
| 4 | 2.1.01 | Fornecedores | LIABILITY |
