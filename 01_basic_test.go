package ccprotocol_test

import (
	"testing"

	. "github.com/hrntknr/claudecodeprotocol"
	"github.com/hrntknr/claudecodeprotocol/utils"
)

// Basic flow for a simple text response
func TestSimpleTextResponse(t *testing.T) {
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
				Content:  []IsContentBlock{TextBlock{ContentBlockBase: ContentBlockBase{Type: BlockText}, Text: "Hello!"}},
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
			Result:            "Hello!",
			SessionID:         utils.AnyString,
			TotalCostUSD:      utils.AnyNumber,
			Usage:             utils.AnyMap,
			ModelUsage:        utils.AnyMap,
			PermissionDenials: []PermissionDenial{},
			UUID:              utils.AnyString,
		}),
	)
}

// Single Bash tool use flow
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

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "run echo tool-use-test-output"},
	}))
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
				Content:  []IsContentBlock{TextBlock{ContentBlockBase: ContentBlockBase{Type: BlockText}, Text: "The command printed: tool-use-test-output"}},
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
			Result:            "The command printed: tool-use-test-output",
			SessionID:         utils.AnyString,
			TotalCostUSD:      utils.AnyNumber,
			Usage:             utils.AnyMap,
			ModelUsage:        utils.AnyMap,
			PermissionDenials: []PermissionDenial{},
			UUID:              utils.AnyString,
		}),
	)
}

// Two-step sequential tool use flow
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

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "run two echo commands"},
	}))
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
				Content:  []IsContentBlock{TextBlock{ContentBlockBase: ContentBlockBase{Type: BlockText}, Text: "Both commands completed successfully."}},
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
			Result:            "Both commands completed successfully.",
			SessionID:         utils.AnyString,
			TotalCostUSD:      utils.AnyNumber,
			Usage:             utils.AnyMap,
			ModelUsage:        utils.AnyMap,
			PermissionDenials: []PermissionDenial{},
			UUID:              utils.AnyString,
		}),
	)
}
