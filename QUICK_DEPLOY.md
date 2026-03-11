# ⚡ Deploy Rápido - Digna

## 🚀 Em 5 Minutos na Sua VPS

### 1. Pré-requisitos
```bash
# Instalar Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Instalar Git e Curl
sudo apt-get update && sudo apt-get install -y git curl
```

### 2. Deploy Automático (RECOMENDADO)
```bash
# Baixar e executar
curl -O https://raw.githubusercontent.com/providentia/digna/main/deploy.sh
chmod +x deploy.sh
./deploy.sh
```

### 3. Ou Deploy Manual
```bash
# Clonar
git clone https://github.com/providentia/digna.git
cd digna

# Configurar
cp .env.example .env
# Edite .env se necessário

# Rodar
docker-compose up -d

# Verificar
curl http://localhost:8090/health
```

## 📍 URLs Importantes
- **Aplicação:** `http://SEU_IP:8090`
- **Health check:** `http://SEU_IP:8090/health`
- **Login padrão:** `cafe_digna` / `cd0123`

## 🔧 Comandos Úteis
```bash
# Logs
docker-compose logs -f

# Parar
docker-compose stop

# Reiniciar  
docker-compose restart

# Backup dados
./scripts/deploy/backup.sh

# Atualizar
./deploy.sh --update
```

## ⚠️ Primeiros Passos
1. Acesse `http://SEU_IP:8090`
2. Faça login com `cafe_digna` / `cd0123`
3. Configure backup automático:
   ```bash
   # Backup diário às 2AM (mantém 7 dias)
   0 2 * * * /opt/digna/scripts/deploy/backup.sh
   ```

## 📞 Precisa de Ajuda?
- Verifique `docs/DEPLOYMENT.md` para guia completo
- Consulte `docker-compose logs` para erros
- Issues no GitHub para problemas

---

**✅ Pronto!** Seu Digna está rodando em produção. 🎉