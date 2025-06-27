package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	appsv1 "k8s.io/api/apps/v1"
	coordinationv1 "k8s.io/api/coordination/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"
)

var (
	managerKubeconfig     string
	managerNamespace      string
	disableLeaderElection bool
	leaderElectionID      string
	metricsAddr           string
)

// managerCmd represents the manager command
var managerCmd = &cobra.Command{
	Use:   "manager",
	Short: "Run the controller manager with leader election and metrics",
	Long: `Run a comprehensive controller manager that controls both informer and controller.

Examples:
  k8s-controller manager                                # Run with default settings
  k8s-controller manager --disable-leader-election     # Run without leader election
  k8s-controller manager --metrics-addr=:9090          # Custom metrics port`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runManager(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(managerCmd)

	managerCmd.Flags().StringVar(&managerKubeconfig, "kubeconfig", "", "path to kubeconfig file")
	managerCmd.Flags().StringVar(&managerNamespace, "namespace", "default", "namespace to watch")
	managerCmd.Flags().BoolVar(&disableLeaderElection, "disable-leader-election", false, "disable leader election")
	managerCmd.Flags().StringVar(&leaderElectionID, "leader-election-id", "k8s-controller-manager", "leader election lease name")
	managerCmd.Flags().StringVar(&metricsAddr, "metrics-addr", ":8080", "address for metrics server")
}

// ManagerReconciler handles deployment events for the manager
type ManagerReconciler struct {
	client.Client
}

func (r *ManagerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var deployment appsv1.Deployment
	if err := r.Get(ctx, req.NamespacedName, &deployment); err != nil {
		if client.IgnoreNotFound(err) == nil {
			logger.Info("Deployment deleted", "name", req.Name, "namespace", req.Namespace)
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	logger.Info("Deployment event",
		"name", deployment.Name,
		"namespace", deployment.Namespace,
		"replicas", *deployment.Spec.Replicas,
		"ready", deployment.Status.ReadyReplicas)

	return ctrl.Result{}, nil
}

func (r *ManagerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		Complete(r)
}

func runManager() error {
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	// Build kubeconfig path
	kubeconfig := managerKubeconfig
	if kubeconfig == "" {
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		}
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to build config: %w", err)
	}

	// Create scheme
	scheme := runtime.NewScheme()
	appsv1.AddToScheme(scheme)
	coordinationv1.AddToScheme(scheme)

	// Create manager
	mgr, err := ctrl.NewManager(config, ctrl.Options{
		Scheme: scheme,
		Metrics: server.Options{
			BindAddress: metricsAddr,
		},
		LeaderElection:             !disableLeaderElection,
		LeaderElectionID:           leaderElectionID,
		LeaderElectionNamespace:    managerNamespace,
		LeaderElectionResourceLock: "leases",
	})
	if err != nil {
		return fmt.Errorf("failed to create manager: %w", err)
	}

	// Setup controller
	if err := (&ManagerReconciler{
		Client: mgr.GetClient(),
	}).SetupWithManager(mgr); err != nil {
		return fmt.Errorf("failed to setup controller: %w", err)
	}

	ctrl.Log.Info("Starting controller manager",
		"namespace", managerNamespace,
		"leader_election", !disableLeaderElection,
		"metrics_addr", metricsAddr)

	return mgr.Start(ctrl.SetupSignalHandler())
}
