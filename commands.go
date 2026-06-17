package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/KennFatt/konomi/config"
	"github.com/KennFatt/konomi/forgejo"
	"github.com/KennFatt/konomi/output"
)

func runList(cfg *config.Config, client *forgejo.Client) {
	owner, repo, err := forgejo.ParseOwnerRepo(cfg.Args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	pulls, err := client.ListPulls(owner, repo, cfg.State)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	out := openOutput(cfg.Output)
	defer out.Close()

	var writeErr error
	if cfg.Format == "json" {
		writeErr = output.WritePullListJSON(out, pulls)
	} else {
		writeErr = output.WritePullListTable(out, pulls)
	}
	if writeErr != nil {
		fmt.Fprintln(os.Stderr, "error:", writeErr)
		os.Exit(1)
	}
}

func runDetail(cfg *config.Config, client *forgejo.Client) {
	owner, repo, err := forgejo.ParseOwnerRepo(cfg.Args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	index, err := strconv.ParseInt(cfg.Args[1], 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: invalid PR number %q\n", cfg.Args[1])
		os.Exit(1)
	}

	if cfg.Close {
		closePullRequest(client, owner, repo, index, cfg.Reason)
	}

	detail, err := client.GetPullDetail(owner, repo, index)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	out := openOutput(cfg.Output)
	defer out.Close()

	var writeErr error
	if cfg.Format == "json" {
		writeErr = output.WritePullDetailJSON(out, detail, cfg.ReviewsOnly)
	} else {
		writeErr = output.WritePullDetailMarkdown(out, detail, cfg.ReviewsOnly)
	}
	if writeErr != nil {
		fmt.Fprintln(os.Stderr, "error:", writeErr)
		os.Exit(1)
	}
}

func closePullRequest(client *forgejo.Client, owner, repo string, index int64, reason string) {
	if reason != "" {
		_, err := client.CreateIssueComment(owner, repo, index, reason)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error creating comment:", err)
			os.Exit(1)
		}
	}

	_, err := client.EditPull(owner, repo, index, map[string]any{"state": "closed"})
	if err != nil {
		fmt.Fprintln(os.Stderr, "error closing pull request:", err)
		os.Exit(1)
	}
}
