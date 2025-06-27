package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2"
)

var (
	informerKubeconfig string
	namespace          string
)

// informerCmd represents the informer command
var informerCmd = &cobra.Command{
	Use:   "informer",
	Short: "Watch Kubernetes deployment changes",
	Long: `Watch for changes to Kubernetes deployment resources and log events.

Examples:
  k8s-controller informer                           # Watch deployments in default namespace
  k8s-controller informer --namespace=kube-system  # Watch deployments in kube-system`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runInformer(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(informerCmd)

	informerCmd.Flags().StringVar(&informerKubeconfig, "kubeconfig", "", "path to kubeconfig file")
	informerCmd.Flags().StringVar(&namespace, "namespace", "default", "namespace to watch")
}

// runInformer starts the deployment informer
func runInformer() error {
	// Create Kubernetes client
	client, err := createClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Create informer
	informerFactory := informers.NewSharedInformerFactoryWithOptions(
		client,
		30*time.Second,
		informers.WithNamespace(namespace),
	)

	deploymentInformer := informerFactory.Apps().V1().Deployments().Informer()

	// Add simple event handlers
	deploymentInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			deployment := obj.(*appsv1.Deployment)
			klog.Infof("Deployment ADDED: %s/%s", deployment.Namespace, deployment.Name)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			deployment := newObj.(*appsv1.Deployment)
			klog.Infof("Deployment UPDATED: %s/%s", deployment.Namespace, deployment.Name)
		},
		DeleteFunc: func(obj interface{}) {
			deployment := obj.(*appsv1.Deployment)
			klog.Infof("Deployment DELETED: %s/%s", deployment.Namespace, deployment.Name)
		},
	})

	klog.Infof("Starting informer for deployments in namespace: %s", namespace)

	// Start informer
	ctx := context.Background()
	informerFactory.Start(ctx.Done())

	// Wait for cache sync
	if !cache.WaitForCacheSync(ctx.Done(), deploymentInformer.HasSynced) {
		return fmt.Errorf("failed to sync cache")
	}

	klog.Info("Informer running! Press Ctrl+C to stop...")

	// Keep running
	select {}
}

// createClient creates a simple Kubernetes client
func createClient() (kubernetes.Interface, error) {
	if informerKubeconfig == "" {
		if home := homedir.HomeDir(); home != "" {
			informerKubeconfig = filepath.Join(home, ".kube", "config")
		}
	}

	config, err := clientcmd.BuildConfigFromFlags("", informerKubeconfig)
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}
