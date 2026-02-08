//go:generate go run ./cmd/gendoc
package ccprotocol

// ---------------------------------------------------------------------------
// Enum
// ---------------------------------------------------------------------------

type MessageType string

const (
	TypeSystem    MessageType = "system"
	TypeAssistant MessageType = "assistant"
	TypeUser      MessageType = "user"
	TypeResult    MessageType = "result"
)

type MessageSubtype string

const (
	SubtypeInit                 MessageSubtype = "init"
	SubtypeStatus               MessageSubtype = "status"
	SubtypeSuccess              MessageSubtype = "success"
	SubtypeErrorDuringExecution MessageSubtype = "error_during_execution"
)

type MessageRole string

const (
	RoleAssistant MessageRole = "assistant"
	RoleUser      MessageRole = "user"
)

type ContentBlockType string

const (
	BlockText       ContentBlockType = "text"
	BlockToolUse    ContentBlockType = "tool_use"
	BlockThinking   ContentBlockType = "thinking"
	BlockToolResult ContentBlockType = "tool_result"
)

type PermissionMode string

const (
	PermissionBypassPermissions PermissionMode = "bypassPermissions"
	PermissionPlan              PermissionMode = "plan"
)

type AssistantBodyType string

const (
	AssistantBodyTypeMessage AssistantBodyType = "message"
)

// ---------------------------------------------------------------------------
// Interfaces
// ---------------------------------------------------------------------------

// IsMessage is the interface that constrains protocol message types.
type IsMessage interface {
	isMessage()
}

// IsContentBlock is the interface that constrains content block types.
type IsContentBlock interface {
	isContentBlock()
}

// ---------------------------------------------------------------------------
// Base types
// ---------------------------------------------------------------------------

// MessageBase holds fields common to all messages.
type MessageBase struct {
	Type    MessageType    `json:"type"`
	Subtype MessageSubtype `json:"subtype,omitempty"`
}

func (MessageBase) isMessage() {}

// ContentBlockBase holds fields common to all content blocks.
type ContentBlockBase struct {
	Type ContentBlockType `json:"type"`
}

func (ContentBlockBase) isContentBlock() {}

// ---------------------------------------------------------------------------
// Message types
// ---------------------------------------------------------------------------

// # system/init
// The first message output when the CLI starts. Contains initial information such as session ID and version.
// Always output as the first message of a session; the tools field contains the available tools list.
// Also includes extended information such as mcp_servers, model, slash_commands, agents, skills, and plugins.
//
// ```json
// {"type":"system","subtype":"init","cwd":"/path","session_id":"abc","tools":["Bash","Read"],"mcp_servers":[],"model":"claude-opus-4-6","permissionMode":"bypassPermissions","slash_commands":[...],"apiKeySource":"none","claude_code_version":"2.1.37","output_style":"default","agents":[...],"skills":[...],"plugins":[],"uuid":"xxx"}
// ```
type SystemInitMessage struct {
	MessageBase
	CWD               string         `json:"cwd"`                 // Working directory
	SessionID         string         `json:"session_id"`          // Session ID
	Tools             []string       `json:"tools"`               // Available tools list
	MCPServers        []string       `json:"mcp_servers"`         // MCP servers list
	Model             string         `json:"model"`               // Model name
	PermissionMode    PermissionMode `json:"permissionMode"`      // Permission mode
	SlashCommands     []string       `json:"slash_commands"`      // Slash commands list
	APIKeySource      string         `json:"apiKeySource"`        // API key source
	ClaudeCodeVersion string         `json:"claude_code_version"` // CLI version
	OutputStyle       string         `json:"output_style"`        // Output style
	Agents            []string       `json:"agents"`              // Agents list
	Skills            []string       `json:"skills"`              // Skills list
	Plugins           []string       `json:"plugins"`             // Plugins list
	UUID              string         `json:"uuid"`                // Message UUID
}

// # system/status
// A message that notifies system state changes.
// Output when the permission mode changes, etc.; the permissionMode field indicates the current mode.
//
// ```json
// {"type":"system","subtype":"status","status":null,"permissionMode":"plan","uuid":"xxx","session_id":"abc"}
// ```
type SystemStatusMessage struct {
	MessageBase
	Status         string         `json:"status,omitempty"` // Status (null or string)
	PermissionMode PermissionMode `json:"permissionMode"`   // Permission mode
	UUID           string         `json:"uuid"`             // Message UUID
	SessionID      string         `json:"session_id"`       // Session ID
}

