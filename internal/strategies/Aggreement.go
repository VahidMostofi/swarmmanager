package strategies

// Agreement ...
type Agreement struct {
	PropertyToConsider string // ResponseTimesMean,ResponseTimes90Percentile,ResponseTimes95Percentile,ResponseTimes99Percentile,RTToleranceIntervalUBoundc90p95
	Value              float64
}
