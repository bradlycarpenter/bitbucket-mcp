package bitbucket

import (
	"encoding/json"
	"fmt"
)

type Repository struct {
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPrivate   bool   `json:"is_private"`
	UpdatedOn   string `json:"updated_on"`
	Links       struct {
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
	} `json:"links"`
}

type repositoryList struct {
	Values []Repository `json:"values"`
	Next   string       `json:"next"`
}

func (c *Client) ListRepositories() ([]Repository, error) {
	path := fmt.Sprintf("/repositories/%s?pagelen=100&sort=-updated_on", c.workspace)
	body, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var list repositoryList
	if err := json.Unmarshal(body, &list); err != nil {
		return nil, err
	}
	return list.Values, nil
}
