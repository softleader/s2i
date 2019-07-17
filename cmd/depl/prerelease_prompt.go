package main

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"strings"
)

func (c *prereleaseCmd) prompt() error {
	var yesNo string

	if err := ask("Name of image to build", c.Image, &c.Image); err != nil {
		return err
	}

	if err := ask("Force to delete the tag if it already exists?", "y", &yesNo); err != nil {
		return err
	}
	c.Force = strings.ToLower(yesNo) == "y" || strings.ToLower(yesNo) == "yes"

	if err := ask("Docker service id to update image (leave blank if you don't need to update)", c.DockerServiceID, &c.DockerServiceID); err != nil {
		return err
	}

	if err := ask("Do you want to go through all of the questions?", "n", &yesNo); err != nil {
		return err
	}
	if strings.ToLower(yesNo) == "n" || strings.ToLower(yesNo) == "no" {
		b, err := yaml.Marshal(c)
		if err != nil {
			logrus.Printf("%#v", c)
		} else {
			logrus.Println(string(b))
		}
		if err := ask("Is this OK?", "y", &yesNo); err != nil {
			return err
		}
		if strings.ToLower(yesNo) == "y" || strings.ToLower(yesNo) == "yes" {
			return nil
		}
	}

	if err := ask("Name of branch to create to create tag", c.SourceBranch, &c.SourceBranch); err != nil {
		return err
	}

	if err := ask("Designating development stage to build, e.g. 0 for alpha, 1 for beta, 2 for release candidate", c.Stage, &c.Stage); err != nil {
		return err
	}

	if err := ask("Name of the owner (user or org) of the repo to create tag", c.SourceOwner, &c.SourceOwner); err != nil {
		return err
	}

	if err := ask("Name of repo to create to create tag", c.SourceRepo, &c.SourceRepo); err != nil {
		return err
	}

	if !c.SkipTests {
		if err := ask("Config server host:port to run the test", c.ConfigServer, &c.ConfigServer); err != nil {
			return err
		}

		if err := ask("Label of config server to run the test, e.g. sqlServer", c.ConfigLabel, &c.ConfigLabel); err != nil {
			return err
		}
	}

	if err := ask("Deployer host:port to update stack service", c.Deployer, &c.Deployer); err != nil {
		return err
	}

	if err := ask("Username of docker registry for jib to build", c.Auth.Username, &c.Auth.Username); err != nil {
		return err
	}

	if err := ask("Password of docker registry for jib to build", c.Auth.Password, &c.Auth.Password); err != nil {
		return err
	}

	return nil
}
