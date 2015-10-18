package codegen

// GetNewServerTemplate => builds the template string for generating the server file
func GetNewServerTemplate(cc CodeConfig) string {
	return `
	{{$timeStampHeader := .TimeStampHeader true}}
	package webserver

	{{$timeStampHeader}}

	import (
		"log"
		"net/http"
		"net/http/httputil"

		"github.com/go-martini/martini"
		"github.com/martini-contrib/cors"
		"github.com/martini-contrib/gorelic"
		"github.com/martini-contrib/render"
		"github.com/twinj/uuid"
		cfg "{{.GoPath}}/config"
		gapi "github.com/obieq/gson-api"
		routes "{{.GoPath}}/app/routes"
	)

  // API_SERVER_INFO => used to construct relationship urls in responses
	//                    should pull values from config settings
  var API_SERVER_INFO = gapi.JSONApiServerInfo{BaseURL: "http://my.domain", Prefix: "v1"}

	// Server => Wrap the Martini server struct.
	type Server *martini.ClassicMartini

	// NewServer => Constructor
	func NewServer() Server {
		// switch the uuid format
		uuid.SwitchFormat(uuid.CleanHyphen)

		// configure martini
		m := Server(martini.Classic())

		// configure NewRelic
    if cfg.Config.NewRelicKey != "" {
      gorelic.InitNewrelicAgent(cfg.Config.NewRelicKey, "{{.AppName}}-"+cfg.Config.Environment, true)
      m.Use(gorelic.Handler)
    }

		// Martini Renderer
		m.Use(render.Renderer())

		if cfg.Config.Debug {
			// Intercept all HTTP Requests and dump to console
			m.Use(func(req *http.Request) {
				log.Println("-- BEGIN HTTP REQUEST INTERCEPTOR (DON'T RUN IN PRODUCTION!!) ---")
				dump, _ := httputil.DumpRequest(req, true)
				log.Println(string(dump))
				log.Println("--- END HTTP REQUEST INTERCEPTOR (DON'T RUN IN PRODUCTION!!) ---")
			})
		}

		// CORS configuration
		m.Use(cors.Allow(&cors.Options{
			AllowAllOrigins: true,
			//AllowOrigins:     []string{"https://*.foo.com"},
			// NOTE: no need to allow PUT b/c PATCH suffices
			AllowMethods:     []string{"PATCH", "GET", "POST", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "X-AUTH-TOKEN", "X-API-VERSION"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
		}))

		// Default Route, which returns the App's Version and Deployment Environment
	  m.Get("/", func() (int, string) {
			return 200, cfg.Config.Version + "\n" + cfg.Config.Environment
	  })

	  // Controller Routes
		m.Group("/v1", func(r martini.Router) {
			// routes.LoadExampleRoutes(r, API_SERVER_INFO)
		})

		return m
	}`
}
