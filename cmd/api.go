package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	apiKubeconfig string
	apiNamespace  string
	apiPort       string
	informer      cache.SharedIndexInformer
)

// Deployment represents a simple deployment response
type Deployment struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Replicas  int32  `json:"replicas"`
	Ready     int32  `json:"ready"`
}

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start JSON API server for deployments",
	Long:  "Start a simple JSON API server that lists deployments from informer cache",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runAPIServer(); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
	apiCmd.Flags().StringVar(&apiKubeconfig, "kubeconfig", "", "path to kubeconfig file")
	apiCmd.Flags().StringVar(&apiNamespace, "namespace", "default", "namespace to watch")
	apiCmd.Flags().StringVar(&apiPort, "port", "8080", "port to run the API server on")
}

func runAPIServer() error {
	// Setup informer
	if err := setupInformer(); err != nil {
		return err
	}

	// Setup HTTP handler
	http.HandleFunc("/deployments", listDeploymentsHandler)

	fmt.Printf("API server running on http://localhost:%s/deployments\n", apiPort)
	return http.ListenAndServe(":"+apiPort, nil)
}

func setupInformer() error {
	// Create client
	if apiKubeconfig == "" {
		if home := homedir.HomeDir(); home != "" {
			apiKubeconfig = filepath.Join(home, ".kube", "config")
		}
	}

	config, err := clientcmd.BuildConfigFromFlags("", apiKubeconfig)
	if err != nil {
		return err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	// Create informer
	factory := informers.NewSharedInformerFactoryWithOptions(
		client, 30*time.Second, informers.WithNamespace(apiNamespace))

	informer = factory.Apps().V1().Deployments().Informer()

	// Start informer
	factory.Start(context.Background().Done())
	cache.WaitForCacheSync(context.Background().Done(), informer.HasSynced)

	return nil
}

func listDeploymentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get deployments from cache
	var deployments []Deployment
	for _, obj := range informer.GetStore().List() {
		d := obj.(*appsv1.Deployment)
		deployments = append(deployments, Deployment{
			Name:      d.Name,
			Namespace: d.Namespace,
			Replicas:  *d.Spec.Replicas,
			Ready:     d.Status.ReadyReplicas,
		})
	}

	// Return JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deployments)
}
