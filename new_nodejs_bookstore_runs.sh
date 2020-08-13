#!/bin/bash
SLA=250

# workloads=(400_120_0.1_10 400_120_0.3_10 400_120_0.5_10 400_120_0.7_10 400_120_0.9_10 )
# workloads=(550_120_0.1_10 550_120_0.3_10 550_120_0.5_10 550_120_0.7_10 550_120_0.9_10 )
# workloads=(700_120_0.1_10 700_120_0.3_10 700_120_0.5_10 700_120_0.7_10 700_120_0.9_10 )
workloads=(475_120_0.1_10 475_120_0.3_10 475_120_0.5_10 475_120_0.7_10 475_120_0.9_10 625_120_0.1_10 625_120_0.3_10 625_120_0.5_10 625_120_0.7_10 625_120_0.9_10 )

for WORKLOAD in "${workloads[@]}"
do
    echo "working on ${WORKLOAD}" 

    # # Fractional CPU increase, based on initial (estimated) CPU utilization, trying to mee ToleranceInterval, one big fat container
    # go run cmd/swarm-autoconfigure/main.go $WORKLOAD "adfc_utilization_ti_95_${SLA}" AddDifferentFractionalCPUcores -property RTToleranceIntervalUBoundc90p95 -value ${SLA} -amount 1 -indicator utilization    

    # # Fractional CPU increase, based on initial (estimated) CPU utilization, trying to mee ToleranceInterval, multiple containers
    # go run cmd/swarm-autoconfigure/main.go $WORKLOAD "mc_adfc_utilization_ti_95_${SLA}" AddDifferentFractionalCPUcores -property RTToleranceIntervalUBoundc90p95 -value ${SLA} -amount 1 -indicator utilization -mc

    # Consider each path in the graph of routes separately, also use Estimated CPU Utilization as the initial configuration, use ToleranceInterval, single fat container, with step size of 1
    go run cmd/swarm-autoconfigure/main.go 400_120_0.1_10 "ppeus_ti_equal_steps_950_${SLA}" PerPathEUBasedScaling -property RTToleranceIntervalUBoundc90p95 -value ${SLA} -step 1

    # Consider each path in the graph of routes separately, also use Estimated CPU Utilization as the initial configuration, use ToleranceInterval, mutli container, with step size of 1
    go run cmd/swarm-autoconfigure/main.go 400_120_0.1_10 "mc_ppeus_ti_equal_steps_950_${SLA}" PerPathEUBasedScaling -property RTToleranceIntervalUBoundc90p95 -value ${SLA} -step 1 -mv

    # Fractional CPU incrase, amount=0.33, tolerance interval, (its like sharing 1 core between three services equaly) with tolerance interval
    go run cmd/swarm-autoconfigure/main.go $WORKLOAD "afc_0.33_ti_95_${SLA}" AddFractionalCPUcores -property RTToleranceIntervalUBoundc90p95 -value ${SLA} -amount 0.33


done