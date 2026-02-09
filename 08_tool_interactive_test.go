package ccprotocol_test

import (
	"testing"

	. "github.com/hrntknr/claudecodeprotocol"
	"github.com/hrntknr/claudecodeprotocol/utils"
)

// AskUserQuestion tool behavior in non-interactive mode
func TestToolUseAskUserQuestion(t *testing.T) {
	stub := &utils.StubAPIServer{Responses: utils.WithInit(
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
	)}
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
		utils.MustJSON(SystemInitMessage{
			MessageBase:       MessageBase{Type: TypeSystem, Subtype: SubtypeInit},
			CWD:               utils.AnyString,
			SessionID:         utils.AnyString,
			Tools:             utils.AnyStringSlice,
			MCPServers:        utils.AnyStringSlice,
			Model:             utils.AnyString,
			PermissionMode:    PermissionBypassPermissions,
			SlashCommands:     utils.AnyStringSlice,
			APIKeySource:      utils.AnyString,
			ClaudeCodeVersion: utils.AnyString,
			OutputStyle:       utils.AnyString,
			Agents:            utils.AnyStringSlice,
			Skills:            utils.AnyStringSlice,
			Plugins:           utils.AnyStringSlice,
			UUID:              utils.AnyString,
		}),
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
			MessageBase:       MessageBase{Type: TypeResult, Subtype: SubtypeSuccess},
			IsError:           false,
			DurationMs:        utils.AnyNumber,
			DurationApiMs:     utils.AnyNumber,
			NumTurns:          utils.AnyNumber,
			Result:            "You chose Go. Let me proceed with Go.",
			SessionID:         utils.AnyString,
			TotalCostUSD:      utils.AnyNumber,
			Usage:             utils.AnyMap,
			ModelUsage:        utils.AnyMap,
			PermissionDenials: []PermissionDenial{{
				ToolName:  "AskUserQuestion",
				ToolUseID: utils.AnyString,
				ToolInput: utils.AnyMap,
			}},
			UUID: utils.AnyString,
		}),
	)
}
