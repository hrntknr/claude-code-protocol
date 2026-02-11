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
