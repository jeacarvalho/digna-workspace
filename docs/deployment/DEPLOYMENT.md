🚀 Guia de Deploy em Produção - Projeto Digna
Última atualização: 27/03/2026
Status: ✅ PRONTO PARA PRODUÇÃO + ECOSSISTEMA DE 4 MÓDULOS

📋 Visão Geral

Este documento descreve o processo completo de deploy do projeto Digna em ambiente de produção (VPS). O sistema foi containerizado com Docker e utiliza variáveis de ambiente para configuração.

**ATUALIZAÇÃO 27/03/2026:** Este guia foi expandido para suportar o **Ecossistema Digna de 4 Módulos** (PDF v1.0), incluindo:
- Módulo 1: digna ERP (núcleo)
- Módulo 2: Motor de Indicadores (BCB/IBGE APIs)
- Módulo 3: Portal de Oportunidades (match de crédito)
- Módulo 4: Rede Digna (marketplace solidário)
- Sistema Transversal: Ajuda Educativa (RF-30)

🎯 Objetivos do Deploy

- **Automação completa** - Script único para deploy em VPS
- **Persistência de dados** - Bancos SQLite em volumes externos
- **Configuração via ambiente** - Zero hardcoding no código
- **Backup/restore** - Sistema robusto de backup dos dados
- **Manutenção simplificada** - Comandos simples para operações
- **Suporte a múltiplos módulos** - Deploy unificado do ecossistema

🏗️ Arquitetura de Deploy

VPS (Ubuntu/Debian)
├── Docker + docker-compose
├── Digna Container (Ecossistema Completo)
│   ├── digna ERP (Módulo 1)
│   ├── Motor de Indicadores (Módulo 2)
│   ├── Portal de Oportunidades (Módulo 3)
│   ├── Rede Digna (Módulo 4)
│   ├── Sistema de Ajuda (RF-30)
│   └── Configuração via .env
└── Volumes Externos
    ├── data/entities/ (SQLite databases por entidade)
    ├── data/central.db (Banco central - indicadores, programas, help_topics)
    └── backups/ (Backups automáticos)

📁 Estrutura de Arquivos de Deploy

scripts/deploy/
├── vps_deploy.sh          # Script principal de deploy
├── backup.sh              # Backup de bancos SQLite
├── restore.sh             # Restauração de backup
└── validate_deployment.sh # Validação do deploy

deploy.sh                  # Wrapper script (raiz do projeto)
docker-compose.yml         # Configuração dev/docker-compose
docker-compose.prod.yml    # Configuração produção
.env.example               # Template variáveis ambiente
Dockerfile                 # Build da aplicação

🚀 Deploy Rápido (VPS Nova)

Pré-requisitos na VPS

# 1. Docker instalado
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# 2. Git instalado
sudo apt-get update && sudo apt-get install -y git curl

# 3. Recursos mínimos recomendados
# - 2GB RAM (mínimo), 4GB RAM (recomendado)
# - 20GB SSD (mínimo), 50GB SSD (recomendado)
# - Ubuntu 20.04+ ou Debian 11+

Deploy Automático (Recomendado)

# 1. Baixar script de deploy
curl -O https://raw.githubusercontent.com/providentia/digna/main/deploy.sh
chmod +x deploy.sh

# 2. Executar deploy
./deploy.sh

# 3. Para atualizar posteriormente
./deploy.sh --update

Deploy Manual (Passo a Passo)

# 1. Clonar repositório
git clone https://github.com/providentia/digna.git
cd digna

# 2. Configurar ambiente
cp .env.example .env
# Editar .env conforme necessário

# 3. Build e execução
docker-compose build
docker-compose up -d

# 4. Verificar status
docker-compose ps
curl http://localhost:8090/health

⚙️ Configuração de Ambiente (.env)

Arquivo `.env` mínimo

# Server Configuration
DIGNA_PORT=8090

# Data Directory Configuration
DIGNA_DATA_DIR=/var/lib/digna/data

# Logging Configuration
DIGNA_LOG_LEVEL=info

# Docker Compose Configuration
COMPOSE_PROJECT_NAME=digna

