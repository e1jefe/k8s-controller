package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	logLevel string
	pretty   bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "zerolog-demo",
	Short: "A simple Go application demonstrating zerolog log levels",
	Long: `A demonstration application that shows how to use zerolog with different log levels:
- TRACE: Most detailed logging for debugging
- DEBUG: Debug information
- INFO: General operational messages
- WARN: Warning messages
- ERROR: Error conditions`,
	Run: func(cmd *cobra.Command, args []string) {
		setupLogger()
		demonstrateLogLevels()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&logLevel, "level", "l", "info", "Log level (trace, debug, info, warn, error)")
	rootCmd.Flags().BoolVarP(&pretty, "pretty", "p", false, "Enable pretty console output")
}

func setupLogger() {
	// Configure zerolog
	if pretty {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	}

	// Parse log level
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		log.Error().Err(err).Msg("Invalid log level, using INFO")
		level = zerolog.InfoLevel
	}

	// Set global log level
	zerolog.SetGlobalLevel(level)

	log.Info().
		Str("level", level.String()).
		Bool("pretty", pretty).
		Msg("Logger initialized")
}

func demonstrateLogLevels() {
	fmt.Println("\n=== Zerolog Log Levels Demonstration ===")
	fmt.Printf("Current log level: %s\n", logLevel)
	fmt.Printf("Pretty output: %t\n\n", pretty)

	// Trace level logging
	log.Trace().Msg("This is a TRACE message - most detailed debugging info")

	// Debug level logging
	log.Debug().Msg("This is a DEBUG message - debugging information")

	// Info level logging
	log.Info().Msg("This is an INFO message - general operational info")

	// Warn level logging
	log.Warn().Msg("This is a WARN message - warning about high memory usage")

	// Error level logging
	log.Error().Msg("This is an ERROR message - error condition occurred")

	fmt.Println("\n=== Log Level Filtering ===")
	fmt.Printf("Note: Only messages at or above '%s' level are displayed\n", logLevel)
	fmt.Println("Try running with different log levels:")
	fmt.Println("  --level trace   (shows all messages)")
	fmt.Println("  --level debug   (shows debug, info, warn, error)")
	fmt.Println("  --level info    (shows info, warn, error)")
	fmt.Println("  --level warn    (shows warn, error)")
	fmt.Println("  --level error   (shows only error)")
	fmt.Println("\nUse --pretty flag for colorized console output")
}
