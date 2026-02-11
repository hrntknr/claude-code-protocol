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
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		// Turn 1
		utils.TextResponse("Ready."),
		// Turn 2 (after mode change, CLI re-inits then processes user message)
		utils.TextResponse("Acknowledged."),
	}}
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
		defaultInitPattern(),
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
	s.Send(utils.MustJSON(ControlRequestMessage{
		MessageBase: MessageBase{Type: TypeControlRequest},
		RequestID:   "test-perm-001",
		Request: ControlRequest{
			Subtype: ControlSetPermissionMode,
			Mode:    "plan",
		},
	}))

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
			Response: ControlResponseBody{
				Subtype:   utils.AnyString,
				RequestID: utils.AnyString,
			},
		}),
		// New system/init with updated permission mode
		defaultInitPattern(func(m *SystemInitMessage) { m.PermissionMode = PermissionPlan }),
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
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		// Turn 1
		utils.TextResponse("Ready."),
		// Turn 2 (after model change, CLI re-inits then processes user message)
		utils.TextResponse("Acknowledged."),
	}}
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
		defaultInitPattern(),
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
	s.Send(utils.MustJSON(ControlRequestMessage{
		MessageBase: MessageBase{Type: TypeControlRequest},
		RequestID:   "test-model-001",
		Request: ControlRequest{
			Subtype: ControlSetModel,
			Model:   "sonnet",
		},
	}))

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
			Response: ControlResponseBody{
				Subtype:   utils.AnyString,
				RequestID: utils.AnyString,
			},
		}),
		// New system/init (model field is the resolved full name)
		defaultInitPattern(),
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