# Ecossistema Configuration [NOVO - 27/03/2026]
DIGNA_ECOSYSTEM_ENABLED=true
DIGNA_INDICATORS_ENABLED=true
DIGNA_PORTAL_ENABLED=true
DIGNA_REDE_ENABLED=false  # Requer massa crítica

# Motor de Indicadores [NOVO - Módulo 2]
DIGNA_INDICATORS_CACHE_TTL=86400  # 24 horas em segundos
DIGNA_INDICATORS_BCB_ENABLED=true
DIGNA_INDICATORS_IBGE_ENABLED=true

# Portal de Oportunidades [NOVO - Módulo 3]
DIGNA_PORTAL_MVP_PROGRAMS=3  # Acredita, Pronampe, Niterói
DIGNA_PORTAL_MATCH_ENABLED=true

# Sistema de Ajuda Educativa [NOVO - RF-30]
DIGNA_HELP_ENABLED=true
DIGNA_HELP_CACHE_TTL=3600  # 1 hora em segundos

Variáveis de Ambiente Disponíveis

| Variável | Descrição | Padrão | Obrigatório |
|----------|-----------|--------|-------------|
| DIGNA_PORT | Porta HTTP da aplicação | 8090 | Não |
| DIGNA_DATA_DIR | Diretório para bancos SQLite | ./data/entities | Não |
| DIGNA_LOG_LEVEL | Nível de log (debug, info, warn, error) | info | Não |
| COMPOSE_PROJECT_NAME | Nome do projeto Docker | digna | Não |
| DIGNA_ECOSYSTEM_ENABLED | Habilitar ecossistema completo | true | Não |
| DIGNA_INDICATORS_ENABLED | Habilitar Motor de Indicadores | true | Não |
| DIGNA_PORTAL_ENABLED | Habilitar Portal de Oportunidades | true | Não |
| DIGNA_REDE_ENABLED | Habilitar Rede Digna | false | Não |
| DIGNA_HELP_ENABLED | Habilitar Sistema de Ajuda | true | Não |

Configuração Avançada

# Para alta disponibilidade
DIGNA_PORT=8080
DIGNA_DATA_DIR=/mnt/ssd/digna/data
DIGNA_LOG_LEVEL=warn

# Com proxy reverso (Nginx)
DIGNA_PORT=3000  # Interno no container
# Nginx redireciona porta 80/443 para 3000

# Para desenvolvimento
DIGNA_LOG_LEVEL=debug
DIGNA_HELP_CACHE_TTL=0  # Sem cache para desenvolvimento

💾 Persistência de Dados

Estrutura de Dados

/var/lib/digna/data/          # Volume Docker
├── central.db                # Banco central (indicadores, programas, help_topics)
├── entities/
│   ├── entity_abc123.db      # Banco da entidade 1
│   ├── entity_def456.db      # Banco da entidade 2
│   └── entity_ghi789.db      # Banco da entidade 3
└── backups/
    ├── digna_backup_20260327_020000.tar.gz
    └── digna_backup_20260328_020000.tar.gz

Backup Automático

# Backup manual
./scripts/deploy/backup.sh

# Backup com retenção personalizada
./scripts/deploy/backup.sh --output-dir=/backups --keep-days=30

# Agendamento via cron (backup diário às 2AM)
0 2 * * * /opt/digna/scripts/deploy/backup.sh --keep-days=7

# Backup inclui:
# - central.db (indicadores, programas, help_topics)
# - Todos os entity_*.db (dados das entidades)
# - Configurações do ecossistema

Restauração de Backup

# Listar backups disponíveis
./scripts/deploy/restore.sh

# Restaurar backup específico
./scripts/deploy/restore.sh --backup-file=/backups/digna_backup_20260327_020000.tar.gz

# Dry run (teste sem alterações)
./scripts/deploy/restore.sh --dry-run

# Restauração preserva:
# - Todos os módulos do ecossistema
# - Histórico de indicadores econômicos
# - Programas de financiamento cadastrados
# - Tópicos de ajuda educativa

🔧 Manutenção e Operações

Comandos Comuns

# Ver logs da aplicação
docker-compose logs -f

# Parar aplicação
docker-compose stop

# Reiniciar aplicação
docker-compose restart

# Rebuild e restart
docker-compose up -d --build

