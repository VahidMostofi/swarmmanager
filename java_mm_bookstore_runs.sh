#!/bin/bash

go run cmd/swarm-autoconfigure/main.go 300_120_0.3_10 cpu95_50 CPUUsageIncrease -property CPUUsage95Percentile -threshold 50
go run cmd/swarm-autoconfigure/main.go 300_120_0.3_10 cpu95_60 CPUUsageIncrease -property CPUUsage95Percentile -threshold 60
go run cmd/swarm-autoconfigure/main.go 300_120_0.3_10 cpu95_70 CPUUsageIncrease -property CPUUsage95Percentile -threshold 70
go run cmd/swarm-autoconfigure/main.go 300_120_0.3_10 cpu95_80 CPUUsageIncrease -property CPUUsage95Percentile -threshold 80
go run cmd/swarm-autoconfigure/main.go 300_120_0.3_10 cpu95_90 CPUUsageIncrease -property CPUUsage95Percentile -threshold 90