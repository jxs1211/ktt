package client

import (
	"context"
	"ktt/backend/types"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/k8sgpt-ai/k8sgpt/pkg/ai"
	"github.com/spf13/viper"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd/api"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func makeApiClient() *APIClient {
	flags := genericclioptions.NewConfigFlags(UsePersistentConfig)
	flags.KubeConfig = &kubeconfigPath
	// config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// c, err := kubernetes.NewForConfig(config)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	client := New(
		NewConfig(flags),
	)

	return client
}

func TestConnection(t *testing.T) {
	tests := []struct {
		name      string
		config    string
		ctx       context.Context
		apiClient *APIClient
		want      types.JSResp
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ClientService{
				ctx:       tt.ctx,
				apiClient: tt.apiClient,
			}
			if got := s.TestConnection(tt.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClientService.TestConnection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClientService_validate(t *testing.T) {
	tests := []struct {
		name      string
		ctx       context.Context
		apiClient *APIClient
		config    string
		want      *clientcmdapi.Config
		wantErr   bool
	}{
		{
			name:      "base",
			ctx:       context.Background(),
			apiClient: makeApiClient(),
			config: `apiVersion: v1
kind: Config
clusters:
- cluster:
    server: https://1.1.1.1
  name: production
contexts:
- context:
    cluster: production
    user: production
  name: production
current-context: production
users:
- name: production
  user:
    auth-provider:
      name: gcp`,
			want: &clientcmdapi.Config{
				CurrentContext: "production",
				Extensions:     map[string]k8sruntime.Object{},
				Preferences: api.Preferences{
					Extensions: map[string]k8sruntime.Object{},
				},
				Contexts: map[string]*clientcmdapi.Context{
					"production": {
						Cluster:    "production",
						AuthInfo:   "production",
						Extensions: map[string]k8sruntime.Object{},
					},
				},
				AuthInfos: map[string]*clientcmdapi.AuthInfo{
					"production": {
						AuthProvider: &clientcmdapi.AuthProviderConfig{
							Name: "gcp",
						},
						Extensions: map[string]k8sruntime.Object{},
					},
				},
				Clusters: map[string]*clientcmdapi.Cluster{
					"production": {
						Server:     "https://1.1.1.1",
						Extensions: map[string]k8sruntime.Object{},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ClientService{
				ctx:       tt.ctx,
				apiClient: tt.apiClient,
			}
			got, err := s.validate(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClientService.validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			diff := cmp.Diff(got, tt.want)
			t.Log("diff: ", diff)
			if len(diff) != 0 {
				t.Errorf("The diff of ClientService.validate() = %v", diff)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClientService.validate() = %+v, \nwant %+v", got, tt.want)
			}
		})
	}
}

func TestClientService_LoadConfig(t *testing.T) {
	type fields struct {
		ctx       context.Context
		apiClient *APIClient
	}
	type args struct {
		configContent string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   types.JSResp
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ClientService{
				ctx:       tt.fields.ctx,
				apiClient: tt.fields.apiClient,
			}
			if got := s.LoadConfig(tt.args.configContent); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClientService.LoadConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClientService_Start(t *testing.T) {
	type fields struct {
		ctx       context.Context
		apiClient *APIClient
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ClientService{
				ctx:       tt.fields.ctx,
				apiClient: tt.fields.apiClient,
			}
			s.Start(tt.args.ctx)
		})
	}
}

func TestClientService_GetLocalConfig(t *testing.T) {
	type fields struct {
		ctx       context.Context
		apiClient *APIClient
	}
	tests := []struct {
		name   string
		fields fields
		want   types.JSResp
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ClientService{
				ctx:       tt.fields.ctx,
				apiClient: tt.fields.apiClient,
			}
			if got := s.GetLocalConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClientService.GetLocalConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClientService_GetContexts(t *testing.T) {
	type fields struct {
		ctx       context.Context
		apiClient *APIClient
	}
	tests := []struct {
		name   string
		fields fields
		want   types.JSResp
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ClientService{
				ctx:       tt.fields.ctx,
				apiClient: tt.fields.apiClient,
			}
			if got := s.GetContexts(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClientService.GetContexts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClientService_getContexts(t *testing.T) {
	type fields struct {
		ctx       context.Context
		apiClient *APIClient
	}
	tests := []struct {
		name    string
		fields  fields
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ClientService{
				ctx:       tt.fields.ctx,
				apiClient: tt.fields.apiClient,
			}
			got, err := s.getContexts()
			if (err != nil) != tt.wantErr {
				t.Errorf("ClientService.getContexts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClientService.getContexts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClientService_GetClusters(t *testing.T) {
	type fields struct {
		ctx       context.Context
		apiClient *APIClient
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ClientService{
				ctx:       tt.fields.ctx,
				apiClient: tt.fields.apiClient,
			}
			if got := s.GetClusters(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClientService.GetClusters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClientService_CurrentContext(t *testing.T) {
	type fields struct {
		ctx       context.Context
		apiClient *APIClient
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ClientService{
				ctx:       tt.fields.ctx,
				apiClient: tt.fields.apiClient,
			}
			if got := s.CurrentContext(); got != tt.want {
				t.Errorf("ClientService.CurrentContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func TestClientService_Analyze(t *testing.T) {
// 	tests := []struct {
// 		name   string
// 		ctx       context.Context
// 		apiClient *APIClient
// 		cluster string
// 		filters []string
// 		want   types.JSResp
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := &ClientService{
// 				ctx:       tt.ctx,
// 				apiClient: tt.apiClient,
// 			}
// 			if got := s.Analyze(tt.cluster, tt.filters); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("ClientService.Analyze() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func TestClientService_getApiResources(t *testing.T) {
	type fields struct {
		ctx       context.Context
		apiClient *APIClient
	}
	tests := []struct {
		name    string
		fields  fields
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ClientService{
				ctx:       tt.fields.ctx,
				apiClient: tt.fields.apiClient,
			}
			got, err := s.getApiResources()
			if (err != nil) != tt.wantErr {
				t.Errorf("ClientService.getApiResources() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClientService.getApiResources() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClientService_getAvailableFilteredResources(t *testing.T) {
	tests := []struct {
		name      string
		ctx       context.Context
		apiClient *APIClient
		want      []string
		wantErr   bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ClientService{
				ctx:       tt.ctx,
				apiClient: tt.apiClient,
			}
			got, err := s.getAvailableFilteredResources()
			if (err != nil) != tt.wantErr {
				t.Errorf("ClientService.getAvailableFilteredResources() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClientService.getAvailableFilteredResources() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClientService_GetAvailableFilteredResources(t *testing.T) {
	type fields struct {
		ctx       context.Context
		apiClient *APIClient
	}
	tests := []struct {
		name   string
		fields fields
		want   types.JSResp
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ClientService{
				ctx:       tt.fields.ctx,
				apiClient: tt.fields.apiClient,
			}
			if got := s.GetAvailableFilteredResources(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClientService.GetAvailableFilteredResources() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestViper(t *testing.T) {
	var configAI ai.AIConfiguration
	err := viper.UnmarshalKey("ai", &configAI)
	t.Logf("AI config: %+v", configAI)
	if err != nil {
		t.Fatal(err)
	}
	viper.Set("key2", ai.AIConfiguration{
		DefaultProvider: "google",
	})
	err = viper.UnmarshalKey("key", &configAI)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("AI config: %+v", configAI)
}

func TestClientService_loadConfigFrom(t *testing.T) {
	tests := []struct {
		name      string
		ctx       context.Context
		apiClient *APIClient
		path      string
		wantErr   bool
	}{
		{
			name:      "base",
			ctx:       context.Background(),
			apiClient: makeApiClient(),
			path:      kubeconfigPath,
		},
		{
			name:      "configPathNotExist",
			ctx:       context.Background(),
			apiClient: makeApiClient(),
			path:      "/path/not/exist",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ClientService{
				ctx:       tt.ctx,
				apiClient: tt.apiClient,
			}
			if err := s.loadConfigFrom(tt.path); (err != nil) != tt.wantErr {
				t.Errorf("ClientService.loadConfigFrom() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
