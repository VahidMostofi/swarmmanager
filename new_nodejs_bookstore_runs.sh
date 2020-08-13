#!/bin/bash
SLA=250

# workloads=(400_120_0.1_10 400_120_0.3_10 400_120_0.5_10 400_120_0.7_10 400_120_0.9_10 )
# workloads=(550_120_0.1_10 550_120_0.3_10 550_120_0.5_10 550_120_0.7_10 550_120_0.9_10 )
# workloads=(700_120_0.1_10 700_120_0.3_10 700_120_0.5_10 700_120_0.7_10 700_120_0.9_10 )
workloads=(475_120_0.1_10 475_120_0.3_10 475_120_0.5_10 475_120_0.7_10 475_120_0.9_10 625_120_0.1_10 625_120_0.3_10 625_120_0.5_10 625_120_0.7_10 625_120_0.9_10 )

for WORKLOAD in "${workloads[@]}"
do
    echo "working on ${WORKLOAD}" 

    # Fractional CPU increase, based on initial (estimated) CPU utilization, trying to mee ToleranceInterval, one big fat container
    go run cmd/swarm-autoconfigure/main.go $WORKLOAD "adfc_utilization_ti_95_${SLA}" AddDifferentFractionalCPUcores -property RTToleranceIntervalUBoundc90p95 -value ${SLA} -amount 1 -indicator utilization    

    # Fractional CPU increase, based on initial (estimated) CPU utilization, trying to mee ToleranceInterval, multiple containers
    go run cmd/swarm-autoconfigure/main.go $WORKLOAD "mc_adfc_utilization_ti_95_${SLA}" AddDifferentFractionalCPUcores -property RTToleranceIntervalUBoundc90p95 -value ${SLA} -amount 1 -indicator utilization -mc

done