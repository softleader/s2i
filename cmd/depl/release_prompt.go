package main

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/depl/pkg/prompt"
)

func releaseQuestions(c *releaseCmd) error {
	if err := prompt.AskRequired("Name of image to build", c.Image, &c.Image); err != nil {
		return err
	}

	ok, err := prompt.Confirm(logrus.StandardLogger(), c)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	if err := prompt.AskRequired("Name of branch to create to create tag", c.SourceBranch, &c.SourceBranch); err != nil {
		return err
	}

	if err := prompt.AskRequired("Name of the owner (user or org) of the repo to create tag", c.SourceOwner, &c.SourceOwner); err != nil {
		return err
	}

	if err := prompt.AskRequired("Name of repo to create to create tag", c.SourceRepo, &c.SourceRepo); err != nil {
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
