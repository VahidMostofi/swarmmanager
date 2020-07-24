#!/bin/bash
# CPUUsage95Percentile
# go run cmd/swarm-autoconfigure/main.go 300_180_0.3_10 cpu95_50 CPUUsageIncrease -property CPUUsage95Percentile -threshold 50
# go run cmd/swarm-autoconfigure/main.go 300_180_0.3_10 cpu95_60 CPUUsageIncrease -property CPUUsage95Percentile -threshold 60
go run cmd/swarm-autoconfigure/main.go 300_180_0.3_10 cpu95_70 CPUUsageIncrease -property CPUUsage95Percentile -threshold 70
go run cmd/swarm-autoconfigure/main.go 300_180_0.3_10 cpu95_80 CPUUsageIncrease -property CPUUsage95Percentile -threshold 80
go run cmd/swarm-autoconfigure/main.go 300_180_0.3_10 cpu95_90 CPUUsageIncrease -property CPUUsage95Percentile -threshold 90

# CPUUsage90Percentile
go run cmd/swarm-autoconfigure/main.go 300_180_0.3_10 cpu90_50 CPUUsageIncrease -property CPUUsage90Percentile -threshold 50
go run cmd/swarm-autoconfigure/main.go 300_180_0.3_10 cpu90_60 CPUUsageIncrease -property CPUUsage90Percentile -threshold 60
go run cmd/swarm-autoconfigure/main.go 300_180_0.3_10 cpu90_70 CPUUsageIncrease -property CPUUsage90Percentile -threshold 70
go run cmd/swarm-autoconfigure/main.go 300_180_0.3_10 cpu90_80 CPUUsageIncrease -property CPUUsage90Percentile -threshold 80
go run cmd/swarm-autoconfigure/main.go 300_180_0.3_10 cpu90_90 CPUUsageIncrease -property CPUUsage90Percentile -threshold 90

# CPUUsageMean
go run cmd/swarm-autoconfigure/main.go 300_180_0.3_10 cpu_mean_50 CPUUsageIncrease -property CPUUsageMean -threshold 50
go run cmd/swarm-autoconfigure/main.go 300_180_0.3_10 cpu_mean_60 CPUUsageIncrease -property CPUUsageMean -threshold 60
go run cmd/swarm-autoconfigure/main.go 300_180_0.3_10 cpu_mean_70 CPUUsageIncrease -property CPUUsageMean -threshold 70
go run cmd/swarm-autoconfigure/main.go 300_180_0.3_10 cpu_mean_80 CPUUsageIncrease -property CPUUsageMean -threshold 80
go run cmd/swarm-autoconfigure/main.go 300_180_0.3_10 cpu_mean_90 CPUUsageIncrease -property CPUUsageMean -threshold 90

# # ResponseTimeSimpleIncrease
# go run cmd/swarm-autoconfigure/main.go 300_45_0.5_10 rtsi_95_300_B ResponseTimeSimpleIncrease -property ResponseTimes95Percentile -value 300

# PredefinedSearch
# go run cmd/swarm-autoconfigure/main.go 250_45_0.5_10 predefined_rtsi_90_350 PredefinedSearch

# Signle
# go run cmd/swarm-autoconfigure/main.go 300_45_0.3_10 after_predefined_rtsi_95_350_B Single