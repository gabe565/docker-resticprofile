package sqlite

import (
	"context"
	"os/exec"

	"github.com/gabe565/docker-restic/internal/clix"
	"github.com/urfave/cli/v3"
)

const ArgFile = "file"

func New() *cli.Command {
	return &cli.Command{
		Name:  "sqlite",
		Usage: "Dump a SQLite database",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:      ArgFile,
				UsageText: "SQLite database file",
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "dry-run",
				Usage:   "Dry run",
				Sources: cli.EnvVars("DB_DRY_RUN"),
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			e := exec.CommandContext(ctx, "sqlite3",
				"-bail", cmd.StringArg(ArgFile), ".dump",
			)
			e.Stdin = cmd.Reader
			e.Stdout = cmd.Writer
			e.Stderr = cmd.ErrWriter
			clix.XTrace(cmd.ErrWriter, e.Args)

			if cmd.Bool("dry-run") {
				return nil
			}
			return e.Run()
		},
	}
}
