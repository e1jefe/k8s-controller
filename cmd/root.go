package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "k8s-controller",
	Short: "A Kubernetes resource management tool",
	Long: `A command line tool for managing Kubernetes resources.
This tool provides commands to manage Kubernetes resources with support for
various resource types including pods, deployments, services, and more.

Examples:
  k8s-controller list deployments                    # List deployments in default namespace
  k8s-controller list deployments --kubeconfig ~/.kube/config  # Use specific kubeconfig`,
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
}
