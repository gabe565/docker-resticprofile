package cnpg

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/gabe565/docker-restic/internal/cobrax"
	"github.com/gabe565/docker-restic/internal/dumpdb"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	var mount, host, database, username, password, restrictKey string
	var dryRun bool

	fs := &cobrax.Flags{}
	cmd := &cobra.Command{
		Use:   "cnpg",
		Short: "Dump a CloudNativePG database",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := fs.Resolve(); err != nil {
				return err
			}

			if restrictKey == "" {
				sum := sha256.Sum256([]byte(host + database + username + password + "\n"))
				restrictKey = hex.EncodeToString(sum[:])
			}

			args = append([]string{
				"--clean",
				"--if-exists",
				"--no-owner",
				"--restrict-key=" + restrictKey,
				"--host=" + host,
				"--username=" + username,
				"--dbname=" + database,
			}, args...)

			return dumpdb.RunCmd(cmd, "pg_dump", args, &dumpdb.RunOpts{
				Envs:   []string{"PGPASSWORD=" + password},
				DryRun: dryRun,
			})
		},
	}

	fs.FlagSet = cmd.Flags()
	fs.String(&mount, "secret-mount", "", "/postgresql-app", "Directory where secrets are mounted")
	fs.String(&host, dumpdb.FlagHost, "H", "", "Database host",
		cobrax.Env("DB_HOST"), cobrax.SecretFile(&mount, "host"))
	fs.String(&database, dumpdb.FlagDatabase, "d", "", "Database name",
		cobrax.Env("DB_DATABASE"), cobrax.SecretFile(&mount, "dbname"))
	fs.String(&username, dumpdb.FlagUsername, "u", "", "Database user",
		cobrax.Env("DB_USERNAME"), cobrax.SecretFile(&mount, "username"))
	fs.String(&password, dumpdb.FlagPassword, "p", "", "Database password",
		cobrax.Env("DB_PASSWORD"), cobrax.SecretFile(&mount, "password"))
	fs.String(&restrictKey, "restrict-key", "", "", "pg_dump restrict key",
		cobrax.Env("PG_RESTRICT_KEY"))
	fs.Bool(&dryRun, dumpdb.FlagDryRun, "", false, "Dry run",
		cobrax.Env("DB_DRY_RUN"))

	return cmd
}
