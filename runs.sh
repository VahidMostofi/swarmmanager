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

go run cmd/swarm-autoconfigure/main.go 500_80_0.65_10 test_532 Single


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

