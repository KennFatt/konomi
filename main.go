package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/KennFatt/forgejo-konomi/config"
	"github.com/KennFatt/forgejo-konomi/forgejo"
	"github.com/KennFatt/forgejo-konomi/output"
)

func main() {
	cfg := config.Parse()

	if cfg.ShowHelp {
		flag.Usage()
		return
	}

	// If no positional args, show help and exit.
	if len(cfg.Args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	// Validate config when we actually need to make API calls.
	if err := cfg.Validate(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		fmt.Fprintln(os.Stderr, "Run 'konomi --help' for usage.")
		os.Exit(1)
	}

	client := forgejo.New(cfg.URL, cfg.Token)

	switch len(cfg.Args) {
	case 1:
		runList(cfg, client)
	case 2:
		runDetail(cfg, client)
	default:
		fmt.Fprintln(os.Stderr, "error: too many arguments")
		fmt.Fprintln(os.Stderr, "Run 'konomi --help' for usage.")
		os.Exit(1)
	}
}

func openOutput(path string) io.WriteCloser {
	if path == "" {
		return nopCloser{os.Stdout}
	}
	f, err := os.Create(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: cannot create output file:", err)
		os.Exit(1)
	}
	return f
}

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

	detail, err := client.GetPullDetail(owner, repo, index)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	out := openOutput(cfg.Output)
	defer out.Close()

	var writeErr error
	if cfg.Format == "json" {
		writeErr = output.WritePullDetailJSON(out, detail)
	} else {
		writeErr = output.WritePullDetailMarkdown(out, detail)
	}
	if writeErr != nil {
		fmt.Fprintln(os.Stderr, "error:", writeErr)
		os.Exit(1)
	}
}

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }
