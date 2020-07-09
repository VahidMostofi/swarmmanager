# #!/bin/bash
# # CPUUsage95Percentile
# go run cmd/swarm-autoconfigure/main.go cpu95_50 CPUUsageIncrease -property CPUUsage95Percentile -threshold 50
# go run cmd/swarm-autoconfigure/main.go cpu95_60 CPUUsageIncrease -property CPUUsage95Percentile -threshold 60
# go run cmd/swarm-autoconfigure/main.go cpu95_70 CPUUsageIncrease -property CPUUsage95Percentile -threshold 70
# go run cmd/swarm-autoconfigure/main.go cpu95_80 CPUUsageIncrease -property CPUUsage95Percentile -threshold 80
# go run cmd/swarm-autoconfigure/main.go cpu95_90 CPUUsageIncrease -property CPUUsage95Percentile -threshold 90

# # CPUUsage90Percentile
# go run cmd/swarm-autoconfigure/main.go cpu90_50 CPUUsageIncrease -property CPUUsage90Percentile -threshold 50
# go run cmd/swarm-autoconfigure/main.go cpu90_60 CPUUsageIncrease -property CPUUsage90Percentile -threshold 60
# go run cmd/swarm-autoconfigure/main.go cpu90_70 CPUUsageIncrease -property CPUUsage90Percentile -threshold 70
# go run cmd/swarm-autoconfigure/main.go cpu90_80 CPUUsageIncrease -property CPUUsage90Percentile -threshold 80
# go run cmd/swarm-autoconfigure/main.go cpu90_90 CPUUsageIncrease -property CPUUsage90Percentile -threshold 90

# # CPUUsage90Percentile
go run cmd/swarm-autoconfigure/main.go cpu_mean_50 CPUUsageIncrease -property CPUUsageMean -threshold 50
go run cmd/swarm-autoconfigure/main.go cpu_mean_60 CPUUsageIncrease -property CPUUsageMean -threshold 60
go run cmd/swarm-autoconfigure/main.go cpu_mean_70 CPUUsageIncrease -property CPUUsageMean -threshold 70
go run cmd/swarm-autoconfigure/main.go cpu_mean_80 CPUUsageIncrease -property CPUUsageMean -threshold 80
go run cmd/swarm-autoconfigure/main.go cpu_mean_90 CPUUsageIncrease -property CPUUsageMean -threshold 90

# # ResponseTimeSimpleIncrease
# go run cmd/swarm-autoconfigure/main.go rtsi_95_300 ResponseTimeSimpleIncrease -property ResponseTimes95Percentile -value 300