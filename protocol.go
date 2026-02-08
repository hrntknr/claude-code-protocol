// Package ccprotocol documents the Claude Code CLI stream-json protocol.
//
// Claude Code CLI は --output-format stream-json で起動すると、
// stdout に JSONL 形式でメッセージを出力し、
// --input-format stream-json で stdin から JSONL メッセージを受け付ける。
//
// 各メッセージは [MessageType] で大分類され、
// [Subtype] や [ContentBlockType] でさらに細分化される。
//
// テストでは [MustJSON] と [Message] を使ってアサーションパターンを構築する。
package ccprotocol

import "encoding/json"

//go:generate go run ./cmd/gendoc

// ---------------------------------------------------------------------------
// Enum: MessageType — メッセージの大分類
// ---------------------------------------------------------------------------

// MessageType はプロトコルメッセージの type フィールドの値。
type MessageType string

const (
	// TypeSystem はCLIのライフサイクルイベント（初期化、状態変更）。
	TypeSystem MessageType = "system"
	// TypeAssistant はモデルからの応答メッセージ。
	TypeAssistant MessageType = "assistant"
	// TypeUser はユーザ入力またはツール実行結果。
	TypeUser MessageType = "user"
	// TypeResult はターンの完了通知。
	TypeResult MessageType = "result"
)

// ---------------------------------------------------------------------------
// Enum: Subtype — メッセージのサブタイプ
// ---------------------------------------------------------------------------

// Subtype はプロトコルメッセージの subtype フィールドの値。
// system メッセージと result メッセージで使用される。
type Subtype string

const (
	// SubtypeInit はCLI起動時の初期化メッセージ (type=system)。
	SubtypeInit Subtype = "init"
	// SubtypeStatus はシステム状態の変更通知 (type=system)。
	SubtypeStatus Subtype = "status"
	// SubtypeSuccess はターンの正常完了 (type=result)。
	SubtypeSuccess Subtype = "success"
	// SubtypeErrorDuringExecution はAPIエラーによるターン終了 (type=result)。
	SubtypeErrorDuringExecution Subtype = "error_during_execution"
)

// ---------------------------------------------------------------------------
// Enum: ContentBlockType — コンテンツブロックの種別
// ---------------------------------------------------------------------------

// ContentBlockType はコンテンツブロックの type フィールドの値。
// assistant メッセージと user メッセージの content 配列内で使用される。
type ContentBlockType string

const (
	// BlockText はテキスト応答ブロック。モデルの自然言語応答を含む。
	BlockText ContentBlockType = "text"
	// BlockToolUse はツール呼び出しブロック。モデルがツール実行を要求する。
	BlockToolUse ContentBlockType = "tool_use"
	// BlockThinking は拡張思考ブロック。モデルの内部推論過程を含む。
	BlockThinking ContentBlockType = "thinking"
	// BlockToolResult はツール実行結果ブロック。CLIがツール実行結果を報告する。
	BlockToolResult ContentBlockType = "tool_result"
)

// ---------------------------------------------------------------------------
// Schema: Message — プロトコルメッセージの統合表現
// ---------------------------------------------------------------------------

// Message はプロトコルメッセージの統合表現。
// スキーマドキュメントの定義ソースであり、テストのアサーションパターンにも使用する。
// 全フィールドに omitempty を付与しており、[MustJSON] でJSON化すると
// 設定されたフィールドのみが出力される（部分一致パターンとして機能する）。
//
// アサーションパターンの構築には NewMessageXxx コンストラクタ関数を使用する。
// 各コンストラクタにメッセージ種別ごとのドキュメントが記載されている。
type Message struct {
	Type              MessageType        `json:"type,omitempty"`               // メッセージ種別
	Subtype           Subtype            `json:"subtype,omitempty"`            // サブタイプ (system, result で使用)
	Message           *MessageBody       `json:"message,omitempty"`            // メッセージ本体 (assistant, user で使用)
	Result            string             `json:"result,omitempty"`             // 最後のテキストブロックの内容 (result で使用)
	IsError           bool               `json:"is_error,omitempty"`           // エラー時 true (result で使用)
	PermissionMode    string             `json:"permissionMode,omitempty"`     // パーミッションモード (system/status で使用)
	PermissionDenials []PermissionDenial `json:"permission_denials,omitempty"` // パーミッション拒否されたツール (result で使用)
	Errors            []ResultError      `json:"errors,omitempty"`             // APIエラー配列 (result で使用)
}

// MessageBody は assistant/user メッセージの内部構造。
type MessageBody struct {
	Role    string         `json:"role,omitempty"`    // "assistant" または "user"
	Content []ContentBlock `json:"content,omitempty"` // コンテンツブロックの配列
}

// ---------------------------------------------------------------------------
// Schema: ContentBlock — コンテンツブロックの統合表現
// ---------------------------------------------------------------------------

