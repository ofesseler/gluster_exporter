package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"

	"github.com/ofesseler/gluster_exporter/collector"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	program = "gluster_exporter"
)

var (
	glusterPath = kingpin.Flag(
		"gluster_executable_path",
		"Path to gluster executable",
	).Default("").String()
	glusterVolumes = kingpin.Flag(
		"volumes",
		fmt.Sprintf("Comma separated volume names: vol1,vol2,vol3. Default is '%v' to scrape all metrics", "_all"),
	).Default("_all").String()
	listenAddress = kingpin.Flag(
		"listen-address",
		"Address to listen on web interface and telemetry.",
	).Default(":9189").String()
	metricPath = kingpin.Flag(
		"metrics-path",
		"Path under which to expose metrics.",
	).Default("/metrics").String()
	profile = kingpin.Flag(
		"profile",
		"When profiling reports in gluster are enabled, set '--profile' to get more metrics, Default disable",
	).Default("false").Bool()
	quota = kingpin.Flag(
		"quota",
		"When quota in gluster are enabled and configured, set '--quota' to get quota metrics, Default disable",
	).Default("false").Bool()
	mount = kingpin.Flag(
		"mount",
		"set '--mount' to get mount metrics, Default disable",
	).Default("false").Bool()
	peer = kingpin.Flag(
		"peer",
		"set '--peer' to get peer metrics, Default disable",
	).Default("false").Bool()
	authUser = kingpin.Flag(
		"auth.user",
		"Username for basic auth.",
	).Default("").String()
	authPasswd = kingpin.Flag(
		"auth.passwd",
		"Password for basic auth.",
	).Default("").String()
)

func init() {
	prometheus.MustRegister(version.NewCollector("gluster_exporter"))
}

type basicAuthHandler struct {
	handler  http.HandlerFunc
	user     string
	password string
}

func (h *basicAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, password, ok := r.BasicAuth()
	if !ok || password != h.password || user != h.user {
		w.Header().Set("WWW-Authenticate", "Basic realm=\"metrics\"")
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	h.handler(w, r)
	return
}

func hasUserAndPassword() bool {
	return *authUser != "" && *authPasswd != ""
}

func filter(filters map[string]bool, name string, flag bool) bool {
	if len(filters) > 0 {
		return flag && filters[name]
	}
	return flag
}

func handler(w http.ResponseWriter, r *http.Request) {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("While trying to get Hostname error happened: %v", err)
	}

	params := r.URL.Query()["collect[]"]
	log.Debugln("collect query:", params)

	// prometheus query with params in prometheus.yml
	// like
	// - job_name: 'mysql performance'
	//	 scrape_interval: 1m
	//	 static_configs:
	//	   - targets:
	//	     - '192.168.1.2:9104'
	//   params:
	//	   collect[]:
	//	     - profile
	//	     - quota
	filters := make(map[string]bool)
	if len(params) > 0 {
		for _, param := range params {
			filters[param] = true
		}
	}

	collect := collector.Collect{
		Profile: filter(filters, "profile", *profile),
		Quota:   filter(filters, "quota", *quota),
		Mount:   filter(filters, "mount", *mount),
		Peer:    filter(filters, "peer", *peer),
	}

	registry := prometheus.NewRegistry()
	registry.MustRegister(collector.New(hostname, *glusterPath, *glusterVolumes, collect))

	gatherers := prometheus.Gatherers{
		prometheus.DefaultGatherer,
		registry,
	}

	handler := promhttp.HandlerFor(gatherers, promhttp.HandlerOpts{})
	if hasUserAndPassword() {
		handler = &basicAuthHandler{
			handler:  promhttp.HandlerFor(gatherers, promhttp.HandlerOpts{}).ServeHTTP,
			user:     *authUser,
			password: *authPasswd,
		}
		log.Info("Use AUTH")
	}

	handler.ServeHTTP(w, r)
}

func main() {

	log.AddFlags(kingpin.CommandLine)
	kingpin.Version(version.Print(program))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	// landingPage contains the HTML served at '/'.
	var landingPage = []byte(`<html>
	<head><title>GlusterFS Exporter</title></head>
	<body>
	<h1>GlusterFS Exporter</h1>
	<p><a href='` + *metricPath + `'>Metrics</a></p>
	</body>
	</html>
	`)

	log.Infoln("Starting gluster_exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	log.Infof("profile: %t, quota: %t, mount: %t, peer: %t", *profile, *quota, *mount, *peer)

	http.HandleFunc(*metricPath, prometheus.InstrumentHandlerFunc("metrics", handler))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(landingPage)
	})

	log.Infoln("Listening on", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
