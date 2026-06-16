# konomi

**konomi** is a single-binary CLI tool that gathers pull request information from a [Forgejo](https://forgejo.org/) repository and outputs it as structured JSON or Markdown. Instead of clicking through the web UI to review PRs, konomi brings all the details - commits, file diffs, reviews, and comments - straight to your terminal or a file.

> The name is inspired by **Konomi Okonogi** from [Kaiju No. 8](https://kaiju-no8.net/).

## Installation

```bash
go install github.com/KennFatt/konomi@latest
```

Or build locally (see [Building from Source](#building-from-source)).

## Quick Start

```bash
# Set your Forgejo instance and token
export FORGEJO_URL=https://forgejo.example.com
export FORGEJO_TOKEN=your_api_token_here

# List open pull requests
konomi owner/repo

# List closed pull requests
konomi owner/repo --state closed

# Show PR #42 details (Markdown)
konomi owner/repo 42

# Show PR #42 details as JSON
konomi --format json owner/repo 42

# Write output to a file
konomi owner/repo 42 --output pr42.md
konomi --format json owner/repo 42 --output pr42.json

# Show only reviews and comments
konomi owner/repo 42 --reviews-only

# Reviews-only as JSON
konomi owner/repo 42 --reviews-only --format json
```

## Usage

```
konomi [flags] <owner/repo>               List pull requests
konomi [flags] <owner/repo> <pr-number>   Show pull request details
```

### Flags

| Flag | Env Variable | Default | Description |
|------|-------------|---------|-------------|
| `--url` | `FORGEJO_URL` | - | Forgejo instance URL *(required)* |
| `--token` | `FORGEJO_TOKEN` | - | Forgejo API token *(required)* |
| `--state` | - | `open` | PR state filter: `open`, `closed`, `all` |
| `--format` | - | `markdown` | Output format: `markdown`, `json` |
| `--output` | - | - | Write output to file instead of stdout |
| `--reviews-only` | - | `false` | Only show reviews and comments in output |
| `--help`, `-h` | - | - | Show help message |

> Flags must be placed **before** positional arguments due to Go's `flag` package behavior.

## What Information Is Collected

### Pull Request Detail (Markdown)

- **Metadata**: PR number, title, author, timestamps, branch from/to, change stats, LoC, URL
- **Description**: The PR body
- **Commits**: Each commit with short hash, URL, message, and author
- **Changed files**: Per-file status (added/modified/deleted) with additions/deletions
- **Reviews**: Review state (approved/changes requested/comment/pending/dismissed), body, and inline comments with diff context
  - Inline comments show `[resolved by <user>]` when the comment has been resolved via the Forgejo UI
- **General comments**: Non-review discussion on the PR

### JSON

The same data structured as JSON:
- `pull_requests` - full PR object from the API
- `commits` - list of commits
- `files` - list of changed files with stats
- `reviews` - list of reviews, each with its inline `Comments`
- `comments` - general (non-review) comments on the PR

When `--reviews-only` is used, JSON output contains only:
- `reviews` - reviews with inline comments (each comment includes a `resolver` field `null` if unresolved)
- `comments` - general comments

## Architecture

```
konomi/
├── main.go                 Entry point & argument dispatch
├── config/                 CLI flags + environment variable loading
├── forgejo/                Forgejo API client
│   ├── types.go            Response structs matching the API
│   ├── client.go           HTTP client with auth & pagination
│   ├── pulls.go            PR, commit, and file endpoints
│   └── reviews.go          Review, review comment, and issue comment endpoints
└── output/                 Output formatters
    ├── table.go            Tabular output for PR listing
    ├── markdown.go         Markdown output for PR detail
    └── json.go             JSON output for both modes
```

### Dependencies

**Zero external dependencies.** Only Go standard library:
- `flag` - CLI argument parsing
- `net/http` - HTTP client
- `encoding/json` - JSON serialization
- `text/tabwriter` - Table formatting

## Building from Source

```bash
git clone https://github.com/KennFatt/konomi
cd konomi
make build
```

Or directly:

```bash
go build -ldflags="-s -w" -o bin/konomi .
```

Requires Go 1.21+.
