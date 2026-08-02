// Package node detects available JavaScript/Node runtimes on PATH.
package node

import (
	"bytes"
	"errors"
	"os/exec"
	"runtime"
)

// Runtime is a named executable detected on PATH.
type Runtime struct {
	// Name is the executable name, e.g. "node", "bun", "deno", "npx".
	Name string
	// Available is true if the executable was found on PATH.
	Available bool
	// Version is the reported version string, if detectable.
	Version string
}

// Detect returns the available JS runtimes from the known list.
func Detect() []Runtime {
	runtimes := []Runtime{
		{Name: "node"},
		{Name: "bun"},
		{Name: "deno"},
		{Name: "npx"},
		{Name: "npm"},
		{Name: "pnpm"},
		{Name: "yarn"},
	}
	var out []Runtime
	for i := range runtimes {
		r := &runtimes[i]
		path, err := exec.LookPath(r.Name)
		if err != nil {
			if errors.Is(err, exec.ErrDot) {
				continue
			}
			if runtime.GOOS == "windows" && err.Error() == "file does not exist: "+r.Name {
				continue
			}
			continue
		}
		if path == "" {
			continue
		}
		r.Available = true
		ver, _ := exec.Command(r.Name, "--version").Output()
		if len(ver) > 0 {
			r.Version = string(bytes.TrimSpace(ver))
		}
		out = append(out, *r)
	}
	return out
}
