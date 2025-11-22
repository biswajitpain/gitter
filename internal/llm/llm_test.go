package llm_test

import (
	"encoding/json"
	"github.com/biswajitpain/gitter/internal/config"
	"github.com/biswajitpain/gitter/internal/llm"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewLLMClient(t *testing.T) {
	// Test case for OpenAI
	openAIConfig := config.Config{Provider: "openai", APIKey: "test-key"}
	client, err := llm.NewLLMClient(openAIConfig)
	if err != nil {
		t.Fatalf("NewLLMClient with openai config failed: %v", err)
	}
	if _, ok := client.(*llm.OpenAIClient); !ok {
		t.Errorf("NewLLMClient did not return an OpenAIClient for provider 'openai'")
	}

	// Test case for no provider
	noProviderConfig := config.Config{}
	_, err = llm.NewLLMClient(noProviderConfig)
	if err == nil {
		t.Errorf("NewLLMClient with no provider should have returned an error, but it didn't")
	}

	// Test case for unsupported provider
	unsupportedConfig := config.Config{Provider: "unsupported-llm"}
	_, err = llm.NewLLMClient(unsupportedConfig)
	if err == nil {
		t.Errorf("NewLLMClient with unsupported provider should have returned an error, but it didn't")
	}
}

func TestOpenAIClient_GenerateCommitMessage(t *testing.T) {
	// The expected response from the mock OpenAI API
	// This struct is not exported, so we redefine it for the test.
	type openAIResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	mockResponse := openAIResponse{
		Choices: []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		}{
			{
				Message: struct {
					Content string `json:"content"`
				}{
					Content: "feat: add new feature",
				},
			},
		},
	}
	mockRespBytes, _ := json.Marshal(mockResponse)

	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-key" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(mockRespBytes)
	}))
	defer server.Close()

	// Create an OpenAIClient and point it to the mock server
	client := &llm.OpenAIClient{
		APIKey:  "test-key",
		BaseURL: server.URL, // Use the mock server's URL
	}

	// Call the method we want to test
	generatedMessage, err := client.GenerateCommitMessage("diff", "user message")
	if err != nil {
		t.Fatalf("GenerateCommitMessage failed: %v", err)
	}

	// Check if the generated message is what we expect from the mock response
	expectedMessage := "feat: add new feature"
	if generatedMessage != expectedMessage {
		t.Errorf("Generated message is '%s', want '%s'", generatedMessage, expectedMessage)
	}

	// Test the error case where the API key is missing.
	clientNoKey := &llm.OpenAIClient{}
	_, err = clientNoKey.GenerateCommitMessage("diff", "user message")
	if err == nil {
		t.Error("Expected an error when API key is missing, but got nil")
	}
}

// Note: A more complete test for GenerateCommitMessage would require refactoring
// the function to accept an http.Client and a URL, allowing us to inject the
// test server's client and URL. The current implementation has a hardcoded URL,
// making direct testing of the HTTP request logic difficult without more extensive mocks.
