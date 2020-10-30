#!/bin/bash
VUSs=(75 100 125 150 175 )
for VUS in "${VUSs[@]}"
do
    # BNV2 Bottleneck Version 2, ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5
    go run main.go autoconfig \
        --appname bookstore_nodejs \
        --config configurations/bookstore_nodejs.yaml \
        --workload  "${VUS}_110_0.33_0.33_0.34" \
        --testName bnv2-250-2.0-mc-c-0.5 \
        bnv2 \
        --property RTToleranceIntervalUBoundc90p95  \
        --value 250 \
        --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
        --mc \
        --stepsize 2.0

    # BNV2 Bottleneck Version 2, ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5
    go run main.go autoconfig \
        --appname bookstore_nodejs \
        --config configurations/bookstore_nodejs.yaml \
        --workload  "${VUS}_110_0.33_0.33_0.34" \
        --testName bnv2-250-1.0-mc-c-0.5 \
        bnv2 \
        --property RTToleranceIntervalUBoundc90p95  \
        --value 250 \
        --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
        --mc \
        --stepsize 1.0

    # BNV2 Bottleneck Version 2, ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5
    go run main.go autoconfig \
        --appname bookstore_nodejs \
        --config configurations/bookstore_nodejs.yaml \
        --workload  "${VUS}_110_0.33_0.33_0.34" \
        --testName bnv2-250-0.5-mc-c-0.5 \
        bnv2 \
        --property RTToleranceIntervalUBoundc90p95  \
        --value 250 \
        --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
        --mc \
        --stepsize 0.5

    # BNV2 Bottleneck Version 2, ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5
    go run main.go autoconfig \
        --appname bookstore_nodejs \
        --config configurations/bookstore_nodejs.yaml \
        --workload  "${VUS}_110_0.33_0.33_0.34" \
        --testName bnv2-250-0.2-mc-c-0.5 \
        bnv2 \
        --property RTToleranceIntervalUBoundc90p95  \
        --value 250 \
        --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
        --mc \
        --stepsize 0.2

    # BottleNeck Versoin 1, ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5
    go run main.go autoconfig \
        --appname bookstore_nodejs \
        --config configurations/bookstore_nodejs.yaml \
        --workload  "${VUS}_110_0.33_0.33_0.34" \
        --testName bnv1-250-2.0-mc-c-0.5 \
        bnv1 \
        --property RTToleranceIntervalUBoundc90p95  \
        --value 250 \
        --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
        --mc \
        --stepsize 2.0 \
        --constantinit 0.5

    # BottleNeck Versoin 1, ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5
    go run main.go autoconfig \
        --appname bookstore_nodejs \
        --config configurations/bookstore_nodejs.yaml \
        --workload  "${VUS}_110_0.33_0.33_0.34" \
        --testName bnv1-250-0.2-mc-c-0.5 \
        bnv1 \
        --property RTToleranceIntervalUBoundc90p95  \
        --value 250 \
        --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
        --mc \
        --stepsize 0.2 \
        --constantinit 0.5

    # BottleNeck Versoin 1, ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5
    go run main.go autoconfig \
        --appname bookstore_nodejs \
        --config configurations/bookstore_nodejs.yaml \
        --workload  "${VUS}_110_0.33_0.33_0.34" \
        --testName bnv1-250-0.5-mc-c-0.5 \
        bnv1 \
        --property RTToleranceIntervalUBoundc90p95  \
        --value 250 \
        --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
        --mc \
        --stepsize 0.5 \
        --constantinit 0.5

    # BottleNeck Versoin 1, ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5
    go run main.go autoconfig \
        --appname bookstore_nodejs \
        --config configurations/bookstore_nodejs.yaml \
        --workload  "${VUS}_110_0.33_0.33_0.34" \
        --testName bnv1-250-1.0-mc-c-0.5 \
        bnv1 \
        --property RTToleranceIntervalUBoundc90p95  \
        --value 250 \
        --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
        --mc \
        --stepsize 1.0 \
        --constantinit 0.5

