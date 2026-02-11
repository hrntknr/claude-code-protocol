package ccprotocol_test

import (
	. "github.com/hrntknr/claudecodeprotocol"
	"github.com/hrntknr/claudecodeprotocol/utils"
)

// defaultInitPattern returns a SystemInitMessage assertion pattern
// appropriate for the current CLI version. All version-specific fields
// (tracked in utils.FieldMinVersion) are included only when the
// installed CLI version is new enough.
//
// Override specific fields via functional options:
//
//	utils.MustJSON(defaultInitPattern(func(m *SystemInitMessage) {
//	    m.PermissionMode = PermissionDefault
//	}))
func defaultInitPattern(opts ...func(*SystemInitMessage)) SystemInitMessage {
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
		FastModeState:     utils.AnyString,
	}
	for _, o := range opts {
		o(&m)
	}
	return m
}
