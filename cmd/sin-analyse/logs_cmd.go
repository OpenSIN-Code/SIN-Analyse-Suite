package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/OpenSIN-Code/SIN-Analyse-Suite/cmd/sin-analyse/internal/logs"
	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs [path]",
	Short: "Analyse a log file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		focus, _ := cmd.Flags().GetString("focus")
		since, _ := cmd.Flags().GetString("since")
		j, _ := cmd.Flags().GetBool("json")
		res, err := logs.Analyze(path, logs.Options{Focus: focus, Since: since})
		if err != nil {
			return err
		}
		if j {
			return json.NewEncoder(os.Stdout).Encode(res)
		}
		fmt.Println(res.ToMarkdown())
		return nil
	},
}

func init() {
	logsCmd.Flags().String("focus", "", "Focus on specific keyword")
	logsCmd.Flags().String("since", "", "Filter entries since timestamp")
	logsCmd.Flags().Bool("json", false, "Output JSON")
}
