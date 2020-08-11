#!/bin/bash
SLA=250
#workloads=(300_120_0.3_10 400_120_0.5_10 420_120_0.7_10 500_120_0.65_10 600_120_0.3_10)
workloads=(400_120_0.5_10 420_120_0.7_10 500_120_0.65_10 600_120_0.3_10 700_120_0.5_10 )

for WORKLOAD in "${workloads[@]}"
do
    echo "working on ${WORKLOAD}" 

    # # CPUUsage95Percentile
    # go run cmd/swarm-autoconfigure/main.go $WORKLOAD cpu95_50 CPUUsageIncrease -property CPUUsage95Percentile -threshold 50
    # go run cmd/swarm-autoconfigure/main.go $WORKLOAD cpu95_60 CPUUsageIncrease -property CPUUsage95Percentile -threshold 60
    # go run cmd/swarm-autoconfigure/main.go $WORKLOAD cpu95_70 CPUUsageIncrease -property CPUUsage95Percentile -threshold 70
    # go run cmd/swarm-autoconfigure/main.go $WORKLOAD cpu95_80 CPUUsageIncrease -property CPUUsage95Percentile -threshold 80
    # go run cmd/swarm-autoconfigure/main.go $WORKLOAD cpu95_90 CPUUsageIncrease -property CPUUsage95Percentile -threshold 90

    # # CPUUsage90Percentile
    # go run cmd/swarm-autoconfigure/main.go $WORKLOAD cpu90_50 CPUUsageIncrease -property CPUUsage90Percentile -threshold 50
    # go run cmd/swarm-autoconfigure/main.go $WORKLOAD cpu90_60 CPUUsageIncrease -property CPUUsage90Percentile -threshold 60
    # go run cmd/swarm-autoconfigure/main.go $WORKLOAD cpu90_70 CPUUsageIncrease -property CPUUsage90Percentile -threshold 70
    # go run cmd/swarm-autoconfigure/main.go $WORKLOAD cpu90_80 CPUUsageIncrease -property CPUUsage90Percentile -threshold 80
    # go run cmd/swarm-autoconfigure/main.go $WORKLOAD cpu90_90 CPUUsageIncrease -property CPUUsage90Percentile -threshold 90

    # # CPUUsageMean
    # go run cmd/swarm-autoconfigure/main.go $WORKLOAD cpu_mean_50 CPUUsageIncrease -property CPUUsageMean -threshold 50
    # go run cmd/swarm-autoconfigure/main.go $WORKLOAD cpu_mean_60 CPUUsageIncrease -property CPUUsageMean -threshold 60
    # go run cmd/swarm-autoconfigure/main.go $WORKLOAD cpu_mean_70 CPUUsageIncrease -property CPUUsageMean -threshold 70
    # go run cmd/swarm-autoconfigure/main.go $WORKLOAD cpu_mean_80 CPUUsageIncrease -property CPUUsageMean -threshold 80
    # go run cmd/swarm-autoconfigure/main.go $WORKLOAD cpu_mean_90 CPUUsageIncrease -property CPUUsageMean -threshold 90

    # # # Fractional CPU incrase, amount=0.5 ToleranceInterval
    # # go run cmd/swarm-autoconfigure/main.go $WORKLOAD "afc_0.5_ti_95_${SLA}" AddFractionalCPUcores -property RTToleranceIntervalUBoundc90p95 -value ${SLA} -amount 0.5

    # # # Fractional CPU incrase, amount=1.0 ToleranceInterval
    # # go run cmd/swarm-autoconfigure/main.go $WORKLOAD "afc_1_95_ti_${SLA}" AddFractionalCPUcores -property RTToleranceIntervalUBoundc90p95 -value ${SLA} -amount 1

    # # Fractional CPU increase, based on initial (mean) demand ToleranceInterval
    # go run cmd/swarm-autoconfigure/main.go $WORKLOAD "adfc_demand_ti_95_${SLA}" AddDifferentFractionalCPUcores -property RTToleranceIntervalUBoundc90p95 -value ${SLA} -amount 1 -indicator demand

    # # Fractional CPU increase, based on initial (mean) cpu utilization ToleranceInterval
    # go run cmd/swarm-autoconfigure/main.go $WORKLOAD "adfc_utilization_ti_95_${SLA}" AddDifferentFractionalCPUcores -property RTToleranceIntervalUBoundc90p95 -value ${SLA} -amount 1 -indicator utilization

    # # # Fractional CPU incrase, amount=0.5 ResponseTime
    # # go run cmd/swarm-autoconfigure/main.go $WORKLOAD "afc_0.5_rt_95_${SLA}" AddFractionalCPUcores -property ResponseTimes95Percentile -value ${SLA} -amount 0.5

    # # # Fractional CPU incrase, amount=1.0 ResponseTime
    # # go run cmd/swarm-autoconfigure/main.go $WORKLOAD "afc_1_rt_95_${SLA}" AddFractionalCPUcores -property ResponseTimes95Percentile -value ${SLA} -amount 1

    # # Fractional CPU increase, based on initial (mean) demand ResponseTime
    # go run cmd/swarm-autoconfigure/main.go $WORKLOAD "adfc_demand_rt_95_${SLA}" AddDifferentFractionalCPUcores -property ResponseTimes95Percentile -value ${SLA} -amount 1 -indicator demand

    # # Fractional CPU increase, based on initial (mean) cpu utilization ResponseTime
    # go run cmd/swarm-autoconfigure/main.go $WORKLOAD "adfc_utilization_rt_95_${SLA}" AddDifferentFractionalCPUcores -property ResponseTimes95Percentile -value ${SLA} -amount 1 -indicator utilization

    # ###################################################################### mc stands for multiple containers
    # # Fractional CPU increase, based on initial (mean) demand ResponseTime
    # go run cmd/swarm-autoconfigure/main.go $WORKLOAD "mc_adfc_demand_rt_95_${SLA}" AddDifferentFractionalCPUcores -property ResponseTimes95Percentile -value ${SLA} -amount 1 -indicator demand

    # # Fractional CPU increase, based on initial (mean) cpu utilization ResponseTime
    # go run cmd/swarm-autoconfigure/main.go $WORKLOAD "mc_adfc_utilization_rt_95_${SLA}" AddDifferentFractionalCPUcores -property ResponseTimes95Percentile -value ${SLA} -amount 1 -indicator utilization

    # # # Fractional CPU increase, based on initial (mean) demand ToleranceInterval
    # go run cmd/swarm-autoconfigure/main.go $WORKLOAD "mc_adfc_demand_ti_95_${SLA}" AddDifferentFractionalCPUcores -property RTToleranceIntervalUBoundc90p95 -value ${SLA} -amount 1 -indicator demand

    # # # Fractional CPU increase, based on initial (mean) cpu utilization ToleranceInterval
    # go run cmd/swarm-autoconfigure/main.go $WORKLOAD "mc_adfc_utilization_ti_95_${SLA}" AddDifferentFractionalCPUcores -property RTToleranceIntervalUBoundc90p95 -value ${SLA} -amount 1 -indicator utilization    

    # # Fractional CPU increase, based on initial (mean) demand ToleranceInterval sharing 2 cores betweens services
    # go run cmd/swarm-autoconfigure/main.go $WORKLOAD "adfc_demand_ti_95_${SLA}_2" AddDifferentFractionalCPUcores -property RTToleranceIntervalUBoundc90p95 -value ${SLA} -amount 2 -indicator demand

    # # Fractional CPU increase, based on initial (mean) cpu utilization ToleranceInterval  sharing 2 cores betweens services
    # go run cmd/swarm-autoconfigure/main.go $WORKLOAD "adfc_utilization_ti_95_${SLA}_2" AddDifferentFractionalCPUcores -property RTToleranceIntervalUBoundc90p95 -value ${SLA} -amount 2 -indicator utilization

    # # Fractional CPU increase, based on initial (mean) demand ResponseTime  sharing 2 cores betweens services
    # go run cmd/swarm-autoconfigure/main.go $WORKLOAD "adfc_demand_rt_95_${SLA}_2" AddDifferentFractionalCPUcores -property ResponseTimes95Percentile -value ${SLA} -amount 2 -indicator demand

    # # Fractional CPU increase, based on initial (mean) cpu utilization ResponseTime  sharing 2 cores betweens services
    # go run cmd/swarm-autoconfigure/main.go $WORKLOAD "adfc_utilization_rt_95_${SLA}_2" AddDifferentFractionalCPUcores -property ResponseTimes95Percentile -value ${SLA} -amount 2 -indicator utilization

    # ###################################################################### mc stands for multiple containers
    # # Fractional CPU increase, based on initial (mean) demand ResponseTime  sharing 2 cores betweens services
    # go run cmd/swarm-autoconfigure/main.go $WORKLOAD "mc_adfc_demand_rt_95_${SLA}_2" AddDifferentFractionalCPUcores -property ResponseTimes95Percentile -value ${SLA} -amount 2 -indicator demand

    # # Fractional CPU increase, based on initial (mean) cpu utilization ResponseTime  sharing 2 cores betweens services
    # go run cmd/swarm-autoconfigure/main.go $WORKLOAD "mc_adfc_utilization_rt_95_${SLA}_2" AddDifferentFractionalCPUcores -property ResponseTimes95Percentile -value ${SLA} -amount 2 -indicator utilization

    # # # Fractional CPU increase, based on initial (mean) demand ToleranceInterval  sharing 2 cores betweens services
    # go run cmd/swarm-autoconfigure/main.go $WORKLOAD "mc_adfc_demand_ti_95_${SLA}_2" AddDifferentFractionalCPUcores -property RTToleranceIntervalUBoundc90p95 -value ${SLA} -amount 2 -indicator demand

    # # # Fractional CPU increase, based on initial (mean) cpu utilization ToleranceInterval  sharing 2 cores betweens services
    # go run cmd/swarm-autoconfigure/main.go $WORKLOAD "mc_adfc_utilization_ti_95_${SLA}_2" AddDifferentFractionalCPUcores -property RTToleranceIntervalUBoundc90p95 -value ${SLA} -amount 2 -indicator utilization    

    # Fractional CPU incrase, amount=0.33 ResponseTime (its like sharing 1 core between three services equaly) with response time
    go run cmd/swarm-autoconfigure/main.go $WORKLOAD "afc_0.33_rt_95_${SLA}" AddFractionalCPUcores -property ResponseTimes95Percentile -value ${SLA} -amount 0.3

    # Fractional CPU incrase, amount=0.33 ResponseTime (its like sharing 1 core between three services equaly) with tolerance interval
    go run cmd/swarm-autoconfigure/main.go $WORKLOAD "afc_0.33_ti_95_${SLA}" AddFractionalCPUcores -property RTToleranceIntervalUBoundc90p95 -value ${SLA} -amount 0.3

done