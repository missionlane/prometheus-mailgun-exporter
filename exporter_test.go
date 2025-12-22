package main

import (
	"errors"
	"strings"
	"testing"

	"github.com/mailgun/mailgun-go/v5/mtypes"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"prometheus-mailgun-exporter/mocks"
)

// Helper function to create a uint64 pointer
func uint64Ptr(v uint64) *uint64 {
	return &v
}

func TestNewExporterWithClient(t *testing.T) {
	mockClient := mocks.NewMockMailgunClient(t)
	exporter := NewExporterWithClient(mockClient)

	assert.NotNil(t, exporter)
	assert.NotNil(t, exporter.mg)
	assert.NotNil(t, exporter.up)
	assert.NotNil(t, exporter.acceptedTotal)
}

func TestExporter_Describe(t *testing.T) {
	mockClient := mocks.NewMockMailgunClient(t)
	exporter := NewExporterWithClient(mockClient)

	ch := make(chan *prometheus.Desc, 20)
	exporter.Describe(ch)
	close(ch)

	var descs []*prometheus.Desc
	for desc := range ch {
		descs = append(descs, desc)
	}

	// Should have 11 descriptors
	assert.Len(t, descs, 11)
}

func TestExporter_Collect_Success(t *testing.T) {
	mockClient := mocks.NewMockMailgunClient(t)

	// Setup mock expectations
	domains := []mtypes.Domain{
		{Name: "example.com", State: "active"},
		{Name: "test.com", State: "disabled"},
	}
	mockClient.On("GetDomains", mock.Anything).Return(domains, nil)

	metrics := &mtypes.Metrics{
		AcceptedIncomingCount:        uint64Ptr(100),
		AcceptedOutgoingCount:        uint64Ptr(200),
		ClickedCount:                 uint64Ptr(50),
		ComplainedCount:              uint64Ptr(5),
		DeliveredHTTPCount:           uint64Ptr(80),
		DeliveredSMTPCount:           uint64Ptr(120),
		HardBouncesCount:             uint64Ptr(10),
		DelayedBounceCount:           uint64Ptr(2),
		SuppressedBouncesCount:       uint64Ptr(3),
		SuppressedComplaintsCount:    uint64Ptr(1),
		SuppressedUnsubscribedCount:  uint64Ptr(4),
		TemporaryFailedESPBlockCount: uint64Ptr(7),
		OpenedCount:                  uint64Ptr(150),
		StoredCount:                  uint64Ptr(25),
		UnsubscribedCount:            uint64Ptr(15),
	}
	mockClient.On("GetMetrics", mock.Anything, "example.com").Return(metrics, nil)
	mockClient.On("GetMetrics", mock.Anything, "test.com").Return(metrics, nil)

	exporter := NewExporterWithClient(mockClient)

	ch := make(chan prometheus.Metric, 100)
	exporter.Collect(ch)
	close(ch)

	var collectedMetrics []prometheus.Metric
	for m := range ch {
		collectedMetrics = append(collectedMetrics, m)
	}

	// Per domain: state(1) + accepted(2) + clicked(1) + complained(1) + delivered(2) +
	//             failed_permanent(5) + failed_temporary(1) + opened(1) + stored(1) + unsubscribed(1) = 16
	// 2 domains * 16 = 32, plus 1 'up' metric = 33
	// Should have metrics for both domains plus the 'up' metric
	assert.Len(t, collectedMetrics, 33)
}

func TestExporter_Collect_GetDomainsError(t *testing.T) {
	mockClient := mocks.NewMockMailgunClient(t)

	// Setup mock to return error
	mockClient.On("GetDomains", mock.Anything).Return(nil, errors.New("API error"))

	exporter := NewExporterWithClient(mockClient)

	ch := make(chan prometheus.Metric, 10)
	exporter.Collect(ch)
	close(ch)

	var collectedMetrics []prometheus.Metric
	for m := range ch {
		collectedMetrics = append(collectedMetrics, m)
	}

	// Should only have the 'up' metric with value 0
	assert.Len(t, collectedMetrics, 1)
}

func TestExporter_Collect_GetMetricsError(t *testing.T) {
	mockClient := mocks.NewMockMailgunClient(t)

	// Setup mock expectations
	domains := []mtypes.Domain{
		{Name: "example.com", State: "active"},
	}
	mockClient.On("GetDomains", mock.Anything).Return(domains, nil)
	mockClient.On("GetMetrics", mock.Anything, "example.com").Return(nil, errors.New("metrics error"))

	exporter := NewExporterWithClient(mockClient)

	ch := make(chan prometheus.Metric, 10)
	exporter.Collect(ch)
	close(ch)

	var collectedMetrics []prometheus.Metric
	for m := range ch {
		collectedMetrics = append(collectedMetrics, m)
	}

	// Should have state metric + up metric (metrics collection failed but didn't stop)
	assert.Len(t, collectedMetrics, 2)
}

func TestExporter_Collect_DomainState(t *testing.T) {
	mockClient := mocks.NewMockMailgunClient(t)

	domains := []mtypes.Domain{
		{Name: "active.com", State: "active"},
		{Name: "disabled.com", State: "disabled"},
	}
	mockClient.On("GetDomains", mock.Anything).Return(domains, nil)
	mockClient.On("GetMetrics", mock.Anything, mock.Anything).Return(&mtypes.Metrics{}, nil)

	exporter := NewExporterWithClient(mockClient)

	ch := make(chan prometheus.Metric, 100)
	exporter.Collect(ch)
	close(ch)

	// Collect all metrics
	var stateMetrics []prometheus.Metric
	for m := range ch {
		desc := m.Desc()
		if desc.String() != "" && strings.Contains(desc.String(), "domain_state") {
			stateMetrics = append(stateMetrics, m)
		}
	}

	// Should have 2 state metrics (one per domain)
	assert.Len(t, stateMetrics, 2)
}

func TestExporter_Collect_NilMetricValues(t *testing.T) {
	mockClient := mocks.NewMockMailgunClient(t)

	domains := []mtypes.Domain{
		{Name: "example.com", State: "active"},
	}
	mockClient.On("GetDomains", mock.Anything).Return(domains, nil)
	// Return metrics with all nil values
	mockClient.On("GetMetrics", mock.Anything, "example.com").Return(&mtypes.Metrics{}, nil)

	exporter := NewExporterWithClient(mockClient)

	ch := make(chan prometheus.Metric, 100)
	// Should not panic with nil metric values
	assert.NotPanics(t, func() {
		exporter.Collect(ch)
		close(ch)
	})
}
