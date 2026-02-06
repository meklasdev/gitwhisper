package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const SystemPrompt = `You are a commit message generator.
Generate a Semantic Commit Message (Conventional Commits) based on the provided git diff.
Format: <type>(<scope>): <subject>
Types: feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert.
Rules:
1. Use the imperative mood.
2. Limit the subject line to 50 characters if possible.
3. Do NOT include any explanations, markdown code blocks, or conversation.
4. Output strict raw text of the commit message only.
`

type LLMEngine interface {
	GenerateCommitMessage(ctx context.Context, diff string) (string, error)
}

func NewEngine() (LLMEngine, error) {
	provider := viper.GetString("ai.provider")
	switch provider {
	case "openai":
		return NewOpenAIClient(
			viper.GetString("ai.openai_api_key"),
			viper.GetString("ai.openai_model"),
			"https://api.openai.com/v1/chat/completions",
		), nil
	case "openai-compatible":
		return NewOpenAIClient(
			viper.GetString("ai.openai_api_key"),
			viper.GetString("ai.openai_model"),
			viper.GetString("ai.openai_base_url"),
		), nil
	case "ollama":
		return NewOllamaClient(viper.GetString("ai.ollama_endpoint"), viper.GetString("ai.ollama_model")), nil
	case "gemini":
		return NewGeminiClient(viper.GetString("ai.gemini_api_key"), viper.GetString("ai.gemini_model")), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

// OpenAI / Compatible Client
type OpenAIClient struct {
	apiKey  string
	model   string
	baseURL string
}

func NewOpenAIClient(apiKey, model, baseURL string) *OpenAIClient {
	if model == "" {
		model = "gpt-3.5-turbo"
	}
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1/chat/completions"
	}
	return &OpenAIClient{apiKey: apiKey, model: model, baseURL: baseURL}
}

func (c *OpenAIClient) GenerateCommitMessage(ctx context.Context, diff string) (string, error) {
	reqBody := map[string]interface{}{
		"model": c.model,
		"messages": []map[string]string{
			{"role": "system", "content": SystemPrompt},
			{"role": "user", "content": diff},
		},
		"temperature": 0.7,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", errors.New("no response from API")
	}

	return strings.TrimSpace(result.Choices[0].Message.Content), nil
}

// Ollama Client
type OllamaClient struct {
	endpoint string
	model    string
}

func NewOllamaClient(endpoint, model string) *OllamaClient {
	if endpoint == "" {
		endpoint = "http://localhost:11434/api/generate"
	} else if !strings.HasSuffix(endpoint, "/api/generate") {
		endpoint = strings.TrimRight(endpoint, "/") + "/api/generate"
	}
	if model == "" {
		model = "mistral"
	}
	return &OllamaClient{endpoint: endpoint, model: model}
}

func (c *OllamaClient) GenerateCommitMessage(ctx context.Context, diff string) (string, error) {
	reqBody := map[string]interface{}{
		"model":  c.model,
		"prompt": diff,
		"system": SystemPrompt,
		"stream": false,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama API error: %s", string(body))
	}

	var result struct {
		Response string `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return strings.TrimSpace(result.Response), nil
}

// Gemini Client
type GeminiClient struct {
	apiKey string
	model  string
}

func NewGeminiClient(apiKey, model string) *GeminiClient {
	if model == "" {
		model = "gemini-pro"
	}
	return &GeminiClient{apiKey: apiKey, model: model}
}

func (c *GeminiClient) GenerateCommitMessage(ctx context.Context, diff string) (string, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", c.model, c.apiKey)

	// Gemini API expects a specific structure
	reqBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{"text": SystemPrompt + "\n\nDiff:\n" + diff},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature": 0.7,
		},
	}

	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Gemini API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", errors.New("no response content from Gemini")
	}

	return strings.TrimSpace(result.Candidates[0].Content.Parts[0].Text), nil
}
