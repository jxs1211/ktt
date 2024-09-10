package k8sutil

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	clientv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
)

func NewConfig() (*rest.Config, error) {
	var config *rest.Config
	var err error
	var kubeconfig string

	// Try in-cluster config first
	config, err = rest.InClusterConfig()
	if err != nil {
		// If in-cluster fails, try kubeconfig
		if flag.Lookup("kubeconfig") == nil {
			flag.StringVar(&kubeconfig, "kubeconfig", os.Getenv("KUBECONFIG"), "location to your kubeconfig file")
		}
		flag.Parse() // Ensure flags are parsed
		if kubeconfig == "" {
			kubeconfig = clientcmd.RecommendedHomeFile
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create config: %v", err)
		}
	}
	if config == nil {
		return nil, fmt.Errorf("failed to create kubernetes config")
	}
	return config, nil
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func NewDynamicClient(config *rest.Config) (*dynamic.DynamicClient, error) {
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewClient(config *rest.Config) (*kubernetes.Clientset, error) {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

type LabelSelector string

const (
	LabelServiceSelector LabelSelector = "service=lts-eks-cluster"
)

// String implements the flag.Value interface
func (ls *LabelSelector) String() string {
	if ls == nil {
		return ""
	}

	return string(*ls)
}

// Set implements the flag.Value interface.
func (ls *LabelSelector) Set(value string) error {
	if _, err := labels.Parse(value); err != nil {
		return err
	}

	*ls = LabelSelector(value)
	return nil
}

func (ls *LabelSelector) Map() map[string]string {
	res := make(map[string]string)
	if ls == nil {
		return res
	}
	labels := string(*ls)
	if len(labels) == 0 {
		return res
	}
	lbs := strings.Split(labels, ",")
	for _, v := range lbs {
		kv := strings.Split(v, "=")
		res[kv[0]] = kv[1]
	}
	return res
}

func LabelRequirement(label map[string]string) []labels.Requirement {
	res := make([]labels.Requirement, 0, len(label))
	if len(label) == 0 {
		return res
	}
	for k, v := range label {
		req, _ := labels.NewRequirement(k, selection.Equals, []string{v})
		res = append(res, *req)
	}
	return res
}

func GetClusterName(ctx context.Context, clientset kubernetes.Interface) (string, error) {
	nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return "", err
	}
	items := nodes.Items
	if len(items) == 0 {
		return "", errors.New("no nodes found")
	}
	for i := 0; i < len(items); i++ {
		if clusterName, ok := items[i].ObjectMeta.Labels["alpha.eksctl.io/cluster-name"]; ok {
			return clusterName, nil
		}
	}
	return "", errors.New("can't found cluster name")
}

func CreateOrUpdateEndpoints(ctx context.Context, eclient clientv1.EndpointsInterface, eps *v1.Endpoints) error {
	// As stated in the RetryOnConflict's documentation, the returned error shouldn't be wrapped.
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		_, err := eclient.Get(ctx, eps.Name, metav1.GetOptions{})
		if err != nil {
			if !apierrors.IsNotFound(err) {
				return err
			}
			_, err = eclient.Create(ctx, eps, metav1.CreateOptions{})
			return err
		}
		_, err = eclient.Update(ctx, eps, metav1.UpdateOptions{})
		return err
	})
}

func LabelValue(labels map[string]string, key string) string {
	if v, ok := labels[key]; ok {
		return v
	}
	return ""
}

func NameWith(projectName string) string {
	return fmt.Sprintf("%s-kubelet", projectName)
}

func NameAsK8SFormat(input string) string {
	// Replace underscores with hyphens
	output := strings.ReplaceAll(input, "_", "-")
	// Convert uppercase to lowercase
	output = strings.ToLower(output)
	return output
}

func Contains(slice []string, element string) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

func RemoveDuplicates(slice []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

func MatchUUID(input string) string {
	if len(input) == 0 {
		return ""
	}
	pattern := `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`
	regex := regexp.MustCompile(pattern)
	res := regex.FindString(input)
	return res
}

// nodeAddress returns the provided node's address, based on the priority:
// 1. NodeInternalIP
// 2. NodeExternalIP
//
// Copied from github.com/prometheus/prometheus/discovery/kubernetes/node.go
func nodeAddress(node *v1.Node) (string, map[v1.NodeAddressType][]string, error) {
	m := map[v1.NodeAddressType][]string{}
	for _, a := range node.Status.Addresses {
		m[a.Type] = append(m[a.Type], a.Address)
	}

	if addresses, ok := m[v1.NodeInternalIP]; ok {
		return addresses[0], m, nil
	}
	if addresses, ok := m[v1.NodeExternalIP]; ok {
		return addresses[0], m, nil
	}
	return "", m, fmt.Errorf("host address unknown")
}

func NodeAddresses(nodes []*v1.Node) ([]v1.EndpointAddress, []error) {
	addresses := make([]v1.EndpointAddress, 0)
	errs := make([]error, 0)
	if nodes == nil {
		return addresses, errs
	}
	for _, n := range nodes {
		address, _, err := nodeAddress(n)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to determine hostname for node (%s): %w", n.Name, err))
			continue
		}
		addresses = append(addresses, v1.EndpointAddress{
			IP: address,
			TargetRef: &v1.ObjectReference{
				Kind:       "Node",
				Name:       n.Name,
				UID:        n.UID,
				APIVersion: n.APIVersion,
			},
		})
	}

	return addresses, errs
}
