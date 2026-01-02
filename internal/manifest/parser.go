package manifest

import (
	"os"

	"go.yaml.in/yaml/v4"
)

func ParseManifest(path string) (*Manifest, error) {
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

	validate := newValidator()

	if err = validate.Struct(&m); err != nil {
		return nil, err
	}

	return &m, nil
}
