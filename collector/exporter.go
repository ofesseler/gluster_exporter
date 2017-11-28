package collector

import (
	"time"
	"strings"

	"github.com/prometheus/common/log"
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics name parts.
const (
	// Default all Volumes
	allVolumes = "_all"

	// Namespace
	namespace = "gluster"
	// Subsystem(s).
	exporter = "exporter"
)

// Metric descriptors.
var (
	scrapeDurationDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, exporter, "collector_duration_seconds"),
		"Collector time duration.",
		[]string{"collector"}, nil,
	)
)

// Collect defines which metrics we should collect
type Collect struct {
	Base 		bool
	Profile		bool
	Quota 		bool
	Mount 		bool
	Peer 		bool
}

type Exporter struct {
	hostname 		string
	glusterPath 	string
	volumes 		[]string
	collect 		Collect
	error 			prometheus.Gauge
	totalScrapes	prometheus.Counter
	scrapeErrors 	*prometheus.CounterVec
	glusterUp		prometheus.Gauge
}

// returns a new GlusterFS exporter
func New(hostname string, glusterPath string, volumeString string, collect Collect) *Exporter {

	gfsPath, err := getGlusterBinary(glusterPath)
	if err != nil {
		log.Errorf("Given Gluster path %v has err: %v", glusterPath, err)
	}

	volumes := strings.Split(volumeString, ",")
	if len(volumes) < 1 {
		log.Infof("No volumes given. Proceeding without volume information. Volumes: %v", volumeString)
	}

	return &Exporter{
		hostname: hostname,
		glusterPath: gfsPath,
		volumes: volumes,
		collect: collect,
		totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: 	namespace,
			Subsystem: 	exporter,
			Name: 		"scrapes_total",
			Help:		"Total number of times GlusterFS was scraped for metrics.",
		}),
		scrapeErrors: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: 	namespace,
			Subsystem: 	exporter,
			Name: 		"scrape_errors_total",
			Help: 		"Total number of times an error occurred scraping a GlusterFS.",
		}, []string{"collector"}),
		error: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: exporter,
			Name:      "last_scrape_error",
			Help:      "Whether the last scrape of metrics from GlusterFS resulted in an error (1 for error, 0 for success).",
		}),
		glusterUp: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "up",
			Help:      "Whether the GlusterFS server is up.",
		}),
	}
}

// Describe implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {

	metricCh := make(chan prometheus.Metric)
	doneCh := make(chan struct{})

	go func() {
		for m := range metricCh {
			ch <- m.Desc()
		}
		close(doneCh)
	}()

	e.Collect(metricCh)
	close(metricCh)
	<-doneCh
}

// Collect implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.scrape(ch)

	ch <- e.totalScrapes
	ch <- e.error
	e.scrapeErrors.Collect(ch)
	ch <- e.glusterUp
}

func (e *Exporter) scrape(ch chan<- prometheus.Metric) {
	e.totalScrapes.Inc()
	var err error

	scrapeTime := time.Now()

	// if can get volume info, glusterFS is UP(1), or Down(0)
	_, err = ExecVolumeInfo()
	if err != nil {
		e.glusterUp.Set(0)
	}
	e.glusterUp.Set(1)

	// default collect volume info as Base Metrics
	e.collect.Base = true

	if e.collect.Base {
		// Base Gluster Info Scrape
		scrapeTime := time.Now()
		if err = ScrapeGlobalVolumeStatus(e.volumes, allVolumes, ch); err != nil {
			log.Errorln("Error scraping for collect.global_status:", err)
			e.scrapeErrors.WithLabelValues("collect.global_status").Inc()
			e.error.Set(1)
		}
		ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, time.Since(scrapeTime).Seconds(), "collect.global_status")
	}

	// Peer Info Scrape
	if e.collect.Peer {
		scrapeTime = time.Now()
		if err = ScrapePeerStatus(ch); err != nil {
			log.Errorln("Error scraping for collect.peer_status: ", err)
			e.scrapeErrors.WithLabelValues("collect.peer_status").Inc()
			e.error.Set(1)
		}
		ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, time.Since(scrapeTime).Seconds(), "collect.peer_status")
	}

	// Mount Scrape
	if e.collect.Mount {
		scrapeTime = time.Now()
		if err = ScrapeVolumeMountStatus(e.scrapeErrors, ch); err != nil {
			log.Errorln("Error scraping for collect.mount_status:", err)
			e.scrapeErrors.WithLabelValues("collect.mount_status").Inc()
			e.error.Set(1)
		}
		ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, time.Since(scrapeTime).Seconds(), "collect.mount_status")
	}

	// Profile Scrape
	if e.collect.Profile {
		scrapeTime = time.Now()
		if err = ScrapeProfileStatus(e.volumes, allVolumes, e.hostname, e.scrapeErrors, ch); err != nil {
			log.Errorln("Error scraping for collect.profile_status:", err)
			e.scrapeErrors.WithLabelValues("collect.profile_status").Inc()
			e.error.Set(1)
		}
		ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, time.Since(scrapeTime).Seconds(), "collect.profile_status")
	}

	// Quota Scrape
	if e.collect.Quota {
		scrapeTime = time.Now()
		if err = ScrapeQuotaStatus(e.volumes, allVolumes, e.scrapeErrors, ch); err != nil {
			log.Errorln("Error scraping for collect.quota_status:", err)
			e.scrapeErrors.WithLabelValues("collect.quota_status").Inc()
			e.error.Set(1)
		}
		ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, time.Since(scrapeTime).Seconds(), "collect.quota_status")
	}

}
