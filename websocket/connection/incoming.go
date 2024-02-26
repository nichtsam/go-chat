package connection

type (
	IntentType     string
	IncomingIntent struct {
		IntentType IntentType `json:"intent"`
		Raw        []byte     `json:"-"`
	}
)

var Intent = struct {
	SEND_MESSAGE,
	SECOND_INTENT,
	THIRE_INTENT IntentType
}{
	SEND_MESSAGE: "send_message",
}
