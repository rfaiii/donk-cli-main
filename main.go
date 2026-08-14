// Package main is the entry point for the DONK CLI.
//
//	@title			DONK API
//	@version		1.1.5
//	@description	DONK is a terminal-based AI coding assistant. This API is served over a Unix socket (or Windows named pipe) and provides programmatic access to workspaces, sessions, agents, LSP, MCP, and more.
//	@contact.name	DONK
//	@contact.url	https://donk-cli.com
//	@license.name	MIT
//	@license.url	https://github.com/richavery/donk-cli/blob/main/LICENSE
//	@BasePath		/v1
package main

import (
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/richavery/donk-cli/internal/cmd"
	_ "github.com/richavery/donk-cli/internal/dns"
)

func main() {
	if os.Getenv("DONK_PROFILE") != "" {
		go func() {
			slog.Info("Serving pprof at localhost:6060")
			if httpErr := http.ListenAndServe("localhost:6060", nil); httpErr != nil {
				slog.Error("Failed to pprof listen", "error", httpErr)
			}
		}()
	}

	cmd.Execute()
}
