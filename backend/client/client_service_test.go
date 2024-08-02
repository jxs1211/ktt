package client

import (
	"context"
	"ktt/backend/types"
	"reflect"
	"testing"

	"github.com/k8sgpt-ai/k8sgpt/pkg/ai"
	"github.com/spf13/viper"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func TestNewClientService(t *testing.T) {
	tests := []struct {
		name string
		want *ClientService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClientService(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClientService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConnection(t *testing.T) {
	type fields struct {
		ctx       context.Context
		apiClient *APIClient
	}
	type args struct {
		config string
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
			if got := s.TestConnection(tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClientService.TestConnection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClientService_validate(t *testing.T) {
	type fields struct {
		ctx       context.Context
		apiClient *APIClient
	}
	type args struct {
		config string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *clientcmdapi.Config
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
			got, err := s.validate(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClientService.validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClientService.validate() = %v, want %v", got, tt.want)
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

func TestClientService_Analyze(t *testing.T) {
	type fields struct {
		ctx       context.Context
		apiClient *APIClient
	}
	type args struct {
		cluster string
		filters []string
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
			if got := s.Analyze(tt.args.cluster, tt.args.filters); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClientService.Analyze() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
