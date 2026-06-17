package forgejo

import (
	"encoding/json"
	"fmt"
)

// ListReviews returns all reviews for a pull request.
func (c *Client) ListReviews(owner, repo string, index int64) ([]PullReview, error) {
	body, err := c.getArray(
		fmt.Sprintf("/repos/%s/%s/pulls/%d/reviews", owner, repo, index), nil,
	)
	if err != nil {
		return nil, fmt.Errorf("list reviews: %w", err)
	}

	var reviews []PullReview
	if err := json.Unmarshal(body, &reviews); err != nil {
		return nil, fmt.Errorf("decode reviews: %w", err)
	}
	return reviews, nil
}

// GetReviewComments returns all inline comments for a specific review.
func (c *Client) GetReviewComments(owner, repo string, index, reviewID int64) ([]PullReviewComment, error) {
	body, err := c.doGet(
		fmt.Sprintf("/repos/%s/%s/pulls/%d/reviews/%d/comments", owner, repo, index, reviewID), nil,
	)
	if err != nil {
		return nil, fmt.Errorf("get review comments: %w", err)
	}

	var comments []PullReviewComment
	if err := json.Unmarshal(body, &comments); err != nil {
		return nil, fmt.Errorf("decode review comments: %w", err)
	}
	return comments, nil
}

// ListIssueComments returns general (non-review) comments for a pull request.
// Forgejo treats PRs as a type of issue, so the issue comments endpoint works.
func (c *Client) ListIssueComments(owner, repo string, index int64) ([]Comment, error) {
	body, err := c.getArray(
		fmt.Sprintf("/repos/%s/%s/issues/%d/comments", owner, repo, index), nil,
	)
	if err != nil {
		return nil, fmt.Errorf("list issue comments: %w", err)
	}

	var comments []Comment
	if err := json.Unmarshal(body, &comments); err != nil {
		return nil, fmt.Errorf("decode issue comments: %w", err)
	}
	return comments, nil
}

// GetPullDetail fetches all information for a single pull request.
// The fetched data is aggregated from multiple endpoints.
func (c *Client) GetPullDetail(owner, repo string, index int64) (*PullDetail, error) {
	pr, err := c.GetPull(owner, repo, index)
	if err != nil {
		return nil, err
	}

	type partial struct {
		commits  []Commit
		files    []ChangedFile
		reviews  []PullReview
		comments []Comment
		err      error
	}

	ch := make(chan partial, 4)

	go func() {
		commits, err := c.GetPullCommits(owner, repo, index)
		ch <- partial{commits: commits, err: err}
	}()

	go func() {
		files, err := c.GetPullFiles(owner, repo, index)
		ch <- partial{files: files, err: err}
	}()

	go func() {
		reviews, err := c.ListReviews(owner, repo, index)
		ch <- partial{reviews: reviews, err: err}
	}()

	go func() {
		comments, err := c.ListIssueComments(owner, repo, index)
		ch <- partial{comments: comments, err: err}
	}()

	detail := &PullDetail{PullRequest: pr}

	for range 4 {
		r := <-ch
		if r.err != nil {
			return nil, r.err
		}
		if r.commits != nil {
			detail.Commits = r.commits
		}
		if r.files != nil {
			detail.Files = r.files
		}
		if r.reviews != nil {
			detail.Reviews = make([]ReviewWithComments, len(r.reviews))
			for i, rev := range r.reviews {
				detail.Reviews[i] = ReviewWithComments{Review: rev}
			}
		}
		if r.comments != nil {
			detail.Comments = r.comments
		}
	}

	for i, rwc := range detail.Reviews {
		comments, err := c.GetReviewComments(owner, repo, index, rwc.Review.ID)
		if err != nil {
			return nil, fmt.Errorf("get comments for review %d: %w", rwc.Review.ID, err)
		}
		detail.Reviews[i].Comments = comments
	}

	return detail, nil
}
