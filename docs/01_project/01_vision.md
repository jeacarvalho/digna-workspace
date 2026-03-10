#### title: Visão Estratégica
status: implemented
version: 1.4
last_updated: 2026-03-08

### Visão Estratégica - Digna

--------------------------------------------------------------------------------

#### 1. Introdução
O **Digna** é um ecossistema de soberania financeira desenhado para a Economia Solidária (EES). Ele não é um "ERP" (Enterprise Resource Planning) imposto de cima para baixo com lógica extrativista; é uma **Tecnologia Social e um Protocolo de Emancipação**. 

Seu propósito é transformar a contabilidade — historicamente vista como um fardo burocrático — em um subproduto invisível da operação diária, atuando simultaneamente como uma ferramenta pedagógica que transforma a atividade produtiva em cidadania digital.

Além de proteger e educar o trabalhador, o Digna atua como uma **Ponte Tecnológica Institucional**, conectando a realidade informal da base produtiva às exigências da Receita Federal e do Conselho Federal de Contabilidade (CFC), através da viabilização da "Contabilidade Popular" gerida por contadores parceiros.

##### 1.1 Declaração de Posicionamento Estratégico

    [ PARA ]        Grupos informais ("Sonhos"), Cooperativas e Associações de EES, 
                    muitas vezes afetados pela baixa literacia digital.
                    
    [ QUE ]         Enfrentam exclusão financeira, dificuldade na formação de preço, 
                    burocracia do CADSOL e falta de ferramentas de autogestão.
                    
    [ O DIGNA É ]   Uma infraestrutura contábil "Local-First" e uma Tecnologia Social 
                    pedagógica.
                    
    [ QUE GERA ]    Contabilidade invisível, valoração do tempo de trabalho (ITG 2002),
                    educação gerencial e soberania de dados isolados.
                    
    [ DIFERENTE DE] ERPs tradicionais comerciais ou planilhas vulneráveis, que usam 
                    lógica de mercado de capitais e não respeitam os laços de 
                    solidariedade e o tempo de amadurecimento do grupo.
                    
    [ ENVOLVENDO ]  Contadores Sociais e parceiros (CRCs) que, através de um painel 
                    multi-tenant, conseguem dar escala à contabilidade popular, 
                    deixando de ser "digitadores de notas" para serem consultores.

--------------------------------------------------------------------------------

#### 2. Pilares de Design (As Leis Sociotécnicas do Sistema)

##### Pilar 1: Soberania do Dado e Poder de Saída (Exit Power)
O dado não pertence à "nuvem" de uma corporação, pertence à entidade produtiva. O dado reside em um arquivo SQLite isolado fisicamente por empreendimento. O usuário detém o poder absoluto de auditar, copiar ou sair do sistema levando toda a sua história com ele.

##### Pilar 2: Contabilidade Invisível e Tradução Cultural
A interface humana (Frontend) foca na ação coloquial (vender, comprar, trabalhar) e atua como uma barreira contra jargões contábeis. O débito e o crédito (Partidas Dobradas) são subprodutos gerados automaticamente pelo Motor Lume no backend.

##### Pilar 3: Primazia do Trabalho (ITG 2002)
O sistema inverte a lógica capitalista: o suor (tempo/horas trabalhadas) vale tanto ou mais que o capital investido (R$). O tempo registrado (em minutos/int64) constitui o Capital Social de Trabalho e é a base para o rateio justo de sobras.

##### Pilar 4: Transição Institucional Gradual (Sem Burocracia Forçada)
O Digna respeita o tempo social do grupo. Ele atua como um facilitador da conformidade (CADSOL/Sinaes), gerando atas e relatórios, mas não impõe a formalização precoce a grupos informais ("Sonhos") que ainda estão construindo sua confiança e coesão política.

##### Pilar 5: Ferramenta Pedagógica e Design Participativo
O software ensina enquanto é operado. Ele auxilia visualmente o trabalhador na formação correta do seu preço (custo de insumos + hora trabalhada). Todo o seu desenvolvimento deve ser validado *com* os trabalhadores e Incubadoras (ITCPs).

##### Pilar 6: Aliança Contábil e Escala Fiscal (Novo)
O sistema não exclui o contador; ele o eleva. O Digna é a materialização tecnológica da norma ITG 2002 do Conselho Federal de Contabilidade (CFC). Ao automatizar a digitação, o sistema permite que "Contadores Sociais" voluntários ou de baixo custo atendam dezenas de cooperativas simultaneamente de forma viável, transformando conformidade legal em inclusão.

--------------------------------------------------------------------------------

#### 3. Princípios Centrais de Operação

O trabalhador da Economia Solidária **não faz contabilidade tradicional**, ele pratica a autogestão. O Contador parceiro **não digita notas**, ele audita e orienta.

    AÇÃO HUMANA (Trabalhadores e Contadores)       AÇÃO DO DIGNA (Motor Lume)
    +----------------------------------------+     +--------------------------------------------+
    | 1. Vende seu produto na feira          | ==> | Gera partida dobrada (D:Ativo / C:Receita) |
    | 2. Compra sementes e insumos           | ==> | Registra despesa e baixa no estoque        |
    | 3. Trabalha 4 horas na produção        | ==> | Valora o tempo como Capital Social         |
    | 4. Reúne-se em Assembleia              | ==> | Gera Ata em Markdown (Hash SHA256)         |
    | 5. Decide dividir os ganhos            | ==> | Calcula as Reservas (15%) e o Rateio       |
    | 6. Contador acessa Painel Multi-tenant | ==> | Valida ITG 2002 e Exporta Arquivos Fiscais |
    +----------------------------------------+     +--------------------------------------------+

--------------------------------------------------------------------------------

#### 4. Roadmap Estratégico de Longo Prazo

##### Fase 0: Demonstração e Validação Cultural (Atual)
*   **Foco:** O grupo informal passa a operar com rigor contábil, mas com interface amigável.
*   **Milestones:** Motor Lume Exato + PDV Pedagógico + Testes de Usabilidade em Campo com as EES.

##### Fase 1: Integração, O Trilho da Formalização e o Contador Social 
*   **Foco:** Oferecer os benefícios do Estado sem o peso da burocracia e aproximar a classe contábil.
*   **Milestones:** Integração Gov.br + Dossiê CADSOL automático + **Painel do Contador (Accountant Dashboard Multi-tenant) + Exportação de Lotes Fiscais (SPED)** + Criação do programa "Imposto de Renda Solidário" com Faculdades de Contabilidade/CRCs.

##### Fase 2: Finanças Solidárias e Territoriais
*   **Foco:** Gerenciar riquezas além da moeda oficial (Real R$).
*   **Milestones:** Integração tecnológica com Bancos Comunitários de Desenvolvimento (BCDs) + Moedas Sociais Locais + Estoque Substantivo (Troca de Sementes, Animais e Horas).

##### Fase 3: Intercooperação e Escala
*   **Foco:** Uma rede nacional de apoio mútuo (O 6º Princípio do Cooperativismo).
*   **Milestones:** Marketplace B2B fechado para EES + Score de Crédito Social baseado no trabalho + Integração com BNDES e políticas públicas via Serpro.
```