Tipo: Infraestrutura (Production Deploy)
Módulo: root (cmd/digna) e infraestrutura
Objetivo: Preparar o sistema para deploy em produção através de containerização com Docker e externalização de configurações via variáveis de ambiente.
Decisões: A aplicação deve rodar em um container Docker otimizado para Go. As configurações (porta, diretório de dados, nível de log) devem sair do código (hardcoded) e passar a usar variáveis de ambiente (.env). O diretório dos bancos de dados SQLite (`data/entities`) deve ser configurado para montagem em volume externo, preservando a Soberania do Dado.

### 📝 Descrição da Tarefa: Preparação para Produção (Docker e .env)
*   **Requisito Funcional (RF):** RNF-01 (Isolamento por Tenant / Soberania) e RNF-09 (Manutenibilidade).
*   **Sprint Relacionada:** Sprint 17 (Production Deploy).

Para esta tarefa, você deve carregar e seguir estritamente as instruções das seguintes skills em docs/skills/:
1. [developing-digna-backend]
2. [managing-sovereign-data]

*   **Anti-Float:** Se envolver cálculos de valor ou tempo, use estritamente int64. Proibido float.
*   **Cache-Proof:** Se houver interface, o template deve ser _simple.html carregado via ParseFiles no Handler.
*   **Soberania:** Garanta que a operação respeite o isolamento do arquivo .db do tenant atual.

---
**🎯 Objetivo da Tarefa**
Preparar o Digna para implantação em servidores reais (Marco 05 - Production Deploy). O sistema atualmente depende de execuções locais e configurações fixas. Precisamos criar um `Dockerfile` multi-stage otimizado para gerar um binário leve, um arquivo `docker-compose.yml` para orquestração fácil e modificar a inicialização (`main.go`) para ler as configurações primárias a partir de variáveis de ambiente.

**📁 Estrutura de Output Esperada**
* `Dockerfile` (na raiz do projeto)
* `docker-compose.yml` (na raiz do projeto)
* `.env.example` (na raiz do projeto)
* `cmd/digna/main.go` (atualizado para ler env vars)
* `pkg/config/config.go` (novo pacote para gerenciar leitura de variáveis de ambiente)

**🛠️ Tarefas de Implementação**
1. **Sistema de Configuração:** Criar o pacote `config` para carregar variáveis de ambiente OBRIGATÓRIAS (ex: `DIGNA_PORT`, `DIGNA_DATA_DIR`, `DIGNA_LOG_LEVEL`). Se `DIGNA_DATA_DIR` não for informado, fazer fallback para o padrão atual (`./data/entities`).
2. **Atualização do Ponto de Entrada:** Modificar a inicialização do servidor no `ui_web` e/ou `main.go` para consumir essas variáveis, garantindo que o diretório base para os arquivos `.sqlite` seja dinâmico.
3. **Dockerfile Multi-stage:** Criar um Dockerfile para Go 1.22+. 
   - Estágio de build: Baixar dependências, rodar build estático com CGO_ENABLED=1 (necessário para o mattn/go-sqlite3).
   - Estágio final: Usar uma imagem base leve (ex: debian-slim ou alpine), copiar o binário e os arquivos de template/estáticos.
4. **Docker Compose:** Criar um `docker-compose.yml` que suba a aplicação mapeando a pasta `./data` local para o volume `/var/lib/digna/data` no container, assegurando que os bancos de dados não sejam perdidos ao recriar o container (Soberania do Dado).

**✅ Critérios de Aceite (Definition of Done)**
- [ ] O sistema compila com sucesso via comando `docker build`.
- [ ] O comando `docker-compose up` sobe o servidor HTTP na porta definida no `.env`.
- [ ] A criação de uma nova entidade ou registro de venda cria o arquivo `.sqlite` persistido corretamente no volume montado fora do container.
- [ ] Os templates `_simple.html` são encontrados e renderizados corretamente por dentro do container.

---
1. Código fonte seguindo Clean Architecture (Domain -> Service -> Handler).
2. Testes unitários com TDD provando a lógica.
3. Atualização sugerida para o próximo Session Log.

Pode iniciar a análise e propor o plano de implementação?
```