done


VUSs=(75 )
for VUS in "${VUSs[@]}"
do
    # BNV2 Bottleneck Version 2, ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5
    go run main.go autoconfig \
        --appname bookstore_nodejs \
        --config configurations/bookstore_nodejs.yaml \
        --workload  "${VUS}_110_0.33_0.33_0.34" \
        --testName bnv2-250-2.0-mc-c-0.5-r95 \
        bnv2 \
        --property ResponseTimes95Percentile  \
        --value 250 \
        --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
        --mc \
        --stepsize 2.0

    # BNV2 Bottleneck Version 2, ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5
    go run main.go autoconfig \
        --appname bookstore_nodejs \
        --config configurations/bookstore_nodejs.yaml \
        --workload  "${VUS}_110_0.33_0.33_0.34" \
        --testName bnv2-250-1.0-mc-c-0.5-r95 \
        bnv2 \
        --property ResponseTimes95Percentile  \
        --value 250 \
        --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
        --mc \
        --stepsize 1.0

    # BNV2 Bottleneck Version 2, ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5
    go run main.go autoconfig \
        --appname bookstore_nodejs \
        --config configurations/bookstore_nodejs.yaml \
        --workload  "${VUS}_110_0.33_0.33_0.34" \
        --testName bnv2-250-0.5-mc-c-0.5-r95 \
        bnv2 \
        --property ResponseTimes95Percentile  \
        --value 250 \
        --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
        --mc \
        --stepsize 0.5

    # BNV2 Bottleneck Version 2, ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5
    go run main.go autoconfig \
        --appname bookstore_nodejs \
        --config configurations/bookstore_nodejs.yaml \
        --workload  "${VUS}_110_0.33_0.33_0.34" \
        --testName bnv2-250-0.2-mc-c-0.5-r95 \
        bnv2 \
        --property ResponseTimes95Percentile  \
        --value 250 \
        --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
        --mc \
        --stepsize 0.2

    # BottleNeck Versoin 1, ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5
    go run main.go autoconfig \
        --appname bookstore_nodejs \
        --config configurations/bookstore_nodejs.yaml \
        --workload  "${VUS}_110_0.33_0.33_0.34" \
        --testName bnv1-250-2.0-mc-c-0.5-r95 \
        bnv1 \
        --property ResponseTimes95Percentile  \
        --value 250 \
        --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
        --mc \
        --stepsize 2.0 \
        --constantinit 0.5

    # BottleNeck Versoin 1, ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5
    go run main.go autoconfig \
        --appname bookstore_nodejs \
        --config configurations/bookstore_nodejs.yaml \
        --workload  "${VUS}_110_0.33_0.33_0.34" \
        --testName bnv1-250-0.2-mc-c-0.5-r95 \
        bnv1 \
        --property ResponseTimes95Percentile  \
        --value 250 \
        --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
        --mc \
        --stepsize 0.2 \
        --constantinit 0.5

    # BottleNeck Versoin 1, ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5
    go run main.go autoconfig \
        --appname bookstore_nodejs \
        --config configurations/bookstore_nodejs.yaml \
        --workload  "${VUS}_110_0.33_0.33_0.34" \
        --testName bnv1-250-0.5-mc-c-0.5-r95 \
        bnv1 \
        --property ResponseTimes95Percentile  \
        --value 250 \
        --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
        --mc \
        --stepsize 0.5 \
        --constantinit 0.5

    # BottleNeck Versoin 1, ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5
    go run main.go autoconfig \
        --appname bookstore_nodejs \
        --config configurations/bookstore_nodejs.yaml \
        --workload  "${VUS}_110_0.33_0.33_0.34" \
        --testName bnv1-250-1.0-mc-c-0.5-r95 \
        bnv1 \
        --property ResponseTimes95Percentile  \
        --value 250 \
        --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
        --mc \
        --stepsize 1.0 \
        --constantinit 0.5

done
