package workload

// ResponseTimeCollector ....
type ResponseTimeCollector interface {
	GetResponseTimes(string) ([]float64, error) // works with miliseconds
}

// RequestCountCollector ....
type RequestCountCollector interface {
	GetRequestCount(string) (int, error) // works with miliseconds
}
