package server

import (
	"context"
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"bitbucket_mcp/internal/bitbucket"
)

// Server wraps the MCP server and holds the Bitbucket client.
type Server struct {
	mcp    *mcp.Server
	client *bitbucket.Client
}

// New creates a new MCP server wired to the given Bitbucket client.
func New(client *bitbucket.Client) *Server {
	s := &Server{client: client}

	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "bitbucket-mcp",
		Version: "1.0.0",
	}, nil)

	s.mcp = mcpServer
	s.registerTools()
	return s
}

// Run starts the server over the stdio transport and blocks until done.
func (s *Server) Run(ctx context.Context) error {
	return s.mcp.Run(ctx, &mcp.StdioTransport{})
}

// --- tool input types ---

type listPRsInput struct {
	RepoSlug string `json:"repo_slug" jsonschema:"required,the repository slug"`
	State    string `json:"state,omitempty" jsonschema:"the PR state filter (OPEN, MERGED, DECLINED, SUPERSEDED); defaults to OPEN"`
}

type prInput struct {
	RepoSlug string  `json:"repo_slug" jsonschema:"required,the repository slug"`
	PRID     float64 `json:"pr_id" jsonschema:"required,the pull request ID"`
}

type createPRInput struct {
	RepoSlug          string `json:"repo_slug" jsonschema:"required,the repository slug"`
	Title             string `json:"title" jsonschema:"required,the PR title"`
	SourceBranch      string `json:"source_branch" jsonschema:"required,the source branch name"`
	DestinationBranch string `json:"destination_branch" jsonschema:"required,the destination branch name"`
	Description       string `json:"description,omitempty" jsonschema:"optional PR description in markdown"`
}

type updatePRInput struct {
	RepoSlug    string  `json:"repo_slug" jsonschema:"required,the repository slug"`
	PRID        float64 `json:"pr_id" jsonschema:"required,the pull request ID"`
	Title       string  `json:"title,omitempty" jsonschema:"optional new title for the PR"`
	Description string  `json:"description,omitempty" jsonschema:"optional new description for the PR (markdown)"`
}

type listBranchesInput struct {
	RepoSlug string `json:"repo_slug" jsonschema:"required,the repository slug"`
}

type listRepositoriesInput struct{}

type compareBranchesInput struct {
	RepoSlug          string `json:"repo_slug" jsonschema:"required,the repository slug"`
	SourceBranch      string `json:"source_branch" jsonschema:"required,the source branch name"`
	DestinationBranch string `json:"destination_branch" jsonschema:"required,the destination branch name"`
}

