package client

import (
	"context"
	"errors"
	"path"

	"github.com/k8sgpt-ai/k8sgpt/pkg/analysis"
	"github.com/spf13/viper"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/util/homedir"

	"ktt/backend/types"
)

var (
	kubeconfigPath = path.Join(homedir.HomeDir(), ".KTT", "ktt-kube-config")
)

type ClientService struct {
	ctx       context.Context
	apiClient *APIClient
}

func NewClientService() *ClientService {
	return &ClientService{}
}

func (s *ClientService) TestConnection(config string) types.JSResp {
	_, err := s.validate(config)
	if err != nil {
		return types.FailedResp(err.Error())
	}
	return types.JSResp{
		Success: true,
	}
}

func (s *ClientService) validate(config string) (*clientcmdapi.Config, error) {
	if len(config) == 0 {
		return nil, ErrEmptyConfig
	}
	content := []byte(config)
	apiConfig, err := clientcmd.Load(content)
	if err != nil {
		return nil, err
	}
	err = clientcmd.Validate(*apiConfig)
	if err != nil {
		return nil, err
	}
	return apiConfig, nil
}

func (s *ClientService) LoadConfig(configContent string) types.JSResp {
	userConfig, err := s.validate(configContent)
	if err != nil {
		return types.FailedResp(err.Error())
	}
	// write to file
	err = clientcmd.WriteToFile(*userConfig, kubeconfigPath)
	if err != nil {
		return types.FailedResp(err.Error())
	}
	flags := genericclioptions.NewConfigFlags(UsePersistentConfig)
	flags.KubeConfig = &kubeconfigPath
	s.apiClient = New(
		NewConfig(flags),
	)
	return types.JSResp{
		Success: true,
		Data:    s.GetContexts(),
	}
}

func (s *ClientService) Start(ctx context.Context) {
	s.ctx = ctx
}

func (s *ClientService) GetContexts() []string {
	config, err := s.apiClient.config.RawConfig()
	if err != nil {
		return nil
	}
	ctxs := config.Contexts
	slice := make([]string, 0, len(ctxs))
	for k := range ctxs {
		slice = append(slice, k)
	}
	return slice
}

func (s *ClientService) GetClusters() []string {
	ctxs, err := s.apiClient.config.Contexts()
	if err != nil {
		return nil
	}
	slice := make([]string, 0, len(ctxs))
	for _, ctx := range ctxs {
		slice = append(slice, ctx.Cluster)
	}
	return slice
}

func (s *ClientService) CurrentContext() string {
	return s.apiClient.ActiveContext()
}

func (s *ClientService) Analyze(cluster string, filters []string) types.JSResp {
	if len(cluster) == 0 {
		return types.FailedResp("cluster is empty")
	}
	viper.Set("kubecontext", cluster)
	viper.Set("kubeconfig", kubeconfigPath)

	analyzer, err := analysis.NewAnalysis(
		"noopai", "english", filters, "", "", false, false, 1, false, false, []string{},
	)
	if err != nil {
		return types.FailedResp(err.Error())
	}
	defer analyzer.Close()
	analyzer.RunAnalysis()
	if len(analyzer.Errors) > 0 {
		errs := make([]error, 0, len(analyzer.Errors))
		for _, err := range analyzer.Errors {
			errs = append(errs, errors.New(err))
		}
		return types.FailedResp(errors.Join(errs...).Error())
	}
	return types.JSResp{
		Success: true,
		Data:    analyzer.Results,
	}
}
