// Package node provides a Node tool and package script runner.
package node

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"charm.land/fantasy"
	"github.com/rfaiii/donk-cli-main/internal/agent/tools"
	"github.com/rfaiii/donk-cli-main/internal/permission"
)

// NodeParams are the parameters for the node tool.
type NodeParams struct {
	Description         string `json:"description" description:"A brief description of what the command does, try to keep it under 30 characters or so"`
	Command             string `json:"command" description:"The Node/JS command to execute"`
	Args                []string `json:"args,omitempty" description:"Optional arguments to append to the command"`
	WorkingDir          string `json:"working_dir,omitempty" description:"The working directory to execute the command in (defaults to current directory)"`
	RunInBackground     bool   `json:"run_in_background,omitempty" description:"Set to true (boolean) to run this command in the background. Use job_output to read the output later."`
	AutoBackgroundAfter int    `json:"auto_background_after,omitempty" description:"Seconds to wait before automatically moving the command to a background job (default: 60)"`
}

// NpmParams are the parameters for the npm tool.
type NpmParams struct {
	Description         string   `json:"description" description:"A brief description of what the command does, try to keep it under 30 characters or so"`
	Script              string   `json:"script" description:"The npm/pnpm/yarn script to run, e.g. build"`
	Arguments           []string `json:"arguments,omitempty" description:"Optional arguments to pass to the script"`
	WorkingDir          string   `json:"working_dir,omitempty" description:"The working directory to execute the command in (defaults to current directory)"`
	RunInBackground     bool     `json:"run_in_background,omitempty" description:"Set to true (boolean) to run this command in the background. Use job_output to read the output later."`
	AutoBackgroundAfter int      `json:"auto_background_after,omitempty" description:"Seconds to wait before automatically moving the command to a background job (default: 60)"`
	PackageManager      string   `json:"package_manager,omitempty" description:"Optional package manager to use. Defaults to npm. Supported: npm, pnpm, yarn"`
}

// PackageScript represents a script entry from a package.json.
type PackageScript struct {
	Name        string
	Description string
}

