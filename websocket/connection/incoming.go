package connection

type (
	IntentType     string
	IncomingIntent struct {
		IntentType IntentType `json:"intent"`
		Raw        []byte     `json:"-"`
	}
)

var Intent = struct {
	SEND_MESSAGE IntentType
}{
	SEND_MESSAGE: "send_message",
}
