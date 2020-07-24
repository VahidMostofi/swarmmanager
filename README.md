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