package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// deploymentsCmd represents the deployments command
var deploymentsCmd = &cobra.Command{
	Use:   "deployments",
	Short: "List deployments in the cluster",
	Long: `List all deployments in the specified namespace or all namespaces.
This command shows deployment status including replicas, ready replicas,
and deployment age.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getKubernetesClient()
		if err != nil {
			fmt.Printf("Error creating Kubernetes client: %v\n", err)
			os.Exit(1)
		}

		var ns string
		if allNamespaces {
			ns = ""
		} else if namespace != "" {
			ns = namespace
		} else {
			ns = getCurrentNamespace()
		}

		deployments, err := client.AppsV1().Deployments(ns).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			fmt.Printf("Error listing deployments: %v\n", err)
			os.Exit(1)
		}

		if len(deployments.Items) == 0 {
			if allNamespaces {
				fmt.Println("No deployments found in any namespace")
			} else {
				fmt.Printf("No deployments found in namespace '%s'\n", ns)
			}
			return
		}

		// Create a tabwriter for formatted output
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		if allNamespaces {
			fmt.Fprintln(w, "NAMESPACE\tNAME\tREADY\tUP-TO-DATE\tAVAILABLE\tAGE")
		} else {
			fmt.Fprintln(w, "NAME\tREADY\tUP-TO-DATE\tAVAILABLE\tAGE")
		}

		for _, deployment := range deployments.Items {
			age := metav1.Now().Sub(deployment.CreationTimestamp.Time).Truncate(1)

			ready := deployment.Status.ReadyReplicas
			desired := *deployment.Spec.Replicas
			upToDate := deployment.Status.UpdatedReplicas
			available := deployment.Status.AvailableReplicas

			if allNamespaces {
				fmt.Fprintf(w, "%s\t%s\t%d/%d\t%d\t%d\t%s\n",
					deployment.Namespace, deployment.Name,
					ready, desired, upToDate, available, age)
			} else {
				fmt.Fprintf(w, "%s\t%d/%d\t%d\t%d\t%s\n",
					deployment.Name, ready, desired,
					upToDate, available, age)
			}
		}
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(deploymentsCmd)

	deploymentsCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace to list deployments from")
	deploymentsCmd.Flags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "List deployments from all namespaces")
}
