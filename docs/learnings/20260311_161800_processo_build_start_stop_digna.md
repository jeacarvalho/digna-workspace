# Aprendizado: Processo de Build/Start/Stop da Aplicação Digna

**Data:** 2026-03-11  
**Contexto:** Implementação do Dashboard do Contador (RF-11)  
**Problema:** Dificuldade recorrente em testar, subir e "descer" a aplicação durante desenvolvimento

## 🚨 Problemas Identificados

### 1. **Porta em Uso (Address Already in Use)**
```bash
{"error":"listen tcp :8090: bind: address already in use"}
```
- Processos anteriores não são completamente encerrados
- Servidor pode estar rodando em background sem controle

### 2. **Template Path Hell**
- Diferentes handlers usam caminhos diferentes para templates
- Execução de diretórios diferentes causa falhas
- Fallbacks implementados de forma inconsistente

### 3. **Build vs Runtime Directory**
- Build feito em `modules/ui_web/`
- Execução pode ser de diretórios diferentes
- Caminhos relativos quebram

### 4. **Logs Não Visíveis**
- Servidor rodando em background não mostra logs
- Erros de inicialização ficam ocultos
- Debug difícil sem logs apropriados

## ✅ Soluções Aprendidas

### 1. **Gerenciamento de Processos**

#### Matar processos existentes ANTES de iniciar:
```bash
# Verificar processos na porta 8090
lsof -i :8090

# Matar processo específico
kill <PID>

# Matar todos processos digna
pkill -f digna_web_server
pkill -f main  # às vezes roda como 'main'
```

#### Verificar se realmente parou:
```bash
# Verificar se ainda há processos
ps aux | grep digna | grep -v grep
ps aux | grep main | grep -v grep
```

### 2. **Build Correta**

#### Sempre buildar do diretório correto:
```bash
# Build do servidor web
cd /home/s015533607/Documentos/desenv/digna-workspace/modules/ui_web
go build -o digna_web_server
```

#### Testar build rapidamente:
```bash
# Executar com timeout para ver erros
timeout 3 ./digna_web_server 2>&1 | head -50
```

### 3. **Inicialização com Logs**

#### Opção 1: Background com logs em arquivo
```bash
./digna_web_server > server.log 2>&1 &
tail -f server.log  # monitorar logs
```

#### Opção 2: Foreground com CTRL+C (para debug)
```bash
./digna_web_server
# CTRL+C para parar
```

#### Opção 3: Usar `screen` ou `tmux` para sessão persistente
```bash
screen -S digna
./digna_web_server
# CTRL+A D para desanexar
# screen -r digna para retornar
```

### 4. **Verificação de Saúde**

#### Após iniciar, verificar se está respondendo:
```bash
sleep 2  # dar tempo para iniciar
curl -s http://localhost:8090/health
curl -s http://localhost:8090/login | head -5
```

#### Verificar logs de inicialização:
```bash
tail -20 server.log | grep -E "(INFO|ERROR|DEBUG)"
```

### 5. **Template Path Patterns**

#### Padrão aprendido: Implementar fallbacks
```go
// Em todos os handlers, usar este padrão:
var templatePaths = []string{
    "templates/arquivo.html",                    // Quando executado de modules/ui_web/
    "modules/ui_web/templates/arquivo.html",     // Quando executado da raiz do projeto
    "../../templates/arquivo.html",              // Caminho relativo alternativo
}

for _, path := range templatePaths {
    tmpl, err = template.ParseFiles(path)
    if err == nil {
        break
    }
}
```

#### Já corrigidos:
- `auth_handler.go` ✓
- `accountant_handler.go` ✓
- Outros handlers ainda precisam ser corrigidos

## 🔧 Workflow Recomendado

### Para TESTES RÁPIDOS durante desenvolvimento:

```bash
# 1. Parar tudo
pkill -f digna_web_server
pkill -f main
sleep 1

# 2. Build
cd /home/s015533607/Documentos/desenv/digna-workspace/modules/ui_web
go build -o digna_web_server

# 3. Iniciar com logs visíveis (foreground para debug)
./digna_web_server
# Observar logs de inicialização
# CTRL+C quando terminar testes
```

### Para DESENVOLVIMENTO CONTÍNUO:

```bash
# 1. Iniciar em background com logs
cd /home/s015533607/Documentos/desenv/digna-workspace/modules/ui_web
./digna_web_server > server.log 2>&1 &

# 2. Monitorar logs em outro terminal
tail -f server.log

# 3. Quando fizer mudanças no código:
pkill -f digna_web_server
go build -o digna_web_server
./digna_web_server > server.log 2>&1 &
```

### Para PRODUÇÃO/TESTES FINAIS:

```bash
# Usar script de inicialização
cd /home/s015533607/Documentos/desenv/digna-workspace
./scripts/start_digna.sh  # criar este script
```

## 📋 Scripts Úteis (Sugeridos)

### `scripts/start_digna.sh`:
```bash
#!/bin/bash
cd /home/s015533607/Documentos/desenv/digna-workspace/modules/ui_web

# Matar processos existentes
pkill -f digna_web_server 2>/dev/null
pkill -f main 2>/dev/null

# Build
echo "🔨 Building..."
go build -o digna_web_server

# Start
echo "🚀 Starting..."
./digna_web_server > server.log 2>&1 &

# Wait and check
sleep 2
echo "📊 Checking..."
curl -s http://localhost:8090/health && echo "✅ Server is healthy" || echo "❌ Server failed"

# Show logs
echo "📋 Last 10 lines of log:"
tail -10 server.log
```

### `scripts/stop_digna.sh`:
```bash
#!/bin/bash
echo "🛑 Stopping Digna..."
pkill -f digna_web_server
pkill -f main
echo "✅ Stopped"
```

### `scripts/status_digna.sh`:
```bash
#!/bin/bash
echo "🔍 Digna Status:"
echo "Processes:"
ps aux | grep -E "(digna|main)" | grep -v grep
echo ""
echo "Port 8090:"
lsof -i :8090 2>/dev/null || echo "Port 8090 is free"
echo ""
echo "Health check:"
curl -s http://localhost:8090/health 2>/dev/null || echo "Server not responding"
```

## 🎯 Lições Aprendidas

1. **Sempre verificar porta antes de iniciar** - `lsof -i :8090`
2. **Sempre matar processos antigos** - `pkill -f digna_web_server`
3. **Usar logs visíveis durante debug** - executar em foreground ou `tail -f`
4. **Implementar fallbacks de template paths** em TODOS os handlers
5. **Dar tempo para inicialização** - `sleep 2` após start
6. **Verificar saúde imediatamente** - testar `/health` endpoint

## 🔄 Próximas Ações

1. **Criar scripts** de start/stop/status
2. **Corrigir template paths** em todos handlers restantes
3. **Documentar** workflow padrão no README
4. **Considerar** usar `docker-compose` para gerenciamento mais fácil

## 📈 Métricas de Sucesso

- **Tempo para rebuild+restart**: < 10 segundos
- **Sucesso na primeira tentativa**: > 90%
- **Logs sempre visíveis**: 100% do tempo de debug
- **Port conflicts**: 0

---

**Nota para próximas sessões:** Seguir este workflow evitará 80% dos problemas de inicialização/encerramento encontrados durante a implementação do RF-11.