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

// IsMessage はプロトコルメッセージ型を制約するインターフェース。
type IsMessage interface {
	isMessage()
}

// IsContentBlock はコンテンツブロック型を制約するインターフェース。
type IsContentBlock interface {
	isContentBlock()
}

// ---------------------------------------------------------------------------
// Base types
// ---------------------------------------------------------------------------

// MessageBase はすべてのメッセージに共通するフィールドを持つ。
type MessageBase struct {
	Type    MessageType    `json:"type"`
	Subtype MessageSubtype `json:"subtype,omitempty"`
}

func (MessageBase) isMessage() {}

// ContentBlockBase はすべてのコンテンツブロックに共通するフィールドを持つ。
type ContentBlockBase struct {
	Type ContentBlockType `json:"type"`
}

func (ContentBlockBase) isContentBlock() {}

// ---------------------------------------------------------------------------
// Message types
// ---------------------------------------------------------------------------

// # system/init
// CLI起動時に最初に出力されるメッセージ。セッションIDやバージョンなどの初期情報を含む。
// 必ずセッションの最初のメッセージとして出力され、tools フィールドに利用可能なツール一覧が含まれる。
// mcp_servers, model, slash_commands, agents, skills, plugins などの拡張情報も含む。
//
// ```json
// {"type":"system","subtype":"init","cwd":"/path","session_id":"abc","tools":["Bash","Read"],"mcp_servers":[],"model":"claude-opus-4-6","permissionMode":"bypassPermissions","slash_commands":[...],"apiKeySource":"none","claude_code_version":"2.1.37","output_style":"default","agents":[...],"skills":[...],"plugins":[],"uuid":"xxx"}
// ```
type SystemInitMessage struct {
	MessageBase
	CWD               string         `json:"cwd"`                 // 作業ディレクトリ
	SessionID         string         `json:"session_id"`          // セッションID
	Tools             []string       `json:"tools"`               // 利用可能ツール一覧
	MCPServers        []string       `json:"mcp_servers"`         // MCP サーバー一覧
	Model             string         `json:"model"`               // 使用モデル名
	PermissionMode    PermissionMode `json:"permissionMode"`      // パーミッションモード
	SlashCommands     []string       `json:"slash_commands"`      // スラッシュコマンド一覧
	APIKeySource      string         `json:"apiKeySource"`        // API キーのソース
	ClaudeCodeVersion string         `json:"claude_code_version"` // CLI バージョン
	OutputStyle       string         `json:"output_style"`        // 出力スタイル
	Agents            []string       `json:"agents"`              // エージェント一覧
	Skills            []string       `json:"skills"`              // スキル一覧
	Plugins           []string       `json:"plugins"`             // プラグイン一覧
	UUID              string         `json:"uuid"`                // メッセージ固有ID
}

// # system/status
// システム状態の変更を通知するメッセージ。
// パーミッションモードの変更時などに出力され、permissionMode フィールドで現在のモードが通知される。
//
// ```json
// {"type":"system","subtype":"status","status":null,"permissionMode":"plan","uuid":"xxx","session_id":"abc"}
// ```
type SystemStatusMessage struct {
	MessageBase
	Status         string         `json:"status,omitempty"` // ステータス (null or string)
	PermissionMode PermissionMode `json:"permissionMode"`   // パーミッションモード
	UUID           string         `json:"uuid"`             // メッセージ固有ID
	SessionID      string         `json:"session_id"`       // セッションID
}

