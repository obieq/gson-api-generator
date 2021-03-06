package codegen

// GetNewGitIgnoreTemplate => builds the template string for generating the git ignore file
func GetNewGitIgnoreTemplate(cc CodeConfig) string {
	return `
# See http://help.github.com/ignore-files/ for more about ignoring files.
#
# If you find yourself ignoring temporary files generated by your text editor
# or operating system, you probably want to add a global ignore instead:
#   git config --global core.excludesfile '~/.gitignore_global'

# Ignore tags
/tags

# Ignore tmp
/tmp

# Ignore test coverage files
*.coverprofile
*coverage.out

# Ignore swap files
*.swp
*.swo

# Ignore config files
# /config.json`
}
