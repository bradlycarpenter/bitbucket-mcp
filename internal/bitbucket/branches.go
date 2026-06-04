package bitbucket

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type Branch struct {
	Name   string `json:"name"`
	Target struct {
		Hash string `json:"hash"`
	} `json:"target"`
}

type branchList struct {
	Values []Branch `json:"values"`
	Next   string   `json:"next"`
}

func (c *Client) ListBranches(repoSlug string) ([]Branch, error) {
	path := fmt.Sprintf("/repositories/%s/%s/refs/branches?pagelen=100", c.workspace, repoSlug)
	body, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var list branchList
	if err := json.Unmarshal(body, &list); err != nil {
		return nil, err
	}
	return list.Values, nil
}

func (c *Client) CompareBranchesDiff(repoSlug, source, destination string) (string, error) {
	spec := url.PathEscape(source) + ".." + url.PathEscape(destination)
	path := fmt.Sprintf("/repositories/%s/%s/diff/%s", c.workspace, repoSlug, spec)
	body, err := c.getRaw(path)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (c *Client) CompareBranchesDiffstat(repoSlug, source, destination string) ([]DiffstatEntry, error) {
	spec := url.PathEscape(source) + ".." + url.PathEscape(destination)
	path := fmt.Sprintf("/repositories/%s/%s/diffstat/%s?pagelen=100", c.workspace, repoSlug, spec)
	body, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var list DiffstatList
	if err := json.Unmarshal(body, &list); err != nil {
		return nil, err
	}
	return list.Values, nil
}

func (c *Client) CompareBranchesCommits(repoSlug, source, destination string) ([]Commit, error) {
	path := fmt.Sprintf(
		"/repositories/%s/%s/commits?include=%s&exclude=%s&pagelen=100",
		c.workspace, repoSlug,
		url.QueryEscape(source), url.QueryEscape(destination),
	)
	body, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var list CommitList
	if err := json.Unmarshal(body, &list); err != nil {
		return nil, err
	}
	return list.Values, nil
}