// # assistant
// モデルの応答を含むメッセージ。各コンテンツブロックは個別の assistant メッセージとして出力される。
// content 配列には text, tool_use, thinking のいずれかのブロックが含まれる。
//
// #### assistant(text)
//
// テキスト応答のコンテンツブロック。text フィールドに応答テキストが入る。
//
// ```json
// {"type":"assistant","message":{"content":[{"type":"text","text":"Hello!"}],"id":"msg_001","model":"claude-sonnet-4-5-20250929","role":"assistant","stop_reason":null,"stop_sequence":null,"type":"message","usage":{"input_tokens":10,"output_tokens":1}},"parent_tool_use_id":null,"session_id":"abc","uuid":"xxx"}
// ```
//
// #### assistant(tool_use)
//
// ツール呼び出しのコンテンツブロック。name フィールドにツール名、input フィールドにパラメータが入る。
//
// ```json
// {"type":"assistant","message":{"content":[{"type":"tool_use","id":"toolu_001","name":"Bash","input":{"command":"echo hello"}}],"id":"msg_001","model":"claude-sonnet-4-5-20250929","role":"assistant","stop_reason":null,"stop_sequence":null,"type":"message","usage":{"input_tokens":10,"output_tokens":1}},"parent_tool_use_id":null,"session_id":"abc","uuid":"xxx"}
// ```
//
// #### assistant(thinking)
//
// 拡張思考のコンテンツブロック。thinking フィールドに思考内容が入る。
//
// ```json
// {"type":"assistant","message":{"content":[{"type":"thinking","thinking":"Let me think...","signature":""}],"id":"msg_001","model":"claude-sonnet-4-5-20250929","role":"assistant","stop_reason":null,"stop_sequence":null,"type":"message","usage":{"input_tokens":10,"output_tokens":1}},"parent_tool_use_id":null,"session_id":"abc","uuid":"xxx"}
// ```
type AssistantMessage struct {
	MessageBase
	Message         AssistantBody `json:"message"`
	ParentToolUseID string        `json:"parent_tool_use_id,omitempty"` // 親ツール呼び出しID (null or string)
	SessionID       string        `json:"session_id"`                   // セッションID
	UUID            string        `json:"uuid"`                         // メッセージ固有ID
}

// # user
// ユーザメッセージ。入力（stdin→CLI）と出力（CLI→stdout）の両方で使用される。
// 入力の場合は content が文字列、出力（ツール実行結果）の場合は content がブロック配列となる。
//
// ```json
// {"type":"user","message":{"role":"user","content":"say hello"}}
// ```
//
// #### user(tool_result)
//
// CLIがツール実行結果を報告するメッセージ。content は tool_result ブロックの配列。
// 各ブロックは対応する tool_use の ID を tool_use_id フィールドに持つ。
// parent_tool_use_id, session_id, uuid, tool_use_result がトップレベルに含まれる。
//
// ```json
// {"type":"user","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"toolu_001","content":"command output"}]},"parent_tool_use_id":null,"session_id":"abc","uuid":"xxx","tool_use_result":{}}
// ```
type UserTextMessage struct {
	MessageBase
	Message UserTextBody `json:"message"`
}

// UserToolResultMessage はCLIがツール実行結果を報告する user メッセージ。
// ドキュメントは UserTextMessage の # user セクションを参照。
type UserToolResultMessage struct {
	MessageBase
	Message         UserToolResultBody `json:"message"`
	ParentToolUseID string             `json:"parent_tool_use_id,omitempty"` // 親ツール呼び出しID (null or string)
	SessionID       string             `json:"session_id"`                   // セッションID
	UUID            string             `json:"uuid"`                         // メッセージ固有ID
	ToolUseResult   any                `json:"tool_use_result"`              // ツール実行結果 (map or string or null)
}

// # result/success
// ターンの正常完了を示すメッセージ。1ターンの処理が正常に終了したことを示す。
// result フィールドには最後のテキストブロックの内容が入る。
// permission_denials は常に存在し、空の場合は空配列 [] となる。
//
// ```json
// {"type":"result","subtype":"success","is_error":false,"duration_ms":55,"duration_api_ms":12,"num_turns":1,"result":"Hello!","stop_reason":null,"session_id":"abc","total_cost_usd":0.00055,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":"xxx"}
// ```
type ResultSuccessMessage struct {
	MessageBase
	IsError           bool               `json:"is_error"`              // エラー時 true
	DurationMs        float64            `json:"duration_ms"`           // 総所要時間 (ms)
	DurationApiMs     float64            `json:"duration_api_ms"`       // API所要時間 (ms)
	NumTurns          float64            `json:"num_turns"`             // ターン数
	Result            string             `json:"result"`                // 最後のテキストブロック内容
	StopReason        string             `json:"stop_reason,omitempty"` // 停止理由 (null or string)
	SessionID         string             `json:"session_id"`            // セッションID
	TotalCostUSD      float64            `json:"total_cost_usd"`        // 総コスト (USD)
	Usage             map[string]any     `json:"usage"`                 // トークン使用量
	ModelUsage        map[string]any     `json:"modelUsage"`            // モデル別使用量
	PermissionDenials []PermissionDenial `json:"permission_denials"`    // パーミッション拒否 (常に存在)
	UUID              string             `json:"uuid"`                  // メッセージ固有ID
}

