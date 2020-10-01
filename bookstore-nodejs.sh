# FINDING DEMANDS FOR STAR
# go run main.go autoconfig \
#     --appname bookstore_nodejs \
#     --config configurations/bookstore_nodejs.yaml \
#     --workload  1_1300_0.33_0.33_0.34 \
#     --testName demands \
#     demands \
#     --duration 1200 \
#     --resultpath /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml

# VUSs=(100 125 150 175 200 )
VUSs=(100 125 175 200 )
# VUSs=(150 )
for VUS in "${VUSs[@]}"
do
    # # CUI mean 50%
    # go run main.go autoconfig \
    #     --appname bookstore_nodejs \
    #     --config configurations/bookstore_nodejs.yaml \
    #     --workload  "${VUS}_110_0.33_0.33_0.34" \
    #     --testName cui-250-mean-50 \
    #     cui \
    #     --cpuStat  CPUUsageMean \
    #     --cpuThreshold 50

    # # CUI 90 pecentile of CPU Utilization 50%
    # go run main.go autoconfig \
    #     --appname bookstore_nodejs \
    #     --config configurations/bookstore_nodejs.yaml \
    #     --workload  "${VUS}_110_0.33_0.33_0.34" \
    #     --testName cui-250-90-50 \
    #     cui \
    #     --cpuStat CPUUsage90Percentile  \
    #     --cpuThreshold 50

    # # PPAU ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5 with estimated Utilization init
    # go run main.go autoconfig \
    #     --appname bookstore_nodejs \
    #     --config configurations/bookstore_nodejs.yaml \
    #     --workload  "${VUS}_110_0.33_0.33_0.34" \
    #     --testName ppau-250-0.5-mc \
    #     ppau \
    #     --property RTToleranceIntervalUBoundc90p95  \
    #     --value 250 \
    #     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
    #     --mc \
    #     --stepsize 0.5

    # # PPAU ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5 with constant init
    # go run main.go autoconfig \
    #     --appname bookstore_nodejs \
    #     --config configurations/bookstore_nodejs.yaml \
    #     --workload  "${VUS}_110_0.33_0.33_0.34" \
    #     --testName ppau-250-0.5-mc-c-0.5 \
    #     ppau \
    #     --property RTToleranceIntervalUBoundc90p95  \
    #     --value 250 \
    #     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
    #     --mc \
    #     --stepsize 0.5 \
    #     --constantinit 0.5

    # # PPE ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5 with constant init
    # go run main.go autoconfig \
    #     --appname bookstore_nodejs \
    #     --config configurations/bookstore_nodejs.yaml \
    #     --workload  "${VUS}_110_0.33_0.33_0.34" \
    #     --testName ppe-250-0.5-mc-c-0.5 \
    #     ppe \
    #     --property RTToleranceIntervalUBoundc90p95  \
    #     --value 250 \
    #     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
    #     --mc \
    #     --stepsize 0.5 \
    #     --constantinit 0.5

    # # PPE ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5 with estimated Utilization init
    # go run main.go autoconfig \
    #     --appname bookstore_nodejs \
    #     --config configurations/bookstore_nodejs.yaml \
    #     --workload  "${VUS}_110_0.33_0.33_0.34" \
    #     --testName ppe-250-0.5-mc \
    #     ppe \
    #     --property RTToleranceIntervalUBoundc90p95  \
    #     --value 250 \
    #     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
    #     --mc \
    #     --stepsize 0.5

    # PPAU ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 1.0 with estimated Utilization init
    go run main.go autoconfig \
        --appname bookstore_nodejs \
        --config configurations/bookstore_nodejs.yaml \
        --workload  "${VUS}_110_0.33_0.33_0.34" \
        --testName ppau-250-1.0-mc \
        ppau \
        --property RTToleranceIntervalUBoundc90p95  \
        --value 250 \
        --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
        --mc \
        --stepsize 1.0

    # PPAU ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 1.0 with constant init
    go run main.go autoconfig \
        --appname bookstore_nodejs \
        --config configurations/bookstore_nodejs.yaml \
        --workload  "${VUS}_110_0.33_0.33_0.34" \
        --testName ppau-250-1.0-mc-c-0.5 \
        ppau \
        --property RTToleranceIntervalUBoundc90p95  \
        --value 250 \
        --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
        --mc \
        --stepsize 1.0 \
        --constantinit 0.5

    # PPE ToleranceIntervalc90p95 must be less thatn 250ms with stepsize =1.0 with constant init
    go run main.go autoconfig \
        --appname bookstore_nodejs \
        --config configurations/bookstore_nodejs.yaml \
        --workload  "${VUS}_110_0.33_0.33_0.34" \
        --testName ppe-250-1.0-mc-c-0.5 \
        ppe \
        --property RTToleranceIntervalUBoundc90p95  \
        --value 250 \
        --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
        --mc \
        --stepsize 1.0 \
        --constantinit 0.5

    # PPE ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 1.0 with estimated Utilization init
    go run main.go autoconfig \
        --appname bookstore_nodejs \
        --config configurations/bookstore_nodejs.yaml \
        --workload  "${VUS}_110_0.33_0.33_0.34" \
        --testName ppe-250-1.0-mc \
        ppe \
        --property RTToleranceIntervalUBoundc90p95  \
        --value 250 \
        --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
        --mc \
        --stepsize 1.0

    # # BottleNeck Versoin 1, ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5
    # go run main.go autoconfig \
    #     --appname bookstore_nodejs \
    #     --config configurations/bookstore_nodejs.yaml \
    #     --workload  "${VUS}_110_0.33_0.33_0.34" \
    #     --testName bnv1-250-0.5-mc \
    #     bnv1 \
    #     --property RTToleranceIntervalUBoundc90p95  \
    #     --value 250 \
    #     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
    #     --mc \
    #     --stepsize 0.5

    # # BNV1 Bottleneck Version 1, ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 1.0
    # go run main.go autoconfig \
    #     --appname bookstore_nodejs \
    #     --config configurations/bookstore_nodejs.yaml \
    #     --workload  "${VUS}_110_0.33_0.33_0.34" \
    #     --testName bnv1-250-1.0-mc \
    #     bnv1 \
    #     --property RTToleranceIntervalUBoundc90p95  \
    #     --value 250 \
    #     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
    #     --mc \
    #     --stepsize 1.0
done