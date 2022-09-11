package cmd

import (
	"fmt"
	"os"

	"github.com/nipeharefa/postgre-yaml-role/pkg/types"
	"github.com/nipeharefa/postgre-yaml-role/pkg/yaml"
	"github.com/spf13/cobra"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate YAML files without applying them",
	Long: `Validate YAML files without applying them to a PostgreSQL database.
Use this command to check if your YAML manifests are valid before applying them.`,
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

		// Print validation results
		fmt.Printf("Validated %d resources:\n", len(resources))
		for _, resource := range resources {
			switch r := resource.(type) {
			case types.Role:
				fmt.Printf("- Role: %s\n", r.Metadata.Name)
			case types.User:
				fmt.Printf("- User: %s\n", r.Metadata.Name)
			case types.Grant:
				fmt.Printf("- Grant: %s\n", r.Metadata.Name)
			default:
				fmt.Printf("- Unknown resource type\n")
			}
		}

		fmt.Println("All resources are valid")
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.Flags().StringP("file", "f", "", "YAML file or directory to validate")
}