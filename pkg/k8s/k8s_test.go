package k8s

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testclient "k8s.io/client-go/kubernetes/fake"
)

var fakeNode = &corev1.Node{
	Spec: corev1.NodeSpec{
		Taints: []corev1.Taint{},
	},
}

func TestGetFlyteDeploymentCount(t *testing.T) {
	ctx := context.Background()
	client := testclient.NewSimpleClientset()
	c, err := GetFlyteDeployment(ctx, client.CoreV1())
	assert.Nil(t, err)
	assert.Equal(t, 0, len(c.Items))
}

func TestGetNodeTaintStatus(t *testing.T) {
	t.Run("Check node taint with success", func(t *testing.T) {
		ctx := context.Background()
		client := testclient.NewSimpleClientset()
		fakeNode.SetName("master")
		_, err := client.CoreV1().Nodes().Create(ctx, fakeNode, v1.CreateOptions{})
		if err != nil {
			t.Error(err)
		}
		c, err := GetNodeTaintStatus(ctx, client.CoreV1().Nodes())
		assert.Nil(t, err)
		assert.Equal(t, false, c)
	})
	t.Run("Check node taint with fail", func(t *testing.T) {
		ctx := context.Background()
		client := testclient.NewSimpleClientset()
		fakeNode.SetName("master")
		_, err := client.CoreV1().Nodes().Create(ctx, fakeNode, v1.CreateOptions{})
		if err != nil {
			t.Error(err)
		}
		node, err := client.CoreV1().Nodes().Get(ctx, "master", v1.GetOptions{})
		if err != nil {
			t.Error(err)
		}
		node.Spec.Taints = append(node.Spec.Taints, corev1.Taint{
			Effect: "NoSchedule",
			Key:    "node.kubernetes.io/disk-pressure",
		})
		_, err = client.CoreV1().Nodes().Update(ctx, node, v1.UpdateOptions{})
		if err != nil {
			t.Error(err)
		}
		c, err := GetNodeTaintStatus(ctx, client.CoreV1().Nodes())
		assert.Nil(t, err)
		assert.Equal(t, true, c)
	})
}

func TestGetK8sClient(t *testing.T) {
	content := `
apiVersion: v1
clusters:
- cluster:
    server: https://localhost:8080
    extensions:
    - name: client.authentication.k8s.io/exec
      extension:
        audience: foo
        other: bar
  name: foo-cluster
contexts:
- context:
    cluster: foo-cluster
    user: foo-user
    namespace: bar
  name: foo-context
current-context: foo-context
kind: Config
users:
- name: foo-user
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1alpha1
      args:
      - arg-1
      - arg-2
      command: foo-command
      provideClusterInfo: true
`
	tmpfile, err := ioutil.TempFile("", "kubeconfig")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(tmpfile.Name())
	if err := ioutil.WriteFile(tmpfile.Name(), []byte(content), os.ModePerm); err != nil {
		t.Error(err)
	}
	t.Run("Create client from config", func(t *testing.T) {
		client := testclient.NewSimpleClientset()
		Client = client
		c, err := GetK8sClient(tmpfile.Name(), "https://localhost:8080")
		assert.Nil(t, err)
		assert.NotNil(t, c)
	})
	t.Run("Create client from config", func(t *testing.T) {
		Client = nil
		client, err := GetK8sClient(tmpfile.Name(), "https://localhost:8080")
		assert.Nil(t, err)
		assert.NotNil(t, client)
	})

}
