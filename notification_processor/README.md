## Notification Processor

The Notification processor communicates with the MQTT broker and obtains a list of sensor keys. It then subscribes to
each of the sensor keys and listens for any changes upon which it will notify a user.

Currently, it notifies a single user via SMS using the [Twilio API library](https://github.com/carlosdp/twiliogo).

### Config Properties

| Name | Required | Description
|------|----------|------------|
|TwilioAccountSid | Yes | The AccountSID from your twilio account |
|TwilioAuthToken | Yes | The AuthToken from your twilio account |
|TwilioFromPhone | Yes | A managed phone number from your twilio account |
|BrokerUrl | Yes | The url of the broker key-value store |
|To        | No | The phone number of the user to receive the alerts. It should have the country code specified. E.g. `+1##########`
|LogLevel  | No | The logging level of the process. Default is INFO. Options are: `INFO, ERROR, FATAL, DEBUG`  