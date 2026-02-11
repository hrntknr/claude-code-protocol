package ccprotocol_test

import (
	"testing"

	. "github.com/hrntknr/claudecodeprotocol"
)

// ---------------------------------------------------------------------------
// DecodeMessage — each message type
// ---------------------------------------------------------------------------

func TestDecodeMessage_SystemInit(t *testing.T) {
	data := []byte(`{"type":"system","subtype":"init","cwd":"/tmp","session_id":"s1","tools":["Bash"],"mcp_servers":[],"model":"claude-opus-4-6","permissionMode":"bypassPermissions","slash_commands":[],"apiKeySource":"none","claude_code_version":"2.1.37","output_style":"default","agents":[],"skills":[],"plugins":[],"uuid":"u1","fast_mode_state":"off"}`)

	msg, err := DecodeMessage(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m, ok := msg.(*SystemInitMessage)
	if !ok {
		t.Fatalf("expected *SystemInitMessage, got %T", msg)
	}
	if m.Type != TypeSystem {
		t.Errorf("Type = %q, want %q", m.Type, TypeSystem)
	}
	if m.Subtype != SubtypeInit {
		t.Errorf("Subtype = %q, want %q", m.Subtype, SubtypeInit)
	}
	if m.CWD != "/tmp" {
		t.Errorf("CWD = %q, want %q", m.CWD, "/tmp")
	}
	if m.SessionID != "s1" {
		t.Errorf("SessionID = %q, want %q", m.SessionID, "s1")
	}
	if m.Model != "claude-opus-4-6" {
		t.Errorf("Model = %q, want %q", m.Model, "claude-opus-4-6")
	}
	if m.PermissionMode != PermissionBypassPermissions {
		t.Errorf("PermissionMode = %q, want %q", m.PermissionMode, PermissionBypassPermissions)
	}
	if m.FastModeState != FastModeOff {
		t.Errorf("FastModeState = %q, want %q", m.FastModeState, FastModeOff)
	}
}

func TestDecodeMessage_SystemStatus(t *testing.T) {
	data := []byte(`{"type":"system","subtype":"status","status":null,"permissionMode":"plan","uuid":"u2","session_id":"s2"}`)

	msg, err := DecodeMessage(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m, ok := msg.(*SystemStatusMessage)
	if !ok {
		t.Fatalf("expected *SystemStatusMessage, got %T", msg)
	}
	if m.PermissionMode != PermissionPlan {
		t.Errorf("PermissionMode = %q, want %q", m.PermissionMode, PermissionPlan)
	}
	if m.UUID != "u2" {
		t.Errorf("UUID = %q, want %q", m.UUID, "u2")
	}
}

func TestDecodeMessage_AssistantText(t *testing.T) {
	data := []byte(`{"type":"assistant","message":{"content":[{"type":"text","text":"Hello!"}],"id":"msg_001","model":"claude-sonnet-4-5-20250929","role":"assistant","stop_reason":null,"stop_sequence":null,"type":"message","usage":{"input_tokens":10,"output_tokens":1}},"parent_tool_use_id":null,"session_id":"abc","uuid":"xxx"}`)

	msg, err := DecodeMessage(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m, ok := msg.(*AssistantMessage)
	if !ok {
		t.Fatalf("expected *AssistantMessage, got %T", msg)
	}
	if m.Type != TypeAssistant {
		t.Errorf("Type = %q, want %q", m.Type, TypeAssistant)
	}
	if m.Message.ID != "msg_001" {
		t.Errorf("Message.ID = %q, want %q", m.Message.ID, "msg_001")
	}
	if m.Message.Role != RoleAssistant {
		t.Errorf("Message.Role = %q, want %q", m.Message.Role, RoleAssistant)
	}
	if len(m.Message.Content) != 1 {
		t.Fatalf("len(Content) = %d, want 1", len(m.Message.Content))
	}
	tb, ok := m.Message.Content[0].(TextBlock)
	if !ok {
		t.Fatalf("Content[0] type = %T, want TextBlock", m.Message.Content[0])
	}
	if tb.Text != "Hello!" {
		t.Errorf("TextBlock.Text = %q, want %q", tb.Text, "Hello!")
	}
}

func TestDecodeMessage_AssistantToolUse(t *testing.T) {
	data := []byte(`{"type":"assistant","message":{"content":[{"type":"tool_use","id":"toolu_001","name":"Bash","input":{"command":"echo hello"}}],"id":"msg_001","model":"claude-sonnet-4-5-20250929","role":"assistant","stop_reason":null,"stop_sequence":null,"type":"message","usage":{"input_tokens":10,"output_tokens":1}},"parent_tool_use_id":null,"session_id":"abc","uuid":"xxx"}`)

	msg, err := DecodeMessage(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m, ok := msg.(*AssistantMessage)
	if !ok {
		t.Fatalf("expected *AssistantMessage, got %T", msg)
	}
	if len(m.Message.Content) != 1 {
		t.Fatalf("len(Content) = %d, want 1", len(m.Message.Content))
	}
	tu, ok := m.Message.Content[0].(ToolUseBlock)
	if !ok {
		t.Fatalf("Content[0] type = %T, want ToolUseBlock", m.Message.Content[0])
	}
	if tu.ID != "toolu_001" {
		t.Errorf("ToolUseBlock.ID = %q, want %q", tu.ID, "toolu_001")
	}
	if tu.Name != "Bash" {
		t.Errorf("ToolUseBlock.Name = %q, want %q", tu.Name, "Bash")
	}
	cmd, _ := tu.Input["command"].(string)
	if cmd != "echo hello" {
		t.Errorf("ToolUseBlock.Input[command] = %q, want %q", cmd, "echo hello")
	}
}

func TestDecodeMessage_AssistantThinking(t *testing.T) {
	data := []byte(`{"type":"assistant","message":{"content":[{"type":"thinking","thinking":"Let me think...","signature":""}],"id":"msg_001","model":"claude-sonnet-4-5-20250929","role":"assistant","stop_reason":null,"stop_sequence":null,"type":"message","usage":{"input_tokens":10,"output_tokens":1}},"parent_tool_use_id":null,"session_id":"abc","uuid":"xxx"}`)

	msg, err := DecodeMessage(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m, ok := msg.(*AssistantMessage)
	if !ok {
		t.Fatalf("expected *AssistantMessage, got %T", msg)
	}
	if len(m.Message.Content) != 1 {
		t.Fatalf("len(Content) = %d, want 1", len(m.Message.Content))
	}
	th, ok := m.Message.Content[0].(ThinkingBlock)
	if !ok {
		t.Fatalf("Content[0] type = %T, want ThinkingBlock", m.Message.Content[0])
	}
	if th.Thinking != "Let me think..." {
		t.Errorf("ThinkingBlock.Thinking = %q, want %q", th.Thinking, "Let me think...")
	}
}

func TestDecodeMessage_UserText(t *testing.T) {
	data := []byte(`{"type":"user","message":{"role":"user","content":"say hello"}}`)

	msg, err := DecodeMessage(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m, ok := msg.(*UserTextMessage)
	if !ok {
		t.Fatalf("expected *UserTextMessage, got %T", msg)
	}
	if m.Message.Role != RoleUser {
		t.Errorf("Message.Role = %q, want %q", m.Message.Role, RoleUser)
	}
	if m.Message.Content != "say hello" {
		t.Errorf("Message.Content = %q, want %q", m.Message.Content, "say hello")
	}
}

func TestDecodeMessage_UserToolResult(t *testing.T) {
	data := []byte(`{"type":"user","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"toolu_001","content":"command output"}]},"parent_tool_use_id":null,"session_id":"abc","uuid":"xxx","tool_use_result":{}}`)

	msg, err := DecodeMessage(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m, ok := msg.(*UserToolResultMessage)
	if !ok {
		t.Fatalf("expected *UserToolResultMessage, got %T", msg)
	}
	if m.SessionID != "abc" {
		t.Errorf("SessionID = %q, want %q", m.SessionID, "abc")
	}
	if len(m.Message.Content) != 1 {
		t.Fatalf("len(Content) = %d, want 1", len(m.Message.Content))
	}
	if m.Message.Content[0].ToolUseID != "toolu_001" {
		t.Errorf("ToolUseID = %q, want %q", m.Message.Content[0].ToolUseID, "toolu_001")
	}
}

func TestDecodeMessage_UserReplay(t *testing.T) {
	data := []byte(`{"type":"user","message":{"role":"user","content":"hello"},"session_id":"abc","parent_tool_use_id":null,"uuid":"xxx","isReplay":true}`)

	msg, err := DecodeMessage(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m, ok := msg.(*UserReplayMessage)
	if !ok {
		t.Fatalf("expected *UserReplayMessage, got %T", msg)
	}
	if !m.IsReplay {
		t.Error("IsReplay = false, want true")
	}
	if m.Message.Content != "hello" {
		t.Errorf("Message.Content = %q, want %q", m.Message.Content, "hello")
	}
}

func TestDecodeMessage_ResultSuccess(t *testing.T) {
	data := []byte(`{"type":"result","subtype":"success","is_error":false,"duration_ms":55,"duration_api_ms":12,"num_turns":1,"result":"Hello!","stop_reason":null,"session_id":"abc","total_cost_usd":0.00055,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":"xxx"}`)

	msg, err := DecodeMessage(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m, ok := msg.(*ResultSuccessMessage)
	if !ok {
		t.Fatalf("expected *ResultSuccessMessage, got %T", msg)
	}
	if m.Result != "Hello!" {
		t.Errorf("Result = %q, want %q", m.Result, "Hello!")
	}
	if m.DurationMs != 55 {
		t.Errorf("DurationMs = %v, want 55", m.DurationMs)
	}
	if m.NumTurns != 1 {
		t.Errorf("NumTurns = %v, want 1", m.NumTurns)
	}
}

func TestDecodeMessage_ResultError(t *testing.T) {
	data := []byte(`{"type":"result","subtype":"error_during_execution","is_error":false,"duration_ms":52,"duration_api_ms":18,"num_turns":1,"session_id":"abc","total_cost_usd":0,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":"xxx","errors":["error message"]}`)

	msg, err := DecodeMessage(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m, ok := msg.(*ResultErrorMessage)
	if !ok {
		t.Fatalf("expected *ResultErrorMessage, got %T", msg)
	}
	if len(m.Errors) != 1 || m.Errors[0] != "error message" {
		t.Errorf("Errors = %v, want [\"error message\"]", m.Errors)
	}
}

func TestDecodeMessage_ResultMaxTurns(t *testing.T) {
	data := []byte(`{"type":"result","subtype":"error_max_turns","is_error":false,"duration_ms":184,"duration_api_ms":30,"num_turns":2,"stop_reason":null,"session_id":"abc","total_cost_usd":0.00055,"usage":{},"modelUsage":{},"permission_denials":[],"uuid":"xxx","errors":[]}`)

	msg, err := DecodeMessage(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m, ok := msg.(*ResultMaxTurnsMessage)
	if !ok {
		t.Fatalf("expected *ResultMaxTurnsMessage, got %T", msg)
	}
	if m.Subtype != SubtypeErrorMaxTurns {
		t.Errorf("Subtype = %q, want %q", m.Subtype, SubtypeErrorMaxTurns)
	}
	if m.NumTurns != 2 {
		t.Errorf("NumTurns = %v, want 2", m.NumTurns)
	}
}

func TestDecodeMessage_StreamEvent(t *testing.T) {
	data := []byte(`{"type":"stream_event","event":{"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Hello"}},"session_id":"abc","parent_tool_use_id":null,"uuid":"xxx"}`)

	msg, err := DecodeMessage(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m, ok := msg.(*StreamEventMessage)
	if !ok {
		t.Fatalf("expected *StreamEventMessage, got %T", msg)
	}
	if m.Event["type"] != "content_block_delta" {
		t.Errorf("Event[type] = %v, want %q", m.Event["type"], "content_block_delta")
	}
	if m.SessionID != "abc" {
		t.Errorf("SessionID = %q, want %q", m.SessionID, "abc")
	}
}

func TestDecodeMessage_ControlRequest(t *testing.T) {
	data := []byte(`{"type":"control_request","request_id":"aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee","request":{"subtype":"set_permission_mode","mode":"plan"}}`)

	msg, err := DecodeMessage(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m, ok := msg.(*ControlRequestMessage)
	if !ok {
		t.Fatalf("expected *ControlRequestMessage, got %T", msg)
	}
	if m.RequestID != "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee" {
		t.Errorf("RequestID = %q, want %q", m.RequestID, "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee")
	}
	if m.Request.Subtype != ControlSetPermissionMode {
		t.Errorf("Request.Subtype = %q, want %q", m.Request.Subtype, ControlSetPermissionMode)
	}
	if m.Request.Mode != "plan" {
		t.Errorf("Request.Mode = %q, want %q", m.Request.Mode, "plan")
	}
}

func TestDecodeMessage_ControlResponse(t *testing.T) {
	data := []byte(`{"type":"control_response","response":{"subtype":"success","request_id":"aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee","response":{"mode":"plan"}}}`)

	msg, err := DecodeMessage(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m, ok := msg.(*ControlResponseMessage)
	if !ok {
		t.Fatalf("expected *ControlResponseMessage, got %T", msg)
	}
	if m.Response.Subtype != "success" {
		t.Errorf("Response.Subtype = %q, want %q", m.Response.Subtype, "success")
	}
	if m.Response.RequestID != "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee" {
		t.Errorf("Response.RequestID = %q, want %q", m.Response.RequestID, "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee")
	}
	// Response is any -> map[string]any
	resp, ok := m.Response.Response.(map[string]any)
	if !ok {
		t.Fatalf("Response.Response type = %T, want map[string]any", m.Response.Response)
	}
	if resp["mode"] != "plan" {
		t.Errorf("Response.Response[mode] = %v, want %q", resp["mode"], "plan")
	}
}

// ---------------------------------------------------------------------------
// DecodeContentBlock — each block type
// ---------------------------------------------------------------------------

func TestDecodeContentBlock_Text(t *testing.T) {
	data := []byte(`{"type":"text","text":"Hello!"}`)

	block, err := DecodeContentBlock(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tb, ok := block.(TextBlock)
	if !ok {
		t.Fatalf("expected TextBlock, got %T", block)
	}
	if tb.Text != "Hello!" {
		t.Errorf("Text = %q, want %q", tb.Text, "Hello!")
	}
	if tb.Type != BlockText {
		t.Errorf("Type = %q, want %q", tb.Type, BlockText)
	}
}

func TestDecodeContentBlock_ToolUse(t *testing.T) {
	data := []byte(`{"type":"tool_use","id":"toolu_001","name":"Bash","input":{"command":"echo hello"}}`)

	block, err := DecodeContentBlock(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tu, ok := block.(ToolUseBlock)
	if !ok {
		t.Fatalf("expected ToolUseBlock, got %T", block)
	}
	if tu.ID != "toolu_001" {
		t.Errorf("ID = %q, want %q", tu.ID, "toolu_001")
	}
	if tu.Name != "Bash" {
		t.Errorf("Name = %q, want %q", tu.Name, "Bash")
	}
}

func TestDecodeContentBlock_Thinking(t *testing.T) {
	data := []byte(`{"type":"thinking","thinking":"Let me think...","signature":"sig123"}`)

	block, err := DecodeContentBlock(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	th, ok := block.(ThinkingBlock)
	if !ok {
		t.Fatalf("expected ThinkingBlock, got %T", block)
	}
	if th.Thinking != "Let me think..." {
		t.Errorf("Thinking = %q, want %q", th.Thinking, "Let me think...")
	}
	if th.Signature != "sig123" {
		t.Errorf("Signature = %q, want %q", th.Signature, "sig123")
	}
}

func TestDecodeContentBlock_ToolResult(t *testing.T) {
	data := []byte(`{"type":"tool_result","tool_use_id":"toolu_001","content":"command output"}`)

	block, err := DecodeContentBlock(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tr, ok := block.(ToolResultBlock)
	if !ok {
		t.Fatalf("expected ToolResultBlock, got %T", block)
	}
	if tr.ToolUseID != "toolu_001" {
		t.Errorf("ToolUseID = %q, want %q", tr.ToolUseID, "toolu_001")
	}
	if tr.Content != "command output" {
		t.Errorf("Content = %v, want %q", tr.Content, "command output")
	}
}

// ---------------------------------------------------------------------------
// Error cases
// ---------------------------------------------------------------------------

func TestDecodeMessage_UnknownType(t *testing.T) {
	data := []byte(`{"type":"unknown_type"}`)
	_, err := DecodeMessage(data)
	if err == nil {
		t.Fatal("expected error for unknown type, got nil")
	}
}

func TestDecodeMessage_UnknownSystemSubtype(t *testing.T) {
	data := []byte(`{"type":"system","subtype":"unknown_sub"}`)
	_, err := DecodeMessage(data)
	if err == nil {
		t.Fatal("expected error for unknown system subtype, got nil")
	}
}

func TestDecodeMessage_UnknownResultSubtype(t *testing.T) {
	data := []byte(`{"type":"result","subtype":"unknown_sub"}`)
	_, err := DecodeMessage(data)
	if err == nil {
		t.Fatal("expected error for unknown result subtype, got nil")
	}
}

func TestDecodeContentBlock_UnknownType(t *testing.T) {
	data := []byte(`{"type":"unknown_block"}`)
	_, err := DecodeContentBlock(data)
	if err == nil {
		t.Fatal("expected error for unknown content block type, got nil")
	}
}

func TestDecodeMessage_InvalidJSON(t *testing.T) {
	data := []byte(`{invalid json}`)
	_, err := DecodeMessage(data)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestDecodeContentBlock_InvalidJSON(t *testing.T) {
	data := []byte(`{invalid json}`)
	_, err := DecodeContentBlock(data)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

// ---------------------------------------------------------------------------
// AssistantBody UnmarshalJSON — mixed content blocks
// ---------------------------------------------------------------------------

func TestDecodeMessage_AssistantMixedContent(t *testing.T) {
	data := []byte(`{"type":"assistant","message":{"content":[{"type":"thinking","thinking":"hmm","signature":""},{"type":"text","text":"answer"},{"type":"tool_use","id":"t1","name":"Read","input":{"file":"a.go"}}],"id":"msg_002","model":"claude-opus-4-6","role":"assistant","stop_reason":null,"stop_sequence":null,"type":"message","usage":{"input_tokens":5,"output_tokens":3}},"session_id":"s1","uuid":"u1"}`)

	msg, err := DecodeMessage(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m, ok := msg.(*AssistantMessage)
	if !ok {
		t.Fatalf("expected *AssistantMessage, got %T", msg)
	}
	if len(m.Message.Content) != 3 {
		t.Fatalf("len(Content) = %d, want 3", len(m.Message.Content))
	}

	if _, ok := m.Message.Content[0].(ThinkingBlock); !ok {
		t.Errorf("Content[0] type = %T, want ThinkingBlock", m.Message.Content[0])
	}
	if _, ok := m.Message.Content[1].(TextBlock); !ok {
		t.Errorf("Content[1] type = %T, want TextBlock", m.Message.Content[1])
	}
	if _, ok := m.Message.Content[2].(ToolUseBlock); !ok {
		t.Errorf("Content[2] type = %T, want ToolUseBlock", m.Message.Content[2])
	}
}
