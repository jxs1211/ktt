package watch

import (
	"context"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"testing"
	"time"

	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"

	k8sutil "ktt/backend/utils/k8s"
)

func Test_watchingTypesString(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := watchingTypesString(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("watchingTypesString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func createKubeClient(t *testing.T) *kubernetes.Clientset {
	config, err := k8sutil.NewConfig()
	if err != nil {
		t.Fatal(err)
	}
	client, err := k8sutil.NewClient(config)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func TestCreateKubeClient(t *testing.T) {
	c := createKubeClient(t)
	if c == nil {
		t.Fatal("client is nil")
	}
	t.Log(c)
}

const podManifest = `
apiVersion: v1
kind: Pod
metadata:
  name: cleaner
	namespace: local
spec:
  containers:
  - name: nginx
    image: nginx:latest
    ports:
    - containerPort: 80
`

const wrongPodManifest = `
apiVersion: v1
kind: Pod
metadata:
  name: cleaner
	namespace: local
spec:
  containers:
  - name: nginx
    image: nginx:latest2
    ports:
    - containerPort: 80
`

func createObject(w *Watcher, ns, name, imageName string) (runtime.Object, error) {
	return w.kclient.CoreV1().Pods(ns).Create(context.Background(), &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "nginx",
					Image: imageName,
					Ports: []v1.ContainerPort{
						{
							ContainerPort: 80,
						},
					},
				},
			},
		},
	}, metav1.CreateOptions{})
}

func TestWatcher_Run(t *testing.T) {
	client := createKubeClient(t)
	logr, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}
	ns, name := "default", "cleaner"
	w := NewWatcher(
		"test", client,
		"pod", ns, name,
		watchingTypesString(), logr)

	go w.Run(1)
	_, err = createObject(w, ns, name, "nginx:latest")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(5 * time.Second)
	// deletion
	err = w.kclient.CoreV1().Pods(ns).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(2 * time.Second)
	w.Stop()

	go w.Run(1)
	_, err = createObject(w, ns, name, "nginx:latest2")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(5 * time.Second)
	// deletion
	err = w.kclient.CoreV1().Pods(ns).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(2 * time.Second)
	w.Stop()
}

func TestWatchRunInAdvance(t *testing.T) {
	client := createKubeClient(t)
	logr, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}
	ns, name := "default", "cleaner"
	w := NewWatcher(
		"test", client,
		"pod", ns, name,
		watchingTypesString(), logr)

	go w.Run(1)
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	sig := <-term
	logr.Info("exiting gracefully... ", zap.String("signal", sig.String()))
	w.Stop()

	time.Sleep(1 * time.Second)
	go w.Run(1)

	sig = <-term
	logr.Info("exiting gracefully... ", zap.String("signal", sig.String()))
	w.Stop()
}
