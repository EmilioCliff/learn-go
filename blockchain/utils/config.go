package utils

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	PORT              string        `mapstructure:"PORT"`
	ENVIRONMENT       string        `mapstructure:"ENVIRONMENT"`
	NEIGHBORS         []string      `mapstructure:"NEIGHBORS"`
	MINING_DIFFICULTY int           `mapstructure:"MINING_DIFFICULTY"`
	MINING_SENDER     string        `mapstructure:"MINING_SENDER"`
	MINING_REWARD     float32       `mapstructure:"MINING_REWARD"`
	MINING_TIMER      time.Duration `mapstructure:"MINING_TIMER"`
	HOST              string        `mapstructure:"HOST"`
}

func LoanConfig() (Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	setDefaults()

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found, using environment variables")
		} else {
			return Config{}, fmt.Errorf("failed to read config: %s", err.Error())
		}
	}

	var config Config

	return config, viper.Unmarshal(&config)
}

func setDefaults() {
	viper.SetDefault("PORT", "")
	viper.SetDefault("ENVIRONMENT", "")
	viper.SetDefault("NEIGHBORS", []string{})
	viper.SetDefault("MINING_DIFFICULTY", 3)
	viper.SetDefault("MINING_SENDER", "THE_BLOCKCHAIN")
	viper.SetDefault("MINING_REWARD", 1.0)
	viper.SetDefault("MINING_TIMER", 10*time.Second)
	viper.SetDefault("HOST", "localhost")
}
