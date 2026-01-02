package manifest

type Default struct {
	Tags         []Tag         `yaml:"tags" validate:"dive,tag"`
	Link         *Link         `yaml:"link" validate:"omitempty,link"`
	Language     *Language     `yaml:"language" validate:"omitempty,language"`
	Contributors []Contributor `yaml:"contributors" validate:"dive"`
}
