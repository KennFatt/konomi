package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/KennFatt/konomi/forgejo"
)

// WritePullDetailMarkdown renders the full pull request detail as Markdown.
// When reviewsOnly is true, only reviews and comments sections are included.
func WritePullDetailMarkdown(w io.Writer, detail *forgejo.PullDetail, reviewsOnly bool) error {
	pr := detail.PullRequest

	if !reviewsOnly {
		writeHeader(w, pr)
		writeMetadata(w, pr)
		writeDescription(w, pr)
		writeCommits(w, detail.Commits)
		writeFiles(w, detail.Files)
	}
	writeReviews(w, detail.Reviews)
	writeComments(w, detail.Comments)

	return nil
}

func writeHeader(w io.Writer, pr *forgejo.PullRequest) {
	fmt.Fprintf(w, "# PR #%d: %s\n\n", pr.Number, pr.Title)
	fmt.Fprintf(w, "**State:** %s", pr.State)
	if pr.Draft {
		fmt.Fprint(w, " (draft)")
	}
	if pr.Merged {
		fmt.Fprintf(w, " (merged by **%s**", pr.MergedBy.Login)
		if pr.MergedAt != nil {
			fmt.Fprintf(w, " on %s", pr.MergedAt.Format("2006-01-02 15:04"))
		}
		fmt.Fprint(w, ")")
	}
	fmt.Fprint(w, "\n\n")

	if repo := pr.Base.Repo; repo != nil {
		fmt.Fprintf(w, "**Repository:** [%s](%s)\n\n", repo.FullName, repo.HTMLURL)
	}
}

func writeMetadata(w io.Writer, pr *forgejo.PullRequest) {
	fmt.Fprint(w, "## Metadata\n\n")
	fmt.Fprintf(w, "| Field | Value |\n")
	fmt.Fprintf(w, "|-------|-------|\n")
	fmt.Fprintf(w, "| PR Number | #%d |\n", pr.Number)
	fmt.Fprintf(w, "| Title | %s |\n", pr.Title)
	fmt.Fprintf(w, "| Author | %s |\n", pr.User.Login)
	fmt.Fprintf(w, "| Created | %s |\n", pr.Created.Format("2006-01-02 15:04 MST"))
	if pr.Updated.After(pr.Created) {
		fmt.Fprintf(w, "| Updated | %s |\n", pr.Updated.Format("2006-01-02 15:04 MST"))
	}
	if pr.Closed != nil && !pr.Closed.IsZero() {
		fmt.Fprintf(w, "| Closed | %s |\n", pr.Closed.Format("2006-01-02 15:04 MST"))
	}
	fmt.Fprintf(w, "| Branch (from) | `%s` (%s) |\n", pr.Head.Ref, pr.Head.SHA[:8])
	fmt.Fprintf(w, "| Branch (to) | `%s` (%s) |\n", pr.Base.Ref, pr.Base.SHA[:8])
	fmt.Fprintf(w, "| Changed files | %d |\n", pr.ChangedFiles)
	fmt.Fprintf(w, "| Additions | %d |\n", pr.Additions)
	fmt.Fprintf(w, "| Deletions | %d |\n", pr.Deletions)
	fmt.Fprintf(w, "| Total comments | %d |\n", pr.Comments)
	fmt.Fprintf(w, "| Review comments | %d |\n", pr.ReviewComments)
	fmt.Fprintf(w, "| URL | %s |\n", pr.HTMLURL)
	if pr.MergeCommitSHA != "" {
		fmt.Fprintf(w, "| Merge commit | `%s` |\n", pr.MergeCommitSHA)
	}
	fmt.Fprint(w, "\n")
}

func writeDescription(w io.Writer, pr *forgejo.PullRequest) {
	if pr.Body == "" {
		return
	}
	fmt.Fprint(w, "## Description\n\n")
	fmt.Fprint(w, pr.Body)
	fmt.Fprint(w, "\n\n")
}

