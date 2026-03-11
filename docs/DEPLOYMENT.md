# 🚀 Guia de Deploy em Produção - Projeto Digna

**Última atualização:** 11/03/2026  
**Status:** ✅ PRONTO PARA PRODUÇÃO

---

## 📋 Visão Geral

Este documento descreve o processo completo de deploy do projeto Digna em ambiente de produção (VPS). O sistema foi containerizado com Docker e utiliza variáveis de ambiente para configuração.

### 🎯 Objetivos do Deploy
1. **Automação completa** - Script único para deploy em VPS
2. **Persistência de dados** - Bancos SQLite em volumes externos
3. **Configuração via ambiente** - Zero hardcoding no código
4. **Backup/restore** - Sistema robusto de backup dos dados
5. **Manutenção simplificada** - Comandos simples para operações

---

## 🏗️ Arquitetura de Deploy

```
VPS (Ubuntu/Debian)
├── Docker + docker-compose
├── Digna Container
│   ├── Aplicação Go
│   ├── Templates/Static files
│   └── Configuração via .env
└── Volumes Externos
    └── data/entities/ (SQLite databases)
```

---

## 📁 Estrutura de Arquivos de Deploy

```
scripts/deploy/
├── vps_deploy.sh          # Script principal de deploy
├── backup.sh              # Backup de bancos SQLite
└── restore.sh             # Restauração de backup

deploy.sh                  # Wrapper script (raiz do projeto)
docker-compose.yml         # Configuração dev/docker-compose
docker-compose.prod.yml    # Configuração produção
.env.example               # Template variáveis ambiente
Dockerfile                 # Build da aplicação
```

---

## 🚀 Deploy Rápido (VPS Nova)

### Pré-requisitos na VPS
```bash
# 1. Docker instalado
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# 2. Git instalado
sudo apt-get update && sudo apt-get install -y git curl
```

### Deploy Automático (Recomendado)
```bash
# 1. Baixar script de deploy
curl -O https://raw.githubusercontent.com/providentia/digna/main/deploy.sh
chmod +x deploy.sh

# 2. Executar deploy
./deploy.sh

# 3. Para atualizar posteriormente
./deploy.sh --update
```

### Deploy Manual (Passo a Passo)
```bash
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
```

---

## ⚙️ Configuração de Ambiente (.env)

### Arquivo `.env` mínimo
```bash
# Server Configuration
DIGNA_PORT=8090

# Data Directory Configuration
DIGNA_DATA_DIR=/var/lib/digna/data

# Logging Configuration
DIGNA_LOG_LEVEL=info

# Docker Compose Configuration
COMPOSE_PROJECT_NAME=digna
```

### Variáveis de Ambiente Disponíveis

| Variável | Descrição | Padrão | Obrigatório |
|----------|-----------|--------|-------------|
| `DIGNA_PORT` | Porta HTTP da aplicação | `8090` | Não |
| `DIGNA_DATA_DIR` | Diretório para bancos SQLite | `./data/entities` | Não |
| `DIGNA_LOG_LEVEL` | Nível de log (debug, info, warn, error) | `info` | Não |
| `COMPOSE_PROJECT_NAME` | Nome do projeto Docker | `digna` | Não |

### Configuração Avançada
```bash
# Para alta disponibilidade
DIGNA_PORT=8080
DIGNA_DATA_DIR=/mnt/ssd/digna/data
DIGNA_LOG_LEVEL=warn

# Com proxy reverso (Nginx)
DIGNA_PORT=3000  # Interno no container
# Nginx redireciona porta 80/443 para 3000
```

---

## 💾 Persistência de Dados

### Estrutura de Dados
```
/var/lib/digna/data/          # Volume Docker
├── entity_abc123.db          # Banco da entidade 1
├── entity_def456.db          # Banco da entidade 2
└── entity_ghi789.db          # Banco da entidade 3
```

### Backup Automático
```bash
# Backup manual
./scripts/deploy/backup.sh

# Backup com retenção personalizada
./scripts/deploy/backup.sh --output-dir=/backups --keep-days=30

# Agendamento via cron (backup diário às 2AM)
0 2 * * * /opt/digna/scripts/deploy/backup.sh --keep-days=7
```

### Restauração de Backup
```bash
# Listar backups disponíveis
./scripts/deploy/restore.sh

# Restaurar backup específico
./scripts/deploy/restore.sh --backup-file=/backups/digna_backup_20250311_020000.tar.gz

# Dry run (teste sem alterações)
./scripts/deploy/restore.sh --dry-run
```

---

## 🔧 Manutenção e Operações

### Comandos Comuns
```bash
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
```

