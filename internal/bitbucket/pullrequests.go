package bitbucket

import (
	"encoding/json"
	"fmt"
)

type Author struct {
	DisplayName string `json:"display_name"`
	AccountID   string `json:"account_id"`
}

type Reviewer struct {
	DisplayName string `json:"display_name"`
	AccountID   string `json:"account_id"`
	Approved    bool   `json:"approved"`
}

type PR struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	State       string     `json:"state"`
	Author      Author     `json:"author"`
	Reviewers   []Reviewer `json:"reviewers"`
	CreatedOn   string     `json:"created_on"`
	UpdatedOn   string     `json:"updated_on"`
	Links       struct {
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
	} `json:"links"`
}

type PRList struct {
	Values []PR   `json:"values"`
	Next   string `json:"next"`
}

type Commit struct {
	Hash    string `json:"hash"`
	Message string `json:"message"`
	Author  struct {
		Raw  string `json:"raw"`
		User Author `json:"user"`
	} `json:"author"`
	Date string `json:"date"`
}

type CommitList struct {
	Values []Commit `json:"values"`
	Next   string   `json:"next"`
}

type DiffstatEntry struct {
	Status       string `json:"status"`
	LinesAdded   int    `json:"lines_added"`
	LinesRemoved int    `json:"lines_removed"`
	New          struct {
		Path string `json:"path"`
	} `json:"new"`
	Old struct {
		Path string `json:"path"`
	} `json:"old"`
}

type DiffstatList struct {
	Values []DiffstatEntry `json:"values"`
}

type CreatePRRequest struct {
	Title             string
	Description       string // optional, empty string = omit
	SourceBranch      string
	DestinationBranch string
}

func (c *Client) CreatePullRequest(repoSlug string, req CreatePRRequest) (*PR, error) {
	type branchRef struct {
		Name string `json:"name"`
	}
	type branchHolder struct {
		Branch branchRef `json:"branch"`
	}
	type payload struct {
		Title       string       `json:"title"`
		Description string       `json:"description,omitempty"`
		Source      branchHolder `json:"source"`
		Destination branchHolder `json:"destination"`
	}

	p := payload{
		Title:       req.Title,
		Description: req.Description,
		Source:      branchHolder{Branch: branchRef{Name: req.SourceBranch}},
		Destination: branchHolder{Branch: branchRef{Name: req.DestinationBranch}},
	}

	path := fmt.Sprintf("/repositories/%s/%s/pullrequests", c.workspace, repoSlug)
	body, err := c.post(path, p)
	if err != nil {
		return nil, err
	}
	var pr PR
	if err := json.Unmarshal(body, &pr); err != nil {
		return nil, err
	}
	return &pr, nil
}

type UpdatePRRequest struct {
	Title       string // optional, empty = don't change
	Description string // optional, empty = don't change
}

func (c *Client) UpdatePullRequest(repoSlug string, prID int, req UpdatePRRequest) (*PR, error) {
	body := map[string]any{}
	if req.Title != "" {
		body["title"] = req.Title
	}
	if req.Description != "" {
		body["description"] = req.Description
	}
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d", c.workspace, repoSlug, prID)
	respBody, err := c.put(path, body)
	if err != nil {
		return nil, err
	}
	var pr PR
	if err := json.Unmarshal(respBody, &pr); err != nil {
		return nil, err
	}
	return &pr, nil
}

func (c *Client) ListPullRequests(repoSlug, state string) ([]PR, error) {
	if state == "" {
		state = "OPEN"
	}
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests?state=%s&pagelen=50", c.workspace, repoSlug, state)
	body, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var list PRList
	if err := json.Unmarshal(body, &list); err != nil {
		return nil, err
	}
	return list.Values, nil
}

func (c *Client) GetPullRequest(repoSlug string, prID int) (*PR, error) {
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d", c.workspace, repoSlug, prID)
	body, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var pr PR
	if err := json.Unmarshal(body, &pr); err != nil {
		return nil, err
	}
	return &pr, nil
}

func (c *Client) ListPRCommits(repoSlug string, prID int) ([]Commit, error) {
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/commits?pagelen=50", c.workspace, repoSlug, prID)
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

func (c *Client) GetPRDiff(repoSlug string, prID int) (string, error) {
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/diff", c.workspace, repoSlug, prID)
	body, err := c.getRaw(path)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (c *Client) GetPRDiffstat(repoSlug string, prID int) ([]DiffstatEntry, error) {
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/diffstat", c.workspace, repoSlug, prID)
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

// --- Comments ---

type CommentContent struct {
	Raw string `json:"raw"`
}

type CommentInline struct {
	Path string `json:"path"`
	From *int   `json:"from"`
	To   *int   `json:"to"`
}

type CommentParent struct {
	ID int `json:"id"`
}

type Comment struct {
	ID        int            `json:"id"`
	User      Author         `json:"user"`
	Content   CommentContent `json:"content"`
	CreatedOn string         `json:"created_on"`
	UpdatedOn string         `json:"updated_on"`
	Deleted   bool           `json:"deleted"`
	Parent    *CommentParent `json:"parent,omitempty"`
	Inline    *CommentInline `json:"inline,omitempty"`
	Links     struct {
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
	} `json:"links"`
}

type commentList struct {
	Values []Comment `json:"values"`
	Next   string    `json:"next"`
}

func (c *Client) ListPRComments(repoSlug string, prID int) ([]Comment, error) {
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/comments?pagelen=100", c.workspace, repoSlug, prID)
	body, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var list commentList
	if err := json.Unmarshal(body, &list); err != nil {
		return nil, err
	}
	return list.Values, nil
}

// --- Activity ---

type ActivityActor struct {
	User Author `json:"user"`
	Date string `json:"date"`
}

type ActivityUpdate struct {
	State  string `json:"state"`
	Title  string `json:"title"`
	Date   string `json:"date"`
	Author Author `json:"author"`
}

type ActivityPRRef struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type Activity struct {
	PullRequest      ActivityPRRef  `json:"pull_request"`
	Comment          *Comment       `json:"comment,omitempty"`
	Approval         *ActivityActor `json:"approval,omitempty"`
	ChangesRequested *ActivityActor `json:"changes_requested,omitempty"`
	Update           *ActivityUpdate `json:"update,omitempty"`
}

type activityList struct {
	Values []Activity `json:"values"`
	Next   string     `json:"next"`
}

func (c *Client) ListPRActivity(repoSlug string, prID int) ([]Activity, error) {
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/activity?pagelen=50", c.workspace, repoSlug, prID)
	body, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var list activityList
	if err := json.Unmarshal(body, &list); err != nil {
		return nil, err
	}
	return list.Values, nil
}