func writeCommits(w io.Writer, commits []forgejo.Commit) {
	fmt.Fprint(w, "## Commits\n\n")
	if len(commits) == 0 {
		fmt.Fprint(w, "_No commits._\n\n")
		return
	}
	fmt.Fprintf(w, "Total: %d commits\n\n", len(commits))
	for _, c := range commits {
		msg := firstLine(c.Commit.Message)
		fmt.Fprintf(w, "- [`%s`](%s) %s", c.SHA[:8], c.HTMLURL, msg)
		if c.Author != nil {
			fmt.Fprintf(w, " — %s", c.Author.Login)
		}
		fmt.Fprint(w, "\n")
	}
	fmt.Fprint(w, "\n")
}

func writeFiles(w io.Writer, files []forgejo.ChangedFile) {
	fmt.Fprint(w, "## Changed Files\n\n")
	if len(files) == 0 {
		fmt.Fprint(w, "_No files changed._\n\n")
		return
	}
	fmt.Fprintf(w, "Total: %d files\n\n", len(files))
	fmt.Fprint(w, "| File | Status | Additions | Deletions |\n")
	fmt.Fprint(w, "|------|--------|-----------|-----------|\n")
	for _, f := range files {
		status := f.Status
		if status == "" {
			status = "modified"
		}
		add := fmt.Sprintf("+%d", f.Additions)
		del := fmt.Sprintf("-%d", f.Deletions)
		fmt.Fprintf(w, "| `%s` | %s | %s | %s |\n", f.Filename, status, add, del)
	}
	fmt.Fprint(w, "\n")
}

func writeReviews(w io.Writer, reviews []forgejo.ReviewWithComments) {
	fmt.Fprint(w, "## Reviews\n\n")
	if len(reviews) == 0 {
		fmt.Fprint(w, "_No reviews._\n\n")
		return
	}
	for _, rwc := range reviews {
		rev := rwc.Review
		state := formatReviewState(rev.State)

		fmt.Fprintf(w, "### Review by %s • %s\n\n", rev.User.Login, state)
		if rev.SubmittedAt.Unix() > 0 {
			fmt.Fprintf(w, "**Submitted:** %s\n\n", rev.SubmittedAt.Format("2006-01-02 15:04"))
		}
		if rev.Body != "" {
			fmt.Fprint(w, rev.Body)
			fmt.Fprint(w, "\n\n")
		}

		if len(rwc.Comments) == 0 {
			continue
		}

		fmt.Fprint(w, "#### Comments\n\n")
		for _, c := range rwc.Comments {
			resolved := ""
			if c.Resolver != nil {
				resolved = fmt.Sprintf(" [resolved by %s]", c.Resolver.Login)
			}
			fmt.Fprintf(w, "- **%s** on `%s` (line %d)%s:\n", c.User.Login, c.Path, c.Position, resolved)
			if c.DiffHunk != "" {
				fmt.Fprint(w, "  ```diff\n")
				for _, line := range strings.Split(c.DiffHunk, "\n") {
					fmt.Fprintf(w, "  %s\n", line)
				}
				fmt.Fprint(w, "  ```\n")
			}
			fmt.Fprintf(w, "  > %s\n\n", c.Body)
		}
	}
}

func writeComments(w io.Writer, comments []forgejo.Comment) {
	fmt.Fprint(w, "## Comments\n\n")
	if len(comments) == 0 {
		fmt.Fprint(w, "_No comments._\n\n")
		return
	}
	for _, c := range comments {
		fmt.Fprintf(w, "### %s — %s\n\n", c.User.Login, c.CreatedAt.Format("2006-01-02 15:04"))
		fmt.Fprint(w, c.Body)
		fmt.Fprint(w, "\n\n")
	}
}

func firstLine(s string) string {
	s = strings.TrimSpace(s)
	if idx := strings.Index(s, "\n"); idx >= 0 {
		return s[:idx]
	}
	return s
}

func formatReviewState(s string) string {
	switch s {
	case "approved":
		return "✅ Approved"
	case "changes_requested":
		return "❌ Changes requested"
	case "pending":
		return "⏳ Pending"
	case "comment":
		return "💬 Comment"
	case "dismissed":
		return "🚫 Dismissed"
	default:
		return s
	}
}
