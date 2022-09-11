package cmd

import (
	"fmt"
	"os"

	"github.com/nipeharefa/postgre-yaml-role/pkg/types"
	"github.com/nipeharefa/postgre-yaml-role/pkg/yaml"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete roles and users from a PostgreSQL database",
	Long: `Delete roles and users from a PostgreSQL database based on YAML files.
Use this command to remove roles, users, and grants defined in YAML manifests.`,
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

		// Delete each resource
		db := GetDB()
		if db == nil {
			fmt.Fprintln(os.Stderr, "Error: Database connection not available")
			os.Exit(1)
		}
		
		for _, resource := range resources {
			switch r := resource.(type) {
			case types.Role:
				fmt.Printf("Deleting role: %s\n", r.Metadata.Name)
				if err := db.DeleteRole(r.Metadata.Name); err != nil {
					fmt.Fprintf(os.Stderr, "Error deleting role %s: %v\n", r.Metadata.Name, err)
					os.Exit(1)
				}
			case types.User:
				fmt.Printf("Deleting user: %s\n", r.Metadata.Name)
				if err := db.DeleteUser(r.Metadata.Name); err != nil {
					fmt.Fprintf(os.Stderr, "Error deleting user %s: %v\n", r.Metadata.Name, err)
					os.Exit(1)
				}
			case types.Grant:
				fmt.Printf("Deleting grant: %s\n", r.Metadata.Name)
				// Grants don't have a direct delete operation in our model
				// You would typically revoke privileges instead
				fmt.Printf("Note: Grants don't have a direct delete operation. You would need to revoke privileges manually.\n")
			}
		}

		fmt.Println("Resources deleted successfully")
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringP("file", "f", "", "YAML file or directory to delete")
}