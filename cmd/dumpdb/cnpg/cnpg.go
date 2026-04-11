package cnpg

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/gabe565/docker-restic/internal/clix"
	"github.com/gabe565/docker-restic/internal/dumpdb"
	"github.com/urfave/cli/v3"
)

func New() *cli.Command {
	var mount string
	return &cli.Command{
		Name:  "cnpg",
		Usage: "Dump a CloudNativePG database",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "secret-mount",
				Usage:       "Directory where secrets are mounted",
				Value:       "/postgresql-app",
				Destination: &mount,
			},
			&cli.StringFlag{
				Name:     dumpdb.FlagHost,
				Usage:    "Database host",
				Aliases:  []string{"h"},
				Required: true,
				Sources:  cli.NewValueSourceChain(cli.EnvVar("DB_HOST"), clix.SecretFile(&mount, "host")),
			},
			&cli.StringFlag{
				Name:     dumpdb.FlagDatabase,
				Usage:    "Database name",
				Aliases:  []string{"d"},
				Required: true,
				Sources:  cli.NewValueSourceChain(cli.EnvVar("DB_DATABASE"), clix.SecretFile(&mount, "dbname")),
			},
			&cli.StringFlag{
				Name:     dumpdb.FlagUsername,
				Usage:    "Database user",
				Aliases:  []string{"u"},
				Required: true,
				Sources:  cli.NewValueSourceChain(cli.EnvVar("DB_USERNAME"), clix.SecretFile(&mount, "username")),
			},
			&cli.StringFlag{
				Name:     dumpdb.FlagPassword,
				Usage:    "Database password",
				Aliases:  []string{"p"},
				Required: true,
				Sources:  cli.NewValueSourceChain(cli.EnvVar("DB_PASSWORD"), clix.SecretFile(&mount, "password")),
			},
			&cli.StringFlag{
				Name:    "restrict-key",
				Usage:   "pg_dump restrict key",
				Sources: cli.NewValueSourceChain(cli.EnvVar("PG_RESTRICT_KEY")),
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

			restrictKey := cmd.String("restrict-key")
			if restrictKey == "" {
				sum := sha256.Sum256([]byte(host + database + username + password + "\n"))
				restrictKey = hex.EncodeToString(sum[:])
			}

			return dumpdb.RunCmd(ctx, cmd, &dumpdb.RunOpts{Envs: []string{"PGPASSWORD=" + password}},
				"pg_dump",
				"--clean",
				"--if-exists",
				"--no-owner",
				"--restrict-key="+restrictKey,
				"--host="+host,
				"--username="+username,
				"--dbname="+database,
			)
		},
	}
}
