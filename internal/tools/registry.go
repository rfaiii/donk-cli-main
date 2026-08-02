// Package tools provides the DONK tool registry, slash commands,
// and tool panel definitions ported from the Rust donk-tui tools module.
package tools

import "github.com/rfaiii/donk-cli-main/internal/commands"

// ToolId identifies a built-in tool panel.
type ToolId string

const (
	ToolIdFiles      ToolId = "files"
	ToolIdTerminal   ToolId = "terminal"
	ToolIdAnimations ToolId = "animations"
	ToolIdShowcase   ToolId = "showcase"
	ToolIdSetup      ToolId = "setup"
	ToolIdRead       ToolId = "read"
	ToolIdSys        ToolId = "sys"
	ToolIdScenes     ToolId = "scenes"
	ToolIdHealth     ToolId = "health"
)

// ToolDefinition describes a tool panel.
type ToolDefinition struct {
	ID          ToolId
	Name        string
	Description string
	Slash       string
	Aliases     []string
}

// ToolRegistry is the authoritative list of built-in tools.
var ToolRegistry = []ToolDefinition{
	{
		ID:          ToolIdFiles,
		Name:        "files",
		Description: "Dual-pane file browser (superfile style)",
		Slash:       "/files",
		Aliases:     []string{"/finder"},
	},
	{
		ID:          ToolIdTerminal,
		Name:        "terminal",
		Description: "Live PTY shell session inside DONK",
		Slash:       "/terminal",
		Aliases:     []string{"/shell"},
	},
	{
		ID:          ToolIdAnimations,
		Name:        "animations",
		Description: "Boot reel, splash, fire, cycling, spring, eyes, matrix, vanish, logos",
		Slash:       "/animations",
		Aliases:     []string{"/anim"},
	},
	{
		ID:          ToolIdShowcase,
		Name:        "showcase",
		Description: "CRUSH ecosystem catalog",
		Slash:       "/showcase",
	},
	{
		ID:          ToolIdSetup,
		Name:        "setup",
		Description: "Huh-style setup form (backend / offline / theme)",
		Slash:       "/setup",
	},
	{
		ID:          ToolIdRead,
		Name:        "read",
		Description: "Glow-style markdown reader",
		Slash:       "/read",
	},
	{
		ID:          ToolIdSys,
		Name:        "sys",
		Description: "btop-style CPU / memory / process monitor",
		Slash:       "/sys",
	},
	{
		ID:          ToolIdScenes,
		Name:        "scenes",
		Description: "Mission Control Pro procedural animation scenes",
		Slash:       "/scenes",
	},
	{
		ID:          ToolIdHealth,
		Name:        "health",
		Description: "Service & endpoint health monitor",
		Slash:       "/health",
	},
}

// ToolIDs returns the ordered list of tool IDs.
func ToolIDs() []ToolId {
	ids := make([]ToolId, len(ToolRegistry))
	for i, def := range ToolRegistry {
		ids[i] = def.ID
	}
	return ids
}

// ToolSlashCommands returns all slash commands including aliases.
func ToolSlashCommands() []string {
	var cmds []string
	for _, def := range ToolRegistry {
		cmds = append(cmds, def.Slash)
		cmds = append(cmds, def.Aliases...)
	}
	return cmds
}

// FindToolBySlash returns the tool definition for a slash command, if any.
func FindToolBySlash(slash string) (*ToolDefinition, bool) {
	for _, def := range ToolRegistry {
		if def.Slash == slash {
			return &def, true
		}
		for _, alias := range def.Aliases {
			if alias == slash {
				return &def, true
			}
		}
	}
	return nil, false
}

// FindToolByID returns the tool definition for an ID, if any.
func FindToolByID(id ToolId) (*ToolDefinition, bool) {
	for i := range ToolRegistry {
		if ToolRegistry[i].ID == id {
			return &ToolRegistry[i], true
		}
	}
	return nil, false
}

// AsCustomCommands converts tool definitions into command-palette entries.
func AsCustomCommands() []commands.CustomCommand {
	var out []commands.CustomCommand
	for _, def := range ToolRegistry {
		entry := commands.CustomCommand{
			ID:   string(def.ID),
			Name: def.Name,
			Content: "# Tool: " + def.Name + "\n" +
				"Slash: `" + def.Slash + "`\n" +
				"Description: " + def.Description + "\n",
		}
		out = append(out, entry)
	}
	return out
}
