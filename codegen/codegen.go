package codegen

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/obieq/gas"
)

const directoryApp = "app"

const fileTypeMain = "main"
const fileTypeGitIgnore = "gitconfig"
const fileTypeConfig = "config"
const fileTypeConfigJson = "config_json"
const fileTypeGodepsJson = "godeps_json"
const fileTypeGoDir = "godir"
const fileTypeDBType = "db_type"
const fileTypeWebServer = "web_server"
const fileTypeController = "controller"
const fileTypeModel = "model"
const fileTypeResource = "resource"
const fileTypeRoute = "route"
const fileTypeSuiteTest = "suite_test"
const fileTypeRequestTest = "request_test"

type CodeConfig struct {
	Action            string
	DBType            string
	EntityName        string
	EntityDefinition  string
	IDType            string
	IncludeTimestamps bool
	GoPath            string
	Debug             bool
}

type Model struct {
	ActiveRecordStructName string
	Schemata               []Schema
}

type Schema struct {
	Name           string
	DasherizedName string
	Type           string
	Tag            string
	Required       bool
	SortOrder      int // the order in which the field was specified in the comman line
}

func BuildJsonTag(name string, objectType string) string {
	var tag string
	name = gas.String(name).Dasherize()
	switch objectType {
	case "model":
		tag = `json:"` + name + `,omitempty"`
	case "resource":
		tag = `json:"` + name + `,omitempty" jsonapi:"name=` + name + `"`
	default:
		log.Panicln("invalid objectType for BuildJsonTag:", objectType)
	}

	return tag
}

func (cc CodeConfig) directoryName(fileType string) string {
	pluralizedFileType := gas.String(fileType).Pluralize()

	switch fileType {
	case fileTypeMain, fileTypeGitIgnore, fileTypeConfigJson, fileTypeGoDir:
		return "./" + cc.GoPath + "/"
	case fileTypeConfig:
		return "./" + cc.GoPath + "/config/"
	case fileTypeGodepsJson:
		return "./" + cc.GoPath + "/Godeps/"
	case fileTypeWebServer, fileTypeSuiteTest:
		return "./" + cc.GoPath + "/web_server/"
	case fileTypeDBType:
		return "./" + cc.GoPath + "/app/models/"
	case fileTypeController, fileTypeModel, fileTypeResource, fileTypeRoute:
		return "./app/" + pluralizedFileType + "/"
	case fileTypeRequestTest:
		return "./web_server/"
	default:
		log.Panicln("Invalid File Type:", fileType)
	}

	return ""
}

func (cc CodeConfig) fileName(fileType string) string {
	switch fileType {
	case fileTypeMain:
		return "main.go"
	case fileTypeGitIgnore:
		return ".gitignore"
	case fileTypeConfig:
		return "config.go"
	case fileTypeConfigJson:
		return "config.json"
	case fileTypeGodepsJson:
		return "Godeps.json"
	case fileTypeGoDir:
		return ".godir"
	case fileTypeWebServer:
		return "server.go"
	case fileTypeDBType:
		return "base_model.go"
	case fileTypeSuiteTest:
		return "web_server_suite_test.go"
	case fileTypeController, fileTypeModel, fileTypeResource, fileTypeRoute:
		return gas.String(cc.EntityName).Underscore() + ".go"
	case fileTypeRequestTest:
		return gas.String(cc.EntityName).Underscore() + "_test.go"
	default:
		log.Panicln("Invalid File Type:", fileType)
	}

	return ""
}

func (cc CodeConfig) filePath(fileType string) string {
	return cc.directoryName(fileType) + cc.fileName(fileType)
}

func (cc CodeConfig) Generate() error {
	if err := cc.generateFile(fileTypeController); err != nil {
		return err
	}
	if err := cc.generateFile(fileTypeModel); err != nil {
		return err
	}
	if err := cc.generateFile(fileTypeResource); err != nil {
		return err
	}
	if err := cc.generateFile(fileTypeRoute); err != nil {
		return err
	}
	if err := cc.generateFile(fileTypeRequestTest); err != nil {
		return err
	}

	return nil
}

