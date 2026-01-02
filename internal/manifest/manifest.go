package manifest

type Manifest struct {
	Version  Version       `yaml:"version" validate:"version"`
	Defaults *Default      `yaml:"defaults"`
	Ignore   []PatternPath `yaml:"ignore" validate:"dive,pattern_path"`
	Sources  []Source      `yaml:"sources" validate:"dive"`
}
