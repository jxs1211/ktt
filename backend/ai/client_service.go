package ai

import (
	"context"
	"errors"
	"ktt/backend/types"
	"ktt/backend/utils/log"

	"github.com/k8sgpt-ai/k8sgpt/pkg/ai"
)

var (
	aiClientMap = map[string]ai.IAI{
		"ollama": &ai.OllamaClient{},
		"llama3": &ai.OllamaClient{},
		"openai": &ai.OpenAIClient{},
		// "azureopenai": &ai.AzureAIClient{},
		// "localai":     &ai.LocalAIClient{},
		// "cohere":      &ai.CohereClient{},
		// "amazonbedrock": &ai.AmazonBedRockClient{},
		// "amazonsagemaker": &ai.AmazonSageMakerAIClient{},
		// "google": &ai.GoogleGenAIClient{},
		// "huggingface": &ai.HuggingfaceClient{},
		// "googlevertexai": &ai.GoogleVertexAIClient{},
		// "oci": &ai.OCIGenAIClient{},
		// "watsonxai": &ai.WatsonxAIClient{},
	}
)

type ClientService struct {
	m   map[string]ai.IAI
	ctx context.Context
}

func NewClientService() *ClientService {
	return &ClientService{
		m: make(map[string]ai.IAI, len(aiClientMap)),
	}
}

func (s *ClientService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *ClientService) GetAIProviders() types.JSResp {
	keys := make([]string, 0, len(aiClientMap))
	for k := range aiClientMap {
		keys = append(keys, k)
	}
	return types.JSResp{Success: true, Data: keys}
}

func (s *ClientService) Configure(provider ai.AIProvider) types.JSResp {
	err := s.configure(provider)
	if err != nil {
		return types.FailedResp(err.Error())
	}
	return types.JSResp{Success: true}
}

func (s *ClientService) configure(provider ai.AIProvider) error {
	log.Info("configure", "provider", provider)
	if _, ok := s.m[provider.Model]; ok {
		// client already exists: verified and cached
		return nil
	}
	c, ok := aiClientMap[provider.Model]
	if !ok {
		return errors.New("client not supported: " + provider.Model)
	}
	err := c.Configure(&provider)
	if err != nil {
		return err
	}
	s.m[provider.Model] = c
	return nil
}

func (s *ClientService) GetCompletion(key string, prompt string) types.JSResp {
	res, err := s.getCompletion(key, prompt)
	if err != nil {
		return types.FailedResp(err.Error())
	}
	return types.JSResp{Success: true, Data: res}
}

func (s *ClientService) getCompletion(key string, prompt string) (string, error) {
	c := s.m[key]
	if c == nil {
		return "", errors.New("client not found by " + key)
	}
	return c.GetCompletion(s.ctx, prompt)
}
