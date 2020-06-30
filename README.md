![build](https://github.com/Whyeasy/jira-exporter/workflows/build/badge.svg)
![status-badge](https://goreportcard.com/badge/github.com/Whyeasy/jira-exporter)
![Github go.mod Go version](https://img.shields.io/github/go-mod/go-version/Whyeasy/jira-exporter)

# jira-exporter

A Prometheus Exporter for Jira Cloud

Currently this exporter retrieves the following metrics:

- Project Info within Jira (Key, Name and ID) `jira_project_info`
- Histogram of MTTR for bugs that are done. `jira_bugs_mttr`
- Bugs total(Project Key, Done) `jira_bugs_total`

## Requirements

### Required

Provide your Jira cloud URI; `--jiraURI <string>` or as env variable `JIRA_URI`

Provide a Jira API Key; `--jiraAPIKey` or as env variable `JIRA_API_KEY`

Provide the Jira user who created the API key; `--jiraAPIUser` or as env variable `JIRA_API_USER`

Provide a comma separated string with labels that mark stories with technical debt; `--jiraTbLabels` or as env variable `JIRA_TB_LABELS`

### Optional

Change listening port of the exporter; `--listenAddress <string>` or as env variable `LISTEN_ADDRESS`. Default = `8080`

Change listening path of the exporter; `--listenPath <string>` or as env variable `LISTEN_PATH`. Default = `/metrics`

To include or exclude projects from the Bugs metrics, please provide a comma separated string with the project keys. Please only provide 1.

Either with `--jiraKeyExclude <string>` or `--jiraKeyInclude <string>`. You can also provide it via env variables `JIRA_KEY_EXCL` or `JIRA_KEY_INCL`.
