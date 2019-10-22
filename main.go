package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/mailgun/mailgun-go/v3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	namespace = "mailgun"
)

// Exporter collects metrics from Mailgun's via their API.
type Exporter struct {
	mg                   *mailgun.MailgunImpl
	up                   *prometheus.Desc
	acceptedTotal        *prometheus.Desc
	clickedTotal         *prometheus.Desc
	complainedTotal      *prometheus.Desc
	deliveredTotal       *prometheus.Desc
	failedPermanentTotal *prometheus.Desc
	failedTemporaryTotal *prometheus.Desc
	openedTotal          *prometheus.Desc
	storedTotal          *prometheus.Desc
	unsubscribedTotal    *prometheus.Desc
	state                *prometheus.Desc
}

func prometheusDomainStatsDesc(metric string, help string) *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName(
			namespace,
			"domain_stats",
			fmt.Sprintf("%s_total", metric),
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
			fmt.Sprintf("%s_total", metric),
		),
		help,
		[]string{"name", "type"},
		nil,
	)
}

// NewExporter returns an initialized exporter.
func NewExporter() *Exporter {
	// NewMailgunFromEnv requires MG_DOMAIN to get set, even though we don't need it for listing all domains
	err := os.Setenv("MG_DOMAIN", "dummy")
	if err != nil {
		log.Fatalln(err)
	}

	mg, err := mailgun.NewMailgunFromEnv()
	APIBase, exists := os.LookupEnv("API_BASE")
	if exists {
		mg.SetAPIBase(APIBase)
	}
	if err != nil {
		log.Fatalln(err)
	}

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
		acceptedTotal: prometheusDomainStatsTypeDesc(
			"accepted",
			"Mailgun accepted the request for incoming/outgoing to send/forward the email and the message has been placed in queue.",
		),
		clickedTotal: prometheusDomainStatsDesc(
			"clicked",
			"The email recipient clicked on a link in the email.",
		),
		complainedTotal: prometheusDomainStatsDesc(
			"complained",
			"The email recipient clicked on the spam complaint button within their email client.",
		),
		deliveredTotal: prometheusDomainStatsTypeDesc(
			"delivered",
			"Mailgun sent the email via HTTP or SMTP and it was accepted by the recipient email server.",
		),
		failedPermanentTotal: prometheusDomainStatsTypeDesc(
			"failed_permanent",
			"All permanently failed emails. Includes bounce, delayed bounce, suppress bounce, suppress complaint, suppress unsubscribe",
		),
		failedTemporaryTotal: prometheusDomainStatsTypeDesc(
			"failed_temporary",
			"All temporary failed emails due to ESP block, that will be retried",
		),
		openedTotal: prometheusDomainStatsDesc(
			"opened",
			"The email recipient opened the email and enabled image viewing.",
		),
		storedTotal: prometheusDomainStatsDesc(
			"stored",
			"The email recipient opened the email and enabled image viewing.",
		),
		unsubscribedTotal: prometheusDomainStatsDesc(
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
	ch <- e.acceptedTotal
	ch <- e.clickedTotal
	ch <- e.complainedTotal
	ch <- e.deliveredTotal
	ch <- e.failedPermanentTotal
	ch <- e.failedTemporaryTotal
	ch <- e.openedTotal
	ch <- e.storedTotal
	ch <- e.unsubscribedTotal
	ch <- e.state
}

// Collect implements prometheus.Collector. It only initiates a scrape of
// Collins if no scrape is currently ongoing. If a scrape of Collins is
// currently ongoing, Collect waits for it to end and then uses its result to
// collect the metrics.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	domains, err := e.listDomains()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(e.up, prometheus.GaugeValue, 0)
		log.Errorf("Scrape of Mailgun's API failed: %s", err)
	}

	for _, info := range domains {
		domain := info.Name

		state := 1
		if info.State != "active" {
			state = 0
		}
		ch <- prometheus.MustNewConstMetric(e.state, prometheus.GaugeValue, float64(state), domain)

		stats, err := getStats(domain)
		if err != nil {
			log.Errorln(err)
		}

		var acceptedTotalIncoming = float64(0)
		var acceptedTotalOutgoing = float64(0)
		var clickedTotal = float64(0)
		var complainedTotal = float64(0)
		var deliveredHttpTotal = float64(0)
		var deliveredSmtpTotal = float64(0)
		var failedPermanentBounce = float64(0)
		var failedPermanentDelayedBounce = float64(0)
		var failedPermanentSuppressBounce = float64(0)
		var failedPermanentSuppressComplaint = float64(0)
		var failedPermanentSuppressUnsubscribe = float64(0)
		var failedTemporaryEspblock = float64(0)
		var openedTotal = float64(0)
		var storedTotal = float64(0)
		var unsubscribedTotal = float64(0)

		for _, stat := range stats {
			acceptedTotalIncoming += float64(stat.Accepted.Incoming)
			acceptedTotalOutgoing += float64(stat.Accepted.Outgoing)
			clickedTotal += float64(stat.Clicked.Total)
			complainedTotal += float64(stat.Complained.Total)
			complainedTotal += float64(stat.Complained.Total)
			deliveredHttpTotal += float64(stat.Delivered.Http)
			deliveredSmtpTotal += float64(stat.Delivered.Smtp)
			failedPermanentBounce += float64(stat.Failed.Permanent.Bounce)
			failedPermanentDelayedBounce += float64(stat.Failed.Permanent.DelayedBounce)
			failedPermanentSuppressBounce += float64(stat.Failed.Permanent.SuppressBounce)
			failedPermanentSuppressComplaint += float64(stat.Failed.Permanent.SuppressComplaint)
			failedPermanentSuppressUnsubscribe += float64(stat.Failed.Permanent.SuppressUnsubscribe)
			failedTemporaryEspblock += float64(stat.Failed.Temporary.Espblock)
			openedTotal += float64(stat.Opened.Total)
			storedTotal += float64(stat.Stored.Total)
			unsubscribedTotal += float64(stat.Unsubscribed.Total)
		}

		// Begin Accepted Total
		ch <- prometheus.MustNewConstMetric(
			e.acceptedTotal,
			prometheus.CounterValue,
			acceptedTotalIncoming,
			domain, "incoming",
		)
		ch <- prometheus.MustNewConstMetric(
			e.acceptedTotal,
			prometheus.CounterValue,
			acceptedTotalOutgoing,
			domain, "outgoing",
		)
		// End Accepted Total

		ch <- prometheus.MustNewConstMetric(e.clickedTotal, prometheus.CounterValue, clickedTotal, domain)
		ch <- prometheus.MustNewConstMetric(e.complainedTotal, prometheus.CounterValue, complainedTotal, domain)

		// Begin Delivered Total
		ch <- prometheus.MustNewConstMetric(
			e.deliveredTotal,
			prometheus.CounterValue,
			deliveredHttpTotal,
			domain, "http",
		)
		ch <- prometheus.MustNewConstMetric(
			e.deliveredTotal,
			prometheus.CounterValue,
			deliveredSmtpTotal,
			domain, "smtp",
		)
		// End Delivered Total

		// Begin Failed Permanent Total
		ch <- prometheus.MustNewConstMetric(
			e.failedPermanentTotal,
			prometheus.CounterValue,
			failedPermanentBounce,
			domain, "bounce",
		)
		ch <- prometheus.MustNewConstMetric(
			e.failedPermanentTotal,
			prometheus.CounterValue,
			failedPermanentDelayedBounce,
			domain, "delayed_bounce",
		)
		ch <- prometheus.MustNewConstMetric(
			e.failedPermanentTotal,
			prometheus.CounterValue,
			failedPermanentSuppressBounce,
			domain, "suppress_bounce",
		)
		ch <- prometheus.MustNewConstMetric(
			e.failedPermanentTotal,
			prometheus.CounterValue,
			failedPermanentSuppressComplaint,
			domain, "suppress_complaint",
		)
		ch <- prometheus.MustNewConstMetric(e.failedPermanentTotal, prometheus.CounterValue,
			failedPermanentSuppressUnsubscribe,
			domain, "suppress_unsubscribe",
		)
		// End Failed Permanent Total

		ch <- prometheus.MustNewConstMetric(
			e.failedTemporaryTotal,
			prometheus.CounterValue,
			failedTemporaryEspblock,
			domain, "esp_block",
		)

		ch <- prometheus.MustNewConstMetric(e.openedTotal, prometheus.CounterValue, openedTotal, domain)
		ch <- prometheus.MustNewConstMetric(e.storedTotal, prometheus.CounterValue, storedTotal, domain)
		ch <- prometheus.MustNewConstMetric(e.unsubscribedTotal, prometheus.CounterValue, unsubscribedTotal, domain)
	}

	ch <- prometheus.MustNewConstMetric(e.up, prometheus.GaugeValue, 1)
}

