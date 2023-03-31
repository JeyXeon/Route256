package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Topic struct {
	Name string `yaml:"name"`
}

type Kafka struct {
	Topics struct {
		OrderStateChange Topic `yaml:"order_state_change"`
	} `yaml:"topics"`
	Brokers []string `yaml:"brokers"`
}

type ConfigStruct struct {
	Port      string `yaml:"port"`
	LomsDbUrl string `yaml:"lomsDbUrl"`
	Kafka     Kafka  `yaml:"kafka"`
}

var ConfigData ConfigStruct

func Init() error {
	rawYAML, err := os.ReadFile("config.yml")
	if err != nil {
		return errors.WithMessage(err, "reading config file")
	}

	err = yaml.Unmarshal(rawYAML, &ConfigData)
	if err != nil {
		return errors.WithMessage(err, "parsing yaml")
	}

	return nil
}
