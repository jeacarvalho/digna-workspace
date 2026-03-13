###### title: Stakeholders e Riscos
status: implemented version: 1.2 last_updated: 2026-03-13
##### Stakeholders e Riscos - Digna

--------------------------------------------------------------------------------

###### 1. Stakeholder Map
###### 1.1 Stakeholders Primários
| Stakeholder | Role |
| ------ | ------ |
| Cooperativas | Usuários principais da operação e autogestão |
| Associações | Gestão coletiva e deliberações em assembleia |
| Grupos informais ("Sonhos") | Entrada no sistema (Estágio DREAM) |
| Contadores Sociais [NOVO] | Auditores e facilitadores da conformidade via Painel Multi-tenant |
| ITCPs e Incubadoras | Apoiadores metodológicos e pedagógicos |

###### 1.2 Stakeholders Institucionais
| Stakeholder | Interest |
| ------ | ------ |
| Ministério do Trabalho | Política pública e CADSOL |
| Senaes | Fortalecimento da economia solidária (SINAES) |
| Serpro | Infraestrutura tecnológica de nuvem soberana |
| CFC / CRCs [NOVO] | Conformidade com a norma ITG 2002 e ampliação da contabilidade formal |
| Receita Federal (RFB) [NOVO] | Conformidade tributária (EFD-Reinf, ECF) e respeito à imunidade/não-incidência do Ato Cooperativo |
| MAPA / Vigilância Sanitária [NOVO] | Adequação sanitária de pequenas agroindústrias via Memorial Técnico Sanitário (MTSE) |
| DREI / Juntas Comerciais [NOVO] | Validade de registros, anonimização de votos (IN 79/2020) e integração de Assinaturas Eletrônicas Gov.br |

###### 1.3 Stakeholders de Suporte
| Stakeholder | Role |
| ------ | ------ |
| Universidades | Pesquisa, extensão e Mutirões de "Fechamento Anual Solidário" |
| ONGs | Apoio territorial e articulação local |
| Incubadoras (ITCPs) | Formação cooperativa, pedagógica e validação de campo |

###### 1.4 Stakeholders de Governança
| Stakeholder | Role |
| ------ | ------ |
| Fundação Providentia | Governança, neutralidade e curadoria do projeto |
| PMC | Decisões técnicas e roadmap de arquitetura |
| Comunidade dev | Evolução do software (The Apache Way) |


--------------------------------------------------------------------------------

###### 2. Risk Register
###### 2.1 Riscos de Projeto
| Risco | Probabilidade | Impacto | Mitigação |
| ------ | ------ | ------ | ------ |
| Perda de dados locais | Medium | High | Backup automatizado (Litestream/Sync) e educação digital |
| Uso incorreto do sistema | Medium | Medium | UX simplificada, pedagógica e testada com o usuário final |
| Bugs contábeis | Low | Critical | Testes unitários contábeis rigorosos (Validação Soma Zero em int64) |
| Complexidade institucional | Medium | High | Documentação clara e transição gradual (Não forçar formalização) |
| Resistência da Classe Contábil [NOVO] | Medium | High | Inserção do CFC/CRCs como aliados e criação do Painel Multi-tenant que poupa tempo do contador |
| Dependência tecnológica | Low | High | Arquitetura Local-First e código Open Source |
| Autuação Fiscal e Bitributação [NOVO] | Medium | Critical | Módulo de *Tax Compliance* operando a EFD-Reinf no backend e assegurando exclusão de Atos Cooperativos no Bloco M da ECF |
| Interdição Sanitária (SIF) [NOVO] | High | High | Gerador automatizado do MTSE (Portaria 393/2021) para viabilizar a licença de comercialização de produtos de origem animal |
| Invalidação de Assembleias [NOVO] | Medium | High | Integração com Assinatura Gov.br para a mesa diretora e sistema de anonimização de votos dos cooperados (Lei 14.063/2020) |

###### 2.2 Riscos Técnicos
| Risco | Probabilidade | Impacto | Mitigação |
| ------ | ------ | ------ | ------ |
| Inconsistência no ledger | Low | Critical | Motor Lume blindado e isolado da UI (Anti-Float) |
| Falhas na sincronização | Medium | Medium | Retry com backoff e arquitetura baseada em Delta tracking |
| Corrupção de banco SQLite | Low | High | PRAGMAs otimizados (WAL mode, foreign_keys) + backup |
| Desatualização do Formato SPED/Fiscal [NOVO] | High | Medium | Módulo de exportação isolado na arquitetura e feedback contínuo da comunidade de Contadores parceiros |
| Mudança nas APIs Estatais (CADSOL/Gov.br) [NOVO] | Medium | Medium | Módulo de integrações isolado com *Circuit Breaker*, mantendo o sistema operante no modo offline em caso de indisponibilidade governamental |

###### 2.3 Riscos de Governança
| Risco | Probabilidade | Impacto | Mitigação |
| ------ | ------ | ------ | ------ |
| Captura corporativa/institucional | Low | High | Modelo Apache, governança via Fundação Providentia e licença livre |
| Fragmentação da comunidade | Medium | Medium | Comunicação ativa, RFCs transparentes e meritocracia técnica |

--------------------------------------------------------------------------------

###### 3. Matriz de Responsabilidade (RACI Adaptada)
| Atividade | Fundação | PMC | Comunidade | Contadores Parceiros [NOVO] |
| ------ | ------ | ------ | ------ | ------ |
| Roadmap estratégico | Decisão | Input | Input | Input (Compliance) |
| Decisões técnicas | Veto | Decisão | Input | - |
| Implementação Core | - | Review | Execução | - |
| Documentação | - | Aprovação | Contribuição | - |
| Conformidade ITG 2002 | Validação | Implementação | - | Auditoria Prática |
| Suporte e Educação | Oversight | - | Execução Local | Orientação Fiscal |
| Conformidade Estatal (MAPA/RFB/MTE) [NOVO] | Validação | Implementação | Execução | Auditoria Fim-a-Fim |


--------------------------------------------------------------------------------
