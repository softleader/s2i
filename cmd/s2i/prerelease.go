package main

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/softleader/s2i/pkg/deployer"
	"github.com/softleader/s2i/pkg/docker"
	"github.com/softleader/s2i/pkg/github"
	"github.com/softleader/s2i/pkg/jib"
	"github.com/softleader/s2i/pkg/test"
	"github.com/spf13/cobra"
	"os"
)

const pluginPrereleaseDesc = `Draft a pre-release to SoftLeader docker swarm ecosystem

建立 pre-release 版本, pre 為此 command 的縮寫, 傳入 '--interactive' 可以開啟互動模式
在互動模式下, tag 若不傳入就會自動的到 GitHub 找出 latest release 並增加一個 patch 版號做為問答預設的 tag:

	$ s2i prerelease TAG
	$ s2i pre -i

pre-release 必須指定 stage, 預設為 '0', 基本上是建議:

	- 0 for alpha
	- 1 for beta
	- 2 for release candidate

你可以透過 '--stage' 調整, 可以傳入任意字串:

	$ s2i pre TAG --stage do.not.use

s2i 會試著從當前目錄收集專案資訊, 你都可以自行傳入做調整:

	- git 資訊: '--source-owner', '--source-repo' 及 '--source-branch'
	- jib 資訊: '--jib-auth-username' 及 '--jib-auth-password'

傳入 '--service-id' 即可在最後自動的更新 SoftLeader Deployer 上的服務
當然你必須先到 SoftLeader Deployer (http://softleader.com.tw:5678) 上查出要更新的 Service ID
或是開啟互動模式來協助你選到 Service ID:

	$ s2i pre TAG --service-id SERVICE_ID

可以使用 '--help' 查看所有選項及其詳細說明

	$ s2i pre -h
`

type prereleaseCmd struct {
	Force           bool
	Interactive     bool
	SourceOwner     string
	SourceRepo      string
	SourceBranch    string
	SkipTests       bool
	UpdateSnapshots bool
	ConfigServer    string
	ConfigLabel     string
	Image           *docker.SoftleaderHubImage
	Stage           string
	Deployer        string
	Auth            *jib.Auth
	ServiceID       string
}

func newPrereleaseCmd() *cobra.Command {
	c := &prereleaseCmd{
		Auth:  &jib.Auth{},
		Image: &docker.SoftleaderHubImage{},
	}
	cmd := &cobra.Command{
		Use:     "prerelease <TAG>",
		Aliases: []string{"pre"},
		Short:   "draft a pre-release version",
		Long:    pluginPrereleaseDesc,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !c.Interactive && len(args) < 1 {
				return errors.New(`accepts 1 arg(s), received 0`)
			}
			if len(args) > 0 {
				c.Image.Tag = args[0]
			}
			if pwd, err := os.Getwd(); err == nil {
				c.SourceOwner, c.SourceRepo = github.Remote(logrus.StandardLogger(), pwd)
				c.Image.Name = c.SourceRepo
				c.SourceBranch = github.Head(logrus.StandardLogger(), pwd)
				c.Auth = jib.GetAuth(logrus.StandardLogger(), pwd)
			}
			if c.Interactive {
				if c.Image.Tag == "" {
					var err error
					c.Image.Tag, err = github.FindNextReleaseVersion(logrus.StandardLogger(), token, c.SourceOwner, c.SourceRepo)
					if err != nil {
						logrus.Debugln(err)
					}
				}
				if err := prereleaseQuestions(c); err != nil {
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
	f.BoolVarP(&c.Force, "force", "f", false, "force to delete the tag if it already exists")
	f.BoolVarP(&c.Interactive, "interactive", "i", false, "interactive prompt")
	f.BoolVar(&c.SkipTests, "skip-tests", false, "skip tests when building image")
	f.BoolVarP(&c.UpdateSnapshots, "update-snapshots", "U", false, "force updated snapshots on remote repositories")
	f.StringVar(&c.SourceOwner, "source-owner", c.SourceOwner, "name of the owner (user or org) of the repo to create tag")
	f.StringVar(&c.SourceRepo, "source-repo", c.SourceRepo, "name of repo to create tag")
	f.StringVar(&c.SourceBranch, "source-branch", c.SourceBranch, "name of branch to create tag")
	f.StringVar(&c.ConfigServer, "config-server", "http://softleader.com.tw:8887", "config server to run the test")
	f.StringVar(&c.ConfigLabel, "config-label", "", "the label of config server to run the test, e.g. sqlServer")
	f.StringVar(&c.Image.Name, "image", c.Image.Name, "name of image to build")
	f.StringVar(&c.Stage, "stage", "0", "designating development stage to build, e.g. 0 for alpha, 1 for beta, 2 for release candidate")
	f.StringVar(&c.Deployer, "deployer", "http://softleader.com.tw:5678", "deployer to deploy")
	f.StringVar(&c.Auth.Username, "jib-auth-username", "", "username of docker registry for jib")
	f.StringVar(&c.Auth.Password, "jib-auth-password", "", "password of docker registry for jib")
	f.StringVar(&c.ServiceID, "service-id", "", "docker swarm service id to update")
	return cmd
}

func (c *prereleaseCmd) run() error {
	if !c.SkipTests {
		if err := test.Run(logrus.StandardLogger(), c.ConfigServer, c.ConfigLabel, c.UpdateSnapshots); err != nil {
			return err
		}
	}
	c.Image.SetPreRelease(c.Stage)
	if c.Auth.IsValid() {
		if err := jib.Build(logrus.StandardLogger(), c.Image, c.Auth, c.UpdateSnapshots); err != nil {
			return err
		}
	} else {
		// 當沒提供 docker registry auth 資訊時, 我們就 build 到 local docker daemon 再推
		// 因為使用者可能已經在 local 的 docker daemon 登入過 hub.softleader.com.tw
		if err := jib.DockerBuild(logrus.StandardLogger(), c.Image, c.UpdateSnapshots); err != nil {
			return err
		}
		if err := docker.Push(logrus.StandardLogger(), c.Image); err != nil {
			return err
		}
		if err := docker.Rmi(logrus.StandardLogger(), c.Image); err != nil {
			return err
		}
	}
	if err := github.CreatePrerelease(logrus.StandardLogger(), token, c.SourceOwner, c.SourceRepo, c.SourceBranch, c.Image.Tag, c.Force); err != nil {
		return err
	}
	if c.ServiceID != "" {
		if err := deployer.UpdateService(logrus.StandardLogger(), "s2i", metadata.String(), c.Deployer, c.ServiceID, c.Image); err != nil {
			return err
		}
	}
	logrus.Printf("Everything is all set, you are good to go.")
	return nil
}
