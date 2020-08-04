#!/bin/bash
SLA=350
#workloads=(300_80_0.3_10 400_80_0.5_10 420_80_0.7_10 500_80_0.65_10 )
workloads=(300_80_0.3_10 )
go run cmd/swarm-autoconfigure/main.go 1_360_0.5_1 finding_demands Single
for WORKLOAD in "${workloads[@]}"
do
    echo "working on ${WORKLOAD}" 
    # ResponseTimeSimpleIncrease
    go run cmd/swarm-autoconfigure/main.go $WORKLOAD "rtsi_95_${SLA}" ResponseTimeSimpleIncrease -property RTToleranceIntervalUBoundc90p95 -value $SLA

    # HybridCPUUtilResponseTimeSimpleIncrease
    go run cmd/swarm-autoconfigure/main.go $WORKLOAD "hybrid_95_${SLA}_90_80" CPUUtil_RT_Hybrid -property RTToleranceIntervalUBoundc90p95 -value $SLA

    # CPUUsage95Percentile
    go run cmd/swarm-autoconfigure/main.go $WORKLOAD cpu95_50 CPUUsageIncrease -property CPUUsage95Percentile -threshold 50
    go run cmd/swarm-autoconfigure/main.go $WORKLOAD cpu95_60 CPUUsageIncrease -property CPUUsage95Percentile -threshold 60
    go run cmd/swarm-autoconfigure/main.go $WORKLOAD cpu95_70 CPUUsageIncrease -property CPUUsage95Percentile -threshold 70
    go run cmd/swarm-autoconfigure/main.go $WORKLOAD cpu95_80 CPUUsageIncrease -property CPUUsage95Percentile -threshold 80
    go run cmd/swarm-autoconfigure/main.go $WORKLOAD cpu95_90 CPUUsageIncrease -property CPUUsage95Percentile -threshold 90

    # CPUUsage90Percentile
    go run cmd/swarm-autoconfigure/main.go $WORKLOAD cpu90_50 CPUUsageIncrease -property CPUUsage90Percentile -threshold 50
    go run cmd/swarm-autoconfigure/main.go $WORKLOAD cpu90_60 CPUUsageIncrease -property CPUUsage90Percentile -threshold 60
    go run cmd/swarm-autoconfigure/main.go $WORKLOAD cpu90_70 CPUUsageIncrease -property CPUUsage90Percentile -threshold 70
    go run cmd/swarm-autoconfigure/main.go $WORKLOAD cpu90_80 CPUUsageIncrease -property CPUUsage90Percentile -threshold 80
    go run cmd/swarm-autoconfigure/main.go $WORKLOAD cpu90_90 CPUUsageIncrease -property CPUUsage90Percentile -threshold 90

    # CPUUsageMean
    go run cmd/swarm-autoconfigure/main.go $WORKLOAD cpu_mean_50 CPUUsageIncrease -property CPUUsageMean -threshold 50
    go run cmd/swarm-autoconfigure/main.go $WORKLOAD cpu_mean_60 CPUUsageIncrease -property CPUUsageMean -threshold 60
    go run cmd/swarm-autoconfigure/main.go $WORKLOAD cpu_mean_70 CPUUsageIncrease -property CPUUsageMean -threshold 70
    go run cmd/swarm-autoconfigure/main.go $WORKLOAD cpu_mean_80 CPUUsageIncrease -property CPUUsageMean -threshold 80
    go run cmd/swarm-autoconfigure/main.go $WORKLOAD cpu_mean_90 CPUUsageIncrease -property CPUUsageMean -threshold 90
done