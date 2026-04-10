package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/gabe565/docker-restic/cmd/dumpdb/cnpg"
	"github.com/gabe565/docker-restic/cmd/dumpdb/mariadb"
	"github.com/gabe565/docker-restic/cmd/dumpdb/mongodb"
	"github.com/gabe565/docker-restic/cmd/dumpdb/sqlite"
	"github.com/urfave/cli/v3"
)

func New() *cli.Command {
	return &cli.Command{
		Name:  "dumpdb",
		Usage: "Database utilities",
		Commands: []*cli.Command{
			cnpg.New(),
			mariadb.New(),
			mongodb.New(),
			sqlite.New(),
		},
	}
}

func main() {
	if err := New().Run(context.Background(), os.Args); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