func (e *Exporter) listDomains() ([]mailgun.Domain, error) {
	it := e.mg.ListDomains(nil)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var page, result []mailgun.Domain
	for it.Next(ctx, &page) {
		result = append(result, page...)
	}

	if it.Err() != nil {
		return nil, it.Err()
	}
	return result, nil
}

func getStats(domain string) ([]mailgun.Stats, error) {
	// Since we are using NewMailgunFromEnv, we need to set MG_DOMAIN before fetching stats for said domain
	err := os.Setenv("MG_DOMAIN", domain)
	if err != nil {
		log.Errorln(err)
	}

	mg, err := mailgun.NewMailgunFromEnv()
	APIBase, exists := os.LookupEnv("API_BASE")
	if exists {
		mg.SetAPIBase(APIBase)
	}
	if err != nil {
		log.Errorln(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	return mg.GetStats(ctx, []string{
		"accepted", "clicked", "complained", "delivered", "failed", "opened", "stored", "unsubscribed",
	}, &mailgun.GetStatOptions{
		Duration: "240m",
	})
}

func main() {
	var (
		listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9616").String()
		metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	)

	log.AddFlags(kingpin.CommandLine)
	kingpin.Version(version.Print("prometheus-mailgun-exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	log.Infoln("Starting Mailgun exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	prometheus.MustRegister(NewExporter())
	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Mailgun Exporter</title></head>
            <body>
            <h1>Mailgun Exporter</h1>
            <p><a href='` + *metricsPath + `'>Metrics</a></p>
            </body>
            </html>`))
	})
	log.Infof("Starting HTTP server on listen address %s and metric path %s", *listenAddress, *metricsPath)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
