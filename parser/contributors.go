package parser

type Contributor struct {
	Name     Name     `yaml:"name" validate:"name"`
	LastName LastName `yaml:"last_name" validate:"last_name"`
	Email    Email    `yaml:"email" validate:"email_addr"`
}
