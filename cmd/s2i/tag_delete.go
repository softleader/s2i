package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/s2i/pkg/github"
	"github.com/spf13/cobra"
	"os"
	"regexp"
)

const pluginTagDeleteDesc = `刪除 tag 及其 release

傳入 '--interactive' 可以開啟互動模式

	$ s2i tag delete TAG..
	$ s2i tag delete TAG.. -i

s2i 會試著從當前目錄收集專案資訊, 你都可以自行傳入做調整:

	- git 資訊: '--source-owner', '--source-repo'

傳入 '--regex' 將以 regular expression 方式過濾 match 的 tag, 並刪除之

	$ slctl s2i tag delete REGEX.. -r

傳入 '--dry-run' 將 "模擬" 刪除, 不會真的作用到 GitHub 上, 通常可用於檢視 regex 是否如預期

	$ slctl s2i tag delete REGEX... -r --dry-run

Example:

	# 以互動的問答方式, 詢問所有可控制的問題
	slctl s2i tag delete -i

	# 在當前目錄的專案中, 刪除名稱 1.0.0 及 1.1.0 的 tag 及 release (完整比對)
	$ slctl s2i tag delete 1.0.0 1.1.0

	# 在當前目錄的專案中, "模擬" 刪除所有名稱為 1 開頭或 2 開頭的 tag 及其 release 
	$ slctl s2i tag delete ^1 ^2 -r --dry-run

	# 刪除指定專案 github.com/me/my-repo 的所有 tag 及其 release
	$ slctl s2i tag delete .+ -r --source-owner me --source-repo my-repo
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
