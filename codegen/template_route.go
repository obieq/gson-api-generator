package codegen

// GetRouteTemplate => builds the template string for generating a model
func GetRouteTemplate(cc CodeConfig) string {
	return `
	{{$timeStampHeader := .TimeStampHeader false}}
	package routes

  {{$timeStampHeader}}

	import (
		"github.com/go-martini/martini"
		c "{{.GoPath}}/app/controllers"
		gapi "github.com/obieq/gson-api"
	)

	var {{.LowerCasedEntityName}}ServerInfo gapi.JSONApiServerInfo

  func get{{.EntityName}}ServerInfo(c martini.Context) {
  	c.Map({{.LowerCasedEntityName}}ServerInfo)
  	c.Next()
  }

	func Load{{.EntityName}}Routes(r martini.Router, serverInfo gapi.JSONApiServerInfo) {
		{{.LowerCasedEntityName}}ServerInfo = serverInfo
		r.Get("/{{.DasherizedEntityName true}}",       get{{.EntityName}}ServerInfo, c.HandleGet{{.PluralizedEntityName}})
		r.Get("/{{.DasherizedEntityName true}}/:id",   get{{.EntityName}}ServerInfo, c.HandleGet{{.EntityName}})
		r.Post("/{{.DasherizedEntityName true}}",      get{{.EntityName}}ServerInfo, c.HandleCreate{{.EntityName}})
		r.Patch("/{{.DasherizedEntityName true}}/:id", get{{.EntityName}}ServerInfo, c.HandleUpdate{{.EntityName}})
		r.Delete("/{{.DasherizedEntityName true}}/:id", c.HandleDelete{{.EntityName}})
	}`
}
