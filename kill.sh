#!/bin/bash

# Portas a serem verificadas
ports=(9876 9877 9878)

# Loop através das portas
for port in "${ports[@]}"; do
    # Obtém os PIDs dos processos usando a porta
    pids=$(lsof -t -i :$port)
    
    # Verifica se os PIDs foram encontrados
    if [ -n "$pids" ]; then
        echo "Matando processo(s) usando a porta $port: PID(s) $pids"
        kill $pids
    else
        echo "Nenhum processo encontrado usando a porta $port"
    fi
done