// ReadPackageScripts reads package.json scripts from the provided directory.
func ReadPackageScripts(dir string) ([]PackageScript, error) {
	if dir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}
	path := filepath.Join(dir, "package.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var pkg struct {
		Scripts map[string]string `json:"scripts"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, err
	}
	var scripts []PackageScript
	for name, desc := range pkg.Scripts {
		scripts = append(scripts, PackageScript{
			Name:        name,
			Description: desc,
		})
	}
	return scripts, nil
}

// DetectPackageManager attempts to detect the package manager for the project.
func DetectPackageManager(dir string) string {
	if dir == "" {
		dir = "."
	}
	if _, err := os.Stat(filepath.Join(dir, "pnpm-lock.yaml")); err == nil {
		return "pnpm"
	}
	if _, err := os.Stat(filepath.Join(dir, "yarn.lock")); err == nil {
		return "yarn"
	}
	if _, err := os.Stat(filepath.Join(dir, "package-lock.json")); err == nil {
		return "npm"
	}
	if _, err := os.Stat(filepath.Join(dir, "bun.lockb")); err == nil {
		return "bun"
	}
	return "npm"
}

const (
	NodeToolName = "node"
	NpmToolName  = "npm"
)

const defaultAutoBackgroundAfter = 60

func nodeDescription() string {
	return "Execute a Node/JavaScript command or script. Output is returned as text."
}

func npmDescription() string {
	return "Run an npm/pnpm/yarn script from package.json. Output is returned as text."
}

// NewNodeTool creates a new node tool.
func NewNodeTool(permissions permission.Service, workingDir string) fantasy.AgentTool {
	return fantasy.NewAgentTool(
		NodeToolName,
		nodeDescription(),
		func(ctx context.Context, params NodeParams, call fantasy.ToolCall) (fantasy.ToolResponse, error) {
			sessionID := tools.GetSessionFromContext(ctx)
			if sessionID == "" {
				return fantasy.ToolResponse{}, fmt.Errorf("session ID is required")
			}
			permissionDescription := fmt.Sprintf("execute %s with the following parameters:\n", NodeToolName)
			execDir := workingDir
			if params.WorkingDir != "" {
				execDir = params.WorkingDir
			}
			p, err := permissions.Request(
				ctx,
				permission.CreatePermissionRequest{
					SessionID:   sessionID,
					ToolCallID:  call.ID,
					Path:        execDir,
					ToolName:    NodeToolName,
					Action:      "execute",
					Description: permissionDescription,
					Params:      params,
				},
			)
			if err != nil {
				return fantasy.ToolResponse{}, err
			}
			if !p {
				return tools.NewPermissionDeniedResponse(), nil
			}
			if params.Command == "" {
				return fantasy.NewTextErrorResponse("missing command"), nil
			}
			cmd := buildNodeCommand(params)
			return runCommand(ctx, cmd, execDir, params.RunInBackground, params.AutoBackgroundAfter)
		},
	)
}

// NewNpmTool creates a new npm tool.
func NewNpmTool(permissions permission.Service, workingDir string) fantasy.AgentTool {
	return fantasy.NewAgentTool(
		NpmToolName,
		npmDescription(),
		func(ctx context.Context, params NpmParams, call fantasy.ToolCall) (fantasy.ToolResponse, error) {
			sessionID := tools.GetSessionFromContext(ctx)
			if sessionID == "" {
				return fantasy.ToolResponse{}, fmt.Errorf("session ID is required")
			}
			if params.Script == "" {
				return fantasy.NewTextErrorResponse("missing script"), nil
			}
			execDir := workingDir
			if params.WorkingDir != "" {
				execDir = params.WorkingDir
			}
			pm := params.PackageManager
			if pm == "" {
				pm = DetectPackageManager(execDir)
			}
			args := append([]string{"run", params.Script}, params.Arguments...)
			cmd := buildPackageCommand(pm, args)
			permissionDescription := fmt.Sprintf("execute %s with the following parameters:\n", pm)
			p, err := permissions.Request(
				ctx,
				permission.CreatePermissionRequest{
					SessionID:   sessionID,
					ToolCallID:  call.ID,
					Path:        execDir,
					ToolName:    NpmToolName,
					Action:      "execute",
					Description: permissionDescription,
					Params:      params,
				},
			)
			if err != nil {
				return fantasy.ToolResponse{}, err
			}
			if !p {
				return tools.NewPermissionDeniedResponse(), nil
			}
			return runCommand(ctx, cmd, execDir, params.RunInBackground, params.AutoBackgroundAfter)
		},
	)
}

func buildNodeCommand(params NodeParams) *exec.Cmd {
	args := append([]string{params.Command}, params.Args...)
	cmd := exec.Command("node", args...)
	if runtime.GOOS == "windows" {
		cmd = exec.Command("node.exe", args...)
	}
	return cmd
}

func buildPackageCommand(pm string, args []string) *exec.Cmd {
	switch pm {
	case "pnpm":
		return exec.Command("pnpm", args...)
	case "yarn":
		return exec.Command("yarn", args...)
	case "bun":
		return exec.Command("bun", args...)
	default:
		return exec.Command("npm", args...)
	}
}

func runCommand(ctx context.Context, cmd *exec.Cmd, dir string, runInBackground bool, autoBackgroundAfter int) (fantasy.ToolResponse, error) {
	if runInBackground {
		return runBackgroundCommand(ctx, cmd, dir, autoBackgroundAfter)
	}
	return runForegroundCommand(ctx, cmd, dir)
}

func runForegroundCommand(ctx context.Context, cmd *exec.Cmd, dir string) (fantasy.ToolResponse, error) {
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fantasy.ToolResponse{}, fmt.Errorf("command failed: %w: %s", err, strings.TrimSpace(string(output)))
	}
	return fantasy.NewTextResponse(string(output)), nil
}

func runBackgroundCommand(ctx context.Context, cmd *exec.Cmd, dir string, autoBackgroundAfter int) (fantasy.ToolResponse, error) {
	// TODO: integrate with shell background job manager if available.
	return fantasy.NewTextResponse("background execution not yet implemented"), nil
}

// HasNode reports whether a `node` executable is available on PATH.
func HasNode() bool {
	_, err := exec.LookPath("node")
	return err == nil
}

// NodeVersion returns the installed Node.js version string, if detectable.
func NodeVersion() string {
	out, err := exec.Command("node", "--version").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// HasPackageManager reports whether the given package manager is available on PATH.
func HasPackageManager(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
