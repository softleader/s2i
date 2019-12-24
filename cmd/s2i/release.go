package main

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/softleader/s2i/pkg/docker"
	"github.com/softleader/s2i/pkg/github"
	"github.com/softleader/s2i/pkg/jenkins"
	"github.com/softleader/s2i/pkg/slack"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

const pluginReleaseDesc = `Draft a release to SoftLeader docker swarm ecosystem

建立 release 版本, 傳入 '--interactive' 可以開啟互動模式
在互動模式下, tag 若不傳入就會自動的到 GitHub 找出 latest release 並增加一個 patch 版號做為問答預設的 tag:

	$ s2i release TAG
	$ s2i release -i

s2i 會試著從當前目錄收集專案資訊, 你都可以自行傳入做調整:

	- git 資訊: '--source-owner', '--source-repo' 及 '--source-branch'

傳入 '--service-id' 即可一併將要更新的 Service ID 傳給 Jenkins Pipeline
當然你必須先到 SoftLeader Deployer (http://softleader.com.tw:5678) 上查出要更新的 Service ID
或是開啟互動模式來協助你選到 Service ID:

	$ s2i release TAG --service-id SERVICE_ID

可以使用 '--help' 查看所有選項及其詳細說明

	$ s2i release -h
`

var (
	hook = regexp.MustCompile(`params.serviceID`)
)

type releaseCmd struct {
	interactive     bool
	promptSize      int
	SourceOwner     string `yaml:"source-owner"`
	SourceRepo      string `yaml:"source-repo"`
	SourceBranch    string `yaml:"source-branch"`
	Image           *docker.SoftleaderHubImage
	Jenkins         string
	Deployer        string
	ServiceID       string `yaml:"service-id"`
	SkipSlack       bool   `yaml:"skip-slack"`
	SlackWebhookURL string `yaml:"slack-webhook-url"`
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
			if !c.interactive && len(args) < 1 {
				return errors.New(`accepts 1 arg(s), received 0`)
			}
			if len(args) > 0 {
				c.Image.Tag = args[0]
			}
			if pwd, err := os.Getwd(); err == nil {
				var t string
				t, c.SourceOwner, c.SourceRepo = github.Remote(logrus.StandardLogger(), pwd)
				if len(t) != 0 {
					token = t // 代表此 repo 是用指定 token clone 的, 因此換掉這次 global 的 token
				}
				c.Image.Name = c.SourceRepo
				c.SourceBranch = github.Head(logrus.StandardLogger(), pwd)
			}
			if c.interactive {
				if c.Image.Tag == "" {
					var err error
					c.Image.Tag, err = github.FindNextReleaseVersion(logrus.StandardLogger(), token, c.SourceOwner, c.SourceRepo)
					if err != nil {
						logrus.Debugln(err)
					}
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
	f.BoolVarP(&c.interactive, "interactive", "i", false, "interactive prompt")
	f.IntVar(&c.promptSize, "interactive-prompt-size", 7, "interactive prompt size")
	f.StringVar(&c.SourceOwner, "source-owner", c.SourceOwner, "name of the owner (user or org) of the repo to create tag")
	f.StringVar(&c.SourceRepo, "source-repo", c.SourceRepo, "name of repo to create tag")
	f.StringVar(&c.SourceBranch, "source-branch", c.SourceBranch, "name of branch to create tag")
	f.StringVar(&c.Image.Name, "image", c.Image.Name, "name of image to build")
	f.StringVar(&c.Jenkins, "jenkins", "https://jenkins.softleader.com.tw", "jenkins to run the pipeline")
	f.StringVar(&c.Deployer, "deployer", "http://softleader.com.tw:5678", "deployer to deploy")
	f.StringVar(&c.ServiceID, "service-id", "", "docker swarm service id to update")
	f.BoolVar(&c.SkipSlack, "skip-slack", false, "skip slack webhook")
	f.StringVar(&c.SlackWebhookURL, "slack-webhook-url", "https://hooks.slack.com/services/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", "slack webhook url")
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
	if c.ServiceID != "" {
		params["serviceID"] = c.ServiceID
	}
	if err := jenkins.Job().BuildWithParameters(c.SourceRepo, params); err != nil {
		return err
	}

	logrus.Printf("Everything is all set, you can check the progress at: %s/job/%s", c.Jenkins, c.SourceRepo)

	if c.ServiceID != "" {
		if ensureJenkinsfileContainsServiceIDHook() {
			if !c.SkipSlack {
				slack.Post(c.SlackWebhookURL, fmt.Sprintf("SIT %s@%s 過版", c.Image.Name, c.Image.Tag))
			}
		}
	}
	return nil
}

// service id 需要 Jenkinsfile 也要配合修改, 如果發現沒有查到關鍵字就提醒一下吧
func ensureJenkinsfileContainsServiceIDHook() (contains bool) {
	pwd, err := os.Getwd()
	if err != nil {
		return
	}
	p := filepath.Join(pwd, "Jenkinsfile")
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return
	}
	jenkinsfile := string(b)
	if !hook.MatchString(jenkinsfile) {
		logrus.Warnf(`not found any hook stage in '%s', auto serviceID update might not work
read more: https://github.com/softleader/softleader-microservice-wiki/wiki/Jenkins-Hook-to-Update-Service-on-Deployer`, p)
		return
	}
	return true
}
