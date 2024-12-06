package i43

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/tools/txtar"
)

func TestI43(t *testing.T) {
	files, err := filepath.Glob(filepath.FromSlash("testdata/txtar/*.txtar"))
	if err != nil {
		t.Fatalf("failed to list testdata files: %s", err)
	}

	t.Logf("Found %q txtars", strings.Join(files, ", "))

	for _, f := range files {
		t.Logf("parsing %q", f)
		arch, err := txtar.ParseFile(f)
		if err != nil {
			t.Errorf("failed to parse txtar file %q: %s", f, err)
		}
		t.Logf("parsed %q txtar: %s", f, arch.Comment)
		dir := t.TempDir()
		for _, af := range arch.Files {
			t.Logf("file %q: %s", af.Name, af.Data)

			name := filepath.Join(dir, af.Name)
			t.Logf("Extracting %q to %q", af.Name, name)
			if err := os.MkdirAll(filepath.Dir(name), 0o777); err != nil {
				t.Fatalf("failed to create dir: %s", err)
			}
			if err := os.WriteFile(name, af.Data, 0o666); err != nil {
				t.Fatalf("failed to write file: %s", err)
			}
		}

	}
}
