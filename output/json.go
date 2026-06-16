package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/KennFatt/konomi/forgejo"
)

// WritePullDetailJSON writes the full pull request detail as indented JSON.
// When reviewsOnly is true, only reviews and comments are included.
func WritePullDetailJSON(w io.Writer, detail *forgejo.PullDetail, reviewsOnly bool) error {
	if reviewsOnly {
		filtered := &struct {
			Reviews  []forgejo.ReviewWithComments `json:"reviews"`
			Comments []forgejo.Comment            `json:"comments"`
		}{
			Reviews:  detail.Reviews,
			Comments: detail.Comments,
		}
		return writeJSON(w, filtered)
	}
	return writeJSON(w, detail)
}

func writeJSON(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}
