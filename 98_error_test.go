package ccprotocol_test

import (
	"testing"

	. "github.com/hrntknr/claudecodeprotocol"
	"github.com/hrntknr/claudecodeprotocol/utils"
)

// Error handling on tool execution failure
func TestToolError(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		// Request 1: Read a non-existent file
		utils.ToolUseResponse("toolu_err_001", "Read", map[string]any{
			"file_path": "/tmp/this-file-does-not-exist-for-test-12345.txt",
		}),
		// Request 2: The API acknowledges the error and responds
		utils.TextResponse("The file does not exist. Let me handle this error."),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "read a missing file"},
	}))
	// The CLI should handle the tool error gracefully.
	// The API receives the error as a tool_result with is_error=true,
	// then returns a normal text response.
	utils.AssertOutput(t, s.Read(),
		defaultInitPattern(),
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				TextBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockText},
					Text:             "The file does not exist. Let me handle this error.",
				},
			}
		}),
		defaultResultPattern(func(m *ResultSuccessMessage) {
			m.Result = "The file does not exist. Let me handle this error."
		}).Assert("result"),
	)
}

// Behavior when receiving API-level SSE error events
func TestAPIError(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		// All requests (including init) get the same error.
		// The CLI should handle the API error.
		utils.ErrorSSEResponse("overloaded_error", "Overloaded"),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "trigger an error"},
	}))
	// Observed: When the API returns an error, the CLI emits a result with
	// subtype "error_during_execution" and an "errors" array containing error
	// message strings (full stack traces). No assistant messages are emitted.
	utils.AssertOutput(t, s.Read(),
		defaultInitPattern(),
		defaultResultErrorPattern(func(m *ResultErrorMessage) {
			m.Errors = []string{"API error: Overloaded"}
		}).Ignore("errors"),
	)
}
