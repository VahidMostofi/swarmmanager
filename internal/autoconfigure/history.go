package autoconfigure

import "github.com/VahidMostofi/swarmmanager/internal/swarm"

// StackHistory ...
type StackHistory struct {
	Name     string        `yaml:"name"`
	Workload interface{}   `yaml:"workload"`
	History  []Information `yaml:"configs"`
}

// Information ...
type Information struct {
	Infomations map[string]ServiceInfo        `yaml:"info"`
	Specs       map[string]swarm.ServiceSpecs `yaml:"specs"`
}

// ServiceInfo ...
type ServiceInfo struct {
	Start int64 `yaml:"start,omitempty"` // miliseconds
	End   int64 `yaml:"end,omitempty"`   // miliseconds

	CPUUsageMean         float64 `yaml:"cpuUsageMean,omitempty"` // Normalized between 0-100
	CPUUsage70Percentile float64 `yaml:"cpuUsage70th,omitempty"`
	CPUUsage75Percentile float64 `yaml:"cpuUsage75th,omitempty"`
	CPUUsage80Percentile float64 `yaml:"cpuUsage80th,omitempty"`
	CPUUsage85Percentile float64 `yaml:"cpuUsage85th,omitempty"`
	CPUUsage90Percentile float64 `yaml:"cpuUsage90th,omitempty"`
	CPUUsage95Percentile float64 `yaml:"cpuUsage95th,omitempty"`
	CPUUsage99Percentile float64 `yaml:"cpuUsage99th,omitempty"`
	NumberOfCores        float64 `yaml:"numberOfCores,omitempty"` // shows number of cores available to each core
	ReplicaCount         int     `yaml:"replicaCount,omitempty"`

	RequestCount    int   `yaml:"requestCount,omitempty"`
	SubTracesCounts []int `yaml:"subTracesCount,omitempty"`

	ResponseTimesMean         float64 `yaml:"responseTimesMean,omitempty"`
	ResponseTimes90Percentile float64 `yaml:"responseTimes90th,omitempty"`
	ResponseTimes95Percentile float64 `yaml:"responseTimes95th,omitempty"`
	ResponseTimes99Percentile float64 `yaml:"responseTimes99th,omitempty"`

	SubTracesResponseTimeMean          []float64 `yaml:"subTracesResponseTimeMean,omitempty"`
	SubTracesResponseTimes90Percentile []float64 `yaml:"subTracesResponseTime90th,omitempty"`
	SubTracesResponseTimes95Percentile []float64 `yaml:"subTracesResponseTime95th,omitempty"`
	SubTracesResponseTimes99Percentile []float64 `yaml:"subTracesResponseTime99th,omitempty"`
}
