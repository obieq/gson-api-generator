package codegen

import (
	"log"
	"strings"
)

func parseResourceProperties(cc CodeConfig) string {
	props := ""
	t := ""
	m := cc.ParseModel()
	for i := range m.Schemata {
		s := m.Schemata[i]
		switch strings.ToLower(s.Type) {
		case "string":
			t = "String"
		case "int":
			t = "Int"
		case "bool":
		case "boolean":
			t = "Bool"
		case "float":
			t = "Float"
		case "time":
			t = "Time"
		default:
			log.Panicln("parseResourceProperties mapping is missing for:", t)
		}
		props += s.Name + " null." + t + " `" + BuildJsonTag(s.Name, "resource") + "`\n"
	}

	if cc.Debug {
		log.Println("GetModelTemplate.props:", props)
		log.Println("TemplateResource.parseResourceProperties:", props)
	}

	return props
}

func mapProperties(cc CodeConfig, isMapFromModel bool) string {
	props := ""
	m := cc.ParseModel()

	for i := range m.Schemata {
		s := m.Schemata[i]
		if isMapFromModel { // map resource to model
			switch strings.ToLower(s.Type) {
			case "string":
				props += `r.` + s.Name + ` = null.StringFrom(m.` + s.Name + `)`
			case "int":
				props += `r.` + s.Name + ` = null.IntFrom(int64(m.` + s.Name + `))`
			case "bool":
			case "boolean":
				props += `r.` + s.Name + ` = null.BoolFrom(m.` + s.Name + `)`
			default:
				log.Panicln("mapToModel mapping is missing for:", s.Type)
			}
		} else { // map resource from model
			switch strings.ToLower(s.Type) {
			case "string":
				props += `if !r.` + s.Name + `.IsZero() { m.` + s.Name + ` = r.` + s.Name + `.String }`
			case "int":
				props += `if !r.` + s.Name + `.IsZero() { m.` + s.Name + ` = int(r.` + s.Name + `.Int64) }`
			case "bool":
			case "boolean":
				props += `if !r.` + s.Name + `.IsZero() { m.` + s.Name + ` = r.` + s.Name + `.Bool }`
			default:
				log.Panicln("mapFromModel mapping is missing for:", s.Type)
			}
		}
		props += "\n"
	}

	return props
}

// GetResourceTemplate => builds the template string for generating a resource
func GetResourceTemplate(cc CodeConfig) string {
	props := parseResourceProperties(cc)

	return `
{{$model := .ParseModel}}
{{$upperCasedEntityName := .UpperCasedEntityName}}
{{$dasherizedAndPluralizedEntityName := .DasherizedEntityName true}}
{{$timeStampHeader := .TimeStampHeader false}}
package resources

{{$timeStampHeader}}

import (
  models "{{.GoPath}}/app/models"
	"github.com/obieq/goar"
	gapi "github.com/obieq/gson-api"
	"gopkg.in/guregu/null.v3"
)

const {{$upperCasedEntityName}}_RESOURCE_TYPE = "{{$dasherizedAndPluralizedEntityName}}"

// {{.EntityName}} Resource
type {{.EntityName}} struct {
	gapi.Resource ` + "`jsonapi:" + `"-"` + "`" + `
	` + props + `
	{{if .IncludeTimestamps}}
		CreatedAt null.Time ` + "`" + BuildJsonTag("CreatedAt", "resource") + "`" + `
		UpdatedAt null.Time ` + "`" + BuildJsonTag("UpdatedAt", "resource") + "`" + `
	{{end}}
}

func (r {{.EntityName}}) GetName() string {
	return {{$upperCasedEntityName}}_RESOURCE_TYPE
}

func (r *{{.EntityName}}) MapFromModel(model interface{}) (err error) {
	m := model.(models.{{.EntityName}})

	if !m.HasErrors() {
		r.ID = m.ID
		` + mapProperties(cc, true) + `
		{{if .IncludeTimestamps}}
		  r.CreatedAt = null.TimeFrom(m.CreatedAt)
			if !m.UpdatedAt.IsZero() {
			  r.UpdatedAt = null.TimeFrom(m.UpdatedAt)
		  }
		{{end}}
	} else {
		r.SetErrors(m.ErrorMap())
	}

	return nil
}

func (r *{{.EntityName}}) MapToModel(model interface{}) (err error) {
	m := model.(*models.{{.EntityName}})

	` + mapProperties(cc, false) + `
	{{if .IncludeTimestamps}}
		if m.CreatedAt.IsZero() { // we're inserting a new record
		  m.ID = r.ID
		}
	{{end}}
	// convert model to an active record model
	goar.ToAR(m)

	return nil
}`
}
