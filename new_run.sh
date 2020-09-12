#!/bin/bash
# ######################################################
# go run main.go autoconfig \
#     --appname bookstore_nodejs \
#     --config configurations/bookstore_nodejs.yaml \
#     --workload  1_100_0.5_5 \
#     --testName adfc-350-1-mc \
#     adfc \
#     --property RTToleranceIntervalUBoundc90p95 \
#     --value 350 \
#     --stepsize 1 \
#     --mc
# "args": ["autoconfig", "--appname", "bookstore_nodejs", "--config", "configurations/bookstore_nodejs.yaml", "--workload","1_100_0.5_5", "--testName", "adfc-350-1-mc", "adfc", "--property", "RTToleranceIntervalUBoundc90p95", "--value", "350", "--stepsize", "1", "--mc"]
# ######################################################
go run main.go autoconfig \
    --appname bookstore_nodejs \
    --config configurations/bookstore_nodejs.yaml \
    --workload  250_80_0.5_10 \
    --testName cui-350-mean-60 \
    cui \
    --cpuStat  CPUUsageMean \
    --cpuThreshold 60
# "args": ["autoconfig", "--appname", "bookstore_nodejs", "--config", "configurations/bookstore_nodejs.yaml", "--workload","250_80_0.5_10", "--testName", "cui-350-mean-60", "cui", "--cpuStat", "CPUUsageMean", "--cpuThreshold", "60"]
# ######################################################