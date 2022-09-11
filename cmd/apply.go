package cmd

import (
	"fmt"
	"os"

	"github.com/nipeharefa/postgre-yaml-role/pkg/types"
	"github.com/nipeharefa/postgre-yaml-role/pkg/yaml"
	"github.com/spf13/cobra"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply a configuration to a PostgreSQL database",
	Long: `Apply a configuration to a PostgreSQL database from YAML files.
Use this command to create or update roles, users, and grants.`,
	Run: func(cmd *cobra.Command, args []string) {
		filename, _ := cmd.Flags().GetString("file")
		if filename == "" {
			fmt.Fprintln(os.Stderr, "Error: --file flag is required")
			os.Exit(1)
		}

		// Check if file is a directory
		var resources []interface{}
		fi, err := os.Stat(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accessing file: %v\n", err)
			os.Exit(1)
		}

		if fi.IsDir() {
			// Parse all YAML files in directory
			resources, err = yaml.ParseDirectory(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing directory: %v\n", err)
				os.Exit(1)
			}
		} else {
			// Parse single YAML file
			resources, err = yaml.ParseFile(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing file: %v\n", err)
				os.Exit(1)
			}
		}

		// Apply each resource
		db := GetDB()
		if db == nil {
			fmt.Fprintln(os.Stderr, "Error: Database connection not available")
			os.Exit(1)
		}
		
		for _, resource := range resources {
			switch r := resource.(type) {
			case types.Role:
				fmt.Printf("Applying role: %s\n", r.Metadata.Name)
				if err := db.CreateRole(r); err != nil {
					fmt.Fprintf(os.Stderr, "Error creating role %s: %v\n", r.Metadata.Name, err)
					os.Exit(1)
				}
			case types.User:
				fmt.Printf("Applying user: %s\n", r.Metadata.Name)
				if err := db.CreateUser(r); err != nil {
					fmt.Fprintf(os.Stderr, "Error creating user %s: %v\n", r.Metadata.Name, err)
					os.Exit(1)
				}
			case types.Grant:
				fmt.Printf("Applying grant: %s\n", r.Metadata.Name)
				if err := db.CreateGrant(r); err != nil {
					fmt.Fprintf(os.Stderr, "Error creating grant %s: %v\n", r.Metadata.Name, err)
					os.Exit(1)
				}
			}
		}

		fmt.Println("Configuration applied successfully")
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.Flags().StringP("file", "f", "", "YAML file or directory to apply")
}