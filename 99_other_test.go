package ccprotocol_test

import (
	"testing"

	. "github.com/hrntknr/claudecodeprotocol"
	"github.com/hrntknr/claudecodeprotocol/utils"
)

// Observe CLI behavior when the API returns stop_reason:"stop_sequence"
// with a non-null stop_sequence value in the message_delta SSE event.
//
// Observed: The CLI does NOT pass through stop_reason or stop_sequence from
// the API's message_delta event. Both fields remain null in the assistant
// message and result message, even though the stub returns
// stop_reason:"stop_sequence" and stop_sequence:"###".
// The CLI uses stop_reason internally (e.g. to trigger tool execution for
// "tool_use") but does not expose it in the stream-json output.
func TestStopSequence(t *testing.T) {
	t.Parallel()
	stub := &utils.StubAPIServer{Responses: [][]utils.SSEEvent{
		utils.StopSequenceTextResponse("Hello", "###"),
	}}
	stub.Start()
	defer stub.Close()

	s := utils.NewSession(t, stub.URL())
	defer s.Close()

	s.Send(utils.MustJSON(UserTextMessage{
		MessageBase: MessageBase{Type: TypeUser},
		Message:     UserTextBody{Role: RoleUser, Content: "test stop sequence"},
	}))
	utils.AssertOutput(t, s.Read(),
		defaultInitPattern(),
		// The assistant message has stop_reason:null and stop_sequence:null
		// (not "stop_sequence" and "###" as returned by the stub API).
		defaultAssistantPattern(func(m *AssistantMessage) {
			m.Message.Content = []IsContentBlock{
				TextBlock{
					ContentBlockBase: ContentBlockBase{Type: BlockText},
					Text:             "Hello",
				},
			}
		}),
		// The result message also has stop_reason:null (StopReason omitted → zero value).
		defaultResultPattern(func(m *ResultSuccessMessage) { m.Result = "Hello" }).Assert("result"),
	)
}
