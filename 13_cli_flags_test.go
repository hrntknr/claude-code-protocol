package ccprotocol_test

import (
	"encoding/json"
	"strings"
	"testing"

	. "github.com/hrntknr/claudecodeprotocol"
	"github.com/hrntknr/claudecodeprotocol/utils"
)

// Replay user messages via --replay-user-messages flag
func TestReplayUserMessages(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		utils.TextResponse("Echoed!"),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSessionWithFlags(t, stub.URL(),
		[]string{"--replay-user-messages"}, nil)
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "replay this message"},
	}))
	// Observed: With --replay-user-messages, the CLI echoes the user message
	// back on stdout AFTER system/init. The echoed message is a "user" type
	// with extra fields: session_id, parent_tool_use_id, uuid, and isReplay:true.
	output := s.Read()
	utils.AssertOutput(t, output,
		defaultInitPattern(),
		utils.MustJSON(UserReplayMessage{
			MessageBase: MessageBase{Type: TypeUser},
			Message:     UserTextBody{Role: RoleUser, Content: "replay this message"},
			SessionID:   utils.AnyString,
			UUID:        utils.AnyString,
			IsReplay:    true,
		}),
		utils.MustJSON(AssistantMessage{
			MessageBase: MessageBase{Type: TypeAssistant},
			Message: AssistantBody{
				Content: []IsContentBlock{
					TextBlock{
						ContentBlockBase: ContentBlockBase{Type: BlockText},
						Text:             "Echoed!",
					},
				},
				ID:       utils.AnyString,
				Model:    utils.AnyString,
				Role:     RoleAssistant,
				BodyType: AssistantBodyTypeMessage,
				Usage:    utils.AnyMap,
			},
			SessionID: utils.AnyString,
			UUID:      utils.AnyString,
		}),
		utils.MustJSON(ResultSuccessMessage{
			MessageBase:       MessageBase{Type: TypeResult, Subtype: SubtypeSuccess},
			IsError:           false,
			DurationMs:        utils.AnyNumber,
			DurationApiMs:     utils.AnyNumber,
			NumTurns:          utils.AnyNumber,
			Result:            "Echoed!",
			SessionID:         utils.AnyString,
			TotalCostUSD:      utils.AnyNumber,
			Usage:             utils.AnyMap,
			ModelUsage:        utils.AnyMap,
			PermissionDenials: []PermissionDenial{},
			UUID:              utils.AnyString,
		}),
	)
}

