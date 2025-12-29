package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	rpsCounter *prometheus.CounterVec
}

func NewMetrics() *Metrics {
	m := &Metrics{
		rpsCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests successfully processed",
			},
			[]string{"statusCode", "route"},
		),
	}

	prometheus.MustRegister(m.rpsCounter)

	return m
}

func (m *Metrics) IncRequest(route string, statusCode string) {
	m.rpsCounter.WithLabelValues(statusCode, route).Inc()
}
