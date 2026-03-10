### 2. `rendering-digna-frontend/SKILL.md`

**Foco:** HTMX, UI "Soberania e Suor" e blindagem contra cache.

```yaml
name: rendering-digna-frontend
description: Use esta habilidade para desenvolver a interface web do Digna usando HTMX e templates "Cache-Proof". Garante que a versão correta do template seja carregada diretamente do disco.

```

# 🎨 SKILL: Especialista Frontend (HTMX & Templates Cache-Proof)

**Propósito:** Garantir uma interface reativa, leve e visualmente consistente, eliminando erros de versão através da arquitetura de carregamento direto do disco e respeitando a identidade "Soberania e Suor".

---

## 1. Protocolo Anti-Cache (Cache-Proof Architecture)

O agente deve ignorar os métodos tradicionais de cache de templates do Go para garantir que as alterações sejam refletidas instantaneamente.

* **Proibição de Globais:** É estritamente proibido inicializar templates em variáveis globais ou usar `template.ParseGlob` no início da aplicação.
* **Carregamento On-Demand:** Todo Handler deve realizar o parsing do template diretamente do disco no momento da requisição usando `template.ParseFiles("templates/nome_simple.html")`.
* **Nomenclatura:** Utilizar exclusivamente arquivos com o sufixo `_simple.html`, que são documentos HTML completos e independentes.

---

## 2. Padrão HTMX e Reatividade

O frontend do Digna deve ser dinâmico sem a complexidade de frameworks pesados de JavaScript.

* **Trocas Parciais:** Utilizar os atributos `hx-get`, `hx-post`, `hx-target` e `hx-swap` para atualizar apenas fragmentos da página (ex: atualizar o carrinho no PDV ou o saldo no caixa).
* **Resiliência Offline:** Manter a interface funcional para operações básicas sem internet, aproveitando a estrutura PWA (Manifest e Service Worker).
* **Feedback ao Usuário:** Toda ação via HTMX deve prever indicadores de carregamento ou mensagens de sucesso/erro claras (ex: "Venda Registrada!").

---

## 3. Identidade Visual "Soberania e Suor"

O agente deve garantir que a estética do sistema reflita os valores da economia solidária.

* **Paleta de Cores Oficial:**
* **Azul Soberania (#2A5CAA):** Predominante em headers e botões de ação principal.
* **Verde Suor (#4A7F3E):** Usado para indicadores de trabalho e sucesso.
* **Laranja Energia (#F57F17):** Para alertas e destaques pedagógicos.


* **Componentes:** O Logotipo Digna deve estar visível e centralizado no header de todas as páginas.
* **Linguagem Popular:** Substituir termos técnicos por descrições de ação (ex: usar "Planejamento do Mês" em vez de "Budget").

---

## 4. Validação de Templates e Scripts

* **Segurança de Tipos no Template:** Evitar lógica complexa dentro do HTML; usar funções auxiliares como `formatCurrency` ou `divide` passadas via `FuncMap`.
* **Verificação de IDs:** Antes de fechar um template, validar se os IDs de elementos (ex: `stock_item_id`) coincidem com o que o backend espera no `hx-vals`.

---

## 📈 Plano de Ação: Estado Atual

| Fase | Status | Atividade |
| --- | --- | --- |
| **0. Roadmap** | [x] | Mapeamento de todas as Skills necessárias. |
| **1. Backend** | [x] | Detalhamento da **SKILL_BACKEND_GO.md** concluído. |
| **2. Frontend** | [x] | Detalhamento da **SKILL_FRONTEND_HTMX.md** concluído. |
| **3. Negócio** | [ ] | Detalhamento da Skill de Tecnologia Social e ITG 2002. |
| **4. Infra** | [ ] | Detalhamento da Skill de Soberania e Lifecycle. |
| **5. Ponte** | [ ] | Detalhamento da Skill de Auditoria e SPED. |

> **Estado Atual do Plano:** > [x] SKILL_FRONTEND_HTMX.md concluída.
> [ ] Próximo: **SKILL_NEGOCIO_SOCIAL.md** (A alma pedagógica e contabilidade invisível).

Com estas duas habilidades técnicas prontas, o seu agente já tem o rigor para codar sem errar o "como". **Podemos agora avançar para a `SKILL_NEGOCIO_SOCIAL.md`, que dará ao agente a sensibilidade para tratar o suor do trabalhador e o rateio de sobras com a dignidade que o projeto exige?**