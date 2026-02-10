package ccprotocol_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	. "github.com/hrntknr/claudecodeprotocol"
	"github.com/hrntknr/claudecodeprotocol/utils"
)

// Allowed tools restriction: only specified tools are permitted
func TestAllowedToolsRestriction(t *testing.T) {
	tmpDir := t.TempDir()
	targetFile := filepath.Join(tmpDir, "restricted.txt")

	stub := &utils.StubAPIServer{Responses: utils.WithInit(
		// API tells CLI to use Write tool (not in allowed list)
		utils.ToolUseResponse("toolu_wr_001", "Write", map[string]any{
			"file_path": targetFile,
			"content":   "should-not-be-written",
		}),
		utils.TextResponse("Write tool was denied."),
	)}
	stub.Start()
	defer stub.Close()

	s := utils.NewSessionWithFlags(t, stub.URL(),
		[]string{"--allowedTools", "Bash,Read"}, nil)
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "write a file"},
	}))
	// Observed: --allowedTools with --dangerously-skip-permissions does NOT
	// restrict tool execution. The tools list in init is unchanged (all tools
	// shown), the Write tool executes normally, and there are no
	// permission_denials. The tool_result has is_error:false.
	// This means --allowedTools is overridden by --dangerously-skip-permissions.
	output := s.Read()
	utils.AssertOutput(t, output,
		utils.MustJSON(SystemInitMessage{
			MessageBase:       MessageBase{Type: TypeSystem, Subtype: SubtypeInit},
			CWD:               utils.AnyString,
			SessionID:         utils.AnyString,
			Tools:             utils.AnyStringSlice,
			MCPServers:        utils.AnyStringSlice,
			Model:             utils.AnyString,
			PermissionMode:    utils.AnyString,
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
						Name:             "Write",
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
		utils.MustJSON(AssistantMessage{
			MessageBase: MessageBase{Type: TypeAssistant},
			Message: AssistantBody{
				Content: []IsContentBlock{
					TextBlock{
						ContentBlockBase: ContentBlockBase{Type: BlockText},
						Text:             "Write tool was denied.",
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

// Disallowed tools restriction: specified tools are denied
func TestDisallowedToolsRestriction(t *testing.T) {
	tmpDir := t.TempDir()
	targetFile := filepath.Join(tmpDir, "denied.txt")

	stub := &utils.StubAPIServer{Responses: utils.WithInit(
		// API tells CLI to use Write tool (in disallowed list)
		utils.ToolUseResponse("toolu_wr_001", "Write", map[string]any{
			"file_path": targetFile,
			"content":   "should-not-be-written",
		}),
		utils.TextResponse("Write tool was denied."),
	)}
	stub.Start()
	defer stub.Close()

	s := utils.NewSessionWithFlags(t, stub.URL(),
		[]string{"--disallowedTools", "Write,Edit"}, nil)
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "write a file"},
	}))
	// Observed: --disallowedTools removes the tools from the tools array in
	// system/init. When the API sends a tool_use for a removed tool, the CLI
	// returns a tool_result with is_error:true and content containing
	// "<tool_use_error>Error: No such tool available: Write</tool_use_error>".
	// Note: This is NOT recorded in permission_denials (it's treated as a
	// "no such tool" error, not a permission denial).
	output := s.Read()

	// Verify Write and Edit are excluded from tools list.
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
				if name == "Write" || name == "Edit" {
					t.Errorf("unexpected %q in tools list with --disallowedTools Write,Edit", name)
				}
			}
			break
		}
	}

	utils.AssertOutput(t, output,
		utils.MustJSON(SystemInitMessage{
			MessageBase:       MessageBase{Type: TypeSystem, Subtype: SubtypeInit},
			CWD:               utils.AnyString,
			SessionID:         utils.AnyString,
			Tools:             utils.AnyStringSlice,
			MCPServers:        utils.AnyStringSlice,
			Model:             utils.AnyString,
			PermissionMode:    utils.AnyString,
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
						Name:             "Write",
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
						Text:             "Write tool was denied.",
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

	// Verify the file was NOT created.
	if _, err := os.Stat(targetFile); err == nil {
		t.Error("file should not have been created when Write tool is disallowed")
	}
}
