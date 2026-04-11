package sqlite

import (
	"os"

	"github.com/gabe565/docker-restic/internal/cobrax"
	"github.com/gabe565/docker-restic/internal/dumpdb"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	var dryRun bool

	fs := &cobrax.Flags{}
	cmd := &cobra.Command{
		Use:   "sqlite file",
		Short: "Dump a SQLite database",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := fs.Resolve(); err != nil {
				return err
			}

			path := args[0]

			if _, err := os.Stat(path); err != nil {
				return err
			}

			return dumpdb.RunCmd(
				cmd, &dumpdb.RunOpts{DryRun: dryRun},
				"sqlite3", "-bail", path, ".dump",
			)
		},
	}

	fs.FlagSet = cmd.Flags()
	fs.Bool(&dryRun, dumpdb.FlagDryRun, "", false, "Dry run",
		cobrax.Env("DB_DRY_RUN"))

	return cmd
}
