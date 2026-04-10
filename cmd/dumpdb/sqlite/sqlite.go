package sqlite

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v3"
)

func New() *cli.Command {
	return &cli.Command{
		Name:  "sqlite",
		Usage: "Dump a SQLite database",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:      "file",
				UsageText: "SQLite database file",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "dry-run",
				Usage:   "Dry run",
				Sources: cli.EnvVars("DRY_RUN"),
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			e := exec.CommandContext(ctx, "sqlite3",
				"-bail", cmd.Args().First(), ".dump",
			)
			e.Stdin = os.Stdin
			e.Stdout = os.Stdout
			e.Stderr = os.Stderr
			fmt.Println("+ ", strings.Join(e.Args, " "))
			if cmd.Bool("dry-run") {
				return nil
			}
			return e.Run()
		},
	}
}