// ContentBlock はメッセージ内のコンテンツブロックの統合表現。
// [ContentBlockType] に応じて使用されるフィールドが異なる。
//
// # type=text
//
// テキスト応答ブロック。モデルの自然言語応答テキストを含む。
//
// JSON例:
//
//	{"type":"text","text":"The answer is 42."}
//
// # type=tool_use
//
// ツール呼び出しブロック。モデルがツールの実行を要求する際に出力される。
//
// JSON例:
//
//	{"type":"tool_use","id":"toolu_001","name":"Bash","input":{"command":"echo hello","description":"Print hello"}}
//
// 観測事項:
//   - --dangerously-skip-permissions なしの場合、パーミッション拒否で tool_result に is_error:true が返る。
//
// # type=thinking
//
// 拡張思考ブロック。モデルの内部推論過程を含む。
//
// JSON例:
//
//	{"type":"thinking","thinking":"Let me think about this step by step..."}
//
// 観測事項:
//   - thinking ブロックは text ブロックより前に、別の assistant メッセージとして出力される。
//
// # type=tool_result
//
// ツール実行結果ブロック。CLIがツール実行結果を報告する。
//
// JSON例:
//
//	{"type":"tool_result","tool_use_id":"toolu_001","content":"result text","is_error":false}
type ContentBlock struct {
	Type      ContentBlockType `json:"type,omitempty"`        // ブロック種別
	Text      string           `json:"text,omitempty"`        // 応答テキスト (text)
	Name      string           `json:"name,omitempty"`        // ツール名 (tool_use)
	ID        string           `json:"id,omitempty"`          // ツール呼び出しID (tool_use)
	Input     map[string]any   `json:"input,omitempty"`       // ツールパラメータ (tool_use)
	Thinking  string           `json:"thinking,omitempty"`    // 思考テキスト (thinking)
	ToolUseID string           `json:"tool_use_id,omitempty"` // 対応するツール呼び出しID (tool_result)
	Content   string           `json:"content,omitempty"`     // 実行結果テキスト (tool_result)
	IsError   bool             `json:"is_error,omitempty"`    // エラー時 true (tool_result)
}

// ---------------------------------------------------------------------------
// Schema: 補助型
// ---------------------------------------------------------------------------

// ResultError は result メッセージ内のエラー詳細。
type ResultError struct {
	Type    string `json:"type,omitempty"`    // エラー種別 (例: "overloaded_error")
	Message string `json:"message,omitempty"` // エラーメッセージ
}

// PermissionDenial はパーミッション拒否されたツールの情報。
type PermissionDenial struct {
	ToolName string `json:"tool_name,omitempty"` // 拒否されたツール名
}

// ---------------------------------------------------------------------------
// ヘルパー関数
// ---------------------------------------------------------------------------

// MustJSON は v をJSON文字列に変換する。
// テストのアサーションパターン構築に使用する。
//
// 使用例:
//
//	MustJSON(NewMessageSystemInit())
//	// => `{"type":"system","subtype":"init"}`
func MustJSON(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic("MustJSON: " + err.Error())
	}
	return string(b)
}

// ---------------------------------------------------------------------------
// コンストラクタ関数 — アサーションパターンの構築
// ---------------------------------------------------------------------------

// NewMessageSystemInit は system/init メッセージのアサーションパターンを生成する。
// # system/init
// CLI起動時に最初に出力されるメッセージ。セッションIDやバージョンなどの初期情報を含む。
// 必ずセッションの最初のメッセージとして出力され、tools フィールドに利用可能なツール一覧が含まれる。
// Example:
// ```json
//	{"type":"system","subtype":"init","apiKeyStatus":"valid","cwd":"/path","sessionId":"abc","version":"1.0.0","tools":["Bash","Read"]}
// ```
func NewMessageSystemInit() Message {
	return Message{Type: TypeSystem, Subtype: SubtypeInit}
}

// NewMessageSystemStatus は system/status メッセージのアサーションパターンを生成する。
// # system/status
// システム状態の変更を通知するメッセージ。
// EnterPlanMode ツール実行直後に出力され、PermissionMode フィールドで現在のモードが通知される。
// Example:
// ```json
//	{"type":"system","subtype":"status","permissionMode":"plan"}
// ```
func NewMessageSystemStatus(permissionMode string) Message {
	return Message{Type: TypeSystem, Subtype: SubtypeStatus, PermissionMode: permissionMode}
}

// NewMessageAssistantText はテキスト応答の assistant メッセージのアサーションパターンを生成する。
// # assistant/text
// モデルのテキスト応答を含むメッセージ。
// 各コンテンツブロックは個別の assistant メッセージとして出力され、
// 複数テキストブロックがある場合はブロックごとに別メッセージになる。
// text が空文字列の場合、Text フィールドは省略される（部分一致用）。
// Example:
// ```json
//	{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"Hello!"}]}}
// ```
func NewMessageAssistantText(text string) Message {
	return Message{Type: TypeAssistant, Message: &MessageBody{Content: []ContentBlock{{Type: BlockText, Text: text}}}}
}

