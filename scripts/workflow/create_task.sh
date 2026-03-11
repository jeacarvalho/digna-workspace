#!/bin/bash
# create_task.sh - Cria uma nova tarefa no sistema
# Uso: ./create_task.sh "Nome da Tarefa" [opcional: módulo]

set -e

echo "📝 Criando nova tarefa..."
echo "=========================="

# Verificar se há sessão ativa
if [ ! -d "work_in_progress/current_session" ]; then
    echo "❌ Nenhuma sessão ativa encontrada."
    echo "💡 Execute primeiro: ./start_session.sh"
    exit 1
fi

# Parâmetros
TASK_NAME="${1}"
MODULE="${2:-ui_web}"
TASK_ID=$(date +%Y%m%d_%H%M%S)
TASK_DIR="work_in_progress/tasks/task_${TASK_ID}"

# Validar nome da tarefa
if [ -z "$TASK_NAME" ]; then
    echo "❌ Nome da tarefa é obrigatório."
    echo "💡 Uso: ./create_task.sh \"Nome da Tarefa\" [módulo]"
    exit 1
fi

# Criar diretório da tarefa
mkdir -p "${TASK_DIR}"
mkdir -p "${TASK_DIR}/task_artifacts"

echo "✅ Diretório criado: ${TASK_DIR}"

# Copiar template de prompt
if [ -f "work_in_progress/task_template/task_prompt.md" ]; then
    cp "work_in_progress/task_template/task_prompt.md" "${TASK_DIR}/task_prompt.md"
    
	# Personalizar template
	sed -i "s|\[NOME_DA_TAREFA\]|${TASK_NAME}|g" "${TASK_DIR}/task_prompt.md"
	sed -i "s|\[DATA\]|$(date +%d/%m/%Y)|g" "${TASK_DIR}/task_prompt.md"
	sed -i "s|\[módulo1, módulo2, \.\.\.\]|${MODULE}|g" "${TASK_DIR}/task_prompt.md"
    
    echo "✅ Template de prompt criado: ${TASK_DIR}/task_prompt.md"
else
    echo "⚠️  Template não encontrado, criando prompt básico..."
    
    cat > "${TASK_DIR}/task_prompt.md" << EOF
# 📋 TAREFA: ${TASK_NAME}

**Data:** $(date +%d/%m/%Y)
**Módulo:** ${MODULE}

---

## 🎯 OBJETIVO

[Descreva o objetivo da tarefa]

---

## 📋 REQUISITOS

### Funcionais
- [ ] [Requisito 1]

### Técnicos  
- [ ] Seguir padrões do projeto Digna
- [ ] Implementar testes unitários
- [ ] Atualizar documentação

---

**Status:** PENDENTE
**Última atualização:** $(date +%d/%m/%Y)
EOF
fi

# Criar arquivo de metadados da tarefa
cat > "${TASK_DIR}/task_metadata" << EOF
TASK_ID=${TASK_ID}
TASK_NAME="${TASK_NAME}"
MODULE="${MODULE}"
CREATED_AT=$(date +%s)
STATUS="pending"
EOF

echo ""
echo "✅ TAREFA CRIADA COM SUCESSO!"
echo "=============================="
echo "ID: ${TASK_ID}"
echo "Nome: ${TASK_NAME}"
echo "Módulo: ${MODULE}"
echo "Diretório: ${TASK_DIR}"
echo ""
echo "📋 PRÓXIMOS PASSOS:"
echo "1. Editar o prompt: vim ${TASK_DIR}/task_prompt.md"
echo "2. Iniciar processamento: ./process_task.sh --task=${TASK_ID} --checklist"
echo "3. Ou processar diretamente: ./process_task.sh --task=${TASK_ID} --execute"
echo ""
echo "💡 Dica: Use 'ls work_in_progress/tasks/' para listar todas as tarefas."