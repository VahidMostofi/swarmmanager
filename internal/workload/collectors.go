package workload

// ResponseTimeCollector ....
type ResponseTimeCollector interface {
	GetRequestResponseTimes(string) ([]float64, error)
	GetServiceDetails(string) (map[string]map[string][]float64, error) // should be ServiceTimeDetails
}

// RequestCountCollector ....
type RequestCountCollector interface {
	GetRequestCount(string) (int, error)
	GetRequestNames() []string
}