# Remover tudo (cuidado!)
docker-compose down -v

# Acessar shell do container
docker-compose exec digna sh

# Ver status dos módulos do ecossistema
curl http://localhost:8090/health/ecosystem

Monitoramento

# Health check geral
curl http://localhost:8090/health

# Health check por módulo [NOVO - 27/03/2026]
curl http://localhost:8090/health/erp          # Módulo 1
curl http://localhost:8090/health/indicators   # Módulo 2
curl http://localhost:8090/health/portal       # Módulo 3
curl http://localhost:8090/health/rede         # Módulo 4
curl http://localhost:8090/health/help         # RF-30

# Readiness check
curl http://localhost:8090/ready

# Status do container
docker-compose ps
docker-compose top

# Métricas do ecossistema [NOVO - 27/03/2026]
curl http://localhost:8090/metrics/ecosystem
# Retorna: entidades_ativas, indicadores_coletados, matches_realizados, etc.

Limpeza

# Limpar containers parados
docker-compose rm -f

# Limpar imagens não utilizadas
docker image prune -a

# Limpar volumes não utilizados
docker volume prune

# Limpar cache de indicadores (se necessário)
curl -X POST http://localhost:8090/admin/indicators/cache/clear

📊 Escalabilidade

Recursos do Container

# Em docker-compose.prod.yml
deploy:
  resources:
    limits:
      memory: 512M
      cpus: '1.0'
    reservations:
      memory: 256M
      cpus: '0.5'

# Para ecossistema completo [NOVO - 27/03/2026]
deploy:
  resources:
    limits:
      memory: 1024M  # Mais memória para múltiplos módulos
      cpus: '2.0'
    reservations:
      memory: 512M
      cpus: '1.0'

Otimizações para Produção

- Use `docker-compose.prod.yml` para limites de recursos
- Configure swap na VPS para picos de memória
- Use volume externo para dados (não bind mount)
- Configure log rotation para logs do Docker
- Use rede interna para segurança
- **NOVO:** Configure cache Redis para indicadores (opcional, alta escala)

🔒 Segurança

Melhores Práticas

- Não exponha a porta diretamente - Use Nginx/Apache como proxy
- Configure firewall - Apenas portas necessárias
- Use HTTPS - Certificado SSL via Let's Encrypt
- Atualize regularmente - Docker, sistema, aplicação
- Backup frequente - Dados são críticos
- **NOVO:** Rate limiting para APIs externas (BCB, IBGE)

Exemplo Nginx como Proxy

