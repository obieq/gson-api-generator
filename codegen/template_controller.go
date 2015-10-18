package codegen

// GetControllerTemplate => builds the template string for generating a controller
func GetControllerTemplate(cc CodeConfig) string {
	return `
package controllers

{{.TimeStampHeader false}}

import (
  "io/ioutil"
	"net/http"

  "github.com/go-martini/martini"
	"github.com/manyminds/api2go/jsonapi"
  "github.com/martini-contrib/render"
  goar "github.com/obieq/goar"
  gapi "github.com/obieq/gson-api"
  models "{{.GoPath}}/app/models"
  resources "{{.GoPath}}/app/resources"
)

func HandleGet{{.PluralizedEntityName}}(r render.Render, serverInfo gapi.JSONApiServerInfo) {
var instances []resources.{{.EntityName}}
var jsonApiError *gapi.JsonApiError

dbModels := make([]models.{{.EntityName}}, 0)
err := models.{{.EntityName}}{}.ToActiveRecord().All(&dbModels, nil)

// map the models to resources
if err == nil {
	instances = make([]resources.{{.EntityName}}, len(dbModels))
	for i, m := range dbModels {
		instance := resources.{{.EntityName}}{}
		instance.MapFromModel(m)
		instances[i] = instance
	}
} else {
	jsonApiError = &gapi.JsonApiError{Status: "400", Detail: err.Error()}
}

gapi.HandleIndexResponse(serverInfo, jsonApiError, instances, r)
}

func HandleGet{{.EntityName}}(args martini.Params, r render.Render, serverInfo gapi.JSONApiServerInfo) {
var instance resources.{{.EntityName}}
var jsonApiError *gapi.JsonApiError

var dbModel models.{{.EntityName}}
err := goar.ToAR(&models.{{.EntityName}}{}).Find(args["id"], &dbModel)

// map the model to the resource
if err == nil {
	instance = resources.{{.EntityName}}{}
	instance.MapFromModel(dbModel)
} else {
	if err.Error() == "record not found" {
	  jsonApiError = &gapi.JsonApiError{Status: "404", Detail: err.Error()}
	} else {
	  jsonApiError = &gapi.JsonApiError{Status: "400", Detail: err.Error()}
	}
}
gapi.HandleGetResponse(serverInfo, jsonApiError, instance, r)
}

func HandleCreate{{.EntityName}}(request *http.Request, r render.Render, serverInfo gapi.JSONApiServerInfo) {
var resource resources.{{.EntityName}}
var success bool = false
var jsonApiError *gapi.JsonApiError
var err error

// map the resource to the model
m := models.{{.EntityName}}{}

// unmarshal json api request
defer request.Body.Close()
body, _ := ioutil.ReadAll(request.Body)
err = jsonapi.UnmarshalFromJSON(body, &resource)

if err == nil {
// map resource to model
resource.MapToModel(&m)

// persist the model
success, err = m.Save() // TODO: implement patch?
}

// map the model to the resource
if err == nil {
	resource = resources.{{.EntityName}}{}
	resource.MapFromModel(m)
} else {
	jsonApiError = &gapi.JsonApiError{Status: "400", Detail: err.Error()}
}

// process result
gapi.HandlePostResponse(serverInfo, success, jsonApiError, &resource, r)
}

func HandleUpdate{{.EntityName}}(args martini.Params, request *http.Request, r render.Render, serverInfo gapi.JSONApiServerInfo) {
var resource resources.{{.EntityName}}
var success bool = false
var jsonApiError *gapi.JsonApiError
var err error

var dbModel models.{{.EntityName}}
err = models.{{.EntityName}}{}.ToActiveRecord().Find(args["id"], &dbModel)

if err == nil {
	// unmarshal json api request
	defer request.Body.Close()
	body, _ := ioutil.ReadAll(request.Body)
	err = jsonapi.UnmarshalFromJSON(body, &resource)

	// update properties
	if err == nil {
	  resource.MapToModel(&dbModel)

	  // persist changes
	  success, err = dbModel.Save()
  }

	// map the model to the resource
	if err == nil {
		resource = resources.{{.EntityName}}{}
		resource.MapFromModel(dbModel)
	} else {
		jsonApiError = &gapi.JsonApiError{Status: "400", Detail: err.Error()}
	}

	// process result
	gapi.HandlePatchResponse(serverInfo, success, jsonApiError, &resource, r)
} else { // get failed, so re-use the get response method, which properly handles the error condition
	jsonApiError = &gapi.JsonApiError{Status: "400", Detail: err.Error()}
	gapi.HandleGetResponse(serverInfo, jsonApiError, dbModel, r)
}
}

func HandleDelete{{.EntityName}}(args martini.Params, r render.Render) {
var jsonApiError *gapi.JsonApiError
model := &models.{{.EntityName}}{}
{{if .IsIDString}}
  model.ID = args["id"]
{{else}}
  id, _ := strconv.Atoi(args["id"]) // How do you want to handle a conversion error?
	model.ID = id
{{end}}

if err := goar.ToAR(model).Delete(); err != nil {
	jsonApiError = &gapi.JsonApiError{Status: "400", Detail: err.Error()}
}

gapi.HandleDeleteResponse(jsonApiError, r)
}`
}
