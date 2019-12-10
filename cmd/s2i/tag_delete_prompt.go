package main

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/s2i/pkg/prompt"
)

const (
	sep = " "
)

func tagDeleteQuestions(c *tagDeleteCmd) error {
	if err := prompt.AskTagMatcherStrategy("Choose tag matcher strategy", &c.TagMatcherStrategy); err != nil {
		return err
	}

	if err := prompt.AskArrayRequired("Tags to delete (use space to separate each tags if more than one tag)", c.Tags, &c.Tags, sep); err != nil {
		return err
	}

	if err := prompt.AskYesNoBool("Simulate tag deletion", c.DryRun, &c.DryRun); err != nil {
		return err
	}

	if err := prompt.AskRequired("Name of the owner (user or org) to delete tag", c.SourceOwner, &c.SourceOwner); err != nil {
		return err
	}

	if err := prompt.AskRequired("Name of repo to delete tag", c.SourceRepo, &c.SourceRepo); err != nil {
		return err
	}

	ok, err := prompt.Confirm(logrus.StandardLogger(), c)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	logrus.Println("That's try again!")
	return tagDeleteQuestions(c)
}
