// > Stdin input messages for mid-session configuration changes.
// > Uses `control_request`/`control_response` message types.
package ccprotocol_test

import (
	"testing"

	. "github.com/hrntknr/claudecodeprotocol"
	"github.com/hrntknr/claudecodeprotocol/utils"
)

// Permission mode change via set_permission_mode control request
func TestControlSetPermissionMode(t *testing.T) {
	stub := &utils.StubAPIServer{Responses: utils.WithInit(
		// Turn 1
		utils.TextResponse("Ready."),
		// Turn 2 (after mode change, CLI re-inits then processes user message)
		utils.TextResponse("Acknowledged."),
	)}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	// Turn 1: normal flow
	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "hello"},
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
		utils.MustJSON(ResultSuccessMessage{
			MessageBase:       MessageBase{Type: TypeResult, Subtype: SubtypeSuccess},
			IsError:           false,
			DurationMs:        utils.AnyNumber,
			DurationApiMs:     utils.AnyNumber,
			NumTurns:          utils.AnyNumber,
			Result:            "Ready.",
			SessionID:         utils.AnyString,
			TotalCostUSD:      utils.AnyNumber,
			Usage:             utils.AnyMap,
			ModelUsage:        utils.AnyMap,
			PermissionDenials: []PermissionDenial{},
			UUID:              utils.AnyString,
		}),
	)

	// Send control_request to change permission mode
	s.Send(`{"type":"control_request","request_id":"test-perm-001","request":{"subtype":"set_permission_mode","mode":"plan"}}`)

	// Turn 2: control_response + new system/init + normal flow
	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "hi"},
	}))
	// Observed: The CLI emits control_response(s) with the request_id, then a
	// new system/init with the updated permissionMode, then processes the user
	// message normally. The result text matches the last response from the stub.
	utils.AssertOutput(t, s.Read(),
		// control_response with success
		utils.MustJSON(ControlResponseMessage{
			MessageBase: MessageBase{Type: TypeControlResponse},
			Response:    utils.AnyMap,
		}),
		// New system/init with updated permission mode
		utils.MustJSON(SystemInitMessage{
			MessageBase:       MessageBase{Type: TypeSystem, Subtype: SubtypeInit},
			CWD:               utils.AnyString,
			SessionID:         utils.AnyString,
			Tools:             utils.AnyStringSlice,
			MCPServers:        utils.AnyStringSlice,
			Model:             utils.AnyString,
			PermissionMode:    PermissionPlan,
			SlashCommands:     utils.AnyStringSlice,
			APIKeySource:      utils.AnyString,
			ClaudeCodeVersion: utils.AnyString,
			OutputStyle:       utils.AnyString,
			Agents:            utils.AnyStringSlice,
			Skills:            utils.AnyStringSlice,
			Plugins:           utils.AnyStringSlice,
			UUID:              utils.AnyString,
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

// Model change via set_model control request
func TestControlSetModel(t *testing.T) {
	stub := &utils.StubAPIServer{Responses: utils.WithInit(
		// Turn 1
		utils.TextResponse("Ready."),
		// Turn 2 (after model change, CLI re-inits then processes user message)
		utils.TextResponse("Acknowledged."),
	)}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	// Turn 1: normal flow
	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "hello"},
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
		utils.MustJSON(ResultSuccessMessage{
			MessageBase:       MessageBase{Type: TypeResult, Subtype: SubtypeSuccess},
			IsError:           false,
			DurationMs:        utils.AnyNumber,
			DurationApiMs:     utils.AnyNumber,
			NumTurns:          utils.AnyNumber,
			Result:            "Ready.",
			SessionID:         utils.AnyString,
			TotalCostUSD:      utils.AnyNumber,
			Usage:             utils.AnyMap,
			ModelUsage:        utils.AnyMap,
			PermissionDenials: []PermissionDenial{},
			UUID:              utils.AnyString,
		}),
	)

	// Send control_request to change model
	s.Send(`{"type":"control_request","request_id":"test-model-001","request":{"subtype":"set_model","model":"sonnet"}}`)

	// Turn 2: control_response + new system/init + normal flow
	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "hi"},
	}))
	// Observed: The CLI emits control_response with the request_id, then a
	// new system/init with the updated model, then processes the user message.
	utils.AssertOutput(t, s.Read(),
		// control_response with success
		utils.MustJSON(ControlResponseMessage{
			MessageBase: MessageBase{Type: TypeControlResponse},
			Response:    utils.AnyMap,
		}),
		// New system/init (model field is the resolved full name)
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
