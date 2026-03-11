Tipo: Feature (Formalização Institucional e Dossiê CADSOL)
Módulo: modules/legal_facade, modules/ui_web
Objetivo: Implementar a geração e exportação automática do Dossiê CADSOL e Atas de Assembleia em Markdown com Hash SHA256.
Decisões: A geração de documentos não deve forçar a formalização (respeito ao tempo do grupo). O sistema só deve permitir a exportação oficial do Dossiê para entidades que cumpram os critérios de formalização (mínimo de 3 decisões registradas no banco de dados isolado). O formato de saída será Markdown embutido com o Hash SHA256 para comprovar integridade perante o Estado.

### 📝 Descrição da Tarefa: Geração do Dossiê CADSOL e Exportação de Atas
*   **Requisito Funcional (RF):** RF-04 (Dossiê de Formalização - CADSOL).
*   **Sprint Relacionada:** Sprint 18 (Fase 2 - O Trilho da Formalização).

Para esta tarefa, você deve carregar e seguir estritamente as instruções das seguintes skills em docs/skills/:
1. [developing-digna-backend]
2. [rendering-digna-frontend]
3. [managing-sovereign-data]

*   **Anti-Float:** Se envolver cálculos de valor ou tempo, use estritamente int64. Proibido float.
*   **Cache-Proof:** Se houver interface, o template deve ser _simple.html carregado via ParseFiles no Handler.
*   **Soberania:** Garanta que a operação respeite o isolamento do arquivo .db do tenant atual.

---
**🎯 Objetivo da Tarefa**
Implementar o motor de geração de documentos dentro do módulo `legal_facade`, completando o Requisito RF-04. A funcionalidade permitirá que cooperativas e grupos extraiam seu histórico de deliberações em Assembleia (Decisions) na forma de um Dossiê em Markdown. O documento gerado DEVE conter um Hash criptográfico SHA256 para atestar sua imutabilidade técnica, permitindo comprovação de autogestão para o Ministério do Trabalho e Emprego (SINAES/CADSOL).

**📁 Estrutura de Output Esperada**
* `modules/legal_facade/internal/document/generator.go` (Lógica de compilação de atas e Hash)
* `modules/ui_web/internal/handler/legal_handler.go` (Controlador HTTP para a UI)
* `modules/ui_web/templates/legal_dossier_simple.html` (Interface visual HTMX + Tailwind)

**🛠️ Tarefas de Implementação**
1. **Gerador de Documentos (`legal_facade`):** Desenvolver o serviço que busca as decisões cadastradas no `DecisionRepository` isolado da entidade e formata um documento em texto puro (Markdown). 
2. **Garantia de Integridade (Hash SHA256):** O arquivo `.md` resultante deve computar um Hash SHA256 de todo o seu conteúdo (decisões, datas, nomes da entidade) e embutir este Hash no final do texto para validação de auditoria pública.
3. **Regra de Transição (Autogestão Gradual):** A lógica no backend (Handler/Service) deve validar se a entidade já atingiu os critérios de formalização. Se houver menos de 3 decisões de Assembleia registradas, o botão de "Gerar Dossiê Oficial CADSOL" na view deve estar inativo ou retornar uma mensagem pedagógica via HTMX orientando o grupo a registrar mais decisões antes de se formalizar.
4. **Interface HTMX (Cache-Proof):** Criar a rota `GET /legal/dossier` que carrega a página `legal_dossier_simple.html` (herdando o design "Soberania e Suor"). A tela terá um botão grande acessível para acionar o download do `.md`.

**✅ Critérios de Aceite (Definition of Done)**
- [ ] O usuário consegue baixar com sucesso um arquivo `.md` formatado contendo as decisões e informações sociais da entidade, caso cumpra os requisitos mínimos de formalização.
- [ ] O documento gerado possui um carimbo verificável de Hash SHA256.
- [ ] Entidades no estágio incipiente (`DREAM`) com histórico vazio ou < 3 decisões são pedagogicamente informadas na tela de que precisam de mais assembleias, preservando a regra de não forçar a burocracia governamental.
- [ ] A Soberania do Dado é totalmente preservada, acessando apenas os dados do banco `.sqlite` instanciado para o usuário atual da sessão.

---
1. Código fonte seguindo Clean Architecture (Domain -> Service -> Handler).
2. Testes unitários com TDD provando a lógica.
3. Atualização sugerida para o próximo Session Log.

Pode iniciar a análise e propor o plano de implementação?
```