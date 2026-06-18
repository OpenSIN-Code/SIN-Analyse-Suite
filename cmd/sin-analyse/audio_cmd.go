package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/OpenSIN-Code/SIN-Analyse-Suite/cmd/sin-analyse/internal/audio"
	"github.com/spf13/cobra"
)

var audioCmd = &cobra.Command{
	Use:   "audio [path]",
	Short: "Analyse an audio file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		lang, _ := cmd.Flags().GetString("language")
		j, _ := cmd.Flags().GetBool("json")
		res, err := audio.Analyze(path, audio.Options{Language: lang})
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
	audioCmd.Flags().String("language", "auto", "Source language (auto, en, de, fr)")
	audioCmd.Flags().Bool("json", false, "Output JSON")
}
