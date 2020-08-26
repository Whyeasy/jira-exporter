package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//ExpressionResult struct holds the structure to unmarshal the Jira API Response
type ExpressionResult struct {
	Value map[string]map[string]float64 `json:"value"`
	Meta  struct {
		Issues struct {
			Jql struct {
				StartAt    int `json:"startAt"`
				MaxResults int `json:"maxResults"`
				Count      int `json:"count"`
				TotalCount int `json:"totalCount"`
			} `json:"jql"`
		} `json:"issues"`
	} `json:"meta"`
}

func (c *Client) expression(opt *ListOptions, expression string, jql string) (*ExpressionResult, error) {

	url := fmt.Sprintf("%s/rest/api/3/expression/eval", c.jiraURI)

	query := fmt.Sprintf(`{"expression": "%s", "context": {"issues": { "jql": { "maxResults": %d, "startAt": %d, "query": "%s" }}}}`, expression, opt.MaxResult, opt.StartAt, jql)

	jsonStr := []byte(query)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
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

	var result ExpressionResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

//DoExpression Executes the expression function and takes care of the pagination.
func (c *Client) DoExpression(max int, exp string, jql string) ([]*ExpressionResult, error) {
	var results []*ExpressionResult

	startAt := 0

	for {

		result, err := c.expression(&ListOptions{
			MaxResult: max,
			StartAt:   startAt,
		},
			exp,
			jql)

		if err != nil {
			return nil, err
		}

		if result.Meta.Issues.Jql.Count == 0 {
			break
		}

		results = append(results, result)

		startAt += 1000

	}

	return results, nil
}
