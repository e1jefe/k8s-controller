package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	logLevel string
	verbose  bool
	quiet    bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pflag-demo",
	Short: "A simple Go application demonstrating pflag with log levels",
	Long: `pflag-demo is a CLI application that demonstrates the use of pflag
for handling command-line flags, specifically for setting log levels.

You can set log levels using --log-level flag or use convenient shortcuts
like --verbose for debug level and --quiet for error level.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Handle conflicting flags
		if verbose && quiet {
			fmt.Fprintf(os.Stderr, "Error: cannot use both --verbose and --quiet flags together\n")
			os.Exit(1)
		}

		// Set log level based on flags
		if verbose {
			logLevel = "debug"
		} else if quiet {
			logLevel = "error"
		}

		// Validate log level
		validLevels := []string{"debug", "info", "warn", "error"}
		logLevel = strings.ToLower(logLevel)

		valid := false
		for _, level := range validLevels {
			if logLevel == level {
				valid = true
				break
			}
		}

		if !valid {
			fmt.Fprintf(os.Stderr, "Error: invalid log level '%s'. Valid levels are: %s\n",
				logLevel, strings.Join(validLevels, ", "))
			os.Exit(1)
		}

		// Configure and demonstrate logging
		setupLogging(logLevel)

		// Demo the application
		runDemo()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Define flags
	rootCmd.Flags().StringVarP(&logLevel, "log-level", "l", "info", "Set the log level (debug, info, warn, error)")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging (equivalent to --log-level=debug)")
	rootCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Enable quiet mode (equivalent to --log-level=error)")
}

// setupLogging configures the log output based on level
func setupLogging(level string) {
	switch level {
	case "debug":
		log.SetPrefix("[DEBUG] ")
	case "info":
		log.SetPrefix("[INFO] ")
	case "warn":
		log.SetPrefix("[WARN] ")
	case "error":
		log.SetPrefix("[ERROR] ")
	}

	fmt.Printf("Log level set to: %s\n", strings.ToUpper(level))
}

// runDemo demonstrates different log levels
func runDemo() {
	fmt.Println("\n--- Application Demo ---")

	switch logLevel {
	case "debug":
		log.Println("Debug message: Application started with detailed logging")
		log.Println("Info message: Processing data...")
		log.Println("Warn message: This is a warning")
		log.Println("Error message: This is an error")
	case "info":
		log.Println("Info message: Processing data...")
		log.Println("Warn message: This is a warning")
		log.Println("Error message: This is an error")
	case "warn":
		log.Println("Warn message: This is a warning")
		log.Println("Error message: This is an error")
	case "error":
		log.Println("Error message: This is an error")
	}
}
