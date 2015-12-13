package notifiers

import (
	twilio "github.com/carlosdp/twiliogo"
	"github.com/wfernandes/homesec/logging"
)

type Notifier interface {
	Notify(string) error
}

type Twilio struct {
	client *twilio.TwilioClient
	from   string
	to     string
}

func NewTwilio(accountSid, authToken, from, to string) *Twilio {
	twilioClient := twilio.NewClient(accountSid, authToken)

	return &Twilio{
		client: twilioClient,
		from:   from,
		to:     to,
	}
}

func (t *Twilio) Notify(body string) error {
	resp, err := twilio.NewMessage(t.client, t.from, t.to, twilio.Body(body))
	logging.Log.Debugf("Twilio Response: %#v", resp)
	return err
}
