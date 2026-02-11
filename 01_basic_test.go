package ccprotocol_test

import (
	"testing"

	. "github.com/hrntknr/claudecodeprotocol"
	"github.com/hrntknr/claudecodeprotocol/utils"
)

// Basic flow for a simple text response
func TestSimpleTextResponse(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		utils.TextResponse("Hello!"),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "say hello"},
	}))
	utils.AssertOutput(t, s.Read(),
		defaultInitPattern(),
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				TextBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockText},
					Text:             "Hello!",
				},
			}
		}),
		defaultResultPattern(func(m *ResultSuccessMessage) {
			m.Result = "Hello!"
		}),
	)
}

// Single Bash tool use flow
func TestToolUseBash(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		// Request 1: API tells CLI to run a Bash command
		utils.ToolUseResponse("toolu_bash_001", "Bash", map[string]any{
			"command":     "echo tool-use-test-output",
			"description": "Print test output",
		}),
		// Request 2: After tool execution, API returns final text
		utils.TextResponse("The command printed: tool-use-test-output"),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "run echo tool-use-test-output"},
	}))
	utils.AssertOutput(t, s.Read(),
		defaultInitPattern(),
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				TextBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockText},
					Text:             "The command printed: tool-use-test-output",
				},
			}
		}),
		defaultResultPattern(func(m *ResultSuccessMessage) {
			m.Result = "The command printed: tool-use-test-output"
		}),
	)
}

// Two-step sequential tool use flow
func TestToolUseMultiStep(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		// Request 1: First Bash tool use
		utils.ToolUseResponse("toolu_bash_001", "Bash", map[string]any{
			"command":     "echo step-one",
			"description": "First step",
		}),
		// Request 2: Second Bash tool use
		utils.ToolUseResponse("toolu_bash_002", "Bash", map[string]any{
			"command":     "echo step-two",
			"description": "Second step",
		}),
		// Request 3: Final text response
		utils.TextResponse("Both commands completed successfully."),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "run two echo commands"},
	}))
	utils.AssertOutput(t, s.Read(),
		defaultInitPattern(),
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				TextBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockText},
					Text:             "Both commands completed successfully.",
				},
			}
		}),
		defaultResultPattern(func(m *ResultSuccessMessage) {
			m.Result = "Both commands completed successfully."
		}),
	)
}
