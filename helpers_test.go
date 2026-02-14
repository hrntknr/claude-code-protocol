package ccprotocol_test

import (
	. "github.com/hrntknr/claudecodeprotocol"
	"github.com/hrntknr/claudecodeprotocol/utils"
)

// defaultInitPattern returns a SystemInitMessage JSON assertion pattern
func defaultInitPattern(opts ...func(*SystemInitMessage)) utils.Pattern {
	m := SystemInitMessage{
		MessageBase:       MessageBase{Type: TypeSystem, Subtype: SubtypeInit},
		CWD:               "/home/user/project",
		SessionID:         "session-abc123",
		Tools:             []string{"Bash", "Read", "Write", "Edit", "Glob", "Grep"},
		MCPServers:        []string{},
		Model:             "claude-sonnet-4-5-20250929",
		PermissionMode:    PermissionBypassPermissions,
		SlashCommands:     []string{},
		APIKeySource:      "env_variable",
		ClaudeCodeVersion: "2.1.0",
		OutputStyle:       "default",
		Agents:            []string{},
		Skills:            []string{},
		Plugins:           []string{},
		UUID:              "uuid-abc123",
		FastModeState:     FastModeOff,
	}
	for _, o := range opts {
		o(&m)
	}
	return utils.NewPattern(utils.MustJSONVersioned(m),
		"cwd", "session_id", "tools", "mcp_servers", "model",
		"slash_commands", "apiKeySource", "claude_code_version",
		"output_style", "agents", "skills", "plugins", "uuid",
	)
}

// defaultAssistantPattern returns an AssistantMessage JSON assertion pattern
func defaultAssistantPattern(opts ...func(*AssistantMessage)) utils.Pattern {
	m := AssistantMessage{
		MessageBase: MessageBase{Type: TypeAssistant},
		Message: AssistantBody{
			ID:       "msg_stub_001",
			Model:    "claude-sonnet-4-5-20250929",
			Role:     RoleAssistant,
			BodyType: AssistantBodyTypeMessage,
			Usage:    map[string]any{"input_tokens": float64(10), "output_tokens": float64(1)},
		},
		SessionID: "session-abc123",
		UUID:      "uuid-abc123",
	}
	for _, o := range opts {
		o(&m)
	}
	return utils.NewPattern(utils.MustJSONVersioned(m),
		"session_id", "uuid", "message.id", "message.model", "message.usage",
	)
}

// defaultResultPattern returns a ResultSuccessMessage JSON assertion pattern
func defaultResultPattern(opts ...func(*ResultSuccessMessage)) utils.Pattern {
	m := ResultSuccessMessage{
		MessageBase:       MessageBase{Type: TypeResult, Subtype: SubtypeSuccess},
		DurationMs:        100,
		DurationApiMs:     50,
		NumTurns:          1,
		Result:            "Hello!",
		SessionID:         "session-abc123",
		TotalCostUSD:      0.001,
		UUID:              "uuid-abc123",
		Usage:             map[string]any{"input_tokens": float64(10), "output_tokens": float64(1)},
		ModelUsage:        map[string]any{"claude-sonnet-4-5-20250929": map[string]any{"input_tokens": float64(10), "output_tokens": float64(1)}},
		PermissionDenials: []PermissionDenial{},
	}
	for _, o := range opts {
		o(&m)
	}
	return utils.NewPattern(utils.MustJSONVersioned(m),
		"duration_ms", "duration_api_ms", "num_turns", "total_cost_usd",
		"session_id", "uuid", "usage", "modelUsage", "result",
	)
}

// defaultUserToolResultPattern returns a UserToolResultMessage JSON assertion pattern
func defaultUserToolResultPattern(opts ...func(*UserToolResultMessage)) utils.Pattern {
	m := UserToolResultMessage{
		MessageBase:   MessageBase{Type: TypeUser},
		SessionID:     "session-abc123",
		UUID:          "uuid-abc123",
		ToolUseResult: map[string]any{"stdout": "command output"},
		Message:       UserToolResultBody{Role: RoleUser},
	}
	for _, o := range opts {
		o(&m)
	}
	return utils.NewPattern(utils.MustJSONVersioned(m),
		"session_id", "uuid", "tool_use_result",
	)
}

