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
	"ktt/backend/utils/tool"
)

var (
	kubeconfigPath        = path.Join(strutil.RootPath(), "ktt-kube-config")
	defaultMaxConcurrency = 20
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
	cluster, aiBackend, model, baseURL string,
	filters []string, explain, aggregate, anonymize bool,
) types.JSResp {
	results, err := s.analyze(
		cluster, aiBackend, model, baseURL,
		filters, explain, aggregate, anonymize,
	)
	if err != nil {
		return types.FailedResp(err.Error())
	}
	return types.JSResp{
		Success: true,
		Data:    results,
	}
}

func (s *ClientService) GetErrorsCount(
	cluster, aiBackend, model, baseURL string,
	filters []string, explain, aggregate, anonymize bool,
) types.JSResp {
	count, err := s.getErrorsCount(
		cluster, aiBackend, model, baseURL,
		filters, explain, aggregate, anonymize,
	)
	if err != nil {
		return types.FailedResp(err.Error())
	}
	return types.JSResp{
		Success: true,
		Data:    count,
	}
}

func (s *ClientService) getErrorsCount(
	cluster, aiBackend, model, baseURL string,
	filters []string, explain, aggregate, anonymize bool,
) (int, error) {
	results, err := s.analyze(
		cluster, aiBackend, model, baseURL,
		filters, explain, aggregate, anonymize,
	)
	if err != nil {
		return 0, err
	}
	return len(results), nil
}

func (s *ClientService) analyze(
	cluster, aiBackend, model, baseURL string,
	filters []string, explain, aggregate, anonymize bool,
) ([]Result, error) {
	defer tool.TrackTime("clientService.analyze")()
	if len(cluster) == 0 {
		return nil, errors.New("no cluster available, add any first")
	}
	viper.Set("kubecontext", cluster)
	viper.Set("kubeconfig", kubeconfigPath)
	log.Println("explain: ", explain)
	if explain {
		if len(baseURL) <= 0 && aiBackend == "localai" {
			return nil, errors.New("baseURL is required for localai")
		}
		viper.Set("ai", ai.AIConfiguration{
			DefaultProvider: aiBackend,
			Providers: []ai.AIProvider{
				ai.AIProvider{
					Name:    aiBackend,
					Model:   model,
					BaseURL: baseURL,
				},
			},
		})
	}

	analyzer, err := analysis.NewAnalysis(
		aiBackend, "english", filters, "", "", false,
		explain, defaultMaxConcurrency, false, false, []string{},
	)
	if err != nil {
		return nil, err
	}
	defer analyzer.Close()
	analyzer.RunAnalysis()

	if len(analyzer.Errors) > 0 {
		errs := make([]error, 0, len(analyzer.Errors))
		for _, err := range analyzer.Errors {
			errs = append(errs, errors.New(err))
		}
		return nil, errors.Join(errs...)
	}
	if analyzer.Results == nil {
		analyzer.Results = []common.Result{}
	}
	if explain {
		if err := analyzer.GetAIResults("json", anonymize); err != nil {
			return nil, err
		}
	}
	results := parseDetail(analyzer.Results)
	return results, nil
}

func parseDetail(results []common.Result) []Result {
	res := make([]Result, 0, len(results))
	for _, result := range results {
		parts := strings.Split(result.Details, "\n\nSolution: \n")
		detail := Detail{}
		if len(parts) > 0 {
			detail.Error = strings.TrimPrefix(parts[0], "Error: ")
		}
		if len(parts) > 1 {
			solutionSteps := strings.Split(parts[1], "\n")
			for _, step := range solutionSteps {
				if step != "" {
					detail.Solution = append(detail.Solution, strings.TrimSpace(step))
				}
			}
		}
		if len(detail.Solution) == 0 {
			detail.Solution = []string{}
		}
		res = append(res, Result{
			Kind:         result.Kind,
			Name:         result.Name,
			Error:        result.Error,
			Details:      detail,
			ParentObject: result.ParentObject,
		})
	}
	return res
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

func (s *ClientService) CheckConnectivity(name string) types.JSResp {
	err := s.apiClient.SwitchContext(name)
	if err != nil {
		return types.FailedResp(err.Error())
	}
	return types.JSResp{
		Success: true,
		Data:    "connected to cluster",
	}
}

func (s *ClientService) GetClusterInfo(
	cluster, aiBackend, model, baseUrl string,
	filters []string, explain, aggregate, anonymize bool,
) types.JSResp {
	info, err := s.apiClient.ServerVersion()
	if err != nil {
		return types.FailedResp(err.Error())
	}
	count, err := s.getErrorsCount(
		cluster, aiBackend, model, baseUrl,
		filters, explain, aggregate, anonymize,
	)
	if err != nil {
		return types.FailedResp(err.Error())
	}
	clusterInfo := ClusterInfo{
		VersionInfo: *info,
		ErrorsCount: count,
	}
	return types.JSResp{
		Success: true,
		Data:    clusterInfo,
	}
}
