package manifest

import (
	"strings"
	"testing"
)

func TestParseManifestValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		wantTag string
	}{
		{name: "testdata/invalid_version.yaml", wantTag: "version"},
		{name: "testdata/invalid_link.yaml", wantTag: "link"},
		{name: "testdata/invalid_language.yaml", wantTag: "language"},
		{name: "testdata/invalid_email.yaml", wantTag: "email_addr"},
		{name: "testdata/invalid_tag.yaml", wantTag: "tag"},
		{name: "testdata/invalid_pattern_path.yaml", wantTag: "pattern_path"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseManifest(tt.name)
			if err == nil {
				t.Fatalf("expected validation error for %s", tt.name)
			}
			if !strings.Contains(err.Error(), tt.wantTag) {
				t.Fatalf("expected error to mention %q, got %v", tt.wantTag, err)
			}
		})
	}
}
