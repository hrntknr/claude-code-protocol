package ccprotocol_test

import (
	"testing"

	. "github.com/hrntknr/claudecodeprotocol"
	"github.com/hrntknr/claudecodeprotocol/utils"
)

// Plan mode transition via the EnterPlanMode tool
func TestToolUseEnterPlanMode(t *testing.T) {
	stub := &utils.StubAPIServer{Responses: utils.WithInit(
		// Request 1: EnterPlanMode tool_use
		utils.ToolUseResponse("toolu_plan_001", "EnterPlanMode", map[string]any{}),
		// Request 2: Final text after plan mode transition
		utils.TextResponse("I have entered plan mode. Let me explore the codebase."),
	)}
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
		utils.MustJSON(defaultInitPattern()),
		utils.MustJSON(AssistantMessage{
			MessageBase: MessageBase{Type: TypeAssistant},
			Message: AssistantBody{
				Content: []IsContentBlock{
					ToolUseBlock{
						ContentBlockBase: ContentBlockBase{Type: BlockToolUse},
						ID:               utils.AnyString,
						Name:             "EnterPlanMode",
						Input:            utils.AnyMap,
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
		utils.MustJSON(SystemStatusMessage{
			MessageBase:    MessageBase{Type: TypeSystem, Subtype: SubtypeStatus},
			PermissionMode: PermissionPlan,
			UUID:           utils.AnyString,
			SessionID:      utils.AnyString,
		}),
		utils.MustJSON(UserToolResultMessage{
			MessageBase: MessageBase{Type: TypeUser},
			Message: UserToolResultBody{
				Role: RoleUser,
				Content: []ToolResultBlock{{
					ContentBlockBase: ContentBlockBase{Type: BlockToolResult},
					ToolUseID:        utils.AnyString,
					Content:          utils.AnyString,
				}},
			},
			SessionID:     utils.AnyString,
			UUID:          utils.AnyString,
			ToolUseResult: utils.AnyString,
		}),
		utils.MustJSON(AssistantMessage{
			MessageBase: MessageBase{Type: TypeAssistant},
			Message: AssistantBody{
				Content: []IsContentBlock{
					TextBlock{
						ContentBlockBase: ContentBlockBase{Type: BlockText},
						Text:             "I have entered plan mode. Let me explore the codebase.",
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
			Result:            "I have entered plan mode. Let me explore the codebase.",
			SessionID:         utils.AnyString,
			TotalCostUSD:      utils.AnyNumber,
			Usage:             utils.AnyMap,
			ModelUsage:        utils.AnyMap,
			PermissionDenials: []PermissionDenial{},
			UUID:              utils.AnyString,
		}),
	)
}

// ExitPlanMode success via --permission-prompt-tool stdio.
// The test explicitly sends a control_response on stdin to document
// the bidirectional protocol: control_request (stdout) → control_response (stdin).
func TestExitPlanModeSuccess(t *testing.T) {
	stub := &utils.StubAPIServer{Responses: utils.WithInit(
		// Request 1: ExitPlanMode tool_use
		utils.ToolUseResponse("toolu_exit_ok_001", "ExitPlanMode", map[string]any{}),
		// Request 2: Final text after plan approval
		utils.TextResponse("Plan approved, proceeding."),
	)}
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
		utils.MustJSON(defaultInitPattern(func(m *SystemInitMessage) { m.PermissionMode = PermissionDefault })),
		utils.MustJSON(AssistantMessage{
			MessageBase: MessageBase{Type: TypeAssistant},
			Message: AssistantBody{
				Content: []IsContentBlock{
					ToolUseBlock{
						ContentBlockBase: ContentBlockBase{Type: BlockToolUse},
						ID:               utils.AnyString,
						Name:             "ExitPlanMode",
						Input:            utils.AnyMap,
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
		// stdout: CLI asks for permission
		utils.MustJSON(ControlRequestMessage{
			MessageBase: MessageBase{Type: TypeControlRequest},
			RequestID:   utils.AnyString,
			Request: ControlRequest{
				Subtype:   ControlCanUseTool,
				ToolName:  "ExitPlanMode",
				Input:     utils.AnyMap,
				ToolUseID: utils.AnyString,
			},
		}),
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
		utils.MustJSON(UserToolResultMessage{
			MessageBase: MessageBase{Type: TypeUser},
			Message: UserToolResultBody{
				Role: RoleUser,
				Content: []ToolResultBlock{{
					ContentBlockBase: ContentBlockBase{Type: BlockToolResult},
					ToolUseID:        utils.AnyString,
					Content:          utils.AnyString,
				}},
			},
			SessionID:     utils.AnyString,
			UUID:          utils.AnyString,
			ToolUseResult: utils.AnyString,
		}),
		utils.MustJSON(AssistantMessage{
			MessageBase: MessageBase{Type: TypeAssistant},
			Message: AssistantBody{
				Content: []IsContentBlock{
					TextBlock{
						ContentBlockBase: ContentBlockBase{Type: BlockText},
						Text:             "Plan approved, proceeding.",
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
			Result:            "Plan approved, proceeding.",
			SessionID:         utils.AnyString,
			TotalCostUSD:      utils.AnyNumber,
			Usage:             utils.AnyMap,
			ModelUsage:        utils.AnyMap,
			PermissionDenials: []PermissionDenial{},
			UUID:              utils.AnyString,
		}),
	)
}

// ExitPlanMode tool transitions from plan mode to implementation mode
func TestToolUseExitPlanMode(t *testing.T) {
	stub := &utils.StubAPIServer{Responses: utils.WithInit(
		// Request 1: ExitPlanMode
		utils.ToolUseResponse("toolu_exit_001", "ExitPlanMode", map[string]any{}),
		// Request 2: Final text after exiting plan mode
		utils.TextResponse("Plan approved, proceeding with implementation."),
	)}
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
		utils.MustJSON(defaultInitPattern(func(m *SystemInitMessage) { m.PermissionMode = utils.AnyString })),
		utils.MustJSON(AssistantMessage{
			MessageBase: MessageBase{Type: TypeAssistant},
			Message: AssistantBody{
				Content: []IsContentBlock{
					ToolUseBlock{
						ContentBlockBase: ContentBlockBase{Type: BlockToolUse},
						ID:               utils.AnyString,
						Name:             "ExitPlanMode",
						Input:            utils.AnyMap,
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
		utils.MustJSON(UserToolResultMessage{
			MessageBase: MessageBase{Type: TypeUser},
			Message: UserToolResultBody{
				Role: RoleUser,
				Content: []ToolResultBlock{{
					ContentBlockBase: ContentBlockBase{Type: BlockToolResult},
					ToolUseID:        utils.AnyString,
					Content:          utils.AnyString,
					IsError:          true,
				}},
			},
			SessionID:     utils.AnyString,
			UUID:          utils.AnyString,
			ToolUseResult: utils.AnyString,
		}),
		utils.MustJSON(AssistantMessage{
			MessageBase: MessageBase{Type: TypeAssistant},
			Message: AssistantBody{
				Content: []IsContentBlock{
					TextBlock{
						ContentBlockBase: ContentBlockBase{Type: BlockText},
						Text:             "Plan approved, proceeding with implementation.",
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
			MessageBase:   MessageBase{Type: TypeResult, Subtype: SubtypeSuccess},
			IsError:       false,
			DurationMs:    utils.AnyNumber,
			DurationApiMs: utils.AnyNumber,
			NumTurns:      utils.AnyNumber,
			Result:        "Plan approved, proceeding with implementation.",
			SessionID:     utils.AnyString,
			TotalCostUSD:  utils.AnyNumber,
			Usage:         utils.AnyMap,
			ModelUsage:    utils.AnyMap,
			PermissionDenials: []PermissionDenial{{
				ToolName:  "ExitPlanMode",
				ToolUseID: utils.AnyString,
				ToolInput: utils.AnyMap,
			}},
			UUID: utils.AnyString,
		}),
	)
}
