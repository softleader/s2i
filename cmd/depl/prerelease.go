package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/depl/pkg/deployer"
	"github.com/softleader/depl/pkg/github"
	"github.com/softleader/depl/pkg/jib"
	"github.com/softleader/depl/pkg/test"
	"github.com/spf13/cobra"
	"os"
)

const pluginPrereleaseDesc = `Draft a pre-release to SoftLeader docker swarm ecosystem

建立 pre-release 版本, pre 為此 command 的縮寫, 傳入 '--interactive' 可以開啟互動式指令:

	$ depl prerelease TAG
	$ depl pre TAG -i

pre-release 版本的 tag 必須帶著 stage, 預設的 stage 預設為 '0', 基本上是建議:

	- 0 for alpha
	- 1 for beta
	- 2 for release candidate

你可以透過 '--stage' 調整, stage 可以給任意字串: 

	$ depl pre TAG --stage do.not.use

depl 會試著從當前目錄收集專案資訊, 你都可以自行傳入做調整:

	- git 資訊: '--sourceOwner', '--sourceRepo' 及 '--sourceBranch'
	- jib 資訊: '--jib-auth-username' 及 '--jib-auth-password'

傳入 '--docker-service-id' 即可在最後自動的更新 SoftLeader Deployer (http://softleader.com.tw:5678) 上的服務
當然你必須先到 Deployer 上查出該 service id:

	$ depl pre TAG --docker-service-id 0989olwerft

可以使用 '--help' 查看所有選項及其詳細說明

	$ depl pre -h
`

type prereleaseCmd struct {
	Force           bool
	Interactive     bool
	SourceOwner     string
	SourceRepo      string
	SourceBranch    string
	SkipTests       bool
	ConfigServer    string
	ConfigLabel     string
	Image           string
	Tag             string
	Stage           string
	Deployer        string
	Auth            *jib.Auth
	DockerServiceID string
}

func newPrereleaseCmd() *cobra.Command {
	c := &prereleaseCmd{
		Auth: &jib.Auth{},
	}
	cmd := &cobra.Command{
		Use:     "prerelease <TAG>",
		Aliases: []string{"pre"},
		Short:   "draft a pre-release version",
		Long:    pluginPrereleaseDesc,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c.Tag = args[0]
			if pwd, err := os.Getwd(); err == nil {
				c.SourceOwner, c.SourceRepo = github.Remote(logrus.StandardLogger(), pwd)
				c.Image = c.SourceRepo
				c.SourceBranch = github.Head(logrus.StandardLogger(), pwd)
				c.Auth = jib.GetAuth(logrus.StandardLogger(), pwd)
			}
			if c.Interactive {
				if err := c.prompt(); err != nil {
					return err
				}
			}
			return c.run()
		},
	}

	f := cmd.Flags()
	f.BoolVarP(&c.Force, "force", "f", false, "force to delete the tag if it already exists")
	f.BoolVarP(&c.Interactive, "interactive", "i", false, "interactive prompt")
	f.BoolVar(&c.SkipTests, "skip-tests", false, "skip tests when building image")
	f.StringVar(&c.SourceOwner, "source-owner", c.SourceOwner, "name of the owner (user or org) of the repo to create tag")
	f.StringVar(&c.SourceRepo, "source-repo", c.SourceRepo, "name of repo to create to create tag")
	f.StringVar(&c.SourceBranch, "source-branch", c.SourceBranch, "name of branch to create to create tag")
	f.StringVar(&c.ConfigServer, "config-server", "http://192.168.1.88:8887", "config server to run the test")
	f.StringVar(&c.ConfigLabel, "config-label", "", "the label of config server to run the test, e.g. sqlServer")
	f.StringVar(&c.Image, "image", c.Image, "name of image to build")
	f.StringVar(&c.Stage, "stage", "0", "designating development stage to build, e.g. 0 for alpha, 1 for beta, 2 for release candidate")
	f.StringVar(&c.Deployer, "deployer", "http://softleader.com.tw:5678", "deployer to deploy")
	f.StringVar(&c.Auth.Username, "jib-auth-username", "dev", "username of docker registry for jib to build")
	f.StringVar(&c.Auth.Password, "jib-auth-password", "sleader", "password of docker registry for jib to build")
	f.StringVar(&c.DockerServiceID, "docker-service-id", "", "docker service id to update image")
	return cmd
}

func (c *prereleaseCmd) run() error {
	if !c.SkipTests {
		if err := test.Run(logrus.StandardLogger(), c.ConfigServer, c.ConfigLabel); err != nil {
			return err
		}
	}
	tagName := fmt.Sprintf("%s-%s", c.Tag, c.Stage)
	if err := jib.Build(logrus.StandardLogger(), c.Image, tagName, c.Auth); err != nil {
		return err
	}
	if err := github.CreatePrerelease(logrus.StandardLogger(), token, c.SourceOwner, c.SourceRepo, c.SourceBranch, tagName, c.Force); err != nil {
		return err
	}
	if c.DockerServiceID != "" {
		if err := deployer.UpdateService(logrus.StandardLogger(), "depl", metadata.String(), c.Deployer, c.DockerServiceID, c.Image, tagName); err != nil {
			return err
		}
	}
	logrus.Printf("Everything is all set, you are good to go.")
	return nil
}
