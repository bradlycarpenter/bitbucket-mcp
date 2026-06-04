# bitbucket-mcp

An MCP (Model Context Protocol) server that exposes Bitbucket Cloud repository and pull request data as tools consumable by any MCP-compatible client (Claude, Cursor, etc.).

## How it works

The server runs over stdio and authenticates to the Bitbucket Cloud REST API using HTTP Basic Auth (email + API token). Configuration is read from environment variables at startup. Each MCP tool maps directly to one or more Bitbucket API calls and returns the result as JSON text content.

## Tools

| Tool | Description |
|---|---|
| `list_repositories` | List all repositories in the configured workspace |
| `list_branches` | List branches for a repository |
| `list_pull_requests` | List pull requests, optionally filtered by state |
| `get_pull_request` | Get a specific pull request by ID |
| `create_pull_request` | Create a new pull request |
| `update_pull_request` | Update a pull request title and/or description |
| `list_pr_commits` | List commits on a pull request |
| `list_pr_comments` | List all comments on a pull request |
| `list_pr_activity` | Full activity stream for a pull request |
| `get_pr_diff` | Unified diff for a pull request |
| `get_pr_diffstat` | File-level change summary for a pull request |
| `compare_branches_diff` | Unified diff between two branches |
| `compare_branches_diffstat` | File-level diffstat between two branches |
| `compare_branches_commits` | Commits in source branch not in destination branch |

## Configuration

Copy `.env.example` to `.env` and populate the values:

```
BITBUCKET_EMAIL=user@yourcompany.com
BITBUCKET_TOKEN=<api-token>
BITBUCKET_DOMAIN=https://bitbucket.org/<workspace-slug>
```

The workspace slug is extracted from the path component of `BITBUCKET_DOMAIN`.

## Creating an API token

1. Go to https://id.atlassian.com/manage-profile/security/api-tokens
2. Click **Create API token**
3. Give it a descriptive label (e.g. `mcp-server`) and click **Create**
4. Copy the token immediately — it is not shown again

Use the token as `BITBUCKET_TOKEN` and your Atlassian account email as `BITBUCKET_EMAIL`.

### Required scopes

| Scope | Purpose |
|---|---|
| `read:repository:bitbucket` | List repositories and branches |
| `read:pullrequest:bitbucket` | Read pull requests, commits, diffs, comments, and activity |
| `write:pullrequest:bitbucket` | Create and update pull requests |

## Installation

Download the binary for your platform from the [latest release](https://github.com/bradlycarpenter/bitbucket-mcp/releases/latest), then follow the steps for your OS.

<details>
<summary>Linux</summary>

Download `bitbucket-mcp-linux-amd64`, mark it executable, and move it into your PATH:

```sh
chmod +x bitbucket-mcp-linux-amd64
mv bitbucket-mcp-linux-amd64 /usr/local/bin/bitbucket-mcp
```

</details>

<details>
<summary>macOS</summary>

Download the binary for your architecture:

- Apple Silicon: `bitbucket-mcp-darwin-arm64`
- Intel: `bitbucket-mcp-darwin-amd64`

Mark it executable and move it into your PATH:

```sh
chmod +x bitbucket-mcp-darwin-*
mv bitbucket-mcp-darwin-* /usr/local/bin/bitbucket-mcp
```

</details>

<details>
<summary>Windows</summary>

Download `bitbucket-mcp-windows-amd64.exe` and move it to a directory of your choice. Note the full path — you will need it when configuring your MCP client.

</details>

## MCP client configuration

Add an entry to your MCP client config pointing at the downloaded binary. Example for Claude Code:

```json
{
  "mcpServers": {
    "bitbucket": {
      "command": "/usr/local/bin/bitbucket-mcp",
      "env": {
        "BITBUCKET_EMAIL": "user@yourcompany.com",
        "BITBUCKET_TOKEN": "<api-token>",
        "BITBUCKET_DOMAIN": "https://bitbucket.org/<workspace-slug>"
      }
    }
  }
}
```

The server starts automatically when your MCP client launches and communicates over stdio — no separate process management is required.
