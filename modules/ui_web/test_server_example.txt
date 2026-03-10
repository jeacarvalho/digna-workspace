package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			<!DOCTYPE html>
			<html>
			<head>
				<title>Teste Digna</title>
				<script src="https://cdn.tailwindcss.com"></script>
			</head>
			<body class="bg-gray-100 p-8">
				<h1 class="text-3xl font-bold mb-6">Teste Digna Web</h1>
				<div class="space-y-4">
					<a href="http://localhost:8088/" class="block p-4 bg-white rounded shadow hover:bg-blue-50">🏠 Home</a>
					<a href="http://localhost:8088/dashboard" class="block p-4 bg-white rounded shadow hover:bg-blue-50">📊 Dashboard</a>
					<a href="http://localhost:8088/cash?entity_id=test-entity" class="block p-4 bg-white rounded shadow hover:bg-blue-50">💰 Caixa</a>
					<a href="http://localhost:8088/supply?entity_id=test-entity" class="block p-4 bg-white rounded shadow hover:bg-blue-50">📦 Compras/Estoque</a>
					<a href="http://localhost:8088/budget?entity_id=test-entity" class="block p-4 bg-white rounded shadow hover:bg-blue-50">📈 Orçamento</a>
					<a href="http://localhost:8088/accountant?entity_id=test-entity" class="block p-4 bg-white rounded shadow hover:bg-blue-50">👔 Contador Social</a>
				</div>
				<div class="mt-8 p-4 bg-yellow-50 rounded">
					<p class="text-sm text-gray-700">Nota: Alguns módulos podem retornar erro 500 se houver problemas de inicialização.</p>
				</div>
			</body>
			</html>
		`)
	})

	fmt.Println("Servidor de teste rodando em http://localhost:9999")
	log.Fatal(http.ListenAndServe(":9999", nil))
}
