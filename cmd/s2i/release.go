package main

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/softleader/s2i/pkg/docker"
	"github.com/softleader/s2i/pkg/github"
	"github.com/softleader/s2i/pkg/jenkins"
	"github.com/spf13/cobra"
	"os"
)

const pluginReleaseDesc = `Draft a release to SoftLeader docker swarm ecosystem

建立 release 版本, 傳入 '--interactive' 可以開啟互動模式
在互動模式下, tag 若不傳入就會自動的到 GitHub 找出 latest release 並增加一個 patch 版號做為問答預設的 tag:

	$ s2i release TAG
	$ s2i release -i

s2i 會試著從當前目錄收集專案資訊, 你都可以自行傳入做調整:

	- git 資訊: '--source-owner', '--source-repo' 及 '--source-branch'

可以使用 '--help' 查看所有選項及其詳細說明

	$ s2i release -h
`

type releaseCmd struct {
	Interactive  bool
	SourceOwner  string
	SourceRepo   string
	SourceBranch string
	Image        *docker.SoftleaderHubImage
	Jenkins      string
}

func newReleaseCmd() *cobra.Command {
	c := &releaseCmd{
		Image: &docker.SoftleaderHubImage{},
	}
	cmd := &cobra.Command{
		Use:   "release <TAG>",
		Short: "draft a release version",
		Long:  pluginReleaseDesc,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !c.Interactive && len(args) < 1 {
				return errors.New(`accepts 1 arg(s), received 0`)
			}
			if pwd, err := os.Getwd(); err == nil {
				c.SourceOwner, c.SourceRepo = github.Remote(logrus.StandardLogger(), pwd)
				c.Image.Name = c.SourceRepo
				c.SourceBranch = github.Head(logrus.StandardLogger(), pwd)
			}
			if c.Interactive {
				if len(args) < 1 {
					var err error
					c.Image.Tag, err = github.FindNextReleaseVersion(logrus.StandardLogger(), token, c.SourceOwner, c.SourceRepo)
					if err != nil {
						logrus.Debugln(err)
					}
				} else {
					c.Image.Tag = args[0]
				}
				if err := releaseQuestions(c); err != nil {
					return err
				}
			}
			if err := c.Image.CheckValid(); err != nil {
				return err
			}
			return c.run()
		},
	}

	f := cmd.Flags()
	f.BoolVarP(&c.Interactive, "interactive", "i", false, "interactive prompt")
	f.StringVar(&c.SourceOwner, "source-owner", c.SourceOwner, "name of the owner (user or org) of the repo to create tag")
	f.StringVar(&c.SourceRepo, "source-repo", c.SourceRepo, "name of repo to create tag")
	f.StringVar(&c.SourceBranch, "source-branch", c.SourceBranch, "name of branch to create tag")
	f.StringVar(&c.Image.Name, "image", c.Image.Name, "name of image to build")
	f.StringVar(&c.Jenkins, "jenkins", "https://jenkins.softleader.com.tw", "jenkins to run the pipeline")
	return cmd
}

func (c *releaseCmd) run() error {
	if err := github.CreateRelease(logrus.StandardLogger(), token, c.SourceOwner, c.SourceRepo, c.SourceBranch, c.Image.Tag); err != nil {
		return err
	}

	jenkins := jenkins.NewClient(c.Jenkins).
		SetVerbose(verbose).
		SetLogger(logrus.StandardLogger())
	params := make(map[string]string)
	params["tag"] = c.Image.Tag
	if err := jenkins.Job().BuildWithParameters(c.SourceRepo, params); err != nil {
		return err
	}

	logrus.Printf("Everything is all set, you can check the progress at: %s/job/%s", c.Jenkins, c.SourceRepo)
	return nil
}
