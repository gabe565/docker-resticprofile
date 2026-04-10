package mongodb

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gabe565/docker-restic/cmd/dumpdb/internal/clix"
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
				Sources:  cli.NewValueSourceChain(cli.EnvVar("DB_HOST"), clix.SecretFile(&mount, "host")),
			},
			&cli.StringFlag{
				Name:     "dbname",
				Usage:    "Database name",
				Aliases:  []string{"d"},
				Required: true,
				Sources:  cli.NewValueSourceChain(cli.EnvVar("DB_NAME"), clix.SecretFile(&mount, "dbname")),
			},
			&cli.StringFlag{
				Name:     "username",
				Usage:    "Database user",
				Aliases:  []string{"u"},
				Required: true,
				Sources:  cli.NewValueSourceChain(cli.EnvVar("DB_USERNAME"), clix.SecretFile(&mount, "username")),
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
			&cli.StringFlag{
				Name:    "dry-run",
				Usage:   "Dry run",
				Sources: cli.EnvVars("DRY_RUN"),
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			host := cmd.String("host")
			dbname := cmd.String("dbname")
			username := cmd.String("username")
			password := cmd.String("password")
			authDB := cmd.String("authentication-db")

			e := exec.CommandContext(ctx, "mongodump",
				"--authenticationDatabase="+authDB,
				"--host="+host, "--username="+username, "--password="+password, "--db="+dbname,
			)
			e.Stdin = os.Stdin
			e.Stdout = os.Stdout
			e.Stderr = os.Stderr
			fmt.Println("+ ", strings.Replace(strings.Join(e.Args, " "), password, "***", 1))
			if cmd.Bool("dry-run") {
				return nil
			}
			return e.Run()
		},
	}
}
