#!/bin/bash

# PythonPath is the path to python interpretor
PythonPath="/home/vahid/.virtualenvs/with-data/bin/python"

# ScriptPath is the path to python script
ScriptPath="$(pwd)/scripts/mobo_CPU_split_mc.py"

# VUSs=(75 )
# for VUS in "${VUSs[@]}"
# do
#     go run main.go autoconfig \
#         --appname bookstore_nodejs \
#         --config configurations/bookstore_nodejs.yaml \
#         --workload  "${VUS}_110_0.33_0.33_0.34" \
#         --testName mobo \
#         mobo \
#         --python $PythonPath \
#         --script $ScriptPath
# done


# ScriptPath is the path to python script
ScriptPath="$(pwd)/scripts/mobo_CPU_split_mc-star.py"
VUSs=(20 30 40 50 )
for VUS in "${VUSs[@]}"
do
    go run main.go autoconfig \
        --appname muck_star-small \
        --config configurations/muck_star.yaml \
        --workload  "${VUS}_110_0.45_0.25_0.15_0.15" \
        --testName mobo \
        mobo \
        --python $PythonPath \
        --script $ScriptPath
done