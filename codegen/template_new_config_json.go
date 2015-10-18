package codegen

// GetNewConfigJsonTemplate => builds the template string for generating the config json file
func GetNewConfigJsonTemplate(cc CodeConfig) string {
	return `{
	"version": "0.1.0",
    "gson_api_url": "https://api.your-domain-goes-here.com/v1/",
	"port": 4000,
	"new_relic_key": "",
	"debug": false,
	"environment": "test",
	"test": { {{if .IsMSSQL}}
		"mssql": {
		  "aws": {
			"server": "",
			"port": 1433,
			"dbname": "",
			"username": "",
			"password": "",
			"maxidleconnections": 10,
			"maxopenconnections": 100,
			"debug": true
		  }
		}
  {{else if .IsPostgresql}}
        "postgresql": {
		  "aws": {
			"server": "",
			"port": 5432,
			"dbname": "",
			"username": "",
			"password": "",
			"maxidleconnections": 10,
			"maxopenconnections": 100,
			"debug": true
		  }
		}
  {{else if .IsRethinkDB}}
		"rethinkdb": {
		  "aws": {
		  "addresses":"",
		  "dbname": "",
		  "authkey": "",
		  "discoverhosts": false,
		  "maxidleconnections": 10,
		  "maxopenconnections": 100,
		  "debug": true
		}
	}{{end}}
  }
}`
}
