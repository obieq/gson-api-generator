package codegen

// GetNewConfigTemplate => builds the template string for generating the config file
func GetNewConfigTemplate(cc CodeConfig) string {
	return `
	package config

  {{.TimeStampHeader false}}

	import (
		"log"

		"github.com/obieq/gas"
	)

	// Config for the {{.AppName}} web service
	var Config *config

	func init() {
		log.Println("calling {{.AppName}}.config.init()")
		Config = newConfig()
	}

	type config struct {
		gas.Config
		Version     string
		Environment string
		Port        int
		GsonApiUrl  string
		NewRelicKey string
		Debug       bool
	}

	func newConfig() *config {
		c := &config{}
		c.ParseConfigFile("config")
		c.Validate()

		return c
	}

	func (c *config) ParseConfigFile(configFileName string) error {
		err := c.Load("{{.AppName}}", "config", true)

		// get Version
	  c.Version = gas.GetString("version")

		// get Environment
	  c.Environment = gas.GetString("environment")

		// get Port
		c.Port = gas.GetInt("port")

		// get GsonApiUrl
		c.GsonApiUrl = gas.GetString("gson_api_url")

		// get NewRelic Key
    c.NewRelicKey = gas.GetString("new_relic_key")

		// get Debug flag
	  c.Debug = gas.GetBool("debug")

		return err
	}

	func (c *config) Validate() {
		// PANICS
		if c.Version == "" {
		  log.Panicln("{{.AppName}} config error: Version cannot be blank")
	  }
		if c.Environment == "" {
      log.Panicln("{{.AppName}} config error: Environment cannot be blank")
    }
		if c.Port == 0 {
			log.Panicln("{{.AppName}} config error: Port cannot be blank")
		}
		if c.GsonApiUrl == "" {
      log.Panicln("{{.AppName}} config error: Gson Api Url cannot be blank")
    }

		// WARNINGS
		if c.NewRelicKey == "" {
			log.Println("-- WARNING -- NewRelicKey is blank")
		}
	}`
}
