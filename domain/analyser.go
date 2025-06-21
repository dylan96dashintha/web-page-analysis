package domain

type AnalyserRequest struct {
	Url []string `json:"url"`
}

type AnalysisResult struct {
	HTMLVersion  string         `json:"html_version"`
	Title        string         `json:"title"`
	Headings     map[string]int `json:"headings"`
	Link         Link           `json:"link"`
	Inaccessible int            `json:"inaccessible_links"`
	HasLoginForm bool           `json:"has_login_form"`
}

type Link struct {
	InternalLinks         int      `json:"internal_links"`
	ExternalLinks         int      `json:"external_links"`
	InaccessibleLinkCount int      `json:"inaccessible_link_count"`
	InaccessibleLink      []string `json:"inaccessible_link"`
}
