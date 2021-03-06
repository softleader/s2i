package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/s2i/pkg/github"
	"github.com/spf13/cobra"
	"os"
)

const pluginTagDeleteDesc = `刪除 tag 及其 release

傳入 '--interactive' 可以開啟互動模式

	$ s2i tag delete TAG..
	$ s2i tag delete TAG.. -i

s2i 會試著從當前目錄收集專案資訊, 你都可以自行傳入做調整:

	- git 資訊: '--source-owner', '--source-repo'

傳入 '--regex' 將以 regular expression 方式模糊過濾 tag, 並刪除之

	$ slctl s2i tag delete REGEX.. -r

傳入 '--semver' 將以 semantic versioning 2.0.0 方式模糊過濾 tag, 並刪除之
建議可查看 https://devhints.io/semver

	$ slctl s2i tag delete RANGE.. -s

傳入 '--dry-run' 將 "模擬" 刪除, 不會真的作用到 GitHub 上, 通常可用於檢視 regex 是否如預期

	$ slctl s2i tag delete RANGE... -s --dry-run

模糊過濾 flag ('-r' 或 '-s' 等) 使用上請注意: 
- 將會 scan 所有 GitHub 上所有的 tag, 效能自然會比完全比對 tag 來得差
- 判斷先後順序依序為: '-r', '-s'

Example:

	# 以互動的問答方式, 詢問所有可控制的問題
	$ slctl s2i tag delete -i

	# 在當前目錄的專案中, 刪除名稱 1.0.0 及 1.1.0 的 tag 及 release (完整比對)
	$ slctl s2i tag delete 1.0.0 1.1.0

	# 在當前目錄的專案中, "模擬" 刪除所有名稱為 1 開頭或 2 開頭的 tag 及其 release 
	$ slctl s2i tag delete ^1 ^2 -r --dry-run

	# 在當前目錄的專案中, "模擬" 刪除所有小於 2.5.x 開頭的 tag 及其 release 
	$ slctl s2i tag delete "<2.5.x" -s --dry-run

	# 刪除指定專案 github.com/me/my-repo 的所有 tag 及其 release
	$ slctl s2i tag delete .+ -r --source-owner me --source-repo my-repo
`

type tagDeleteCmd struct {
	Tags                      []string
	SourceOwner               string `yaml:"source-owner"`
	SourceRepo                string `yaml:"source-repo"`
	DryRun                    bool   `yaml:"dry-run"`
	Interactive               bool
	github.TagMatcherStrategy `yaml:"tag-matcher-strategy"`
}

func newTagDeleteCmd() *cobra.Command {
	c := &tagDeleteCmd{}
	cmd := &cobra.Command{
		Use:     "delete <TAG>",
		Aliases: []string{"del", "rm"},
		Short:   "delete tag and its release on GitHub",
		Long:    pluginTagDeleteDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(c.SourceOwner) == 0 || len(c.SourceRepo) == 0 {
				if pwd, err := os.Getwd(); err == nil {
					t, owner, repo := github.Remote(logrus.StandardLogger(), pwd)
					if len(c.SourceOwner) != 0 { // 代表此 repo 是用指定 token clone 的, 因此換掉這次 global 的 token
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
	f.StringVar(&c.SourceOwner, "source-owner", c.SourceOwner, "name of the owner (user or org) of the repo to delete tag")
	f.StringVar(&c.SourceRepo, "source-repo", c.SourceRepo, "name of repo to delete tag")
	f.BoolVar(&c.DryRun, "dry-run", false, "simulate tag deletion \"for real\"")
	f.BoolVarP(&c.Regex, "regex", "r", false, "matches tag by regex (bad performance warning, it'll scan over all tags of the repo)")
	f.BoolVarP(&c.SemVer, "semver", "s", false, "matches tag by semantic versioning (bad performance warning, it'll scan over all tags of the repo)")
	return cmd
}

func (c *tagDeleteCmd) run() error {
	if c.Regex {
		matcher, err := github.NewRegexMatcher(c.Tags)
		if err != nil {
			return err
		}
		return github.DeleteMatchesReleasesAndTags(logrus.StandardLogger(), token, c.SourceOwner, c.SourceRepo, matcher, c.DryRun)
	}
	if c.SemVer {
		matcher, err := github.NewSemVerMatcher(c.Tags)
		if err != nil {
			return err
		}
		return github.DeleteMatchesReleasesAndTags(logrus.StandardLogger(), token, c.SourceOwner, c.SourceRepo, matcher, c.DryRun)
	}
	return github.DeleteReleasesAndTags(logrus.StandardLogger(), token, c.SourceOwner, c.SourceRepo, c.Tags, c.DryRun)
}
