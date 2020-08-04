package workload

// ResponseTimeCollector ....
type ResponseTimeCollector interface {
	GetResponseTimes(string) (map[string][]float64, error) // works with miliseconds
}

// RequestCountCollector ....
type RequestCountCollector interface {
	GetRequestCount(string) (map[string]int, error) // works with miliseconds
}
