package k8s

import (
	"context"
	"os"

	"github.com/pkg/errors"
	corev1api "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	appsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	flyteNamespace = "flyte"
)

type K8s interface {
	AppsV1() appsv1.AppsV1Interface
	CoreV1() corev1.CoreV1Interface
}

var Client K8s

// GetK8sClient return the k8s client from sandbox kubeconfig
func GetK8sClient(cfg, master string) (K8s, error) {
	kubeConfigPath := os.ExpandEnv(cfg)
	kubecfg, err := clientcmd.BuildConfigFromFlags(master, kubeConfigPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Error building kubeconfig")
	}
	if Client == nil {
		kubeClient, err := kubernetes.NewForConfig(kubecfg)
		if err != nil {
			return nil, errors.Wrapf(err, "Error building kubernetes clientset")
		}
		return kubeClient, nil
	}
	return Client, nil
}

// GetFlyteDeployment return the pod list from flyte namespace
func GetFlyteDeployment(ctx context.Context, client corev1.CoreV1Interface) (*corev1api.PodList, error) {
	pods, err := client.Pods(flyteNamespace).List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return pods, nil
}

// GetNodeTaintStatus check disk pressure taint in node
func GetNodeTaintStatus(ctx context.Context, client corev1.NodeInterface) (bool, error) {
	nodes, err := client.List(ctx, v1.ListOptions{})
	if err != nil {
		return false, err
	}
	match := 0
	for _, node := range nodes.Items {
		for _, c := range node.Spec.Taints {
			if c.Key == "node.kubernetes.io/disk-pressure" && c.Effect == "NoSchedule" {
				match++
			}
		}
	}
	if match > 0 {
		return true, nil
	}
	return false, nil
}