func (cc CodeConfig) Destroy() error {
	for _, fileType := range []string{fileTypeController, fileTypeModel, fileTypeResource, fileTypeRoute, fileTypeRequestTest} {
		if err := os.Remove(path.Join(cc.filePath(fileType))); err != nil {
			return err
		}
	}

	return nil
}

func (cc CodeConfig) NewProject() error {
	// create file directory
	os.MkdirAll("./"+cc.GoPath+"/lib", 0777)
	os.MkdirAll("./"+cc.GoPath+"/config", 0777)
	os.MkdirAll("./"+cc.GoPath+"/Godeps", 0777)
	os.MkdirAll("./"+cc.GoPath+"/web_server", 0777)
	os.MkdirAll("./"+cc.GoPath+"/app/controllers", 0777)
	os.MkdirAll("./"+cc.GoPath+"/app/models", 0777)
	os.MkdirAll("./"+cc.GoPath+"/app/resources", 0777)
	os.MkdirAll("./"+cc.GoPath+"/app/routes", 0777)

	cc.generateFile(fileTypeMain)
	cc.generateFile(fileTypeGitIgnore)
	cc.generateFile(fileTypeConfig)
	cc.generateFile(fileTypeConfigJson)
	cc.generateFile(fileTypeGodepsJson)
	cc.generateFile(fileTypeGoDir)
	cc.generateFile(fileTypeWebServer)
	cc.generateFile(fileTypeDBType)
	cc.generateFile(fileTypeSuiteTest)
	return nil
}

func (cc CodeConfig) UpperCasedEntityName() string {
	return strings.ToUpper(cc.EntityName)
}

func (cc CodeConfig) LowerCasedEntityName() string {
	return strings.ToLower(cc.EntityName)
}

func (cc CodeConfig) PluralizedEntityName() string {
	return gas.String(cc.EntityName).Pluralize()
}

func (cc CodeConfig) DasherizedEntityName(shouldPluralize bool) string {
	s := gas.String(cc.EntityName).Dasherize()
	if shouldPluralize {
		s = gas.String(s).Pluralize()
	}

	return s
}

func (cc CodeConfig) CommandLineArgs() string {
	args := fmt.Sprintf("%s", os.Args)
	args = args[1 : len(args)-1] // remove first char '[' and last char ']'
	return args
}

func (cc CodeConfig) TimeStampHeader(includeCommandLineArgs bool) string {
	args := ""
	if includeCommandLineArgs {
		args = `
		// * Command Line Args
		// * ` + cc.CommandLineArgs()
	}

	return `// * -----------------------------------------------------------
// *
// * Auto-Generated by gson-api-generator on ` + time.Now().String() + `
// *` + args + `
// *
// * -----------------------------------------------------------`
}

func (cc CodeConfig) AppName() string {
	arr := strings.Split(cc.GoPath, "/")
	return arr[len(arr)-1]
}

func (cc CodeConfig) ParseModel() Model {
	var m Model

	if err := json.Unmarshal([]byte(cc.EntityDefinition), &m); err != nil {
		log.Fatalln("Could not parse model schemata:", err)
	}

	if cc.Debug {
		log.Println("--- Model Schemata ---")
		for i := range m.Schemata {
			schema := &m.Schemata[i]
			schema.DasherizedName = gas.String(schema.Name).Dasherize()
			schema.SortOrder = i
			log.Println(schema.Name+":"+schema.Type+":", schema.SortOrder)
		}
	}

	return m
}

func (cc CodeConfig) ParseRequiredFields() []string {
	requiredFields := []string{}
	m := cc.ParseModel()

	for i := range m.Schemata {
		s := m.Schemata[i]
		if s.Required {
			requiredFields = append(requiredFields, s.Name)
		}
	}

	if cc.Debug {
		log.Println("TemplateModel.parseRequiredFields(:", requiredFields)
	}

	return requiredFields
}

func (cc CodeConfig) IsIDString() bool {
	return cc.IDType == "string"
}

func (cc CodeConfig) IsMSSQL() bool {
	return cc.DBType == "mssql"
}

func (cc CodeConfig) IsPostgresql() bool {
	return cc.DBType == "postgresql"
}

func (cc CodeConfig) IsRethinkDB() bool {
	return cc.DBType == "rethinkdb"
}

