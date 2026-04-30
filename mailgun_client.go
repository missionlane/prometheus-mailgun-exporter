package main

import (
	"context"
	"time"

	"github.com/mailgun/mailgun-go/v5"
	"github.com/mailgun/mailgun-go/v5/mtypes"
)

// MailgunClient defines the interface for Mailgun API operations used by the exporter.
// This interface allows for mocking in tests.
//
//go:generate mockery
type MailgunClient interface {
	// GetDomains returns all domains for the account
	GetDomains(ctx context.Context) ([]mtypes.Domain, error)
	// GetMetrics returns metrics for a specific domain
	GetMetrics(ctx context.Context, domain string) (*mtypes.Metrics, error)
}

// mailgunClientWrapper wraps the mailgun.Client to implement MailgunClient interface
type mailgunClientWrapper struct {
	client *mailgun.Client
}

// NewMailgunClientWrapper creates a new wrapper around the mailgun.Client
func NewMailgunClientWrapper(client *mailgun.Client) MailgunClient {
	return &mailgunClientWrapper{client: client}
}

func (w *mailgunClientWrapper) GetDomains(ctx context.Context) ([]mtypes.Domain, error) {
	it := w.client.ListDomains(nil)

	var page, result []mtypes.Domain
	for it.Next(ctx, &page) {
		result = append(result, page...)
	}

	if it.Err() != nil {
		return nil, it.Err()
	}
	return result, nil
}

func (w *mailgunClientWrapper) GetMetrics(ctx context.Context, domain string) (*mtypes.Metrics, error) {
	opts := mailgun.MetricsOptions{
		End:               mtypes.RFC2822Time(time.Now().UTC()),
		Duration:          "4h",
		Resolution:        mtypes.ResolutionHour,
		IncludeAggregates: true,
		Filter: mtypes.MetricsFilterPredicateGroup{
			BoolGroupAnd: []mtypes.MetricsFilterPredicate{{
				Attribute:     "domain",
				Comparator:    "=",
				LabeledValues: []mtypes.MetricsLabeledValue{{Label: domain, Value: domain}},
			}},
		},
		Metrics: []string{
			"accepted_incoming_count",
			"accepted_outgoing_count",
			"clicked_count",
			"complained_count",
			"delivered_http_count",
			"delivered_smtp_count",
			"permanent_failed_count",
			"temporary_failed_count",
			"opened_count",
			"stored_count",
			"unsubscribed_count",
			"hard_bounces_count",
			"soft_bounces_count",
			"delayed_bounce_count",
			"suppressed_bounces_count",
			"suppressed_complaints_count",
			"suppressed_unsubscribed_count",
			"temporary_failed_esp_block_count",
		},
	}

	iter, err := w.client.ListMetrics(opts)
	if err != nil {
		return nil, err
	}

	var resp mtypes.MetricsResponse
	iter.Next(ctx, &resp)
	if iter.Err() != nil {
		return nil, iter.Err()
	}

	return &resp.Aggregates.Metrics, nil
}
