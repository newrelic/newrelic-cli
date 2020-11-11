package config

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// WithConfig loads and returns the CLI configuration.
func WithConfig(f func(c *Config)) {
	fmt.Printf("\n\n *** WithConfig *** \n\n")

	WithConfigFrom(DefaultConfigDirectory, f)
}

// WithConfigFrom loads and returns the CLI configuration from a specified location.
func WithConfigFrom(configDir string, f func(c *Config)) {
	fmt.Printf("\n\n *** WithConfigFrom *** \n\n")

	c, err := LoadConfig(configDir)
	if err != nil {
		log.Fatal(err)
	}

	f(c)
}

// WithConfiguration loads and returns the CLI configuration from a specified location.
func WithConfiguration(f func(c *viper.Viper)) {
	c, err := Configure()
	if err != nil {
		log.Fatal(err)
	}

	f(c)
}

func Configure() (*viper.Viper, error) {
	cfgViper := viper.New()

	configDir, err := getDefaultConfigDirectory()
	if err != nil {
		log.Fatal(err.Error())
	}

	cfgViper.SetEnvPrefix("NEW_RELIC_CLI")
	cfgViper.SetConfigName("config")
	cfgViper.SetConfigType("json")

	// Set the config file path
	cfgViper.AddConfigPath(configDir)

	// Read in environment variables that
	// match the environment prefix
	cfgViper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {

		fmt.Printf("\n\n *** loadConfig - AllSettings: %v \n\n", cfgViper.AllSettings())

		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Debugf("configuration file not found: %s", configDir)
			log.Debugf("creating empty configuration file")

			err := createConfigFile()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			// Config file was found but another error was produced
			log.Fatal(err.Error())
		}
	}

	hydrateConfig(cfgViper)

	return cfgViper, nil
}

func createConfigFile() error {

	return nil
}

// func hydrateConfig(viperConfig *viper.Viper) *Configuration {
// 	c := &Configuration{
// 		LogLevel: "DEBUG",
// 	}

// 	profiles := viperConfig.Get("")

// 	fmt.Print("\n\n **************************** \n")
// 	fmt.Printf("\n hydrateConfig:  %+v \n", *viperConfig)
// 	fmt.Print("\n **************************** \n\n")
// 	time.Sleep(2 * time.Second)

// 	return c
// }
