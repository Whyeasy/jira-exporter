package internal

//Config struct for holding config for exporter and Gitlab
type Config struct {
	ListenAddress  string
	ListenPath     string
	JiraURI        string
	JiraAPIKey     string
	JiraAPIUser    string
	JiraKeyExclude string
	JiraKeyInclude string
	JiraTbLabels   string
	Interval       string
}
