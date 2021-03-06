package views

const (
	AlertLvlError   = "danger"
	AlertLvlWarning = "warning"
	AlertLvlInfo    = "info"
	AlertLvlSuccess = "success"

	// AlertMsgGeneric is displayed when any random error is
	// encountered by our backend.
	AlertMsgGeneric = "Something went wrong. Please try again, and contact us if the problem persists."
)

type Alert struct {
	Level   string
	Message string
}

type Data struct {
	Alert *Alert
	Yield interface{}
}

func (d *Data) SetAlert(err error) {
	if pErr, ok := err.(PublicError); ok {
		d.Alert = &Alert{
			Level:   AlertLvlError,
			Message: pErr.Public(),
		}
	} else {
		d.Alert = &Alert{
			Level:   AlertLvlError,
			Message: AlertMsgGeneric,
		}
	}
}

type PublicError interface {
	error
	Public() string
}

func (d *Data) AlertError(msg string) {
	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: msg,
	}
}
