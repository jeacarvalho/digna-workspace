#!/bin/bash

echo "🧪 Testando configuração Docker do Digna"
echo "========================================"

# Test 1: Build Docker image
echo "1. Construindo imagem Docker..."
docker build -t digna-test:latest . > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "   ✅ Build Docker bem-sucedido"
else
    echo "   ❌ Falha no build Docker"
    exit 1
fi

# Test 2: Run container with test port
echo "2. Executando container de teste..."
docker run -d --name digna-test-container \
  -p 9090:9090 \
  -e DIGNA_PORT=9090 \
  -e DIGNA_DATA_DIR=/tmp/test-data \
  digna-test:latest > /dev/null 2>&1

sleep 3

# Test 3: Check if container is running
echo "3. Verificando se container está rodando..."
if docker ps | grep -q digna-test-container; then
    echo "   ✅ Container está rodando"
else
    echo "   ❌ Container não está rodando"
    docker logs digna-test-container
    docker rm -f digna-test-container > /dev/null 2>&1
    exit 1
fi

# Test 4: Check health endpoint
echo "4. Testando endpoint de health..."
sleep 2
if curl -s http://localhost:9090/health | grep -q '"status":"ok"'; then
    echo "   ✅ Health check passou"
else
    echo "   ❌ Health check falhou"
    docker logs digna-test-container
    docker rm -f digna-test-container > /dev/null 2>&1
    exit 1
fi

# Test 5: Cleanup
echo "5. Limpando container de teste..."
docker rm -f digna-test-container > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "   ✅ Container removido com sucesso"
else
    echo "   ⚠️  Não foi possível remover container"
fi

echo ""
echo "🎉 Todos os testes passaram!"
echo "A aplicação Digna está pronta para deploy em produção com Docker."