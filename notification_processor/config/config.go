package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/wfernandes/iot/logging"
)

type Config struct {
	TwilioAccountSid string
	TwilioAuthToken  string
	TwilioFromPhone  string

	To        string
	BrokerUrl string
	LogLevel  logging.LogLevel
}

func FromBytes(data []byte) (*Config, error) {
	config := &Config{}
	err := json.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	err = config.validate()
	if err != nil {
		return nil, err
	}

	return config, nil
}

func Configuration(configFilePath string) (*Config, error) {

	configBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	return FromBytes(configBytes)
}

func (c *Config) validate() error {
	if c.TwilioAccountSid == "" {
		return fmt.Errorf("Twilio Account SID is required")
	}

	if c.TwilioAuthToken == "" {
		return fmt.Errorf("Twilio Auth Token is required")
	}

	if c.TwilioFromPhone == "" {
		return fmt.Errorf("Twilio From Phone is required")
	}

	if c.BrokerUrl == "" {
		return fmt.Errorf("MQTT broker url is required")
	}

	if c.LogLevel.String() == "INVALID" {
		c.LogLevel = logging.INFO
	}
	return nil
}
