package ccprotocol_test

import (
	"testing"

	. "github.com/hrntknr/claudecodeprotocol"
	"github.com/hrntknr/claudecodeprotocol/utils"
)

// Bash tool permission approved via --permission-prompt-tool stdio.
// The test explicitly sends a control_response on stdin to document
// the bidirectional protocol: control_request (stdout) → control_response (stdin).
func TestBashPermissionApproved(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		// Request 1: Bash tool_use (dangerous rm command to trigger permission)
		utils.ToolUseResponse("toolu_bash_ok_001", "Bash", map[string]any{
			"command":     "rm -f /tmp/ccprotocol_perm_test_file",
			"description": "Remove test file",
		}),
		// Request 2: Final text after tool execution
		utils.TextResponse("Command executed successfully."),
	}}
	stub.Start()
	defer stub.Close()

	// Use --permission-prompt-tool stdio without auto-handler (nil).
	// Permission responses are sent explicitly via s.Send().
	s := utils.NewSessionWithPermissionHandler(t, stub.URL(), nil)
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "remove the test file"},
	}))

	// Phase 1: Read until the CLI asks for permission via control_request.
	output1 := s.ReadUntil("control_request")
	utils.AssertOutput(t, output1,
		defaultInitPattern(func(m *SystemInitMessage) { m.PermissionMode = PermissionDefault }),
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				ToolUseBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockToolUse},
					ID:               "toolu_stub_001",
					Name:             "Bash",
					Input:            map[string]any{"command": "rm -f /tmp/ccprotocol_perm_test_file", "description": "Remove test file"},
				},
			}
		}).Ignore("message.content.*.id", "message.content.*.input"),
		// stdout: CLI asks for permission to run Bash.
		// The request includes permission_suggestions (suggested rules) and
		// blocked_path (the filesystem path that triggered the check).
		defaultControlRequestPattern(func(m *ControlRequestMessage) {
			m.Request = ControlRequest{
				Subtype:               ControlCanUseTool,
				ToolName:              "Bash",
				Input:                 map[string]any{"command": "rm -f /tmp/ccprotocol_perm_test_file", "description": "Remove test file"},
				ToolUseID:             "toolu_stub_001",
				PermissionSuggestions: []string{"allow:Bash(/tmp/*)"},
				BlockedPath:           "/tmp/ccprotocol_perm_test_file",
			}
		}).Ignore("request.input", "request.tool_use_id", "request.permission_suggestions", "request.blocked_path"),
	)

	// Phase 2: Send control_response on stdin to approve Bash.
	reqID := utils.ExtractRequestID(output1[len(output1)-1])
	s.Send(utils.MustJSON(ControlResponseMessage{
		MessageBase: MessageBase{Type: TypeControlResponse},
		Response: ControlResponseBody{
			Subtype:   "success",
			RequestID: reqID,
			Response: PermissionPayload{
				Behavior: "allow",
				UpdatedInput: map[string]any{
					"command":     "rm -f /tmp/ccprotocol_perm_test_file",
					"description": "Remove test file",
				},
			},
		},
	}))

	// Phase 3: Read the remaining output until result.
	// Note: Bash tool_result includes is_error:false explicitly (unlike
	// ExitPlanMode/AskUserQuestion which omit it). Since the struct uses
	// omitempty, we skip the tool_result assertion and verify success via
	// the result's empty permission_denials.
	output2 := s.Read()
	utils.AssertOutput(t, output2,
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				TextBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockText},
					Text:             "Command executed successfully.",
				},
			}
		}),
		// result: success with empty permission_denials
		defaultResultPattern(func(m *ResultSuccessMessage) {
			m.Result = "Command executed successfully."
		}).Assert("result"),
	)
}

// Bash tool permission denied via --permission-prompt-tool stdio.
// The test explicitly sends a deny control_response on stdin.
func TestBashPermissionDenied(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		// Request 1: Bash tool_use
		utils.ToolUseResponse("toolu_bash_deny_001", "Bash", map[string]any{
			"command":     "rm -rf /",
			"description": "Dangerous command",
		}),
		// Request 2: Final text after denial
		utils.TextResponse("I understand, I will not run that command."),
	}}
	stub.Start()
	defer stub.Close()

	// Use --permission-prompt-tool stdio without auto-handler (nil).
	s := utils.NewSessionWithPermissionHandler(t, stub.URL(), nil)
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "run rm -rf /"},
	}))

	// Phase 1: Read until the CLI asks for permission via control_request.
	output1 := s.ReadUntil("control_request")
	utils.AssertOutput(t, output1,
		defaultInitPattern(func(m *SystemInitMessage) { m.PermissionMode = PermissionDefault }),
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				ToolUseBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockToolUse},
					ID:               "toolu_stub_001",
					Name:             "Bash",
					Input:            map[string]any{"command": "rm -rf /", "description": "Dangerous command"},
				},
			}
		}).Ignore("message.content.*.id", "message.content.*.input"),
		// stdout: CLI asks for permission to run Bash
		defaultControlRequestPattern(func(m *ControlRequestMessage) {
			m.Request = ControlRequest{
				Subtype:               ControlCanUseTool,
				ToolName:              "Bash",
				Input:                 map[string]any{"command": "rm -rf /", "description": "Dangerous command"},
				ToolUseID:             "toolu_stub_001",
				PermissionSuggestions: []string{"allow:Bash(rm)"},
				DecisionReason:        "Command requires permissions",
			}
		}).Ignore("request.input", "request.tool_use_id", "request.permission_suggestions", "request.decision_reason"),
	)

	// Phase 2: Send control_response on stdin to deny Bash.
	reqID := utils.ExtractRequestID(output1[len(output1)-1])
	s.Send(utils.MustJSON(ControlResponseMessage{
		MessageBase: MessageBase{Type: TypeControlResponse},
		Response: ControlResponseBody{
			Subtype:   "success",
			RequestID: reqID,
			Response: PermissionPayload{
				Behavior: "deny",
				Message:  "Denied by test",
			},
		},
	}))

	// Phase 3: Read the remaining output until result.
	output2 := s.Read()
	utils.AssertOutput(t, output2,
		// tool_result with is_error:true (denied)
		defaultUserToolResultPattern(func(m *UserToolResultMessage) {
			m.Message.Content = []ToolResultBlock{{
				ContentBlockBase: ContentBlockBase{Type: BlockToolResult},
				ToolUseID:        "toolu_stub_001",
				Content:          "tool execution output",
				IsError:          true,
			}}
		}).Ignore("message.content.*.tool_use_id", "message.content.*.content"),
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				TextBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockText},
					Text:             "I understand, I will not run that command.",
				},
			}
		}),
		// result: Bash in permission_denials
		defaultResultPattern(func(m *ResultSuccessMessage) {
			m.Result = "I understand, I will not run that command."
			m.PermissionDenials = []PermissionDenial{{
				ToolName:  "Bash",
				ToolUseID: "toolu_stub_001",
				ToolInput: map[string]any{"command": "rm -rf /", "description": "Dangerous command"},
			}}
		}).Assert("result").Ignore("permission_denials.*.tool_use_id", "permission_denials.*.tool_input"),
	)
}
