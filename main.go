package main

import (
	"bytes"
	"context"
	"database/sql"
	"io"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/nipeharefa/postgre-yaml-role/action"
	"github.com/nipeharefa/postgre-yaml-role/lib"
	"github.com/spf13/cobra"
)

type (
	Root struct {
		Kind string `yaml:"kind"`
	}
)

type (
	Action string
)

func main() {
	// var buf bytes.Buffer
	connStr := "user=postgres dbname=postgres sslmode=disable password=password"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	applyCmd := &cobra.Command{
		Use: "apply",
		Run: func(cmd *cobra.Command, args []string) {

			var buf bytes.Buffer
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				io.Copy(&buf, os.Stdin)
			} else {
				term, _ := cmd.Flags().GetString("file")
				f, _ := os.Open(term)
				io.Copy(&buf, f)
				f.Close()
			}

			lib.NewUserKind(db).Parser(context.Background(), &buf, action.ApplyAction)
		},
	}

	applyCmd.Flags().StringP("file", "f", "", "apply changes from yaml")

	rootCmd := &cobra.Command{}
	rootCmd.AddCommand(applyCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

}
