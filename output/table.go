package output

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/KennFatt/forgejo-konomi/forgejo"
)

// WritePullListTable formats a list of pull requests as a table.
func WritePullListTable(w io.Writer, pulls []forgejo.PullRequest) error {
	if len(pulls) == 0 {
		_, err := fmt.Fprintln(w, "No pull requests found.")
		return err
	}

	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)

	fmt.Fprintln(tw, "#\tState\tTitle\tAuthor\tFeedback\tCreated\tURL")
	fmt.Fprintln(tw, "-\t-----\t-----\t------\t--------\t-------\t---")

	for _, pr := range pulls {
		state := pr.State
		if pr.Draft {
			state = "draft"
		}
		created := relativeTime(pr.Created)
		title := truncate(pr.Title, 72)
		author := pr.User.Login

		fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%d\t%s\t%s\n", pr.Number, state, title, author, pr.Comments+pr.ReviewComments, created, pr.HTMLURL)
	}

	return tw.Flush()
}

// WritePullListJSON writes the PR list as a JSON array.
func WritePullListJSON(w io.Writer, pulls []forgejo.PullRequest) error {
	return writeJSON(w, pulls)
}

func relativeTime(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	case d < 30*24*time.Hour:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	default:
		return t.Format("2006-01-02")
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return strings.TrimSpace(s[:max-3]) + "..."
}
