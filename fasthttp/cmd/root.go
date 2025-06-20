package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fasthttp-server",
	Short: "A high-performance HTTP server built with fasthttp",
	Long: `A command line HTTP server application built with fasthttp library.
This tool provides commands to start and manage a high-performance HTTP server.

Examples:
  fasthttp-server server                    # Start server on default port 8080
  fasthttp-server server --port 3000       # Start server on port 3000
  fasthttp-server server -p 9000           # Start server on port 9000`,
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.fasthttp-server.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
