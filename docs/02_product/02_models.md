###### title: Modelos de Domínio e Dados
status: implemented version: 1.4 last_updated: 2026-03-13
##### Modelos - Projeto Digna
**Projeto:**  Sistema de Gestão Contábil e Pedagógica para Economia Solidária

--------------------------------------------------------------------------------

###### 1. Domain Model (Modelo de Domínio)
O domínio do Digna reflete os princípios da autogestão e da contabilidade invisível, priorizando as relações humanas sobre o capital financeiro, mas agora atuando também como  **ponte institucional**  para a conformidade legal.

###### 1.1 Entidades Principais

###### Enterprise (Empreendimento de Economia Solidária - EES)
Representa o coletivo produtivo. Pode transitar gradualmente por três estados, respeitando o tempo político do grupo:
*   **DREAM (Sonho):**  Grupo informal, focado na união produtiva inicial.
*   **INCUBATED (Incubado):**  Em processo de estruturação, recebendo apoio pedagógico (ITCPs, ONGs).
*   **FORMALIZED (Formalizado):**  Cooperativa ou Associação com CNPJ e estatuto base (Pronto para CADSOL e obrigações fiscais).

###### WorkLog (Registro ITG 2002)
Registro de trabalho cooperativo. Converte o suor em Capital Social de Trabalho (mensurado em minutos).

###### Decision (Decisão Democrática)
Decisão coletiva tomada e registrada em Assembleia. Base para a geração das Atas em Markdown (CADSOL).

###### LegalDocument (Dossiê CADSOL e Atas)
*   **Camada Cotidiana:** O livro de atas oficial do grupo. Agrupamento das decisões exportado para provar que a gestão é democrática. Com a nova atualização, o presidente assina pelo celular com a senha oficial do Gov.br, sem precisar pisar no cartório, e o voto secreto de cada um é garantido pelo sistema.
*   **Camada Técnica:** Agrupamento das `Decisions` exportado em formato Markdown/PDF para fins legais, comprovando a gestão democrática perante o Estado. O sistema agora supera o mero uso de Hash SHA256 integrando Assinatura Eletrônica Avançada/Qualificada (ICP-Brasil/Gov.br) para a mesa diretora e assegurando a anonimização em votações (IN DREI nº 79/2020 e Lei nº 14.063/2020).

###### Fund (Fundos Obrigatórios)
Reservas estatutárias e legais blindadas pelo sistema (Ex: Reserva Legal e FATES).

###### EnterpriseAccountant (Vínculo Contábil Temporal) [NOVO]
*   **Camada Cotidiana:** A "chave temporária" que a cooperativa empresta para o contador social arrumar a casa, podendo ser recolhida a qualquer hora.
*   **Camada Técnica:** Relacionamento N:N que estabelece a delegação de acesso temporal entre um `Enterprise` e um Contador parceiro. Atua como um "Filtro de Vigência", liberando acesso ao painel unicamente em modo `Read-Only` para os períodos fiscais previamente autorizados pela autogestão.

###### TaxCompliance & ReinfEvent [NOVO]
*   **Camada Cotidiana:** O mensageiro silencioso do aplicativo. Ele avisa os robôs da Receita Federal sobre pagamentos de frete e serviços para a cooperativa não ser multada, e garante que o dinheiro do "Ato Cooperativo" não sofra a mesma mordida de imposto das grandes empresas.
*   **Camada Técnica:** Motor de geração de XMLs de retenção (assinados via certificado A1/A3) e envio síncrono/assíncrono para Web Services da EFD-Reinf (ex: evento de fechamento R-2099) para alimentar a DCTFWeb. Detém a inteligência contábil de expurgar as receitas de Atos Cooperativos na composição do Bloco M da ECF (e-Lalur/e-Lacs), blindando a isenção prevista na Lei 5.764/71 e LC 214/2025.

