package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var kubeconfig string

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List Kubernetes resources",
	Long: `List Kubernetes resources in the specified namespace.
	
Examples:
  k8s-controller list deployments              # List deployments in default namespace
  k8s-controller list deployments --kubeconfig ~/.kube/config  # Use specific kubeconfig`,
}

// deploymentsCmd represents the deployments subcommand
var deploymentsCmd = &cobra.Command{
	Use:   "deployments",
	Short: "List all deployments in default namespace",
	Long: `List all deployments in the default namespace using the configured kubeconfig.
	
This command will connect to your Kubernetes cluster and display all deployments
in the default namespace with their basic information.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := listDeployments(); err != nil {
			fmt.Printf("Error listing deployments: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.AddCommand(deploymentsCmd)

	// Add kubeconfig flag to the list command
	listCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", "", "path to kubeconfig file (default is $HOME/.kube/config)")
}

// createKubernetesClient creates a Kubernetes client using the provided kubeconfig
func createKubernetesClient() (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error

	if kubeconfig == "" {
		// If no kubeconfig path provided, try default locations
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		}
	}

	// Check if kubeconfig file exists
	if _, err := os.Stat(kubeconfig); err == nil {
		// Use the current context in kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create config from kubeconfig: %w", err)
		}
	} else {
		// Fall back to in-cluster config (for when running inside a pod)
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to create in-cluster config and no valid kubeconfig found: %w", err)
		}
	}

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	return clientset, nil
}

// listDeployments lists all deployments in the default namespace
func listDeployments() error {
	// Create Kubernetes client
	clientset, err := createKubernetesClient()
	if err != nil {
		return err
	}

	// List deployments in default namespace
	deployments, err := clientset.AppsV1().Deployments("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list deployments: %w", err)
	}

	// Display results
	fmt.Printf("Found %d deployments in default namespace:\n\n", len(deployments.Items))

	if len(deployments.Items) == 0 {
		fmt.Println("No deployments found in default namespace.")
		return nil
	}

	// Print header
	fmt.Printf("%-30s %-10s %-10s %-10s %-15s\n", "NAME", "READY", "UP-TO-DATE", "AVAILABLE", "AGE")
	fmt.Println("----------------------------------------------------------------------------------------------")

	// Print deployment details
	for _, deployment := range deployments.Items {
		ready := fmt.Sprintf("%d/%d", deployment.Status.ReadyReplicas, deployment.Status.Replicas)
		upToDate := fmt.Sprintf("%d", deployment.Status.UpdatedReplicas)
		available := fmt.Sprintf("%d", deployment.Status.AvailableReplicas)

		// Calculate age
		age := metav1.Now().Sub(deployment.CreationTimestamp.Time)
		ageStr := ""
		if age.Hours() >= 24 {
			ageStr = fmt.Sprintf("%.0fd", age.Hours()/24)
		} else if age.Hours() >= 1 {
			ageStr = fmt.Sprintf("%.0fh", age.Hours())
		} else {
			ageStr = fmt.Sprintf("%.0fm", age.Minutes())
		}

		fmt.Printf("%-30s %-10s %-10s %-10s %-15s\n",
			deployment.Name,
			ready,
			upToDate,
			available,
			ageStr)
	}

	return nil
}
