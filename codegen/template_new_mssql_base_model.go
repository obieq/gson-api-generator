package codegen

func GetNewMSSQLBaseModelTemplate(cc CodeConfig) string {
	return `
	package models

  {{.TimeStampHeader false}}

	import (
	    arMSSQL "github.com/obieq/goar/db/mssql"
		  cfg "{{.GoPath}}/config"
	)

	type MSSQLBaseModel struct {
		arMSSQL.ArMsSql
	}

	func (m *MSSQLBaseModel) DBConnectionName() string {
    return "aws" // NOTE: this value *could* be pulled from ENV or config file
  }

  func (m *MSSQLBaseModel) DBConnectionEnvironment() string {
		return cfg.Config.Environment // NOTE: this value *should* be pulled from ENV or config file
  }

	func (m *MSSQLBaseModel) BeforeSave() error {
		var err error = nil

		// Custom implementation goes here

		return err
	}`
}
