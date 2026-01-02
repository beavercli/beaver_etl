package manifest


import (
	"strings"

	"go.yaml.in/yaml/v4"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Path string
type Title string
type PatternPath string
type Tag string

func (t *Tag) UnmarshalYAML(node *yaml.Node) error {
	s, err := UnmarshallToLower(node)
	if err != nil {
		return err
	}

	*t = Tag(s)

	return nil
}

type Link string

func (l *Link) UnmarshalYAML(node *yaml.Node) error {
	s, err := UnmarshallToLower(node)
	if err != nil {
		return err
	}

	*l = Link(s)

	return nil
}

type Language string

func (l *Language) UnmarshalYAML(node *yaml.Node) error {
	s, err := UnmarshallToLower(node)
	if err != nil {
		return err
	}

	*l = Language(s)

	return nil
}

type Name string

func (n *Name) UnmarshalYAML(node *yaml.Node) error {
	s, err := UnmarshallToTitle(node)
	if err != nil {
		return err
	}

	*n = Name(s)

	return nil
}

type LastName string

func (n *LastName) UnmarshalYAML(node *yaml.Node) error {
	s, err := UnmarshallToTitle(node)
	if err != nil {
		return err
	}

	*n = LastName(s)

	return nil
}

type Email string

func (e *Email) UnmarshalYAML(node *yaml.Node) error {
	s, err := UnmarshallToLower(node)
	if err != nil {
		return err
	}

	*e = Email(s)

	return nil
}

type SourceType string

const (
	File    SourceType = "file"
	Pattern SourceType = "pattern"
)

type Version int

const (
	FirstVersion Version = 1
)

func UnmarshallToLower(node *yaml.Node) (string, error) {
	var s string
	if err := node.Decode(&s); err != nil {
		return "", err
	}
	return strings.ToLower(s), nil
}

func UnmarshallToTitle(node *yaml.Node) (string, error) {
	var s string
	if err := node.Decode(&s); err != nil {
		return "", err
	}
	s = cases.Title(language.English).String(s)
	return s, nil
}
