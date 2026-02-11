package ccprotocol_test

import (
	"encoding/json"
	"testing"

	. "github.com/hrntknr/claudecodeprotocol"
	"github.com/hrntknr/claudecodeprotocol/utils"
)

// AskUserQuestion tool behavior in non-interactive mode
func TestToolUseAskUserQuestion(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		// Request 1: AskUserQuestion tool_use
		utils.ToolUseResponse("toolu_ask_001", "AskUserQuestion", map[string]any{
			"questions": []any{
				map[string]any{
					"question":    "Which language do you prefer?",
					"header":      "Language",
					"multiSelect": false,
					"options": []any{
						map[string]any{"label": "Go", "description": "The Go language"},
						map[string]any{"label": "Rust", "description": "The Rust language"},
					},
				},
			},
		}),
		// Request 2: Final text (after auto-answer or user interaction)
		utils.TextResponse("You chose Go. Let me proceed with Go."),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "ask me a question"},
	}))
	// Observed: In non-interactive stream-json mode, AskUserQuestion is emitted
	// as an assistant tool_use, then a user tool_result with is_error:true
	// (content "Answer questions?"). The API then returns the final text.
	// The result includes a permission_denials array listing the denied tool.
	utils.AssertOutput(t, s.Read(),
		defaultInitPattern(),
		utils.MustJSON(AssistantMessage{
			MessageBase: MessageBase{Type: TypeAssistant},
			Message: AssistantBody{
				Content: []IsContentBlock{
					ToolUseBlock{
						ContentBlockBase: ContentBlockBase{Type: BlockToolUse},
						ID:               utils.AnyString,
						Name:             "AskUserQuestion",
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
						Text:             "You chose Go. Let me proceed with Go.",
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
			Result:        "You chose Go. Let me proceed with Go.",
			SessionID:     utils.AnyString,
			TotalCostUSD:  utils.AnyNumber,
			Usage:         utils.AnyMap,
			ModelUsage:    utils.AnyMap,
			PermissionDenials: []PermissionDenial{{
				ToolName:  "AskUserQuestion",
				ToolUseID: utils.AnyString,
				ToolInput: utils.AnyMap,
			}},
			UUID: utils.AnyString,
		}),
	)
}

// AskUserQuestion success via --permission-prompt-tool stdio.
// The test explicitly sends a control_response on stdin to document
// the bidirectional protocol: control_request (stdout) → control_response (stdin).
func TestAskUserQuestionSuccess(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		// Request 1: AskUserQuestion tool_use
		utils.ToolUseResponse("toolu_ask_ok_001", "AskUserQuestion", map[string]any{
			"questions": []any{
				map[string]any{
					"question":    "Which color?",
					"header":      "Color",
					"multiSelect": false,
					"options": []any{
						map[string]any{"label": "Red", "description": "Red color"},
						map[string]any{"label": "Blue", "description": "Blue color"},
					},
				},
			},
		}),
		// Request 2: Final text after user answers
		utils.TextResponse("You chose Red."),
	}}
	stub.Start()
	defer stub.Close()

	// Use --permission-prompt-tool stdio without auto-handler (nil).
	// Permission responses are sent explicitly via s.Send().
	s := utils.NewSessionWithPermissionHandler(t, stub.URL(), nil)
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "ask me a question"},
	}))

	// Phase 1: Read until the CLI asks for permission via control_request.
	output1 := s.ReadUntil("control_request")
	utils.AssertOutput(t, output1,
		defaultInitPattern(func(m *SystemInitMessage) { m.PermissionMode = PermissionDefault }),
		utils.MustJSON(AssistantMessage{
			MessageBase: MessageBase{Type: TypeAssistant},
			Message: AssistantBody{
				Content: []IsContentBlock{
					ToolUseBlock{
						ContentBlockBase: ContentBlockBase{Type: BlockToolUse},
						ID:               utils.AnyString,
						Name:             "AskUserQuestion",
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
				ToolName:  "AskUserQuestion",
				Input:     utils.AnyMap,
				ToolUseID: utils.AnyString,
			},
		}),
	)

	// Phase 2: Send control_response on stdin with user's answers.
	reqID := utils.ExtractRequestID(output1[len(output1)-1])
	s.Send(utils.MustJSON(ControlResponseMessage{
		MessageBase: MessageBase{Type: TypeControlResponse},
		Response: ControlResponseBody{
			Subtype:   "success",
			RequestID: reqID,
			Response: PermissionPayload{
				Behavior: "allow",
				UpdatedInput: map[string]any{
					"questions": []any{
						map[string]any{
							"question":    "Which color?",
							"header":      "Color",
							"multiSelect": false,
							"options": []any{
								map[string]any{"label": "Red", "description": "Red color"},
								map[string]any{"label": "Blue", "description": "Blue color"},
							},
						},
					},
					"answers": map[string]any{
						"Which color?": "Red",
					},
				},
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
						Text:             "You chose Red.",
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
			Result:            "You chose Red.",
			SessionID:         utils.AnyString,
			TotalCostUSD:      utils.AnyNumber,
			Usage:             utils.AnyMap,
			ModelUsage:        utils.AnyMap,
			PermissionDenials: []PermissionDenial{},
			UUID:              utils.AnyString,
		}),
	)
}

// AskUserQuestion with --disallowedTools: treated as "no such tool" error, not permission denial
func TestAskUserQuestionDisallowed(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		// Request 1: API sends AskUserQuestion tool_use (but tool is disallowed)
		utils.ToolUseResponse("toolu_ask_dis_001", "AskUserQuestion", map[string]any{
			"questions": []any{
				map[string]any{
					"question":    "Pick a color?",
					"header":      "Color",
					"multiSelect": false,
					"options": []any{
						map[string]any{"label": "Red", "description": "Red color"},
						map[string]any{"label": "Blue", "description": "Blue color"},
					},
				},
			},
		}),
		// Request 2: Final text after error
		utils.TextResponse("Understood, I will not ask questions."),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSessionWithFlags(t, stub.URL(),
		[]string{"--disallowedTools", "AskUserQuestion"}, nil)
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "ask me something"},
	}))
	// Observed: --disallowedTools removes AskUserQuestion from the tools array.
	// When the API sends a tool_use for a removed tool, the CLI returns a
	// tool_result with is_error:true containing "<tool_use_error>Error: No such
	// tool available: AskUserQuestion</tool_use_error>". This is NOT recorded
	// in permission_denials (empty array).
	output := s.Read()

	// Verify AskUserQuestion is excluded from tools list.
	for _, msg := range output {
		var m map[string]any
		if err := json.Unmarshal(msg, &m); err != nil {
			continue
		}
		if m["type"] == "system" && m["subtype"] == "init" {
			tools, ok := m["tools"].([]any)
			if !ok {
				t.Fatal("tools field is not an array")
			}
			for _, tool := range tools {
				name, _ := tool.(string)
				if name == "AskUserQuestion" {
					t.Error("unexpected AskUserQuestion in tools list with --disallowedTools AskUserQuestion")
				}
			}
			break
		}
	}

	utils.AssertOutput(t, output,
		defaultInitPattern(func(m *SystemInitMessage) { m.PermissionMode = utils.AnyString }),
		utils.MustJSON(AssistantMessage{
			MessageBase: MessageBase{Type: TypeAssistant},
			Message: AssistantBody{
				Content: []IsContentBlock{
					ToolUseBlock{
						ContentBlockBase: ContentBlockBase{Type: BlockToolUse},
						ID:               utils.AnyString,
						Name:             "AskUserQuestion",
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
						Text:             "Understood, I will not ask questions.",
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
			Result:            utils.AnyString,
			SessionID:         utils.AnyString,
			TotalCostUSD:      utils.AnyNumber,
			Usage:             utils.AnyMap,
			ModelUsage:        utils.AnyMap,
			PermissionDenials: []PermissionDenial{},
			UUID:              utils.AnyString,
		}),
	)
}

// Multiple AskUserQuestion denials in a single session
func TestAskUserQuestionMultipleDenials(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		// Request 1: First AskUserQuestion
		utils.ToolUseResponse("toolu_ask_m1", "AskUserQuestion", map[string]any{
			"questions": []any{
				map[string]any{
					"question":    "Pick a framework?",
					"header":      "Framework",
					"multiSelect": false,
					"options": []any{
						map[string]any{"label": "React", "description": "React framework"},
						map[string]any{"label": "Vue", "description": "Vue framework"},
					},
				},
			},
		}),
		// Request 2: Second AskUserQuestion
		utils.ToolUseResponse("toolu_ask_m2", "AskUserQuestion", map[string]any{
			"questions": []any{
				map[string]any{
					"question":    "Pick a language?",
					"header":      "Language",
					"multiSelect": false,
					"options": []any{
						map[string]any{"label": "TypeScript", "description": "TypeScript lang"},
						map[string]any{"label": "JavaScript", "description": "JavaScript lang"},
					},
				},
			},
		}),
		// Request 3: Final text after both denials
		utils.TextResponse("I will proceed without asking."),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "ask me two questions"},
	}))
	// Observed: Each AskUserQuestion tool_use produces an assistant(tool_use) →
	// user(tool_result, is_error:true) cycle. Both denials are recorded in
	// the result's permission_denials array.
	utils.AssertOutput(t, s.Read(),
		defaultInitPattern(),
		// First AskUserQuestion cycle
		utils.MustJSON(AssistantMessage{
			MessageBase: MessageBase{Type: TypeAssistant},
			Message: AssistantBody{
				Content: []IsContentBlock{
					ToolUseBlock{
						ContentBlockBase: ContentBlockBase{Type: BlockToolUse},
						ID:               utils.AnyString,
						Name:             "AskUserQuestion",
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
		// Second AskUserQuestion cycle
		utils.MustJSON(AssistantMessage{
			MessageBase: MessageBase{Type: TypeAssistant},
			Message: AssistantBody{
				Content: []IsContentBlock{
					ToolUseBlock{
						ContentBlockBase: ContentBlockBase{Type: BlockToolUse},
						ID:               utils.AnyString,
						Name:             "AskUserQuestion",
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
		// Final text response
		utils.MustJSON(AssistantMessage{
			MessageBase: MessageBase{Type: TypeAssistant},
			Message: AssistantBody{
				Content: []IsContentBlock{
					TextBlock{
						ContentBlockBase: ContentBlockBase{Type: BlockText},
						Text:             "I will proceed without asking.",
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
			Result:        "I will proceed without asking.",
			SessionID:     utils.AnyString,
			TotalCostUSD:  utils.AnyNumber,
			Usage:         utils.AnyMap,
			ModelUsage:    utils.AnyMap,
			PermissionDenials: []PermissionDenial{
				{
					ToolName:  "AskUserQuestion",
					ToolUseID: utils.AnyString,
					ToolInput: utils.AnyMap,
				},
				{
					ToolName:  "AskUserQuestion",
					ToolUseID: utils.AnyString,
					ToolInput: utils.AnyMap,
				},
			},
			UUID: utils.AnyString,
		}),
	)
}

// AskUserQuestion and Bash in parallel: AskUserQuestion denied, Bash succeeds
func TestAskUserQuestionWithParallelTool(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		// Request 1: Parallel tool_use — AskUserQuestion + Bash in one message
		utils.MultiToolUseResponse(
			utils.ToolCall{
				ID:   "toolu_ask_par_001",
				Name: "AskUserQuestion",
				Input: map[string]any{
					"questions": []any{
						map[string]any{
							"question":    "Which DB?",
							"header":      "Database",
							"multiSelect": false,
							"options": []any{
								map[string]any{"label": "PostgreSQL", "description": "PostgreSQL DB"},
								map[string]any{"label": "MySQL", "description": "MySQL DB"},
							},
						},
					},
				},
			},
			utils.ToolCall{
				ID:   "toolu_bash_par_001",
				Name: "Bash",
				Input: map[string]any{
					"command":     "echo parallel-ok",
					"description": "Echo test",
				},
			},
		),
		// Request 2: Final text after parallel results
		utils.TextResponse("Bash succeeded, question was denied."),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "ask a question and run a command"},
	}))
	// Observed: When AskUserQuestion and Bash are in the same parallel tool_use,
	// the CLI splits them into sequential processing. AskUserQuestion is emitted
	// first as a single-tool assistant message → denied tool_result. Then Bash
	// is emitted as a separate assistant message → tool_result with is_error:true
	// containing "Sibling tool call errored" (because a sibling tool errored).
	// Only AskUserQuestion appears in permission_denials.
	utils.AssertOutput(t, s.Read(),
		defaultInitPattern(),
		// AskUserQuestion tool_use (emitted as its own assistant message)
		utils.MustJSON(AssistantMessage{
			MessageBase: MessageBase{Type: TypeAssistant},
			Message: AssistantBody{
				Content: []IsContentBlock{
					ToolUseBlock{
						ContentBlockBase: ContentBlockBase{Type: BlockToolUse},
						ID:               utils.AnyString,
						Name:             "AskUserQuestion",
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
		// AskUserQuestion denied
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
		// Bash tool_use (emitted as a separate assistant message)
		utils.MustJSON(AssistantMessage{
			MessageBase: MessageBase{Type: TypeAssistant},
			Message: AssistantBody{
				Content: []IsContentBlock{
					ToolUseBlock{
						ContentBlockBase: ContentBlockBase{Type: BlockToolUse},
						ID:               utils.AnyString,
						Name:             "Bash",
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
		// Bash tool_result — also errored due to sibling failure
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
		// Final text
		utils.MustJSON(AssistantMessage{
			MessageBase: MessageBase{Type: TypeAssistant},
			Message: AssistantBody{
				Content: []IsContentBlock{
					TextBlock{
						ContentBlockBase: ContentBlockBase{Type: BlockText},
						Text:             "Bash succeeded, question was denied.",
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
			Result:        utils.AnyString,
			SessionID:     utils.AnyString,
			TotalCostUSD:  utils.AnyNumber,
			Usage:         utils.AnyMap,
			ModelUsage:    utils.AnyMap,
			PermissionDenials: []PermissionDenial{{
				ToolName:  "AskUserQuestion",
				ToolUseID: utils.AnyString,
				ToolInput: utils.AnyMap,
			}},
			UUID: utils.AnyString,
		}),
	)
}
