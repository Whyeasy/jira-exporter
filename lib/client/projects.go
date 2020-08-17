package client

import log "github.com/sirupsen/logrus"

//ProjectStats is the struct that holds the data we want from Jira
type ProjectStats struct {
	ID   string
	Key  string
	Name string
}

func getProjects(c *ExporterClient) (*[]ProjectStats, error) {

	var result []ProjectStats

	projects, err := c.jc.ListProjects("software")
	if err != nil {
		return nil, err
	}

	for _, project := range projects.Values {
		result = append(result, ProjectStats{
			ID:   project.ID,
			Key:  project.Key,
			Name: project.Name,
		})
	}

	log.Info("Amount of projects found: ", len(projects.Values))

	return &result, nil
}
