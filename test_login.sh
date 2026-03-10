#!/bin/bash

# Matar processos na porta 8090
fuser -k 8090/tcp 2>/dev/null || true

# Entrar no diretório
cd /home/s015533607/Documentos/desenv/digna-workspace/modules/ui_web

# Compilar
echo "Compilando..."
go build -o digna_server .

# Rodar servidor
echo "Iniciando servidor..."
./digna_server &
SERVER_PID=$!

# Esperar servidor iniciar
sleep 3

# Testar acesso
echo "Testando acesso..."
curl -s http://localhost:8090/ | head -5
echo ""
echo "Testando login..."
curl -s http://localhost:8090/login | head -100

# Parar servidor
kill $SERVER_PID 2>/dev/null || true