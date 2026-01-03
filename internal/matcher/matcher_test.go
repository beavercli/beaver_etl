package matcher

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/beavercli/beaver_etl/internal/manifest"
	"github.com/bradleyjkemp/cupaloy"
)

var (
	updateSnapshots    = flag.Bool("update", false, "update matcher snapshots")
	updateSnapshotsAlt = flag.Bool("update-snapshots", false, "update matcher snapshots")
)

func TestMatchFilesSnapshots(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	paths := []string{
		filepath.Join(root, "docs", "readme.md"),
		filepath.Join(root, "notes", "todo.txt"),
		filepath.Join(root, "src", "main.go"),
		filepath.Join(root, "src", "nested", "util.go"),
		filepath.Join(root, "src", "ignore.tmp"),
	}
	for _, path := range paths {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", path, err)
		}
		if err := os.WriteFile(path, []byte("test\n"), 0o644); err != nil {
			t.Fatalf("write %s: %v", path, err)
		}
	}

	manifestPath := filepath.Join(root, "manifest.yaml")
	manifestBody := `version: 1
sources:
  - type: file
    path: docs/readme.md
    title: Readme
    language: markdown
    tags: [docs, intro]
    link: https://example.com/docs
    contributors:
      - name: Ada
        last_name: Lovelace
        email: ada@example.com
  - type: pattern
    include:
      - src/**/*.go
      - notes/*.txt
    exclude:
      - src/**/ignore.*
    tags: [code]
    link: https://example.com/code
    contributors:
      - name: Grace
        last_name: Hopper
        email: grace@example.com
`
	if err := os.WriteFile(manifestPath, []byte(manifestBody), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	m, err := manifest.ParseManifest(manifestPath)
	if err != nil {
		t.Fatalf("parse manifest: %v", err)
	}

	matches, err := MatchFiles(root, manifestPath, m)
	if err != nil {
		t.Fatalf("match files: %v", err)
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Path < matches[j].Path
	})

	normalized, err := json.MarshalIndent(matches, "", "  ")
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	normalized = append(normalized, '\n')

	snapshotter := cupaloy.New(
		cupaloy.SnapshotFileExtension(".json"),
		cupaloy.ShouldUpdate(func() bool {
			return *updateSnapshots || *updateSnapshotsAlt || os.Getenv("UPDATE_SNAPSHOTS") != ""
		}),
		cupaloy.FailOnUpdate(false),
		cupaloy.SnapshotSubdirectory("testdata"),
	)
	snapshotter.SnapshotT(t, string(normalized))
}
