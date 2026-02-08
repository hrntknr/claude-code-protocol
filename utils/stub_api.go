package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
)

// SSEEvent represents a single SSE event to send to the client.
type SSEEvent struct {
	Event string
	Data  map[string]any
}

// StubAPIServer is a minimal Anthropic Messages API stub.
// Responses is a sequence of SSE event lists, one per API request.
// The first request gets Responses[0], the second gets Responses[1], etc.
// Extra requests beyond the list repeat the last response.
type StubAPIServer struct {
	Responses [][]SSEEvent

	server   *httptest.Server
	mu       sync.Mutex
	reqCount int
}

// Start creates and starts the stub HTTP server.
func (s *StubAPIServer) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/messages", s.handleMessages)
	s.server = httptest.NewServer(mux)
}

// Close shuts down the stub server.
func (s *StubAPIServer) Close() {
	if s.server != nil {
		s.server.Close()
	}
}

// URL returns the base URL of the stub server.
func (s *StubAPIServer) URL() string {
	return s.server.URL
}

func (s *StubAPIServer) handleMessages(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	s.mu.Lock()
	idx := s.reqCount
	if idx >= len(s.Responses) {
		idx = len(s.Responses) - 1
	}
	s.reqCount++
	s.mu.Unlock()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for _, e := range s.Responses[idx] {
		data, err := json.Marshal(e.Data)
		if err != nil {
			return
		}
		fmt.Fprintf(w, "event: %s\ndata: %s\n\n", e.Event, data)
		flusher.Flush()
	}
}

// ToolUseResponse builds an SSE event sequence where the assistant calls a tool.
// The stop_reason is "tool_use" so the CLI will execute the tool and make a follow-up request.
func ToolUseResponse(toolID, toolName string, input map[string]any) []SSEEvent {
	inputJSON, _ := json.Marshal(input)
	return []SSEEvent{
		{
			Event: "message_start",
			Data: map[string]any{
				"type": "message_start",
				"message": map[string]any{
					"id":            "msg_stub_tool",
					"type":          "message",
					"role":          "assistant",
					"content":       []any{},
					"model":         "claude-sonnet-4-5-20250929",
					"stop_reason":   nil,
					"stop_sequence": nil,
					"usage":         map[string]any{"input_tokens": 10, "output_tokens": 1},
				},
			},
		},
		{
			Event: "content_block_start",
			Data: map[string]any{
				"type":  "content_block_start",
				"index": 0,
				"content_block": map[string]any{
					"type":  "tool_use",
					"id":    toolID,
					"name":  toolName,
					"input": map[string]any{},
				},
			},
		},
		{
			Event: "content_block_delta",
			Data: map[string]any{
				"type":  "content_block_delta",
				"index": 0,
				"delta": map[string]any{
					"type":         "input_json_delta",
					"partial_json": string(inputJSON),
				},
			},
		},
		{
			Event: "content_block_stop",
			Data: map[string]any{
				"type":  "content_block_stop",
				"index": 0,
			},
		},
		{
			Event: "message_delta",
			Data: map[string]any{
				"type":  "message_delta",
				"delta": map[string]any{"stop_reason": "tool_use", "stop_sequence": nil},
				"usage": map[string]any{"output_tokens": 20},
			},
		},
		{
			Event: "message_stop",
			Data: map[string]any{
				"type": "message_stop",
			},
		},
	}
}

// TextResponse builds a standard SSE event sequence for a simple text response.
func TextResponse(text string) []SSEEvent {
	return []SSEEvent{
		{
			Event: "message_start",
			Data: map[string]any{
				"type": "message_start",
				"message": map[string]any{
					"id":            "msg_stub_001",
					"type":          "message",
					"role":          "assistant",
					"content":       []any{},
					"model":         "claude-sonnet-4-5-20250929",
					"stop_reason":   nil,
					"stop_sequence": nil,
					"usage":         map[string]any{"input_tokens": 10, "output_tokens": 1},
				},
			},
		},
		{
			Event: "content_block_start",
			Data: map[string]any{
				"type":          "content_block_start",
				"index":         0,
				"content_block": map[string]any{"type": "text", "text": ""},
			},
		},
		{
			Event: "content_block_delta",
			Data: map[string]any{
				"type":  "content_block_delta",
				"index": 0,
				"delta": map[string]any{"type": "text_delta", "text": text},
			},
		},
		{
			Event: "content_block_stop",
			Data: map[string]any{
				"type":  "content_block_stop",
				"index": 0,
			},
		},
		{
			Event: "message_delta",
			Data: map[string]any{
				"type":  "message_delta",
				"delta": map[string]any{"stop_reason": "end_turn", "stop_sequence": nil},
				"usage": map[string]any{"output_tokens": 5},
			},
		},
		{
			Event: "message_stop",
			Data: map[string]any{
				"type": "message_stop",
			},
		},
	}
}