// # assistant
// A message containing the model's response. Each content block is output as a separate assistant message.
// The content array contains blocks of type text, tool_use, or thinking.
//
// #### assistant(text)
//
// A text response content block. The text field contains the response text.
//
// ```json
// {"type":"assistant","message":{"content":[{"type":"text","text":"Hello!"}],"id":"msg_001","model":"claude-sonnet-4-5-20250929","role":"assistant","stop_reason":null,"stop_sequence":null,"type":"message","usage":{"input_tokens":10,"output_tokens":1}},"parent_tool_use_id":null,"session_id":"abc","uuid":"xxx"}
// ```
//
// #### assistant(tool_use)
//
// A tool use content block. The name field contains the tool name and the input field contains the parameters.
//
// ```json
// {"type":"assistant","message":{"content":[{"type":"tool_use","id":"toolu_001","name":"Bash","input":{"command":"echo hello"}}],"id":"msg_001","model":"claude-sonnet-4-5-20250929","role":"assistant","stop_reason":null,"stop_sequence":null,"type":"message","usage":{"input_tokens":10,"output_tokens":1}},"parent_tool_use_id":null,"session_id":"abc","uuid":"xxx"}
// ```
//
// #### assistant(thinking)
//
// An extended thinking content block. The thinking field contains the thinking content.
//
// ```json
// {"type":"assistant","message":{"content":[{"type":"thinking","thinking":"Let me think...","signature":""}],"id":"msg_001","model":"claude-sonnet-4-5-20250929","role":"assistant","stop_reason":null,"stop_sequence":null,"type":"message","usage":{"input_tokens":10,"output_tokens":1}},"parent_tool_use_id":null,"session_id":"abc","uuid":"xxx"}
// ```
type AssistantMessage struct {
	MessageBase
	Message         AssistantBody `json:"message"`
	ParentToolUseID string        `json:"parent_tool_use_id,omitempty"` // Parent tool use ID (null or string)
	SessionID       string        `json:"session_id"`                   // Session ID
	UUID            string        `json:"uuid"`                         // Message UUID
}

// # user
// A user message. Used for both input (stdin -> CLI) and output (CLI -> stdout).
// For input, content is a string; for output (tool execution results), content is a block array.
//
// ```json
// {"type":"user","message":{"role":"user","content":"say hello"}}
// ```
//
// #### user(tool_result)
//
// A message where the CLI reports tool execution results. content is an array of tool_result blocks.
// Each block has the corresponding tool_use ID in the tool_use_id field.
// parent_tool_use_id, session_id, uuid, and tool_use_result are included at the top level.
//
// ```json
// {"type":"user","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"toolu_001","content":"command output"}]},"parent_tool_use_id":null,"session_id":"abc","uuid":"xxx","tool_use_result":{}}
// ```
type UserTextMessage struct {
	MessageBase
	Message UserTextBody `json:"message"`
}

// UserToolResultMessage is a user message where the CLI reports tool execution results.
// See the # user section of UserTextMessage for documentation.
type UserToolResultMessage struct {
	MessageBase
	Message         UserToolResultBody `json:"message"`
	ParentToolUseID string             `json:"parent_tool_use_id,omitempty"` // Parent tool use ID (null or string)
	SessionID       string             `json:"session_id"`                   // Session ID
	UUID            string             `json:"uuid"`                         // Message UUID
	ToolUseResult   any                `json:"tool_use_result"`              // Tool execution result (map or string or null)
}

// # result/success
// A message indicating successful completion of a turn. Indicates that processing of one turn completed normally.
// The result field contains the last text block content.
// permission_denials is always present; when empty, it is an empty array [].
//
// ```json
// {"type":"result","subtype":"success","is_error":false,"duration_ms":55,"duration_api_ms":12,"num_turns":1,"result":"Hello!","stop_reason":null,"session_id":"abc","total_cost_usd":0.00055,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":"xxx"}
// ```
type ResultSuccessMessage struct {
	MessageBase
	IsError           bool               `json:"is_error"`              // true on error
	DurationMs        float64            `json:"duration_ms"`           // Total duration (ms)
	DurationApiMs     float64            `json:"duration_api_ms"`       // API duration (ms)
	NumTurns          float64            `json:"num_turns"`             // Number of turns
	Result            string             `json:"result"`                // Last text block content
	StopReason        string             `json:"stop_reason,omitempty"` // Stop reason (null or string)
	SessionID         string             `json:"session_id"`            // Session ID
	TotalCostUSD      float64            `json:"total_cost_usd"`        // Total cost (USD)
	Usage             map[string]any     `json:"usage"`                 // Token usage
	ModelUsage        map[string]any     `json:"modelUsage"`            // Per-model usage
	PermissionDenials []PermissionDenial `json:"permission_denials"`    // Permission denials (always present)
	UUID              string             `json:"uuid"`                  // Message UUID
}

