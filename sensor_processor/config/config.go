package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)

type Config struct {
	Sensors     map[string]string
	NotifierUrl string
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

	if c.NotifierUrl == "" {
		return fmt.Errorf("Notifier Url required")
	}

	return nil
}

func expectedSensors() []string {
	return []string{"touch", "button"}
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
