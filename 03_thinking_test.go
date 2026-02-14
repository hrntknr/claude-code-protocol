package ccprotocol_test

import (
	"testing"

	. "github.com/hrntknr/claudecodeprotocol"
	"github.com/hrntknr/claudecodeprotocol/utils"
)

// Response containing extended thinking blocks
func TestThinkingResponse(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		utils.ThinkingResponse(
			"Let me think about this step by step...",
			"The answer is 42.",
		),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "what is the answer?"},
	}))
	// Observed: Thinking blocks ARE emitted as a separate assistant message
	// with content[0].type="thinking". Then the text block follows as another
	// assistant message. Result contains only the text.
	utils.AssertOutput(t, s.Read(),
		defaultInitPattern(),
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				ThinkingBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockThinking},
					Thinking:         "Let me think about this step by step...",
					Signature:        "",
				},
			}
		}),
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				TextBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockText},
					Text:             "The answer is 42.",
				},
			}
		}),
		defaultResultPattern(func(m *ResultSuccessMessage) {
			m.Result = "The answer is 42."
		}).Assert("result"),
	)
}

// Flow where tool use follows a thinking block
func TestThinkingWithToolUse(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		// Request 1: thinking + tool_use
		utils.ThinkingAndToolUseResponse(
			"I need to run a command to check...",
			"toolu_think_001", "Bash", map[string]any{
				"command":     "echo thinking-tool-test",
				"description": "Test after thinking",
			},
		),
		// Request 2: final text
		utils.TextResponse("After thinking and running the command: thinking-tool-test"),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "think and then run a command"},
	}))
	utils.AssertOutput(t, s.Read(),
		defaultInitPattern(),
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				TextBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockText},
					Text:             "After thinking and running the command: thinking-tool-test",
				},
			}
		}),
		defaultResultPattern(func(m *ResultSuccessMessage) {
			m.Result = "After thinking and running the command: thinking-tool-test"
		}).Assert("result"),
	)
}
