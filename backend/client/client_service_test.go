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

// func TestClientService_Analyze(t *testing.T) {
// 	type fields struct {
// 		ctx       context.Context
// 		apiClient *APIClient
// 	}
// 	type args struct {
// 		cluster string
// 		filters []string
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   types.JSResp
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := &ClientService{
// 				ctx:       tt.fields.ctx,
// 				apiClient: tt.fields.apiClient,
// 			}
// 			if got := s.Analyze(tt.args.cluster, tt.args.filters); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("ClientService.Analyze() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// const kubeconfig = `
// apiVersion: v1
// clusters:
// - cluster:
//     certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURCVENDQWUyZ0F3SUJBZ0lJSmJFeVlmTEEzUWd3RFFZSktvWklodmNOQVFFTEJRQXdGVEVUTUJFR0ExVUUKQXhNS2EzVmlaWEp1WlhSbGN6QWVGdzB5TkRBM01qRXhNalUwTXpGYUZ3MHpOREEzTVRreE1qVTVNekZhTUJVeApFekFSQmdOVkJBTVRDbXQxWW1WeWJtVjBaWE13Z2dFaU1BMEdDU3FHU0liM0RRRUJBUVVBQTRJQkR3QXdnZ0VLCkFvSUJBUURSWTZ1VTg5Rkl5VFMrUnU3QTZnME9qbzBHcnpTZ3FvQnNSY3VBNGNKVkpHdzBIdnR4UEZLSFFLczkKRXdtNjViOG1GbEtoeUlYclF4cEFPZm9HVnlEYmU4ZGpQVFF4TFBaYjhYTnA5dEIrT2tpQjltUVFpZ2FrVURObwo3aDZXTjdVS1dlYWtyQm1VbHk2L2UraWI2UXY2VjY4clJaZmhCRDhoQlBOYjNSZ25PK2xGOGhpaW1BdWdpWFhzCjh3NDQ2YWswMmJvZTV3b1NKVDZEOFpIZmJRbGtNSjB5dVNjRG1jNXFqVHpNQW5Ga0dkS0R6Um9vUk55MXozQloKZnJUMXpmc3kxNS9wd1d5dHRPNmpVYTFDNDgyL3U5cGRYNnA5UHZDNjB0aXZCNXJIVjJsdXAxU09TYW1YV3RNego5M1RSQ1J5VVN5SjNySEFPUE9YLzlqeHc0WWROQWdNQkFBR2pXVEJYTUE0R0ExVWREd0VCL3dRRUF3SUNwREFQCkJnTlZIUk1CQWY4RUJUQURBUUgvTUIwR0ExVWREZ1FXQkJSWUluR3VXdjlsVno3VHp0UWZSQkhBekRVdEVqQVYKQmdOVkhSRUVEakFNZ2dwcmRXSmxjbTVsZEdWek1BMEdDU3FHU0liM0RRRUJDd1VBQTRJQkFRRER1aENvdnEwZQprYzdEM2k4Z0YyaERQNDArTjN4YmVDN0FsN1RrSU52N281YmM5eHlKeDhuTUZHNUdvT0hURjlwTzN6cFFyN3VPCitGbGZQckN3bCt6a1d0VEJRb2RKT1IwRUl0QkUvcVpXT0o5aHRwNjBQMUwwemREZVVVK0F6ZW5hMS9Sd0VOSEkKNkVQdU91ckpjNXdtSTBzNVZneXVxbjBBTUtBaFJ2R3hlSGx5TFN0N29zb1ZienBIY3laQ21uMUtRdW51emE3eQpUVE40Wi95K2szM0xYUVc2elJqbVZDeW5MRElZNzZYRjNqR2ZzNmxlVFJKTlZqSlN3a05ZL1Z0aUxqVUdZTlE4ClpaL2dMem1uVHMxcU8yalhFQy91ZGJCdDBBY0VXamgzQVZBNmxzdGFUQXI4ZHNXRmFQNWNZS3ZFckV0cVVQb1oKSjR4Q3dVZ0FEckNKCi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K
//     server: https://127.0.0.1:58953
//   name: kind-test
// contexts:
// - context:
//     cluster: kind-test
//     user: kind-test
//   name: kind-test
// current-context: kind-test
// kind: Config
// preferences: {}
// users:
// - name: kind-test
//   user:
//     client-certificate-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURLVENDQWhHZ0F3SUJBZ0lJQjhBRFYwaytldnN3RFFZSktvWklodmNOQVFFTEJRQXdGVEVUTUJFR0ExVUUKQXhNS2EzVmlaWEp1WlhSbGN6QWVGdzB5TkRBM01qRXhNalUwTXpGYUZ3MHlOVEEzTWpFeE1qVTVNemxhTUR3eApIekFkQmdOVkJBb1RGbXQxWW1WaFpHMDZZMngxYzNSbGNpMWhaRzFwYm5NeEdUQVhCZ05WQkFNVEVHdDFZbVZ5CmJtVjBaWE10WVdSdGFXNHdnZ0VpTUEwR0NTcUdTSWIzRFFFQkFRVUFBNElCRHdBd2dnRUtBb0lCQVFEekE4L2YKcHJRaEtWYzBkWnRZZ0ZDSDZGNTdod1M4MlpUekx1QTNvcXJLcUsrQWo1OWRBVWdMZCtQdUVyQnlpaWh5cjR3KwoxQUdSSnczT2xESnNET0FWVHB0aFgxeFQwQldjUEZhUVkxc283Vys3eGdSMmhFZHNrNGJGOVA2dFlIWDgvRkxUCnVWQnRKSHpxQWRkR1NGRkYwWC9hSWNDM1V1bklHckFHcjRVNVh6UHozL1pMQkY1VUJQbXBqN0RwVm5wNXZsT3IKMkI1TXY4c3FVQmlsaXhGbFJ6amRwYlBpWnhpRG1PcWNYR0NoQ3BCSnp5U2ZLMEpBaGhVN0hKajQ2SGJPeUNlNApZaG9kUXRaczZPaVBHYVU2TTJWWTl4akVqRC9ZQlpITlJqdUdBalBFakdrRTByb2JrTHdTSlBCNWRudVczbWh0Ckg4YlAvdjV5Y3pVVjltU2RBZ01CQUFHalZqQlVNQTRHQTFVZER3RUIvd1FFQXdJRm9EQVRCZ05WSFNVRUREQUsKQmdnckJnRUZCUWNEQWpBTUJnTlZIUk1CQWY4RUFqQUFNQjhHQTFVZEl3UVlNQmFBRkZnaWNhNWEvMlZYUHRQTwoxQjlFRWNETU5TMFNNQTBHQ1NxR1NJYjNEUUVCQ3dVQUE0SUJBUUMvM1B4VG1uZ2Q2Sjg2bGN0Q08rRTRtVE1xCndKdFMxRFVwYmdCaXV6dDdjQVE4ZHRNZ3Q4cFIzNTUyc0diUkJpZmRVckw2aXh1eVc1aXFLZzhkV3F5UzF1ZEwKck9ieWFaUUJGaXF1YkIvc3l4cUY4dkNOYmp2blFxQnJVSTI0S25tQ1JJQWtvbU5EaTBvVW5HeFJMc2hYZHFGQQpSU0dBby9Nc1JJLzZKWVQydFJ1cUU2ellCV0g1TVYxN3dhd1hmU0hXbGpFMUd6cjRib2l6cm5GMVA0TkJpRTRLCnQxckNKbjM0TWxneG1RM3pqRDdHV0dZTlFaWnpoM3dzUVZaK01nSk1tV2tkRCtJRERROGpTOXVQYkhQcnpZTFUKSmd5M1hmeEN2dE5MOVZoR3l3dnJBOHNVSGtROUhjeG5JNm41QTlyTXc1dnJiZGE5eHFqbzlPSXQ0ZzhFCi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K
//     client-key-data: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcEFJQkFBS0NBUUVBOHdQUDM2YTBJU2xYTkhXYldJQlFoK2hlZTRjRXZObVU4eTdnTjZLcXlxaXZnSStmClhRRklDM2ZqN2hLd2Nvb29jcStNUHRRQmtTY056cFF5YkF6Z0ZVNmJZVjljVTlBVm5EeFdrR05iS08xdnU4WUUKZG9SSGJKT0d4ZlQrcldCMS9QeFMwN2xRYlNSODZnSFhSa2hSUmRGLzJpSEF0MUxweUJxd0JxK0ZPVjh6ODkvMgpTd1JlVkFUNXFZK3c2Vlo2ZWI1VHE5Z2VUTC9MS2xBWXBZc1JaVWM0M2FXejRtY1lnNWpxbkZ4Z29RcVFTYzhrCm55dENRSVlWT3h5WStPaDJ6c2dudUdJYUhVTFdiT2pvanhtbE9qTmxXUGNZeEl3LzJBV1J6VVk3aGdJenhJeHAKQk5LNkc1QzhFaVR3ZVhaN2x0NW9iUi9Hei83K2NuTTFGZlprblFJREFRQUJBb0lCQVFDTDhoMEswZWlYMUVQWQovUVZKOFJMdlFWenZ6REJwUk8vbkg0NlYyNEo1bEswRTN6REtFWXZZdHVFMjMrSm5BN01KWUpqbU1aYjViVUVoClB1ek9scGVSRGFTamJaUXcyL3NsN3dWMnZ4RG5QOTBCaUtaWFRoUUhyZW1HSkhGcGpNeEZ0VlZKZ2tXVHBOaWUKdGhLdjhjT0Q3b0t2THBMWnY5RExvTHJRRHRJczE4d3FpbHAvcGlpY0FIVlE5cFJTUmRnM3dnL091VG15dzZ1NwpVNnp5SnMvWFhyeDlXT2cvZWZ3bHpaTXBUWXV3UDFNc2dJTnlQUHFYV2tTVFk0M0ZlK2hwa0ZmYlowRCtQQ25wCkgweWNZWk9ITWUveWt0T1FBa1FLUkVUekVPUFdYUWdXNEJWdHd3a3ZxbmM2b2RBdlI3VDF5b1NJU1pibTVhUHUKWVhsalBXd0JBb0dCQVA5M0xnOFV2RkxmU0JvNU9Jc0QwekREcWoyOEsweUdRMDFYNEg4dSsrM3VJWHk0c2VHbAp1OXhjUURuL0xkZDBZbTlGRjZHeW1QeS9EMmxaYkdISFEreTRnbnJvYThMQ1VsUzVVWSt1TjZTZFV1NHRMWHR2Cm9tS2ZhcmEvekZ1TjBTN1pYOVBsWGhGSEFyRmhKUytseEEvcmNCeTlGM1ZIa1d4cC9FNVorQWVkQW9HQkFQT0YKOXNCSXNXbEpBY0MxYWNXOW5SbmFRZXBDNytUMitCM092LzJZbGExdWtkWjJUSlNjTXpzN2podGFjRFdwb29KKwpERDA1RVdobkFySGt5MHJrUnZRN2JMcW9pQkRzV3AwQVg4bnczb2V6QTN4TlorLzRrNk1xL3ZjeGZzc1RidVhyCnJvUmJMWEhBOWY0NjhIamlrUEhKdWtPT1dNL2JpMHVIbHcrNTljRUJBb0dBSkdoYnpCSnNkSE1WMTRib0pBZTcKaFd1ZFM0Y1J0S285MVR2ODVxTlBqQzB2NEpLQjUyS3pUMGhtYnQyTEx2V0xRY2hiQWRSdU1UY1pmeTRiWUNRQgo3aTZ0aWM2dDdPZTh6QTVOTFdqcXpTOE9ycHNKckZuUWpyV2hnOEg1NGVKb09ZRUpReTJoSmwwMGFRc3JQWXNtClVnNS92OWpEQ3hmVUJkdGs5ZFdrbjVFQ2dZRUFsb0ZkRjBuV0cvUDdHVmNGb0Npb1I2b1V3dXZMeC82N0tmRDEKeDQzZU0rbjZTRW1rMnRRTzliVEJCemJGMHVTY2czblRwcUwybDBmUzZvODA5WHhRUGZIY0tTQ3Z0NFhjR3R4NApWeDUrNmU5QllEbVcyMEVPUGZIODBsbk11MFd2YkhwVXlZaHdkYldFVXpPcUc0d3JlTVBzVW9SUGIyZUlsNDJkClZ5TGZPQUVDZ1lCenM1TTdkN2o5SWNtaFhBSTNOSXFRbmxZN3dMUTI1Q0VNaFZjOFlkOVBuWWNlVFFWV3lPazYKSVlqdFFDTmgranR2RnFySHR6K1VHbUhsN3RhVXAxR1hKd0k4ek4zRy9sbGJ5MVVIMkE2N3JVdDVUc25lUVJBUwpzTXZnNkhJUEkrNVRBZFdWYnBOWUNPbHp2NDlFVDZNOUJtTzZDYzRLWnBudEYxZ2pzK1Mwenc9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=
// `

// func TestClientService_getApiResources(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		ctx     context.Context
// 		config  string
// 		want    []string
// 		wantErr bool
// 	}{
// 		{
// 			name:   "base",
// 			ctx:    context.Background(),
// 			config: kubeconfig,
// 			want:   []string{"Pod", "Service"},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := NewClientService()
// 			s.Start(tt.ctx)
// 			s.LoadConfig(tt.config)
// 			got, err := s.getApiResources()
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("ClientService.getApiResources() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("ClientService.getApiResources() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

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
