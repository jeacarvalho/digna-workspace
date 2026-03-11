#!/bin/bash
find_latest_active_task() {
    if [ ! -d "work_in_progress/tasks" ]; then
        echo "Diretório não existe"
        return 1
    fi
    
    shopt -s nullglob
    local task_dirs=(work_in_progress/tasks/task_*)
    shopt -u nullglob
    
    echo "Número de diretórios: ${#task_dirs[@]}"
    
    if [ ${#task_dirs[@]} -eq 0 ]; then
        echo "Nenhum diretório encontrado"
        return 1
    fi
    
    echo "Primeiro diretório: ${task_dirs[0]}"
    return 0
}

echo "Testando..."
find_latest_active_task
