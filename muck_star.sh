
# # #FINDING DEMANDS FOR STAR
# go run main.go autoconfig \
#     --appname muck_star-small \
#     --config configurations/muck_star.yaml \
#     --workload  1_1300_0.45_0.25_0.15_0.15 \
#     --testName demands \
#     demands \
#     --duration 1200 \
#     --resultpath /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star-small_demands.yaml


VUSs=(10 20 30 40 50)
for VUS in "${VUSs[@]}"
do

    # # BNV2 Bottleneck Version 2, ToleranceIntervalc90p95 must be less thatn 300ms with stepsize = 0.5
    # go run main.go autoconfig \
    #     --appname muck_star-small \
    #     --config configurations/muck_star.yaml \
    #     --workload  "${VUS}_110_0.45_0.25_0.15_0.15" \
    #     --testName bnv2-300-2.0-mc-c-0.5 \
    #     bnv2 \
    #     --property RTToleranceIntervalUBoundc90p95  \
    #     --value 300 \
    #     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star-small_demands.yaml \
    #     --mc \
    #     --stepsize 2.0

    # # BNV2 Bottleneck Version 2, ToleranceIntervalc90p95 must be less thatn 300ms with stepsize = 0.5
    # go run main.go autoconfig \
    #     --appname muck_star-small \
    #     --config configurations/muck_star.yaml \
    #     --workload  "${VUS}_110_0.45_0.25_0.15_0.15" \
    #     --testName bnv2-300-1.5-mc-c-0.5 \
    #     bnv2 \
    #     --property RTToleranceIntervalUBoundc90p95  \
    #     --value 300 \
    #     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star-small_demands.yaml \
    #     --mc \
    #     --stepsize 1.5

    # # BNV2 Bottleneck Version 2, ToleranceIntervalc90p95 must be less thatn 300ms with stepsize = 0.5
    # go run main.go autoconfig \
    #     --appname muck_star-small \
    #     --config configurations/muck_star.yaml \
    #     --workload  "${VUS}_110_0.45_0.25_0.15_0.15" \
    #     --testName bnv2-300-1.0-mc-c-0.5 \
    #     bnv2 \
    #     --property RTToleranceIntervalUBoundc90p95  \
    #     --value 300 \
    #     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star-small_demands.yaml \
    #     --mc \
    #     --stepsize 1.0

    # # BNV2 Bottleneck Version 2, ToleranceIntervalc90p95 must be less thatn 300ms with stepsize = 0.5
    # go run main.go autoconfig \
    #     --appname muck_star-small \
    #     --config configurations/muck_star.yaml \
    #     --workload  "${VUS}_110_0.45_0.25_0.15_0.15" \
    #     --testName bnv2-300-0.5-mc-c-0.5 \
    #     bnv2 \
    #     --property RTToleranceIntervalUBoundc90p95  \
    #     --value 300 \
    #     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star-small_demands.yaml \
    #     --mc \
    #     --stepsize 0.5

    # # BNV2 Bottleneck Version 2, ToleranceIntervalc90p95 must be less thatn 300ms with stepsize = 0.2
    # go run main.go autoconfig \
    #     --appname muck_star-small \
    #     --config configurations/muck_star.yaml \
    #     --workload  "${VUS}_110_0.45_0.25_0.15_0.15" \
    #     --testName bnv2-300-0.2-mc-c-0.5 \
    #     bnv2 \
    #     --property RTToleranceIntervalUBoundc90p95  \
    #     --value 300 \
    #     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star-small_demands.yaml \
    #     --mc \
    #     --stepsize 0.2


    # # BottleNeck Versoin 1, ToleranceIntervalc90p95 must be less thatn 300ms with stepsize = 0.5
    # go run main.go autoconfig \
    #     --appname muck_star-small \
    #     --config configurations/muck_star.yaml \
    #     --workload  "${VUS}_110_0.45_0.25_0.15_0.15" \
    #     --testName bnv1-300-2.0-mc-c-0.5 \
    #     bnv1 \
    #     --property RTToleranceIntervalUBoundc90p95  \
    #     --value 300 \
    #     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star-small_demands.yaml \
    #     --mc \
    #     --stepsize 2.0 \
    #     --constantinit 0.5

    # # BottleNeck Versoin 1, ToleranceIntervalc90p95 must be less thatn 300ms with stepsize = 0.5
    # go run main.go autoconfig \
    #     --appname muck_star-small \
    #     --config configurations/muck_star.yaml \
    #     --workload  "${VUS}_110_0.45_0.25_0.15_0.15" \
    #     --testName bnv1-300-1.5-mc-c-0.5 \
    #     bnv1 \
    #     --property RTToleranceIntervalUBoundc90p95  \
    #     --value 300 \
    #     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star-small_demands.yaml \
    #     --mc \
    #     --stepsize 1.5 \
    #     --constantinit 0.5
    
    # # BottleNeck Versoin 1, ToleranceIntervalc90p95 must be less thatn 300ms with stepsize = 0.5
    # go run main.go autoconfig \
    #     --appname muck_star-small \
    #     --config configurations/muck_star.yaml \
    #     --workload  "${VUS}_110_0.45_0.25_0.15_0.15" \
    #     --testName bnv1-300-0.2-mc-c-0.5 \
    #     bnv1 \
    #     --property RTToleranceIntervalUBoundc90p95  \
    #     --value 300 \
    #     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star-small_demands.yaml \
    #     --mc \
    #     --stepsize 0.2 \
    #     --constantinit 0.5

    # # BottleNeck Versoin 1, ToleranceIntervalc90p95 must be less thatn 300ms with stepsize = 0.5
    # go run main.go autoconfig \
    #     --appname muck_star-small \
    #     --config configurations/muck_star.yaml \
    #     --workload  "${VUS}_110_0.45_0.25_0.15_0.15" \
    #     --testName bnv1-300-0.5-mc-c-0.5 \
    #     bnv1 \
    #     --property RTToleranceIntervalUBoundc90p95  \
    #     --value 300 \
    #     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star-small_demands.yaml \
    #     --mc \
    #     --stepsize 0.5 \
    #     --constantinit 0.5

    # # BottleNeck Versoin 1, ToleranceIntervalc90p95 must be less thatn 300ms with stepsize = 0.5
    # go run main.go autoconfig \
    #     --appname muck_star-small \
    #     --config configurations/muck_star.yaml \
    #     --workload  "${VUS}_110_0.45_0.25_0.15_0.15" \
    #     --testName bnv1-300-1.0-mc-c-0.5 \
    #     bnv1 \
    #     --property RTToleranceIntervalUBoundc90p95  \
    #     --value 300 \
    #     --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star-small_demands.yaml \
    #     --mc \
    #     --stepsize 1.0 \
    #     --constantinit 0.5

    # BNV2 Bottleneck Version 2, ToleranceIntervalc90p95 must be less thatn 300ms with stepsize = 0.5
    go run main.go autoconfig \
        --appname muck_star-small \
        --config configurations/muck_star.yaml \
        --workload  "${VUS}_110_0.45_0.25_0.15_0.15" \
        --testName bnv2-300-2.0-mc-c-0.5-r95 \
        bnv2 \
        --property ResponseTimes95Percentile  \
        --value 300 \
        --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star-small_demands.yaml \
        --mc \
        --stepsize 2.0

    # BNV2 Bottleneck Version 2, ToleranceIntervalc90p95 must be less thatn 300ms with stepsize = 0.5
    go run main.go autoconfig \
        --appname muck_star-small \
        --config configurations/muck_star.yaml \
        --workload  "${VUS}_110_0.45_0.25_0.15_0.15" \
        --testName bnv2-300-1.5-mc-c-0.5-r95 \
        bnv2 \
        --property ResponseTimes95Percentile  \
        --value 300 \
        --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star-small_demands.yaml \
        --mc \
        --stepsize 1.5

    # BNV2 Bottleneck Version 2, ToleranceIntervalc90p95 must be less thatn 300ms with stepsize = 0.5
    go run main.go autoconfig \
        --appname muck_star-small \
        --config configurations/muck_star.yaml \
        --workload  "${VUS}_110_0.45_0.25_0.15_0.15" \
        --testName bnv2-300-1.0-mc-c-0.5-r95 \
        bnv2 \
        --property ResponseTimes95Percentile  \
        --value 300 \
        --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star-small_demands.yaml \
        --mc \
        --stepsize 1.0

    # BNV2 Bottleneck Version 2, ToleranceIntervalc90p95 must be less thatn 300ms with stepsize = 0.5
    go run main.go autoconfig \
        --appname muck_star-small \
        --config configurations/muck_star.yaml \
        --workload  "${VUS}_110_0.45_0.25_0.15_0.15" \
        --testName bnv2-300-0.5-mc-c-0.5-r95 \
        bnv2 \
        --property ResponseTimes95Percentile  \
        --value 300 \
        --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star-small_demands.yaml \
        --mc \
        --stepsize 0.5

    # BNV2 Bottleneck Version 2, ToleranceIntervalc90p95 must be less thatn 300ms with stepsize = 0.2
    go run main.go autoconfig \
        --appname muck_star-small \
        --config configurations/muck_star.yaml \
        --workload  "${VUS}_110_0.45_0.25_0.15_0.15" \
        --testName bnv2-300-0.2-mc-c-0.5-r95 \
        bnv2 \
        --property ResponseTimes95Percentile  \
        --value 300 \
        --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star-small_demands.yaml \
        --mc \
        --stepsize 0.2


    # BottleNeck Versoin 1, ToleranceIntervalc90p95 must be less thatn 300ms with stepsize = 0.5
    go run main.go autoconfig \
        --appname muck_star-small \
        --config configurations/muck_star.yaml \
        --workload  "${VUS}_110_0.45_0.25_0.15_0.15" \
        --testName bnv1-300-2.0-mc-c-0.5-r95 \
        bnv1 \
        --property ResponseTimes95Percentile  \
        --value 300 \
        --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star-small_demands.yaml \
        --mc \
        --stepsize 2.0 \
        --constantinit 0.5

    # BottleNeck Versoin 1, ToleranceIntervalc90p95 must be less thatn 300ms with stepsize = 0.5
    go run main.go autoconfig \
        --appname muck_star-small \
        --config configurations/muck_star.yaml \
        --workload  "${VUS}_110_0.45_0.25_0.15_0.15" \
        --testName bnv1-300-1.5-mc-c-0.5-r95 \
        bnv1 \
        --property ResponseTimes95Percentile  \
        --value 300 \
        --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star-small_demands.yaml \
        --mc \
        --stepsize 1.5 \
        --constantinit 0.5
    
    # BottleNeck Versoin 1, ToleranceIntervalc90p95 must be less thatn 300ms with stepsize = 0.5
    go run main.go autoconfig \
        --appname muck_star-small \
        --config configurations/muck_star.yaml \
        --workload  "${VUS}_110_0.45_0.25_0.15_0.15" \
        --testName bnv1-300-0.2-mc-c-0.5-r95 \
        bnv1 \
        --property ResponseTimes95Percentile  \
        --value 300 \
        --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star-small_demands.yaml \
        --mc \
        --stepsize 0.2 \
        --constantinit 0.5

    # BottleNeck Versoin 1, ToleranceIntervalc90p95 must be less thatn 300ms with stepsize = 0.5
    go run main.go autoconfig \
        --appname muck_star-small \
        --config configurations/muck_star.yaml \
        --workload  "${VUS}_110_0.45_0.25_0.15_0.15" \
        --testName bnv1-300-0.5-mc-c-0.5-r95 \
        bnv1 \
        --property ResponseTimes95Percentile  \
        --value 300 \
        --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star-small_demands.yaml \
        --mc \
        --stepsize 0.5 \
        --constantinit 0.5

    # BottleNeck Versoin 1, ToleranceIntervalc90p95 must be less thatn 300ms with stepsize = 0.5
    go run main.go autoconfig \
        --appname muck_star-small \
        --config configurations/muck_star.yaml \
        --workload  "${VUS}_110_0.45_0.25_0.15_0.15" \
        --testName bnv1-300-1.0-mc-c-0.5-r95 \
        bnv1 \
        --property ResponseTimes95Percentile  \
        --value 300 \
        --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/muck_star-small_demands.yaml \
        --mc \
        --stepsize 1.0 \
        --constantinit 0.5
done
