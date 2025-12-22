package main

import (
	"context"

	"github.com/mailgun/mailgun-go/v5"
	"github.com/mailgun/mailgun-go/v5/mtypes"
)

// MailgunClient defines the interface for Mailgun API operations used by the exporter.
// This interface allows for mocking in tests.
//
//go:generate mockery --name=MailgunClient --output=mocks --outpkg=mocks
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
		Duration: "4h",
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

	aggregatedMetrics := &mtypes.Metrics{}

	var resp mtypes.MetricsResponse
	for iter.Next(ctx, &resp) {
		if resp.Aggregates.Metrics.AcceptedIncomingCount != nil {
			aggregatedMetrics = &resp.Aggregates.Metrics
			break
		}
		for _, item := range resp.Items {
			addMetrics(aggregatedMetrics, &item.Metrics)
		}
	}

	if iter.Err() != nil {
		return nil, iter.Err()
	}

	return aggregatedMetrics, nil
}

// addMetrics adds the values from src to dst
func addMetrics(dst, src *mtypes.Metrics) {
	if src.AcceptedIncomingCount != nil {
		if dst.AcceptedIncomingCount == nil {
			dst.AcceptedIncomingCount = new(uint64)
		}
		*dst.AcceptedIncomingCount += *src.AcceptedIncomingCount
	}
	if src.AcceptedOutgoingCount != nil {
		if dst.AcceptedOutgoingCount == nil {
			dst.AcceptedOutgoingCount = new(uint64)
		}
		*dst.AcceptedOutgoingCount += *src.AcceptedOutgoingCount
	}
	if src.ClickedCount != nil {
		if dst.ClickedCount == nil {
			dst.ClickedCount = new(uint64)
		}
		*dst.ClickedCount += *src.ClickedCount
	}
	if src.ComplainedCount != nil {
		if dst.ComplainedCount == nil {
			dst.ComplainedCount = new(uint64)
		}
		*dst.ComplainedCount += *src.ComplainedCount
	}
	if src.DeliveredHTTPCount != nil {
		if dst.DeliveredHTTPCount == nil {
			dst.DeliveredHTTPCount = new(uint64)
		}
		*dst.DeliveredHTTPCount += *src.DeliveredHTTPCount
	}
	if src.DeliveredSMTPCount != nil {
		if dst.DeliveredSMTPCount == nil {
			dst.DeliveredSMTPCount = new(uint64)
		}
		*dst.DeliveredSMTPCount += *src.DeliveredSMTPCount
	}
	if src.PermanentFailedCount != nil {
		if dst.PermanentFailedCount == nil {
			dst.PermanentFailedCount = new(uint64)
		}
		*dst.PermanentFailedCount += *src.PermanentFailedCount
	}
	if src.TemporaryFailedCount != nil {
		if dst.TemporaryFailedCount == nil {
			dst.TemporaryFailedCount = new(uint64)
		}
		*dst.TemporaryFailedCount += *src.TemporaryFailedCount
	}
	if src.OpenedCount != nil {
		if dst.OpenedCount == nil {
			dst.OpenedCount = new(uint64)
		}
		*dst.OpenedCount += *src.OpenedCount
	}
	if src.StoredCount != nil {
		if dst.StoredCount == nil {
			dst.StoredCount = new(uint64)
		}
		*dst.StoredCount += *src.StoredCount
	}
	if src.UnsubscribedCount != nil {
		if dst.UnsubscribedCount == nil {
			dst.UnsubscribedCount = new(uint64)
		}
		*dst.UnsubscribedCount += *src.UnsubscribedCount
	}
	if src.HardBouncesCount != nil {
		if dst.HardBouncesCount == nil {
			dst.HardBouncesCount = new(uint64)
		}
		*dst.HardBouncesCount += *src.HardBouncesCount
	}
	if src.SoftBouncesCount != nil {
		if dst.SoftBouncesCount == nil {
			dst.SoftBouncesCount = new(uint64)
		}
		*dst.SoftBouncesCount += *src.SoftBouncesCount
	}
	if src.DelayedBounceCount != nil {
		if dst.DelayedBounceCount == nil {
			dst.DelayedBounceCount = new(uint64)
		}
		*dst.DelayedBounceCount += *src.DelayedBounceCount
	}
	if src.SuppressedBouncesCount != nil {
		if dst.SuppressedBouncesCount == nil {
			dst.SuppressedBouncesCount = new(uint64)
		}
		*dst.SuppressedBouncesCount += *src.SuppressedBouncesCount
	}
	if src.SuppressedComplaintsCount != nil {
		if dst.SuppressedComplaintsCount == nil {
			dst.SuppressedComplaintsCount = new(uint64)
		}
		*dst.SuppressedComplaintsCount += *src.SuppressedComplaintsCount
	}
	if src.SuppressedUnsubscribedCount != nil {
		if dst.SuppressedUnsubscribedCount == nil {
			dst.SuppressedUnsubscribedCount = new(uint64)
		}
		*dst.SuppressedUnsubscribedCount += *src.SuppressedUnsubscribedCount
	}
	if src.TemporaryFailedESPBlockCount != nil {
		if dst.TemporaryFailedESPBlockCount == nil {
			dst.TemporaryFailedESPBlockCount = new(uint64)
		}
		*dst.TemporaryFailedESPBlockCount += *src.TemporaryFailedESPBlockCount
	}
}
