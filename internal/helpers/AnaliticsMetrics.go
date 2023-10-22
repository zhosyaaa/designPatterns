package helpers

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
)

type MetricsAdapter interface {
	RegisterCounter(name string) error
	IncrementCounter(name string)
	RegisterTimer(name string) error
	RecordTime(name string, duration float64)
}

type CounterMetricsSource struct {
	counters map[string]int
}

func NewCounterMetricsSource() *CounterMetricsSource {
	return &CounterMetricsSource{
		counters: make(map[string]int),
	}
}

func (c *CounterMetricsSource) RegisterCounter(name string) error {
	c.counters[name] = 0
	return nil
}

func (c *CounterMetricsSource) IncrementCounter(name string) {
	c.counters[name]++
}

type TimerMetricsSource struct {
	timers map[string]float64
}

func NewTimerMetricsSource() *TimerMetricsSource {
	return &TimerMetricsSource{
		timers: make(map[string]float64),
	}
}

func (t *TimerMetricsSource) RegisterTimer(name string) error {
	t.timers[name] = 0
	return nil
}

func (t *TimerMetricsSource) RecordTime(name string, duration float64) {
	t.timers[name] = duration
}

type PrometheusAdapter struct {
	counters map[string]prometheus.Counter
	timers   map[string]prometheus.Histogram
}

func NewPrometheusAdapter() *PrometheusAdapter {
	return &PrometheusAdapter{
		counters: make(map[string]prometheus.Counter),
		timers:   make(map[string]prometheus.Histogram),
	}
}

func (p *PrometheusAdapter) RegisterCounter(name string) error {
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: name,
		Help: fmt.Sprintf("Counter for %s", name),
	})
	if err := prometheus.Register(counter); err != nil {
		return err
	}
	p.counters[name] = counter
	return nil
}

func (p *PrometheusAdapter) IncrementCounter(name string) {
	if counter, ok := p.counters[name]; ok {
		counter.Inc()
	}
}

func (p *PrometheusAdapter) RegisterTimer(name string) error {
	timer := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: name,
		Help: fmt.Sprintf("Timer for %s", name),
	})
	if err := prometheus.Register(timer); err != nil {
		return err
	}
	p.timers[name] = timer
	return nil
}

func (p *PrometheusAdapter) RecordTime(name string, duration float64) {
	if timer, ok := p.timers[name]; ok {
		timer.Observe(duration)
	}
}
