package node

import (
	"os"
	"testing"
)

func TestReadPackageScriptsReadsPackageJSON(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(dir+"/package.json", []byte(`{"scripts":{"build":"tsc","test":"go test ./..."}}`), 0644); err != nil {
		t.Fatalf("write package.json: %v", err)
	}
	scripts, err := ReadPackageScripts(dir)
	if err != nil {
		t.Fatalf("ReadPackageScripts: %v", err)
	}
	if len(scripts) != 2 {
		t.Fatalf("expected 2 scripts, got %d", len(scripts))
	}
	names := map[string]bool{}
	for _, s := range scripts {
		names[s.Name] = true
	}
	if !names["build"] || !names["test"] {
		t.Fatalf("expected build and test scripts, got %v", names)
	}
}

func TestDetectPackageManagerByLockfile(t *testing.T) {
	for _, tt := range []struct {
		name     string
		files    []string
		expected string
	}{
		{"pnpm", []string{"pnpm-lock.yaml"}, "pnpm"},
		{"yarn", []string{"yarn.lock"}, "yarn"},
		{"npm", []string{"package-lock.json"}, "npm"},
		{"bun", []string{"bun.lockb"}, "bun"},
		{"default", []string{}, "npm"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			for _, f := range tt.files {
				if err := os.WriteFile(dir+"/"+f, []byte("lock"), 0644); err != nil {
					t.Fatalf("write %s: %v", f, err)
				}
			}
			got := DetectPackageManager(dir)
			if got != tt.expected {
				t.Fatalf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}

func TestHasNode(t *testing.T) {
	if !HasNode() {
		t.Log("node not installed; skipping strict assertions")
	}
}

func TestNodeVersion(t *testing.T) {
	v := NodeVersion()
	if v == "" {
		t.Log("node not installed; version empty")
		return
	}
	if len(v) < 2 || v[0] != 'v' {
		t.Fatalf("expected version starting with v, got %s", v)
	}
}

func TestHasPackageManager(t *testing.T) {
	for _, name := range []string{"npm", "node"} {
		if !HasPackageManager(name) {
			t.Logf("%s not installed", name)
		}
	}
}
