package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
)

// SSEEvent represents a single SSE event to send to the client.
type SSEEvent struct {
	Event string
	Data  map[string]any
}

// ToolCall describes a single tool invocation for use in MultiToolUseResponse.
type ToolCall struct {
	ID    string
	Name  string
	Input map[string]any
}

// RecordedRequest stores the decoded body of each API request.
type RecordedRequest struct {
	Body map[string]any
}

// StubAPIServer is a minimal Anthropic Messages API stub.
// Responses is a sequence of SSE event lists, one per API request.
// The first request gets Responses[0], the second gets Responses[1], etc.
// Extra requests beyond the list repeat the last response.
//
// StaticPages maps URL paths to static HTML content, served as GET requests.
// This allows testing tools like WebFetch that need to fetch external URLs.
type StubAPIServer struct {
	Responses   [][]SSEEvent
	StaticPages map[string]string

	server   *httptest.Server
	mu       sync.Mutex
	reqCount int
	requests []RecordedRequest
}

// Start creates and starts the stub HTTP server.
func (s *StubAPIServer) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/messages", s.handleMessages)
	mux.HandleFunc("GET /static/", s.handleStatic)
	s.server = httptest.NewServer(mux)
}

func (s *StubAPIServer) handleStatic(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if s.StaticPages != nil {
		if content, ok := s.StaticPages[path]; ok {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, content)
			return
		}
	}
	http.NotFound(w, r)
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

// Requests returns a copy of all recorded requests.
func (s *StubAPIServer) Requests() []RecordedRequest {
	s.mu.Lock()
	defer s.mu.Unlock()
	cp := make([]RecordedRequest, len(s.requests))
	copy(cp, s.requests)
	return cp
}

// RequestCount returns the number of requests received so far.
func (s *StubAPIServer) RequestCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.reqCount
}

func (s *StubAPIServer) handleMessages(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	// Record request body before consuming it.
	bodyBytes, _ := io.ReadAll(r.Body)
	var body map[string]any
	if err := json.Unmarshal(bodyBytes, &body); err == nil {
		s.mu.Lock()
		s.requests = append(s.requests, RecordedRequest{Body: body})
		s.mu.Unlock()
	}

	// Auto-detect haiku init requests by checking the model field.
	// The CLI makes internal haiku requests (quota check, file-change detection)
	// before the main request. These are answered with a dummy "ok" response
	// without consuming from the Responses queue.
	if model, _ := body["model"].(string); strings.Contains(model, "haiku") {
		s.writeSSE(w, flusher, TextResponse("ok"))
		return
	}

	s.mu.Lock()
	idx := s.reqCount
	if idx >= len(s.Responses) {
		idx = len(s.Responses) - 1
	}
	s.reqCount++
	s.mu.Unlock()

	s.writeSSE(w, flusher, s.Responses[idx])
}

func (s *StubAPIServer) writeSSE(w http.ResponseWriter, flusher http.Flusher, events []SSEEvent) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for _, e := range events {
		data, err := json.Marshal(e.Data)
		if err != nil {
			return
		}
		fmt.Fprintf(w, "event: %s\ndata: %s\n\n", e.Event, data)
		flusher.Flush()
	}
}

// ---------------------------------------------------------------------------
// Response helpers
// ---------------------------------------------------------------------------

// messageStartEvent returns the common message_start SSE event.
func messageStartEvent() SSEEvent {
	return SSEEvent{
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
	}
}

// messageDeltaEvent returns a message_delta SSE event with the given stop_reason.
func messageDeltaEvent(stopReason string) SSEEvent {
	return SSEEvent{
		Event: "message_delta",
		Data: map[string]any{
			"type":  "message_delta",
			"delta": map[string]any{"stop_reason": stopReason, "stop_sequence": nil},
			"usage": map[string]any{"output_tokens": 20},
		},
	}
}

// messageStopEvent returns the common message_stop SSE event.
func messageStopEvent() SSEEvent {
	return SSEEvent{
		Event: "message_stop",
		Data:  map[string]any{"type": "message_stop"},
	}
}

// textBlockEvents returns content_block_start/delta/stop for a text block at the given index.
func textBlockEvents(index int, text string) []SSEEvent {
	return []SSEEvent{
		{
			Event: "content_block_start",
			Data: map[string]any{
				"type":          "content_block_start",
				"index":         index,
				"content_block": map[string]any{"type": "text", "text": ""},
			},
		},
		{
			Event: "content_block_delta",
			Data: map[string]any{
				"type":  "content_block_delta",
				"index": index,
				"delta": map[string]any{"type": "text_delta", "text": text},
			},
		},
		{
			Event: "content_block_stop",
			Data: map[string]any{
				"type":  "content_block_stop",
				"index": index,
			},
		},
	}
}

// toolUseBlockEvents returns content_block_start/delta/stop for a tool_use block at the given index.
func toolUseBlockEvents(index int, toolID, toolName string, input map[string]any) []SSEEvent {
	inputJSON, _ := json.Marshal(input)
	return []SSEEvent{
		{
			Event: "content_block_start",
			Data: map[string]any{
				"type":  "content_block_start",
				"index": index,
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
				"index": index,
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
				"index": index,
			},
		},
	}
}

