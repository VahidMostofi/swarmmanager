
# # #FINDING DEMANDS FOR STAR
# go run main.go autoconfig \
#     --appname muck_star \
#     --config configurations/muck_star.yaml \
#     --workload  1_1300_0.2_0.2_0.2_0.2_0.2 \
#     --testName demands \
#     demands \
#     --duration 1200 \
#     --resultpath /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star_demands.yaml

# # #CUI mean 50%
# go run main.go autoconfig \
#     --appname muck_star \
#     --config configurations/muck_star.yaml \
#     --workload  50_110_0.3_0.2_0.1_0.1_0.3 \
#     --testName cui-250-mean-50 \
#     cui \
#     --cpuStat  CPUUsageMean \
#     --cpuThreshold 50

# # #CUI 90 pecentile of CPU Utilization 50%
# go run main.go autoconfig \
#     --appname muck_star \
#     --config configurations/muck_star.yaml \
#     --workload  50_110_0.3_0.2_0.1_0.1_0.3 \
#     --testName cui-250-90-50 \
#     cui \
#     --cpuStat CPUUsage90Percentile  \
#     --cpuThreshold 50

# #PPEU ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5
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

# # #PPEU ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 1.0
# go run main.go autoconfig \
#     --appname muck_star \
#     --config configurations/muck_star.yaml \
#     --workload  50_110_0.3_0.2_0.1_0.1_0.3 \
#     --testName ppeu-250-1.0-mc \
#     ppeu \
#     --property RTToleranceIntervalUBoundc90p95  \
#     --value 250 \
#     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star_demands.yaml \
#     --mc \
#     --stepsize 1.0

# # #BNV1 Bottleneck Version 1, ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5
# go run main.go autoconfig \
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

# # #BNV1 Bottleneck Version 1, ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 1.0
# go run main.go autoconfig \
#     --appname muck_star \
#     --config configurations/muck_star.yaml \
#     --workload  50_110_0.3_0.2_0.1_0.1_0.3 \
#     --testName bnv1-250-1.0-mc \
#     bnv1 \
#     --property RTToleranceIntervalUBoundc90p95  \
#     --value 250 \
#     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star_demands.yaml \
#     --mc \
#     --stepsize 1.0

# #BNV1 Bottleneck Version 1, ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 2.0
# go run main.go autoconfig \
#     --appname muck_star \
#     --config configurations/muck_star.yaml \
#     --workload  50_110_0.3_0.2_0.1_0.1_0.3 \
#     --testName bnv1-250-2.0-mc \
#     bnv1 \
#     --property RTToleranceIntervalUBoundc90p95  \
#     --value 250 \
#     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star_demands.yaml \
#     --mc \
#     --stepsize 2.0

# # #BNV1 Bottleneck Version 1, ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 3.0
# go run main.go autoconfig \
#     --appname muck_star \
#     --config configurations/muck_star.yaml \
#     --workload  50_110_0.3_0.2_0.1_0.1_0.3 \
#     --testName bnv1-250-3.0-mc \
#     bnv1 \
#     --property RTToleranceIntervalUBoundc90p95  \
#     --value 250 \
#     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star_demands.yaml \
#     --mc \
#     --stepsize 3.0

# #BNV2 Bottleneck Version 2, ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5
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

# # #BNV2 Bottleneck Version 2, ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 1.0
# go run main.go autoconfig \
#     --appname muck_star \
#     --config configurations/muck_star.yaml \
#     --workload  50_110_0.3_0.2_0.1_0.1_0.3 \
#     --testName bnv2-250-1.0-mc \
#     bnv2 \
#     --property RTToleranceIntervalUBoundc90p95  \
#     --value 250 \
#     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star_demands.yaml \
#     --mc \
#     --stepsize 1.0

# # #PPAU ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5
# go run main.go autoconfig \
#     --appname muck_star \
#     --config configurations/muck_star.yaml \
#     --workload  50_110_0.3_0.2_0.1_0.1_0.3 \
#     --testName ppau-250-0.5-mc \
#     ppau \
#     --property RTToleranceIntervalUBoundc90p95  \
#     --value 250 \
#     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star_demands.yaml \
#     --mc \
#     --stepsize 0.5

# # #PPAU ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 1.0
# go run main.go autoconfig \
#     --appname muck_star \
#     --config configurations/muck_star.yaml \
#     --workload  50_110_0.3_0.2_0.1_0.1_0.3 \
#     --testName ppau-250-1.0-mc \
#     ppau \
#     --property RTToleranceIntervalUBoundc90p95  \
#     --value 250 \
#     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star_demands.yaml \
#     --mc \
#     --stepsize 1.0

# VUSs=(30 40 60 70 80 90 )
VUSs=(20 )
for VUS in "${VUSs[@]}"
do
    #PPAU ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5
    go run main.go autoconfig \
        --appname muck_star \
        --config configurations/muck_star.yaml \
        --workload  "${VUS}_110_0.3_0.2_0.1_0.1_0.3" \
        --testName ppau-250-0.5-mc \
        ppau \
        --property RTToleranceIntervalUBoundc90p95  \
        --value 250 \
        --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star_demands.yaml \
        --mc \
        --stepsize 0.5
    
    # #PPAU ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 1.0
    # go run main.go autoconfig \
    #     --appname muck_star \
    #     --config configurations/muck_star.yaml \
    #     --workload  "${VUS}_110_0.3_0.2_0.1_0.1_0.3" \
    #     --testName ppau-250-1.0-mc \
    #     ppau \
    #     --property RTToleranceIntervalUBoundc90p95  \
    #     --value 250 \
    #     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star_demands.yaml \
    #     --mc \
    #     --stepsize 1.0
done