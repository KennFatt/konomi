package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

// Config holds all CLI and environment configuration.
type Config struct {
	URL         string
	Token       string
	State       string
	Format      string
	Output      string
	ReviewsOnly bool
	Args        []string

	ShowHelp bool
}

// Parse reads configuration from flags and environment variables.
func Parse() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.URL, "url", "", "Forgejo instance URL (env: FORGEJO_URL)")
	flag.StringVar(&cfg.Token, "token", "", "Forgejo API token (env: FORGEJO_TOKEN)")
	flag.StringVar(&cfg.State, "state", "open", "PR state filter: open, closed, all")
	flag.StringVar(&cfg.Format, "format", "markdown", "Output format: markdown, json")
	flag.StringVar(&cfg.Output, "output", "", "Write output to file instead of stdout")
	flag.BoolVar(&cfg.ReviewsOnly, "reviews-only", false, "Only show reviews and comments in output")

	help := flag.Bool("help", false, "Show this help message")
	flag.BoolVar(help, "h", false, "Show this help message")

	flag.Usage = printUsage
	flag.Parse()

	if *help {
		cfg.ShowHelp = true
		return cfg
	}

	// Environment overrides for unset flags.
	if cfg.URL == "" {
		cfg.URL = os.Getenv("FORGEJO_URL")
	}
	if cfg.Token == "" {
		cfg.Token = os.Getenv("FORGEJO_TOKEN")
	}

	cfg.Args = flag.Args()

	return cfg
}

// Validate checks that required configuration is present.
func (c *Config) Validate() error {
	var errs []error

	if c.URL == "" {
		errs = append(errs, errors.New("Forgejo URL is required: set FORGEJO_URL or --url"))
	}
	if c.Token == "" {
		errs = append(errs, errors.New("Forgejo API token is required: set FORGEJO_TOKEN or --token"))
	}
	switch c.State {
	case "open", "closed", "all":
	default:
		errs = append(errs, fmt.Errorf("invalid state %q: must be open, closed, or all", c.State))
	}
	switch c.Format {
	case "markdown", "json":
	default:
		errs = append(errs, fmt.Errorf("invalid format %q: must be markdown or json", c.Format))
	}

	return errors.Join(errs...)
}

func printUsage() {
	fmt.Fprint(os.Stderr, `konomi — Forgejo pull request information collector

Usage:
  konomi [flags] <owner/repo>               List pull requests
  konomi [flags] <owner/repo> <pr-number>   Show pull request details

Arguments:
  owner/repo     Repository in owner/name format (e.g., owner/repo)
  pr-number      Pull request number (e.g., 42)

Flags:
  --url          string   Forgejo instance URL (env: FORGEJO_URL)
  --token        string   Forgejo API token (env: FORGEJO_TOKEN)
  --state        string   PR state filter: open, closed, all (default: open)
  --format       string   Output format: markdown, json (default: markdown)
  --output       string   Write output to file instead of stdout
  --reviews-only          Only show reviews and comments in output
  --help, -h              Show this help message

Environment:
  FORGEJO_URL     Forgejo instance URL
  FORGEJO_TOKEN   Forgejo API token

Examples:
  konomi owner/repo
  konomi owner/repo --state closed
  konomi owner/repo 42
  konomi owner/repo 42 --format json --output pr42.json
`)
}
