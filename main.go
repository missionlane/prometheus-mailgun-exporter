package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/mailgun/mailgun-go/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	namespace = "mailgun"
)

// Exporter collects metrics from Mailgun's via their API.
type Exporter struct {
	mg              MailgunClient
	up              *prometheus.Desc
	accepted        *prometheus.Desc
	clicked         *prometheus.Desc
	complained      *prometheus.Desc
	delivered       *prometheus.Desc
	failedPermanent *prometheus.Desc
	failedTemporary *prometheus.Desc
	opened          *prometheus.Desc
	stored          *prometheus.Desc
	unsubscribed    *prometheus.Desc
	state           *prometheus.Desc
}

func prometheusDomainStatsDesc(metric string, help string) *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName(
			namespace,
			"domain_stats",
			metric,
		),
		help,
		[]string{"name"},
		nil,
	)
}

func prometheusDomainStatsTypeDesc(metric string, help string) *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName(
			namespace,
			"domain_stats",
			metric,
		),
		help,
		[]string{"name", "type"},
		nil,
	)
}

// NewExporter returns an initialized exporter.
func NewExporter() *Exporter {
	apiKey := os.Getenv("MG_API_KEY")
	if apiKey == "" {
		log.Fatal().Msg("MG_API_KEY environment variable is required")
	}

	mg := mailgun.NewMailgun(apiKey)
	if apiBase, exists := os.LookupEnv("API_BASE"); exists {
		if err := mg.SetAPIBase(apiBase); err != nil {
			log.Fatal().Err(err).Msgf("Failed to set API base: %v", err)
		}
	}

	return NewExporterWithClient(NewMailgunClientWrapper(mg))
}

// NewExporterWithClient returns an initialized exporter with a custom client.
// This is primarily used for testing with mock clients.
func NewExporterWithClient(mg MailgunClient) *Exporter {
	return &Exporter{
		mg: mg,
		up: prometheus.NewDesc(
			prometheus.BuildFQName(
				"mailgun",
				"",
				"up",
			),
			"'1' if the last scrape of Mailgun's API was successful, '0' otherwise.",
			nil,
			nil,
		),
		accepted: prometheusDomainStatsTypeDesc(
			"accepted",
			"Mailgun accepted the request for incoming/outgoing to send/forward the email and the message has been placed in queue.",
		),
		clicked: prometheusDomainStatsDesc(
			"clicked",
			"The email recipient clicked on a link in the email.",
		),
		complained: prometheusDomainStatsDesc(
			"complained",
			"The email recipient clicked on the spam complaint button within their email client.",
		),
		delivered: prometheusDomainStatsTypeDesc(
			"delivered",
			"Mailgun sent the email via HTTP or SMTP and it was accepted by the recipient email server.",
		),
		failedPermanent: prometheusDomainStatsTypeDesc(
			"failed_permanent",
			"All permanently failed emails. Includes bounce, delayed bounce, suppress bounce, suppress complaint, suppress unsubscribe",
		),
		failedTemporary: prometheusDomainStatsTypeDesc(
			"failed_temporary",
			"All temporary failed emails due to ESP block, that will be retried",
		),
		opened: prometheusDomainStatsDesc(
			"opened",
			"The email recipient opened the email and enabled image viewing.",
		),
		stored: prometheusDomainStatsDesc(
			"stored",
			"The email recipient opened the email and enabled image viewing.",
		),
		unsubscribed: prometheusDomainStatsDesc(
			"unsubscribed",
			"The email recipient clicked on the unsubscribe link.",
		),
		state: prometheus.NewDesc(
			prometheus.BuildFQName(
				namespace,
				"domain",
				"state",
			),
			"Is the domain active (1) or disabled (0)",
			[]string{"name"},
			nil,
		),
	}
}

// Describe implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.up
	ch <- e.accepted
	ch <- e.clicked
	ch <- e.complained
	ch <- e.delivered
	ch <- e.failedPermanent
	ch <- e.failedTemporary
	ch <- e.opened
	ch <- e.stored
	ch <- e.unsubscribed
	ch <- e.state
}

