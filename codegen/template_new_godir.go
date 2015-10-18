package codegen

// GetNewGoDirTemplate => builds the template string for generating the godir file
func GetNewGoDirTemplate(cc CodeConfig) string {
	return `
	{{.GoPath}}`
}