server {
    listen 80;
    server_name digna.example.com;
    
    location / {
        proxy_pass http://localhost:8090;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
    
    # Rate limiting para APIs [NOVO - 27/03/2026]
    location /api/indicators/ {
        limit_req zone=indicators burst=10 nodelay;
        proxy_pass http://localhost:8090;
    }
}

# Rate limiting config (nginx.conf)
http {
    limit_req_zone $binary_remote_addr zone=indicators:10m rate=10r/s;
}

🐛 Troubleshooting

Problemas Comuns

1. Container não inicia

# Ver logs
docker-compose logs

# Verificar portas em uso
sudo netstat -tulpn | grep :8090

# Testar build manual
docker-compose build --no-cache

2. Erro de permissão nos dados

# Corrigir permissões
sudo chown -R 1000:1000 /var/lib/digna/data

# No container
docker-compose exec digna chown -R digna:digna /var/lib/digna/data

3. Health check falha

# Verificar se aplicação responde
curl -v http://localhost:8090/health

# Verificar logs da aplicação
docker-compose logs digna

# Reiniciar container
docker-compose restart

# Verificar módulos individuais [NOVO - 27/03/2026]
curl http://localhost:8090/health/indicators  # Verificar Motor
curl http://localhost:8090/health/portal      # Verificar Portal

4. Espaço em disco

# Verificar uso
df -h /var/lib/digna

# Limpar backups antigos
find /var/backups/digna -name "*.tar.gz" -mtime +30 -delete

# Limpar logs do Docker
docker system prune -a

# Limpar cache de indicadores [NOVO - 27/03/2026]
du -sh /var/lib/digna/data/central.db
# Se > 1GB, considerar limpeza de indicadores antigos

5. APIs externas indisponíveis [NOVO - 27/03/2026]

# Verificar status das APIs
curl http://localhost:8090/admin/indicators/status

# Forçar refresh manual
curl -X POST http://localhost:8090/admin/indicators/refresh

# Verificar circuit breaker
curl http://localhost:8090/admin/indicators/circuit-breaker

📈 Monitoramento e Métricas

Métricas Importantes

- Uso de CPU/Memória - `docker stats`
- Logs da aplicação - `docker-compose logs`
- Health checks - Monitorar endpoint `/health`
- Espaço em disco - Monitorar volume de dados
- Backups - Verificar execução regular
- **NOVO:** Coleta de indicadores - Frequência e sucesso
- **NOVO:** Matches de crédito - Quantidade e taxa de sucesso
- **NOVO:** Tópicos de ajuda - Visualizações e tópicos mais acessados

Ferramentas Recomendadas

- Docker Compose - Gerenciamento de containers
- Cron - Agendamento de backups
- Logrotate - Rotação de logs
- Monitoramento - Prometheus + Grafana (opcional)
- **NOVO:** Dashboard do ecossistema - `/admin/dashboard`

🔄 Atualização da Aplicação

Fluxo de Atualização

# 1. Parar aplicação
docker-compose stop

# 2. Backup dos dados
./scripts/deploy/backup.sh

# 3. Atualizar código
git pull origin main

# 4. Rebuild
docker-compose build --no-cache

# 5. Iniciar
docker-compose up -d

# 6. Verificar
curl http://localhost:8090/health

# 7. Verificar ecossistema [NOVO - 27/03/2026]
curl http://localhost:8090/health/ecosystem

Rollback (em caso de problemas)

# 1. Parar aplicação
docker-compose stop

# 2. Restaurar backup
./scripts/deploy/restore.sh --backup-file=/backups/digna_backup_YYYYMMDD_HHMMSS.tar.gz

# 3. Voltar para versão anterior do código
git checkout <commit-hash>

# 4. Rebuild e start
docker-compose up -d --build

# 5. Verificar rollback
curl http://localhost:8090/health

📞 Suporte e Contato

Recursos

- Documentação: `docs/` no repositório
- Issues: GitHub Issues do projeto
- Scripts: `scripts/deploy/` para operações
- **NOVO:** Dashboard administrativo: `/admin`

Checklist Pós-Deploy

- [ ] Aplicação responde em `http://SEU_IP:8090`
- [ ] Health check passa (`/health`)
- [ ] Health check dos módulos passa (`/health/ecosystem`)
- [ ] Dados persistem após restart
- [ ] Backup configurado e testado
- [ ] Logs estão sendo gerados
- [ ] Monitoramento configurado
- [ ] **NOVO:** Motor de Indicadores coletando dados
- [ ] **NOVO:** Sistema de Ajuda acessível (`/help`)
- [ ] **NOVO:** Backup do central.db incluído

🎯 Conclusão

O sistema Digna está pronto para deploy em produção com:

- ✅ Containerização completa com Docker
- ✅ Persistência robusta de dados SQLite
- ✅ Backup/restore automático
- ✅ Configuração via ambiente
- ✅ Scripts de automação para VPS
- ✅ Documentação completa para operação
- ✅ **NOVO:** Suporte ao Ecossistema de 4 Módulos (PDF v1.0)
- ✅ **NOVO:** Sistema de Ajuda Educativa (RF-30)
- ✅ **NOVO:** Health checks por módulo

Para começar, execute `./deploy.sh` na sua VPS ou siga o guia manual acima.

Próximos passos recomendados:

1. Teste o deploy em ambiente staging
2. Configure backup automático via cron
3. Implemente proxy reverso (Nginx/Apache)
4. Configure monitoramento básico
5. **NOVO:** Valide coleta de indicadores econômicos
6. **NOVO:** Popule tópicos de ajuda educativa
7. **NOVO:** Teste match de crédito com entidades reais

**Status:** ✅ PRONTO PARA PRODUÇÃO + ECOSSISTEMA COMPLETO
**Última Atualização:** 27/03/2026
**Versão:** 2.0 (Ecossistema de 4 Módulos)
