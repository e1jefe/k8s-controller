package cmd

import (
	"context"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

var (
	cfg       *rest.Config
	k8sClient kubernetes.Interface
	testEnv   *envtest.Environment
	ctx       context.Context
	cancel    context.CancelFunc
)

func TestInformer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Informer Suite")
}

var _ = BeforeSuite(func() {
	ctx, cancel = context.WithCancel(context.Background())

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{},
		ErrorIfCRDPathMissing: false,
	}

	var err error
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	k8sClient, err = kubernetes.NewForConfig(cfg)
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())
})

var _ = AfterSuite(func() {
	cancel()
	By("tearing down the test environment")
	if testEnv != nil {
		err := testEnv.Stop()
		Expect(err).NotTo(HaveOccurred())
	}
})

var _ = Describe("Deployment Informer", func() {

	Context("when testing basic functionality", func() {
		It("should create kubernetes client successfully", func() {
			// Test that we can create a client with the test config
			client, err := kubernetes.NewForConfig(cfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(client).NotTo(BeNil())

			// Test basic connectivity
			_, err = client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should create and configure informer successfully", func() {
			// Create informer factory
			informerFactory := informers.NewSharedInformerFactoryWithOptions(
				k8sClient,
				30*time.Second,
				informers.WithNamespace("default"),
			)

			// Get deployment informer
			deploymentInformer := informerFactory.Apps().V1().Deployments().Informer()
			Expect(deploymentInformer).NotTo(BeNil())

			// Test that we can add event handlers
			eventReceived := make(chan string, 1)
			deploymentInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
				AddFunc: func(obj interface{}) {
					deployment := obj.(*appsv1.Deployment)
					eventReceived <- "ADDED:" + deployment.Name
				},
			})

			// Create a simple test deployment to verify the informer works
			deployment := &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-deployment",
					Namespace: "default",
				},
				Spec: appsv1.DeploymentSpec{
					Replicas: int32Ptr(1),
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app": "test",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app": "test",
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "nginx",
									Image: "nginx:1.14.2",
								},
							},
						},
					},
				},
			}

			// Start informer in a separate goroutine
			stopCh := make(chan struct{})
			defer close(stopCh)

			go informerFactory.Start(stopCh)

			// Wait for cache sync
			Eventually(func() bool {
				return cache.WaitForCacheSync(stopCh, deploymentInformer.HasSynced)
			}, 10*time.Second).Should(BeTrue())

			// Create deployment
			_, err := k8sClient.AppsV1().Deployments("default").Create(ctx, deployment, metav1.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())

			// Wait for event (with timeout)
			Eventually(eventReceived, 10*time.Second).Should(Receive(Equal("ADDED:test-deployment")))

			// Clean up
			err = k8sClient.AppsV1().Deployments("default").Delete(ctx, "test-deployment", metav1.DeleteOptions{})
			Expect(err).NotTo(HaveOccurred())
		})
	})
})

// Helper function to get pointer to int32
func int32Ptr(i int32) *int32 { return &i }
