package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type versionCmd struct {
	full bool
}

func newVersionCmd() *cobra.Command {
	c := &versionCmd{}
	cmd := &cobra.Command{
		Use:   "version",
		Short: "print s2i version",
		Long:  "print s2i version",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.run()
		},
	}

	f := cmd.Flags()
	f.BoolVar(&c.full, "full", false, "print full version number and commit hash")

	return cmd
}

func (c *versionCmd) run() error {
	if c.full {
		logrus.Infoln(metadata.FullString())
	} else {
		logrus.Infoln(metadata.String())
	}
	return nil
}