###### SanitaryDossier (MTSE) [NOVO]
*   **Camada Cotidiana:** O gerador automatizado do "caderno da fábrica", ajudando a pequena agroindústria (queijaria, casa de mel, pescados) a desenhar sua planta e procedimentos de higiene para conquistar o selo de venda oficial sem precisar de consultorias caras.
*   **Camada Técnica:** Estrutura de dados que parametriza fluxogramas de maquinário, capacidade diária e potabilidade de água, exportando automaticamente o Memorial Técnico Sanitário de Estabelecimento (MTSE) em conformidade estrita com o RIISPOA (Portaria MAPA nº 393/2021) para peticionamento no sistema SEI.

--------------------------------------------------------------------------------

###### 2. Data Model (Schema v1)
O banco de dados é instanciado fisicamente de forma isolada por Enterprise (Soberania do Dado local).  *Nota Arquitetural:*  O Painel do Contador Social  **não possui**  banco de transações próprio; ele atua apenas consumindo dados em modo de leitura (Read-Only) dos micro-databases autorizados. E como um contador social pode estar cuidando de várias entidades, as funcionalidades a ele relacionadas devem estar preparadas para pesquisar em diversos db sqlite utilizando o filtro do `EnterpriseAccountant`.

--------------------------------------------------------------------------------

###### 3. Algoritmos de Negócio e Compliance

###### 3.4 Algoritmo de Formalização Gradual e Integração CADSOL
**Objetivo:**  Avaliar a maturidade institucional para permitir a transição DREAM -> FORMALIZED.
**Critérios Automatizados:**
*  Mínimo de 3 registros de Decision (Assembleias realizadas provando autogestão).
*  Mínimo de 1 membro ativo com histórico de WorkLog.
*  Criação automática do Dossiê Hash SHA256 com validação de Assinaturas ICP-Brasil/Gov.br.
*  **Integração Automática (MTE):** Consumo via integração (módulo `integrations`) das futuras APIs do SINAES/CADSOL (Decreto nº 12.784/2025), substituindo envios em papel pela interoperabilidade digital assim que o grupo atingir o estágio FORMALIZED.

###### 3.5 Algoritmo de Tradução Fiscal (Ponte do Contador) [NOVO]
**Objetivo:**  Converter a contabilidade social gerada invisivelmente pelo produtor em linguagem de conformidade estatal (SPED/Lotes Fiscais/Reinf).  **Processo:**
1. Painel do Contador solicita dados de um Period fechado via filtro temporal `EnterpriseAccountant`.
2. O algoritmo compila todas as entries de soma zero (respeitando a ITG 2002).
3. Mapeia as contas locais amigáveis ("Gaveta") para o Plano de Contas Referencial da Receita Federal.
4. **Isolamento do Ato Cooperativo:** O sistema segrega as transações com terceiros não cooperados (destinando sobras ao FATES) e estrutura a base de cálculo tributável. As receitas do trabalho intercooperativo são expurgadas automaticamente para preenchimento na Parte A do e-Lalur/e-Lacs (Bloco M da ECF).
5. Gera o pacote CSV/SPED/XMLs da EFD-Reinf e salva o evento de exportação na tabela fiscal_exports.

--------------------------------------------------------------------------------

###### 4. Seed Data (Carga Inicial Padrão)
Toda nova base SQLite de um EES nasce com este plano de contas enxuto e adaptado:
| ID | Código | Nome Amigável | Natureza Contábil (Invisível) | Mapeamento Fiscal [NOVO] |
| ------ | ------ | ------ | ------ | ------ |
| 1 | 1.1.01 | Gaveta / Caixa | ASSET (Ativo) | Disponibilidades (Ativo) |
| 2 | 3.1.01 | Nossas Vendas | REVENUE (Receita) | Receita Bruta |
| 3 | 1.1.02 | Banco / Conta | ASSET (Ativo) | Contas Bancárias |
| 4 | 2.1.01 | Quem Fornece | LIABILITY (Passivo) | Fornecedores a Pagar |
| 5 | 3.2.01 | Fundo FATES | EQUITY (Patrimônio Líquido) | Reservas Estatutárias |
