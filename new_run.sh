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
# go run main.go autoconfig \
#     --appname bookstore_nodejs \
#     --config configurations/bookstore_nodejs.yaml \
#     --workload  250_90_0.5_10 \
#     --testName cui-350-mean-60 \
#     cui \
#     --cpuStat  CPUUsageMean \
#     --cpuThreshold 50
# "args": ["autoconfig", "--appname", "bookstore_nodejs", "--config", "configurations/bookstore_nodejs.yaml", "--workload","250_80_0.5_10", "--testName", "cui-350-mean-60", "cui", "--cpuStat", "CPUUsageMean", "--cpuThreshold", "60"]
# ######################################################
# go run main.go autoconfig \
#     --appname muck_two_layers \
#     --config configurations/muck_two_layers.yaml \
#     --workload  100_110_0.2_0.2_0.2_0.2_0.2 \
#     --testName cui-350-mean-50 \
#     cui \
#     --cpuStat  CPUUsageMean \
#     --cpuThreshold 50
# ######################################################
# #FINDING DEMANDS FOR STAR
go run main.go autoconfig \
    --appname muck_star \
    --config configurations/muck_star.yaml \
    --workload  1_1400_0.2_0.2_0.2_0.2_0.2 \
    --testName demands \
    demands \
    --duration 1200 \
    --resultpath /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star_demands.yaml
# #THESE ARE FOR STAR ARCHITECTURE
# go run main.go autoconfig \
#     --appname muck_star \
#     --config configurations/muck_star.yaml \
#     --workload  50_110_0.3_0.2_0.1_0.1_0.3 \
#     --testName cui-350-mean-50 \
#     cui \
#     --cpuStat  CPUUsageMean \
#     --cpuThreshold 50
# go run main.go autoconfig \
#     --appname muck_star \
#     --config configurations/muck_star.yaml \
#     --workload  60_110_0.3_0.2_0.1_0.1_0.3 \
#     --testName cui-350-mean-50 \
#     cui \
#     --cpuStat  CPUUsageMean \
#     --cpuThreshold 50
# go run main.go autoconfig \
#     --appname muck_star \
#     --config configurations/muck_star.yaml \
#     --workload  70_110_0.3_0.2_0.1_0.1_0.3 \
#     --testName cui-350-mean-50 \
#     cui \
#     --cpuStat  CPUUsageMean \
#     --cpuThreshold 50
# go run main.go autoconfig \
#     --appname muck_star \
#     --config configurations/muck_star.yaml \
#     --workload  80_110_0.3_0.2_0.1_0.1_0.3 \
#     --testName cui-350-mean-50 \
#     cui \
#     --cpuStat  CPUUsageMean \
#     --cpuThreshold 50
# go run main.go autoconfig \
#     --appname muck_star \
#     --config configurations/muck_star.yaml \
#     --workload  90_110_0.3_0.2_0.1_0.1_0.3 \
#     --testName cui-350-mean-50 \
#     cui \
#     --cpuStat  CPUUsageMean \
#     --cpuThreshold 50
# #################################################
