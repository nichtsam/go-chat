package connection

type (
	EventType     string
	OutgoingEvent struct {
		EventType EventType
		Payload   interface{}
	}
)

var Event = struct {
	WRITE_MESSAGE EventType
}{
	WRITE_MESSAGE: "write_message",
}