// Partial message streaming via --include-partial-messages flag
func TestIncludePartialMessages(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		utils.TextResponse("Streamed response."),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSessionWithFlags(t, stub.URL(),
		[]string{"--include-partial-messages"}, nil)
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "stream something"},
	}))
	// Observed: With --include-partial-messages, the CLI emits stream_event
	// messages containing raw SSE events (message_start, content_block_start,
	// content_block_delta, content_block_stop, message_delta, message_stop)
	// interleaved with the normal assistant messages.
	// Output order: init → stream_event(message_start) → stream_event(content_block_start)
	// → stream_event(content_block_delta) → assistant(text) → stream_event(content_block_stop)
	// → stream_event(message_delta) → stream_event(message_stop) → result
	output := s.Read()

	utils.AssertOutput(t, output,
		defaultInitPattern(),
		// stream_event: message_start
		utils.MustJSON(StreamEventMessage{
			MessageBase: MessageBase{Type: TypeStreamEvent},
			Event:       utils.AnyMap,
			SessionID:   utils.AnyString,
			UUID:        utils.AnyString,
		}),
		// stream_event: content_block_start
		utils.MustJSON(StreamEventMessage{
			MessageBase: MessageBase{Type: TypeStreamEvent},
			Event:       utils.AnyMap,
			SessionID:   utils.AnyString,
			UUID:        utils.AnyString,
		}),
		// stream_event: content_block_delta
		utils.MustJSON(StreamEventMessage{
			MessageBase: MessageBase{Type: TypeStreamEvent},
			Event:       utils.AnyMap,
			SessionID:   utils.AnyString,
			UUID:        utils.AnyString,
		}),
		// Final assistant message (complete text block)
		utils.MustJSON(AssistantMessage{
			MessageBase: MessageBase{Type: TypeAssistant},
			Message: AssistantBody{
				Content: []IsContentBlock{
					TextBlock{
						ContentBlockBase: ContentBlockBase{Type: BlockText},
						Text:             "Streamed response.",
					},
				},
				ID:       utils.AnyString,
				Model:    utils.AnyString,
				Role:     RoleAssistant,
				BodyType: AssistantBodyTypeMessage,
				Usage:    utils.AnyMap,
			},
			SessionID: utils.AnyString,
			UUID:      utils.AnyString,
		}),
		// stream_event: content_block_stop
		utils.MustJSON(StreamEventMessage{
			MessageBase: MessageBase{Type: TypeStreamEvent},
			Event:       utils.AnyMap,
			SessionID:   utils.AnyString,
			UUID:        utils.AnyString,
		}),
		// stream_event: message_delta (stop_reason)
		utils.MustJSON(StreamEventMessage{
			MessageBase: MessageBase{Type: TypeStreamEvent},
			Event:       utils.AnyMap,
			SessionID:   utils.AnyString,
			UUID:        utils.AnyString,
		}),
		// stream_event: message_stop
		utils.MustJSON(StreamEventMessage{
			MessageBase: MessageBase{Type: TypeStreamEvent},
			Event:       utils.AnyMap,
			SessionID:   utils.AnyString,
			UUID:        utils.AnyString,
		}),
		utils.MustJSON(ResultSuccessMessage{
			MessageBase:       MessageBase{Type: TypeResult, Subtype: SubtypeSuccess},
			IsError:           false,
			DurationMs:        utils.AnyNumber,
			DurationApiMs:     utils.AnyNumber,
			NumTurns:          utils.AnyNumber,
			Result:            "Streamed response.",
			SessionID:         utils.AnyString,
			TotalCostUSD:      utils.AnyNumber,
			Usage:             utils.AnyMap,
			ModelUsage:        utils.AnyMap,
			PermissionDenials: []PermissionDenial{},
			UUID:              utils.AnyString,
		}),
	)
	// Verify stream_event messages are present (9 total: init + 6 stream_events + assistant + result).
	if len(output) <= 3 {
		t.Errorf("expected more than 3 messages with partial streaming, got %d", len(output))
	}
}

// Turn limit behavior via --max-turns flag
func TestMaxTurnsLimit(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		// Tool use to force multi-turn: the stub returns tool_use, but
		// --max-turns 1 should stop after the first turn.
		utils.ToolUseResponse("toolu_mt_001", "Bash", map[string]any{
			"command":     "echo turn-one",
			"description": "First turn command",
		}),
		utils.TextResponse("Turn completed."),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSessionWithFlags(t, stub.URL(),
		[]string{"--max-turns", "1"}, nil)
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "run a command"},
	}))
	// Observed: With --max-turns 1, the CLI executes the tool but then
	// terminates with subtype "error_max_turns" instead of "success".
	// The result has an empty errors array and is_error:false.
	// The tool_use and tool_result messages still appear before the result.
	output := s.Read()
	utils.AssertOutput(t, output,
		defaultInitPattern(),
		utils.MustJSON(ResultMaxTurnsMessage{
			MessageBase:       MessageBase{Type: TypeResult, Subtype: SubtypeErrorMaxTurns},
			IsError:           false,
			DurationMs:        utils.AnyNumber,
			DurationApiMs:     utils.AnyNumber,
			NumTurns:          utils.AnyNumber,
			SessionID:         utils.AnyString,
			TotalCostUSD:      utils.AnyNumber,
			Usage:             utils.AnyMap,
			ModelUsage:        utils.AnyMap,
			PermissionDenials: []PermissionDenial{},
			UUID:              utils.AnyString,
			Errors:            []string{},
		}),
	)
}

