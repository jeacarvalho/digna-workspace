Tipo: Infraestrutura (Production Deploy)
Módulo: root (cmd/digna) e infraestrutura
Objetivo: Preparar o sistema para deploy em produção através de containerização com Docker e externalização de configurações via variáveis de ambiente.
Decisões: A aplicação deve rodar em um container Docker otimizado para Go. As configurações (porta, diretório de dados, nível de log) devem sair do código (hardcoded) e passar a usar variáveis de ambiente (.env). O diretório dos bancos de dados SQLite (`data/entities`) deve ser configurado para montagem em volume externo, preservando a Soberania do Dado.

### 📝 Descrição da Tarefa: Preparação para Produção (Docker e .env)
O projeto já foi configurado para subir via docker em uma VPS. Mas como estou pagando por uso, sempre desligo a máquina e ao ligar preciso fazer tudo de novo.
Assim, monte um script shell que faça todos os passos automaticamente ao entrar. Instale docker-compose se não existir, git clone do repositório remoto e todo o build e subida do conteiner.
Pode iniciar a análise e propor o plano de implementação?
```