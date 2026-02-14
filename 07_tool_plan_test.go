package ccprotocol_test

import (
	"testing"

	. "github.com/hrntknr/claudecodeprotocol"
	"github.com/hrntknr/claudecodeprotocol/utils"
)

// Plan mode transition via the EnterPlanMode tool
func TestToolUseEnterPlanMode(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		// Request 1: EnterPlanMode tool_use
		utils.ToolUseResponse("toolu_plan_001", "EnterPlanMode", map[string]any{}),
		// Request 2: Final text after plan mode transition
		utils.TextResponse("I have entered plan mode. Let me explore the codebase."),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "plan the implementation"},
	}))
	// Observed: EnterPlanMode emits the tool_use as an assistant message,
	// then a system status message with permissionMode:"plan", then the
	// user tool_result with plan mode instructions, then the final text.
	utils.AssertOutput(t, s.Read(),
		defaultInitPattern(),
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				ToolUseBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockToolUse},
					ID:               "toolu_stub_001",
					Name:             "EnterPlanMode",
					Input:            map[string]any{"command": "echo hello", "description": "Example"},
				},
			}
		}).Ignore("message.content.*.id", "message.content.*.input"),
		defaultSystemStatusPattern(func(m *SystemStatusMessage) {
			m.PermissionMode = PermissionPlan
		}),
		defaultUserToolResultPattern(func(m *UserToolResultMessage) {
			m.Message.Content = []ToolResultBlock{{
				ContentBlockBase: ContentBlockBase{Type: BlockToolResult},
				ToolUseID:        "toolu_stub_001",
				Content:          "tool execution output",
			}}
		}).Ignore("message.content.*.tool_use_id", "message.content.*.content"),
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				TextBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockText},
					Text:             "I have entered plan mode. Let me explore the codebase.",
				},
			}
		}),
		defaultResultPattern(func(m *ResultSuccessMessage) {
			m.Result = "I have entered plan mode. Let me explore the codebase."
		}).Assert("result"),
	)
}

// ExitPlanMode success via --permission-prompt-tool stdio.
// The test explicitly sends a control_response on stdin to document
// the bidirectional protocol: control_request (stdout) → control_response (stdin).
func TestExitPlanModeSuccess(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		// Request 1: ExitPlanMode tool_use
		utils.ToolUseResponse("toolu_exit_ok_001", "ExitPlanMode", map[string]any{}),
		// Request 2: Final text after plan approval
		utils.TextResponse("Plan approved, proceeding."),
	}}
	stub.Start()
	defer stub.Close()

	// Use --permission-prompt-tool stdio without auto-handler (nil).
	// Permission responses are sent explicitly via s.Send().
	s := utils.NewSessionWithPermissionHandler(t, stub.URL(), nil)
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "exit plan mode"},
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
					Name:             "ExitPlanMode",
					Input:            map[string]any{"command": "echo hello", "description": "Example"},
				},
			}
		}).Ignore("message.content.*.id", "message.content.*.input"),
		// stdout: CLI asks for permission
		defaultControlRequestPattern(func(m *ControlRequestMessage) {
			m.Request = ControlRequest{
				Subtype:   ControlCanUseTool,
				ToolName:  "ExitPlanMode",
				Input:     map[string]any{"command": "echo hello", "description": "Example"},
				ToolUseID: "toolu_stub_001",
			}
		}).Ignore("request.input", "request.tool_use_id"),
	)

	// Phase 2: Send control_response on stdin to approve ExitPlanMode.
	reqID := utils.ExtractRequestID(output1[len(output1)-1])
	s.Send(utils.MustJSON(ControlResponseMessage{
		MessageBase: MessageBase{Type: TypeControlResponse},
		Response: ControlResponseBody{
			Subtype:   "success",
			RequestID: reqID,
			Response: PermissionPayload{
				Behavior:     "allow",
				UpdatedInput: map[string]any{},
			},
		},
	}))

	// Phase 3: Read the remaining output until result.
	output2 := s.Read()
	utils.AssertOutput(t, output2,
		// tool_result with is_error absent (success)
		defaultUserToolResultPattern(func(m *UserToolResultMessage) {
			m.Message.Content = []ToolResultBlock{{
				ContentBlockBase: ContentBlockBase{Type: BlockToolResult},
				ToolUseID:        "toolu_stub_001",
				Content:          "tool execution output",
			}}
		}).Ignore("message.content.*.tool_use_id", "message.content.*.content"),
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				TextBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockText},
					Text:             "Plan approved, proceeding.",
				},
			}
		}),
		defaultResultPattern(func(m *ResultSuccessMessage) {
			m.Result = "Plan approved, proceeding."
		}).Assert("result"),
	)
}

// ExitPlanMode tool transitions from plan mode to implementation mode
func TestToolUseExitPlanMode(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		// Request 1: ExitPlanMode
		utils.ToolUseResponse("toolu_exit_001", "ExitPlanMode", map[string]any{}),
		// Request 2: Final text after exiting plan mode
		utils.TextResponse("Plan approved, proceeding with implementation."),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSessionWithFlags(t, stub.URL(),
		[]string{"--permission-mode", "plan"}, nil)
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "exit plan mode"},
	}))
	// Observed: In non-interactive stream-json mode, ExitPlanMode requires user
	// approval which cannot be obtained. The CLI treats it as a permission denial:
	// tool_use → tool_result(is_error:true, "Exit plan mode?") → final text.
	// The result includes ExitPlanMode in permission_denials.
	// No system/status message is emitted (unlike EnterPlanMode).
	output := s.Read()
	utils.AssertOutput(t, output,
		defaultInitPattern(func(m *SystemInitMessage) { m.PermissionMode = PermissionBypassPermissions }).Ignore("permissionMode"),
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				ToolUseBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockToolUse},
					ID:               "toolu_stub_001",
					Name:             "ExitPlanMode",
					Input:            map[string]any{"command": "echo hello", "description": "Example"},
				},
			}
		}).Ignore("message.content.*.id", "message.content.*.input"),
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
					Text:             "Plan approved, proceeding with implementation.",
				},
			}
		}),
		defaultResultPattern(func(m *ResultSuccessMessage) {
			m.Result = "Plan approved, proceeding with implementation."
			m.PermissionDenials = []PermissionDenial{{
				ToolName:  "ExitPlanMode",
				ToolUseID: "toolu_stub_001",
				ToolInput: map[string]any{"key": "value"},
			}}
		}).Assert("result").Ignore("permission_denials.*.tool_use_id", "permission_denials.*.tool_input"),
	)
}
