package dumpdb

import (
	"context"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/gabe565/docker-restic/internal/clix"
	"github.com/urfave/cli/v3"
)

const (
	FlagHost     = "host"
	FlagDatabase = "database"
	FlagUsername = "username"
	FlagPassword = "password"
	FlagDryRun   = "dry-run"
)

type RunOpts struct {
	Envs   []string
	Redact []string
}

func RunCmd(ctx context.Context, cmd *cli.Command, opts *RunOpts, name string, args ...string) error {
	if opts == nil {
		opts = &RunOpts{}
	}

	e := exec.CommandContext(ctx, name, args...)
	e.Env = append(os.Environ(), opts.Envs...)
	e.Stdin = cmd.Reader
	e.Stdout = cmd.Writer
	e.Stderr = cmd.ErrWriter

	xtrace := clix.XTraceString(e.Args)
	for _, v := range opts.Redact {
		xtrace = strings.ReplaceAll(xtrace, v, "***")
	}
	_, _ = io.WriteString(cmd.ErrWriter, xtrace)

	if cmd.Bool(FlagDryRun) {
		return nil
	}

	return e.Run()
}
