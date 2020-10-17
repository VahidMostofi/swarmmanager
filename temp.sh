# # CUI 90 pecentile of CPU Utilization 50%
# go run main.go autoconfig \
#     --appname bookstore_nodejs \
#     --config configurations/bookstore_nodejs.yaml \
#     --workload  "150_3000_0.33_0.33_0.34" \
#     --testName cui-250-90-50 \
#     cui \
#     --cpuStat CPUUsage90Percentile  \
#     --cpuThreshold 50

# # PPAU ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5 with estimated Utilization init
# go run main.go autoconfig \
#     --appname bookstore_nodejs \
#     --config configurations/bookstore_nodejs.yaml \
#     --workload  "150_30000_0.33_0.33_0.34" \
#     --testName ppau-250-0.5-mc-liveupdate \
#     ppau \
#     --property RTToleranceIntervalUBoundc90p95  \
#     --value 250 \
#     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
#     --mc \
#     --stepsize 0.5 \
#      --constantinit 0.5
# while true
# do
# go run main.go autoconfig \
#     --appname bookstore_nodejs \
#     --config configurations/bookstore_nodejs.yaml \
#     --workload  "125_60_0.33_0.33_0.34" \
#     --testName brute-force \
#     brute

# sleep 1
# done

# go run main.go autoconfig \
#     --appname bookstore_nodejs \
#     --config configurations/bookstore_nodejs.yaml \
#     --workload  "125_110_0.33_0.33_0.34" \
#     --testName signle_6.5_1.5_2.5 \
#     single entry 2.5 auth 6.5 books 1.5


# #BNV2 Bottleneck Version 2, ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5
# go run main.go autoconfig \
#     --appname bookstore_nodejs \
#     --config configurations/bookstore_nodejs.yaml \
#     --workload  75_110_0.33_0.33_0.34 \
#     --testName bnv2-250-2.0-mc-c-1.0 \
#     bnv2 \
#     --property RTToleranceIntervalUBoundc90p95  \
#     --value 250 \
#     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
#     --mc \
#     --stepsize 2.0

# # BottleNeck Versoin 1, ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5
# go run main.go autoconfig \
#     --appname bookstore_nodejs \
#     --config configurations/bookstore_nodejs.yaml \
#     --workload  75_110_0.33_0.33_0.34 \
#     --testName bnv1-250-0.2-mc-c-1.0 \
#     bnv1 \
#     --property RTToleranceIntervalUBoundc90p95  \
#     --value 250 \
#     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
#     --mc \
#     --stepsize 0.2 \
#     --constantinit 1


# # BNV2 Bottleneck Version 2, ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5
# go run main.go autoconfig \
#     --appname bookstore_nodejs \
#     --config configurations/bookstore_nodejs.yaml \
#     --workload  125_110_0.33_0.33_0.34 \
#     --testName bnv2-250-2.0-mc-c-0.5-0.02 \
#     bnv2 \
#     --property RTToleranceIntervalUBoundc90p95  \
#     --value 250 \
#     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
#     --mc \
#     --stepsize 2.0 \
#     --minstepsize 0.02


