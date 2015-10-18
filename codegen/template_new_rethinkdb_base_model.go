package codegen

func GetNewRethinkDBBaseModelTemplate(cc CodeConfig) string {
	return `
	package models

  {{.TimeStampHeader false}}

	import (
	    arRethinkDB "github.com/obieq/goar/db/rethinkdb"
		  cfg "{{.GoPath}}/config"
	)

	type RethinkDBBaseModel struct {
		arRethinkDB.ArRethinkDb
	}

	func (m *RethinkDBBaseModel) DBConnectionName() string {
    return "aws" // NOTE: this value *could* be pulled from ENV or config file
  }

  func (m *RethinkDBBaseModel) DBConnectionEnvironment() string {
		return cfg.Config.Environment // NOTE: this value *should* be pulled from ENV or config file
  }

	func (m *RethinkDBBaseModel) BeforeSave() error {
		var err error = nil

		// Custom implementation goes here

		return err
	}`
}
