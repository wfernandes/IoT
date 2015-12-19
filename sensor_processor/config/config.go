package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/wfernandes/iot/logging"
)

type Config struct {
	Sensors     map[string]string
	NotifierUrl string
	BrokerUrl   string
	LogLevel    logging.LogLevel
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

	for k, v := range c.Sensors {
		if !validSensorType(v) {
			return fmt.Errorf("Invalid sensor type: %s", v)
		}

		if !validSensorPin(k) {
			return fmt.Errorf("Invalid sensor pin: %s", k)
		}
	}

	if c.LogLevel.String() == "INVALID" {
		c.LogLevel = logging.INFO
	}

	if c.NotifierUrl == "" {
		return fmt.Errorf("Notifier Url required")
	}

	if c.BrokerUrl == "" {
		return fmt.Errorf("Broker Url required")
	}

	return nil
}

func expectedSensors() []string {
	return []string{"touch", "sound"}
}

func validSensorType(value string) bool {
	for _, sensor := range expectedSensors() {
		if sensor == value {
			return true
		}
	}
	return false
}

func validSensorPin(pin string) bool {
	_, err := strconv.Atoi(pin)
	if err != nil {
		return false
	}
	return true
}
