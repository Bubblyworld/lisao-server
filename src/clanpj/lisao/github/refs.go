package github

type Object struct {
	Type string `json:"type"`
	Sha  string `json:"sha"`
	Url  string `json:"url"`
}

type Ref struct {
	Ref    string `json:"ref"`
	NodeID string `json:"node_id"`
	Url    string `json:"url"`
}

type Refs []Ref

type Ref404 struct {
	Message          string `json:"message"`
	DocumentationUrl string `json:"documentation_url"`
}

// GetRefs returns all refs for the client's configured repo.
func (c *Client) GetRefs() (Refs, error) {
	apiUrl := c.repoApiUrl() + "/git/refs"
	req, err := c.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	statusMap := StatusMap{
		200: &Refs{},
	}

	refs, err := c.DoToJSON(req, statusMap)
	if err != nil {
		return nil, err
	}

	return *refs.(*Refs), nil
}

// TODO(guy) add support for fetching just tags/commits and individual.
