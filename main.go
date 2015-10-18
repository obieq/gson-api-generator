package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/obieq/gson-api-generator/codegen"
)

const actionDestroy = "destroy"
const actionGenerate = "generate"
const actionNew = "new"

func main() {
	var action, dbType, entityName, entityDefinition, idType, goPath string
	var includeTimestamps, debug bool
	var pCount int

	start := time.Now()

	flag.StringVar(&entityName, "entityName", "", "Name of the entity to be generated or destroyed")
	flag.StringVar(&action, "action", "", "Action types: 'new', 'generate', 'destroy'")
	flag.StringVar(&dbType, "dbType", "", "Options: 'mssql', 'postgresql', 'rethinkdb'")
	flag.StringVar(&entityDefinition, "entityDefinition", "", "Entity definition as stringified json")
	flag.StringVar(&idType, "idType", "", "What's the database primary key type? Options: 'string', 'int'")
	flag.BoolVar(&includeTimestamps, "includeTimestamps", false, "Resource should map created-at and updated-at timestamps")
	flag.StringVar(&goPath, "goPath", "", "Go Path for your application. EX: github.com/obieq/goar")
	flag.BoolVar(&debug, "debug", false, "Run in debug mode")
	flag.IntVar(&pCount, "p", 4, "Parallell running for code generator")
	flag.Parse()

	if debug {
		log.Println("--- Command Line Args ---")
		log.Println(os.Args)
	}

	runtime.GOMAXPROCS(pCount)

	codeConfig := &codegen.CodeConfig{
		Action:            action,
		DBType:            dbType,
		EntityName:        entityName,
		EntityDefinition:  entityDefinition,
		IDType:            idType,
		IncludeTimestamps: includeTimestamps,
		GoPath:            goPath,
		Debug:             debug,
	}

	log.Println(codeConfig)

	validateAction(action)

	switch action {
	case actionGenerate:
		validateEntityName(entityName)
		validateEntityDefinition(entityDefinition)
		validateDBType(dbType)
		validateIDType(idType)
		validateGoPath(goPath)
		if err := codeConfig.Generate(); err != nil {
			log.Fatalln("Error occured during Generate():", err)
		}
		formatCodes("app", debug)
	case actionDestroy:
		validateEntityName(entityName)
		if err := codeConfig.Destroy(); err != nil {
			log.Fatalln("Error occured during Destroy():", err)
		}
	case actionNew:
		validateGoPath(goPath)
		validateDBType(dbType)
		if err := codeConfig.NewProject(); err != nil {
			log.Fatalln("Error occured during NewProject():", err)
		}
		formatCodes("./"+goPath, debug)
	default:
		log.Fatalln("Invalid action:", codeConfig.Action)
	}

	log.Println("Finished in", time.Since(start))
}

func formatCodes(pkg string, debug bool) {
	if debug {
		log.Println("-- running gofmt *.go for:", pkg, " ---")
	}

	var out bytes.Buffer
	cmd := exec.Command("gofmt", "-w", pkg)
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		log.Println(out.String())
		log.Fatalln("gofmt failed:", err)
	} else {
		if debug {
			log.Println("--- gofmt was successful ---")
		}
	}
}

func printUsages(message ...interface{}) {
	for _, x := range message {
		fmt.Println(x)
	}
	fmt.Println("\nUsage:")
	flag.PrintDefaults()
	log.Panic("exiting....")
}

func validateAction(action string) {
	if action == "" {
		printUsages("Please provide the type of action to take.  Options are: 'destroy', 'generate' or 'new'.")
		return
	}
}

func validateProjectName(projectName string) {
	if projectName == "" {
		printUsages("Please provide the project name.")
		return
	}
}

func validateDBType(dbType string) {
	if dbType == "" {
		printUsages("Please provide the database type ('mssql', 'postgresql', 'rethinkdb').")
		return
	}
}

func validateEntityName(entityName string) {
	if entityName == "" {
		printUsages("Please provide the entity name.")
		return
	}
}

func validateEntityDefinition(entityDefinition string) {
	if entityDefinition == "" {
		printUsages("Please provide the entity definition in the form of stringified json.")
		return
	}
}

func validateIDType(idType string) {
	if idType == "" {
		printUsages("Please provide the id type ('string', 'int').")
		return
	}
}

func validateGoPath(goPath string) {
	if goPath == "" {
		printUsages("Please provide the Go Path for your application")
		return
	}
}
