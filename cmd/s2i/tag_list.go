package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/s2i/pkg/github"
	"github.com/spf13/cobra"
	"os"
)

const pluginTagListDesc = `列出 tag 名稱, 發佈時間及發佈人員

傳入 '--interactive' 可以開啟互動模式

	$ s2i tag list TAG..
	$ s2i tag list TAG.. -i

s2i 會試著從當前目錄收集專案資訊, 你都可以自行傳入做調整:

	- git 資訊: '--source-owner', '--source-repo'

傳入 '--regex' 將以 regular expression 方式過濾 match 的 tag, 並列出之

	$ slctl s2i tag list ^1. -r
`

type tagListCmd struct {
	Tags        []string
	SourceOwner string
	SourceRepo  string
	Interactive bool
	Regex       bool
	SemVer      bool
}

func newTagListCmd() *cobra.Command {
	c := &tagListCmd{}
	cmd := &cobra.Command{
		Use:   "list <TAG...>",
		Short: "list tags on GitHub",
		Long:  pluginTagListDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(c.SourceOwner) == 0 || len(c.SourceRepo) == 0 {
				if pwd, err := os.Getwd(); err == nil {
					t, owner, repo := github.Remote(logrus.StandardLogger(), pwd)
					if len(t) != 0 { // 代表此 repo 是用指定 token clone 的, 因此換掉這次 global 的 token
						token = t
					}
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
				if err := tagListQuestions(c); err != nil {
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
	f.StringVar(&c.SourceOwner, "source-owner", c.SourceOwner, "name of the owner (user or org) of the repo to list tag")
	f.StringVar(&c.SourceRepo, "source-repo", c.SourceRepo, "name of repo to list tag")
	f.BoolVarP(&c.Regex, "regex", "r", false, "matches tag by regex (bad performance warning, it'll scan over all tags of the repo)")
	f.BoolVarP(&c.SemVer, "semver", "s", false, "matches tag by semantic versioning (bad performance warning, it'll scan over all tags of the repo)")
	return cmd
}

func (c *tagListCmd) run() error {
	if c.Regex {
		matcher := github.NewRegexMatcher(c.Tags)
		return github.ListReleaseByMatcher(logrus.StandardLogger(), token, c.SourceOwner, c.SourceRepo, matcher)
	}
	if c.SemVer {
		matcher := github.NewSemVerMatcher(c.Tags)
		return github.ListReleaseByMatcher(logrus.StandardLogger(), token, c.SourceOwner, c.SourceRepo, matcher)
	}
	return github.ListRelease(logrus.StandardLogger(), token, c.SourceOwner, c.SourceRepo, c.Tags)
}
