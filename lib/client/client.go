package client

import (
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
}

//New returns a new Client for connecting to Jira
func New(c internal.Config) *ExporterClient {
	return &ExporterClient{
		jc:             jira.NewClient(c.JiraAPIKey, c.JiraAPIUser, c.JiraURI),
		logsFilter:     c.JiraTbLabels,
		jiraKeyExclude: c.JiraKeyExclude,
		jiraKeyInclude: c.JiraKeyInclude,
	}
}

//GetStats retrieves data from API to create metrics from.
func (c *ExporterClient) GetStats() (*Stats, error) {

	projects, err := getProjects(c)
	if err != nil {
		return nil, err
	}

	bugsMTTR, err := getBugsMTTR(c)
	if err != nil {
		return nil, err
	}

	bugs, err := getBugs(c)
	if err != nil {
		return nil, err
	}

	tbMTTR, err := getTbMTTR(c)
	if err != nil {
		return nil, err
	}

	tb, err := getTb(c)
	if err != nil {
		return nil, err
	}

	return &Stats{
		Projects: projects,
		BugsMTTR: bugsMTTR,
		Bugs:     bugs,
		TbMTTR:   tbMTTR,
		Tb:       tb,
	}, nil
}
