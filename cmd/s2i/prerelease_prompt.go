package main

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/s2i/pkg/prompt"
)

func prereleaseQuestions(c *prereleaseCmd) error {
	if err := prompt.AskRequired("Name of image to build", c.Image.Name, &c.Image.Name); err != nil {
		return err
	}

	if err := prompt.AskRequired("Tag of image to build", c.Image.Tag, &c.Image.Tag); err != nil {
		return err
	}

	if err := prompt.AskYesNo("Force to delete the tag if it already exists?", "y", &c.Force); err != nil {
		return err
	}

	if err := prompt.AskYesNo("Force to check for updated snapshots on remote repositories?", "n", &c.UpdateSnapshots); err != nil {
		return err
	}

	if err := prompt.AskServiceID(logrus.StandardLogger(), "s2i", metadata.String(), c.Deployer, c.Image.Name, c.ServiceID, c.promptSize, &c.ServiceID); err != nil {
		return err
	}

	ok, err := prompt.Confirm(logrus.StandardLogger(), c)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	if err := prompt.AskRequired("Name of the owner (user or org) to create tag", c.SourceOwner, &c.SourceOwner); err != nil {
		return err
	}

	if err := prompt.AskRequired("Name of repo to create tag", c.SourceRepo, &c.SourceRepo); err != nil {
		return err
	}

	if err := prompt.AskRequired("Name of branch to create tag", c.SourceBranch, &c.SourceBranch); err != nil {
		return err
	}

	if err := prompt.AskRequired("Designating development stage to build, e.g. 0 for alpha, 1 for beta, 2 for release candidate", c.Stage, &c.Stage); err != nil {
		return err
	}

	if !c.SkipTests {
		if err := prompt.AskRequired("Config server host:port to run the test", c.ConfigServer, &c.ConfigServer); err != nil {
			return err
		}

		if err := prompt.Ask("Label of config server to run the test, e.g. sqlServer", c.ConfigLabel, &c.ConfigLabel); err != nil {
			return err
		}
	}

	if err := prompt.AskRequired("Username of docker registry for jib to build", c.Auth.Username, &c.Auth.Username); err != nil {
		return err
	}

	if err := prompt.AskRequired("Password of docker registry for jib to build", c.Auth.Password, &c.Auth.Password); err != nil {
		return err
	}

	if err := prompt.AskIntRequired("Strategy to ship source, 0 for auto, 1 for jib, 2 for docker", c.ShipStrategy, &c.ShipStrategy); err != nil {
		return err
	}

	ok, err = prompt.Confirm(logrus.StandardLogger(), c)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	logrus.Println("That's try again!")
	return prereleaseQuestions(c)
}
