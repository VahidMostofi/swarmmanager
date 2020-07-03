package resource

import (
	"fmt"
	"time"
)

// Utilization contains information about resource utilization for each resource(container/service)
// CPUUtilizationsAtTime contains the CPU utilization at each timestamp (nano). The resource usage is in percent.
// Due to error in gathering info, it may be more than 100%
type Utilization struct {
	ResourceID            string
	CPUUtilizationsAtTime map[int64]float64
}

// NewResourceUtilization is the constructor
func NewResourceUtilization(resourceID string) *Utilization {
	return &Utilization{
		ResourceID:            resourceID,
		CPUUtilizationsAtTime: make(map[int64]float64),
	}
}

// AddCPUUsage handles a new recorded stat
func (ru *Utilization) AddCPUUsage(percent float64, timestamp int64) {
	ru.CPUUtilizationsAtTime[timestamp] = percent
}

// GetResourceRecordingRate returns ResourceRecordingRate (per second)
func (ru *Utilization) GetResourceRecordingRate() float64 {
	var count float64

	var start int64 = time.Now().UnixNano()
	var end int64 = 0
	for time := range ru.CPUUtilizationsAtTime {
		if time < start {
			start = time
		}
		if time > end {
			end = time
		}
	}
	count = float64(len(ru.CPUUtilizationsAtTime))
	return (count) / ((float64(end) - float64(start)) / 1e9)
}

// Print for debug purposes
func (ru *Utilization) Print() {
	fmt.Println(ru.ResourceID)
	for time, value := range ru.CPUUtilizationsAtTime {
		fmt.Println(time, value)
	}
}

// GetCPUPercentile ...
// func (ru *Utilization) GetCPUPercentile(percentile int) float64 {
// 	return 0
// }

// GetCPUAverage ...
func (ru *Utilization) GetCPUAverage(start, end int64) float64 {
	var total float64
	var count float64
	for timestamp, percent := range ru.CPUUtilizationsAtTime {
		if timestamp >= start && timestamp <= end {
			count++
			total += percent
		}
	}
	if count == 0 {
		return -1
	}
	return total / count
}
