package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Service struct {
	Url         string `yaml:"url"`
	RateSeconds int    `yaml:"rateSeconds"`
	Tokens      int    `yaml:"tokens"`
}

type ConfigStruct struct {
	Port          string `yaml:"port"`
	MetricsPort   string `yaml:"metricsPort"`
	CheckoutDbUrl string `yaml:"checkoutDbUrl"`
	Token         string `yaml:"token"`
	Services      struct {
		Loms           Service `yaml:"loms"`
		ProductService Service `yaml:"productService"`
	} `yaml:"services"`
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