### Monitoramento
```bash
# Health check
curl http://localhost:8090/health

# Readiness check
curl http://localhost:8090/ready

# Status do container
docker-compose ps
docker-compose top
```

### Limpeza
```bash
# Limpar containers parados
docker-compose rm -f

# Limpar imagens não utilizadas
docker image prune -a

# Limpar volumes não utilizados
docker volume prune
```

---

## 📊 Escalabilidade

### Recursos do Container
```yaml
# Em docker-compose.prod.yml
deploy:
  resources:
    limits:
      memory: 512M
      cpus: '1.0'
    reservations:
      memory: 256M
      cpus: '0.5'
```

### Otimizações para Produção
1. **Use `docker-compose.prod.yml`** para limites de recursos
2. **Configure swap** na VPS para picos de memória
3. **Use volume externo** para dados (não bind mount)
4. **Configure log rotation** para logs do Docker
5. **Use rede interna** para segurança

---

## 🔒 Segurança

### Melhores Práticas
1. **Não exponha a porta diretamente** - Use Nginx/Apache como proxy
2. **Configure firewall** - Apenas portas necessárias
3. **Use HTTPS** - Certificado SSL via Let's Encrypt
4. **Atualize regularmente** - Docker, sistema, aplicação
5. **Backup frequente** - Dados são críticos

### Exemplo Nginx como Proxy
```nginx
server {
    listen 80;
    server_name digna.example.com;
    
    location / {
        proxy_pass http://localhost:8090;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

---

## 🐛 Troubleshooting

### Problemas Comuns

#### 1. Container não inicia
```bash
# Ver logs
docker-compose logs

# Verificar portas em uso
sudo netstat -tulpn | grep :8090

# Testar build manual
docker-compose build --no-cache
```

#### 2. Erro de permissão nos dados
```bash
# Corrigir permissões
sudo chown -R 1000:1000 /var/lib/digna/data

# No container
docker-compose exec digna chown -R digna:digna /var/lib/digna/data
```

#### 3. Health check falha
```bash
# Verificar se aplicação responde
curl -v http://localhost:8090/health

# Verificar logs da aplicação
docker-compose logs digna

# Reiniciar container
docker-compose restart
```

#### 4. Espaço em disco
```bash
# Verificar uso
df -h /var/lib/digna

# Limpar backups antigos
find /var/backups/digna -name "*.tar.gz" -mtime +30 -delete

# Limpar logs do Docker
docker system prune -a
```

---

## 📈 Monitoramento e Métricas

### Métricas Importantes
1. **Uso de CPU/Memória** - `docker stats`
2. **Logs da aplicação** - `docker-compose logs`
3. **Health checks** - Monitorar endpoint `/health`
4. **Espaço em disco** - Monitorar volume de dados
5. **Backups** - Verificar execução regular

### Ferramentas Recomendadas
- **Docker Compose** - Gerenciamento de containers
- **Cron** - Agendamento de backups
- **Logrotate** - Rotação de logs
- **Monitoramento** - Prometheus + Grafana (opcional)

---

## 🔄 Atualização da Aplicação

### Fluxo de Atualização
```bash
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
```

### Rollback (em caso de problemas)
```bash
# 1. Parar aplicação
docker-compose stop

# 2. Restaurar backup
./scripts/deploy/restore.sh --backup-file=/backups/digna_backup_YYYYMMDD_HHMMSS.tar.gz

# 3. Voltar para versão anterior do código
git checkout <commit-hash>

# 4. Rebuild e start
docker-compose up -d --build
```

---

## 📞 Suporte e Contato

### Recursos
- **Documentação:** `docs/` no repositório
- **Issues:** GitHub Issues do projeto
- **Scripts:** `scripts/deploy/` para operações

### Checklist Pós-Deploy
- [ ] Aplicação responde em `http://localhost:8090`
- [ ] Health check passa (`/health`)
- [ ] Dados persistem após restart
- [ ] Backup configurado e testado
- [ ] Logs estão sendo gerados
- [ ] Monitoramento configurado

---

## 🎯 Conclusão

O sistema Digna está pronto para deploy em produção com:

✅ **Containerização completa** com Docker  
✅ **Persistência robusta** de dados SQLite  
✅ **Backup/restore automático**  
✅ **Configuração via ambiente**  
✅ **Scripts de automação** para VPS  
✅ **Documentação completa** para operação

Para começar, execute `./deploy.sh` na sua VPS ou siga o guia manual acima.

**Próximos passos recomendados:**
1. Teste o deploy em ambiente staging
2. Configure backup automático via cron
3. Implemente proxy reverso (Nginx/Apache)
4. Configure monitoramento básico
5. Documente procedimentos específicos da sua organização