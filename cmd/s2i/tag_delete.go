package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/s2i/pkg/github"
	"github.com/spf13/cobra"
	"os"
	"regexp"
)

const pluginTagDeleteDesc = `
`

type tagDeleteCmd struct {
	Tags        []string
	SourceOwner string
	SourceRepo  string
	DryRun      bool
	Interactive bool
	Regex       bool
}

func newTagDeleteCmd() *cobra.Command {
	c := &tagDeleteCmd{}
	cmd := &cobra.Command{
		Use:     "delete <TAG>",
		Aliases: []string{"del", "rm"},
		Short:   "delete github tag and its release",
		Long:    pluginTagDeleteDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(c.SourceOwner) == 0 || len(c.SourceRepo) == 0 {
				if pwd, err := os.Getwd(); err == nil {
					owner, repo := github.Remote(logrus.StandardLogger(), pwd)
					if len(c.SourceOwner) == 0 {
						c.SourceOwner = owner
					}
					if len(c.SourceRepo) == 0 {
						c.SourceRepo = repo
					}
				}
			}
			c.Tags = args
			if c.Interactive {
				if err := tagDeleteQuestions(c); err != nil {
					return err
				}
			}
			if len := len(c.Tags); len == 0 {
				return fmt.Errorf("requires at least 1 arg(s), only received %v", len)
			}
			return c.run()
		},
	}

	f := cmd.Flags()
	f.BoolVarP(&c.Interactive, "interactive", "i", false, "interactive prompt")
	f.StringVar(&c.SourceOwner, "source-owner", c.SourceOwner, "name of the owner (user or org) of the repo to create tag")
	f.StringVar(&c.SourceRepo, "source-repo", c.SourceRepo, "name of repo to create tag")
	f.BoolVar(&c.DryRun, "dry-run", false, "simulate an tag delete \"for real\"")
	f.BoolVarP(&c.Regex, "regex", "r", false, "matches tag by regex (bad performance warning)")
	return cmd
}

func (c *tagDeleteCmd) run() error {
	if c.Regex {
		var regex []*regexp.Regexp
		for _, tag := range c.Tags {
			regex = append(regex, regexp.MustCompile(tag))
		}
		return github.DeleteReleasesAndTagsByRegex(logrus.StandardLogger(), token, c.SourceOwner, c.SourceRepo, regex, c.DryRun)
	}
	return github.DeleteReleasesAndTags(logrus.StandardLogger(), token, c.SourceOwner, c.SourceRepo, c.Tags, c.DryRun)
}
