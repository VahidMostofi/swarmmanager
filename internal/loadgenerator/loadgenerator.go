package loadgenerator

// LoadGenerator has the general behavior of any load generator
type LoadGenerator interface {
	Prepare(map[string]string) error
	Start(map[string]string) error
	Stop(map[string]string) error
	GetFeedback(map[string]string) (map[string]interface{}, error)
}
