package main

import (
	"errors"
	"flag"
	"fmt"

	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	"github.com/whyeasy/jira-exporter/internal"
	"github.com/whyeasy/jira-exporter/lib/client"
	"github.com/whyeasy/jira-exporter/lib/collector"
)

var (
	config internal.Config
)

func init() {
	flag.StringVar(&config.ListenAddress, "listenAddress", os.Getenv("LISTEN_ADDRESS"), "Port address of exporter to run on")
	flag.StringVar(&config.ListenPath, "listenPath", os.Getenv("LISTEN_PATH"), "Path where metrics will be exposed")
	flag.StringVar(&config.JiraURI, "jiraURI", os.Getenv("JIRA_URI"), "URI to Gitlab instance to monitor")
	flag.StringVar(&config.JiraAPIKey, "jiraAPIKey", os.Getenv("JIRA_API_KEY"), "API Key to access the Gitlab instance")
	flag.StringVar(&config.JiraAPIUser, "jiraAPIUser", os.Getenv("JIRA_API_USER"), "API User which created the key to access the Gitlab instance")
	flag.StringVar(&config.JiraKeyExclude, "jiraKeyExclude", os.Getenv("JIRA_KEY_EXCL"), "Comma seperated string with Project Keys to exclude from queries.")
	flag.StringVar(&config.JiraKeyInclude, "jiraKeyInclude", os.Getenv("JIRA_KEY_INCL"), "Comma seperated string with Project Keys to include in queries.")
	flag.StringVar(&config.JiraTbLabels, "jiraTbLabels", os.Getenv("JIRA_TB_LABELS"), "Comma seperated string with Label(s) that define tech-debt issues")
}

func main() {
	if err := parseConfig(); err != nil {
		log.Error(err)
		flag.Usage()
		os.Exit(2)
	}

	log.Info("Starting Jira Exporter")

	client := client.New(config)
	coll := collector.New(client)
	prometheus.MustRegister(coll)

	log.Info("Start serving metrics")

	http.Handle(config.ListenPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`<html>
			<head><title>Jira Exporter</title></head>
			<body>
			<h1>Jira Exporter</h1>
			<p><a href="` + config.ListenPath + `">Metrics</a></p>
			</body>
			</html>`))
		if err != nil {
			log.Error(err)
		}
	})
	log.Fatal(http.ListenAndServe(":"+config.ListenAddress, nil))
}

func parseConfig() error {
	flag.Parse()
	required := []string{"jiraURI", "jiraAPIKey", "jiraAPIUser", "jiraTbLabels"}
	var err error
	flag.VisitAll(func(f *flag.Flag) {
		for _, r := range required {
			if r == f.Name && (f.Value.String() == "" || f.Value.String() == "0") {
				err = fmt.Errorf("%v is empty", f.Usage)
			}
		}
		if f.Name == "listenAddress" && (f.Value.String() == "" || f.Value.String() == "0") {
			err = f.Value.Set("8080")
			if err != nil {
				log.Error(err)
			}
		}
		if f.Name == "listenPath" && (f.Value.String() == "" || f.Value.String() == "0") {
			err = f.Value.Set("/metrics")
			if err != nil {
				log.Error(err)
			}
		}
	})

	if config.JiraKeyExclude != "" && config.JiraKeyInclude != "" {
		err = errors.New("Please provide only project keys to exclude OR include, not both")
	}

	return err
}
