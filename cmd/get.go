package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Display one or many PostgreSQL resources",
	Long: `Display one or many PostgreSQL resources including roles and users.
Use this command to list existing roles and users in the database.`,
}

func init() {
	rootCmd.AddCommand(getCmd)
}

// rolesCmd represents the roles subcommand
var rolesCmd = &cobra.Command{
	Use:   "roles",
	Short: "List roles",
	Long:  `List all roles in the PostgreSQL database.`,
	Run: func(cmd *cobra.Command, args []string) {
		db := GetDB()
		if db == nil {
			fmt.Fprintln(os.Stderr, "Error: Database connection not available")
			os.Exit(1)
		}
		
		roles, err := db.ListRoles()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing roles: %v\n", err)
			os.Exit(1)
		}

		// Print roles in a table format
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		if _, err := fmt.Fprintln(w, "NAME"); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to output: %v\n", err)
			os.Exit(1)
		}
		for _, role := range roles {
			if _, err := fmt.Fprintln(w, role); err != nil {
				fmt.Fprintf(os.Stderr, "Error writing to output: %v\n", err)
				os.Exit(1)
			}
		}
		if err := w.Flush(); err != nil {
			fmt.Fprintf(os.Stderr, "Error flushing output: %v\n", err)
			os.Exit(1)
		}
	},
}

// usersCmd represents the users subcommand
var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "List users",
	Long:  `List all users in the PostgreSQL database.`,
	Run: func(cmd *cobra.Command, args []string) {
		db := GetDB()
		if db == nil {
			fmt.Fprintln(os.Stderr, "Error: Database connection not available")
			os.Exit(1)
		}
		
		users, err := db.ListUsers()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing users: %v\n", err)
			os.Exit(1)
		}

		// Print users in a table format
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		if _, err := fmt.Fprintln(w, "NAME"); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to output: %v\n", err)
			os.Exit(1)
		}
		for _, user := range users {
			if _, err := fmt.Fprintln(w, user); err != nil {
				fmt.Fprintf(os.Stderr, "Error writing to output: %v\n", err)
				os.Exit(1)
			}
		}
		if err := w.Flush(); err != nil {
			fmt.Fprintf(os.Stderr, "Error flushing output: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	getCmd.AddCommand(rolesCmd)
	getCmd.AddCommand(usersCmd)
}