// NewMessageAssistantToolUse はツール呼び出しの assistant メッセージのアサーションパターンを生成する。
// # assistant/tool_use
// モデルがツール実行を要求するメッセージ。
// 各コンテンツブロックは個別の assistant メッセージとして出力される。
// Example:
// ```json
//	{"type":"assistant","message":{"role":"assistant","content":[{"type":"tool_use","id":"toolu_001","name":"Bash","input":{"command":"echo hello"}}]}}
// ```
func NewMessageAssistantToolUse(name string) Message {
	return Message{Type: TypeAssistant, Message: &MessageBody{Content: []ContentBlock{{Type: BlockToolUse, Name: name}}}}
}

// NewMessageAssistantThinking は拡張思考ブロックの assistant メッセージのアサーションパターンを生成する。
// # assistant/thinking
// モデルの内部推論過程を含むメッセージ。
// thinking ブロックは text ブロックより前に、別の assistant メッセージとして順序通り出力される。
// Example:
// ```json
//	{"type":"assistant","message":{"role":"assistant","content":[{"type":"thinking","thinking":"Let me think..."}]}}
// ```
func NewMessageAssistantThinking(thinking string) Message {
	return Message{Type: TypeAssistant, Message: &MessageBody{Content: []ContentBlock{{Type: BlockThinking, Thinking: thinking}}}}
}

// NewMessageUserToolResult はツール実行結果の user メッセージのアサーションパターンを生成する。
// # user/tool_result
// CLIがツール実行結果を報告する user メッセージ。stdout に出力される。
// ContentBlock の ID に対応する tool_use_id を持つ。
// Example:
// ```json
//	{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"toolu_001","content":"command output"}]}}
// ```
func NewMessageUserToolResult() Message {
	return Message{Type: TypeUser, Message: &MessageBody{Content: []ContentBlock{{Type: BlockToolResult}}}}
}

// NewMessageUserToolResultError はエラーのツール実行結果の user メッセージのアサーションパターンを生成する。
// # user/tool_result, is_error=true
// ツール実行失敗またはパーミッション拒否時に出力されるエラー付きツール実行結果。
// Example:
// ```json
//	{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"toolu_001","is_error":true,"content":"Error: file not found"}]}}
// ```
func NewMessageUserToolResultError() Message {
	return Message{Type: TypeUser, Message: &MessageBody{Content: []ContentBlock{{Type: BlockToolResult, IsError: true}}}}
}

// NewMessageResultSuccess は正常完了の result メッセージのアサーションパターンを生成する。
// # result/success
// ターンの正常完了を示すメッセージ。Read() はこのメッセージ受信時に読み取りを終了する。
// Result フィールドには最後のテキストブロックの内容が入る。
// result が空文字列の場合、Result フィールドは省略される（部分一致用）。
// Example:
// ```json
//	{"type":"result","subtype":"success","result":"Hello!","is_error":false,"duration_ms":1234,"duration_api_ms":1000}
// ```
func NewMessageResultSuccess(result string) Message {
	return Message{Type: TypeResult, Subtype: SubtypeSuccess, Result: result}
}

// NewMessageResultSuccessIsError は is_error:true の result/success のアサーションパターンを生成する。
// # result/success, is_error=true
// max_tokens による応答打ち切り時に出力される。Subtype は "success" だが IsError が true になる。
// Example:
// ```json
//	{"type":"result","subtype":"success","is_error":true,"result":"truncated response..."}
// ```
func NewMessageResultSuccessIsError() Message {
	return Message{Type: TypeResult, Subtype: SubtypeSuccess, IsError: true}
}

// NewMessageResultSuccessWithDenials は permission_denials 付きの result/success のアサーションパターンを生成する。
// # result/success, permission_denials
// 非インタラクティブモードでパーミッション拒否されたツールがある場合に出力される。
// Example:
// ```json
//	{"type":"result","subtype":"success","permission_denials":[{"tool_name":"AskUserQuestion"}]}
// ```
func NewMessageResultSuccessWithDenials(denials ...PermissionDenial) Message {
	return Message{Type: TypeResult, Subtype: SubtypeSuccess, PermissionDenials: denials}
}

// NewMessageResultErrorDuringExecution はAPIエラーの result メッセージのアサーションパターンを生成する。
// # result/error_during_execution
// APIがSSEエラーイベントを返した場合のターン終了メッセージ。
// assistant メッセージは出力されず、Errors 配列にエラー種別とメッセージが格納される。
// Example:
// ```json
//	{"type":"result","subtype":"error_during_execution","errors":[{"type":"overloaded_error","message":"Overloaded"}]}
// ```
func NewMessageResultErrorDuringExecution() Message {
	return Message{Type: TypeResult, Subtype: SubtypeErrorDuringExecution}
}
