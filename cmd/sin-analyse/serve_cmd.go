package main

import (
	"log"

	"github.com/OpenSIN-Code/SIN-Analyse-Suite/cmd/sin-analyse/internal/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start MCP server",
	RunE: func(cmd *cobra.Command, args []string) error {
		s := server.NewMCPServer("SIN-Analyse-Suite", "1.0.0")
		if err := mcp.RegisterAll(s); err != nil {
			return err
		}
		log.Println("SIN-Analyse-Suite MCP server ready (stdio)")
		return server.ServeStdio(s)
	},
}
