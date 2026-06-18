package main

import (
	"os"
	"encoding/json"
	"fmt"
	"github.com/OpenSIN-Code/SIN-Analyse-Suite/cmd/sin-analyse/internal/data"
	"github.com/spf13/cobra"
)

var dataCmd = &cobra.Command{
	Use:   "data [path]",
	Short: "Analyse a data file (CSV)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		ddl, _ := cmd.Flags().GetBool("ddl")
		sample, _ := cmd.Flags().GetInt("sample")
		j, _ := cmd.Flags().GetBool("json")
		res, err := data.Analyze(path, data.Options{DDL: ddl, SampleSize: sample})
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
	dataCmd.Flags().Bool("ddl", false, "Generate CREATE TABLE DDL")
	dataCmd.Flags().Int("sample", 5, "Number of sample rows")
	dataCmd.Flags().Bool("json", false, "Output JSON")
}
