package tools

import (
	"testing"
)

func TestToolRegistryOrder(t *testing.T) {
	if len(ToolRegistry) == 0 {
		t.Fatal("expected non-empty tool registry")
	}
}

func TestToolIDsReturnsAll(t *testing.T) {
	ids := ToolIDs()
	if len(ids) != len(ToolRegistry) {
		t.Fatalf("expected %d ids, got %d", len(ToolRegistry), len(ids))
	}
}

func TestToolIDsMatchesRegistryOrder(t *testing.T) {
	ids := ToolIDs()
	for i, id := range ids {
		if id != ToolRegistry[i].ID {
			t.Fatalf("tool ids mismatch at index %d: %s != %s", i, id, ToolRegistry[i].ID)
		}
	}
}

func TestToolSlashCommandsIncludesAliases(t *testing.T) {
	cmds := ToolSlashCommands()
	if len(cmds) < len(ToolRegistry) {
		t.Fatalf("expected at least %d slash commands, got %d", len(ToolRegistry), len(cmds))
	}
	found := false
	for _, cmd := range cmds {
		if cmd == "/finder" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected /finder alias in slash commands")
	}
}

func TestFindToolBySlash(t *testing.T) {
	def, ok := FindToolBySlash("/files")
	if !ok {
		t.Fatal("expected /files tool")
	}
	if def.ID != ToolIdFiles {
		t.Fatalf("expected files id, got %s", def.ID)
	}
	_, ok = FindToolBySlash("/does-not-exist")
	if ok {
		t.Fatal("did not expect unknown slash command")
	}
}

func TestFindToolBySlashAlias(t *testing.T) {
	def, ok := FindToolBySlash("/finder")
	if !ok {
		t.Fatal("expected /finder alias")
	}
	if def.ID != ToolIdFiles {
		t.Fatalf("expected files id via alias, got %s", def.ID)
	}
}

func TestFindToolByID(t *testing.T) {
	def, ok := FindToolByID(ToolIdHealth)
	if !ok {
		t.Fatal("expected health tool")
	}
	if def.Slash != "/health" {
		t.Fatalf("expected /health, got %s", def.Slash)
	}
}

func TestAsCustomCommands(t *testing.T) {
	cmds := AsCustomCommands()
	if len(cmds) != len(ToolRegistry) {
		t.Fatalf("expected %d commands, got %d", len(ToolRegistry), len(cmds))
	}
	for _, cmd := range cmds {
		if cmd.Name == "" {
			t.Fatal("expected non-empty command name")
		}
	}
}
