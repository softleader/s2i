package main

import (
	"github.com/spf13/cobra"
)

const pluginTagDesc = `更方便的管理 GitHub 上的 tags
`

func neTagCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tag",
		Short: "manage tags on GitHub",
		Long:  pluginTagDesc,
	}
	cmd.AddCommand(
		newTagListCmd(),
		newTagDeleteCmd(),
	)
	return cmd
}
