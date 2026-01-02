package parser

import "github.com/go-playground/validator/v10"

func newValidator() *validator.Validate {
	validate := validator.New()
	registerScalarAliases(validate)
	return validate
}

func registerScalarAliases(validate *validator.Validate) {
	validate.RegisterAlias("version", "gte=1")
	validate.RegisterAlias("path", "min=1")
	validate.RegisterAlias("title", "min=1")
	validate.RegisterAlias("pattern_path", "min=1")
	validate.RegisterAlias("tag", "min=1")
	validate.RegisterAlias("link", "url")
	validate.RegisterAlias("language", "min=1")
	validate.RegisterAlias("name", "min=1")
	validate.RegisterAlias("last_name", "min=1")
	validate.RegisterAlias("email_addr", "email")
	validate.RegisterAlias("source_type", "oneof=file pattern")
}
