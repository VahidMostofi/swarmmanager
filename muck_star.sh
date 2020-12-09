#!/bin/bash
# VUSs=(10 20 30 40 )
VUSs=(  30  )
for VUS in "${VUSs[@]}"
do
    # # BNV2, ResponseTimes95Percentile must be less thatn 300ms with stepsize = 2.0
    # go run main.go autoconfig \
    #     --appname muck_star-small \
    #     --config configurations/muck_star.yaml \
    #     --workload  "${VUS}_110_0.45_0.25_0.15_0.15" \
    #     --testName bnv2-300-2.0-mc-c-0.5 \
    #     bnv2 \
    #     --property ResponseTimes95Percentile  \
    #     --value 300 \
    #     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star-small_demands.yaml \
    #     --mc \
    #     --stepsize 2.0

    # # BNV2, ResponseTimes95Percentile must be less thatn 300ms with stepsize = 1.0
    # go run main.go autoconfig \
    #     --appname muck_star-small \
    #     --config configurations/muck_star.yaml \
    #     --workload  "${VUS}_110_0.45_0.25_0.15_0.15" \
    #     --testName bnv2-300-1.0-mc-c-0.5 \
    #     bnv2 \
    #     --property ResponseTimes95Percentile  \
    #     --value 300 \
    #     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star-small_demands.yaml \
    #     --mc \
    #     --stepsize 1.0

    # BNV1, ResponseTimes95Percentile must be less thatn 300ms with stepsize = 2.0
    go run main.go autoconfig \
        --appname muck_star-small \
        --config configurations/muck_star.yaml \
        --workload  "${VUS}_110_0.45_0.25_0.15_0.15" \
        --testName bnv1-300-2.0-mc-c-0.5 \
        bnv1 \
        --property ResponseTimes95Percentile  \
        --value 300 \
        --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star-small_demands.yaml \
        --mc \
        --stepsize 2.0 \
        --constantinit 0.5

    # # BNV1, ResponseTimes95Percentile must be less thatn 300ms with stepsize = 1.0
    # go run main.go autoconfig \
    #     --appname muck_star-small \
    #     --config configurations/muck_star.yaml \
    #     --workload  "${VUS}_110_0.45_0.25_0.15_0.15" \
    #     --testName bnv1-300-1.0-mc-c-0.5 \
    #     bnv1 \
    #     --property ResponseTimes95Percentile  \
    #     --value 300 \
    #     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star-small_demands.yaml \
    #     --mc \
    #     --stepsize 1.0 \
    #     --constantinit 0.5
done
