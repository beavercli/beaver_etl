package manifest


import (
	"encoding/json"
	"flag"
	"os"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
)

var (
	updateSnapshots    = flag.Bool("update", false, "update parser snapshots")
	updateSnapshotsAlt = flag.Bool("update-snapshots", false, "update parser snapshots")
)

func TestParseManifestSnapshots(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "testdata/valid_beaver.yaml"},
		{name: "testdata/minimal_file.yaml"},
		{name: "testdata/defaults_and_pattern.yaml"},
		{name: "testdata/mixed_sources_optional.yaml"},
		{name: "testdata/unmarshal_casing.yaml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := ParseManifest(tt.name)
			if err != nil {
				t.Fatalf("cannot parse the config %s: %v", tt.name, err)
			}

			normalized, err := json.MarshalIndent(m, "", "  ")
			if err != nil {
				t.Fatalf("marshal %s: %v", tt.name, err)
			}
			normalized = append(normalized, '\n')

			snapshotter := cupaloy.New(
				cupaloy.ShouldUpdate(func() bool {
					return *updateSnapshots || *updateSnapshotsAlt || os.Getenv("UPDATE_SNAPSHOTS") != ""
				}),
				cupaloy.FailOnUpdate(false),
				cupaloy.SnapshotSubdirectory("testdata"),
			)
			snapshotter.SnapshotT(t, string(normalized))
		})
	}
}