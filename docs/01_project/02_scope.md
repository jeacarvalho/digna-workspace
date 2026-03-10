
#### title: Escopo do Produto
status: implemented
version: 1.1
last_updated: 2026-03-08

### Escopo do Produto - Digna

--------------------------------------------------------------------------------

#### 1. Tipo de Produto
Infraestrutura contábil, institucional e pedagógica para economia solidária, atuando também como **ponte tecnológica** entre os empreendimentos informais/formais e os profissionais de contabilidade (Contadores Sociais).

--------------------------------------------------------------------------------

#### 2. Capacidades Principais

##### 2.1 PDV Operacional
Registro simplificado de vendas e compras sem jargões contábeis.
**Funcionalidades:**
*   Keyboard numérico para entrada de valores
*   Seleção de produtos e visualização pedagógica de preço
*   Geração automática de partidas dobradas invisíveis

##### 2.2 Ledger Contábil (Motor Lume)
Motor de partidas dobradas automático.
**Funcionalidades:**
*   Validação rigorosa de soma zero
*   Histórico de lançamentos e imutabilidade
*   Balancete e demonstrativos gerenciais

##### 2.3 Registro de Trabalho
Captura de horas de trabalho cooperativo (Primazia do Trabalho / ITG 2002).
**Funcionalidades:**
*   Cronômetro em tempo real
*   Registro manual de minutos
*   Conversão do tempo em capital social de trabalho (int64)

##### 2.4 Documentação Institucional (Legal Facade)
Geração automática de documentos respeitando o ritmo do grupo.
**Funcionalidades:**
*   Atas de Assembleia (Markdown com Hash SHA256)
*   Relatórios de impacto social e prestação de contas visual
*   Dossiês de formalização (CADSOL/SINAES)

##### 2.5 Sincronização Offline-First
Operação resiliente mesmo sem internet em áreas rurais ou periféricas.
**Funcionalidades:**
*   Delta tracking local (banco SQLite isolado)
*   Sync apenas com dados agregados (preservação de privacidade)
*   Marketplace B2B

##### 2.6 Painel do Contador Social (Multi-tenant) [NOVO]
Interface de auditoria e conformidade legal focada em ganho de escala para contadores parceiros.
**Funcionalidades:**
*   Visão consolidada e alertas de múltiplos empreendimentos (Multi-tenant)
*   Auditoria de conformidade com a norma ITG 2002 (CFC) e fundos obrigatórios
*   Exportação de Lotes Fiscais/SPED para softwares contábeis externos

--------------------------------------------------------------------------------

#### 3. Capacidades Excluídas

Para manter a simplicidade e a essência da Tecnologia Social, o Digna **não** é:
*   **ERP corporativo complexo:** Não visa o mercado de capitais tradicional.
*   **Folha de pagamento empresarial (RH):** Não realiza cálculos trabalhistas da CLT.
*   **Contabilidade fiscal complexa:** O sistema **não calcula impostos (Guia do Simples, ICMS, etc)** mensalmente; ele foca no gerencial e exporta os dados íntegros para o Contador parceiro processar em seu próprio sistema.
*   **Gestão financeira bancária:** Não é um gerenciador de internet banking.
*   **Emissor de NF-e/DF-e:** (Fora do Core Lume, dependendo de integrações de terceiros futuras).

--------------------------------------------------------------------------------

#### 4. Limites do Sistema (Matriz de Escopo)

| Função | Status | Justificativa Técnica e Social |
| ------ | ------ | ------------------------------ |
| Contabilidade Gerencial | ✅ No Escopo | Core do sistema (Soberania Financeira e PDV) |
| ITG 2002 (CFC) | ✅ No Escopo | Requisito legal de valorização do trabalho |
| Painel Multi-tenant (Contador) | ✅ No Escopo (Novo)| Permite que contadores parceiros deem escala ao atendimento |
| Exportação Fiscal (SPED) | ✅ No Escopo (Novo)| A ponte que acaba com a "digitação de notas" pelo contador |
| Cálculo Direto de Impostos | ❌ Fora do Escopo | Responsabilidade do Contador; evita sobrecarga burocrática no sistema |
```
