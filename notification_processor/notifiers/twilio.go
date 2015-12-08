package notifiers

import twilio "github.com/carlosdp/twiliogo"

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
	// TODO: Log the response
	_, err := twilio.NewMessage(t.client, t.from, t.to, twilio.Body(body))
	return err
}
