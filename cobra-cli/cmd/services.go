package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// servicesCmd represents the services command
var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "List services in the cluster",
	Long: `List all services in the specified namespace or all namespaces.
This command shows service information including type, cluster IP,
external IP, ports, and age.`,
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

		services, err := client.CoreV1().Services(ns).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			fmt.Printf("Error listing services: %v\n", err)
			os.Exit(1)
		}

		if len(services.Items) == 0 {
			if allNamespaces {
				fmt.Println("No services found in any namespace")
			} else {
				fmt.Printf("No services found in namespace '%s'\n", ns)
			}
			return
		}

		// Create a tabwriter for formatted output
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		if allNamespaces {
			fmt.Fprintln(w, "NAMESPACE\tNAME\tTYPE\tCLUSTER-IP\tEXTERNAL-IP\tPORT(S)\tAGE")
		} else {
			fmt.Fprintln(w, "NAME\tTYPE\tCLUSTER-IP\tEXTERNAL-IP\tPORT(S)\tAGE")
		}

		for _, service := range services.Items {
			age := metav1.Now().Sub(service.CreationTimestamp.Time).Truncate(1)

			// Format external IPs
			externalIPs := "<none>"
			if len(service.Status.LoadBalancer.Ingress) > 0 {
				var ips []string
				for _, ingress := range service.Status.LoadBalancer.Ingress {
					if ingress.IP != "" {
						ips = append(ips, ingress.IP)
					} else if ingress.Hostname != "" {
						ips = append(ips, ingress.Hostname)
					}
				}
				if len(ips) > 0 {
					externalIPs = strings.Join(ips, ",")
				}
			} else if len(service.Spec.ExternalIPs) > 0 {
				externalIPs = strings.Join(service.Spec.ExternalIPs, ",")
			}

			// Format ports
			var ports []string
			for _, port := range service.Spec.Ports {
				if port.NodePort != 0 {
					ports = append(ports, fmt.Sprintf("%d:%d/%s", port.Port, port.NodePort, port.Protocol))
				} else {
					ports = append(ports, fmt.Sprintf("%d/%s", port.Port, port.Protocol))
				}
			}
			portString := strings.Join(ports, ",")

			if allNamespaces {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
					service.Namespace, service.Name, service.Spec.Type,
					service.Spec.ClusterIP, externalIPs, portString, age)
			} else {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
					service.Name, service.Spec.Type, service.Spec.ClusterIP,
					externalIPs, portString, age)
			}
		}
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(servicesCmd)

	servicesCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace to list services from")
	servicesCmd.Flags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "List services from all namespaces")
}
