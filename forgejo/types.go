package forgejo

import "time"

// User represents a Forgejo user.
type User struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	FullName  string `json:"full_name"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
}

// PRBranchInfo contains branch information for a pull request.
type PRBranchInfo struct {
	Label  string      `json:"label"`
	Ref    string      `json:"ref"`
	SHA    string      `json:"sha"`
	RepoID int64       `json:"repo_id"`
	Repo   *Repository `json:"repo"`
}

// Repository represents a Forgejo repository.
type Repository struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	FullName      string `json:"full_name"`
	Owner         *User  `json:"owner"`
	HTMLURL       string `json:"html_url"`
	CloneURL      string `json:"clone_url"`
	DefaultBranch string `json:"default_branch"`
	Description   string `json:"description"`
}

// PullRequest represents a Forgejo pull request.
type PullRequest struct {
	ID             int64         `json:"id"`
	Number         int64         `json:"number"`
	Title          string        `json:"title"`
	Body           string        `json:"body"`
	User           *User         `json:"user"`
	State          string        `json:"state"`
	Created        time.Time     `json:"created_at"`
	Updated        time.Time     `json:"updated_at"`
	Closed         *time.Time    `json:"closed_at"`
	Merged         bool          `json:"merged"`
	MergedAt       *time.Time    `json:"merged_at"`
	MergedBy       *User         `json:"merged_by"`
	Base           *PRBranchInfo `json:"base"`
	Head           *PRBranchInfo `json:"head"`
	Additions      int           `json:"additions"`
	Deletions      int           `json:"deletions"`
	ChangedFiles   int           `json:"changed_files"`
	Comments       int           `json:"comments"`
	ReviewComments int           `json:"review_comments"`
	HTMLURL        string        `json:"html_url"`
	URL            string        `json:"url"`
	DiffURL        string        `json:"diff_url"`
	PatchURL       string        `json:"patch_url"`
	Mergeable      bool          `json:"mergeable"`
	MergeCommitSHA string        `json:"merge_commit_sha"`
	Draft          bool          `json:"draft"`
	MergeBase      string        `json:"merge_base"`
}

// PullReview represents a pull request review.
type PullReview struct {
	ID            int64     `json:"id"`
	User          *User     `json:"user"`
	State         string    `json:"state"`
	Body          string    `json:"body"`
	CommitID      string    `json:"commit_id"`
	SubmittedAt   time.Time `json:"submitted_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	CommentsCount int64     `json:"comments_count"`
	Stale         bool      `json:"stale"`
	Dismissed     bool      `json:"dismissed"`
	HTMLURL       string    `json:"html_url"`
	Official      bool      `json:"official"`
}

// PullReviewComment represents a comment on a pull request review.
type PullReviewComment struct {
	ID        int64     `json:"id"`
	Body      string    `json:"body"`
	Path      string    `json:"path"`
	Position  int       `json:"position"`
	DiffHunk  string    `json:"diff_hunk"`
	CommitID  string    `json:"commit_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      *User     `json:"user"`
	ReviewID  int64     `json:"pull_request_review_id"`
}

// Comment represents a comment on an issue or pull request.
type Comment struct {
	ID             int64     `json:"id"`
	Body           string    `json:"body"`
	User           *User     `json:"user"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	PullRequestURL string    `json:"pull_request_url"`
	IssueURL       string    `json:"issue_url"`
	HTMLURL        string    `json:"html_url"`
}

// RepoCommit contains the metadata of a commit.
type RepoCommit struct {
	Message   string      `json:"message"`
	URL       string      `json:"url"`
	Author    *CommitUser `json:"author"`
	Committer *CommitUser `json:"committer"`
}

// CommitUser represents the author or committer of a commit.
type CommitUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Date  string `json:"date"`
}

// Commit represents a Git commit linked to a pull request.
type Commit struct {
	SHA       string      `json:"sha"`
	Commit    *RepoCommit `json:"commit"`
	Author    *User       `json:"author"`
	Committer *User       `json:"committer"`
	Created   time.Time   `json:"created"`
	HTMLURL   string      `json:"html_url"`
}

// ChangedFile represents a file changed in a pull request.
type ChangedFile struct {
	Filename         string `json:"filename"`
	PreviousFilename string `json:"previous_filename"`
	Status           string `json:"status"`
	Additions        int64  `json:"additions"`
	Deletions        int64  `json:"deletions"`
	Changes          int64  `json:"changes"`
	RawURL           string `json:"raw_url"`
	ContentsURL      string `json:"contents_url"`
	HTMLURL          string `json:"html_url"`
}

// ReviewWithComments pairs a review with its inline comments.
type ReviewWithComments struct {
	Review   PullReview          `json:"review"`
	Comments []PullReviewComment `json:"comments"`
}

// PullDetail holds all gathered data about a single pull request.
type PullDetail struct {
	PullRequest *PullRequest         `json:"pull_request"`
	Commits     []Commit             `json:"commits"`
	Files       []ChangedFile        `json:"files"`
	Reviews     []ReviewWithComments `json:"reviews"`
	Comments    []Comment            `json:"comments"`
}
