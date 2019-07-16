package main

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/depl/pkg/deployer"
	"github.com/softleader/depl/pkg/github"
	"github.com/softleader/depl/pkg/jib"
	"github.com/softleader/depl/pkg/test"
	"github.com/spf13/cobra"
	"os"
)

type buildCmd struct {
	force           bool
	sourceOwner     string
	sourceRepo      string
	sourceBranch    string
	skipTests       bool
	configServer    string
	configLabel     string
	image           string
	tag             string
	deployer        string
	auth            *jib.Auth
	dockerServiceID string
}

func newBuildCmd() *cobra.Command {
	c := &buildCmd{
		auth: &jib.Auth{},
	}
	cmd := &cobra.Command{
		Use:   "build",
		Short: "build source to image and publish to deployer",
		Long:  "build source to image and publish to deployer",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c.tag = args[0]
			if pwd, err := os.Getwd(); err == nil {
				c.sourceOwner, c.sourceRepo = github.Remote(logrus.StandardLogger(), pwd)
				c.image = c.sourceRepo
				c.sourceBranch = github.Head(logrus.StandardLogger(), pwd)
				c.auth = jib.GetAuth(logrus.StandardLogger(), pwd)
			}
			if err := c.prompt(); err != nil {
				return err
			}
			return c.run()
		},
	}

	f := cmd.Flags()
	f.BoolVarP(&c.force, "force", "f", false, "force to delete the tag if it already exists")
	f.BoolVar(&c.skipTests, "skip-tests", false, "skip tests when building image")
	f.StringVar(&c.sourceOwner, "source-owner", c.sourceOwner, "name of the owner (user or org) of the repo to create tag")
	f.StringVar(&c.sourceRepo, "source-repo", c.sourceRepo, "name of repo to create to create tag")
	f.StringVar(&c.sourceBranch, "source-branch", c.sourceBranch, "name of branch to create to create tag")
	f.StringVar(&c.configServer, "config-server", "http://192.168.1.88:8887", "config server to run the test")
	f.StringVar(&c.configLabel, "config-label", "", "the label of config server to run the test, e.g. sqlServer")
	f.StringVar(&c.image, "image", c.image, "name of image to build")
	f.StringVar(&c.deployer, "deployer", "http://softleader.com.tw:5678", "deployer to deploy")
	f.StringVar(&c.auth.Username, "jib-username", "dev", "username of docker registry for jib to build")
	f.StringVar(&c.auth.Password, "jib-password", "sleader", "password of docker registry for jib to build")
	f.StringVar(&c.dockerServiceID, "docker-service-id", "", "docker service id to update image")
	return cmd
}

func (c *buildCmd) run() error {
	if !c.skipTests {
		if err := test.Run(logrus.StandardLogger(), c.configServer, c.configLabel); err != nil {
			return err
		}
	}
	if err := jib.Build(logrus.StandardLogger(), c.image, c.tag, c.auth); err != nil {
		return err
	}
	if err := github.CreateRelease(logrus.StandardLogger(), token, c.sourceOwner, c.sourceRepo, c.sourceBranch, c.tag, c.force); err != nil {
		return err
	}
	if c.dockerServiceID != "" {
		if err := deployer.UpdateService(logrus.StandardLogger(), "depl", metadata.String(), c.deployer, c.dockerServiceID, c.image, c.tag); err != nil {
			return err
		}
	}
	logrus.Printf("Everything is all set, you are good to go.")
	return nil
}