// # result/error_during_execution
// A turn-ending message when the API returns an error.
// In addition to the same common fields as result/success, the errors array contains error message strings.
//
// ```json
// {"type":"result","subtype":"error_during_execution","is_error":false,"duration_ms":52,"duration_api_ms":18,"num_turns":1,"session_id":"abc","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":"xxx","errors":["error message"]}
// ```
type ResultErrorMessage struct {
	MessageBase
	IsError           bool               `json:"is_error"`           // Error flag
	DurationMs        float64            `json:"duration_ms"`        // Total duration (ms)
	DurationApiMs     float64            `json:"duration_api_ms"`    // API duration (ms)
	NumTurns          float64            `json:"num_turns"`          // Number of turns
	SessionID         string             `json:"session_id"`         // Session ID
	TotalCostUSD      float64            `json:"total_cost_usd"`     // Total cost (USD)
	Usage             map[string]any     `json:"usage"`              // Token usage
	ModelUsage        map[string]any     `json:"modelUsage"`         // Per-model usage
	PermissionDenials []PermissionDenial `json:"permission_denials"` // Permission denials (always present)
	UUID              string             `json:"uuid"`               // Message UUID
	Errors            []string           `json:"errors"`             // Error array
}

// ---------------------------------------------------------------------------
// Content block types
// ---------------------------------------------------------------------------

// TextBlock is a text response content block.
type TextBlock struct {
	ContentBlockBase
	Text string `json:"text"`
}

// ToolUseBlock is a tool use content block.
type ToolUseBlock struct {
	ContentBlockBase
	ID    string         `json:"id"`    // Tool use ID
	Name  string         `json:"name"`  // Tool name
	Input map[string]any `json:"input"` // Tool parameters
}

// ThinkingBlock is an extended thinking content block.
type ThinkingBlock struct {
	ContentBlockBase
	Thinking  string `json:"thinking"`
	Signature string `json:"signature"` // Signature (empty string)
}

// ToolResultBlock is a tool execution result content block.
type ToolResultBlock struct {
	ContentBlockBase
	ToolUseID string `json:"tool_use_id"`        // Corresponding tool use ID
	Content   any    `json:"content"`            // Execution result (array or string)
	IsError   bool   `json:"is_error,omitempty"` // true on error
}

// ---------------------------------------------------------------------------
// Other
// ---------------------------------------------------------------------------

// AssistantBody is the body of an assistant message.
type AssistantBody struct {
	Content      []IsContentBlock  `json:"content"`
	ID           string            `json:"id"`    // Message ID
	Model        string            `json:"model"` // Model name
	Role         MessageRole       `json:"role"`
	StopReason   string            `json:"stop_reason,omitempty"`   // Stop reason (null or string)
	StopSequence string            `json:"stop_sequence,omitempty"` // Stop sequence (null or string)
	BodyType     AssistantBodyType `json:"type"`                    // Always "message"
	Usage        map[string]any    `json:"usage"`                   // Token usage
}

// UserTextBody is the body of a user text message.
type UserTextBody struct {
	Role    MessageRole `json:"role"`
	Content string      `json:"content"`
}

// UserToolResultBody is the body of a user tool result message.
type UserToolResultBody struct {
	Role    MessageRole       `json:"role"` // Always "user"
	Content []ToolResultBlock `json:"content"`
}

// PermissionDenial holds information about a denied tool.
type PermissionDenial struct {
	ToolName  string         `json:"tool_name"`   // Tool name
	ToolUseID string         `json:"tool_use_id"` // Tool use ID
	ToolInput map[string]any `json:"tool_input"`  // Tool input parameters
}
