package parser

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v4"
)

type Manifest struct {
	Version  Version       `yaml:"version"`
	Defaults *Default      `yaml:"defaults"`
	Ignore   []PatternPath `yaml:"ignore"`
	Sources  []Source      `yaml:"sources"`
}
type Default struct {
	Tags         []Tag         `yaml:"tags"`
	Link         *Link         `yaml:"link"`
	Language     *Language     `yaml:"language"`
	Contributors []Contributor `yaml:"contributors"`
}
type Source struct {
	Type SourceType `yaml:"type"`
	Spec SourceSpec
}

func (s *Source) UnmarshalYAML(node *yaml.Node) error {
	var t struct {
		Type SourceType `yaml:"type"`
	}
	if err := node.Decode(&t); err != nil {
		return err
	}
	switch t.Type {
	case File:
		var fs FileSource
		if err := node.Decode(&fs); err != nil {
			return err
		}
		s.Type, s.Spec = t.Type, &fs
	case Pattern:
		var ps PatternSource
		if err := node.Decode(&ps); err != nil {
			return err
		}
		s.Type, s.Spec = t.Type, &ps
	default:
		return fmt.Errorf("Provided type %s is not supported", t.Type)
	}

	return nil
}

type SourceType string

const (
	File    SourceType = "file"
	Pattern SourceType = "pattern"
)

type SourceSpec interface {
	Kind() SourceType
	Validate() error
}
type FileSource struct {
	Title        *Title        `yaml:"title"`
	Path         Path          `yaml:"path"`
	Tags         []Tag         `yaml:"tags"`
	Link         *Link         `yaml:"link"`
	Contributors []Contributor `yaml:"contributors"`
}

func (f *FileSource) Kind() SourceType {
	return File
}
func (f *FileSource) Validate() error {
	return nil
}

type PatternSource struct {
	Include      []PatternPath `yaml:"include"`
	Exclude      []PatternPath `yaml:"exclude"`
	Tags         []Tag         `yaml:"tags"`
	Link         *Link         `yaml:"link"`
	Contributors []Contributor `yaml:"contributors"`
}

func (p *PatternSource) Kind() SourceType {
	return Pattern
}
func (p *PatternSource) Validate() error {
	return nil
}

type Version int

const (
	FirstVersion Version = 1
)

type Path string
type Title string
type PatternPath string
type Tag string
type Link string
type Language string
type Contributor struct {
	Name     Name     `yaml:"name"`
	LastName LastName `yaml:"last_name"`
	Email    Email    `yaml:"email"`
}
type Name string
type LastName string
type Email string

func ParseBeaverConfig(path string) (*Manifest, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dec := yaml.NewDecoder(f)
	dec.KnownFields(true)

	var m Manifest
	if err := dec.Decode(&m); err != nil {
		return nil, err
	}

	return &m, nil
}
