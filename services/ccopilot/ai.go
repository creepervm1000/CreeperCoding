package ccopilot

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"creepercoding.dev/modules/httplib"
	"creepercoding.dev/modules/log"
	"creepercoding.dev/modules/setting"
)

type chatMessage struct {
	Role       string      `json:"role"`
	Content    string      `json:"content"`
	ToolCalls  []toolCall  `json:"tool_calls,omitempty"`
	ToolCallID string      `json:"tool_call_id,omitempty"`
}

type chatRequest struct {
	Model       string           `json:"model"`
	Messages    []chatMessage    `json:"messages"`
	MaxTokens   int64            `json:"max_tokens"`
	Temperature float64          `json:"temperature"`
	Stream      bool             `json:"stream"`
	Tools       []toolDefinition `json:"tools,omitempty"`
	ToolChoice  any              `json:"tool_choice,omitempty"`
}

type chatResponseChoice struct {
	Index        int         `json:"index"`
	Message      chatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

type chatResponse struct {
	Choices []chatResponseChoice `json:"choices"`
	Error   *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func queryAI(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	return queryAIMessages(ctx, []chatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}, nil)
}

func queryAIMessages(ctx context.Context, messages []chatMessage, tools []toolDefinition) (string, error) {
	endpoint := strings.TrimRight(setting.Config().Ccopilot.Endpoint.Value(ctx), "/")
	apiKey := setting.Config().Ccopilot.APIKey.Value(ctx)
	modelName := setting.Config().Ccopilot.ModelName.Value(ctx)
	maxTokens := setting.Config().Ccopilot.MaxTokens.Value(ctx)

	chatURL := endpoint + "/chat/completions"

	req := chatRequest{
		Model:       modelName,
		Messages:    messages,
		MaxTokens:   maxTokens,
		Temperature: 0.7,
		Stream:      false,
	}
	if len(tools) > 0 {
		req.Tools = tools
		req.ToolChoice = "auto"
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq := httplib.NewRequest(chatURL, http.MethodPost)
	httpReq.SetContext(ctx)
	httpReq.Header("Authorization", "Bearer "+apiKey)
	httpReq.Header("Content-Type", "application/json")
	httpReq.Body(body)

	resp, err := httpReq.Response()
	if err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var chatResp chatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if chatResp.Error != nil {
		return "", fmt.Errorf("API error: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("API returned no choices")
	}

	// Serialize the whole choice back to JSON so caller can inspect tool_calls
	resultBytes, err := json.Marshal(chatResp.Choices[0])
	if err != nil {
		return "", fmt.Errorf("failed to marshal choice: %w", err)
	}

	log.Debug("ccopilot: AI response received, finish_reason=%s", chatResp.Choices[0].FinishReason)
	return string(resultBytes), nil
}