func (cc CodeConfig) ParseOrderedKeys(includeCreatedAt bool, includeUpdatedAt bool) string {
	// s := `map[int]string{` 0: "name", 1: "description", 2: "created-at"}
	schemata := cc.ParseModel().Schemata
	length := len(schemata)
	s := `map[int]string{`
	for i, schema := range schemata {
		s += strconv.Itoa(schema.SortOrder) + `: "` + schema.DasherizedName + `"`
		if i < length {
			s += ", "
		}
	}
	if includeCreatedAt {
		s += strconv.Itoa(length) + `: "created-at"`
		length++
	}
	if includeUpdatedAt {
		if includeCreatedAt {
			s += ", "
		}
		s += strconv.Itoa(length) + `: "updated-at"`
		length++
	}
	s += "}"
	return s
}

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func (cc CodeConfig) generateFile(fileType string) error {
	// verify file type is acceptable
	found := false
	for _, s := range []string{
		fileTypeMain, fileTypeGitIgnore, fileTypeConfig, fileTypeConfigJson, fileTypeGodepsJson, fileTypeWebServer,
		fileTypeDBType, fileTypeGoDir, fileTypeController, fileTypeModel, fileTypeResource, fileTypeRoute,
		fileTypeSuiteTest, fileTypeRequestTest} {
		if s == fileType {
			found = true
		}
	}
	if !found {
		log.Fatal("Invalid File Type:", fileType)
	}

	// determine directory name
	directoryName := cc.directoryName(fileType)

	// create directories
	found, err := exists(directoryName)
	if err != nil {
		return err
	}

	if !found {
		os.MkdirAll(directoryName, 0777)
	}

	file, err := os.Create(path.Join(cc.filePath(fileType)))
	if err != nil {
		return err
	}
	w := bufio.NewWriter(file)

	defer func() {
		w.Flush()
		file.Close()
	}()

	var tmpl *template.Template

	switch fileType {
	case fileTypeMain:
		tmpl, err = template.New(fileType).Parse(GetNewMainTemplate(cc))
	case fileTypeGitIgnore:
		tmpl, err = template.New(fileType).Parse(GetNewGitIgnoreTemplate(cc))
	case fileTypeConfig:
		tmpl, err = template.New(fileType).Parse(GetNewConfigTemplate(cc))
	case fileTypeConfigJson:
		tmpl, err = template.New(fileType).Parse(GetNewConfigJsonTemplate(cc))
	case fileTypeGodepsJson:
		tmpl, err = template.New(fileType).Parse(GetNewGodepsJsonTemplate(cc))
	case fileTypeGoDir:
		tmpl, err = template.New(fileType).Parse(GetNewGoDirTemplate(cc))
	case fileTypeDBType:
		switch cc.DBType {
		case "mssql":
			tmpl, err = template.New(fileType).Parse(GetNewMSSQLBaseModelTemplate(cc))
		case "rethinkdb":
			tmpl, err = template.New(fileType).Parse(GetNewRethinkDBBaseModelTemplate(cc))
		case "postgresql":
			tmpl, err = template.New(fileType).Parse(GetNewPostgresqlBaseModelTemplate(cc))
		default:
			log.Panicln("no base model generator exists for:", cc)
		}
	case fileTypeWebServer:
		tmpl, err = template.New(fileType).Parse(GetNewServerTemplate(cc))
	case fileTypeSuiteTest:
		tmpl, err = template.New(fileType).Parse(GetNewTestSuiteTemplate(cc))
	case fileTypeController:
		tmpl, err = template.New(fileType).Parse(GetControllerTemplate(cc))
	case fileTypeModel:
		tmpl, err = template.New(fileType).Parse(GetModelTemplate(cc))
	case fileTypeResource:
		tmpl, err = template.New(fileType).Parse(GetResourceTemplate(cc))
	case fileTypeRoute:
		tmpl, err = template.New(fileType).Parse(GetRouteTemplate(cc))
	case fileTypeRequestTest:
		tmpl, err = template.New(fileType).Parse(GetRequestTestTemplate(cc))
	default:
		log.Panicln("Invalid File Type:", fileType)
	}

	err = tmpl.Execute(w, cc)

	return err
}
