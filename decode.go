package ccprotocol

import (
	"encoding/json"
	"fmt"
)

// DecodeMessage decodes JSON into the correct concrete message type based on
// the "type" and "subtype" fields. It returns a pointer to the concrete struct.
func DecodeMessage(data []byte) (IsMessage, error) {
	var base MessageBase
	if err := json.Unmarshal(data, &base); err != nil {
		return nil, fmt.Errorf("decode message base: %w", err)
	}

	switch base.Type {
	case TypeSystem:
		return decodeSystemMessage(data, base.Subtype)
	case TypeAssistant:
		var m AssistantMessage
		if err := json.Unmarshal(data, &m); err != nil {
			return nil, fmt.Errorf("decode assistant message: %w", err)
		}
		return &m, nil
	case TypeUser:
		return decodeUserMessage(data)
	case TypeResult:
		return decodeResultMessage(data, base.Subtype)
	case TypeStreamEvent:
		var m StreamEventMessage
		if err := json.Unmarshal(data, &m); err != nil {
			return nil, fmt.Errorf("decode stream_event message: %w", err)
		}
		return &m, nil
	case TypeControlRequest:
		var m ControlRequestMessage
		if err := json.Unmarshal(data, &m); err != nil {
			return nil, fmt.Errorf("decode control_request message: %w", err)
		}
		return &m, nil
	case TypeControlResponse:
		var m ControlResponseMessage
		if err := json.Unmarshal(data, &m); err != nil {
			return nil, fmt.Errorf("decode control_response message: %w", err)
		}
		return &m, nil
	default:
		return nil, fmt.Errorf("unknown message type: %q", base.Type)
	}
}

// decodeSystemMessage dispatches on the subtype for system messages.
func decodeSystemMessage(data []byte, subtype MessageSubtype) (IsMessage, error) {
	switch subtype {
	case SubtypeInit:
		var m SystemInitMessage
		if err := json.Unmarshal(data, &m); err != nil {
			return nil, fmt.Errorf("decode system/init message: %w", err)
		}
		return &m, nil
	case SubtypeStatus:
		var m SystemStatusMessage
		if err := json.Unmarshal(data, &m); err != nil {
			return nil, fmt.Errorf("decode system/status message: %w", err)
		}
		return &m, nil
	default:
		return nil, fmt.Errorf("unknown system subtype: %q", subtype)
	}
}

// decodeUserMessage distinguishes between UserReplayMessage, UserToolResultMessage,
// and UserTextMessage by inspecting the raw JSON fields.
func decodeUserMessage(data []byte) (IsMessage, error) {
	// Peek at discriminating fields without fully unmarshaling.
	var peek struct {
		IsReplay bool            `json:"isReplay"`
		Message  json.RawMessage `json:"message"`
	}
	if err := json.Unmarshal(data, &peek); err != nil {
		return nil, fmt.Errorf("decode user message (peek): %w", err)
	}

	if peek.IsReplay {
		var m UserReplayMessage
		if err := json.Unmarshal(data, &m); err != nil {
			return nil, fmt.Errorf("decode user replay message: %w", err)
		}
		return &m, nil
	}

	// Check if message.content is a JSON array (tool_result) or string (text).
	var msgPeek struct {
		Content json.RawMessage `json:"content"`
	}
	if err := json.Unmarshal(peek.Message, &msgPeek); err != nil {
		return nil, fmt.Errorf("decode user message body (peek): %w", err)
	}

	if len(msgPeek.Content) > 0 && msgPeek.Content[0] == '[' {
		var m UserToolResultMessage
		if err := json.Unmarshal(data, &m); err != nil {
			return nil, fmt.Errorf("decode user tool_result message: %w", err)
		}
		return &m, nil
	}

	var m UserTextMessage
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("decode user text message: %w", err)
	}
	return &m, nil
}

// decodeResultMessage dispatches on the subtype for result messages.
func decodeResultMessage(data []byte, subtype MessageSubtype) (IsMessage, error) {
	switch subtype {
	case SubtypeSuccess:
		var m ResultSuccessMessage
		if err := json.Unmarshal(data, &m); err != nil {
			return nil, fmt.Errorf("decode result/success message: %w", err)
		}
		return &m, nil
	case SubtypeErrorDuringExecution:
		var m ResultErrorMessage
		if err := json.Unmarshal(data, &m); err != nil {
			return nil, fmt.Errorf("decode result/error message: %w", err)
		}
		return &m, nil
	case SubtypeErrorMaxTurns:
		var m ResultMaxTurnsMessage
		if err := json.Unmarshal(data, &m); err != nil {
			return nil, fmt.Errorf("decode result/error_max_turns message: %w", err)
		}
		return &m, nil
	default:
		return nil, fmt.Errorf("unknown result subtype: %q", subtype)
	}
}

// DecodeContentBlock decodes JSON into the correct content block type based on
// the "type" field. It returns a value (not a pointer).
func DecodeContentBlock(data []byte) (IsContentBlock, error) {
	var base ContentBlockBase
	if err := json.Unmarshal(data, &base); err != nil {
		return nil, fmt.Errorf("decode content block base: %w", err)
	}

	switch base.Type {
	case BlockText:
		var b TextBlock
		if err := json.Unmarshal(data, &b); err != nil {
			return nil, fmt.Errorf("decode text block: %w", err)
		}
		return b, nil
	case BlockToolUse:
		var b ToolUseBlock
		if err := json.Unmarshal(data, &b); err != nil {
			return nil, fmt.Errorf("decode tool_use block: %w", err)
		}
		return b, nil
	case BlockThinking:
		var b ThinkingBlock
		if err := json.Unmarshal(data, &b); err != nil {
			return nil, fmt.Errorf("decode thinking block: %w", err)
		}
		return b, nil
	case BlockToolResult:
		var b ToolResultBlock
		if err := json.Unmarshal(data, &b); err != nil {
			return nil, fmt.Errorf("decode tool_result block: %w", err)
		}
		return b, nil
	default:
		return nil, fmt.Errorf("unknown content block type: %q", base.Type)
	}
}

// UnmarshalJSON implements json.Unmarshaler for AssistantBody, handling the
// polymorphic Content field via DecodeContentBlock.
func (b *AssistantBody) UnmarshalJSON(data []byte) error {
	// Use an alias to avoid infinite recursion.
	type Alias AssistantBody
	var raw struct {
		Alias
		Content []json.RawMessage `json:"content"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("decode assistant body: %w", err)
	}

	*b = AssistantBody(raw.Alias)
	b.Content = make([]IsContentBlock, len(raw.Content))
	for i, c := range raw.Content {
		block, err := DecodeContentBlock(c)
		if err != nil {
			return fmt.Errorf("decode assistant body content[%d]: %w", i, err)
		}
		b.Content[i] = block
	}
	return nil
}
