### 5. `auditing-fiscal-compliance/SKILL.md`

**Foco:** Accountant Dashboard, SPED e Aliança Contábil.

```yaml
name: auditing-fiscal-compliance
description: Use esta habilidade para operar o Painel do Contador e exportações fiscais. Garante acesso Read-Only e mapeamento para o Plano de Contas Referencial.

```
# 🏛️ SKILL: Auditoria e Aliança Contábil (Dashboard & SPED)

**Propósito:** Atuar como o motor de tradução entre a contabilidade social invisível e as obrigações fiscais do Estado, facilitando o trabalho de "Contadores Sociais" e garantindo a conformidade com a ITG 2002.

---

## 1. Painel do Contador Social (Multi-tenant & Read-Only)

O agente deve gerenciar a interface de auditoria garantindo que o profissional parceiro tenha visão ampla sem comprometer a autonomia do produtor.

* **Acesso Estritamente Read-Only:** Toda conexão estabelecida para fins de auditoria no módulo `accountant_dashboard` **deve** utilizar obrigatoriamente o modo `?mode=ro` no SQLite.
* **Visão Multi-tenant:** O agente deve permitir que um contador parceiro visualize o status de fechamento e conformidade de múltiplos empreendimentos (arquivos `.db` isolados) em uma única tela consolidada.
* **Independência de Dados:** O painel não deve possuir banco de transações próprio; ele deve consumir os micro-databases autorizados em tempo real.

## 2. Motor de Tradução Fiscal (SPED & CSV)

O agente deve ser capaz de traduzir a linguagem coloquial do PDV para os leiautes técnicos exigidos pela Receita Federal.

* **Mapeamento de Contas:** Realizar o vínculo automático entre as "contas amigáveis" (ex: Gaveta) e o Plano de Contas Referencial do SPED (ex: Disponibilidades).
* **Geração de Lotes Fiscais:** Compilar as *Entries* de um período fechado e gerar arquivos estruturados (SPED, CSV ou TXT) para importação em softwares contábeis comerciais.
* **Blindagem Fiscal do Core:** O agente está proibido de codificar cálculos de impostos dentro do Motor Lume; essa responsabilidade é delegada ao software do contador, alimentado pelos dados íntegros exportados pelo Digna.

## 3. Integridade e Auditoria de Norma (ITG 2002)

O sistema deve atuar como um validador preventivo para o contador.

* **Validação de Soma Zero:** Antes de qualquer exportação, o agente deve validar que cada lançamento exportado mantém o equilíbrio rigoroso de débitos e créditos.
* **Check de Fundos:** Auditar visualmente se as travas de Reserva Legal (10%) e FATES (5%) foram aplicadas corretamente conforme a norma ITG 2002.
* **Hash de Exportação:** Registrar o hash SHA256 de cada lote exportado na tabela `fiscal_exports` para evitar envios duplicados e garantir a rastreabilidade.

---

## 📈 Plano de Ação: Estado Atual

| Fase | Status | Atividade |
| --- | --- | --- |
| **0. Roadmap** | [x] | Mapeamento de todas as Skills necessárias. |
| **1. Backend** | [x] | Detalhamento da **SKILL_BACKEND_GO.md** concluído. |
| **2. Frontend** | [x] | Detalhamento da **SKILL_FRONTEND_HTMX.md** concluído. |
| **3. Negócio** | [x] | Detalhamento da **SKILL_NEGOCIO_SOCIAL.md** concluído. |
| **4. Infra** | [x] | Detalhamento da **SKILL_SOBERANIA_DATA.md** concluído. |
| **5. Ponte** | [x] | Detalhamento da **SKILL_AUDITORIA_FISCAL.md** concluído. |

> **Estado Atual do Plano:**
> [x] **Todas as Skills foram detalhadas e concluídas.**
> **Próximo Passo:** Consolidar este ecossistema de habilidades no seu ambiente de desenvolvimento.

---
