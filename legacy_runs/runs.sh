#!/bin/bash
# CPUUsage95Percentile
# go run cmd/swarm-autoconfigure/main.go 500_80_0.65_10 cpu95_50 CPUUsageIncrease -property CPUUsage95Percentile -threshold 50
# go run cmd/swarm-autoconfigure/main.go 500_80_0.65_10 cpu95_60 CPUUsageIncrease -property CPUUsage95Percentile -threshold 60
# go run cmd/swarm-autoconfigure/main.go 500_80_0.65_10 cpu95_70 CPUUsageIncrease -property CPUUsage95Percentile -threshold 70
# go run cmd/swarm-autoconfigure/main.go 500_80_0.65_10 cpu95_80 CPUUsageIncrease -property CPUUsage95Percentile -threshold 80
# go run cmd/swarm-autoconfigure/main.go 500_80_0.65_10 cpu95_90 CPUUsageIncrease -property CPUUsage95Percentile -threshold 90

# CPUUsage90Percentile
# go run cmd/swarm-autoconfigure/main.go 500_80_0.65_10 cpu90_50 CPUUsageIncrease -property CPUUsage90Percentile -threshold 50
# go run cmd/swarm-autoconfigure/main.go 500_80_0.65_10 cpu90_60 CPUUsageIncrease -property CPUUsage90Percentile -threshold 60
# go run cmd/swarm-autoconfigure/main.go 500_80_0.65_10 cpu90_70 CPUUsageIncrease -property CPUUsage90Percentile -threshold 70
# go run cmd/swarm-autoconfigure/main.go 500_80_0.65_10 cpu90_80 CPUUsageIncrease -property CPUUsage90Percentile -threshold 80
# go run cmd/swarm-autoconfigure/main.go 500_80_0.65_10 cpu90_90 CPUUsageIncrease -property CPUUsage90Percentile -threshold 90

# CPUUsageMean
# go run cmd/swarm-autoconfigure/main.go 500_80_0.65_10 cpu_mean_50 CPUUsageIncrease -property CPUUsageMean -threshold 50
# go run cmd/swarm-autoconfigure/main.go 500_80_0.65_10 cpu_mean_60 CPUUsageIncrease -property CPUUsageMean -threshold 60
# go run cmd/swarm-autoconfigure/main.go 500_80_0.65_10 cpu_mean_70 CPUUsageIncrease -property CPUUsageMean -threshold 70
# go run cmd/swarm-autoconfigure/main.go 500_80_0.65_10 cpu_mean_80 CPUUsageIncrease -property CPUUsageMean -threshold 80
# go run cmd/swarm-autoconfigure/main.go 500_80_0.65_10 cpu_mean_90 CPUUsageIncrease -property CPUUsageMean -threshold 90

# ResponseTimeSimpleIncrease
# go run cmd/swarm-autoconfigure/main.go 300_80_0.3_10 rtsi_95_350 ResponseTimeSimpleIncrease -property RTToleranceIntervalUBoundc90p95 -value 350
# go run cmd/swarm-autoconfigure/main.go 400_80_0.5_10 rtsi_95_350 ResponseTimeSimpleIncrease -property RTToleranceIntervalUBoundc90p95 -value 350
# go run cmd/swarm-autoconfigure/main.go 420_80_0.7_10 rtsi_95_350 ResponseTimeSimpleIncrease -property RTToleranceIntervalUBoundc90p95 -value 350
# go run cmd/swarm-autoconfigure/main.go 500_80_0.65_10 rtsi_95_350 ResponseTimeSimpleIncrease -property RTToleranceIntervalUBoundc90p95 -value 350

# HybridCPUUtilResponseTimeSimpleIncrease
# go run cmd/swarm-autoconfigure/main.go 300_80_0.3_10 hybrid_95_350 CPUUtil_RT_Hybrid -property RTToleranceIntervalUBoundc90p95 -value 350
# go run cmd/swarm-autoconfigure/main.go 400_80_0.5_10 hybrid_95_350 CPUUtil_RT_Hybrid -property RTToleranceIntervalUBoundc90p95 -value 350
# go run cmd/swarm-autoconfigure/main.go 420_80_0.7_10 hybrid_95_350 CPUUtil_RT_Hybrid -property RTToleranceIntervalUBoundc90p95 -value 350
# go run cmd/swarm-autoconfigure/main.go 500_80_0.65_10 hybrid_95_350 CPUUtil_RT_Hybrid -property RTToleranceIntervalUBoundc90p95 -value 350

# PredefinedSearch
# go run cmd/swarm-autoconfigure/main.go 250_45_0.5_10 predefined_rtsi_90_350 PredefinedSearch

# Signle
# go run cmd/swarm-autoconfigure/main.go 300_45_0.3_10 after_predefined_rtsi_95_350_B Single
# go run cmd/swarm-autoconfigure/main.go 20_49_0.5_1 test Single
# go run cmd/swarm-autoconfigure/main.go 500_80_0.65_10 test_532 Single


# Finidng Demands:

# go run cmd/swarm-autoconfigure/main.go 1_360_0.5_1 finding_demands Single


# MOBO
# go run cmd/swarm-autoconfigure/main.go 300_80_0.3_10 mobo_244 MOBO auth 2 books 4 gateway 4
# go run cmd/swarm-autoconfigure/main.go 300_80_0.3_10 mobo_233 MOBO auth 2 books 3 gateway 3

# go run cmd/swarm-autoconfigure/main.go 400_80_0.5_10 mobo_344 MOBO auth 3 books 4 gateway 4
# go run cmd/swarm-autoconfigure/main.go 400_80_0.5_10 mobo_333 MOBO auth 3 books 3 gateway 3

