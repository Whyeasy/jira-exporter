package jira

import (
	"net/http"
	"strings"
)

//Client struct holds the data we need for the Jira Client
type Client struct {
	jiraURI     string
	client      *http.Client
	jiraAPIKey  string
	jiraAPIUser string
}

//ListOptions provides additional options during the request.
type ListOptions struct {
	StartAt   int
	MaxResult int
}

//NewClient Creates a new client to communicate with the Jira Cloud API.
func NewClient(api string, user string, baseURI string) *Client {

	if !strings.HasSuffix(baseURI, "/") {
		baseURI += "/"
	}

	return &Client{
		client:      &http.Client{},
		jiraAPIKey:  api,
		jiraAPIUser: user,
		jiraURI:     baseURI,
	}
}
