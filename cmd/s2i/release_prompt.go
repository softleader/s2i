package main

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/s2i/pkg/prompt"
)

func releaseQuestions(c *releaseCmd) error {
	if err := prompt.AskRequired("Name of image to build", c.Image.Name, &c.Image.Name); err != nil {
		return err
	}

	if err := prompt.AskRequired("Tag of image to build", c.Image.Tag, &c.Image.Tag); err != nil {
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

	if err := prompt.AskRequired("Name of branch to create to create tag", c.SourceBranch, &c.SourceBranch); err != nil {
		return err
	}

	if err := prompt.AskRequired("Jenkins to run the pipeline", c.Jenkins, &c.Jenkins); err != nil {
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
	return releaseQuestions(c)
}
