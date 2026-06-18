package main

import (
	"os"
	"encoding/json"
	"fmt"
	"github.com/OpenSIN-Code/SIN-Analyse-Suite/cmd/sin-analyse/internal/video"
	"github.com/spf13/cobra"
)

var videoCmd = &cobra.Command{
	Use:   "video [path]",
	Short: "Analyse a video file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		interval, _ := cmd.Flags().GetInt("keyframe-interval")
		j, _ := cmd.Flags().GetBool("json")
		res, err := video.Analyze(path, video.Options{KeyframeInterval: interval})
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
	videoCmd.Flags().Int("keyframe-interval", 30, "Keyframe interval (frames)")
	videoCmd.Flags().Bool("json", false, "Output JSON")
}
