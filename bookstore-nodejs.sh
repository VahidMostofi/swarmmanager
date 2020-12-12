#!/bin/bash
VUSs=(75 100 125 150 )
for VUS in "${VUSs[@]}"
do
    # # BNV2, 95 percentile of respones time must be less thatn 250ms with stepsize = 2.0
    # go run main.go autoconfig \
    #     --appname bookstore_nodejs \
    #     --config configurations/bookstore_nodejs.yaml \
    #     --workload  "${VUS}_110_0.33_0.33_0.34" \
    #     --testName bnv2-250-2.0-mc-c-1.0 \
    #     bnv2 \
    #     --property ResponseTimes95Percentile  \
    #     --value 250 \
    #     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
    #     --mc \
    #     --stepsize 2.0 \
    #     --constantinit 1.0

    # # BNV2, 95 percentile of respones time must be less thatn 250ms with stepsize = 2.0
    # go run main.go autoconfig \
    #     --appname bookstore_nodejs \
    #     --config configurations/bookstore_nodejs.yaml \
    #     --workload  "${VUS}_110_0.33_0.33_0.34" \
    #     --testName bnv2-250-2.0-mc-c-0.5 \
    #     bnv2 \
    #     --property ResponseTimes95Percentile  \
    #     --value 250 \
    #     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
    #     --mc \
    #     --stepsize 2.0 \
    #     --constantinit 0.5

    # # BNV1, 95 percentile of respones time must be less thatn 250ms with stepsize = 1.0
    # go run main.go autoconfig \
    #     --appname bookstore_nodejs \
    #     --config configurations/bookstore_nodejs.yaml \
    #     --workload  "${VUS}_110_0.33_0.33_0.34" \
    #     --testName bnv1-250-1.0-mc-c-1.0 \
    #     bnv1 \
    #     --property ResponseTimes95Percentile  \
    #     --value 250 \
    #     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
    #     --mc \
    #     --stepsize 1.0 \
    #     --constantinit 1.0

    # # BNV1, 95 percentile of respones time must be less thatn 250ms with stepsize = 1.0
    # go run main.go autoconfig \
    #     --appname bookstore_nodejs \
    #     --config configurations/bookstore_nodejs.yaml \
    #     --workload  "${VUS}_110_0.33_0.33_0.34" \
    #     --testName bnv1-250-1.0-mc-c-0.5 \
    #     bnv1 \
    #     --property ResponseTimes95Percentile  \
    #     --value 250 \
    #     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
    #     --mc \
    #     --stepsize 1.0 \
    #     --constantinit 0.5

    # # PythonPath is the path to python interpretor
    # PythonPath="/home/vahid/.virtualenvs/with-data/bin/python"

    # # ScriptPath is the path to python script
    # ScriptPath="$(pwd)/scripts/mobo_CPU_split_mc.py"

    # go run main.go autoconfig \
    #     --appname bookstore_nodejs \
    #     --config configurations/bookstore_nodejs.yaml \
    #     --workload  "${VUS}_110_0.33_0.33_0.34" \
    #     --testName mobo \
    #     mobo \
    #     --python $PythonPath \
    #     --script $ScriptPath

    # BNV2, 95 percentile of respones time must be less thatn 250ms with stepsize = 2.0
    go run main.go autoconfig \
        --appname bookstore_nodejs \
        --config configurations/bookstore_nodejs.yaml \
        --workload  "${VUS}_110_0.33_0.33_0.34" \
        --testName bnv2-250-1.0-mc-c-1.0 \
        bnv2 \
        --property ResponseTimes95Percentile  \
        --value 250 \
        --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
        --mc \
        --stepsize 1.0 \
        --constantinit 1.0


    # BNV1, 95 percentile of respones time must be less thatn 250ms with stepsize = 2.0
    go run main.go autoconfig \
        --appname bookstore_nodejs \
        --config configurations/bookstore_nodejs.yaml \
        --workload  "${VUS}_110_0.33_0.33_0.34" \
        --testName bnv1-250-2.0-mc-c-1.0 \
        bnv1 \
        --property ResponseTimes95Percentile  \
        --value 250 \
        --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
        --mc \
        --stepsize 2.0 \
        --constantinit 1.0


done
