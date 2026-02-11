package ccprotocol_test

import (
	. "github.com/hrntknr/claudecodeprotocol"
	"github.com/hrntknr/claudecodeprotocol/utils"
)

// defaultInitPattern returns a SystemInitMessage JSON assertion pattern
func defaultInitPattern(opts ...func(*SystemInitMessage)) string {
	m := SystemInitMessage{
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
		FastModeState:     FastModeOff,
	}
	for _, o := range opts {
		o(&m)
	}
	return utils.MustJSONVersioned(m)
}

// defaultAssistantPattern returns an AssistantMessage JSON assertion pattern
func defaultAssistantPattern(opts ...func(*AssistantMessage)) string {
	m := AssistantMessage{
		MessageBase: MessageBase{Type: TypeAssistant},
		Message: AssistantBody{
			ID:       utils.AnyString,
			Model:    utils.AnyString,
			Role:     RoleAssistant,
			BodyType: AssistantBodyTypeMessage,
			Usage:    utils.AnyMap,
		},
		SessionID: utils.AnyString,
		UUID:      utils.AnyString,
	}
	for _, o := range opts {
		o(&m)
	}
	return utils.MustJSONVersioned(m)
}

// defaultResultPattern returns a ResultSuccessMessage JSON assertion pattern
func defaultResultPattern(opts ...func(*ResultSuccessMessage)) string {
	m := ResultSuccessMessage{
		MessageBase:       MessageBase{Type: TypeResult, Subtype: SubtypeSuccess},
		DurationMs:        utils.AnyNumber,
		DurationApiMs:     utils.AnyNumber,
		NumTurns:          utils.AnyNumber,
		TotalCostUSD:      utils.AnyNumber,
		SessionID:         utils.AnyString,
		UUID:              utils.AnyString,
		Usage:             utils.AnyMap,
		ModelUsage:        utils.AnyMap,
		PermissionDenials: []PermissionDenial{},
	}
	for _, o := range opts {
		o(&m)
	}
	return utils.MustJSONVersioned(m)
}

// defaultUserToolResultPattern returns a UserToolResultMessage JSON assertion pattern
func defaultUserToolResultPattern(opts ...func(*UserToolResultMessage)) string {
	m := UserToolResultMessage{
		MessageBase:   MessageBase{Type: TypeUser},
		SessionID:     utils.AnyString,
		UUID:          utils.AnyString,
		ToolUseResult: utils.AnyString,
		Message:       UserToolResultBody{Role: RoleUser},
	}
	for _, o := range opts {
		o(&m)
	}
	return utils.MustJSONVersioned(m)
}

// defaultResultErrorPattern returns a ResultErrorMessage JSON assertion pattern
func defaultResultErrorPattern(opts ...func(*ResultErrorMessage)) string {
	m := ResultErrorMessage{
		MessageBase:       MessageBase{Type: TypeResult, Subtype: SubtypeErrorDuringExecution},
		DurationMs:        utils.AnyNumber,
		DurationApiMs:     utils.AnyNumber,
		NumTurns:          utils.AnyNumber,
		TotalCostUSD:      utils.AnyNumber,
		SessionID:         utils.AnyString,
		UUID:              utils.AnyString,
		Usage:             utils.AnyMap,
		ModelUsage:        utils.AnyMap,
		PermissionDenials: []PermissionDenial{},
	}
	for _, o := range opts {
		o(&m)
	}
	return utils.MustJSONVersioned(m)
}

// defaultResultMaxTurnsPattern returns a ResultMaxTurnsMessage JSON assertion pattern
func defaultResultMaxTurnsPattern(opts ...func(*ResultMaxTurnsMessage)) string {
	m := ResultMaxTurnsMessage{
		MessageBase:       MessageBase{Type: TypeResult, Subtype: SubtypeErrorMaxTurns},
		DurationMs:        utils.AnyNumber,
		DurationApiMs:     utils.AnyNumber,
		NumTurns:          utils.AnyNumber,
		TotalCostUSD:      utils.AnyNumber,
		SessionID:         utils.AnyString,
		UUID:              utils.AnyString,
		Usage:             utils.AnyMap,
		ModelUsage:        utils.AnyMap,
		PermissionDenials: []PermissionDenial{},
		Errors:            []string{},
	}
	for _, o := range opts {
		o(&m)
	}
	return utils.MustJSONVersioned(m)
}

// defaultStreamEventPattern returns a StreamEventMessage JSON assertion pattern
func defaultStreamEventPattern(opts ...func(*StreamEventMessage)) string {
	m := StreamEventMessage{
		MessageBase: MessageBase{Type: TypeStreamEvent},
		Event:       utils.AnyMap,
		SessionID:   utils.AnyString,
		UUID:        utils.AnyString,
	}
	for _, o := range opts {
		o(&m)
	}
	return utils.MustJSONVersioned(m)
}

// defaultSystemStatusPattern returns a SystemStatusMessage JSON assertion pattern
func defaultSystemStatusPattern(opts ...func(*SystemStatusMessage)) string {
	m := SystemStatusMessage{
		MessageBase: MessageBase{Type: TypeSystem, Subtype: SubtypeStatus},
		UUID:        utils.AnyString,
		SessionID:   utils.AnyString,
	}
	for _, o := range opts {
		o(&m)
	}
	return utils.MustJSONVersioned(m)
}

// defaultUserReplayPattern returns a UserReplayMessage JSON assertion pattern
func defaultUserReplayPattern(opts ...func(*UserReplayMessage)) string {
	m := UserReplayMessage{
		MessageBase: MessageBase{Type: TypeUser},
		IsReplay:    true,
		SessionID:   utils.AnyString,
		UUID:        utils.AnyString,
		Message:     UserTextBody{Role: RoleUser},
	}
	for _, o := range opts {
		o(&m)
	}
	return utils.MustJSONVersioned(m)
}

// defaultControlResponsePattern returns a ControlResponseMessage JSON assertion pattern
func defaultControlResponsePattern(opts ...func(*ControlResponseMessage)) string {
	m := ControlResponseMessage{
		MessageBase: MessageBase{Type: TypeControlResponse},
		Response: ControlResponseBody{
			Subtype:   utils.AnyString,
			RequestID: utils.AnyString,
		},
	}
	for _, o := range opts {
		o(&m)
	}
	return utils.MustJSONVersioned(m)
}

// defaultControlRequestPattern returns a ControlRequestMessage JSON assertion pattern
func defaultControlRequestPattern(opts ...func(*ControlRequestMessage)) string {
	m := ControlRequestMessage{
		MessageBase: MessageBase{Type: TypeControlRequest},
		RequestID:   utils.AnyString,
	}
	for _, o := range opts {
		o(&m)
	}
	return utils.MustJSONVersioned(m)
}
