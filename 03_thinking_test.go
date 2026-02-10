package ccprotocol_test

import (
	"testing"

	. "github.com/hrntknr/claudecodeprotocol"
	"github.com/hrntknr/claudecodeprotocol/utils"
)

// Response containing extended thinking blocks
func TestThinkingResponse(t *testing.T) {
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		utils.ThinkingResponse(
			"Let me think about this step by step...",
			"The answer is 42.",
		),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "what is the answer?"},
	}))
	// Observed: Thinking blocks ARE emitted as a separate assistant message
	// with content[0].type="thinking". Then the text block follows as another
	// assistant message. Result contains only the text.
	utils.AssertOutput(t, s.Read(),
		utils.MustJSON(defaultInitPattern()),
		utils.MustJSON(AssistantMessage{
			MessageBase: MessageBase{Type: TypeAssistant},
			Message: AssistantBody{
				Content: []IsContentBlock{
					ThinkingBlock{
						ContentBlockBase: ContentBlockBase{Type: BlockThinking},
						Thinking:         "Let me think about this step by step...",
						Signature:        "",
					},
				},
				ID:       utils.AnyString,
				Model:    utils.AnyString,
				Role:     RoleAssistant,
				BodyType: AssistantBodyTypeMessage,
				Usage:    utils.AnyMap,
			},
			SessionID: utils.AnyString,
			UUID:      utils.AnyString,
		}),
		utils.MustJSON(AssistantMessage{
			MessageBase: MessageBase{Type: TypeAssistant},
			Message: AssistantBody{
				Content: []IsContentBlock{
					TextBlock{
						ContentBlockBase: ContentBlockBase{Type: BlockText},
						Text:             "The answer is 42.",
					},
				},
				ID:       utils.AnyString,
				Model:    utils.AnyString,
				Role:     RoleAssistant,
				BodyType: AssistantBodyTypeMessage,
				Usage:    utils.AnyMap,
			},
			SessionID: utils.AnyString,
			UUID:      utils.AnyString,
		}),
		utils.MustJSON(ResultSuccessMessage{
			MessageBase:       MessageBase{Type: TypeResult, Subtype: SubtypeSuccess},
			IsError:           false,
			DurationMs:        utils.AnyNumber,
			DurationApiMs:     utils.AnyNumber,
			NumTurns:          utils.AnyNumber,
			Result:            "The answer is 42.",
			SessionID:         utils.AnyString,
			TotalCostUSD:      utils.AnyNumber,
			Usage:             utils.AnyMap,
			ModelUsage:        utils.AnyMap,
			PermissionDenials: []PermissionDenial{},
			UUID:              utils.AnyString,
		}),
	)
}

// Flow where tool use follows a thinking block
func TestThinkingWithToolUse(t *testing.T) {
	stub := &utils.StubAPIServer{Responses: utils.WithInit(
		// Request 1: thinking + tool_use
		utils.ThinkingAndToolUseResponse(
			"I need to run a command to check...",
			"toolu_think_001", "Bash", map[string]any{
				"command":     "echo thinking-tool-test",
				"description": "Test after thinking",
			},
		),
		// Request 2: final text
		utils.TextResponse("After thinking and running the command: thinking-tool-test"),
	)}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "think and then run a command"},
	}))
	utils.AssertOutput(t, s.Read(),
		utils.MustJSON(defaultInitPattern()),
		utils.MustJSON(AssistantMessage{
			MessageBase: MessageBase{Type: TypeAssistant},
			Message: AssistantBody{
				Content: []IsContentBlock{
					TextBlock{
						ContentBlockBase: ContentBlockBase{Type: BlockText},
						Text:             "After thinking and running the command: thinking-tool-test",
					},
				},
				ID:       utils.AnyString,
				Model:    utils.AnyString,
				Role:     RoleAssistant,
				BodyType: AssistantBodyTypeMessage,
				Usage:    utils.AnyMap,
			},
			SessionID: utils.AnyString,
			UUID:      utils.AnyString,
		}),
		utils.MustJSON(ResultSuccessMessage{
			MessageBase:       MessageBase{Type: TypeResult, Subtype: SubtypeSuccess},
			IsError:           false,
			DurationMs:        utils.AnyNumber,
			DurationApiMs:     utils.AnyNumber,
			NumTurns:          utils.AnyNumber,
			Result:            "After thinking and running the command: thinking-tool-test",
			SessionID:         utils.AnyString,
			TotalCostUSD:      utils.AnyNumber,
			Usage:             utils.AnyMap,
			ModelUsage:        utils.AnyMap,
			PermissionDenials: []PermissionDenial{},
			UUID:              utils.AnyString,
		}),
	)
}
