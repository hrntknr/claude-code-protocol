package ccprotocol_test

import (
	"testing"

	"github.com/hrntknr/ccprotocol/utils"
)

func TestSimpleTextResponse(t *testing.T) {
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		utils.TextResponse("Hello!"),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(`{"type":"user","message":{"role":"user","content":"say hello"}}`)
	utils.AssertOutput(t, s.Read(),
		`{"type":"system", "subtype":"init"}`,
		`{"type":"assistant", "message":{"content":[{"type":"text","text":"Hello!"}]}}`,
		`{"type":"result", "subtype":"success", "result":"Hello!"}`,
	)
}

// TestToolUseBash verifies the message flow when the assistant uses the Bash tool once.
// Expected: system → assistant (final text only) → result
// Intermediate tool_use/tool_result messages are NOT emitted to stdout.
func TestToolUseBash(t *testing.T) {
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

	s.Send(`{"type":"user","message":{"role":"user","content":"run echo tool-use-test-output"}}`)
	utils.AssertOutput(t, s.Read(),
		`{"type":"system", "subtype":"init"}`,
		`{"type":"assistant", "message":{"content":[{"type":"text","text":"The command printed: tool-use-test-output"}]}}`,
		`{"type":"result", "subtype":"success"}`,
	)
}

// TestToolUseMultiStep verifies the message flow with two sequential tool uses.
// The API returns tool_use twice before returning the final text.
func TestToolUseMultiStep(t *testing.T) {
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

	s.Send(`{"type":"user","message":{"role":"user","content":"run two echo commands"}}`)
	utils.AssertOutput(t, s.Read(),
		`{"type":"system", "subtype":"init"}`,
		`{"type":"assistant", "message":{"content":[{"type":"text","text":"Both commands completed successfully."}]}}`,
		`{"type":"result", "subtype":"success"}`,
	)
}