// # result/error_during_execution
// APIがエラーを返した場合のターン終了メッセージ。
// result/success と同じ共通フィールドに加え、errors 配列にエラーメッセージ文字列が格納される。
//
// ```json
// {"type":"result","subtype":"error_during_execution","is_error":false,"duration_ms":52,"duration_api_ms":18,"num_turns":1,"session_id":"abc","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":"xxx","errors":["error message"]}
// ```
type ResultErrorMessage struct {
	MessageBase
	IsError           bool               `json:"is_error"`           // エラーフラグ
	DurationMs        float64            `json:"duration_ms"`        // 総所要時間 (ms)
	DurationApiMs     float64            `json:"duration_api_ms"`    // API所要時間 (ms)
	NumTurns          float64            `json:"num_turns"`          // ターン数
	SessionID         string             `json:"session_id"`         // セッションID
	TotalCostUSD      float64            `json:"total_cost_usd"`     // 総コスト (USD)
	Usage             map[string]any     `json:"usage"`              // トークン使用量
	ModelUsage        map[string]any     `json:"modelUsage"`         // モデル別使用量
	PermissionDenials []PermissionDenial `json:"permission_denials"` // パーミッション拒否 (常に存在)
	UUID              string             `json:"uuid"`               // メッセージ固有ID
	Errors            []string           `json:"errors"`             // エラー配列
}

// ---------------------------------------------------------------------------
// Content block types
// ---------------------------------------------------------------------------

// TextBlock はテキスト応答のコンテンツブロック。
type TextBlock struct {
	ContentBlockBase
	Text string `json:"text"`
}

// ToolUseBlock はツール呼び出しのコンテンツブロック。
type ToolUseBlock struct {
	ContentBlockBase
	ID    string         `json:"id"`    // ツール呼び出しID
	Name  string         `json:"name"`  // ツール名
	Input map[string]any `json:"input"` // ツールパラメータ
}

// ThinkingBlock は拡張思考のコンテンツブロック。
type ThinkingBlock struct {
	ContentBlockBase
	Thinking  string `json:"thinking"`
	Signature string `json:"signature"` // 署名 (空文字列)
}

// ToolResultBlock はツール実行結果のコンテンツブロック。
type ToolResultBlock struct {
	ContentBlockBase
	ToolUseID string `json:"tool_use_id"`        // 対応するツール呼び出しID
	Content   any    `json:"content"`            // 実行結果 (配列 or 文字列)
	IsError   bool   `json:"is_error,omitempty"` // エラー時 true
}

// ---------------------------------------------------------------------------
// Other
// ---------------------------------------------------------------------------

// AssistantBody は assistant メッセージの本体。
type AssistantBody struct {
	Content      []IsContentBlock  `json:"content"`
	ID           string            `json:"id"`    // メッセージID
	Model        string            `json:"model"` // モデル名
	Role         MessageRole       `json:"role"`
	StopReason   string            `json:"stop_reason,omitempty"`   // 停止理由 (null or string)
	StopSequence string            `json:"stop_sequence,omitempty"` // 停止シーケンス (null or string)
	BodyType     AssistantBodyType `json:"type"`                    // 常に "message"
	Usage        map[string]any    `json:"usage"`                   // トークン使用量
}

// UserTextBody は user テキストメッセージの本体。
type UserTextBody struct {
	Role    MessageRole `json:"role"`
	Content string      `json:"content"`
}

// UserToolResultBody は user ツール実行結果メッセージの本体。
type UserToolResultBody struct {
	Role    MessageRole       `json:"role"` // 常に "user"
	Content []ToolResultBlock `json:"content"`
}

// PermissionDenial は拒否されたツールの情報。
type PermissionDenial struct {
	ToolName  string         `json:"tool_name"`   // ツール名
	ToolUseID string         `json:"tool_use_id"` // ツール呼び出しID
	ToolInput map[string]any `json:"tool_input"`  // ツール入力パラメータ
}
