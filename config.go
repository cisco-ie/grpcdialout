package main

import (
	"log"

	"github.com/spf13/viper"
)

//ConfigStruct holds config info
type ConfigStruct struct {
	Kafka KafkaStruct `json:"kafka"`
	Raw   bool        `json:"raw"`
	Dump  bool        `json:"dump"`
	File  string      `json:"filename"`
}

//KafkaStruct hold kafka config
type KafkaStruct struct {
	Brokers []string `json:"broker"`
	Topic   string   `json:"topic"`
}

//Configuration is the global config object
var Configuration ConfigStruct

//ConfigLoader loads config from file
func ConfigLoader() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetDefault("raw", false)
	viper.SetDefault("dump", false)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("fatal error config file: %s", err)
	}
	err = viper.Unmarshal(&Configuration)
	if err != nil {
		log.Fatalf("enable to decode into struct, %v", err)
	}

}
