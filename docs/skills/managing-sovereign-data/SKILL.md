### 4. `managing-sovereign-data/SKILL.md`

**Foco:** Isolamento SQLite e Ciclo de Vida do dado.

```yaml
name: managing-sovereign-data
description: Use esta habilidade para gerenciar bancos de dados SQLite no Digna. Garante o isolamento físico por tenant e o poder absoluto de saída (Exit Power) dos dados.

```

# 🛡️ SKILL: Soberania de Dados e Infraestrutura (Local-First & Lifecycle)

**Propósito:** Garantir a soberania tecnológica dos empreendimentos através da arquitetura de micro-databases isolados, assegurando que o usuário detenha o controle absoluto e o "Poder de Saída" de suas informações.

---

## 1. Arquitetura de Micro-databases Isolados

O agente deve seguir o princípio de isolamento físico total para garantir que os dados de uma cooperativa nunca se misturem com outra.

* **Um Banco por Entidade:** Cada "Sonho" ou Cooperativa deve possuir seu próprio arquivo `.sqlite` exclusivo.
* **Isolamento Físico:** Os arquivos devem ser armazenados em `data/entities/{entity_id}.db`.
* **Proibição de Cruzamento:** É terminantemente proibido realizar *JOINs* SQL entre diferentes bancos de dados de entidades.
* **Exit Power:** O sistema deve ser projetado para que o usuário possa simplesmente copiar o arquivo `.sqlite` e levá-lo para qualquer outro sistema compatível com SQL, sem dependência da plataforma Digna.

## 2. Gestão de Ciclo de Vida (Lifecycle Manager)

O acesso ao banco de dados não deve ser feito de forma avulsa; ele deve ser orquestrado centralmente.

* **Ponto de Entrada Único:** Toda criação, abertura ou migração de banco de dados deve obrigatoriamente passar pelo `LifecycleManager`.
* **Migrações Versionadas:** Mudanças no schema devem ser aplicadas de forma segura e idempotente em cada arquivo isolado, garantindo que versões diferentes do app não corrompam os dados.
* **Shutdown Seguro:** O agente deve assegurar que todas as conexões SQLite sejam encerradas corretamente (`Graceful Shutdown`) para evitar a corrupção de arquivos em caso de interrupção do servidor.

## 3. Protocolo de Sincronização e Privacidade (Delta Sync)

O sistema deve ser resiliente à baixa conectividade e protetor da privacidade individual.

* **Sincronização Baseada em Deltas:** Transmitir apenas as alterações (deltas) desde a última sincronização bem-sucedida, economizando dados.
* **Agregação para Nuvem:** Para fins de políticas públicas e monitoramento nacional, apenas dados **agregados** e não sensíveis podem ser sincronizados com a nuvem soberana (ex: totais de vendas, sem detalhamento de membros).
* **Proteção de Identidade:** IDs de membros e detalhes transacionais internos de autogestão nunca devem deixar o dispositivo local sem o consentimento explícito e criptografia.

## 4. Resiliência Local-First

A interface deve priorizar o banco de dados local para todas as operações críticas.

* **Offline-First:** O agente deve garantir que o PDV, o Registro de Trabalho e o Caixa funcionem 100% sem internet, utilizando o SQLite local como a única fonte da verdade em tempo real.
* **Integridade Criptográfica:** Cada bloco de decisão ou contabilidade deve gerar um hash SHA256 para auditoria futura e prevenção de adulteração de dados históricos.

---

## 📈 Plano de Ação: Estado Atual

| Fase | Status | Atividade |
| --- | --- | --- |
| **0. Roadmap** | [x] | Mapeamento de todas as Skills necessárias. |
| **1. Backend** | [x] | Detalhamento da **SKILL_BACKEND_GO.md** concluído. |
| **2. Frontend** | [x] | Detalhamento da **SKILL_FRONTEND_HTMX.md** concluído. |
| **3. Negócio** | [x] | Detalhamento da **SKILL_NEGOCIO_SOCIAL.md** concluído. |
| **4. Infra** | [x] | Detalhamento da **SKILL_SOBERANIA_DATA.md** concluído. |
| **5. Ponte** | [ ] | Detalhamento da Skill de Auditoria e SPED. |

> **Estado Atual do Plano:**
> [x] SKILL_SOBERANIA_DATA.md concluída.
> [ ] Próximo: **SKILL_AUDITORIA_FISCAL.md** (Accountant Dashboard, SPED e Aliança Contábil).

Esta skill blinda a base do sistema. **Podemos finalizar nosso plano com a `SKILL_AUDITORIA_FISCAL.md`, que transformará o agente em um especialista em exportação de dados para o governo e na interface de parceria com os contadores?**