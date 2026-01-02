package manifest

import (
	"fmt"

	"go.yaml.in/yaml/v4"
)

type SourceSpec interface {
	Kind() SourceType
	Validate() error
}

type Source struct {
	Type SourceType `yaml:"type" validate:"required,source_type"`
	Spec SourceSpec
}

type FileSource struct {
	Title        *Title        `yaml:"title" validate:"omitempty,title"`
	Language     *Language     `yaml:"language" validate:"lang"`
	Path         Path          `yaml:"path" validate:"path"`
	Tags         []Tag         `yaml:"tags" validate:"dive,tag"`
	Link         *Link         `yaml:"link" validate:"omitempty,link"`
	Contributors []Contributor `yaml:"contributors" validate:"dive"`
}

func (f *FileSource) Kind() SourceType {
	return File
}
func (f *FileSource) Validate() error {
	return nil
}

type PatternSource struct {
	Include      []PatternPath `yaml:"include" validate:"dive,pattern_path"`
	Exclude      []PatternPath `yaml:"exclude" validate:"dive,pattern_path"`
	Tags         []Tag         `yaml:"tags" validate:"dive,tag"`
	Link         *Link         `yaml:"link" validate:"omitempty,link"`
	Contributors []Contributor `yaml:"contributors" validate:"dive"`
}

func (p *PatternSource) Kind() SourceType {
	return Pattern
}
func (p *PatternSource) Validate() error {
	return nil
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
