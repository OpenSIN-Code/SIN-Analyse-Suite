package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func printJSON(v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(data)
	return err
}

func printMarkdown(md string) {
	fmt.Fprintln(os.Stdout, md)
}