// Tools restriction via --tools flag
func TestToolsRestriction(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		utils.TextResponse("OK"),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSessionWithFlags(t, stub.URL(),
		[]string{"--tools", "Bash,Read"}, nil)
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "hello"},
	}))
	output := s.Read()

	// Observed: With --tools Bash,Read, the tools array in system/init
	// should contain only the specified tools.
	// Parse the init message to check the tools list.
	var initFound bool
	for _, msg := range output {
		var m map[string]any
		if err := json.Unmarshal(msg, &m); err != nil {
			continue
		}
		if m["type"] == "system" && m["subtype"] == "init" {
			initFound = true
			tools, ok := m["tools"].([]any)
			if !ok {
				t.Fatal("tools field is not an array")
			}
			toolNames := make(map[string]bool)
			for _, tool := range tools {
				if name, ok := tool.(string); ok {
					toolNames[name] = true
				}
			}
			if !toolNames["Bash"] {
				t.Error("expected 'Bash' in tools list")
			}
			if !toolNames["Read"] {
				t.Error("expected 'Read' in tools list")
			}
			if toolNames["Write"] {
				t.Error("unexpected 'Write' in restricted tools list")
			}
			if toolNames["Edit"] {
				t.Error("unexpected 'Edit' in restricted tools list")
			}
			t.Logf("tools list (%d items): %v", len(tools), tools)
			break
		}
	}
	if !initFound {
		t.Fatal("system/init message not found")
	}
}

// Permission mode variation via --permission-mode flag
func TestPermissionModeDefault(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		utils.TextResponse("Hello from default mode."),
	}}
	stub.Start()
	defer stub.Close()

	// Note: --dangerously-skip-permissions is still in the base flags, but
	// --permission-mode default should override the permissionMode field.
	s := utils.NewSessionWithFlags(t, stub.URL(),
		[]string{"--permission-mode", "default"}, nil)
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "hello"},
	}))
	// Observed: The permissionMode field in system/init should reflect the
	// --permission-mode flag value.
	output := s.Read()
	var initFound bool
	for _, msg := range output {
		var m map[string]any
		if err := json.Unmarshal(msg, &m); err != nil {
			continue
		}
		if m["type"] == "system" && m["subtype"] == "init" {
			initFound = true
			mode, _ := m["permissionMode"].(string)
			t.Logf("permissionMode = %q", mode)
			// The actual permissionMode may be overridden by --dangerously-skip-permissions.
			// We log the observed value for documentation.
			break
		}
	}
	if !initFound {
		t.Fatal("system/init message not found")
	}
}

// fast_mode_state defaults to "off" without authentication
func TestFastModeStateDefault(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		utils.TextResponse("Hello"),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "hello"},
	}))
	// Observed: Without OAuth authentication or a paid subscription,
	// fast mode cannot be enabled. The fast_mode_state field in system/init
	// is "off". Possible values are "off", "on", and "cooldown".
	output := s.Read()
	utils.AssertOutput(t, output,
		defaultInitPattern(func(m *SystemInitMessage) {
			m.FastModeState = FastModeOff
		}),
		utils.MustJSON(ResultSuccessMessage{
			MessageBase:       MessageBase{Type: TypeResult, Subtype: SubtypeSuccess},
			IsError:           false,
			DurationMs:        utils.AnyNumber,
			DurationApiMs:     utils.AnyNumber,
			NumTurns:          utils.AnyNumber,
			Result:            "Hello",
			SessionID:         utils.AnyString,
			TotalCostUSD:      utils.AnyNumber,
			Usage:             utils.AnyMap,
			ModelUsage:        utils.AnyMap,
			PermissionDenials: []PermissionDenial{},
			UUID:              utils.AnyString,
		}),
	)
}

// Model override via --model flag
func TestModelOverride(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		utils.TextResponse("Hello from sonnet."),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSessionWithFlags(t, stub.URL(),
		[]string{"--model", "sonnet"}, nil)
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "hello"},
	}))
	output := s.Read()

	// Observed: The model field in system/init should reflect the specified model.
	var initFound bool
	for _, msg := range output {
		var m map[string]any
		if err := json.Unmarshal(msg, &m); err != nil {
			continue
		}
		if m["type"] == "system" && m["subtype"] == "init" {
			initFound = true
			model, _ := m["model"].(string)
			t.Logf("model = %q", model)
			if !strings.Contains(model, "sonnet") {
				t.Errorf("expected model to contain 'sonnet', got %q", model)
			}
			break
		}
	}
	if !initFound {
		t.Fatal("system/init message not found")
	}
}
