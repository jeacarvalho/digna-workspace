### 3. `applying-solidarity-logic/SKILL.md`

**Foco:** Tradução cultural, ITG 2002 e pedagogia social.

```yaml
name: applying-solidarity-logic
description: Use esta habilidade ao implementar ou modificar regras de negócio, interfaces ou algoritmos que envolvam a Economia Solidária. Ela garante a conformidade com a ITG 2002, a valorização do trabalho e a pedagogia social do sistema Digna.

```

# ✊ Instruções de Negócio Social e Tecnologia Social

Esta habilidade orienta o agente a agir como um facilitador da autogestão e da soberania financeira, protegendo a alma do projeto contra lógicas extrativistas ou puramente capitalistas.

## 1. Tradução Cultural (Contabilidade Invisível)

O agente deve atuar como um tradutor entre a lida diária do trabalhador e o rigor fiscal.

* **Ação vs. Lançamento:** Toda funcionalidade deve ser projetada a partir da ação coloquial (ex: "Vendi na feira") e convertida silenciosamente em partidas dobradas no backend.
* **Barreira de Jargão:** É terminantemente proibido o uso de termos técnicos como "Ativo", "Passivo" ou "Exercício Fiscal" na interface do produtor. Use "Dinheiro na Gaveta", "Contas a Pagar" e "Fechamento do Mês".

## 2. Primazia do Trabalho (O Motor de Suor)

No Digna, o tempo é a unidade de medida da dignidade e do capital social.

* **Valoração do Tempo:** Todo registro de trabalho deve ser capturado estritamente em minutos (`int64`), servindo de lastro para o rateio justo.
* **Capital Social de Trabalho:** Garantir que o suor (horas trabalhadas) tenha peso fundamental no cálculo das sobras, refletindo o princípio de que o trabalho vale mais que o capital investido.

## 3. Algoritmos de Justiça e Reservas (ITG 2002)

Implementar o rigor do CFC (Conselho Federal de Contabilidade) de forma automatizada e transparente.

* **Reservas Obrigatórias:** Antes de qualquer rateio de sobras, o sistema deve bloquear mandatoriamente **10% para Reserva Legal** e **5% para o Fundo FATES** (Assistência Técnica e Educacional).
* **Transparência Algorítmica:** Ao realizar cálculos de distribuição, o agente deve propor visualizações didáticas (gráficos ou tabelas simples) que permitam a compreensão e aprovação democrática em Assembleia Geral.

## 4. Design Participativo e Pedagógico

O software deve atuar como uma ferramenta de ensino e emancipação digital.

* **Apoio na Formação de Preço:** Propor interfaces que auxiliem visualmente o trabalhador a entender que o preço do produto deve cobrir o custo do insumo mais o valor justo da hora trabalhada.
* **Transição Respeitosa (DREAM -> FORMALIZED):** Respeitar o tempo social do coletivo. Não forçar a formalização (CNPJ) antes que o grupo demonstre maturidade política, validada pelo registro de ao menos 3 decisões em Assembleia.

## 📋 Checklist de Validação para o Agente

* [ ] Esta funcionalidade utiliza linguagem popular em vez de jargão contábil?
* [ ] O tempo de trabalho está sendo tratado como Capital Social em `int64`?
* [ ] As reservas de 10% (Legal) e 5% (FATES) foram protegidas antes do rateio?
* [ ] A interface ajuda o usuário a entender o impacto financeiro da sua ação?

---
