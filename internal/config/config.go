package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DbConn       string `yaml:"database_url,omitempty"`
	JwtSecretKey string `yaml:"jwt_secret_key"`
}

func MustConfig(path string, validator *validator.Validate) *Config {
	const f = "NewConfig "

	config := &Config{}
	err := cleanenv.ReadConfig(path, config)
	if err != nil {
		panic(f + err.Error())
	}

	err = validator.Struct(config)
	if err != nil {
		panic(f + err.Error())
	}

	return config
}
