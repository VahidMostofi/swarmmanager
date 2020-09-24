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
    --workload  1_150_0.2_0.2_0.2_0.2_0.2 \
    --testName demands \
    demands \
    --duration 120 \
    --resultpath /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star_demands.yaml
#THESE ARE FOR STAR ARCHITECTURE
# go run main.go autoconfig \
#     --appname muck_star \
#     --config configurations/muck_star.yaml \
#     --workload  50_110_0.3_0.2_0.1_0.1_0.3 \
#     --testName cui-250-mean-50 \
#     cui \
#     --cpuStat  CPUUsageMean \
#     --cpuThreshold 50
# go run main.go autoconfig \
#     --appname muck_star \
#     --config configurations/muck_star.yaml \
#     --workload  50_110_0.3_0.2_0.1_0.1_0.3 \
#     --testName cui-250-90-50 \
#     cui \
#     --cpuStat CPUUsage90Percentile  \
#     --cpuThreshold 50
# go run main.go autoconfig \
#     --appname muck_star \
#     --config configurations/muck_star.yaml \
#     --workload  60_110_0.3_0.2_0.1_0.1_0.3 \
#     --testName cui-250-mean-50 \
#     cui \
#     --cpuStat  CPUUsageMean \
#     --cpuThreshold 50

# go run main.go autoconfig \
#     --appname muck_star \
#     --config configurations/muck_star.yaml \
#     --workload  60_110_0.3_0.2_0.1_0.1_0.3 \
#     --testName cui-250-90-50 \
#     cui \
#     --cpuStat CPUUsage90Percentile  \
#     --cpuThreshold 50

# go run main.go autoconfig \
#     --appname muck_star \
#     --config configurations/muck_star.yaml \
#     --workload  70_110_0.3_0.2_0.1_0.1_0.3 \
#     --testName cui-250-mean-50 \
#     cui \
#     --cpuStat  CPUUsageMean \
#     --cpuThreshold 50

# go run main.go autoconfig \
#     --appname muck_star \
#     --config configurations/muck_star.yaml \
#     --workload  70_110_0.3_0.2_0.1_0.1_0.3 \
#     --testName cui-250-90-50 \
#     cui \
#     --cpuStat CPUUsage90Percentile  \
#     --cpuThreshold 50

# go run main.go autoconfig \
#     --appname muck_star \
#     --config configurations/muck_star.yaml \
#     --workload  80_110_0.3_0.2_0.1_0.1_0.3 \
#     --testName cui-250-mean-50 \
#     cui \
#     --cpuStat  CPUUsageMean \
#     --cpuThreshold 50

# go run main.go autoconfig \
#     --appname muck_star \
#     --config configurations/muck_star.yaml \
#     --workload  80_110_0.3_0.2_0.1_0.1_0.3 \
#     --testName cui-250-90-50 \
#     cui \
#     --cpuStat CPUUsage90Percentile  \
#     --cpuThreshold 50

# go run main.go autoconfig \
#     --appname muck_star \
#     --config configurations/muck_star.yaml \
#     --workload  90_110_0.3_0.2_0.1_0.1_0.3 \
#     --testName cui-250-mean-50 \
#     cui \
#     --cpuStat  CPUUsageMean \
#     --cpuThreshold 50

# go run main.go autoconfig \
#     --appname muck_star \
#     --config configurations/muck_star.yaml \
#     --workload  90_110_0.3_0.2_0.1_0.1_0.3 \
#     --testName cui-250-90-50 \
#     cui \
#     --cpuStat CPUUsage90Percentile  \
#     --cpuThreshold 50
# #################################################
# #FINDING the configuration using PPEU approach with stepSize = 0.5
# go run main.go autoconfig \
#     --appname muck_star \
#     --config configurations/muck_star.yaml \
#     --workload  50_110_0.3_0.2_0.1_0.1_0.3 \
#     --testName ppeu-250-0.5-mc \
#     ppeu \
#     --property RTToleranceIntervalUBoundc90p95  \
#     --value 250 \
#     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star_demands.yaml \
#     --mc \
#     --stepsize 0.5

# go run main.go autoconfig \
#     --appname muck_star \
#     --config configurations/muck_star.yaml \
#     --workload  60_110_0.3_0.2_0.1_0.1_0.3 \
#     --testName ppeu-250-0.5-mc \
#     ppeu \
#     --property RTToleranceIntervalUBoundc90p95  \
#     --value 250 \
#     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star_demands.yaml \
#     --mc \
#     --stepsize 0.5

