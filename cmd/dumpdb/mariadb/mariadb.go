package mariadb

import (
	"context"

	"github.com/gabe565/docker-restic/internal/clix"
	"github.com/gabe565/docker-restic/internal/dumpdb"
	"github.com/urfave/cli/v3"
)

func New() *cli.Command {
	var mount string
	return &cli.Command{
		Name:  "mariadb",
		Usage: "Dump a MariaDB database",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "secret-mount",
				Usage:       "Directory where secrets are mounted",
				Value:       "/mariadb",
				Destination: &mount,
			},
			&cli.StringFlag{
				Name:     dumpdb.FlagHost,
				Usage:    "Database host",
				Aliases:  []string{"h"},
				Required: true,
				Sources:  cli.NewValueSourceChain(cli.EnvVar("DB_HOST")),
			},
			&cli.StringFlag{
				Name:     dumpdb.FlagDatabase,
				Usage:    "Database name",
				Aliases:  []string{"d"},
				Required: true,
				Sources:  cli.NewValueSourceChain(cli.EnvVar("DB_DATABASE")),
			},
			&cli.StringFlag{
				Name:     dumpdb.FlagUsername,
				Usage:    "Database user",
				Aliases:  []string{"u"},
				Required: true,
				Sources:  cli.NewValueSourceChain(cli.EnvVar("DB_USERNAME")),
			},
			&cli.StringFlag{
				Name:     dumpdb.FlagPassword,
				Usage:    "Database password",
				Aliases:  []string{"p"},
				Required: true,
				Sources:  cli.NewValueSourceChain(cli.EnvVar("DB_PASSWORD"), clix.SecretFile(&mount, "mariadb-password")),
			},
			&cli.BoolFlag{
				Name:    dumpdb.FlagDryRun,
				Usage:   "Dry run",
				Sources: cli.EnvVars("DB_DRY_RUN"),
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			host := cmd.String(dumpdb.FlagHost)
			database := cmd.String(dumpdb.FlagDatabase)
			username := cmd.String(dumpdb.FlagUsername)
			password := cmd.String(dumpdb.FlagPassword)

			return dumpdb.RunCmd(ctx, cmd, &dumpdb.RunOpts{Envs: []string{"MYSQL_PWD=" + password}},
				"mariadb-dump",
				"--add-drop-table",
				"--skip-dump-date",
				"--single-transaction",
				"--host="+host,
				"--user="+username,
				database,
			)
		},
	}
}
