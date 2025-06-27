package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var controllerKubeconfig string

// controllerCmd represents the controller command
var controllerCmd = &cobra.Command{
	Use:   "controller",
	Short: "Run a controller-runtime based controller with event logging",
	Long:  `Run a controller that watches Kubernetes Deployments and logs each event.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runController(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(controllerCmd)
	controllerCmd.Flags().StringVar(&controllerKubeconfig, "kubeconfig", "", "path to kubeconfig file")
}

// DeploymentReconciler reconciles Deployment objects
type DeploymentReconciler struct {
	client.Client
}

// Reconcile logs each event received through the informer
func (r *DeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var deployment appsv1.Deployment
	if err := r.Get(ctx, req.NamespacedName, &deployment); err != nil {
		if client.IgnoreNotFound(err) == nil {
			logger.Info("Deployment deleted",
				"name", req.Name,
				"namespace", req.Namespace,
				"time", time.Now().Format(time.RFC3339))
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Log the event with essential information
	logger.Info("Deployment event",
		"name", deployment.Name,
		"namespace", deployment.Namespace,
		"replicas", *deployment.Spec.Replicas,
		"ready", deployment.Status.ReadyReplicas,
		"image", getImage(&deployment),
		"time", time.Now().Format(time.RFC3339))

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller
func (r *DeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		Complete(r)
}

// getImage returns the first container image
func getImage(deployment *appsv1.Deployment) string {
	if len(deployment.Spec.Template.Spec.Containers) > 0 {
		return deployment.Spec.Template.Spec.Containers[0].Image
	}
	return "unknown"
}

// runController starts the controller
func runController() error {
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	// Build kubeconfig path
	kubeconfig := controllerKubeconfig
	if kubeconfig == "" {
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		}
	}

	// Create config and manager
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to build config: %w", err)
	}

	scheme := runtime.NewScheme()
	appsv1.AddToScheme(scheme)

	mgr, err := ctrl.NewManager(config, ctrl.Options{Scheme: scheme})
	if err != nil {
		return fmt.Errorf("failed to create manager: %w", err)
	}

	// Setup controller
	if err := (&DeploymentReconciler{
		Client: mgr.GetClient(),
	}).SetupWithManager(mgr); err != nil {
		return fmt.Errorf("failed to setup controller: %w", err)
	}

	log.Log.Info("Starting controller - watching Deployment events...")
	return mgr.Start(ctrl.SetupSignalHandler())
}
