package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mailgun/mailgun-go/v5"
	"github.com/mailgun/mailgun-go/v5/mtypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMailgunClientWrapper_GetMetrics_RequestShape(t *testing.T) {
	var captured mtypes.MetricsRequest

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// assert.* is required here — require.* calls t.FailNow() which must not
		// be called from a goroutine other than the test goroutine.
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/v1/analytics/metrics", r.URL.Path)

		assert.NoError(t, json.NewDecoder(r.Body).Decode(&captured))

		accepted := uint64(42)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mtypes.MetricsResponse{
			Aggregates: mtypes.MetricsAggregates{
				Metrics: mtypes.Metrics{
					AcceptedIncomingCount: &accepted,
				},
			},
		})
	}))
	defer srv.Close()

	mg := mailgun.NewMailgun("test-key")
	require.NoError(t, mg.SetAPIBase(srv.URL))

	before := time.Now().UTC()
	metrics, err := NewMailgunClientWrapper(mg).GetMetrics(context.Background(), "example.com")
	after := time.Now().UTC()

	require.NoError(t, err)

	// End must be set to approximately now — the zero-value (year 0001) caused
	// the analytics endpoint to return empty data for every domain.
	endTime := time.Time(captured.End)
	assert.False(t, endTime.IsZero(), "End must not be zero-value time")
	assert.False(t, endTime.Before(before.Add(-time.Second)), "End should be >= call start")
	assert.False(t, endTime.After(after.Add(time.Second)), "End should be <= call end")

	assert.Equal(t, "4h", captured.Duration)
	assert.Equal(t, mtypes.ResolutionHour, captured.Resolution)
	assert.True(t, captured.IncludeAggregates, "IncludeAggregates must be true so server returns totals")

	require.Len(t, captured.Filter.BoolGroupAnd, 1)
	pred := captured.Filter.BoolGroupAnd[0]
	assert.Equal(t, "domain", pred.Attribute)
	assert.Equal(t, "=", pred.Comparator)
	require.Len(t, pred.LabeledValues, 1)
	assert.Equal(t, "example.com", pred.LabeledValues[0].Value)

	require.NotNil(t, metrics.AcceptedIncomingCount)
	assert.Equal(t, uint64(42), *metrics.AcceptedIncomingCount)
}

func TestMailgunClientWrapper_GetMetrics_EmptyResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mtypes.MetricsResponse{})
	}))
	defer srv.Close()

	mg := mailgun.NewMailgun("test-key")
	require.NoError(t, mg.SetAPIBase(srv.URL))

	metrics, err := NewMailgunClientWrapper(mg).GetMetrics(context.Background(), "example.com")

	require.NoError(t, err)
	// All fields nil — Collect handles this with the "no data" warning
	assert.Nil(t, metrics.AcceptedIncomingCount)
	assert.Nil(t, metrics.AcceptedOutgoingCount)
}

func TestMailgunClientWrapper_GetMetrics_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	mg := mailgun.NewMailgun("test-key")
	require.NoError(t, mg.SetAPIBase(srv.URL))

	_, err := NewMailgunClientWrapper(mg).GetMetrics(context.Background(), "example.com")
	assert.Error(t, err)
}
