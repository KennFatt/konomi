package forgejo

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ListPulls returns pull requests for a repository, filtered by state.
// state is one of "open", "closed", "all".
func (c *Client) ListPulls(owner, repo, state string) ([]PullRequest, error) {
	body, err := c.getArray("/repos/"+owner+"/"+repo+"/pulls", map[string]string{
		"state": state,
	})
	if err != nil {
		return nil, fmt.Errorf("list pulls: %w", err)
	}

	var pulls []PullRequest
	if err := json.Unmarshal(body, &pulls); err != nil {
		return nil, fmt.Errorf("decode pulls: %w", err)
	}
	return pulls, nil
}

// GetPull returns a single pull request.
func (c *Client) GetPull(owner, repo string, index int64) (*PullRequest, error) {
	body, err := c.doGet(fmt.Sprintf("/repos/%s/%s/pulls/%d", owner, repo, index), nil)
	if err != nil {
		return nil, fmt.Errorf("get pull: %w", err)
	}

	var pr PullRequest
	if err := json.Unmarshal(body, &pr); err != nil {
		return nil, fmt.Errorf("decode pull: %w", err)
	}
	return &pr, nil
}

// GetPullCommits returns commits for a pull request.
func (c *Client) GetPullCommits(owner, repo string, index int64) ([]Commit, error) {
	body, err := c.getArray(
		fmt.Sprintf("/repos/%s/%s/pulls/%d/commits", owner, repo, index),
		map[string]string{"verification": "false", "files": "false"},
	)
	if err != nil {
		return nil, fmt.Errorf("get pull commits: %w", err)
	}

	var commits []Commit
	if err := json.Unmarshal(body, &commits); err != nil {
		return nil, fmt.Errorf("decode commits: %w", err)
	}
	return commits, nil
}

// GetPullFiles returns changed files for a pull request.
func (c *Client) GetPullFiles(owner, repo string, index int64) ([]ChangedFile, error) {
	body, err := c.getArray(
		fmt.Sprintf("/repos/%s/%s/pulls/%d/files", owner, repo, index), nil,
	)
	if err != nil {
		return nil, fmt.Errorf("get pull files: %w", err)
	}

	var files []ChangedFile
	if err := json.Unmarshal(body, &files); err != nil {
		return nil, fmt.Errorf("decode files: %w", err)
	}
	return files, nil
}

// RepoInfo returns basic repository metadata.
func (c *Client) RepoInfo(owner, repo string) (*Repository, error) {
	body, err := c.doGet("/repos/"+owner+"/"+repo, nil)
	if err != nil {
		return nil, fmt.Errorf("get repo: %w", err)
	}
	var r Repository
	if err := json.Unmarshal(body, &r); err != nil {
		return nil, fmt.Errorf("decode repo: %w", err)
	}
	return &r, nil
}

// ParseOwnerRepo splits "owner/repo" and validates.
func ParseOwnerRepo(s string) (owner, repo string, err error) {
	parts := strings.SplitN(s, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid repository format %q, expected owner/repo", s)
	}
	return parts[0], parts[1], nil
}
