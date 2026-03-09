package handler

// supplyTemplates contém templates HTML embutidos para fallback
const supplyTemplates = `
{{define "supply_dashboard.html"}}
<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
</head>
<body class="bg-gray-100 min-h-screen">
    <nav class="bg-blue-600 text-white shadow-lg sticky top-0 z-50">
        <div class="container mx-auto px-4 py-3">
            <div class="flex items-center justify-between">
                <div class="flex items-center space-x-2">
                    <a href="/" class="text-white hover:text-blue-100">
                        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
                        </svg>
                    </a>
                    <span class="text-xl font-bold">Compras e Estoque</span>
                </div>
            </div>
        </div>
    </nav>

    <main class="container mx-auto px-4 py-6 max-w-4xl">
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
            <a href="/supply/purchase" class="bg-white rounded-xl shadow-md p-6 hover:shadow-lg transition text-center">
                <div class="text-blue-500 mb-3">
                    <svg class="w-12 h-12 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v13m0-13V6a2 2 0 112 2h-2zm0 0V5.5A2.5 2.5 0 109.5 8H12zm-7 4h14M5 12a2 2 0 110-4h14a2 2 0 110 4M5 12v7a2 2 0 002 2h10a2 2 0 002-2v-7"></path>
                    </svg>
                </div>
                <h3 class="text-lg font-semibold text-gray-800">Nova Compra</h3>
                <p class="text-gray-600 text-sm mt-2">Registrar compra de material</p>
            </a>

            <a href="/supply/suppliers" class="bg-white rounded-xl shadow-md p-6 hover:shadow-lg transition text-center">
                <div class="text-green-500 mb-3">
                    <svg class="w-12 h-12 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"></path>
                    </svg>
                </div>
                <h3 class="text-lg font-semibold text-gray-800">Fornecedores</h3>
                <p class="text-gray-600 text-sm mt-2">Gerenciar quem fornece</p>
            </a>

            <a href="/supply/stock" class="bg-white rounded-xl shadow-md p-6 hover:shadow-lg transition text-center">
                <div class="text-amber-500 mb-3">
                    <svg class="w-12 h-12 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"></path>
                    </svg>
                </div>
                <h3 class="text-lg font-semibold text-gray-800">Meu Estoque</h3>
                <p class="text-gray-600 text-sm mt-2">Controlar materiais e produtos</p>
            </a>
        </div>

        <div class="bg-white rounded-xl shadow-md p-6">
            <h2 class="text-2xl font-bold text-gray-800 mb-4">Bem-vindo à Gestão de Compras</h2>
            <p class="text-gray-600 mb-4">
                Aqui você pode registrar suas compras de materiais, gerenciar fornecedores e controlar seu estoque.
                O sistema cuida automaticamente da contabilidade por trás das cenas.
            </p>
            
            <div class="bg-blue-50 border-l-4 border-blue-500 p-4 rounded">
                <p class="text-blue-700 text-sm">
                    <strong>Dica:</strong> Ao registrar uma compra, o sistema automaticamente atualiza seu estoque
                    e registra a transação financeira. Você só precisa informar o que comprou, de quem e por quanto.
                </p>
            </div>
        </div>
    </main>
</body>
</html>
{{end}}

{{define "supply_purchase.html"}}
<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <script>
        function calculateTotal() {
            const quantity = parseInt(document.getElementById('quantity').value) || 0;
            const unitCost = parseInt(document.getElementById('unit_cost').value) || 0;
            const total = quantity * unitCost;
            document.getElementById('total_cost').value = total;
            document.getElementById('total_display').textContent = 'R$ ' + (total / 100).toFixed(2);
        }
        
        function submitPurchase() {
            const form = document.getElementById('purchase-form');
            htmx.ajax('POST', '/api/supply/purchase', {
                values: htmx.values(form),
                target: '#purchase-result',
                swap: 'innerHTML'
            }).then(() => {
                form.reset();
                calculateTotal();
            });
            return false;
        }
    </script>
</head>
<body class="bg-gray-100 min-h-screen">
    <nav class="bg-blue-600 text-white shadow-lg sticky top-0 z-50">
        <div class="container mx-auto px-4 py-3">
            <div class="flex items-center justify-between">
                <div class="flex items-center space-x-2">
                    <a href="/supply" class="text-white hover:text-blue-100">
                        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
                        </svg>
                    </a>
                    <span class="text-xl font-bold">Nova Compra</span>
                </div>
            </div>
        </div>
    </nav>

    <main class="container mx-auto px-4 py-6 max-w-md">
        <div id="purchase-result" class="mb-4"></div>

        <div class="bg-white rounded-2xl shadow-lg p-6">
            <h2 class="text-xl font-bold text-gray-800 mb-6">O que você comprou?</h2>
            
            <form id="purchase-form" onsubmit="return submitPurchase()">
                <div class="mb-4">
                    <label class="block text-sm font-medium text-gray-700 mb-2">De quem você comprou?</label>
                    <select name="supplier_id" required class="w-full p-3 border-2 border-gray-200 rounded-xl text-lg focus:border-blue-500 focus:outline-none">
                        <option value="">Selecione um fornecedor</option>
                        {{range .Suppliers}}
                        <option value="{{.ID}}">{{.Name}}</option>
                        {{end}}
                    </select>
                    <p class="text-gray-500 text-xs mt-1">Não encontrou? <a href="/supply/suppliers" class="text-blue-500 hover:underline">Cadastre um novo fornecedor</a></p>
                </div>

                <div class="mb-4">
                    <label class="block text-sm font-medium text-gray-700 mb-2">O que você comprou?</label>
                    <select name="stock_item_id" required class="w-full p-3 border-2 border-gray-200 rounded-xl text-lg focus:border-blue-500 focus:outline-none">
                        <option value="">Selecione um item</option>
                        {{range .StockItems}}
                        <option value="{{.ID}}">{{.Name}} ({{stockItemTypeLabel .Type}})</option>
                        {{end}}
                    </select>
                    <p class="text-gray-500 text-xs mt-1">Não encontrou? <a href="/supply/stock" class="text-blue-500 hover:underline">Cadastre um novo item</a></p>
                </div>

                <div class="grid grid-cols-2 gap-4 mb-4">
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">Quantidade</label>
                        <input type="number" id="quantity" name="quantity" min="1" required 
                               oninput="calculateTotal()"
                               class="w-full p-3 border-2 border-gray-200 rounded-xl text-lg focus:border-blue-500 focus:outline-none">
                    </div>
                    
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">Valor unitário (R$)</label>
                        <div class="relative">
                            <span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500">R$</span>
                            <input type="number" id="unit_cost" name="unit_cost" min="1" step="1" required 
                                   oninput="calculateTotal()"
                                   class="w-full p-3 pl-10 border-2 border-gray-200 rounded-xl text-lg focus:border-blue-500 focus:outline-none"
                                   placeholder="0,00">
                        </div>
                    </div>
                </div>

                <div class="mb-4">
                    <label class="block text-sm font-medium text-gray-700 mb-2">Forma de pagamento</label>
                    <div class="flex space-x-4">
                        <label class="flex items-center">
                            <input type="radio" name="payment_type" value="CASH" checked class="mr-2">
                            <span>À vista</span>
                        </label>
                        <label class="flex items-center">
                            <input type="radio" name="payment_type" value="CREDIT" class="mr-2">
                            <span>A prazo</span>
                        </label>
                    </div>
                </div>

                <div class="mb-6 p-4 bg-blue-50 rounded-xl">
                    <label class="block text-sm font-medium text-gray-700 mb-2">Valor total da compra</label>
                    <div class="text-3xl font-bold text-blue-600" id="total_display">R$ 0,00</div>
                    <input type="hidden" id="total_cost" name="total_cost" value="0">
                </div>

                <button type="submit" 
                        class="w-full bg-blue-500 hover:bg-blue-600 text-white text-xl font-bold py-4 rounded-xl shadow-lg transition transform active:scale-95">
                    REGISTRAR COMPRA
                </button>
            </form>
        </div>

        <div class="text-center text-gray-500 text-sm mt-6">
            <p>O sistema atualizará automaticamente seu estoque e registrará a transação financeira.</p>
        </div>
    </main>
</body>
</html>
{{end}}

{{define "supply_suppliers.html"}}
<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <script>
        function submitSupplier() {
            const form = document.getElementById('supplier-form');
            htmx.ajax('POST', '/api/supply/supplier', {
                values: htmx.values(form),
                target: '#supplier-result',
                swap: 'innerHTML'
            }).then(() => {
                form.reset();
                // Recarregar lista de fornecedores
                htmx.ajax('GET', '/api/supply/supplier', {
                    target: '#suppliers-list',
                    swap: 'innerHTML'
                });
            });
            return false;
        }
    </script>
</head>
<body class="bg-gray-100 min-h-screen">
    <nav class="bg-green-600 text-white shadow-lg sticky top-0 z-50">
        <div class="container mx-auto px-4 py-3">
            <div class="flex items-center justify-between">
                <div class="flex items-center space-x-2">
                    <a href="/supply" class="text-white hover:text-green-100">
                        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
                        </svg>
                    </a>
                    <span class="text-xl font-bold">Fornecedores</span>
                </div>
            </div>
        </div>
    </nav>

    <main class="container mx-auto px-4 py-6 max-w-4xl">
        <div id="supplier-result" class="mb-4"></div>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <!-- Formulário de cadastro -->
            <div class="bg-white rounded-2xl shadow-lg p-6">
                <h2 class="text-xl font-bold text-gray-800 mb-6">Cadastrar Novo Fornecedor</h2>
                
                <form id="supplier-form" onsubmit="return submitSupplier()">
                    <div class="mb-4">
                        <label class="block text-sm font-medium text-gray-700 mb-2">Nome do fornecedor</label>
                        <input type="text" name="name" required 
                               class="w-full p-3 border-2 border-gray-200 rounded-xl text-lg focus:border-green-500 focus:outline-none"
                               placeholder="Ex: João das Velas">
                    </div>

                    <div class="mb-6">
                        <label class="block text-sm font-medium text-gray-700 mb-2">Informações de contato</label>
                        <textarea name="contact_info" rows="3"
                               class="w-full p-3 border-2 border-gray-200 rounded-xl text-lg focus:border-green-500 focus:outline-none"
                               placeholder="Telefone, email, endereço..."></textarea>
                    </div>

                    <button type="submit" 
                            class="w-full bg-green-500 hover:bg-green-600 text-white text-xl font-bold py-4 rounded-xl shadow-lg transition transform active:scale-95">
                        SALVAR FORNECEDOR
                    </button>
                </form>
            </div>

            <!-- Lista de fornecedores -->
            <div class="bg-white rounded-2xl shadow-lg p-6">
                <h2 class="text-xl font-bold text-gray-800 mb-6">Meus Fornecedores</h2>
                
                <div id="suppliers-list">
                    {{if .Suppliers}}
                    <div class="space-y-4">
                        {{range .Suppliers}}
                        <div class="border border-gray-200 rounded-xl p-4 hover:bg-gray-50 transition">
                            <div class="flex justify-between items-start">
                                <div>
                                    <h3 class="font-semibold text-gray-800">{{.Name}}</h3>
                                    {{if .ContactInfo}}
                                    <p class="text-gray-600 text-sm mt-1">{{.ContactInfo}}</p>
                                    {{end}}
                                </div>
                                <span class="text-xs text-gray-500">{{.CreatedAt.Format "02/01/2006"}}</span>
                            </div>
                        </div>
                        {{end}}
                    </div>
                    {{else}}
                    <div class="text-center py-8 text-gray-500">
                        <svg class="w-12 h-12 mx-auto text-gray-300 mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"></path>
                        </svg>
                        <p>Nenhum fornecedor cadastrado ainda.</p>
                        <p class="text-sm mt-1">Cadastre seu primeiro fornecedor usando o formulário ao lado.</p>
                    </div>
                    {{end}}
                </div>
            </div>
        </div>
    </main>
</body>
</html>
{{end}}

{{define "supply_stock.html"}}
<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <script>
        function submitStockItem() {
            const form = document.getElementById('stock-item-form');
            htmx.ajax('POST', '/api/supply/stock-item', {
                values: htmx.values(form),
                target: '#stock-item-result',
                swap: 'innerHTML'
            }).then(() => {
                form.reset();
                // Recarregar lista de itens
                htmx.ajax('GET', '/api/supply/stock-item', {
                    target: '#stock-items-list',
                    swap: 'innerHTML'
                });
            });
            return false;
        }
    </script>
</head>
<body class="bg-gray-100 min-h-screen">
    <nav class="bg-amber-600 text-white shadow-lg sticky top-0 z-50">
        <div class="container mx-auto px-4 py-3">
            <div class="flex items-center justify-between">
                <div class="flex items-center space-x-2">
                    <a href="/supply" class="text-white hover:text-amber-100">
                        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
                        </svg>
                    </a>
                    <span class="text-xl font-bold">Meu Estoque</span>
                </div>
            </div>
        </div>
    </nav>

    <main class="container mx-auto px-4 py-6 max-w-4xl">
        <div id="stock-item-result" class="mb-4"></div>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <!-- Formulário de cadastro -->
            <div class="bg-white rounded-2xl shadow-lg p-6">
                <h2 class="text-xl font-bold text-gray-800 mb-6">Cadastrar Novo Item</h2>
                
                <form id="stock-item-form" onsubmit="return submitStockItem()">
                    <div class="mb-4">
                        <label class="block text-sm font-medium text-gray-700 mb-2">Nome do item</label>
                        <input type="text" name="name" required 
                               class="w-full p-3 border-2 border-gray-200 rounded-xl text-lg focus:border-amber-500 focus:outline-none"
                               placeholder="Ex: Cera de Abelha">
                    </div>

                    <div class="mb-4">
                        <label class="block text-sm font-medium text-gray-700 mb-2">Tipo do item</label>
                        <select name="type" required class="w-full p-3 border-2 border-gray-200 rounded-xl text-lg focus:border-amber-500 focus:outline-none">
                            <option value="">Selecione o tipo</option>
                            <option value="INSUMO">Insumo (matéria-prima)</option>
                            <option value="PRODUTO">Produto (para venda)</option>
                            <option value="MERCADORIA">Mercadoria (para revenda)</option>
                        </select>
                        <p class="text-gray-500 text-xs mt-1">
                            <strong>Insumo:</strong> matéria-prima para produção<br>
                            <strong>Produto:</strong> produto acabado para venda<br>
                            <strong>Mercadoria:</strong> produto para revenda
                        </p>
                    </div>

                    <div class="grid grid-cols-2 gap-4 mb-4">
                        <div>
                            <label class="block text-sm font-medium text-gray-700 mb-2">Quantidade inicial</label>
                            <input type="number" name="quantity" min="0" required 
                                   class="w-full p-3 border-2 border-gray-200 rounded-xl text-lg focus:border-amber-500 focus:outline-none">
                        </div>
                        
                        <div>
                            <label class="block text-sm font-medium text-gray-700 mb-2">Quantidade mínima</label>
                            <input type="number" name="min_quantity" min="0" required 
                                   class="w-full p-3 border-2 border-gray-200 rounded-xl text-lg focus:border-amber-500 focus:outline-none">
                        </div>
                    </div>

                    <div class="mb-6">
                        <label class="block text-sm font-medium text-gray-700 mb-2">Custo unitário (R$)</label>
                        <div class="relative">
                            <span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500">R$</span>
                            <input type="number" name="unit_cost" min="1" step="1" required 
                                   class="w-full p-3 pl-10 border-2 border-gray-200 rounded-xl text-lg focus:border-amber-500 focus:outline-none"
                                   placeholder="0,00">
                        </div>
                    </div>

                    <button type="submit" 
                            class="w-full bg-amber-500 hover:bg-amber-600 text-white text-xl font-bold py-4 rounded-xl shadow-lg transition transform active:scale-95">
                        SALVAR ITEM
                    </button>
                </form>
            </div>

            <!-- Lista de itens -->
            <div class="bg-white rounded-2xl shadow-lg p-6">
                <h2 class="text-xl font-bold text-gray-800 mb-6">Itens em Estoque</h2>
                
                <div id="stock-items-list">
                    {{if .StockItems}}
                    <div class="space-y-4">
                        {{range .StockItems}}
                        <div class="border border-gray-200 rounded-xl p-4 hover:bg-gray-50 transition {{if isBelowMinimum .Quantity .MinQuantity}}border-red-200 bg-red-50{{end}}">
                            <div class="flex justify-between items-start">
                                <div>
                                    <div class="flex items-center space-x-2">
                                        <h3 class="font-semibold text-gray-800">{{.Name}}</h3>
                                        <span class="text-xs px-2 py-1 rounded-full {{if eq .Type "INSUMO"}}bg-blue-100 text-blue-800{{else if eq .Type "PRODUTO"}}bg-green-100 text-green-800{{else}}bg-purple-100 text-purple-800{{end}}">
                                            {{stockItemTypeLabel .Type}}
                                        </span>
                                    </div>
                                    <div class="mt-2 grid grid-cols-3 gap-4 text-sm">
                                        <div>
                                            <span class="text-gray-500">Quantidade:</span>
                                            <span class="font-semibold ml-1 {{if isBelowMinimum .Quantity .MinQuantity}}text-red-600{{else}}text-gray-800{{end}}">
                                                {{.Quantity}}
                                            </span>
                                        </div>
                                        <div>
                                            <span class="text-gray-500">Mínimo:</span>
                                            <span class="font-semibold ml-1">{{.MinQuantity}}</span>
                                        </div>
                                        <div>
                                            <span class="text-gray-500">Custo:</span>
                                            <span class="font-semibold ml-1">R$ {{formatCurrency .UnitCost}}</span>
                                        </div>
                                    </div>
                                    {{if isBelowMinimum .Quantity .MinQuantity}}
                                    <div class="mt-2 text-xs text-red-600 font-medium">
                                        ⚠️ Estoque abaixo do mínimo
                                    </div>
                                    {{end}}
                                </div>
                                <span class="text-xs text-gray-500">{{.CreatedAt.Format "02/01/2006"}}</span>
                            </div>
                        </div>
                        {{end}}
                    </div>
                    {{else}}
                    <div class="text-center py-8 text-gray-500">
                        <svg class="w-12 h-12 mx-auto text-gray-300 mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"></path>
                        </svg>
                        <p>Nenhum item cadastrado ainda.</p>
                        <p class="text-sm mt-1">Cadastre seu primeiro item usando o formulário ao lado.</p>
                    </div>
                    {{end}}
                </div>
            </div>
        </div>
    </main>
</body>
</html>
{{end}}`
