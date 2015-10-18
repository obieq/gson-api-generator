package codegen

func GetNewPostgresqlBaseModelTemplate(cc CodeConfig) string {
	return `
	package models

  {{.TimeStampHeader false}}

	import (
	    arPostgresql "github.com/obieq/goar/db/postgresql"
		  cfg "{{.GoPath}}/config"
	)

	type PostgresqlBaseModel struct {
		arPostgresql.ArPostgres
	}

	func (m *PostgresqlBaseModel) DBConnectionName() string {
    return "aws" // NOTE: this value *could* be pulled from ENV or config file
  }

  func (m *PostgresqlBaseModel) DBConnectionEnvironment() string {
		return cfg.Config.Environment // NOTE: this value *should* be pulled from ENV or config file
  }

	func (m *PostgresqlBaseModel) BeforeSave() error {
		var err error = nil

		// Custom implementation goes here

		return err
	}`
}
