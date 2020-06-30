package client

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

//Bugs struct is the data we want to have from Jira
type Bugs struct {
	ProjectKey string
	Total      float64
	Done       bool
}

//BugsHistogram struct is the data we want to have from Jira
type BugsHistogram struct {
	Count      uint64
	Sum        float64
	Bucket     map[float64]uint64
	ProjectKey string
}

func getBugsMTTR(c *ExporterClient) (*[]BugsHistogram, error) {

	var jql string

	switch {
	case c.jiraKeyExclude != "":
		jql = fmt.Sprintf("issuetype = Bug AND resolutiondate != null AND project NOT IN (%s)", c.jiraKeyExclude)
	case c.jiraKeyInclude != "":
		jql = fmt.Sprintf("issuetype = Bug AND resolutiondate != null AND project NOT IN (%s)", c.jiraKeyInclude)
	default:
		jql = "issuetype = Bug AND resolutiondate != null"
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

	promBucket := prometheus.ExponentialBuckets(7200, 2, 8)

	var results []BugsHistogram

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

		results = append(results, BugsHistogram{
			ProjectKey: project,
			Count:      uint64(len(issues)),
			Sum:        sum,
			Bucket:     bucket,
		})
	}

	return &results, nil
}

func getBugs(c *ExporterClient) (*[]Bugs, error) {

	var jql string

	switch {
	case c.jiraKeyExclude != "":
		jql = fmt.Sprintf("issuetype = Bug AND project NOT IN (%s)", c.jiraKeyExclude)
	case c.jiraKeyInclude != "":
		jql = fmt.Sprintf("issuetype = Bug AND project NOT IN (%s)", c.jiraKeyInclude)
	default:
		jql = "issuetype = Bug"
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

	var results []Bugs

	for project, counter := range issueTypesByProjects {

		results = append(results, Bugs{
			ProjectKey: project,
			Total:      counter.doneCounter,
			Done:       true,
		})
		results = append(results, Bugs{
			ProjectKey: project,
			Total:      counter.notDoneCounter,
			Done:       false,
		})
	}

	return &results, nil
}
