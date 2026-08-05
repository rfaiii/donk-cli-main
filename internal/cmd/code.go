package cmd

import (
	"fmt"
	"log/slog"
	"os"

	tea "charm.land/bubbletea/v2"
	uv "github.com/charmbracelet/ultraviolet"
	"github.com/richavery/donk-cli/internal/event"
	"github.com/richavery/donk-cli/internal/ui/common"
	"github.com/richavery/donk-cli/internal/ui/model"
	"github.com/spf13/cobra"
)

var codeCmd = &cobra.Command{
	Use:   "code",
	Short: "Start a coder-only session",
	Long:  `Start a DONK session in coder mode with coding-focused tools and system prompt.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ws, cleanup, err := setupWorkspaceWithProgressBar(cmd)
		if err != nil {
			return err
		}
		defer cleanup()

		event.AppInitialized()

		if err := ws.InitCoderAgent(cmd.Context()); err != nil {
			return fmt.Errorf("failed to initialize coder agent: %w", err)
		}

		com := common.DefaultCommon(ws)
		m := model.New(com, "", false)

		var env uv.Environ = os.Environ()
		program := tea.NewProgram(
			m,
			tea.WithEnvironment(env),
			tea.WithContext(cmd.Context()),
		)
		go ws.Subscribe(program)

		if _, err := program.Run(); err != nil {
			event.Error(err)
			slog.Error("TUI run error", "error", err)
			return fmt.Errorf("Donk crashed: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(codeCmd)
}
