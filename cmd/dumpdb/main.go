package main

import (
	"log/slog"
	"os"

	"github.com/gabe565/docker-restic/cmd/dumpdb/cnpg"
	"github.com/gabe565/docker-restic/cmd/dumpdb/mariadb"
	"github.com/gabe565/docker-restic/cmd/dumpdb/mongodb"
	"github.com/gabe565/docker-restic/cmd/dumpdb/sqlite"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dumpdb",
		Short: "Database utilities",
	}

	cmd.AddCommand(
		cnpg.New(),
		mariadb.New(),
		mongodb.New(),
		sqlite.New(),
	)

	return cmd
}

func main() {
	if err := New().Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
