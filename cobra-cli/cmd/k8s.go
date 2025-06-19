package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// getKubernetesClient creates and returns a Kubernetes client
func getKubernetesClient() (*kubernetes.Clientset, error) {
	// Try in-cluster config first (if running in a pod)
	config, err := rest.InClusterConfig()
	if err != nil {
		// Fall back to kubeconfig file
		var kubeconfig string
		if home := homeDir(); home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		}

		// Use KUBECONFIG environment variable if set
		if kubeconfigEnv := os.Getenv("KUBECONFIG"); kubeconfigEnv != "" {
			kubeconfig = kubeconfigEnv
		}

		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create kubernetes config: %v", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %v", err)
	}

	return clientset, nil
}

// homeDir returns the home directory for the executing user
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

// getCurrentNamespace returns the current namespace from kubeconfig context
func getCurrentNamespace() string {
	var kubeconfig string
	if home := homeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	if kubeconfigEnv := os.Getenv("KUBECONFIG"); kubeconfigEnv != "" {
		kubeconfig = kubeconfigEnv
	}

	config, err := clientcmd.LoadFromFile(kubeconfig)
	if err != nil {
		return "default"
	}

	currentContext := config.CurrentContext
	if context, exists := config.Contexts[currentContext]; exists {
		if context.Namespace != "" {
			return context.Namespace
		}
	}

	return "default"
}
