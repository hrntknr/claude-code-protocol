// > Requires: `CLAUDE_CODE_ENABLE_TASKS=1`
package ccprotocol_test

import (
	"encoding/json"
	"testing"

	. "github.com/hrntknr/claudecodeprotocol"
	"github.com/hrntknr/claudecodeprotocol/utils"
)

const enableTasksEnv = "CLAUDE_CODE_ENABLE_TASKS=1"

// Task management tools replace TodoWrite
func TestTaskToolsEnabled(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		utils.TextResponse("OK"),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSessionWithEnv(t, stub.URL(), []string{enableTasksEnv})
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "hello"},
	}))
	output := s.Read()

	// Observed: With CLAUDE_CODE_ENABLE_TASKS=1, the tools array in system/init
	// should include task management tools.
	var initFound bool
	for _, msg := range output {
		var m map[string]any
		if err := json.Unmarshal(msg, &m); err != nil {
			continue
		}
		if m["type"] == "system" && m["subtype"] == "init" {
			initFound = true
			tools, ok := m["tools"].([]any)
			if !ok {
				t.Fatal("tools field is not an array")
			}
			toolNames := make(map[string]bool)
			for _, tool := range tools {
				if name, ok := tool.(string); ok {
					toolNames[name] = true
				}
			}
			for _, expected := range []string{"TaskCreate", "TaskUpdate", "TaskList", "TaskGet"} {
				if !toolNames[expected] {
					t.Errorf("expected %q in tools list when CLAUDE_CODE_ENABLE_TASKS=1", expected)
				}
			}
			t.Logf("tools list (%d items): %v", len(tools), tools)
			break
		}
	}
	if !initFound {
		t.Fatal("system/init message not found")
	}
}

// TaskCreate tool invocation
func TestToolUseTaskCreate(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		utils.ToolUseResponse("toolu_tc_001", "TaskCreate", map[string]any{
			"subject":     "Test task",
			"description": "A test task for protocol observation",
			"activeForm":  "Creating test task",
		}),
		utils.TextResponse("Task created."),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSessionWithEnv(t, stub.URL(), []string{enableTasksEnv})
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "create a task"},
	}))
	utils.AssertOutput(t, s.Read(),
		defaultInitPattern(),
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				ToolUseBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockToolUse},
					ID:               utils.AnyString,
					Name:             "TaskCreate",
					Input:            utils.AnyMap,
				},
			}
		}),
		defaultUserToolResultPattern(func(m *UserToolResultMessage) {
			m.Message.Content = []ToolResultBlock{{
				ContentBlockBase: ContentBlockBase{Type: BlockToolResult},
				ToolUseID:        utils.AnyString,
				Content:          utils.AnyString,
			}}
		}),
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				TextBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockText},
					Text:             "Task created.",
				},
			}
		}),
		defaultResultPattern(func(m *ResultSuccessMessage) {
			m.Result = "Task created."
		}),
	)
}

// TaskList tool invocation
func TestToolUseTaskList(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		utils.ToolUseResponse("toolu_tl_001", "TaskList", map[string]any{}),
		utils.TextResponse("No tasks found."),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSessionWithEnv(t, stub.URL(), []string{enableTasksEnv})
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "list all tasks"},
	}))
	utils.AssertOutput(t, s.Read(),
		defaultInitPattern(),
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				ToolUseBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockToolUse},
					ID:               utils.AnyString,
					Name:             "TaskList",
					Input:            utils.AnyMap,
				},
			}
		}),
		defaultUserToolResultPattern(func(m *UserToolResultMessage) {
			m.Message.Content = []ToolResultBlock{{
				ContentBlockBase: ContentBlockBase{Type: BlockToolResult},
				ToolUseID:        utils.AnyString,
				Content:          utils.AnyString,
			}}
		}),
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				TextBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockText},
					Text:             "No tasks found.",
				},
			}
		}),
		defaultResultPattern(func(m *ResultSuccessMessage) {
			m.Result = "No tasks found."
		}),
	)
}

// TaskGet tool invocation
func TestToolUseTaskGet(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		utils.ToolUseResponse("toolu_tg_001", "TaskGet", map[string]any{
			"taskId": "1",
		}),
		utils.TextResponse("Task details retrieved."),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSessionWithEnv(t, stub.URL(), []string{enableTasksEnv})
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "get task 1"},
	}))
	utils.AssertOutput(t, s.Read(),
		defaultInitPattern(),
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				ToolUseBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockToolUse},
					ID:               utils.AnyString,
					Name:             "TaskGet",
					Input:            utils.AnyMap,
				},
			}
		}),
		defaultUserToolResultPattern(func(m *UserToolResultMessage) {
			m.Message.Content = []ToolResultBlock{{
				ContentBlockBase: ContentBlockBase{Type: BlockToolResult},
				ToolUseID:        utils.AnyString,
				Content:          utils.AnyString,
			}}
		}),
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				TextBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockText},
					Text:             "Task details retrieved.",
				},
			}
		}),
		defaultResultPattern(func(m *ResultSuccessMessage) {
			m.Result = "Task details retrieved."
		}),
	)
}

// TaskUpdate tool invocation
func TestToolUseTaskUpdate(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		utils.ToolUseResponse("toolu_tu_001", "TaskUpdate", map[string]any{
			"taskId": "1",
			"status": "completed",
		}),
		utils.TextResponse("Task updated."),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSessionWithEnv(t, stub.URL(), []string{enableTasksEnv})
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "mark task 1 as completed"},
	}))
	utils.AssertOutput(t, s.Read(),
		defaultInitPattern(),
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				ToolUseBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockToolUse},
					ID:               utils.AnyString,
					Name:             "TaskUpdate",
					Input:            utils.AnyMap,
				},
			}
		}),
		defaultUserToolResultPattern(func(m *UserToolResultMessage) {
			m.Message.Content = []ToolResultBlock{{
				ContentBlockBase: ContentBlockBase{Type: BlockToolResult},
				ToolUseID:        utils.AnyString,
				Content:          utils.AnyString,
			}}
		}),
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				TextBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockText},
					Text:             "Task updated.",
				},
			}
		}),
		defaultResultPattern(func(m *ResultSuccessMessage) {
			m.Result = "Task updated."
		}),
	)
}