# VUSs=(75 100 125)
# for VUS in "${VUSs[@]}"
# do
    # # BNV2 Bottleneck Version 2, ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5
    # go run main.go autoconfig \
    #     --appname bookstore_nodejs \
    #     --config configurations/bookstore_nodejs.yaml \
    #     --workload  "${VUS}_110_0.33_0.33_0.34" \
    #     --testName bnv2-250-2.0-mc-c-0.5 \
    #     bnv2 \
    #     --property RTToleranceIntervalUBoundc90p95  \
    #     --value 250 \
    #     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
    #     --mc \
    #     --stepsize 2.0

    # # BNV2 Bottleneck Version 2, ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5
    # go run main.go autoconfig \
    #     --appname bookstore_nodejs \
    #     --config configurations/bookstore_nodejs.yaml \
    #     --workload  "${VUS}_110_0.33_0.33_0.34" \
    #     --testName bnv2-250-2.0-mc-c-0.5-0.02 \
    #     bnv2 \
    #     --property RTToleranceIntervalUBoundc90p95  \
    #     --value 250 \
    #     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
    #     --mc \
    #     --stepsize 2.0 \
    #     --minstepsize 0.02

    # # BNV2 Bottleneck Version 2, ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5
    # go run main.go autoconfig \
    #     --appname bookstore_nodejs \
    #     --config configurations/bookstore_nodejs.yaml \
    #     --workload  "${VUS}_110_0.33_0.33_0.34" \
    #     --testName bnv2-250-4.0-mc-c-0.5-0.02 \
    #     bnv2 \
    #     --property RTToleranceIntervalUBoundc90p95  \
    #     --value 250 \
    #     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
    #     --mc \
    #     --stepsize 4.0 \
    #     --minstepsize 0.02

    # # BottleNeck Versoin 1, ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5
    # go run main.go autoconfig \
    #     --appname bookstore_nodejs \
    #     --config configurations/bookstore_nodejs.yaml \
    #     --workload  "${VUS}_110_0.33_0.33_0.34" \
    #     --testName bnv1-250-0.2-mc-c-0.5 \
    #     bnv1 \
    #     --property RTToleranceIntervalUBoundc90p95  \
    #     --value 250 \
    #     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
    #     --mc \
    #     --stepsize 0.2 \
    #     --constantinit 0.5

    # # BottleNeck Versoin 1, ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5
    # go run main.go autoconfig \
    #     --appname bookstore_nodejs \
    #     --config configurations/bookstore_nodejs.yaml \
    #     --workload  "${VUS}_110_0.33_0.33_0.34" \
    #     --testName bnv1-250-0.5-mc-c-0.5 \
    #     bnv1 \
    #     --property RTToleranceIntervalUBoundc90p95  \
    #     --value 250 \
    #     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
    #     --mc \
    #     --stepsize 0.5 \
    #     --constantinit 0.5

    # # BottleNeck Versoin 1, ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5
    # go run main.go autoconfig \
    #     --appname bookstore_nodejs \
    #     --config configurations/bookstore_nodejs.yaml \
    #     --workload  "${VUS}_110_0.33_0.33_0.34" \
    #     --testName bnv1-250-1.0-mc-c-0.5 \
    #     bnv1 \
    #     --property RTToleranceIntervalUBoundc90p95  \
    #     --value 250 \
    #     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
    #     --mc \
    #     --stepsize 1.0 \
    #     --constantinit 0.5

# done

# VUSs=(10 20 30 )
# for VUS in "${VUSs[@]}"
# do
#     # BNV2 Bottleneck Version 2, ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5
#     go run main.go autoconfig \
#         --appname muck_star \
#         --config configurations/muck_star.yaml \
#         --workload  "${VUS}_110_0.3_0.2_0.1_0.1_0.3" \
#         --testName bnv2-250-1.0-mc-c-0.5 \
#         bnv2 \
#         --property RTToleranceIntervalUBoundc90p95  \
#         --value 250 \
#         --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star_demands.yaml \
#         --mc \
#         --stepsize 1.0

#     # # BottleNeck Versoin 1, ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5
#     # go run main.go autoconfig \
#     #     --appname muck_star \
#     #     --config configurations/muck_star.yaml \
#     #     --workload  "${VUS}_110_0.3_0.2_0.1_0.1_0.3" \
#     #     --testName bnv1-300-0.2-mc-c-0.5 \
#     #     bnv1 \
#     #     --property RTToleranceIntervalUBoundc90p95  \
#     #     --value 300 \
#     #     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star_demands.yaml \
#     #     --mc \
#     #     --stepsize 0.2 \
#     #     --constantinit 0.5

#     # # BottleNeck Versoin 1, ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5
#     # go run main.go autoconfig \
#     #     --appname muck_star \
#     #     --config configurations/muck_star.yaml \
#     #     --workload  "${VUS}_110_0.3_0.2_0.1_0.1_0.3" \
#     #     --testName bnv1-300-0.5-mc-c-0.5 \
#     #     bnv1 \
#     #     --property RTToleranceIntervalUBoundc90p95  \
#     #     --value 300 \
#     #     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star_demands.yaml \
#     #     --mc \
#     #     --stepsize 0.5 \
#     #     --constantinit 0.5

#     # # BottleNeck Versoin 1, ToleranceIntervalc90p95 must be less thatn 250ms with stepsize = 0.5
#     # go run main.go autoconfig \
#     #     --appname muck_star \
#     #     --config configurations/muck_star.yaml \
#     #     --workload  "${VUS}_110_0.3_0.2_0.1_0.1_0.3" \
#     #     --testName bnv1-300-1.0-mc-c-0.5 \
#     #     bnv1 \
#     #     --property RTToleranceIntervalUBoundc90p95  \
#     #     --value 300 \
#     #     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star_demands.yaml \
#     #     --mc \
#     #     --stepsize 1.0 \
#     #     --constantinit 0.5
# done