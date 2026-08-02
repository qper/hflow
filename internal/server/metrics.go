package server

// Metrics provides simple Prometheus-style text output for deployment monitoring.
type Metrics struct{}

func NewMetrics() *Metrics { return &Metrics{} }

func (m *Metrics) Render() string {
	return `# HELP habitflow_up Application availability
# TYPE habitflow_up gauge
habitflow_up 1
# HELP habitflow_requests_total Total HTTP requests
# TYPE habitflow_requests_total counter
habitflow_requests_total 0
`
}
