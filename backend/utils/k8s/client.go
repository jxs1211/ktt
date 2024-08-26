package k8s

import (
	"flag"
	"os"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func NewConfig() (*rest.Config, error) {
	var config *rest.Config
	config, err := rest.InClusterConfig()
	if err != nil {
		kubeconfig := flag.String("kubeconfig", GetEnv("KUBECONFIG", clientcmd.RecommendedConfigDir), "location to your kubeconfig file")
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			// handle error
			return nil, err
		}
	}
	return config, nil
}
