package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	TwilioAccountSid string
	TwilioAuthToken  string
	TwilioFromPhone  string

	To   string
	Port uint

	BrokerUrl string
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
	return nil
}
