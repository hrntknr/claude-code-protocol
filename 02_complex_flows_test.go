package ccprotocol_test

import (
	"testing"

	. "github.com/hrntknr/claudecodeprotocol"
	"github.com/hrntknr/claudecodeprotocol/utils"
)

// Text and tool use in the same response
func TestTextAndToolUseInSameResponse(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		// Request 1: text + tool_use in one response
		utils.TextAndToolUseResponse(
			"Let me check that.",
			"toolu_combo_001", "Bash", map[string]any{
				"command":     "echo combined-test",
				"description": "Combined test",
			},
		),
		// Request 2: final text after tool execution
		utils.TextResponse("Done. The output was: combined-test"),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "check and run combined"},
	}))
	utils.AssertOutput(t, s.Read(),
		defaultInitPattern(),
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				TextBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockText},
					Text:             "Done. The output was: combined-test",
				},
			}
		}),
		defaultResultPattern(func(m *ResultSuccessMessage) {
			m.Result = "Done. The output was: combined-test"
		}).Assert("result"),
	)
}

// Parallel invocation of multiple tools
func TestParallelToolUse(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		// Request 1: two tool_use blocks in one response
		utils.MultiToolUseResponse(
			utils.ToolCall{
				ID:   "toolu_par_001",
				Name: "Bash",
				Input: map[string]any{
					"command":     "echo parallel-one",
					"description": "First parallel",
				},
			},
			utils.ToolCall{
				ID:   "toolu_par_002",
				Name: "Bash",
				Input: map[string]any{
					"command":     "echo parallel-two",
					"description": "Second parallel",
				},
			},
		),
		// Request 2: final text after both tool executions
		utils.TextResponse("Both commands ran: parallel-one and parallel-two"),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "run two commands in parallel"},
	}))
	utils.AssertOutput(t, s.Read(),
		defaultInitPattern(),
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				TextBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockText},
					Text:             "Both commands ran: parallel-one and parallel-two",
				},
			}
		}),
		defaultResultPattern(func(m *ResultSuccessMessage) {
			m.Result = "Both commands ran: parallel-one and parallel-two"
		}).Assert("result"),
	)
}

// Multi-turn conversation within the same session
func TestMultiTurnConversation(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		utils.TextResponse("First answer."),
		utils.TextResponse("Second answer."),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	// Turn 1
	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "first question"},
	}))
	utils.AssertOutput(t, s.Read(),
		defaultInitPattern(),
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				TextBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockText},
					Text:             "First answer.",
				},
			}
		}),
		defaultResultPattern(func(m *ResultSuccessMessage) {
			m.Result = "First answer."
		}).Assert("result"),
	)

	// Turn 2
	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "second question"},
	}))
	utils.AssertOutput(t, s.Read(),
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				TextBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockText},
					Text:             "Second answer.",
				},
			}
		}),
		defaultResultPattern(func(m *ResultSuccessMessage) {
			m.Result = "Second answer."
		}).Assert("result"),
	)
}

// Behavior when response is truncated by max_tokens
func TestMaxTokensStopReason(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		utils.MaxTokensTextResponse("This response was truncated because it hit the max"),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "generate a very long response"},
	}))
	// Observed: The CLI retries the max_tokens response multiple times,
	// emitting the truncated text alternating with a synthetic "API Error"
	// message about the max output token limit. Eventually it produces a
	// result with subtype "success" but is_error true.
	utils.AssertOutput(t, s.Read(),
		defaultInitPattern(),
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				TextBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockText},
					Text:             "This response was truncated because it hit the max",
				},
			}
		}).Ignore("message.content.*.text"),
		defaultResultPattern(func(m *ResultSuccessMessage) {
			m.IsError = true
			m.StopReason = "end_turn"
		}).Ignore("stop_reason"),
	)
}

// Output format for responses containing multiple text blocks
func TestMultipleTextBlocks(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		utils.MultiTextResponse("First paragraph.", "Second paragraph."),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "write two paragraphs"},
	}))
	// Observed: Each text content block is emitted as a separate assistant
	// message. The result contains only the LAST text block's content.
	utils.AssertOutput(t, s.Read(),
		defaultInitPattern(),
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				TextBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockText},
					Text:             "First paragraph.",
				},
			}
		}),
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				TextBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockText},
					Text:             "Second paragraph.",
				},
			}
		}),
		defaultResultPattern(func(m *ResultSuccessMessage) {
			m.Result = "Second paragraph."
		}).Assert("result"),
	)
}