# go run main.go autoconfig \
#     --appname muck_star \
#     --config configurations/muck_star.yaml \
#     --workload  70_110_0.3_0.2_0.1_0.1_0.3 \
#     --testName ppeu-250-0.5-mc \
#     ppeu \
#     --property RTToleranceIntervalUBoundc90p95  \
#     --value 250 \
#     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star_demands.yaml \
#     --mc \
#     --stepsize 0.5

# go run main.go autoconfig \
#     --appname muck_star \
#     --config configurations/muck_star.yaml \
#     --workload  80_110_0.3_0.2_0.1_0.1_0.3 \
#     --testName ppeu-250-0.5-mc \
#     ppeu \
#     --property RTToleranceIntervalUBoundc90p95  \
#     --value 250 \
#     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star_demands.yaml \
#     --mc \
#     --stepsize 0.5

# go run main.go autoconfig \
#     --appname muck_star \
#     --config configurations/muck_star.yaml \
#     --workload  90_110_0.3_0.2_0.1_0.1_0.3 \
#     --testName ppeu-250-0.5-mc \
#     ppeu \
#     --property RTToleranceIntervalUBoundc90p95  \
#     --value 250 \
#     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star_demands.yaml \
#     --mc \
#     --stepsize 0.5


# go run main.go autoconfig \ I don't want to use this!!!
#     --appname muck_star \
#     --config configurations/muck_star.yaml \
#     --workload  50_110_0.3_0.2_0.1_0.1_0.3 \
#     --testName bnv1-250-0.5-mc \
#     bnv1 \
#     --property RTToleranceIntervalUBoundc90p95  \
#     --value 250 \
#     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star_demands.yaml \
#     --mc \
#     --stepsize 0.5

# go run main.go autoconfig \
#     --appname muck_star \
#     --config configurations/muck_star.yaml \
#     --workload  50_110_0.3_0.2_0.1_0.1_0.3 \
#     --testName bnv2-250-0.5-mc \
#     bnv2 \
#     --property RTToleranceIntervalUBoundc90p95  \
#     --value 250 \
#     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star_demands.yaml \
#     --mc \
#     --stepsize 0.5

# go run main.go autoconfig \
#     --appname muck_star \
#     --config configurations/muck_star.yaml \
#     --workload  50_110_0.3_0.2_0.1_0.1_0.3 \
#     --testName bnv2-250-1-mc \
#     bnv2 \
#     --property RTToleranceIntervalUBoundc90p95  \
#     --value 250 \
#     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star_demands.yaml \
#     --mc \
#     --stepsize 1

# go run main.go autoconfig \
#     --appname muck_star \
#     --config configurations/muck_star.yaml \
#     --workload  50_110_0.3_0.2_0.1_0.1_0.3 \
#     --testName bnv1-200-1-mc \
#     bnv1 \
#     --property RTToleranceIntervalUBoundc90p95  \
#     --value 200 \
#     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star_demands.yaml \
#     --mc \
#     --stepsize 1

# go run main.go autoconfig \ 
#     --appname muck_star \
#     --config configurations/muck_star.yaml \
#     --workload  50_110_0.3_0.2_0.1_0.1_0.3 \
#     --testName bnv1-250-2-mc \
#     bnv1 \
#     --property RTToleranceIntervalUBoundc90p95  \
#     --value 250 \
#     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star_demands.yaml \
#     --mc \
#     --stepsize 2

VUSs=(30 40 50 60 70 80 90 )
for VUS in "${VUSs[@]}"
do
    go run main.go autoconfig \ 
        --appname muck_star \
        --config configurations/muck_star.yaml \
        --workload  "${VUS}_110_0.3_0.2_0.1_0.1_0.3" \
        --testName bnv1-250-1-mc \
        bnv1 \
        --property RTToleranceIntervalUBoundc90p95  \
        --value 250 \
        --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star_demands.yaml \
        --mc \
        --stepsize 1

    go run main.go autoconfig \ 
        --appname muck_star \
        --config configurations/muck_star.yaml \
        --workload  "${VUS}_110_0.3_0.2_0.1_0.1_0.3" \
        --testName bnv1-250-2-mc \
        bnv1 \
        --property RTToleranceIntervalUBoundc90p95  \
        --value 250 \
        --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star_demands.yaml \
        --mc \
        --stepsize 2

done