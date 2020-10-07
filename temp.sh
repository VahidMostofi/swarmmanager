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
#     --constantinit 0.5
while true
do
go run main.go autoconfig \
    --appname bookstore_nodejs \
    --config configurations/bookstore_nodejs.yaml \
    --workload  "125_60_0.33_0.33_0.34" \
    --testName brute-force \
    brute

sleep 1
done
