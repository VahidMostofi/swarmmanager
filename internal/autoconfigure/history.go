package autoconfigure

// ServiceInfo ...
type ServiceInfo struct {
	Start int64 // miliseconds
	End   int64 // miliseconds

	CPUUsageMean         float64 // Normalized between 0-100
	CPUUsage90Percentile float64
	CPUUsage95Percentile float64
	CPUUsage99Percentile float64
	NumberOfCores        float64 // shows number of cores available to each core
	ReplicaCount         int

	RequestCount    int
	SubTracesCounts []int

	ResponseTimesMean         float64
	ResponseTimes90Percentile float64
	ResponseTimes95Percentile float64
	ResponseTimes99Percentile float64

	SubTracesResponseTimeMean          []float64
	SubTracesResponseTimes90Percentile []float64
	SubTracesResponseTimes95Percentile []float64
	SubTracesResponseTimes99Percentile []float64
}
