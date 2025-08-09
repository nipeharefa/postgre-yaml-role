package cmd

import (
	"fmt"
	"os"

	"github.com/nipeharefa/postgre-yaml-role/pkg/config"
	"github.com/nipeharefa/postgre-yaml-role/pkg/postgres"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	db      *postgres.DB
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pgctl",
	Short: "A Kubernetes-inspired CLI for managing PostgreSQL roles and users",
	Long: `pgctl is a CLI tool that allows you to manage PostgreSQL roles and users
using YAML manifests, similar to how kubectl works with Kubernetes resources.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pgctl/config.yaml)")
	if err := viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config")); err != nil {
		fmt.Fprintf(os.Stderr, "Error binding config flag: %v\n", err)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// If config file is specified via flag, set it in viper
	if cfgFile != "" {
		viper.Set("config", cfgFile)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Connect to database
	db, err = postgres.NewDB(cfg.ConnectionString())
	if err != nil {
		// We won't exit here since some commands don't require a database connection
		// fmt.Fprintf(os.Stderr, "Warning: Error connecting to database: %v\n", err)
		_ = err // Explicitly ignore the error
	}
}

// GetDB returns the database connection
func GetDB() *postgres.DB {
	return db
}