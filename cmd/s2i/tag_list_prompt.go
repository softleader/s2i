package main

import (
	"github.com/sirupsen/logrus"
	"github.com/softleader/s2i/pkg/prompt"
)

func tagListQuestions(c *tagListCmd) error {
	if err := prompt.AskArrayRequired("Tags to list (use space to separate each tags if more than one tag)", c.Tags, &c.Tags, sep); err != nil {
		return err
	}

	if err := prompt.AskYesNoBool("Use regex to matches tags?", c.Regex, &c.Regex); err != nil {
		return err
	}

	ok, err := prompt.Confirm(logrus.StandardLogger(), c)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	if err := prompt.AskRequired("Name of the owner (user or org) to delete tag", c.SourceOwner, &c.SourceOwner); err != nil {
		return err
	}

	if err := prompt.AskRequired("Name of repo to delete tag", c.SourceRepo, &c.SourceRepo); err != nil {
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
	return tagListQuestions(c)
}