// defaultResultErrorPattern returns a ResultErrorMessage JSON assertion pattern
func defaultResultErrorPattern(opts ...func(*ResultErrorMessage)) utils.Pattern {
	m := ResultErrorMessage{
		MessageBase:       MessageBase{Type: TypeResult, Subtype: SubtypeErrorDuringExecution},
		DurationMs:        100,
		DurationApiMs:     50,
		NumTurns:          1,
		SessionID:         "session-abc123",
		TotalCostUSD:      0,
		UUID:              "uuid-abc123",
		Usage:             map[string]any{"input_tokens": float64(10), "output_tokens": float64(1)},
		ModelUsage:        map[string]any{"claude-sonnet-4-5-20250929": map[string]any{"input_tokens": float64(10), "output_tokens": float64(1)}},
		PermissionDenials: []PermissionDenial{},
	}
	for _, o := range opts {
		o(&m)
	}
	return utils.NewPattern(utils.MustJSONVersioned(m),
		"duration_ms", "duration_api_ms", "num_turns", "total_cost_usd",
		"session_id", "uuid", "usage", "modelUsage",
	)
}

// defaultResultMaxTurnsPattern returns a ResultMaxTurnsMessage JSON assertion pattern
func defaultResultMaxTurnsPattern(opts ...func(*ResultMaxTurnsMessage)) utils.Pattern {
	m := ResultMaxTurnsMessage{
		MessageBase:       MessageBase{Type: TypeResult, Subtype: SubtypeErrorMaxTurns},
		DurationMs:        100,
		DurationApiMs:     50,
		NumTurns:          1,
		SessionID:         "session-abc123",
		TotalCostUSD:      0.001,
		UUID:              "uuid-abc123",
		Usage:             map[string]any{"input_tokens": float64(10), "output_tokens": float64(1)},
		ModelUsage:        map[string]any{"claude-sonnet-4-5-20250929": map[string]any{"input_tokens": float64(10), "output_tokens": float64(1)}},
		PermissionDenials: []PermissionDenial{},
		Errors:            []string{},
	}
	for _, o := range opts {
		o(&m)
	}
	return utils.NewPattern(utils.MustJSONVersioned(m),
		"duration_ms", "duration_api_ms", "num_turns", "total_cost_usd",
		"session_id", "uuid", "usage", "modelUsage",
	)
}

// defaultStreamEventPattern returns a StreamEventMessage JSON assertion pattern
func defaultStreamEventPattern(opts ...func(*StreamEventMessage)) utils.Pattern {
	m := StreamEventMessage{
		MessageBase: MessageBase{Type: TypeStreamEvent},
		Event:       map[string]any{"type": "content_block_delta", "index": float64(0), "delta": map[string]any{"type": "text_delta", "text": "Hello"}},
		SessionID:   "session-abc123",
		UUID:        "uuid-abc123",
	}
	for _, o := range opts {
		o(&m)
	}
	return utils.NewPattern(utils.MustJSONVersioned(m),
		"event", "session_id", "uuid",
	)
}

// defaultSystemStatusPattern returns a SystemStatusMessage JSON assertion pattern
func defaultSystemStatusPattern(opts ...func(*SystemStatusMessage)) utils.Pattern {
	m := SystemStatusMessage{
		MessageBase: MessageBase{Type: TypeSystem, Subtype: SubtypeStatus},
		UUID:        "uuid-abc123",
		SessionID:   "session-abc123",
	}
	for _, o := range opts {
		o(&m)
	}
	return utils.NewPattern(utils.MustJSONVersioned(m),
		"uuid", "session_id",
	)
}

// defaultUserReplayPattern returns a UserReplayMessage JSON assertion pattern
func defaultUserReplayPattern(opts ...func(*UserReplayMessage)) utils.Pattern {
	m := UserReplayMessage{
		MessageBase: MessageBase{Type: TypeUser},
		IsReplay:    true,
		SessionID:   "session-abc123",
		UUID:        "uuid-abc123",
		Message:     UserTextBody{Role: RoleUser},
	}
	for _, o := range opts {
		o(&m)
	}
	return utils.NewPattern(utils.MustJSONVersioned(m),
		"session_id", "uuid",
	)
}

// defaultControlResponsePattern returns a ControlResponseMessage JSON assertion pattern
func defaultControlResponsePattern(opts ...func(*ControlResponseMessage)) utils.Pattern {
	m := ControlResponseMessage{
		MessageBase: MessageBase{Type: TypeControlResponse},
		Response: ControlResponseBody{
			Subtype:   "success",
			RequestID: "request-abc123",
		},
	}
	for _, o := range opts {
		o(&m)
	}
	return utils.NewPattern(utils.MustJSONVersioned(m),
		"response.subtype", "response.request_id",
	)
}

// defaultControlRequestPattern returns a ControlRequestMessage JSON assertion pattern
func defaultControlRequestPattern(opts ...func(*ControlRequestMessage)) utils.Pattern {
	m := ControlRequestMessage{
		MessageBase: MessageBase{Type: TypeControlRequest},
		RequestID:   "request-abc123",
	}
	for _, o := range opts {
		o(&m)
	}
	return utils.NewPattern(utils.MustJSONVersioned(m),
		"request_id",
	)
}
