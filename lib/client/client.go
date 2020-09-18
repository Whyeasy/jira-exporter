package client

import (
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/whyeasy/jira-exporter/internal"
	"github.com/whyeasy/jira-exporter/lib/jira"
)

//Stats struct is the list of expected results to export
type Stats struct {
	Projects *[]ProjectStats
	BugsMTTR *[]BugsHistogram
	Bugs     *[]Bugs
	TbMTTR   *[]TechDebtHistogram
	Tb       *[]TechDebt
}

//ExporterClient contains Jira information for connecting
type ExporterClient struct {
	jc             *jira.Client
	logsFilter     string
	jiraKeyInclude string
	jiraKeyExclude string
	interval       time.Duration
}

//New returns a new Client for connecting to Jira
func New(c internal.Config) *ExporterClient {

	convertedTime, _ := strconv.ParseInt(c.Interval, 10, 64)

	exporter := &ExporterClient{
		jc:             jira.NewClient(c.JiraAPIKey, c.JiraAPIUser, c.JiraURI),
		logsFilter:     c.JiraTbLabels,
		jiraKeyExclude: c.JiraKeyExclude,
		jiraKeyInclude: c.JiraKeyInclude,
		interval:       time.Duration(convertedTime),
	}

	exporter.startFetchData()

	return exporter
}

// CachedStats is to store scraped data for caching purposes.
var CachedStats *Stats = &Stats{
	Projects: &[]ProjectStats{},
	Bugs:     &[]Bugs{},
	BugsMTTR: &[]BugsHistogram{},
	Tb:       &[]TechDebt{},
	TbMTTR:   &[]TechDebtHistogram{},
}

//GetStats retrieves data from API to create metrics from.
func (c *ExporterClient) GetStats() (*Stats, error) {

	return CachedStats, nil
}

func (c *ExporterClient) getData() error {

	projects, err := getProjects(c)
	if err != nil {
		return err
	}

	bugsMTTR, err := getBugsMTTR(c)
	if err != nil {
		return err
	}

	bugs, err := getBugs(c)
	if err != nil {
		return err
	}

	tbMTTR, err := getTbMTTR(c)
	if err != nil {
		return err
	}

	tb, err := getTb(c)
	if err != nil {
		return err
	}

	CachedStats = &Stats{
		Projects: projects,
		BugsMTTR: bugsMTTR,
		Bugs:     bugs,
		TbMTTR:   tbMTTR,
		Tb:       tb,
	}

	log.Info("New data retrieved")

	return nil
}

func (c *ExporterClient) startFetchData() {

	// Do initial call to have data from the start.
	go func() {
		err := c.getData()
		if err != nil {
			log.Error("Scraping failed.")
		}
	}()

	ticker := time.NewTicker(c.interval * time.Second)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				err := c.getData()
				if err != nil {
					log.Error("Scraping failed.")
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
