package ai

import (
	"context"
	"fmt"
	"ktt/backend/types"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/k8sgpt-ai/k8sgpt/pkg/ai"
	"github.com/ollama/ollama/api"
)

func TestClientService_Configure(t *testing.T) {
	tests := []struct {
		name     string
		m        map[string]ai.IAI
		ctx      context.Context
		provider ai.AIProvider
		want     types.JSResp
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ClientService{
				m:   tt.m,
				ctx: tt.ctx,
			}
			if got := s.Configure(tt.provider); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClientService.Configure() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClientService_configure(t *testing.T) {
	ollamaClient := &ai.OllamaClient{}
	if err := ollamaClient.Configure(&ai.AIProvider{
		Name: "ollama", Model: "llama3", BaseURL: "http://localhost:11434",
	}); err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name     string
		m        map[string]ai.IAI
		ctx      context.Context
		provider ai.AIProvider
		wantErr  bool
	}{
		{
			name: "client not supported",
			m:    map[string]ai.IAI{},
			ctx:  context.Background(),
			provider: ai.AIProvider{
				Model: "not-exist",
			},
			wantErr: true,
		},
		{
			name: "base",
			m:    map[string]ai.IAI{},
			ctx:  context.Background(),
			provider: ai.AIProvider{
				Name:    "localai",
				Model:   "llama3",
				BaseURL: "http://localhost:11434",
			},
		},
		{
			name: "client already exists in cache",
			m: map[string]ai.IAI{
				"ollama": ollamaClient,
			},
			ctx: context.Background(),
			provider: ai.AIProvider{
				Name:    "localai",
				Model:   "llama3",
				BaseURL: "http://localhost:11434",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ClientService{
				m:   tt.m,
				ctx: tt.ctx,
			}
			if err := s.configure(tt.provider); (err != nil) != tt.wantErr {
				t.Errorf("ClientService.configure() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClientService_GetCompletion(t *testing.T) {
	tests := []struct {
		name     string
		provider ai.AIProvider
		ctx      context.Context
		key      string
		prompt   string
		want     types.JSResp
	}{
		{
			name: "base",
			provider: ai.AIProvider{
				Name:    "localai",
				Model:   "llama3",
				BaseURL: "http://localhost:11434",
			},
			ctx: context.Background(),
			key: "llama3",
			// prompt: wrongNginxErrMsg,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewClientService()
			s.SetContext(tt.ctx)
			err := s.configure(tt.provider)
			if err != nil {
				t.Fatal(err)
			}
			got := s.GetCompletion(tt.key, tt.prompt)
			if !got.Success {
				t.Fatal(got.Msg)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClientService.GetCompletion() = %v, want %v", got, tt.want)
			}
			t.Log(got)
		})
	}
}

func TestClientService_getCompletion(t *testing.T) {
	t.Log("in")
	tests := []struct {
		name     string
		provider ai.AIProvider
		ctx      context.Context
		key      string
		prompt   string
		want     string
		wantErr  error
	}{
		{
			name: "base",
			provider: ai.AIProvider{
				Name:    "localai",
				Model:   "llama3.2",
				BaseURL: "http://localhost:11434",
			},
			ctx:    context.Background(),
			key:    "llama3",
			prompt: wrongNginxErrMsg,
			want:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewClientService()
			s.SetContext(tt.ctx)
			err := s.configure(tt.provider)
			if err != nil {
				t.Fatal(err)
			}
			got, err := s.getCompletion(tt.key, tt.prompt)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClientService.GetCompletion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOllamaClient(t *testing.T) {
	// Create a new Ollama client
	baseURL := "http://localhost:11434"
	baseClientURL, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
	}

	httpClient := http.DefaultClient

	c := api.NewClient(baseClientURL, httpClient)

	// Set up the generation request
	req := &api.GenerateRequest{
		Model:  "llama3.2",
		Prompt: "Why is the sky blue?",
		Stream: new(bool),
		Options: map[string]interface{}{
			"temperature": float32(0),
			"top_p":       float32(0),
		},
	}

	completion := ""
	respFunc := func(resp api.GenerateResponse) error {
		completion = resp.Response
		return nil
	}
	err = c.Generate(context.Background(), req, respFunc)
	// Make the generate request
	if err != nil {
		log.Fatalf("Error generating response: %v", err)
	}

	// Print the generated response
	fmt.Println("Generated response:")
	fmt.Println(completion)
}
