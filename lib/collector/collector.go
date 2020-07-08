package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/whyeasy/jira-exporter/lib/client"
)

//Collector struct for holding Prometheus Desc and Exporter Client
type Collector struct {
	up     *prometheus.Desc
	client *client.ExporterClient

	projectInfo *prometheus.Desc

	bugsMTTR *prometheus.Desc
	bugs     *prometheus.Desc

	tbMTTR *prometheus.Desc
	tb     *prometheus.Desc
}

//New creates a new Collecotor with Prometheus descriptors
func New(c *client.ExporterClient) *Collector {
	log.Info("Creating collector")
	return &Collector{
		up:     prometheus.NewDesc("jira_up", "Whether Jira scrap was successful", nil, nil),
		client: c,

		projectInfo: prometheus.NewDesc("jira_project_info", "General information about projects", []string{"project_key", "project_id", "project_name"}, nil),

		bugsMTTR: prometheus.NewDesc("jira_bugs_mttr", "Histogram metric which contains the duration between reporting and completing a bug within Jira", []string{"project_key"}, nil),
		bugs:     prometheus.NewDesc("jira_bugs_total", "Total amount of bugs within a project", []string{"project_key", "done"}, nil),

		tbMTTR: prometheus.NewDesc("jira_tb_mttr", "Histogram metric which contains the duration between creating and finishing stories marked as technical debt within Jira", []string{"project_key"}, nil),
		tb:     prometheus.NewDesc("jira_tb_total", "Total amount of stories within a project marked as technical debt within Jira", []string{"project_key", "done"}, nil),
	}
}

//Describe the metrics that are collected.
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up

	ch <- c.projectInfo

	ch <- c.bugsMTTR
	ch <- c.bugs

	ch <- c.tbMTTR
}

//Collect gathers the metrics that are exported.
func (c *Collector) Collect(ch chan<- prometheus.Metric) {

	log.Info("Running scrape")

	if stats, err := c.client.GetStats(); err != nil {
		log.Error(err)
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 0)
	} else {
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 1)

		collectProjectInfo(c, ch, stats)

		collectBugsMTTR(c, ch, stats)

		collectBugs(c, ch, stats)

		collectTbMTTR(c, ch, stats)

		collectTb(c, ch, stats)

		log.Info("Scrape Complete")
	}
}

func collectProjectInfo(c *Collector, ch chan<- prometheus.Metric, stats *client.Stats) {
	for _, project := range *stats.Projects {
		ch <- prometheus.MustNewConstMetric(c.projectInfo, prometheus.GaugeValue, 1, project.Key, project.ID, project.Name)
	}
}

func collectBugsMTTR(c *Collector, ch chan<- prometheus.Metric, stats *client.Stats) {
	for _, histogram := range *stats.BugsMTTR {
		ch <- prometheus.MustNewConstHistogram(c.bugsMTTR, histogram.Count, histogram.Sum, histogram.Bucket, histogram.ProjectKey)
	}
}

func collectBugs(c *Collector, ch chan<- prometheus.Metric, stats *client.Stats) {
	for _, bugs := range *stats.Bugs {
		ch <- prometheus.MustNewConstMetric(c.bugs, prometheus.GaugeValue, bugs.Total, bugs.ProjectKey, strconv.FormatBool(bugs.Done))
	}
}

func collectTbMTTR(c *Collector, ch chan<- prometheus.Metric, stats *client.Stats) {
	for _, histogram := range *stats.TbMTTR {
		ch <- prometheus.MustNewConstHistogram(c.tbMTTR, histogram.Count, histogram.Sum, histogram.Bucket, histogram.ProjectKey)
	}
}

func collectTb(c *Collector, ch chan<- prometheus.Metric, stats *client.Stats) {
	for _, tb := range *stats.Tb {
		ch <- prometheus.MustNewConstMetric(c.tb, prometheus.GaugeValue, tb.Total, tb.ProjectKey, strconv.FormatBool(tb.Done))
	}
}