// thinkingBlockEvents returns content_block_start/delta/stop for a thinking block at the given index.
func thinkingBlockEvents(index int, thinking string) []SSEEvent {
	return []SSEEvent{
		{
			Event: "content_block_start",
			Data: map[string]any{
				"type":  "content_block_start",
				"index": index,
				"content_block": map[string]any{
					"type":     "thinking",
					"thinking": "",
				},
			},
		},
		{
			Event: "content_block_delta",
			Data: map[string]any{
				"type":  "content_block_delta",
				"index": index,
				"delta": map[string]any{
					"type":     "thinking_delta",
					"thinking": thinking,
				},
			},
		},
		{
			Event: "content_block_stop",
			Data: map[string]any{
				"type":  "content_block_stop",
				"index": index,
			},
		},
	}
}

// ---------------------------------------------------------------------------
// Public response builders
// ---------------------------------------------------------------------------

// ToolUseResponse builds an SSE event sequence where the assistant calls a tool.
// The stop_reason is "tool_use" so the CLI will execute the tool and make a follow-up request.
func ToolUseResponse(toolID, toolName string, input map[string]any) []SSEEvent {
	events := []SSEEvent{messageStartEvent()}
	events = append(events, toolUseBlockEvents(0, toolID, toolName, input)...)
	events = append(events, messageDeltaEvent("tool_use"), messageStopEvent())
	return events
}

// TextResponse builds a standard SSE event sequence for a simple text response.
func TextResponse(text string) []SSEEvent {
	events := []SSEEvent{messageStartEvent()}
	events = append(events, textBlockEvents(0, text)...)
	events = append(events, messageDeltaEvent("end_turn"), messageStopEvent())
	return events
}

// TextAndToolUseResponse builds an SSE event sequence with a text block (index 0)
// followed by a tool_use block (index 1) in the same message.
func TextAndToolUseResponse(text, toolID, toolName string, input map[string]any) []SSEEvent {
	events := []SSEEvent{messageStartEvent()}
	events = append(events, textBlockEvents(0, text)...)
	events = append(events, toolUseBlockEvents(1, toolID, toolName, input)...)
	events = append(events, messageDeltaEvent("tool_use"), messageStopEvent())
	return events
}

// MultiToolUseResponse builds an SSE event sequence with multiple tool_use blocks.
func MultiToolUseResponse(tools ...ToolCall) []SSEEvent {
	events := []SSEEvent{messageStartEvent()}
	for i, tc := range tools {
		events = append(events, toolUseBlockEvents(i, tc.ID, tc.Name, tc.Input)...)
	}
	events = append(events, messageDeltaEvent("tool_use"), messageStopEvent())
	return events
}

// ThinkingResponse builds an SSE event sequence with a thinking block (index 0)
// followed by a text block (index 1).
func ThinkingResponse(thinking, text string) []SSEEvent {
	events := []SSEEvent{messageStartEvent()}
	events = append(events, thinkingBlockEvents(0, thinking)...)
	events = append(events, textBlockEvents(1, text)...)
	events = append(events, messageDeltaEvent("end_turn"), messageStopEvent())
	return events
}

// ThinkingAndToolUseResponse builds an SSE event sequence with a thinking block (index 0)
// followed by a tool_use block (index 1).
func ThinkingAndToolUseResponse(thinking, toolID, toolName string, input map[string]any) []SSEEvent {
	events := []SSEEvent{messageStartEvent()}
	events = append(events, thinkingBlockEvents(0, thinking)...)
	events = append(events, toolUseBlockEvents(1, toolID, toolName, input)...)
	events = append(events, messageDeltaEvent("tool_use"), messageStopEvent())
	return events
}

// MaxTokensTextResponse builds an SSE event sequence for a text response
// that was truncated by the max_tokens limit (stop_reason: "max_tokens").
func MaxTokensTextResponse(text string) []SSEEvent {
	events := []SSEEvent{messageStartEvent()}
	events = append(events, textBlockEvents(0, text)...)
	events = append(events, messageDeltaEvent("max_tokens"), messageStopEvent())
	return events
}

// MultiTextResponse builds an SSE event sequence with multiple text blocks
// in a single response message.
func MultiTextResponse(texts ...string) []SSEEvent {
	events := []SSEEvent{messageStartEvent()}
	for i, text := range texts {
		events = append(events, textBlockEvents(i, text)...)
	}
	events = append(events, messageDeltaEvent("end_turn"), messageStopEvent())
	return events
}

// StopSequenceTextResponse builds an SSE event sequence for a text response
// that was stopped by a stop_sequence (stop_reason: "stop_sequence").
func StopSequenceTextResponse(text, stopSequence string) []SSEEvent {
	events := []SSEEvent{messageStartEvent()}
	events = append(events, textBlockEvents(0, text)...)
	events = append(events, SSEEvent{
		Event: "message_delta",
		Data: map[string]any{
			"type":  "message_delta",
			"delta": map[string]any{"stop_reason": "stop_sequence", "stop_sequence": stopSequence},
			"usage": map[string]any{"output_tokens": 20},
		},
	}, messageStopEvent())
	return events
}

// ErrorSSEResponse builds a single SSE error event (API-level error, not tool error).
func ErrorSSEResponse(errorType, message string) []SSEEvent {
	return []SSEEvent{
		{
			Event: "error",
			Data: map[string]any{
				"type": "error",
				"error": map[string]any{
					"type":    errorType,
					"message": message,
				},
			},
		},
	}
}
