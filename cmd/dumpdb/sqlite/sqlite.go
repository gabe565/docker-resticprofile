package sqlite

import (
	"context"

	"github.com/gabe565/docker-restic/internal/dumpdb"
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
				Name:    dumpdb.FlagDryRun,
				Usage:   "Dry run",
				Sources: cli.EnvVars("DB_DRY_RUN"),
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return dumpdb.RunCmd(ctx, cmd, nil,
				"sqlite3", "-bail", cmd.StringArg(ArgFile), ".dump",
			)
		},
	}
}
