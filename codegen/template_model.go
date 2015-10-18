package codegen

import "log"

func parseModelProperties(cc CodeConfig) string {
	props := ""
	m := cc.ParseModel()

	for i := range m.Schemata {
		s := m.Schemata[i]
		if s.Tag > "" { // prepend the Tag value with a whitespace
			s.Tag = " " + s.Tag
		}
		props += s.Name + " " + s.Type + " `" + BuildJsonTag(s.Name, "model") + s.Tag + "`\n"
	}

	if cc.Debug {
		log.Println("TemplateModel.parseModelProperties:", props)
	}

	return props
}

// GetModelTemplate => builds the template string for generating a model
func GetModelTemplate(cc CodeConfig) string {
	// var pluralizedEntityName = gas.String(entityName).Pluralize()
	props := parseModelProperties(cc)

	return `
{{$model := .ParseModel}}
{{$timeStampHeader := .TimeStampHeader true}}
package models

{{$timeStampHeader}}

import goar "github.com/obieq/goar"

type {{.EntityName}} struct {
  {{$model.ActiveRecordStructName}}
  ` + props + `
}

func (model {{.EntityName}}) ToActiveRecord() *{{.EntityName}} {
	return goar.ToAR(&model).(*{{.EntityName}})
}

func (m *{{.EntityName}}) Validate() {
  {{range $rf := .ParseRequiredFields}}m.Validation.Required("{{$rf}}", m.{{$rf}})
  {{end}}}`
}
