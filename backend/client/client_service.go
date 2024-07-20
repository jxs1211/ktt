package client

import (
	"context"
	"ktt/backend/types"
	"path"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	kubeconfigPath = path.Join(homedir.HomeDir(), "KTT", "ktt-kube-config")
)

type ClientService struct {
	// client *APIClient
	client *kubernetes.Clientset
	ctx    context.Context
}

func NewClientService() *ClientService {
	return &ClientService{}
}

// func (s *ClientService) Load(configPath string) {
// 	k8sFlags := genericclioptions.NewConfigFlags(UsePersistentConfig)
// 	k8sFlags.KubeConfig = &configPath
// 	k8sCfg := NewConfig(k8sFlags)
// 	client := New(k8sCfg)
// 	s.client = client
// }

func (s *ClientService) LoadConfig(configContent string) types.JSResp {
	content := []byte(configContent)
	config, err := clientcmd.Load(content)
	if err != nil {
		return types.FailedResp(err.Error())
	}
	err = clientcmd.Validate(*config)
	if err != nil {
		return types.FailedResp(err.Error())
	}
	// write to file
	err = clientcmd.WriteToFile(*config, kubeconfigPath)
	if err != nil {
		return types.FailedResp(err.Error())
	}
	restClientConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return types.FailedResp(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(restClientConfig)
	if err != nil {
		return types.FailedResp(err.Error())
	}
	s.client = clientset

	return types.JSResp{
		Success: true,
	}
}

func (s *ClientService) Start(ctx context.Context) {
	s.ctx = ctx
}
