package codegen

// GetNewMainTemplate => builds the template string for generating the main file
func GetNewMainTemplate(cc CodeConfig) string {
	return `
	package main

	{{.TimeStampHeader true}}

	import (
		"strconv"

		cfg "{{.GoPath}}/config"
		web "{{.GoPath}}/web_server"
	)

	func main() {
		server := web.NewServer()
		server.RunOnAddr(":" + strconv.Itoa(cfg.Config.Port))
	}`
}
