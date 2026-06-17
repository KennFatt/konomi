package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/KennFatt/konomi/config"
	"github.com/KennFatt/konomi/forgejo"
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

	// Close flag requires a PR number (exactly one positional arg after repo).
	if cfg.Close && len(cfg.Args) != 2 {
		fmt.Fprintln(os.Stderr, "error: --close flag requires <owner/repo> <pr-number>")
		os.Exit(1)
	}

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

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }
