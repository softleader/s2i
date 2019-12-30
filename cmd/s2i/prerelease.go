package main

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/softleader/s2i/pkg/deployer"
	"github.com/softleader/s2i/pkg/docker"
	"github.com/softleader/s2i/pkg/github"
	"github.com/softleader/s2i/pkg/jib"
	"github.com/softleader/s2i/pkg/mvn"
	"github.com/softleader/s2i/pkg/slack"
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

如果你當前的專案並非 maven 專案 (如 nodejs), 請務必使用 multi-stage builds 來建構 source code
(https://docs.docker.com/develop/develop-images/multistage-build/)
s2i 會自動判斷 multi-stage build 等專案建構條件, 在 jib 及 docker 之間自動的挑選 shipping source 的策略
你也可以傳入 '--ship-source' 來指定策略:

	- 0 for auto-detect (default)
	- 1 for jib
	- 2 for docker

	$ s2i pre TAG -S 1

可以使用 '--help' 查看所有選項及其詳細說明

	$ s2i pre -h
`

type prereleaseCmd struct {
	Force           bool
	interactive     bool
	promptSize      int
	SourceOwner     string `yaml:"source-owner"`
	SourceRepo      string `yaml:"source-repo"`
	SourceBranch    string `yaml:"source-branch"`
	SkipTests       bool   `yaml:"skip-tests"`
	SkipDraft       bool   `yaml:"skip-draft"`
	UpdateSnapshots bool   `yaml:"update-snapshots"`
	ConfigServer    string `yaml:"config-server"`
	ConfigLabel     string `yaml:"config-label"`
	Image           *docker.SoftleaderHubImage
	Stage           string
	Deployer        string
	Auth            *jib.Auth
	ServiceID       string `yaml:"service-id"`
	ShipStrategy    int    `yaml:"build-strategy"`
	SkipSlack       bool   `yaml:"skip-slack"`
	SlackWebhookURL string `yaml:"slack-webhook-url"`
	pwd             string
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
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if !c.interactive && len(args) < 1 {
				return errors.New(`accepts 1 arg(s), received 0`)
			}
			if len(args) > 0 {
				c.Image.Tag = args[0]
			}
			if c.pwd, err = os.Getwd(); err == nil {
				var t string
				t, c.SourceOwner, c.SourceRepo = github.Remote(logrus.StandardLogger(), c.pwd)
				if len(t) != 0 { // 代表此 repo 是用指定 token clone 的, 因此換掉這次 global 的 token
					token = t
				}
				c.Image.Name = c.SourceRepo
				c.SourceBranch = github.Head(logrus.StandardLogger(), c.pwd)
				c.Auth = jib.GetAuth(logrus.StandardLogger(), c.pwd)
			}
			if c.interactive {
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
	f.BoolVarP(&c.interactive, "interactive", "i", false, "interactive prompt")
	f.IntVar(&c.promptSize, "interactive-prompt-size", 7, "interactive prompt size")
	f.BoolVar(&c.SkipTests, "skip-tests", false, "skip tests when building image")
	f.BoolVar(&c.SkipDraft, "skip-draft", false, "skip draft pre-release tag")
	f.BoolVarP(&c.UpdateSnapshots, "update-snapshots", "U", false, "force to check for updated snapshots on remote repositories")
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
	f.IntVarP(&c.ShipStrategy, "ship-strategy", "S", 0, "specify how to ship source, 0 for auto-detect, 1 for jib, 2 for docker")
	f.BoolVar(&c.SkipSlack, "skip-slack", false, "skip slack webhook")
	f.StringVar(&c.SlackWebhookURL, "slack-webhook-url", "", "slack webhook url, will override the webhook url cache")
	return cmd
}

func (c *prereleaseCmd) run() (err error) {
	if !c.SkipTests {
		if err := mvn.Test(logrus.StandardLogger(), c.ConfigServer, c.ConfigLabel, c.UpdateSnapshots); err != nil {
			return err
		}
	}
	c.Image.SetPreRelease(c.Stage)

	if err := c.ship(); err != nil {
		return err
	}

	var release *github.Release
	if !c.SkipDraft {
		if release, err = github.CreatePrerelease(logrus.StandardLogger(), token, c.SourceOwner, c.SourceRepo, c.SourceBranch, c.Image.Tag, c.Force); err != nil {
			return err
		}
	}
	if c.ServiceID != "" {
		if err := deployer.UpdateService(logrus.StandardLogger(), "s2i", metadata.String(), c.Deployer, c.ServiceID, c.Image); err != nil {
			return err
		}
	}
	if !c.SkipSlack {
		if err := slack.Post(logrus.StandardLogger(), metadata, release, c.SlackWebhookURL, c.Image); err != nil {
			logrus.Debugf("failed posting slack webhook: %s", err)
		}
	}
	logrus.Printf("Everything is all set, you are good to go.")
	return nil
}

func (c *prereleaseCmd) ship() error {
	if c.ShipStrategy == 1 { // jib
		return c.jibRelease()
	}

	if c.ShipStrategy == 2 { // docker
		return c.dockerRelease()
	}

	// auto
	if err := c.jibRelease(); err == nil {
		return nil
	}
	return c.dockerRelease()
}

func (c *prereleaseCmd) jibRelease() error {
	if c.Auth.IsValid() {
		return jib.Build(logrus.StandardLogger(), c.Image, c.Auth, c.UpdateSnapshots)
	}
	// 當沒提供 docker registry auth 資訊時, 我們就 build 到 local docker daemon 再推
	// 因為使用者可能已經在 local 的 docker daemon 登入過 hub.softleader.com.tw
	if err := jib.DockerBuild(logrus.StandardLogger(), c.Image, c.UpdateSnapshots); err != nil {
		return err
	}
	return c.dockerPublish()
}

func (c *prereleaseCmd) dockerRelease() error {
	if !docker.ContainsMultiStageBuilds(logrus.StandardLogger(), c.pwd) { // 如果不是 multi-stage build 才 build build 看
		if err := mvn.Package(logrus.StandardLogger(), c.UpdateSnapshots); err != nil {
			return err
		}
	}
	if err := docker.Build(logrus.StandardLogger(), c.Image); err != nil {
		return err
	}
	return c.dockerPublish()
}

func (c *prereleaseCmd) dockerPublish() error {
	if err := docker.Push(logrus.StandardLogger(), c.Image); err != nil {
		return err
	}
	return docker.Rmi(logrus.StandardLogger(), c.Image)
}