# go run cmd/swarm-autoconfigure/main.go 420_80_0.7_10 mobo_434 MOBO auth 4 books 3 gateway 4
# go run cmd/swarm-autoconfigure/main.go 420_80_0.7_10 mobo_333 MOBO auth 3 books 3 gateway 3

# go run cmd/swarm-autoconfigure/main.go 500_80_0.65_10 mobo_535 MOBO auth 5 books 3 gateway 5
# go run cmd/swarm-autoconfigure/main.go 500_80_0.65_10 mobo_434 MOBO auth 4 books 3 gateway 4



# go run cmd/swarm-autoconfigure/main.go 300_80_0.3_10 rtsi_95_350 AddFractionalCPUcores -property RTToleranceIntervalUBoundc90p95 -value 350 -amount 0.5

# go run cmd/swarm-autoconfigure/main.go 300_120_0.3_10 proportional_to_deman_95_350 AddFractionalCPUcores -property RTToleranceIntervalUBoundc90p95 -value 350 -amount 1

# go run cmd/swarm-autoconfigure/main.go 300_120_0.3_10 proportional_to_utilizations_95_350 AddFractionalCPUcores -property RTToleranceIntervalUBoundc90p95 -value 350 -amount 1



#go run cmd/swarm-autoconfigure/main.go 1_120_1_0.001 auth_only_111 Single


# go run cmd/swarm-autoconfigure/main.go 300_120_0.3_10 afc_0.5_95_500 AddFractionalCPUcores -property RTToleranceIntervalUBoundc90p95 -value 500 -amount 0.5
# go run cmd/swarm-autoconfigure/main.go 400_120_0.5_10 afc_0.5_95_500 AddFractionalCPUcores -property RTToleranceIntervalUBoundc90p95 -value 500 -amount 0.5
# go run cmd/swarm-autoconfigure/main.go 420_120_0.7_10 afc_0.5_95_500 AddFractionalCPUcores -property RTToleranceIntervalUBoundc90p95 -value 500 -amount 0.5
# go run cmd/swarm-autoconfigure/main.go 500_120_0.65_10 afc_0.5_95_500 AddFractionalCPUcores -property RTToleranceIntervalUBoundc90p95 -value 500 -amount 0.5

# go run cmd/swarm-autoconfigure/main.go 300_120_0.3_10 afc_1_95_500 AddFractionalCPUcores -property RTToleranceIntervalUBoundc90p95 -value 500 -amount 1
# go run cmd/swarm-autoconfigure/main.go 400_120_0.5_10 afc_1_95_500 AddFractionalCPUcores -property RTToleranceIntervalUBoundc90p95 -value 500 -amount 1
# go run cmd/swarm-autoconfigure/main.go 420_120_0.7_10 afc_1_95_500 AddFractionalCPUcores -property RTToleranceIntervalUBoundc90p95 -value 500 -amount 1
# go run cmd/swarm-autoconfigure/main.go 500_120_0.65_10 afc_1_95_500 AddFractionalCPUcores -property RTToleranceIntervalUBoundc90p95 -value 500 -amount 1

# AddDifferentFractionalCPUcores
# go run cmd/swarm-autoconfigure/main.go 300_120_0.3_10 adfc_demand_95_500 AddDifferentFractionalCPUcores -property RTToleranceIntervalUBoundc90p95 -value 500 -amount 1 -indicator demand
# go run cmd/swarm-autoconfigure/main.go 400_120_0.5_10 adfc_demand_95_500 AddDifferentFractionalCPUcores -property RTToleranceIntervalUBoundc90p95 -value 500 -amount 1 -indicator demand
# go run cmd/swarm-autoconfigure/main.go 420_120_0.7_10 adfc_demand_95_500 AddDifferentFractionalCPUcores -property RTToleranceIntervalUBoundc90p95 -value 500 -amount 1 -indicator demand
# go run cmd/swarm-autoconfigure/main.go 500_120_0.65_10 adfc_demand_95_500 AddDifferentFractionalCPUcores -property RTToleranceIntervalUBoundc90p95 -value 500 -amount 1 -indicator demand

# go run cmd/swarm-autoconfigure/main.go 300_120_0.3_10 adfc_utilization_95_500 AddDifferentFractionalCPUcores -property RTToleranceIntervalUBoundc90p95 -value 500 -amount 1 -indicator utilization
# go run cmd/swarm-autoconfigure/main.go 400_120_0.5_10 adfc_utilization_95_500 AddDifferentFractionalCPUcores -property RTToleranceIntervalUBoundc90p95 -value 500 -amount 1 -indicator utilization
# go run cmd/swarm-autoconfigure/main.go 420_120_0.7_10 adfc_utilization_95_500 AddDifferentFractionalCPUcores -property RTToleranceIntervalUBoundc90p95 -value 500 -amount 1 -indicator utilization
# go run cmd/swarm-autoconfigure/main.go 500_120_0.65_10 adfc_utilization_95_500 AddDifferentFractionalCPUcores -property RTToleranceIntervalUBoundc90p95 -value 500 -amount 1 -indicator utilization


#go run cmd/swarm-autoconfigure/main.go 420_120_0.7_10 something Single
# go run cmd/swarm-autoconfigure/main.go 420_120_0.7_10 something_else Single

# go run cmd/swarm-autoconfigure/main.go 400_120_0.1_10 mc_ppeus_ti_equal_steps_950_250 PerPathEUBasedScaling -property RTToleranceIntervalUBoundc90p95 -value 250 -step 1


# # the values for auth, books and gateway are not being used! they are here for backward compatiblity
# go run cmd/swarm-autoconfigure/main.go 400_120_0.5_10 mobo_CPU_split_mc MOBO auth 1 books 1 gateway 1