package handler

// budgetTemplates contém templates HTML embutidos para orçamento
const budgetTemplates = `
{{define "budget_dashboard.html"}}
<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <script>
        function submitBudgetPlan() {
            const form = document.getElementById('budget-plan-form');
            htmx.ajax('POST', '/api/budget/plan', {
                values: htmx.values(form),
                target: '#budget-result',
                swap: 'innerHTML'
            }).then(() => {
                form.reset();
                // Recarregar relatório
                htmx.ajax('GET', '/budget/report?period=' + document.getElementById('period').value, {
                    target: '#budget-report',
                    swap: 'innerHTML'
                });
            });
            return false;
        }

        function loadReport(period) {
            htmx.ajax('GET', '/budget/report?period=' + period, {
                target: '#budget-report',
                swap: 'innerHTML'
            });
        }

        function deletePlan(planId) {
            if (confirm('Tem certeza que quer remover este planejamento?')) {
                htmx.ajax('DELETE', '/api/budget/plan/' + planId, {
                    target: '#budget-result',
                    swap: 'innerHTML'
                }).then(() => {
                    // Recarregar relatório
                    const period = document.getElementById('period').value;
                    htmx.ajax('GET', '/budget/report?period=' + period, {
                        target: '#budget-report',
                        swap: 'innerHTML'
                    });
                });
            }
        }

        function formatCurrency(value) {
            return 'R$ ' + (value / 100).toFixed(2).replace('.', ',');
        }

        function getStatusColor(status) {
            switch(status) {
                case 'SAFE': return 'bg-green-100 text-green-800';
                case 'WARNING': return 'bg-yellow-100 text-yellow-800';
                case 'EXCEEDED': return 'bg-red-100 text-red-800';
                default: return 'bg-gray-100 text-gray-800';
            }
        }

        function getStatusLabel(status) {
            switch(status) {
                case 'SAFE': return 'Dentro do planejado';
                case 'WARNING': return 'Atenção: perto do limite';
                case 'EXCEEDED': return 'Ultrapassou o planejado';
                default: return status;
            }
        }

        function getProgressColor(percentage) {
            if (percentage <= 70) return 'bg-green-500';
            if (percentage <= 100) return 'bg-yellow-500';
            return 'bg-red-500';
        }
    </script>
</head>
<body class="bg-gray-100 min-h-screen">
    <nav class="bg-purple-600 text-white shadow-lg sticky top-0 z-50">
        <div class="container mx-auto px-4 py-3">
            <div class="flex items-center justify-between">
                <div class="flex items-center space-x-2">
                    <a href="/" class="text-white hover:text-purple-100">
                        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
                        </svg>
                    </a>
                    <span class="text-xl font-bold">Planejamento do Mês</span>
                </div>
            </div>
        </div>
    </nav>

    <main class="container mx-auto px-4 py-6 max-w-6xl">
        <div id="budget-result" class="mb-4"></div>

        <!-- Resumo Geral -->
        <div class="bg-white rounded-2xl shadow-lg p-6 mb-6">
            <div class="flex justify-between items-center mb-4">
                <h2 class="text-xl font-bold text-gray-800">Resumo do Mês</h2>
                <div class="flex items-center space-x-2">
                    <select id="period" onchange="loadReport(this.value)" class="p-2 border border-gray-300 rounded-lg">
                        {{range .Periods}}
                        <option value="{{.}}" {{if eq . $.Period}}selected{{end}}>{{.}}</option>
                        {{end}}
                    </select>
                </div>
            </div>

            {{if .Summary}}
            <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
                <div class="text-center p-4 bg-gray-50 rounded-xl">
                    <div class="text-2xl font-bold text-gray-800">{{formatCurrency .Summary.TotalPlanned}}</div>
                    <div class="text-sm text-gray-600">Planejamos gastar</div>
                </div>
                <div class="text-center p-4 bg-gray-50 rounded-xl">
                    <div class="text-2xl font-bold text-gray-800">{{formatCurrency .Summary.TotalExecuted}}</div>
                    <div class="text-sm text-gray-600">Já gastamos</div>
                </div>
                <div class="text-center p-4 bg-gray-50 rounded-xl">
                    <div class="text-2xl font-bold {{if gt .Summary.TotalRemaining 0}}text-green-600{{else}}text-red-600{{end}}">
                        {{formatCurrency .Summary.TotalRemaining}}
                    </div>
                    <div class="text-sm text-gray-600">Falta gastar</div>
                </div>
                <div class="text-center p-4 rounded-xl {{getStatusColor .Summary.OverallStatus}}">
                    <div class="text-2xl font-bold">{{.Summary.Percentage}}%</div>
                    <div class="text-sm">{{getStatusLabel .Summary.OverallStatus}}</div>
                </div>
            </div>
            {{else}}
            <div class="text-center py-8 text-gray-500">
                <p>Nenhum planejamento registrado para este mês.</p>
                <p class="text-sm mt-1">Comece planejando seus gastos usando o formulário abaixo.</p>
            </div>
            {{end}}
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <!-- Formulário de Planejamento -->
            <div class="bg-white rounded-2xl shadow-lg p-6">
                <h2 class="text-xl font-bold text-gray-800 mb-6">O que planejamos gastar?</h2>
                
                <form id="budget-plan-form" onsubmit="return submitBudgetPlan()">
                    <input type="hidden" name="entity_id" value="{{.EntityID}}">
                    
                    <div class="mb-4">
                        <label class="block text-sm font-medium text-gray-700 mb-2">Mês</label>
                        <input type="text" name="period" value="{{.Period}}" readonly
                               class="w-full p-3 border-2 border-gray-200 rounded-xl text-lg bg-gray-50">
                    </div>

                    <div class="mb-4">
                        <label class="block text-sm font-medium text-gray-700 mb-2">Em que vamos gastar?</label>
                        <select name="category" required class="w-full p-3 border-2 border-gray-200 rounded-xl text-lg focus:border-purple-500 focus:outline-none">
                            <option value="">Selecione uma categoria</option>
                            {{range .Categories}}
                            <option value="{{.ID}}">{{.Label}}</option>
                            {{end}}
                        </select>
                    </div>

                    <div class="mb-4">
                        <label class="block text-sm font-medium text-gray-700 mb-2">Quanto planejamos gastar? (R$)</label>
                        <div class="relative">
                            <span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500">R$</span>
                            <input type="number" name="planned" min="1" step="1" required 
                                   class="w-full p-3 pl-10 border-2 border-gray-200 rounded-xl text-lg focus:border-purple-500 focus:outline-none"
                                   placeholder="0,00">
                        </div>
                    </div>

                    <div class="mb-6">
                        <label class="block text-sm font-medium text-gray-700 mb-2">Observação (opcional)</label>
                        <textarea name="description" rows="2"
                               class="w-full p-3 border-2 border-gray-200 rounded-xl text-lg focus:border-purple-500 focus:outline-none"
                               placeholder="Ex: Compra de cera de abelha para velas"></textarea>
                    </div>

                    <button type="submit" 
                            class="w-full bg-purple-500 hover:bg-purple-600 text-white text-xl font-bold py-4 rounded-xl shadow-lg transition transform active:scale-95">
                        SALVAR PLANEJAMENTO
                    </button>
                </form>
            </div>

            <!-- Relatório de Execução -->
            <div class="bg-white rounded-2xl shadow-lg p-6">
                <h2 class="text-xl font-bold text-gray-800 mb-6">Como está nosso planejamento?</h2>
                
                <div id="budget-report">
                    {{if .Executions}}
                    <div class="space-y-4">
                        {{range .Executions}}
                        <div class="border border-gray-200 rounded-xl p-4 hover:bg-gray-50 transition">
                            <div class="flex justify-between items-start mb-2">
                                <div>
                                    <h3 class="font-semibold text-gray-800">{{.Category}}</h3>
                                    {{if .Description}}
                                    <p class="text-gray-600 text-sm mt-1">{{.Description}}</p>
                                    {{end}}
                                </div>
                                <button onclick="deletePlan('{{.PlanID}}')" 
                                        class="text-red-500 hover:text-red-700 text-sm">
                                    Remover
                                </button>
                            </div>

                            <div class="mb-3">
                                <div class="flex justify-between text-sm text-gray-600 mb-1">
                                    <span>Planejado: {{formatCurrency .Planned}}</span>
                                    <span>Gasto: {{formatCurrency .Executed}}</span>
                                </div>
                                <div class="w-full bg-gray-200 rounded-full h-2">
                                    <div class="{{getProgressColor .Percentage}} h-2 rounded-full" 
                                         style="width: {{if le .Percentage 100}}{{.Percentage}}{{else}}100{{end}}%"></div>
                                </div>
                                <div class="flex justify-between text-xs text-gray-500 mt-1">
                                    <span>0%</span>
                                    <span>{{.Percentage}}%</span>
                                    <span>100%</span>
                                </div>
                            </div>

                            <div class="flex justify-between items-center">
                                <span class="text-sm {{if gt .Remaining 0}}text-green-600{{else}}text-red-600{{end}}">
                                    {{if gt .Remaining 0}}
                                    Falta: {{formatCurrency .Remaining}}
                                    {{else}}
                                    Ultrapassou: {{formatCurrency (multiply -1 .Remaining)}}
                                    {{end}}
                                </span>
                                <span class="text-xs px-2 py-1 rounded-full {{getStatusColor .AlertStatus}}">
                                    {{getStatusLabel .AlertStatus}}
                                </span>
                            </div>
                        </div>
                        {{end}}
                    </div>
                    {{else}}
                    <div class="text-center py-8 text-gray-500">
                        <svg class="w-12 h-12 mx-auto text-gray-300 mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
                        </svg>
                        <p>Nenhum planejamento registrado ainda.</p>
                        <p class="text-sm mt-1">Cadastre seu primeiro planejamento usando o formulário ao lado.</p>
                    </div>
                    {{end}}
                </div>
            </div>
        </div>

        <div class="text-center text-gray-500 text-sm mt-6">
            <p>O sistema acompanha automaticamente seus gastos reais e alerta quando estiver perto do limite planejado.</p>
        </div>
    </main>
</body>
</html>
{{end}}

{{define "budget_report.html"}}
<div class="space-y-4">
    {{if .Executions}}
        {{range .Executions}}
        <div class="border border-gray-200 rounded-xl p-4 hover:bg-gray-50 transition">
            <div class="flex justify-between items-start mb-2">
                <div>
                    <h3 class="font-semibold text-gray-800">{{.Category}}</h3>
                    {{if .Description}}
                    <p class="text-gray-600 text-sm mt-1">{{.Description}}</p>
                    {{end}}
                </div>
                <button onclick="deletePlan('{{.PlanID}}')" 
                        class="text-red-500 hover:text-red-700 text-sm">
                    Remover
                </button>
            </div>

            <div class="mb-3">
                <div class="flex justify-between text-sm text-gray-600 mb-1">
                    <span>Planejado: {{formatCurrency .Planned}}</span>
                    <span>Gasto: {{formatCurrency .Executed}}</span>
                </div>
                <div class="w-full bg-gray-200 rounded-full h-2">
                    <div class="{{getProgressColor .Percentage}} h-2 rounded-full" 
                         style="width: {{if le .Percentage 100}}{{.Percentage}}{{else}}100{{end}}%"></div>
                </div>
                <div class="flex justify-between text-xs text-gray-500 mt-1">
                    <span>0%</span>
                    <span>{{.Percentage}}%</span>
                    <span>100%</span>
                </div>
            </div>

            <div class="flex justify-between items-center">
                <span class="text-sm {{if gt .Remaining 0}}text-green-600{{else}}text-red-600{{end}}">
                    {{if gt .Remaining 0}}
                    Falta: {{formatCurrency .Remaining}}
                    {{else}}
                    Ultrapassou: {{formatCurrency (multiply -1 .Remaining)}}
                    {{end}}
                </span>
                <span class="text-xs px-2 py-1 rounded-full {{getStatusColor .AlertStatus}}">
                    {{getStatusLabel .AlertStatus}}
                </span>
            </div>
        </div>
        {{end}}
    {{else}}
    <div class="text-center py-8 text-gray-500">
        <svg class="w-12 h-12 mx-auto text-gray-300 mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
        </svg>
        <p>Nenhum planejamento registrado para {{.Period}}.</p>
    </div>
    {{end}}
</div>
{{end}}`