// registerTools adds all Bitbucket PR tools to the MCP server.
func (s *Server) registerTools() {
	// 1. list_pull_requests
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "list_pull_requests",
		Description: "List pull requests for a Bitbucket repository",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input listPRsInput) (*mcp.CallToolResult, any, error) {
		prs, err := s.client.ListPullRequests(input.RepoSlug, input.State)
		if err != nil {
			return nil, nil, err
		}
		b, err := json.Marshal(prs)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(b)}},
		}, nil, nil
	})

	// 2. get_pull_request
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_pull_request",
		Description: "Get a specific pull request by ID",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input prInput) (*mcp.CallToolResult, any, error) {
		pr, err := s.client.GetPullRequest(input.RepoSlug, int(input.PRID))
		if err != nil {
			return nil, nil, err
		}
		b, err := json.Marshal(pr)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(b)}},
		}, nil, nil
	})

	// 3. list_pr_commits
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "list_pr_commits",
		Description: "List commits for a pull request",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input prInput) (*mcp.CallToolResult, any, error) {
		commits, err := s.client.ListPRCommits(input.RepoSlug, int(input.PRID))
		if err != nil {
			return nil, nil, err
		}
		b, err := json.Marshal(commits)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(b)}},
		}, nil, nil
	})

	// 4. get_pr_diff
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_pr_diff",
		Description: "Get the unified diff for a pull request",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input prInput) (*mcp.CallToolResult, any, error) {
		diff, err := s.client.GetPRDiff(input.RepoSlug, int(input.PRID))
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: diff}},
		}, nil, nil
	})

	// 5. get_pr_diffstat
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_pr_diffstat",
		Description: "Get the diffstat (file-level change summary) for a pull request",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input prInput) (*mcp.CallToolResult, any, error) {
		entries, err := s.client.GetPRDiffstat(input.RepoSlug, int(input.PRID))
		if err != nil {
			return nil, nil, err
		}
		b, err := json.Marshal(entries)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(b)}},
		}, nil, nil
	})

	// 6. create_pull_request
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "create_pull_request",
		Description: "Create a new pull request in a Bitbucket repository",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input createPRInput) (*mcp.CallToolResult, any, error) {
		pr, err := s.client.CreatePullRequest(input.RepoSlug, bitbucket.CreatePRRequest{
			Title:             input.Title,
			Description:       input.Description,
			SourceBranch:      input.SourceBranch,
			DestinationBranch: input.DestinationBranch,
		})
		if err != nil {
			return nil, nil, err
		}
		b, err := json.Marshal(pr)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(b)}},
		}, nil, nil
	})

	// 7. list_branches
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "list_branches",
		Description: "List branches for a Bitbucket repository",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input listBranchesInput) (*mcp.CallToolResult, any, error) {
		branches, err := s.client.ListBranches(input.RepoSlug)
		if err != nil {
			return nil, nil, err
		}
		b, err := json.Marshal(branches)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(b)}},
		}, nil, nil
	})

	// 8. list_repositories
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "list_repositories",
		Description: "List all repositories in the configured Bitbucket workspace",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input listRepositoriesInput) (*mcp.CallToolResult, any, error) {
		repos, err := s.client.ListRepositories()
		if err != nil {
			return nil, nil, err
		}
		b, err := json.Marshal(repos)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(b)}},
		}, nil, nil
	})

	// 9. compare_branches_diff
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "compare_branches_diff",
		Description: "Get the raw unified diff between two branches in a repository",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input compareBranchesInput) (*mcp.CallToolResult, any, error) {
		diff, err := s.client.CompareBranchesDiff(input.RepoSlug, input.SourceBranch, input.DestinationBranch)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: diff}},
		}, nil, nil
	})

	// 10. compare_branches_diffstat
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "compare_branches_diffstat",
		Description: "Get the file-level diffstat between two branches in a repository",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input compareBranchesInput) (*mcp.CallToolResult, any, error) {
		entries, err := s.client.CompareBranchesDiffstat(input.RepoSlug, input.SourceBranch, input.DestinationBranch)
		if err != nil {
			return nil, nil, err
		}
		b, err := json.Marshal(entries)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(b)}},
		}, nil, nil
	})

	// 11. compare_branches_commits
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "compare_branches_commits",
		Description: "List commits reachable from the source branch but not the destination branch",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input compareBranchesInput) (*mcp.CallToolResult, any, error) {
		commits, err := s.client.CompareBranchesCommits(input.RepoSlug, input.SourceBranch, input.DestinationBranch)
		if err != nil {
			return nil, nil, err
		}
		b, err := json.Marshal(commits)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(b)}},
		}, nil, nil
	})

	// 13. list_pr_comments
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "list_pr_comments",
		Description: "List all comments on a pull request (both general and inline code comments)",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input prInput) (*mcp.CallToolResult, any, error) {
		comments, err := s.client.ListPRComments(input.RepoSlug, int(input.PRID))
		if err != nil {
			return nil, nil, err
		}
		b, err := json.Marshal(comments)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(b)}},
		}, nil, nil
	})

	// 14. list_pr_activity
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "list_pr_activity",
		Description: "List the full activity stream for a pull request (comments, approvals, changes-requested, source-branch updates)",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input prInput) (*mcp.CallToolResult, any, error) {
		activity, err := s.client.ListPRActivity(input.RepoSlug, int(input.PRID))
		if err != nil {
			return nil, nil, err
		}
		b, err := json.Marshal(activity)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(b)}},
		}, nil, nil
	})

	// 12. update_pull_request
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "update_pull_request",
		Description: "Update an existing pull request's title and/or description",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input updatePRInput) (*mcp.CallToolResult, any, error) {
		pr, err := s.client.UpdatePullRequest(input.RepoSlug, int(input.PRID), bitbucket.UpdatePRRequest{
			Title:       input.Title,
			Description: input.Description,
		})
		if err != nil {
			return nil, nil, err
		}
		b, err := json.Marshal(pr)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(b)}},
		}, nil, nil
	})
}
