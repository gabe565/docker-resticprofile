package mongodb

import (
	"context"
	"os/exec"

	"github.com/gabe565/docker-restic/internal/clix"
	"github.com/urfave/cli/v3"
)

func New() *cli.Command {
	var mount string
	return &cli.Command{
		Name:  "mongodb",
		Usage: "Dump a MongDB database",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "secret-mount",
				Usage:       "Directory where secrets are mounted",
				Value:       "/mongodb",
				Destination: &mount,
			},
			&cli.StringFlag{
				Name:     "host",
				Usage:    "Database host",
				Aliases:  []string{"h"},
				Required: true,
				Sources:  cli.NewValueSourceChain(cli.EnvVar("DB_HOST")),
			},
			&cli.StringFlag{
				Name:     "database",
				Usage:    "Database name",
				Aliases:  []string{"d"},
				Required: true,
				Sources:  cli.NewValueSourceChain(cli.EnvVar("DB_DATABASE")),
			},
			&cli.StringFlag{
				Name:     "username",
				Usage:    "Database user",
				Aliases:  []string{"u"},
				Required: true,
				Sources:  cli.NewValueSourceChain(cli.EnvVar("DB_USERNAME")),
			},
			&cli.StringFlag{
				Name:     "password",
				Usage:    "Database password",
				Aliases:  []string{"p"},
				Required: true,
				Sources:  cli.NewValueSourceChain(cli.EnvVar("DB_PASSWORD"), clix.SecretFile(&mount, "mongodb-passwords")),
			},
			&cli.StringFlag{
				Name:    "authentication-db",
				Usage:   "Authentication database",
				Sources: cli.NewValueSourceChain(cli.EnvVar("AUTHENTICATION_DB")),
			},
			&cli.BoolFlag{
				Name:    "dry-run",
				Usage:   "Dry run",
				Sources: cli.EnvVars("DB_DRY_RUN"),
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			host := cmd.String("host")
			database := cmd.String("database")
			username := cmd.String("username")
			password := cmd.String("password")
			authDB := cmd.String("authentication-db")

			e := exec.CommandContext(ctx, "mongodump",
				"--archive", "--authenticationDatabase="+authDB,
				"--host="+host, "--username="+username, "--db="+database,
			)
			e.Stdin = cmd.Reader
			e.Stdout = cmd.Writer
			e.Stderr = cmd.ErrWriter
			clix.XTrace(cmd.ErrWriter, append(e.Args, "--password=***"))
			e.Args = append(e.Args, "--password="+password)

			if cmd.Bool("dry-run") {
				return nil
			}
			return e.Run()
		},
	}
}
