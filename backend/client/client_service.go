package client

import (
	"context"
	"errors"
	"log"
	"os"
	"path"
	"strings"

	"github.com/k8sgpt-ai/k8sgpt/pkg/ai"
	"github.com/k8sgpt-ai/k8sgpt/pkg/analysis"
	"github.com/k8sgpt-ai/k8sgpt/pkg/analyzer"
	"github.com/k8sgpt-ai/k8sgpt/pkg/common"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"ktt/backend/types"
	sliceutil "ktt/backend/utils/slice"
	strutil "ktt/backend/utils/string"
)

var (
	kubeconfigPath = path.Join(strutil.RootPath(), "ktt-kube-config")
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
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return types.FailedResp(err.Error())
	}
	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		return types.FailedResp(err.Error())
	}
	s.apiClient = New(
		NewConfig(flags),
	)
	s.apiClient.setClient(c)
	ctxs, err := s.getContexts()
	if err != nil {
		return types.FailedResp(err.Error())
	}
	return types.JSResp{
		Success: true,
		Data:    ctxs,
	}
}

func (s *ClientService) Start(ctx context.Context) {
	s.ctx = ctx
}

func (s *ClientService) GetLocalConfig() types.JSResp {
	kubeconfigBytes, err := os.ReadFile(kubeconfigPath)
	if err != nil {
		return types.FailedResp(err.Error())
	}
	_, err = s.validate(string(kubeconfigBytes))
	if err != nil {
		return types.FailedResp(err.Error())
	}
	return types.JSResp{
		Success: true,
		Data:    string(kubeconfigBytes),
	}
}

func (s *ClientService) GetContexts() types.JSResp {
	ctxs, err := s.getContexts()
	if err != nil {
		return types.FailedResp(err.Error())
	}
	return types.JSResp{
		Success: true,
		Data:    ctxs,
	}
}

func (s *ClientService) getContexts() ([]string, error) {
	config, err := s.apiClient.config.RawConfig()
	if err != nil {
		return nil, err
	}
	ctxs := config.Contexts
	slice := make([]string, 0, len(ctxs))
	for k := range ctxs {
		if !strings.Contains(k, config.CurrentContext) {
			slice = append(slice, k)
		}
	}
	// keep current context at the top
	res := make([]string, 0, len(ctxs))
	res = append(res, config.CurrentContext)
	res = append(res, slice...)
	return res, nil
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

// feat: support to analyze with ai
func (s *ClientService) Analyze(
	cluster, aiBackend, model string,
	filters []string, explain bool) types.JSResp {
	if len(cluster) == 0 {
		return types.FailedResp("cluster is empty")
	}
	viper.Set("kubecontext", cluster)
	viper.Set("kubeconfig", kubeconfigPath)
	log.Println("explain: ", explain)
	if explain {
		viper.Set("ai", ai.AIConfiguration{
			DefaultProvider: aiBackend,
			Providers: []ai.AIProvider{
				ai.AIProvider{
					Name:  aiBackend,
					Model: model,
				},
			},
		})
	}

	analyzer, err := analysis.NewAnalysis(
		aiBackend, "english", filters, "", "", false,
		explain, 1, false, false, []string{},
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
	if analyzer.Results == nil {
		analyzer.Results = []common.Result{}
	}
	anonymize := true
	if explain {
		if err := analyzer.GetAIResults("json", anonymize); err != nil {
			return types.FailedResp(err.Error())
		}
	}
	return types.JSResp{
		Success: true,
		Data:    analyzer.Results,
	}
}

// feat: implement api-resources API
// func (s *ClientService) GetApiResources() types.JSResp {
// 	resources, err := s.getApiResources()
// 	if err != nil {
// 		return types.FailedResp(err.Error())
// 	}
// 	return types.JSResp{
// 		Success: true,
// 		Data:    resources,
// 	}
// }

// feat: implement api-resources API
func (s *ClientService) getApiResources() ([]string, error) {
	lists, err := s.apiClient.client.Discovery().ServerPreferredResources()
	if err != nil {
		return nil, err
	}
	res := make([]string, 0, len(lists))
	for _, list := range lists {
		if len(list.APIResources) == 0 {
			continue
		}
		_, err := schema.ParseGroupVersion(list.GroupVersion)
		if err != nil {
			continue
		}
		for _, resource := range list.APIResources {
			if len(resource.Verbs) == 0 {
				continue
			}
			res = append(res, resource.Kind)
		}
	}
	return res, nil
}

// feat: implement available api-resources API
func (s *ClientService) getAvailableFilteredResources() ([]string, error) {
	apiResources, err := s.getApiResources()
	if err != nil {
		return nil, err
	}
	coreResources, _, _ := analyzer.ListFilters()
	mergedRes := sliceutil.FindCommonElements(coreResources, apiResources)
	return mergedRes, nil
}

// feat: implement available api-resources API
func (s *ClientService) GetAvailableFilteredResources() types.JSResp {
	resources, err := s.getAvailableFilteredResources()
	if err != nil {
		return types.FailedResp(err.Error())
	}
	return types.JSResp{
		Success: true,
		Data:    resources,
	}
}
