package client

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

//TechDebt struct is the data we want to have from Jira
type TechDebt struct {
	ProjectKey string
	Total      float64
	Done       bool
}

//TechDebtHistogram struct for holding jira data we want
type TechDebtHistogram struct {
	Count      uint64
	Sum        float64
	Bucket     map[float64]uint64
	ProjectKey string
}

func getTbMTTR(c *ExporterClient) (*[]TechDebtHistogram, error) {

	var jql string

	switch {
	case c.jiraKeyExclude != "":
		jql = fmt.Sprintf("issuetype != Epic AND resolutiondate != null AND project NOT IN (%s) AND labels IN (%s)", c.jiraKeyExclude, c.logsFilter)
	case c.jiraKeyInclude != "":
		jql = fmt.Sprintf("issuetype != Epic AND resolutiondate != null AND project NOT IN (%s) AND labels IN (%s)", c.jiraKeyInclude, c.logsFilter)
	default:
		jql = fmt.Sprintf("issuetype != Epic AND resolutiondate != null AND labels IN (%s)", c.logsFilter)
	}

	apiResults, err := c.jc.DoExpression(
		"issues.reduce((result, issue) => result.set(issue.project.key, (result[issue.project.key] || {}).set(issue.key, ((result[issue.project.key] || {})[issue.key] || 0) + issue.resolutionDate.getTime() - issue.created.getTime())), new Map())",
		jql)
	if err != nil {
		return nil, err
	}

	issueMTTRByProjects := make(map[string][]float64)

	// Paginate through list of expression response
	for _, apiResult := range apiResults {

		// Loop through issues in expression value response.
		for project, issues := range apiResult.Value {

			for _, issue := range issues {

				issueMTTRByProjects[project] = append(issueMTTRByProjects[project], issue)

			}

		}
	}

	promBucket := prometheus.ExponentialBuckets(7200, 2, 10)

	var results []TechDebtHistogram

	for project, issues := range issueMTTRByProjects {

		var sum float64

		bucket := make(map[float64]uint64)

		for _, value := range issues {
			//Convert to seconds.
			value = value / 1000
			sum += value

			for _, x := range promBucket {
				if value <= x {
					bucket[x]++
				}
			}
		}

		results = append(results, TechDebtHistogram{
			ProjectKey: project,
			Count:      uint64(len(issues)),
			Sum:        sum,
			Bucket:     bucket,
		})
	}

	return &results, nil
}

func getTb(c *ExporterClient) (*[]TechDebt, error) {

	var jql string

	switch {
	case c.jiraKeyExclude != "":
		jql = fmt.Sprintf("issuetype != Epic AND project NOT IN (%s) AND labels IN (%s)", c.jiraKeyExclude, c.logsFilter)
	case c.jiraKeyInclude != "":
		jql = fmt.Sprintf("issuetype != Epic AND project NOT IN (%s) AND labels IN (%s)", c.jiraKeyInclude, c.logsFilter)
	default:
		jql = fmt.Sprintf("issuetype != Epic AND labels IN (%s)", c.logsFilter)
	}

	apiResults, err := c.jc.DoExpression(
		"issues.reduce((result, issue) => result.set(issue.project.key, (result[issue.project.key] || {}).set(issue.status.name == 'Done' ? 'DONE' : 'NOT_DONE', ((result[issue.project.key] || {})[issue.status.name == 'Done' ? 'DONE' : 'NOT_DONE'] || 0) + 1)), new Map())",
		jql)
	if err != nil {
		return nil, err
	}

	type issueTypeCounter struct {
		doneCounter    float64
		notDoneCounter float64
	}
	issueTypesByProjects := make(map[string]issueTypeCounter)

	for _, apiResult := range apiResults {

		for project, issues := range apiResult.Value {

			projectIssues, ok := issueTypesByProjects[project]
			if !ok {
				projectIssues = issueTypeCounter{}
				issueTypesByProjects[project] = projectIssues
			}
			projectIssues.doneCounter = projectIssues.doneCounter + issues["DONE"]
			projectIssues.notDoneCounter = projectIssues.notDoneCounter + issues["NOT_DONE"]

			issueTypesByProjects[project] = projectIssues
		}
	}

	var results []TechDebt

	for project, counter := range issueTypesByProjects {

		results = append(results, TechDebt{
			ProjectKey: project,
			Total:      counter.doneCounter,
			Done:       true,
		})
		results = append(results, TechDebt{
			ProjectKey: project,
			Total:      counter.notDoneCounter,
			Done:       false,
		})
	}

	return &results, nil
}
