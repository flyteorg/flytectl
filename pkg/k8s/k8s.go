package k8s

import (
	"context"
	"os"

	"github.com/pkg/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	appsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

type K8s interface {
	AppsV1() appsv1.AppsV1Interface
	CoreV1() corev1.CoreV1Interface
}

type FlyteK8s struct {
	*kubernetes.Clientset
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

func GetCountOfReadyDeployment(ctx context.Context, client appsv1.AppsV1Interface) (int64, error) {
	deployments, err := client.Deployments("flyte").List(ctx, v1.ListOptions{})
	if err != nil {
		return 0, err
	}

	var count int64
	for _, dep := range deployments.Items {
		if dep.Status.AvailableReplicas == 1 {
			count++
			continue
		}
	}
	return count, nil
}

func GetFlyteDeploymentCount(ctx context.Context, client appsv1.AppsV1Interface) (int64, error) {
	deployments, err := client.Deployments("flyte").List(ctx, v1.ListOptions{})
	if err != nil {
		return 0, err
	}
	return int64(len(deployments.Items)), nil
}

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
