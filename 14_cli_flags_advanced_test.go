package ccprotocol_test

import (
	"encoding/json"
	"strings"
	"testing"

	. "github.com/hrntknr/claudecodeprotocol"
	"github.com/hrntknr/claudecodeprotocol/utils"
)

// Custom system prompt via --system-prompt flag
func TestSystemPrompt(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		utils.TextResponse("Hello from custom system prompt."),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSessionWithFlags(t, stub.URL(),
		[]string{"--system-prompt", "You are a test bot"}, nil)
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "hello"},
	}))
	output := s.Read()

	// Verify normal init/result flow
	utils.AssertOutput(t, output,
		defaultInitPattern(),
		defaultResultPattern(func(m *ResultSuccessMessage) {
			m.Result = "Hello from custom system prompt."
		}),
	)

	// Verify the API request contains the custom system prompt
	reqs := stub.Requests()
	var found bool
	for _, req := range reqs {
		model, _ := req.Body["model"].(string)
		if strings.Contains(model, "haiku") {
			continue
		}
		system, _ := req.Body["system"]
		systemJSON, _ := json.Marshal(system)
		systemStr := string(systemJSON)
		if strings.Contains(systemStr, "You are a test bot") {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected API request system field to contain 'You are a test bot'")
	}
}

// Append to default system prompt via --append-system-prompt flag
func TestAppendSystemPrompt(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		utils.TextResponse("Hello with appended prompt."),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSessionWithFlags(t, stub.URL(),
		[]string{"--append-system-prompt", "EXTRA_MARKER_FOR_TEST"}, nil)
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "hello"},
	}))
	output := s.Read()

	// Verify normal init/result flow
	utils.AssertOutput(t, output,
		defaultInitPattern(),
		defaultResultPattern(func(m *ResultSuccessMessage) {
			m.Result = "Hello with appended prompt."
		}),
	)

	// Verify the API request contains the appended marker
	reqs := stub.Requests()
	var foundMarker bool
	for _, req := range reqs {
		model, _ := req.Body["model"].(string)
		if strings.Contains(model, "haiku") {
			continue
		}
		system, _ := req.Body["system"]
		systemJSON, _ := json.Marshal(system)
		systemStr := string(systemJSON)
		if strings.Contains(systemStr, "EXTRA_MARKER_FOR_TEST") {
			foundMarker = true
			break
		}
	}
	if !foundMarker {
		t.Error("expected API request system field to contain 'EXTRA_MARKER_FOR_TEST'")
	}
}

// Session ID override via --session-id flag
func TestSessionID(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		utils.TextResponse("Hello with custom session."),
	}}
	stub.Start()
	defer stub.Close()

	sessionID := "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
	s := utils.NewSessionWithFlags(t, stub.URL(),
		[]string{"--session-id", sessionID}, nil)
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "hello"},
	}))
	output := s.Read()

	// Verify the init message contains the specified session_id
	utils.AssertOutput(t, output,
		defaultInitPattern(func(m *SystemInitMessage) {
			m.SessionID = sessionID
		}),
		defaultResultPattern(func(m *ResultSuccessMessage) {
			m.Result = "Hello with custom session."
			m.SessionID = sessionID
		}),
	)
}
