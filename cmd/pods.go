package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var namespace string
var allNamespaces bool

// podsCmd represents the pods command
var podsCmd = &cobra.Command{
	Use:   "pods",
	Short: "List pods in the cluster",
	Long: `List all pods in the specified namespace or all namespaces.
This command connects to your Kubernetes cluster and displays
information about running pods including their status, age, and restarts.`,
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

		pods, err := client.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			fmt.Printf("Error listing pods: %v\n", err)
			os.Exit(1)
		}

		if len(pods.Items) == 0 {
			if allNamespaces {
				fmt.Println("No pods found in any namespace")
			} else {
				fmt.Printf("No pods found in namespace '%s'\n", ns)
			}
			return
		}

		// Create a tabwriter for formatted output
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		if allNamespaces {
			fmt.Fprintln(w, "NAMESPACE\tNAME\tREADY\tSTATUS\tRESTARTS\tAGE")
		} else {
			fmt.Fprintln(w, "NAME\tREADY\tSTATUS\tRESTARTS\tAGE")
		}

		for _, pod := range pods.Items {
			readyCount := 0
			totalCount := len(pod.Status.ContainerStatuses)
			restarts := int32(0)

			for _, containerStatus := range pod.Status.ContainerStatuses {
				if containerStatus.Ready {
					readyCount++
				}
				restarts += containerStatus.RestartCount
			}

			age := metav1.Now().Sub(pod.CreationTimestamp.Time).Truncate(1)

			if allNamespaces {
				fmt.Fprintf(w, "%s\t%s\t%d/%d\t%s\t%d\t%s\n",
					pod.Namespace, pod.Name, readyCount, totalCount,
					pod.Status.Phase, restarts, age)
			} else {
				fmt.Fprintf(w, "%s\t%d/%d\t%s\t%d\t%s\n",
					pod.Name, readyCount, totalCount,
					pod.Status.Phase, restarts, age)
			}
		}
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(podsCmd)

	podsCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace to list pods from")
	podsCmd.Flags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "List pods from all namespaces")
}
