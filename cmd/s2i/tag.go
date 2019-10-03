package main

import (
	"github.com/spf13/cobra"
)

const pluginTagDesc = `
`

func neTagCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tag",
		Aliases: []string{"pre"},
		Short:   "tag",
		Long:    pluginTagDesc,
	}
	cmd.AddCommand(
		newTagDeleteCmd(),
	)
	return cmd
}
