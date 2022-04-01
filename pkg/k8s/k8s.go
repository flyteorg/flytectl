package k8s

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type K8s interface {
	CoreV1() corev1.CoreV1Interface
}

//go:generate mockery -name=ContextOps -case=underscore
type ContextOps interface {
	CopyContext(srcConfigAccess clientcmd.ConfigAccess, srcCtxName, targetCtxName string) error
	RemoveContext(ctxName string) error
}

// K8sContextManager context manager implementing ContextOps
type K8sContextManager struct {
	configAccess clientcmd.ConfigAccess
}

func NewK8sContextManager() ContextOps {
	if ContextMgr != nil {
		return ContextMgr
	}
	ContextMgr = &K8sContextManager{
		configAccess: clientcmd.NewDefaultPathOptions(),
	}
	return ContextMgr
}

var Client K8s
var ContextMgr ContextOps

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

// CopyKubeContext copies context srcCtxName part of srcConfigAccess to targetCtxName part of targetConfigAccess.
func (k *K8sContextManager) CopyContext(srcConfigAccess clientcmd.ConfigAccess, srcCtxName, targetCtxName string) error {
	_, err := k.configAccess.GetStartingConfig()
	if err != nil {
		return err
	}

	fromStartingConfig, err := srcConfigAccess.GetStartingConfig()
	if err != nil {
		return err
	}
	_, exists := fromStartingConfig.Contexts[srcCtxName]
	if !exists {
		return fmt.Errorf("context %v doesn't exist", srcCtxName)
	}

	toStartingConfig, err := k.configAccess.GetStartingConfig()
	if err != nil {
		return err
	}

	_, exists = toStartingConfig.Contexts[targetCtxName]
	if exists {
		fmt.Printf("context %v already exist. Overwriting it\n", targetCtxName)
	} else {
		toStartingConfig.Contexts[targetCtxName] = clientcmdapi.NewContext()
	}

	toStartingConfig.Clusters[targetCtxName] = fromStartingConfig.Clusters[srcCtxName]
	toStartingConfig.Clusters[targetCtxName].LocationOfOrigin = k.configAccess.GetDefaultFilename()
	toStartingConfig.AuthInfos[targetCtxName] = fromStartingConfig.AuthInfos[srcCtxName]
	toStartingConfig.AuthInfos[targetCtxName].LocationOfOrigin = k.configAccess.GetDefaultFilename()
	toStartingConfig.Contexts[targetCtxName].Cluster = targetCtxName
	toStartingConfig.Contexts[targetCtxName].AuthInfo = targetCtxName
	toStartingConfig.CurrentContext = targetCtxName
	if err := clientcmd.ModifyConfig(k.configAccess, *toStartingConfig, true); err != nil {
		return err
	}

	fmt.Printf("context modified for %q and switched over to it.\n", targetCtxName)
	return nil
}

// RemoveKubeContext removes the contextToRemove from the kubeContext pointed to be fromConfigAccess
func (k *K8sContextManager) RemoveContext(ctxName string) error {
	fromStartingConfig, err := k.configAccess.GetStartingConfig()
	if err != nil {
		return err
	}
	_, exists := fromStartingConfig.Contexts[ctxName]
	if !exists {
		return fmt.Errorf("context %v doesn't exist", ctxName)
	}

	delete(fromStartingConfig.Clusters, ctxName)
	delete(fromStartingConfig.AuthInfos, ctxName)
	delete(fromStartingConfig.Contexts, ctxName)
	fromStartingConfig.CurrentContext = ""

	if err := clientcmd.ModifyConfig(k.configAccess, *fromStartingConfig, true); err != nil {
		return err
	}

	fmt.Printf("context removed for %q.\n", ctxName)
	return nil
}
