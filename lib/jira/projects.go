package jira

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//ProjectListResponse is the struct to unmarshal the Jira API Response
type ProjectListResponse struct {
	Self       string `json:"self"`
	NextPage   string `json:"nextPage"`
	MaxResults int    `json:"maxResults"`
	StartAt    int    `json:"startAt"`
	Total      int    `json:"total"`
	IsLast     bool   `json:"isLast"`
	Values     []struct {
		Self       string `json:"self"`
		ID         string `json:"id"`
		Key        string `json:"key"`
		Name       string `json:"name"`
		AvatarUrls struct {
			Four8X48  string `json:"48x48"`
			Two4X24   string `json:"24x24"`
			One6X16   string `json:"16x16"`
			Three2X32 string `json:"32x32"`
		} `json:"avatarUrls"`
		ProjectCategory struct {
			Self        string `json:"self"`
			ID          string `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"projectCategory"`
		Simplified bool   `json:"simplified"`
		Style      string `json:"style"`
		Insight    struct {
			TotalIssueCount     int    `json:"totalIssueCount"`
			LastIssueUpdateTime string `json:"lastIssueUpdateTime"`
		} `json:"insight"`
	} `json:"values"`
}

//ListProjects requests all projects from Jira, paginated.
func (c *Client) ListProjects() (*ProjectListResponse, error) {

	url := fmt.Sprintf("%s/rest/api/3/project/search", c.jiraURI)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(c.jiraAPIUser, c.jiraAPIKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var projects ProjectListResponse
	err = json.Unmarshal(body, &projects)
	if err != nil {
		return nil, err
	}

	return &projects, nil

}
