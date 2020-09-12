


```
go run cmd/swarm-autoconfigure/main.go <NAME> CPUUsageIncrease -property CPUUsage90Percentile -threshold 90
```

## How to prepare R-server
It is required to have the R-server running on localhost on port 6311

install these 2 packages:
```
install.packages("tolerance")
install.packages("Rserve",,"http://rforge.net")
```

load and start R-serve
```
library(Rserve)
Rserve(args="--no-save")
```
for more info refer to [R Serve docs](https://www.rforge.net/Rserve/doc.html).


## How to configure for a new Host
1. Update these values in config.yml
   * jaeger-host
   * host
   * available-cpu-count

## How to configure for a new system
1. Update these values in config.yml
   * system-name
   * docker-compose-file
   * service-count
   * stack-name
   * jaeger-root-service
   * jaeger-details-file-path
   * k6-script
   * services-to-monitor
2. Add ```jaeger-details-file``` in ```formulas/``` directory. Follow the rules in the README.md file in the ```formulas/``` directory.