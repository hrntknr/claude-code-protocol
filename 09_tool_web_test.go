package ccprotocol_test

import (
	"testing"

	. "github.com/hrntknr/claudecodeprotocol"
	"github.com/hrntknr/claudecodeprotocol/utils"
)

// URL fetching behavior of the WebFetch tool
func TestToolUseWebFetch(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{
		StaticPages: map[string]string{
			"/static/test-page": "<html><body><h1>Test Page</h1><p>Content for WebFetch test.</p></body></html>",
		},
		Responses: [][]utils.SSEEvent{
			// Request 1: WebFetch tool_use
			utils.ToolUseResponse("toolu_wf_001", "WebFetch", map[string]any{
				"url":    "", // placeholder, will be set after Start()
				"prompt": "What does the page say?",
			}),
			// Final text
			utils.TextResponse("The page contains a heading: Test Page"),
		},
	}
	stub.Start()
	defer stub.Close()

	// Update the WebFetch URL to point to our stub's static page.
	// The ToolUseResponse is already built, so we rebuild with the correct URL.
	targetURL := stub.URL() + "/static/test-page"
	stub.Responses = [][]utils.SSEEvent{
		utils.ToolUseResponse("toolu_wf_001", "WebFetch", map[string]any{
			"url":    targetURL,
			"prompt": "What does the page say?",
		}),
		utils.TextResponse("The page contains a heading: Test Page"),
	}

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "fetch the test page"},
	}))
	// Observed: WebFetch upgrades HTTP to HTTPS, causing an SSL error when
	// hitting the plain HTTP stub server. The CLI emits the tool_use, then
	// a user tool_result with is_error:true containing the SSL error.
	// The API then returns the next response as final text.
	utils.AssertOutput(t, s.Read(),
		defaultInitPattern(),
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				ToolUseBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockToolUse},
					ID:               "toolu_stub_001",
					Name:             "WebFetch",
					Input:            map[string]any{"command": "echo hello", "description": "Example"},
				},
			}
		}).Ignore("message.content.*.id", "message.content.*.input"),
		defaultUserToolResultPattern(func(m *UserToolResultMessage) {
			m.Message.Content = []ToolResultBlock{{
				ContentBlockBase: ContentBlockBase{Type: BlockToolResult},
				ToolUseID:        "toolu_stub_001",
				Content:          "tool execution output",
				IsError:          true,
			}}
		}).Ignore("message.content.*.tool_use_id", "message.content.*.content"),
		defaultResultPattern(),
	)
}

// WebSearch tool behavior (may fail due to test environment restrictions)
func TestToolUseWebSearch(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		utils.ToolUseResponse("toolu_ws_001", "WebSearch", map[string]any{
			"query": "protocol test query",
		}),
		utils.TextResponse("Search completed."),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "search the web"},
	}))
	// Observed: WebSearch may return a tool_result with is_error:true if the
	// environment does not support web search (e.g., no valid API key or
	// not in US region). The protocol format is the same as other tool errors.
	output := s.Read()
	utils.AssertOutput(t, output,
		defaultInitPattern(),
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				ToolUseBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockToolUse},
					ID:               "toolu_stub_001",
					Name:             "WebSearch",
					Input:            map[string]any{"command": "echo hello", "description": "Example"},
				},
			}
		}).Ignore("message.content.*.id", "message.content.*.input"),
		defaultUserToolResultPattern(func(m *UserToolResultMessage) {
			m.Message.Content = []ToolResultBlock{{
				ContentBlockBase: ContentBlockBase{Type: BlockToolResult},
				ToolUseID:        "toolu_stub_001",
				Content:          "tool execution output",
			}}
		}).Ignore("message.content.*.tool_use_id", "message.content.*.content"),
		defaultResultPattern(),
	)
}
