package connection

type (
	EventType     string
	OutgoingEvent struct {
		EventType EventType
		Payload   interface{}
	}
)

var Event = struct {
	WRITE_MESSAGE,
	SECOND_EVENT,
	THIRE_EVENT EventType
}{
	WRITE_MESSAGE: "write_message",
}
