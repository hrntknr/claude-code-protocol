package ccprotocol_test

import (
	"testing"

	. "github.com/hrntknr/claudecodeprotocol"
	"github.com/hrntknr/claudecodeprotocol/utils"
)

// URL fetching behavior of the WebFetch tool
func TestToolUseWebFetch(t *testing.T) {
	stub := &utils.StubAPIServer{
		StaticPages: map[string]string{
			"/static/test-page": "<html><body><h1>Test Page</h1><p>Content for WebFetch test.</p></body></html>",
		},
		Responses: utils.WithInit(
			// Request 1: WebFetch tool_use
			utils.ToolUseResponse("toolu_wf_001", "WebFetch", map[string]any{
				"url":    "", // placeholder, will be set after Start()
				"prompt": "What does the page say?",
			}),
			// Extra responses for WebFetch's internal haiku processing
			utils.TextResponse("ok"),
			utils.TextResponse("ok"),
			utils.TextResponse("ok"),
			// Final text
			utils.TextResponse("The page contains a heading: Test Page"),
		),
	}
	stub.Start()
	defer stub.Close()

	// Update the WebFetch URL to point to our stub's static page.
	// The ToolUseResponse is already built, so we rebuild with the correct URL.
	targetURL := stub.URL() + "/static/test-page"
	stub.Responses = utils.WithInit(
		utils.ToolUseResponse("toolu_wf_001", "WebFetch", map[string]any{
			"url":    targetURL,
			"prompt": "What does the page say?",
		}),
		// Extra responses for WebFetch's internal haiku processing
		utils.TextResponse("ok"),
		utils.TextResponse("ok"),
		utils.TextResponse("ok"),
		utils.TextResponse("The page contains a heading: Test Page"),
	)

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "fetch the test page"},
	}))
	// Observed: WebFetch upgrades HTTP to HTTPS, causing an SSL error when
	// hitting the plain HTTP stub server. The CLI emits the tool_use, then
	// a user tool_result with is_error:true containing the SSL error.
	// The API then returns the next response as final text.
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
			FastModeState:     utils.AnyString,
		}),
		utils.MustJSON(AssistantMessage{
			MessageBase: MessageBase{Type: TypeAssistant},
			Message: AssistantBody{
				Content: []IsContentBlock{
					ToolUseBlock{
						ContentBlockBase: ContentBlockBase{Type: BlockToolUse},
						ID:               utils.AnyString,
						Name:             "WebFetch",
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

// WebSearch tool behavior (may fail due to test environment restrictions)
func TestToolUseWebSearch(t *testing.T) {
	stub := &utils.StubAPIServer{Responses: utils.WithInit(
		utils.ToolUseResponse("toolu_ws_001", "WebSearch", map[string]any{
			"query": "protocol test query",
		}),
		utils.TextResponse("Search completed."),
	)}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "search the web"},
	}))
	// Observed: WebSearch may return a tool_result with is_error:true if the
	// environment does not support web search (e.g., no valid API key or
	// not in US region). The protocol format is the same as other tool errors.
	output := s.Read()
	utils.AssertOutput(t, output,
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
			FastModeState:     utils.AnyString,
		}),
		utils.MustJSON(AssistantMessage{
			MessageBase: MessageBase{Type: TypeAssistant},
			Message: AssistantBody{
				Content: []IsContentBlock{
					ToolUseBlock{
						ContentBlockBase: ContentBlockBase{Type: BlockToolUse},
						ID:               utils.AnyString,
						Name:             "WebSearch",
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
				}},
			},
			SessionID:     utils.AnyString,
			UUID:          utils.AnyString,
			ToolUseResult: utils.AnyString,
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