// Collect implements prometheus.Collector. It only initiates a scrape of
// Collins if no scrape is currently ongoing. If a scrape of Collins is
// currently ongoing, Collect waits for it to end and then uses its result to
// collect the metrics.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	domains, err := e.mg.GetDomains(ctx)
	if err != nil {
		ch <- prometheus.MustNewConstMetric(e.up, prometheus.GaugeValue, 0)
		log.Error().Err(err).Msgf("Scrape of Mailgun's API failed: %s", err)
		return
	}

	log.Debug().Int("count", len(domains)).Msg("Fetched domains")

	for _, info := range domains {
		domain := info.Name

		state := 1
		if info.State != "active" {
			state = 0
		}
		ch <- prometheus.MustNewConstMetric(e.state, prometheus.GaugeValue, float64(state), domain)

		metrics, err := e.mg.GetMetrics(ctx, domain)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to get metrics for domain %s", domain)
			continue
		}

		// Helper function to safely get uint64 pointer value
		getVal := func(p *uint64) float64 {
			if p == nil {
				return 0
			}
			return float64(*p)
		}

		if metrics.AcceptedIncomingCount == nil && metrics.AcceptedOutgoingCount == nil {
			log.Warn().Str("domain", domain).
				Msg("Metrics API returned no data; check that your Mailgun plan includes the analytics API, or verify that there was mail activity in the last 4h.")
		}

		ch <- prometheus.MustNewConstMetric(e.accepted, prometheus.GaugeValue, getVal(metrics.AcceptedIncomingCount), domain, "incoming")
		ch <- prometheus.MustNewConstMetric(e.accepted, prometheus.GaugeValue, getVal(metrics.AcceptedOutgoingCount), domain, "outgoing")

		ch <- prometheus.MustNewConstMetric(e.clicked, prometheus.GaugeValue, getVal(metrics.ClickedCount), domain)
		ch <- prometheus.MustNewConstMetric(e.complained, prometheus.GaugeValue, getVal(metrics.ComplainedCount), domain)

		ch <- prometheus.MustNewConstMetric(e.delivered, prometheus.GaugeValue, getVal(metrics.DeliveredHTTPCount), domain, "http")
		ch <- prometheus.MustNewConstMetric(e.delivered, prometheus.GaugeValue, getVal(metrics.DeliveredSMTPCount), domain, "smtp")

		ch <- prometheus.MustNewConstMetric(e.failedPermanent, prometheus.GaugeValue, getVal(metrics.HardBouncesCount), domain, "bounce")
		ch <- prometheus.MustNewConstMetric(e.failedPermanent, prometheus.GaugeValue, getVal(metrics.DelayedBounceCount), domain, "delayed_bounce")
		ch <- prometheus.MustNewConstMetric(e.failedPermanent, prometheus.GaugeValue, getVal(metrics.SuppressedBouncesCount), domain, "suppress_bounce")
		ch <- prometheus.MustNewConstMetric(e.failedPermanent, prometheus.GaugeValue, getVal(metrics.SuppressedComplaintsCount), domain, "suppress_complaint")
		ch <- prometheus.MustNewConstMetric(e.failedPermanent, prometheus.GaugeValue, getVal(metrics.SuppressedUnsubscribedCount), domain, "suppress_unsubscribe")
		ch <- prometheus.MustNewConstMetric(e.failedPermanent, prometheus.GaugeValue, getVal(metrics.SoftBouncesCount), domain, "soft_bounce")

		ch <- prometheus.MustNewConstMetric(e.failedTemporary, prometheus.GaugeValue, getVal(metrics.TemporaryFailedESPBlockCount), domain, "esp_block")

		ch <- prometheus.MustNewConstMetric(e.opened, prometheus.GaugeValue, getVal(metrics.OpenedCount), domain)
		ch <- prometheus.MustNewConstMetric(e.stored, prometheus.GaugeValue, getVal(metrics.StoredCount), domain)
		ch <- prometheus.MustNewConstMetric(e.unsubscribed, prometheus.GaugeValue, getVal(metrics.UnsubscribedCount), domain)
	}

	ch <- prometheus.MustNewConstMetric(e.up, prometheus.GaugeValue, 1)
}

func main() {
	var (
		listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9616").String()
		metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	)

	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	}

	kingpin.Version(version.Print("prometheus-mailgun-exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	log.Info().Msgf("Starting Mailgun exporter %v", version.Info())
	log.Info().Msgf("Build context %v", version.BuildContext())

	prometheus.MustRegister(NewExporter())
	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Mailgun Exporter</title></head>
            <body>
            <h1>Mailgun Exporter</h1>
            <p><a href='` + *metricsPath + `'>Metrics</a></p>
			<p><a href='/healthz'>Health</a></p>
            </body>
            </html>`))
	})
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	log.Info().Msgf("Starting HTTP server on listen address %s and metric path %s", *listenAddress, *metricsPath)

	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		log.Fatal().Err(err).Msgf("%v", err)
	}
}
